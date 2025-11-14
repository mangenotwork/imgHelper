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
	// todo 曲边多边形（Curvilinear Polygon）
	// 核心特征
	//边界组成：由至少 3 个曲线段首尾相连形成的闭合回路，区别于普通多边形的直边
	//曲线类型：边可以是任意曲线，包括：
	//圆弧（如鲁洛三角形）
	//椭圆弧
	//贝塞尔曲线段
	//自由曲线（通过拟合离散点形成）
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

func (rg Range) Value() (int, int, int, int) {
	return rg.X0, rg.Y0, rg.X1, rg.Y1
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
