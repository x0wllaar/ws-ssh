#!/bin/bash
set -exuo pipefail

rm -f rootpage.gen.go

echo package browserproxy > rootpage.gen.go
echo "" >> rootpage.gen.go
echo "" >> rootpage.gen.go

echo const ROOTPAGE_HTML = '`' >> rootpage.gen.go

#https://stackoverflow.com/a/56144740
cat rootpage.gen.html | sed 's/`/`+"`"+`/g' >> rootpage.gen.go

echo '`' >> rootpage.gen.go
gofmt -s -w rootpage.gen.go