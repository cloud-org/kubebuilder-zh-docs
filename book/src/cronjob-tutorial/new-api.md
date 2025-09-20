# 添加一个新的 API

要为一个新的 Kind（你有在关注[上一章](./gvks.md#kinds-and-resources)吧？）及其对应的控制器搭建脚手架，我们可以使用 `kubebuilder create api`：

```bash
kubebuilder create api --group batch --version v1 --kind CronJob
```

在 “Create Resource” 和 “Create Controller” 处按下 `y`。

对于每个 group-version，第一次调用该命令时会为它创建一个目录。

在当前示例中，会创建 [`api/v1/`](https://sigs.k8s.io/kubebuilder/docs/book/src/cronjob-tutorial/testdata/project/api/v1) 目录，对应 `batch.tutorial.kubebuilder.io/v1`（还记得我们一开始的 [`--domain` 设置](cronjob-tutorial.md#scaffolding-out-our-project) 吗？）。

它还为我们的 `CronJob` Kind 添加了一个文件 `api/v1/cronjob_types.go`。每次用不同的 kind 调用该命令时，都会相应地添加一个新文件。

我们先看看“开箱即用”的内容，然后再继续补全。

{{#literatego ./testdata/emptyapi.go}}

了解了基本结构后，我们来把它填充完整！
