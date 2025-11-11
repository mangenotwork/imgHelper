# imgHelper
改库提供图像处理和绘制；设计最小粒度是图层，图层可复用在并发处理场景下效率提升；包含了众多图像的处理算法，绘制则可进行配置方式进行绘制详细用法见Readme......

## 支持

#### 图像处理
- 缩放

#### 图像绘制

#### 图层运算
- 加法
- 减法
- 乘法
- 除法

#### 图像分割

#### 图像压缩

#### 人像相关


## 方法文档

#### 打开图像

- OpenImgFromLocalFile(imgPath string) (image.Image, error) // 从本地文件读取图像
- OpenImgFromReader(rd io.Reader) (image.Image, error) // 从Reader读取图像
- OpenImgFromBytes(data []byte) (image.Image, error)  // 从Bytes读取图像
- OpenImgFromHttpGet(imgUrl string) (image.Image, error) // http get请求下载url图像

#### 创建画布

- NewCanvas(width, height int) *CanvasContext // NewCanvas 透明背景的画布
- NewColorCanvas(width, height int, color color.RGBA) *CanvasContext // 指定颜色背景画布
- NewImgCanvas(resource image.Image) (*CanvasContext, error) // 指定图片背景画布,会使用图片的宽高
- NewImgCanvasFromSize(width, height int, resource image.Image) (*CanvasContext, error) // 指定图片背景画布,自定义宽高，会根据宽高调整背景图大小

#### 缩放

- Scale(src image.Image, targetWidth, targetHeight int) image.Image // 使用双线性插值算法将源图片拉伸或压缩到目标大小



