# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure(2) do |config|

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

   api_key = "c802dd9f43274a0b8a9a3c2ef106fdc1-r"

   # use "hack" or "nohack" as second parameter if you wish to avoid
   # login to gpadmin directly after doing vagrant ssh

   hack_or_nohack = "hack"

   # The below line, run the script to setup the script as per the system
   # requirements to run gpdb.

   config.vm.provision :shell, :path => 'sysprep.sh', :args => [api_key, hack_or_nohack]

end
