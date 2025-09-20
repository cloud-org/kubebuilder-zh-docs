# Kubebuilder v2 与 v3 对比（Legacy：从 v2.0.0+ 布局到 3.0.0+）

本文覆盖从 v2 迁移到 v3 时的所有破坏性变更。

更多（含非破坏性）变更详情可参考
[controller-runtime][controller-runtime]、
[controller-tools][controller-tools]
以及 [kb-releases][kb-releases] 的发布说明。

## 共同变化（Common changes）

v3 项目使用 Go modules，且要求 Go 1.18+；不再支持使用 Dep 管理依赖。

## Kubebuilder

- 引入对插件的初步支持。详见[可扩展 CLI 与脚手架插件：Phase 1][plugins-phase1-design-doc]、
  [Phase 1.5][plugins-phase1-design-doc-1.5] 与 [Phase 2][plugins-phase2-design-doc] 的设计文档；亦可参考[插件章节][plugins-section]。

- `PROJECT` 文件采用了新布局，记录更多资源信息，以便插件在脚手架时做出合理决策。

    另外，`PROJECT` 文件本身也引入版本：`version` 字段表示 `PROJECT` 文件版本；`layout` 字段表示脚手架与主插件版本。

- `gcr.io/kubebuilder/kube-rbac-proxy` 镜像版本从 `0.5.0` 升级到 `0.11.0`（该组件默认开启，用于保护 manager 的请求），以解决安全问题。详情见 [kube-rbac-proxy][kube-rbac-proxy]。

## 新版 `go/v3` 插件要点（TL;DR）

更多细节见 [kubebuilder 发布说明][kb-releases]，核心高亮如下：

<aside class="note">
<h1>默认插件</h1>
由 Kubebuilder v3 脚手架生成的项目默认使用 `go.kubebuilder.io/v3` 插件。
</aside>

- 生成的 API/清单变化：
  * 生成的 CRD 使用 `apiextensions/v1`（`apiextensions/v1beta1` 在 Kubernetes `1.16` 中已弃用）
  * 生成的 Webhook 使用 `admissionregistration.k8s.io/v1`（`v1beta1` 在 Kubernetes `1.16` 中已弃用）
  * 使用 Webhook 时，证书管理切换为 `cert-manager.io/v1`（`v1alpha2` 在 Cert-Manager `0.14` 中弃用，参见[文档][cert-manager-docs]）

- 代码变化：
  * manager 的 `--metrics-addr` 与 `enable-leader-election` 现更名为 `--metrics-bind-address` 与 `--leader-elect`，与 Kubernetes 核心组件命名保持一致。详见 [#1839][issue-1893]
  * 默认添加存活/就绪探针，使用 [`healthz.Ping`][healthz-ping]
  * 新增以 ComponentConfig 方式创建项目的选项，详见[增强提案][enhancement proposal]与[教程][component-config-tutorial]
  * Manager 清单默认使用 `SecurityContext` 以提升安全性，详见 [#1637][issue-1637]
- 其他：
  * 支持 [controller-tools][controller-tools] `v0.9.0`（`go/v2` 为 `v0.3.0`，更早为 `v0.2.5`）
  * 支持 [controller-runtime][controller-runtime] `v0.12.1`（`go/v2` 为 `v0.6.4`，更早为 `v0.5.0`）
  * 支持 [kustomize][kustomize] `v3.8.7`（`go/v2` 为 `v3.5.4`，更早为 `v3.1.0`）
  * 自动下载所需的 Envtest 二进制
  * 最低 Go 版本升至 `1.18`（此前为 `1.13`）

<aside class="note warning">
<h1>项目定制化</h1>

创建项目后你可以自由定制，但除非非常清楚后果，否则不建议偏离推荐布局。

例如，不要随意移动脚手架生成的文件，否则会给后续升级带来阻碍，并可能失去部分 CLI 能力与辅助特性。项目布局详见[基础项目包含什么？][basic-project-doc]

</aside>

## 迁移到 Kubebuilder v3

若希望升级到最新脚手架特性，请参考以下指南，获得最直观的步骤：

<aside class="note warning">
<h1>Apple Silicon（M1）</h1>

`go/v3` 使用的 [kubernetes-sigs/kustomize][kustomize] v3 不提供 Apple Silicon（`darwin/arm64`）可用的二进制。
因此可以直接使用支持该平台的 `go/v4` 插件：

```bash
kubebuilder init --domain my.domain --repo my.domain/guestbook --plugins=go/v4
```

</aside>

- [v2 → v3 迁移指南][migration-guide-v2-to-v3]（推荐）

### 通过手动更新文件

若希望在不改变现有脚手架的前提下使用最新 Kubebuilder CLI，可参考下述“仅更新 PROJECT 版本并切换插件版本”的手动步骤。

该方式复杂、易错且不保证成功；并且不会获得默认脚手架文件中的改进与修复。

你仍可通过 `go/v2` 插件继续使用旧布局（不会把 [controller-runtime][controller-runtime] 与 [controller-tools][controller-tools] 升至 `go/v3` 所用的版本，以避免破坏性变更）。本文也提供了如何手动修改文件以切换到 `go/v3` 插件与依赖版本的说明。

- [通过手动更新文件迁移到 Kubebuilder v3][manually-upgrade]

[plugins-phase1-design-doc]: https://github.com/kubernetes-sigs/kubebuilder/blob/master/designs/extensible-cli-and-scaffolding-plugins-phase-1.md
[plugins-phase1-design-doc-1.5]: https://github.com/kubernetes-sigs/kubebuilder/blob/master/designs/extensible-cli-and-scaffolding-plugins-phase-1-5.md
[plugins-phase2-design-doc]: https://github.com/kubernetes-sigs/kubebuilder/blob/master/designs/extensible-cli-and-scaffolding-plugins-phase-2.md
[plugins-section]: ./../../plugins/plugins.md
[manually-upgrade]: manually_migration_guide_v2_v3.md
[component-config-tutorial]: ../../component-config-tutorial/tutorial.md
[issue-1893]: https://github.com/kubernetes-sigs/kubebuilder/issues/1839
[migration-guide-v2-to-v3]: migration_guide_v2tov3.md
[healthz-ping]: https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/healthz#CheckHandler
[controller-runtime]: https://github.com/kubernetes-sigs/controller-runtime/releases
[controller-tools]: https://github.com/kubernetes-sigs/controller-tools/releases
[kustomize]: https://github.com/kubernetes-sigs/kustomize/releases
[issue-1637]: https://github.com/kubernetes-sigs/kubebuilder/issues/1637
[enhancement proposal]: https://github.com/kubernetes/enhancements/tree/master/keps/sig-cluster-lifecycle/wgs
[cert-manager-docs]: https://cert-manager.io/docs/installation/upgrading/
[kb-releases]: https://github.com/kubernetes-sigs/kubebuilder/releases
[kube-rbac-proxy]: https://github.com/brancz/kube-rbac-proxy/releases
[basic-project-doc]: ../../cronjob-tutorial/basic-project.md
[kustomize]: https://github.com/kubernetes-sigs/kustomize
