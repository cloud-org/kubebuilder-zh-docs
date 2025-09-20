# Scaffold（脚手架）

`+kubebuilder:scaffold` 标记是 Kubebuilder 脚手架体系中的关键部分。它标注了在生成文件中的插入位置：当你脚手架出新的资源（例如控制器、Webhook 或 API）时，Kubebuilder 会在这些位置注入相应代码。
借助该机制，Kubebuilder 能够将新组件无缝集成到项目中，同时不影响用户自定义的代码。

<aside class="note warning">
<H1>若删除或修改 `+kubebuilder:scaffold` 标记</H1>

Kubebuilder CLI 会在预期的文件中查找这些标记以完成代码生成。如果你移动或删除了标记，CLI 将无法插入必要代码，从而导致脚手架流程失败或行为异常。

</aside>

## 工作原理（How It Works）

当你使用 Kubebuilder CLI（如 `kubebuilder create api`）来生成新资源时，CLI 会在关键位置寻找 `+kubebuilder:scaffold` 标记，并把它们当作占位点来插入必要的 import 或注册代码。

## `main.go` 中的示例（Example Usage in `main.go`）

以下展示了 `+kubebuilder:scaffold` 在典型 `main.go` 文件中的用法。为便于说明，假设执行：

```shell
kubebuilder create api --group crew --version v1 --kind Admiral --controller=true --resource=true
```

### 添加新的导入（Imports）

`+kubebuilder:scaffold:imports` 标记允许 Kubebuilder CLI 注入额外的 import（例如新控制器或 Webhook 所需的包）。当我们创建新 API 时，CLI 会在此位置自动添加所需的导入路径。

以单组布局中新建 `Admiral` API 为例，CLI 会在 import 段落中添加 `crewv1 "<repo-path>/api/v1"`：

```go
import (
    "crypto/tls"
    "flag"
    "os"

    // Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
    // to ensure that exec-entrypoint and run can make use of them.
    _ "k8s.io/client-go/plugin/pkg/client/auth"
    ...
    crewv1 "sigs.k8s.io/kubebuilder/testdata/project-v4/api/v1"
    // +kubebuilder:scaffold:imports
)
```

### 注册新的 Scheme（Register a New Scheme）

`+kubebuilder:scaffold:scheme` 标记用于将新创建的 API 版本注册到 runtime scheme，确保这些类型能被 manager 识别。

例如，在创建 Admiral API 之后，CLI 会在 `init()` 函数中注入如下代码：


```go
func init() {
    ...
    utilruntime.Must(crewv1.AddToScheme(scheme))
    // +kubebuilder:scaffold:scheme
}
```

## 设置控制器（Set Up a Controller）

当我们创建新的控制器（如 Admiral）时，Kubebuilder CLI 会借助 `+kubebuilder:scaffold:builder` 标记将控制器的初始化代码注入到 manager。这一标记指示了新控制器的注册位置。

例如，在创建 `AdmiralReconciler` 后，CLI 会添加如下代码将控制器注册到 manager：

```go
if err = (&crewv1.AdmiralReconciler{
    Client: mgr.GetClient(),
    Scheme: mgr.GetScheme(),
}).SetupWithManager(mgr); err != nil {
    setupLog.Error(err, "unable to create controller", "controller", "Admiral")
    os.Exit(1)
}
// +kubebuilder:scaffold:builder
```

`+kubebuilder:scaffold:builder` 标记确保新生成的控制器能正确注册至 manager，从而开始对资源进行调谐。

## `+kubebuilder:scaffold` 标记列表

| 标记                                      | 常见位置                      | 作用                                                                 |
|-------------------------------------------|-------------------------------|----------------------------------------------------------------------|
| `+kubebuilder:scaffold:imports`           | `main.go`                     | 指示在此处为新控制器/Webhook/API 注入 import。                         |
| `+kubebuilder:scaffold:scheme`            | `main.go` 的 `init()`         | 向 runtime scheme 注册 API 版本。                                      |
| `+kubebuilder:scaffold:builder`           | `main.go`                     | 指示在此处向 manager 注册新控制器。                                    |
| `+kubebuilder:scaffold:webhook`           | Webhook 测试相关文件           | 指示在此处添加 Webhook 的初始化函数。                                   |
| `+kubebuilder:scaffold:crdkustomizeresource`| `config/crd`                 | 指示在此处添加 CRD 自定义资源补丁。                                     |
| `+kubebuilder:scaffold:crdkustomizewebhookpatch` | `config/crd`              | 指示在此处添加 CRD Webhook 补丁。                                       |
| `+kubebuilder:scaffold:crdkustomizecainjectionns`| `config/default`          | 指示在此处添加转换 Webhook 的 CA 注入补丁（命名空间）。                   |
| `+kubebuilder:scaffold:crdkustomizecainjectioname`| `config/default`         | 指示在此处添加转换 Webhook 的 CA 注入补丁（名称）。                       |
| （不再支持）`+kubebuilder:scaffold:crdkustomizecainjectionpatch` | `config/crd` | 旧的 Webhook CA 注入补丁位置；现已由上面两个标记替代。                   |
| `+kubebuilder:scaffold:manifestskustomizesamples` | `config/samples`          | 指示在此处注入 Kustomize 示例清单。                                     |
| `+kubebuilder:scaffold:e2e-webhooks-checks` | `test/e2e`                  | 基于已生成的 Webhook 类型添加相应的 e2e 校验。                           |

<aside class="note warning">
<h1>（不再支持）`+kubebuilder:scaffold:crdkustomizecainjectionpatch`</h1>

如果在你的代码中发现该标记，请按以下步骤处理：

1. **从 `config/crd/kustomization.yaml` 中移除 CERTMANAGER 段落：**

   删除 `CERTMANAGER` 段落，避免为 CRD 生成非预期的 CA 注入补丁。确保移除或注释以下内容：

   ```yaml
   # [CERTMANAGER] To enable cert-manager, uncomment all the sections with [CERTMANAGER] prefix.
   # patches here are for enabling the CA injection for each CRD
   #- path: patches/cainjection_in_firstmates.yaml
   # +kubebuilder:scaffold:crdkustomizecainjectionpatch
   ```

2. **在 `config/default/kustomization.yaml` 中确保 CA 注入配置：**

   在 `config/default/kustomization.yaml` 的 `[CERTMANAGER]` 区域，加入以下内容以正确生成 CA 注入：

   **注意：** 必须确保包含以下目标标记：
   - `+kubebuilder:scaffold:crdkustomizecainjectionns`
   - `+kubebuilder:scaffold:crdkustomizecainjectioname`

   ```yaml
   # - source: # Uncomment the following block if you have a ConversionWebhook (--conversion)
   #     kind: Certificate
   #     group: cert-manager.io
   #     version: v1
   #     name: serving-cert # This name should match the one in certificate.yaml
   #     fieldPath: .metadata.namespace # Namespace of the certificate CR
   #   targets: # Do not remove or uncomment the following scaffold marker; required to generate code for target CRD.
   # +kubebuilder:scaffold:crdkustomizecainjectionns
   # - source:
   #     kind: Certificate
   #     group: cert-manager.io
   #     version: v1
   #     name: serving-cert # This name should match the one in certificate.yaml
   #     fieldPath: .metadata.name
   #   targets: # Do not remove or uncomment the following scaffold marker; required to generate code for target CRD.
   # +kubebuilder:scaffold:crdkustomizecainjectioname
   ```

3. **确保 `config/crd/patches` 中仅包含“转换 Webhook”的补丁：**

   `config/crd/patches` 目录及其在 `config/crd/kustomization.yaml` 中的条目应仅包含“转换 Webhook”的补丁。此前曾因缺陷导致为任意 Webhook 生成补丁；应仅保留使用 `--conversion` 选项脚手架出的 Webhook 补丁。

更多指导可参考仓库中的 `testdata/` 示例。

> 备选方案：你也可以使用 [`alpha generate`](./../rescaffold.md) 直接基于最新版本重新生成整个项目；随后仅将你的业务代码增量叠加，以确保获得所有最新修复与改进。

</aside>

<aside class="note">
<h1>自定义标记</h1>

如果你将 Kubebuilder 作为库来[开发自己的插件](./../../plugins/creating-plugins.md)并扩展 CLI 能力，你也可以自定义并使用专属标记。可参考 [kubebuilder/v4/pkg/machinery](https://pkg.go.dev/sigs.k8s.io/kubebuilder/v4/pkg/machinery) 中的工具来创建并管理这些标记。

</aside>


