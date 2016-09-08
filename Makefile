NAME=trello_to_clubhouse
BINPATH=bin

build:
	GOOS=windows GOARCH=386   go build -o $(BINPATH)/$(NAME)_windows_x86.exe ./*.go
	GOOS=windows GOARCH=amd64 go build -o $(BINPATH)/$(NAME)_windows_x64.exe ./*.go
	GOOS=darwin  GOARCH=amd64 go build -o $(BINPATH)/$(NAME)_osx_x64 ./*.go
	GOOS=linux   GOARCH=amd64 go build -o $(BINPATH)/$(NAME)_linux_x64 ./*.go
