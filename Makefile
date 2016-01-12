build: clean
	GOOS=linux GOARCH=amd64 go build -o igcstat_lnx
	GOOS=darwin GOARCH=amd64 go build -o igcstat_osx
	GOOS=windows GOARCH=amd64 go build -o igcstat.exe

clean:
	rm -f igcstat_lnx igcstat_osx igcstat.exe
	rm -f *.csv *.xlsx
