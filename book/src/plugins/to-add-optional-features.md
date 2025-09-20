## 增加可选特性

以下插件可用于生成代码并利用可选特性：

| 插件 | Key | 说明 |
|---|---|---|
| [autoupdate.kubebuilder.io/v1-alpha][autoupdate] | `autoupdate/v1-alpha` | 可选辅助插件，脚手架生成一个定时任务，帮助你的项目自动跟进生态变更，显著降低人工维护成本。 |
| [deploy-image.go.kubebuilder.io/v1-alpha][deploy] | `deploy-image/v1-alpha` | 可选辅助插件，可脚手架 API 与 Controller，并内置代码实现以部署并管理一个镜像（Operand）。 |
| [grafana.kubebuilder.io/v1-alpha][grafana] | `grafana/v1-alpha` | 可选辅助插件，可为 controller-runtime 导出的默认指标脚手架生成 Grafana Dashboard 清单。 |
| [helm.kubebuilder.io/v1-alpha][helm] | `helm/v1-alpha` | 可选辅助插件，可在 `dist` 目录下脚手架生成 Helm Chart 用于项目分发。 |

[grafana]: ./available/grafana-v1-alpha.md
[deploy]: ./available/deploy-image-plugin-v1-alpha.md
[helm]: ./available/helm-v1-alpha.md
[autoupdate]: ./available/autoupdate-v1-alpha.md
