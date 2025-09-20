# Alpha 命令

Kubebuilder 提供了实验性的 Alpha 命令，用于协助项目迁移、脚手架再生成等高级操作。

这些命令通过自动化或半自动化的方式，简化了以往手动且易出错的任务。

<aside class="note warning">
<h1>Alpha 命令为实验性质</h1>

Alpha 命令仍在积极开发中，未来版本可能发生变化或被移除。它们会对你的项目进行本地修改，并可能在执行过程中删除文件。

使用前务必确保你的工作已提交或备份。
</aside>

当前可用的 Alpha 命令包括：

- [`alpha generate`](./../reference/commands/alpha_generate.md) — 使用当前安装的 CLI 版本重新生成项目脚手架
- [`alpha update`](./../reference/commands/alpha_update.md) — 通过脚手架快照执行三方合并以自动化迁移流程

更多信息请查看各命令的专门文档。
