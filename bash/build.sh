echo $GOPATH
echo $PATH

cd /home/apps/go/src/github.com/AzureRelease/boiler-server && go build

rm -rf common
rm -rf controllers
rm -rf routers

rm -f models/*.go
rm -f dba/*.go

rm -r main.go
rm -rf bash