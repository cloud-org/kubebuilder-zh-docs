# 最佳实践（Good Practices）

## 什么是 Operator 的 “Reconciliation”？

用 Kubebuilder 创建项目后，你会在 `cmd/main.go` 看到脚手架代码。该代码初始化一个 [Manager][controller-runtime-manager]，项目基于 [controller-runtime][controller-runtime] 框架。Manager 管理若干[控制器][controllers]，每个控制器提供 reconcile 函数，使资源在集群中不断向期望状态收敛。

“Reconciliation（调谐）” 是一个持续循环，按 Kubernetes 的[控制回路][k8s-control-loop]原理执行必要操作以维持期望状态。更多背景可参考 [Operator 模式][k8s-operator-pattern] 文档。

## 为什么调谐应具备幂等性？

开发 Operator 时，控制器的调谐循环需要是幂等的。遵循[Operator 模式][operator-pattern]，我们实现的[控制器][controllers]应能在集群中不断同步资源直至达到期望状态。幂等的设计有助于正确应对通用或意外事件、顺利处理启动与升级。更多说明见[此处][controller-runtime-topic]。

将调谐逻辑严格绑定到特定事件会违背 Operator 模式与 [controller-runtime][controller-runtime] 的设计原则，可能导致资源卡死、需要人工介入等问题。

## 理解 Kubernetes API 并遵循 API 约定

构建 Operator 通常涉及扩展 Kubernetes API。理解 CRD 与 API 的交互方式至关重要。建议阅读 [Kubebuilder 文档][docs] 中的 Group/Version/Kind 章节，以及 Kubernetes 的 [Operator 模式][operator-pattern] 文档。

## 为什么要遵循 Kubernetes API 约定与标准

遵循 [Kubernetes API 约定与标准][k8s-api-conventions] 对应用与部署至关重要：

- 互操作性：遵循约定可减少兼容性问题，带来一致体验；
- 可维护性：一致的模式/结构便于调试与支持，提高效率；
- 发挥平台能力：在标准框架下更好地利用特性，实现可扩展与高可用；
- 面向未来：与生态演进保持一致，兼容后续更新与特性。

总之，遵循这些约定能显著提升集成、维护、性能与演进能力。

## 为何应避免一个控制器同时管理多个 CRD（例如 “install_all_controller.go”）？

避免让同一个控制器调谐多个 Kind。这通常违背 controller-runtime 的设计，也损害封装、单一职责与内聚性等原则，增加扩展/复用/维护难度。问题包括：

- 复杂性：单控多 CR 会显著增加代码复杂度；
- 可扩展性：易成为瓶颈，降低系统效率与响应性；
- 单一职责：每个控制器聚焦一个职责更稳健；
- 错误隔离：单控多 CR 时，一处错误可能影响所有受管 CR；
- 并发与同步：多 CR 并行易引发竞态与复杂同步（尤其存在依赖关系时）。

因此，通常遵循单一职责：一个 CR 对应一个控制器。

## 推荐使用 Status Conditions

建议按 [K8s API 约定][k8s-api-conventions] 使用 Status Conditions 管理状态，原因包括：

- 标准化：为自定义资源提供统一的状态表示，便于人和工具理解；
- 可读性：多 Condition 组合可表达复杂状态；
- 可扩展：新增特性/状态时易于扩展，而无需重构 API；
- 可观测：便于运维/监控工具跟踪资源状态；
- 兼容性：与生态一致，带来一致的使用体验。

<aside class="note">
<h1> 使用示例 </h1>

可参考 [Deploy Image 插件][deploy-image]：它遵循最佳实践为集群部署并管理镜像（Operand），屏蔽实现复杂度且支持自定义生成代码。其脚手架的 API 与控制器中的调谐逻辑展示了 Status Conditions 的具体用法。

</aside>

[docs]: /cronjob-tutorial/gvks.html
[operator-pattern]: https://kubernetes.io/docs/concepts/extend-kubernetes/operator/
[controllers]: https://kubernetes.io/docs/concepts/architecture/controller/
[controller-runtime-topic]: https://github.com/kubernetes-sigs/controller-runtime/blob/main/FAQ.md#q-how-do-i-have-different-logic-in-my-reconciler-for-different-types-of-events-eg-create-update-delete
[controller-runtime]: https://github.com/kubernetes-sigs/controller-runtime
[deploy-image]: /plugins/available/deploy-image-plugin-v1-alpha.md
[controller-runtime-manager]: https://github.com/kubernetes-sigs/controller-runtime/blob/304027bcbe4b3f6d582180aec5759eb4db3f17fd/pkg/manager/manager.go#L53
[k8s-api-conventions]: https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md
[k8s-control-loop]: https://kubernetes.io/docs/concepts/architecture/controller/
[k8s-operator-pattern]: https://kubernetes.io/docs/concepts/extend-kubernetes/operator/
