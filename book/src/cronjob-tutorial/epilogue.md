# 尾声

到这里，我们已经完成了一个功能相当完备的 CronJob 控制器实现，使用了 Kubebuilder 的大多数特性，并借助 envtest 为控制器编写了测试。

如果你还想继续深入，前往[多版本教程](/multiversion-tutorial/tutorial.md) 学习如何为项目添加新的 API 版本。

此外，你也可以自行尝试以下步骤 —— 我们很快会为它们补充教程：

- 在 `kubectl get` 中添加[额外的打印列][printer-columns]

[printer-columns]: /reference/generating-crd.md#additional-printer-columns
