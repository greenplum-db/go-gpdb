# -*- mode: ruby -*-
# vi: set ft=ruby :

VAGRANT_COMMAND = ARGV[0]

GO_VERSION = "1.7.4"
OS = "linux"
ARCH = "amd64"

Vagrant.configure(2) do |config|

    if VAGRANT_COMMAND == "ssh"
        config.ssh.username = 'gpadmin'
    end

    # You can find all the vagrant boxes at location here
    # https://app.vagrantup.com/bento/

    config.vm.box = "bento/centos-7.4"

    # Make sure the private network is attached to the network adapter
    # vboxnet*, if you don't see vboxnet adapters on your workstation
    # then read instruction on the README
    # https://github.com/ielizaga/piv-go-gpdb
    # It is needed for command center GUI and for transferring files b/w machines

    config.vm.define "gpdb" do |node|
        node.vm.hostname = "gpdb"
        node.vm.network "private_network", ip: "192.168.99.100", name: "vboxnet0"
        node.vm.provider "virtualbox" do |vb|
            vb.name = "gpdb"
            vb.memory = "8196"
        end
   end

   # You can obtain the API key after login to pivotal network website
   # and on the edit profile section, more information on the repo readme
   # at https://github.com/ielizaga/piv-go-gpdb

   API_KEY = "c802dd9f43274a0b8a9a3c2ef106fdc1-r"

   # The below line, run the script to setup the script as per the system
   # requirements to run gpdb.

   config.vm.provision "shell", path: 'scripts/os.prep.sh', args: [GO_VERSION,OS,ARCH]

end