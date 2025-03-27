# Civo Command-Line Client

## Introduction

Civo CLI is a tool to manage your [Civo.com](https://www.civo.com) account from the terminal. The [Civo web control panel](https://www.civo.com/account/) has a user-friendly interface for managing your account, but in case you want to automate or run scripts on your account, or have multiple complex services, the command-line interface outlined here will be useful. This guide will cover the set-up and usage of the Civo CLI tool with examples.

**STATUS:** This project is currently under active development and maintenance.

## Table of contents

- [Introduction](#introduction)
- [Global Options](#global-options)
- [Set-Up](#set-up)
- [Version/Updating CLI](#update)
- [Docker Usage](#docker-usage)
- [API Keys](#api-keys)
- [Instances](#instances)
- [Kubernetes clusters](#kubernetes-clusters)
- [Kubernetes applications](#kubernetes-applications)
- [Domains and Domain Records](#domains-and-domain-records)
- [Firewalls](#firewalls)
- [Networks](#networks)
- [Database Backup and Restore](./doc/DTABASE_BACKUP_RESTORE.md)
- [Object Stores](#object-stores)
- [Object Store Credentials](#object-store-credentials)
- [Load Balancers](#load-balancers)
- [Quota](#quota)
- [Sizes](#sizes)
- [SSH Keys](#ssh-keys)
- [DiskImages](#disk-image)
- [Volumes](#volumes)
- [Teams](#teams)
- [Permissions](#permissions)
- [Region](#region)
- [Enabling shell autocompletion](#enabling-shell-autocompletion)
- [Contributing](#contributing)
- [License](#license)

## Set-up

Civo CLI is built with Go and distributed as binary files, available for multiple operating systems and downloadable from <https://github.com/civo/cli/releases>.

### Installing on macOS

If you have a Mac, you can install it using [Homebrew](https://brew.sh):

```sh
brew tap civo/tools
brew install civo
```

or if you prefer you can run this in the console:

```sh
curl -sL https://civo.com/get | sh
```

### Installing on Windows

Civo CLI is available to download on windows via Chocolatey and Scoop

For installing via Chocolatey you need [Chocolatey](https://chocolatey.org/install) package manager installed on your PC.

- run the following command after confirming Chocolatey on your PC

 ```
 choco install civo-cli
 ```

 and it will install Civo CLI on your PC.

For installing via Scoop you need [Scoop](https://scoop.sh/) installed as a package manager, then:

- add the extras bucket with

```
scoop bucket add extras
```

- install civo with

```
scoop install civo
```

You will also, of course, need a Civo account, for which you can [register here](https://www.civo.com/signup).

### Installing on Linux

For Linux Civo CLI can be installed by various methods.

- Install via the direct shell script:

```sh
curl -sL https://civo.com/get | sh
```

- Install via the brew package manager, as shown in the above instructions for MacOS.

- Install via wget, specifying the [release version](https://github.com/civo/cli/releases) you want.

**_Note that the version in the example below may not be the latest. Specify the version based on the latest available if you are using this method._**

```
wget https://github.com/civo/cli/releases/download/v1.0.40/civo-1.0.40-linux-amd64.tar.gz
tar -xvf civo-1.0.40-linux-amd64.tar.gz
chmod +x civo
mv ./civo /usr/local/bin/
```

-   You can also build the binary, but make sure you have go installed,

```sh
git clone https://github.com/civo/cli.git
cd cli
make
cd ..
cp -r cli $HOME
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

## Version/Updating CLI

Every user receives a reminder to update the CLI once in 24 hours as well as notified to update the CLI version in case of any error. Run `civo update` to update the CLI to the latest version.

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
docker run -it --rm -v $HOME/.civo.json:/.civo.json -v $HOME/.kube/config:/root/.kube/config civo/cli:latest
```

To make usage easier, an alias is recommended. Here's an example how to set one to the same command as would be used if installed directly on the system, and using the Docker image:

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

The civo cli have multiple global options, that you can use, like this:

```
      --config string   config file (default is $HOME/.civo.json)
  -f, --fields string   output fields for custom format output (use -h to determine fields)
  -h, --help            help for civo
  -o, --output string   output format (json/human/custom) (default "human")
      --pretty          Print pretty the json output
      --region string   Choose the region to connect to, if you use this option it will use it over the default region
  -y, --yes             Automatic yes to prompts; assume "yes" as answer to all prompts and run non-interactive
```

## API Keys

#### Introduction

In order to use the command-line tool, you will need to authenticate yourself to the Civo API using a special key. You can find an automatically-generated API key or regenerate a new key at [https://www.civo.com/api](https://www.civo.com/api). The CLI have a global env variable call `CIVO_TOKEN` if this is set, the CLI will use this token instead of the one in the config file. This is useful for scripting situations or CI/CD pipelines. When you set the `CIVO_TOKEN` you will see a new apikey entry with the name `tempKey` in the `civo apikey list` command.

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

By default, the Civo account credentials API Key along with other settings like region will be saved in a file called `.civo.json` in the user home directory. The default location of the file can be changed using the environment variable `CIVO_CONFIG`.

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
  -t, --diskimage string     the instance's disk image name (from 'civo diskimage ls' command)
  -l, --firewall string      the instance's firewall you can use the Name or the ID
  -h, --help                 help for create
  -s, --hostname string      the instance's hostname
  -u, --initialuser string   the instance's initial user
  -r, --network string       the instance's network you can use the Name or the ID
  -p, --publicip string      This should be either "none" or "create" (default "create")
  -i, --size string          the instance's size (from 'civo instance size' command)
  -k, --sshkey string        the instance's ssh key you can use the Name or the ID
  -g, --tags string          the instance's tags
  -w, --wait                 wait until the instance's is ready

```

Example usage:

```sh
$ civo instance create --hostname=api-demo.test --size g3.small  --diskimage=ubuntu-focal --initialuser=demo-user
  The instance api-demo.test has been created

$ civo instance show api-demo.test
              ID : 112f2407-fb89-443e-bd0e-5ddabc4682c6
        Hostname : api-demo.test
          Status : ACTIVE
            Size : g3.small
       Cpu Cores : 1
             Ram : 2048
        SSD disk : 25
          Region : LON1
      Network ID : 28244c7d-b1b9-48cf-9727-aebb3493aaac

   ID : ubuntu-bionic
     Snapshot ID :
    Initial User : demo-user
Initial Password : demo-user
         SSH Key :
     Firewall ID : c9e14ae8-b8eb-4bae-a687-9da4637233da
            Tags :
      Created At : Mon, 01 Jan 0001 00:00:00 UTC
      Private IP : 192.168.1.7
       Public IP : 74.220.21.246

----------------------------- NOTES -----------------------------
```

You will be able to see the instance's details by running `civo instance show api-demo.test` as above.

#### Disk images and instance sizes

You can view the Disk images by running `civo diskimage ls`

```sh
$ civo diskimage ls
+--------------------------------------+-----------------+----------------+-----------+--------------+
| ID                                   | Name            | Version        | State     | Distribution |
+--------------------------------------+-----------------+----------------+-----------+--------------+
| 9b661c46-ac4f-46e1-9f3d-aaacde9b4fec | debian-9        |              9 | available | debian       |
+--------------------------------------+-----------------+----------------+-----------+--------------+
| d927ad2f-5073-4ed6-b2eb-b8e61aef29a8 | ubuntu-focal    |          20.04 | available | ubuntu       |
+--------------------------------------+-----------------+----------------+-----------+--------------+
| a4204155-a876-43fa-b4d6-ea2af8774560 | debian-10       |             10 | available | debian       |
+--------------------------------------+-----------------+----------------+-----------+--------------+
| eda67ea0-4282-4945-9b7b-d3e1cba1d987 | ubuntu-jammy    |          22.04 | available | ubuntu       |
+--------------------------------------+-----------------+----------------+-----------+--------------+
| 170db96f-8458-44aa-83ca-0c31fb81a835 | rocky-9-1       |            9.1 | available | rocky        |
+--------------------------------------+-----------------+----------------+-----------+--------------+
| 25fbbd96-d5ec-4d08-9c75-a5e154dabf9b | debian-11       |             11 | available | debian       |
+--------------------------------------+-----------------+----------------+-----------+--------------+
| 21613daa-a66b-44fc-87f5-b6db566d8f91 | ubuntu-cuda11-8 | 22.04-cuda11-8 | available | ubuntu       |
+--------------------------------------+-----------------+----------------+-----------+--------------+
| ffb6fd93-cb06-4e8d-8058-46003b78e2ff | talos-v1.2.8    | 1.25.5         | available | talos        |
+--------------------------------------+-----------------+----------------+-----------+--------------+
| 9a16a77e-1a1f-45c8-87fd-6d1a19eeaac9 | talos-v1.5.0    | 1.27.0         | available | talos        |
+--------------------------------------+-----------------+----------------+-----------+--------------+
| 13232803-0928-4634-9ab8-476bff29ef1b | ubuntu-cuda12-2 | 22.04-cuda12-2 | available | ubuntu       |
+--------------------------------------+-----------------+----------------+-----------+--------------+
```

You can view the instance sizes list by running `civo size ls`

```sh
$ civo size ls
+--------------------+--------------------------------+------------+-----+---------+-----+
| Name               | Description                    | Type       | CPU | RAM     | SSD |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g3.xsmall          | Extra Small                    | Instance   |   1 |    1024 |  25 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g3.small           | Small                          | Instance   |   1 |    2048 |  25 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g3.medium          | Medium                         | Instance   |   2 |    4096 |  50 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g3.large           | Large                          | Instance   |   4 |    8192 | 100 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g3.xlarge          | Extra Large                    | Instance   |   6 |   16384 | 150 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g3.2xlarge         | 2X Large                       | Instance   |   8 |   32768 | 200 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4s.kube.xsmall    | Extra Small - Standard         | Kubernetes |   1 |    1024 |  30 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4s.kube.small     | Small - Standard               | Kubernetes |   1 |    2048 |  40 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4s.kube.medium    | Medium - Standard              | Kubernetes |   2 |    4096 |  50 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4s.kube.large     | Large - Standard               | Kubernetes |   4 |    8192 |  60 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4p.kube.small     | Small - Performance            | Kubernetes |   4 |   16384 |  60 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4p.kube.medium    | Medium - Performance           | Kubernetes |   8 |   32768 |  80 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4p.kube.large     | Large - Performance            | Kubernetes |  16 |   65536 | 120 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4p.kube.xlarge    | Extra Large - Performance      | Kubernetes |  32 |  131072 | 180 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4c.kube.small     | Small - CPU optimized          | Kubernetes |   8 |   16384 |  60 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4c.kube.medium    | Medium - CPU optimized         | Kubernetes |  16 |   32768 |  80 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4c.kube.large     | Large - CPU optimized          | Kubernetes |  32 |   65536 | 120 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4c.kube.xlarge    | Extra Large - CPU optimized    | Kubernetes |  64 |  131072 | 180 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4m.kube.small     | Small - RAM optimized          | Kubernetes |   2 |   16384 |  60 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4m.kube.medium    | Medium - RAM optimized         | Kubernetes |   4 |   32768 |  80 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4m.kube.large     | Large - RAM optimized          | Kubernetes |   8 |   65536 | 120 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4m.kube.xlarge    | Extra Large - RAM optimized    | Kubernetes |  16 |  131072 | 180 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g3.db.xsmall       | Extra Small                    | Database   |   1 |    2048 |  20 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g3.db.small        | Small                          | Database   |   2 |    4096 |  40 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g3.db.medium       | Medium                         | Database   |   4 |    8192 |  80 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g3.db.large        | Large                          | Database   |   6 |   16384 | 160 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g3.db.xlarge       | Extra Large                    | Database   |   8 |   32768 | 320 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g3.db.2xlarge      | Double Extra Large             | Database   |  10 |   65536 | 640 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g3.kf.small        | Small - CPU optimized          | KfCluster  |   4 |   16384 |  60 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g3.kf.medium       | Medium - CPU optimized         | KfCluster  |   8 |   32768 |  80 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g3.kf.large        | Large - CPU optimized          | KfCluster  |  16 |   65536 | 120 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g3.kf.xlarge       | Extra Large - CPU optimized    | KfCluster  |  32 |  131072 | 180 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4s.xsmall         | xSmall - Standard              | Instance   |   1 |    1024 |  25 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4s.small          | Small - Standard               | Instance   |   1 |    2048 |  25 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4s.medium         | Medium - Standard              | Instance   |   2 |    4096 |  50 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4s.large          | Large - Standard               | Instance   |   4 |    8192 | 100 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4s.xlarge         | Extra Large - Standard         | Instance   |   6 |   16384 | 150 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4s.2xlarge        | 2X Large - Standard            | Instance   |   8 |   32768 | 200 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4g.kube.small     | Small - Nvidia A100 80GB       | Kubernetes |  12 |  131072 | 100 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4g.kube.medium    | Medium - Nvidia A100 80GB      | Kubernetes |  24 |  262144 | 100 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4g.kube.large     | Large - Nvidia A100 80GB       | Kubernetes |  48 |  524288 | 100 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4g.kube.xlarge    | Extra Large - Nvidia A100 80GB | Kubernetes |  96 | 1048576 | 100 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4g.small          | Small - Nvidia A100 80GB       | Instance   |  12 |  131072 | 200 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4g.medium         | Medium - Nvidia A100 80GB      | Instance   |  24 |  262144 | 200 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4g.large          | Large - Nvidia A100 80GB       | Instance   |  48 |  524288 | 400 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4g.xlarge         | Extra Large - Nvidia A100 80GB | Instance   |  96 | 1048576 | 400 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4g.40.small       | Small - Nvidia A100 40GB       | Instance   |   8 |   65536 | 200 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4g.40.medium      | Medium - Nvidia A100 40GB      | Instance   |  16 |  131072 | 400 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4g.40.large       | Large - Nvidia A100 40GB       | Instance   |  32 |  262133 | 400 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4g.40.xlarge      | Extra Large - Nvidia A100 40GB | Instance   |  64 |  524288 | 400 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4g.40.kube.small  | Small - Nvidia A100 40GB       | Kubernetes |   8 |   65536 | 200 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4g.40.kube.medium | Medium - Nvidia A100 40GB      | Kubernetes |  16 |  131072 | 400 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4g.40.kube.large  | Large - Nvidia A100 40GB       | Kubernetes |  32 |  262133 | 400 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4g.40.kube.xlarge | Extra Large - Nvidia A100 40GB | Kubernetes |  64 |  524288 | 400 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| an.g1.l40s.x1      | Small - Nvidia L40S 40GB       | Instance   |  12 |  131072 | 200 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| an.g1.l40s.x2      | Medium - Nvidia L40S 40GB      | Instance   |  24 |  262133 | 200 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| an.g1.l40s.x4      | Large - Nvidia L40S 40GB       | Instance   |  48 |  524288 | 400 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| an.g1.l40s.x8      | Extra Large - Nvidia L40S 40GB | Instance   |  96 | 1048576 | 400 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| an.g1.l40s.kube.x1 | Small - Nvidia L40S 40GB       | Kubernetes |  12 |  131072 | 200 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| an.g1.l40s.kube.x2 | Medium - Nvidia L40S 40GB      | Kubernetes |  24 |  262133 | 200 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| an.g1.l40s.kube.x4 | Large - Nvidia L40S 40GB       | Kubernetes |  48 |  524288 | 400 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| an.g1.l40s.kube.x8 | Extra Large - Nvidia L40S 40GB | Kubernetes |  96 | 1048576 | 400 |
+--------------------+--------------------------------+------------+-----+---------+-----+

```

#### Viewing the Default User Password For an Instance

You can view the default user's password for an instance by running `civo instance password ID/hostname`

```sh
$ civo instance password api-demo.test
The instance api-demo.test (112f2407-fb89-443e-bd0e-5ddabc4682c6) has the password BrbXNW2RUYLe (and user demo-user)
```

You can also run this command with the option `-o` and `-f` to get only the password output, useful for scripting situations:

```sh
$ civo instance password api-demo.test -o custom -f Password
BrbXNW2RUYLe
```

#### Viewing Instance Public IP Address

If an instance has a public IP address configured, you can display it using `civo instance public-ip ID/hostname`:

```sh
$ civo instance show api-demo.test -o custom -f public_ip
74.220.21.246
```

The above example uses `-o` and `-f` to display only the IP address in the output.

#### Setting Firewalls

Instances can make use of separately-configured firewalls. By default, an instance is created with default all port open firewall rules set. If you want to secure your instances more, so you will need to configure some rules (see [Firewalls](#firewalls) for more information). Once you have configured the rules you can check it by running `civo firewall ls`

```sh
$ civo firewall ls
+--------------------------------------+------------------------------------+---------+-------------+-----------------+
| ID                                   | Name                               | Network | Total rules | Total Instances |
+--------------------------------------+------------------------------------+---------+-------------+-----------------+
| c9e14ae8-b8eb-4bae-a687-9da4637233da | Default (all open)                 | Default |           3 |               0 |
| f79db64d-41f0-4be0-ae80-ce4499164319 | Kubernetes cluster: Demo           | Default |           2 |               0 |
+--------------------------------------+------------------------------------+---------+-------------+-----------------+
```

You can take the firewall ID and use it to associate a firewall with an instance, use the command `civo instance firewall ID/hostname firewall_id`. For example:

```sh
$ civo instance firewall api-demo.test f79db64d-41f0-4be0-ae80-ce4499164319
Set the firewall for the instance api-demo.test (112f2407-fb89-443e-bd0e-5ddabc4682c6) to Kubernetes cluster: Demo (f79db64d-41f0-4be0-ae80-ce4499164319)
```

#### Listing Instances

You can list all instances associated with a particular API key by running `civo instance list`.

#### Listing Instances sizes

You can list all instances sizes by running `civo instance size`.

#### Rebooting/Restarting Instances

A user can reboot an instance at any time, for example to fix a crashed piece of software. Simply run `civo instance reboot instanceID/hostname`. You will see a confirmation message:

```sh
$ civo instance reboot api-demo.test
 Rebooting api-demo.test. Use 'civo instance show api-demo.test' to see the current status.
```

If you prefer a hard reboot, you can run `civo instance hard-reboot instanceID/hostname` instead.

#### Removing Instances

You can use a command to remove an instance from your account. This is immediate, so use with caution! Any snapshots taken of the instance, as well as any mapped storage, will remain.
Usage: `civo instance remove instanceID/hostname`. For example:

```sh
$ civo instance remove api-demo.test
The instance (api-demo.test) has been deleted
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
$ civo instance tags api-demo.test 'ubuntu demo'
 The instance api-demo.test (b5f82fa7-b8e4-44aa-9dda-df02dab71d6c) has been tagged with ubuntu demo. Use 'civo instance show api-demo.test' to see the current tags.'
$ civo instance show api-demo.test
              ID : b5f82fa7-b8e4-44aa-9dda-df02dab71d6c
        Hostname : api-demo.test
          Status : SHUTOFF
            Size : g3.small
       Cpu Cores : 1
             Ram : 2048
        SSD disk : 25
          Region : LON1
      Network ID : 28244c7d-b1b9-48cf-9727-aebb3493aaac
   Disk image ID : ubuntu-bionic
     Snapshot ID :
    Initial User : demo-user
Initial Password : demo-user
         SSH Key :
     Firewall ID : c9e14ae8-b8eb-4bae-a687-9da4637233da
            Tags : ubuntu demo
      Created At : Mon, 01 Jan 0001 00:00:00 UTC
      Private IP : 192.168.1.7
       Public IP : 74.220.16.23

----------------------------- NOTES -----------------------------
```

#### Updating Instance Information

In case you need to rename an instance or add notes, you can do so with the `instance update` command as follows:

```sh
$ civo instance update api-demo.test --hostname api-demo-renamed.test --notes 'Hello, world!'
The instance api-demo-renamed.test (b5f82fa7-b8e4-44aa-9dda-df02dab71d6c) has been updated

$ civo instance show api-demo-renamed.test
              ID : 715f95d1-3cee-4a3c-8759-f9b49eec34c4
        Hostname : api-demo-renamed.test
          Status : ACTIVE
            Size : g3.small
       Cpu Cores : 1
             Ram : 2048
        SSD disk : 25
          Region : LON1
      Network ID : 28244c7d-b1b9-48cf-9727-aebb3493aaac
   Disk image ID : ubuntu-bionic
     Snapshot ID :
    Initial User : demo-user
Initial Password : demo-user
         SSH Key :
     Firewall ID : c9e14ae8-b8eb-4bae-a687-9da4637233da
            Tags :
      Created At : Mon, 01 Jan 0001 00:00:00 UTC
      Private IP : 192.168.1.7
       Public IP : 74.220.21.246

----------------------------- NOTES -----------------------------

Hello, world!
```

You can leave out either the `--name` or `--notes` switch if you only want to update one of the fields.

#### Upgrading (Resizing) an Instance

Provided you have room in your Civo quota, you can upgrade any instance up in size. You can upgrade an instance by using `civo instance upgrade instanceID/hostname new_size` where `new_size` is from the list of sizes at `civo sizes ls`:

```sh
$ civo instance upgrade api-demo-renamed.test g3.medium
 The instance api-demo-renamed.test (9579d478-a09e-4196-a08c-f52545a90fea) is being upgraded to g3.medium

$ civo instance show api-demo-renamed.test
          Status : ACTIVE
            Size : g3.medium
       Cpu Cores : 2
             Ram : 4096
        SSD disk : 50
          Region : LON1
      Network ID : 28244c7d-b1b9-48cf-9727-aebb3493aaac
   Disk image ID : ubuntu-bionic
     Snapshot ID :
    Initial User : demo-user
Initial Password : demo-user
         SSH Key :
     Firewall ID : c9e14ae8-b8eb-4bae-a687-9da4637233da
            Tags : ubuntu, demo
      Created At : Mon, 01 Jan 0001 00:00:00 UTC
      Private IP : 192.168.1.9
       Public IP : 74.220.17.71

----------------------------- NOTES -----------------------------

Hello, world!
```

Please note that resizing can take a few minutes.


### Recovery Mode

Recovery mode allows you to boot your instance into a recovery environment for troubleshooting.

```sh
# Enable recovery mode for an instance
civo instance recovery enable INSTANCE_ID/HOSTNAME

# Disable recovery mode for an instance
civo instance recovery disable INSTANCE_ID/HOSTNAME

# Check recovery mode status
civo instance recovery-status INSTANCE_ID/HOSTNAME
```

The recovery-status command supports custom output formats with the following fields:
* id - The instance ID
* hostname - The instance hostname
* status - Current recovery mode status


### VNC Access

The VNC command allows you to access your instance through a browser-based VNC console.

```sh
# Open VNC console (default duration)
civo instance vnc INSTANCE_ID/HOSTNAME

# Open VNC console with custom duration
civo instance vnc INSTANCE_ID/HOSTNAME --duration 2h
```

The `--duration` flag accepts Go's duration format:
* "30m" - 30 minutes
* "1h" - 1 hour
* "24h" - 24 hours

When executed, the command will:
1. Enable VNC access for the specified instance
2. Generate a secure VNC URL
3. Automatically open the console in your default browser
4. Attempt to connect for up to 35 seconds before timing out


## Kubernetes clusters

#### Introduction

You can manage Kubernetes clusters on Civo using the Kubernetes subcommands.

#### List clusters

To see your created clusters, call `civo kubernetes list`:

```sh
$ civo kubernetes list
+--------------------------------------+----------------+--------+-------+-------+--------+
| ID                                   | Name           | Region | Nodes | Pools | Status |
+--------------------------------------+----------------+--------+-------+-------+--------+
| 5604340f-caa3-4ac1-adb7-40c863fe5639 | falling-sunset | NYC1   |     2 |     1 | ACTIVE |
+--------------------------------------+----------------+--------+-------+-------+--------+
```

#### Listing kubernetes sizes

You can list all kubernetes sizes by running `civo kubernetes size`.

```sh
$ civo kubernetes size
+--------------------+--------------------------------+------------+-----------+---------+--------+------------+
| Name               | Description                    | Type       | CPU Cores | RAM MB  | SSD GB | Selectable |
+--------------------+--------------------------------+------------+-----------+---------+--------+------------+
| g4s.kube.xsmall    | Extra Small - Standard         | Kubernetes |         1 |    1024 |     30 | Yes        |
+--------------------+--------------------------------+------------+-----------+---------+--------+------------+
| g4s.kube.small     | Small - Standard               | Kubernetes |         1 |    2048 |     40 | Yes        |
+--------------------+--------------------------------+------------+-----------+---------+--------+------------+
| g4s.kube.medium    | Medium - Standard              | Kubernetes |         2 |    4096 |     50 | Yes        |
+--------------------+--------------------------------+------------+-----------+---------+--------+------------+
| g4s.kube.large     | Large - Standard               | Kubernetes |         4 |    8192 |     60 | Yes        |
+--------------------+--------------------------------+------------+-----------+---------+--------+------------+
| g4p.kube.small     | Small - Performance            | Kubernetes |         4 |   16384 |     60 | Yes        |
+--------------------+--------------------------------+------------+-----------+---------+--------+------------+
| g4p.kube.medium    | Medium - Performance           | Kubernetes |         8 |   32768 |     80 | Yes        |
+--------------------+--------------------------------+------------+-----------+---------+--------+------------+
| g4p.kube.large     | Large - Performance            | Kubernetes |        16 |   65536 |    120 | Yes        |
+--------------------+--------------------------------+------------+-----------+---------+--------+------------+
| g4p.kube.xlarge    | Extra Large - Performance      | Kubernetes |        32 |  131072 |    180 | Yes        |
+--------------------+--------------------------------+------------+-----------+---------+--------+------------+
| g4c.kube.small     | Small - CPU optimized          | Kubernetes |         8 |   16384 |     60 | Yes        |
+--------------------+--------------------------------+------------+-----------+---------+--------+------------+
| g4c.kube.medium    | Medium - CPU optimized         | Kubernetes |        16 |   32768 |     80 | Yes        |
+--------------------+--------------------------------+------------+-----------+---------+--------+------------+
| g4c.kube.large     | Large - CPU optimized          | Kubernetes |        32 |   65536 |    120 | Yes        |
+--------------------+--------------------------------+------------+-----------+---------+--------+------------+
| g4c.kube.xlarge    | Extra Large - CPU optimized    | Kubernetes |        64 |  131072 |    180 | Yes        |
+--------------------+--------------------------------+------------+-----------+---------+--------+------------+
| g4m.kube.small     | Small - RAM optimized          | Kubernetes |         2 |   16384 |     60 | Yes        |
+--------------------+--------------------------------+------------+-----------+---------+--------+------------+
| g4m.kube.medium    | Medium - RAM optimized         | Kubernetes |         4 |   32768 |     80 | Yes        |
+--------------------+--------------------------------+------------+-----------+---------+--------+------------+
| g4m.kube.large     | Large - RAM optimized          | Kubernetes |         8 |   65536 |    120 | Yes        |
+--------------------+--------------------------------+------------+-----------+---------+--------+------------+
| g4m.kube.xlarge    | Extra Large - RAM optimized    | Kubernetes |        16 |  131072 |    180 | Yes        |
+--------------------+--------------------------------+------------+-----------+---------+--------+------------+
| g4g.kube.small     | Small - Nvidia A100 80GB       | Kubernetes |        12 |  131072 |    100 | Yes        |
+--------------------+--------------------------------+------------+-----------+---------+--------+------------+
| g4g.kube.medium    | Medium - Nvidia A100 80GB      | Kubernetes |        24 |  262144 |    100 | Yes        |
+--------------------+--------------------------------+------------+-----------+---------+--------+------------+
| g4g.kube.large     | Large - Nvidia A100 80GB       | Kubernetes |        48 |  524288 |    100 | Yes        |
+--------------------+--------------------------------+------------+-----------+---------+--------+------------+
| g4g.kube.xlarge    | Extra Large - Nvidia A100 80GB | Kubernetes |        96 | 1048576 |    100 | Yes        |
+--------------------+--------------------------------+------------+-----------+---------+--------+------------+
| g4g.40.kube.small  | Small - Nvidia A100 40GB       | Kubernetes |         8 |   65536 |    200 | Yes        |
+--------------------+--------------------------------+------------+-----------+---------+--------+------------+
| g4g.40.kube.medium | Medium - Nvidia A100 40GB      | Kubernetes |        16 |  131072 |    400 | Yes        |
+--------------------+--------------------------------+------------+-----------+---------+--------+------------+
| g4g.40.kube.large  | Large - Nvidia A100 40GB       | Kubernetes |        32 |  262133 |    400 | Yes        |
+--------------------+--------------------------------+------------+-----------+---------+--------+------------+
| g4g.40.kube.xlarge | Extra Large - Nvidia A100 40GB | Kubernetes |        64 |  524288 |    400 | Yes        |
+--------------------+--------------------------------+------------+-----------+---------+--------+------------+
| an.g1.l40s.kube.x1 | Small - Nvidia L40S 40GB       | Kubernetes |        12 |  131072 |    200 | Yes        |
+--------------------+--------------------------------+------------+-----------+---------+--------+------------+
| an.g1.l40s.kube.x2 | Medium - Nvidia L40S 40GB      | Kubernetes |        24 |  262133 |    200 | Yes        |
+--------------------+--------------------------------+------------+-----------+---------+--------+------------+
| an.g1.l40s.kube.x4 | Large - Nvidia L40S 40GB       | Kubernetes |        48 |  524288 |    400 | Yes        |
+--------------------+--------------------------------+------------+-----------+---------+--------+------------+
| an.g1.l40s.kube.x8 | Extra Large - Nvidia L40S 40GB | Kubernetes |        96 | 1048576 |    400 | Yes        |
+--------------------+--------------------------------+------------+-----------+---------+--------+------------+

```

#### Create a cluster

You can create a cluster by running `civo kubernetes create` with a cluster name parameter, as well as any options you provide:

```bash
-a, --applications string          optional, use names shown by running 'civo kubernetes applications ls'
-p, --cni-plugin string            optional, possible options: flannel,cilium. (default "flannel")
-c, --create-firewall              optional, create a firewall for the cluster with all open ports
-e, --existing-firewall string     optional, ID of existing firewall to use
-u, --firewall-rules string        optional, can be used if the --create-firewall flag is set, semicolon-separated list of ports to open (default "default")
-h, --help                         help for create
-m, --merge                        merge the config with existing kubeconfig if it already exists.
-t, --network string               the name of the network to use in the creation (default "default")
-n, --nodes int                    the number of nodes to create (the master also acts as a node). (default 3)
-r, --remove-applications string   optional, remove default application names shown by running  'civo kubernetes applications ls'
 --save                         save the config
-s, --size string                  the size of nodes to create. (default "g4s.kube.medium")
 --switch                       switch context to newly-created cluster
-v, --version string               the k3s version to use on the cluster. Defaults to the latest. Example - 'civo k3s create --version 1.21.2+k3s1' (default "latest")
-w, --wait                         a simple flag (e.g. --wait) that will cause the CLI to spin and wait for the cluster to be ACTIVE
```

*Note*

- The '--create-firewall' will open the ports 80,443 and 6443 in the firewall if '--firewall-rules' is not used.
- The '--create-firewall' and '--existing-firewall' flags are mutually exclusive. You can't use them together.
- The '--firewall-rules' flag need to be used with '--create-firewall'.
- The '--firewall-rules' flag can accept:
  - You can pass 'all' to open all ports.
  - An optional end port using 'start_port-end_port' format (e.g. 8000-8100)
  - An optional CIDR notation (e.g. 0.0.0.0/0)
  - When no CIDR notation is provided, the port will get 0.0.0.0/0 (open to public) as default CIDR notation
  - When a CIDR notation is provided without slash and number segment, it will default to /32
  - Within a rule, you can use comma separator for multiple ports to have same CIDR notation
  - To separate between rules, you can use semicolon symbol and wrap everything in double quotes (see below)
    So the following would all be valid:
  - "80,443,6443:0.0.0.0/0;8080:1.2.3.4" (open 80,443,6443 to public and 8080 just for 1.2.3.4/32)
  - "80,443,6443;6000-6500:4.4.4.4/24" (open 80,443,6443 to public and 6000 to 6500 just for 4.4.4.4/24)

```sh
$ civo kubernetes create my-first-cluster
Created Kubernetes cluster my-first-cluster
```

#### Adding pools the cluster

You can add more pools to your cluster live (obviously 1 is the minimum) while the cluster is running. It takes the name of the cluster (or the ID), the parameter of `--nodes` which is the new number of nodes to run and `--size` which is the size of the pool,
if `--node` and `--size` are not specified, the default values will be used.

```sh
civo kubernetes node-pool create my-first-cluster
The pool (8064c7) was added to the cluster (my-first-cluster)
```

#### Scaling pools in the cluster

You can scale a pool in your cluster, while the cluster is running. It takes the name of the cluster (or the ID), the pool ID (at least 6 characters), and `--nodes` which is the new number of nodes in the pool. You can get the pool ID details form `civo k3s show my-first-cluster`

```sh
civo kubernetes node-pool scale my-first-cluster pool-id --nodes 5
```

#### Deleting pools in the cluster

You can delete a pool in your cluster. It takes the name of the cluster (or the ID) and the pool ID (at least 6 characters)

```sh
civo kubernetes node-pool delete my-first-cluster pool-id
```

#### Recycling nodes in your cluster

If you need to recycle a particular node in your cluster for any reason, you can use the `recycle` command. This requires the name or ID of the cluster, and the name of the node you wish to recycle preceded by `--node=`:

```sh
$ civo k8s recycle my-first-cluster --node=k3s-my-first-cluster-5ae1561e-node-pool-a56f
The node (k3s-my-first-cluster-5ae1561e-node-pool-a56f) was recycled
```

_Note:_ When a node is recycled, it is fully deleted. The recycle command does not [drain](https://kubernetes.io/docs/tasks/administer-cluster/safely-drain-node/) a node, it simply deletes it before building a new node and attaching it to a cluster. It is intended for scenarios where the node itself develops an issue and must be replaced with a new one.

#### Viewing or Saving the cluster configuration

To output a cluster's configuration information, you can invoke `civo kubernetes config cluster-name`. This will output the `kubeconfig` file to the screen.

You can save a cluster's configuration to your local `~/.kube/config` file. This requires `kubectl` to be installed. Usage:

```sh
civo kubernetes config my-first-cluster -s
Saved config to ~/.kube/config
```

If you already have a `~/.kube/config` file, any cluster configuration that is saved will be _overwritten_ unless you also pass the `--merge` option. If you have multiple cluster configurations, merging allows you to switch contexts at will. If you prefer to save the configuration in another place, just use the parameter `--local-path` or `-p` and the path. If you use `--switch` the cli will automatically change the kubernetes context to the new cluster.

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

By default, `Traefik` is bundled in with `k3s`. If you want to set up a cluster without Traefik, you can use the `remove-applications` option in the creation command to start a cluster without it:

```sh
civo kubernetes create --remove-applications=Traefik-v2-nodeport --nodes=2 --wait
```

The command uses the application name as displayed by running `civo kubernetes applications ls`.

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
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| Name                        | Version           | Category     | Plans                                                                                      | Dependencies                      |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| acorn                       | 0.7.1             | management   | Disabled Auto-TLS, Enabled Auto-TLS                                                        |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| ambassador-edge-stack       | 3.8.0             | architecture |                                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| apisix-ingress-controller   | 1.6.0             | architecture |                                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| argo-rollouts               | v1.4.1            | ci_cd        |                                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| argo-workflows              | v3.0.3            | ci_cd        |                                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| argocd                      | v2.11.4           | ci_cd        |                                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| aspnet                      | 5.0.5             | management   |                                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| atmo                        | 0.2.2             | architecture |                                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| bitwarden-passwordless-dev  | 1.0.74            | management   |                                                                                            | cert-manager                      |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| bitwarden-unified           | 2024.1.0          | management   |                                                                                            | cert-manager                      |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| blackbox-exporter           | 5.8.2             | monitoring   |                                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| cerbos                      | 0.34.0            | security     |                                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| cert-manager                | v1.15.1           | architecture |                                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| chaos-mesh                  | (default)         | ci_cd        |                                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| civo-cluster-autoscaler     | v1.25.0           | management   |                                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| code-server                 | 3.11.1            | management   | 1GB, 2GB, 5GB                                                                              | cert-manager                      |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| dapr                        | 1.11.0            | architecture |                                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| devtron                     |                   | ci_cd        |                                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| dynamic-pv-scaler           | 0.1.0             | storage      |                                                                                            | prometheus-operator               |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| edp                         | 3.8.1             | ci_cd        | GitHub, GitLab                                                                             | argocd, tekton, traefik2-nodeport |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| enroute-onestep             | 0.8.0             | architecture |                                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| epinio                      | v1.11.1           | management   |                                                                                            | cert-manager                      |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| falco                       |              0.27 | security     |                                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| ferretdb                    | 1.22.0            | database     | 5GB, 10GB, 20GB                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| flagsmith                   | 0.39.0            | ci_cd        |                                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| flux                        | Latest            | ci_cd        |                                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| gatekeeper                  | 3.4.0             | security     |                                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| ghost                       | 5.87.1            | management   | 5GB, 10GB, 15GB                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| gimlet                      | v0.17.2           | management   |                                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| gitea                       | 1.12.5            | management   |                                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| gitlab                      | 16.2.3            | management   | Community edition, Enterprise edition                                                      | cert-manager                      |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| gloo-edge                   | latest            | architecture |                                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| haproxy                     |               1.5 | architecture |                                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| helm                        | 2.16.5            | management   |                                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| istio                       | multiple          | architecture | Istio Latest, Istio v1.10.1, Istio v1.9.5, Istio v1.8.6                                    |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| jenkins                     | 2.452.2           | ci_cd        | 5GB, 10GB, 20GB                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| joomla                      |               5.2 | management   | 5GB, 10GB, 20GB                                                                            | mariadb:5GB, cert-manager         |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| kafka                       | 2.7.0             | architecture |                                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| keda                        | 2.14.0            | architecture |                                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| keptn                       | 0.8.7             | ci_cd        |                                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| ketch                       | 0.7.0             | management   |                                                                                            | cert-manager                      |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| keycloak                    | 25.0.1            | security     |                                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| kong-ingress-controller     | 3.0.1             | architecture |                                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| kriten                      | 4.0.1             | architecture |                                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| kube-hunter                 | latest            | security     |                                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| kubeclarity                 | Latest            | security     |                                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| kubefirst                   | v2.4.6            | management   |                                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| kubeflow                    | 1.6.0-rc.1        | management   |                                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| kubenav                     | 3.1.0             | management   |                                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| kubernetes-dashboard        | v2.4.0            | management   |                                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| kubernetes-external-secrets | 0.9.14            | security     |                                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| kubesphere                  | v3.4.1            | management   |                                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| kubevela                    | 1.0.1             | management   |                                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| kubevious                   | v1.0.3            | management   |                                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| kubewarden                  | v1.14.0           | security     |                                                                                            | cert-manager                      |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| kyverno                     | 1.8.2             | security     |                                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| linkerd                     | Latest            | architecture | Linkerd Minimal, Linkerd & Jaeger, Linkerd with Dashboard, Linkerd with Dashboard & Jaeger |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| litmuschaos                 | 2.0.0             | ci_cd        |                                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| loki                        | v2.9.11           | monitoring   |                                                                                            | prometheus-operator               |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| longhorn                    | 1.3.2             | storage      |                                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| maesh                       | v1.4.5            | architecture |                                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| mariadb                     | 11.4.2            | database     | 5GB, 10GB, 20GB                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| metrics-server              | (default)         | architecture |                                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| minio                       | 2019-08-29        | storage      | 5GB, 10GB, 20GB                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| mongodb                     | 4.2.12            | database     | 5GB, 10GB, 20GB                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| netdata                     | Latest            | monitoring   |                                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| neuvector                   | latest            | security     |                                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| nextcloud                   | 29.0.3            | management   | 5GB, 10GB, 20GB                                                                            | mariadb:5GB, cert-manager         |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| nexus3                      | 3.30.1            | ci_cd        |                                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| nginx                       | latest            | architecture |                                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| ngrok-ingress-controller    | 0.12.1            | architecture |                                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| nodered                     | 1.2.7             | architecture | 5GB, 10GB, 20GB                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| openfaas                    | 0.18.0            | architecture |                                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| openobserve                 | v0.8.1            | monitoring   |                                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| otomi                       | 1.0.0             | management   |                                                                                            | metrics-server                    |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| paralus                     | 0.2.3             | security     |                                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| parca                       | v0.15.0           | monitoring   |                                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| percona-mysql               | 1.11.0            | database     | 10GB, 20GB, 50GB                                                                           |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| permission-manager          | 1.6.0             | management   |                                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| pmm                         | 2.36.0            | monitoring   |                                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| polaris                     | 7.3.2             | security     |                                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| portainer                   | 2.16.2            | management   | Community edition, Business edition                                                        |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| postgresql                  |              16.3 | database     | 5GB, 10GB, 20GB                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| projectsveltos              | v0.29.1           | management   |                                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| prometheus-operator         | 0.38.1            | monitoring   |                                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| pyroscope                   | 0.2.5             | monitoring   |                                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| rabbitmq                    | 3.13.4-management | architecture |                                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| rancher                     | 2.8.5             | management   |                                                                                            | cert-manager                      |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| redis                       | 7.2-alpine        | database     |                                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| rekor                       | 0.2.2             | architecture |                                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| reloader                    | v0.0.125          | architecture |                                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| rqlite                      | 8.26.6            | database     | 1GB, 5GB, 10GB, 3Replicas                                                                  |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| sealed-secrets              | v0.26.2           | architecture |                                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| selenium                    |                 4 | ci_cd        |                                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| shipa                       | 1.7.2             | management   |                                                                                            | cert-manager                      |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| siglens-oss                 | 0.1.22            | monitoring   |                                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| spinkube                    | v0.1.0            | management   |                                                                                            | cert-manager                      |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| system-upgrade-controller   | v0.6.2            | architecture |                                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| tekton                      | v0.50.5           | ci_cd        |                                                                                            | metrics-server                    |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| traefik2-loadbalancer       | 2.11.0            | architecture |                                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| traefik2-nodeport           | 2.9.4             | architecture |                                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| unifi-controller            | stable-6          | management   |                                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| uptime-kuma                 | 1.23.3            | monitoring   | 1GB, 2GB                                                                                   |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| vault                       | 0.24.0            | security     | Standalone, Dev, High-Availability                                                         |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| volcano                     | v1.9.0            | management   |                                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| weavescope                  | v1.13.2           | monitoring   |                                                                                            |                                   |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
| wordpress                   | 5.6.2             | management   | 5GB, 10GB, 20GB                                                                            | mariadb:5GB                       |
+-----------------------------+-------------------+--------------+--------------------------------------------------------------------------------------------+-----------------------------------+
```

#### Show Applications details when is installed in the cluster

This option will be allow you to see the post-install instruction of every app installed in the cluster

```sh
civo kubernetes application show Traefik apps-demo-cluster
```

the first parameter is for the name of the app and the second is the name of the cluster

#### Installing Applications Onto a New Cluster

To specify applications to install onto a new cluster, list them at cluster creation by specifying their `name` from the list above:

```sh
civo kubernetes create apps-demo-cluster --nodes=2  --applications=Redis,Linkerd
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

**Note:** Applications like `metrics-server` are installed by default on any new cluster you create. If you don't need to install the default applications, use the `--remove-applications` flag as below:

```sh
civo kubernetes create apps-demo-cluster --nodes=2 --applications=Redis,Linkerd --remove-applications=metrics-server
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

We don't offer registration of domains names, this is purely for hosting the DNS. If you're looking to buy a domain name, we recommend [LCN.com](https://www.lcn.com/) for their excellent friendly support and very competitive prices.

#### Set Up a New Domain

Any user can add a domain name (that has been registered elsewhere) to be managed by Civo.com. You should adjust the nameservers of your domain (through your registrar) to point to `ns0.civo.com` and `ns1.civo.com`.

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
The domain record called www with ID 4e181dde-bde8-4744-8984-067f957a7d59 was deleted
```

## Firewalls

#### Introduction

You can configure custom firewall rules for your instances using the Firewall component of Civo CLI. These are freely configurable, however customers should be careful to not lock out their own access to their instances. By default, all ports are closed for custom firewalls.

Firewalls can be configured with rules, and they can be made to apply to your chosen instance(s) with subsequent commands.

#### Configuring a New Firewall

To create a new Firewall, use `civo firewall create new_firewall_name`:

```sh
$ civo firewall create civocli_demo
Created a firewall called civocli_demo with ID ab2a25d7-edd4-4ecd-95c4-58cb6bc402de
```

You can also create a firewall without any default rules by using the flag `-r` or `--create-rules` set to `false`. In both cases, the usage is like:

```bash
civo firewall create new_firewall_name --create-rules=false

```

You will then be able to **configure rules** that allow connections to and from your instance by adding a new rule using `civo firewall rule create firewall_id` with the required and your choice of optional parameters, listed here and used in an example below:

```text
Options:
-a, --action string      the action of the rule can be allow or deny (default is allow) (default "allow")
-c, --cidr string        the CIDR of the rule you can use (e.g. -c 10.10.10.1/32,148.2.6.120/32) (default "0.0.0.0/0")
-d, --direction string   the direction of the rule can be ingress or egress (default is ingress) (default "ingress")
-e, --endport string     the end port of the rule
-h, --help               help for create
-l, --label string       a string that will be the displayed as the name/reference for this rule
-p, --protocol string    the protocol choice (TCP, UDP, ICMP) (default "TCP")
-s, --startport string   the start port of the rule
```

Example usage:

```sh
$ civo firewall rule create civocli_demo --startport=22 --direction=ingress --label='SSH access for CLI demo' -a allow
 New rule SSH access for CLI demo created

$ civo firewall rule list civocli_demo
+--------------------------------------+-----------+----------+------------+----------+--------+-----------+-------------------------+
| ID                                   | Direction | Protocol | Start Port | End Port | Action | Cidr      | Label                   |
+--------------------------------------+-----------+----------+------------+----------+--------+-----------+-------------------------+
| 74fff0d7-0ba4-497b-bbc1-83179b4e3b23 | ingress   | tcp      |         22 |       22 | allow  | 0.0.0.0/0 | SSH access for CLI demo |
+--------------------------------------+-----------+----------+------------+----------+-----------+----------------------------------+
```

You can see all active rules for a particular firewall by calling `civo firewall rule ls firewall_id`, where `firewall_id` is the UUID of your particular firewall.

#### Managing Firewalls

You can see an overview of your firewalls using `civo firewall list` showing you which firewalls have been configured with rules, and whether any of your instances are using a given firewall, such as in this case where the firewall we have just configured has the one rule, but no instances using it.

```sh
$ civo firewall list
+--------------------------------------+--------------------+---------+-------------+-----------------+----------------+--------------------+
| ID                                   | Name               | Network | Total rules | Total Instances | Total Clusters | Total LoadBalancer |
+--------------------------------------+--------------------+---------+-------------+-----------------+----------------+--------------------+
| 3ac0681d-2f71-4921-ae20-019272d9248b | Default (all open) | Default |           3 |               0 |              0 |                  0 |
+--------------------------------------+--------------+-------------+----------------+----------------+----------------+--------------------+
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

#### Listing Networks

To list all the networks you can run `civo network ls`

```sh
$ civo network ls
+--------------------------------------+----------+--------+---------+
| ID                                   | Label    | Region | Default |
+--------------------------------------+----------+--------+---------+
| 28244c7d-b1b9-48cf-9727-aebb3493aaac | Default  | LON1   | true    |
| fa21edfa-c089-421c-8008-0c0c7784386a | test     | LON1   | false   |
| 35ec87e7-fbd2-4ee8-849a-f88d7363e23f | cli-demo | LON1   | false   |
+--------------------------------------+----------+--------+---------+
```

#### Removing Networks

Removal of a network, provided you do not need it and your applications do not depend on routing through it, is simple - you call `civo network remove network_ID`:

```sh
$ civo network remove 74b69006-ea59-46a0-96c4-63f5bfa290e1
Removed the network cli-demo with ID 74b69006-ea59-46a0-96c4-63f5bfa290e1
```

## Object Stores

#### Introduction

Object stores are S3-compatible data storage structures on Civo. Through creating object stores in your account, you can manage unstructured data within the size limits of each object store and subject to your quota.

#### Listing Object Stores

You can run `civo objectstore ls` to get the list of all object stores in your account.

```console
$ civo objectstore ls
+--------+----------+------+---------------------------+--------+
| ID     | Name     | Size | Object Store Endpoint     | Status |
+--------+----------+------+---------------------------+--------+
| 699e42 | cli-demo |  500 | objectstore.lon1.civo.com | ready  |
| 083e59 | test     |  500 | objectstore.lon1.civo.com | ready  |
+--------+----------+------+---------------------------+--------+
```

#### Creating Object Stores

You can create an object store by running `civo objectstore create` with an object store name parameter. You can also specify the size by using the `--size` flag. If you don't specify the size, it will be set to 500GB by default. If you do not specify a `--region` flag, the object store will be created in your currently-active region.

```sh
$ civo objectstore create cli-demo

Creating Object Store cli-demo in LON1
To check the status of the Object Store run: civo objectstore show cli-demo
```

#### Updating Object Stores

Provided you have room in your Civo quota, you can update the size of an object store by running `civo objectstore update` with the object store name and the new size.

```console
$ civo objectstore update cli-demo --size=1000

The Object Store with ID 699e42a7-918b-42f7-ac22-fb9869e835ad was updated to size: 1000 GB
```

#### Deleting Object Stores

You can delete an object store by running `civo objectstore delete` with the object store name.

```console
$ civo objectstore delete cli-demo

Warning: Are you sure you want to delete the cli-demo Object Store (y/N) ? y
The Object Store (cli-demo) has been deleted
```

## Object Store Credentials

#### Introduction

Access to object stores is controlled by credentials management. When a new object store is created, a default administrative set of credentials is created with it. You can create other credentials for object stores you create, and export credential information for use in applications.

#### Listing Object Store Credentials

You can run `civo objectstore credential ls` to get the list of all object store credentials in your account.

```sh
$ civo objectstore credential ls

+----------------------+----------------------+--------+
| Name                 | Access Key           | Status |
+----------------------+----------------------+--------+
| cli-demo-xazd-7de9b2 | 3W-X-Y-X-Y-X-Y-X-Y-X | ready  |
| cli-demo-fa5d-fe8d50 | QW-X-Y-X-Y-X-Y-X-Y-N | ready  |
+----------------------+----------------------+--------+
```

#### Creating Object Store Credentials

You can create an object store credential by running `civo objectstore credential create` with a credential name parameter. You can also specify an access key and a secret access key with `----access-key` and `--secret-key` flags respectively. If no flag is provided, we'll generate both the keys for you.

```sh
$ civo objectstore credential create cli-demo-credential --access-key=YOUR_ACCESS_KEY --secret-key=YOUR_SECRET_KEY
Creating Object Store Credential cli-demo-credential in LON1
```

#### Exporting Credentials

You can export the credentials by running `civo objectstore credential export` with the access key for your credential. We support export in 2 formats - `env` and `s3cfg`. For more information on credential export formats, refer to the help text found at `civo objectstore credential export --help`.

If no `--format` flag is passed, we export it to `env` by default.

```sh
$ civo objectstore credential export --access-key=YOUR_ACCESS_KEY

# Tip: You can redirect output with (>> ~/.zshrc) to add these to Zsh's startup automatically
export AWS_ACCESS_KEY_ID=YOUR_ACCESS_KEY
export AWS_SECRET_ACCESS_KEY=YOUR_SECRET_KEY
export AWS_DEFAULT_REGION=LON1
export AWS_HOST=https://objectstore.lon1.civo.com
```

#### Deleting Object Store Credentials

You can delete an object store credential by running `civo objectstore credential delete` with the credential name as a parameter.

```sh
$ civo objectstore credential delete cli-demo-fa5d-7de9b2

Warning: Are you sure you want to delete the cli-demo-fa5d-7de9b2 Object Store Credential (y/N) ? y
The Object Store Credential (cli-demo-fa5d-7de9b2) has been deleted
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
+--------------------+--------------------------------+------------+-----+---------+-----+
| Name               | Description                    | Type       | CPU | RAM     | SSD |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g3.xsmall          | Extra Small                    | Instance   |   1 |    1024 |  25 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g3.small           | Small                          | Instance   |   1 |    2048 |  25 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g3.medium          | Medium                         | Instance   |   2 |    4096 |  50 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g3.large           | Large                          | Instance   |   4 |    8192 | 100 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g3.xlarge          | Extra Large                    | Instance   |   6 |   16384 | 150 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g3.2xlarge         | 2X Large                       | Instance   |   8 |   32768 | 200 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4s.kube.xsmall    | Extra Small - Standard         | Kubernetes |   1 |    1024 |  30 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4s.kube.small     | Small - Standard               | Kubernetes |   1 |    2048 |  40 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4s.kube.medium    | Medium - Standard              | Kubernetes |   2 |    4096 |  50 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4s.kube.large     | Large - Standard               | Kubernetes |   4 |    8192 |  60 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4p.kube.small     | Small - Performance            | Kubernetes |   4 |   16384 |  60 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4p.kube.medium    | Medium - Performance           | Kubernetes |   8 |   32768 |  80 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4p.kube.large     | Large - Performance            | Kubernetes |  16 |   65536 | 120 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4p.kube.xlarge    | Extra Large - Performance      | Kubernetes |  32 |  131072 | 180 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4c.kube.small     | Small - CPU optimized          | Kubernetes |   8 |   16384 |  60 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4c.kube.medium    | Medium - CPU optimized         | Kubernetes |  16 |   32768 |  80 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4c.kube.large     | Large - CPU optimized          | Kubernetes |  32 |   65536 | 120 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4c.kube.xlarge    | Extra Large - CPU optimized    | Kubernetes |  64 |  131072 | 180 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4m.kube.small     | Small - RAM optimized          | Kubernetes |   2 |   16384 |  60 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4m.kube.medium    | Medium - RAM optimized         | Kubernetes |   4 |   32768 |  80 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4m.kube.large     | Large - RAM optimized          | Kubernetes |   8 |   65536 | 120 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4m.kube.xlarge    | Extra Large - RAM optimized    | Kubernetes |  16 |  131072 | 180 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g3.db.xsmall       | Extra Small                    | Database   |   1 |    2048 |  20 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g3.db.small        | Small                          | Database   |   2 |    4096 |  40 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g3.db.medium       | Medium                         | Database   |   4 |    8192 |  80 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g3.db.large        | Large                          | Database   |   6 |   16384 | 160 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g3.db.xlarge       | Extra Large                    | Database   |   8 |   32768 | 320 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g3.db.2xlarge      | Double Extra Large             | Database   |  10 |   65536 | 640 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g3.kf.small        | Small - CPU optimized          | KfCluster  |   4 |   16384 |  60 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g3.kf.medium       | Medium - CPU optimized         | KfCluster  |   8 |   32768 |  80 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g3.kf.large        | Large - CPU optimized          | KfCluster  |  16 |   65536 | 120 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g3.kf.xlarge       | Extra Large - CPU optimized    | KfCluster  |  32 |  131072 | 180 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4s.xsmall         | xSmall - Standard              | Instance   |   1 |    1024 |  25 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4s.small          | Small - Standard               | Instance   |   1 |    2048 |  25 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4s.medium         | Medium - Standard              | Instance   |   2 |    4096 |  50 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4s.large          | Large - Standard               | Instance   |   4 |    8192 | 100 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4s.xlarge         | Extra Large - Standard         | Instance   |   6 |   16384 | 150 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4s.2xlarge        | 2X Large - Standard            | Instance   |   8 |   32768 | 200 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4g.kube.small     | Small - Nvidia A100 80GB       | Kubernetes |  12 |  131072 | 100 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4g.kube.medium    | Medium - Nvidia A100 80GB      | Kubernetes |  24 |  262144 | 100 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4g.kube.large     | Large - Nvidia A100 80GB       | Kubernetes |  48 |  524288 | 100 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4g.kube.xlarge    | Extra Large - Nvidia A100 80GB | Kubernetes |  96 | 1048576 | 100 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4g.small          | Small - Nvidia A100 80GB       | Instance   |  12 |  131072 | 200 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4g.medium         | Medium - Nvidia A100 80GB      | Instance   |  24 |  262144 | 200 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4g.large          | Large - Nvidia A100 80GB       | Instance   |  48 |  524288 | 400 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4g.xlarge         | Extra Large - Nvidia A100 80GB | Instance   |  96 | 1048576 | 400 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4g.40.small       | Small - Nvidia A100 40GB       | Instance   |   8 |   65536 | 200 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4g.40.medium      | Medium - Nvidia A100 40GB      | Instance   |  16 |  131072 | 400 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4g.40.large       | Large - Nvidia A100 40GB       | Instance   |  32 |  262133 | 400 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4g.40.xlarge      | Extra Large - Nvidia A100 40GB | Instance   |  64 |  524288 | 400 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4g.40.kube.small  | Small - Nvidia A100 40GB       | Kubernetes |   8 |   65536 | 200 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4g.40.kube.medium | Medium - Nvidia A100 40GB      | Kubernetes |  16 |  131072 | 400 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4g.40.kube.large  | Large - Nvidia A100 40GB       | Kubernetes |  32 |  262133 | 400 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4g.40.kube.xlarge | Extra Large - Nvidia A100 40GB | Kubernetes |  64 |  524288 | 400 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| an.g1.l40s.x1      | Small - Nvidia L40S 40GB       | Instance   |  12 |  131072 | 200 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| an.g1.l40s.x2      | Medium - Nvidia L40S 40GB      | Instance   |  24 |  262133 | 200 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| an.g1.l40s.x4      | Large - Nvidia L40S 40GB       | Instance   |  48 |  524288 | 400 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| an.g1.l40s.x8      | Extra Large - Nvidia L40S 40GB | Instance   |  96 | 1048576 | 400 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| an.g1.l40s.kube.x1 | Small - Nvidia L40S 40GB       | Kubernetes |  12 |  131072 | 200 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| an.g1.l40s.kube.x2 | Medium - Nvidia L40S 40GB      | Kubernetes |  24 |  262133 | 200 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| an.g1.l40s.kube.x4 | Large - Nvidia L40S 40GB       | Kubernetes |  48 |  524288 | 400 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| an.g1.l40s.kube.x8 | Extra Large - Nvidia L40S 40GB | Kubernetes |  96 | 1048576 | 400 |
+--------------------+--------------------------------+------------+-----+---------+-----+
```

This command is useful for getting the name of the instance type if you do not remember it - you will need to specify the instance size name when creating an instance using the CLI tool.

Also you can use `--filter` to filter the result by the type, the avalible option are (instance, kubernetes) like this:

```sh
$ civo sizes list --filter kubernetes
+--------------------+--------------------------------+------------+-----+---------+-----+
| Name               | Description                    | Type       | CPU | RAM     | SSD |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4s.kube.xsmall    | Extra Small - Standard         | Kubernetes |   1 |    1024 |  30 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4s.kube.small     | Small - Standard               | Kubernetes |   1 |    2048 |  40 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4s.kube.medium    | Medium - Standard              | Kubernetes |   2 |    4096 |  50 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4s.kube.large     | Large - Standard               | Kubernetes |   4 |    8192 |  60 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4p.kube.small     | Small - Performance            | Kubernetes |   4 |   16384 |  60 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4p.kube.medium    | Medium - Performance           | Kubernetes |   8 |   32768 |  80 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4p.kube.large     | Large - Performance            | Kubernetes |  16 |   65536 | 120 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4p.kube.xlarge    | Extra Large - Performance      | Kubernetes |  32 |  131072 | 180 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4c.kube.small     | Small - CPU optimized          | Kubernetes |   8 |   16384 |  60 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4c.kube.medium    | Medium - CPU optimized         | Kubernetes |  16 |   32768 |  80 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4c.kube.large     | Large - CPU optimized          | Kubernetes |  32 |   65536 | 120 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4c.kube.xlarge    | Extra Large - CPU optimized    | Kubernetes |  64 |  131072 | 180 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4m.kube.small     | Small - RAM optimized          | Kubernetes |   2 |   16384 |  60 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4m.kube.medium    | Medium - RAM optimized         | Kubernetes |   4 |   32768 |  80 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4m.kube.large     | Large - RAM optimized          | Kubernetes |   8 |   65536 | 120 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4m.kube.xlarge    | Extra Large - RAM optimized    | Kubernetes |  16 |  131072 | 180 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4g.kube.small     | Small - Nvidia A100 80GB       | Kubernetes |  12 |  131072 | 100 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4g.kube.medium    | Medium - Nvidia A100 80GB      | Kubernetes |  24 |  262144 | 100 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4g.kube.large     | Large - Nvidia A100 80GB       | Kubernetes |  48 |  524288 | 100 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4g.kube.xlarge    | Extra Large - Nvidia A100 80GB | Kubernetes |  96 | 1048576 | 100 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4g.40.kube.small  | Small - Nvidia A100 40GB       | Kubernetes |   8 |   65536 | 200 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4g.40.kube.medium | Medium - Nvidia A100 40GB      | Kubernetes |  16 |  131072 | 400 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4g.40.kube.large  | Large - Nvidia A100 40GB       | Kubernetes |  32 |  262133 | 400 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| g4g.40.kube.xlarge | Extra Large - Nvidia A100 40GB | Kubernetes |  64 |  524288 | 400 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| an.g1.l40s.kube.x1 | Small - Nvidia L40S 40GB       | Kubernetes |  12 |  131072 | 200 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| an.g1.l40s.kube.x2 | Medium - Nvidia L40S 40GB      | Kubernetes |  24 |  262133 | 200 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| an.g1.l40s.kube.x4 | Large - Nvidia L40S 40GB       | Kubernetes |  48 |  524288 | 400 |
+--------------------+--------------------------------+------------+-----+---------+-----+
| an.g1.l40s.kube.x8 | Extra Large - Nvidia L40S 40GB | Kubernetes |  96 | 1048576 | 400 |
+--------------------+--------------------------------+------------+-----+---------+-----+
```

## SSH Keys

#### Introduction

To manage the SSH keys for an account that are used to log in to cloud instances, the Civo CLI tool provides the following commands. You would need to [
generate a new key](https://www.civo.com/learn/ssh-key-basics) according to your particular circumstances, if you do not have a suitable SSH key yet.

#### Uploading a New SSH Key

You will need the path to your public SSH Key to upload a new key to Civo. The usage is as follows: `civo ssh create NAME --key /path/to/FILENAME`

#### Listing Your SSH Keys

You will be able to list the SSH keys known for the current account holder by invoking `civo ssh list`:

```sh
$ civo sshkeys ls
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

## Disk Image

#### Introduction

Civo instances are built from a disk image. Currently there centos, debian and ubuntu are supported.In order to create an instance the diskimage ID is needed that can be found by running `civo diskimage ls`

#### Listing Available Disk Images

```sh
$ civo diskimage ls
+--------------------------------------+-----------------+----------------+-----------+--------------+
| ID                                   | Name            | Version        | State     | Distribution |
+--------------------------------------+-----------------+----------------+-----------+--------------+
| 21613daa-a66b-44fc-87f5-b6db566d8f91 | ubuntu-cuda11-8 | 22.04-cuda11-8 | available | ubuntu       |
+--------------------------------------+-----------------+----------------+-----------+--------------+
| a4204155-a876-43fa-b4d6-ea2af8774560 | debian-10       |             10 | available | debian       |
+--------------------------------------+-----------------+----------------+-----------+--------------+
| eda67ea0-4282-4945-9b7b-d3e1cba1d987 | ubuntu-jammy    |          22.04 | available | ubuntu       |
+--------------------------------------+-----------------+----------------+-----------+--------------+
| d927ad2f-5073-4ed6-b2eb-b8e61aef29a8 | ubuntu-focal    |          20.04 | available | ubuntu       |
+--------------------------------------+-----------------+----------------+-----------+--------------+
| 13232803-0928-4634-9ab8-476bff29ef1b | ubuntu-cuda12-2 | 22.04-cuda12-2 | available | ubuntu       |
+--------------------------------------+-----------------+----------------+-----------+--------------+
| ffb6fd93-cb06-4e8d-8058-46003b78e2ff | talos-v1.2.8    | 1.25.5         | available | talos        |
+--------------------------------------+-----------------+----------------+-----------+--------------+
| 9a16a77e-1a1f-45c8-87fd-6d1a19eeaac9 | talos-v1.5.0    | 1.27.0         | available | talos        |
+--------------------------------------+-----------------+----------------+-----------+--------------+
| 25fbbd96-d5ec-4d08-9c75-a5e154dabf9b | debian-11       |             11 | available | debian       |
+--------------------------------------+-----------------+----------------+-----------+--------------+
| 170db96f-8458-44aa-83ca-0c31fb81a835 | rocky-9-1       |            9.1 | available | rocky        |
+--------------------------------------+-----------------+----------------+-----------+--------------+
| 9b661c46-ac4f-46e1-9f3d-aaacde9b4fec | debian-9        |              9 | available | debian       |
+--------------------------------------+-----------------+----------------+-----------+--------------+
```

## Volumes

#### Introduction

Volumes are flexible-size additional storage for instances. By creating and associating a Volume with an instance, an additional virtual disk will be made available for backups or database files that can then moved to another instance.

Volumes for Kubernetes clusters do not have to be created – a cluster object for a PersistentVolume and PersistentVolumeClaim will create the volume in your account.

Volumes take disk space on your account's quota, and can only be created up to this quota limit. For more information about the quota system, see [Quota](#quota).

#### Creating a Volume

You can create a new volume by calling `civo volume create NAME SIZE(GB)`:

```text
Options:
  -h, --help             help for create
  -t, --network string   The network name/ID where the volume will be created (default "default")
  -s, --size-gb int      The new size in GB (required)
```

```sh
$ civo volume create CLI-demo-volume -s 25
Created a volume called CLI-demo-volume with ID 59076ec8-edba-4071-80d0-e9cfcce37b12
```

#### Attaching a Volume to an Instance

Mounting (Attaching) a volume onto an instance will allow that instance to use the volume as a drive:

```sh
$ civo volume attach CLI-demo-volume api-demo.test
The volume called CLI-demo-volume with ID 59076ec8-edba-4071-80d0-e9cfcce37b12 was attached to the instance api-demo.test
```

If this is a newly-created volume, you would need to partition, format and mount the volume. For more information, [see the Learn guide here](https://www.civo.com/learn/configuring-block-storage-on-civo).
Note: You can only attach a volume to one instance at a time.

#### Detaching a Volume From an Instance

If you want to detach a volume to move it to another instance, or are just finished with it, you can detach it once it's been [unmounted](https://www.civo.com/learn/configuring-block-storage-on-civo) using `civo volume detach volume_id`:

```sh
$ civo volume detach CLI-demo-volume
The volume called CLI-demo-volume with ID 59076ec8-edba-4071-80d0-e9cfcce37b12 was detached
```

#### Listing Volumes

You can get an overall view of your volumes, their sizes and status by using `civo volume list`:

```sh
$ civo volume ls
+--------------------------------------+------------------------------------------+---------+----------+----------+-------+-------------+-----------+
| ID                                   | Name                                     | Network | Cluster  | Instance | Size  | Mount Point | Status    |
+--------------------------------------+------------------------------------------+---------+----------+----------+-------+-------------+-----------+
| 59076ec8-edba-4071-80d0-e9cfcce37b12 | CLI-demo-volume                          | Default |          |          | 25 GB |             | available |
+--------------------------------------+------------------------------------------+---------+----------+----------+-------+-------------+-----------+
```

#### Deleting Volumes

To free up quota and therefore the amount to be billed to your account, you can delete a volume through `civo volume delete volume_id`. This deletion is immediate:

```sh
$ civo volume delete CLI-demo-volume
The volume called CLI-demo-volume with ID 59076ec8-edba-4071-80d0-e9cfcce37b12 was deleted
```

If a Kubernetes volume is showing with a status of `dangling` it can be deleted to release the quota and prevent further billing by running `civo volume delete <VOLUME-NAME> --region <REGION-NAME>`.

## Teams

Teams are a grouping of users, each member of a team having one or more permissions, or roles. When a user logs in, they don't have to select which team to use - only which account they want to act within. The permissions available are the total set of permissions they have across the teams in that account, combined.

#### List all teams

You can run `civo teams ls` to get the list of all teams

```sh
$ civo teams ls
+--------------------------------------+------------------+
| ID                                   | Name             |
+--------------------------------------+------------------+
| 39a755c3-87b4-4d7b-9b05-a8dfb8672081 | Support          |
| 3f233c4f-b931-44f2-9417-acf79b7f1a2d | Developers       |
| 5fdff566-a77a-4343-9c41-8ad714ec8593 | Management       |
| ad54-4614-8740-8604216b7ad4          | Technical Staff  |
| d1bc7dd6-2930-4c8d-b9df-0b9b2bfbf68c | Site Reliability |
| e7ec2649-ad54-4614-8740-8604216b7ad4 | Owners           |
+--------------------------------------+------------------+
```

#### Create a new team

To create a new team in your account, the cmd you need to run is `civo teams create <NEW-TEAM-NAME>` and a new team will be created with the given name.

```sh
$ civo teams create Community
Created a team called Community with team ID 475a087b-bec8-4a66-ac14-95bc09bd8d1e
```

#### Rename a team

To rename a team, you need to run the cmd `civo teams rename <OLD-TEAM-NAME> <NEW-TEAM-NAME>`

```sh
$ civo teams rename Community Advocacy
The team with ID 475a087b-bec8-4a66-ac14-95bc09bd8d1e was renamed to Advocacy
```

#### Delete a team

To delete a team, you need to run the cmd `civo teams delete <TEAM-NAME>`

```sh
$ civo teams delete Advocacy
Warning: Are you sure you want to delete the Advocacy team  (y/N) ? y
The team (Advocacy) has been deleted
```

## Permissions

Each member of a team is assigned one or more permissions, or roles. The permissions available are the total set of permissions they have across the teams in that account, combined.

#### List all permissions

You have to run the cmd `civo permissions ls` to list down all the available permissions.

```sh
$ civo permissions ls
+-------------------------+----------------------------+------------------------------------------------------------------------------------------------+
| Code                    | Name                       | Description                                                                                    |
+-------------------------+----------------------------+------------------------------------------------------------------------------------------------+
| *.*                     | Owner                      | Can perform any action                                                                         |
| organisation.owner      | Organisation Owner         | Can administer organisation details                                                            |
| billing.*               | Billing Administrator      | Can view/change billing details and see invoices                                               |
| billing.details_viewer  | Billing Details Viewer     | Can view billing details                                                                       |
| billing.invoices_viewer | Billing Invoices Viewer    | Can view invoices                                                                              |
| billing.details_changer | Billing Details Changer    | Can change billing details                                                                     |
| team.*                  | Team Administrator         | Can administer teams and their members                                                         |
| team.viewer             | Team Viewer                | Can view existing teams and their members                                                      |
| team.creater            | Team Creater               | Can create new teams, as well as add/remove team members                                       |
| team.updater            | Team Updater               | Can update team details, as well as add/remove team members                                    |
| team.deleter            | Team Deleter               | Can remove teams (and therefore all team members)                                              |
| kubernetes.*            | Kubernetes Administrator   | Can view/change Kubernetes clusters                                                            |
| kubernetes.viewer       | Kubernetes Viewer          | Can view existing Kubernetes clusters                                                          |
| kubernetes.creater      | Kubernetes Creater         | Can create new Kubernetes clusters                                                             |
| kubernetes.updater      | Kubernetes Updater         | Can make changes to existing Kubernetes clusters, such as scaling or marketplace installations |
| kubernetes.deleter      | Kubernetes Deleter         | Can delete Kubernetes clusters                                                                 |
| firewall.*              | Firewall Administrator     | Can view/change firewalls                                                                      |
| firewall.viewer         | Firewall Viewer            | Can view existing firewalls and their rules                                                    |
| firewall.creater        | Firewall Creater           | Can create new firewalls                                                                       |
| firewall.updater        | Firewall Updater           | Can update existing firewalls and change their rules                                           |
| firewall.deleter        | Firewall Deleter           | Can delete firewalls                                                                           |
| compute.*               | Compute Administrator      | Can view/change Compute instances                                                              |
| compute.viewer          | Compute Viewer             | Can view existing Compute instances                                                            |
| compute.creater         | Compute Creater            | Can create new Compute instances                                                               |
| compute.updater         | Compute Updater            | Can update (e.g. scale, rename) existing Compute instances                                     |
| compute.deleter         | Compute Deleter            | Can delete Compute instances                                                                   |
| loadbalancer.*          | Loadbalancer Administrator | Can view/change Loadbalancers                                                                  |
| loadbalancer.viewer     | Loadbalancer Viewer        | Can view existing Loadbalancers                                                                |
| loadbalancer.creater    | Loadbalancer Creater       | Can create new Loadbalancers                                                                   |
| loadbalancer.updater    | Loadbalancer Updater       | Can update (e.g. scale, rename) existing Loadbalancers                                         |
| loadbalancer.deleter    | Loadbalancer Deleter       | Can delete Loadbalancers                                                                       |
| dns.*                   | DNS Administrator          | Can view/change DNS domain names and records                                                   |
| dns.viewer              | DNS Viewer                 | Can view existing DNS domain names and records                                                 |
| dns.creater             | DNS Creater                | Can create new DNS domain names and records                                                    |
| dns.updater             | DNS Updater                | Can update existing DNS domain names and records                                               |
| dns.deleter             | DNS Deleter                | Can delete DNS domain names and records                                                        |
| network.*               | Network Administrator      | Can view/change networks                                                                       |
| network.viewer          | Network Viewer             | Can view existing networks                                                                     |
| network.creater         | Network Creater            | Can create new networks                                                                        |
| network.updater         | Network Updater            | Can rename existing networks                                                                   |
| network.deleter         | Network Deleter            | Can delete networks                                                                            |
| quota.manager           | Quota Manager              | Can request changes to the account quota                                                       |
| volume.*                | Volume Administrator       | Can view/change volumes                                                                        |
| volume.viewer           | Volume Viewer              | Can view volumes                                                                               |
| volume.creater          | Volume Creater             | Can create new volumes                                                                         |
| volume.updater          | Volume Updater             | Can change existing volumes (e.g. manage attachments)                                          |
| volume.deleter          | Volume Deleter             | Can delete volumes                                                                             |
| sshkey.*                | SSH Keys Administrator     | Can view/change SSH keys                                                                       |
| sshkey.viewer           | SSH Keys Viewer            | Can view existing SSH keys                                                                     |
| sshkey.creater          | SSH Keys Creater           | Can upload new SSH keys                                                                        |
| sshkey.updater          | SSH Keys Updater           | Can change existing SSH keys                                                                   |
| sshkey.deleter          | SSH Keys Deleter           | Can delete SSH keys                                                                            |
| webhook.*               | Webhook Administrator      | Can view/change webhooks                                                                       |
| webhook.viewer          | Webhook Viewer             | Can view existing webhooks                                                                     |
| webhook.creater         | Webhook Creater            | Can create new webhooks                                                                        |
| webhook.updater         | Webhook Updater            | Can change existing webhooks                                                                   |
| webhook.deleter         | Webhook Deleter            | Can delete webhooks                                                                            |
+-------------------------+----------------------------+------------------------------------------------------------------------------------------------+
```

## Region

As Civo grows, more regions for your instances will become available. You can run `civo region ls` to list the regions available. Block storage (Volumes) is region-specific, so if you configure an instance in one region, any volumes you wish to attach to that instance would have to be in the same region.

#### List all region

You can run `civo region ls` to get the list of all region

```sh
civo region ls
+--------+--------+----------------+---------+
| Code   | Name   | Country        | Current |
+--------+--------+----------------+---------+
| nyc1   | nyc1   | United States  |         |
+--------+--------+----------------+---------+
| phx1   | phx1   | United States  |         |
+--------+--------+----------------+---------+
| dg-exm | dg-exm | United Kingdom |         |
+--------+--------+----------------+---------+
| fra1   | fra1   | Germany        |         |
+--------+--------+----------------+---------+
| lon1   | lon1   | United Kingdom | <=====  |
+--------+--------+----------------+---------+
```

#### Change region

To change the region the only cmd you need run is `civo region current <REGION-CODE>` and you will see a message like this:

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

-   Source the completion script in your `~/.bashrc` file:

    ```shell
    echo 'source <(civo completion bash)' >>~/.bashrc
    ```

-   Add the completion script to the `/etc/bash_completion.d` directory:

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

Bug reports and pull requests are welcome on GitHub at <https://github.com/civo/cli>.

## License

The code is available as open source under the terms of the [Apache License 2.0](https://opensource.org/licenses/Apache-2.0).

## Thanks to all the contributors ❤

 <a href = "https://github.com/civo/cli/graphs/contributors">
   <img src = "https://contrib.rocks/image?repo=civo/cli"/>
 </a>
