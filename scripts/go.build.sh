#!/usr/bin/env bash
#!/bin/bash
set -e
source /vagrant/functions.h

abort() { 
  log "FAILED"
  exit 1
}

log "Testing: Internet Connection" 	
if ! wget -q --tries=2 --timeout=5 --spider http://google.com ; then abort; fi

log "Testing: GO Binary" 	
if ! [ -d "/usr/local/go" ]; then

	log "GO Binary Does Not Exist"
	log "Downloading GO Binary: $GOVERSION"
	if ! wget https://storage.googleapis.com/golang/go$GO_VERSION.$OS-$ARCH.tar.gz -O /tmp/go.tar.gz -q; then abort; fi

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
else
	# todo: version checking and updgrade / downgrade	
	log "GO Binary Version Installed: " go_version
fi

