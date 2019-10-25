# Installing and Running gpdb CLI on your own VM

This guide has been tested on fresh CentOS7 and RHEL7 instances spun up in Google Cloud Platform.

## Prepare your environment

Make sure you are logged in as a non-root user (e.g. gpadmin) which has the necessary privileges to `sudo yum install` the greenplum-db packages.

## Install and configure go-gpdb

```sh
sudo yum install -y sshpass wget

cd
wget https://github.com/pivotal-gss/go-gpdb/releases/download/v3.3.0/gpdb
chmod +x gpdb

wget https://raw.githubusercontent.com/pivotal-gss/go-gpdb/master/gpdb/config.yml
sed -i "s/gpadmin/$USER/g" config.yml
sed -i "s|/data|data|g" config.yml
sed -i "s|/usr/local/src|src|g" config.yml

# Specify a single node GPDB installation on localhost
printf $HOSTNAME > hostfile
```

## Run go-gpdb

For the open source gpdb release, run:
`./gpdb download -v 6.1.0 --github --install`

For the closed source pivnet gpdb release, run:
`./gpdb download -v 6.1.0 --install`

## Use your new gpdb cluster! 😉🥳

For further information on the `gpdb` tool, check out [README.md](README.md).

## Troubleshooting

+ If you run into ssh issues, make sure `ssh $HOSTNAME` works. If not, you can run `ssh-keygen -f ~/.ssh/id_rsa -N '' && cat .ssh/id_rsa.pub >> .ssh/authorized_keys`.
+ Here is a sample script to set up all of the packages and permissions needed on your host for gpdb installation: [os.prep.sh](scripts/os.prep.sh).
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

The tool will detect that there is multiple host and create a cluster, no extra steps needed
