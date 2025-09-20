# 插件版本管理

| 名称 | 示例 | 描述 |
|---|---|---|
| Kubebuilder 版本 | `v2.2.0`, `v2.3.0`, `v2.3.1`, `v4.2.0` | Kubebuilder 项目的打标签版本，代表本仓库源码的变更。二进制请见 [releases][kb-releases]。 |
| 项目版本（Project version） | `"1"`, `"2"`, `"3"` | Project version 定义 `PROJECT` 配置文件的 schema，即 `PROJECT` 中的 `version` 字段。 |
| 插件版本（Plugin version） | `v2`, `v3`, `v4` | 某个插件自身的版本以及其生成的脚手架版本。版本体现在插件键上，例如 `go.kubebuilder.io/v2`。详见[设计文档][cli-plugins-versioning]。 |

### 版本递增（Incrementing versions）

关于 Kubebuilder 发布版本的规范，参见 [semver][semver]。

仅当 PROJECT 文件的 schema 自身发生破坏性变更时，才应提升 Project version。Go 脚手架或 Kubebuilder CLI 的改动并不会影响 Project version。

类似地，引入新的插件版本往往只会带来 Kubebuilder 的次版本发布，因为 CLI 本身并未发生破坏性变更。只有当我们移除旧插件版本的支持时，才会对 Kubebuilder 本身构成破坏性变更。更多细节见插件设计文档的[版本管理章节][cli-plugins-versioning]。

## 对插件引入变更

只有当改动会破坏由旧版本插件脚手架的项目时，才需要提升插件版本。一旦 `vX` 稳定（不再带有 `alpha`/`beta` 后缀），应创建一个新包并在其中提供 `v(X+1)-alpha` 版本的插件。通常做法是“语义上复制”：`cp -r pkg/plugins/golang/vX pkg/plugins/golang/v(X+1)`，然后更新版本号与路径。随后所有破坏性变更都应只在新包中进行；`vX` 版本不再接受破坏性变更。

另外，你必须在 PR 中向 Kubebuilder Book 的 [migrations][migrations] 部分补充迁移指南，详细说明用户如何从 `vX` 升级到 `v(X+1)-alpha`。

<aside class="note">
<h1>示例</h1>

默认情况下，Kubebuilder 使用 `go.kubebuilder.io/v4` 脚手架项目。

假设你新增了一个特性：在 `init` 脚手架生成的 `main.go` 中加入了一个新 marker，后续 `create api` 会依赖此 marker 更新该文件。对于已用 `go.kubebuilder.io/v4` 创建的项目，如果不先手工更新，就会触发错误。

因此，这属于对 `go.kubebuilder.io` 插件的破坏性变更，只能合并到 `v5-alpha` 插件版本中（该插件包应已存在）。

</aside>

[semver]: https://semver.org/
[migrations]: ../migrations.md
[kb-releases]:https://github.com/kubernetes-sigs/kubebuilder/releases
[design-doc]: ./extending
[cli-plugins-versioning]:./extending#plugin-versioning
