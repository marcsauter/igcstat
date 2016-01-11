build:
	GOOS=linux GOARCH=amd64 go build -o igcstat_lnx
	GOOS=darwin GOARCH=amd64 go build -o igcstat_osx
	GOOS=windows GOARCH=amd64 go build -o igcstat.exe
