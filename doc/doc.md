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
- NewImgCanvas(resource image.Image) *CanvasContext // 指定图片背景画布,会使用图片的宽高
- NewImgCanvasFromSize(width, height int, resource image.Image) *CanvasContext // 指定图片背景画布,自定义宽高，会根据宽高调整背景图大小
- NewImgCanvasFromRange(rg Range, resource image.Image) *CanvasContext // 指定图片背景画布,取指定范围
- CanvasFromLocalImg(imgPath string) *CanvasContext // 指定本地一张图片作为画布的背景,画布的大小会使用图片的宽高

```

#### 画布方法
```
- CanvasContext.GetErr() error // 获取画布的错误
- CanvasContext.Ext(fn func(ctx *CanvasContext) error) *CanvasContext // 执行传入绘制的方法(操作ops)并接收绘制产生的错误
- CanvasContext.AddLayer(layer Layer) *CanvasContext  // 按顺序添加图层
- CanvasContext.Addition(layer Layer) *CanvasContext // 将当前图层加法添加到画布上，也就是与当前画布做加法
- CanvasContext.Subtraction(layer Layer, soft ...bool) *CanvasContext // 将当前图层减法添加到画布上，也就是与当前画布做减法, 可选参数soft true:柔和减法效果
- CanvasContext.Multiplication(layer Layer, normalize ...any) *CanvasContext // 将当前图层乘法添加到画布上，也就是与当前画布做乘法, 可选参数： - normalize: 是否归一化（默认true） - scale: 缩放因子（默认255，值越小保留的小值越多，如128会增强低亮度
- CanvasContext.Division(layer Layer, opts ...any) *CanvasContext // 将当前图层除法添加到画布上，也就是与当前画布做除法, 可选参数：- normalize: 是否归一化（默认true，结果映射到 0~255 范围）- scale: 缩放因子（默认255，影响归一化强度） - zeroVal: 图层像素为0时的替代结果（默认255，避免除零错误）
- CanvasContext.AND(layer Layer) *CanvasContext // 将当前图层与画布进行逻辑运算 - 与（AND）
- CanvasContext.OR(layer Layer) *CanvasContext // 将当前图层与画布进行逻辑运算 - 或（OR）
- CanvasContext.XOR(layer Layer) *CanvasContext // 将当前图层与画布进行逻辑运算 - 异或（XOR）
- CanvasContext.NOT(layer Layer) *CanvasContext // 将当前图层与画布进行逻辑运算 - 非（NOT）
- CanvasContext.SaveToFile(filePath string) error // 保存在本地文件
- CanvasContext.Print() // 在终端打印当前画布每个像素点的颜色值
```

#### 图层 - 图像图层与方法

```
- NewImgLayer(src image.Image, rg Range) *ImgLayer // 新建图像图层
- ImgLayerFromLocalFile(imgPath string, rg Range) (*ImgLayer, error) // 从本地打开一张图片作为图层放在画布的指定范围
- ImgLayerFromFromReader(rd io.Reader, rg Range) (*ImgLayer, error) // 从IO中打开一张图片作为图层放在画布的指定范围
- ImgLayer.Draw(ctx *CanvasContext) error // 执行将当前图像图层绘制到画布上
- ImgLayer.GetResource() image.Image // 获取当前图像图层的图像资源
- ImgLayer.GetXY() (int, int, int, int) // 获取当前图像图层的矩形范围
- ImgLayer.Scale(targetWidth, targetHeight int) error // 将当前图像图层进行缩放
- ImgLayer.Save(filePath string) error // 将当前图像图层保存到文件
- ImgLayer.Ext(fn func(ctx *CanvasContext) error) *ImgLayer // 执行传入绘制的方法(操作ops)并接收绘制产生的错误
- ImgLayer.Translation(dx, dy int) *ImgLayer // 将资源图像在图层上进行平移

```

#### 图层 - 文本图层与方法

todo...

#### 画布图层体系内使用Ext执行图像处理

```
例如灰度处理
- OpsGray() // 灰度处理

画布和图层的 Ext()使用, 

画布如: imgHelper.CanvasFromLocalImg("./test.png").Ext(imgHelper.OpsGray())

图层如: imgHelper.ImgLayerFromLocalFile("./test.png", imgHelper.Range{}).Ext(imgHelper.OpsGray())

```

#### 图像处理

- 灰度 Gray
```
- Gray(src image.Image) image.Image 
- OpsGray() // 画布和图层体系使用
```

- 二值图 BinaryImg
```
- BinaryImg(src image.Image, thresholdVal ...int) image.Image // 参数: thresholdVal阈值，通过这个阈值来划分二值，默认为128
- OpsBinaryImg(thresholdVal ...int) // 画布和图层体系使用
```

- 图像点的亮度调整 Brightness
```
- Brightness(src image.Image, brightnessVal int) image.Image
- OpsBrightness(brightnessVal int) // 画布和图层体系使用
```

- 图像转置 Transposition
```
- Transposition(src image.Image) image.Image 
- OpsTransposition() // 画布和图层体系使用
```

- 图像镜像 Mirror
```
- MirrorHorizontal(src image.Image) image.Image // 图像水平镜像
- MirrorVertical(src image.Image) image.Image // 图像垂直镜像
- OpsMirrorHorizontal() // 画布和图层体系使用
- OpsMirrorVertical() // 画布和图层体系使用
```

- 图像浮雕 Relief
```
- Relief(src image.Image) image.Image
- OpsRelief() // 画布和图层体系使用
```

- 裁剪 Crop
```
- Crop(src image.Image, x0, y0, x1, y1 int) image.Image // 矩形裁剪
- CropCircle(src image.Image, cx, cy, r int) image.Image // 圆形裁剪
- CropTriangle(src image.Image, x1, y1, x2, y2, x3, y3 int) image.Image // 三角形裁剪
- CropPolygon(src image.Image, points ...int) (image.Image, error) // 多边形裁剪
- OpsCrop(rg RangeValue) // 参数 rg RangeValues是范围（矩形，圆，三角形，多边形，曲边多边形） 画布和图层体系使用
```

- 马赛克 Mosaic
```
- Mosaic(src image.Image, x0, y0, x1, y1 int, blockSize int) image.Image // 矩形马赛克
- MosaicCircle(src image.Image, cx, cy, r int, blockSize int) image.Image // 圆形范围马赛克
- MosaicTriangle(src image.Image, x1, y1, x2, y2, x3, y3 int, blockSize int) image.Image // 三角形范围马赛克
- MosaicPolygon(src image.Image, blockSize int, points ...int) (image.Image, error) // 多边形范围马赛克
- OpsMosaic(rg RangeValue, blockSize int) // 参数 rg RangeValues是范围（矩形，圆，三角形，多边形，曲边多边形） 画布和图层体系使用
```

- 旋转 Rotate
```
- Rotate(src image.Image, angle float64) image.Image // angle是旋转度
- OpsRotate(angle float64) // 画布和图层体系使用
- OpsRotate90() // 画布和图层体系使用
- OpsRotate180() // 画布和图层体系使用
- OpsRotate270() // 画布和图层体系使用
```

- 伸缩 Scale
```
- Scale(src image.Image, targetWidth, targetHeight int) image.Image // 使用双线性插值算法将源图片拉伸或压缩到目标大小
- ScaleNearestNeighbor(src image.Image, targetWidth, targetHeight int) image.Image // 最近邻插值，速度快，但可能会导致图像出现锯齿。
- ScaleCatmullRom(src image.Image, targetWidth, targetHeight int) image.Image // 插值，质量较高，但速度较慢。
- OpsScale(targetWidth, targetHeight int) // 画布和图层体系使用
- OpsScaleNearestNeighbor(targetWidth, targetHeight int) // 画布和图层体系使用
- OpsScaleCatmullRom(targetWidth, targetHeight int) // 画布和图层体系使用
```

- 图像颜色反转 ColorReversal
```
- ColorReversal(src image.Image) image.Image 
- OpsColorReversal() // 画布和图层体系使用
```

- 图像腐蚀 Corrosion
```
- Corrosion(src image.Image) image.Image 
- OpsCorrosion() // 画布和图层体系使用
```

- 图像膨胀 Dilation
```
- Dilation(src image.Image) image.Image
- OpsDilation() // 画布和图层体系使用
```

- 图像的开运算 Opening
```
- Opening(src image.Image) image.Image
- OpsOpening()
```

- 图像的闭运算 Closing
```
- Closing(src image.Image) image.Image
- OpsClosing() // 画布和图层体系使用
```

- 调整色相 Hue 
```
- Hue(src image.Image, hueAdjustment float64) image.Image //  hueAdjustment : 色相调整值
- OpsHue(hueAdjustment float64) // 画布和图层体系使用
```

- 图像调整饱和度 Saturation
```
- Saturation(src image.Image, saturationAdjustment float64) image.Image // saturationAdjustment: 调整饱和度的值
- OpsSaturation(saturationAdjustment float64) // 画布和图层体系使用
```

- 调整色彩平衡 ColorBalance
```
- AdjustColorBalance(src image.Image, rAdjustment, gAdjustment, bAdjustment int) image.Image 
- OpsAdjustColorBalance(rAdjustment, gAdjustment, bAdjustment int) // 画布和图层体系使用
```

- 调整对比度 Contrast
```
- AdjustContrast(src image.Image, contrast float64) image.Image // contrast : 对比度调整值
- OpsAdjustContrast(contrast float64) // 画布和图层体系使用
```

- 调整锐度 Sharpness 
```
- AdjustSharpness(src image.Image, sharpness float64) image.Image // sharpness:锐度调整值
- OpsAdjustSharpness(sharpness float64) // 画布和图层体系使用
```

- 调整色阶 ColorScale 
```
- AdjustColorScale(src image.Image, blackPoint, whitePoint, gamma float64) image.Image // - blackPoint : 黑点  - whitePoint : 白点 - gamma : 伽马校正
- OpsAdjustColorScale(blackPoint, whitePoint, gamma float64) // 画布和图层体系使用
```

- 调整曝光度 Exposure 
```
- AdjustExposure(src image.Image, exposure float64) image.Image
- OpsAdjustExposure(exposure float64) // 画布和图层体系使用
```

- 调整色温 ColorTemperature
```
- ColorTemperature(src image.Image, temperature float64) image.Image
- OpsColorTemperature(temperature float64) // 画布和图层体系使用
```

- 调整色调 ColorTone
```
- ColorTone(src image.Image, adjustmentValue float64) image.Image
- OpsColorTone(adjustmentValue float64) // 画布和图层体系使用
```

- 图像降噪 Denoise
```
- Denoise(src image.Image, sigma float64) image.Image
- GaussianBlur1D(src image.Image, sigma float64) image.Image // 一维高斯模糊 可以用于降噪
- OpsDenoise(sigma float64) // 画布和图层体系使用
- OpsGaussianBlur1D(sigma float64) // 画布和图层体系使用
```

- 图像细化 Thinning
```
- Thinning(src image.Image) image.Image // 针对文本图像进行细化处理
- OpsThinning() // 画布和图层体系使用
```

- 仿射变换 Transform
```
- RigidTransform(img image.Image, angle, scale, tx, ty float64) *image.RGBA // 刚性变换（旋转、缩放、平移）
- OpsRigidTransform(angle, scale, tx, ty float64) 
- AffineTransform(img image.Image, mat [6]float64) *image.RGBA // 仿射变换
- OpsAffineTransform(mat [6]float64)
- PerspectiveTransform(img image.Image, mat [9]float64) *image.RGBA // 透视变换
- OpsPerspectiveTransform(mat [9]float64)
- AffineTransform23(img *image.RGBA, matrix [2][3]float64) *image.RGBA // 仿射变换通过 2x3 矩阵实现
- OpsAffineTransform23(matrix [2][3]float64)
- PerspectiveTransform33(img *image.RGBA, matrix [3][3]float64) *image.RGBA // 透视变换通过 3x3 矩阵实现
- OpsPerspectiveTransform33(matrix [3][3]float64)
```

- 彩色图像的平滑处理 SmoothProcessing
```
- SmoothProcessing(src image.Image, kernelSize int) image.Image  // kernelSize : 平滑处理的核大小
- OpsSmoothProcessing(kernelSize int)
```


