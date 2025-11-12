package imgHelper

import (
	drawx "golang.org/x/image/draw"
	"image"
	"image/draw"
)

// opsScale 伸缩操作
type opsScale struct {
	TargetWidth  int
	TargetHeight int
}

func OpsScale(targetWidth, targetHeight int) func(ctx *CanvasContext) error {
	ops := &opsScale{
		TargetWidth:  targetWidth,
		TargetHeight: targetHeight,
	}
	return ops.scale
}

func OpsScaleNearestNeighbor(targetWidth, targetHeight int) func(ctx *CanvasContext) error {
	ops := &opsScale{
		TargetWidth:  targetWidth,
		TargetHeight: targetHeight,
	}
	return ops.scaleNearestNeighbor
}

func OpsScaleCatmullRom(targetWidth, targetHeight int) func(ctx *CanvasContext) error {
	ops := &opsScale{
		TargetWidth:  targetWidth,
		TargetHeight: targetHeight,
	}
	return ops.scaleCatmullRom
}

func (layer *opsScale) scale(ctx *CanvasContext) error {
	ctx.Dst = Scale(ctx.Dst, layer.TargetWidth, layer.TargetHeight).(*image.RGBA)
	return nil
}

func (layer *opsScale) scaleNearestNeighbor(ctx *CanvasContext) error {
	ctx.Dst = ScaleNearestNeighbor(ctx.Dst, layer.TargetWidth, layer.TargetHeight).(*image.RGBA)
	return nil
}

func (layer *opsScale) scaleCatmullRom(ctx *CanvasContext) error {
	ctx.Dst = ScaleCatmullRom(ctx.Dst, layer.TargetWidth, layer.TargetHeight).(*image.RGBA)
	return nil
}

// Scale 使用双线性插值算法将源图片拉伸或压缩到目标大小
func Scale(src image.Image, targetWidth, targetHeight int) image.Image {
	dst := image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))
	drawx.ApproxBiLinear.Scale(dst, dst.Bounds(), src, src.Bounds(), draw.Over, nil)
	return dst
}

// ScaleNearestNeighbor 最近邻插值，速度快，但可能会导致图像出现锯齿。
func ScaleNearestNeighbor(src image.Image, targetWidth, targetHeight int) image.Image {
	dst := image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))
	drawx.NearestNeighbor.Scale(dst, dst.Bounds(), src, src.Bounds(), draw.Over, nil)
	return dst
}

// ScaleCatmullRom 插值，质量较高，但速度较慢。
func ScaleCatmullRom(src image.Image, targetWidth, targetHeight int) image.Image {
	dst := image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))
	drawx.CatmullRom.Scale(dst, dst.Bounds(), src, src.Bounds(), draw.Over, nil)
	return dst
}
