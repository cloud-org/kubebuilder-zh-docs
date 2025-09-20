# 教程：构建 CronJob

太多教程一上来就是生硬的场景或玩具应用，只能讲清基础，然后在复杂内容面前戛然而止。本教程不同：我们会用 Kubebuilder 贯穿（几乎）全谱系的复杂度，从简单开始，逐步构建到一个功能相当完备的示例。

让我们假设（是的，这有点点“设定”）我们已经厌倦了 Kubernetes 中非 Kubebuilder 实现的 CronJob 控制器的维护负担，想用 Kubebuilder 重新实现它。

CronJob 控制器（双关非本意）的工作是在 Kubernetes 集群上以固定间隔运行一次性任务。它是构建在 Job 控制器之上的，而 Job 控制器的任务是将一次性任务运行一次并确保完成。

我们不会顺带重写 Job 控制器，而是把这当作一个学习机会，看看如何与外部类型进行交互。

<aside class="note">

<h1>跟着做 vs 快速跳转</h1>

需要注意的是，本教程的大部分内容是由位于书籍源码目录中的“可文学化”的 Go 文件生成的：
[docs/book/src/cronjob-tutorial/testdata][tutorial-source]。完整且可运行的项目位于 [project][tutorial-project-source]，而中间产物则直接位于 [testdata][tutorial-source] 目录下。

[tutorial-source]: https://github.com/kubernetes-sigs/kubebuilder/tree/master/docs/book/src/cronjob-tutorial/testdata

[tutorial-project-source]: https://github.com/kubernetes-sigs/kubebuilder/tree/master/docs/book/src/cronjob-tutorial/testdata/project

</aside>

## 为项目搭建脚手架

如在[快速开始](../quick-start.md)中所述，我们需要为新项目搭建脚手架。请先确认你已经[安装了 Kubebuilder](../quick-start.md#installation)，然后为新项目创建脚手架：

```bash
# 创建项目目录，然后执行 init 命令
mkdir project
cd project
# 我们使用 tutorial.kubebuilder.io 作为域名，
# 因此所有 API 组都将是 <group>.tutorial.kubebuilder.io。
kubebuilder init --domain tutorial.kubebuilder.io --repo tutorial.kubebuilder.io/project
```

<aside class="note">

项目名称默认取当前工作目录名。你可以通过 `--project-name=<dns1123-label-string>` 指定不同的项目名。

</aside>

现在我们已经就位，让我们看看 Kubebuilder 到目前为止为我们搭了些什么……

<aside class="note">

<h1>在 <code>$GOPATH</code> 中开发</h1>

如果项目是在 [`GOPATH`][GOPATH-golang-docs] 中初始化的，被隐式调用的 `go mod init` 会为你插入模块路径。否则必须设置 `--repo=<module path>`。

如果对模块系统不熟悉，可以阅读 [Go modules 的博文][go-modules-blogpost]。

</aside>

[GOPATH-golang-docs]: https://golang.org/doc/code.html#GOPATH
[go-modules-blogpost]: https://blog.golang.org/using-go-modules
