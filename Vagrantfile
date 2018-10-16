# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure(2) do |config|

    HN  = ENV['HN'] || "go-gpdb"
    IP  = ENV['IP'] || "192.168.99.100"
    API  = ENV['API'] || ""
    
    # You can find all the vagrant boxes at location here
    # https://app.vagrantup.com/bento/

    config.vm.box = "bento/centos-7.4"

    # Make sure the private network is attached to the network adapter
    # vboxnet*, if you don't see vboxnet adapters on your workstation
    # then read instruction on the README
    # https://github.com/ielizaga/piv-go-gpdb
    # It is needed for command center GUI and for transferring files b/w machines

    config.vm.define HN do |node|
        node.vm.hostname = HN
        node.vm.network "private_network", ip: IP, name: "vboxnet0"

        node.vm.provision :hosts do |provisioner|
           provisioner.autoconfigure = true
           provisioner.sync_hosts = true
           provisioner.add_localhost_hostnames = false
         end
         
        node.vm.provider :virtualbox do |vb|
            vb.name = HN
            vb.memory = "8196"
        end
   end

# Optional --
# If the .vagrant/machines/<hostname> folder is empty -- or we request provisioning, then prompt for the API Key

#   if Dir.glob("#{File.dirname(__FILE__)}/.vagrant/machines/#{HN}/*").empty? || ARGV[1] == '--provision'
#     puts "Enter Your Pivnet UAA Token."
#     puts "A NULL Entry Will Prompt When Requried."
#     print "Pivnet API: "
#     API = STDIN.gets.chomp
#     print API
#     print "\n"
     config.vm.provision :hosts
     config.vm.provision :shell, path: 'scripts/os.prep.sh', args: [API]
     config.vm.provision :shell, path: 'scripts/go.build.sh'
#   end

  if ARGV[0] == "ssh"
      config.ssh.username = 'gpadmin'
  end

end