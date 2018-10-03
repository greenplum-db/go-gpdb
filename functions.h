# Separating script output from echo input...
log()
{
	printf "[`basename $0 .sh`]: %s...\n" "$*" 
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