# 实现转换

转换模型确定后，就该实际实现转换函数了。我们将为 CronJob API 的 `v1`（Hub）到 `v2`（Spoke）创建一个转换 webhook，见：

```go
kubebuilder create webhook --group batch --version v1 --kind CronJob --conversion --spoke v2
```

上述命令会在 `cronjob_types.go` 旁边生成 `cronjob_conversion.go` 文件，以避免在主类型文件中堆积额外函数。

## Hub...

首先实现 hub。我们选择 v1 作为 hub：

{{#literatego ./testdata/project/api/v1/cronjob_conversion.go}}

## ...以及 Spokes

然后实现 spoke，即 v2 版本：

{{#literatego ./testdata/project/api/v2/cronjob_conversion.go}}

现在转换已经就位，我们只需要把 main 连起来以提供该 webhook！
