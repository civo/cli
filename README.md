# Civo CLI (2020+)

## The plan

The plan is very simple, to write a new version of the Ruby Civo CLI in Golang and deploy it using a simple `curl https://get... | sh` type script or Homebrew for Macs.

We want it to stay interface compatible where possible, with all the aliases currently configured - but want to enhance it with things like custom formatting of output in to JSON or in to custom string formats (e.g. `-o "Hostname|Size"`).

## External libraries

Golang has no shortage of external libraries for various parts of this, but the ones currently planned to be used are:

### Cobra CLI Library

* https://ordina-jworks.github.io/development/2018/10/20/make-your-own-cli-with-golang-and-cobra.html
* https://www.bradcypert.com/testing-a-cobra-cli-in-go/
* https://www.linode.com/docs/development/go/using-cobra/

### Other libraries

* https://github.com/briandowns/spinner
* go get github.com/fatih/color
* https://github.com/olekukonko/tablewriter
* https://github.com/spf13/viper

## Progress

- ✅ ~~Makefile for cross-platform builds~~
- ✅ ~~API Key management~~
- ✅ ~~Regions~~
- ✅ ~~Quotas~~
- ✅ ~~Sizes~~
- ✅ ~~Instances~~
- ✅ ~~Domain names~~
- ✅ ~~Domain records~~
- ✅ ~~Firewalls~~
- ✅ ~~Load balancers~~
- ✅ ~~SSH keys~~
- ✅ ~~Networks~~
- ✅ ~~Snapshots~~
- ✅ ~~Volumes~~
- ✅ ~~Templates~~
- ✅ ~~Kubernetes Clusters~~
- ✅ ~~Kubernetes Applications~~

- `curl | bash` installation mechanism
- Homebrew

## Enabling shell autocompletion

The civo completion script for Bash can be generated with the command civo completion bash. Sourcing the completion script in your shell enables civo autocompletion.

However, the completion script depends on bash-completion, which means that you have to install this software first (you can test if you have bash-completion already installed by running `type _init_completion`).

## Install bash-completion

bash-completion is provided by many package managers (see [here](https://github.com/scop/bash-completion#installation)). You can install it with `apt-get install bash-completion` or `yum install bash-completion`, etc.

The above commands create `/usr/share/bash-completion/bash_completion`, which is the main script of bash-completion. Depending on your package manager, you have to manually source this file in your `~/.bashrc` file.

To find out, reload your shell and run `type _init_completion`. If the command succeeds, you're already set, otherwise add the following to your `~/.bashrc` file:

```shell
source /usr/share/bash-completion/bash_completion
```

Reload your shell and verify that bash-completion is correctly installed by typing `type _init_completion`.

### Enable civo autocompletion

You now need to ensure that the civo completion script gets sourced in all your shell sessions. There are two ways in which you can do this:

- Source the completion script in your `~/.bashrc` file:

    ```shell
    echo 'source <(civo completion bash)' >>~/.bashrc
    ```
- Add the completion script to the `/etc/bash_completion.d` directory:

    ```shell
    civo completion bash >/etc/bash_completion.d/civo
    ```
If you have an alias for civo, you can extend shell completion to work with that alias:

```shell
echo 'alias c=civo' >>~/.bashrc
echo 'complete -F __start_civo c' >>~/.bashrc
```

## Install zsh-completion

The civo completion script for Zsh can be generated with the command `civo completion zsh`. Sourcing the completion script in your shell enables civo autocompletion.

To do so in all your shell sessions, add the following to your `~/.zshrc` file:

```shell
source <(civo completion zsh)
```

If you have an alias for civo, you can extend shell completion to work with that alias:

```shell
echo 'alias c=civo' >>~/.zshrc
echo 'complete -F __start_civo c' >>~/.zshrc
```
    
After reloading your shell, civo autocompletion should be working.

If you get an error like `complete:13: command not found: compdef`, then add the following to the beginning of your `~/.zshrc` file:

```shell
autoload -Uz compinit
compinit
```

To set the civo completion code for zsh to autoload on startup yo can run this command.
```bash
civo completion zsh > "${fpath[1]}/_civo"
```