# 创建事件（Events）

在控制器的 Reconcile 函数中发布 *Event* 通常很有用：它让用户或自动化流程能够了解某个对象上发生了什么，并据此作出反应。

可通过 `$ kubectl describe <资源类型> <资源名>` 查看某对象的近期事件，或通过 `$ kubectl get events` 查看全局事件列表。

<aside class="warning">
<h1>仅在必要场景下触发事件</h1>

请勿为所有操作都触发事件。事件过多会带来糟糕的体验，使集群使用者难以从噪音中筛选出可执行的信息。参见 [Kubernetes API 约定][Events]。

</aside>

## 编写事件（Writing Events）

事件的函数原型：

```go
Event(object runtime.Object, eventtype, reason, message string)
```

- `object`：事件关联的对象。
- `eventtype`：事件类型，为 *Normal* 或 *Warning*（[更多][Event-Example]）。
- `reason`：事件原因。建议短小唯一、采用 `UpperCamelCase`，便于自动化流程在 switch 中处理（[更多][Reason-Example]）。
- `message`：展示给人看的详细描述（[更多][Message-Example]）。



<aside class="note">
<h1>示例</h1>

下面示例展示如何触发一个事件：

```go
	// The following implementation will raise an event
	r.Recorder.Eventf(cr, "Warning", "Deleting",
		"Custom Resource %s is being deleted from the namespace %s",
		cr.Name, cr.Namespace)
```

</aside>

### 如何在控制器中触发事件？

在控制器的调谐流程中，你可以使用 [EventRecorder][Events] 发布事件。通过在 Manager 上调用 `GetRecorder(name string)` 可创建对应的 recorder。下面演示如何修改 `cmd/main.go`：

```go
	if err := (&controller.MyKindReconciler{
		Client:   mgr.GetClient(),
		Scheme:   mgr.GetScheme(),
		// Note that we added the following line:
		Recorder: mgr.GetEventRecorderFor("mykind-controller"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "MyKind")
		os.Exit(1)
	}
```

### 在控制器中接入 EventRecorder

为触发事件，控制器需要持有 `record.EventRecorder`：
```go
import (
	...
	"k8s.io/client-go/tools/record"
	...
)
// MyKindReconciler reconciles a MyKind object
type MyKindReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	// See that we added the following code to allow us to pass the record.EventRecorder
	Recorder record.EventRecorder
}
```
### 将 EventRecorder 传入控制器

仍以 `cmd/main.go` 为例，向控制器构造体传入 recorder：

```go
	if err := (&controller.MyKindReconciler{
		Client:   mgr.GetClient(),
		Scheme:   mgr.GetScheme(),
		// Note that we added the following line:
		Recorder: mgr.GetEventRecorderFor("mykind-controller"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "MyKind")
		os.Exit(1)
	}
```

### 授权所需权限（RBAC）

还需为项目授予创建事件的权限。在控制器上添加如下 [RBAC][rbac-markers] 标记：

```go
...
// +kubebuilder:rbac:groups=core,resources=events,verbs=create;patch
...
func (r *MyKindReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
```

然后执行 `$ make manifests` 更新 `config/rbac/role.yaml` 中的规则。

[Events]: https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#events
[Event-Example]: https://github.com/kubernetes/api/blob/6c11c9e4685cc62e4ddc8d4aaa824c46150c9148/core/v1/types.go#L6019-L6024
[Reason-Example]: https://github.com/kubernetes/api/blob/6c11c9e4685cc62e4ddc8d4aaa824c46150c9148/core/v1/types.go#L6048
[Message-Example]: https://github.com/kubernetes/api/blob/6c11c9e4685cc62e4ddc8d4aaa824c46150c9148/core/v1/types.go#L6053
[rbac-markers]: ./markers/rbac.md
