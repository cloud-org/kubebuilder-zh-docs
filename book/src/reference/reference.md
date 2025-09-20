# 参考（Reference）

  - [生成 CRD](generating-crd.md)
  - [使用 Finalizer](using-finalizers.md)
    Finalizer 允许在对象从 Kubernetes 集群中被删除之前执行自定义清理逻辑。
  - [监听资源（Watching Resources）](watching-resources.md)
    监听集群中的资源变化，并据此采取动作。
      - [监听由控制器“拥有”的二级资源](watching-resources/secondary-owned-resources.md)
      - [监听非本控制器“拥有”的二级资源](watching-resources/secondary-resources-not-owned.md)
      - [使用 Predicates 精细化 Watch](watching-resources/predicates-with-watch.md)
  - [Kind 集群](kind.md)
  - [什么是 Webhook？](webhook-overview.md)
    Webhook 是 HTTP 回调。Kubernetes 中常见的三类：1）准入（admission）2）CRD 转换 3）鉴权（authorization）。
    - [Admission Webhook](admission-webhook.md)
      Admission Webhook 在 API Server 接受对象前，用于变更或校验资源。
  - [用于配置/代码生成的标记（Markers）](markers.md)

      - [CRD 生成](markers/crd.md)
      - [CRD 校验](markers/crd-validation.md)
      - [Webhook](markers/webhook.md)
      - [Object/DeepCopy](markers/object.md)
      - [RBAC](markers/rbac.md)
      - [脚手架（Scaffold）](markers/scaffold.md)

  - [使用 Pprof 进行监控](pprof-tutorial.md)
  - [controller-gen CLI](controller-gen.md)
  - [命令自动补全（completion）](completion.md)
  - [构建产物（Artifacts）](artifacts.md)
  - [平台支持](platform.md)

  - [子模块布局（Sub-Module Layouts）](submodule-layouts.md)
  - [使用外部资源 / API](using_an_external_resource.md)

  - [指标（Metrics）](metrics.md)
      - [指标参考](metrics-reference.md)

  - [CLI 插件](../plugins/plugins.md)
