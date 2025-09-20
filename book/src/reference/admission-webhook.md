# Admission Webhooks

[Admission Webhook](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/#what-are-admission-webhooks) 是一种用于接收并处理准入请求的 HTTP 回调，返回相应的准入响应。

Kubernetes 提供两类 Admission Webhook：

- **Mutating Admission Webhook（变更型）**：
  在对象被持久化前（创建或更新时）修改对象。常用于为资源设置默认值（例如为用户未指定的 Deployment 字段赋默认值），或注入 sidecar 容器。

- **Validating Admission Webhook（校验型）**：
  在对象被持久化前（创建或更新时）进行校验。它能实现比纯 schema 校验更复杂的逻辑，例如跨字段校验或镜像白名单等。

默认情况下，apiserver 不会对 Webhook 端进行自我认证。如果你需要在 Webhook 侧认证 apiserver，可以配置 apiserver 使用 Basic Auth、Bearer Token 或证书进行认证。详见
[官方文档](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/#authenticate-apiservers)。

<aside class="note">
<H1>执行顺序</H1>

**Validating Webhook 总是在所有 Mutating Webhook 之后执行**，因此无需担心在你的校验通过后，另一个 Mutating Webhook 再次更改了对象。

</aside>

## 在 Admission Webhook 中处理资源 Status

<aside class="warning">
<H1>不要在 Mutating Webhook 中修改 status</H1>

**不能使用 Mutating Admission Webhook 来修改或设置资源的 status。**
当控制器第一次观察到一个新对象时，应在控制器内设置其初始 status。

</aside>

### 原因说明

#### Mutating Admission Webhook 的职责

Mutating Webhook 主要用于拦截并修改关于对象创建、变更或删除的请求。尽管它可以修改对象的规范（spec），但直接修改 status 并非标准做法，且常常带来意外结果。

```go
// MutatingWebhookConfiguration 允许修改对象
// 但直接修改 status 可能导致非预期行为
type MutatingWebhookConfiguration struct {
    ...
}
```

#### 设置初始 Status

对于自定义控制器而言，理解“设置初始 status”的概念至关重要。该初始化通常在控制器内部完成：当控制器（通常通过 watch）发现受管资源的新实例时，应由控制器赋予该资源一个初始的 status。

```go
// 自定义控制器的调谐函数示例
func (r *ReconcileMyResource) Reconcile(request reconcile.Request) (reconcile.Result, error) {
    // ...
    // 发现新实例时设置初始 status
    instance.Status = SomeInitialStatus
    // ...
}
```

#### Status 子资源

在 Kubernetes 的自定义资源中，spec（期望状态）与 status（观察状态）是明确分离的。为 CRD 启用 /status 子资源会将 `status` 与 `spec` 分离到各自的 API 端点。
这保证了用户发起的修改（例如更新 spec）与系统驱动的变更（例如更新 status）互不干扰。因此，试图在一次修改 spec 的操作中利用 Mutating Webhook 去更改 status，往往不会得到预期结果。

```yaml
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: myresources.mygroup.mydomain
spec:
  ...
  subresources:
    status: {} # 启用 /status 子资源
```

#### 结论

虽然在某些极端场景下 Mutating Webhook 似乎能顺带修改 status，但这既不通用，也不被推荐。将 status 更新逻辑放在控制器中处理，仍是最佳实践。
