# 部署 cert-manager

我们建议使用 [cert-manager](https://github.com/cert-manager/cert-manager) 为 Webhook 服务器签发证书。只要能把证书放到期望的位置，其他方案也同样可行。

你可以按照 [cert-manager 文档](https://cert-manager.io/docs/installation/) 进行安装。

cert-manager 还有一个名为 [CA Injector](https://cert-manager.io/docs/concepts/ca-injector/) 的组件，负责将 CA bundle 注入到 [`MutatingWebhookConfiguration`](https://pkg.go.dev/k8s.io/api/admissionregistration/v1#MutatingWebhookConfiguration) / [`ValidatingWebhookConfiguration`](https://pkg.go.dev/k8s.io/api/admissionregistration/v1#ValidatingWebhookConfiguration) 中。

为实现这一点，你需要在 [`MutatingWebhookConfiguration`](https://pkg.go.dev/k8s.io/api/admissionregistration/v1#MutatingWebhookConfiguration) / [`ValidatingWebhookConfiguration`](https://pkg.go.dev/k8s.io/api/admissionregistration/v1#ValidatingWebhookConfiguration) 对象上添加键为 `cert-manager.io/inject-ca-from` 的注解。该注解的值应指向一个已存在的 [certificate request 实例](https://cert-manager.io/docs/concepts/certificaterequest/)，格式为 `<certificate-namespace>/<certificate-name>`。

下面是我们用于给 [`MutatingWebhookConfiguration`](https://pkg.go.dev/k8s.io/api/admissionregistration/v1#MutatingWebhookConfiguration) / [`ValidatingWebhookConfiguration`](https://pkg.go.dev/k8s.io/api/admissionregistration/v1#ValidatingWebhookConfiguration) 对象添加注解的 [kustomize](https://github.com/kubernetes-sigs/kustomize) 补丁。
