# go/v3 与 go/v4 对比

本文覆盖从使用 `go/v3` 插件（自 `2021-04-28` 起为默认脚手架）构建的项目迁移到新版 `go/v4` 插件时的所有破坏性变更。

更多（含非破坏性）变更详情可参考：

- [controller-runtime][controller-runtime]
- [controller-tools][controller-tools]
- [kustomize][kustomize-release]
- [kb-releases][kb-releases] 的发布说明。

## 共同变化（Common changes）

- `go/v4` 项目使用 Kustomize v5x（不再是 v3x）。
- `config/` 目录下若干清单已调整，去除了 Kustomize 的废弃用法（例如环境变量）。
- `config/samples` 下新增 `kustomization.yaml`，可通过 `kustomize build config/samples` 简单灵活地生成样例清单。
- 增加对 Apple Silicon M1（darwin/arm64）的支持。
- 移除对 Kubernetes `v1beta1` 版 CRD/Webhook API 的支持（自 k8s 1.22 起废弃）。
- 不再脚手架引用 `"k8s.io/api/admission/v1beta1"` 的 webhook 测试文件（该 API 自 k8s 1.25 起不再提供）；默认改为 `"k8s.io/api/admission/v1"`（自 k8s 1.20 可用）。
- 不再保证兼容 k8s `< 1.16`。
- 布局调整以贴合社区对[标准 Go 项目结构][standard-go-project]的诉求：API 置于 `api/`，控制器置于 `internal/`，`main.go` 置于 `cmd/`。

<aside class="note">
<H1>关于 go/v4 插件</H1>

详见[go/v4 插件章节][go/v4-doc]。

</aside>

## 新版 `go/v4` 插件要点（TL;DR）

更多细节见 [kubebuilder 发布说明][kb-releases]，核心高亮如下：

<aside class="note warning">
<h1>项目定制化</h1>

你可根据需要定制项目，但除非非常清楚后果，否则不建议偏离推荐布局。

例如，不要随意移动脚手架生成的文件，否则会给后续升级带来阻碍，并可能失去部分 CLI 能力与辅助特性。项目布局详见[基础项目包含什么？][basic-project-doc]

</aside>

## 迁移到 Kubebuilder go/v4

若希望升级到最新脚手架特性，请参考以下指南获取最直观的步骤，帮助你获得全部改进：

- [从 go/v3 迁移到 go/v4][migration-guide-gov3-to-gov4]（推荐）

### 通过手动更新文件

若希望在不改变现有脚手架的前提下使用最新 Kubebuilder CLI，可参考下述“仅更新 PROJECT 版本并切换插件版本”的手动步骤。

该方式复杂、易错且不保证成功；并且不会获得默认脚手架文件中的改进与修复。

- [通过手动更新文件迁移到 go/v4][manually-upgrade]

[plugins-phase1-design-doc]: https://github.com/kubernetes-sigs/kubebuilder/blob/master/designs/extensible-cli-and-scaffolding-plugins-phase-1.md
[plugins-phase1-design-doc-1.5]: https://github.com/kubernetes-sigs/kubebuilder/blob/master/designs/extensible-cli-and-scaffolding-plugins-phase-1-5.md
[plugins-phase2-design-doc]: https://github.com/kubernetes-sigs/kubebuilder/blob/master/designs/extensible-cli-and-scaffolding-plugins-phase-2.md
[plugins-section]: ./../plugins/plugins.md
[kustomize]: https://github.com/kubernetes-sigs/kustomize/releases/tag/kustomize%2Fv4.0.0
[go/v4-doc]: ./../plugins/available/go-v4-plugin.md
[migration-guide-gov3-to-gov4]: migration_guide_gov3_to_gov4.md
[manually-upgrade]: manually_migration_guide_gov3_to_gov4.md
[basic-project-doc]: ./../cronjob-tutorial/basic-project.md
[standard-go-project]: https://github.com/golang-standards/project-layout
[controller-runtime]: https://github.com/kubernetes-sigs/controller-runtime
[controller-tools]: https://github.com/kubernetes-sigs/controller-tools
[kustomize-release]: https://github.com/kubernetes-sigs/kustomize/releases/tag/kustomize%2Fv5.0.0
[kb-releases]: https://github.com/kubernetes-sigs/kubebuilder/releases
