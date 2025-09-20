# 部署 Admission Webhook

## cert-manager

你需要按照[这里](./cert-manager.md)的说明安装 cert-manager 组件。

## 构建镜像

在本地运行以下命令构建镜像：

```bash
make docker-build docker-push IMG=<some-registry>/<project-name>:tag
```

<aside class="note">
<h1> 使用 Kind </h1>

考虑把 Kind 融入到你的工作流中，以获得更快、更高效的本地开发与 CI 体验。注意：如果你使用的是 Kind 集群，则无需将镜像推送到远程镜像仓库，可以直接把本地镜像加载到指定的 Kind 集群：

```bash
kind load docker-image <your-image-name>:tag --name <your-kind-cluster-name>
```

了解更多请参见：[用于开发与 CI 的 Kind](./../reference/kind.md)

</aside>


## 部署 Webhook

你需要通过 kustomize 启用 webhook 与 cert-manager 的配置。
`config/default/kustomization.yaml` 现在应如下所示：

```yaml
{{#include ./testdata/project/config/default/kustomization.yaml}}
```

而 `config/crd/kustomization.yaml` 现在应如下所示：

```yaml
{{#include ./testdata/project/config/crd/kustomization.yaml}}
```

现在你可以将其部署到集群：

```bash
make deploy IMG=<some-registry>/<project-name>:tag
```

稍等片刻，直到 webhook Pod 启动并签发好证书。通常在 1 分钟内完成。

现在可以创建一个合法的 CronJob 来测试你的 webhooks；创建应当能够顺利完成。

```bash
kubectl create -f config/samples/batch_v1_cronjob.yaml
```

你也可以尝试创建一个非法的 CronJob（例如使用格式错误的 schedule 字段）。此时应看到创建失败并返回校验错误。

<aside class="note warning">

<h1>引导（Bootstrapping）问题</h1>

如果你在同一个集群里为 Pod 部署了 webhook，需要注意引导问题：webhook Pod 的创建请求会被发送到它自己，但此时它尚未启动。

为避免该问题，你可以在 Kubernetes 1.9+ 中使用 [namespaceSelector]，或在 1.15+ 中使用 [objectSelector]，以跳过自身。

</aside>

[namespaceSelector]: https://github.com/kubernetes/api/blob/kubernetes-1.14.5/admissionregistration/v1beta1/types.go#L189-L233
[objectSelector]: https://github.com/kubernetes/api/blob/kubernetes-1.15.2/admissionregistration/v1beta1/types.go#L262-L274
