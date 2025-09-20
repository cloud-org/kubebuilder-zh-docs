# 从 v2 迁移到 v3

在继续之前，请先了解 [Kubebuilder v2 与 v3 的差异][v2vsv3]。

请确保已按[安装指南][quick-start]安装所需组件。

推荐的迁移方式是：新建一个 v3 项目，然后将 API 与调谐（reconciliation）代码拷贝过去。这样最终得到的项目就是原生的 v3 布局。
在某些情况下，也可以“就地升级”（复用 v2 的项目布局，同时升级 [controller-runtime][controller-runtime] 与 [controller-tools][controller-tools]）。

## 初始化 v3 项目

<aside class="note">
<h1>项目名</h1>

以下示例使用 `migration-project` 作为项目名，`tutorial.kubebuilder.io` 作为域名；请按你的实际情况替换。

</aside>

新建一个以项目名命名的目录。注意该名称会用于脚手架中，默认影响 manager Pod 的名称以及其部署的 Namespace：

```bash
$ mkdir migration-project-name
$ cd migration-project-name
```

初始化 v3 项目前，若不在 `GOPATH` 内，建议先初始化一个新的 Go 模块（在 `GOPATH` 内虽非必须，但仍推荐）：

```bash
go mod init tutorial.kubebuilder.io/migration-project
```

<aside class="note">
<h1>项目模块</h1>

模块名位于项目根目录的 `go.mod`：

```
module tutorial.kubebuilder.io/migration-project
```

</aside>

然后使用 kubebuilder 完成初始化：

```bash
kubebuilder init --domain tutorial.kubebuilder.io
```

<aside class="note">
<h1>项目域名</h1>

域名位于 PROJECT 文件：

```yaml
...
domain: tutorial.kubebuilder.io
...
```
</aside>

## 迁移 API 与 Controller

接下来重新脚手架 API 类型与控制器。

<aside class="note">
<h1>同时脚手架 API 与 Controller</h1>

以下示例假设需要同时生成 API 与 Controller，实际需以旧项目当初的脚手架选择为准。

</aside>

```bash
kubebuilder create api --group batch --version v1 --kind CronJob
```

### 迁移 API

<aside class="note">
<h1>若使用多 Group</h1>

在迁移 API 与 Controller 之前，先运行 `kubebuilder edit --multigroup=true` 以启用多 Group 支持。详见 [multi-group][multi-group]。

</aside>

现在，把旧项目中的 `api/v1/<kind>_types.go` 拷贝到新项目中。

这些文件在新插件中没有功能性修改，因此可以直接用旧文件覆盖新生成的文件。若存在格式差异，也可以只拷贝类型定义本身。

### 迁移 Controller

将旧项目中的 `controllers/cronjob_controller.go` 迁移到新项目。此处存在一个破坏性变化，且可能出现一些格式差异。

新的 `Reconcile` 方法现在将 `context` 作为入参，而不再需要 `context.Background()`。你可以将旧控制器中的其它逻辑复制到新脚手架的方法中，将：

```go
func (r *CronJobReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
    ctx := context.Background()
    log := r.Log.WithValues("cronjob", req.NamespacedName)
```

替换为：

```go
func (r *CronJobReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("cronjob", req.NamespacedName)
```

<aside class="note warning">
<h1>controller-runtime 升级的破坏性变更</h1>

请查看 [controller-runtime 0.8.0+ 的发布说明][controller-runtime] 中的破坏性变更。

</aside>

## 迁移 Webhook

<aside class="note">
<h1>可跳过</h1>

如果项目未使用任何 Webhook，可跳过本节。

</aside>

为 CRD（CronJob）脚手架 Webhook。需要带上 `--defaulting` 与 `--programmatic-validation`（示例项目用到了默认化与校验 Webhook）：

```bash
kubebuilder create webhook --group batch --version v1 --kind CronJob --defaulting --programmatic-validation
```

然后，将旧项目中的 `api/v1/<kind>_webhook.go` 拷贝到新项目中。

## 其他

如果 v2 的 `main.go` 有手工改动，需要迁移到新项目的 `main.go` 中。同时确保所有需要的 scheme 都已注册。

若 `config` 目录下存在新增清单，同步迁移它们。

如有需要，请更新 Makefile 中的镜像名等配置。

## 验证

最后，运行 `make` 与 `make docker-build` 以确认一切正常。

[v2vsv3]: v2vsv3.md
[quick-start]: /quick-start.md#installation
[controller-tools]: https://github.com/kubernetes-sigs/controller-tools/releases
[controller-runtime]: https://github.com/kubernetes-sigs/controller-runtime/releases
[multi-group]: /migration/multi-group.md
