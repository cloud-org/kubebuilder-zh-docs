# 部署与测试

在测试转换之前，我们需要在 CRD 中启用它：

Kubebuilder 会在 `config` 目录下生成 Kubernetes 清单，默认禁用 webhook 相关内容。要启用它们，我们需要：

- 在 `config/crd/kustomization.yaml` 文件中启用 `patches/webhook_in_<kind>.yaml` 与 `patches/cainjection_in_<kind>.yaml`

- 在 `config/default/kustomization.yaml` 的 `bases` 段落下启用 `../certmanager` 与 `../webhook` 目录

- 在 `config/default/kustomization.yaml` 文件中启用 `CERTMANAGER` 段落下的全部变量

此外，如果 Makefile 中存在 `CRD_OPTIONS` 变量，我们需要将其设置为仅 `"crd"`，去掉 `trivialVersions` 选项（这确保我们确实[为每个版本生成校验][ref-multiver]，而不是告诉 Kubernetes 它们相同）：

```makefile
CRD_OPTIONS ?= "crd"
```

现在代码修改与清单都已就位，让我们把它部署到集群并进行测试。

除非你有其他证书管理方案，否则你需要安装 [cert-manager](../cronjob-tutorial/cert-manager.md)（版本 `0.9.0+`）。Kubebuilder 团队已经用版本 [0.9.0-alpha.0](https://github.com/cert-manager/cert-manager/releases/tag/v0.9.0-alpha.0) 验证过本教程中的步骤。

当证书相关内容准备就绪后，我们可以像平常一样运行 `make install deploy`，将所有组件（CRD、controller-manager 部署）部署到集群。

## 测试

当所有组件在集群上运行且已启用转换后，我们可以通过请求不同版本来测试转换。

我们基于 v1 版本创建一个 v2 版本（放在 `config/samples` 下）

```yaml
{{#include ./testdata/project/config/samples/batch_v2_cronjob.yaml}}
```

然后在集群中创建它：

```shell
kubectl apply -f config/samples/batch_v2_cronjob.yaml
```

如果一切正确，应能创建成功，并且我们应当能使用 v2 资源来获取它：

```shell
kubectl get cronjobs.v2.batch.tutorial.kubebuilder.io -o yaml
```

```yaml
{{#include ./testdata/project/config/samples/batch_v2_cronjob.yaml}}
```

以及 v1 资源：

```shell
kubectl get cronjobs.v1.batch.tutorial.kubebuilder.io -o yaml
```
```yaml
{{#include ./testdata/project/config/samples/batch_v1_cronjob.yaml}}
```

两者都应被正确填充，并分别与我们的 v2 与 v1 示例等价。注意它们的 API 版本不同。

最后，稍等片刻，你会注意到即便我们的控制器是基于 v1 API 版本编写的，CronJob 依然会持续进行调谐。

<aside class="note">

<h1>kubectl 与首选版本</h1>

在 Go 代码中访问 API 类型时，我们会通过该版本的 Go 类型（例如 `batchv2.CronJob`）来请求特定版本。

你可能注意到，上面调用 kubectl 的方式看起来与平时有些不同——它们指定的是“组-版本-资源（group-version-resource）”，而不仅仅是资源。

当我们执行 `kubectl get cronjob` 时，kubectl 需要确定这对应哪个组-版本-资源。为此，它会使用“发现 API（discovery API）”来确定 `cronjob` 资源的首选版本。对于 CRD，这通常是最新的稳定版本（具体细节见 [CRD 文档][CRD-version-pref]）。

随着我们对 CronJob 的更新，这意味着 `kubectl get cronjob` 将获取 `batch/v2` 组-版本。

如果希望指定精确版本，可以像上面那样使用 `kubectl get resource.version.group`。

在脚本中你应始终使用“完全限定的组-版本-资源”语法。`kubectl get resource` 是为人类、自我意识的机器人以及其他能自行识别新版本的智能体准备的；`kubectl get resource.version.group` 才是为其他一切准备的。

</aside>

## 故障排查

[排查步骤](/TODO.md)

[ref-multiver]: /reference/generating-crd.md#multiple-versions "Generating CRDs: Multiple Versions"

[crd-version-pref]: https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definition-versioning/#version-priority "Versions in CustomResourceDefinitions"
