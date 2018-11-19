# Constants (well, sorta..)
SCRIPT="$0"
PASS="[√]"
FAIL="[X]"

# Banner Messages
function banner() {
	echo
	printf '%.s=' {1..60}
    printf "\n%s\n"  "$*" 
	printf '%.s=' {1..60}
	printf '\n'
}
	
# Compare Version Strings
function compare_versions() {
    if [[ $1 == $2 ]]
    then
        return 0
    fi

    local IFS=.
    local i ver1=($1) ver2=($2)

    # fill empty fields in ver1 with zeros
    for ((i=${#ver1[@]}; i<${#ver2[@]}; i++))
    do
        ver1[i]=0
    done

    for ((i=0; i<${#ver1[@]}; i++))
    do
        if [[ -z ${ver2[i]} ]]
        then
            # fill empty fields in ver2 with zeros
            ver2[i]=0
        fi
        if ((10#${ver1[i]} > 10#${ver2[i]}))
        then
            return 1
        fi
        if ((10#${ver1[i]} < 10#${ver2[i]}))
        then
            return 2
        fi
    done
    return 0
}

# Return current GO version
function go_version() {
    version=$(/usr/local/go/bin/go version)
    regex="([0-9].[0-9].[0-9])"
    if [[ $version =~ $regex ]]; then 
         echo ${BASH_REMATCH[1]}
    fi
}

# Log out including Script Name
function log() {
  printf "%s...\n" "$*" 
}

# Usage `spinner pid display_text`
# Call spinner last background execution with `spinner $!`
# Include script name with `basename $0 .sh`
# Will return the exit status of the watched process. 
function spinner() { 
    local pid=$1
    local delay=0.5
    local spinstr='\|/-'
    while [ "$(ps a | awk '{print $1}' | grep $pid)" ]; do
        local temp=${spinstr#?}
        printf "[%c] $2..." "$spinstr"
        local spinstr=$temp${spinstr%"$temp"}
        sleep $delay
        printf "\r"
    done
	wait $pid
	exit_code=$?
	if [[ $exit_code -eq 0 ]]; then
	    printf "[√] $2...\n"
	else
	    printf "[X] $2...\n"
		return $exit_code
	fi
}


# Call progressBar on known length with `progressBar 19.5 20`
# Change the bar lenght with the 3rd varable `progressBar 19.5 20 60`

function progressBar() {

  local p_progress=${1:-0}
  local p_total=${2:-100}
  local p_barwidth=${3:-40}

  p_percent=$(printf %.2f $(echo "$p_progress/$p_total*100" | bc -l))
  p_fill=$(printf "%.0f" $(echo "$p_barwidth*($p_progress/$p_total)" | bc -l ))
  p_empty=$(echo "$p_barwidth-$p_fill" | bc -l )

  printf "\r["
  if [ $p_fill -gt 0 ]; then printf "%.s▉" $(seq 1 $p_fill); fi
  if [ $p_empty -gt 0 ]; then printf "%.s░" $(seq 1 $p_empty); fi
  printf "] $p_percent %%"
  if [ $p_empty -eq 0 ]; then sleep .5; printf "\n"; fi
}  

# Use with source to load variables...
# source <(parse_yaml /vagrant/gpdb/config.yml)
# echo $ENV_DIR

function parse_yaml() {
   local s='[[:space:]]+' 
   local w='[a-zA-Z0-9_]*'
   sed -E \
	   -e "0,/^---$/d" \
	   -e "s|$s#.*||g" \
	   -e "s|$s(.*):[^:\/\/](.*$)|\1=\"\2\"|g" \
       -e "s|($w):|#\1|g" $1 
}

