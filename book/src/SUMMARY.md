# 目录

[简介](./introduction.md)

[架构](./architecture.md)

[快速开始](./quick-start.md)

[入门](./getting-started.md)

[版本兼容性与可支持性](./versions_compatibility_supportability.md)
---

- [教程：构建 CronJob](cronjob-tutorial/cronjob-tutorial.md)

  - [基础项目包含什么？](./cronjob-tutorial/basic-project.md)
  - [每段旅程都有起点，每个程序都有 main](./cronjob-tutorial/empty-main.md)
  - [Groups、Versions 与 Kind，哇哦！](./cronjob-tutorial/gvks.md)
  - [添加一个新的 API](./cronjob-tutorial/new-api.md)
  - [设计一个 API](./cronjob-tutorial/api-design.md)

    - [简短插曲：其他这些东西是什么？](./cronjob-tutorial/other-api-files.md)

  - [控制器包含什么？](./cronjob-tutorial/controller-overview.md)
  - [实现一个控制器](./cronjob-tutorial/controller-implementation.md)

    - [你刚才提到 main？](./cronjob-tutorial/main-revisited.md)

  - [实现 defaulting/validating Webhook](./cronjob-tutorial/webhook-implementation.md)
  - [运行并部署控制器](./cronjob-tutorial/running.md)

    - [部署 cert-manager](./cronjob-tutorial/cert-manager.md)
    - [部署 Webhook](./cronjob-tutorial/running-webhook.md)

  - [编写测试](./cronjob-tutorial/writing-tests.md)

  - [尾声](./cronjob-tutorial/epilogue.md)

- [教程：多版本 API](./multiversion-tutorial/tutorial.md)

  - [变更内容](./multiversion-tutorial/api-changes.md)
  - [Hub、Spoke 以及其他轮式隐喻](./multiversion-tutorial/conversion-concepts.md)
  - [实现转换](./multiversion-tutorial/conversion.md)

    - [并配置 Webhook](./multiversion-tutorial/webhooks.md)

  - [部署与测试](./multiversion-tutorial/deployment.md)

---

- [迁移](./migrations.md)

  - [旧版（<= v3.0.0 之前）](./migration/legacy.md)
    - [Kubebuilder v1 与 v2](migration/legacy/v1vsv2.md)

      - [迁移指南](./migration/legacy/migration_guide_v1tov2.md)

    - [Kubebuilder v2 与 v3](migration/legacy/v2vsv3.md)

      - [迁移指南](migration/legacy/migration_guide_v2tov3.md)
      - [通过更新文件进行迁移](migration/legacy/manually_migration_guide_v2_v3.md)
  - [从 v3.0.0 起（带插件）](./migration/v3-plugins.md)
    - [go/v3 与 go/v4](migration/v3vsv4.md)

      - [迁移指南](migration/migration_guide_gov3_to_gov4.md)
      - [通过更新文件进行迁移](migration/manually_migration_guide_gov3_to_gov4.md)
  - [单组到多组](./migration/multi-group.md)

- [Alpha 命令](./reference/alpha_commands.md)

  - [alpha generate](./reference/commands/alpha_generate.md)
  - [alpha update](./reference/commands/alpha_update.md)

---

- [参考](./reference/reference.md)

  - [生成 CRD](./reference/generating-crd.md)
  - [使用 Finalizer](./reference/using-finalizers.md)
  - [最佳实践](./reference/good-practices.md)
  - [触发事件](./reference/raising-events.md)
  - [监视资源](./reference/watching-resources.md)
    - [被拥有的资源](./reference/watching-resources/secondary-owned-resources.md)
    - [非拥有的资源](./reference/watching-resources/secondary-resources-not-owned.md)
    - [使用谓词](./reference/watching-resources/predicates-with-watch.md)
  - [用于开发与 CI 的 Kind](reference/kind.md)
  - [什么是 Webhook？](reference/webhook-overview.md)
    - [准入 Webhook](reference/admission-webhook.md)
  - [用于配置/代码生成的 Marker](./reference/markers.md)

    - [CRD 生成](./reference/markers/crd.md)
    - [CRD 校验](./reference/markers/crd-validation.md)
    - [CRD 处理](./reference/markers/crd-processing.md)
    - [Webhook](./reference/markers/webhook.md)
    - [对象/DeepCopy](./reference/markers/object.md)
    - [RBAC](./reference/markers/rbac.md)
    - [Scaffold](./reference/markers/scaffold.md)

  - [controller-gen CLI](./reference/controller-gen.md)
  - [completion](./reference/completion.md)
  - [构建产物](./reference/artifacts.md)
  - [平台支持](./reference/platform.md)
  - [使用 pprof 进行监控](./reference/pprof-tutorial.md)

  - [Manager 与 CRD 作用域](./reference/scopes.md)

  - [子模块布局](./reference/submodule-layouts.md)
  - [使用外部资源/API](./reference/using_an_external_resource.md)

  - [配置 EnvTest](./reference/envtest.md)

  - [指标](./reference/metrics.md)

    - [参考](./reference/metrics-reference.md)

  - [项目配置](./reference/project-config.md)

---

- [插件][plugins]

  - [可用插件](./plugins/available-plugins.md)
    - [autoupdate/v1-alpha](./plugins/available/autoupdate-v1-alpha.md)
    - [deploy-image/v1-alpha](./plugins/available/deploy-image-plugin-v1-alpha.md)
    - [go/v4](./plugins/available/go-v4-plugin.md)
    - [grafana/v1-alpha](./plugins/available/grafana-v1-alpha.md)
    - [helm/v1-alpha](./plugins/available/helm-v1-alpha.md)
    - [kustomize/v2](./plugins/available/kustomize-v2.md)
  - [扩展](./plugins/extending.md)
    - [CLI 与插件](./plugins/extending/extending_cli_features_and_plugins.md)
    - [外部插件](./plugins/extending/external-plugins.md)
    - [E2E 测试](./plugins/extending/testing-plugins.md)
  - [插件版本管理](./plugins/plugins-versioning.md)


---

[常见问题](./faq.md)

[plugins]: ./plugins/plugins.md
