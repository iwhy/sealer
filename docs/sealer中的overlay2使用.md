**https://zhuanlan.zhihu.com/p/393534044** 

**https://github.com/alibaba/sealer/blob/main/utils/mount/overlay2.go**

>写时复制 COPY on Write

COPY on Write在很多地方都会用到，比如git的存储，去修改一个文件并非真的修改，
而是把文件拷贝出来修改，读的时候读取上层的修改层，同样比如在ceph的后段存储等系统中都存在应用。

好处是可以非常方便的进行回滚，以及存储复用，比如给虚拟机做快照，那我们并不需要把整个系统盘复制一份，
而只需要用指针指一下快照的地方就好。
容器镜像同理，写时复制可以让多个应用共享一个基础镜像。

>集群镜像中的写时复制

集群镜像把整个集群打包，用写时复制的方法可以帮助分布式应用共享k8s基础镜像。
以及可以很方便的把各种分布式应用镜像进行融合。

> 操作系统中使用overlay2

`mount -t overlay overlay -o lowerdir=./lower,upperdir=./upper,workdir=./work ./merged`

overlay文件系统分为lowerdir、upperdir、merged， 对外统一展示为merged，uperdir和lower的同名文件会被upperdir覆盖

workdir必须和upperdir是mount在同一个文件系统下， 而lower不是必须的

`mount -t overlay overlay -o lowerdir=/lower1:/lower2:/lower3,upperdir=/upper,workdir=/work /merged`
lowerdir可以是多个目录(前面覆盖后面)，如果没有upperdir,那merged是read only.

overlay只支持两层，upper文件系统通常是可写的；lower文件系统则是只读，这就表示着，当我们对 overlay 文件系统做任何的变更，都只会修改 upper 文件系统中的文件