package imgHelper

import (
	"image"
)

// Layer 图层
type Layer interface {
	Draw(ctx *CanvasContext) error // 绘制实现
	Save(filePath string) error
	GetResource() image.Image
	GetXY() (int, int, int, int) // 依次是 x0,y0,x1,y1
}

type RangeType string

const (
	RangeRectType     RangeType = "rect"
	RangeCircleType   RangeType = "circle"
	RangeTriangleType RangeType = "triangle"
	RangePolygonType  RangeType = "polygon"
)

type RangeValue interface {
	Type() RangeType
}

// Range 矩形范围 用于图层绘制在画布的指定位置
type Range struct {
	X0 int
	Y0 int
	X1 int
	Y1 int
}

func (Range) Type() RangeType {
	return RangeRectType
}

// RangeCircle 圆形范围 用于图层绘制在画布的指定位置
type RangeCircle struct {
	Cx int // 圆心在源图像中的坐标
	Cy int // 圆心在源图像中的坐标
	R  int // 圆的半径
}

func (RangeCircle) Type() RangeType {
	return RangeCircleType
}

// RangeTriangle 三角形范围
type RangeTriangle struct {
	X0, Y0 int
	X1, Y1 int
	X2, Y2 int
}

func (RangeTriangle) Type() RangeType {
	return RangeTriangleType
}

type Point struct {
	X, Y int
}

// RangePolygon 多边形范围
type RangePolygon struct {
	Points []Point
}

func (RangePolygon) Type() RangeType {
	return RangePolygonType
}
