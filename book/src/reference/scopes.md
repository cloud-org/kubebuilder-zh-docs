# 管理器（Operator）与 CRD 的作用域（Scopes）

本节介绍 Kubebuilder 项目中运行与资源层面的作用域配置。Kubernetes 中的 Manager（“Operator”）可以限定在某个命名空间或整个集群范围内，从而影响其对资源的监听与管理方式。

同时，CRD 也可定义为命名空间级或集群级，这会影响其在集群中的可见范围。

## 配置 Manager 的作用域

可根据所需管理的资源，选择不同的作用域：

### （默认）监听全部命名空间

默认情况下，若未指定命名空间，manager 将监听所有命名空间：

```go
mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
...
})
```

### 监听单个命名空间

如需限定到单个命名空间，可设置相应的 Cache 配置：

```go
mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
...
   Cache: cache.Options{
      DefaultNamespaces: map[string]cache.Config{"operator-namespace": cache.Config{}},
   },
})
```

### 监听多个命名空间

也可通过 [Cache Config][CacheConfig] 指定多个命名空间：

```go
mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
...
Cache: cache.Options{
    DefaultNamespaces: map[string]cache.Config{
        "operator-namespace1": cache.Config{},
        "operator-namespace2": cache.Config{},
        },
    },
})
```

## 配置 CRD 的作用域

CRD 的作用域决定其仅在部分命名空间可见，还是在整个集群可见。

### 命名空间级（Namespace-scoped）CRD

当需要将资源隔离到特定命名空间时，可选择命名空间级 CRD，有助于按团队或应用进行划分。
但需注意：由于 CRD 的特殊性，验证新版本并不直接。需要设计合理的版本与转换策略（参见 [Kubebuilder 多版本教程][kubebuilder-multiversion-tutorial]），并协调由哪一个 manager 实例负责转换（参见 [Kubernetes 官方文档][k8s-crd-conversion]）。
此外，为确保在预期范围内生效，Mutating/Validating Webhook 的配置也应考虑命名空间作用域，从而支持更可控、分阶段的发布。

### 集群级（Cluster-scoped）CRD

对于需要在整个集群访问与管理的资源（例如共享配置或全局资源），应选择集群级 CRD。

#### 配置 CRD 的作用域

**在创建 API 时**

CRD 的作用域会在生成清单时确定。Kubebuilder 的 API 创建命令支持该配置。

默认情况下，生成的 API 对应 CRD 为命名空间级；若需集群级，请使用 `--namespaced=false`，例如：

```shell
kubebuilder create api --group cache --version v1alpha1 --kind Memcached --resource=true --controller=true --namespaced=false
```

上述命令会生成集群级 CRD，意味着它可在所有命名空间访问与管理。

**更新已有 API**

在创建 API 之后仍可调整作用域。若想将 CRD 配置为集群级，可在 Go 类型定义上方添加 `+kubebuilder:resource:scope=Cluster` 标记。例如：

```go
//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster,shortName=mc

...
```

设置标记后，运行 `make manifests` 以生成文件。该命令会调用 [`controller-gen`][controller-tools]，依据 Go 文件中的标记生成 CRD 清单。

生成的清单会正确体现作用域（Cluster 或 Namespaced），无需手动修改 YAML。

[controller-tools]: https://sigs.k8s.io/controller-tools
[CacheConfig]: https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/cache#Config
[kubebuilder-multiversion-tutorial]: https://book.kubebuilder.io/multiversion-tutorial/tutorial
[k8s-crd-conversion]: https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definition-versioning/#webhook-conversion
