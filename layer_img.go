package imgHelper

import (
	"errors"
	"image"
	"image/draw"
	"image/png"
	"io"
	"os"
)

// ImgLayer 图层 - 图片，在画布上绘制图片
type ImgLayer struct {
	Resource image.Image // 图像资源
	X0       int
	Y0       int
	X1       int
	Y1       int
}

func NewImgLayer(src image.Image, rg Range) *ImgLayer {
	layer := &ImgLayer{
		Resource: src,
		X0:       rg.X0,
		Y0:       rg.Y0,
		X1:       rg.X1,
		Y1:       rg.Y1,
	}
	if layer.X1 == 0 {
		layer.X1 = rg.X0 + src.Bounds().Dx()
	}
	if layer.Y1 == 0 {
		layer.Y1 = rg.Y0 + src.Bounds().Dy()
	}
	return layer
}

// ImgLayerFromLocalFile 从本地打开一张图片作为图层放在画布的指定范围
func ImgLayerFromLocalFile(imgPath string, rg Range) (*ImgLayer, error) {
	resource, err := OpenImgFromLocalFile(imgPath)
	if err != nil {
		return nil, err
	}
	layer := &ImgLayer{
		Resource: resource,
		X0:       rg.X0,
		Y0:       rg.Y0,
		X1:       rg.X1,
		Y1:       rg.Y1,
	}
	if layer.X1 == 0 {
		layer.X1 = rg.X0 + resource.Bounds().Dx()
	}
	if layer.Y1 == 0 {
		layer.Y1 = rg.Y0 + resource.Bounds().Dy()
	}
	return layer, nil
}

func ImgLayerFromFromReader(rd io.Reader, rg Range) (*ImgLayer, error) {
	resource, err := OpenImgFromReader(rd)
	if err != nil {
		return nil, err
	}
	return &ImgLayer{
		Resource: resource,
		X0:       rg.X0,
		Y0:       rg.Y0,
		X1:       rg.X1,
		Y1:       rg.Y1,
	}, nil
}

// Draw 执行将当前图像图层绘制到画布上
func (imgLayer *ImgLayer) Draw(ctx *CanvasContext) error {
	if imgLayer.X1 == 0 {
		imgLayer.X1 = imgLayer.X0 + imgLayer.Resource.Bounds().Dx()
	}
	if imgLayer.Y1 == 0 {
		imgLayer.Y1 = imgLayer.Y0 + imgLayer.Resource.Bounds().Dy()
	}
	draw.Draw(
		ctx.Dst,
		image.Rect(imgLayer.X0, imgLayer.Y0, imgLayer.X1, imgLayer.Y1),
		imgLayer.Resource,
		image.Point{},
		draw.Over,
	)
	return nil
}

// GetResource 获取当前图像图层的图像资源
func (imgLayer *ImgLayer) GetResource() image.Image {
	return imgLayer.Resource
}

// GetXY 获取当前图像图层的矩形范围
func (imgLayer *ImgLayer) GetXY() (int, int, int, int) {
	return imgLayer.X0, imgLayer.Y0, imgLayer.X1, imgLayer.Y1
}

// Scale 将当前图像图层进行缩放
func (imgLayer *ImgLayer) Scale(targetWidth, targetHeight int) error {
	imgLayer.Resource = Scale(imgLayer.Resource, targetWidth, targetHeight)
	imgLayer.X1 = imgLayer.X0 + imgLayer.Resource.Bounds().Dx()
	imgLayer.Y1 = imgLayer.Y0 + imgLayer.Resource.Bounds().Dy()
	return nil
}

// Save 将当前图像图层保存到文件
func (imgLayer *ImgLayer) Save(filePath string) error {
	outputFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func() {
		_ = outputFile.Close()
	}()

	// todo 判断图片类型，根据类型进行存储
	return png.Encode(outputFile, imgLayer.Resource)
}

// Ext 执行传入绘制的方法(操作ops)并接收绘制产生的错误
// 只要是实现了  fn func(ctx *CanvasContext) error 方法就可以调用此方法
func (imgLayer *ImgLayer) Ext(fn func(ctx *CanvasContext) error) *ImgLayer {
	bounds := imgLayer.Resource.Bounds()
	rgbaImg := image.NewRGBA(bounds)
	draw.Draw(rgbaImg, bounds, imgLayer.Resource, bounds.Min, draw.Src)
	nowImgLayerCtx := &CanvasContext{
		Dst: rgbaImg,
	}
	nowImgLayerCtx.Err = errors.Join(nowImgLayerCtx.Err, fn(nowImgLayerCtx))
	imgLayer.Resource = nowImgLayerCtx.Dst
	return imgLayer
}

// Translation 将资源图像在图层上进行平移
func (imgLayer *ImgLayer) Translation(dx, dy int) *ImgLayer {
	bounds := imgLayer.Resource.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	draw.Draw(dst, dst.Bounds(), &image.Uniform{C: image.Transparent}, image.Point{}, draw.Src)
	for y := 0; y < bounds.Dy(); y++ {
		for x := 0; x < bounds.Dx(); x++ {
			newX := x + dx
			newY := y + dy
			if newX >= 0 && newX < bounds.Dx() && newY >= 0 && newY < bounds.Dy() {
				dst.Set(newX, newY, imgLayer.Resource.At(x, y))
			}
		}
	}
	imgLayer.Resource = dst
	return imgLayer
}
