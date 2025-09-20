# 做些改动

Kubernetes API 中一个相当常见的变更是：把原先非结构化或以特殊字符串格式存储的数据，改为结构化数据。我们的 `schedule` 字段就非常符合这一点——当前在 `v1` 中，它长这样：

```yaml
schedule: "*/1 * * * *"
```

这是一个教科书式的“特殊字符串格式”的例子（除非你是 Unix 管理员，否则可读性不佳）。

让我们把它变得更结构化一些。依据我们的 [CronJob 代码][cronjob-sched-code]，我们支持“标准”的 Cron 格式。

在 Kubernetes 中，所有版本之间必须能够安全地往返转换。这意味着，如果我们从版本 1 转换到版本 2，再转换回版本 1，就不能丢失信息。因此，我们对 API 做的任何变更都必须与 v1 所支持的内容兼容，同时还需要确保在 v2 中新增的任何内容在 v1 中也能得到支持。在某些情况下，这意味着需要向 v1 添加新字段，但在我们的场景中，由于没有新增功能，因此无需这么做。

牢记上述要求，我们把上面的例子转换为略微更结构化的形式：

```yaml
schedule:
  minute: */1
```

现在，至少每个字段都有了标签，同时仍然可以轻松支持每个字段的不同语法。

为完成此变更，我们需要一个新的 API 版本，就叫它 v2：

```shell
kubebuilder create api --group batch --version v2 --kind CronJob
```

在 “Create Resource” 处选择 `y`，在 “Create Controller” 处选择 `n`。

现在，复制现有类型，然后做相应修改：

{{#literatego ./testdata/project/api/v2/cronjob_types.go}}

## 存储版本（Storage Versions）

{{#literatego ./testdata/project/api/v1/cronjob_types.go}}

既然类型已经就位，接下来我们需要设置转换……

[cronjob-sched-code]: ./multiversion-tutorial/testdata/project/api/v2/cronjob_types.go "CronJob Code"
