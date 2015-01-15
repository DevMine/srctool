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
	golint ${GOPATH}/src/${PKG}

deps:
	go get -u -v github.com/codegangsta/cli

cover:
	go test -cover ${PKG}/...

clean:
	rm -f ./${EXEC}
