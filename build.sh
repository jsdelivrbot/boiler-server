fuser -k 8080/tcp

export GOPATH=/home/apps/go
export PATH=$PATH:$GOPATH/bin:/usr/local/go/bin

echo $GOPATH
echo $PATH

cd /home/apps/go/src/github.com/AzureRelease/boiler-server && go build

cd /home/apps/go/src/github.com/AzureRelease/boiler-server && nohup ./boiler-server &

rm -rf common
rm -rf controllers
rm -rf models
rm -f dba/*.go
rm -rf routers

rm main.go
rm build.sh