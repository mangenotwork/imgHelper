package imgHelper

import (
	"image"
	"image/color"
	"math"
)

// RotateLayer 旋转图层
type RotateLayer struct {
	Angle float64
}

func NewRotateLayer(angle float64) *RotateLayer {
	return &RotateLayer{
		Angle: angle,
	}
}

// Rotate90Layer 旋转90度的图层
func Rotate90Layer() func(ctx *CanvasContext) error {
	return NewRotateLayer(90).Draw
}

func Rotate180Layer() func(ctx *CanvasContext) error {
	return NewRotateLayer(180).Draw
}

func Rotate270Layer() func(ctx *CanvasContext) error {
	return NewRotateLayer(270).Draw
}

func (layer *RotateLayer) Draw(ctx *CanvasContext) error {
	ctx.Dst = Rotate(ctx.Dst, layer.Angle).(*image.RGBA)
	return nil
}

// Rotate 图像旋转
func Rotate(src image.Image, angle float64) image.Image {
	srcBounds := src.Bounds()
	srcWidth := srcBounds.Dx()
	srcHeight := srcBounds.Dy()

	rad := angle * math.Pi / 180
	cos := math.Cos(rad)
	sin := math.Sin(rad)

	x1 := math.Abs(float64(srcWidth)*cos) + math.Abs(float64(srcHeight)*sin)
	y1 := math.Abs(float64(srcWidth)*sin) + math.Abs(float64(srcHeight)*cos)

	dstWidth := int(math.Ceil(x1))
	dstHeight := int(math.Ceil(y1))
	dst := image.NewRGBA(image.Rect(0, 0, dstWidth, dstHeight))

	// 计算旋转中心
	srcCenterX := float64(srcWidth) / 2
	srcCenterY := float64(srcHeight) / 2
	dstCenterX := float64(dstWidth) / 2
	dstCenterY := float64(dstHeight) / 2

	// 遍历目标图像的每个像素
	for y := 0; y < dstHeight; y++ {
		for x := 0; x < dstWidth; x++ {
			// 计算目标像素相对于旋转中心的坐标
			dx := float64(x) - dstCenterX
			dy := float64(y) - dstCenterY

			// 逆向旋转得到源图像中的坐标
			srcX := cos*dx + sin*dy + srcCenterX
			srcY := -sin*dx + cos*dy + srcCenterY

			// 检查源坐标是否在源图像范围内
			if srcX >= 0 && srcX < float64(srcWidth) && srcY >= 0 && srcY < float64(srcHeight) {
				// 双线性插值
				x0 := int(math.Floor(srcX))
				y0 := int(math.Floor(srcY))
				x1 := x0 + 1
				y1 := y0 + 1

				if x1 >= srcWidth {
					x1 = srcWidth - 1
				}
				if y1 >= srcHeight {
					y1 = srcHeight - 1
				}

				srcColor00 := src.At(x0, y0)
				srcColor01 := src.At(x0, y1)
				srcColor10 := src.At(x1, y0)
				srcColor11 := src.At(x1, y1)

				// 计算插值权重
				u := srcX - float64(x0)
				v := srcY - float64(y0)

				// 双线性插值计算颜色
				r0, g0, b0, a0 := interpolateColor(srcColor00, srcColor10, u)
				r1, g1, b1, a1 := interpolateColor(srcColor01, srcColor11, u)
				r, g, b, a := interpolateColor(color.RGBA{r0, g0, b0, a0}, color.RGBA{r1, g1, b1, a1}, v)

				dst.Set(x, y, color.RGBA{r, g, b, a})
			}
		}
	}
	return dst
}
