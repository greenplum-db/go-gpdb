# Installing and Running gpdb CLI on your own VM

This guide has been tested on fresh CentOS7 and RHEL7 instances spun up in Google Cloud Platform.

## Prepare your environment

Configure a non-root user (e.g. gpadmin) and grant it the necessary privileges to `sudo yum install` the greenplum-db packages to the locations specified in config.yml in the next section.

```sh
yum install -y sshpass sudo wget
```

## Install and configure go-gpdb

```sh
# Download the latest gpdb installer release
cd ~
curl -s https://api.github.com/repos/greenplum-db/go-gpdb/releases/latest \
      | grep "browser_download_url.*gpdb" \
      | grep -v "browser_download_url.*datalab" \
      | cut -d : -f 2,3 | tr -d \" \
      | wget -qi - -O gpdb

chmod +x gpdb

wget https://raw.githubusercontent.com/pivotal-gss/go-gpdb/master/gpdb/config.yml
# For datalab created VM, the below changes are not needed
# For your own VM, make sure the directory is created and your current user has permission 
sed -i "s/gpadmin/${USER}/g" config.yml
mkdir -p data src
sed -i "s|/data|data|g" config.yml
sed -i "s|/usr/local/src|src|g" config.yml

# Specify a single node GPDB installation on localhost
printf "${HOSTNAME}" > hostfile
```

## Run go-gpdb

+ For the open source gpdb release, run:
`./gpdb download -v 6.1.0 --github --install`
+ For the closed source [PivNet](https://network.pivotal.io/) gpdb release, run:
`./gpdb download -v 6.1.0 --install`

## Use your new gpdb cluster! 😉🥳

For further information on the `gpdb` tool, check out [README.md](README.md).

## Multi Host Greenplum Cluster

If you have multiple VM host and requires the tool to create a multi node Greenplum cluster, then ensure you add all the host information on the `/etc/hosts` of the master node (or on the node where you are going to run the tool). Also please check the steps 3 of the troubleshooting section below on how to setup your `/etc/hosts`
 
## Troubleshooting

+ If you run into ssh issues, make sure `ssh $HOSTNAME` works. If not, you can run `ssh-keygen -f ~/.ssh/id_rsa -N '' && cat .ssh/id_rsa.pub >> .ssh/authorized_keys`.
+ Some warnings such as "Master open file limit is 1024 should be >= 65535" can be resolved with this script: [os.prep.sh](../scripts/os.prep.sh).
+ The first line of the `/etc/hosts` is always ignored and it can be anything, your hostname should start from the second line and its usually where the localhost information should start.

**NOTE:** Comments on `/etc/hosts` are not ignored by the tool, its a drawback so ensure there is no comments on /etc/hosts after the first line

Your `/etc/hosts` would have to look something like this.
```sh
$ cat /etc/hosts
127.0.0.1 localhost
192.168.99.100 gpdb-m
```

If you have multiple hosts and want to create a cluster, then add them like this below

```sh
$ cat /etc/hosts
127.0.0.1 localhost
192.168.99.100 gpdb-m
192.168.99.101 gpdb1-m
192.168.99.102 gpdb2-m
```

The tool will detect that there are multiple hosts and create a cluster, no extra steps needed.

## Demo 

[![demo](https://img.youtube.com/vi/q5v6ac2lbd4/0.jpg)](https://www.youtube.com/watch?v=q5v6ac2lbd4)
