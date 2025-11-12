# imgHelper
改库提供图像处理和绘制；
设计画布，所有图像处理和绘制都在画布上进行；
设计最小粒度是图层，图层可复用在并发处理场景下效率提升；
设计图像操作方法ops，画布和图层都能使用；
包含了众多图像的处理算法，绘制则可进行配置方式进行绘制详细用法见Readme......

#### 设计

- canvas ： 画布, 图层绘制到画布上，支持各种绘制方法
- layer ：图层，进行图像处理和绘制最小单元,含有图像资源，文本，几何绘制，蒙层，操作画布
- ops ：图像处理操作，可以操作图层和画布

#### 先来一个例子

打开一张图片作为画布，然后打开第二张图缩放到200*200，旋转66度，绘制到画布的x0,y0=100,100的位置

```go

```

#### todo 

- 缩放图层
- 图层减法
- 图层乘法
- 图层除法
- 图层逻辑运算
- 打印当前画布的rgba
- 记录画布的图层,可以打印当前画布的图层

## 支持

#### 图像处理
- 缩放 Scale [v]

#### 图像绘制

#### 图层运算
- 加法
- 减法
- 乘法
- 除法

#### 图像分割

#### 图像压缩

#### 人像相关

## 简单例子

#### 创建一个画布读取一个图片放在创建的图层上进行缩放，最终输出图片
```go
func main() {
    cas := imgHelper.NewCanvas(400, 400) // 创建400*400的画布
    
    testImg, err := imgHelper.OpenImgFromLocalFile("./test.png") // 本地读取test.png图片
    if err != nil {
        log.Fatal(err)
    }
    
    imgLayer := &imgHelper.ImgLayer{ // 创建一个图层,从右上角x0y0=100*100的位置开始绘制
        Resource: imgHelper.Scale(testImg, 100, 100), // 将资源单独进行缩放到100*100
        X0:       100,
        Y0:       100,
    }
    _ = imgLayer.Scale(200, 200) //将图层整体缩放到200*200
    
    cas.AddLayer(imgLayer)       // 图层添加到画布
    _, _ = cas.Do()              // 执行绘制
    cas.SaveToFile("./case.png") // 存储到本地
}
```

#### 读取一个图片放在创建的图层上进行缩放，直接输出这个图层为图片
```go
func main() {
    src, err := imgHelper.OpenImgFromLocalFile("./test.png")
    if err != nil {
        log.Fatal(err)
    }
    imgLayer := imgHelper.ImgLayer{
        Resource: src,
    }
    _ = imgLayer.Scale(200, 200)
    imgLayer.Save("./case.png")
}
```

## 使用文档

#### 打开图像
```
- OpenImgFromLocalFile(imgPath string) (image.Image, error) // 从本地文件读取图像
- OpenImgFromReader(rd io.Reader) (image.Image, error) // 从Reader读取图像
- OpenImgFromBytes(data []byte) (image.Image, error)  // 从Bytes读取图像
- OpenImgFromHttpGet(imgUrl string) (image.Image, error) // http get请求下载url图像
```

#### 创建画布
```
- NewCanvas(width, height int) *CanvasContext // NewCanvas 透明背景的画布
- NewColorCanvas(width, height int, color color.RGBA) *CanvasContext // 指定颜色背景画布
- NewImgCanvas(resource image.Image) (*CanvasContext, error) // 指定图片背景画布,会使用图片的宽高
- NewImgCanvasFromSize(width, height int, resource image.Image) (*CanvasContext, error) // 指定图片背景画布,自定义宽高，会根据宽高调整背景图大小
```

#### 缩放
```
- Scale(src image.Image, targetWidth, targetHeight int) image.Image // 使用双线性插值算法将源图片拉伸或压缩到目标大小
```


