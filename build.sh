fuser -k 8080/tcp

echo $GOPATH
echo $PATH

cd /home/apps/go/src/github.com/AzureRelease/boiler-server && go build

cd /home/apps/go/src/github.com/AzureRelease/boiler-server && nohup ./boiler-server &

rm -rf common
rm -rf controllers
rm -rf models
rm -rf dba
rm -rf routers

rm main.go
rm build.sh