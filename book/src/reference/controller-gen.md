# controller-gen CLI

Kubebuilder 使用
[controller-gen](https://sigs.k8s.io/controller-tools/cmd/controller-gen)
来生成实用代码与 Kubernetes YAML。生成行为由 Go 代码中的特殊[“标记注释”](/reference/markers.md)控制。

controller-gen 由不同的“生成器”（指定生成什么）与“输出规则”（指定输出位置与方式）组成。

二者均通过以[标记格式](/reference/markers.md)书写的命令行选项进行配置。

例如，下述命令：

```shell
controller-gen paths=./... crd:trivialVersions=true rbac:roleName=controller-perms output:crd:artifacts:config=config/crd/bases
```

会生成 CRD 与 RBAC；其中 CRD YAML 被放入 `config/crd/bases`。RBAC 使用默认输出规则（`config/rbac`）。该命令会遍历当前目录树中的所有包（遵循 Go `...` 通配符的规则）。

## 生成器（Generators）

每个生成器通过一个 CLI 选项进行配置。你可以在一次 `controller-gen` 调用中启用多个生成器。

{{#markerdocs CLI: generators}}

## 输出规则（Output Rules）

输出规则决定某个生成器如何输出产物。总会存在一个全局“兜底”输出规则（`output:<rule>`），也可以为某个生成器单独覆盖（`output:<generator>:<rule>`）。

<aside class="note">

<h1>默认规则</h1>

当未显式指定兜底规则时，会采用每个生成器的默认规则：YAML 输出到 `config/<generator>`，代码类产物保持在其所属包中。

等价写法为：对每个生成器使用 `output:<generator>:artifacts:config=config/<generator>`。

一旦显式给出“兜底”规则，则会覆盖默认规则。

例如：指定 `crd rbac:roleName=controller-perms output:crd:stdout` 时，CRD 输出到标准输出，而 RBAC 仍写入 `config/rbac`。若再加一个全局兜底规则（如 `crd rbac:roleName=controller-perms output:crd:stdout output:none`），则 CRD 仍输出到标准输出，而其余产物会被丢弃（/dev/null），因为兜底规则已显式指定。

</aside>

为简洁起见，下方省略了逐生成器的输出规则写法（`output:<generator>:<rule>`）。它们与此处列出的全局兜底选项等价。

{{#markerdocs CLI: output rules (optionally as output:<generator>:...)}}

## 其他选项（Other Options）

{{#markerdocs CLI: generic}}
