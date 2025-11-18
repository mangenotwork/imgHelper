package imgHelper

import (
	"image"
	"image/color"
	"image/draw"
	"math"
)

// Gray 灰度处理
func Gray(src image.Image) image.Image {
	bounds := src.Bounds()
	grayImg := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := src.At(x, y).RGBA()
			// 计算灰度值
			gray := uint8(math.Round(float64(r>>8)*0.299 + float64(g>>8)*0.587 + float64(b>>8)*0.114))
			grayImg.Set(x, y, color.RGBA{R: gray, G: gray, B: gray, A: uint8(a >> 8)})
		}
	}
	return grayImg
}

// OpsGray 灰度处理操作
func OpsGray() func(ctx *CanvasContext) error {
	return func(ctx *CanvasContext) error {
		ctx.Dst = Gray(ctx.Dst).(*image.RGBA)
		return nil
	}
}

// BinaryImg 二值图
// 参数: thresholdVal阈值，通过这个阈值来划分二值，默认为128
func BinaryImg(src image.Image, thresholdVal ...int) image.Image {
	threshold := 128
	if len(thresholdVal) > 0 {
		threshold = thresholdVal[0]
	}
	bounds := src.Bounds()
	binaryImg := image.NewGray(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			gray := color.GrayModel.Convert(src.At(x, y)).(color.Gray)
			if gray.Y > uint8(threshold) {
				// 大于阈值的像素设为白色
				binaryImg.SetGray(x, y, color.Gray{Y: 255})
			} else {
				// 小于阈值的像素设为黑色
				binaryImg.SetGray(x, y, color.Gray{Y: 0})
			}
		}
	}
	return binaryImg
}

// OpsBinaryImg 二值图操作
// 参数: thresholdVal阈值，通过这个阈值来划分二值，默认为128
func OpsBinaryImg(thresholdVal ...int) func(ctx *CanvasContext) error {
	return func(ctx *CanvasContext) error {
		binaryImg := BinaryImg(ctx.Dst, thresholdVal...)
		bounds := binaryImg.Bounds()
		rgbaImg := image.NewRGBA(bounds)
		draw.Draw(rgbaImg, bounds, binaryImg, bounds.Min, draw.Src)
		ctx.Dst = rgbaImg
		return nil
	}
}

// Brightness 图像点的亮度调整
func Brightness(src image.Image, brightnessVal int) image.Image {
	bounds := src.Bounds()
	newImg := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := src.At(x, y).RGBA()
			newR := int(r>>8) + brightnessVal
			newG := int(g>>8) + brightnessVal
			newB := int(b>>8) + brightnessVal
			if newR < 0 {
				newR = 0
			} else if newR > 255 {
				newR = 255
			}
			if newG < 0 {
				newG = 0
			} else if newG > 255 {
				newG = 255
			}
			if newB < 0 {
				newB = 0
			} else if newB > 255 {
				newB = 255
			}
			newImg.Set(x, y, color.RGBA{R: uint8(newR), G: uint8(newG), B: uint8(newB), A: uint8(a >> 8)})
		}
	}
	return newImg
}

// OpsBrightness 图像点的亮度调整操作
func OpsBrightness(brightnessVal int) func(ctx *CanvasContext) error {
	return func(ctx *CanvasContext) error {
		ctx.Dst = Brightness(ctx.Dst, brightnessVal).(*image.RGBA)
		return nil
	}
}

// Transposition 图像转置
func Transposition(src image.Image) image.Image {
	bounds := src.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	dst := image.NewRGBA(image.Rect(0, 0, height, width))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			dst.Set(y, x, src.At(x, y))
		}
	}
	return dst
}

// OpsTransposition 图像转置操作
func OpsTransposition() func(ctx *CanvasContext) error {
	return func(ctx *CanvasContext) error {
		ctx.Dst = Transposition(ctx.Dst).(*image.RGBA)
		return nil
	}
}

// MirrorHorizontal 图像水平镜像
func MirrorHorizontal(src image.Image) image.Image {
	bounds := src.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	dst := image.NewRGBA(bounds)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			dst.Set(width-1-x, y, src.At(x, y))
		}
	}
	return dst
}

// MirrorVertical 图像垂直镜像
func MirrorVertical(src image.Image) image.Image {
	bounds := src.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	dst := image.NewRGBA(bounds)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			dst.Set(x, height-1-y, src.At(x, y))
		}
	}
	return dst
}

// OpsMirrorHorizontal 图像水平镜像操作
func OpsMirrorHorizontal() func(ctx *CanvasContext) error {
	return func(ctx *CanvasContext) error {
		ctx.Dst = MirrorHorizontal(ctx.Dst).(*image.RGBA)
		return nil
	}
}

// OpsMirrorVertical 图像垂直镜像操作
func OpsMirrorVertical() func(ctx *CanvasContext) error {
	return func(ctx *CanvasContext) error {
		ctx.Dst = MirrorVertical(ctx.Dst).(*image.RGBA)
		return nil
	}
}

// Relief 图像浮雕
func Relief(src image.Image) image.Image {
	bounds := src.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	dst := image.NewRGBA(bounds)
	for y := 0; y < height-1; y++ {
		for x := 0; x < width-1; x++ {
			currentPixel := src.At(x, y)
			nextPixel := src.At(x+1, y+1)
			r1, g1, b1, _ := currentPixel.RGBA()
			r2, g2, b2, _ := nextPixel.RGBA()
			r := int(r1/256) - int(r2/256) + 128
			g := int(g1/256) - int(g2/256) + 128
			b := int(b1/256) - int(b2/256) + 128
			if r < 0 {
				r = 0
			} else if r > 255 {
				r = 255
			}
			if g < 0 {
				g = 0
			} else if g > 255 {
				g = 255
			}
			if b < 0 {
				b = 0
			} else if b > 255 {
				b = 255
			}
			dst.Set(x, y, color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255})
		}
	}
	return dst
}

// OpsRelief 浮雕操作
func OpsRelief() func(ctx *CanvasContext) error {
	return func(ctx *CanvasContext) error {
		ctx.Dst = Relief(ctx.Dst).(*image.RGBA)
		return nil
	}
}

// ColorReversal 图像颜色反转
func ColorReversal(src image.Image) image.Image {
	bounds := src.Bounds()
	dst := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := src.At(x, y).RGBA()
			r = r / 256
			g = g / 256
			b = b / 256
			a = a / 256
			// 反转颜色
			r = 255 - r
			g = 255 - g
			b = 255 - b
			dst.Set(x, y, color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)})
		}
	}
	return dst
}

// OpsColorReversal 图像颜色反转操作
func OpsColorReversal() func(ctx *CanvasContext) error {
	return func(ctx *CanvasContext) error {
		ctx.Dst = ColorReversal(ctx.Dst).(*image.RGBA)
		return nil
	}
}

// Corrosion 图像腐蚀
func Corrosion(src image.Image) image.Image {
	bounds := src.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	dst := image.NewRGBA(bounds)
	structuringElement := [3][3]bool{
		{true, true, true},
		{true, true, true},
		{true, true, true},
	}
	for y := 1; y < height-1; y++ {
		for x := 1; x < width-1; x++ {
			allForeground := true
			for ky := 0; ky < 3; ky++ {
				for kx := 0; kx < 3; kx++ {
					if structuringElement[ky][kx] {
						nx := x + kx - 1
						ny := y + ky - 1
						r, _, _, _ := src.At(nx, ny).RGBA()
						if r/256 < 128 {
							allForeground = false
							break
						}
					}
				}
				if !allForeground {
					break
				}
			}
			if allForeground {
				dst.Set(x, y, color.RGBA{R: 255, G: 255, B: 255, A: 255})
			} else {
				dst.Set(x, y, color.RGBA{R: 0, G: 0, B: 0, A: 255})
			}
		}
	}
	return dst
}

// OpsCorrosion 图像腐蚀操作
func OpsCorrosion() func(ctx *CanvasContext) error {
	return func(ctx *CanvasContext) error {
		ctx.Dst = Corrosion(ctx.Dst).(*image.RGBA)
		return nil
	}
}

// Dilation 图像膨胀
func Dilation(src image.Image) image.Image {
	src = BinaryImg(src) // 先二值化,提升膨胀效果
	bounds := src.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	dst := image.NewRGBA(bounds)
	structuringElement := [3][3]bool{
		{true, true, true},
		{true, true, true},
		{true, true, true},
	}
	for y := 1; y < height-1; y++ {
		for x := 1; x < width-1; x++ {
			anyForeground := false
			for ky := 0; ky < 3; ky++ {
				for kx := 0; kx < 3; kx++ {
					if structuringElement[ky][kx] {
						nx := x + kx - 1
						ny := y + ky - 1
						r, _, _, _ := src.At(nx, ny).RGBA()
						if r/256 >= 128 {
							anyForeground = true
							break
						}
					}
				}
				if anyForeground {
					break
				}
			}
			if anyForeground {
				dst.Set(x, y, color.RGBA{R: 255, G: 255, B: 255, A: 255})
			} else {
				dst.Set(x, y, color.RGBA{R: 0, G: 0, B: 0, A: 255})
			}
		}
	}
	return dst
}

// OpsDilation 图像膨胀操作
func OpsDilation() func(ctx *CanvasContext) error {
	return func(ctx *CanvasContext) error {
		ctx.Dst = Dilation(ctx.Dst).(*image.RGBA)
		return nil
	}
}

// Opening 图像的开运算
// 图像的开运算（Opening）是一种形态学操作，它是先对图像进行腐蚀操作，然后再进行膨胀操作。开运算可以去除图像中的小物体、分离物体以及平滑物体的边界
func Opening(src image.Image) image.Image {
	src = Corrosion(src)
	src = Dilation(src)
	return src
}

// OpsOpening 图像的开运算操作
func OpsOpening() func(ctx *CanvasContext) error {
	return func(ctx *CanvasContext) error {
		ctx.Dst = Opening(ctx.Dst).(*image.RGBA)
		return nil
	}
}

// Closing 图像的闭运算
// 图像的闭运算（Closing）是一种形态学操作，它是先对图像进行膨胀操作，然后再进行腐蚀操作。闭运算常用于填充物体内的小孔、连接邻近的物体等。
func Closing(src image.Image) image.Image {
	src = Dilation(src)
	src = Corrosion(src)
	return src
}

// OpsClosing 图像的闭运算操作
func OpsClosing() func(ctx *CanvasContext) error {
	return func(ctx *CanvasContext) error {
		ctx.Dst = Closing(ctx.Dst).(*image.RGBA)
		return nil
	}
}

// Hue 调整色相
// 通过将 RGB 颜色空间转换为 HSV（Hue, Saturation, Value）颜色空间，调整色相（Hue）值后再转换回 RGB 颜色空间
// hueAdjustment : 色相调整值
func Hue(src image.Image, hueAdjustment float64) image.Image {
	bounds := src.Bounds()
	dst := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := src.At(x, y).RGBA()
			r8 := uint8(r / 256)
			g8 := uint8(g / 256)
			b8 := uint8(b / 256)
			a8 := uint8(a / 256)
			h, s, v := RGBToHSV(r8, g8, b8)
			h = math.Mod(h+hueAdjustment, 360)
			if h < 0 {
				h += 360
			}
			rFloat, gFloat, bFloat := HSVToRGB(h, s, v)
			dst.Set(x, y, color.RGBA{
				R: rFloat,
				G: gFloat,
				B: bFloat,
				A: a8,
			})
		}
	}
	return dst
}

// OpsHue 调整色相操作
// hueAdjustment : 色相调整值
func OpsHue(hueAdjustment float64) func(ctx *CanvasContext) error {
	return func(ctx *CanvasContext) error {
		ctx.Dst = Hue(ctx.Dst, hueAdjustment).(*image.RGBA)
		return nil
	}
}
