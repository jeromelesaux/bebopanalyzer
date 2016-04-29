CC=go
RM=rm
MV=mv


SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')
#GOPATH=$(SOURCEDIR)/
GOOS=linux
GOARCH=amd64
#GOARCH=arm


EXEC=bebopanalyzer

VERSION=1.2.0
BUILD_TIME=`date +%FT%T%z`
PACKAGES := fmt path/filepath


LIBS= 

LDFLAGS=	

.DEFAULT_GOAL: $(EXEC)

$(EXEC): organize $(SOURCES)
		@echo "    Compilation des sources ${BUILD_TIME}"
		@GOPATH=$(PWD)/../.. GOOS=${GOOS} GOARCH=${GOARCH} go build ${LDFLAGS} -o ${EXEC}-${VERSION} $(SOURCEDIR)/main.go
		@echo "    ${EXEC}-${VERSION} generated."

deps: init
		@echo "    Download packages"
		@$(foreach element,$(PACKAGES),go get -d -v $(element);)

organize: deps
		@echo "    Go FMT"
		@$(foreach element,$(SOURCES),go fmt $(element);)

init: clean
		@echo "    Init of the project"

execute:
		./${EXEC}-${VERSION}

clean:
		@if [ -f "${EXEC}-${VERSION}" ] ; then rm ${EXEC}-${VERSION} ; fi
		@echo "    Nettoyage effectuee"