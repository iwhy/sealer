# Sealer roadmap

## v0.1.3

* general e2e tests
* support build, apply and push a CloudImage

## **v0.2.0**

* support cluster ENV and global config
* support ARM base image
* registry cache docker image when pull images
* base image support more kubernetes versions and auto make rootfs
* official dashboard image
* image merge feature
* **yingfan 2021/8/27** containerd作为运行是，因为k8s 1.20之后放弃dockershim
* **yingfan  2021/8/27** 支持aws，其它共有云

## future versions

* cluster upgrade, backup, and restore
* application upgrade
* multi cloud provider
* multi runtime
* rootfs mount filesystem
* sealer hub UI

## official registry opensource cloud images

- [ ] dashboard, https://github.com/kubernetes/dashboard
- [ ] prometheus stack, https://github.com/prometheus-operator/kube-prometheus
- [ ] loki stack
- [ ] mysql
- [ ] redis
- [ ] rocketmq
- [ ] zookeeper
- [ ] minio
- [ ] openEBS
- [ ] rook ceph
- [ ] kubeflow
- [ ] kafka
- [ ] cassandra
- [ ] cockroachDB
- [ ] postgreSQL
- [ ] tiDB
- [ ] istio
- [ ] dapr
- [ ] ingress
- [ ] gitea/drone/harbor, devops stack

# LONG-TERM


**yingfan  2021/8/27**
就目前来讲，sealer还有很多不完善的功能
，但是已经具备基础核心的一键部署功能：  
做成套件：
一键安装，尝试版，免费， 只有k8s和dashboard，，
一键安装 istio，基础版，完整版
一键安装devops战
一键安装监控 ETL，promethus，grafana等等
一键安装 MQ
一键安装 redis，mysql





