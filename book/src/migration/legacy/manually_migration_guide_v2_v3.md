# 通过手动更新文件从 v2 迁移到 v3

在继续之前，请先了解 [Kubebuilder v2 与 v3 的差异][migration-v2vsv3]。

请确保已按[安装指南](/quick-start.md#installation)安装所需组件。

本文描述手动升级项目配置版本并开启插件化版本所需的步骤。

注意：这种方式更复杂、容易出错且无法保证成功；同时你也得不到默认脚手架文件中的改进与修复。

通常仅在你对项目做了大量定制、严重偏离推荐脚手架时才建议走手动。继续前务必阅读[项目定制化][project-customizations]的提示。与其手动硬迁移，不如先收敛项目结构到推荐布局，会更有利于长期维护与升级。

推荐优先采用[从 v2 迁移到 v3][migration-guide-v2-to-v3]的“新建项目+迁移代码”的方式。

## 将项目配置版本从 "2" 升级到 "3"

在不同项目配置版本之间迁移，意味着需要在 `init` 命令生成的 `PROJECT` 文件中进行字段的新增、删除与修改。

`PROJECT` 文件采用了新布局，会记录更多资源信息，以便插件在脚手架时做出合理决策。

此外，`PROJECT` 文件本身也引入版本：`version` 字段表示 `PROJECT` 文件版本；`layout` 字段表示脚手架与主插件版本。

### 迁移步骤

以下为需要对 `PROJECT`（位于根目录）进行的手工修改。其目的在于补上 Kubebuilder 生成该文件时会写入的信息。

#### 新增 `projectName`

项目名为项目目录的小写名：

```yaml
...
projectName: example
...
```

#### 新增 `layout`

与旧版本等价的默认插件布局为 `go.kubebuilder.io/v2`：

```yaml
...
layout:
- go.kubebuilder.io/v2
...
```

#### 更新 `version`

`version` 表示项目布局版本，更新为 `"3"`：

```yaml
...
version: "3"
...
```

#### 补充资源信息

`resources` 属性表示项目已脚手架出来的资源清单。

为项目中的每个资源补充以下信息：

##### 添加 Kubernetes API 版本：`resources[entry].api.crdVersion: v1beta1`

```yaml
...
resources:
- api:
    ...
    crdVersion: v1beta1
  domain: my.domain
  group: webapp
  kind: Guestbook
  ...
```

##### 添加 CRD 作用域：`resources[entry].api.namespaced: true`（集群级则为 false）

```yaml
...
resources:
- api:
    ...
    namespaced: true
  group: webapp
  kind: Guestbook
  ...
```

##### 若该 API 脚手架了控制器，则添加 `resources[entry].controller: true`

```yaml
...
resources:
- api:
    ...
  controller: true
  group: webapp
  kind: Guestbook
```

##### 为资源添加域名，例如 `resources[entry].domain: testproject.org`
通常使用项目域名；若为核心类型或外部类型，规则见下方说明：

```yaml
...
resources:
- api:
    ...
  domain: testproject.org
  group: webapp
  kind: Guestbook
```

<aside class="note">
<h1>支持范围</h1>

Kubebuilder 默认仅支持核心类型与项目内脚手架出的 API；若不手工调整，无法直接处理外部类型。

  对核心类型，domain 为 `k8s.io` 或为空。

  对外部类型可留空。该能力尚未正式支持，最佳实践未定，详见 issue [#1999][issue-1999]。

</aside>

仅当核心类型在其 Kubernetes API 组的 scheme 定义中 `Domain` 不为空时，才需在项目中设置 `domain`。
例如：[apps/v1](https://github.com/kubernetes/api/blob/v0.19.7/apps/v1/register.go#L26) 中 Kind 的域为空；而 [authentication/v1](https://github.com/kubernetes/api/blob/v0.19.7/authentication/v1/register.go#L26) 的域为 `k8s.io`。

核心类型与其 domain 参考表：

| Core Type | Domain |
|----------|:-------------:|
| admission | "k8s.io" |
| admissionregistration | "k8s.io" |
| apps | empty |
| auditregistration | "k8s.io" |
| apiextensions | "k8s.io" |
| authentication | "k8s.io" |
| authorization | "k8s.io" |
| autoscaling | empty |
| batch | empty |
| certificates | "k8s.io" |
| coordination | "k8s.io" |
| core | empty |
| events | "k8s.io" |
| extensions | empty |
| imagepolicy | "k8s.io" |
| networking | "k8s.io" |
| node | "k8s.io" |
| metrics | "k8s.io" |
| policy | empty |
| rbac.authorization | "k8s.io" |
| scheduling | "k8s.io" |
| setting | "k8s.io" |
| storage | "k8s.io" |

示例：通过 `create api --group apps --version v1 --kind Deployment --controller=true --resource=false --make=false` 为核心类型 Deployment 脚手架控制器：

```yaml
- controller: true
  group: apps
  kind: Deployment
  path: k8s.io/api/apps/v1
  version: v1
```

##### 添加 `resources[entry].path`（API 的 import 路径）

<aside class="note">
<h1>Path</h1>

若未脚手架 API，而只是为已存在的（外部或核心）类型添加控制器，则可不填 path。

Kubebuilder 默认仅支持核心类型与项目内脚手架出的 API；若不手工调整，无法直接处理外部类型。

path 即 Go 代码中导入该 API 的 import 路径。

</aside>

```yaml
...
resources:
- api:
    ...
  ...
  group: webapp
  kind: Guestbook
  path: example/api/v1
```

##### 若项目使用 Webhook，则为每类 Webhook 添加 `resources[entry].webhooks.[type]: true`，并设置 `resources[entry].webhooks.webhookVersion: v1beta1`

<aside class="note">
<h1>Webhooks</h1>

可选类型：`defaulting`、`validation`、`conversion`。按项目中实际脚手架的类型填写。

`Kubebuilder v2` 脚手架 Webhook 使用的 Kubernetes API 版本为 `v1beta1`，因此统一设置 `webhookVersion: v1beta1`。

</aside>

```yaml
resources:
- api:
    ...
  ...
  group: webapp
  kind: Guestbook
  webhooks:
    defaulting: true
    validation: true
    webhookVersion: v1beta1
```

#### 检查 PROJECT 文件

确保使用 Kubebuilder v3 CLI 生成清单时，你的 `PROJECT` 文件包含一致的信息。

以 QuickStart 为例，手动升级后、使用 `go.kubebuilder.io/v2` 的 `PROJECT` 文件类似：

```yaml
domain: my.domain
layout:
- go.kubebuilder.io/v2
projectName: example
repo: example
resources:
- api:
    crdVersion: v1
    namespaced: true
  controller: true
  domain: my.domain
  group: webapp
  kind: Guestbook
  path: example/api/v1
  version: v1
version: "3"
```

你可以通过下方示例对比 `version 2` 与 `version 3` 在 `go.kubebuilder.io/v2` 布局下的差异（示例涉及多个 API 与 Webhook）：

**Example (Project version 2)**

```yaml
domain: testproject.org
repo: sigs.k8s.io/kubebuilder/example
resources:
- group: crew
  kind: Captain
  version: v1
- group: crew
  kind: FirstMate
  version: v1
- group: crew
  kind: Admiral
  version: v1
version: "2"
```

**Example (Project version 3)**

```yaml
domain: testproject.org
layout:
- go.kubebuilder.io/v2
projectName: example
repo: sigs.k8s.io/kubebuilder/example
resources:
- api:
    crdVersion: v1
    namespaced: true
  controller: true
  domain: testproject.org
  group: crew
  kind: Captain
  path: example/api/v1
  version: v1
  webhooks:
    defaulting: true
    validation: true
    webhookVersion: v1
- api:
    crdVersion: v1
    namespaced: true
  controller: true
  domain: testproject.org
  group: crew
  kind: FirstMate
  path: example/api/v1
  version: v1
  webhooks:
    conversion: true
    webhookVersion: v1
- api:
    crdVersion: v1
  controller: true
  domain: testproject.org
  group: crew
  kind: Admiral
  path: example/api/v1
  plural: admirales
  version: v1
  webhooks:
    defaulting: true
    webhookVersion: v1
version: "3"
```

### 验证

以上步骤仅更新了代表项目配置的 `PROJECT` 文件，它只对 CLI 生效，不应影响项目运行行为。

没有“自动验证是否正确更新配置”的办法。最佳做法是用相同的 API、Controller 与 Webhook 新建一个 v3 项目，对比其生成的配置与手动修改后的配置。

若上述过程有误，后续使用 CLI 时可能会遇到问题。


## 将项目切换为使用 go/v3 插件

在项目[插件][plugins-doc]之间迁移，意味着对 `init`、`create` 等插件支持的命令所创建的文件执行新增、删除与修改。
每个插件可支持一个或多个项目配置版本；请先将项目配置升级到目标插件支持的最新版本，再切换插件版本。

以下为手工修改项目布局以启用 `go/v3` 插件的步骤。注意，这无法覆盖已生成脚手架中的所有缺陷修复。

<aside class="note warning">
<h1>弃用的 API</h1>

以下步骤不会迁移这些已弃用的 API 版本：`apiextensions.k8s.io/v1beta1`、`admissionregistration.k8s.io/v1beta1`、`cert-manager.io/v1alpha2`。

</aside>

### 迁移步骤

#### 在 PROJECT 中更新插件版本

更新 `layout` 之前，请先完成项目版本升级到 `3`。随后将 `layout` 改为 `go.kubebuilder.io/v3`：

```yaml
domain: my.domain
layout:
- go.kubebuilder.io/v3
...
```

#### 升级 Go 版本与依赖

在 `go.mod` 中使用 Go `1.18`（至少满足示例版本），并对齐以下依赖版本：

```go
module example

go 1.18

require (
    github.com/onsi/ginkgo/v2 v2.1.4
    github.com/onsi/gomega v1.19.0
    k8s.io/api v0.24.0
    k8s.io/apimachinery v0.24.0
    k8s.io/client-go v0.24.0
    sigs.k8s.io/controller-runtime v0.12.1
)

```

#### Update the golang image

In the Dockerfile, replace:

```
# Build the manager binary
FROM docker.io/golang:1.13 as builder
```

With:
```
# Build the manager binary
FROM docker.io/golang:1.16 as builder
```

####  Update your Makefile

##### To allow controller-gen to scaffold the nw Kubernetes APIs

To allow `controller-gen` and the scaffolding tool to use the new API versions, replace:

```
CRD_OPTIONS ?= "crd:trivialVersions=true"
```

With:

```
CRD_OPTIONS ?= "crd"
```

##### To allow automatic downloads

To allow downloading the newer versions of the Kubernetes binaries required by Envtest into the `testbin/` directory of your project instead of the global setup, replace:

```makefile
# Run tests
test: generate fmt vet manifests
	go test ./... -coverprofile cover.out
```

With:

```makefile
# Setting SHELL to bash allows bash commands to be executed by recipes.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

ENVTEST_ASSETS_DIR=$(shell pwd)/testbin
test: manifests generate fmt vet ## Run tests.
	mkdir -p ${ENVTEST_ASSETS_DIR}
	test -f ${ENVTEST_ASSETS_DIR}/setup-envtest.sh || curl -sSLo ${ENVTEST_ASSETS_DIR}/setup-envtest.sh https://raw.githubusercontent.com/kubernetes-sigs/controller-runtime/v0.8.3/hack/setup-envtest.sh
	source ${ENVTEST_ASSETS_DIR}/setup-envtest.sh; fetch_envtest_tools $(ENVTEST_ASSETS_DIR); setup_envtest_env $(ENVTEST_ASSETS_DIR); go test ./... -coverprofile cover.out
```

<aside class="note">
<h1>Envtest binaries</h1>

The Kubernetes binaries that are required for the Envtest were upgraded from `1.16.4` to `1.22.1`.
You can still install them globally by following [these installation instructions][doc-envtest].

</aside>

##### To upgrade `controller-gen` and `kustomize` dependencies versions used

To upgrade the `controller-gen` and `kustomize` version used to generate the manifests replace:

```
# find or download controller-gen
# download controller-gen if necessary
controller-gen:
ifeq (, $(shell which controller-gen))
	@{ \
	set -e ;\
	CONTROLLER_GEN_TMP_DIR=$$(mktemp -d) ;\
	cd $$CONTROLLER_GEN_TMP_DIR ;\
	go mod init tmp ;\
	go get sigs.k8s.io/controller-tools/cmd/controller-gen@v0.2.5 ;\
	rm -rf $$CONTROLLER_GEN_TMP_DIR ;\
	}
CONTROLLER_GEN=$(GOBIN)/controller-gen
else
CONTROLLER_GEN=$(shell which controller-gen)
endif
```

With:

```
##@ Build Dependencies

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

## Tool Binaries
KUSTOMIZE ?= $(LOCALBIN)/kustomize
CONTROLLER_GEN ?= $(LOCALBIN)/controller-gen
ENVTEST ?= $(LOCALBIN)/setup-envtest

## Tool Versions
KUSTOMIZE_VERSION ?= v3.8.7
CONTROLLER_TOOLS_VERSION ?= v0.9.0

KUSTOMIZE_INSTALL_SCRIPT ?= "https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/hack/install_kustomize.sh"
.PHONY: kustomize
kustomize: $(KUSTOMIZE) ## Download kustomize locally if necessary.
$(KUSTOMIZE): $(LOCALBIN)
	test -s $(LOCALBIN)/kustomize || { curl -Ss $(KUSTOMIZE_INSTALL_SCRIPT) | bash -s -- $(subst v,,$(KUSTOMIZE_VERSION)) $(LOCALBIN); }

.PHONY: controller-gen
controller-gen: $(CONTROLLER_GEN) ## Download controller-gen locally if necessary.
$(CONTROLLER_GEN): $(LOCALBIN)
	test -s $(LOCALBIN)/controller-gen || GOBIN=$(LOCALBIN) go install sigs.k8s.io/controller-tools/cmd/controller-gen@$(CONTROLLER_TOOLS_VERSION)

.PHONY: envtest
envtest: $(ENVTEST) ## Download envtest-setup locally if necessary.
$(ENVTEST): $(LOCALBIN)
	test -s $(LOCALBIN)/setup-envtest || GOBIN=$(LOCALBIN) go install sigs.k8s.io/controller-runtime/tools/setup-envtest@latest
```

And then, to make your project use the `kustomize` version defined in the Makefile, replace all usage of `kustomize` with `$(KUSTOMIZE)`

<aside class="note">
<h1>Makefile</h1>

You can check all changes applied to the Makefile by looking in the samples projects generated in the `testdata` directory of the Kubebuilder repository or by just by creating a new project with the Kubebuilder CLI.

</aside>

#### 更新控制器

<aside class="note warning">
<h1>Controller-runtime version updated has breaking changes</h1>

Check [sigs.k8s.io/controller-runtime release docs from 0.7.0+ version][controller-releases] for breaking changes.

</aside>

Replace:

```go
func (r *<MyKind>Reconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
    ctx := context.Background()
    log := r.Log.WithValues("cronjob", req.NamespacedName)
```

With:

```go
func (r *<MyKind>Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
    log := r.Log.WithValues("cronjob", req.NamespacedName)
```

#### Update your controller and webhook test suite

<aside class="note warning">
<h1>Ginkgo V2 version update has breaking changes</h1>

见 [Ginkgo V2 迁移指南](https://onsi.github.io/ginkgo/MIGRATING_TO_V2) 获取破坏性变更详情。

</aside>

Replace:

```go
	. "github.com/onsi/ginkgo"
```

With:

```go
	. "github.com/onsi/ginkgo/v2"
```

同时调整你的测试用例：

For Controller Suite:

```go
	RunSpecsWithDefaultAndCustomReporters(t,
		"Controller Suite",
		[]Reporter{printer.NewlineReporter{}})
```

With:

```go
	RunSpecs(t, "Controller Suite")
```

For Webhook Suite:

```go
	RunSpecsWithDefaultAndCustomReporters(t,
		"Webhook Suite",
		[]Reporter{printer.NewlineReporter{}})
```

With:

```go
	RunSpecs(t, "Webhook Suite")
```

最后，从 `BeforeSuite` 中移除超时参数：

Replace:

```go
var _ = BeforeSuite(func(done Done) {
	....
}, 60)
```

With


```go
var _ = BeforeSuite(func(done Done) {
	....
})
```



#### 调整 Logger，使用 flag 选项

在 `main.go` 中将如下内容：

```go
flag.Parse()

ctrl.SetLogger(zap.New(zap.UseDevMode(true)))
```

替换为：

```go
opts := zap.Options{
	Development: true,
}
opts.BindFlags(flag.CommandLine)
flag.Parse()

ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))
```

#### 重命名 manager 参数

`--metrics-addr` 与 `enable-leader-election` 改为 `--metrics-bind-address` 与 `--leader-elect`，以与 Kubernetes 核心组件保持一致。详见 [#1839][issue-1893]。

在 `main.go` 中，将：


```go
func main() {
	var metricsAddr string
	var enableLeaderElection bool
	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
```

替换为：
```go
func main() {
	var metricsAddr string
	var enableLeaderElection bool
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
```

随后在 `config/default/manager_auth_proxy_patch.yaml` 与 `config/default/manager.yaml` 中同步重命名：

```yaml
- name: manager
args:
- "--health-probe-bind-address=:8081"
- "--metrics-bind-address=127.0.0.1:8080"
- "--leader-elect"
```

#### 验证

最后，运行 `make` 与 `make docker-build` 确认一切正常。

## 移除对已弃用 Kubernetes API 版本的使用

<aside class="note">
<h1>继续前阅读</h1>

请先了解 [CRD 版本化][custom-resource-definition-versioning]。

</aside>

以下步骤描述如何移除对已弃用 API 的使用：`apiextensions.k8s.io/v1beta1`、`admissionregistration.k8s.io/v1beta1`、`cert-manager.io/v1alpha2`。

Kubebuilder CLI 不支持“同一项目同时脚手架两代 Kubernetes API”的情况，例如既有 `apiextensions.k8s.io/v1beta1` 又有 `v1` 的 CRD。

<aside class="note">
<h1>Cert Manager API</h1>

当你使用 `admissionregistration.k8s.io/v1` 脚手架 webhook 时，默认会在清单中使用 `cert-manager.io/v1`。

</aside>

首先更新 `PROJECT` 文件，将 `api.crdVersion:v1beta` 与 `webhooks.WebhookVersion:v1beta` 改为 `api.crdVersion:v1` 与 `webhooks.WebhookVersion:v1`，例如：

```yaml
domain: my.domain
layout: go.kubebuilder.io/v3
projectName: example
repo: example
resources:
- api:
    crdVersion: v1
    namespaced: true
  group: webapp
  kind: Guestbook
  version: v1
  webhooks:
    defaulting: true
    webhookVersion: v1
version: "3"
```

你可以尝试通过 `--force` 重新生成 API（CRD）与 Webhook 的清单。

<aside class="note warning">
<h1>重建前请注意</h1>

工具会对文件进行“重新脚手架”，这意味着你会丢失其中已有的自定义内容。

执行前务必备份，或使用 `git` 对比以找回本地变更。

</aside>

随后对相同 Group/Kind/Version 使用 `kubebuilder create api` 与 `kubebuilder create webhook` 并加上 `--force`，分别重建 CRD 与 Webhook 清单。


[migration-guide-v2-to-v3]: migration_guide_v2tov3.md
[envtest]: https://book.kubebuilder.io/reference/testing/envtest.html
[controller-releases]: https://github.com/kubernetes-sigs/controller-runtime/releases
[issue-1893]: https://github.com/kubernetes-sigs/kubebuilder/issues/1839
[plugins-doc]: /reference/cli-plugins.md
[migration-v2vsv3]: v2vsv3.md
[custom-resource-definition-versioning]: https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definition-versioning/
[issue-1999]: https://github.com/kubernetes-sigs/kubebuilder/issues/1999
[project-customizations]: v2vsv3.md#project-customizations
[doc-envtest]:/reference/envtest.md
