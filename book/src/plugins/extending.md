# 扩展 Kubebuilder

Kubebuilder 提供可扩展的插件架构用于脚手架生成项目。通过插件，你可以自定义 CLI 行为或集成新特性。

## 概览

你可以通过自定义插件扩展 Kubebuilder 的 CLI，以便：

- 构建新的脚手架。
- 增强已有脚手架。
- 为脚手架系统添加新的命令与功能。

这种灵活性让你可以按照具体需求搭建定制化的项目基线。

<aside class="note">
<h1>为什么采用 Kubebuilder 风格？</h1>

Kubebuilder 与 Operator SDK 都广泛采用并基于 [controller-runtime][controller-runtime]。二者都支持使用[Operator 模式][operator-pattern]构建解决方案，并遵循通用标准。

采用这些标准能带来显著收益：共同维护通用能力、复用社区贡献，使你能把精力聚焦在插件与场景的特定需求上。同时，你也能复用这些项目现在或将来提供的自定义插件与选项。

</aside>

## 扩展方式

扩展 Kubebuilder 主要有两种途径：

1. 扩展 CLI 能力与插件：
   基于已有插件进行二次开发以[扩展其能力][extending-cli]。当一个工具已受益于 Kubebuilder 的脚手架体系、你仅需补齐特定能力时很有用。
   例如 [Operator SDK][sdk] 复用了 [kustomize 插件][kustomize-plugin]，从而为 Ansible/Helm 等语言提供支持，使项目只需维护语言相关的差异部分。

2. 编写外部插件：
   构建独立二进制的插件，可用任意语言实现，但需遵循 Kubebuilder 识别的执行约定。参见[创建外部插件][external-plugins]。

想进一步了解如何扩展 Kubebuilder，请阅读：

- [CLI 与插件](./extending/extending_cli_features_and_plugins.md)：如何扩展 CLI 与插件。
- [外部插件](./extending/external-plugins.md)：如何创建独立插件。
- [E2E 测试](./extending/testing-plugins.md)：如何确保插件如期工作。

[extending-cli]: ./extending/extending_cli_features_and_plugins.md
[external-plugins]: ./extending/external-plugins.md
[sdk]: https://github.com/operator-framework/operator-sdk
[kustomize-plugin]: ./available/kustomize-v2.md
[controller-runtime]: https://github.com/kubernetes-sigs/controller-runtime
[operator-pattern]: https://kubernetes.io/docs/concepts/extend-kubernetes/operator/
