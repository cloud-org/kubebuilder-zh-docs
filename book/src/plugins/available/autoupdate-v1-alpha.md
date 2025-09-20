# AutoUpdate（`autoupdate/v1-alpha`）

让你的 Kubebuilder 项目与最新改进保持同步不应该是件苦差事。经过少量设置，每当有新的 Kubebuilder 版本可用时，你都可以收到自动 Pull Request 建议——让项目保持维护良好、安全，并与生态变化保持一致。

该自动化使用带“三方合并策略”的 [`kubebuilder alpha update`][alpha-update-command] 命令来刷新项目脚手架，并通过一个 GitHub Actions 工作流包装：它会打开一个带 Pull Request 对比链接的 Issue，方便你创建 PR 并进行审阅。

<aside class="warning">
<h1>保护你的分支</h1>

该工作流默认仅会把合并结果创建并推送到名为 `kubebuilder-update-from-<from-version>-to-<to-version>` 的分支。

为保证代码库安全，请使用分支保护规则以确保变更不会在未经适当审查的情况下被推送或合并。

</aside>

## 何时使用

- 当你的项目没有过多偏离默认脚手架（请务必阅读此处的自定义注意事项：https://book.kubebuilder.io/versions_compatibility_supportability#project-customizations）
- 当你希望降低保持项目更新与良好维护的负担
- 当你希望借助 AI 的指引，了解保持项目最新所需的变更并解决冲突

## 如何使用

- 为现有项目添加 `autoupdate` 插件：

```shell
kubebuilder edit --plugins="autoupdate.kubebuilder.io/v1-alpha"
```

- 创建启用 `autoupdate` 插件的新项目：

```shell
kubebuilder init --plugins=go/v4,autoupdate/v1-alpha
```

## 工作原理

该操作会生成一个运行 [kubebuilder alpha update][alpha-update-command] 命令的 GitHub Actions 工作流。每当有新版本发布时，工作流都会自动打开一个带 PR 对比链接的 Issue，方便你创建 PR 并进行审阅，例如：

<img width="638" height="482" alt="Example Issue" src="https://github.com/user-attachments/assets/589fd16b-7709-4cd5-b169-fd53d69790d4" />

默认情况下，生成的工作流会使用 `--use-gh-models` 参数以利用 [AI models][ai-models] 帮助你理解所需变更。你会获得一份简洁的变更文件列表，以加快审阅，例如：

<img width="582" height="646" alt="Screenshot 2025-08-26 at 13 40 53" src="https://github.com/user-attachments/assets/d460a5af-5ca4-4dd5-afb8-7330dd6de148" />

如发生冲突，AI 生成的评论会指出并提供后续步骤，例如：

<img width="600" height="188" alt="Conflicts" src="https://github.com/user-attachments/assets/2142887a-730c-499a-94df-c717f09ab600" />

### 工作流细节

该工作流每周检查一次新版本；如有新版本，将创建带 PR 对比链接的 Issue，以便你创建 PR 并审阅。工作流调用的命令如下：

```shell
	# 更多信息参见：https://kubebuilder.io/reference/commands/alpha_update
    - name: Run kubebuilder alpha update
      run: |
		# 使用指定参数执行更新命令。
		# --force：即使出现冲突也完成合并，保留冲突标记
		# --push：将结果分支自动推送到 'origin'
		# --restore-path：在 squash 时保留指定路径（例如 CI 工作流文件）
		# --open-gh-issue：创建 Issue
		# --use-gh-models：在创建的 Issue 中添加 AI 生成的评论，给出脚手架变更概览及（如有）冲突解决指引
        kubebuilder alpha update \
          --force \
          --push \
          --restore-path .github/workflows \
          --open-gh-issue \
          --use-gh-models
```

[alpha-update-command]: ./../../reference/commands/alpha_update.md
[ai-models]: https://docs.github.com/en/github-models/about-github-models
