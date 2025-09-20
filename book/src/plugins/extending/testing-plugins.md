# 编写 E2E 测试

可以参考 [Kubebuilder/v4/test/e2e/utils][utils-kb] 包，其中提供了功能丰富的 `TestContext`：

- [NewTestContext][new-context] 用于定义：
  - 临时项目目录；
  - 临时 controller-manager 镜像；
  - [Kubectl 执行方法][kubectl-ktc]；
  - CLI 可执行文件（`kubebuilder`、`operator-sdk` 或你扩展的 CLI）。

定义完成后，即可使用 `TestContext`：

1. 搭建测试环境：
   - 清理环境并创建临时目录，见 [Prepare][prepare-method]；
   - 安装前置 CRD，见 [InstallCertManager][cert-manager-install]、[InstallPrometheusManager][prometheus-manager-install]。

2. 校验插件行为：
   - 触发插件绑定的子命令，见 [Init][init-subcommand]、[CreateAPI][create-api-subcommand]；
   - 使用 [PluginUtil][plugin-util] 校验脚手架输出，见 [InsertCode][insert-code]、[ReplaceInFile][replace-in-file]、[UncommentCode][uncomment-code]。

3. 验证脚手架工程可工作：
   - 执行 `Makefile` 中的目标，见 [Make][make-command]；
   - 临时加载被测控制器镜像到 Kind，见 [LoadImageToKindCluster][load-image-to-kind]；
   - 使用 Kubectl 验证运行中的资源，见 [Kubectl][kubectl-ktc]。

4. 清理测试资源：
   - 卸载前置 CRD，见 [UninstallPrometheusOperManager][uninstall-prometheus-manager]；
   - 删除临时目录，见 [Destroy][destroy-method]。

参考：
- [operator-sdk e2e 测试][sdk-e2e-tests]
- [kubebuilder e2e 测试][kb-e2e-tests]

## 生成测试样例

查看由你的插件生成的样例项目内容非常直接。

例如 Kubebuilder 基于不同插件生成[样例项目][kb-samples]以验证布局。

你也可以用 `TestContext` 生成由插件脚手架的项目目录结构。用到的命令与[扩展 CLI 能力与插件][extending-cli]中类似。

以下演示使用 `go/v4` 插件创建样例项目的一般流程（其中 `kbc` 是 `TestContext` 实例）：

- 初始化一个项目：
  ```go
  By("initializing a project")
  err = kbc.Init(
      "--plugins", "go/v4",
      "--project-version", "3",
      "--domain", kbc.Domain,
      "--fetch-deps=false",
  )
  Expect(err).NotTo(HaveOccurred(), "Failed to initialize a project")
  ```

- 定义 API：
  ```go
  By("creating API definition")
  err = kbc.CreateAPI(
      "--group", kbc.Group,
      "--version", kbc.Version,
      "--kind", kbc.Kind,
      "--namespaced",
      "--resource",
      "--controller",
      "--make=false",
  )
  Expect(err).NotTo(HaveOccurred(), "Failed to create an API")
  ```

- 脚手架生成 webhook 配置：
  ```go
  By("scaffolding mutating and validating webhooks")
  err = kbc.CreateWebhook(
      "--group", kbc.Group,
      "--version", kbc.Version,
      "--kind", kbc.Kind,
      "--defaulting",
      "--programmatic-validation",
  )
  Expect(err).NotTo(HaveOccurred(), "Failed to create an webhook")
  ```

[cert-manager-install]: https://pkg.go.dev/sigs.k8s.io/kubebuilder/v4/test/e2e/utils#TestContext.InstallCertManager
[create-api-subcommand]: https://pkg.go.dev/sigs.k8s.io/kubebuilder/v4/test/e2e/utils#TestContext.CreateAPI
[destroy-method]: https://pkg.go.dev/sigs.k8s.io/kubebuilder/v4/test/e2e/utils#TestContext.Destroy
[extending-cli]: ./extending_cli_features_and_plugins.md
[init-subcommand]: https://pkg.go.dev/sigs.k8s.io/kubebuilder/v4/test/e2e/utils#TestContext.Init
[insert-code]: https://pkg.go.dev/sigs.k8s.io/kubebuilder/v4/pkg/plugin/util#InsertCode
[kb-e2e-tests]: https://github.com/kubernetes-sigs/kubebuilder/tree/book-v4/test/e2e
[kb-samples]: https://github.com/kubernetes-sigs/kubebuilder/tree/book-v4/testdata
[kubectl-ktc]: https://pkg.go.dev/sigs.k8s.io/kubebuilder/v4/test/e2e/utils#Kubectl
[load-image-to-kind]: https://pkg.go.dev/sigs.k8s.io/kubebuilder/v4/test/e2e/utils#TestContext.LoadImageToKindCluster
[make-command]: https://pkg.go.dev/sigs.k8s.io/kubebuilder/v4/test/e2e/utils#TestContext.Make
[new-context]: https://pkg.go.dev/sigs.k8s.io/kubebuilder/v4/test/e2e/utils#NewTestContext
[plugin-util]: https://pkg.go.dev/sigs.k8s.io/kubebuilder/v4/pkg/plugin/util
[prepare-method]: https://pkg.go.dev/sigs.k8s.io/kubebuilder/v4/test/e2e/utils#TestContext.Prepare
[prometheus-manager-install]: https://pkg.go.dev/sigs.k8s.io/kubebuilder/v4/test/e2e/utils#TestContext.InstallPrometheusOperManager
[replace-in-file]: https://pkg.go.dev/sigs.k8s.io/kubebuilder/v4/pkg/plugin/util#ReplaceInFile
[sdk-e2e-tests]: https://github.com/operator-framework/operator-sdk/tree/master/test/e2e/go
[uncomment-code]: https://pkg.go.dev/sigs.k8s.io/kubebuilder/v4/pkg/plugin/util#UncommentCode
[uninstall-prometheus-manager]: https://pkg.go.dev/sigs.k8s.io/kubebuilder/v4/test/e2e/utils#TestContext.UninstallPrometheusOperManager
[utils-kb]: https://github.com/kubernetes-sigs/kubebuilder/tree/book-v4/test/e2e/utils
