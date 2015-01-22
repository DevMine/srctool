PKG  = github.com/DevMine/srctool
EXEC = srctool

all: check test build

install:
	go install ${PKG}

build:
	go build -o ${EXEC} ${PKG}

test:
	go test ${PKG}/...

check:
	go vet ${PKG}/...
	golint ${PKG}/...

deps:
	go get -u -v github.com/codegangsta/cli
	go get -u -v github.com/DevMine/srcanlzr
	go get -u -v github.com/gilliek/go-xterm256/xterm256
	go get -u -v github.com/mitchellh/ioprogress
	go get -u -v golang.org/x/crypto/ssh/terminal
	go get -u -v github.com/DevMine/repotool/model

dev-deps:
	go get -u github.com/golang/lint/golint

cover:
	go test -cover ${PKG}/...

clean:
	rm -f ./${EXEC}
