# 生成 CRD（Generating CRDs）

Kubebuilder 使用名为 [`controller-gen`][controller-tools] 的工具来生成实用代码与 Kubernetes 对象 YAML（例如 CRD）。

它依赖源码中的特殊“标记注释”（以 `// +` 开头）来为字段、类型与包提供额外元信息。针对 CRD，相关标记通常写在你的 `_types.go` 文件中。更多标记说明请参考[标记参考文档][marker-ref]。

Kubebuilder 提供了一个 `make` 目标来运行 controller-gen 以生成 CRD：`make manifests`。

执行 `make manifests` 后，你会在 `config/crd/bases` 目录下看到生成的 CRD。`make manifests` 还会生成其他若干产物——详见[标记参考文档][marker-ref]。

## 校验（Validation）

CRD 在其 `validation` 段落中通过 [OpenAPI v3 schema][openapi-schema] 支持[声明式校验][kube-validation]。

通常，[校验相关标记](./markers/crd-validation.md)可以加在字段或类型上。若校验逻辑较复杂、需要复用，或需要校验切片元素，建议定义一个新的类型以承载你的校验描述。

例如：

```go
type ToySpec struct {
	// +kubebuilder:validation:MaxLength=15
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name,omitempty"`

	// +kubebuilder:validation:MaxItems=500
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:UniqueItems=true
	Knights []string `json:"knights,omitempty"`

	Alias   Alias   `json:"alias,omitempty"`
	Rank    Rank    `json:"rank"`
}

// +kubebuilder:validation:Enum=Lion;Wolf;Dragon
type Alias string

// +kubebuilder:validation:Minimum=1
// +kubebuilder:validation:Maximum=3
// +kubebuilder:validation:ExclusiveMaximum=false
type Rank int32

```

## 自定义输出列（Additional Printer Columns）

自 Kubernetes 1.11 起，`kubectl get` 可以向服务端询问应显示哪些列。对 CRD 而言，这使其能像内建资源一样，在 `kubectl get` 中展示更贴合类型的信息。

展示哪些信息由 CRD 的 [additionalPrinterColumns 字段][kube-additional-printer-columns]控制，而该字段又由你在 Go 类型上标注的 [`+kubebuilder:printcolumn`][crd-markers] 标记决定。

例如，下面示例为之前的校验示例添加几列，显示 `alias`、`rank` 与 `knights` 的信息：

```go
// +kubebuilder:printcolumn:name="Alias",type=string,JSONPath=`.spec.alias`
// +kubebuilder:printcolumn:name="Rank",type=integer,JSONPath=`.spec.rank`
// +kubebuilder:printcolumn:name="Bravely Run Away",type=boolean,JSONPath=`.spec.knights[?(@ == "Sir Robin")]`,description="when danger rears its ugly head, he bravely turned his tail and fled",priority=10
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type Toy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ToySpec   `json:"spec,omitempty"`
	Status ToyStatus `json:"status,omitempty"`
}

```

## 子资源（Subresources）

自 Kubernetes 1.13 起，CRD 可以选择实现 `/status` 与 `/scale` [子资源][kube-subresources]。

一般建议：凡是具有 `status` 字段的资源，都应启用 `/status` 子资源。

上述两个子资源均有对应的[标记][crd-markers]。

### Status

使用 `+kubebuilder:subresource:status` 启用 `status` 子资源。启用后，对主资源的更新不会直接修改其 `status`；同样，对 status 子资源的更新也只能修改 `status` 字段。

例如：

```go
// +kubebuilder:subresource:status
type Toy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ToySpec   `json:"spec,omitempty"`
	Status ToyStatus `json:"status,omitempty"`
}
```

### Scale

使用 `+kubebuilder:subresource:scale` 启用 `scale` 子资源。启用后，用户可以对你的资源使用 `kubectl scale`。若 `selectorpath` 指向标签选择器的字符串形式，HPA 也能自动伸缩你的资源。

例如：

```go
type CustomSetSpec struct {
	Replicas *int32 `json:"replicas"`
}

type CustomSetStatus struct {
	Replicas int32 `json:"replicas"`
    Selector string `json:"selector"` // this must be the string form of the selector
}


// +kubebuilder:subresource:status
// +kubebuilder:subresource:scale:specpath=.spec.replicas,statuspath=.status.replicas,selectorpath=.status.selector
type CustomSet struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CustomSetSpec   `json:"spec,omitempty"`
	Status CustomSetStatus `json:"status,omitempty"`
}
```

## 多版本（Multiple Versions）

自 Kubernetes 1.13 起，你可以在同一个 CRD 中定义某个 Kind 的多个版本，并通过 Webhook 在版本间进行转换。

更多细节见[多版本教程](/multiversion-tutorial/tutorial.md)。

出于与旧版 Kubernetes 的兼容性考虑，Kubebuilder 默认不会为不同版本生成不同的校验规则。

如需启用，请修改 Makefile 中的选项：若使用 v1beta CRD，将 `CRD_OPTIONS ?= "crd:trivialVersions=true,preserveUnknownFields=false` 改为 `CRD_OPTIONS ?= crd:preserveUnknownFields=false`；若使用 v1（推荐），则为 `CRD_OPTIONS ?= crd`。

随后，可使用 `+kubebuilder:storageversion` [标记][crd-markers] 指定由 API Server 用于持久化数据的 [GVK](/cronjob-tutorial/gvks.md "Group-Version-Kind")。

## 实现细节（Under the hood）

Kubebuilder 通过脚手架提供了运行 `controller-gen` 的 make 规则；若本地尚无该可执行文件，会使用 Go Modules 的 `go install` 自动安装。

你也可以直接运行 `controller-gen` 来观察其行为。

controller-gen 的每个“生成器”都通过命令行选项进行控制（语法与标记一致）。同时它也支持不同的输出“规则”，用于控制产物的输出位置与形式。如下所示为 `manifests` 规则（为示例简化为仅生成 CRD）：

```makefile
# Generate manifests for CRDs
manifests: controller-gen
	$(CONTROLLER_GEN) rbac:roleName=manager-role crd webhook paths="./..." output:crd:artifacts:config=config/crd/bases
```

它通过 `output:crd:artifacts` 输出规则将与 CRD 相关的配置类（非代码）产物输出至 `config/crd/bases`，而非 `config/crd`。

想要查看 `controller-gen` 的所有生成器与选项，运行：

```shell
controller-gen -h
```

or, for more details:

```shell
$ controller-gen -hhh
```

[marker-ref]: ./markers.md "Markers for Config/Code Generation"

[kube-validation]: https://kubernetes.io/docs/tasks/access-kubernetes-api/custom-resources/custom-resource-definitions/#validation "Custom Resource Definitions: Validation"

[openapi-schema]: https://github.com/OAI/OpenAPI-Specification/blob/master/versions/3.0.0.md#schemaObject "OpenAPI v3"

[kube-additional-printer-columns]: https://kubernetes.io/docs/tasks/access-kubernetes-api/custom-resources/custom-resource-definitions/#additional-printer-columns "Custom Resource Definitions: Additional Printer Columns"

[kube-subresources]: https://kubernetes.io/docs/tasks/access-kubernetes-api/custom-resources/custom-resource-definitions/#status-subresource "Custom Resource Definitions: Status Subresource"

[crd-markers]: ./markers/crd.md "CRD Generation"

[controller-tools]: https://sigs.k8s.io/controller-tools "Controller Tools"
