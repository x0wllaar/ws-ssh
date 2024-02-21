#!/bin/bash
set -exuo pipefail

rm -rf ./build
mkdir -p ./build
cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" ./build
cp ./rootpage.html ./build/rootpage.html

GOOS=js GOARCH=wasm go build -o ./build/bpwasm.wasm ./bpwasm/
cat ./build/bpwasm.wasm | base64 -w 0 | tr -d "\n" > ./build/bpwasm.wasm.b64
perl -pe 's/"WASM_EXEC_JS_INSERTED_HERE"/`cat .\/build\/wasm_exec.js`/ge' -i ./build/rootpage.html

perl -pe 's/WASM_BASE_64_HERE/`cat .\/build\/bpwasm.wasm.b64`/ge' -i ./build/rootpage.html