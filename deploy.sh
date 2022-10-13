IP=$1
DOWNLOAD_PATH="/home/sls_server/downloads"
NAME="sls_server_bin"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $NAME
scp ./$NAME root@$IP:$DOWNLOAD_PATH
scp ./env.json root@$IP:$DOWNLOAD_PATH
ssh root@$IP "cd $DOWNLOAD_PATH && mv $NAME .. && mv env.json .. && cd .. && ./$NAME && nohup ./bin > sls_server.log 2>&1 & exit"
rm $NAME