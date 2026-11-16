BINARY := ryvst

GO ?= go
GOFLAGS ?=

.PHONY: all build run install clean

all: build

build:
	@${GO} build ${GOFLAGS} -o ./${BINARY} ./cmd/ryvst
	@echo "Build ./${BINARY}"

run: build
	@./${BINARY}

install: build
	@mkdir -p $HOME/bin
	install ${BINARY} ${HOME}/bin/${BINARY}

clean:
	@rm -f ${BINARY}
	@echo "Removed ./${BINARY}"

