# 使用 Finalizer（Using Finalizers）

`Finalizer` 允许控制器实现异步的“删除前”钩子。举例来说，如果你的每个自定义对象在外部系统中都对应着某个资源（例如对象存储的桶），当该对象在 Kubernetes 中被删除时，你希望同步清理外部资源，此时即可借助 Finalizer 实现。

关于 Finalizer 的更多背景请参阅 [Kubernetes 参考文档](https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definitions/#finalizers)。下文展示如何在控制器的 `Reconcile` 方法中注册并触发删除前的处理逻辑。

关键点：Finalizer 会让对对象的“删除”变为一次“更新”，即为对象打上删除时间戳。对象上存在删除时间戳意味着其处于删除流程中。否则（没有 Finalizer 时），删除表现为一次调谐（reconcile）里对象已从缓存中缺失的情况。

要点摘录：
- 当对象未被删除且尚未注册 Finalizer 时，需添加 Finalizer 并更新该对象。
- 当对象进入删除流程且 Finalizer 仍存在时，执行删除前逻辑，随后移除 Finalizer 并更新对象。
- 删除前逻辑应具备幂等性。

{{#literatego ../cronjob-tutorial/testdata/finalizer_example.go}}
