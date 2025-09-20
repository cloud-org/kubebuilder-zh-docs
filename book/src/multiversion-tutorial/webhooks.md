# 设置 webhooks

转换逻辑已就绪，剩下的就是让 controller-runtime 知道我们的转换。

## Webhook 设置……

{{#literatego ./testdata/project/internal/webhook/v1/cronjob_webhook.go}}

## ……以及 `main.go`

同样，我们现有的 main 文件也足够了：

{{#literatego ./testdata/project/cmd/main.go}}

一切就绪！接下来就是测试我们的 webhooks。
