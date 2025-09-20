# 运行并部署控制器

### 可选
如果你选择修改了 API 定义，那么在继续之前，请生成 CR/CRD 等清单：
```bash
make manifests
```

为了测试控制器，我们可以在本地连接到集群运行它。但在此之前，需要按照[快速开始](/quick-start.md)安装我们的 CRD。必要时，这会使用 controller-tools 自动更新 YAML 清单：

```bash
make install
```

<aside class="note">

<h1>注解过长错误</h1>

如果在应用 CRD 时遇到 `metadata.annotations` 超过 262144 字节限制导致的错误，请参考[常见问题](/faq#the-error-too-long-must-have-at-most-262144-bytes-is-faced-when-i-run-make-install-to-apply-the-crd-manifests-how-to-solve-it-why-this-error-is-faced)中的对应条目。

</aside>

现在 CRD 已安装好，我们可以连接到集群运行控制器了。它会使用我们连接集群所用的凭证，因此暂时不需要担心 RBAC。

<aside class="note">

<h1>本地运行 webhooks</h1>

如果你希望在本地运行 webhooks，需要为其生成服务证书，并将其放在正确目录（默认 `/tmp/k8s-webhook-server/serving-certs/tls.{crt,key}`）。

如果你没有运行本地 API server，还需要想办法把远程集群的流量代理到本地的 webhook server。因此，我们通常建议在本地的代码-运行-测试循环中禁用 webhooks，就像下面所做的那样。

</aside>

在另一个终端中运行：

```bash
export ENABLE_WEBHOOKS=false
make run
```

你应当能看到控制器的启动日志，但此时它还不会做任何事情。

接下来我们需要一个 CronJob 来测试。把示例写到 `config/samples/batch_v1_cronjob.yaml`，然后使用它：

```yaml
{{#include ./testdata/project/config/samples/batch_v1_cronjob.yaml}}
```

```bash
kubectl create -f config/samples/batch_v1_cronjob.yaml
```

此时你应该能看到一系列活动。如果观察这些变化，应能看到 cronjob 正在运行并更新状态：

```bash
kubectl get cronjob.batch.tutorial.kubebuilder.io -o yaml
kubectl get job
```

确认它已正常工作后，我们可以将其在集群中运行。停止 `make run`，然后执行：

```bash
make docker-build docker-push IMG=<some-registry>/<project-name>:tag
make deploy IMG=<some-registry>/<project-name>:tag
```

<aside class="note">
<h1>镜像仓库权限</h1>

该镜像应发布到你指定的个人镜像仓库；你的工作环境需要具备拉取该镜像的权限。如果上述命令无法执行，请确保你对该仓库拥有正确权限。

考虑把 Kind 融入到你的工作流中，以获得更快、更高效的本地开发与 CI 体验。注意：如果使用的是 Kind 集群，无需将镜像推送到远程仓库，可以直接把本地镜像加载到指定的 Kind 集群：

```bash
kind load docker-image <your-image-name>:tag --name <your-kind-cluster-name>
```

了解更多请参见：[用于开发与 CI 的 Kind](./../reference/kind.md)

<h1>RBAC 错误</h1>

如果遇到 RBAC 错误，你可能需要为自己授予 cluster-admin 权限，或以 admin 身份登录。可参阅 [在 GKE v1.11.x 及更早版本集群上使用 Kubernetes RBAC 的前置条件][pre-rbc-gke]（这可能正是你的情形）。

</aside>

如果像之前那样再次列出 cronjobs，我们应该能看到控制器又在正常工作了！

[pre-rbc-gke]: https://cloud.google.com/kubernetes-engine/docs/how-to/role-based-access-control#iam-rolebinding-bootstrap
