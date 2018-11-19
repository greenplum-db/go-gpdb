#!/usr/bin/env bash
#!/bin/bash

set -e
source /vagrant/scripts/functions.h
# if ! [ -f /vagrant/gpdb/config.yml ]; then abort "MISSING YAML"; fi
source <(parse_yaml /vagrant/gpdb/config.yml)
set +e

abort() {
	log "$FAIL Return Code: [$1]"
	exit $1
}

cleanup() {
	banner "Cleanup"
	
	for package in "${src[@]}"
	do
		{ rm -rf $package & } &>/dev/null
		spinner $! "Removing Temporary File: $package"
	done
	
	if ! $build_complete; then
		banner "Completed"
		log "$PASS /vagrant/scripts/go.build.sh"
	else
		banner "Failed"
		log "$FAIL /vagrant/scripts/go.build.sh"
	fi
	echo
}

trap cleanup EXIT

go_install() {	
	# Download GO
	{ wget -q https://storage.googleapis.com/golang/go$GO_BUILD.$OS-$ARCH.tar.gz -O $BASE_DIR/go.tar.gz & } &>/dev/null
	spinner $! "Downloading GO Binary: $GO_BUILD"
	if [ $? -ne 0 ]; then wait $!; abort $?; fi

	# Extract
	{ tar -C "/usr/local" -xzf $BASE_DIR/go.tar.gz & } &>/dev/null
	spinner $! "Extracting: $GO_BUILD"
	if [ $? -ne 0 ]; then wait $!; abort $?; fi
		
	{ rm -rf "$BASE_DIR/go.tar.gz" & } &>/dev/null
	spinner $! "Removing Temporary File: $BASE_DIR/go.tar.gz"
	
	# Notify	
	log "$PASS GO Binary Version Installed: $GO_BUILD"

	# Update Environment Variables if it doesn't exist
	if ! ( grep -q "# GOLANG" /etc/profile.d/gpdb.profile.sh &>/dev/null ); then
	    {
	        echo '# GOLANG'
	        echo 'export GOROOT=/usr/local/go'
	        echo 'export GOPATH='$BASE_DIR
	        echo 'export PATH=$PATH:$GOROOT/bin:$GOPATH/bin'
	    } >> /etc/profile.d/gpdb.profile.sh
		spinner $! "Update Environment Variables"
		if [ $? -ne 0 ]; then wait $!; abort $?; fi

	fi 	
}


banner "Configuration"

# Internet Connetivity
{ wget -q --tries=2 --timeout=5 --spider http://google.com & } &>/dev/null 
spinner $! "Internet Connection"
if [ $? -ne 0 ]; then wait $!; abort $?; fi

# YAML: BASE_DIR
{ mkdir -p "$BASE_DIR" && test -w "$BASE_DIR" & } &>/dev/null 
spinner $! "YAML: BASE_DIR: $BASE_DIR" 
if [ $? -ne 0 ]; then wait $!; abort $?; fi

# Hostname
# { ping -c 1 $MASTER_HOST & } &>/dev/null 
{ ping -c 1 `hostname` & } &>/dev/null 
spinner $! "YAML: HOSTNAME: `hostname`"
if [ $? -ne 0 ]; then wait $!; abort $?; fi

banner "GOLANG Installation"

# GO Binaries
if ! [ -d "/usr/local/go" ]; then
	go_install
else
	# Compare 1 : 2 [EQ 0; GT 1; LT 2) 
	compare_versions $(go_version) $GO_BUILD

	if [ $? -lt 2 ]; then
		log "$PASS GO Binary Version Required: $GO_BUILD (Installed: $(go_version))" 
	else
		log "$FAIL GO Binary Version Required: $GO_BUILD (Installed: $(go_version))"

		# Backup Exisitng Build
		{ mv /usr/local/go /usr/local/go.$(go_version) & } &>/dev/null
		spinner $! "Backing Up Existing Build"
		if [ $? -ne 0 ]; then wait $!; abort $?; fi
		
		# Call Installer
		go_install
	fi
fi

source /etc/profile.d/gpdb.profile.sh		

# Download program dependencies -- eventually, we will check for updates with git
banner "Program Dependencies"
src[0]='github.com/op/go-logging'
src[1]='github.com/codingconcepts/env' 	
for package in "${src[@]}"
do
	{ go get $package & } &>/dev/null
	spinner $! "$package"
	if [ $? -ne 0 ]; then wait $!; abort $?; fi
done

# Compile Time...
banner "Compiling"
{ go build /vagrant/gpdb/gpdb.go & } &>/dev/null
spinner $! "Building"
if [ $? -ne 0 ]; then wait $!; abort $?; fi

{ (mkdir -p $GOPATH/bin/ && mv -f gpdb $GOPATH/bin/) & } &>/dev/null
spinner $! "Installing"
if [ $? -ne 0 ]; then wait $!; abort $?; fi

build_complete=true

exit 0



