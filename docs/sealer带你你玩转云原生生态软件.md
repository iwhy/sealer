**https://zhuanlan.zhihu.com/p/389552044 **
# 安装: 命令行和Clusterfile
可以把整个集群打成像docker镜像一样的镜像包，安装生态任何软件都可以sealer run xxx搞定，保障整个集群纬度的交付一致性。
是专有云和离线交付利器，当然普通开发者可以用它来快速实践云原生生态软件，比如你要安装k8s 或者prometheus或者高可用的mysql都可以做到一键搞定。

sealer不存在helm的镜像不一致的问题，所有依赖都会被打包，而且整个集群整体打包。
```shell
sealer run kuberetes:v1.19.9 \
--master 192.168.0.2,192.168.0.3,192.168.0.4 \
--node  192.168.0.5,192.168.0.6 -p 123456 \
-m 3 -n 3
```

对比sealos有没有发现命令行更简约更干净. 其中kuberentes:v1.19.9我们称之为集群镜像，它很神奇，和Docker镜像类似本质是一坨安装整个集群所需要的所有文件的集合，
在sealos里面可能就是个tar包，而sealer里面做了分层和兼容docker registry的设计，意味着我们可以把这个集群镜像放到docker registry里面进行交付。

对接公有云更简单，安装时只需要指定机器数量,指定AK和SK。
想定义更多参数？ 定义Clusterfile即可：


# sealer定义Kubefile,实现集群镜像交付
安装只是sealer的一个部分，sealer是一个《集群镜像》的实现，也就是如何通过一定的技术手段把整个集群打包！
这一点上相比sealos是一个质的提升,sealos只是一个安装工具。
sealer赋予了用户Build的能力，用非常简单的方式让用户进行自定义集群镜像:
![file](http://jeff.spring4all.com/Fk_pI9Ok3X-UiSXziYJt6ZER--r0)

我们想定义一个包含mysql ELK redis wordpress的集群，并把所有的依赖打包在一起，sealer就可以通过极简单方式帮助你做到这一点：
> 1. Kubefile定义集群镜像
```shell
FROM kuberentes:v1.19.9 # 集群镜像基础镜像，sealer官方提供 
COPY mysql . # mysql 编排文件 
COPY ELK . 
COPY redis . 
COPY wordpress . 
CMD kubectl apply -f . # 集群启动后执行的命令
```

> 2. Build自定义镜像

`sealer build -t mysql-redis-elk:latest .`

然后需要部署一个新集群只需要

`sealer run mysql-redis-elk:latest --master 192.168.0.2 -p 123456`

这个集群run完就包含了mysql redis等
你还可以把集群镜像推送到私有镜像仓库中：
`sealer login http://hub.docker.io -u xxx -p xxx sealer push mysql-redis-elk:latest`
还可以pull下来save成tar到客户环境中load:
```sealer pull mysql-redis-elk:latest 
sealer save -o mysql-redis-elk.tar mysql-redis-elk:latest 
sealer load -i mysql-redis-elk.tar # 客户离线环境
```

**从此以后，使用变得非常简单**
大部分云原生生态软件落地：
```
sealer run rook:latest 
sealer run prometheus:latest 
sealer run ingress:latest 
sealer run istio:latest
```
一切变得如此简单...


# sealer设计思想
sealer设计是极其优秀的，其实把整个集群制作成镜像并非一件简单的事，sealer的牛掰之处就在于把复杂的事通过优雅的设计让其足够大道至简，这也是几乎我所有产品的特点，牺牲复杂度换取的功能宁可不要。

> kubefile设计

这是核心亮点，它以一个非常简单的用户接口让用户实现自定义集群镜像的能力。
用怎样的一种描述语言可以用来描述整个集群所需要的文件，并且还要简单？
在sealer诞生之前这其实是个复杂的问题，受Dockerfile启发，为何不把单机容器镜像上升到集群纬度？

![file](http://jeff.spring4all.com/FgxTluVQKaltW03KqoUI01alpXQR)

于是便有了Kubefile。

docker可以通过Dockerfile构建一个docker镜像，使用compose就可以运行容器。

sealer通过Kubefile构建一个CloudImage,使用Clusterfile启动整个集群。

这是个非常亮眼的想法和设计。
### 那Kubefile中应该包含哪些指令？

### FROM 
`FROM kubernetes:v1.19.9`
FROM指定基础镜像，它可以是一个很干净的k8s基础镜像，也可以一个用户已经打包了一些服务的自定义镜像，对于使用者来说不需要关心里面的细节，就像用docker时不用关心centos rootfs里面有哪些文件一样。

### COPY && RUN
```
COPY my-chart . 
RUN wget helm.sh/download/helm-v3 && mv helm-v3 /usr/bin/helm
```
COPY指令可以像Docker一样把build工作目录的文件拷贝到集群镜像中。
RUN指令会在Build的时候执行，执行的过程中产生的文件都会缓存到集群镜像的一个layer中，比如上面在集群镜像中打包helm二进制

### CMD
`CMD helm install app my-chart`
CMD指令在k8s集群拉起后执行，它可以有多个。

这样在Build的过程中sealer会拉起一个临时的k8s集群，然后在里面执行Kubefile定义的指令，最终把这些指令产生的所有文件打包。

> 容器镜像缓存设计
缓存容器镜像可不是件轻松的事，这其中有一些难点问题:

如何知道分布式软件中有哪些容器镜像，因为我们需要把这些镜像缓存下来，
不管是扫描用户的yaml文件还是用helm template之后扫描都是不完美的，
首先不能确定用户的编排方式是什么，
其次有些软件甚至不把镜像地址写在编排文件中，而是通过自己的程序去拉起。
无法保证build成功运行就一定没问题。

容器镜像是需要被存储到私有仓库中打包在集群镜像里，那容器镜像仓库地址势必和编排文件中写的不一样，
特别是怎么保证用户alwayPull的时候还是能够在私有仓库中下载到镜像。

这里就体现sealer build的过程起一个**临时k8s集群**的优势了，最终集群会让docker去pull镜像，我们在pull的过程中拦截镜像并缓存，透明的支持了容器镜像存储

![file](http://jeff.spring4all.com/FiPujq6bfFvktTHJLntMC1LWv09e)

如此做到了你Build的产物一致性会非常好，到了其它环境部署几乎无需更改。

> 配置文件管理

很多交付场景会有大量的业务配置文件要向外透出，sealer可以非常友好的让用户透出这些配置到Clusterfile中。
典型的情况是用户希望集群镜像里面的helm values能够在部署时修改。

用户只需要在Clusterfile中定义一个Config即可：
```shell
apiVersion: http://sealer.aliyun.com/v1alpha1 
kind: Config 
metadata: name: mysql-values.yaml 
spec: path: etc/mysql-chart/values.yaml 
data: | mysql-user: root mysql-passwd: xxx
```
data中的内容就会覆盖掉默认的mysql chart的values

> 插件机制

还有一些场景比如希望通过sealer去修改主机名，或者升级内核，或者同步时间这些“本不该”由sealer去做的事情，
那么我们可以启用插件的方式来完成，以修改主机名插件为例：

```shell
apiVersion: http://sealer.aliyun.com/v1alpha1 
kind: Plugin 
metadata: 
  name: HOSTNAME 
spec: 
  data: | 192.168.0.2 master-0 192.168.0.3 master-1 192.168.0.4 master-2 192.168.0.5 node-0 192.168.0.6 node-1 192.168.0.7 node-2
```
只需要定义上面插件就可以帮助用户把集群中节点的主机名修改中data中定义的名字
当然还有一些其它的插件如打标签插件，执行shell命令插件等.

>不同runtime支持

未来你可以FROM k3s FROM k0s FROM ACK等等，而完全不用关心他们之间的安装差异。

>对接公有云

现在很多用户都希望在云端运行自己的集群镜像，sealer自带对接公有云能力，sealer自己实现的基础设施管理器，得益于我们更精细的退避重试机制，30s即可完成基础设施构建(阿里云6节点)性能是同类工具中的佼佼者，且API调用次数大大降低，配置兼容Clusterfile。

> 何种场景适合使用sealer

如果你要整体交付你的分布式SaaS，请用sealer

如果你要集成多个分布式服务在一起，如数据库消息队列或者微服务运行时，请用sealer

如果你要安装一个分布式应用如mysql主备集群，请用sealer

如果你需要安装/管理一个kubernetes高可用集群，请用sealer

如果你要初始化多个数据中心，保持多个数据中心状态强一致，请用sealer

如果你需要在公有云上实现上述场景，请用sealer

| 心得

sealer最值得我们自豪的地方是把复杂的东西变简单了，将近使用了一年的时间去思考User Interface，怎么在用户视角看不损失功能还能简单，这非常难，Kubefile的设计草稿被推翻了不知道多少次，千锤百炼最终打造出sealer这个项目。

还想说一下sealer和sealos的渊源，其实sealos是我很早开源的一个很受欢迎的项目，一步一步迭代，真的把安装k8s集群这件事做到了接近完美。然而有很多理由让我必须要做一个大的改变了：

sealos背后其实有一整套自动化构建离线包能力的平台，但是这些东西非常专用，基本是给我们自己发布新的离线包使用的而一般的开发者根本没有办法复用到这些能力，如何“优雅的开放这些能力”一直是我思考的。sealer完美的给出了答案！

可能用过sealos的知道sealos有个install命令，可以安装其它app如prometheus ingress dashboard这些，然而这块的设计我一直不满，但是又找不到更优雅的设计，只能说sealos的app包是没有技术含量的，首先镜像靠load，这样yaml里面always pull就凉了，其次打包麻烦，需要用户自己save镜像再tar，简直太low，而sealer一个Kubefile完全解决，这得益于底层镜像缓存技术的创新。

sealos不会产生一个生态，很简单我们做东西给开发者使用，是一个一对多的关系，而sealer的出现，任何人都可以成为生产者和消费者，生态的崛起才有可能。

sealos的代码现在看来简直就是一坨SHIT(自我批评一下)，很早的时候我只重视User interface而不重视用户看不到的地方，想着一个安装工具而已，随手写写，只要命令好用谁管你里面是什么样，PR只要回归测试没问题我也基本都全合并，现在证明这样是大错特错！这种心态会让你失去对自己作品的爱最终只能灭亡或者重构，好在sealer算是涅槃重生了。

我在怎么把事情变简单的道路上一直有着很多的探索和很多非常棒的想法，并且最终把它变为现实，sealos就是个例子，然而那只能是个业余兴趣项目，而sealer不一样，融入了更多的思考，以及整个社区的共同努力。当初大家还以为集群镜像只是个空洞的概念，而今天大家都可以切实体验到它变成了现实!