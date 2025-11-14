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
