package imgHelper

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
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
	canvasContext.Dst = dst
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

func (ctx *CanvasContext) GetErr() error {
	return ctx.Err
}

// Ext 执行传入绘制的方法(操作ops)并接收绘制产生的错误
// 只要是实现了  fn func(ctx *CanvasContext) error 方法就可以调用此方法
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

			ctx.Dst.Set(tx, ty, color.RGBA{R: newR, G: newG, B: newB, A: newA})
		}
	}
	return ctx
}

// Subtraction 将当前图层减法添加到画布上，也就是与当前画布做减法
// 可以用于检测图像中的变化或突出差异
// 可选参数soft true:柔和减法效果
func (ctx *CanvasContext) Subtraction(layer Layer, soft ...bool) *CanvasContext {
	src := layer.GetResource()
	srcBounds := src.Bounds()
	x0, y0, _, _ := layer.GetXY()
	dstBounds := ctx.Dst.Bounds()

	softFlag := false
	if len(soft) > 0 {
		softFlag = soft[0]
	}

	for y := srcBounds.Min.Y; y < srcBounds.Max.Y; y++ {
		for x := srcBounds.Min.X; x < srcBounds.Max.X; x++ {

			tx := x0 + (x - srcBounds.Min.X)
			ty := y0 + (y - srcBounds.Min.Y)
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

			var newR, newG, newB, newA uint8

			if softFlag {
				// 缩放差值：用 (dst - src) 的一半，不足则取0（保留一定亮度）
				newR = max(0, (dstR8-srcR8)/2+128) // 加128偏移，避免过暗
				newG = max(0, (dstG8-srcG8)/2+128)
				newB = max(0, (dstB8-srcB8)/2+128)
				newA = max(0, (dstA8-srcA8)/2+128)

			} else {
				// 取绝对值：保留差异（无论谁减谁）
				newR = abs(dstR8 - srcR8)
				newG = abs(dstG8 - srcG8)
				newB = abs(dstB8 - srcB8)
				newA = abs(dstA8 - srcA8)
				if newA == 0 { // 保留一定亮度
					newA = 128
				}
			}

			ctx.Dst.Set(tx, ty, color.RGBA{R: newR, G: newG, B: newB, A: newA})
		}
	}
	return ctx
}

// Multiplication 将当前图层乘法添加到画布上，也就是与当前画布做乘法
// 可以用于掩膜操作，图像合成，亮度调整
// 可选参数：
//   - normalize: 是否归一化（默认true）
//   - scale: 缩放因子（默认255，值越小保留的小值越多，如128会增强低亮度
func (ctx *CanvasContext) Multiplication(layer Layer, normalize ...any) *CanvasContext {
	src := layer.GetResource()
	srcBounds := src.Bounds()
	x0, y0, _, _ := layer.GetXY()
	dstBounds := ctx.Dst.Bounds()

	normalizeFlag := true
	if len(normalize) > 0 {
		if flag, ok := normalize[0].(bool); ok {
			normalizeFlag = flag
		}
	}

	scale := 255
	if len(normalize) > 1 {
		if s, ok := normalize[1].(int); ok && s > 0 {
			scale = s
		}
	}

	for y := srcBounds.Min.Y; y < srcBounds.Max.Y; y++ {
		for x := srcBounds.Min.X; x < srcBounds.Max.X; x++ {
			tx := x0 + (x - srcBounds.Min.X)
			ty := y0 + (y - srcBounds.Min.Y)
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

			var newR, newG, newB, newA uint8
			if normalizeFlag {
				// 四舍五入（+scale/2 再整除），减少小值被截断为0
				newR = uint8((int(dstR8)*int(srcR8) + scale/2) / scale)
				newG = uint8((int(dstG8)*int(srcG8) + scale/2) / scale)
				newB = uint8((int(dstB8)*int(srcB8) + scale/2) / scale)
				newA = uint8((int(dstA8)*int(srcA8) + scale/2) / scale)

				// 极暗区域（接近0）强制保留一点亮度，避免完全消失
				if newR < 5 {
					newR = min(5, newR+1) // 最低亮度1-5
				}
				if newG < 5 {
					newG = min(5, newG+1)
				}
				if newB < 5 {
					newB = min(5, newB+1)
				}
				if newA < 5 {
					newA = min(5, newA+1)
				}
			} else {
				// 截断模式：保留原逻辑，但也做低亮度保护
				newR = min(255, uint8(int(dstR8)*int(srcR8)))
				newG = min(255, uint8(int(dstG8)*int(srcG8)))
				newB = min(255, uint8(int(dstB8)*int(srcB8)))
				newA = min(255, uint8(int(dstA8)*int(srcA8)))
			}

			ctx.Dst.Set(tx, ty, color.RGBA{R: newR, G: newG, B: newB, A: newA})
		}
	}
	return ctx
}

// Division 将当前图层除法添加到画布上，也就是与当前画布做除法
// 可以用于光照归一化，比值分析，去雾处理
// 可选参数：
//   - normalize: 是否归一化（默认true，结果映射到 0~255 范围）
//   - scale: 缩放因子（默认255，影响归一化强度）
//   - zeroVal: 图层像素为0时的替代结果（默认255，避免除零错误）
func (ctx *CanvasContext) Division(layer Layer, opts ...any) *CanvasContext {
	src := layer.GetResource()
	srcBounds := src.Bounds()
	x0, y0, _, _ := layer.GetXY()
	dstBounds := ctx.Dst.Bounds()

	normalizeFlag := true
	if len(opts) > 0 {
		if flag, ok := opts[0].(bool); ok {
			normalizeFlag = flag
		}
	}

	scale := 255
	if len(opts) > 1 {
		if s, ok := opts[1].(int); ok && s > 0 {
			scale = s
		}
	}

	zeroVal := 255
	if len(opts) > 2 {
		if z, ok := opts[2].(int); ok {
			zeroVal = clamp(z, 0, 255)
		}
	}

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

			newR := ctx.calcDiv(dstR8, srcR8, normalizeFlag, scale, zeroVal)
			newG := ctx.calcDiv(dstG8, srcG8, normalizeFlag, scale, zeroVal)
			newB := ctx.calcDiv(dstB8, srcB8, normalizeFlag, scale, zeroVal)
			newA := ctx.calcDiv(dstA8, srcA8, normalizeFlag, scale, zeroVal)

			ctx.Dst.Set(tx, ty, color.RGBA{R: newR, G: newG, B: newB, A: newA})
		}
	}
	return ctx
}

func (ctx *CanvasContext) calcDiv(dst, src uint8, normalize bool, scale, zeroVal int) uint8 {
	// 处理除零：图层像素为0时返回替代值
	if src == 0 {
		return uint8(zeroVal)
	}

	dstInt := int(dst)
	srcInt := int(src)

	var res int
	if normalize {
		// 归一化模式：(dst * scale) / src（四舍五入），映射到0-255
		res = (dstInt*scale + srcInt/2) / srcInt // +srcInt/2实现四舍五入
	} else {
		// 非归一化模式：直接除法（dst / src），结果可能很小
		res = dstInt / srcInt
	}

	return uint8(clamp(res, 0, 255))
}

// AND 将当前图层与画布进行逻辑运算 - 与（AND）
func (ctx *CanvasContext) AND(layer Layer) *CanvasContext {
	src := layer.GetResource()
	srcBounds := src.Bounds()
	x0, y0, _, _ := layer.GetXY()
	dstBounds := ctx.Dst.Bounds()

	for y := srcBounds.Min.Y; y < srcBounds.Max.Y; y++ {
		for x := srcBounds.Min.X; x < srcBounds.Max.X; x++ {

			tx := x0 + (x - srcBounds.Min.X)
			ty := y0 + (y - srcBounds.Min.Y)
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

			newR := srcR8 & dstR8
			newG := srcG8 & dstG8
			newB := srcB8 & dstB8
			newA := srcA8 & dstA8

			ctx.Dst.Set(tx, ty, color.RGBA{R: newR, G: newG, B: newB, A: newA})
		}
	}
	return ctx
}

// OR 将当前图层与画布进行逻辑运算 - 或（OR）
func (ctx *CanvasContext) OR(layer Layer) *CanvasContext {
	src := layer.GetResource()
	srcBounds := src.Bounds()
	x0, y0, _, _ := layer.GetXY()
	dstBounds := ctx.Dst.Bounds()

	for y := srcBounds.Min.Y; y < srcBounds.Max.Y; y++ {
		for x := srcBounds.Min.X; x < srcBounds.Max.X; x++ {

			tx := x0 + (x - srcBounds.Min.X)
			ty := y0 + (y - srcBounds.Min.Y)
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

			// 按位或运算：对应位有一个为1则结果为1
			newR := srcR8 | dstR8
			newG := srcG8 | dstG8
			newB := srcB8 | dstB8
			newA := srcA8 | dstA8

			ctx.Dst.Set(tx, ty, color.RGBA{R: newR, G: newG, B: newB, A: newA})
		}
	}
	return ctx
}

// XOR 将当前图层与画布进行逻辑运算 - 异或（XOR）
func (ctx *CanvasContext) XOR(layer Layer) *CanvasContext {
	src := layer.GetResource()
	srcBounds := src.Bounds()
	x0, y0, _, _ := layer.GetXY()
	dstBounds := ctx.Dst.Bounds()

	for y := srcBounds.Min.Y; y < srcBounds.Max.Y; y++ {
		for x := srcBounds.Min.X; x < srcBounds.Max.X; x++ {

			tx := x0 + (x - srcBounds.Min.X)
			ty := y0 + (y - srcBounds.Min.Y)
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

			// 按位异或运算：对应位不同则为1，相同则为0
			newR := srcR8 ^ dstR8
			newG := srcG8 ^ dstG8
			newB := srcB8 ^ dstB8
			newA := srcA8 ^ dstA8

			ctx.Dst.Set(tx, ty, color.RGBA{R: newR, G: newG, B: newB, A: newA})
		}
	}
	return ctx
}

// NOT 将当前图层与画布进行逻辑运算 - 非（NOT）
func (ctx *CanvasContext) NOT(layer Layer) *CanvasContext {
	src := layer.GetResource()
	srcBounds := src.Bounds()
	x0, y0, _, _ := layer.GetXY()
	dstBounds := ctx.Dst.Bounds()

	for y := srcBounds.Min.Y; y < srcBounds.Max.Y; y++ {
		for x := srcBounds.Min.X; x < srcBounds.Max.X; x++ {
			tx := x0 + (x - srcBounds.Min.X)
			ty := y0 + (y - srcBounds.Min.Y)

			if tx < dstBounds.Min.X || tx >= dstBounds.Max.X ||
				ty < dstBounds.Min.Y || ty >= dstBounds.Max.Y {
				continue
			}

			srcR, srcG, srcB, srcA := src.At(x, y).RGBA()
			srcR8 := uint8(srcR >> 8)
			srcG8 := uint8(srcG >> 8)
			srcB8 := uint8(srcB >> 8)
			srcA8 := uint8(srcA >> 8)

			// 对每个通道执行按位非运算（0变1，1变0）
			newR := ^srcR8 // 等价于 0xff ^ srcR8 或 255 - srcR8
			newG := ^srcG8
			newB := ^srcB8
			newA := ^srcA8

			ctx.Dst.Set(tx, ty, color.RGBA{R: newR, G: newG, B: newB, A: newA})
		}
	}
	return ctx
}

// SaveToFile 保存在本地文件
func (ctx *CanvasContext) SaveToFile(filePath string) error {
	if ctx.Err != nil {
		return ctx.Err
	}
	outputFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func() {
		_ = outputFile.Close()
	}()

	// todo 判断图片类型，根据类型进行存储
	err = png.Encode(outputFile, ctx.Dst)
	if err != nil {
		return err
	}
	return nil
}

// Print 在终端打印当前画布每个像素点的颜色值
func (ctx *CanvasContext) Print() {
	for y := ctx.Dst.Bounds().Min.Y; y < ctx.Dst.Bounds().Max.Y; y++ {
		for x := ctx.Dst.Bounds().Min.X; x < ctx.Dst.Bounds().Max.X; x++ {
			colorVal := ctx.Dst.At(x, y)
			r, g, b, a := colorVal.RGBA()
			fmt.Printf("Pixel at (%d, %d): R=%d, G=%d, B=%d, A=%d", x, y, r>>8, g>>8, b>>8, a>>8)
		}
	}
}
