package imgHelper

import (
	drawx "golang.org/x/image/draw"
	"image"
	"image/draw"
)

// Scale 使用双线性插值算法将源图片拉伸或压缩到目标大小
func Scale(src image.Image, targetWidth, targetHeight int) image.Image {
	dst := image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))
	drawx.ApproxBiLinear.Scale(dst, dst.Bounds(), src, src.Bounds(), draw.Over, nil)
	return dst
}

// todo 缩放
// draw.NearestNeighbor：最近邻插值，速度快，但可能会导致图像出现锯齿。
// draw.ApproxBiLinear：近似双线性插值，速度比 draw.BiLinear 快，但质量稍低。
// draw.CatmullRom：Catmull-Rom 插值，质量较高，但速度较慢。

// todo  Gray 灰度处理， 遍历图像的每个像素点并进行灰度化处理

// todo Brightness 图像点的亮度调整

//
