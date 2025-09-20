# “轮毂与辐条”以及其他轮子隐喻

现在我们有两个不同的版本，用户可以请求任意一个版本，因此我们必须定义一种在版本之间进行转换的方式。对于 CRD，这通过 webhook 来完成，类似于我们在[基础教程](/cronjob-tutorial/webhook-implementation.md)中定义的 defaulting 与 validating webhooks。与之前一样，controller-runtime 会帮助我们把细节串起来，我们只需要实现实际的转换逻辑。

不过在这之前，我们需要先理解 controller-runtime 如何看待版本，具体来说：

## 完全图可不够“航海”

一种简单的定义转换的方法是：为每一对版本之间都定义转换函数。然后，每当需要转换时，我们查找相应函数并调用它来完成转换。

当只有两个版本时这么做没问题，但如果有 4 种类型呢？8 种类型呢？那将会有很多很多转换函数。

因此，controller-runtime 用“轮毂-辐条（hub-and-spoke）”模型来表示转换：我们将某一个版本标记为“hub”，其他所有版本只需定义到该 hub 的转换以及从该 hub 的转换：

<div class="diagrams">
{{#include ./complete-graph-8.svg}}
<div>变为</div>
{{#include ./hub-spoke-graph.svg}}
</div>

随后，如果需要在两个非 hub 版本之间进行转换，我们先转换到 hub 版本，再转换到目标版本：

<div class="diagrams">
{{#include ./conversion-diagram.svg}}
</div>

这减少了我们需要定义的转换函数数量，并且该模型参考了 Kubernetes 内部的做法。

## 这和 Webhook 有什么关系？

当 API 客户端（如 kubectl 或你的控制器）请求你的资源的某个特定版本时，Kubernetes API server 需要返回该版本的结果。然而，该版本可能与 API server 存储的版本不一致。

这种情况下，API server 需要知道如何在期望版本与存储版本之间进行转换。由于 CRD 的转换不是内建的，Kubernetes API server 会调用一个 webhook 来完成转换。对于 Kubebuilder，这个 webhook 由 controller-runtime 实现，它执行我们上面讨论的 hub-and-spoke 转换。

现在转换模型已经明晰，我们就可以实际实现转换了。
