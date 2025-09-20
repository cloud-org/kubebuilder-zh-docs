# Groups、Versions 与 Kinds，哇哦！

在开始我们的 API 之前，先简单聊聊术语。

在 Kubernetes 中谈论 API 时，我们经常用到四个术语：groups、versions、kinds 和 resources。

## Groups 与 Versions

Kubernetes 中的 API Group 只是相关功能的一个集合。每个 group 拥有一个或多个 version，顾名思义，它们允许我们随时间改变 API 的工作方式。

## Kinds 与 Resources

每个 API group-version 包含一个或多个 API 类型，我们称之为 Kind。一个 Kind 可以在不同版本间改变其形式，但每种形式都必须能够以某种方式存储其他形式的全部数据（我们可以把数据放在字段里，或是放在注解中）。这意味着使用较旧的 API 版本不会导致较新的数据丢失或损坏。更多信息参见 [Kubernetes API 指南](https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md)。

你也会偶尔听到 resource 这个词。resource 只是某个 Kind 在 API 中的一种使用。通常，Kind 与 resource 是一一对应的。例如，`pods` 这个 resource 对应 `Pod` 这个 Kind。然而，有时相同的 Kind 可能由多个 resource 返回。例如，`Scale` 这个 Kind 由所有的 scale 子资源返回，如 `deployments/scale` 或 `replicasets/scale`。这正是 Kubernetes HorizontalPodAutoscaler 能与不同资源交互的原因。不过，对于 CRD，每个 Kind 只会对应单个 resource。

请注意，resource 总是小写，并且按照惯例是 Kind 的小写形式。

## 那这与 Go 如何对应？

当我们引用某个特定 group-version 下的 kind 时，我们称之为 GroupVersionKind，简称 GVK。资源的情况类似，简称 GVR。正如我们很快会看到的，每个 GVK 都对应包中的某个根 Go 类型。

现在术语已经讲清，我们终于可以创建 API 了！

## 那我们如何创建 API？

在下一节 [Adding a new API](../cronjob-tutorial/new-api.html) 中，我们将看看 `kubebuilder create api` 这条命令是如何帮助我们创建自定义 API 的。

该命令的目标是为我们的 Kind 创建 Custom Resource（CR）与 Custom Resource Definition（CRD）。欲了解更多，请参见：[使用 CustomResourceDefinition 扩展 Kubernetes API][kubernetes-extend-api]。

## 但为什么要创建 API 呢？

新的 API 是我们让 Kubernetes 理解自定义对象的方式。Go 结构体被用于生成 CRD，CRD 包含了我们数据的 schema，以及诸如新类型叫什么之类的跟踪信息。随后我们就可以创建自定义对象的实例，它们将由我们的[控制器][controllers]进行管理。

我们的 API 与 resource 代表了我们在集群中的解决方案。基本上，CRD 是对自定义对象的定义，而 CR 则是其一个实例。

## 有个例子吗？

想象一个经典场景：我们的目标是让一个应用及其数据库在 Kubernetes 平台上运行。那么，一个 CRD 可以表示 App，另一个 CRD 可以表示 DB。用一个 CRD 描述 App、另一个 CRD 描述 DB，不会破坏封装、单一职责和内聚性等概念。破坏这些概念可能会带来意想不到的副作用，比如难以扩展、复用或维护，等等。

这样，我们可以创建一个 App 的 CRD，对应的控制器负责创建包含该 App 的 Deployment、为其创建可访问的 Service 等。同理，我们可以创建一个表示 DB 的 CRD，并部署一个控制器来管理 DB 实例。

## 呃，那 Scheme 又是什么？

我们之前看到的 `Scheme` 只是用来跟踪某个 GVK 对应哪个 Go 类型的一种方式（不要被它的 [godocs](https://pkg.go.dev/k8s.io/apimachinery/pkg/runtime?tab=doc#Scheme) 吓到）。

例如，假设我们将 `"tutorial.kubebuilder.io/api/v1".CronJob{}` 标记为属于 `batch.tutorial.kubebuilder.io/v1` API 组（这也就隐含了它的 Kind 是 `CronJob`）。

随后，当 API server 返回如下 JSON 时，我们就能据此构造一个新的 `&CronJob{}`：

```json
{
    "kind": "CronJob",
    "apiVersion": "batch.tutorial.kubebuilder.io/v1",
    ...
}
```

或者在我们提交 `&CronJob{}` 更新时，正确地查找其 group-version。

[kubernetes-extend-api]: https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definitions/
[controllers]: ../cronjob-tutorial/controller-overview.md
