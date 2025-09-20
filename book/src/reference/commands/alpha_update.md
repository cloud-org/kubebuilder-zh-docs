# 使用 `alpha update` 升级项目

## 概述

`kubebuilder alpha update` 通过 Git 的三方合并将你的项目脚手架升级到更新的 Kubebuilder 版本。它会为旧版本与新版本分别重建干净的脚手架，将你当前的代码合并进新的脚手架，并生成一个便于审阅的输出分支。它负责繁重工作，让你专注于审阅与解决冲突，而不是重复应用你的代码。

默认情况下，最终结果会在专用输出分支上被压缩为单个提交（squash）。若希望保留完整历史（不 squash），请使用 `--show-commits`。

<aside class="note">
<H1> 自动化该流程 </H1>

你可以使用 [AutoUpdate 插件][autoupdate-plugin] 减少保持项目最新的负担：当有新版本可用时，它会在计划任务中自动运行 `kubebuilder alpha update`。

此外，你还可以借助 [AI models][ai-gh-models] 了解保持项目最新所需的变更，并在发生冲突时获得解决方案。

</aside>

## 适用场景

在以下情况下使用该命令：

- 希望迁移到更新的 Kubebuilder 版本或插件布局
- 希望在独立分支上审阅脚手架变更
- 希望专注于解决合并冲突（而非重复应用自定义代码）

## 工作原理

你需要告知工具“目标版本”以及当前项目所在的分支。它会重建两个脚手架，并通过“三方合并”将你的代码并入新脚手架，最后给出一个可供审阅与安全合并的输出分支。你可以决定是只保留一个干净提交、保留完整历史，还是自动推送到远端。

### 第一步：检测版本
- 读取 `PROJECT` 文件或命令行参数
- 通过 `PROJECT` 文件中的 `cliVersion` 字段（若存在）确定“来源版本”
- 确定“目标版本”（默认最新发布版本）
- 选择当前代码所在分支（默认 `main`）

### 第二步：创建脚手架
命令会创建三个临时分支：
- Ancestor：来自旧版本的干净项目脚手架
- Original：你的当前代码快照
- Upgrade：来自新版本的干净脚手架

### 第三步：执行三方合并
- 使用 Git 的三方合并将 Original（你的代码）合并到 Upgrade（新脚手架）
- 在引入上游变更的同时保留你的自定义
- 如果发生冲突：
    - 默认：停止并让你手动解决
    - 使用 `--force`：即使存在冲突标记也继续提交（适合自动化）
- 运行 `make manifests generate fmt vet lint-fix` 进行整理

### 第四步：写入输出分支
- 默认情况下，所有变更会在一个安全的输出分支上被压缩为单个提交：`kubebuilder-update-from-<from-version>-to-<to-version>`
- 你可以调整行为：
    - `--show-commits`：保留完整历史
    - `--restore-path`：在 squash 模式下，从基分支恢复特定文件（例如 CI 配置）
    - `--output-branch`：自定义输出分支名
    - `--push`：自动推送结果到 `origin`
    - `--git-config`：设置 Git 配置
    - `--open-gh-issue`：创建带检查清单与对比链接的 GitHub Issue（需要 `gh`）
    - `--use-gh-models`：使用 `gh models` 向 Issue 添加 AI 概览评论

### 第五步：清理
- 输出分支就绪后，所有临时工作分支会被删除
- 你会得到一个干净的分支，可用于测试、审阅并合并回主分支

## 如何使用（命令）

Run from your project root:

```shell
kubebuilder alpha update
```

固定版本与基分支：

```shell
kubebuilder alpha update \
--from-version v4.5.2 \
--to-version   v4.6.0 \
--from-branch  main
```
适合自动化（即使发生冲突也继续）：

```shell
kubebuilder alpha update --force
```

保留完整历史而非 squash：
```
kubebuilder alpha update --from-version v4.5.0 --to-version v4.7.0 --force --show-commits
```

默认 squash，但保留基分支中的 CI/workflows：

```shell
kubebuilder alpha update --force \
--restore-path .github/workflows \
--restore-path docs
```

使用自定义输出分支名：

```shell
kubebuilder alpha update --force \
--output-branch upgrade/kb-to-v4.7.0
```

执行更新并将结果推送到 origin：

```shell
kubebuilder alpha update --from-version v4.6.0 --to-version v4.7.0 --force --push
```

## 处理冲突（`--force` 与默认行为）

使用 `--force` 时，即使存在冲突 Git 也会完成合并。提交中会包含如下冲突标记：

```shell
<<<<<<< HEAD
Your changes
=======
Incoming changes
>>>>>>> (original)
```

This allows you to run the command in CI or cron jobs without manual intervention.

- Without `--force`: the command stops on the merge branch and prints guidance; no commit is created.
- With `--force`: the merge is committed (merge or output branch) and contains the markers.

After you fix conflicts, always run:

```shell
make manifests generate fmt vet lint-fix
# or
make all
```

## Using with GitHub Issues (`--open-gh-issue`) and AI (`--use-gh-models`) assistance

Pass `--open-gh-issue` to have the command create a GitHub **Issue** in your repository
to assist with the update. Also, if you also pass `--use-gh-models`, the tool posts a follow-up comment
on that Issue with an AI-generated overview of the most important changes plus brief conflict-resolution
guidance.

### Examples

Create an Issue with a compare link:
```shell
kubebuilder alpha update --open-gh-issue
```

Create an Issue **and** add an AI summary:
```shell
kubebuilder alpha update --open-gh-issue --use-gh-models
```

### What you’ll see

The command opens an Issue that links to the diff so you can create the PR and review it, for example:

<img width="638" height="482" alt="Example Issue" src="https://github.com/user-attachments/assets/589fd16b-7709-4cd5-b169-fd53d69790d4" />

With `--use-gh-models`, an AI comment highlights key changes and suggests how to resolve any conflicts:

<img width="740" height="425" alt="Comment" src="https://github.com/user-attachments/assets/fb5f214e-be0e-43b8-a3fb-b5744ac8f66e" />

Moreover, AI models are used to help you understand what changes are needed to keep your project up to date,
and to suggest resolutions if conflicts are encountered, as in the following example:

### Automation

This integrates cleanly with automation. The [`autoupdate.kubebuilder.io/v1-alpha`][autoupdate-plugin] plugin can scaffold a GitHub Actions workflow that runs the command on a schedule (e.g., weekly). When a new Kubebuilder release is available, it opens an Issue with a compare link so you can create the PR and review it.

## Changing Extra Git configs only during the run (does not change your ~/.gitconfig)_

By default, `kubebuilder alpha update` applies safe Git configs:
`merge.renameLimit=999999`, `diff.renameLimit=999999`, `merge.conflictStyle=merge`
You can add more, or disable them.

- **Add more on top of defaults**
```shell
kubebuilder alpha update \
  --git-config rerere.enabled=true
```

- **Disable defaults entirely**
```shell
kubebuilder alpha update --git-config disable
```

- **Disable defaults and set your own**

```shell
kubebuilder alpha update \
  --git-config disable \
  --git-config rerere.enabled=true
```

<aside class="note warning">
<h1>你可能需要先升级项目</h1>

该命令底层会调用 `kubebuilder alpha generate`。我们支持由 v4.5.0+ 创建的项目；如果你的项目更老，请先运行一次 `kubebuilder alpha generate` 以现代化脚手架，此后你即可使用 `kubebuilder alpha update` 进行后续升级。

使用 Kubebuilder v4.6.0+ 创建的项目会在 `PROJECT` 文件中包含 `cliVersion`，我们会用它来选择正确的 CLI 进行再生成。

</aside>

## 参数

| Flag               | 描述                                                                                                                                                                                                                             |
|--------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `--force`          | 即使发生合并冲突也继续。将带有冲突标记的文件提交（适合 CI/定时任务）                                                                                                                       |
| `--from-branch`    | 当前项目代码所在的 Git 分支。默认 `main`                                                                                                                                                                    |
| `--from-version`   | 迁移来源的 Kubebuilder 版本（例如 `v4.6.0`）。未设置时尽可能从 `PROJECT` 读取                                                                                                                          |
| `--git-config`     | 可重复。以 `-c key=value` 形式传入 Git 配置。默认（未指定时）：`-c merge.renameLimit=999999 -c diff.renameLimit=999999`。你的配置会叠加其上。若要禁用默认值，添加 `--git-config disable` |
| `--open-gh-issue`  | 更新完成后创建带预填检查清单与对比链接的 GitHub Issue（需要 `gh`）                                                                                                                          |
| `--output-branch`  | 输出分支名称。默认：`kubebuilder-update-from-<from-version>-to-<to-version>`                                                                                                                                           |
| `--push`           | 更新完成后将输出分支推送到 `origin` 远端                                                                                                                                                               |
| `--restore-path`   | 可重复。在 squash 模式下，从基分支保留的路径（例如 `.github/workflows`）。与 `--show-commits` 不兼容                                                                                                 |
| `--show-commits`   | 保留完整历史（不 squash）。与 `--restore-path` 不兼容                                                                                                                                                            |
| `--to-version`     | 迁移目标的 Kubebuilder 版本（例如 `v4.7.0`）。未设置时默认最新可用版本                                                                                                                              |
| `--use-gh-models`  | 使用 `gh models` 作为 Issue 评论发布 AI 概览。需要 `gh` 与 `gh-models` 扩展。仅当同时设置 `--open-gh-issue` 时有效                                                                                    |
| `-h, --help`       | Show help for this command.                                                                                                                                                                                                             |

## 演示

<iframe width="560" height="315" src="https://www.youtube.com/embed/J8zonID__8k?si=WC-FXOHX0mCjph71" title="YouTube video player" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share" referrerpolicy="strict-origin-when-cross-origin" allowfullscreen></iframe>

<aside class="note">
<h1>关于此演示</h1>

该视频录制于 Kubebuilder `v7.0.1` 版本。自那以后命令已获得改进，因此当前行为可能与演示中略有不同。

</aside>

## 更多资源

- [AutoUpdate Plugin][autoupdate-plugin]
- [Design proposal for update automation][design-proposal]
- [Project configuration reference][project-config]

[project-config]: ../../reference/project-config.md
[autoupdate-plugin]: ./../../plugins/available/autoupdate-v1-alpha.md
[design-proposal]: ./../../../../../designs/update_action.md
[ai-gh-models]: https://docs.github.com/en/github-models/about-github-models
