#!/bin/bash
source /vagrant/functions.h

# Modify OS Values 
log "Upgrading the semaphore values"
echo "kernel.sem = 250 512000 100 2048" >> /etc/sysctl.conf > /dev/null
log "Reloading sysctl"
sysctl -p

# Manage User Accounts
log "Setting root password"
echo "root:changeme" | chpasswd

log "Creating gpadmin user"
useradd gpadmin

log "Setting gpadmin password"
echo "gpadmin:changeme" | chpasswd

log "Adding gpadmin to wheel group"
usermod -a -G wheel gpadmin

log "Allow passwordless sudo for wheel group"
cp /etc/sudoers /etc/sudoers.bak
sed --in-place 's/^#\s*\(%wheel\s\+ALL=(ALL)\s\+NOPASSWD:\s\+ALL\)/\1/' /etc/sudoers

log "Copy the bashrc to gpadmin user folder"
cp ~/.bashrc /home/gpadmin

log "Updating the vagrant user bashrc to auto login to gpadmin"
if [ "${2}" == "hack" ]; then
    log "Requested to auto login to gpadmin during vagrant ssh"
    {
        echo "sudo su - gpadmin"
        echo "exit"
    } >> /home/vagrant/.bashrc
else
     log "Skipping the step to auto login as per the request"
fi

# Manage Software
log "Cleaning RPM cache"
sed -i 's/gpgcheck=1/gpgcheck=0/g' /etc/yum.repos.d/*
sudo yum -q -e 0 clean all

log "Installing RPMs"
sudo yum -y -q -e 0 install ed unzip tar git strace gdb vim-enhanced wget m4 > /tmp/yum.out

log "Changing the permission of /usr/local"
chmod 777 /usr/local

log "Create /data/ directory"
mkdir -p /data
chown gpadmin:gpadmin /data

log "Running deployment script: go.build.sh"
source /vagrant/go.build.sh

log "The vagrant setup is complete"
