package imgHelper

import (
	"errors"
	"image"
	"image/draw"
	"image/png"
	"io"
	"log"
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

// NewImgLayerFromLocalFile 从本地打开一张图片作为图层放在画布的指定范围
func NewImgLayerFromLocalFile(imgPath string, rg Range) (*ImgLayer, error) {
	resource, err := OpenImgFromLocalFile(imgPath)
	if err != nil {
		log.Fatal(err)
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

func NewImgLayerFromFromReader(rd io.Reader, rg Range) (*ImgLayer, error) {
	resource, err := OpenImgFromReader(rd)
	if err != nil {
		log.Fatal(err)
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

func (imgLayer *ImgLayer) GetResource() image.Image {
	return imgLayer.Resource
}

func (imgLayer *ImgLayer) GetXY() (int, int, int, int) {
	return imgLayer.X0, imgLayer.Y0, imgLayer.X1, imgLayer.Y1
}

func (imgLayer *ImgLayer) Scale(targetWidth, targetHeight int) error {
	imgLayer.Resource = Scale(imgLayer.Resource, targetWidth, targetHeight)
	imgLayer.X1 = imgLayer.X0 + imgLayer.Resource.Bounds().Dx()
	imgLayer.Y1 = imgLayer.Y0 + imgLayer.Resource.Bounds().Dy()
	return nil
}

func (imgLayer *ImgLayer) Save(filePath string) {
	outputFile, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = outputFile.Close()
	}()

	// todo 判断图片类型，根据类型进行存储
	err = png.Encode(outputFile, imgLayer.Resource)
	if err != nil {
		log.Fatal(err)
	}
}

func (imgLayer *ImgLayer) Ext(fn func(ctx *CanvasContext) error) *ImgLayer {
	nowImgLayerCtx := &CanvasContext{
		Dst: imgLayer.Resource.(*image.RGBA),
	}
	nowImgLayerCtx.Err = errors.Join(nowImgLayerCtx.Err, fn(nowImgLayerCtx))
	imgLayer.Resource = nowImgLayerCtx.Dst
	return imgLayer
}
