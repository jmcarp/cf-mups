# cf-mups

*cf plugin for setting individual credentials on user-provided services*

## Installation

    go get github.com/jmcarp/cf-mups
    cf install-plugin $GOPATH/bin/cf-mups

## Usage

    cf mups <service_name> <credential-name> <credential-value>

    cf mups creds secret '"shh"'
    cf mups creds secret '{"super": "secret"}'
