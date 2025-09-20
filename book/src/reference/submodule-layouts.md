# 子模块布局（Sub-Module Layouts）

本节介绍如何将脚手架生成的项目调整为“API 与 Controller 各自拥有独立 `go.mod`”的布局。

子模块布局（某种意义上可视作 [Monorepo][monorepo] 的一种特例）主要用于在不引入不必要的传递依赖的前提下复用 API，以便外部项目在仅消费 API 时不会被不应暴露的依赖污染。

<aside class="note">
<h1>Using External Resources/APIs</h1>

如果你希望编写 Controller 去操作与调谐（reconcile）另一项目所拥有的类型（CRD），或是 Kubernetes 核心 API 的类型，请参考：[Using an external Resources/API](./using_an_external_resource.md)。

</aside>

## 概览（Overview）

将 API 与 Controller 拆分为不同的 `go.mod` 模块，适用于如下场景：

- 有企业版 Operator 需要复用社区版的 API；
- 有众多（可能是外部的）模块依赖该 API，需要严格限制传递依赖范围；
- 降低当该 API 被其他项目引用时所带来的传递依赖影响；
- 希望将 API 的发布生命周期与 Controller 的发布生命周期分离管理；
- 希望模块化而不想把代码拆到多个仓库。

但这也会带来一些权衡，使其不太适合作为通用默认做法或插件默认布局：

- Go 官方并不推荐单仓库内使用多个模块，[多模块布局一般不被鼓励][multi-module-repositories]；
- 你随时可以将 API 抽取到一个独立仓库，这往往更利于明确跨仓库的版本管理与发布流程；
- 至少需要一条 [replace 指令][replace-directives] 来进行本地替换：要么使用 `go.work`（这引入 2 个文件并可能需要设置环境变量，在没有 GO_WORK 的构建环境中尤为明显），要么在 `go.mod` 里使用 `replace`（每次发布前后都要手动增删）。

<aside class="note warning">
<h1>维护成本的影响</h1>

一旦偏离 Kubebuilder 标准的 `PROJECT` 配置或其插件提供的扩展布局，上游的变更可能与此处的自定义模块结构产生冲突，从而带来额外维护成本。

将代码拆成多个仓库/多个模块会产生持续成本：需要明确模块间的版本依赖、分阶段升级等。对中小型项目而言，“一个仓库 + 一个模块”往往是性价比最高的选择。

除非你非常清楚自己在做什么，否则不建议偏离推荐布局。偏离后还可能失去对某些 CLI 功能与辅助项的使用。关于基础项目布局，参见文档：[What's in a basic project?][basic-project-doc]

</aside>

## 调整你的项目（Adjusting your Project）

下面的步骤将以脚手架生成的 API 为起点，逐步改造成子模块布局。

以下示例假设你在 `GOPATH` 下创建了项目：

```shell
kubebuilder init
```

并创建了 API 与 Controller：

```shell
kubebuilder create api --group operator --version v1alpha1 --kind Sample --resource --controller --make
```

### 为 API 创建第二个模块（Creating a second module for your API）

有了基础布局后，我们来启用多模块：

1. 进入 `api/v1alpha1`
2. 执行 `go mod init` 创建新的子模块
3. 执行 `go mod tidy` 解析依赖

你的 API 模块的 `go.mod` 可能如下：

```go.mod
module YOUR_GO_PATH/test-operator/api/v1alpha1

go 1.21.0

require (
        k8s.io/apimachinery v0.28.4
        sigs.k8s.io/controller-runtime v0.16.3
)

require (
        github.com/go-logr/logr v1.2.4 // indirect
        github.com/gogo/protobuf v1.3.2 // indirect
        github.com/google/gofuzz v1.2.0 // indirect
        github.com/json-iterator/go v1.1.12 // indirect
        github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
        github.com/modern-go/reflect2 v1.0.2 // indirect
        golang.org/x/net v0.17.0 // indirect
        golang.org/x/text v0.13.0 // indirect
        gopkg.in/inf.v0 v0.9.1 // indirect
        gopkg.in/yaml.v2 v2.4.0 // indirect
        k8s.io/klog/v2 v2.100.1 // indirect
        k8s.io/utils v0.0.0-20230406110748-d93618cff8a2 // indirect
        sigs.k8s.io/json v0.0.0-20221116044647-bc3834ca7abd // indirect
        sigs.k8s.io/structured-merge-diff/v4 v4.2.3 // indirect
)
```

如上所示，它仅包含 `apimachinery` 与 `controller-runtime` 等 API 所需依赖；你在 Controller 模块声明的依赖不会被一并带入为间接依赖。

### 开发期使用 replace 指令（Using replace directives for development）

在 Operator 根目录解析主模块时，如果使用 VCS 路径，可能会遇到类似错误：

```shell
go mod tidy
go: finding module for package YOUR_GO_PATH/test-operator/api/v1alpha1
YOUR_GO_PATH/test-operator imports
	YOUR_GO_PATH/test-operator/api/v1alpha1: cannot find module providing package YOUR_GO_PATH/test-operator/api/v1alpha1: module YOUR_GO_PATH/test-operator/api/v1alpha1: git ls-remote -q origin in LOCALVCSPATH: exit status 128:
	remote: Repository not found.
	fatal: repository 'https://YOUR_GO_PATH/test-operator/' not found
```

原因在于你尚未把模块推送到 VCS，主模块在解析时不再能以包的方式直接访问 API 类型，只能从模块解析，因此会失败。

解决方法是告诉 Go 工具链将 API 模块 `replace` 成你本地路径。可选两种方式：基于 go modules，或基于 go workspaces。

#### 基于 go modules（Using go modules）

在主模块的 `go.mod` 中添加 replace：

```shell
go mod edit -require YOUR_GO_PATH/test-operator/api/v1alpha1@v0.0.0 # Only if you didn't already resolve the module
go mod edit -replace YOUR_GO_PATH/test-operator/api/v1alpha1@v0.0.0=./api/v1alpha1
go mod tidy
```

注意这里使用了占位版本 `v0.0.0`。若你的 API 模块已发布过，也可以使用真实版本，但前提是该版本已可从 VCS 获取。

<aside class="note warning">
<h1>对 Controller 发版的影响</h1>

由于主模块的 `go.mod` 中包含 replace，发布 Controller 时务必先删除它：

```shell
go mod edit -dropreplace YOUR_GO_PATH/test-operator/api/v1alpha1
go mod tidy
```

</aside>

#### 基于 go workspaces（Using go workspaces）

若使用 go workspace，则无需直接改 `go.mod`，而是依赖工作区：

在项目根目录执行 `go work init` 初始化 workspace。

随后把两个模块加入 workspace：
```shell
go work use . # This includes the main module with the controller
go work use api/v1alpha1 # This is the API submodule
go work sync
```

这样 `go run`、`go build` 等命令会遵循 workspace，从而优先使用本地解析。你可以在本地直接开发而无需先发布模块。

一般不建议把 `go.work` 提交到仓库，应在 `.gitignore` 中忽略：

```gitignore
go.work
go.work.sum
```

若发布流程中存在 `go.work`，务必设置环境变量 `GOWORK=off`（可通过 `go env GOWORK` 验证）以免影响发布。

#### 调整 Dockerfile（Adjusting the Dockerfile）

构建 Controller 镜像时，Kubebuilder 默认并不了解多模块布局。你需要手动把新的 API 模块加入依赖下载步骤：

```dockerfile
# Build the manager binary
FROM docker.io/golang:1.20 as builder
ARG TARGETOS
ARG TARGETARCH

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# Copy the Go Sub-Module manifests
COPY api/v1alpha1/go.mod api/go.mod
COPY api/v1alpha1/go.sum api/go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY cmd/main.go cmd/main.go
COPY api/ api/
COPY internal/controller/ internal/controller/

# Build
# the GOARCH has not a default value to allow the binary be built according to the host where the command
# was called. For example, if we call make docker-build in a local env which has the Apple Silicon M1 SO
# the docker BUILDPLATFORM arg will be linux/arm64 when for Apple x86 it will be linux/amd64. Therefore,
# by leaving it empty we can ensure that the container and binary shipped on it will have the same platform.
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -a -o manager cmd/main.go

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/manager .
USER 65532:65532

ENTRYPOINT ["/manager"]
```

### 创建新的 API 与 Controller 版本（Creating a new API and controller release）

由于你调整了默认布局，在发布第一个版本之前，请先了解[单仓库/多模块发布流程][multi-module-repositories]（仓库中不同子目录各有一个 `go.mod`）。

假设只有一个 API，发布流程可能如下：

```sh
git commit
git tag v1.0.0 # this is your main module release
git tag api/v1.0.0 # this is your api release
go mod edit -require YOUR_GO_PATH/test-operator/api@v1.0.0 # now we depend on the api module in the main module
go mod edit -dropreplace YOUR_GO_PATH/test-operator/api/v1alpha1 # this will drop the replace directive for local development in case you use go modules, meaning the sources from the VCS will be used instead of the ones in your monorepo checked out locally.
git push origin main v1.0.0 api/v1.0.0
```

完成后，模块即可从 VCS 获取，本地开发无需再保留 `replace`。若后续继续在本地迭代，请相应地恢复 `replace` 以便本地联调。

### 复用已抽出的 API 模块（Reusing your extracted API module）

当你希望在另一个 kubebuilder 项目中复用该 API 模块时，请参考：[Using an external Type](./using_an_external_resource.md)。
在“Edit the API files”那一步，引入依赖即可：

```shell
go get YOUR_GO_PATH/test-operator/api@v1.0.0
```

随后按指南继续使用。

[basic-project-doc]: ./../cronjob-tutorial/basic-project.md
[monorepo]: https://en.wikipedia.org/wiki/Monorepo
[replace-directives]: https://go.dev/ref/mod#go-mod-file-replace
[multi-module-repositories]: https://github.com/golang/go/wiki/Modules#faqs--multi-module-repositories
