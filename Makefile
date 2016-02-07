SRCDIR=.
DISTDIR=$(SRCDIR)/dist

SOURCES := $(shell find $(SOURCEDIR) -name '*.go') 
LNXBIN=igcstat-linux-amd64
OSXBIN=igcstat-darwin-amd64
WINBIN=igcstat-windows-amd64.exe

all: $(LNXBIN) $(OSXBIN) $(WINBIN)

$(LNXBIN): $(SOURCES)
	GOOS=linux GOARCH=amd64 go build -o ${LNXBIN}

$(OSXBIN): $(SOURCES)
	GOOS=darwin GOARCH=amd64 go build -o ${OSXBIN}

$(WINBIN): $(SOURCES)
	GOOS=windows GOARCH=amd64 go build -o ${WINBIN}

dist: clean all
	if [ ! -d ${DISTDIR} ] ; then mkdir ${DISTDIR} ; fi
	cp ${LNXBIN} ${DISTDIR}
	cp ${OSXBIN} ${DISTDIR}
	cp ${WINBIN} ${DISTDIR}

clean:
	rm -f ${LNXBIN} ${OSXBIN} ${WINBIN}
	rm -f *.csv *.xlsx

.PHONY: all dist clean
