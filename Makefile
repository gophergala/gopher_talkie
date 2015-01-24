GO := `which go`
GOFMT := `which gofmt`
GOVET := ./scripts/vet

all: app server

app: 
	@/bin/bash ./scripts/build app

server:
	@/bin/bash ./scripts/build server

test: app server
	@/bin/bash ./scripts/test app
	@/bin/bash ./scripts/test server

check:
	@./.hooks/pre-commit

vet:
	@git ls-files | grep '.go$$' | while read i; do $(GO) vet $$i 2>&1; done | grep -v exit\ status | grep -v pb.go | grep -v Error\ call

format:
	@git ls-files | grep '.go$$' | xargs $(GOFMT) -w -s

.PHONY = all app server test check vet format
