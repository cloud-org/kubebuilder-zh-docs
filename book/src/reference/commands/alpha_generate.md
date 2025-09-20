# 使用 `alpha generate` 重新生成项目

## 概述

`kubebuilder alpha generate` 命令会使用当前安装的 CLI 与插件版本为你的项目重新生成脚手架。

它会基于 [PROJECT][project-config] 文件中指定的配置重新生成完整脚手架。这使你能够应用 Kubebuilder 新版本引入的最新布局变更、插件特性与代码生成改进。

你可以选择在原地重新生成（覆盖现有文件），或输出到其他目录以进行差异比较与手动集成。

<aside class="note warning">
<h1>脚手架再生成会删除文件</h1>
在原地执行时，该命令会删除除 `.git` 与 `PROJECT` 外的所有文件。

在运行该命令前，请务必备份你的项目或使用版本控制。
</aside>

## 适用场景

当 Kubebuilder 引入新变更时，你可以使用 `kubebuilder alpha generate` 升级项目脚手架。这包括插件更新（例如 `go.kubebuilder.io/v3` → `go.kubebuilder.io/v4`）或 CLI 版本更新（例如 4.3.1 → 最新）。

当你想要：

- 让项目使用最新布局或插件版本
- 重新生成脚手架以包含最近的变更
- 将当前脚手架与最新版本进行比较并手动应用更新
- 创建一个干净的脚手架以审阅或测试变更

当你希望完全掌控升级流程时，请使用该命令。如果项目由较旧的 CLI 版本创建且不支持 `alpha update`，该命令也很有用。

这种方式允许你对比当前分支与上游脚手架更新（例如主分支）之间的差异，并帮助你将自定义代码覆盖到新脚手架之上。

<aside class="note tip">
<h1>需要更自动化的迁移？</h1>

如果你希望用更少的手动工作升级项目脚手架，试试 [`kubebuilder alpha update`](./alpha_update.md)。

它会使用三方合并自动保留你的代码并应用最新脚手架变更。如果 `alpha update` 尚不适用于你的项目，或你更偏好手动处理变更，则使用 `alpha generate`。

</aside>

## 如何使用

### 将当前项目升级到已安装的 CLI 版本（最新脚手架）

```sh
kubebuilder alpha generate
```

运行该命令后，项目会在原地重新生成脚手架。你可以将本地变更与主分支对比以查看更新内容，并按需将自定义代码叠加上去。

### 在新目录生成脚手架

使用 `--input-dir` 与 `--output-dir` 指定输入与输出路径。

```sh
kubebuilder alpha generate \
  --input-dir=/path/to/existing/project \
  --output-dir=/path/to/new/project
```

执行后，你可以在指定的输出目录中查看生成的脚手架。

### 参数

| Flag            | 描述                                                                 |
|------------------|-------------------------------------------------------------------------|
| `--input-dir`    | 含 `PROJECT` 文件的目录路径。默认 CWD。原地模式下会删除除 `.git` 与 `PROJECT` 外的所有文件 |
| `--output-dir`   | 输出脚手架的目录。未设置时在原地重新生成 |
| `--plugins`      | 本次生成所使用的插件键 |
| `-h, --help`     | 显示帮助 |


## 更多资源

- [工作方式演示视频](https://youtu.be/7997RIbx8kw?si=ODYMud5lLycz7osp)
- [设计提案文档](../../../../../designs/helper_to_upgrade_projects_by_rescaffolding.md)

[example]: ../../../../../testdata/project-v4-with-plugins/PROJECT
[project-config]: ../../reference/project-config.md
