# 插件

Kubebuilder 的架构从根本上是基于插件的。
这种设计使 Kubebuilder CLI 能在保持对旧版本向后兼容的同时演进；允许用户按需启用或禁用特性，并能与外部工具无缝集成。

通过利用插件，项目可以扩展 Kubebuilder，并将其作为库来支持新的功能，或实现贴合用户需求的自定义脚手架。这种灵活性允许维护者在 Kubebuilder 的基础上构建，适配特定用例，同时受益于其强大的脚手架引擎。

插件具备以下关键优势：

- 兼容性：确保旧的布局与项目结构在新版本下仍能工作
- 可定制：允许用户按需启用或禁用特性（例如 [Grafana][grafana-plugin] 与 [Deploy Image][deploy-image] 插件）
- 可扩展：便于集成第三方工具与希望提供自有[外部插件][external-plugins]的项目，这些插件可与 Kubebuilder 协同使用，以修改和增强项目脚手架或引入新功能

例如，使用多个全局插件初始化项目：

```sh
kubebuilder init --plugins=pluginA,pluginB,pluginC
```

例如，使用特定插件应用自定义脚手架：

```sh
kubebuilder create api --plugins=pluginA,pluginB,pluginC
OR
kubebuilder create webhook --plugins=pluginA,pluginB,pluginC
OR
kubebuilder edit --plugins=pluginA,pluginB,pluginC
```

本节将介绍可用插件、如何扩展 Kubebuilder，以及如何在遵循相同布局结构的前提下创建你自己的插件。

- [可用插件](./available-plugins.md)
- [扩展](./extending.md)
- [插件版本管理](./plugins-versioning.md)

[extending-cli]: extending.md
[grafana-plugin]: ./available/grafana-v1-alpha.md
[deploy-image]: ./available/deploy-image-plugin-v1-alpha.md
[external-plugins]: ./extending/external-plugins.md
