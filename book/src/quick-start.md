# 快速开始

本快速开始将涵盖：

- [创建项目](#create-a-project)
- [创建 API](#create-an-api)
- [本地运行](#test-it-out)
- [在集群中运行](#run-it-on-the-cluster)

## 前置条件

- [go](https://go.dev/dl/) 版本 v1.24.5+
- [docker](https://docs.docker.com/install/) 版本 17.03+
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) 版本 v1.11.3+
- 可访问一个 Kubernetes v1.11.3+ 集群

<aside class="note">
<h1>版本兼容性与可支持性</h1>

请确保你查看了[指南](./versions_compatibility_supportability.md)。

</aside>

## 安装

安装 [kubebuilder](https://sigs.k8s.io/kubebuilder)：

```bash
# download kubebuilder and install locally.
curl -L -o kubebuilder "https://go.kubebuilder.io/dl/latest/$(go env GOOS)/$(go env GOARCH)"
chmod +x kubebuilder && sudo mv kubebuilder /usr/local/bin/
```

<aside class="note">
<h1>使用 master 分支</h1>

你可以通过克隆仓库并运行 `make install` 来生成二进制，从而使用 master 分支。请按照[贡献指南](https://github.com/kubernetes-sigs/kubebuilder/blob/master/CONTRIBUTING.md#how-to-build-kubebuilder-locally)中“如何在本地构建 Kubebuilder”一节的步骤进行。

</aside>

<aside class="note">
<h1>启用 shell 自动补全</h1>

Kubebuilder 通过命令 `kubebuilder completion <bash|fish|powershell|zsh>` 提供自动补全支持，可减少大量输入。更多信息见[completion](./reference/completion.md) 文档。

</aside>

## 创建项目

创建一个目录，并在其中运行 init 命令来初始化新项目。示例如下：

```bash
mkdir -p ~/projects/guestbook
cd ~/projects/guestbook
kubebuilder init --domain my.domain --repo my.domain/guestbook
```

<aside class="note">
<h1>在 $GOPATH 中开发</h1>

如果你的项目在 [`GOPATH`][GOPATH-golang-docs] 中初始化，被隐式调用的 `go mod init` 会为你插入模块路径；否则必须设置 `--repo=<module path>`。

如果对模块系统不熟悉，请阅读 [Go modules 的博文][go-modules-blogpost]。

</aside>

## 创建 API

运行以下命令创建一个新的 API（group/version 为 `webapp/v1`）以及其上的新 Kind（CRD）`Guestbook`：

```bash
kubebuilder create api --group webapp --version v1 --kind Guestbook
```

<aside class="note">
<h1>交互选项</h1>

如果在 Create Resource [y/n] 与 Create Controller [y/n] 处都输入 `y`，则会创建定义 API 的 `api/v1/guestbook_types.go` 文件，以及实现该 Kind（CRD）调谐业务逻辑的 `internal/controllers/guestbook_controller.go` 文件。

</aside>

（可选）编辑 API 定义与调谐业务逻辑。更多信息见 [设计一个 API](/cronjob-tutorial/api-design.md) 与[控制器包含什么](cronjob-tutorial/controller-overview.md)。

如果编辑了 API 定义，请生成诸如自定义资源（CR）或自定义资源定义（CRD）等清单：

```bash
make manifests
```

<details><summary>点击查看示例。<tt>(api/v1/guestbook_types.go)</tt></summary>
<p>

```go
// GuestbookSpec defines the desired state of Guestbook
type GuestbookSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Quantity of instances
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=10
	Size int32 `json:"size"`

	// Name of the ConfigMap for GuestbookSpec's configuration
	// +kubebuilder:validation:MaxLength=15
	// +kubebuilder:validation:MinLength=1
	ConfigMapName string `json:"configMapName"`

	// +kubebuilder:validation:Enum=Phone;Address;Name
	Type string `json:"type,omitempty"`
}

// GuestbookStatus defines the observed state of Guestbook
type GuestbookStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// PodName of the active Guestbook node.
	Active string `json:"active"`

	// PodNames of the standby Guestbook nodes.
	Standby []string `json:"standby"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster

// Guestbook is the Schema for the guestbooks API
type Guestbook struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GuestbookSpec   `json:"spec,omitempty"`
	Status GuestbookStatus `json:"status,omitempty"`
}
```

</p>
</details>


<aside class="note">
<h1> `+kubebuilder` markers </h1>

`+kubebuilder` are [markers][markers] processed by [controller-gen][controller-gen]
to generate CRDs and RBAC. Kubebuilder also provides [scaffolding markers][scaffolding-markers]
将代码注入到已有文件并简化常见任务。示例见 `cmd/main.go`。

</aside>

## 试运行

你需要一个 Kubernetes 集群来作为运行目标。可以使用 [KinD][kind] 获取一个用于测试的本地集群，或者针对远程集群运行。

<aside class="note">
<h1>使用的上下文</h1>

你的控制器会自动使用 kubeconfig 文件中的当前上下文（即 `kubectl cluster-info` 显示的那个集群）。

</aside>

在集群中安装 CRD：

```bash
make install
```

为获得快速反馈与代码级调试，运行控制器（它会在前台运行；如需保持运行，请切换到新的终端）：

```bash
make run
```

## 安装自定义资源的实例

如果你在 Create Resource [y/n] 处输入了 `y`，则在 samples 中已为该 CRD 创建了一个 CR（如果你更改过 API 定义，请先编辑样例）：

```bash
kubectl apply -k config/samples/
```

## 在集群中运行

当控制器准备好进行打包并在其他集群中测试时：

构建并将镜像推送到 `IMG` 指定的位置：

```bash
make docker-build docker-push IMG=<some-registry>/<project-name>:tag
```

使用 `IMG` 指定的镜像将控制器部署到集群：

```bash
make deploy IMG=<some-registry>/<project-name>:tag
```

<aside class="note">
<h1>镜像仓库权限</h1>

该镜像应发布到你指定的个人镜像仓库；你的工作环境需要具备拉取该镜像的权限。如果上述命令无法执行，请确保你对该仓库拥有正确权限。

考虑将 [Kind][kind] 融入到工作流中，以获得更快、更高效的本地开发与 CI 体验。注意：如果你使用的是 [Kind][kind] 集群，则无需将镜像推送到远程仓库；你可以直接把本地镜像加载到指定的 [Kind][kind] 集群：

```bash
kind load docker-image <your-image-name>:tag --name <your-kind-cluster-name>
```

强烈建议在开发与 CI 用途中使用 [Kind][kind]。更多信息见：[用于开发与 CI 的 Kind](./reference/kind.md)

<h1>RBAC 错误</h1>

如果遇到 RBAC 错误，你可能需要为自己授予 cluster-admin 权限，或以 admin 身份登录。可参阅[在 GKE v1.11.x 及更早版本集群上使用 Kubernetes RBAC 的前置条件][pre-rbc-gke]（这可能正是你的情形）。

</aside>

## 卸载 CRD

从集群删除你的 CRD：

```bash
make uninstall
```

## 取消部署控制器

从集群中取消部署控制器：

```bash
make undeploy
```
## 使用插件

Kubebuilder 的设计基于[插件][plugins]，你可以使用[可用插件][available-plugins]为项目添加可选特性。

### 生成用于管理镜像的 API 与控制器

例如，你可以使用 [deploy-image 插件][deploy-image-v1-alpha] 生成一个用于管理容器镜像的 API 与控制器：

```bash
kubebuilder create api --group webapp --version v1alpha1 --kind Busybox --image=busybox:1.36.1 --plugins="deploy-image/v1-alpha"
```

该命令会生成：

- `api/v1alpha1/busybox_types.go` 中的 API 定义
- `internal/controllers/busybox_controller.go` 中的控制器逻辑
- `internal/controllers/busybox_controller_test.go` 中的测试脚手架（使用 [EnvTest][envtest] 进行集成式测试）

<aside class="note">
<h1> 参考与示例 </h1>

你可以参考 [DeployImage 插件][deploy-image-v1-alpha] 的代码来创建你的项目。它遵循 Kubernetes 约定与推荐最佳实践。

</aside>

### 让你的项目与生态变化保持同步

Kubebuilder 提供了 [AutoUpdate 插件][autoupdate-v1-alpha]，帮助你的项目与最新的生态变化保持一致。当有新版本发布时，该插件会打开一个包含 Pull Request 对比链接的 Issue。你可以审阅更新，并在需要时使用 [GitHub AI models][ai-gh-models] 来理解保持项目最新所需的变更。

```bash
kubebuilder edit --plugins="autoupdate/v1-alpha"
```

该命令会在 `.github/workflows/autoupdate.yml` 生成一个 GitHub workflow 文件。

## 下一步

- 继续阅读[入门指南][getting-started]（不超过 30 分钟），以打下坚实基础。
- 随后深入[CronJob 教程][cronjob-tutorial]，通过开发示例项目加深理解。
- 在设计你自己的 API 与项目之前，确保你理解 [Groups、Versions 与 Kinds，哇哦！][gkv-doc] 中关于 API 与 Group 的概念。

[pre-rbc-gke]: https://cloud.google.com/kubernetes-engine/docs/how-to/role-based-access-control#iam-rolebinding-bootstrap
[cronjob-tutorial]: https://book.kubebuilder.io/cronjob-tutorial/cronjob-tutorial.html
[GOPATH-golang-docs]: https://go.dev/doc/code.html#GOPATH
[go-modules-blogpost]: https://blog.go.dev/using-go-modules
[architecture-concept-diagram]: architecture.md
[kustomize]: https://github.com/kubernetes-sigs/kustomize
[getting-started]: getting-started.md
[plugins]: plugins/plugins.md
[available-plugins]: plugins/available-plugins.md
[envtest]: ./reference/envtest.md
[autoupdate-v1-alpha]: plugins/available/autoupdate-v1-alpha.md
[deploy-image-v1-alpha]: plugins/available/deploy-image-plugin-v1-alpha.md
[gkv-doc]: cronjob-tutorial/gvks.md
[kind]: https://sigs.k8s.io/kind
[markers]: reference/markers.md
[controller-gen]: https://sigs.k8s.io/controller-tools/cmd/controller-gen
[scaffolding-markers]: reference/markers/scaffold.md
[ai-gh-models]: https://docs.github.com/en/github-models/about-github-models
