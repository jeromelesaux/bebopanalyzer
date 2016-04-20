CC=go
RM=rm
MV=mv


SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')
#GOPATH=$(SOURCEDIR)/
GOOS=linux
GOARCH=amd64
#GOARCH=arm


EXEC=BebopAnalyzer

VERSION=1.1.0
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
		@go get -d -v $(PACKAGES)

organize: deps
		@echo "    Go FMT"
		@$(foreach element,$(SOURCES),go fmt $(element);)

init: clean
		@echo "    Init of the project"
		@export GOPATH=$(PWD)
		$(shell export GOPATH=$(pwd))

execute:
		./${EXEC}-${VERSION}

clean:
		@rm -f *.o core
		@if [ -f "${EXEC}-${VERSION}" ] ; then rm ${EXEC}-${VERSION} ; fi
		@echo "    Nettoyage effectuee"

