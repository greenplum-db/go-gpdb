# -*- mode: ruby -*-
# vi: set ft=ruby :

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

    _host = ENV['HOST'] || "go-gpdb"
    _ip = ENV['IP'] || "192.168.99.100"

    config.vm.define _host do |node|
        node.vm.hostname = _host
        node.vm.network "private_network", ip: _ip, name: "vboxnet0"

        node.vm.provision :hosts do |provisioner|
           provisioner.autoconfigure = true
           provisioner.sync_hosts = true
           provisioner.add_localhost_hostnames = false
         end
         
        node.vm.provider :virtualbox do |vb|
            vb.name = _host
            vb.memory = "8196"
        end
   end
   
   config.vm.provision :hosts
   config.vm.provision :shell, path: 'scripts/os.prep.sh'
   config.vm.provision :shell, inline: "cd /vagrant && git update-index --assume-unchanged UAA.token"
   #config.vm.provision "shell", path: 'scripts/go.build.sh', run: "always"
   
end