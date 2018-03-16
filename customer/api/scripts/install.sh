#!/usr/bin/env bash

set -e

build=${PWD}
workspace=$build/.build
version='0.0.1'

name='github.com/pdrosos/hyperledger-fabric-demo/customer/api'
package=$workspace/src/$name
goversion='go1.10'exit
constraint='.'
ldflags="-X ${name}.Version=${version}"

PATH=$PATH:/usr/local/go/bin

function main() {
    info
    compile
}

function compile() {
    trap 'previous_command=$this_command; this_command=$BASH_COMMAND' DEBUG

    if [ -z $DEBUG ]; then
        trap "rm -Rf $workspace" EXIT INT SIGTERM
    fi

    if [ -d $workspace ]; then
        echo "---> Setup"
        echo "  Removing dirty build directory"
        echo ""

        rm -Rf $workspace
    fi

    mkdir -p $build/bin
    mkdir -p $package

    cp -R $build/* $package

    cd $package

    GOBIN=$build/bin GOPATH=$workspace go install -ldflags "$ldflags" -tags $constraint ./...
}

function info() {
    echo "---> Environment"
    echo "  Name:      $name"
    echo "  Version:   $version"
    echo "  Go:        $goversion"
    echo "  Workspace: $workspace"
    echo "  Bin:       $build/bin"
    echo "  Constraint $constraint"
    echo ""
}

main
