#!/bin/bash
source /vagrant/scripts/functions.h

abort() {
	log "$FAIL Return Code: [$1]"
	exit $1
}

## OS Parameters
banner "OS Parameters"

{ echo "kernel.sem = 250 512000 100 2048" >> /etc/sysctl.conf & } &> /dev/null
spinner $! "Upgrading the semaphore values"
if [[ $? -ne 0 ]]; then wait $!; abort $?; fi
	
{ sysctl -p & } &> /dev/null
spinner $! "Reloading sysctl"
if [[ $? -ne 0 ]]; then wait $!; abort $?; fi

## User Accounts
banner "User Accounts"

{ echo "root:changeme" | chpasswd & } &> /dev/null
spinner $! "Setting root password"

{ getent group wheel || groupadd wheel & } &> /dev/null
spinner $! "Creating user group if it doesn't exists"

{ useradd -s /bin/bash -m gpadmin --groups wheel & } &> /dev/null
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
    echo 'export GPDB_SEGMENTS='$2
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

{ sudo yum -y -q -e 0 install bc ed gdb git m4 strace tar unzip vim-enhanced wget sshpass & } &> /tmp/yum.out
spinner $! "Installing RPMs"

## Default Directories & Permissions
banner "Default Directories & Permissions"
{ 
  sudo chmod 777 /usr/local
  sudo chmod 777 /usr/local/src
  sudo chmod 777 /usr/local/bin
  sudo chmod 777 /etc/profile.d/gpdb.profile.sh
} &>> /tmp/yum.out
spinner $! "Changing the permission of /usr/local"

{
  sudo mkdir -p /data/master
  sudo mkdir -p /data/primary
  sudo mkdir -p /data/mirror
  sudo chown gpadmin:gpadmin /data/master
  sudo chown gpadmin:gpadmin /data/primary
  sudo chown gpadmin:gpadmin /data/mirror
} &>> /tmp/yum.out
spinner $! "Creating & Changing ownership of database directories"

banner "OS Setup Complete"