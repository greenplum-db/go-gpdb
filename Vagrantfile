# -*- mode: ruby -*-
# vi: set ft=ruby :

#load 'UAA.token'

Vagrant.configure(2) do |config|

    VAGRANT_COMMAND = ARGV[0]
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

   config.vm.provision "shell", path: 'scripts/os.prep.sh'
   config.vm.provision "shell", inline: 'cd /vagrant && git reset UAA.token' 
   #config.vm.provision "shell", path: 'scripts/go.build.sh', run: "always"
   
end