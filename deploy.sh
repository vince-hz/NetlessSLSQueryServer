IP=$1
PORT="8080"
SERVER_DIR_PATH="/home/sls_server"
DOWNLOAD_PATH="$SERVER_DIR_PATH/downloads"
NAME="sls_server_bin"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $NAME
scp ./$NAME root@$IP:$DOWNLOAD_PATH
scp ./env.json root@$IP:$DOWNLOAD_PATH
ssh -f root@$IP "cd $DOWNLOAD_PATH && mv $NAME .. && mv env.json .. && netstat -tunlp|grep $PORT | awk '{print \$7}' | grep -o '[0-9]*' | xargs -r kill && cd $SERVER_DIR_PATH && nohup ./$NAME release > ./sls_server.log 2>&1 &"
rm $NAME