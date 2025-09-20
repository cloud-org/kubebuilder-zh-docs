# 简短插曲：其他这些东西是什么？

如果你瞥了一眼 [`api/v1/`](https://sigs.k8s.io/kubebuilder/docs/book/src/cronjob-tutorial/testdata/project/api/v1) 目录下的其他文件，可能会注意到除了 `cronjob_types.go` 之外还有两个文件：`groupversion_info.go` 与 `zz_generated.deepcopy.go`。

这两个文件都不需要手动编辑（前者保持不变，后者是自动生成的），但了解它们的内容是有帮助的。

## `groupversion_info.go`

`groupversion_info.go` 包含了关于 group-version 的通用元数据：

{{#literatego ./testdata/project/api/v1/groupversion_info.go}}

## `zz_generated.deepcopy.go`

`zz_generated.deepcopy.go` 包含了前面提到的 `runtime.Object` 接口的自动生成实现，它将我们所有的根类型标记为代表某个 Kind。

`runtime.Object` 接口的核心是一个深拷贝方法 `DeepCopyObject`。

controller-tools 中的 `object` 生成器还会为每个根类型及其所有子类型生成另外两个实用方法：`DeepCopy` 与 `DeepCopyInto`。
