# 基础项目包含什么？

当为一个新项目搭建脚手架时，Kubebuilder 会为我们提供一些基础样板代码。

## 构建基础设施

首先，是用于构建项目的基础设施：

<details><summary><code>go.mod</code>：与项目匹配的新 Go 模块，包含基础依赖</summary>

```go
{{#include ./testdata/project/go.mod}}
```
</details>

<details><summary><code>Makefile</code>：用于构建与部署控制器的 Make 目标</summary>

```makefile
{{#include ./testdata/project/Makefile}}
```
</details>

<details><summary><code>PROJECT</code>：用于搭建新组件的 Kubebuilder 元数据</summary>

```yaml
{{#include ./testdata/project/PROJECT}}
```
</details>

## 启动配置

我们还会在 [`config/`](https://github.com/kubernetes-sigs/kubebuilder/tree/master/docs/book/src/cronjob-tutorial/testdata/project/config) 目录下获得启动配置。当前它只包含将控制器在集群中启动所需的 [Kustomize](https://sigs.k8s.io/kustomize) YAML 定义，但一旦开始编写控制器，它还会包含我们的 CustomResourceDefinition、RBAC 配置以及 WebhookConfiguration。

[`config/default`](https://github.com/kubernetes-sigs/kubebuilder/tree/master/docs/book/src/cronjob-tutorial/testdata/project/config/default) 中包含一个用于以标准配置启动控制器的 [Kustomize base](https://github.com/kubernetes-sigs/kubebuilder/blob/master/docs/book/src/cronjob-tutorial/testdata/project/config/default/kustomization.yaml)。

其它每个目录都包含不同的配置内容，并被重构为各自的 base：

- [`config/manager`](https://github.com/kubernetes-sigs/kubebuilder/tree/master/docs/book/src/cronjob-tutorial/testdata/project/config/manager)：将控制器作为 Pod 在集群中启动

- [`config/rbac`](https://github.com/kubernetes-sigs/kubebuilder/tree/master/docs/book/src/cronjob-tutorial/testdata/project/config/rbac)：在其专用 ServiceAccount 下运行控制器所需的权限

## 入口点

最后但同样重要的是，Kubebuilder 会为我们的项目搭建基本的入口点：`main.go`。接下来我们看看它……
