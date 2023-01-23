# Introduction [![Go Version](https://img.shields.io/badge/go-v1.19.4-green.svg?style=flat-square)](https://golang.org/dl/) [![MIT License](https://img.shields.io/badge/License-MIT-red.svg?style=flat-square)](https://github.com/ielizaga/piv-go-gpdb/blob/master/LICENSE)

This repo helps in making it ease to setup greenplum database on your local environment. You can provision single node deployment or multi node deployment cluster and test the functionality of the product

This repository is split into two parts

+ Datalab cli 
    > The datalab cli helps to create & manage vagrant VM's provisioning.
+ Gpdb cli
    > The gpdb cli helps to download, install, remove and manage the software of GPDB / GPCC.)

If you are interested in creating demo database / table and loading fake data to those table for testing purpose, check out the [project mock-data](https://github.com/pivotal-gss/mock-data)
    
Table of Contents
=================

   * [Table of Contents](#table-of-contents)
   * [Prerequisite](#prerequisite)
        * [VirtualBox](#virtualbox)
        * [Vagrant](#vagrant)
   * [Installation](#installation)
   * [Tools](#tools)
        * [Gpdb cli](#gpdb-cli)
        * [Datalab cli](#datalab-cli)
        * [UnitTest](#unittest)
        * [Vagrant](#vagrant-1)
   * [Demo (installation)](#demo-installation)
   * [Developers / Contributor's](#developers--contributors)


# Prerequisite

The "go-gpdb" software needs the below two tools pre-installed on your machine for it to work. Please follow the below instruction on how to setup the prerequisite

### VirtualBox

+ If not already installed, [download](http://download.virtualbox.org/virtualbox/5.1.22/VirtualBox-5.1.22-115126-OSX.dmg) and install VirtualBox or you can use brew to install virtual box using the command
    ```
    brew update
    brew cask install virtualbox
    ```
+ Once you virtualbox installation is complete ensure you have two interfaces (namely vboxnet0/1 is seen) on your MAC,
    ```
    IRFALI123:Vagrant fali$ ifconfig
    vboxnet0: flags=8842<BROADCAST,RUNNING,SIMPLEX,MULTICAST> mtu 1500
        ether 0a:00:27:00:00:00
    vboxnet1: flags=8943<UP,BROADCAST,RUNNING,PROMISC,SIMPLEX,MULTICAST> mtu 1500
        ether 0a:00:27:00:00:01
        inet 192.168.11.1 netmask 0xffffff00 broadcast 192.168.11.255
    ```
    if its not shown then refer to the [link](http://islandora.ca/content/fixing-missing-vboxnet0) for Virtualbox version lower than 5 and if your virtualbox version is 5 and above follow this [link](https://luppeng.wordpress.com/2017/07/17/enabling-virtualbox-host-only-adapter-on-mac-os-x/) on how to set those two interfaces up

### Vagrant

+ On your MAC install vagrant using the below command ( if vagrant executable is not already installed )
    ```
    brew update
    brew cask install vagrant
    ```
    If you have already installed vagrant ensure you are running the latest version of vagrant, to update your vagrant run
    ```
    brew update
    brew cask reinstall vagrant
    ```
+ Once the vagrant is installed, install the vagrant plugin using the below command

    ```
    vagrant plugin install vagrant-hosts
    ```

# Installation

If you have your own VM and would like to use that, then follow [this instruction](https://github.com/pivotal-gss/go-gpdb/blob/master/gpdb/INSTALL.md) on how to setup and use the go-gpdb. If you like to also setup a VM and also install gpdb then follow the set of instruction below

+ Ensure you have completed all the [Prerequisite](https://github.com/pivotal-gss/go-gpdb#prerequisite) as mentioned on the [Prerequisite](https://github.com/pivotal-gss/go-gpdb#prerequisite) section above.
+ Download the latest version of the [datalab CLI](https://github.com/pivotal-gss/go-gpdb/releases/latest) and follow the steps mentioned on the [installation](https://github.com/pivotal-gss/go-gpdb/tree/master/datalab#installation) and [setup](https://github.com/pivotal-gss/go-gpdb/tree/master/datalab#setup) section to setup the datalab CLI
+ Provision the VM's using the steps mentioned [here](https://github.com/pivotal-gss/go-gpdb/tree/master/datalab#create)
+ Once successfully provisioned, you can now start using the gpdb cli to download, install & manage gpdb software, refer the [example](https://github.com/pivotal-gss/go-gpdb/tree/master/gpdb#example) for more information. **NOTE:** gpdb cli is automaticalled installed on the provision VM's

# Tools
 
### Gpdb cli

+ Please check the gpdb cli [README](https://github.com/pivotal-gss/go-gpdb/blob/master/gpdb/README.md) for details on how to install & use the gpdb cli.

### Datalab cli

+ Please check the datalab cli [README](https://github.com/pivotal-gss/go-gpdb/blob/master/datalab/README.md) for details on how to install & use the datalab cli.

### UnitTest

+ Please check the [README](https://github.com/pivotal-gss/go-gpdb/blob/master/test/README.md) on how to run the unit test case.

### Vagrant 

+ If you wish to install vagrant manually using the Vagrant file & don't want to use the datalab cli, please follow the instruction mentioned [here](https://github.com/pivotal-gss/go-gpdb/tree/master/scripts/README.md) for all the options.

# Demo (installation)

 [![asciicast](https://asciinema.org/a/riz2i1xjPC0GNo0jso8V93JST.svg)](https://asciinema.org/a/riz2i1xjPC0GNo0jso8V93JST)
 
# Developers / Contributor's

Please read the [section](https://github.com/pivotal-gss/go-gpdb/blob/master/DEVLEOPER.md) on how to setup the environment to test and hack this tool
