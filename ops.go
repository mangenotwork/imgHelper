package imgHelper

import (
	"image"
	"image/color"
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
