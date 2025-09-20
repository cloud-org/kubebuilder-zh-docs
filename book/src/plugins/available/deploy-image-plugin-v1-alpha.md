# Deploy Image 插件（deploy-image/v1-alpha）

`deploy-image` 插件允许用户创建用于在集群中部署与管理容器镜像的[控制器][controller-runtime]与自定义资源，遵循 Kubernetes 最佳实践。它简化了部署镜像的复杂性，同时允许用户按需自定义项目。

使用该插件，你将获得：

- 一个在集群中部署与管理 Operand（镜像）的控制器实现
- 使用 [ENVTEST][envtest] 的测试，以验证调谐逻辑
- 已填充必要规格的自定义资源样例
- 在 manager 中用于管理 Operand（镜像）的环境变量支持

<aside class="note">
<h1>示例</h1>

在 Kubebuilder 项目的 [testdata][testdata] 目录下的 `project-v4-with-plugins` 中可以查看使用该插件生成的脚手架示例。

`Memcached` API 及其控制器是通过以下命令生成的：

```shell
kubebuilder create api \
  --group example.com \
  --version v1alpha1 \
  --kind Memcached \
  --image=memcached:memcached:1.6.26-alpine3.19 \
  --image-container-command="memcached,--memory-limit=64,-o,modern,-v" \
  --image-container-port="11211" \
  --run-as-user="1001" \
  --plugins="deploy-image/v1-alpha"
```

`Busybox` API 则通过以下命令创建：

```shell
kubebuilder create api \
  --group example.com \
  --version v1alpha1 \
  --kind Busybox \
  --image=busybox:1.36.1 \
  --plugins="deploy-image/v1-alpha"
```
</aside>


## 何时使用？

- 该插件非常适合刚开始接触 Kubernetes Operator 的用户
- 它帮助用户使用[Operator 模式][operator-pattern] 来部署并管理镜像（Operand）
- 如果你在寻找一种快速高效的方式来搭建自定义控制器并管理容器镜像，该插件是上佳选择

## 如何使用？

1. 初始化项目：
   使用 `kubebuilder init` 创建新项目后，你可以使用该插件创建 API。在继续之前，请确保已完成[快速开始][quick-start]。

2. 创建 API：
   使用该插件，你可以[创建 API][create-apis] 以指定要在集群上部署的镜像（Operand）。你还可以通过参数选择性地指定命令、端口与安全上下文：

   示例命令：
   ```sh
   kubebuilder create api --group example.com --version v1alpha1 --kind Memcached --image=memcached:1.6.15-alpine --image-container-command="memcached,--memory-limit=64,modern,-v" --image-container-port="11211" --run-as-user="1001" --plugins="deploy-image/v1-alpha"
   ```

<aside class="warning">
<h1>关于 make run 的说明：</h1>

当本地使用 `make run` 运行项目时，提供的 Operand 镜像会以环境变量的形式写入 `config/manager/manager.yaml`。

请在本地运行项目前导出该环境变量，例如：

```shell
export MEMCACHED_IMAGE="memcached:1.4.36-alpine"
```

</aside>

## 子命令

`deploy-image` 插件包含以下子命令：

- `create api`：使用该命令为管理容器镜像生成 API 与控制器代码

## 受影响的文件

当使用该插件的 `create api` 命令时，除了 Kubebuilder 现有的脚手架外，以下文件会受到影响：

- `controllers/*_controller_test.go`：为控制器生成测试
- `controllers/*_suite_test.go`：生成或更新测试套件
- `api/<version>/*_types.go`：生成 API 规格
- `config/samples/*_.yaml`：为自定义资源生成默认值
- `main.go`：更新以添加控制器初始化
- `config/manager/manager.yaml`：更新以包含用于存储镜像的环境变量

## 更多资源

- 查看此[视频][video]了解其工作方式

[video]: https://youtu.be/UwPuRjjnMjY
[operator-pattern]: https://kubernetes.io/docs/concepts/extend-kubernetes/operator/
[controller-runtime]: https://github.com/kubernetes-sigs/controller-runtime
[testdata]: https://github.com/kubernetes-sigs/kubebuilder/tree/master/testdata/project-v4-with-plugins
[envtest]: ./../../reference/envtest.md
[quick-start]: ./../../quick-start.md
[create-apis]: ../../cronjob-tutorial/new-api.md
