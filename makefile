CC=go
RM=rm
MV=mv


SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')
#GOPATH=$(SOURCEDIR)/
GOOS=linux
GOARCH=amd64
#GOARCH=arm
GOARM=7


EXEC=bebopanalyzer

VERSION=1.9
BUILD_TIME=`date +%FT%T%z`
PACKAGES := fmt path/filepath github.com/metakeule/fmtdate github.com/ptrv/go-gpx github.com/gorilla/mux github.com/gorilla/mux


LIBS= 

LDFLAGS=	

.DEFAULT_GOAL:= $(EXEC)

$(EXEC): organize $(SOURCES)
		@echo "    Compilation des sources ${BUILD_TIME}"
		@if  [ "arm" = "${GOARCH}" ]; then\
		    GOPATH=$(PWD)/../.. GOOS=${GOOS} GOARCH=${GOARCH} GOARM=${GOARM} go build ${LDFLAGS} -o ${EXEC}-${VERSION}.${GOOS}-${GOARCH} $(SOURCEDIR)/main.go;\
		else\
            GOPATH=$(PWD)/../.. GOOS=${GOOS} GOARCH=${GOARCH} GOARM=${GOARM} go build ${LDFLAGS} -o ${EXEC}-${VERSION}.${GOOS}-${GOARCH} $(SOURCEDIR)/main.go;\
        fi
		@echo "    ${EXEC}-${VERSION}.${GOOS}-${GOARCH} generated."

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
		@if [ -f "${EXEC}-${VERSION}.${GOOS}-${GOARCH}" ] ; then rm ${EXEC}-${VERSION}.${GOOS}-${GOARCH} ; fi
		@echo "    Nettoyage effectuee"

package:  ${EXEC}
		@zip -r ${EXEC}-${GOOS}-${GOARCH}-${VERSION}.zip ./${EXEC}-${VERSION} resources
		@echo "    Archive ${EXEC}-${GOOS}-${GOARCH}-${VERSION}.zip created"

audit:   ${EXEC}
		@go tool vet -all -shadow ./
		@echo "    Audit effectue"

swagger:
	@echo "Generate swagger json file specs"
	@GOPATH=$(PWD)/../.. go run ${GOPATH}/src/github.com/go-swagger/go-swagger/cmd/swagger/swagger.go generate spec -m -b ./routes > resources/swagger.json
	@echo "Specs generate at resources/swagger.json"
