# 常见问题（FAQ）

<aside class="note">
<h1> Controller-Runtime 常见问题 </h1>

Kubebuilder 构建于 [controller-runtime](https://github.com/kubernetes-sigs/controller-runtime) 和 [controller-tools](https://github.com/kubernetes-sigs/controller-tools) 之上。我们也建议你查看 [Controller-Runtime FAQ 页面](https://github.com/kubernetes-sigs/controller-runtime/blob/main/FAQ.md)。
</aside>


## 在初始化项目时通过 domain 参数传入的值（例如 `kubebuilder init --domain example.com`）有什么作用？

创建项目后，通常你会希望扩展 Kubernetes API，并定义由你的项目拥有的新 API。因此，该 domain 值会被记录在定义项目配置的 [PROJECT][project-file-def] 文件中，并作为域名用于创建 API 端点。请确保你理解[Groups、Versions 与 Kinds，哇哦！][gvk] 中的概念。

domain 用于作为 group 的后缀，用来直观地表示资源组的类别。例如，如果设置了 `--domain=example.com`：
```
kubebuilder init --domain example.com --repo xxx --plugins=go/v4
kubebuilder create api --group mygroup --version v1beta1 --kind Mykind
```
那么最终的资源组将是 `mygroup.example.com`。

> 如果没有设置 domain 字段，默认值为 `my.domain`。

## 我想自定义项目使用 [klog][klog]，而不是 controller-runtime 提供的 [zap][zap]。如何将 `klog` 或其他 logger 用作项目的日志器？

在 `main.go` 中你可以把：
```go
    opts := zap.Options{
    Development: true,
    }
    opts.BindFlags(flag.CommandLine)
    flag.Parse()

    ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))
```
替换为：
```go
    flag.Parse()
	ctrl.SetLogger(klog.NewKlogr())
```

## 执行 `make run` 后，我看到类似 “unable to find leader election namespace: not running in-cluster...” 的错误

你可以启用 leader election。不过，如果你在本地使用 `make run` 目标测试项目（该命令会让 manager 在集群外运行），那么你可能还需要设置创建 leader election 资源的命名空间，如下所示：
```go
mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                  scheme,
		MetricsBindAddress:      metricsAddr,
		Port:                    9443,
		HealthProbeBindAddress:  probeAddr,
		LeaderElection:          enableLeaderElection,
		LeaderElectionID:        "14be1926.testproject.org",
		LeaderElectionNamespace: "<project-name>-system",
```

如果你在集群中通过 `make deploy` 目标运行项目，则可能不希望添加此选项。因此，你可以使用环境变量自定义该行为，仅在开发时添加此选项，例如：

```go
    leaderElectionNS := ""
	if os.Getenv("ENABLE_LEADER_ELECTION_NAMESPACE") != "false" {
		leaderElectionNS = "<project-name>-system"
	}

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                  scheme,
		MetricsBindAddress:      metricsAddr,
		Port:                    9443,
		HealthProbeBindAddress:  probeAddr,
		LeaderElection:          enableLeaderElection,
		LeaderElectionNamespace: leaderElectionNS,
		LeaderElectionID:        "14be1926.testproject.org",
		...
```

## 在旧版本 Kubernetes 上部署项目时遇到错误 “open /var/run/secrets/kubernetes.io/serviceaccount/token: permission denied”，如何解决？

如果你遇到如下错误：
```
1.6656687258729894e+09  ERROR   controller-runtime.client.config        unable to get kubeconfig        {"error": "open /var/run/secrets/kubernetes.io/serviceaccount/token: permission denied"}
sigs.k8s.io/controller-runtime/pkg/client/config.GetConfigOrDie
        /go/pkg/mod/sigs.k8s.io/controller-runtime@v0.13.0/pkg/client/config/config.go:153
main.main
        /workspace/main.go:68
runtime.main
        /usr/local/go/src/runtime/proc.go:250
```
当你在 Kubernetes 较旧版本（可能 <= 1.21）上运行项目时，这可能由[该问题][permission-issue]导致，原因是挂载的 token 文件权限为 `0600`，解决方案见[此 PR][permission-PR]。临时解决办法是：

在 manager.yaml 中添加 `fsGroup`：
```yaml
securityContext:
        runAsNonRoot: true
        fsGroup: 65532 # 添加该 fsGroup 以使 token 文件可读
```
不过请注意，该问题已被修复；若你将项目部署在更高版本（可能 >= 1.22），则不会出现此问题。

## 运行 `make install` 应用 CRD 清单时出现 `Too long: must have at most 262144 bytes` 错误。如何解决？为什么会出现该错误？

尝试运行 `make install` 应用 CRD 清单时，可能会遇到 `Too long: must have at most 262144 bytes` 错误。该错误源于 Kubernetes API 实施的大小限制。注意：`make install` 目标会使用 `kubectl apply -f -` 应用 `config/crd` 下的 CRD 清单。因此，当使用 apply 命令时，API 会为对象添加包含完整先前配置的 `last-applied-configuration` 注解。如果该配置过大，就会超出允许的字节大小。（[更多信息][k8s-obj-creation]）

理想情况下，使用 client-side apply 看似完美，因为不需要把完整对象配置作为注解（last-applied-configuration）存储在服务端。然而，需要注意的是，目前 controller-gen 与 kubebuilder 尚不支持该特性。更多内容参见：[Controller-tool 讨论][controller-tool-pr]。

因此，你可以使用以下方式之一来规避该问题：

**移除 CRD 中的描述（description）：**

你的 CRD 是由 [controller-gen][controller-gen] 生成的。通过使用 `maxDescLen=0` 选项来移除描述，可以减小大小，从而可能解决该问题。为此，你可以按以下示例修改 Makefile，然后调用 `make manifest` 目标以重新生成不包含描述的 CRD，如下所示：

```shell

 .PHONY: manifests
 manifests: controller-gen ## Generate WebhookConfiguration, ClusterRole and CustomResourceDefinition objects.
     # 注意：在默认脚手架中加入了 maxDescLen=0 选项以解决 “Too long: must have at most 262144 bytes” 问题。
     # 使用 kubectl apply 创建/更新资源时，K8s API 会创建注解以存储资源的最新版本（kubectl.kubernetes.io/last-applied-configuration）。
     # 该注解有大小限制，如果 CRD 过大且描述很多，就会导致失败。
	$(CONTROLLER_GEN) rbac:roleName=manager-role crd:maxDescLen=0 webhook paths="./..." output:crd:artifacts:config=config/crd/bases
```
**重新设计你的 API：**

你可以审视 API 的设计，看看是否违反了例如单一职责原则而导致规格过多，从而考虑对其进行重构。

## 如何高效地校验和解析 CRD 中的字段？

为提升用户体验，编写 CRD 时建议使用 [OpenAPI v3 schema](https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.0.md#schemaObject) 进行校验。不过，这种方式有时需要额外的解析步骤。
例如，考虑如下代码：
```go
type StructName struct {
	// +kubebuilder:validation:Format=date-time
	TimeField string `json:"timeField,omitempty"`
}
```

### 这种情况下会发生什么？

- 如果用户尝试以非法的 timeField 值创建 CRD，Kubernetes API 会返回错误提示。
- 对于开发者，字符串值在使用前需要手动解析。

### 有更好的方式吗？

为了同时提供更好的用户体验与更顺畅的开发体验，建议使用诸如 [`metav1.Time`](https://pkg.go.dev/k8s.io/apimachinery@v0.31.1/pkg/apis/meta/v1#Time) 这样的预定义类型。例如：
```go
type StructName struct {
	TimeField metav1.Time `json:"timeField,omitempty"`
}
```

### 这种情况下会发生什么？

- 对非法的 `timeField` 值，用户仍会从 Kubernetes API 获得错误提示。
- 开发者可以直接使用已解析的 TimeField，而无需额外解析，从而降低错误并提升效率。



[k8s-obj-creation]: https://kubernetes.io/docs/tasks/manage-kubernetes-objects/declarative-config/#how-to-create-objects
[gvk]: ./cronjob-tutorial/gvks.md
[project-file-def]: ./reference/project-config.md
[klog]: https://github.com/kubernetes/klog
[zap]: https://github.com/uber-go/zap
[permission-issue]: https://github.com/kubernetes/kubernetes/issues/82573
[permission-PR]: https://github.com/kubernetes/kubernetes/pull/89193
[controller-gen]: ./reference/controller-gen.html
[controller-tool-pr]: https://github.com/kubernetes-sigs/controller-tools/pull/536
