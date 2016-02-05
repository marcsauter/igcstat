all: igcstat_lnx igcstat_osx igcstat.exe

igcstat_lnx:
	GOOS=linux GOARCH=amd64 go build -o igcstat_lnx

igcstat_osx:
	GOOS=darwin GOARCH=amd64 go build -o igcstat_osx

igcstat.exe:
	GOOS=windows GOARCH=amd64 go build -o igcstat.exe

clean:
	rm -f igcstat_lnx igcstat_osx igcstat.exe
	rm -f *.csv *.xlsx

.PHONY: all clean
