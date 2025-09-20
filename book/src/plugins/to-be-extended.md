## 供扩展使用

以下插件适用于其他工具及[外部插件][external-plugins]，用于扩展 Kubebuilder 的功能。

你可以使用 kustomize 插件来脚手架生成 `config/` 下的 kustomize 文件；基础语言插件负责生成 Golang 相关文件。
这样你就能为其它语言创建自己的插件（例如 [Operator-SDK][sdk] 让用户可以使用 Ansible/Helm），或是在其上叠加更多能力。

例如 [Operator-SDK][sdk] 提供了与 [OLM][olm] 集成的插件，为项目添加了其自有的能力。

| 插件 | Key | 说明 |
|---|---|---|
| [kustomize.common.kubebuilder.io/v2][kustomize-plugin] | `kustomize/v2` | 负责脚手架生成 `config/` 目录下的全部 [kustomize][kustomize] 文件 |
| `base.go.kubebuilder.io/v4` | `base/v4` | 负责脚手架生成所有 Golang 相关文件。该插件与其它插件组合后形成 `go/v4` |

[kustomize]: https://kustomize.io/
[sdk]: https://github.com/operator-framework/operator-sdk
[olm]: https://olm.operatorframework.io/
[kustomize-plugin]: ./available/kustomize-v2.md
[external-plugins]: ./extending/external-plugins.md
