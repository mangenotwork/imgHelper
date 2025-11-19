package imgHelper

import (
	"image"
	"image/color"
	"math"
)

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

// Saturation 图像调整饱和度
// 通过将 RGB 颜色空间转换为 HSV 颜色空间，调整饱和度值后再转换回 RGB 颜色空间来完成
// saturationAdjustment: 调整饱和度的值
func Saturation(src image.Image, saturationAdjustment float64) image.Image {
	bounds := src.Bounds()
	dst := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := src.At(x, y).RGBA()
			r = r / 256
			g = g / 256
			b = b / 256
			h, s, v := RGBToHSV(uint8(r), uint8(g), uint8(b))
			s += saturationAdjustment
			if s < 0 {
				s = 0
			} else if s > 1 {
				s = 1
			}
			r1, g1, b1 := HSVToRGB(h, s, v)
			dst.Set(x, y, color.RGBA{R: r1, G: g1, B: b1, A: uint8(a)})
		}
	}
	return dst
}

// OpsSaturation 图像调整饱和度操作
func OpsSaturation(saturationAdjustment float64) func(ctx *CanvasContext) error {
	return func(ctx *CanvasContext) error {
		ctx.Dst = Saturation(ctx.Dst, saturationAdjustment).(*image.RGBA)
		return nil
	}
}

// AdjustColorBalance 调整色彩平衡
// 分别对图像中红、绿、蓝三个通道的值进行调整
func AdjustColorBalance(src image.Image, rAdjustment, gAdjustment, bAdjustment int) image.Image {
	bounds := src.Bounds()
	dst := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := src.At(x, y).RGBA()
			r = r / 256
			g = g / 256
			b = b / 256
			newR := int(r) + rAdjustment
			if newR < 0 {
				newR = 0
			} else if newR > 255 {
				newR = 255
			}
			newG := int(g) + gAdjustment
			if newG < 0 {
				newG = 0
			} else if newG > 255 {
				newG = 255
			}
			newB := int(b) + bAdjustment
			if newB < 0 {
				newB = 0
			} else if newB > 255 {
				newB = 255
			}
			dst.Set(x, y, color.RGBA{R: uint8(newR), G: uint8(newG), B: uint8(newB), A: uint8(a)})
		}
	}
	return dst
}

// OpsAdjustColorBalance 调整色彩平衡操作
func OpsAdjustColorBalance(rAdjustment, gAdjustment, bAdjustment int) func(ctx *CanvasContext) error {
	return func(ctx *CanvasContext) error {
		ctx.Dst = AdjustColorBalance(ctx.Dst, rAdjustment, gAdjustment, bAdjustment).(*image.RGBA)
		return nil
	}
}

// AdjustContrast 调整对比度
// contrast : 对比度调整值
func AdjustContrast(src image.Image, contrast float64) image.Image {
	bounds := src.Bounds()
	dst := image.NewRGBA(bounds)
	factor := (259 * (contrast + 255)) / (255 * (259 - contrast))
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := src.At(x, y).RGBA()
			r = r / 256
			g = g / 256
			b = b / 256
			newR := int(factor*(float64(r)-128) + 128)
			newG := int(factor*(float64(g)-128) + 128)
			newB := int(factor*(float64(b)-128) + 128)
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
			dst.Set(x, y, color.RGBA{R: uint8(newR), G: uint8(newG), B: uint8(newB), A: uint8(a)})
		}
	}
	return dst
}

// OpsAdjustContrast 调整对比度操作
func OpsAdjustContrast(contrast float64) func(ctx *CanvasContext) error {
	return func(ctx *CanvasContext) error {
		ctx.Dst = AdjustContrast(ctx.Dst, contrast).(*image.RGBA)
		return nil
	}
}

// AdjustSharpness 调整锐度
// 使用拉普拉斯算子进行图像锐化处理
// sharpness:锐度调整值
func AdjustSharpness(src image.Image, sharpness float64) image.Image {
	// 拉普拉斯算子
	var laplacianKernel = [][]int{
		{0, -1, 0},
		{-1, 4, -1},
		{0, -1, 0},
	}
	bounds := src.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	dst := image.NewRGBA(bounds)
	for y := 1; y < height-1; y++ {
		for x := 1; x < width-1; x++ {
			var rSum, gSum, bSum int
			for ky := -1; ky <= 1; ky++ {
				for kx := -1; kx <= 1; kx++ {
					r, g, b, _ := src.At(x+kx, y+ky).RGBA()
					r = r / 256
					g = g / 256
					b = b / 256
					kernelValue := laplacianKernel[ky+1][kx+1]
					rSum += int(r) * kernelValue
					gSum += int(g) * kernelValue
					bSum += int(b) * kernelValue
				}
			}
			r, g, b, a := src.At(x, y).RGBA()
			r = r / 256
			g = g / 256
			b = b / 256
			newR := int(r) + int(float64(rSum)*sharpness)
			newG := int(g) + int(float64(gSum)*sharpness)
			newB := int(b) + int(float64(bSum)*sharpness)
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
			dst.Set(x, y, color.RGBA{R: uint8(newR), G: uint8(newG), B: uint8(newB), A: uint8(a)})
		}
	}
	// 处理边缘像素复制原始像素值
	for y := 0; y < height; y++ {
		dst.Set(0, y, src.At(0, y))
		dst.Set(width-1, y, src.At(width-1, y))
	}
	for x := 0; x < width; x++ {
		dst.Set(x, 0, src.At(x, 0))
		dst.Set(x, height-1, src.At(x, height-1))
	}
	return dst
}

// OpsAdjustSharpness 调整锐度操作
func OpsAdjustSharpness(sharpness float64) func(ctx *CanvasContext) error {
	return func(ctx *CanvasContext) error {
		ctx.Dst = AdjustSharpness(ctx.Dst, sharpness).(*image.RGBA)
		return nil
	}
}

// AdjustColorScale 调整色阶
// blackPoint : 黑点
// whitePoint : 白点
// gamma : 伽马校正
func AdjustColorScale(src image.Image, blackPoint, whitePoint, gamma float64) image.Image {
	// 调整单个通道的色阶
	adjustChannel := func(value, blackPoint, whitePoint, gamma float64) float64 {
		// 将输入值限制在黑点和白点之间
		if value < blackPoint {
			value = 0
		} else if value > whitePoint {
			value = 255
		} else {
			// 线性映射到 0 - 255 范围
			value = (value - blackPoint) / (whitePoint - blackPoint) * 255
		}
		// 应用伽马校正
		if gamma != 1 {
			value = 255 * math.Pow(value/255, 1/gamma)
		}
		return value
	}

	bounds := src.Bounds()
	dst := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := src.At(x, y).RGBA()
			r = r / 256
			g = g / 256
			b = b / 256
			newR := adjustChannel(float64(r), blackPoint, whitePoint, gamma)
			newG := adjustChannel(float64(g), blackPoint, whitePoint, gamma)
			newB := adjustChannel(float64(b), blackPoint, whitePoint, gamma)
			dst.Set(x, y, color.RGBA{R: uint8(newR), G: uint8(newG), B: uint8(newB), A: uint8(a)})
		}
	}
	return dst
}

// OpsAdjustColorScale 调整色阶操作
// blackPoint : 黑点
// whitePoint : 白点
// gamma : 伽马校正
func OpsAdjustColorScale(blackPoint, whitePoint, gamma float64) func(ctx *CanvasContext) error {
	return func(ctx *CanvasContext) error {
		ctx.Dst = AdjustColorScale(ctx.Dst, blackPoint, whitePoint, gamma).(*image.RGBA)
		return nil
	}
}

// AdjustExposure 调整曝光度
func AdjustExposure(src image.Image, exposure float64) image.Image {
	bounds := src.Bounds()
	dst := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := src.At(x, y).RGBA()
			r = r / 256
			g = g / 256
			b = b / 256
			newR := int(math.Min(255, math.Max(0, float64(r)*math.Pow(2, exposure))))
			newG := int(math.Min(255, math.Max(0, float64(g)*math.Pow(2, exposure))))
			newB := int(math.Min(255, math.Max(0, float64(b)*math.Pow(2, exposure))))
			dst.Set(x, y, color.RGBA{R: uint8(newR), G: uint8(newG), B: uint8(newB), A: uint8(a)})
		}
	}
	return dst
}

// OpsAdjustExposure 调整曝光度操作
func OpsAdjustExposure(exposure float64) func(ctx *CanvasContext) error {
	return func(ctx *CanvasContext) error {
		ctx.Dst = AdjustExposure(ctx.Dst, exposure).(*image.RGBA)
		return nil
	}
}

// ColorTemperature 调整色温
// 通过调整图像中 RGB 颜色通道的比例来模拟不同的色温效果
// temperature:色温值
func ColorTemperature(src image.Image, temperature float64) image.Image {
	bounds := src.Bounds()
	dst := image.NewRGBA(bounds)
	rGain, gGain, bGain := calculateColorGains(temperature)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := src.At(x, y).RGBA()
			r = r / 256
			g = g / 256
			b = b / 256
			newR := int(math.Min(255, math.Max(0, float64(r)*rGain)))
			newG := int(math.Min(255, math.Max(0, float64(g)*gGain)))
			newB := int(math.Min(255, math.Max(0, float64(b)*bGain)))
			dst.Set(x, y, color.RGBA{R: uint8(newR), G: uint8(newG), B: uint8(newB), A: uint8(a)})
		}
	}
	return dst
}

// calculateColorGains 根据色温计算 RGB 增益
func calculateColorGains(temperature float64) (float64, float64, float64) {
	temperature = math.Max(1000, math.Min(40000, temperature)) / 100
	var r, g, b float64
	// 计算红色增益
	if temperature <= 66 {
		r = 255
	} else {
		r = temperature - 60
		r = 329.698727446 * math.Pow(r, -0.1332047592)
		if r < 0 {
			r = 0
		}
		if r > 255 {
			r = 255
		}
	}
	// 计算绿色增益
	if temperature <= 66 {
		g = temperature
		g = 99.4708025861*math.Log(g) - 161.1195681661
		if g < 0 {
			g = 0
		}
		if g > 255 {
			g = 255
		}
	} else {
		g = temperature - 60
		g = 288.1221695283 * math.Pow(g, -0.0755148492)
		if g < 0 {
			g = 0
		}
		if g > 255 {
			g = 255
		}
	}
	// 计算蓝色增益
	if temperature >= 66 {
		b = 255
	} else {
		if temperature <= 19 {
			b = 0
		} else {
			b = temperature - 10
			b = 138.5177312231*math.Log(b) - 305.0447927307
			if b < 0 {
				b = 0
			}
			if b > 255 {
				b = 255
			}
		}
	}
	// 归一化增益
	maxV := math.Max(r, math.Max(g, b))
	rGain := r / maxV
	gGain := g / maxV
	bGain := b / maxV
	return rGain, gGain, bGain
}

// OpsColorTemperature 调整色温
func OpsColorTemperature(temperature float64) func(ctx *CanvasContext) error {
	return func(ctx *CanvasContext) error {
		ctx.Dst = ColorTemperature(ctx.Dst, temperature).(*image.RGBA)
		return nil
	}
}

// ColorTone 调整色调
func ColorTone(src image.Image, adjustmentValue float64) image.Image {
	bounds := src.Bounds()
	dst := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := src.At(x, y).RGBA()
			r = r / 256
			g = g / 256
			b = b / 256
			h, s, v := RGBToHSV(uint8(r), uint8(g), uint8(b))
			h = math.Mod(h+adjustmentValue, 360)
			if h < 0 {
				h += 360
			}
			r1, g1, b1 := HSVToRGB(h, s, v)
			dst.Set(x, y, color.RGBA{R: r1, G: g1, B: b1, A: uint8(a)})
		}
	}
	return dst
}

// OpsColorTone 调整色调
func OpsColorTone(adjustmentValue float64) func(ctx *CanvasContext) error {
	return func(ctx *CanvasContext) error {
		ctx.Dst = ColorTone(ctx.Dst, adjustmentValue).(*image.RGBA)
		return nil
	}
}

// Denoise 图像降噪
// 通过调整高斯核的大小来降噪
// sigma:  ∑ (参数是降噪程度)
func Denoise(src image.Image, sigma float64) image.Image {
	return GaussianBlur1D(src, sigma)
}

// OpsDenoise 图像降噪
func OpsDenoise(sigma float64) func(ctx *CanvasContext) error {
	return func(ctx *CanvasContext) error {
		ctx.Dst = Denoise(ctx.Dst, sigma).(*image.RGBA)
		return nil
	}
}
