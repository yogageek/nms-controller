#!/bin/bash
# 宣告使用 /bin/bash

#1.0.0 groups api只提供imec使用
#1.0.1 groups api提供所有&普羅米修斯重構
#1.0.2 config api重構 打post後關閉服務(在k8s上會自動重啟)
#撰寫if dev選擇部署

#執行 ./dockerbu.sh true or false

# 手動修改佈署參數
VERSION=1.0.2
IMAGE=controller
# NAME=nms-controller-dev
NAME=nms-controller
KILLDEPLOY=$1
#打開則部屬正式環境
PLATFORM=release

echo "VERSION: ${VERSION} !"
echo "部屬名稱: ${NAME} !"
echo "佈署環境: ${PLATFORM} !"
echo "砍掉重部: ${KILLDEPLOY} !"


# 在函數中取得傳入參數
function redeploy() {   
    kubectl delete deployment.apps/${NAME}
    kubectl delete service/${NAME}-svc
    kubectl delete ingress/${NAME}-ing
    kubectl apply -f service.yaml
    kubectl apply -f pod.yaml
    kubectl apply -f ingress.yaml
    # kubectl apply -f podlihong.yaml
}

function getPodId(){    
    echo "取得podid"
    POD_ID=$(kubectl -o name get pods | grep ${NAME})
    echo ${POD_ID}
}

function deployPod() {
    getPodId    
    kubectl delete ${POD_ID}
}

function showLog() {
    getPodId
    kubectl logs ${POD_ID} -f
}

# -----------------run---------------

# docker build -t iiicondor/${IMAGE}:${VERSION} .
# docker push iiicondor/${IMAGE}:${VERSION}

# 選擇部屬檔
if [[ $PLATFORM = "release" ]]
then
    echo "--------進入dep資料夾..."
    cd ~/Desktop/nms-controller/dep    
else
    echo "--------進入dep_dev資料夾..."    
    cd ~/Desktop/nms-controller/dep_dev
fi

# 是否全部重部
if [[ $KILLDEPLOY = "true" ]]
then
    echo "全部重部..."
    redeploy
    echo "wait for 10 seconds..."
    sleep 10
    showLog
else
    echo "只刪除POD..."
    deployPod
    echo "wait for 10 seconds..."
    sleep 10
    showLog
fi






#部屬步驟

### 建立image
# docker build -t iiicondor/controller:1.0.0 .
# docker push iiicondor/controller:1.0.0

### 推上k8s
# 1.修改pod.yaml版本
# IMAGE: iiicondor/controller:1.0.0
# 2.apply
# kubectl apply -f pod.yaml -n nms

# (首次部屬才需要)
# kubectl apply -f service.yaml -n nms
# kubectl apply -f ingress.yaml -n nms

# 更換版本
# kubectl delete deployment.apps/nms-controller
