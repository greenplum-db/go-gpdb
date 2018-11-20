# -*- mode: ruby -*-
# vi: set ft=ruby :

# Includes
require "ipaddr"

# EVN PIVNET DEFAULTS
@pivnet_token = ENV['UAA_API_TOKEN'] || ""

# EVN VM DEFAULTS
@vm_cpus = ENV['VM_CPUS'].to_i || 2
@vm_memory = ENV['VM_MEMORY'].to_i || 2048

# ENV APPLICATION DEFAULTS
@subnet = ENV['GPDB_SUBNET'] || "192.168.99.100"
@hostname = ENV['GPDB_HOSTNAME'] || "go-gpdb"
@segments = ENV['GPDB_SEGMENTS'].to_i || 0

# Define a Template for Building All Our VMs.
def build_vm( config, hostname, ip )
    config.vm.define hostname do |node|
      node.vm.hostname = hostname
      node.vm.network :private_network, :ip => ip

      node.vm.provision :hosts do |provisioner|
         provisioner.autoconfigure = true
         provisioner.sync_hosts = true
         provisioner.add_localhost_hostnames = false
       end

      node.vm.provider :virtualbox do |vb|
          vb.name = hostname
          vb.gui = false
          vb.cpus = 2
          vb.memory = "2048"
      end
    end
end

# All Vagrant configuration is done below. 
Vagrant.configure("2") do |config|

  @ip = IPAddr.new @subnet

  # You can find all the vagrant boxes at location here
  # https://app.vagrantup.com/bento/
  config.vm.box = "bento/centos-7.4"
  
  # If "vagrant ssh", login As gpadmin, without hacking the vagrant profile
  if ARGV[0] == "ssh"
      config.ssh.username = 'gpadmin'
  end

  # If the token is empty and we are deploying, prompt for a token...
  if (@pivnet_token.to_s.empty?) && (ARGV[0] == 'up')
    puts "","UAA_API_TOKEN Environment Variable Not Found..."
    puts "A NULL entry will deploy, but a token is required for authentication.",""
    print "UAA API TOKEN: "
    @pivnet_token = STDIN.gets.chomp
  end

  puts "","UAA_API_TOKEN: #{@pivnet_token}"
  puts "Master Hostname: #{@hostname}"
  puts "Master IP: #{@ip}",""

  # Master Node:
  build_vm( config, @hostname, "#{@ip}" )

  if (ARGV[0] == 'up')
  puts "Segments Hosts: #{@segments}"
  end 
  
  # Segment Nodes:
  (1..@segments).each do |i|
    @ip = @ip.succ
    puts "[#{i}] Segment Hostname: #{@hostname}-#{i}"
    puts "[#{i}] Segment IP: #{@ip}",""
    build_vm( config, "#{@hostname}-#{i}", "#{@ip}")
  end

  # Provisioning 
  config.vm.provision :shell, path: 'scripts/os.prep.sh', args: [@pivnet_token]

  # Developer Mode
  # config.vm.provision :shell, path: 'scripts/go.build.sh'

end