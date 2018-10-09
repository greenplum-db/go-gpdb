#!/bin/bash
source /vagrant/scripts/functions.h

# Modify OS Parameters
log "Upgrading the semaphore values"
echo "kernel.sem = 250 512000 100 2048" >> /etc/sysctl.conf > /dev/null
log "Reloading sysctl"
sysctl -p

# Setup User Accounts
log "Setting root password"
echo "root:changeme" | chpasswd

log "Creating gpadmin user"
useradd -m gpadmin --groups wheel

log "Setting gpadmin password"
echo "gpadmin:changeme" | chpasswd

log "Configuring SSH"
cp -pr /home/vagrant/.ssh /home/gpadmin/
chown -R gpadmin:gpadmin /home/gpadmin

log "Allow passwordless sudo for gpadmin"
echo "%gpadmin ALL=(ALL) NOPASSWD: ALL" > /etc/sudoers.d/gpadmin

# Manage Software
log "Cleaning RPM cache"
sed -i 's/gpgcheck=1/gpgcheck=0/g' /etc/yum.repos.d/*
sudo yum -q -e 0 clean all

log "Installing RPMs"
sudo yum -y -q -e 0 install ed unzip tar git strace gdb vim-enhanced wget m4 > /tmp/yum.out

log "Changing the permission of /usr/local"
chmod 777 /usr/local

log "Ignoring Modifications to UAA Token File"
cd /vagrant && git update-index --assume-unchanged UAA.token

log "OS Setup Complete"


