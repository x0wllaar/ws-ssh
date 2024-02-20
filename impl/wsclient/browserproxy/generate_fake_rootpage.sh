#!/bin/bash
set -exuo pipefail

rm -f rootpage.gen.go

cp -v ./rootpage.go.fake ./rootpage.gen.go