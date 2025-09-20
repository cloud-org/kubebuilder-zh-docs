# 在本地开发与 CI 中使用 Kind

## 为什么用 Kind

- **搭建迅速：** 本地启动多节点集群通常不到 1 分钟。
- **销毁快捷：** 数秒内即可销毁集群，提升迭代效率。
- **直接使用本地镜像：** 无需推送到远端仓库即可部署。
- **轻量高效：** 适合本地开发与 CI/CD 场景。

这里仅覆盖使用 kind 集群的基础内容。更多细节请参阅
[kind 官方文档](https://kind.sigs.k8s.io/)。

## 安装（Installation）

按照[安装指南](https://kind.sigs.k8s.io/#installation-and-usage)安装 `kind`。

## 创建集群（Create a Cluster）

创建最简单的 kind 集群：

```bash
kind create cluster
```

如需自定义集群，可提供额外配置。下面是一个示例 `kind` 配置：

```yaml
{{#include ./kind-config.yaml}}
```

使用上述配置，执行以下命令将创建一个包含 1 个控制面与 3 个工作节点的 k8s v1.17.2 集群：

```bash
kind create cluster --config hack/kind-config.yaml --image=kindest/node:v1.17.2
```

可以通过 `--image` 指定目标集群版本，例如 `--image=kindest/node:v1.17.2`。支持的版本见
[镜像标签列表](https://hub.docker.com/r/kindest/node/tags)。

## 向集群加载本地镜像（Load Docker Image）

本地开发时，可将镜像直接加载到 kind 集群，无需使用镜像仓库：

```bash
kind load docker-image your-image-name:your-tag
```

更多信息见：[将本地镜像加载到 kind 集群](https://kind.sigs.k8s.io/docs/user/quick-start/#loading-an-image-into-your-cluster)。

## 删除集群（Delete a Cluster）

```bash
kind delete cluster
```
