package imgHelper

import (
	drawx "golang.org/x/image/draw"
	"image"
	"image/color"
	"image/draw"
)

// Scale 使用双线性插值算法将源图片拉伸或压缩到目标大小
func Scale(src image.Image, targetWidth, targetHeight int) image.Image {
	dst := image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))
	drawx.ApproxBiLinear.Scale(dst, dst.Bounds(), src, src.Bounds(), draw.Over, nil)
	return dst
}

// 双线性插值计算颜色
func interpolateColor(c1, c2 color.Color, t float64) (uint8, uint8, uint8, uint8) {
	r1, g1, b1, a1 := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()
	r := uint8((1-t)*float64(r1>>8) + t*float64(r2>>8))
	g := uint8((1-t)*float64(g1>>8) + t*float64(g2>>8))
	b := uint8((1-t)*float64(b1>>8) + t*float64(b2>>8))
	a := uint8((1-t)*float64(a1>>8) + t*float64(a2>>8))
	return r, g, b, a
}

// todo 缩放
// draw.NearestNeighbor：最近邻插值，速度快，但可能会导致图像出现锯齿。
// draw.ApproxBiLinear：近似双线性插值，速度比 draw.BiLinear 快，但质量稍低。
// draw.CatmullRom：Catmull-Rom 插值，质量较高，但速度较慢。

// todo  Gray 灰度处理， 遍历图像的每个像素点并进行灰度化处理

// todo Brightness 图像点的亮度调整

//
