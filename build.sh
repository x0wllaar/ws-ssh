#!/bin/bash
set -exuo pipefail

cp ./browserproxy/build/rootpage.html ./impl/wsclient/browserproxy/rootpage.gen.html
pushd ./impl/wsclient/browserproxy/
./generate_fake_rootpage.sh
popd

pushd ./browserproxy
./build_rootpage.sh 
popd

cp ./browserproxy/build/rootpage.html ./impl/wsclient/browserproxy/rootpage.gen.html
pushd ./impl/wsclient/browserproxy/
./generate_rootpage.sh
popd

go build -o ./ws-ssh ./cli

