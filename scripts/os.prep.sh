#!/bin/bash
source /vagrant/scripts/functions.h

## OS Parameters
banner "OS Parameters"

{ echo "kernel.sem = 250 512000 100 2048" >> /etc/sysctl.conf & } &> /dev/null
spinner $! "Upgrading the semaphore values"
if [ $? -ne 0 ]; then wait $!; abort $?; fi
	
{ sysctl -p & } &> /dev/null
spinner $! "Reloading sysctl"
if [ $? -ne 0 ]; then wait $!; abort $?; fi

## User Accounts
banner "User Accounts"

{ echo "root:changeme" | chpasswd & } &> /dev/null
spinner $! "Setting root password"

{ useradd -m gpadmin --groups wheel & } &> /dev/null
spinner $! "Creating gpadmin user"

{ echo "gpadmin:changeme" | chpasswd & } &> /dev/null
spinner $! "Setting gpadmin password"

{ cp -pr /home/vagrant/.ssh /home/gpadmin/ 
  chown -R gpadmin:gpadmin /home/gpadmin & 
} &> /dev/null
spinner $! "Configuring SSH"

{ echo "%gpadmin ALL=(ALL) NOPASSWD: ALL" > /etc/sudoers.d/gpadmin & } &> /dev/null
spinner $! "Allow passwordless sudo for gpadmin"

{ 
  { echo '# GO-GPDB'
    echo 'export UAA_API_TOKEN='$1
    echo 'source /vagrant/scripts/functions.h'
    echo 'source <(parse_yaml /vagrant/gpdb/config.yml)'
  } >> /etc/profile.d/gpdb.profile.sh
  chmod +x /etc/profile.d/gpdb.profile.sh
} &> /dev/null
spinner $! "Update Environment Variables"

## Package Installation
banner "Package Installation"

{ 
  sed -i 's/gpgcheck=1/gpgcheck=0/g' /etc/yum.repos.d/* 
  sudo yum -q -e 0 clean all 
} &> /dev/null
spinner $! "Cleaning RPM cache"

{ sudo yum -y -q -e 0 install bc ed gdb git m4 strace tar unzip vim-enhanced wget & } &> /tmp/yum.out
spinner $! "Installing RPMs"

## Default Permisions
banner "Permissions"

{ 
  chmod 777 /usr/local
  chmod 777 /usr/local/src
} &> /tmp/yum.out
spinner $! "Changing the permission of /usr/local"


banner "OS Setup Complete"


