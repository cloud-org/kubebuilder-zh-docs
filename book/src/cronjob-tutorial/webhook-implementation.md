# 实现 defaulting/validating Webhook

如果你想为 CRD 实现[准入 Webhook](../reference/admission-webhook.md)，你只需要实现 `CustomDefaulter` 和/或 `CustomValidator` 接口。

其余工作由 Kubebuilder 替你完成，例如：

1. 创建 webhook 服务器
1. 确保服务器已添加进 manager
1. 为你的 webhooks 创建处理器
1. 将每个处理器注册到服务器上的某条路径

首先，我们为 CRD（CronJob）搭建 webhook 脚手架。由于测试项目会使用 defaulting 与 validating webhooks，我们需要带上 `--defaulting` 与 `--programmatic-validation` 参数执行以下命令：

```bash
kubebuilder create webhook --group batch --version v1 --kind CronJob --defaulting --programmatic-validation
```

这会为你生成 webhook 函数，并在 `main.go` 中把你的 webhook 注册到 manager。

{{#literatego ./testdata/project/internal/webhook/v1/cronjob_webhook.go}}
