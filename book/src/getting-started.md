# 入门

我们将创建一个示例项目来展示其工作方式。该示例将：

- 调谐一个 Memcached CR——它代表一个在集群中部署/由集群管理的 Memcached 实例
- 使用 Memcached 镜像创建一个 Deployment
- 不允许实例数超过 CR 中定义的 size
- 更新 Memcached CR 的状态

<aside class="note">
<h1>为什么是 Operator？</h1>

遵循[Operator 模式][k8s-operator-pattern]，我们不仅能够提供所有预期的资源，还可以在运行时以编程方式、动态地对它们进行管理。举个例子：如果有人不小心修改了配置或误删了资源，Operator 可以在无人干预的情况下将其修复。

</aside>

<aside class="note">
<h1>跟着做 vs 快速跳转</h1>

需要注意的是，本教程的大部分内容由可文学化的 Go 文件生成，它们构成了一个可运行项目，位于书籍源码目录：
[docs/book/src/getting-started/testdata/project][tutorial-source]。

[tutorial-source]: https://github.com/kubernetes-sigs/kubebuilder/tree/master/docs/book/src/cronjob-tutorial/testdata/project

</aside>

## 创建项目

首先，为你的项目创建并进入一个目录。然后使用 `kubebuilder` 初始化：

```shell
mkdir $GOPATH/memcached-operator
cd $GOPATH/memcached-operator
kubebuilder init --domain=example.com
```

<aside class="note">
<h1>在 $GOPATH 中开发</h1>

如果你的项目在 [`GOPATH`][GOPATH-golang-docs] 中初始化，被隐式调用的 `go mod init` 会为你插入模块路径；否则必须设置 `--repo=<module path>`。

如果对模块系统不熟悉，请阅读 [Go modules 的博文][go-modules-blogpost]。

</aside>

## 创建 Memcached API（CRD）

接下来，我们将创建负责在集群上部署并管理 Memcached 实例的 API。

```shell
kubebuilder create api --group cache --version v1alpha1 --kind Memcached
```

### 理解 API

该命令的主要目标是为 Memcached 这个 Kind 生成自定义资源（CR）与自定义资源定义（CRD）。它会创建 group 为 `cache.example.com`、version 为 `v1alpha1` 的 API，从而唯一标识 Memcached Kind 的新 CRD。借助 Kubebuilder，我们可以定义代表我们在平台上方案的 API 与对象。

虽然本示例中仅添加了一种资源的 Kind，但我们可以根据需要拥有任意数量的 `Group` 与 `Kind`。为便于理解，可以将 CRD 看作自定义对象的“定义”，而 CR 则是其“实例”。

<aside class="note">
<h1> 请确保你查看 </h1>

[Groups、Versions 与 Kinds，哇哦！][group-kind-oh-my]

</aside>

### 定义我们的 API

#### 定义规格（Spec）

现在，我们将定义集群中每个 Memcached 资源实例可以采用的值。在本示例中，我们允许通过以下方式配置实例数量：

```go
type MemcachedSpec struct {
	...
	// +kubebuilder:validation:Minimum=0
	// +required
	Size *int32 `json:"size,omitempty"`
}
```

#### 定义 Status

我们还希望跟踪为管理 Memcached CR 所进行操作的状态。这使我们能够像使用 Kubernetes API 中的任何资源那样，校验自定义资源对我们 API 的描述，并判断一切是否成功，或是否遇到错误。

```go
// MemcachedStatus defines the observed state of Memcached
type MemcachedStatus struct {
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,1,rep,name=conditions"`
}
```

<aside class="note">
<h1> Status Conditions </h1>

Kubernetes 制定了相应约定，因此我们在此使用 Status Conditions。我们希望自定义 API 与控制器像 Kubernetes 资源及其控制器那样工作，遵循这些标准以确保一致、直观的体验。

请务必查看：[Kubernetes API Conventions](https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#typical-status-properties)
</aside>


#### 标记（Markers）与校验

此外，我们希望对自定义资源中的值进行校验以确保其有效。为此，我们将使用[标记][markers]，例如 `+kubebuilder:validation:Minimum=1`。

现在，来看我们完整的示例。

{{#literatego ./getting-started/testdata/project/api/v1alpha1/memcached_types.go}}

#### 生成包含规格与校验的清单

生成所有必需文件：

1. 运行 `make generate` 在 `api/v1alpha1/zz_generated.deepcopy.go` 中生成 DeepCopy 实现。

2. 然后运行 `make manifests` 在 `config/crd/bases` 下生成 CRD 清单，并在 `config/samples` 下生成其示例。

这两个命令都会使用 [controller-gen][controller-gen]，但分别使用不同的参数来生成代码与清单。

<details><summary><code>config/crd/bases/cache.example.com_memcacheds.yaml</code>: Our Memcached CRD</summary>

```yaml
{{#include ./getting-started/testdata/project/config/crd/bases/cache.example.com_memcacheds.yaml}}
```

</details>

#### 自定义资源示例

`config/samples` 目录下的清单是可应用到集群的自定义资源示例。在本例中，将该资源应用到集群会生成一个副本数为 1 的 Deployment（见 `size: 1`）。

```yaml
{{#include ./getting-started/testdata/project/config/samples/cache_v1alpha1_memcached.yaml}}
```

### 调谐（Reconcile）流程

简单来说，Kubernetes 允许我们声明系统的期望状态，然后其控制器会持续观察集群并采取操作，以确保实际状态与期望状态一致。对于我们的自定义 API 与控制器，过程也是类似的。记住：我们是在扩展 Kubernetes 的行为与 API 以满足特定需求。

在控制器中，我们将实现一个调谐流程。

本质上，调谐流程以循环方式工作：持续检查条件并执行必要操作，直到达到期望状态。该流程会一直运行，直到系统中的所有条件与我们的实现所定义的期望状态一致。

下面是一个伪代码示例：

```go
reconcile App {

  // Check if a Deployment for the app exists, if not, create one
  // If there's an error, then restart from the beginning of the reconcile
  if err != nil {
    return reconcile.Result{}, err
  }

  // Check if a Service for the app exists, if not, create one
  // If there's an error, then restart from the beginning of the reconcile
  if err != nil {
    return reconcile.Result{}, err
  }

  // Look for Database CR/CRD
  // Check the Database Deployment's replicas size
  // If deployment.replicas size doesn't match cr.size, then update it
  // Then, restart from the beginning of the reconcile. For example, by returning `reconcile.Result{Requeue: true}, nil`.
  if err != nil {
    return reconcile.Result{Requeue: true}, nil
  }
  ...

  // If at the end of the loop:
  // Everything was executed successfully, and the reconcile can stop
  return reconcile.Result{}, nil

}
```

<aside class="note">
<h1> 返回选项 </h1>

以下是可用于重新开始调谐的一些返回选项：

- 携带错误：

```go
return ctrl.Result{}, err
```
- 不带错误：

```go
return ctrl.Result{Requeue: true}, nil
```

- 因此，要停止调谐：

```go
return ctrl.Result{}, nil
```

- 在 X 时间后再次调谐：

```go
return ctrl.Result{RequeueAfter: nextRun.Sub(r.Now())}, nil
```

</aside>

#### 放到本示例的上下文中

当我们将示例自定义资源（CR）应用到集群（例如 `kubectl apply -f config/sample/cache_v1alpha1_memcached.yaml`）时，我们希望确保会为 Memcached 镜像创建一个 Deployment，且其副本数与 CR 中定义的一致。

为实现这一点，我们首先需要实现一个操作：检查集群中是否已存在该 Memcached 实例对应的 Deployment；如果不存在，控制器将据此创建 Deployment。因此，调谐流程必须包含一项操作来确保该期望状态被持续维持。该操作大致包括：

```go
	// Check if the deployment already exists, if not create a new one
	found := &appsv1.Deployment{}
	err = r.Get(ctx, types.NamespacedName{Name: memcached.Name, Namespace: memcached.Namespace}, found)
	if err != nil && apierrors.IsNotFound(err) {
		// Define a new deployment
		dep := r.deploymentForMemcached()
		// Create the Deployment on the cluster
		if err = r.Create(ctx, dep); err != nil {
            log.Error(err, "Failed to create new Deployment",
            "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
            return ctrl.Result{}, err
        }
		...
	}
```

接着需要注意，`deploymentForMemcached()` 函数需要定义并返回应在集群上创建的 Deployment。该函数应构造具备必要规格的 Deployment 对象，如下例所示：

```go
    dep := &appsv1.Deployment{
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image:           "memcached:1.6.26-alpine3.19",
						Name:            "memcached",
						ImagePullPolicy: corev1.PullIfNotPresent,
						Ports: []corev1.ContainerPort{{
							ContainerPort: 11211,
							Name:          "memcached",
						}},
						Command: []string{"memcached", "--memory-limit=64", "-o", "modern", "-v"},
					}},
				},
			},
		},
	}
```

此外，我们需要实现一个机制，以校验集群中的 Memcached 副本数是否与 CR 中指定的期望值一致。如果不一致，调谐过程必须更新集群以确保一致性。这意味着：无论何时在集群上创建或更新 Memcached Kind 的 CR，控制器都会持续调谐，直到实际副本数与期望值一致。如下例所示：

```go
	...
	size := memcached.Spec.Size
	if *found.Spec.Replicas != size {
		found.Spec.Replicas = &size
		if err = r.Update(ctx, found); err != nil {
			log.Error(err, "Failed to update Deployment",
				"Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
            return ctrl.Result{}, err
        }
    ...
```

现在，你可以查看负责管理 Memcached Kind 自定义资源的完整控制器。该控制器确保集群中的期望状态得以维持，从而保证 Memcached 实例始终以用户指定的副本数运行。

<details><summary><code>internal/controller/memcached_controller.go</code>: Our Controller Implementation </summary>

```go
{{#include ./getting-started/testdata/project/internal/controller/memcached_controller.go}}
```
</details>

### 深入控制器实现

#### 配置 Manager 监听资源

核心思想是监听对控制器重要的资源。当控制器关注的资源发生变化时，Watch 会触发控制器的调谐循环，以确保资源的实际状态与控制器逻辑定义的期望状态相匹配。

注意我们如何配置 Manager 来监控 Memcached Kind 的自定义资源（CR）的创建、更新或删除等事件，以及控制器所管理并拥有的 Deployment 的任何变化：

```go
// SetupWithManager sets up the controller with the Manager.
// The Deployment is also watched to ensure its
// desired state in the cluster.
func (r *MemcachedReconciler) SetupWithManager(mgr ctrl.Manager) error {
    return ctrl.NewControllerManagedBy(mgr).
		// Watch the Memcached Custom Resource and trigger reconciliation whenever it
		//is created, updated, or deleted
		For(&cachev1alpha1.Memcached{}).
		// Watch the Deployment managed by the Memcached controller. If any changes occur to the Deployment
        // owned and managed by this controller, it will trigger reconciliation, ensuring that the cluster
        // state aligns with the desired state.
		Owns(&appsv1.Deployment{}).
		Complete(r)
    }
```

#### 但是，Manager 如何知道哪些资源归它所有？

我们并不希望控制器去监听集群中的所有 Deployment 并触发调谐循环；我们只希望在运行 Memcached 实例的那个特定 Deployment 发生变化时才触发。例如，如果有人误删了我们的 Deployment 或修改了其副本数，我们希望触发调谐以使其回到期望状态。

Manager 之所以知道应该观察哪个 Deployment，是因为我们设置了 `ownerRef`（Owner Reference）：

```go
if err := ctrl.SetControllerReference(memcached, dep, r.Scheme); err != nil {
    return nil, err
}
```

<aside class="note">

<h1><code>ownerRef</code> 与级联事件</h1>

ownerRef 不仅能让我们观察特定资源的变化，还很关键的一点是：当我们从集群中删除 Memcached 自定义资源（CR）时，我们希望其所拥有的所有资源也能在级联事件中被自动删除。

这能确保当父资源（Memcached CR）被移除时，所有关联资源（如 Deployment、Service 等）也会被清理，从而保持集群整洁一致。

更多信息请参见 Kubernetes 文档：[Owners and Dependents](https://kubernetes.io/docs/concepts/overview/working-with-objects/owners-dependents/)。

</aside>

### 授予权限

确保控制器拥有管理其资源所需的权限（例如创建、获取、更新、列出）非常重要。

[RBAC 权限][k8s-rbac] 现在通过 [RBAC 标记][rbac-markers] 配置，这些标记用于生成并更新 `config/rbac/` 中的清单文件。它们可以（且应当）定义在每个控制器的 `Reconcile()` 方法上，如下示例所示：

```go
// +kubebuilder:rbac:groups=cache.example.com,resources=memcacheds,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cache.example.com,resources=memcacheds/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=cache.example.com,resources=memcacheds/finalizers,verbs=update
// +kubebuilder:rbac:groups=core,resources=events,verbs=create;patch
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch
```

修改控制器后，运行 `make manifests` 命令。这将促使 [controller-gen][controller-gen] 刷新 `config/rbac` 下的文件。

<details><summary><code>config/rbac/role.yaml</code>: Our RBAC Role generated </summary>

```yaml
{{#include ./getting-started/testdata/project/config/rbac/role.yaml}}
```
</details>

### Manager（main.go）

`cmd/main.go` 中的 [Manager][manager] 负责管理应用中的各个控制器。

<details><summary><code>cmd/main.go</code>: Our main.go </summary>

```go
{{#include ./getting-started/testdata/project/cmd/main.go}}
```
</details>

### 使用 Kubebuilder 插件生成额外选项

现在你已经更好地理解了如何创建自己的 API 与控制器，让我们在该项目中引入 [`autoupdate.kubebuilder.io/v1-alpha`][autoupdate-plugin] 插件，以便你的项目能跟随最新的 Kubebuilder 版本脚手架变化保持更新，并由此采纳生态中的改进。

```shell
kubebuilder edit --plugins="autoupdate/v1-alpha"
```

查看 `.github/workflows/auto-update.yml` 文件了解其工作方式。

### 在集群中验证项目

此时你可以参考快速开始中定义的步骤在集群中验证该项目，见：[Run It On the Cluster](./quick-start#run-it-on-the-cluster)

## 下一步

- 若想更深入地开发你的方案，考虑阅读 [CronJob 教程][cronjob-tutorial]
- 有关优化方法的思路，请参考[最佳实践][best-practices]

<aside class="note">
<h1> 使用 Deploy Image 插件生成 API 与源码 </h1>

既然已经更为熟悉，你可能想看看 [Deploy Image][deploy-image] 插件。该插件允许用户为在集群上部署和管理 Operand（镜像）搭建 API/Controller 脚手架。它会提供与本指南类似的脚手架，并附带额外特性，例如为你的控制器实现的测试。

</aside>

[k8s-operator-pattern]: https://kubernetes.io/docs/concepts/extend-kubernetes/operator/
[controller-runtime]: https://github.com/kubernetes-sigs/controller-runtime
[group-kind-oh-my]: ./cronjob-tutorial/gvks.md
[controller-gen]: ./reference/controller-gen.md
[markers]: ./reference/markers.md
[rbac-markers]: ./reference/markers/rbac.md
[k8s-rbac]: https://kubernetes.io/docs/reference/access-authn-authz/rbac/
[manager]: https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/manager
[options-manager]: https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/manager#Options
[quick-start]: ./quick-start.md
[best-practices]: ./reference/good-practices.md
[cronjob-tutorial]: https://book.kubebuilder.io/cronjob-tutorial/cronjob-tutorial.html
[deploy-image]: ./plugins/available/deploy-image-plugin-v1-alpha.md
[GOPATH-golang-docs]: https://golang.org/doc/code.html#GOPATH
[go-modules-blogpost]: https://blog.golang.org/using-go-modules
[autoupdate-plugin]: ./plugins/available/autoupdate-v1-alpha.md
