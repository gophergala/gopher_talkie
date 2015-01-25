GO := `which go`
GOFMT := `which gofmt`
GOVET := ./scripts/vet

all: cli server

cli: 
	@/bin/bash ./scripts/build talkie

server:
	@/bin/bash ./scripts/build server
	
test:
	@/bin/bash ./scripts/test common
	@/bin/bash ./scripts/test audio
	@/bin/bash ./scripts/test crypto
	
check:
	@./.hooks/pre-commit

vet:
	@/bin/bash ./scripts/vet

format:
	@git ls-files | grep '.go$$' | xargs $(GOFMT) -w -s

deps:
	@/bin/bash ./scripts/deps

clean:
	@rm -rf .godeps/pkg/*

.PHONY = all cli test check vet format deps clean
