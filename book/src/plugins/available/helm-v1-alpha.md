# Helm 插件（`helm/v1-alpha`）

Helm 插件是一个可选插件，用于脚手架生成 Helm Chart，便于你通过 Helm 分发项目。

在默认脚手架下，用户可以先生成包含全部清单的打包文件：

```bash
make build-installer IMG=<some-registry>/<project-name:tag>
```

随后，项目使用者可以直接应用该打包文件安装：

```bash
kubectl apply -f https://raw.githubusercontent.com/<org>/project-v4/<tag or branch>/dist/install.yaml
```

不过在很多场景，你可能更希望提供 Helm Chart 的分发方式。这时就可以使用本插件在 `dist` 目录下生成 Helm Chart。

<aside class="note">
<h1>示例</h1>

你可以在 Kubebuilder 仓库根目录的 [testdata][testdata] 目录下，查看 `project-v4-with-plugins` 示例了解用法。

</aside>

## 何时使用

- 你希望向用户提供 Helm Chart 来安装和管理你的项目。
- 你需要用最新的项目变更同步更新 `dist/chart/` 下已生成的 Helm Chart：
  - 生成新清单后，使用 `edit` 子命令同步 Helm Chart。
  - 重要：如果你通过 [DeployImage][deployImage-plugin] 插件创建了 Webhook 或 API，
    需在（运行过 `make manifests` 之后）使用 `--force` 标志执行 `edit`，以基于最新清单重新生成 Helm Chart 的 values；
    若你曾定制过 `dist/chart/values.yaml` 和 `templates/manager/manager.yaml`，则在强制更新后需要手动把你的定制重新套上去。

## 如何使用

### 基本用法

Helm 插件挂载在 `edit` 子命令上，因为 `helm/v1-alpha` 依赖于先完成 Go 项目的脚手架。

```sh

# 初始化一个新项目
kubebuilder init

# 在已有项目上启用/更新 Helm Chart（先生成 config/ 下的清单）
make manifests
kubebuilder edit --plugins=helm/v1-alpha
```
<aside class="note">
  <h1>使用 edit 同步最新变更到 Helm Chart</h1>

  当项目内容变更后，先运行 `make manifests`，再执行
  `kubebuilder edit --plugins=helm/v1-alpha` 更新 Helm Chart。

  注意：除非加上 `--force`，以下文件默认不会被覆盖更新：

  <pre>
  dist/chart/
  ├── values.yaml
  └── templates/
      └── manager/
          └── manager.yaml
  </pre>

  而 `chart/Chart.yaml`、`chart/templates/_helpers.tpl` 与 `chart/.helmignore` 在初次创建后也不会自动更新（除非你删除它们）。

</aside>

## 子命令

Helm 插件实现了以下子命令：

- edit（`$ kubebuilder edit [OPTIONS]`）

## 影响的文件

本插件会创建或更新以下脚手架：

- `dist/chart/*`

[testdata]: https://github.com/kubernetes-sigs/kubebuilder/tree/master/testdata/project-v4-with-plugins
[deployImage-plugin]: ./deploy-image-plugin-v1-alpha.md
