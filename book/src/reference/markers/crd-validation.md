# CRD 校验（CRD Validation）

以下标记用于控制针对相应类型或字段生成的 CRD 校验 schema。每个标记大致对应一个 OpenAPI/JSON schema 选项。

示例参见[生成 CRD](/reference/generating-crd.md)。

<aside class="note">
<h1>关于文档中的标记分组</h1>

某些标记看起来似乎重复。实际上，这些标记会按照使用上下文（例如：字段、类型、数组元素）进行分组。例如 `+kubebuilder:validation:Enum` 既可用于单个字段，也可用于数组元素。文档通过分组来体现该复用能力。

</aside>


{{#markerdocs CRD validation}}
