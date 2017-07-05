CC=go
RM=rm
MV=mv


SOURCEDIR=.
SOURCES:= $(shell find $(SOURCEDIR) -name '*.go')
GOARM=7
#GOARCH=amd64


EXEC=bebopanalyzer

VERSION=1.10
BUILD_TIME=`date +%FT%T%z`
PACKAGES:= github.com/metakeule/fmtdate github.com/ptrv/go-gpx github.com/gorilla/mux github.com/gorilla/mux


LIBS= 

LDFLAGS=	

.DEFAULT_GOAL:= $(EXEC)

$(EXEC): organize $(SOURCES)
		@echo "    Compilation des sources ${BUILD_TIME}"
		@if  [ "arm" = "${GOARCH}" ]; then\
		   GOOS=${GOOS} GOARCH=${GOARCH} GOARM=${GOARM} go build ${LDFLAGS} -o ${EXEC} $(SOURCEDIR)/main.go;\
		else\
           GOOS=${GOOS} GOARCH=${GOARCH} GOARM=${GOARM} go build ${LDFLAGS} -o ${EXEC} $(SOURCEDIR)/main.go;\
        fi
		@echo "    ${EXEC} version:${VERSION}.${GOOS}-${GOARCH} generated."

deps: init
		@echo "    Download packages"
		@$(foreach element,$(PACKAGES),go get -d -v $(element);)

organize: deps
		@echo "    Go FMT"
		@$(foreach element,$(SOURCES),go fmt $(element);)

init: clean
		@echo "    Init of the project"

execute:
		./${EXEC} conf.json

clean:
		@if [ -f "${EXEC}" ] ; then rm ${EXEC} ; fi
		@echo "    Nettoyage effectuee"

package:  ${EXEC} swagger
		@zip -r ${EXEC}-${GOOS}-${GOARCH}-${VERSION}.zip ./${EXEC} resources
		@echo "    Archive ${EXEC}-${GOOS}-${GOARCH}-${VERSION}.zip created"

audit:   ${EXEC}
		@go tool vet -all -shadow ./
		@echo "    Audit effectue"

swagger:
	@echo "Generate swagger json file specs"
	@GOOS=linux GOARCH=amd64 go run ${GOPATH}/src/github.com/go-swagger/go-swagger/cmd/swagger/swagger.go generate spec -m -b ./routes > resources/swagger.json
	@echo "Specs generate at resources/swagger.json"
