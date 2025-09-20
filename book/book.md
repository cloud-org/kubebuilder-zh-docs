# Kubebuilder 文档站部署说明

本文说明如何在本地构建 `docs/book` 文档站，并将其部署到独立站点。命令与脚本保持仓库现状，不额外引入新依赖。

## 1. 环境准备

- 安装 Go ≥ 1.23（建议与仓库中 `GO_VERSION=1.23.0` 保持一致）。
- 确保 `curl`、`tar`、`unzip` 可用，脚本会下载 mdBook 与 `controller-gen`。
- 若在 CI 或容器中运行，可使用 `gimme` 预装 Go 版本，以兼容 `install-and-build.sh` 的逻辑。

## 2. 本地构建

所有命令均需在仓库根目录执行：

```bash
GO_VERSION=1.23.0 ./docs/book/install-and-build.sh
```

脚本会：

- 根据当前平台下载 mdBook 0.4.40 至临时目录。
- 安装 `sigs.k8s.io/controller-tools/cmd/controller-gen@v0.19.0` 到 `docs/book/functions`。
- 运行自定义预处理器 `litgo.sh` 与 `markerdocs.sh`。
- 输出静态站点到 `docs/book/book`。

## 3. 本地预览

构建成功后，可使用任意静态文件服务器本地预览，例如：

```bash
npx serve docs/book/book
```

或使用 Netlify CLI 复现线上行为：

```bash
netlify dev --dir docs/book/book
```

## 4. 独立部署流程

### 4.1 直接同步静态文件

1. 执行构建脚本得到 `docs/book/book` 内容。
2. 同步该目录下的静态文件至目标托管（如 GitHub Pages、S3、Cloudflare Pages、Vercel）。
3. 若平台支持自定义 Header/Redirect，可将 `docs/book/functions` 中的逻辑迁移到对应的 Edge Function 或重写规则，以保留下载链接重定向能力。

### 4.2 Netlify 单独站点

若想在新域名上继续使用 Netlify：

1. 新建 Netlify 站点，仓库指向本项目但使用独立分支或 `docs/book` 子目录。
2. 在站点设置中配置：
   - **Base directory**: `docs/book`
   - **Build command**: `GO_VERSION=1.23.0 ./install-and-build.sh`
   - **Publish directory**: `docs/book/book`
   - **Functions directory**: `docs/book/functions`
3. 绑定新域名，配置 HTTPS 即可上线。

## 5. CI/CD 建议

- 在新站点仓库或分支上添加构建流程（如 GitHub Actions），执行同一脚本并发布工件。
- 若构建频次高，建议封装一个包含 Go 与 mdBook 的轻量容器镜像，缩短脚本下载时间。
- 可以引入链接检查、Markdown lint 或 `mdbook build --dest-dir` 的 dry run 以提前发现断链或预处理器错误。

## 6. 常见问题排查

- **mdBook 下载失败**：检查是否被代理阻断，可预先下载并缓存到镜像或内部文件仓库。
- **预处理脚本报错**：确认 `litgo.sh` 与 `markerdocs.sh` 拥有执行权限，并在 Bash 环境下运行。
- **Go 依赖缺失**：确保 `GOBIN` 写入路径在 `PATH` 中，或手动导出 `export PATH=$(go env GOBIN):$PATH`。

完成上述步骤后，你即可在本地验证文档站，并将同样的构建产物部署到新的独立网站。
