# 扩展 CLI 能力与插件

Kubebuilder 提供可扩展的插件架构用于脚手架生成项目。通过插件，你可以自定义 CLI 行为或集成新特性。

本文介绍如何扩展 CLI 能力、创建自定义插件以及对多个插件进行组合（Bundle）。

## 创建自定义插件

要创建自定义插件，你需要实现 [Kubebuilder Plugin 接口][plugin-interface]。

该接口允许你的插件挂接 Kubebuilder 的命令（如 `init`、`create api`、`create webhook` 等），在执行时注入自定义逻辑。

### 自定义插件示例

你可以创建一个通过 [Bundle Plugin](#bundle-plugin) 同时生成“语言相关脚手架 + 配置文件”的插件。下面示例将 Golang 插件与 Kustomize 插件进行组合：

```go
import (
    kustomizecommonv2 "sigs.k8s.io/kubebuilder/v4/pkg/plugins/common/kustomize/v2"
    golangv4 "sigs.k8s.io/kubebuilder/v4/pkg/plugins/golang/v4"
)

mylanguagev1Bundle, _ := plugin.NewBundle(
    plugin.WithName("mylanguage.kubebuilder.io"),
    plugin.WithVersion(plugin.Version{Number: 1}),
    plugin.WithPlugins(kustomizecommonv2.Plugin{}, mylanguagev1.Plugin{}),
)
```

这种组合可以：通过 Kustomize 提供通用配置基线，同时由 `mylanguagev1` 生成语言相关文件。

此外，你也可以借助 `create api` 与 `create webhook` 子命令，脚手架生成特定资源（如 CRD 与 Controller）。

### 插件子命令

插件需要实现在子命令调用时被执行的代码。你可以实现 [Plugin 接口][plugin-interface] 来创建新插件。

除基础能力 `Base` 外，插件还应实现 [`SubcommandMetadata`][plugin-subc-metadata] 接口以便经由 CLI 运行。可以选择自定义目标命令的帮助信息（若不提供则保留 [cobra][cobra] 默认帮助）。

Kubebuilder CLI 插件将脚手架与 CLI 能力封装为 Go 类型，由 `kubebuilder` 可执行文件（或任意导入该插件的可执行文件）运行。插件会配置以下命令之一的执行：

- `init`：初始化项目结构。
- `create api`：脚手架生成新的 API 与控制器。
- `create webhook`：脚手架生成新的 Webhook。
- `edit`：编辑项目结构。

示例：使用自定义插件运行 `init`：

```sh
kubebuilder init --plugins=mylanguage.kubebuilder.io/v1
```

这会使用 `mylanguage` 插件初始化项目。

### 插件键（Plugin Key）

插件以 `<name>/<version>` 形式标识。指定插件有两种方式：

- 通过命令行设置：`kubebuilder init --plugins=<plugin key>`；
- 在脚手架生成的 [PROJECT 配置文件][project-file-config] 中设置 `layout: <plugin key>`（除 `init` 外，其它命令会读取该值并据此选择插件）。

默认情况下，`<plugin key>` 形如 `go.kubebuilder.io/vX`（X 为整数）。完整实现可参考 Kubebuilder 内置 [`go.kubebuilder.io`][kb-go-plugin] 插件。

### 插件命名

插件名必须符合 DNS1123 标签规范，且应使用完全限定名（例如带 `.example.com` 后缀）。例如 `go.kubebuilder.io`。限定名可避免命名冲突。

### 插件版本

插件的 `Version()` 返回 [`plugin.Version`][plugin-version-type]，包含一个整数版本与可选阶段字符串（`alpha` 或 `beta`）。

- 不同整数表示不兼容版本；
- 阶段说明稳定性：`alpha` 变更频繁，`beta` 仅小改动（如修复）。

### 模板与样板（Boilerplates）

Kubebuilder 内置插件通过模板生成代码文件。例如 `go/v4` 在初始化时会用模板脚手架生成 `go.mod`。

在自定义插件中，你可以基于 [machinery 库][machinery] 定义模板与文件生成逻辑：

- 定义文件 I/O 行为；
- 向模板中添加 [markers][markers-scaffold]；
- 指定模板内容并执行脚手架生成。

示例：`go/v4` 通过实现 machinery 接口对象来脚手架生成 `go.mod`，其原始模板在 `Template.SetTemplateDefaults` 的 `TemplateBody` 字段中定义：

```go
/*
Copyright 2022 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package templates

import (
	"sigs.k8s.io/kubebuilder/v4/pkg/machinery"
)

var _ machinery.Template = &GoMod{}

// GoMod scaffolds a file that defines the project dependencies
type GoMod struct {
	machinery.TemplateMixin
	machinery.RepositoryMixin

	ControllerRuntimeVersion string
}

// SetTemplateDefaults implements machinery.Template
func (f *GoMod) SetTemplateDefaults() error {
	if f.Path == "" {
		f.Path = "go.mod"
	}

	f.TemplateBody = goModTemplate

	f.IfExistsAction = machinery.OverwriteFile

	return nil
}

const goModTemplate = `module {{ .Repo }}

go 1.24.5

require (
	sigs.k8s.io/controller-runtime {{ .ControllerRuntimeVersion }}
)
`
```

随后，该对象会被传入脚手架执行：

```go
// Scaffold implements cmdutil.Scaffolder
func (s *initScaffolder) Scaffold() error {
    log.Println("Writing scaffold for you to edit...")

    scaffold := machinery.NewScaffold(s.fs,
        machinery.WithConfig(s.config),
    )

    ...

    return scaffold.Execute(
        ...
        &templates.GoMod{
            ControllerRuntimeVersion: ControllerRuntimeVersion,
        },
        ...
    )
}
```

#### 覆写已存在文件

当子命令执行时，如果你希望覆写已有文件，可以在模板定义中设置：

```go
f.IfExistsAction = machinery.OverwriteFile
```

借助这些选项，你的插件可以接管并调整 Kubebuilder 默认脚手架生成的文件。

## 定制已有脚手架

Kubebuilder 提供了实用函数帮助你修改默认脚手架。借助[插件工具集][plugin-utils]，你可以在文件中插入、替换或追加内容：

- 插入内容：在目标位置添加内容；
- 替换内容：查找并替换指定片段；
- 追加内容：在文件末尾追加，不影响既有内容。

示例：使用 `InsertCode` 向文件内注入自定义内容：

```go
pluginutil.InsertCode(filename, target, code)
```

更多细节可参考 [Kubebuilder 插件工具集][kb-utils]。

## Bundle Plugin

可将多个插件打包为 Bundle，以组合执行更复杂的脚手架流程：

```go
myPluginBundle, _ := plugin.NewBundle(
    plugin.WithName("myplugin.example.com"),
    plugin.WithVersion(plugin.Version{Number: 1}),
    plugin.WithPlugins(pluginA.Plugin{}, pluginB.Plugin{}, pluginC.Plugin{}),
)
```

上述 Bundle 会按顺序为各插件执行 `init`：

1. pluginA
2. pluginB
3. pluginC

运行命令：

```sh
kubebuilder init --plugins=myplugin.example.com/v1
```

## CLI 系统

插件由 [`CLI`][cli] 对象运行：它将插件类型映射到子命令并调用插件方法。例如，向 `CLI` 注入一个 `Init` 插件并调用 `CLI.Run()`，就会在 `kubebuilder init` 时依次调用该插件的 [SubcommandMetadata][plugin-sub-command]、[UpdatesMetadata][plugin-update-meta] 与 `Run`，并传入用户参数。

示例程序：

```go
package cli

import (
    log "log/slog"
    "github.com/spf13/cobra"

    "sigs.k8s.io/kubebuilder/v4/pkg/cli"
    cfgv3 "sigs.k8s.io/kubebuilder/v4/pkg/config/v3"
    "sigs.k8s.io/kubebuilder/v4/pkg/plugin"
    kustomizecommonv2 "sigs.k8s.io/kubebuilder/v4/pkg/plugins/common/kustomize/v2"
    "sigs.k8s.io/kubebuilder/v4/pkg/plugins/golang"
    deployimage "sigs.k8s.io/kubebuilder/v4/pkg/plugins/golang/deploy-image/v1alpha1"
    golangv4 "sigs.k8s.io/kubebuilder/v4/pkg/plugins/golang/v4"

)

var (
    // 你的 CLI 里可能会有的命令
    commands = []*cobra.Command{
        myExampleCommand.NewCmd(),
    }
    alphaCommands = []*cobra.Command{
        myExampleAlphaCommand.NewCmd(),
    }
)

// GetPluginsCLI 返回带插件的 CLI，用于你的 CLI 二进制
func GetPluginsCLI() (*cli.CLI) {
    // 组合插件：Kubebuilder go/v4 的 Golang 项目脚手架
    gov3Bundle, _ := plugin.NewBundleWithOptions(plugin.WithName(golang.DefaultNameQualifier),
        plugin.WithVersion(plugin.Version{Number: 3}),
        plugin.WithPlugins(kustomizecommonv2.Plugin{}, golangv4.Plugin{}),
    )


    c, err := cli.New(
        // 你的 CLI 名称
        cli.WithCommandName("example-cli"),

        // 你的 CLI 版本
        cli.WithVersion(versionString()),

        // 注册可用于脚手架的插件（示例使用 Kubebuilder 提供的插件）
        cli.WithPlugins(
            gov3Bundle,
            &deployimage.Plugin{},
        ),

        // 设置默认插件（未指定时使用）
        cli.WithDefaultPlugins(cfgv3.Version, gov3Bundle),

        // 设置默认的项目配置版本（未通过 --project-version 指定时）
        cli.WithDefaultProjectVersion(cfgv3.Version),

        // 添加自定义命令
        cli.WithExtraCommands(commands...),

        // 添加自定义 alpha 命令
        cli.WithExtraAlphaCommands(alphaCommands...),

        // 开启自动补全
        cli.WithCompletion(),
    )
    if err != nil {
        log.Fatal(err)
    }

    return c
}

// versionString 返回 CLI 版本
func versionString() string {
    // return your binary project version
}
```

该程序的运行方式示例：

默认行为：

```sh
# 使用默认的 Init 插件（例如 "go.example.com/v1"）初始化项目，
# 该键会自动写入 PROJECT 配置文件
$ my-bin-builder init

# 读取配置文件中的键，使用 "go.example.com/v1" 的 CreateAPI 与 CreateWebhook
$ my-bin-builder create api [flags]
$ my-bin-builder create webhook [flags]
```

通过 `--plugins` 指定插件：

```sh
# 使用 "ansible.example.com/v1" 的 Init 插件初始化项目，并写入配置文件
$ my-bin-builder init --plugins ansible

# 读取配置文件中的键，使用 "ansible.example.com/v1" 的 CreateAPI 与 CreateWebhook
$ my-bin-builder create api [flags]
$ my-bin-builder create webhook [flags]
```

### 在 PROJECT 文件中跟踪输入

CLI 负责管理[PROJECT 配置文件][project-file-config]，用于记录由 CLI 脚手架生成的项目信息。

扩展 Kubebuilder 时，建议你的工具或[外部插件][external-plugin]正确读写该文件以追踪关键信息：

- 便于其它工具与插件正确集成；
- 便于基于已追踪的数据进行“二次脚手架”（如使用[Alpha 命令](./../../reference/alpha_commands.md)升级项目结构）。

例如，插件可以据此判断是否支持当前项目的布局，并基于已记录的输入参数重新执行命令。

#### 示例

使用 [Deploy Image][deploy-image] 插件为 API 及其控制器脚手架：

```sh
kubebuilder create api --group example.com --version v1alpha1 --kind Memcached --image=memcached:memcached:1.6.26-alpine3.19 --image-container-command="memcached,--memory-limit=64,-o,modern,-v" --image-container-port="11211" --run-as-user="1001" --plugins="deploy-image/v1-alpha" --make=false
```

PROJECT 文件中将新增：

```yaml
...
plugins:
  deploy-image.go.kubebuilder.io/v1-alpha:
    resources:
    - domain: testproject.org
      group: example.com
      kind: Memcached
      options:
        containerCommand: memcached,--memory-limit=64,-o,modern,-v
        containerPort: "11211"
        image: memcached:memcached:1.6.26-alpine3.19
        runAsUser: "1001"
      version: v1alpha1
    - domain: testproject.org
      group: example.com
      kind: Busybox
      options:
        image: busybox:1.36.1
      version: v1alpha1
...
```

通过检查 PROJECT 文件，就能了解插件如何被使用、传入了哪些参数。这样不仅可复现命令执行，也便于开发依赖这些信息的能力或插件。

[sdk]: https://github.com/operator-framework/operator-sdk
[plugin-interface]: https://pkg.go.dev/sigs.k8s.io/kubebuilder/v4/pkg/plugin
[machinery]: https://github.com/kubernetes-sigs/kubebuilder/tree/master/pkg/machinery
[plugin-subc-metadata]: https://pkg.go.dev/sigs.k8s.io/kubebuilder/v4/pkg/plugin#SubcommandMetadata
[plugin-version-type]: https://pkg.go.dev/sigs.k8s.io/kubebuilder/v4/pkg/plugin#Version
[bundle-plugin-doc]: https://pkg.go.dev/sigs.k8s.io/kubebuilder/v4/pkg/plugin#Bundle
[deprecate-plugin-doc]: https://pkg.go.dev/sigs.k8s.io/kubebuilder/v4/pkg/plugin#Deprecated
[plugin-sub-command]: https://pkg.go.dev/sigs.k8s.io/kubebuilder/v4/pkg/plugin#Subcommand
[plugin-update-meta]: https://pkg.go.dev/sigs.k8s.io/kubebuilder/v4/pkg/plugin#UpdatesMetadata
[plugin-utils]: https://pkg.go.dev/sigs.k8s.io/kubebuilder/v4/pkg/plugin/util
[markers-scaffold]: ./../../reference/markers/scaffold.md
[kb-utils]: https://github.com/kubernetes-sigs/kubebuilder/blob/book-v4/pkg/plugin/util/util.go
[project-file-config]: ./../../reference/project-config.md
[cli]: https://github.com/kubernetes-sigs/kubebuilder/tree/book-v4/pkg/cli
[kb-go-plugin]: https://github.com/kubernetes-sigs/kubebuilder/tree/book-v4/pkg/plugins/golang/v4
[cobra]: https://github.com/spf13/cobra
[external-plugin]: external-plugins.md
[deploy-image]: ./../available/deploy-image-plugin-v1-alpha.md
[upgrade-assistant]: ./../../reference/rescaffold.md
