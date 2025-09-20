# 使用 Predicates 精细化 Watch

在编写控制器时，使用 **Predicates** 来过滤事件、控制何时触发调谐往往很有帮助。

[Predicates][predicates-doc] 允许基于事件（创建、更新、删除）和资源字段（标签、注解、状态等）定义触发条件。借助 **[Predicates][predicates-doc]**，可以让控制器仅对关心的变化做出响应。

当需要精确限定哪些变化应触发调谐时，Predicates 尤其有用：它能避免无谓的调谐，让控制器只对真正相关的变更做出反应。

## 何时使用 Predicates

**适用场景：**

- 忽略不相关的变更，例如不影响业务字段的更新；
- 仅对带特定标签/注解的资源触发调谐；
- 监听外部资源时仅对特定变化作出反应。

## 示例：使用 Predicates 过滤更新事件

设想我们只在 **`Busybox`** 的特定字段变化（例如 `spec.size`）时触发 **`BackupBusybox`** 控制器调谐，忽略其它变化（如 status 更新）。

### 定义 Predicate

如下定义仅在 **`Busybox`** 发生“有意义”的更新时允许调谐：

```go
import (
    "sigs.k8s.io/controller-runtime/pkg/predicate"
    "sigs.k8s.io/controller-runtime/pkg/event"
)

// 仅在 Busybox 的 spec.size 变化时触发调谐
updatePred := predicate.Funcs{
    // 仅当 spec.size 发生变化时允许更新事件通过
    UpdateFunc: func(e event.UpdateEvent) bool {
        oldObj := e.ObjectOld.(*examplecomv1alpha1.Busybox)
        newObj := e.ObjectNew.(*examplecomv1alpha1.Busybox)

    // 仅当 spec.size 字段变化时返回 true
        return oldObj.Spec.Size != newObj.Spec.Size
    },

    // 放行创建事件
    CreateFunc: func(e event.CreateEvent) bool {
        return true
    },

    // 放行删除事件
    DeleteFunc: func(e event.DeleteEvent) bool {
        return true
    },

    // 放行通用事件（如外部触发）
    GenericFunc: func(e event.GenericEvent) bool {
        return true
    },
}
```

### 说明

在本例中：
- 仅当 **`spec.size`** 发生变化时 **`UpdateFunc`** 才返回 `true`，其余 `spec` 变更（注解等）会被忽略；
- **`CreateFunc`**、**`DeleteFunc`**、**`GenericFunc`** 返回 `true`，意味着这三类事件依旧会触发调谐。

这样可确保控制器仅在 **`spec.size`** 被修改时进行调谐，忽略与业务无关的其它变更。

### 在 `Watches` 中使用 Predicates

Now, we apply this predicate in the **`Watches()`** method of
the **`BackupBusyboxReconciler`** to trigger reconciliation only for relevant events:

```go
// SetupWithManager 配置控制器。控制器会监听主资源 BackupBusybox 与 Busybox，并应用 predicates。
func (r *BackupBusyboxReconciler) SetupWithManager(mgr ctrl.Manager) error {
    return ctrl.NewControllerManagedBy(mgr).
        For(&examplecomv1alpha1.BackupBusybox{}).  // 监听主资源（BackupBusybox）
        Watches(
            &examplecomv1alpha1.Busybox{},  // 监听 Busybox CR
            handler.EnqueueRequestsFromMapFunc(func(ctx context.Context, obj client.Object) []reconcile.Request {
                return []reconcile.Request{
                    {
                        NamespacedName: types.NamespacedName{
                            Name:      "backupbusybox",    // 对应的 BackupBusybox 资源
                            Namespace: obj.GetNamespace(),  // 使用 Busybox 的命名空间
                        },
                    },
                }
            }),
            builder.WithPredicates(updatePred),  // 应用 Predicate
        ).  // 当 Busybox 变化且满足条件时触发调谐
        Complete(r)
}
```

### 说明

- **[`builder.WithPredicates(updatePred)`][predicates-doc]**：应用谓词，确保仅当 **`Busybox`** 的 **`spec.size`** 变化时才触发调谐。
- **其他事件**：控制器仍会响应 `Create`、`Delete` 与 `Generic` 事件。

[predicates-doc]: https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/source#WithPredicates
