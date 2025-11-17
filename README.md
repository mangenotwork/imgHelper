# imgHelper
该库提供图像处理和绘制；
设计画布，所有图像处理和绘制都在画布上进行；
设计最小粒度是图层，图层可复用在并发处理场景下效率提升；
设计图像操作方法ops，画布和图层都能使用；
包含了众多图像的处理算法，绘制则可进行配置方式进行绘制详细用法见Readme......

### 设计

- canvas ： 画布, 图层绘制到画布上，支持各种绘制方法
- layer ：图层，进行图像处理和绘制最小单元,含有图像资源，文本，几何绘制，蒙层，操作画布
- ops ：图像处理操作，可以操作图层和画布
- func : 各种图像处理方法函数，输入 image.Image 处理后 输出 image.Image


### 使用文档

[打开文档](https://github.com/mangenotwork/imgHelper/blob/main/doc/doc.md "点击跳转到文档")

### 先来一个例子

打开一张图片作为画布，然后打开第二张图为一个图层缩放到100*200，旋转66度，绘制到画布的x0,y0=100,100的位置，然后画布整体旋转90度

```go
func case13() {
	// 创建画布
	cas := imgHelper.CanvasFromLocalImg("./test.png")
	// 打开一张图片作为图层
	imgLayer, err := imgHelper.ImgLayerFromLocalFile("./case6.png", imgHelper.Range{X0: 100, Y0: 100})
	if err != nil {
		log.Fatal(err)
	}
	// 图层执行缩放和旋转操作
	imgLayer.Ext(imgHelper.OpsScale(100, 200)).Ext(imgHelper.OpsRotate(66))
	// 图层画绘制到画布上并保存到图片文件
        err = cas.AddLayer(imgLayer).Ext(imgHelper.OpsRotate90()).SaveToFile("./case13.png")
	if err != nil {
		log.Fatal(err)
	}
}
```

### 调用灵活

该库提供很多图像处理方法函数，可以将图层的Resource进行各种处理

```go
func case14() {
    cas := imgHelper.CanvasFromLocalImg("./test.png")
	imgLayer, err := imgHelper.ImgLayerFromLocalFile("./case6.png", imgHelper.Range{X0: 100, Y0: 100})
	if err != nil {
		log.Fatal(err)
	}
	// 缩放函数将图层的图片资源进行缩放
	imgHelper.Scale(imgLayer.Resource, 100, 200)
	// 旋转函数将图层的图片资源进行旋转
	imgHelper.Rotate(imgLayer.Resource, 66)
	// 图层画绘制到画布上并保存到图片文件
	err = cas.AddLayer(imgLayer).Ext(imgHelper.OpsRotate90()).SaveToFile("./case14.png")
	if err != nil {
		log.Fatal(err)
	}
}
```

### 只想简单处理图片 - 直接操作图层

读取一个图片为图层进行缩放，直接输出这个图层为图片

```go
func case5() {
    imgLayer, err := imgHelper.ImgLayerFromLocalFile("./test.png", imgHelper.Range{})
    if err != nil {
        log.Fatal(err)
    }
    imgLayer.Ext(imgHelper.OpsScale(100, 100)).Save("./case5.png")
}
```

### 只想使用该库的图像处理方法

读取一个图片进行缩放到100*100然后保存

```go
func case15() {
    file, err := os.Open("./test.png")
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()
    src, err := png.Decode(file)
    if err != nil {
        log.Fatal(err)
    }
    // 调用图像缩放方法
    dst := imgHelper.Scale(src, 100, 100)
    outputFile, err := os.Create("./case15.png")
    if err != nil {
        log.Fatal(err)
    }
    defer outputFile.Close()
    err = png.Encode(outputFile, dst)
    if err != nil {
        log.Fatal(err)
    }
}
```

### todo 还支持绘制

几何绘制

文本绘制

### todo 制点特别的

绘制渐变文本

绘制蒙层

### todo 该库还有很多功能，如图像分割，人像处理，相似度计算......

### todo gif相关的处理这里也有

### todo 更多的请见使用文档或引用该库进行体验

### todo 如果没用满足你的需求也没关系，你可以拉下代码自行修改

### 该库的依赖如下，协议均为MIT放心使用和学习，对这个项目感兴趣可以找我学习交流，邮箱: 2912882908@qq.com

```
require golang.org/x/image v0.32.0
```

#### todo 

- 记录画布的图层,可以打印当前画布的图层

## 支持

#### 图像处理
- 灰度 Gray
- 缩放 Scale  
- 

#### 图像绘制

#### 图层运算
- 加法
- 减法
- 乘法
- 除法

#### 图像分割

#### 图像压缩

#### 人像相关






