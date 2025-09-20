# 设计一个 API

在 Kubernetes 中，我们在设计 API 时有一些规则。具体来说，所有序列化字段必须使用 `camelCase`，因此我们通过 JSON 结构体标签来指定这一点。我们也可以使用 `omitempty` 结构体标签在字段为空时省略序列化。

字段可以使用大多数原始类型。数字是个例外：出于 API 兼容性考虑，我们接受三种数字形式：用于整数的 `int32` 与 `int64`，以及用于小数的 `resource.Quantity`。

<details><summary>等等，什么是 Quantity？</summary>

Quantity 是一种用于小数的特殊表示法，具有明确固定的表示，使其在不同机器之间更具可移植性。你很可能在 Kubernetes 中为 Pod 指定资源请求与限制时见过它。

从概念上看，它类似于浮点数：包含有效数、基数和指数。其可序列化且便于阅读的人类可读格式使用整数与后缀来表示数值，就像我们描述计算机存储的方式一样。

例如，`2m` 在十进制表示中等于 `0.002`。`2Ki` 在十进制中表示 `2048`，而 `2K` 在十进制中表示 `2000`。如果我们需要表示小数部分，可以切换到允许使用整数的后缀：`2.5` 可写作 `2500m`。

支持两种基：10 和 2（分别称为十进制与二进制）。十进制基使用“常规”的 SI 后缀（例如 `M` 与 `K`），而二进制基使用 “mebi” 表示法（例如 `Mi` 与 `Ki`）。可参见 [megabytes vs mebibytes](https://en.wikipedia.org/wiki/Binary_prefix)。

</details>

还有一个我们会用到的特殊类型：`metav1.Time`。它与 `time.Time` 的功能相同，但具有固定且可移植的序列化格式。

介绍到这里，让我们看看 CronJob 对象长什么样！

{{#literatego ./testdata/project/api/v1/cronjob_types.go}}

现在我们已有 API，需要编写一个控制器来真正实现其功能。
