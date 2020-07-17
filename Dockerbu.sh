# 宣告使用 /bin/bash
#!/bin/bash

image=controller
version=1.0.0

echo "controller version=" ${version}
docker build -t iiicondor/${image}:${version} . 
docker push iiicondor/${image}:${version}

# cd ~/Desktop/nms-controller/dep
# kubectl delete deployment.apps/${image}
# kubectl delete service/${image}
# kubectl apply -f service.yaml 
# kubectl apply -f pod.yaml 
# kubectl apply -f podlihong.yaml 

echo "wait for 10 seconds..."
sleep 10

POD_ID=$(kubectl -o name get pods | grep ${image})
kubectl delete POD_ID
kubectl logs ${POD_ID} -f

#部屬步驟

### 建立image
# docker build -t iiicondor/controller:1.0.0 .
# docker push iiicondor/controller:1.0.0

### 推上k8s
# 1.修改pod.yaml版本
# image: iiicondor/controller:1.0.0
# 2.apply
# kubectl apply -f pod.yaml -n nms

# (首次部屬才需要)
# kubectl apply -f service.yaml -n nms
# kubectl apply -f ingress.yaml -n nms