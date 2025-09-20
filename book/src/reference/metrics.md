# 指标（Metrics）

默认情况下，controller-runtime 会构建全局 Prometheus 注册表，并为每个控制器发布[一组性能指标](/reference/metrics-reference.md)。


<aside class="note warning">
<h1>重要：如果你仍在使用 `kube-rbac-proxy`</h1>

请尽快停止使用 `gcr.io/kubebuilder/kube-rbac-proxy` 镜像。若未来无法拉取，该镜像会导致项目受影响甚至不可用。

**`gcr.io/kubebuilder/` 下的镜像将在 2025 年初起不可用。**

- **使用 Kubebuilder `v3.14` 及以下版本创建的项目**通常使用 [kube-rbac-proxy](https://github.com/brancz/kube-rbac-proxy) 保护 metrics 端点。建议升级到最新版本，或应用等效修改。

- **而使用 Kubebuilder `v4.1.0` 及以上创建的项目**默认通过 controller-runtime 的 [WithAuthenticationAndAuthorization](https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.18.4/pkg/metrics/filters#WithAuthenticationAndAuthorization) 启用 `authn/authz` 提供类似保护。

如需继续使用 [kube-rbac-proxy](https://github.com/brancz/kube-rbac-proxy)，必须切换到其他镜像源。

> For further information, see: [kubebuilder/discussions/3907](https://github.com/kubernetes-sigs/kubebuilder/discussions/3907)

</aside>

## 指标配置（Metrics Configuration）

查看 `config/default/kustomization.yaml` 可知默认已暴露 metrics：

```yaml
# [METRICS] Expose the controller manager metrics service.
- metrics_service.yaml
```

```yaml
patches:
   # [METRICS] The following patch will enable the metrics endpoint using HTTPS and the port :8443.
   # More info: https://book.kubebuilder.io/reference/metrics
   - path: manager_metrics_patch.yaml
     target:
        kind: Deployment
```

随后可在 `cmd/main.go` 中查看 metrics server 的配置：

```go
// Metrics endpoint is enabled in 'config/default/kustomization.yaml'. The Metrics options configure the server.
// For more info: https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/metrics/server
Metrics: metricsserver.Options{
   ...
},
```

## 在 Kubebuilder 中消费控制器指标

你可以使用 `curl` 或 Prometheus 等 HTTP 客户端访问控制器暴露的指标。

但在此之前，请确保客户端具备访问 `/metrics` 端点所需的 **RBAC 权限**。

### 授权访问指标端点

Kubebuilder 在如下位置脚手架了一个拥有读取权限的 `ClusterRole`：

```
config/rbac/metrics_reader_role.yaml
```

该文件包含了允许访问 metrics 端点所需的 RBAC 规则。

<aside class="note">
<H1>该 ClusterRole 仅为辅助</H1>

Kubebuilder **不会默认生成 RoleBinding/ClusterRoleBinding**，以避免：

- 误绑到错误的 ServiceAccount；
- 在受限环境中误授访问；
- 在多团队/多租户集群中造成冲突。

</aside>

#### 创建 ClusterRoleBinding

可通过 `kubectl` 创建绑定：

```bash
kubectl create clusterrolebinding metrics \
  --clusterrole=<project-prefix>-metrics-reader \
  --serviceaccount=<namespace>:<service-account-name>
```

或使用清单：

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: allow-metrics-access
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: metrics-reader
subjects:
- kind: ServiceAccount
  name: controller-manager
  namespace: system # Replace 'system' with your controller-manager's namespace
```

<aside class="note">
<H1>为何需手动绑定：</H1>

Kubebuilder 避免默认生成绑定，原因包括：
 - 可能绑定到错误的 ServiceAccount；
 - 可能在不需要时授予访问；
 - 在受限或多租户集群中可能造成问题。
   该设计在安全与灵活间取舍，因此需由使用者手动绑定。
</aside>

### 测试指标端点（通过 Curl Pod）

如需手动测试访问 metrics 端点，可执行：

- 创建 RoleBinding

```bash
kubectl create clusterrolebinding <project-name>-metrics-binding \
  --clusterrole=<project-name>-metrics-reader \
  --serviceaccount=<project-name>-system:<project-name>-controller-manager
```

- 生成 Token

```bash
export TOKEN=$(kubectl create token <project-name>-controller-manager -n <project-name>-system)
echo $TOKEN
```

- Launch Curl Pod

```bash
kubectl run curl-metrics --rm -it --restart=Never \
  --image=curlimages/curl:7.87.0 -n <project-name>-system -- /bin/sh
```

- 调用 Metrics 端点

在 Pod 内使用：

```bash
curl -v -k -H "Authorization: Bearer $TOKEN" \
  https://<project-name>-controller-manager-metrics-service.<project-name>-system.svc.cluster.local:8443/metrics
```

<aside class="note">
<H1>注意</H1>

- 将 `<project-name>`、`<namespace>`、`<service-account-name>` 替换为实际值；
- 如未跳过验证（`-k`），请确保 TLS 已启用且证书有效；

下一节将介绍保护 metrics 端点的可选方案。
</aside>

## 指标保护与可选方案

未加保护的 metrics 端点可能向未授权用户暴露敏感数据（系统性能、应用行为、运维指标等），从而带来安全风险。

### 使用 authn/authz（默认启用）

为降低风险，Kubebuilder 项目通过认证（authn）与鉴权（authz）保护 metrics 端点，确保仅授权用户/服务账号可访问敏感指标。

过去常使用 [kube-rbac-proxy](https://github.com/brancz/kube-rbac-proxy) 进行保护；新版本已不再使用。自 `v4.1.0` 起，项目默认通过 controller-runtime 的 [WithAuthenticationAndAuthorization](https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/metrics/server) 启用并保护 metrics 端点。

因此，你会看到如下配置：

- In the `cmd/main.go`:

```go
if secureMetrics {
  ...
  metricsServerOptions.FilterProvider = filters.WithAuthenticationAndAuthorization
}
```

该配置通过 FilterProvider 对 metrics 端点实施认证与鉴权，确保仅具有相应权限的实体可访问。

- In the `config/rbac/kustomization.yaml`:

```yaml
# The following RBAC configurations are used to protect
# the metrics endpoint with authn/authz. These configurations
# ensure that only authorized users and service accounts
# can access the metrics endpoint.
- metrics_auth_role.yaml
- metrics_auth_role_binding.yaml
- metrics_reader_role.yaml
```

这样，只有使用相应 `ServiceAccount` token 的 Pod 才能读取 metrics。示例：

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: metrics-consumer
  namespace: system
spec:
  # Use the scaffolded service account name to allow authn/authz
  serviceAccountName: controller-manager
  containers:
  - name: metrics-consumer
    image: curlimages/curl:latest
    command: ["/bin/sh"]
    args:
      - "-c"
      - >
        while true;
        do
          # Note here that we are passing the token obtained from the ServiceAccount to curl the metrics endpoint
          curl -s -k -H "Authorization: Bearer $(cat /var/run/secrets/kubernetes.io/serviceaccount/token)"
          https://controller-manager-metrics-service.system.svc.cluster.local:8443/metrics;
          sleep 60;
        done
```

### （推荐）在生产环境启用证书（默认关闭）

<aside class="warning">
<h1>为何默认不启用？</h1>

该选项引入对 CertManager 的依赖。为保持项目轻量与易上手，默认不启用。

</aside>

<aside class="warning">
<h1>生产环境建议启用</h1>

`cmd/main.go` 中默认使用 **controller-runtime** 提供的能力自动生成自签证书以保护 metrics server，便于开发测试；但**不建议**用于生产。

这些证书用于保护传输层（TLS）。默认启用的 `authn/authz` 则承担应用层凭据。当你将指标集成至 Prometheus 等系统时，可使用这些证书加固通信。

</aside>

自 Kubebuilder `4.4.0` 起，脚手架包含使用 [CertManager](https://cert-manager.io/) 管理证书以保护 metrics server 的逻辑。按以下步骤可启用：

1. **在 `config/default/kustomization.yaml` 启用 Cert-Manager**：
    - 取消注释 cert-manager 资源：

      ```yaml
      - ../certmanager
      ```

2. **启用在 `config/default/kustomization.yaml` 中用于挂载证书的 Patch**：
    - 取消注释 `cert_metrics_manager_patch.yaml`，在 Manager 的 Deployment 中挂载 `serving-cert`：

      ```yaml
      # Uncomment the patches line if you enable Metrics and CertManager
      # [METRICS-WITH-CERTS] To enable metrics protected with certManager, uncomment the following line.
      # This patch will protect the metrics with certManager self-signed certs.
      - path: cert_metrics_manager_patch.yaml
        target:
          kind: Deployment
      ```
3. **在 `config/default/kustomization.yaml` 中启用为 Metrics Server 配置证书的 replacements**：
    - 取消注释下方 replacements 块，为 `config/certmanager` 下的证书正确设置 DNS 名称：

      ```yaml
      # [CERTMANAGER] To enable cert-manager, uncomment all sections with 'CERTMANAGER' prefix.
      # Uncomment the following replacements to add the cert-manager CA injection annotations
      #replacements:
      # - source: # Uncomment the following block to enable certificates for metrics
      #     kind: Service
      #     version: v1
      #     name: controller-manager-metrics-service
      #     fieldPath: metadata.name
      #   targets:
      #     - select:
      #         kind: Certificate
      #         group: cert-manager.io
      #         version: v1
      #         name: metrics-certs
      #       fieldPaths:
      #         - spec.dnsNames.0
      #         - spec.dnsNames.1
      #       options:
      #         delimiter: '.'
      #         index: 0
      #         create: true
      #
      # - source:
      #     kind: Service
      #     version: v1
      #     name: controller-manager-metrics-service
      #     fieldPath: metadata.namespace
      #   targets:
      #     - select:
      #         kind: Certificate
      #         group: cert-manager.io
      #         version: v1
      #         name: metrics-certs
      #       fieldPaths:
      #         - spec.dnsNames.0
      #         - spec.dnsNames.1
      #       options:
      #         delimiter: '.'
      #         index: 1
      #         create: true
      #
      ```

4. **在 `config/prometheus/kustomization.yaml` 中启用 `ServiceMonitor` 的证书配置**：
    - 添加或取消注释 `ServiceMonitor` 的 patch，以使用 cert-manager 管理的 Secret 并启用证书校验：

      ```yaml
      # [PROMETHEUS-WITH-CERTS] The following patch configures the ServiceMonitor in ../prometheus
      # to securely reference certificates created and managed by cert-manager.
      # Additionally, ensure that you uncomment the [METRICS WITH CERTMANAGER] patch under config/default/kustomization.yaml
      # to mount the "metrics-server-cert" secret in the Manager Deployment.
      patches:
        - path: monitor_tls_patch.yaml
          target:
            kind: ServiceMonitor
      ```

    > **NOTE** that the `ServiceMonitor` patch above will ensure that if you enable the Prometheus integration,
    it will securely reference the certificates created and managed by CertManager. But it will **not** enable the
    integration with Prometheus. To enable the integration with Prometheus, you need uncomment the `#- ../certmanager`
    in the `config/default/kustomization.yaml`. For more information, see [Exporting Metrics for Prometheus](#exporting-metrics-for-prometheus).

### **(Optional)** By using Network Policy (Disabled by default)

NetworkPolicy acts as a basic firewall for pods within a Kubernetes cluster, controlling traffic
flow at the IP address or port level. However, it doesn't handle `authn/authz`.

Uncomment the following line in the `config/default/kustomization.yaml`:

```
# [NETWORK POLICY] Protect the /metrics endpoint and Webhook Server with NetworkPolicy.
# Only Pod(s) running a namespace labeled with 'metrics: enabled' will be able to gather the metrics.
# Only CR(s) which uses webhooks and applied on namespaces labeled 'webhooks: enabled' will be able to work properly.
#- ../network-policy
```

## Exporting Metrics for Prometheus

使用 Prometheus Operator 导出指标的步骤：

1. 安装 Prometheus 与 Prometheus Operator。
   若无自建监控系统，生产环境建议使用 [kube-prometheus](https://github.com/coreos/kube-prometheus#installing)。
   若仅用于试验，可只安装 Prometheus 与 Prometheus Operator。

2. 在 `config/default/kustomization.yaml` 中取消注释 `- ../prometheus`，以创建 `ServiceMonitor` 并启用指标导出：

```yaml
# [PROMETHEUS] To enable prometheus monitor, uncomment all sections with 'PROMETHEUS'.
- ../prometheus
```

注意：当你将项目安装到集群时会创建 `ServiceMonitor` 用于导出指标。可通过 `kubectl get ServiceMonitor -n <project>-system` 检查，例如：

```
$ kubectl get ServiceMonitor -n monitor-system
NAME                                         AGE
monitor-controller-manager-metrics-monitor   2m8s
```

<aside class="warning">
<h2>使用 Prometheus Operator 时请确保权限完备</h2>

默认情况下，Prometheus Operator 的 RBAC 规则仅在 `default` 与 `kube-system` 命名空间启用。参考其文档了解如何[通过 `.jsonnet` 配置监控其他命名空间](https://github.com/prometheus-operator/kube-prometheus/blob/main/docs/monitoring-other-namespaces.md)。

或通过 RBAC 授予 Prometheus Operator 监控其他命名空间的权限，参见：[为 Prometheus pods 启用 RBAC 规则](https://github.com/prometheus-operator/prometheus-operator/blob/main/Documentation/user-guides/getting-started.md#enable-rbac-rules-for-prometheus-pods)。
</aside>

另外，指标默认通过 `8443` 端口导出。你可以在 Prometheus 控制台中通过 `{namespace="<project>-system"}` 查询该命名空间导出的指标：

<img width="1680" alt="Screenshot 2019-10-02 at 13 07 13" src="https://user-images.githubusercontent.com/7708031/66042888-a497da80-e515-11e9-9d77-d8a9fc1159a5.png">

## 发布自定义指标

如果希望从控制器发布更多指标，可使用 `controller-runtime/pkg/metrics` 的全局注册表。

一种常见方式是在控制器包中将采集器声明为全局变量，并在 `init()` 中注册：

For example:

```go
import (
    "github.com/prometheus/client_golang/prometheus"
    "sigs.k8s.io/controller-runtime/pkg/metrics"
)

var (
    goobers = prometheus.NewCounter(
        prometheus.CounterOpts{
            Name: "goobers_total",
            Help: "Number of goobers processed",
        },
    )
    gooberFailures = prometheus.NewCounter(
        prometheus.CounterOpts{
            Name: "goober_failures_total",
            Help: "Number of failed goobers",
        },
    )
)

func init() {
    // Register custom metrics with the global prometheus registry
    metrics.Registry.MustRegister(goobers, gooberFailures)
}
```

随后可在调谐循环中任意位置对这些采集器写入数据；在 operator 代码中的任意位置均可读取与评估这些指标。

<aside class="note">
<h1>在 Prometheus UI 中查看指标</h1>

要在 Prometheus UI 中查看指标，需要配置 Prometheus 实例根据标签选择相应的 ServiceMonitor。

</aside>

上述指标可被 Prometheus 或其他 OpenMetrics 系统抓取。

![Screen Shot 2021-06-14 at 10 15 59 AM](https://user-images.githubusercontent.com/37827279/121932262-8843cd80-ccf9-11eb-9c8e-98d0eda80169.png)

<aside class="note">
<h1>Controller-Runtime Auth/Authz 已知限制与注意事项</h1>

目前存在一些限制：`cache TTL`、`anonymous access` 与 `timeouts` 等设置为硬编码，无法精细调优，可能引发性能或安全顾虑；缺少诸如关键路径（如 `/healthz`）的 `alwaysAllow` 与 `alwaysAllowGroups`（如 `system:masters`）等配置支持，可能带来运维挑战；对 `kube-apiserver` 的稳定连接依赖较强，在网络不稳时可能出现指标中断，导致关键时刻丢失重要监控数据。

已就此在 controller-runtime 提交[改进议题](https://github.com/kubernetes-sigs/controller-runtime/issues/2781)。
</aside>
