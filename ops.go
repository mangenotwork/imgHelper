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

// GaussianBlur1D  一维高斯模糊
// sigma : 降噪程度
func GaussianBlur1D(src image.Image, sigma float64) image.Image {
	bounds := src.Bounds()
	result := image.NewRGBA(bounds)
	kernel := generateGaussianKernel(sigma)
	// 水平方向模糊
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			var rSum, gSum, bSum, aSum float64
			kernelSize := len(kernel)
			halfKernelSize := kernelSize / 2
			for i := -halfKernelSize; i <= halfKernelSize; i++ {
				newX := x + i
				if newX < bounds.Min.X {
					newX = bounds.Min.X
				} else if newX >= bounds.Max.X {
					newX = bounds.Max.X - 1
				}
				r, g, b, a := src.At(newX, y).RGBA()
				r = r / 256
				g = g / 256
				b = b / 256
				rSum += float64(r) * kernel[i+halfKernelSize]
				gSum += float64(g) * kernel[i+halfKernelSize]
				bSum += float64(b) * kernel[i+halfKernelSize]
				aSum += float64(a) * kernel[i+halfKernelSize]
			}
			result.Set(x, y, color.RGBA{R: uint8(rSum), G: uint8(gSum), B: uint8(bSum), A: uint8(aSum)})
		}
	}
	// 垂直方向模糊
	temp := image.NewRGBA(bounds)
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			var rSum, gSum, bSum, aSum float64
			kernelSize := len(kernel)
			halfKernelSize := kernelSize / 2
			for i := -halfKernelSize; i <= halfKernelSize; i++ {
				newY := y + i
				if newY < bounds.Min.Y {
					newY = bounds.Min.Y
				} else if newY >= bounds.Max.Y {
					newY = bounds.Max.Y - 1
				}
				r, g, b, a := result.At(x, newY).RGBA()
				r = r / 256
				g = g / 256
				b = b / 256
				rSum += float64(r) * kernel[i+halfKernelSize]
				gSum += float64(g) * kernel[i+halfKernelSize]
				bSum += float64(b) * kernel[i+halfKernelSize]
				aSum += float64(a) * kernel[i+halfKernelSize]
			}
			temp.Set(x, y, color.RGBA{R: uint8(rSum), G: uint8(gSum), B: uint8(bSum), A: uint8(aSum)})
		}
	}
	return temp
}

// OpsGaussianBlur1D  一维高斯模糊
func OpsGaussianBlur1D(sigma float64) func(ctx *CanvasContext) error {
	return func(ctx *CanvasContext) error {
		ctx.Dst = GaussianBlur1D(ctx.Dst, sigma).(*image.RGBA)
		return nil
	}
}

// Thinning 图像细化
// 针对文本图像进行细化处理
func Thinning(src image.Image) image.Image {
	bin := binaryImgForText(src)
	height, width := len(bin), len(bin[0])
	// 复制原图避免修改输入
	thinned := make([][]uint8, height)
	for y := range bin {
		thinned[y] = make([]uint8, width)
		copy(thinned[y], bin[y])
	}
	for {
		deleted := false // 标记本轮是否删除像素
		// 标记符合条件的像素
		toDelete1 := make([][]bool, height)
		for y := 0; y < height; y++ {
			toDelete1[y] = make([]bool, width)
			for x := 0; x < width; x++ {
				if thinned[y][x] != 255 {
					continue // 跳过背景
				}
				n := countForeground(thinned, x, y)
				c := countConnections(thinned, x, y)
				// 端点保护：N(p1)=1时不删除
				if n == 1 {
					continue
				}
				// 原条件：2 ≤ N(p1) ≤ 6 且 C(p1)=1
				if n < 2 || n > 6 || c != 1 {
					continue
				}
				// 条件4-5：p2*p4*p6=0 且 p4*p6*p8=0
				p2 := getNeighbor(thinned, x, y, 0)
				p4 := getNeighbor(thinned, x, y, 2)
				p6 := getNeighbor(thinned, x, y, 4)
				p8 := getNeighbor(thinned, x, y, 6)
				if p2*p4*p6 == 0 && p4*p6*p8 == 0 {
					toDelete1[y][x] = true
					deleted = true
				}
			}
		}
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				if toDelete1[y][x] {
					thinned[y][x] = 0
				}
			}
		}

		// 标记符合条件的像素
		toDelete2 := make([][]bool, height)
		for y := 0; y < height; y++ {
			toDelete2[y] = make([]bool, width)
			for x := 0; x < width; x++ {
				if thinned[y][x] != 255 {
					continue // 跳过背景
				}
				n := countForeground(thinned, x, y)
				c := countConnections(thinned, x, y)
				// 【新增】端点保护：N(p1)=1时不删除
				if n == 1 {
					continue
				}
				// 原条件：2 ≤ N(p1) ≤ 6 且 C(p1)=1
				if n < 2 || n > 6 || c != 1 {
					continue
				}
				// 条件4-5：p2*p4*p8=0 且 p2*p6*p8=0
				p2 := getNeighbor(thinned, x, y, 0)
				p4 := getNeighbor(thinned, x, y, 2)
				p6 := getNeighbor(thinned, x, y, 4)
				p8 := getNeighbor(thinned, x, y, 6)
				if p2*p4*p8 == 0 && p2*p6*p8 == 0 {
					toDelete2[y][x] = true
					deleted = true
				}
			}
		}
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				if toDelete2[y][x] {
					thinned[y][x] = 0
				}
			}
		}
		if !deleted {
			break
		}
	}
	height, width = len(thinned), len(thinned[0])
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			val := thinned[y][x]
			img.SetRGBA(x, y, color.RGBA{R: val, G: val, B: val, A: 255}) // 灰度图（骨架为白）
		}
	}
	return img
}

// OpsThinning 图像细化
func OpsThinning() func(ctx *CanvasContext) error {
	return func(ctx *CanvasContext) error {
		ctx.Dst = Thinning(ctx.Dst).(*image.RGBA)
		return nil
	}
}

// SmoothProcessing 彩色图像的平滑处理
// kernelSize : 平滑处理的核大小
func SmoothProcessing(src image.Image, kernelSize int) image.Image {
	bounds := src.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	dst := image.NewRGBA(bounds)
	halfKernel := kernelSize / 2
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			var rSum, gSum, bSum, count int
			for ky := -halfKernel; ky <= halfKernel; ky++ {
				for kx := -halfKernel; kx <= halfKernel; kx++ {
					nx := x + kx
					ny := y + ky
					if nx >= 0 && nx < width && ny >= 0 && ny < height {
						r, g, b, _ := src.At(nx, ny).RGBA()
						rSum += int(r / 256)
						gSum += int(g / 256)
						bSum += int(b / 256)
						count++
					}
				}
			}
			if count > 0 {
				r := uint8(rSum / count)
				g := uint8(gSum / count)
				b := uint8(bSum / count)
				dst.Set(x, y, color.RGBA{R: r, G: g, B: b, A: 255})
			}
		}
	}
	return dst
}

// OpsSmoothProcessing 彩色图像的平滑处理
func OpsSmoothProcessing(kernelSize int) func(ctx *CanvasContext) error {
	return func(ctx *CanvasContext) error {
		ctx.Dst = SmoothProcessing(ctx.Dst, kernelSize).(*image.RGBA)
		return nil
	}
}
