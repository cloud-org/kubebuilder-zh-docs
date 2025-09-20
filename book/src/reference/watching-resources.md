# 监听资源（Watching Resources）

在扩展 Kubernetes API 时，我们希望方案的行为与 Kubernetes 本身保持一致。以 `Deployment` 为例，其由一个控制器管理：当集群中发生创建、更新、删除等事件时，控制器触发调谐以使资源状态与期望一致。

类似地，开发控制器时，我们需要监听与方案相关的资源变化；无论是创建、更新还是删除，都应触发调谐循环以采取相应动作并保持一致性。

[controller-runtime][controller-runtime] 提供了多种监听与管理资源的方式。

## 主资源（Primary Resources）

**主资源** 是控制器直接负责管理的资源。例如，为 `MyApp` 创建了 CRD，则相应控制器负责管理 `MyApp` 实例。

在这种情况下，`MyApp` 是该控制器的主资源，调谐循环的目标就是维持这些主资源的期望状态。

使用 Kubebuilder 创建新 API 时，会脚手架如下默认代码，保证控制器通过 `For()` 监听该 API 的创建、更新与删除事件：

该设置确保当 API 实例被创建、更新或删除时，都会触发调谐：

```go
// Watches the primary resource (e.g., MyApp) for create, update, delete events
if err := ctrl.NewControllerManagedBy(mgr).
   For(&<YourAPISpec>{}). <-- See there that the Controller is For this API
   Complete(r); err != nil {
   return err
}
```

## 二级资源（Secondary Resources）

控制器通常还需管理 **二级资源**，即为支撑主资源在集群中运行所需的各类资源。

二级资源的变化会直接影响主资源，因此控制器需相应地监听并调谐它们。

### 由控制器“拥有”的二级资源

当二级资源（如 `Service`、`ConfigMap`、`Deployment`）被控制器 `Owned` 时，意味着它们由该控制器创建并通过 [OwnerReferences][owner-ref-k8s-docs] 与主资源关联。

For example, if we have a controller to manage our CR(s) of the Kind `MyApp`
on the cluster, which represents our application solution, all resources required
to ensure that `MyApp` is up and running with the desired number of instances
will be **Secondary Resources**. The code responsible for creating, deleting,
and updating these resources will be part of the `MyApp` Controller.
We would add the appropriate [OwnerReferences][owner-ref-k8s-docs]
using the [controllerutil.SetControllerReference][cr-owner-ref-doc]
function to indicate that these resources are owned by the same controller
responsible for managing `MyApp` instances, which will be reconciled by the `MyAppReconciler`.

此外，当主资源被删除时，Kubernetes 的垃圾回收会级联删除关联的二级资源。

### 非本控制器“拥有”的二级资源

二级资源既可能来自本项目，也可能来自其他项目，与主资源相关，但并非由本控制器创建或管理。

For example, if we have a CRD that represents a backup solution (i.e. `MyBackup`) for our `MyApp`,
it might need to watch changes in the `MyApp` resource to trigger reconciliation in `MyBackup`
to ensure the desired state. Similarly, `MyApp`'s behavior might also be impacted by
CRDs/APIs defined in other projects.

在这两种情况下，即便它们不是 `MyAppController` 的 `Owned` 资源，仍被视为二级资源。

In Kubebuilder, resources that are not defined in the project itself and are not
a **Core Type** (those not defined in the Kubernetes API) are called **External Types**.

An **External Type** refers to a resource that is not defined in your
project but one that you need to watch and respond to.
For example, if **Operator A** manages a `MyApp` CRD for application deployment,
and **Operator B** handles backups, **Operator B** can watch the `MyApp` CRD as an external type
to trigger backup operations based on changes in `MyApp`.

In this scenario, **Operator B** could define a `BackupConfig` CRD that relies on the state of `MyApp`.
By treating `MyApp` as a **Secondary Resource**, **Operator B** can watch and reconcile changes in **Operator A**'s `MyApp`,
ensuring that backup processes are initiated whenever `MyApp` is updated or scaled.

## 监听资源的一般思路

Whether a resource is defined within your project or comes from an external project, the concept of **Primary**
and **Secondary Resources** remains the same:
- The **Primary Resource** is the resource the controller is primarily responsible for managing.
- **Secondary Resources** are those that are required to ensure the primary resource works as desired.

Therefore, regardless of whether the resource was defined by your project or by another project,
your controller can watch, reconcile, and manage changes to these resources as needed.

## Why does watching the secondary resources matter?

When building a Kubernetes controller, it’s crucial to not only focus
on **Primary Resources** but also to monitor **Secondary Resources**.
Failing to track these resources can lead to inconsistencies in your
controller's behavior and the overall cluster state.

Secondary resources may not be directly managed by your controller,
but changes to these resources can still significantly
impact the primary resource and your controller's functionality.
Here are the key reasons why it's important to watch them:

- **Ensuring Consistency**:
    - Secondary resources (e.g., child objects or external dependencies) may diverge from their desired state.
    For instance, a secondary resource may be modified or deleted, causing the system to fall out of sync.
    - Watching secondary resources ensures that any changes are detected immediately, allowing the controller to
    reconcile and restore the desired state.

- **Avoiding Random Self-Healing**:
    - Without watching secondary resources, the controller may "heal" itself only upon restart or when specific events
    are triggered. This can cause unpredictable or delayed reactions to issues.
    - Monitoring secondary resources ensures that inconsistencies are addressed promptly, rather than waiting for a
    controller restart or external event to trigger reconciliation.

- **Effective Lifecycle Management**:
    - Secondary resources might not be owned by the controller directly, but their state still impacts the behavior
    of primary resources. Without watching these, you risk leaving orphaned or outdated resources.
    - Watching non-owned secondary resources lets the controller respond to lifecycle events (create, update, delete)
    that might affect the primary resource, ensuring consistent behavior across the system.

示例见：[监听非 Owned 的二级资源](./watching-resources/secondary-resources-not-owned.md#configuration-example)。

## 为何不直接用 `RequeueAfter X` 代替监听？

Kubernetes 控制器本质上是**事件驱动**的：调谐循环通常由资源的创建、更新、删除等事件触发。相较于固定周期的 `RequeueAfter` 轮询，事件驱动更高效、更及时，能在需要时才行动，兼顾性能与效率。

In many cases, **watching resources** is the preferred approach for ensuring Kubernetes resources
remain in the desired state. It is more efficient, responsive, and aligns with Kubernetes' event-driven architecture.
However, there are scenarios where `RequeueAfter` is appropriate and necessary, particularly for managing external
systems that do not emit events or for handling resources that take time to converge, such as long-running processes.
Relying solely on `RequeueAfter` for all scenarios can lead to unnecessary overhead and
delayed reactions. Therefore, it is essential to prioritize **event-driven reconciliation** by configuring
your controller to **watch resources** whenever possible, and reserving `RequeueAfter` for situations
where periodic checks are required.

### 何时应使用 `RequeueAfter X`

While `RequeueAfter` is not the primary method for triggering reconciliations, there are specific cases where it is
necessary, such as:

- 观察无事件外部系统：例如外部数据库、三方服务等不产生活动事件的对象，可用 `RequeueAfter` 周期性检查。
- 基于时间的操作：如轮换密钥、证书续期等需按固定间隔进行的任务。
- 处理错误/延迟：当资源需要时间自愈时，`RequeueAfter` 可避免持续触发调谐，改为延时再查。

## 使用 Predicates

在更复杂的场景中，可使用 [Predicates][cr-predicates] 精细化触发条件：按特定字段、标签或注解的变化过滤事件，使控制器只对相关事件响应并保持高效。

[controller-runtime]: https://github.com/kubernetes-sigs/controller-runtime
[owner-ref-k8s-docs]: https://kubernetes.io/docs/concepts/overview/working-with-objects/owners-dependents/
[cr-predicates]: https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/predicate
[secondary-resources-doc]: watching-resources/secondary-owned-resources
[predicates-with-external-type-doc]: watching-resources/predicates-with-watch
[cr-owner-ref-doc]: https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/controller/controllerutil#SetOwnerReference
