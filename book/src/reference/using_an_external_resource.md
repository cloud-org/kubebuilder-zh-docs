# 使用外部资源（External Resources）

在某些场景下，你的项目需要处理并非由自身 API 定义的资源。这些外部资源主要分为两类：

- **Core Types（核心类型）**：由 Kubernetes 本身定义的 API 类型，如 `Pod`、`Service`、`Deployment` 等。
- **External Types（外部类型）**：由其他项目定义的 API 类型，例如其他方案所定义的 CRD。

## 管理 External Types

### 为 External Type 创建控制器

在不脚手架资源定义的前提下，你可以为外部类型创建控制器：使用 `create api` 并带上 `--resource=false`，同时通过 `--external-api-path` 与 `--external-api-domain` 指定外部 API 类型所在路径与域名。这样会为项目外的类型（例如由其他 Operator 管理的 CRD）生成控制器。

命令示例：

```shell
kubebuilder create api --group <theirgroup> --version <theirversion> --kind <theirKind> --controller --resource=false --external-api-path=<their Golang path import> --external-api-domain=<theirdomain>
```

- `--external-api-path`：外部类型的 Go import 路径。
- `--external-api-domain`：外部类型的 domain。该值用于生成 RBAC 时构造完整的 API 组名（如 `apiGroups: <group>.<domain>`）。

例如，若需要管理 Cert Manager 的 Certificates：

```shell
kubebuilder create api --group certmanager --version v1 --kind Certificate --controller=true --resource=false --external-api-path=github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1 --external-api-domain=io
```

由此生成的 RBAC [标记][markers-rbac]：

```go
// +kubebuilder:rbac:groups=cert-manager.io,resources=certificates,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cert-manager.io,resources=certificates/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=cert-manager.io,resources=certificates/finalizers,verbs=update
```

对应的 RBAC 角色：

```ymal
- apiGroups:
  - cert-manager.io
  resources:
  - certificates
  verbs:
  - create
  - delete
  - get
  - list
  - patch
  - update
  - watch
- apiGroups:
  - cert-manager.io
  resources:
  - certificates/finalizers
  verbs:
  - update
- apiGroups:
  - cert-manager.io
  resources:
  - certificates/status
  verbs:
  - get
  - patch
  - update
```

这会为外部类型生成控制器，但不会生成资源定义（因为该类型定义在外部项目）。

### 为 External Type 创建 Webhook

示例：

```shell
kubebuilder create webhook --group certmanager --version v1 --kind Issuer --defaulting --programmatic-validation --external-api-path=github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1 --external-api-domain=cert-manager.io
```

## 管理 Core Types

Kubernetes 的核心 API 类型（如 `Pod`、`Service`、`Deployment`）由系统预先定义。要在不脚手架资源定义的情况下为这些核心类型创建控制器，请参考下表中的组名，并指定版本与 Kind。

| Group                    | K8s API Group            |
|---------------------------|------------------------------------|
| admission                 | k8s.io/admission                  |
| admissionregistration      | k8s.io/admissionregistration      |
| apps                      | apps                              |
| auditregistration          | k8s.io/auditregistration          |
| apiextensions              | k8s.io/apiextensions              |
| authentication             | k8s.io/authentication             |
| authorization              | k8s.io/authorization              |
| autoscaling                | autoscaling                       |
| batch                     | batch                             |
| certificates               | k8s.io/certificates               |
| coordination               | k8s.io/coordination               |
| core                      | core                              |
| events                    | k8s.io/events                     |
| extensions                | extensions                        |
| imagepolicy               | k8s.io/imagepolicy                |
| networking                | k8s.io/networking                 |
| node                      | k8s.io/node                       |
| metrics                   | k8s.io/metrics                    |
| policy                    | policy                            |
| rbac.authorization        | k8s.io/rbac.authorization         |
| scheduling                | k8s.io/scheduling                 |
| setting                   | k8s.io/setting                    |
| storage                   | k8s.io/storage                    |

为 `Pod` 创建控制器的命令示例：

```shell
kubebuilder create api --group core --version v1 --kind Pod --controller=true --resource=false
```

为 `Deployment` 创建控制器：

```sh
create api --group apps --version v1 --kind Deployment --controller=true --resource=false
```

由此生成的 RBAC [标记][markers-rbac]：

```go
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps,resources=deployments/finalizers,verbs=update
```

对应的 RBAC 角色：

```yaml
- apiGroups:
  - apps
  resources:
  - deployments
  verbs:
  - create
  - delete
  - get
  - list
  - patch
  - update
  - watch
- apiGroups:
  - apps
  resources:
  - deployments/finalizers
  verbs:
  - update
- apiGroups:
  - apps
  resources:
  - deployments/status
  verbs:
  - get
  - patch
  - update
```

这会为核心类型（如 `corev1.Pod`）生成控制器，但不会生成资源定义（该类型已由 Kubernetes API 定义）。

### 为 Core Type 创建 Webhook

与创建控制器类似，使用核心类型的信息来创建 Webhook。示例：

```go
kubebuilder create webhook --group core --version v1 --kind Pod --programmatic-validation
```
[markers-rbac]: ./markers/rbac.md
