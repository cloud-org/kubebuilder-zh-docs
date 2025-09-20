# 监听由控制器“拥有”的二级资源

在 Kubernetes 控制器中，通常同时管理 **主资源（Primary Resources）** 与 **二级资源（Secondary Resources）**。主资源是控制器直接负责的对象；二级资源则由控制器创建与管理，用于支撑主资源的运行。

本节介绍如何管理被控制器 `Owned` 的二级资源。示例涵盖：

- 在主资源（`Busybox`）与二级资源（`Deployment`）之间设置 [OwnerReference][cr-owner-ref-doc]，以确保生命周期正确关联；
- 在 `SetupWithManager()` 中通过 `Owns()` 让控制器监听该二级资源。由于 `Deployment` 由 `Busybox` 控制器创建并管理，因此属于其 Owned 资源。

## 设置 OwnerReference

要将二级资源（`Deployment`）的生命周期与主资源（`Busybox`）关联，需要在二级资源上设置 [OwnerReference][cr-owner-ref-doc]。这样，当主资源被删除时，Kubernetes 会级联删除二级资源。

controller-runtime 提供了 [controllerutil.SetControllerReference][cr-owner-ref-doc] 来设置该关系。

### 设置 OwnerReference 示例

Below, we create the `Deployment` and set the Owner reference between the `Busybox` custom resource and the `Deployment` using `controllerutil.SetControllerReference()`.

```go
// deploymentForBusybox returns a Deployment object for Busybox
func (r *BusyboxReconciler) deploymentForBusybox(busybox *examplecomv1alpha1.Busybox) *appsv1.Deployment {
    replicas := busybox.Spec.Size

    dep := &appsv1.Deployment{
        ObjectMeta: metav1.ObjectMeta{
            Name:      busybox.Name,
            Namespace: busybox.Namespace,
        },
        Spec: appsv1.DeploymentSpec{
            Replicas: &replicas,
            Selector: &metav1.LabelSelector{
                MatchLabels: map[string]string{"app": busybox.Name},
            },
            Template: metav1.PodTemplateSpec{
                ObjectMeta: metav1.ObjectMeta{
                    Labels: map[string]string{"app": busybox.Name},
                },
                Spec: corev1.PodSpec{
                    Containers: []corev1.Container{
                        {
                            Name:  "busybox",
                            Image: "busybox:latest",
                        },
                    },
                },
            },
        },
    }

    // 为 Deployment 设置 ownerRef，保证 Busybox 被删除时它也会被删除
    controllerutil.SetControllerReference(busybox, dep, r.Scheme)
    return dep
}
```

### 说明

设置 `OwnerReference` 后，当 `Busybox` 被删除时，`Deployment` 也会被自动清理。控制器也可据此监听 `Deployment` 的变化，确保副本数等期望状态得以维持。

例如，若有人将 `Deployment` 的副本数改为 3，而 `Busybox` CR 期望为 1，控制器会在调谐中将其缩回到 1。

**Reconcile 函数示例**

```go
// Reconcile handles the main reconciliation loop for Busybox and the Deployment
func (r *BusyboxReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
    log := logf.FromContext(ctx)

    // Fetch the Busybox instance
    busybox := &examplecomv1alpha1.Busybox{}
    if err := r.Get(ctx, req.NamespacedName, busybox); err != nil {
        if apierrors.IsNotFound(err) {
            log.Info("Busybox resource not found. Ignoring since it must be deleted")
            return ctrl.Result{}, nil
        }
        log.Error(err, "Failed to get Busybox")
        return ctrl.Result{}, err
    }

    // Check if the Deployment already exists, if not create a new one
    found := &appsv1.Deployment{}
    err := r.Get(ctx, types.NamespacedName{Name: busybox.Name, Namespace: busybox.Namespace}, found)
    if err != nil && apierrors.IsNotFound(err) {
        // Define a new Deployment
        dep := r.deploymentForBusybox(busybox)
        log.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
        if err := r.Create(ctx, dep); err != nil {
            log.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
            return ctrl.Result{}, err
        }
        // Requeue the request to ensure the Deployment is created
        return ctrl.Result{RequeueAfter: time.Minute}, nil
    } else if err != nil {
        log.Error(err, "Failed to get Deployment")
        return ctrl.Result{}, err
    }

    // Ensure the Deployment size matches the desired state
    size := busybox.Spec.Size
    if *found.Spec.Replicas != size {
        found.Spec.Replicas = &size
        if err := r.Update(ctx, found); err != nil {
            log.Error(err, "Failed to update Deployment size", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
            return ctrl.Result{}, err
        }
        // Requeue the request to ensure the correct state is achieved
        return ctrl.Result{Requeue: true}, nil
    }

    // Update Busybox status to reflect that the Deployment is available
    busybox.Status.AvailableReplicas = found.Status.AvailableReplicas
    if err := r.Status().Update(ctx, busybox); err != nil {
        log.Error(err, "Failed to update Busybox status")
        return ctrl.Result{}, err
    }

    return ctrl.Result{}, nil
}
```

## Watching Secondary Resources

To ensure that changes to the secondary resource (such as the `Deployment`) trigger
a reconciliation of the primary resource (`Busybox`), we configure the controller
to watch both resources.

The `Owns()` method allows you to specify secondary resources
that the controller should monitor. This way, the controller will
automatically reconcile the primary resource whenever the secondary
resource changes (e.g., is updated or deleted).

### Example: Configuring `SetupWithManager` to Watch Secondary Resources

```go
// SetupWithManager sets up the controller with the Manager.
// The controller will watch both the Busybox primary resource and the Deployment secondary resource.
func (r *BusyboxReconciler) SetupWithManager(mgr ctrl.Manager) error {
    return ctrl.NewControllerManagedBy(mgr).
        For(&examplecomv1alpha1.Busybox{}).  // Watch the primary resource
        Owns(&appsv1.Deployment{}).          // Watch the secondary resource (Deployment)
        Complete(r)
}
```

## Ensuring the Right Permissions

Kubebuilder uses [markers][markers] to define RBAC permissions
required by the controller. In order for the controller to
properly watch and manage both the primary (`Busybox`) and secondary (`Deployment`)
resources, it must have the appropriate permissions granted;
i.e. to `watch`, `get`, `list`, `create`, `update`, and `delete` permissions for those resources.

### Example: RBAC Markers

Before the `Reconcile` method, we need to define the appropriate RBAC markers.
These markers will be used by [controller-gen][controller-gen] to generate the necessary
roles and permissions when you run `make manifests`.

```go
// +kubebuilder:rbac:groups=example.com,resources=busyboxes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
```

- The first marker gives the controller permission to manage the `Busybox` custom resource (the primary resource).
- The second marker grants the controller permission to manage `Deployment` resources (the secondary resource).

Note that we are granting permissions to `watch` the resources.

[owner-ref-k8s-docs]: https://kubernetes.io/docs/concepts/overview/working-with-objects/owners-dependents/
[cr-owner-ref-doc]: https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/controller/controllerutil#SetOwnerReference
[controller-gen]: ./../controller-gen.md
[markers]:./../markers/rbac.md
