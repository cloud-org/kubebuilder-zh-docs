# 为 Kubebuilder 创建外部插件

## 概览

Kubebuilder 的能力可以通过外部插件扩展。外部插件是可执行文件（可用任意语言实现），需遵循 Kubebuilder 识别的执行契约。Kubebuilder 通过 `stdin`/`stdout` 与插件交互。

## 为什么使用外部插件？

外部插件让第三方方案维护者将其工具与 Kubebuilder 集成。与 Kubebuilder 自身插件类似，外部插件是“可选启用”的，赋予用户工具选择的灵活性。把插件放在其自有仓库中，有助于与其 CI 流水线同步演进，并在其职责边界内管理变更。

如需此类集成，建议与你所依赖的第三方方案维护者协作。Kubebuilder 维护者也乐于支持扩展其能力。

## 如何编写外部插件

Kubebuilder 与外部插件通过标准 I/O 通信。只要遵循 [PluginRequest][code-plugin-external] 与 [PluginResponse][code-plugin-external] 结构，任何语言均可实现。

### PluginRequest

`PluginRequest` 包含从 CLI 收集的参数与此前已执行插件的输出。Kubebuilder 会通过 `stdin` 以 JSON 发送给外部插件。

示例（执行 `kubebuilder init --plugins sampleexternalplugin/v1 --domain my.domain` 时发送的 `PluginRequest`）：

```json
{
  "apiVersion": "v1alpha1",
  "args": ["--domain", "my.domain"],
  "command": "init",
  "universe": {}
}
```

### PluginResponse

`PluginResponse` 用于描述外部插件对项目所做的修改。Kubebuilder 通过 `stdout` 读取 JSON 格式的返回值。

示例 `PluginResponse`：

```json
{
  "apiVersion": "v1alpha1",
  "command": "init",
  "metadata": {
    "description": "The `init` subcommand initializes a project via Kubebuilder. It scaffolds a single file: `initFile`.",
    "examples": "kubebuilder init --plugins sampleexternalplugin/v1 --domain my.domain"
  },
  "universe": {
    "initFile": "A file created with the `init` subcommand."
  },
  "error": false,
  "errorMsgs": []
}
```

<aside>
<H1> </H1>

请不要在外部插件中直接向 `stdout` 打印日志。由于与 Kubebuilder 的通信依赖结构化 JSON 的 `stdin`/`stdout`，任何意外输出（如调试日志）都可能导致解析失败。需要日志时请写入文件。

</aside>

## 如何使用外部插件

### 前置条件

- Kubebuilder CLI 版本 > 3.11.0
- 外部插件的可执行文件
- 配置插件查找路径：使用 `${EXTERNAL_PLUGINS_PATH}`，或采用默认的系统路径：
  - Linux：`$HOME/.config/kubebuilder/plugins/${name}/${version}/${name}`
  - macOS：`~/Library/Application Support/kubebuilder/plugins/${name}/${version}/${name}`

示例：Linux 上名为 `foo.acme.io`、版本 `v2` 的插件路径为 `$HOME/.config/kubebuilder/plugins/foo.acme.io/v2/foo.acme.io`。

### 支持的子命令

外部插件可支持以下子命令：
- `init`：项目初始化
- `create api`：脚手架生成 Kubernetes API 定义
- `create webhook`：脚手架生成 Kubernetes Webhook
- `edit`：更新项目配置

可选的增强子命令：
- `metadata`：配合 `--help` 提供描述与示例
- `flags`：声明支持的参数，便于提前做参数校验

<aside class="note">
<h1>关于 `flags` 子命令</h1>

外部插件实现 `flags` 子命令后，可以在执行前由 Kubebuilder 预先校验不被支持的参数，从而尽早失败。若未实现 `flags`，Kubebuilder 会将所有参数原样传递给外部插件，由插件自行处理非法参数。

</aside>

### 配置插件路径

设置环境变量 `$EXTERNAL_PLUGINS_PATH` 指定自定义插件二进制路径：

```sh
export EXTERNAL_PLUGINS_PATH=<custom-path>
```

否则 Kubebuilder 会根据操作系统在默认路径下查找插件。

### CLI 命令示例

```sh
# 使用名为 `sampleplugin` 的外部插件初始化项目
kubebuilder init --plugins sampleplugin/v1

# 查看该外部插件的 init 子命令帮助
kubebuilder init --plugins sampleplugin/v1 --help

# 使用自定义参数 `number` 创建 API
kubebuilder create api --plugins sampleplugin/v1 --number 2

# 使用自定义参数 `hooked` 创建 webhook
kubebuilder create webhook --plugins sampleplugin/v1 --hooked

# 使用外部插件更新项目配置
kubebuilder edit --plugins sampleplugin/v1

# 以链式顺序同时使用 v1 与 v2 两个外部插件创建 API
kubebuilder create api --plugins sampleplugin/v1,sampleplugin/v2

# 使用 go/v4 创建 API 后，再链式传递给外部插件处理
kubebuilder create api --plugins go/v4,sampleplugin/v1
```

## 延伸阅读

- 一个用 Go 编写的[外部插件示例](https://github.com/kubernetes-sigs/kubebuilder/tree/master/docs/book/src/simple-external-plugin-tutorial/testdata/sampleexternalplugin/v1)
- 一个用 Python 编写的[外部插件示例](https://github.com/rashmigottipati/POC-Phase2-Plugins)
- 一个用 JavaScript 编写的[外部插件示例](https://github.com/Eileen-Yu/kb-js-plugin)

[code-plugin-external]: https://github.com/kubernetes-sigs/kubebuilder/blob/book-v4/pkg/plugin/external/types.go
