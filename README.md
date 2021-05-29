# Civo Command-Line Client

## Introduction

Civo CLI is a tool to manage your [Civo.com](https://www.civo.com) account from the terminal. The [Civo web control panel](https://www.civo.com/account/) has a user-friendly interface for managing your account, but in case you want to automate or run scripts on your account, or have multiple complex services, the command-line interface outlined here will be useful. This guide will cover the set-up and usage of the Civo CLI tool with examples.

**STATUS:** This project is currently under active development and maintenance.

## Table of contents

- [Introduction](#introduction)
- [Global Options](#global-options)
- [Set-Up](#set-up)
- [Docker Usage](#docker-usage)
- [API Keys](#api-keys)
- [Instances](#instances)
- [Kubernetes clusters](#kubernetes-clusters)
- [Kubernetes applications](#kubernetes-applications)
- [Domains and Domain Records](#domains-and-domain-records)
- [Firewalls](#firewalls)
- [Networks](#networks)
- [Load Balancers](#load-balancers)
- [Quota](#quota)
- [Sizes](#sizes)
- [Snapshots](#snapshots)
- [SSH Keys](#ssh-keys)
- [Templates](#templates)
- [Volumes](#volumes)
- [Region](#region)
- [Enabling shell autocompletion](#enabling-shell-autocompletion)
- [Contributing](#contributing)
- [License](#license)

## Set-up

Civo CLI is built with Go and distributed as binary files, available for multiple operating systems and downloadable from https://github.com/civo/cli/releases.

### Installing on macOS

If you have a Mac, you can install it using [Homebrew](https://brew.sh):

```sh
brew tap civo/tools
brew install civo
```

or if you prefer you can run this in the console:

```sh
$ curl -sL https://civo.com/get | sh
```

### Installing on Windows

Civo CLI is available to download on windows via Chocolatey and Scoop

For installing via Chocolatey you need [Chocolatey](https://chocolatey.org/install) package manager installed on your PC.
- run `choco install civo-cli` and it will install Civo CLI on your PC.

For installing via Scoop you need [Scoop](https://scoop.sh/) installed as a package manager, then:
- add the extras bucket with `scoop bucket add extras`
- install civo with `scoop install civo`.

You will also, of course, need a Civo account, for which you can [register here](https://www.civo.com/signup).


### Installing on Linux

For Linux Civo CLI can be installed by various methods.

* Install via brew, as shows in above step.

* Install via wget, specifying the release version you want: 

```
wget https://github.com/civo/cli/releases/download/v0.7.6/civo-0.7.6-linux-amd64.tar.gz
tar -xvf civo-0.7.6-linux-amd64.tar.gz
chmod +x civo
mv ./civo /usr/local/bin/
```

* You can also build the binary, but make sure you have go installed,


```sh
git clone https://github.com/civo/cli.git
cd cli
make
cd ..
cp -r cli ./$HOME
export PATH="$HOME/cli:$PATH"
```
With this, we have installed the Civo CLI successfully. Check it is working by running any of the following commands.

**Note:** For the first time when you are running, make sure you set your current region. Check [Region](#region) for more information.


### Running the Civo CLI tool and getting help

To use the tool, simply run `civo` with your chosen options. You can find context-sensitive help for commands and their options by invoking the `help` or `-h` command:
`civo help`,
`civo instance help`,
`civo instance create help`
and so on. The main components of Civo CLI are outlined in the following sections.

## Docker Usage

Civo's CLI utility can also run within a Docker container, if you prefer to keep your base OS clean.

To run, you generally will want to map the API key for persistence.

```sh
touch $HOME/.civo.json
docker run -it --rm -v $HOME/.civo.json:/.civo.json civo/cli:latest
```

You can also use the Kubernetes options of the CLI. Kubectl is included inside our image, to use it you just need to mount the configuration in the container.

```sh
touch $HOME/.civo.json
mkdir $HOME/.kube/
touch $HOME/.kube/config
docker run -it --rm -v $HOME/.civo.json:/.civo.json -v $HOME/.kube/config:$HOME/.kube/config civo/cli:latest
```

To make usage easier, an alias is recommended.  Here's an example how to set one to the same command as would be used if installed directly on the system, and using the Docker image:

Ubuntu etc:
```sh
alias civo="docker run -it --rm -v $HOME/.civo.json:/.civo.json civo/cli:latest"
# Maybe put the above line in ~/.bash_profile or ~/.zshrc
civo sshkey list
civo instance list
civo instance create --size g2.xsmall
civo k8s list
```

For Fedora users:
```sh
alias civo="docker run -it --rm -v $HOME/.civo.json:/.civo.json:Z -v $HOME/.kube/config:$HOME/.kube/config:Z civo/cli:latest"
```

Here's an example how to set an alias and get started with Kubernetes.

```sh
alias civo="docker run -it --rm -v $HOME/.civo.json:/.civo.json -v $HOME/.kube/config:$HOME/.kube/config civo/cli:latest"
# Maybe put the above line in ~/.bash_profile or ~/.zshrc
civo sshkey list
civo instance list
civo instance create --size g2.xsmall
civo k8s list
```

## Global Options

The civo cli have mutiple global options, that you can use, like this:

```
      --config string   config file (default is $HOME/.civo.json)
  -f, --fields string   output fields for custom format output (use -h to determine fields)
  -h, --help            help for civo
  -o, --output string   output format (json/human/custom) (default "human")
      --pretty          Print pretty the json output
      --region string   Choose the region to connect to, if you use this option it will use it over the default region
  -y, --yes             Automatic yes to prompts; assume "yes" as answer to all prompts and run non-interactivel
```

## API Keys

#### Introduction

In order to use the command-line tool, you will need to authenticate yourself to the Civo API using a special key. You can find an automatically-generated API key or regenerate a new key at [https://www.civo.com/api](https://www.civo.com/api).

#### Adding a current API Key to your account

You can add the API Key to the CLI tool through the API Keys command.
`civo apikey add apikey_name apikey` such as:

```sh
$ civo apikey add Demo_Test_Key DAb75oyqVeaE7BI6Aa74FaRSP0E2tMZXkDWLC9wNQdcpGfH51r
  Saved the API Key Demo_Test_Key as DAb75oyqVeaE7BI6Aa74FaRSP0E2tMZXkDWLC9wNQdcpGfH51r
```

As you can have multiple API keys stored to handle multiple accounts, you will need to tell which key the tool should use to authenticate with `civo apikey current [apikey_name]`. This sets your chosen API key as the default key to use for any subsequent commands:

```sh
$ civo apikey current Demo_Test_Key
  Set the default API Key to be Demo_Test_Key
```

#### Managing and listing API keys

You can list all stored API keys in your configuration by invoking `civo apikey list` or remove one by name by using `civo apikey remove apikey_name`.

To see the secret key you can use `civo apikey show` which will show only the default key, to see others just use `civo apikey show NAME`

```sh
civo apikey list
+--------------+---------+
| Name         | Default |
+--------------+---------+
| my_username  | <=====  |
+--------------+---------+
```

```sh
civo apikey show my_username
+-------------+------------+
| Name        | Key        |
+-------------+------------+
| my_username | secret_key |
+-------------+------------+
```

## Instances

#### Introduction

An instance is a virtual server running on the Civo cloud platform. They can be of variable size and you can run any number of them up to your quota on your account.

#### Creating an instance

You can create an instance by running `civo instance create` with a hostname parameter, as well as any options you provide:

```sh
Options:
  -s, --hostname string      the instance's hostname
  -u, --initialuser string   the instance's initial user
  -r, --network string       the instance's network you can use the Name or the ID
  -p, --publicip string      the instance's public ip
  -e, --region string        the region code identifier to have your instance built in
  -i, --size string          the instance's size
  -n, --snapshot string      the instance's snapshot
  -k, --sshkey string        the instance's ssh key you can use the Name or the ID
  -g, --tags string          the instance's tags
  -t, --template string      the instance's template
  -w, --wait                 wait until the instance's is ready
```

Example usage:

```sh
$ civo instance create --hostname=api-demo.test --size g2.small --template=811a8dfb-8202-49ad-b1ef-1e6320b20497 --initialuser=demo-user
 Created instance api-demo.test

$ civo instance show api-demo.test
                 ID : c31e3262-b3f4-467a-b4d2-7980e02d181f
           Hostname : api-demo.test
Openstack Server ID : 986529e8-4063-4576-8952-1cbcd5743e44
             Status : ACTIVE
               Size : g2.small
             Region : lon1
         Network ID : cc3aceb1-66a8-4060-ae35-36e862ba4f82
        Template ID : fffbe2e5-0dd8-476b-b480-cb7c9fccbe39
        Snapshot ID :
       Initial User : demo-user
            SSH Key : bcdd0589-7543-459a-8452-b5fe2252c170
        Firewall ID : default
               Tags :
         Created At : Thu, 18 Jun 2020 03:35:24 +0100
         Private IP : 10.173.156.59
          Public IP : 172.31.6.7 => 185.136.234.101

----------------------------- NOTES -----------------------------
```

You will be able to see the instance's details by running `civo instance show api-demo.test` as above.

#### Viewing the Default User Password For an Instance

You can view the default user's password for an instance by running `civo instance password ID/hostname`

```sh
$ civo instance password api-demo.test
The password for user civo on api-demo.test is 5OaGxNhaN11pLeWB
```

You can also run this command with the option `-o` and `-f` to get only the password output, useful for scripting situations:

```sh
$ civo instance password api-demo.test -o custom -f Password
5OaGxNhaN11pLeWB
```

#### Viewing Instance Public IP Address

If an instance has a public IP address configured, you can display it using `civo instance public-ip ID/hostname`:

```sh
$ civo instance public-ip api-demo.test -o custom -f PublicIP
91.211.152.100
```

The above example uses `-o` and `-f` to display only the IP address in the output.

#### Setting Firewalls

Instances can make use of separately-configured firewalls. By default, an instance is created with no firewall rules set, so you will need to configure some rules (see [Firewalls](#firewalls) for more information).

To associate a firewall with an instance, use the command `civo instance firewall ID/hostname firewall_id`. For example:

```sh
$ civo instance firewall api-demo.test firewall_1
Set api-demo.test to use firewall firewall_1
```

#### Listing Instances

You can list all instances associated with a particular API key by running `civo instance list`.

#### Listing Instances sizes

You can list all instances sizes by running `civo instance size`.

#### Moving a Public IP Between Instances

Given two instances, one with a public IP and one without, you can move the public IP by `civo instance move-ip instance ip_address`:

```sh
$ civo instance move_ip cli-private-ip-demo.test 123.234.123.255`
 Moved public IP 123.234.123.255 to instance cli-private-ip-demo.test
```

#### Rebooting/Restarting Instances

A user can reboot an instance at any time, for example to fix a crashed piece of software. Simply run `civo instance reboot instanceID/hostname`. You will see a confirmation message:

```sh
$ civo instance reboot api-demo.test
 Rebooting api-demo.test. Use 'civo instance show api-demo.test' to see the current status.
```

If you prefer a soft reboot, you can run `civo instance soft-reboot instanceID/hostname` instead.

#### Removing Instances

You can use a command to remove an instance from your account. This is immediate, so use with caution! Any snapshots taken of the instance, as well as any mapped storage, will remain.
Usage: `civo instance remove instanceID/hostname`. For example:

```sh
$ civo instance remove api-demo.test
 Removing instance api-demo.test
```

#### Stopping (Shutting Down) and Starting Instances

You can shut down an instance at any time by running `civo instance stop instanceID/hostname`:

```sh
$ civo instance stop api-demo.test
 Stopping api-demo.test. Use 'civo instance show api-demo.test' to see the current status.
```

Any shut-down instance on your account can be powered back up with `civo instance start instanceID/hostname`:

```sh
$ civo instance start api-demo.test
 Starting api-demo.test. Use 'civo instance show api-demo.test' to see the current status.
```

#### (Re)Tagging an Instance

Tags can be useful in distinguishing and managing your instances. You can retag an instance using `civo instance tags instanceID/hostname 'tag1 tag2 tag3...'` as follows:

```sh
$ civo instance tags api-demo.test 'ubuntu demo web'
 Updated tags on api-demo.test. Use 'civo instance show api-demo.test' to see the current tags.'
$ civo instance show api-demo.test
                ID : 715f95d1-3cee-4a3c-8759-f9b49eec34c4
          Hostname : api-demo.test
              Tags : ubuntu, demo, web
              Size : Small - 2GB RAM, 1 CPU Core, 25GB SSD Disk
            Status : ACTIVE
        Private IP : 10.250.199.4
         Public IP : 172.31.2.164 => 91.211.152.100
           Network : Default (10.250.199.0/24)
          Firewall :  (rules: )
            Region : lon1
      Initial User : api-demouser
      OpenStack ID : 7c89f7de-2b29-4178-a2e5-55bdaa5c4c21
       Template ID : 811a8dfb-8202-49ad-b1ef-1e6320b20497
       Snapshot ID :

----------------------------- NOTES -----------------------------
```

#### Updating Instance Information

In case you need to rename an instance or add notes, you can do so with the `instance update` command as follows:

```sh
$ civo instance update api-demo.test --name api-demo-renamed.test --notes 'Hello, world!'
 Instance 715f95d1-3cee-4a3c-8759-f9b49eec34c4 now named api-demo-renamed.test
 Instance 715f95d1-3cee-4a3c-8759-f9b49eec34c4 notes are now: Hello, world!
$ civo instance show api-demo-renamed.test
                ID : 715f95d1-3cee-4a3c-8759-f9b49eec34c4
          Hostname : api-demo-renamed.test
              Tags : ubuntu, demo, web
              Size : Small - 2GB RAM, 1 CPU Core, 25GB SSD Disk
            Status : ACTIVE
        Private IP : 10.250.199.4
         Public IP : 172.31.2.164 => 91.211.152.100
           Network : Default (10.250.199.0/24)
          Firewall :  (rules: )
            Region : lon1
      Initial User : api-demouser
      OpenStack ID : 7c89f7de-2b29-4178-a2e5-55bdaa5c4c21
       Template ID : 811a8dfb-8202-49ad-b1ef-1e6320b20497
       Snapshot ID :

----------------------------- NOTES -----------------------------

Hello, world!
```

You can leave out either the ``--name`` or `--notes` switch if you only want to update one of the fields.

#### Upgrading (Resizing) an Instance

Provided you have room in your Civo quota, you can upgrade any instance up in size. You can upgrade an instance by using `civo instance upgrade instanceID/hostname new_size` where `new_size` is from the list of sizes at `civo sizes`:

```sh
$ civo instance upgrade api-demo-renamed.test g2.medium
 Resizing api-demo-renamed.test to g2.medium. Use 'civo instance show api-demo-renamed.test' to see the current status.

$ civo instance show api-demo-renamed.test
                ID : 715f95d1-3cee-4a3c-8759-f9b49eec34c4
          Hostname : api-demo-renamed.test
              Tags : ubuntu, demo, web
              Size : Medium - 4GB RAM, 2 CPU Cores, 50GB SSD Disk
            Status : ACTIVE
        Private IP : 10.250.199.4
         Public IP : 172.31.2.164 => 91.211.152.100
           Network : Default (10.250.199.0/24)
          Firewall :  (rules: )
            Region : lon1
      Initial User : api-demouser
  Initial Password : [randomly-assigned-password-here]
      OpenStack ID : 7c89f7de-2b29-4178-a2e5-55bdaa5c4c21
       Template ID : 811a8dfb-8202-49ad-b1ef-1e6320b20497
       Snapshot ID :

----------------------------- NOTES -----------------------------

Hello, world!
```

Please note that resizing can take a few minutes.

## Kubernetes clusters

#### Introduction

You can manage Kubernetes clusters on Civo using the Kubernetes subcommands.

#### List clusters

To see your created clusters, call `civo kubernetes list`:

#### Listing kubernetes sizes

You can list all kubernetes sizes by running `civo kubernetes size`.

```sh
$ civo kubernetes list
+--------------------------------------+----------------+--------+-------+-------+--------+
| ID                                   | Name           | Region | Nodes | Pools | Status |
+--------------------------------------+----------------+--------+-------+-------+--------+
| 5604340f-caa3-4ac1-adb7-40c863fe5639 | falling-sunset | NYC1   |     2 |     1 | ACTIVE |
+--------------------------------------+----------------+--------+-------+-------+--------+
```

#### Create a cluster

You can create a cluster by running `civo kubernetes create` with a cluster name parameter, as well as any options you provide:

```bash
  -a, --applications string          optional, use names shown by running 'civo kubernetes applications ls'
  -h, --help                         help for create
  -m, --merge                        merge the config with existing kubeconfig if it already exists.
  -t, --network string               the name of the network to use in the creation (default "default")
  -n, --nodes int                    the number of nodes to create (the master also acts as a node). (default 3)
  -r, --remove-applications string   optional, remove default application names shown by running  'civo kubernetes applications ls'
      --save                         save the config
  -s, --size string                  the size of nodes to create. (default "g3.k3s.medium")
      --switch                       switch context to newly-created cluster
  -v, --version string               the k3s version to use on the cluster. Defaults to the latest. (default "latest")
  -w, --wait                         a simple flag (e.g. --wait) that will cause the CLI to spin and wait for the cluster to be ACTIVE
```

```sh
$ civo kubernetes create my-first-cluster
Created Kubernetes cluster my-first-cluster
```

#### Adding pools the cluster

You can add more pools to your cluster live (obviously 1 is the minimum) while the cluster is running. It takes the name of the cluster (or the ID), the parameter of `--nodes` which is the new number of nodes to run and `--size` which is the size of the pool,
if `--node` and `--size` are not expecified, the default values will be used.

```sh
civo kubernetes node-pool create my-first-cluster
```

#### Scaling pools in the cluster

You can scale a pool in your cluster, while the cluster is running. It takes the name of the cluster (or the ID), the pool ID, and `--nodes` which is the new number of nodes in the pool

```sh
civo kubernetes node-pool scale my-first-cluster pool-id --nodes 5
```

#### Deleting pools in the cluster

You can delete a pool in your cluster. It takes the name of the cluster (or the ID) and the pool ID

```sh
civo kubernetes node-pool delete my-first-cluster pool-id
```

#### Recycling nodes in your cluster

If you need to recycle a particular node in your cluster for any reason, you can use the `recycle` command. This requires the name or ID of the cluster, and the name of the node you wish to recycle preceded by `--node=`:

```sh
$ civo k8s recycle my-first-cluster --node=kube-node-e9ca
The node (kube-node-e9ca) was recycled
```

*Note:* When a node is recycled, it is fully deleted. The recycle command does not [drain](https://kubernetes.io/docs/tasks/administer-cluster/safely-drain-node/) a node, it simply deletes it before building a new node and attaching it to a cluster. It is intended for scenarios where the node itself develops an issue and must be replaced with a new one.

#### Viewing or Saving the cluster configuration

To output a cluster's configuration information, you can invoke `civo kubernetes config cluster-name`. This will output the `kubeconfig` file to the screen.

You can save a cluster's configuration to your local `~/.kube/config` file. This requires `kubectl` to be installed. Usage:

```sh
civo kubernetes config my-first-cluster -s
Saved config to ~/.kube/config
```

If you already have a `~/.kube/config` file, any cluster configuration that is saved will be *overwritten* unless you also pass the `--merge` option. If you have multiple cluster configurations,  merging allows you to switch contexts at will. If you prefer to save the configuration in another place, just use the parameter `--local-path` or `-p` and the path. If you use `--switch` the cli will automatically change the kubernetes context to the new cluster.

```sh
civo kubernetes config my-first-cluster -s --merge
Merged with main kubernetes config: /root/.kube/config
Saved configuration to: /root/.kube/config

Access your cluster with:
kubectl config use-context my-first-cluster
kubectl get node
```

#### Renaming the cluster

Although the name isn't used anywhere except for in the list of clusters (e.g. it's not in any way written in to the cluster), if you wish to rename a cluster you can do so with:

```sh
civo kubernetes rename my-first-cluster --name="Production"
Kubernetes cluster my-first-cluster is now named Production
```

#### Starting a cluster without default applications

By default, `Traefik` is bundled in with `k3s` to act as the ingress controller. If you want to set up a cluster without `traefik`, you can use the `remove-applications` option in the creation command to start a cluster without it:

```sh
civo kubernetes create --remove-applications=Traefik --nodes=2 --wait
```

#### Removing the cluster

If you're completely finished with a cluster you can delete it with:

```sh
civo kubernetes remove my-first-cluster
```

## Kubernetes Applications

#### Introduction

You can install applications from the [Applications Marketplace](https://github.com/civo/kubernetes-marketplace/) through the command-line interface. The installation depends on whether you are creating a new cluster or adding applications to an existing cluster.

#### Listing Available Applications

To get an up-to-date list of available applications on the Marketplace, run `civo kubernetes apps list`. At the time of writing, the list looked like this:

```text
+---------------------------+-------------+--------------+-----------------+-----------------------------+
| Name                      | Version     | Category     | Plans           | Dependencies                |
+---------------------------+-------------+--------------+-----------------+-----------------------------+
| cert-manager              | v0.11.0     | architecture |                 | Helm                        |
| docker-registry           | ALPHA       | architecture |                 | Helm, cert-manager, Traefik |
| haproxy                   | 1.4.6       | architecture |                 |                             |
| Helm                      | 2.16.5      | management   |                 |                             |
| Jenkins                   | 2.190.1     | ci_cd        | 5GB, 10GB, 20GB | Longhorn                    |
| KubeDB                    | v0.12.0-r1  | database     |                 | Longhorn                    |
| Kubeless                  | 1.0.5       | architecture |                 |                             |
| kubernetes-dashboard      | v2.0.0      | management   |                 |                             |
| Linkerd                   | 2.5.0       | architecture |                 |                             |
| Longhorn                  | 0.7.0       | storage      |                 |                             |
| Maesh                     | Latest      | architecture |                 | Helm                        |
| MariaDB                   | 10.4.7      | database     | 5GB, 10GB, 20GB | Longhorn                    |
| metrics-server            | (default)   | architecture |                 |                             |
| MinIO                     | 2019-08-29  | storage      | 5GB, 10GB, 20GB | Longhorn                    |
| MongoDB                   | 4.2.0       | database     | 5GB, 10GB, 20GB | Longhorn                    |
| OpenFaaS                  | 0.18.0      | architecture |                 | Helm                        |
| Portainer                 | beta        | management   |                 |                             |
| PostgreSQL                |        11.5 | database     | 5GB, 10GB, 20GB | Longhorn                    |
| prometheus-operator       | 0.35.0      | monitoring   |                 |                             |
| Rancher                   | v2.3.0      | management   |                 |                             |
| Redis                     |         3.2 | database     |                 |                             |
| sealed-secrets            | v0.12.4     | architecture |                 |                             |
| Selenium                  | 3.141.59-r1 | ci_cd        |                 |                             |
| system-upgrade-controller | v0.6.2      | architecture |                 |                             |
| Tekton                    | v0.14.0     | ci_cd        |                 |                             |
| Traefik                   | (default)   | architecture |                 |                             |
+---------------------------+-------------+--------------+-----------------+-----------------------------+
```
#### Show Applications details when is installed in the cluster

This option will be allow you to see the post-install instruction of every app installed in the cluster

```sh
$ civo kubernetes application show Traefik apps-demo-cluster
```

the first parameter is for the name of the app and the second is the name of the cluster


#### Installing Applications Onto a New Cluster

To specify applications to install onto a new cluster, list them at cluster creation by specifying their `name` from the list above:

```sh
$ civo kubernetes create apps-demo-cluster --nodes=2  --applications=Redis,Linkerd
```

Now, if you take a look at the cluster's details, you will see the newly-installed applications listed:

```sh
$ civo kubernetes show apps-demo
                ID : 1199efbe-e2a5-4d25-a32f-0b7aa50082b2
              Name : apps-demo-cluster
           # Nodes : 2
              Size : g2.medium
            Status : ACTIVE
           Version : 0.8.1
      API Endpoint : https://[Cluster-IP]:6443
         Master IP : [Cluster-IP]
      DNS A record : 1199efbe-e2a5-4d25-a32f-0b7aa50082b2.k8s.civo.com

Nodes:
+------------------+----------------+--------+
| Name             | IP             | Status |
+------------------+----------------+--------+
| kube-master-1e91 |      (IP)      | ACTIVE |
| kube-node-e678   |      (IP)      | ACTIVE |
+------------------+----------------+--------+

Installed marketplace applications:
+---------+-----------+-----------+--------------+
| Name    | Version   | Installed | Category     |
+---------+-----------+-----------+--------------+
| Traefik | (default) | Yes       | architecture |
| Linkerd | 2.5.0     | Yes       | architecture |
| Redis   | 3.2       | Yes       | database     |
+---------+-----------+-----------+--------------+
```

#### Installing Applications to an Existing Cluster

If you want to add a new application to an existing cluster, you can do so by running the `civo applications` command specifying the cluster and the app(s) you wish to add:

```sh
$ civo kubernetes applications add Longhorn --cluster=apps-demo
Added Longhorn 0.5.0 to Kubernetes cluster apps-demo-cluster
```

#### Installing Applications That Require Plans

Some applications, specifically database apps, require a storage plan that you can specify at installation time from the list of plan options. If you do not provide a plan for an application that requires one, the CLI will notify you and suggest a default size:

```sh
$ civo kubernetes applications add mariadb --cluster=apps-demo
You requested to add MariaDB but didn't select a plan. Please choose one... (5GB, 10GB, 20GB) [5GB]: 10GB
Thank you, next time you could use "MariaDB:10GB" to choose automatically
Added MariaDB 10.4.7 to Kubernetes cluster apps-demo-cluster
```

## Domains and Domain Records

#### Introduction

We host reverse DNS for all instances automatically. If you'd like to manage forward (normal) DNS for your domains, you can do that for free within your account.

This section is effectively split in to two parts: 1) Managing domain names themselves, and 2) Managing records within those domain names.

We don't offer registration of domains names, this is purely for hosting the DNS. If you're looking to buy a domain name, we recommend  [LCN.com](https://www.lcn.com/)  for their excellent friendly support and very competitive prices.

#### Set Up a New Domain

Any user can add a domain name (that has been registered elsewhere) to be managed by Civo.com. You should adjust the nameservers of your domain (through your registrar) to point to  `ns0.civo.com`  and  `ns1.civo.com`.

The command to set up a new domain is `civo domain create domainname`:

```sh
$ civo domain create civoclidemo.xyz
Created a domain called civoclidemo.xyz with ID 418181b2-fcd2-46a2-ba7f-c843c331e79b
```

You can then proceed to add DNS records to this domain.

#### List Domain Names
To see your created domains, call `civo domain list`:

```sh
$ civo domain list
+--------------------------------------+-----------------+
| ID                                   | Name            |
+--------------------------------------+-----------------+
| 418181b2-fcd2-46a2-ba7f-c843c331e79b | civoclidemo.xyz |
+--------------------------------------+-----------------+
```

#### Deleting a Domain

If you choose to delete a domain, you can call `civo domain remove domain_id` and have the system immediately remove the domain and any associated DNS records. This removal is immediate, so use with caution.

#### Creating a DNS Record

A DNS record creation command takes a number of options in the format `civo domain record create domain_id [options]` and the options are this.

```text
Options:
-n, --name string    the name of the record
-p, --priority int   the priority of record only for MX record
-t, --ttl int        The TTL of the record (default 600)
-e, --type string    type of the record (A, CNAME, TXT, SRV, MX)
-v, --value string   the value of the record
```

Usage is as follows:

```sh
$ civo domain record create civoclidemo.xyz -n www -t 600 -e A -v 192.168.1.1

Created a record www1 for civoclidemo.xyz with a TTL of 600 seconds and with a priority of 0 with ID 4e181dde-bde8-4744-8984-067f957a7d59
```

#### Listing DNS Records
You can get an overview of all records you have created for a particular domain by requesting `civo domain record list domain.name`:

```sh
$ civo domain record list civoclidemo.xyz
+--------------------------------------+------+---------------------+-------------+------+----------+
| ID                                   | Type | Name                | Value       | TTL  | Priority |
+--------------------------------------+------+---------------------+-------------+------+----------+
| 4e181dde-bde8-4744-8984-067f957a7d59 | A    | www.civoclidemo.xyz | 192.168.1.1 | 1000 | 0        |
+--------------------------------------+------+---------------------+-------------+------+----------+
```
#### Deleting a DNS Record
You can remove a particular DNS record from a domain you own by requesting `civo domain record remove record_id`. This immediately removes the associated record, so use with caution:

```sh
$ civo domain record remove 4e181dde-bde8-4744-8984-067f957a7d59
The domain record called www with ID 4e181dde-bde8-4744-8984-067f957a7d59 was delete
```

## Firewalls

#### Introduction

You can configure custom firewall rules for your instances using the Firewall component of Civo CLI. These are freely configurable, however customers should be careful to not lock out their own access to their instances. By default, all ports are closed for custom firewalls.

Firewalls can be configured with rules, and they can be made to apply to your chosen instance(s) with subsequent commands.

#### Configuring a New Firewall

To create a new Firewall, use `civo firewall create new_firewall_name`:

```sh
$ civo firewall create civocli_demo
 Created firewall civocli_demo
```

You will then be able to **configure rules** that allow connections to and from your instance by adding a new rule using `civo firewall rule create firewall_id` with the required and your choice of optional parameters, listed here and used in an example below:

```text
Options:
-c, --cidr string Array  the CIDR of the rule you can use (e.g. -c 10.10.10.1/32, 10.10.10.2/32)
-d, --direction string   the direction of the rule need to be ingress
-e, --endport string     the end port of the rule
-h, --help               help for create
-l, --label string       a string that will be the displayed as the name/reference for this rule
-p, --protocol string    the protocol choice (from: TCP, UDP, ICMP)
-s, --startport string   the start port of the rule
```

Example usage:
```sh
$ civo firewall rule create civocli_demo --startport=22 --direction=ingress --label='SSH access for CLI demo'
 New rule SSH access for CLI demo created

$ civo firewall rule list civocli_demo
+--------------------------------------+-----------+----------+------------+----------+-----------+-------------------------+
| ID                                   | Direction | Protocol | Start Port | End Port | Cidr      | Label                   |
+--------------------------------------+-----------+----------+------------+----------+-----------+-------------------------+
| 00270e70-0e1b-498e-9a21-9bcc65736811 | ingress   | tcp      |         22 |          | 0.0.0.0/0 | SSH access for CLI demo |
+--------------------------------------+-----------+----------+------------+----------+-----------+-------------------------+
```

You can see all active rules for a particular firewall by calling `civo firewall rule firewall_id`, where `firewall_id` is the UUID of your particular firewall.

#### Managing Firewalls

You can see an overview of your firewalls using `civo firewall list` showing you which firewalls have been configured with rules, and whether any of your instances are using a given firewall, such as in this case where the firewall we have just configured has the one rule, but no instances using it.

```sh
$ civo firewall list
+--------------------------------------+--------------+-------------+----------------+--------+
| ID                                   | Name         | Total rules | Total Instances | Region |
+--------------------------------------+--------------+-------------+----------------+--------+
| 232d91e9-1550-4c96-bcb6-e9dfecd3e9ee | civocli_demo |           4 |              3 | lon1   |
+--------------------------------------+--------------+-------------+----------------+--------+
```

To configure an instance to use a particular firewall, see [Instances/Setting firewalls elsewhere in this guide](#setting-firewalls).

To get more detail about the specific rule(s) of a particular firewall, you can use `civo firewall rule list firewall_id`.

#### Deleting Firewall Rules and Firewalls

You can remove a firewall rule simply by calling `civo firewall rule remove firewall_id rule_id` - confirming the Firewall ID to delete a particular rule from - as follows:

```sh
$ civo firewall rule remove 09f8d85b-0cf1-4dcf-a472-ba247fb4be21 4070f87b-e6c6-4208-91c5-fc4bc72c1587
  Removed Firewall rule 4070f87b-e6c6-4208-91c5-fc4bc72c1587

$ civo firewall rule list 09f8d85b-0cf1-4dcf-a472-ba247fb4be21
```

Similarly, you can delete a firewall itself by calling `civo firewall remove firewall_id`:

```sh
$ civo firewall remove 09f8d85b-0cf1-4dcf-a472-ba247fb4be21
  Removed firewall 09f8d85b-0cf1-4dcf-a472-ba247fb4be21

$ civo firewall list
```

## Networks

#### Introduction

Civo allows for true private networking if you want to isolate instances from each other. For example, you could set up three instances, keeping one as a [
](https://en.wikipedia.org/wiki/Bastion_host) and load balancer, with instances acting as e.g. a database server and a separate application server, both with private IPs only.

#### Viewing Networks

You can list your currently-configured networks by calling `civo network list`. This will show the network ID, name label and its CIDR range.

#### Creating Networks

You can create a new private network using `civo network create network_label`:

```sh
$ civo network create cli-demo
Create a private network called cli-demo with ID 74b69006-ea59-46a0-96c4-63f5bfa290e1
```

#### Removing Networks

Removal of a network, provided you do not need it and your applications do not depend on routing through it, is simple - you call `civo network remove network_ID`:

```sh
$ civo network remove 74b69006-ea59-46a0-96c4-63f5bfa290e1
Removed the network cli-demo with ID 74b69006-ea59-46a0-96c4-63f5bfa290e1
```

## Load Balancers

#### Introduction

Civo supports load balancing for your instances, allowing you to spread web traffic between them to maximize availability. You can view details about load balancers you may have running, create new ones, update information and even remove them from the command line.

#### Viewing Load Balancers

You can list currently-active load balancers by calling `civo loadbalancer list`. This will draw a table detailing the unique ID, hostname, protocol, port, TLS certificate information, backend check path and connection information.

#### Creating Load Balancers

Create a new load balancer by calling `civo loadbalancer create` as well as any options you provide. The options are:

```text
Options:
-b, --backends stringArray         Specify a backend instance to associate with the load balancer. Takes instance_id, protocol and port in the format --backend=instance:instance-id|instance-name,protocol:http,port:80
-t, --fail_timeout int             Timeout in seconds to consider a backend to have failed. Defaults to 30 (default 30)
-l, --health_check_path string     URL to check for a valid (2xx/3xx) HTTP status on the backends. Defaults to /
-h, --help                         help for create
-e, --hostname string              If not supplied, will be in format loadbalancer-uuid.civo.com
-i, --ignore_invalid_backend_tls   Should self-signed/invalid certificates be ignored from backend servers? Defaults to true (default true)
-x, --max_connections int          Maximum concurrent connections to each backend. Defaults to 10 (default 10)
-m, --max_request_size int         Maximum request content size, in MB. Defaults to 20 (default 20)
    --policy string                <least_conn | random | round_robin | ip_hash> - Balancing policy to choose backends
-r, --port int                     Listening port. Defaults to 80 to match default http protocol (default 80)
-p, --protocol string              Either http or https. If you specify https then you must also provide the next two fields
-c, --tls_certificate string       TLS certificate in Base64-encoded PEM. Required if --protocol is https
-k, --tls_key string               TLS certificate in Base64-encoded PEM. Required if --protocol is https
```

```sh
$ civo loadbalancer create
Created a new Load Balancer with hostname loadbalancer-01da06bc-40ef-4d4c-bb68-d0765d288b54.civo.com
```

#### Updating Load Balancers

Updating an existing load balancer takes the same options as creation, with the syntax being `civo loadbalancer update ID [options]`. For example, we can update the hostname of the load balancer created above using `--hostname`:

```sh
$ civo loadbalancer update 01da06bc-40ef-4d4c-bb68-d0765d288b54 --hostname="civo-demo-loadbalancer.civo.com"
Updated Load Balancer
```

#### Removing Load Balancers

Removing a load balancer is done through calling `civo loadbalancer remove loadbalancer_id`. Please note that this change is immediate:

```sh
$ civo loadbalancer remove 01da06bc-40ef-4d4c-bb68-d0765d288b54
Removed the load balancer civo-demo-loadbalancer.civo.com with ID 01da06bc-40ef-4d4c-bb68-d0765d288b54
```

## Quota

All customers joining Civo will have a default quota applied to their account. The quota has nothing to do with charges or payments, but with the limits on the amount of simultaneous resources you can use. You can view the state of your quota at any time by running `civo quota show`. Here is my current quota usage at the time of writing:

```sh
$ civo quota show
+------------------+-------+-------+
| Item             | Usage | Limit |
+------------------+-------+-------+
| Instances        | 4     | 16    |
| CPU cores        | 5     | 16    |
| RAM MB           | 7168  | 32768 |
| Disk GB          | 150   | 400   |
| Volumes          | 4     | 16    |
| Snapshots        | 1     | 48    |
| Public IPs       | 4     | 16    |
| Subnets          | 1     | 10    |
| Private networks | 1     | 10    |
| Firewalls        | 1     | 16    |
| Firewall rules   | 1     | 160   |
+------------------+-------+-------+
Any items in red are at least 80% of your limit
```

If you have a legitimate need for a quota increase, visit the [Quota page](https://www.civo.com/account/quota) to place your request - we won't unreasonably withhold any increase, it's just in place so we can control the rate of growth of our platform and so that errand scripts using our API don't suddenly exhaust our available resources.

## Sizes

Civo instances come in a variety of sizes depending on your need and budget. You can get details of the sizes of instances available by calling `civo sizes list`. You will get something along the lines of the following, the result will be depends or your selected region:

```sh
$ civo sizes list
+----------------+-------------+------------+-----+----------+-----------+------------+
| Name           | Description | Type       | CPU | RAM (MB) | Disk (GB) | Selectable |
+----------------+-------------+------------+-----+----------+-----------+------------+
| g3.xsmall      | Extra Small | Instance   |   1 |     1024 |        25 | Yes        |
| g3.small       | Small       | Instance   |   1 |     2048 |        25 | Yes        |
| g3.medium      | Medium      | Instance   |   2 |     4096 |        50 | Yes        |
| g3.large       | Large       | Instance   |   4 |     8192 |       100 | Yes        |
| g3.xlarge      | Extra Large | Instance   |   6 |    16384 |       150 | Yes        |
| g3.2xlarge     | 2X Large    | Instance   |   8 |    32768 |       200 | Yes        |
| g3.k3s.xsmall  | Extra Small | Kubernetes |   1 |     1024 |        25 | Yes        |
| g3.k3s.small   | Small       | Kubernetes |   1 |     2048 |        25 | Yes        |
| g3.k3s.medium  | Medium      | Kubernetes |   2 |     4096 |        25 | Yes        |
| g3.k3s.large   | Large       | Kubernetes |   4 |     8192 |        25 | Yes        |
| g3.k3s.xlarge  | Extra Large | Kubernetes |   6 |    16384 |        25 | Yes        |
| g3.k3s.2xlarge | 2X Large    | Kubernetes |   8 |    32768 |        10 | Yes        |
+----------------+-------------+------------+-----+----------+-----------+------------+
```

This command is useful for getting the name of the instance type if you do not remember it - you will need to specify the instance size name when creating an instance using the CLI tool.

Also you can use `--filter` to filter the result by the type, the avalible option are (instance, kubernetes) like this:

```sh
$ civo sizes list --filter kubernetes
+----------------+-------------+------------+-----+----------+-----------+------------+
| Name           | Description | Type       | CPU | RAM (MB) | Disk (GB) | Selectable |
+----------------+-------------+------------+-----+----------+-----------+------------+
| g3.k3s.xsmall  | Extra Small | Kubernetes |   1 |     1024 |        25 | Yes        |
| g3.k3s.small   | Small       | Kubernetes |   1 |     2048 |        25 | Yes        |
| g3.k3s.medium  | Medium      | Kubernetes |   2 |     4096 |        25 | Yes        |
| g3.k3s.large   | Large       | Kubernetes |   4 |     8192 |        25 | Yes        |
| g3.k3s.xlarge  | Extra Large | Kubernetes |   6 |    16384 |        25 | Yes        |
| g3.k3s.2xlarge | 2X Large    | Kubernetes |   8 |    32768 |        10 | Yes        |
+----------------+-------------+------------+-----+----------+-----------+------------+
```
## Snapshots

#### Introduction

Snapshots are a clever way to back up your instances. A snapshot is an exact copy of the instance's virtual hard drive at the moment of creation. At any point, you can restore an instance to the state it was at snapshot creation, or use snapshots to build new instances that are configured exactly the same as other servers you host.

As snapshot storage is chargeable (see [
Quota](#quota)), at any time these can be deleted by you. They can also be scheduled rather than immediately created, and if desired repeated at the same schedule each week (although the repeated snapshot will overwrite itself each week, not keep multiple weekly snapshots).

#### Creating Snapshots

You can create a snapshot from an existing instance on the command line by using `civo snapshot create snapshot_name instance_id`
For a one-off snapshot that's all you will need:

```sh
$ civo snapshot create CLI-demo-snapshot 715f95d1-3cee-4a3c-8759-f9b49eec34c4
Created snapshot CLI-demo-snapshot with ID d6d7704b-3402-44d0-aeb1-09875f71d168
```

For scheduled snapshots, include the `-c '0 * * * *'` switch, where the `'0 * * * *'` string is in `cron` format.

Creating snapshots is not instant, and will take a while depending on the size of the instance being backed up. You will be able to monitor the status of your snapshot by listing your snapshots as described below.

#### Listing Snapshots

You can view all your currently-stored snapshots and a bit of information about them by running `civo snapshot list`:

```sh
$ civo snapshot list
+--------------------------------------+-------------------+------+---------------------+---------+------------+-------------------------------+-------------------------------+-------------------------------+
| ID                                   | Name              | Size | Hostname            | State   | Cron       | Schedule                      | RequestedAt                   | CompletedAt                   |
+--------------------------------------+-------------------+------+---------------------+---------+------------+-------------------------------+-------------------------------+-------------------------------+
| d593263e-2433-4c82-aad3-81c2d6ddaa09 | CLI-demo-snapshot | 0 GB | www1                | pending | 55 0 * * * | Thu, 18 Jun 2020 00:55:00 CDT | Mon, 01 Jan 0001 00:00:00 UTC | Mon, 01 Jan 0001 00:00:00 UTC |
+--------------------------------------+-------------------+------+---------------------+---------+------------+-------------------------------+-------------------------------+-------------------------------+
```

#### Removing Snapshots

Snapshots that are not associated with an instance can be removed using `civo snapshot remove snapshot_id` as follows:

```sh
$ civo snapshot remove d6d7704b-3402-44d0-aeb1-09875f71d168
Removed snapshot CLI-demo-snapshot with ID d6d7704b-3402-44d0-aeb1-09875f71d168
```

If an instance was created from a snapshot, you will not be able to remove the snapshot itself.

## SSH Keys

#### Introduction

To manage the SSH keys for an account that are used to log in to cloud instances, the Civo CLI tool provides the following commands. You would need to [
generate a new key](https://www.civo.com/learn/ssh-key-basics) according to your particular circumstances, if you do not have a suitable SSH key yet.

#### Uploading a New SSH Key

You will need the path to your public SSH Key to upload a new key to Civo. The usage is as follows: `civo ssh create NAME --key /path/to/FILENAME`

#### Listing Your SSH Keys

You will be able to list the SSH keys known for the current account holder by invoking `civo ssh list`:

```sh
$ civo sshkeys
+--------------------------------------+------------------+----------------------------------------------------+
| ID                                   | Name             | Fingerprint                                        |
+--------------------------------------+------------------+----------------------------------------------------+
| 8aa45fea-a395-471c-93a6-27485a8429f3 | civo_cli_demo    | SHA256:[Unique SSH Fingerprint]                    |
+--------------------------------------+------------------+----------------------------------------------------+
```

#### Removing a SSH Key

You can delete a SSH key by calling `remove` for it by ID:
```sh
$ civo ssh remove 531d0998-4152-410a-af20-0cccb1c7c73b
Removed SSH key cli-demo with ID 531d0998-4152-410a-af20-0cccb1c7c73b
```

## Templates

#### Introduction

Civo instances are built from a template that specifies a disk image. Templates can contain the bare-bones OS install such as Ubuntu or Debian, or custom pre-configured operating systems that you can create yourself from a bootable volume. This allows you to speedily deploy pre-configured instances.

#### Listing Available Template Images

A simple list of available templates, both globally-defined ones and user-configured account-specific templates, can be seen by running `civo template list`:

```sh
$ civo template list
+--------------------------------------+----------------+----------------+--------------------------------------+----------------------------------------------------+-------------+------------------+
| ID                                   | Code           | Name           | Image ID                             | Short Description                                  | Description | Default Username |
+--------------------------------------+----------------+----------------+--------------------------------------+----------------------------------------------------+-------------+------------------+
| 458ae900-30e0-4ade-bd68-d137d57d4e47 | centos-7       | CentOS 7       | e17ec38a-1e77-4c45-bef3-569567c9b169 | CentOS 7 - aiming to be compatible with RHEL 7     |             | centos           |
| 033c35a0-a8c3-4518-8114-d156a4d4c512 | debian-stretch | Debian Stretch | 2ffff07e-6953-4864-8ce9-1f754d70de31 | Debian v9 (Stretch), current stable Debian release |             | admin            |
| b0d30599-898a-4072-86a1-6ed2965320d9 | ubuntu-16.04   | Ubuntu 16.04   | 8b4d81e0-6283-4ea3-bbc4-478df568024e | Ubuntu 16.04                                       |             | ubuntu           |
| 811a8dfb-8202-49ad-b1ef-1e6320b20497 | ubuntu-18.04   | Ubuntu 18.04   | e4838e89-f086-41a1-86b2-60bc4b0a259e | Ubuntu 18.04                                       |             | ubuntu           |
| fffbe2e5-0dd8-476b-b480-cb7c9fccbe39 | debian-buster  | Debian Buster  | 38686161-ba25-4899-ac0a-54eaf35239c0 | Debian v10 (Buster), latest stable Debian release  |             | admin            |
+--------------------------------------+----------------+----------------+--------------------------------------+----------------------------------------------------+-------------+------------------+
```

#### Viewing Details of a Template

Detailed information about a template can be obtained via the CLI using `civo template show template_ID`.


#### Creating a Template

You can convert a **bootable** Volume (virtual disk) of an instance, or alternatively use an existing image ID, to create a template. The options for the `civo template create` command are:

```text
Options:
-i, --cloudconfig string         The path of the cloud config
-c, --code string                The code name of the template, this can't change after creation
-u, --default-username string    The default username of the template
-d, --description string         Add a description
-h, --help                       help for create
-m, --image-id string            The image id for the template
-n, --name string                The name of the template
-s, --short-description string   Add a short description
-v, --volume-id string           The volume id for the template
```

```sh
$ civo template create -n="cli-demo" -v=1427e49f-d159-4421-b6cc-34c43775764b --description="This is a demo template made from a CoreOS image" --short-description="CoreOS CLI demo"
	Created template cli-demo
```

#### Updating Template Information

Once you have  created a custom template, you can update information that allows for the easy identification and management of the template. Usage is `civo template update template_id [options]`:

```text
Options:
-i, --cloudconfig string         The cloud config
-u, --default-username string    The default username of the template
-d, --description string         Add a description
-h, --help                       help for update
-n, --name string                The name of the template
-s, --short-description string   Add a short description
```

#### Removing a Template

Removing an account-specific template is done using the `template remove template_id` command:

```sh
$ civo template remove 1427e22f-d149-4421-b6ab-34c43754224c
```

Please note that template removal is immediate! Use with caution.

## Volumes

#### Introduction

Volumes are flexible-size additional storage for instances. By creating and associating a Volume with an instance, an additional virtual disk will be made available for backups or database files that can then moved to another instance.

Volumes take disk space on your account's quota, and can only be created up to this quota limit. For more information about the quota system, see [Quota](#quota).

#### Creating a Volume

You can create a new volume by calling `civo volume create NAME SIZE(GB)`:

```text
Options:
-b, --bootable      Mark the volume as bootable
-h, --help          help for create
-s, --size-gb int   The new size in GB (required)
```

```sh
$ civo volume create CLI-demo-volume -s 25
Created a new 25GB volume called CLI-demo-volume with ID 9b232ffa-7e05-45a4-85d8-d3643e68952e
```

#### Attaching a Volume to an Instance

Mounting (Attaching) a volume onto an instance will allow that instance to use the volume as a drive:

```sh
$ civo volume attach 9b232ffa-7e05-45a4-85d8-d3643e68952e 715f95d1-3cee-4a3c-8759-f9b49eec34c4
Attached volume CLI-demo-volume with ID 9b232ffa-7e05-45a4-85d8-d3643e68952e to api-demo.test
```

If this is a newly-created volume, you would need to partition, format and mount the volume. For more information, [see the Learn guide here](https://www.civo.com/learn/configuring-block-storage-on-civo).
Note: You can only attach a volume to one instance at a time.

#### Detaching a Volume From an Instance

If you want to detach a volume to move it to another instance, or are just finished with it, you can detach it once it's been [unmounted](https://www.civo.com/learn/configuring-block-storage-on-civo) using `civo volume detach volume_id`:

```sh
$ civo volume detach 9b232ffa-7e05-45a4-85d8-d3643e68952e
Detached volume CLI-demo-volume with ID 9b232ffa-7e05-45a4-85d8-d3643e68952e
```

#### Listing Volumes

You can get an overall view of your volumes, their sizes and status by using `civo volume list`.

#### Resizing Volumes

An un-attached volume can be resized if you need extra space. This is done by calling `civo volume resize volume_id -s new_size` where `new-size` is in gigabytes:

```sh
$ civo volume resize 9b232ffa-7e05-45a4-85d8-d3643e68952e -s 30
Resized volume CLI-demo-volume with ID 9b232ffa-7e05-45a4-85d8-d3643e68952e to be 30GB
```

#### Deleting Volumes

To free up quota and therefore the amount to be billed to your account, you can delete a volume through `civo volume delete volume_id`. This deletion is immediate:

```sh
$ civo volume delete 9b232ffa-7e05-45a4-85d8-d3643e68952e
Removed volume CLI-demo-volume with ID 9b232ffa-7e05-45a4-85d8-d3643e68952e (was 30GB)
$ civo volume list
```


## Region
As Civo grows, more regions for your instances will become available. You can run `civo region ls` to list the regions available. Block storage (Volumes) is region-specific, so if you configure an instance in one region, any volumes you wish to attach to that instance would have to be in the same region.

#### List all region

You can run `civo region ls` to get the list of all region
```sh
+------+-------------+----------------+---------+
| Code | Name        | Country        | Current |
+------+-------------+----------------+---------+
| NYC1 | New York 1  | United States  | <=====  |
| SVG1 | Stevenage 1 | United Kingdom |         |
+------+-------------+----------------+---------+
```

#### Change region

To change the region the only cmd you need run is `civo region current <REGION-CODE>` and you will see a message like this
```sh
 civo region current NYC1

```
The default region was set to (New York 1) NYC1

#### Use region in non-interactive mode

To set the region in non-interactive mode, you only need pass to the command this `--region <REGION-CODE>` like this

```sh
civo kubernetes create production-01 -n 4 --wait --region NYC1
```

## Enabling shell autocompletion

The civo binary is delivered with the support for bash, zsh, powershell and fish, and you can use in this way

`civo completion [bash|zsh|powershell|fish]`

Sourcing the completion script in your shell enables civo autocompletion.

However, the completion script depends on bash-completion, which means that you have to install this software first (you can test if you have bash-completion already installed by running `type _init_completion`).

### Install bash-completion

bash-completion is provided by many package managers (see [here](https://github.com/scop/bash-completion#installation)). You can install it with `apt-get install bash-completion` or `yum install bash-completion`, etc.

The above commands create `/usr/share/bash-completion/bash_completion`, which is the main script of bash-completion. Depending on your package manager, you have to manually source this file in your `~/.bashrc` file.

To find out, reload your shell and run `type _init_completion`. If the command succeeds, you're already set, otherwise add the following to your `~/.bashrc` file:

```shell
source /usr/share/bash-completion/bash_completion
```

Reload your shell and verify that bash-completion is correctly installed by typing `type _init_completion`.

#### Enable civo autocompletion

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

### Install zsh-completion

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

To set the civo completion code for zsh to auto-load on start up yo can run this command.
```bash
civo completion zsh > "${fpath[1]}/_civo"
```

## Contributing

Bug reports and pull requests are welcome on GitHub at https://github.com/civo/cli.

## License

The code is available as open source under the terms of the [MIT License](https://opensource.org/licenses/MIT).
