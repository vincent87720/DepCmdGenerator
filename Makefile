WINDOWS := windows
PROJNAME := DepCmdGenerator
MAINPKG := github.com/vincent87720/$(PROJNAME)/cmd/$(PROJNAME)

all: test build

build: rmdir release cpbin cpjson

rmdir:
	rm -rf ./bin/* ./test/*

cpbin:
	cp -r ./bin/* test

cpjson:
	cp ./settings.json test/$(WINDOWS)/

##########BUILD##########
.PHONY: buildwindows
buildwindows:
	GOOS=$(WINDOWS) GOARCH=amd64 go build -o bin/$(WINDOWS)/$(PROJNAME).exe $(MAINPKG)

.PHONY: release
release: buildwindows

##########RUN##########
run:
	go run $(MAINPKG)