# Separating script output from echo input...
log()
{
	printf "[`basename $0 .sh`]: %s...\n" "$*" 
}

function go_version {
    version=$(go version)
    regex="([0-9].[0-9].[0-9])"
    if [[ $version =~ $regex ]]; then 
         echo ${BASH_REMATCH[1]}
    fi
}

# SPINNER for Long Running Processes...
# Call on last executed with `spinner $!`
spinner()
{
    local pid=$1
    local delay=0.75
    local spinstr='|/-\'
    while [ "$(ps a | awk '{print $1}' | grep $pid)" ]; do
        local temp=${spinstr#?}
        printf " [%c]  " "$spinstr"
        local spinstr=$temp${spinstr%"$temp"}
        sleep $delay
        printf "\b\b\b\b\b\b"
    done
    printf "    \b\b\b\b"
    printf "[√]\n"
}