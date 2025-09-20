# 平台支持（Platforms Supported）

Kubebuilder 生成的方案默认可在多平台或特定平台运行，取决于你对工作负载的构建与配置方式。本文指导你按需正确配置项目。

## 概览（Overview）

要支持特定或多种平台，需确保工作负载所用镜像已针对目标平台构建。注意，目标平台未必与开发环境一致，而是你的方案实际运行与发布的目标环境。建议构建多平台镜像，以便在不同操作系统与架构的集群中通用。

## 如何声明/支持目标平台

以下说明为单平台或多平台/多架构提供支持需要做的工作。

### 1）构建支持目标平台的工作负载镜像

用于 Pod/Deployment 的镜像必须支持目标平台。可用 [docker manifest inspect <image>][docker-manifest] 查看镜像的多平台 ManifestList，例如：

```shell
$ docker manifest inspect myregistry/example/myimage:v0.0.1
{
   "schemaVersion": 2,
   "mediaType": "application/vnd.docker.distribution.manifest.list.v2+json",
   "manifests": [
      {
         "mediaType": "application/vnd.docker.distribution.manifest.v2+json",
         "size": 739,
         "digest": "sha256:a274a1a2af811a1daf3fd6b48ff3d08feb757c2c3f3e98c59c7f85e550a99a32",
         "platform": {
            "architecture": "arm64",
            "os": "linux"
         }
      },
      {
         "mediaType": "application/vnd.docker.distribution.manifest.v2+json",
         "size": 739,
         "digest": "sha256:d801c41875f12ffd8211fffef2b3a3d1a301d99f149488d31f245676fa8bc5d9",
         "platform": {
            "architecture": "amd64",
            "os": "linux"
         }
      },
      {
         "mediaType": "application/vnd.docker.distribution.manifest.v2+json",
         "size": 739,
         "digest": "sha256:f4423c8667edb5372fb0eafb6ec599bae8212e75b87f67da3286f0291b4c8732",
         "platform": {
            "architecture": "s390x",
            "os": "linux"
         }
      },
      {
         "mediaType": "application/vnd.docker.distribution.manifest.v2+json",
         "size": 739,
         "digest": "sha256:621288f6573c012d7cf6642f6d9ab20dbaa35de3be6ac2c7a718257ec3aff333",
         "platform": {
            "architecture": "ppc64le",
            "os": "linux"
         }
      },
   ]
}
```

### 2）（最佳实践）配置与平台匹配的 nodeAffinity 表达式

Kubernetes 提供了 [nodeAffinity][node-affinity] 机制，用于限定 Pod 可调度到的节点集合。在多平台（异构）集群中，这对于保证正确的调度行为尤为重要。

**Kubernetes 清单示例**

```yaml
affinity:
  nodeAffinity:
    requiredDuringSchedulingIgnoredDuringExecution:
      nodeSelectorTerms:
      - matchExpressions:
        - key: kubernetes.io/arch
          operator: In
          values:
          - amd64
          - arm64
          - ppc64le
          - s390x
        - key: kubernetes.io/os
            operator: In
            values:
              - linux
```

**Golang 示例**

```go
Template: corev1.PodTemplateSpec{
    ...
    Spec: corev1.PodSpec{
        Affinity: &corev1.Affinity{
            NodeAffinity: &corev1.NodeAffinity{
                RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
                    NodeSelectorTerms: []corev1.NodeSelectorTerm{
                        {
                            MatchExpressions: []corev1.NodeSelectorRequirement{
                                {
                                    Key:      "kubernetes.io/arch",
                                    Operator: "In",
                                    Values:   []string{"amd64"},
                                },
                                {
                                    Key:      "kubernetes.io/os",
                                    Operator: "In",
                                    Values:   []string{"linux"},
                                },
                            },
                        },
                    },
                },
            },
        },
        SecurityContext: &corev1.PodSecurityContext{
            ...
        },
        Containers: []corev1.Container{{
            ...
        }},
    },
```

<aside class="note">
<h1> 示例 </h1>

可参考 Deploy Image 插件生成的示例代码（[更多](../plugins/available/deploy-image-plugin-v1-alpha.md)）。

</aside>

## 产出支持多平台的项目

可使用 [`docker buildx`][buildx] 结合仿真（[QEMU](https://www.qemu.org/)）来构建 manager 的多平台镜像。Kubebuilder 新版本脚手架默认包含 `docker-buildx` 目标。

**使用示例**

```shell
$ make docker-buildx IMG=myregistry/myoperator:v0.0.1
```

注意：需确保项目中所有镜像与工作负载均满足上述多平台支持要求，并为所有工作负载正确配置 [nodeAffinity][node-affinity]。因此请在 `config/manager/manager.yaml` 中取消注释如下示例：

```yaml
# TODO(user): Uncomment the following code to configure the nodeAffinity expression
# according to the platforms which are supported by your solution.
# It is considered best practice to support multiple architectures. You can
# build your manager image using the makefile target docker-buildx.
# affinity:
#   nodeAffinity:
#     requiredDuringSchedulingIgnoredDuringExecution:
#       nodeSelectorTerms:
#         - matchExpressions:
#           - key: kubernetes.io/arch
#             operator: In
#             values:
#               - amd64
#               - arm64
#               - ppc64le
#               - s390x
#           - key: kubernetes.io/os
#             operator: In
#             values:
#               - linux
```

<aside class="note">
<h1>面向发版构建镜像</h1>

通常建议自动化发版流程，确保镜像始终针对相同的平台构建。Goreleaser 同样支持 [docker buildx][buildx]，详见其[文档][goreleaser-buildx]。

你也可以在 GitHub Actions、Prow 或其他 CI 中配置多平台构建；或使用 `docker manifest create` 等手段达成同样目标。

默认使用 Docker 与脚手架的目标时，不需要在 Dockerfile 中显式设置 GOOS/GOARCH。若你要深度自定义构建流程，可参考 Go 的[环境变量文档](https://go.dev/doc/install/source#environment)。

</aside>

## 默认会创建哪些（工作负载）镜像？

Projects created with the Kubebuilder CLI have two workloads which are:

### Manager

运行 manager 的容器定义在 `config/manager/manager.yaml`。该镜像由脚手架生成的 Dockerfile 构建，包含本项目的二进制，默认通过 `go build -a -o manager main.go` 生成。

注意：执行 `make docker-build` 或 `make docker-build IMG=myregistry/myprojectname:<tag>` 时，会在本机构建镜像，其平台通常为 linux/amd64 或 linux/arm64。

<aside class="note">
<h1>macOS</h1>

在 macOS 环境下，Docker 也会按 linux/$arch 处理。例如在 macOS 上运行 Kind，节点最终会被打上 `kubernetes.io/os=linux` 标签。

</aside>

[node-affinity]: https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/#node-affinity
[docker-manifest]: https://docs.docker.com/engine/reference/commandline/manifest/
[buildx]: https://docs.docker.com/build/buildx/
[goreleaser-buildx]: https://goreleaser.com/customization/docker/#use-a-specific-builder-with-docker-buildx
