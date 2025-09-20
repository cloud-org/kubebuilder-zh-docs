# [默认脚手架] Kustomize v2

Kustomize 插件允许你脚手架生成与语言基础插件 `base.go.kubebuilder.io/v4` 搭配使用的全部 kustomize 清单。
对于 `go/v4`（默认脚手架）创建的项目，本插件会在 `config/` 目录下生成相应清单。

注意：[Operator-sdk][sdk] 这类项目会把 Kubebuilder 作为库使用，并提供 Ansible、Helm 等语言支持。
Kustomize 插件让它们能方便地维护统一配置，并确保不同语言下配置一致。对于希望在默认脚手架之上做“二次加工”的场景，本插件也很有帮助：
我们无需在所有语言插件里手工同步配置，从而还可以创建可复用在多个项目和语言中的“辅助”插件。

<aside class="note">
<h1>示例</h1>

你可以在 Kubebuilder 项目的 [testdata][testdata] 目录下，查看 `project-v4-*` 示例的 `config/` 目录以了解 kustomize 内容。

</aside>

## 适用场景（When to use it）

- 你希望为自有语言插件脚手架生成 kustomize 配置清单；
- 你需要 Apple Silicon（`darwin/arm64`）支持（kustomize 4.x 之前该平台没有官方二进制）；
- 你希望尝试并使用 kustomize v4、v5 的新语法与能力（参见[发布说明 v4][release-notes-v4]、[发布说明 v5][release-notes-v5]）；
- 你不需要面向 Kubernetes 集群版本 < `1.22` 的兼容性（kustomize v4 的一些新特性在 `kubectl < 1.22` 下不受支持，可能不可用）；
- 你不依赖资源字段中的特殊 URL；
- 你希望使用 [replacements][kustomize-replacements]，因为 [vars][kustomize-vars] 已被废弃且可能很快被移除。

## 如何使用

若要声明你的语言插件应使用 kustomize，可使用 [Bundle Plugin][bundle] 指定：由“语言插件 + kustomize 配置插件”组合实现：

```go
import (
    ...
    kustomizecommonv2 "sigs.k8s.io/kubebuilder/v4/pkg/plugins/common/kustomize/v2"
    golangv4 "sigs.k8s.io/kubebuilder/v4/pkg/plugins/golang/v4"
    ...
)

// 组合插件：用于 Kubebuilder go/v4 的 Golang 项目脚手架
// 下面的代码通过组合创建了一个带名称与版本的新插件；
// 你可以声明 1 个或多个插件共同组成一个组合插件
gov3Bundle, _ := plugin.NewBundle(plugin.WithName(golang.DefaultNameQualifier),
    plugin.WithVersion(plugin.Version{Number: 3}),
    plugin.WithPlugins(kustomizecommonv2.Plugin{}, golangv4.Plugin{}), // 生成 config/ 与全部 kustomize 文件
    // 同时生成 Golang 文件与语言专属内容（如 go.mod、apis、controllers）
)
```

你也可以单独使用 kustomize/v2：

```sh
kubebuilder init --plugins=kustomize/v2
$ ls -la
total 24
drwxr-xr-x   6 camilamacedo86  staff  192 31 Mar 09:56 .
drwxr-xr-x  11 camilamacedo86  staff  352 29 Mar 21:23 ..
-rw-------   1 camilamacedo86  staff  129 26 Mar 12:01 .dockerignore
-rw-------   1 camilamacedo86  staff  367 26 Mar 12:01 .gitignore
-rw-------   1 camilamacedo86  staff   94 31 Mar 09:56 PROJECT
drwx------   6 camilamacedo86  staff  192 31 Mar 09:56 config
```

或与基础语言插件组合：

```sh
# 与 go/v4 等价的组合脚手架，但显式使用 kustomize/v2
kubebuilder init --plugins=kustomize/v2,base.go.kubebuilder.io/v4 --domain example.org --repo example.org/guestbook-operator
```

## 子命令

kustomize 插件实现了以下子命令：

* init（`$ kubebuilder init [OPTIONS]`）
* create api（`$ kubebuilder create api [OPTIONS]`）
* create webhook（`$ kubebuilder create api [OPTIONS]`）

<aside class="note">
<h1>Create API 与 Webhook</h1>

`create api` 的实现会为每个 API 生成对应的 kustomize 清单，详见 [kustomize-create-api][kustomize-create-api]。`create webhook` 同理。

</aside>

## 影响的文件

本插件会创建或更新以下脚手架：

* `config/*`

## 延伸阅读

* kustomize 插件[实现代码](https://github.com/kubernetes-sigs/kubebuilder/tree/master/pkg/plugins/common/kustomize)
* [kustomize 文档][kustomize-docs]
* [kustomize 仓库][kustomize-github]
* kustomize v5.0.0 的[发布说明][release-notes-v5]
* kustomize v4.0.0 的[发布说明][release-notes-v4]
* 也可对比样例 `project-v3` 与 `project-v4` 的 `config/` 目录，了解默认清单语法的差异

[sdk]:https://github.com/operator-framework/operator-sdk
[testdata]: https://github.com/kubernetes-sigs/kubebuilder/tree/master/testdata/
[bundle]: https://github.com/kubernetes-sigs/kubebuilder/blob/master/pkg/plugin/bundle.go
[kustomize-create-api]: https://github.com/kubernetes-sigs/kubebuilder/blob/master/pkg/plugins/common/kustomize/v2/scaffolds/api.go#L72-L84
[kustomize-docs]: https://kustomize.io/
[kustomize-github]: https://github.com/kubernetes-sigs/kustomize
[kustomize-replacements]: https://kubectl.docs.kubernetes.io/references/kustomize/kustomization/replacements/
[kustomize-vars]: https://kubectl.docs.kubernetes.io/references/kustomize/kustomization/vars/
[release-notes-v5]: https://github.com/kubernetes-sigs/kustomize/releases/tag/kustomize%2Fv5.0.0
[release-notes-v4]: https://github.com/kubernetes-sigs/kustomize/releases/tag/kustomize%2Fv4.0.0

