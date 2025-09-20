# Kubebuilder v1 与 v2 对比（Legacy：从 v1.0.0+ 到 v2.0.0）

本文覆盖从 v1 迁移到 v2 的所有破坏性变更。

更多（含非破坏性）变更详情可参考
[controller-runtime](https://github.com/kubernetes-sigs/controller-runtime/releases)、
[controller-tools](https://github.com/kubernetes-sigs/controller-tools/releases)
与 [kubebuilder](https://github.com/kubernetes-sigs/kubebuilder/releases)
的发布说明。

## 共同变化（Common changes）

v2 项目改用 Go Modules；在 Go 1.13 发布前，kubebuilder 仍兼容 `dep`。

## controller-runtime

- `Client.List` 采用函数式可选项（`List(ctx, list, ...option)`）替代 `List(ctx, ListOptions, list)`。
- `Client.DeleteAllOf` 新增至 `Client` 接口。

- 指标（metrics）默认开启。

- `pkg/runtime` 下部分包位置已调整，旧位置标记为弃用，并会在 controller-runtime v1.0.0 前移除。详见 [godocs][pkg-runtime-godoc]。

#### 与 Webhook 相关

- 移除 Webhook 的自动证书生成与自注册。请使用 controller-tools 生成 Webhook 配置；若需证书生成，推荐使用
[cert-manager](https://github.com/cert-manager/cert-manager)。Kubebuilder v2 会为你脚手架出 cert-manager 的配置，详见
[Webhook 教程](/cronjob-tutorial/webhook-implementation.md)。

- `builder` 包现在为控制器与 Webhook 分别提供构造器，便于选择运行内容。

## controller-tools

v2 重写了生成器框架。大多数场景下用法不变，但也存在破坏性变更。详见[标记文档](/reference/markers.md)。

## Kubebuilder

- v2 引入更简化的项目布局。设计文档见
https://github.com/kubernetes-sigs/kubebuilder/blob/master/designs/simplified-scaffolding.md

- v1 中 manager 以 `StatefulSet` 部署；v2 中改为 `Deployment`。

- 新增 `kubebuilder create webhook` 命令用于脚手架 变更/校验/转换 Webhook，替代 `kubebuilder alpha webhook`。

- v2 使用 `distroless/static` 作为基础镜像（替代 Ubuntu），以降低镜像体积与攻击面。

- v2 需要 kustomize v3.1.0+。

[LeaderElectionRunable]: https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/manager?tab=doc#LeaderElectionRunnable
[pkg-runtime-godoc]: https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/runtime?tab=doc
