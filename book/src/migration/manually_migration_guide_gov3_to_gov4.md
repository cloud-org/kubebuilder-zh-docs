# 通过手动更新文件从 go/v3 迁移到 go/v4

在继续之前，请先了解 [Kubebuilder go/v3 与 go/v4 的差异][v3vsv4]。

请确保已按[安装指南][quick-start]安装所需组件。

本文描述如何手动修改 PROJECT 配置以开始使用 `go/v4`。

注意：这种方式更复杂、容易出错且无法保证成功；同时你也得不到默认脚手架文件中的改进与修复。

通常仅在你对项目做了大量定制、严重偏离推荐脚手架时才建议走手动。继续前务必阅读[项目定制化][project-customizations]的提示。与其手动硬迁移，不如先收敛项目结构到推荐布局，会更有利于长期维护与升级。

推荐优先采用[从 go/v3 迁移到 go/v4][migration-guide-gov3-to-gov4]的“新建项目+迁移代码”的方式。

## 将 PROJECT 的布局从 "go/v3" 迁移到 "go/v4"

更新 `PROJECT` 文件（记录资源与插件信息，供脚手架决策）。其中 `layout` 字段指明脚手架与主插件版本。

### 迁移步骤

#### 在 PROJECT 中调整 layout 版本

以下为需要对 `PROJECT`（位于根目录）进行的手工修改。其目的在于补上 Kubebuilder 生成该文件时会写入的信息。

将：

```yaml
layout:
- go.kubebuilder.io/v3
```

替换为：

```yaml
layout:
- go.kubebuilder.io/v4

```

#### 布局变化

##### 新布局：

- 目录 `apis` 重命名为 `api`
- `controllers` 目录移至新目录 `internal` 且改为单数 `controller`
- 根目录下的 `main.go` 移至新目录 `cmd`

因此，布局会变为：

```sh
...
├── cmd
│ └── main.go
├── internal
│ └── controller
└── api
```

##### 迁移到新布局：

- 新建目录 `cmd`，将 `main.go` 移入其中
- 若项目启用 multi-group，API 原本位于 `apis`，需重命名为 `api`
- 将 `controllers` 目录移动到 `internal` 下并重命名为 `controller`
- 更新 import：
  - 修改 `main.go` 的导入路径，使其引用 `internal/controller` 下的新路径

**接着，更新脚手架相关路径**

- 更新 Dockerfile，确保包含：

```
COPY cmd/main.go cmd/main.go
COPY api/ api/
COPY internal/controller/ internal/controller/
```

然后将：

```
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -a -o manager main.go

```

替换为：

```
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -a -o manager cmd/main.go
```

- 更新 Make 目标以构建并运行 manager，将：

```
.PHONY: build
build: manifests generate fmt vet ## Build manager binary.
	go build -o bin/manager main.go

.PHONY: run
run: manifests generate fmt vet ## Run a controller from your host.
	go run ./main.go
```

替换为：

```
.PHONY: build
build: manifests generate fmt vet ## Build manager binary.
	go build -o bin/manager cmd/main.go

.PHONY: run
run: manifests generate fmt vet ## Run a controller from your host.
	go run ./cmd/main.go
```

- 更新 `internal/controller/suite_test.go` 中 `CRDDirectoryPaths` 的路径：

将：

```
CRDDirectoryPaths:     []string{filepath.Join("..", "config", "crd", "bases")},
```

替换为：

```
CRDDirectoryPaths:     []string{filepath.Join("..", "..", "config", "crd", "bases")},
```

注意：若项目为多 Group（`multigroup:true`），则需要再多加一层，即 `"..", "..", "..",`。

#### 同步更新 PROJECT 中的路径

`PROJECT` 会跟踪项目中所有 API 的路径。确认它们已指向 `api/...`，例如：

更新前：
```
  group: crew
  kind: Captain
  path: sigs.k8s.io/kubebuilder/testdata/project-v4/apis/crew/v1
```

更新后：
```

  group: crew
  kind: Captain
  path: sigs.k8s.io/kubebuilder/testdata/project-v4/api/crew/v1
```

### 用新变更更新 kustomize 清单

- 将 `config/` 下的清单与 `go/v4` 默认脚手架保持一致（可参考 `testdata/project-v4/config/`）
- 在 `config/samples` 下新增 `kustomization.yaml`，聚合该目录中的 CR 样例（参考 `testdata/project-v4/config/samples/kustomization.yaml`）

<aside class="warning">
<h1>关于 `config/` 下脚手架文件的变更</h1>

由于切换到 `go/v4`，将不再使用 Kustomize v3x。你可以对比 `testdata/project-v3/config/` 与 `testdata/project-v4/config/`。

若项目最初由 Kubebuilder 3.0.0 生成，其脚手架可能在后续 v3 发布中有非破坏性调整，或因依赖（如 [controller-runtime][controller-runtime]、[controller-tools][controller-tools]）变化而更新。

</aside>

### 若项目包含 Webhook

在 Webhook 测试文件中，将 `admissionv1beta1 "k8s.io/api/admission/v1beta1"` 替换为 `admissionv1 "k8s.io/api/admission/v1"`。

### Makefile 更新

参考对应版本的 `testdata` 示例更新 Makefile（如 `testdata/project-v4/Makefile`）。

### 依赖更新

参考对应版本 `testdata` 中的 `go.mod`（如 `testdata/project-v4/go.mod`）更新你的 `go.mod`，随后运行 `go mod tidy` 以确保依赖最新且无编译问题。

### 验证

以上步骤旨在让你的项目手动追平 `go/v4` 插件在脚手架与布局上的变更。

没有“自动验证是否正确更新 PROJECT”的办法。最佳做法是用 `go/v4` 插件新建一个同等规模的项目（例如执行 `kubebuilder init --domain tutorial.kubebuilder.io plugins=go/v4`，并生成相同的 API、Controller 与 Webhook），对比其生成的配置与手动修改后的配置。

全部更新完成后，建议至少执行：

- `make manifests`（更新 Makefile 后，用最新 controller-gen 重新生成）
- `make all`（确保能构建并完成所有操作）

[v3vsv4]: v3vsv4.md
[quick-start]: ./../quick-start.md#installation
[migration-guide-gov3-to-gov4]: migration_guide_gov3_to_gov4.md
[controller-tools]: https://github.com/kubernetes-sigs/controller-tools/releases
[controller-runtime]: https://github.com/kubernetes-sigs/controller-runtime/releases
[multi-group]: multi-group.md
