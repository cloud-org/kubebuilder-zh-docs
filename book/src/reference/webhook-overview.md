# Webhook 概览

Webhook 是一种阻塞式的 HTTP 回调机制：当特定事件发生时，实现了 Webhook 的系统会向目标端发送 HTTP 请求并等待响应。

在 Kubernetes 生态中，主要存在三类 Webhook：
- [Admission Webhook](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/#admission-webhooks)
- [Authorization Webhook](https://kubernetes.io/docs/reference/access-authn-authz/webhook/)
- [CRD Conversion Webhook](https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definition-versioning/#webhook-conversion)

[controller-runtime](https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/webhook?tab=doc)
库目前支持 Admission Webhook 与 CRD Conversion Webhook。

Kubernetes 自 1.9（beta）起支持动态 Admission Webhook；自 1.15（beta）起支持 Conversion Webhook。
