@startuml
start
:func Build(...)
sealer/build/local_builder.go;
note right
    初始化->获取build流程->按步骤执行流程
end note

:func initBuilder(...)
  sealer/build/local_builder.go;
  note right
    从kuberfile中设置相关参数
  end note

:func GetBuildPipeLine() -> :func for piplines( f())
  sealer/build/local_builder.go;
  note right
    从下面子流程中，获得pieline列表
  end note

split
:func for piplines( f())
sealer/build/local_builder.go;
note left
    执行pipline
end note
stop

split again
  :func InitImageSpec()
  sealer/build/local_builder.go;
  note right
    用parser解析kuberfile
    获得Image.Spec.Layers[0]
  end note

  :func PullBaseImageNotExist()
  sealer/build/local_builder.go;
  note right
    用ImageService.PullIfNotExist
    拉取layer[0],拉取基础镜像
  end note

  :func ExecBuild()
  sealer/build/local_builder.go;
  note right
    先updateBuilderLayers(me.Image)
    再for newlayers( cache,
    execCopyLayer,execOtherLayer)
  end note

  :func UpdateImageMetadata()
  sealer/build/local_builder.go;
  note right
      保存文件为image，分三步
  end note

  split
    :func setClusterFileToImage()
    sealer/build/local_builder.go;
    note right
     GetRawClusterFile()
     要么读取文件到scratch镜像中()
     要么读取Clusterfile定义的文件到layer中
     要么通过方法image.GetClusterFileFromImage
     从定义的基础镜像读取集群文件
     最后添加ImageAnnotations
    end note

  :func squashBaseImageLayerIntoCurrentImage()
  sealer/build/local_builder.go;
  note right
   把newlayer追加到baselayer中
  end note

  :func updateImageIDAndSaveImage()
    sealer/build/local_builder.go;
  note right
   生成ImageID
   ImageStore.Save()
  end note



end
@enduml