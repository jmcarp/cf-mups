# cf-mups

*cf plugin for setting individual credentials on user-provided services*

## Installation

    go get github.com/jmcarp/cf-mups
    cf install-plugin $GOPATH/bin/cf-mups

## Usage

    cf mups SERVICE_NAME CREDENTIAL_NAME CREDENTIAL_VALUE
