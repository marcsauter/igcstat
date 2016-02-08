SRCDIR = .
DISTDIR = $(SRCDIR)/dist

SOURCES := $(shell find $(SRCDIR) -name '*.go') 
BINARY = igcstat
LNXBIN = igcstat-linux-amd64
OSXBIN = igcstat-darwin-amd64
WINBIN = igcstat-windows-amd64.exe
LNXDIST = dist/igcstat-linux-amd64.tar
OSXDIST = dist/igcstat-osx-amd64.tar
WINDIST = dist/igcstat-windows-amd64.zip

all: $(LNXBIN) $(OSXBIN) $(WINBIN)

$(LNXBIN): $(SOURCES)
	GOOS=linux GOARCH=amd64 go build -o ${LNXBIN}

$(OSXBIN): $(SOURCES)
	GOOS=darwin GOARCH=amd64 go build -o ${OSXBIN}

$(WINBIN): $(SOURCES)
	GOOS=windows GOARCH=amd64 go build -o ${WINBIN}

dist: clean all
	if [ ! -d ${DISTDIR} ] ; then mkdir ${DISTDIR} ; fi
	tar --transform="flags=r;s|${LNXBIN}|${BINARY}|" -cf ${LNXDIST} ${LNXBIN}
	tar --transform="flags=r;s|${OSXBIN}|${BINARY}|" -cf ${OSXDIST} ${OSXBIN}
	zip ${WINDIST} -j scripts ${WINBIN} scripts/igcstat.cmd
	printf "@ ${WINBIN}\n@=${BINARY}.exe\n" | zipnote -w ${WINDIST}

deps:
	go get -u ./...

clean:
	rm -f ${LNXBIN} ${OSXBIN} ${WINBIN} ${LNXDIST} ${OSXDIST} ${WINDIST}
	rm -f *.csv *.xlsx

.PHONY: all dist deps clean
