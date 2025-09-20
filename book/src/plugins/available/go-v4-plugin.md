# go/v4 (go.kubebuilder.io/v4)

（默认脚手架）

Kubebuilder 在初始化项目时指定 `--plugins=go/v4` 后将使用该插件进行脚手架生成。
该插件通过 [Bundle Plugin][bundle] 组合了 `kustomize.common.kubebuilder.io/v2` 与 `base.go.kubebuilder.io/v4`，
用于生成一套项目模板，便于你构建成组的 [controllers][controller-runtime]。

按照[快速开始][quickstart]创建项目时，默认即会使用该插件。

<aside class="note">
<h1>示例</h1>

你可以在 Kubebuilder 仓库根目录下的 [testdata][testdata] 中，查看以 `project-v4-<options>` 命名的示例工程来了解该插件的用法。

</aside>

## 如何使用？

创建一个启用 `go/v4` 插件的新项目，可使用如下命令：

```sh
kubebuilder init --domain tutorial.kubebuilder.io --repo tutorial.kubebuilder.io/project --plugins=go/v4
```

## 支持的子命令

- Init - `kubebuilder init [OPTIONS]`
- Edit - `kubebuilder edit [OPTIONS]`
- Create API - `kubebuilder create api [OPTIONS]`
- Create Webhook - `kubebuilder create webhook [OPTIONS]`

## 延伸阅读

- 查看插件组合方式：Kubebuilder 源码中的 [main.go][plugins-main]
- 查看基础 Go 插件实现：[`base.go.kubebuilder.io/v4`][v4-plugin]
- 查看 Kustomize/v2 插件实现：[Kustomize/v2][kustomize-plugin]
- 了解更多控制器知识：[controller-runtime][controller-runtime]

[controller-runtime]: https://github.com/kubernetes-sigs/controller-runtime
[quickstart]: ./../../quick-start.md
[testdata]: https://github.com/kubernetes-sigs/kubebuilder/tree/master/testdata
[plugins-main]: ./../../../../../cmd/main.go
[kustomize-plugin]: ./../../plugins/available/kustomize-v2.md
[kustomize]: https://github.com/kubernetes-sigs/kustomize
[standard-go-project]: https://github.com/golang-standards/project-layout
[v4-plugin]: ./../../../../../pkg/plugins/golang/v4
[migration-guide-doc]: ./../../migration/migration_guide_gov3_to_gov4.md
[project-doc]: ./../../reference/project-config.md
[bundle]: ./../../../../../pkg/plugin/bundle.go
