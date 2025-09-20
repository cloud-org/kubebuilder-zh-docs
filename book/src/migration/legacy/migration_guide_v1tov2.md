# 从 v1 迁移到 v2

在继续之前，请先了解 [Kubebuilder v1 与 v2 的差异](./v1vsv2.md)。

请确保已按[安装指南](/quick-start.md#installation)安装所需组件。

推荐的迁移方式是：新建一个 v2 项目，然后将 API 与调谐（reconciliation）代码拷贝过去。这样最终得到的项目就是原生的 v2 布局。
在某些情况下，也可以“就地升级”（复用 v1 的项目布局，升级 controller-runtime 与 controller-tools）。

下面以一个 v1 项目为例迁移到 Kubebuilder v2。最终效果应与[示例 v2 项目][v2-project]一致。

## 准备工作

首先确认 Group、Version、Kind 与 Domain。

先看一个 v1 项目的目录结构：

```
pkg/
├── apis
│   ├── addtoscheme_batch_v1.go
│   ├── apis.go
│   └── batch
│       ├── group.go
│       └── v1
│           ├── cronjob_types.go
│           ├── cronjob_types_test.go
│           ├── doc.go
│           ├── register.go
│           ├── v1_suite_test.go
│           └── zz_generated.deepcopy.go
├── controller
└── webhook
```

所有 API 信息都在 `pkg/apis/batch` 下，可以在那里找到所需信息。

In `cronjob_types.go`, we can find

```go
type CronJob struct {...}
```

In `register.go`, we can find

```go
SchemeGroupVersion = schema.GroupVersion{Group: "batch.tutorial.kubebuilder.io", Version: "v1"}
```

据此可知 Kind 为 `CronJob`，Group/Version 为 `batch.tutorial.kubebuilder.io/v1`。

## 初始化 v2 项目

现在初始化 v2 项目。在此之前，若不在 `GOPATH` 下，先初始化一个新的 Go 模块：

```bash
go mod init tutorial.kubebuilder.io/project
```

随后用 kubebuilder 完成项目初始化：

```bash
kubebuilder init --domain tutorial.kubebuilder.io
```

## 迁移 API 与 Controller

接下来重新脚手架 API 类型与控制器。因为两者都需要，交互提示时分别选择生成 API 与 Controller：

```bash
kubebuilder create api --group batch --version v1 --kind CronJob
```

如果你使用多 Group，需要做一些手工迁移，详见[/migration/multi-group.md](/migration/multi-group.md)。

### 迁移 API

将 `pkg/apis/batch/v1/cronjob_types.go` 中的类型拷贝到 `api/v1/cronjob_types.go`。仅需要复制 `Spec` 与 `Status` 字段的实现。

可以把 `+k8s:deepcopy-gen:interfaces=...` 标记（在 Kubebuilder 中已[弃用](/reference/markers/object.md)）替换为 `+kubebuilder:object:root=true`。

以下标记已无需保留（它们来自非常老的 Kubebuilder 版本）：

```go
// +genclient
// +k8s:openapi-gen=true
```

API 类型应类似：

```go
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// CronJob is the Schema for the cronjobs API
type CronJob struct {...}

// +kubebuilder:object:root=true

// CronJobList contains a list of CronJob
type CronJobList struct {...}
```

### 迁移 Controller

将 `pkg/controller/cronjob/cronjob_controller.go` 中的调谐器代码迁移到 `controllers/cronjob_controller.go`：

We'll need to copy
- the fields from the `ReconcileCronJob` struct to `CronJobReconciler`
- the contents of the `Reconcile` function
- the [rbac related markers](/reference/markers/rbac.md) to the new file.
- the code under `func add(mgr manager.Manager, r reconcile.Reconciler) error`
to `func SetupWithManager`

## 迁移 Webhook

如果项目未使用 Webhook，可跳过本节。

### 核心类型与外部 CRD 的 Webhook

若需要为 Kubernetes 核心类型（例如 Pod）或你不拥有的外部 CRD 配置 Webhook，可参考
[controller-runtime 的内置类型示例][builtin-type-example]。Kubebuilder 对此类场景不会脚手架太多内容，但可直接使用 controller-runtime 的能力。

### 为自有 CRD 脚手架 Webhook

为 CronJob 脚手架 Webhook。示例项目使用了默认化与校验 Webhook，因此需要带上 `--defaulting` 与 `--programmatic-validation`：

```bash
kubebuilder create webhook --group batch --version v1 --kind CronJob --defaulting --programmatic-validation
```

根据需要配置 Webhook 的 CRD 数量，可能需要用不同的 GVK 重复执行以上命令。

随后为每个 Webhook 复制逻辑。对于验证型 Webhook，可将
`pkg/default_server/cronjob/validating/cronjob_create_handler.go` 中 `func validatingCronJobFn`
的内容复制到 `api/v1/cronjob_webhook.go` 的 `func ValidateCreate`（更新时对应 `ValidateUpdate`）。

同样地，把 `func mutatingCronJobFn` 的逻辑复制到 `func Default`。

### Webhook 标记（Markers）

在 v2 中脚手架 Webhook 时会添加如下标记：

```
// These are v2 markers

// This is for the mutating webhook
// +kubebuilder:webhook:path=/mutate-batch-tutorial-kubebuilder-io-v1-cronjob,mutating=true,failurePolicy=fail,groups=batch.tutorial.kubebuilder.io,resources=cronjobs,verbs=create;update,versions=v1,name=mcronjob.kb.io

...

// This is for the validating webhook
// +kubebuilder:webhook:path=/validate-batch-tutorial-kubebuilder-io-v1-cronjob,mutating=false,failurePolicy=fail,groups=batch.tutorial.kubebuilder.io,resources=cronjobs,verbs=create;update,versions=v1,name=vcronjob.kb.io
```

默认动词为 `verbs=create;update`。请根据需要调整。例如仅需在创建时校验，则改为 `verbs=create`。

同时确认 `failure-policy` 是否符合预期。

如下标记已不再需要（它们用于“自部署证书配置”，而该机制在 v2 中移除）：

```go
// v1 markers
// +kubebuilder:webhook:port=9876,cert-dir=/tmp/cert
// +kubebuilder:webhook:service=test-system:webhook-service,selector=app:webhook-server
// +kubebuilder:webhook:secret=test-system:webhook-server-secret
// +kubebuilder:webhook:mutating-webhook-config-name=test-mutating-webhook-cfg
// +kubebuilder:webhook:validating-webhook-config-name=test-validating-webhook-cfg
```

在 v1 中，同一段内可能以多个标记表示一个 Webhook；在 v2 中，每个 Webhook 必须由单一标记表示。

## 其他

若 v1 的 `main.go` 有手工改动，需要迁移到新的 `main.go`，并确保所有需要的 scheme 都已注册。

`config` 目录下新增的清单同样需要迁移。

如有需要，更新 Makefile 中的镜像名。

## 验证

最后，运行 `make` 与 `make docker-build` 确认一切正常。

[v2-project]: https://github.com/kubernetes-sigs/kubebuilder/tree/master/docs/book/src/cronjob-tutorial/testdata/project
[builtin-type-example]: https://sigs.k8s.io/controller-runtime/examples/builtins
