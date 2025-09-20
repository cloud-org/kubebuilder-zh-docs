# 监听非本控制器“拥有”的二级资源

在某些场景下，控制器需要监听并响应那些并非由自身创建或管理（即非 `Owned`）的资源的变化——这些资源通常由其他控制器创建与维护。

以下示例展示控制器如何监测并调谐其未直接管理的资源。这适用于任何非 `Owned` 的资源，包括由其他控制器或项目管理、在独立进程中调谐的 **核心类型（Core Types）** 或 **自定义资源（CR）**。

例如，有两个自定义资源 `Busybox` 与 `BackupBusybox`。如果希望 `Busybox` 的变化触发 `BackupBusybox` 控制器的调谐，则可让 `BackupBusybox` 控制器去监听 `Busybox` 的变化。

### 示例：监听非 Owned 的 Busybox 以调谐 BackupBusybox

假设某控制器负责管理 `BackupBusybox`，但也需要关注集群中的 `Busybox` 变化。我们只希望当 `Busybox` 启用了备份能力时，才触发调谐。

- **为何要监听二级资源？**
    - `BackupBusybox` 控制器不创建/不拥有 `Busybox`，但后者的更新与删除会直接影响其主资源（`BackupBusybox`）。
    - 通过只监听具有特定标签的 `Busybox` 实例，可确保仅对相关对象执行必要动作（如备份）。

### 配置示例

如下配置使 `BackupBusyboxReconciler` 监听 `Busybox` 的变化，并触发对 `BackupBusybox` 的调谐：

```go
// SetupWithManager sets up the controller with the Manager.
// The controller will watch both the BackupBusybox primary resource and the Busybox resource.
func (r *BackupBusyboxReconciler) SetupWithManager(mgr ctrl.Manager) error {
    return ctrl.NewControllerManagedBy(mgr).
        For(&examplecomv1alpha1.BackupBusybox{}).  // Watch the primary resource (BackupBusybox)
        Watches(
            &examplecomv1alpha1.Busybox{},  // Watch the Busybox CR
            handler.EnqueueRequestsFromMapFunc(func(ctx context.Context, obj client.Object) []reconcile.Request {
                // Trigger reconciliation for the BackupBusybox in the same namespace
                return []reconcile.Request{
                    {
                        NamespacedName: types.NamespacedName{
                            Name:      "backupbusybox",  // Reconcile the associated BackupBusybox resource
                            Namespace: obj.GetNamespace(),  // Use the namespace of the changed Busybox
                        },
                    },
                }
            }),
        ).  // Trigger reconciliation when the Busybox resource changes
        Complete(r)
}
```

进一步，我们可以只针对带有特定标签的 `Busybox` 触发调谐：

```go
// SetupWithManager sets up the controller with the Manager.
// The controller will watch both the BackupBusybox primary resource and the Busybox resource, filtering by a label.
func (r *BackupBusyboxReconciler) SetupWithManager(mgr ctrl.Manager) error {
    return ctrl.NewControllerManagedBy(mgr).
        For(&examplecomv1alpha1.BackupBusybox{}).  // Watch the primary resource (BackupBusybox)
        Watches(
            &examplecomv1alpha1.Busybox{},  // Watch the Busybox CR
            handler.EnqueueRequestsFromMapFunc(func(ctx context.Context, obj client.Object) []reconcile.Request {
                // 检查 Busybox 是否带有 'backup-enable: "true"' 标签
                if val, ok := obj.GetLabels()["backup-enable"]; ok && val == "true" {
                    // 若命中该标签，则触发 BackupBusybox 的调谐
                    return []reconcile.Request{
                        {
                            NamespacedName: types.NamespacedName{
                                Name:      "backupbusybox",  // Reconcile the associated BackupBusybox resource
                                Namespace: obj.GetNamespace(),  // Use the namespace of the changed Busybox
                            },
                        },
                    }
                }
                // 未命中标签时不触发
                return []reconcile.Request{}
            }),
        ).  // Trigger reconciliation when the labeled Busybox resource changes
        Complete(r)
}
```
