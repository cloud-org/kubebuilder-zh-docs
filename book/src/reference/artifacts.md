# 构建产物（Artifacts）

为了测试你的控制器，你需要包含相关二进制的压缩包（tarballs）：

```shell
./bin/k8s/
└── 1.25.0-darwin-amd64
    ├── etcd
    ├── kube-apiserver
    └── kubectl
```

这些压缩包由 [controller-tools](https://github.com/kubernetes-sigs/controller-tools) 发布，可用版本列表见：[envtest-releases.yaml](https://github.com/kubernetes-sigs/controller-tools/blob/main/envtest-releases.yaml)。

当你运行 `make envtest` 或 `make test` 时，所需压缩包会被自动下载并配置到你的项目中。

<aside class="note">
<h1>设置 ENVTEST 工具</h1>

有关 Makefile 目标 `setup-envtest` 所用工具的更多信息，请参阅 controller-runtime 仓库的 [README](https://github.com/kubernetes-sigs/controller-runtime/blob/main/tools/setup-envtest/README.md)。同时也可参考 Kubebuilder 的 [ENVTEST 文档][env-test-doc]。

</aside>


<aside class="note warning">
<h1>重要：请停止使用 https://storage.googleapis.com/kubebuilder-tools</h1>

**旧地址 [https://storage.googleapis.com/kubebuilder-tools](https://storage.googleapis.com/kubebuilder-tools) 的产物已废弃，维护者不再支持、构建或保证其发布。**

自 k8s `1.28` 起，[ENVTEST][env-test-doc] 二进制迁移至新地址：[envtest-releases.yaml](https://github.com/kubernetes-sigs/controller-tools/blob/main/envtest-releases.yaml)。在 k8s `1.29.3` 之后，用于测试控制器的二进制将不再出现在旧地址。

**新产物仅在新地址发布。**

**请确保你的项目使用新地址。**为保证下载能力，请使用 controller-runtime `v0.19.0` 的 `setup-envtest`。**对 Kubebuilder 用户而言，此更新是透明的。**

针对 k8s `1.31` 的 [ENVTEST][env-test-doc] 产物仅在：[Controller Tools Releases][controller-gen] 提供。

你也可查看 Kubebuilder 脚手架的 Makefile，可以看到 envtest 的配置与各版本 controller-runtime 保持一致。从 `release-0.19` 起会自动从正确位置下载产物，**因此 Kubebuilder 用户不会受影响。**

```shell
## Tool Binaries
..
ENVTEST ?= $(LOCALBIN)/setup-envtest
...

## Tool Versions
...
#ENVTEST_VERSION is the version of controller-runtime release branch to fetch the envtest setup script (i.e. release-0.20)
ENVTEST_VERSION ?= $(shell go list -m -f "{{ .Version }}" sigs.k8s.io/controller-runtime | awk -F'[v.]' '{printf "release-%d.%d", $$2, $$3}')
#ENVTEST_K8S_VERSION is the version of Kubernetes to use for setting up ENVTEST binaries (i.e. 1.31)
ENVTEST_K8S_VERSION ?= $(shell go list -m -f "{{ .Version }}" k8s.io/api | awk -F'[v.]' '{printf "1.%d", $$3}')
...
.PHONY: setup-envtest
setup-envtest: envtest ## Download the binaries required for ENVTEST in the local bin directory.
	@echo "Setting up envtest binaries for Kubernetes version $(ENVTEST_K8S_VERSION)..."
	@$(ENVTEST) use $(ENVTEST_K8S_VERSION) --bin-dir $(LOCALBIN) -p path || { \
		echo "Error: Failed to set up envtest binaries for version $(ENVTEST_K8S_VERSION)."; \
		exit 1; \
	}

.PHONY: envtest
envtest: $(ENVTEST) ## Download setup-envtest locally if necessary.
$(ENVTEST): $(LOCALBIN)
	$(call go-install-tool,$(ENVTEST),sigs.k8s.io/controller-runtime/tools/setup-envtest,$(ENVTEST_VERSION))
```

</aside>

[env-test-doc]: ./envtest.md
[controller-runtime]: https://github.com/kubernetes-sigs/controller-runtime
[controller-gen]: https://github.com/kubernetes-sigs/controller-tools/releases
