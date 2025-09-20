**注意：** 性子急的读者可直接前往 [快速开始](quick-start.md)。

**正在使用 Kubebuilder 的 v1 或 v2 版本？**
**请查看旧版文档：[v1](https://book-v1.book.kubebuilder.io)、[v2](https://book-v2.book.kubebuilder.io) 或 [v3](https://book-v3.book.kubebuilder.io)**

## 适用读者

#### Kubernetes 用户

Kubernetes 的用户将通过学习 API 设计与实现背后的基本概念，获得对 Kubernetes 更深入的理解。本书将教读者如何开发自己的 Kubernetes API，以及核心 Kubernetes API 的设计原则。

包括：

- Kubernetes API 与资源的结构
- API 版本化语义
- 自愈
- 垃圾回收与 Finalizer
- 声明式 vs 命令式 API
- 基于电平（Level-Based）vs 基于边沿（Edge-Base）API
- 资源 vs 子资源

#### Kubernetes API 扩展开发者

API 扩展开发者将学习实现典型 Kubernetes API 的原则与概念，以及用于快速落地的简洁工具与库。本书还涵盖扩展开发者常见的陷阱与误区。

包括：

- 如何将多个事件批量进一次调谐（reconciliation）调用
- 如何配置周期性调谐
- 即将推出
    - 何时使用 lister 缓存 vs 实时查询
    - 垃圾回收 vs Finalizer
    - 如何使用声明式 vs Webhook 校验
    - 如何实现 API 版本化

## 为什么选择 Kubernetes API

Kubernetes API 为遵循一致且丰富结构的对象提供了一致且定义良好的端点。

这种方式催生了用于处理 Kubernetes API 的丰富工具与库生态。

用户通过将对象声明为 *YAML* 或 *JSON* 配置，并使用通用工具来管理这些对象，从而与 API 交互。

将服务构建为 Kubernetes API 相较于传统 REST 具有诸多优势，包括：

* 托管的 API 端点、存储与校验。
* 丰富的工具与 CLI，例如 `kubectl` 与 `kustomize`。
* 支持认证（AuthN）与细粒度授权（AuthZ）。
* 通过 API 版本化与转换支持 API 演进。
* 便于构建自适应/自愈的 API，能够在无需用户干预的情况下持续响应系统状态变化。
* 以 Kubernetes 作为托管运行环境。

开发者可以构建并发布自己的 Kubernetes API，以安装到正在运行的 Kubernetes 集群中。

## 贡献

如果你希望为本书或代码做出贡献，请先阅读我们的[贡献指南](https://github.com/kubernetes-sigs/kubebuilder/blob/master/CONTRIBUTING.md)。

## 资源

* 代码仓库：[sigs.k8s.io/kubebuilder](https://sigs.k8s.io/kubebuilder)

* Slack 频道：[#kubebuilder](http://slack.k8s.io/#kubebuilder)

* Google 讨论组：
  [kubebuilder@googlegroups.com](https://groups.google.com/forum/#!forum/kubebuilder)
