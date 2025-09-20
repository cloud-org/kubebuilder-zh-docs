# 用于配置/代码生成的标记（Markers）

Kubebuilder 使用
[controller-gen](/reference/controller-gen.md)
来生成实用代码与 Kubernetes YAML。生成行为由 Go 代码中的特殊“标记注释”控制。

“标记注释”是以加号开头的单行注释，后跟标记名称，并可选带有该标记的配置：

```go
// +kubebuilder:validation:Optional
// +kubebuilder:validation:MaxItems=2
// +kubebuilder:printcolumn:JSONPath=".status.replicas",name=Replicas,type=string
```

<aside class="note">
<h1><code>// +optional</code> 与 <code>// +kubebuilder:validation:Optional</code> 的区别</h1>

controller-gen 同时支持二者（可参考 `controller-gen crd -www` 的输出）。这两个标记都可以加在字段上。

不过，`+kubebuilder:validation:Optional` 也可以加在包级别，使其作用于包内所有字段。

若只使用 controller-gen，它们有些冗余；但如果你还使用其他生成器，或希望下游开发者为你的 API 自行生成客户端，那么也应包含 `+optional`。

在 1.x 里，获得 `+optional` 的最稳妥方式是使用 `omitempty`。

</aside>

See each subsection for information about different types of code and YAML
generation.

## 在 Kubebuilder 中生成代码与产物

Kubebuilder 项目通常使用两个与 controller-gen 相关的 `make` 目标：

- `make manifests` 生成 Kubernetes 对象 YAML，例如
  [CustomResourceDefinitions](./markers/crd.md)、
  [WebhookConfigurations](./markers/webhook.md) 与 [RBAC 角色](./markers/rbac.md)。

- `make generate` 生成代码，例如 [runtime.Object/DeepCopy 的实现](./markers/object.md)。

完整概览请见[生成 CRD](./generating-crd.md)。

## 标记语法（Marker Syntax）

精确语法可参阅
[controller-tools 的 godocs](https://pkg.go.dev/sigs.k8s.io/controller-tools/pkg/markers?tab=doc)。

一般而言，标记可分为：

- 空标记（Empty，`+kubebuilder:validation:Optional`）：类似命令行里的布尔开关，仅标注即可开启某行为。

- 匿名标记（Anonymous，`+kubebuilder:validation:MaxItems=2`）：接收一个无名参数。

- 多选项标记（Multi-option，
  `+kubebuilder:printcolumn:JSONPath=".status.replicas",name=Replicas,type=string`）：
  接收一个或多个具名参数。第一个参数与标记名以冒号分隔，其后参数以逗号分隔。参数顺序无关，且部分参数可选。

标记参数可以是字符串、整型、布尔、切片或这些类型的映射。字符串、整型和布尔值遵循 Go 语法：

```go
// +kubebuilder:validation:ExclusiveMaximum=false
// +kubebuilder:validation:Format="date-time"
// +kubebuilder:validation:Maximum=42
```

为方便起见，在简单场景下字符串可省略引号（不建议在除单词外的场景使用）：

```go
// +kubebuilder:validation:Type=string
```

切片可以使用花括号加逗号分隔：

```go
// +kubebuilder:webhooks:Enum={"crackers, Gromit, we forgot the crackers!","not even wensleydale?"}
```

或在简单场景下使用分号分隔：

```go
// +kubebuilder:validation:Enum=Wallace;Gromit;Chicken
```

映射以字符串为键、任意类型为值（等价于 `map[string]interface{}`）。使用花括号包裹（`{}`），键值以冒号分隔（`:`），键值对之间以逗号分隔：

```go
// +kubebuilder:default={magic: {numero: 42, stringified: forty-two}}
```
