# 编写控制器测试

测试 Kubernetes 控制器是一个很大的主题，而 kubebuilder 为你生成的测试样板相对精简。

为了带你了解 Kubebuilder 生成的控制器的集成测试模式，我们将重温第一篇教程中构建的 CronJob，并为其编写一个简单测试。

基本方法是：在生成的 `suite_test.go` 文件中，使用 envtest 创建一个本地 Kubernetes API server，实例化并运行你的控制器；随后编写额外的 `*_test.go` 文件，使用 [Ginkgo](http://onsi.github.io/ginkgo) 对其进行测试。

如果你想调整 envtest 集群的配置，请参见[为集成测试配置 envtest](../reference/envtest.md) 章节以及 [`envtest` 文档](https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/envtest?tab=doc)。

## 测试环境准备

{{#literatego ../cronjob-tutorial/testdata/project/internal/controller/suite_test.go}}

## 测试控制器行为

{{#literatego ../cronjob-tutorial/testdata/project/internal/controller/cronjob_controller_test.go}}

上面的 Status 更新示例展示了一个针对带下游对象的自定义 Kind 的通用测试策略。到这里，你应该已经掌握了以下用于测试控制器行为的方法：

- 在 envtest 集群上运行你的控制器
- 为创建测试对象编写桩代码（stubs）
- 只改变对象的某些部分，以测试特定的控制器行为

<aside class="note">
<h1>示例</h1>

你可以使用 [DeployImage](../plugins/available/deploy-image-plugin-v1-alpha.md) 插件来查看示例。该插件允许用户按照指南和最佳实践，为在集群上部署和管理 Operand（镜像）脚手架 API/Controller。它抽象了实现这一目标的复杂性，同时允许用户自定义生成的代码。

因此，你会看到为该控制器生成了一个使用 ENV TEST 的测试，其目标是确保 Deployment 被成功创建。你可以在 `testdata` 目录下查看其代码实现示例，配合 [DeployImage](../plugins/available/deploy-image-plugin-v1-alpha.md) 的样例，[点此](https://github.com/kubernetes-sigs/kubebuilder/blob/master/testdata/project-v4-with-plugins/internal/controller/busybox_controller_test.go)。

</aside>
