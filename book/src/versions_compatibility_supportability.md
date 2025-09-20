# 版本兼容性与可支持性

Kubebuilder 创建的项目包含一个 `Makefile`，用于安装在项目创建时定义版本的工具。主要包含以下工具：

- [kustomize](https://github.com/kubernetes-sigs/kustomize)
- [controller-gen](https://github.com/kubernetes-sigs/controller-tools)
- [setup-envtest](https://github.com/kubernetes-sigs/controller-runtime/tree/main/tools/setup-envtest)

此外，这些项目还包含一个 `go.mod` 文件用于指定依赖版本。Kubebuilder 依赖于 [controller-runtime](https://github.com/kubernetes-sigs/controller-runtime) 以及它的 Go 与 Kubernetes 依赖。因此，`Makefile` 与 `go.mod` 中定义的版本是已测试、受支持且被推荐的版本。

Kubebuilder 的每个次版本都会与特定的 client-go 次版本进行测试。尽管某个 Kubebuilder 次版本可能与其他 client-go 次版本或其他工具兼容，但这种兼容性并不保证、也不受支持或测试覆盖。

Kubebuilder 所需的最低 Go 版本由其依赖项中所需的最高最低 Go 版本决定。这通常与相应的 `k8s.io/*` 依赖所要求的最低 Go 版本保持一致。

兼容的 `k8s.io/*` 版本、client-go 版本和最低 Go 版本可在每个 [标签版本](https://github.com/kubernetes-sigs/kubebuilder/tags) 的项目脚手架 `go.mod` 文件中找到。

示例：对于 `4.1.1` 版本，最低 Go 版本兼容性为 `1.22`。你可以参考该标签版本 [v4.1.1](https://github.com/kubernetes-sigs/kubebuilder/tree/v4.1.1/testdata) 的 testdata 目录中的示例，例如 `project-v4` 的 [go.mod](https://github.com/kubernetes-sigs/kubebuilder/blob/v4.1.1/testdata/project-v4/go.mod#L3) 文件。你也可以通过查看 [Makefile](https://github.com/kubernetes-sigs/kubebuilder/blob/v4.1.1/testdata/project-v4/Makefile#L160-L165) 来检查该版本所支持并经过测试的工具版本。

## 支持的操作系统

当前，Kubebuilder 官方支持 macOS 与 Linux 平台。如果你使用的是 Windows 系统，可能会遇到问题。欢迎提交贡献以支持 Windows。

<aside class="note warning">
<h1>项目自定义</h1>

通过 CLI 创建项目后，你可以按需进行自定义。但请注意，除非非常确定，否则不建议偏离建议的项目布局。

例如，不要随意移动脚手架文件，否则将来会很难升级项目。你也可能失去使用某些 CLI 功能与辅助工具的能力。有关项目布局的更多信息，请参阅文档：[基础项目包含什么？][basic-project-doc]

</aside>

[basic-project-doc]: ./cronjob-tutorial/basic-project.md
