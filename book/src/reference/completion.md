# 启用命令行自动补全（Shell Autocompletion）
可以通过 `kubebuilder completion [bash|fish|powershell|zsh]` 生成 Kubebuilder 的自动补全脚本。
在你的 Shell 中 source 该脚本即可启用自动补全。

<aside class="note">
<h1>Bash 前置条件</h1>

Bash 的自动补全脚本依赖 [bash-completion](https://github.com/scop/bash-completion)。请先安装该软件（可自行检查系统是否已安装）。同时确保 Bash 版本 >= 4.1。

</aside>


- 安装完成后，将 `/usr/local/bin/bash` 加入 `/etc/shells`：

    `echo "/usr/local/bin/bash" | sudo tee -a /etc/shells`

- 切换当前用户的默认 Shell：

    `chsh -s /usr/local/bin/bash`

- 在 `~/.bash_profile` 或 `~/.bashrc` 中加入：

```
# kubebuilder autocompletion
if [ -f /usr/local/share/bash-completion/bash_completion ]; then
  . /usr/local/share/bash-completion/bash_completion
fi
. <(kubebuilder completion bash)
```
- 重启终端或对上述文件执行 `source` 使其生效。

<aside class="note">
<h1>Zsh</h1>

`zsh` 的配置流程与上述类似。

</aside>

<aside class="note">
<h1>Fish</h1>

`source (kubebuilder completion fish | psub)`

</aside>
