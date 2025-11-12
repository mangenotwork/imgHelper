package imgHelper

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
)

// 设计: 画布只接收图层，能实现图层接口的实体都能被绘制在画布上;

// CanvasContext 画布上下文
type CanvasContext struct {
	// 存放画布的rgba
	Dst *image.RGBA

	// 为了方便每个执行都能链式调用，所以这里设计用errors.Join接收多个错误
	// 在最终IO输出的时候抛出
	Err error

	// todo 记录画布的图层
}

// NewCanvas 透明背景的画布
func NewCanvas(width, height int) *CanvasContext {
	imgContext := &CanvasContext{
		Dst: image.NewRGBA(image.Rect(0, 0, width, height)),
		//LayerList: make([]Layer, 0),
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
		Dst: image.NewRGBA(image.Rect(0, 0, width, height)),
	}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			canvasContext.Dst.Set(x, y, color)
		}
	}
	return canvasContext
}

// NewImgCanvas 指定图片背景画布,会使用图片的宽高
func NewImgCanvas(resource image.Image) *CanvasContext {
	bounds := resource.Bounds()
	canvasContext := &CanvasContext{
		Dst: image.NewRGBA(bounds),
	}
	draw.Draw(canvasContext.Dst, bounds, resource, bounds.Min, draw.Over)
	return canvasContext
}

// NewImgCanvasFromSize 指定图片背景画布,自定义宽高，会根据宽高调整背景图大小
func NewImgCanvasFromSize(width, height int, resource image.Image) *CanvasContext {
	canvasContext := NewCanvas(width, height)
	imgLayer := &ImgLayer{
		Resource: resource,
	}
	canvasContext.Err = errors.Join(canvasContext.Err, imgLayer.Scale(width, height))
	canvasContext.Ext(imgLayer.Draw)
	return canvasContext
}

// NewImgCanvasFromRange 指定图片背景画布,取指定范围
func NewImgCanvasFromRange(rg Range, resource image.Image) *CanvasContext {
	canvasContext := &CanvasContext{}
	if rg.X0 >= rg.X1 || rg.Y0 >= rg.Y1 {
		err := fmt.Errorf("无效的 Range：X0 >= X1 或 Y0 >= Y1（%+v）", rg)
		canvasContext.Err = errors.Join(canvasContext.Err, err)
		return canvasContext
	}
	imgBounds := resource.Bounds()
	if rg.X0 < imgBounds.Min.X || rg.X1 > imgBounds.Max.X || rg.Y0 < imgBounds.Min.Y || rg.Y1 > imgBounds.Max.Y {
		err := fmt.Errorf("range 超出图片范围：图片边界 %v，Range %+v", imgBounds, rg)
		canvasContext.Err = errors.Join(canvasContext.Err, err)
		return canvasContext
	}
	width := rg.X1 - rg.X0
	height := rg.Y1 - rg.Y0
	dstBounds := image.Rect(0, 0, width, height)
	dst := image.NewRGBA(dstBounds)
	srcRect := image.Rect(rg.X0, rg.Y0, rg.X1, rg.Y1)
	draw.Draw(dst, dstBounds, resource, srcRect.Min, draw.Src)
	return canvasContext
}

// CanvasFromLocalImg 指定本地一张图片作为画布的背景,画布的大小会使用图片的宽高
func CanvasFromLocalImg(imgPath string) *CanvasContext {
	canvasContext := &CanvasContext{}
	resource, err := OpenImgFromLocalFile(imgPath)
	if err != nil {
		canvasContext.Err = errors.Join(canvasContext.Err, err)
		return canvasContext
	}
	bounds := resource.Bounds()
	canvasContext.Dst = image.NewRGBA(bounds)
	draw.Draw(canvasContext.Dst, bounds, resource, bounds.Min, draw.Over)
	return canvasContext
}

// Ext 执行传入绘制的方法并接收绘制产生的错误
// 只要是实现了  fn func(ctx *CanvasContext) error 方法的图层都可以调用此方法
// 返回画布上下文已支持链式调用
func (ctx *CanvasContext) Ext(fn func(ctx *CanvasContext) error) *CanvasContext {
	ctx.Err = errors.Join(ctx.Err, fn(ctx))
	return ctx
}

// AddLayer 按顺序添加图层
func (ctx *CanvasContext) AddLayer(layer Layer) *CanvasContext {
	ctx.Ext(layer.Draw)
	return ctx
}

// Addition 将当前图层加法添加到画布上，也就是与当前画布做加法
// 可以用于合成图像或增加图像的亮度
func (ctx *CanvasContext) Addition(layer Layer) *CanvasContext {
	src := layer.GetResource()
	srcBounds := src.Bounds()
	x0, y0, _, _ := layer.GetXY()
	dstBounds := ctx.Dst.Bounds()

	for y := srcBounds.Min.Y; y < srcBounds.Max.Y; y++ {
		for x := srcBounds.Min.X; x < srcBounds.Max.X; x++ {
			tx := x0 + (x - srcBounds.Min.X)
			ty := y0 + (y - srcBounds.Min.Y)

			// 越界检查
			if tx < dstBounds.Min.X || tx >= dstBounds.Max.X ||
				ty < dstBounds.Min.Y || ty >= dstBounds.Max.Y {
				continue
			}

			srcR, srcG, srcB, srcA := src.At(x, y).RGBA()
			dstR, dstG, dstB, dstA := ctx.Dst.At(tx, ty).RGBA()

			srcR8 := uint8(srcR >> 8)
			srcG8 := uint8(srcG >> 8)
			srcB8 := uint8(srcB >> 8)
			srcA8 := uint8(srcA >> 8)

			dstR8 := uint8(dstR >> 8)
			dstG8 := uint8(dstG >> 8)
			dstB8 := uint8(dstB >> 8)
			dstA8 := uint8(dstA >> 8)

			newR := uint8((int(srcR8) + int(dstR8)) / 2)
			newG := uint8((int(srcG8) + int(dstG8)) / 2)
			newB := uint8((int(srcB8) + int(dstB8)) / 2)
			newA := uint8((int(srcA8) + int(dstA8)) / 2)

			ctx.Dst.Set(tx, ty, color.RGBA{newR, newG, newB, newA})
		}
	}
	return ctx
}

// 将当前图层减法添加到画布上，也就是与当前画布做减法
// 可以用于检测图像中的变化或突出差异
//subtraction

// 将当前图层乘法添加到画布上，也就是与当前画布做乘法
// 可以用于掩膜操作，图像合成，亮度调整
//multiplication

// 将当前图层除法添加到画布上，也就是与当前画布做除法
// 可以用于光照归一化，比值分析，去雾处理
//division

// 将当前图层与画布进行逻辑运算 - 与（AND）
//AND

// 将当前图层与画布进行逻辑运算 - 或（OR）
//OR

// 将当前图层与画布进行逻辑运算 - 异或（XOR）
//XOR

// 将当前图层与画布进行逻辑运算 - 非（NOT）
//NOT

// SaveToFile 保存在本地文件
func (ctx *CanvasContext) SaveToFile(filePath string) error {
	if ctx.Err != nil {
		return ctx.Err
	}
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
		return err
	}
	return nil
}
