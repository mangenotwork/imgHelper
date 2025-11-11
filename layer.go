package imgHelper

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
)

//type LayerType int

//const (
//	Image    LayerType = 1 // 绘制图片
//	Text     LayerType = 2 // 绘制文字
//	Geometry LayerType = 3 // 绘制几何
//)

// CanvasContext 画布上下文
type CanvasContext struct {
	Dst       *image.RGBA
	LayerList []Layer
}

func (ctx *CanvasContext) SaveToFile(filePath string) {
	outputFile, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = outputFile.Close()
	}()

	// todo 判断图片类型，根据类型进行存储
	err = png.Encode(outputFile, ctx.Dst)
	if err != nil {
		log.Fatal(err)
	}
}

// NewCanvas 透明背景的画布
func NewCanvas(width, height int) *CanvasContext {
	imgContext := &CanvasContext{
		Dst:       image.NewRGBA(image.Rect(0, 0, width, height)),
		LayerList: make([]Layer, 0),
	}
	transparent := color.RGBA{}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			imgContext.Dst.Set(x, y, transparent)
		}
	}
	return imgContext
}

// NewColorCanvas 指定颜色背景画布
func NewColorCanvas(width, height int, color color.RGBA) *CanvasContext {
	canvasContext := &CanvasContext{
		Dst:       image.NewRGBA(image.Rect(0, 0, width, height)),
		LayerList: make([]Layer, 0),
	}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			canvasContext.Dst.Set(x, y, color)
		}
	}
	return canvasContext
}

// NewImgCanvas 指定图片背景画布,会使用图片的宽高
func NewImgCanvas(resource image.Image) (*CanvasContext, error) {
	bounds := resource.Bounds()
	canvasContext := &CanvasContext{
		Dst:       image.NewRGBA(bounds),
		LayerList: make([]Layer, 0),
	}
	draw.Draw(canvasContext.Dst, bounds, resource, bounds.Min, draw.Over)
	return canvasContext, nil
}

// NewImgCanvasFromSize 指定图片背景画布,自定义宽高，会根据宽高调整背景图大小
func NewImgCanvasFromSize(width, height int, resource image.Image) (*CanvasContext, error) {
	canvasContext := NewCanvas(width, height)
	imgLayer := &ImgLayer{
		Resource: resource,
	}
	_ = imgLayer.Scale(width, height)
	canvasContext.LayerList = append(canvasContext.LayerList, imgLayer)
	for _, v := range canvasContext.LayerList {
		_ = v.Draw(canvasContext)
	}
	return canvasContext, nil
}

// Layer 图层
type Layer interface {
	Draw(ctx *CanvasContext) error //
	Save(filePath string)
}

// ImgLayer 图层 - 图片，在画布上绘制图片
type ImgLayer struct {
	Resource image.Image // 图像资源
	X0       int
	Y0       int
	X1       int
	Y1       int
}

func (imgLayer *ImgLayer) Draw(ctx *CanvasContext) error {
	draw.Draw(
		ctx.Dst,
		image.Rect(imgLayer.X0, imgLayer.Y0, imgLayer.X1, imgLayer.Y1),
		imgLayer.Resource,
		image.Point{},
		draw.Over,
	)
	return nil
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

// TextLayer 图层 - 文字，在画布上绘制文字
type TextLayer struct {
	Str               string // 绘制的内容
	Size              float64
	DPI               float64
	Colour            color.Color
	MaxWidth          int           // 字的最大宽度
	FontGradient      bool          // 字体渐变
	FontGradientColor []color.Color // 字体渐变颜色
	FontAlign         string        // 字体
}

func (textLayer *TextLayer) Draw(ctx *CanvasContext) error {
	return nil
}
