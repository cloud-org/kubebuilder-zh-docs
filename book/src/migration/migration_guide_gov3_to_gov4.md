# 从 go/v3 迁移到 go/v4

在继续之前，请先了解 [Kubebuilder go/v3 与 go/v4 的差异][v3vsv4]。

请确保已按[安装指南][quick-start]安装所需组件。

推荐的迁移方式是：新建一个 `go/v4` 项目，然后将 API 与调谐（reconciliation）代码拷贝过去。这样最终得到的项目就是原生的 `go/v4` 布局（最新版本）。

<aside class="note warning">
<h1>升级助手：`alpha generate` 命令</h1>

也可以使用 `kubebuilder alpha generate [OPTIONS]` 对项目进行“再脚手架（re-scaffold）”。
你可以运行 `kubebuilder alpha generate --plugins=go/v4`，基于 [PROJECT][project-file] 文件配置使用 `go/v4` 重新生成工程骨架。（详见 ./../reference/rescaffold.md）

</aside>

不过在某些场景下，也可以“就地升级”（复用 go/v3 的项目布局，手动升级 PROJECT 与脚手架）。详见[手动更新文件从 go/v3 迁移到 go/v4][manually-upgrade]。

## 初始化 go/v4 项目

<aside class="note">
<h1>项目名</h1>

以下示例中项目名使用 `migration-project`，域名使用 `tutorial.kubebuilder.io`。请按你的实际情况替换。

</aside>

新建一个以项目名命名的目录。注意该名称会用于脚手架中，默认影响 manager Pod 的名称以及其部署的 Namespace：

```bash
$ mkdir migration-project-name
$ cd migration-project-name
```

现在初始化 go/v4 项目。在进入这一步前，如果不在 `GOPATH` 内，建议先初始化一个新的 Go 模块（在 `GOPATH` 内虽然技术上非必须，但仍然推荐）：

```bash
go mod init tutorial.kubebuilder.io/migration-project
```

<aside class="note">
<h1>模块名</h1>

模块名位于项目根目录的 `go.mod`：

```
module tutorial.kubebuilder.io/migration-project
```

</aside>

随后使用 kubebuilder 完成初始化：

```bash
kubebuilder init --domain tutorial.kubebuilder.io --plugins=go/v4
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

接下来，重新脚手架 API 类型与控制器。

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

将旧项目的 `controllers/cronjob_controller.go` 迁移到新项目的 `internal/controller/cronjob_controller.go`。

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

如果 v3 的 `main.go` 有手工改动，需要将其迁移到新项目的 `main.go` 中。同时确保所需的 controller-runtime `schemes` 全部完成注册。

若 `config` 目录下存在新增清单，同步迁移它们。注意 `go/v4` 使用 Kustomize v5 而非 v4，因此你若在 `config` 中做过定制，需要确认其兼容 v5，必要时按新版本修复不兼容之处。

在 v4 中，Kustomize 的安装方式由 bash 脚本改为 `go install`。请在 Makefile 中将 `kustomize` 依赖改为：
```
.PHONY: kustomize
kustomize: $(KUSTOMIZE) ## Download kustomize locally if necessary. If wrong version is installed, it will be removed before downloading.
$(KUSTOMIZE): $(LOCALBIN)
	@if test -x $(LOCALBIN)/kustomize && ! $(LOCALBIN)/kustomize version | grep -q $(KUSTOMIZE_VERSION); then \
		echo "$(LOCALBIN)/kustomize version is not expected $(KUSTOMIZE_VERSION). Removing it before installing."; \
		rm -rf $(LOCALBIN)/kustomize; \
	fi
	test -s $(LOCALBIN)/kustomize || GOBIN=$(LOCALBIN) GO111MODULE=on go install sigs.k8s.io/kustomize/kustomize/v5@$(KUSTOMIZE_VERSION)
```

如有需要，请同步更新 Makefile 中的镜像名等配置。

## 验证

最后，运行 `make` 与 `make docker-build` 以确认一切正常。

[v3vsv4]: v3vsv4.md
[quick-start]: ./../quick-start.md#installation
[controller-tools]: https://github.com/kubernetes-sigs/controller-tools/releases
[controller-runtime]: https://github.com/kubernetes-sigs/controller-runtime/releases
[multi-group]: multi-group.md
[manually-upgrade]: manually_migration_guide_gov3_to_gov4.md
[project-file]: ../reference/project-config.md
