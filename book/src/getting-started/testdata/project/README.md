# 项目
// TODO(user): Add simple overview of use/purpose

## 描述（Description）
// TODO(user): An in-depth paragraph about your project and overview of use

## 入门（Getting Started）

### 先决条件（Prerequisites）
- go version v1.24.0+
- docker version 17.03+.
- kubectl version v1.11.3+.
- 可访问一个 Kubernetes v1.11.3+ 集群。

### 集群部署（Deploy on the cluster）
**构建并推送镜像到由 `IMG` 指定的位置：**

```sh
make docker-build docker-push IMG=<some-registry>/project:tag
```

注意：镜像需要发布到你指定的镜像仓库；在执行环境中必须有拉取该镜像的权限。若上述命令失败，请确认你对该镜像仓库具备相应权限。

**在集群中安装 CRD：**

```sh
make install
```

**使用 `IMG` 指定的镜像将 Manager 部署到集群：**

```sh
make deploy IMG=<some-registry>/project:tag
```

> 注意：若遇到 RBAC 报错，你可能需要为自己授予 cluster-admin 权限，或以管理员身份登录。

**创建示例实例**
你可以应用 `config/samples` 下的示例：

```sh
kubectl apply -k config/samples/
```

> 注意：请确认示例包含可用于快速验证的默认值。

### 卸载（Uninstall）
**从集群中删除实例（CR）：**

```sh
kubectl delete -k config/samples/
```

**从集群删除 API（CRD）：**

```sh
make uninstall
```

**从集群中卸载控制器：**

```sh
make undeploy
```

## 项目分发（Project Distribution）

下面是向用户发布该项目的几种方式。

### 提供包含全部 YAML 的安装包

1. 为已构建并发布到镜像仓库的镜像构建安装包：

```sh
make build-installer IMG=<some-registry>/project:tag
```

注意：上述 Makefile 目标会在 `dist` 目录生成 `install.yaml`，其中包含通过 Kustomize 构建的全部资源，足以在不包含依赖项的前提下安装本项目。

2. 使用安装包

用户可以直接通过 `kubectl apply -f <YAML 安装包 URL>` 安装项目，例如：

```sh
kubectl apply -f https://raw.githubusercontent.com/<org>/project/<tag or branch>/dist/install.yaml
```

### 提供 Helm Chart

1. 使用可选的 helm 插件构建 Chart

```sh
kubebuilder edit --plugins=helm/v1-alpha
```

2. `dist/chart` 下会生成 Chart，用户可以直接使用该目录中的内容。

注意：若你更改了项目，需要再次运行上述命令以同步最新变更到 Chart。并且，如果你创建了 Webhook，需要为该命令添加 `--force` 标志，并在随后手动把此前对 `dist/chart/values.yaml` 或 `dist/chart/manager/manager.yaml` 的自定义配置重新应用上去。

## 贡献（Contributing）
// TODO(user): Add detailed information on how you would like others to contribute to this project

注意：运行 `make help` 可查看所有可用的 Make 目标说明。

更多信息参见 [Kubebuilder 文档](https://book.kubebuilder.io/introduction.html)

## 许可证（License）

Copyright 2025 The Kubernetes authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
