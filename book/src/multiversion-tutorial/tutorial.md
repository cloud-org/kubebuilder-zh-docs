# 教程：多版本 API

大多数项目最初都会从一个会随版本变化的 alpha API 开始。然而，最终大多数项目都需要迁移到更稳定的 API。一旦 API 稳定，就不能再引入破坏性变更。这正是 API 版本化发挥作用的地方。

让我们对 `CronJob` 的 API 规格做一些变更，并确保我们的 CronJob 项目能支持不同的版本。

如果你还没有，请先阅读基础的 [CronJob 教程](/cronjob-tutorial/cronjob-tutorial.md)。

<aside class="note">

<h1>跟着做 vs 快速跳转</h1>

需要注意的是，本教程的大部分内容由可文学化的 Go 文件生成，它们构成了一个可运行项目，位于书籍源码目录：
[docs/book/src/multiversion-tutorial/testdata/project][tutorial-source]。

[tutorial-source]: https://github.com/kubernetes-sigs/kubebuilder/tree/master/docs/book/src/multiversion-tutorial/testdata/project

</aside>

接下来，让我们明确要做哪些变更……
