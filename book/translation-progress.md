# 文档中文化进度（当前会话）

更新时间：2025-09-21 03:17:48 CST

## 已翻译并提交（docs/book/src 下）

以下为本次会话中已完成中文化并逐个提交的文档列表（包含路径）：

- docs/book/src/SUMMARY.md
- docs/book/src/TODO.md
- docs/book/src/architecture.md
- docs/book/src/cronjob-tutorial/api-design.md
- docs/book/src/cronjob-tutorial/basic-project.md
- docs/book/src/cronjob-tutorial/cert-manager.md
- docs/book/src/cronjob-tutorial/controller-implementation.md
- docs/book/src/cronjob-tutorial/controller-overview.md
- docs/book/src/cronjob-tutorial/cronjob-tutorial.md
- docs/book/src/cronjob-tutorial/empty-main.md
- docs/book/src/cronjob-tutorial/epilogue.md
- docs/book/src/cronjob-tutorial/gvks.md
- docs/book/src/cronjob-tutorial/main-revisited.md
- docs/book/src/cronjob-tutorial/new-api.md
- docs/book/src/cronjob-tutorial/other-api-files.md
- docs/book/src/cronjob-tutorial/running-webhook.md
- docs/book/src/cronjob-tutorial/running.md
- docs/book/src/cronjob-tutorial/webhook-implementation.md
- docs/book/src/cronjob-tutorial/writing-tests.md
- docs/book/src/faq.md
- docs/book/src/getting-started.md
- docs/book/src/introduction.md
- docs/book/src/logos/README.md
- docs/book/src/migrations.md
- docs/book/src/multiversion-tutorial/api-changes.md
- docs/book/src/multiversion-tutorial/conversion-concepts.md
- docs/book/src/multiversion-tutorial/conversion.md
- docs/book/src/multiversion-tutorial/deployment.md
- docs/book/src/multiversion-tutorial/tutorial.md
- docs/book/src/multiversion-tutorial/webhooks.md
- docs/book/src/plugins/available-plugins.md
- docs/book/src/plugins/available/go-v4-plugin.md
- docs/book/src/plugins/available/autoupdate-v1-alpha.md
- docs/book/src/plugins/available/deploy-image-plugin-v1-alpha.md
- docs/book/src/plugins/plugins.md
- docs/book/src/quick-start.md
 - docs/book/src/plugins/available/helm-v1-alpha.md
 - docs/book/src/plugins/available/grafana-v1-alpha.md
 - docs/book/src/plugins/available/kustomize-v2.md
 - docs/book/src/plugins/extending.md
 - docs/book/src/plugins/extending/extending_cli_features_and_plugins.md
 - docs/book/src/plugins/extending/external-plugins.md
 - docs/book/src/plugins/extending/testing-plugins.md
 - docs/book/src/plugins/plugins-versioning.md
 - docs/book/src/plugins/to-scaffold-project.md
 - docs/book/src/plugins/to-add-optional-features.md
 - docs/book/src/plugins/to-be-extended.md
 - docs/book/src/plugins/kustomize-v2.md
 - docs/book/src/migration/legacy.md
 - docs/book/src/migration/v3-plugins.md
 - docs/book/src/migration/multi-group.md
 - docs/book/src/migration/v3vsv4.md
 - docs/book/src/migration/migration_guide_gov3_to_gov4.md
 - docs/book/src/migration/manually_migration_guide_gov3_to_gov4.md
 - docs/book/src/migration/legacy/migration_guide_v1tov2.md
 - docs/book/src/migration/legacy/migration_guide_v2tov3.md
 - docs/book/src/migration/legacy/v1vsv2.md
 - docs/book/src/migration/legacy/v2vsv3.md
 - docs/book/src/migration/legacy/manually_migration_guide_v2_v3.md
- docs/book/src/reference/alpha_commands.md
- docs/book/src/reference/commands/alpha_generate.md
- docs/book/src/reference/commands/alpha_update.md
- docs/book/src/reference/reference.md
- docs/book/src/reference/generating-crd.md
- docs/book/src/reference/using-finalizers.md
- docs/book/src/reference/watching-resources.md
- docs/book/src/reference/watching-resources/secondary-owned-resources.md
- docs/book/src/reference/watching-resources/secondary-resources-not-owned.md
- docs/book/src/reference/watching-resources/predicates-with-watch.md
- docs/book/src/reference/kind.md
- docs/book/src/reference/webhook-overview.md
- docs/book/src/reference/admission-webhook.md
- docs/book/src/reference/markers.md
- docs/book/src/reference/markers/crd.md
- docs/book/src/reference/markers/crd-validation.md
- docs/book/src/reference/markers/webhook.md
- docs/book/src/reference/markers/object.md
- docs/book/src/reference/markers/rbac.md
- docs/book/src/reference/markers/scaffold.md
- docs/book/src/reference/pprof-tutorial.md
- docs/book/src/reference/controller-gen.md
- docs/book/src/reference/completion.md
- docs/book/src/reference/artifacts.md
- docs/book/src/reference/platform.md
- docs/book/src/reference/using_an_external_resource.md
- docs/book/src/reference/raising-events.md
- docs/book/src/reference/good-practices.md
- docs/book/src/reference/scopes.md
- docs/book/src/reference/metrics-reference.md
- docs/book/src/reference/envtest.md
- docs/book/src/reference/metrics.md
- docs/book/src/versions_compatibility_supportability.md
- docs/book/src/reference/project-config.md
- docs/book/src/reference/submodule-layouts.md
 - docs/book/src/reference/markers/crd-processing.md
 - docs/book/src/cronjob-tutorial/testdata/project/README.md
 - docs/book/src/getting-started/testdata/project/README.md
 - docs/book/src/multiversion-tutorial/testdata/project/README.md

## 已是中文，已跳过

- docs/book/src/SUMMARY.md
- docs/book/src/architecture.md
- docs/book/src/introduction.md

## 待处理范围（后续优先级建议）

（当前已清空，剩余内容均为中文或非本会话范围）

说明：
- 上述清单基于对 docs/book/src 的逐文件中文字符检测结果自动生成，已排除已翻译条目。
- P5 为脚手架示例仓库 README，翻译收益较低，建议跳过或置于最后。

## 提交约定

- 每完成一个文件翻译，单独 git commit（示例前缀：🌐 (docs): Translate …）。
- 不改语义，不引入无关内容；保留专有名词与代码/命令原样。
- 不添加任何 AI 协作者信息。
