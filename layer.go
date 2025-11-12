package imgHelper

import (
	"image"
)

// Layer 图层
type Layer interface {
	Draw(ctx *CanvasContext) error // 绘制实现
	Save(filePath string)
	GetResource() image.Image
	GetXY() (int, int, int, int) // 依次是 x0,y0,x1,y1
}

// Range 矩形范围 用于图层显示在画布的指定位置
type Range struct {
	X0 int
	Y0 int
	X1 int
	Y1 int
}
