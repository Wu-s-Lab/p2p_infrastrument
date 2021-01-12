
HOST=`uname`
WORKING_DIR=`pwd`
EDGE_CONF_FILE="edge.conf"
EDGE_CONF_PATH=$WORKING_DIR/$EDGE_CONF_FILE

IPCERT_URL="8.131.232.232:9876"

EDGE_BIN="edge"
if [[ `uname` =~ "MINGW" ]];then
    EDGE_BIN="edge.exe"
    TAP_BIN="tapinstaller.exe"
    TAP_PATH=$WORKING_DIR/$TAP_BIN
fi


# 下载文件
EDGE_PATH=$WORKING_DIR/$EDGE_BIN


if [ $HOST = "Darwin" ];then
    echo "正在使用 macos 系统"
    EDGE_BIN_URL="https://p2p-infra.oss-cn-beijing.aliyuncs.com/edge_v.2.9.0.r810.0bdab26_x86_64-apple-darwin19.6.0"
    if [ ! -f $EDGE_PATH ];then
        curl -o $EDGE_PATH $EDGE_BIN_URL
        chmod +x $EDGE_PATH
    fi

elif [[ $HOST =~ "MINGW" ]];then
    if [[ $HOST =~ "MINGW64" ]] ;then
        echo "正在使用 windows 64 位 系统"
        EDGE_BIN_URL="https://p2p-infra.oss-cn-beijing.aliyuncs.com/edge_v.2.9.0.r810.0bdab26_Windows-10.0.18363.exe"
        TAP_URL="https://p2p-infra.oss-cn-beijing.aliyuncs.com/tap-windows-9.24.2-I601-Win10.exe"
    elif [[ $HOST =~ "MINGW32" ]]; then
        echo "正在使用 windows 32 位 系统, 暂不支持"
        exit 1
        
    else
        echo "未知的windows系统，暂不支持"
        exit 1

    fi

    if [ ! -f $EDGE_PATH ];then
        curl -o $EDGE_PATH $EDGE_BIN_URL
        chmod +x $EDGE_PATH

    fi

    if [ ! -f $TAP_PATH ];then
        curl -o $TAP_PATH $TAP_URL
        chmod +x $TAP_PATH
        #安装tuntap设备
        $TAP_PATH
    fi
elif [ $HOST = "Linux" ];then
    echo "正在使用 Linux 系统"
else 
    echo "正在使用 $HOST 暂不支持"
    exit 1
fi



# 启动 docker

DOCKER_REGISTRY_URL=registry.cn-beijing.aliyuncs.com
DOCKER_DOMAIN=knowonchain
DOCKER_TAG=p2p_node

DOCKER=$DOCKER_REGISTRY_URL/$DOCKER_DOMAIN/$DOCKER_TAG


HOST_SSH_PORT=22222
EDGE_DIR=`pwd`
EDGE_CONF_DIR=`pwd`



docker login --username=杭研院区块链and郑 -p blockchain1  $DOCKER_REGISTRY_URL
docker pull $DOCKER
docker tag $DOCKER $DOCKER_TAG


if [ ! -f $EDGE_CONF_PATH ];then
    echo ""
    echo "=============节点注册================="
    echo "请输入您的名字"
    read NAME
    echo " $NAME 你好, 请输入你的手机号码"
    read PHONE
fi

if [ $HOST = 'Linux' ];then
    docker run -itd \
--network host \
--device /dev/net/tun \
--cap-add NET_ADMIN \
-v /var/run/docker.sock:/var/run/docker.sock \
-v $WORKING_DIR:/edge \
-e HOST=`uname` \
-e NAME=$NAME \
-e PHONE=$PHONE \
-p $HOST_SSH_PORT:22222 \
-e IPCERT_URL=$IPCERT_URL \
--restart always \
--name p2p_node \
$DOCKER_TAG
elif [ $HOST = 'Darwin' ];then
    docker run -itd \
-v /var/run/docker.sock:/var/run/docker.sock \
-v $WORKING_DIR:/edge \
-e HOST=`uname` \
-e NAME=$NAME \
-e PHONE=$PHONE \
-p $HOST_SSH_PORT:22222 \
-e IPCERT_URL=$IPCERT_URL \
--restart always \
--name p2p_node \
$DOCKER_TAG

elif [[ $HOST =~ 'MINGW' ]];then
    docker run -itd \
-v //var/run/docker.sock:/var/run/docker.sock \
-e HOST=`uname` \
-e NAME=$NAME \
-e PHONE=$PHONE \
-p $HOST_SSH_PORT:22222 \
-e IPCERT_URL=$IPCERT_URL \
--restart always \
--name p2p_node \
$DOCKER_TAG

if [ ! -f $EDGE_CONF_PATH ];then
    sleep 6
    docker cp p2p_node:/edge/edge.conf $EDGE_CONF_PATH

    if [ $? -ne 0 ];then
        echo '无法获得配置文件'
        docker stop p2p_node
        docker rm p2p_node
        exit 1
    fi
fi


fi

sleep 5
if [ ! -f $EDGE_CONF_PATH ]; then
    docker stop p2p_node
    echo "身份获取失败，无法取得配置文件,请联系fjl看看问题出哪里啦"
    echo "================日志为=================\n"
    docker logs p2p_node
    echo "======================================\n"


    docker rm p2p_node
    exit 1
fi


echo "控制节点已启动"


#启动HOST 的edge

echo "正在启动p2p网络"

if [ `uname` = 'Darwin' ]; then
    echo "在宿主机启动p2p vpn"
    if [ ! -f $EDGE_PATH ]; then
        echo "找不到edge二进制"
        exit -1
    fi
    echo "请输入管理员密码以使用sudo"
    sudo $EDGE_PATH $EDGE_CONF_PATH -f
fi

if [[ `uname` =~ 'MINGW' ]];then
    echo "在宿主机启动p2p vpn"
    if [ ! -f $EDGE_PATH ]; then
        echo "找不到edge二进制"
        exit -1
    fi
    $EDGE_PATH $EDGE_CONF_PATH -f
fi

