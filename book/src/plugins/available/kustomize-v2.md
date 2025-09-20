# Kustomize v2

（默认脚手架）

Kustomize 插件用于与语言基础插件 `base.go.kubebuilder.io/v4` 搭配，脚手架生成全部 Kustomize 清单。
对于通过 `go/v4`（默认脚手架）创建的项目，它会在 `config/` 目录下生成配置清单。

诸如 [Operator-sdk][sdk] 这类项目会把 Kubebuilder 当作库使用，并提供 Ansible、Helm 等其它语言的选项。
Kustomize 插件帮助它们在不同语言间保持一致的配置；同时也便于编写在默认脚手架之上做改动的插件，
避免在多种语言插件中手工同步更新。同样的思路还能让你创建可复用到不同项目与语言的“辅助”插件。

<aside class="note">
<h1>示例</h1>

你可以在 Kubebuilder 仓库的 [testdata][testdata] 目录下的 `project-v4-*` 示例中，查看 `config/` 目录下的 kustomize 内容。

</aside>

## 如何使用

如果希望你的语言插件使用 kustomize，可使用 [Bundle Plugin][bundle] 指定：由“你的语言插件 + kustomize 配置插件”组合而成，例如：

```go
import (
   ...
   kustomizecommonv2 "sigs.k8s.io/kubebuilder/v4/pkg/plugins/common/kustomize/v2"
   golangv4 "sigs.k8s.io/kubebuilder/v4/pkg/plugins/golang/v4"
   ...
)

// 为 Kubebuilder go/v4 脚手架的 Golang 项目创建组合插件
gov4Bundle, _ := plugin.NewBundle(plugin.WithName(golang.DefaultNameQualifier),
    plugin.WithVersion(plugin.Version{Number: 4}),
    plugin.WithPlugins(kustomizecommonv2.Plugin{}, golangv4.Plugin{}), // 脚手架生成 config/ 与全部 kustomize 文件
)
```

也可以单独使用 kustomize/v2：

```sh
kubebuilder init --plugins=kustomize/v2
$ ls -la
total 24
drwxr-xr-x   6 camilamacedo86  staff  192 31 Mar 09:56 .
drwxr-xr-x  11 camilamacedo86  staff  352 29 Mar 21:23 ..
-rw-------   1 camilamacedo86  staff  129 26 Mar 12:01 .dockerignore
-rw-------   1 camilamacedo86  staff  367 26 Mar 12:01 .gitignore
-rw-------   1 camilamacedo86  staff   94 31 Mar 09:56 PROJECT
drwx------   6 camilamacedo86  staff  192 31 Mar 09:56 config
```

或者与基础语言插件组合使用：

```sh
# 提供与 go/v4 相同的组合脚手架，但显式声明使用 kustomize/v2
kubebuilder init --plugins=kustomize/v2,base.go.kubebuilder.io/v4 --domain example.org --repo example.org/guestbook-operator
```

## 子命令

Kustomize 插件实现了以下子命令：

* init（`$ kubebuilder init [OPTIONS]`）
* create api（`$ kubebuilder create api [OPTIONS]`）
* create webhook（`$ kubebuilder create api [OPTIONS]`）

<aside class="note">
<h1>Create API 与 Webhook</h1>

`create api` 的实现会为每个 API 脚手架生成专属的 kustomize 清单，详见 [kustomize-create-api][kustomize-create-api]；
`create webhook` 同理。

</aside>

## 影响的文件

本插件会创建或更新以下脚手架：

* `config/*`

## 延伸阅读

* kustomize 插件[实现代码](https://github.com/kubernetes-sigs/kubebuilder/tree/master/pkg/plugins/common/kustomize)
* [kustomize 文档][kustomize-docs]
* [kustomize 仓库][kustomize-github]

[sdk]:https://github.com/operator-framework/operator-sdk
[kustomize-docs]: https://kustomize.io/
[kustomize-github]: https://github.com/kubernetes-sigs/kustomize
[kustomize-replacements]: https://kubectl.docs.kubernetes.io/references/kustomize/kustomization/replacements/
[kustomize-vars]: https://kubectl.docs.kubernetes.io/references/kustomize/kustomization/vars/
[release-notes-v5]: https://github.com/kubernetes-sigs/kustomize/releases/tag/kustomize%2Fv5.0.0
[release-notes-v4]: https://github.com/kubernetes-sigs/kustomize/releases/tag/kustomize%2Fv4.0.0
[testdata]: ./../../../../../testdata/
[bundle]: ./../../../../../pkg/plugin/bundle.go
[kustomize-create-api]: ./../../../../../pkg/plugins/common/kustomize/v2/scaffolds/api.go

