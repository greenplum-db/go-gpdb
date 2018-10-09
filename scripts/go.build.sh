#!/usr/bin/env bash
#!/bin/bash
set -e
source /vagrant/scripts/functions.h

OS = "linux"
ARCH = "amd64"
GO_VERSION = "1.7.4"

abort() { 
  log "FAILED"
  exit 1
}

log "Internet Connection" 	
if ! wget -q --tries=2 --timeout=5 --spider http://google.com ; then abort; fi

if ! [ -d "/usr/local/go" ]; then
	
	log "Downloading GO Binary: $GOVERSION"
	wget https://storage.googleapis.com/golang/go$GO_VERSION.$OS-$ARCH.tar.gz -O /tmp/go.tar.gz -q &
	spinner $!
	if $?; then abort; fi

	log "Extracting"
	tar -C "/usr/local" -xzf /tmp/go.tar.gz

	if grep -q "GOROOT" "$HOME/.bashrc"; then
	    log "GOROOT binaries location is already updated on the .bashrc file"
	else
	    touch "$HOME/.bashrc"
	    {
	        echo '# Golang binaries'
	        echo 'export GOROOT=/usr/local/go'
	        echo 'export PATH=$PATH:$GOROOT/bin'
	    } >> "$HOME/.bashrc"
	fi
	
	log "GO Binary Version Installed: " go_version
else
	log "GO Binary Version Installed: " go_version
fi

