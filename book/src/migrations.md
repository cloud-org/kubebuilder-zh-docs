# 迁移

将你的项目脚手架升级以采用 Kubebuilder 的最新变更，可能涉及迁移到新的插件版本（例如 `go.kubebuilder.io/v3` → `go.kubebuilder.io/v4`）或更新的 CLI 工具链。该过程通常包括重新生成脚手架并手动合并你的自定义代码。

本节详细说明了在不同版本的 Kubebuilder 脚手架之间，以及迁移到更复杂的项目布局结构时所需的步骤。

手动方式容易出错，因此 Kubebuilder 引入了新的 alpha 命令来帮助简化迁移过程。

## 手动迁移

传统流程包括：

- 使用最新的 Kubebuilder 版本或插件重新生成项目脚手架
- 手动重新添加自定义逻辑
- 运行项目生成器：

  ```bash
  make generate
  make manifests
  ```

## 了解 PROJECT 文件（自 `v3.0.0` 引入）

Kubebuilder 使用的所有输入都记录在 [PROJECT][project-config] 文件中。如果你使用 CLI 生成脚手架，该文件会记录项目的配置与元数据。

<aside class="note warning">
<h1>项目自定义</h1>

通过 CLI 创建项目后，你可以按需进行自定义。但请注意，除非非常确定，否则不建议偏离建议的项目布局。

例如，不要随意移动脚手架文件，否则将来会很难升级项目。你也可能失去使用某些 CLI 功能与辅助工具的能力。有关项目布局的更多信息，请参阅文档：[基础项目包含什么？][basic-project-doc]

</aside>

## Alpha 迁移命令

Kubebuilder 提供了 alpha 命令来辅助项目升级。

<aside class="note warning">
<h1>自动化过程会删除文件以重新生成</h1>
会删除除 `.git` 与 `PROJECT` 之外的所有文件。
</aside>

### `kubebuilder alpha generate`

使用已安装的 CLI 版本重新生成项目脚手架。

```bash
kubebuilder alpha generate
```

### `kubebuilder alpha update`（自 `v4.7.0` 起可用）

通过执行三方合并来自动化迁移：

- 原始脚手架
- 你当前的自定义版本
- 最新或指定目标脚手架

```bash
kubebuilder alpha update
```

更多详情请参阅[Alpha 命令文档](./reference/alpha_commands.md)。


[project-config]: ./reference/project-config.md
[basic-project-doc]: ./cronjob-tutorial/basic-project.md
