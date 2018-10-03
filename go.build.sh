#!/usr/bin/env bash
#!/bin/bash
set -e

source /vagrant/functions.h

# Core check
TEST_FAILURE=FALSE

# Check if the internet connection is working
wget -q --tries=2 --timeout=5 --spider http://google.com
if [[ $? -eq 0 ]]; then
        log "Internet connection is available: PASSED"
else
        log "Internet connection is available: FAILED"
        TEST_FAILURE=TRUE
fi

# Check if /usr/local/ directory exists & writable ( needed for installing GPDB software )
if [ -d /usr/local ]; then
     log "Directory /usr/local/ exists: PASSED"
else
    log "Directory /usr/local/ exists: FAILED"
    TEST_FAILURE=TRUE
fi

if [ -w /usr/local ]; then
     log "Directory /usr/local/ writable: PASSED"
else
    log "Directory /usr/local/ writable: FAILED"
    TEST_FAILURE=TRUE
fi

# Check if the BASE DIRECTORY exists & writable
BASE_DIR=`grep BASE_DIR /vagrant/config.yml | cut -d':' -f2 | awk '{print $1}'`
if [ -d "$BASE_DIR" ]; then
     log "Base directory $BASE_DIR exists: PASSED"
else
    log "Base directory $BASE_DIR exists: FAILED"
    TEST_FAILURE=TRUE
fi

if [ -w "$BASE_DIR" ]; then
     log "Directory $BASE_DIR writable: PASSED"
else
    log "Directory $BASE_DIR writable: FAILED"
    TEST_FAILURE=TRUE
fi


# Check if the hostname is reachable
host=`grep MASTER_HOST /vagrant/config.yml | cut -d':' -f2 | awk '{print $1}'`
ping -c 1 $host &>/dev/null
if [ $? -eq 0 ]; then
    log "Host $host can be reached: PASSED"
else
   log "Host $host can be reached: FAILED"
   TEST_FAILURE=TRUE
fi

# If any one of the precheck failed, then exit the setup process.
if [ $TEST_FAILURE == "TRUE" ]; then
    log "Pre check failed, exiting...."
	exit 1
fi

# Download and install GO Binaries.
# Setting up go version to download
VERSION="1.7.4"
DFILE="go$VERSION.linux-amd64.tar.gz"

# If the version of go already exit then uninstall it
if [ -d "$HOME/.go" ]; then
        rm -rf $HOME/.go
fi

# Downloading the go tar file
log "Downloading the GO binary $DFILE"
log "Please wait might take few minutes based on your internet connection"
wget https://storage.googleapis.com/golang/$DFILE -O /tmp/go.tar.gz -q
if [ $? -ne 0 ]; then
    log "Download failed! Exiting."
    exit 1
fi

# Extracting the file
log "Extracting ..."
tar -C "$HOME" -xzf /tmp/go.tar.gz
mv "$HOME/go" "$HOME/.go"
chown -R gpadmin:gpadmin "$HOME/.go"

# Update environment information
# Updating the bashrc with the information of GOROOT.
if grep -q "GOROOT" "$HOME/.bashrc";
then
    log "GOROOT binaries location is already updated on the .bashrc file"
else
    touch "$HOME/.bashrc"
    {
        echo '# Golang binaries'
        echo 'export GOROOT=$HOME/.go'
        echo 'export PATH=$PATH:$GOROOT/bin'
    } >> "$HOME/.bashrc"
fi

# Update bashrc with the information of GOPATH.
if grep -q "GOPATH" "$HOME/.bashrc";
then
    log "GOPATH location is already updated on the .bashrc file"
else
    pwd=`pwd`
    touch "$HOME/.bashrc"
    {
        echo '# GOPATH location'
        echo 'export GOPATH='${pwd}
        echo 'export PATH=$PATH:$GOPATH/bin'
    } >> "$HOME/.bashrc"
fi

# Remove the downloaded tar file
rm -f /tmp/go.tar.gz

log "Removing src / pkg directory to pull in the newer version of the code"
rm -rf src/
rm -rf pkg/

# Download program dependencies
log "Downloading program dependencies"

# go-logging package
# YAML package
source "$HOME/.bashrc"
go get github.com/op/go-logging
if [ $? -ne 0 ]; then
    log "Download failed of dependencies (go-logging) package failed. Exiting....."
    exit 1
fi

# YAML package
source "$HOME/.bashrc"
go get gopkg.in/yaml.v2
if [ $? -ne 0 ]; then
    log "Download failed of dependencies (yaml.v2) package failed. Exiting....."
    exit 1
fi

# gpdb source code
go get github.com/ielizaga/piv-go-gpdb
if [ $? -ne 0 ]; then
    echo "Download failed of dependencies (piv-go-gpdb) package failed. Exiting....."
    exit 1
fi

# source "$HOME/.bashrc"

#
# Changing the owner to gpadmin:gpadmin
#
# chown -R gpadmin:gpadmin /home/gpadmin

#
# Build go executable file.
#

log "Compiling the program... "
# Compile the program
go build $GOPATH/src/github.com/ielizaga/piv-go-gpdb/gpdb/gpdb.go
if [ $? -ne 0 ]; then
    log "Cannot build gpdb executable, exiting ....."
    exit 1
fi

# move the binary file to bin directory
if [ ! -d bin ]; then
    mkdir -p $GOPATH/bin/
fi

# move it to bin directory (forcefully, no need to prompt)
mv -f gpdb $GOPATH/bin/

#
# Changing the owner to gpadmin:gpadmin
#
#chown -R gpadmin:gpadmin /home/gpadmin

#
# Success message.
#

log "GPDBInstall Script has been successfully installed"
log "Config file is cached at location: "$HOME/.config.yml
log "Please close this terminal and open up a new terminal to set the environment"
