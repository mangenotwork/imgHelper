package imgHelper

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"math"
	"os"
)

// 几何绘制图层
// 线，圆形，三角形，矩形，多边形，椭圆，扇形,星形,曲线
// todo : 曲边多边形

/*

todo  扩展2D绘制
参考

1. 直线家族
线段 (Line Segment) - 两点间的最短距离，最基本的 1 维图元
射线 (Ray) - 有起点无终点，向单一方向无限延伸
直线 (Infinite Line) - 无起点无终点，两端无限延伸
折线 (Polyline) - 多段线段首尾相连的非封闭图形
线环 (Line Loop) - 首尾相连的封闭折线
多线集合 (MultiLine) - 多组独立的折线或线段

2. 曲线家族
圆弧 (Circular Arc) - 圆的一部分，由圆心、半径和角度范围定义
椭圆弧 (Elliptical Arc) - 椭圆的一部分，SVG 中通过 rx, ry, x-axis-rotation 等参数精确描述
贝塞尔曲线 (Bézier Curve) - 计算机图形学核心曲线：
二次贝塞尔曲线 - 由起点、控制点、终点定义，形成抛物线
三次贝塞尔曲线 - 由起点、两个控制点、终点定义，可创建更复杂的平滑曲线
样条曲线 (Spline) - 通过多个插值点的平滑曲线，用于精确图形设计
螺旋线 (Spiral) - 如阿基米德螺旋、对数螺旋等特殊曲线

圆类图形 (2 维)
圆 (Circle) - 所有点到中心点距离相等的图形，由圆心和半径定义
椭圆 (Ellipse) - 圆的拉伸变形，由中心点、长半轴、短半轴定义
扇形 (Sector) - 由两条半径和一段圆弧组成的 "饼状" 图形
环形 (Annulus) - 两个同心圆之间的区域
月牙形 (Crescent) - 两个圆相交形成的新月状图形

多边形类图形 (2 维)

1. 基础多边形
三角形 (Triangle) - 最简单的封闭多边形，由三个顶点组成
三角形带 (Triangle Strip)、三角形扇 (Triangle Fan) - 优化绘制的三角形序列
四边形 (Quadrilateral) - 由四个顶点组成的多边形
矩形 (Rectangle) - 四个角为直角的四边形
正方形 (Square) - 四边相等的矩形
平行四边形 - 两组对边平行的四边形
菱形 - 四边相等的平行四边形，对角线垂直
梯形 - 只有一组对边平行的四边形
圆角矩形 (Rounded Rectangle) - 四角为圆弧的矩形
一般多边形 (Polygon) - 由 n (n≥3) 个顶点组成的封闭图形
凸多边形 - 所有内角小于 180°，任意两点连线在形内
凹多边形 - 至少有一个内角大于 180°
正多边形 - 所有边等长、所有角相等的多边形

2. 特殊多边形
星形 (Star Polygon) - 具有尖角放射状的多边形，如五角星、六角星等
带孔洞多边形 (Polygon with Holes) - 内部有一个或多个 "洞" 的多边形


复合形状 (Compound Shape) - 多个基本图形通过布尔运算组合：
并集 (Union) - 保留所有图形区域
交集 (Intersection) - 只保留重叠区域
差集 (Subtraction) - 从一个图形中减去另一个图形覆盖的部分
异或 (XOR) - 保留不重叠区域，排除重叠部分

路径 (Path) - 最强大的 2D 绘图元素，可包含：
直线段、曲线段 (贝塞尔曲线、圆弧等) 任意组合
多个子路径和孔洞
SVG 中通过 M (moveTo)、L (lineTo)、C (cubicTo)、Q (quadTo) 等指令定义


*/

type GeometryLayer struct {
	resource image.Image // 图层透明背景
	shapes   []Shape     // 图形集合：存储所有要绘制的几何图形
}

func NewGeometryLayer() *GeometryLayer {
	return &GeometryLayer{}
}

func (gLayer *GeometryLayer) Draw(ctx *CanvasContext) error {
	width := ctx.Dst.Bounds().Dx()
	height := ctx.Dst.Bounds().Dy()
	layerDst := image.NewRGBA(image.Rect(0, 0, width, height))
	transparent := color.RGBA{}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			layerDst.Set(x, y, transparent)
		}
	}
	gLayer.resource = layerDst
	for _, shape := range gLayer.shapes {
		gLayer.resource = shape.Render(gLayer.resource)
	}
	draw.Draw(
		ctx.Dst,
		image.Rect(0, 0, width, height),
		gLayer.resource,
		image.Point{},
		draw.Over,
	)
	return nil
}

func (gLayer *GeometryLayer) GetResource() image.Image {
	maxW, maxH := 0, 0
	for _, shape := range gLayer.shapes {
		w, y := shape.GetWH()
		maxW = max(maxW, w)
		maxH = max(maxH, y)
	}
	_ = gLayer.Draw(NewCanvas(maxW, maxH))
	return gLayer.resource
}

func (gLayer *GeometryLayer) Save(filePath string) error {
	maxW, maxH := 0, 0
	for _, shape := range gLayer.shapes {
		w, y := shape.GetWH()
		maxW = max(maxW, w)
		maxH = max(maxH, y)
	}

	err := gLayer.Draw(NewCanvas(maxW, maxH))
	if err != nil {
		return err
	}

	outputFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func() {
		_ = outputFile.Close()
	}()

	// todo 判断图片类型，根据类型进行存储
	return png.Encode(outputFile, gLayer.resource)
}

func (gLayer *GeometryLayer) GetXY() (int, int) {
	return 0, 0
}

func (gLayer *GeometryLayer) AddShape(s Shape) *GeometryLayer {
	gLayer.shapes = append(gLayer.shapes, s)
	return gLayer
}

type Shape interface {
	Render(src image.Image) image.Image
	GetWH() (int, int)
}

// Line 直线图形：包含自身的属性（起点、终点、颜色）
type Line struct {
	X0, Y0    int        // 起点
	X1, Y1    int        // 终点
	Color     color.RGBA // 自身颜色
	LineWidth int        // 直线粗度（像素数，最小为1）
}

func NewLine(x0, y0, x1, y1 int, c color.RGBA, lineWidth int) *Line {
	return &Line{X0: x0, Y0: y0, X1: x1, Y1: y1, Color: c, LineWidth: lineWidth}
}

func (l *Line) GetWH() (int, int) {
	return l.X1, l.Y1
}

func (l *Line) Render(src image.Image) image.Image {
	dst := image.NewRGBA(image.Rect(0, 0, src.Bounds().Dx(), src.Bounds().Dy()))
	bgColor := color.RGBA{R: 0, G: 0, B: 0, A: 0}
	draw.Draw(dst, dst.Bounds(), &image.Uniform{C: bgColor}, image.Point{}, draw.Over)

	p0x, p0y := float64(l.X0), float64(l.Y0)
	p1x, p1y := float64(l.X1), float64(l.Y1)

	// 计算直线的方向向量和长度
	dx := p1x - p0x
	dy := p1y - p0y
	lineLen := math.Hypot(dx, dy)
	if lineLen < 1e-6 { // 起点终点重合，无需绘制
		return dst
	}

	// 计算直线的单位方向向量
	dirX := dx / lineLen
	dirY := dy / lineLen

	// 计算线条的包围盒（仅遍历线条附近的像素，提升性能）
	halfWidth := float64(l.LineWidth) / 2.0
	minX := math.Min(p0x, p1x) - halfWidth - 1
	maxX := math.Max(p0x, p1x) + halfWidth + 1
	minY := math.Min(p0y, p1y) - halfWidth - 1
	maxY := math.Max(p0y, p1y) + halfWidth + 1

	// 转换为图像整数坐标范围
	bounds := src.Bounds()
	startX := int(math.Max(minX, float64(bounds.Min.X)))
	endX := int(math.Min(maxX, float64(bounds.Max.X)))
	startY := int(math.Max(minY, float64(bounds.Min.Y)))
	endY := int(math.Min(maxY, float64(bounds.Max.Y)))

	// 预计算平方值，减少浮点运算
	halfWidthMin05Sq := math.Pow(halfWidth-0.5, 2)
	halfWidthAdd05Sq := math.Pow(halfWidth+0.5, 2)

	// 遍历包围盒内的所有像素，逐像素计算透明度并绘制
	for y := startY; y <= endY; y++ {
		for x := startX; x <= endX; x++ {
			// 亚像素坐标 - 使用像素中心计算
			px := float64(x) + 0.5
			py := float64(y) + 0.5

			// 计算像素到直线的向量
			pixelToLineX := px - p0x
			pixelToLineY := py - p0y

			// 计算像素在直线方向上的投影长度（用于限制直线段范围）
			projLen := pixelToLineX*dirX + pixelToLineY*dirY
			projLen = math.Max(0, math.Min(projLen, lineLen)) // 限制在0到直线长度之间

			// 计算直线段上距离当前像素最近的点
			nearX := p0x + projLen*dirX
			nearY := p0y + projLen*dirY

			// 优化2：平方距离计算（避免开方，减少浮点误差，提升性能）
			dxNear := px - nearX
			dyNear := py - nearY
			finalDistSq := dxNear*dxNear + dyNear*dyNear
			finalDist := math.Sqrt(finalDistSq) // 仅在需要时开方

			// 计算像素的Alpha透明度（保留你原有的逻辑，仅优化浮点精度）
			var alpha uint8
			if finalDistSq < halfWidthMin05Sq { // 用平方比较，更精准
				alpha = 255
			} else if finalDistSq > halfWidthAdd05Sq {
				alpha = 0
			} else {
				// 浮点精度优化，避免微小误差导致的Alpha抖动
				alphaVal := (halfWidth + 0.5 - finalDist) * 255
				if alphaVal < 0 {
					alphaVal = 0
				} else if alphaVal > 255 {
					alphaVal = 255
				}
				alpha = uint8(alphaVal)
			}

			// 绘制像素（支持任意颜色的Alpha混合）
			if alpha > 0 {
				setPixel(dst, x, y, l.Color, alpha)
			}
		}
	}

	draw.Draw(dst, bounds, src, image.Point{}, draw.Over)

	return dst
}

// SolidCircle 实心圆图形：圆心、半径、颜色
type SolidCircle struct {
	Cx, Cy int        // 圆心坐标
	Radius int        // 圆半径（像素数，最小为1）
	Color  color.RGBA // 圆颜色
}

// NewSolidCircle 创建实心圆
func NewSolidCircle(cx, cy, radius int, c color.RGBA) *SolidCircle {
	if radius < 1 {
		radius = 1
	}
	return &SolidCircle{Cx: cx, Cy: cy, Radius: radius, Color: c}
}

func (s *SolidCircle) GetWH() (int, int) {
	return s.Cx + s.Radius, s.Cy + s.Radius
}

func (s *SolidCircle) Render(src image.Image) image.Image {
	dst := image.NewRGBA(image.Rect(0, 0, src.Bounds().Dx(), src.Bounds().Dy()))
	bgColor := color.RGBA{R: 0, G: 0, B: 0, A: 0}
	draw.Draw(dst, dst.Bounds(), &image.Uniform{C: bgColor}, image.Point{}, draw.Over)

	// 转换为浮点坐标，计算圆形包围盒
	cx, cy := float64(s.Cx), float64(s.Cy)
	radius := float64(s.Radius)
	// 包围盒：圆心±半径-1，避免超出图像边界
	minX := cx - radius - 1
	maxX := cx + radius + 1
	minY := cy - radius - 1
	maxY := cy + radius + 1
	bounds := src.Bounds()
	startX := int(math.Max(minX, float64(bounds.Min.X)))
	endX := int(math.Min(maxX, float64(bounds.Max.X)))
	startY := int(math.Max(minY, float64(bounds.Min.Y)))
	endY := int(math.Min(maxY, float64(bounds.Max.Y)))

	// 预计算平方值，减少浮点运算（抗锯齿过渡范围0.5像素）
	radiusAdd05Sq := math.Pow(radius+0.5, 2)
	radiusMin05Sq := math.Pow(radius-0.5, 2)

	for y := startY; y <= endY; y++ {
		for x := startX; x <= endX; x++ {
			// 亚像素坐标：使用像素中心计算，与直线抗锯齿逻辑一致
			px := float64(x) + 0.5
			py := float64(y) + 0.5

			// 计算像素到圆心的平方距离（避免开方，提升性能）
			dx := px - cx
			dy := py - cy
			distSq := dx*dx + dy*dy
			dist := math.Sqrt(distSq) // 仅在抗锯齿时开方

			// 计算Alpha透明度（抗锯齿逻辑与直线一致）
			var alpha uint8
			if distSq < radiusMin05Sq {
				alpha = 255 // 完全在圆内，不透明
			} else if distSq > radiusAdd05Sq {
				alpha = 0 // 完全在圆外，透明
			} else {
				// 0.5像素过渡范围，线性衰减Alpha
				alphaVal := (radius + 0.5 - dist) * 255
				if alphaVal < 0 {
					alphaVal = 0
				} else if alphaVal > 255 {
					alphaVal = 255
				}
				alpha = uint8(alphaVal)
			}
			if alpha > 0 {
				setPixel(dst, x, y, s.Color, alpha)
			}
		}
	}
	draw.Draw(dst, bounds, src, image.Point{}, draw.Over)
	return dst
}

// OutlineCircle 非实心圆（轮廓圆）图形：圆心、半径、颜色、线条宽度
type OutlineCircle struct {
	Cx, Cy    int        // 圆心坐标
	Radius    int        // 圆半径（像素数，最小为1）
	Color     color.RGBA // 轮廓颜色
	LineWidth int        // 轮廓线条宽度（像素数，最小为1）
}

// NewOutlineCircle 创建非实心圆
func NewOutlineCircle(cx, cy, radius, lineWidth int, c color.RGBA) *OutlineCircle {
	if radius < 1 {
		radius = 1
	}
	if lineWidth < 1 {
		lineWidth = 1
	}
	return &OutlineCircle{Cx: cx, Cy: cy, Radius: radius, Color: c, LineWidth: lineWidth}
}

func (o *OutlineCircle) GetWH() (int, int) {
	return o.Cx + o.Radius + o.LineWidth, o.Cy + o.Radius + o.LineWidth
}

func (o *OutlineCircle) Render(src image.Image) image.Image {
	dst := image.NewRGBA(image.Rect(0, 0, src.Bounds().Dx(), src.Bounds().Dy()))
	bgColor := color.RGBA{R: 0, G: 0, B: 0, A: 0}
	draw.Draw(dst, dst.Bounds(), &image.Uniform{C: bgColor}, image.Point{}, draw.Over)

	cx, cy := float64(o.Cx), float64(o.Cy)
	radius := float64(o.Radius)
	halfWidth := float64(o.LineWidth) / 2.0 // 线条半宽
	// 轮廓的内外半径
	outerRadius := radius + halfWidth
	innerRadius := radius - halfWidth
	// 包围盒：覆盖轮廓所有像素
	minX := cx - outerRadius - 1
	maxX := cx + outerRadius + 1
	minY := cy - outerRadius - 1
	maxY := cy + outerRadius + 1
	bounds := src.Bounds()
	startX := int(math.Max(minX, float64(bounds.Min.X)))
	endX := int(math.Min(maxX, float64(bounds.Max.X)))
	startY := int(math.Max(minY, float64(bounds.Min.Y)))
	endY := int(math.Min(maxY, float64(bounds.Max.Y)))

	innerEdge := innerRadius - 0.5   // 内轮廓抗锯齿左边界
	innerSmooth := innerRadius + 0.5 // 内轮廓抗锯齿右边界
	outerEdge := outerRadius - 0.5   // 外轮廓抗锯齿左边界
	outerSmooth := outerRadius + 0.5 // 外轮廓抗锯齿右边界

	// 遍历包围盒内的像素，逐像素绘制轮廓
	for y := startY; y <= endY; y++ {
		for x := startX; x <= endX; x++ {
			// 亚像素坐标：使用像素中心计算
			px := float64(x) + 0.5
			py := float64(y) + 0.5

			// 计算像素到圆心的实际距离
			dx := px - cx
			dy := py - cy
			dist := math.Sqrt(dx*dx + dy*dy)

			// 计算Alpha透明度（优化后的抗锯齿逻辑）
			var alpha uint8
			switch {
			// 完全在轮廓内部（内外抗锯齿区间之间）
			case dist >= innerSmooth && dist <= outerEdge:
				alpha = 255
			// 内轮廓抗锯齿区间（innerEdge ~ innerSmooth）
			case dist > innerEdge && dist < innerSmooth:
				alphaVal := (dist - innerEdge) / (innerSmooth - innerEdge) * 255
				alpha = uint8(math.Min(math.Max(alphaVal, 0), 255))
			// 外轮廓抗锯齿区间（outerEdge ~ outerSmooth）
			case dist > outerEdge && dist < outerSmooth:
				alphaVal := (outerSmooth - dist) / (outerSmooth - outerEdge) * 255
				alpha = uint8(math.Min(math.Max(alphaVal, 0), 255))
			// 超出轮廓范围
			default:
				alpha = 0
			}
			if alpha > 0 {
				setPixel(dst, x, y, o.Color, alpha)
			}
		}
	}
	draw.Draw(dst, bounds, src, image.Point{}, draw.Over)
	return dst
}

// SolidTriangle 实心三角形：三个顶点坐标、颜色
type SolidTriangle struct {
	X1, Y1 int        // 顶点1
	X2, Y2 int        // 顶点2
	X3, Y3 int        // 顶点3
	Color  color.RGBA // 填充颜色
}

// cross 叉乘计算：(p2-p1) × (p-p1)
// 结果符号：>0 点在向量左侧；<0 右侧；=0 共线
func cross(p1x, p1y, p2x, p2y, px, py float64) float64 {
	return (p2x-p1x)*(py-p1y) - (p2y-p1y)*(px-p1x)
}

// pointToSegmentDist 计算点到线段的最短距离
func pointToSegmentDist(px, py, x1, y1, x2, y2 float64) float64 {
	// 向量
	dx := x2 - x1
	dy := y2 - y1
	// 点到线段起点的向量
	tx := px - x1
	ty := py - y1
	// 投影长度
	proj := tx*dx + ty*dy
	if proj <= 0 {
		return math.Hypot(px-x1, py-y1) // 投影在起点外侧
	}
	// 线段长度平方
	lenSq := dx*dx + dy*dy
	if proj >= lenSq {
		return math.Hypot(px-x2, py-y2) // 投影在终点外侧
	}
	// 投影在线段上，计算垂距
	proj /= lenSq
	closestX := x1 + proj*dx
	closestY := y1 + proj*dy
	return math.Hypot(px-closestX, py-closestY)
}

// triangleAABB 计算三角形的轴对齐包围盒（AABB）
func triangleAABB(x1, y1, x2, y2, x3, y3 float64) (minX, minY, maxX, maxY float64) {
	minX = math.Min(math.Min(x1, x2), x3) - 1
	minY = math.Min(math.Min(y1, y2), y3) - 1
	maxX = math.Max(math.Max(x1, x2), x3) + 1
	maxY = math.Max(math.Max(y1, y2), y3) + 1
	return
}

// NewSolidTriangle 创建实心三角形（参数校验：确保三个顶点不共线，简单校验）
func NewSolidTriangle(x1, y1, x2, y2, x3, y3 int, c color.RGBA) *SolidTriangle {
	// 简单共线校验：叉乘为0则共线，强制调整第三个顶点
	cr := cross(float64(x1), float64(y1), float64(x2), float64(y2), float64(x3), float64(y3))
	if cr == 0 {
		x3 += 1
		y3 += 1
	}
	return &SolidTriangle{
		X1: x1, Y1: y1,
		X2: x2, Y2: y2,
		X3: x3, Y3: y3,
		Color: c,
	}
}

func (s *SolidTriangle) GetWH() (int, int) {
	return max(s.X1, s.X2, s.X3), max(s.Y1, s.Y2, s.Y3)
}

func (s *SolidTriangle) Render(src image.Image) image.Image {
	dst := image.NewRGBA(image.Rect(0, 0, src.Bounds().Dx(), src.Bounds().Dy()))
	bgColor := color.RGBA{R: 0, G: 0, B: 0, A: 0}
	draw.Draw(dst, dst.Bounds(), &image.Uniform{C: bgColor}, image.Point{}, draw.Over)

	x1, y1 := float64(s.X1), float64(s.Y1)
	x2, y2 := float64(s.X2), float64(s.Y2)
	x3, y3 := float64(s.X3), float64(s.Y3)
	minX, minY, maxX, maxY := triangleAABB(x1, y1, x2, y2, x3, y3)

	bounds := src.Bounds()
	startX := int(math.Max(minX, float64(bounds.Min.X)))
	endX := int(math.Min(maxX, float64(bounds.Max.X)))
	startY := int(math.Max(minY, float64(bounds.Min.Y)))
	endY := int(math.Min(maxY, float64(bounds.Max.Y)))

	for y := startY; y <= endY; y++ {
		for x := startX; x <= endX; x++ {
			// 亚像素坐标：像素中心
			px := float64(x) + 0.5
			py := float64(y) + 0.5
			// 判断点是否在三角形内部
			in := isPointInTriangle(px, py, x1, y1, x2, y2, x3, y3)
			if in {
				setPixel(dst, x, y, s.Color, 255) // 内部完全不透明
			} else {
				// 边缘抗锯齿：计算到三条边的最短距离，距离<1则混合Alpha
				d1 := pointToSegmentDist(px, py, x1, y1, x2, y2)
				d2 := pointToSegmentDist(px, py, x2, y2, x3, y3)
				d3 := pointToSegmentDist(px, py, x3, y3, x1, y1)
				minDist := math.Min(math.Min(d1, d2), d3)
				if minDist < 1.0 {
					alpha := uint8((1.0 - minDist) * 255) // 距离越近，Alpha越高
					setPixel(dst, x, y, s.Color, alpha)
				}
			}
		}
	}
	draw.Draw(dst, bounds, src, image.Point{}, draw.Over)
	return dst
}

// OutlineTriangle 非实心三角形（轮廓）：三个顶点、颜色、线条宽度
type OutlineTriangle struct {
	X1, Y1    int        // 顶点1
	X2, Y2    int        // 顶点2
	X3, Y3    int        // 顶点3
	Color     color.RGBA // 轮廓颜色
	LineWidth int        // 轮廓线条宽度（最小1）
}

// NewOutlineTriangle 创建非实心三角形
func NewOutlineTriangle(x1, y1, x2, y2, x3, y3 int, lineWidth int, c color.RGBA) *OutlineTriangle {
	cr := cross(float64(x1), float64(y1), float64(x2), float64(y2), float64(x3), float64(y3))
	if cr == 0 {
		x3 += 1
		y3 += 1
	}
	if lineWidth < 1 {
		lineWidth = 1
	}
	return &OutlineTriangle{
		X1: x1, Y1: y1,
		X2: x2, Y2: y2,
		X3: x3, Y3: y3,
		LineWidth: lineWidth,
		Color:     c,
	}
}

func (o *OutlineTriangle) GetWH() (int, int) {
	return max(o.X1, o.X2, o.X3) + o.LineWidth, max(o.Y1, o.Y2, o.Y3) + o.LineWidth
}

// isPointInOuterVertexFan 判断点是否在三角形**外部**的顶点尖角区域内（核心优化）
// vx,vy: 顶点坐标；v1x,v1y: 相邻顶点1；v2x,v2y: 相邻顶点2；lineWidth: 线宽；px,py: 待判断点
func isPointInOuterVertexFan(vx, vy, v1x, v1y, v2x, v2y, lineWidth, px, py float64) bool {
	// 三角形内部参考点（顶点对边的中点）
	midX := (v1x + v2x) / 2
	midY := (v1y + v2y) / 2

	// 获取两条边的外部单位法向量
	n1x, n1y := getEdgeNormal(vx, vy, v1x, v1y, midX, midY)
	n2x, n2y := getEdgeNormal(vx, vy, v2x, v2y, midX, midY)

	// 外扩后的顶点边界（法向量×线宽）
	p1x := vx + n1x*lineWidth
	p1y := vy + n1y*lineWidth
	p2x := vx + n2x*lineWidth
	p2y := vy + n2y*lineWidth

	// 判断点是否在 [v1->p1->vx->p2->v2] 的扇形区域内
	c1 := cross(vx, vy, p1x, p1y, px, py)
	c2 := cross(vx, vy, p2x, p2y, px, py)
	c3 := cross(vx, vy, v1x, v1y, px, py)
	c4 := cross(vx, vy, v2x, v2y, px, py)

	// 点到顶点的距离阈值（线宽的1.5倍，确保尖角范围足够）
	dist := math.Hypot(px-vx, py-vy)
	if dist > lineWidth*1.5 {
		return false
	}

	// 扇形区域判断：点在两条外扩法向量之间 或 点在原始边之间
	return (c1*c2 <= 0) || (c3*c4 <= 0)
}

// getEdgeNormal 获取线段的单位法向量（向三角形外部）
// v1x,v1y: 线段起点；v2x,v2y: 线段终点；triX,triY: 三角形内任意一点（用于判断法向量方向）
func getEdgeNormal(v1x, v1y, v2x, v2y, triX, triY float64) (nx, ny float64) {
	// 线段的方向向量
	dx := v2x - v1x
	dy := v2y - v1y
	// 两个法向量（垂直于方向向量）
	n1x, n1y := -dy, dx // 顺时针法向量
	n2x, n2y := dy, -dx // 逆时针法向量
	// 归一化
	n1x, n1y = normalize(n1x, n1y)
	n2x, n2y = normalize(n2x, n2y)
	// 判断哪个法向量指向三角形外部（叉乘符号）
	c1 := cross(v1x, v1y, v2x, v2y, v1x+n1x, v1y+n1y)
	c2 := cross(v1x, v1y, v2x, v2y, triX, triY)
	if c1*c2 < 0 {
		return n1x, n1y
	}
	return n2x, n2y
}

// triangleAABBLineWidth 计算三角形的轴对齐包围盒（扩大外扩范围）
func triangleAABBLineWidth(x1, y1, x2, y2, x3, y3, lineWidth float64) (minX, minY, maxX, maxY float64) {
	minX = math.Min(math.Min(x1, x2), x3) - lineWidth - 2
	minY = math.Min(math.Min(y1, y2), y3) - lineWidth - 2
	maxX = math.Max(math.Max(x1, x2), x3) + lineWidth + 2
	maxY = math.Max(math.Max(y1, y2), y3) + lineWidth + 2
	return
}

func (o *OutlineTriangle) Render(src image.Image) image.Image {
	dst := image.NewRGBA(image.Rect(0, 0, src.Bounds().Dx(), src.Bounds().Dy()))
	bgColor := color.RGBA{R: 0, G: 0, B: 0, A: 0}
	draw.Draw(dst, dst.Bounds(), &image.Uniform{C: bgColor}, image.Point{}, draw.Over)

	x1, y1 := float64(o.X1), float64(o.Y1)
	x2, y2 := float64(o.X2), float64(o.Y2)
	x3, y3 := float64(o.X3), float64(o.Y3)
	lineWidth := float64(o.LineWidth)
	halfWidth := lineWidth / 2.0
	smoothRange := 0.5 // 抗锯齿过渡范围
	minSmooth := halfWidth - smoothRange
	maxSmooth := halfWidth + smoothRange

	minX, minY, maxX, maxY := triangleAABBLineWidth(x1, y1, x2, y2, x3, y3, lineWidth)
	bounds := src.Bounds()
	startX := int(math.Max(minX, float64(bounds.Min.X)))
	endX := int(math.Min(maxX, float64(bounds.Max.X)))
	startY := int(math.Max(minY, float64(bounds.Min.Y)))
	endY := int(math.Min(maxY, float64(bounds.Max.Y)))

	for y := startY; y <= endY; y++ {
		for x := startX; x <= endX; x++ {
			// 亚像素坐标（像素中心）
			px := float64(x) + 0.5
			py := float64(y) + 0.5

			// 标记是否需要绘制该像素
			needDraw := false
			// 用于抗锯齿的距离值
			var distForSmooth = math.MaxFloat64

			d1 := pointToSegmentDist(px, py, x1, y1, x2, y2)
			d2 := pointToSegmentDist(px, py, x2, y2, x3, y3)
			d3 := pointToSegmentDist(px, py, x3, y3, x1, y1)
			minSegDist := math.Min(math.Min(d1, d2), d3)
			if minSegDist < maxSmooth {
				needDraw = true
				distForSmooth = minSegDist
			}

			// 顶点1的外尖角区域（相邻顶点2和3）
			if isPointInOuterVertexFan(x1, y1, x2, y2, x3, y3, lineWidth, px, py) {
				needDraw = true
				distForSmooth = math.Min(distForSmooth, math.Hypot(px-x1, py-y1))
			}
			// 顶点2的外尖角区域（相邻顶点1和3）
			if isPointInOuterVertexFan(x2, y2, x1, y1, x3, y3, lineWidth, px, py) {
				needDraw = true
				distForSmooth = math.Min(distForSmooth, math.Hypot(px-x2, py-y2))
			}
			// 顶点3的外尖角区域（相邻顶点1和2）
			if isPointInOuterVertexFan(x3, y3, x1, y1, x2, y2, lineWidth, px, py) {
				needDraw = true
				distForSmooth = math.Min(distForSmooth, math.Hypot(px-x3, py-y3))
			}

			if needDraw {
				var alpha uint8
				if distForSmooth < minSmooth {
					alpha = 255 // 完全在轮廓内，不透明
				} else {
					// 抗锯齿过渡：线性衰减
					alphaVal := (maxSmooth - distForSmooth) / (2 * smoothRange) * 255
					alpha = uint8(math.Min(math.Max(alphaVal, 0), 255))
				}
				setPixel(dst, x, y, o.Color, alpha)
			}
		}
	}
	draw.Draw(dst, bounds, src, image.Point{}, draw.Over)
	return dst
}

// SolidRect 实心矩形：左上角、右下角坐标 + 填充颜色
type SolidRect struct {
	X0, Y0 int        // 左上角坐标
	X1, Y1 int        // 右下角坐标
	Color  color.RGBA // 填充颜色
}

// NewSolidRect 创建实心矩形（自动修正坐标，确保X0<X1、Y0<Y1）
func NewSolidRect(x0, y0, x1, y1 int, c color.RGBA) *SolidRect {
	if x0 > x1 {
		x0, x1 = x1, x0
	}
	if y0 > y1 {
		y0, y1 = y1, y0
	}
	if x0 == x1 {
		x1++
	}
	if y0 == y1 {
		y1++
	}
	return &SolidRect{
		X0: x0, Y0: y0,
		X1: x1, Y1: y1,
		Color: c,
	}
}

func (s *SolidRect) GetWH() (int, int) {
	return s.X1, s.Y1
}

func (s *SolidRect) Render(src image.Image) image.Image {
	dst := image.NewRGBA(src.Bounds())
	draw.Draw(dst, dst.Bounds(), src, src.Bounds().Min, draw.Src)
	// 严格使用整数坐标遍历矩形区域（无亚像素、无抗锯齿）
	// 遍历范围：X从X0到X1-1，Y从Y0到Y1-1（像素的整数边界）
	for y := s.Y0; y < s.Y1; y++ {
		for x := s.X0; x < s.X1; x++ {
			// 直接绘制不透明像素，无任何过渡
			setPixel(dst, x, y, s.Color, 255)
		}
	}

	return dst
}

type OutlineRect struct {
	X0, Y0    int        // 左上角整数坐标
	X1, Y1    int        // 右下角整数坐标
	Color     color.RGBA // 轮廓颜色
	LineWidth int        // 线条宽度（整数，最小1）
}

// NewOutlineRect 创建非实心矩形
func NewOutlineRect(x0, y0, x1, y1 int, lineWidth int, c color.RGBA) *OutlineRect {
	if x0 > x1 {
		x0, x1 = x1, x0
	}
	if y0 > y1 {
		y0, y1 = y1, y0
	}
	if x0 == x1 {
		x1++
	}
	if y0 == y1 {
		y1++
	}
	if lineWidth < 1 {
		lineWidth = 1
	}
	return &OutlineRect{
		X0: x0, Y0: y0,
		X1: x1, Y1: y1,
		LineWidth: lineWidth,
		Color:     c,
	}
}

func (o *OutlineRect) GetWH() (int, int) {
	return o.X1 + o.LineWidth, o.Y1 + o.LineWidth
}

func (o *OutlineRect) Render(src image.Image) image.Image {
	dst := image.NewRGBA(src.Bounds())
	draw.Draw(dst, dst.Bounds(), src, src.Bounds().Min, draw.Src)
	lineWidth := o.LineWidth
	halfW := lineWidth / 2
	// 精确计算四边的轮廓范围（整数坐标，保证四边宽度一致）
	left := o.X0 - halfW
	right := o.X1 + halfW
	top := o.Y0 - halfW
	bottom := o.Y1 + halfW
	// 奇数线宽补偿，确保轮廓居中
	if lineWidth%2 != 0 {
		right++
		bottom++
	}
	for y := top; y < top+lineWidth; y++ {
		for x := left; x < right; x++ {
			setPixel(dst, x, y, o.Color, 255)
		}
	}
	for y := bottom - lineWidth; y < bottom; y++ {
		for x := left; x < right; x++ {
			setPixel(dst, x, y, o.Color, 255)
		}
	}
	for x := left; x < left+lineWidth; x++ {
		for y := top + lineWidth; y < bottom-lineWidth; y++ {
			setPixel(dst, x, y, o.Color, 255)
		}
	}
	for x := right - lineWidth; x < right; x++ {
		for y := top + lineWidth; y < bottom-lineWidth; y++ {
			setPixel(dst, x, y, o.Color, 255)
		}
	}
	return dst
}

// getPolygonAABB 计算多边形的包围盒
// vertices: 多边形顶点数组；返回minX, minY, maxX, maxY
func getPolygonAABB(vertices [][2]int) (int, int, int, int) {
	if len(vertices) == 0 {
		return 0, 0, 0, 0
	}
	minX, minY := vertices[0][0], vertices[0][1]
	maxX, maxY := vertices[0][0], vertices[0][1]
	for _, v := range vertices {
		if v[0] < minX {
			minX = v[0]
		}
		if v[0] > maxX {
			maxX = v[0]
		}
		if v[1] < minY {
			minY = v[1]
		}
		if v[1] > maxY {
			maxY = v[1]
		}
	}
	return minX, minY, maxX, maxY
}

// rayCasting 射线法判断点是否在多边形内部
// px,py: 待判断点；vertices: 多边形顶点数组（[[x0,y0],[x1,y1],...]）
func rayCasting(px, py float64, vertices [][2]int) bool {
	inside := false
	n := len(vertices)
	for i := 0; i < n; i++ {
		j := (i + 1) % n
		vix, viy := float64(vertices[i][0]), float64(vertices[i][1])
		vjx, vjy := float64(vertices[j][0]), float64(vertices[j][1])

		// 检查点是否在边的y范围内
		if (viy > py) != (vjy > py) {
			// 计算交点的x坐标
			t := (py - viy) / (vjy - viy)
			x := vix + t*(vjx-vix)
			if px < x {
				inside = !inside
			}
		}
	}
	return inside
}

// SolidPolygon 实心多边形：顶点数组 + 填充颜色
type SolidPolygon struct {
	Vertices [][2]int   // 顶点数组（[[x0,y0],[x1,y1],...]，至少3个顶点）
	Color    color.RGBA // 填充颜色
}

// NewSolidPolygon 创建实心多边形（参数校验）
func NewSolidPolygon(vertices [][2]int, c color.RGBA) *SolidPolygon {
	if len(vertices) < 3 {
		log.Println("polygon must have at least 3 vertices")
		return nil
	}
	return &SolidPolygon{
		Vertices: vertices,
		Color:    c,
	}
}

func (s *SolidPolygon) GetWH() (int, int) {
	x, y := make([]int, 0), make([]int, 0)
	for _, v := range s.Vertices {
		x = append(x, v[0])
		y = append(y, v[1])
	}
	return maxValue(x...), maxValue(y...)
}

// Render 实心多边形（射线法，锐利边缘，无多余像素）
func (s *SolidPolygon) Render(src image.Image) image.Image {
	// 1. 创建与src同尺寸的目标图像，复制背景
	dst := image.NewRGBA(src.Bounds())
	draw.Draw(dst, dst.Bounds(), src, src.Bounds().Min, draw.Src)

	// 2. 计算包围盒，优化遍历范围
	minX, minY, maxX, maxY := getPolygonAABB(s.Vertices)

	// 3. 遍历包围盒内的像素，仅绘制内部像素
	for y := minY; y < maxY; y++ {
		for x := minX; x < maxX; x++ {
			// 射线法判断是否在内部（整数坐标，锐利边缘）
			if rayCasting(float64(x)+0.5, float64(y)+0.5, s.Vertices) {
				setPixel(dst, x, y, s.Color, 255)
			}
		}
	}
	return dst
}

// OutlinePolygon 非实心多边形：顶点数组 + 轮廓颜色 + 线条宽度
type OutlinePolygon struct {
	Vertices  [][2]int   // 顶点数组（至少3个顶点）
	Color     color.RGBA // 轮廓颜色
	LineWidth int        // 线条宽度（最小1）
}

// NewOutlinePolygon 创建非实心多边形（参数校验）
func NewOutlinePolygon(vertices [][2]int, lineWidth int, c color.RGBA) *OutlinePolygon {
	if len(vertices) < 3 {
		panic("polygon must have at least 3 vertices")
	}
	if lineWidth < 1 {
		lineWidth = 1
	}
	return &OutlinePolygon{
		Vertices:  vertices,
		LineWidth: lineWidth,
		Color:     c,
	}
}

func (o *OutlinePolygon) GetWH() (int, int) {
	x, y := make([]int, 0), make([]int, 0)
	for _, v := range o.Vertices {
		x = append(x, v[0])
		y = append(y, v[1])
	}
	return maxValue(x...) + o.LineWidth, maxValue(y...) + o.LineWidth
}

// 判断点是否在多边形顶点的有效尖角区域
func isPointInPolygonVertexFan(vx, vy int, prevV, nextV [2]int, lineWidth int, px, py float64) bool {
	// 转换为浮点数
	curX, curY := float64(vx), float64(vy)
	prevX, prevY := float64(prevV[0]), float64(prevV[1])
	nextX, nextY := float64(nextV[0]), float64(nextV[1])

	// 计算顶点到前后顶点的向量
	prevToCurX := curX - prevX
	prevToCurY := curY - prevY
	curToNextX := nextX - curX
	curToNextY := nextY - curY

	// 图像坐标系Y轴向下，叉乘符号反转，凸凹性判断需调整
	// 叉乘计算：prevToCur × curToNext
	crossVal := prevToCurX*curToNextY - prevToCurY*curToNextX
	// 图像坐标系的凸凹性：crossVal < 0 为凸角，>0 为凹角
	isConvex := crossVal < 0

	// 计算两条边的单位法向量
	// 边1（prev->cur）的法向量：垂直于边并指向外侧
	normPrevX, normPrevY := -prevToCurY, prevToCurX // 逆时针法向
	normPrevX, normPrevY = normalize(normPrevX, normPrevY)
	// 边2（cur->next）的法向量：垂直于边并指向外侧
	normNextX, normNextY := -curToNextY, curToNextX // 逆时针法向
	normNextX, normNextY = normalize(normNextX, normNextY)

	// 点到顶点的向量
	pToCurX := px - curX
	pToCurY := py - curY

	// 线宽的有效范围
	lineWidthF := float64(lineWidth) * 1.2
	if math.Hypot(pToCurX, pToCurY) > lineWidthF {
		return false
	}

	// 计算点在两条法向量上的投影
	// 投影值：正数表示点在法向量的外侧
	projPrev := pToCurX*normPrevX + pToCurY*normPrevY
	projNext := pToCurX*normNextX + pToCurY*normNextY

	// 向量点积：判断点是否在两条边的夹角内
	dotPrev := (pToCurX)*(-prevToCurX) + (pToCurY)*(-prevToCurY) // 与prev->cur的反向向量点积
	dotNext := (pToCurX)*(curToNextX) + (pToCurY)*(curToNextY)   // 与cur->next的向量点积
	if dotPrev < 0 || dotNext < 0 {
		return false
	}

	// 根据凸凹性判断尖角区域
	if isConvex {
		// 凸角：点需要在两条边的外侧法向量之间（夹角内部）
		return projPrev >= 0 && projNext >= 0
	} else {
		// 凹角：点需要在两条边的外侧法向量之外（夹角外部）
		return projPrev <= 0 || projNext <= 0
	}
}

// Render 非实心多边形渲染
func (o *OutlinePolygon) Render(src image.Image) image.Image {
	dst := image.NewRGBA(src.Bounds())
	draw.Draw(dst, dst.Bounds(), src, src.Bounds().Min, draw.Src)

	lineWidth := o.LineWidth
	halfWidth := float64(lineWidth) / 2.0
	// 线宽为奇数时的像素对齐（+0.5确保中心对齐）
	if lineWidth%2 != 0 {
		halfWidth = math.Floor(float64(lineWidth)/2) + 0.5
	}
	n := len(o.Vertices)

	minX, minY, maxX, maxY := getPolygonAABB(o.Vertices)
	traverseMinX := minX - lineWidth*2
	traverseMinY := minY - lineWidth*2
	traverseMaxX := maxX + lineWidth*2
	traverseMaxY := maxY + lineWidth*2

	// 遍历像素绘制轮廓和尖角
	for y := traverseMinY; y < traverseMaxY; y++ {
		for x := traverseMinX; x < traverseMaxX; x++ {
			px := float64(x) + 0.5 // 像素中心坐标
			py := float64(y) + 0.5
			needDraw := false

			// 检查是否在边的轮廓范围内（线宽的半宽）
			for i := 0; i < n; i++ {
				j := (i + 1) % n
				vix, viy := float64(o.Vertices[i][0]), float64(o.Vertices[i][1])
				vjx, vjy := float64(o.Vertices[j][0]), float64(o.Vertices[j][1])
				if pointToSegmentDist(px, py, vix, viy, vjx, vjy) < halfWidth {
					needDraw = true
					break
				}
			}

			// 检查是否在顶点的有效尖角区域
			if !needDraw {
				for i := 0; i < n; i++ {
					prevIdx := (i - 1 + n) % n
					nextIdx := (i + 1) % n
					curV := o.Vertices[i]
					prevV := o.Vertices[prevIdx]
					nextV := o.Vertices[nextIdx]
					if isPointInPolygonVertexFan(curV[0], curV[1], prevV, nextV, lineWidth, px, py) {
						needDraw = true
						break
					}
				}
			}

			if needDraw {
				setPixel(dst, x, y, o.Color, 255)
			}
		}
	}

	return dst
}

// SolidEllipse 实心椭圆：中心点 + X/Y轴半径 + 填充颜色 + 旋转角度
type SolidEllipse struct {
	CenterX, CenterY int        // 中心点坐标
	RadiusX, RadiusY int        // X轴和Y轴半径
	Color            color.RGBA // 填充颜色
	Rotation         float64    // 旋转角度（弧度）
}

// NewSolidEllipse 创建实心椭圆（参数校验）
func NewSolidEllipse(centerX, centerY, radiusX, radiusY int, rotation float64, c color.RGBA) *SolidEllipse {
	if radiusX <= 0 {
		radiusX = 1
	}
	if radiusY <= 0 {
		radiusY = 1
	}
	return &SolidEllipse{
		CenterX:  centerX,
		CenterY:  centerY,
		RadiusX:  radiusX,
		RadiusY:  radiusY,
		Rotation: rotation,
		Color:    c,
	}
}

func (s *SolidEllipse) GetWH() (int, int) {
	// todo 计算错误修正bug
	return s.CenterX + s.RadiusX, s.CenterY + s.RadiusY
}

// rotatePoint 旋转点坐标
func rotatePoint(x, y int, cx, cy int, angle float64) (int, int) {
	// 转换为相对于中心点的坐标
	dx := float64(x - cx)
	dy := float64(y - cy)

	// 应用旋转
	rotatedX := dx*math.Cos(angle) - dy*math.Sin(angle)
	rotatedY := dx*math.Sin(angle) + dy*math.Cos(angle)

	// 转换回绝对坐标
	return cx + int(rotatedX+0.5), cy + int(rotatedY+0.5)
}

// getEllipseBounds 计算椭圆的边界框
func getEllipseBounds(cx, cy, a, b int, angle float64) (int, int, int, int) {
	if a == 0 || b == 0 {
		return cx, cy, cx, cy
	}

	// 计算椭圆的四个顶点
	vertices := []struct{ x, y int }{
		{cx + a, cy},
		{cx - a, cy},
		{cx, cy + b},
		{cx, cy - b},
	}

	// 旋转顶点
	minX, minY := cx, cy
	maxX, maxY := cx, cy

	for _, v := range vertices {
		rx, ry := rotatePoint(v.x, v.y, cx, cy, angle)
		if rx < minX {
			minX = rx
		}
		if rx > maxX {
			maxX = rx
		}
		if ry < minY {
			minY = ry
		}
		if ry > maxY {
			maxY = ry
		}
	}

	return minX, minY, maxX, maxY
}

// getEllipseAABBWithLineWidth 计算带线宽的椭圆边界框
func getEllipseAABBWithLineWidth(cx, cy, a, b, lineWidth int, angle float64) (int, int, int, int) {
	minX, minY, maxX, maxY := getEllipseBounds(cx, cy, a, b, angle)
	halfWidth := lineWidth / 2
	return minX - halfWidth, minY - halfWidth, maxX + halfWidth, maxY + halfWidth
}

// Render 实心椭圆渲染 - 使用扫描线算法，解决矩形边问题
func (s *SolidEllipse) Render(src image.Image) image.Image {
	dst := image.NewRGBA(src.Bounds())
	draw.Draw(dst, dst.Bounds(), src, src.Bounds().Min, draw.Src)

	// 计算椭圆边界框
	minX, minY, maxX, maxY := getEllipseBounds(s.CenterX, s.CenterY, s.RadiusX, s.RadiusY, s.Rotation)

	// 扩大边界框以处理抗锯齿
	expand := 2
	minX -= expand
	minY -= expand
	maxX += expand
	maxY += expand

	// 预计算旋转相关的常量
	cosAngle := math.Cos(s.Rotation)
	sinAngle := math.Sin(s.Rotation)
	a2 := float64(s.RadiusX * s.RadiusX)
	b2 := float64(s.RadiusY * s.RadiusY)
	cx := float64(s.CenterX)
	cy := float64(s.CenterY)

	// 扫描线算法：逐行绘制椭圆
	for y := minY; y <= maxY; y++ {
		// 计算当前扫描线y相对于椭圆中心的位置
		dy := float64(y) - cy

		// 对于非旋转椭圆，可以直接计算x的范围
		if s.Rotation == 0 {
			// 椭圆方程：x²/a² + y²/b² = 1
			// 解得：x = ±a * sqrt(1 - y²/b²)
			term := 1.0 - (dy*dy)/b2
			if term <= 0 {
				continue // 没有交点
			}
			xRange := float64(s.RadiusX) * math.Sqrt(term)
			leftX := cx - xRange
			rightX := cx + xRange

			// 绘制当前扫描线上的椭圆部分
			drawScanline(dst, y, leftX, rightX, s.Color)
		} else {
			// 对于旋转椭圆，使用二分法查找每个y对应的x范围
			// 定义函数f(x) = ( (x-cx)cosθ + (y-cy)sinθ )² / a² +
			//                ( -(x-cx)sinθ + (y-cy)cosθ )² / b² - 1
			f := func(x float64) float64 {
				dx := x - cx
				rotX := dx*cosAngle + dy*sinAngle
				rotY := -dx*sinAngle + dy*cosAngle
				return (rotX*rotX)/a2 + (rotY*rotY)/b2 - 1.0
			}

			// 使用二分法查找左边界
			leftX := findRoot(f, float64(minX), float64(s.CenterX), 0.001)
			// 使用二分法查找右边界
			rightX := findRoot(f, float64(s.CenterX), float64(maxX), 0.001)

			if leftX != nil && rightX != nil {
				// 绘制当前扫描线上的椭圆部分
				drawScanline(dst, y, *leftX, *rightX, s.Color)
			}
		}
	}

	return dst
}

// drawScanline 绘制扫描线，处理抗锯齿
func drawScanline(dst *image.RGBA, y int, leftX, rightX float64, color color.RGBA) {
	// 转换为整数坐标
	startX := int(math.Floor(leftX))
	endX := int(math.Ceil(rightX))

	// 绘制完全在椭圆内部的像素
	for x := startX + 1; x < endX; x++ {
		setPixel(dst, x, y, color, 255)
	}

	// 处理左边缘的抗锯齿
	if startX >= dst.Bounds().Min.X && startX < dst.Bounds().Max.X {
		alpha := uint8(255 * (1.0 - (leftX - math.Floor(leftX))))
		setPixel(dst, startX, y, color, alpha)
	}

	// 处理右边缘的抗锯齿
	if endX >= dst.Bounds().Min.X && endX < dst.Bounds().Max.X {
		alpha := uint8(255 * (rightX - math.Floor(rightX)))
		setPixel(dst, endX, y, color, alpha)
	}
}

// findRoot 使用二分法查找函数f在区间[a, b]内的根
func findRoot(f func(float64) float64, a, b, epsilon float64) *float64 {
	fa := f(a)
	fb := f(b)

	// 如果两端点函数值同号，说明没有根
	if fa*fb > 0 {
		return nil
	}

	// 二分法迭代
	for i := 0; i < 100; i++ {
		mid := (a + b) / 2
		fmid := f(mid)

		// 如果找到足够接近的根
		if math.Abs(fmid) < epsilon {
			return &mid
		}

		// 调整搜索区间
		if fa*fmid < 0 {
			b = mid
			fb = fmid
		} else {
			a = mid
			fa = fmid
		}
	}

	// 返回最终的中点作为根的近似值
	result := (a + b) / 2
	return &result
}

// OutlineEllipse 非实心椭圆：中心点 + X/Y轴半径 + 轮廓颜色 + 线条宽度 + 旋转角度
type OutlineEllipse struct {
	CenterX, CenterY int        // 中心点坐标
	RadiusX, RadiusY int        // X轴和Y轴半径
	Color            color.RGBA // 轮廓颜色
	LineWidth        int        // 线条宽度（最小1）
	Rotation         float64    // 旋转角度（弧度）
}

// NewOutlineEllipse 创建非实心椭圆（参数校验）
func NewOutlineEllipse(centerX, centerY, radiusX, radiusY, lineWidth int, rotation float64, c color.RGBA) *OutlineEllipse {
	if radiusX <= 0 {
		radiusX = 1
	}
	if radiusY <= 0 {
		radiusY = 1
	}
	if lineWidth < 1 {
		lineWidth = 1
	}
	return &OutlineEllipse{
		CenterX:   centerX,
		CenterY:   centerY,
		RadiusX:   radiusX,
		RadiusY:   radiusY,
		LineWidth: lineWidth,
		Rotation:  rotation,
		Color:     c,
	}
}

func (o *OutlineEllipse) GetWH() (int, int) {
	// todo 计算错误修正bug
	return o.CenterX + o.RadiusX, o.CenterY + o.RadiusY
}

// Render 非实心椭圆渲染 - 使用扫描线算法，解决矩形边问题
func (o *OutlineEllipse) Render(src image.Image) image.Image {
	dst := image.NewRGBA(src.Bounds())
	draw.Draw(dst, dst.Bounds(), src, src.Bounds().Min, draw.Src)

	halfWidth := float64(o.LineWidth) / 2.0
	// 修正：线宽为奇数时的像素对齐
	if o.LineWidth%2 != 0 {
		halfWidth = math.Floor(float64(o.LineWidth)/2) + 0.5
	}

	// 计算带线宽的椭圆边界框
	minX, minY, maxX, maxY := getEllipseAABBWithLineWidth(
		o.CenterX, o.CenterY, o.RadiusX, o.RadiusY, o.LineWidth, o.Rotation)

	// 扩大边界框以处理抗锯齿
	expand := 2
	minX -= expand
	minY -= expand
	maxX += expand
	maxY += expand

	// 预计算旋转相关的常量
	cosAngle := math.Cos(o.Rotation)
	sinAngle := math.Sin(o.Rotation)
	//a2 := float64(o.RadiusX * o.RadiusX)
	//b2 := float64(o.RadiusY * o.RadiusY)
	cx := float64(o.CenterX)
	cy := float64(o.CenterY)

	// 外椭圆和内椭圆的半径
	outerRadiusX := float64(o.RadiusX) + halfWidth
	outerRadiusY := float64(o.RadiusY) + halfWidth
	innerRadiusX := float64(o.RadiusX) - halfWidth
	innerRadiusY := float64(o.RadiusY) - halfWidth

	// 如果内椭圆半径小于0，则绘制实心椭圆
	if innerRadiusX < 0 || innerRadiusY < 0 {
		solidEllipse := &SolidEllipse{
			CenterX:  o.CenterX,
			CenterY:  o.CenterY,
			RadiusX:  o.RadiusX,
			RadiusY:  o.RadiusY,
			Rotation: o.Rotation,
			Color:    o.Color,
		}
		return solidEllipse.Render(src)
	}

	// 扫描线算法：逐行绘制椭圆轮廓
	for y := minY; y <= maxY; y++ {
		// 计算当前扫描线y相对于椭圆中心的位置
		dy := float64(y) - cy

		// 定义外椭圆函数
		outerEllipse := func(x float64) float64 {
			dx := x - cx
			rotX := dx*cosAngle + dy*sinAngle
			rotY := -dx*sinAngle + dy*cosAngle
			return (rotX*rotX)/(outerRadiusX*outerRadiusX) + (rotY*rotY)/(outerRadiusY*outerRadiusY) - 1.0
		}

		// 定义内椭圆函数
		innerEllipse := func(x float64) float64 {
			dx := x - cx
			rotX := dx*cosAngle + dy*sinAngle
			rotY := -dx*sinAngle + dy*cosAngle
			return (rotX*rotX)/(innerRadiusX*innerRadiusX) + (rotY*rotY)/(innerRadiusY*innerRadiusY) - 1.0
		}

		// 查找外椭圆的左右边界
		outerLeft := findRoot(outerEllipse, float64(minX), float64(cx), 0.001)
		outerRight := findRoot(outerEllipse, float64(cx), float64(maxX), 0.001)

		// 查找内椭圆的左右边界
		innerLeft := findRoot(innerEllipse, float64(minX), float64(cx), 0.001)
		innerRight := findRoot(innerEllipse, float64(cx), float64(maxX), 0.001)

		// 绘制当前扫描线上的椭圆轮廓
		drawOutlineScanline(dst, y, outerLeft, outerRight, innerLeft, innerRight, o.Color)
	}

	return dst
}

// drawOutlineScanline 绘制非实心椭圆的扫描线
func drawOutlineScanline(dst *image.RGBA, y int,
	outerLeft, outerRight, innerLeft, innerRight *float64,
	color color.RGBA) {

	// 填充外椭圆和内椭圆之间的区域（先填充，再处理边缘抗锯齿）
	if outerLeft != nil && outerRight != nil {
		leftX := *outerLeft
		rightX := *outerRight

		// 如果有内椭圆，调整填充范围
		if innerLeft != nil && innerRight != nil {
			// 左半部分：外椭圆左边缘到内椭圆左边缘
			drawScanlineSegment(dst, y, leftX, *innerLeft, color)
			// 右半部分：内椭圆右边缘到外椭圆右边缘
			drawScanlineSegment(dst, y, *innerRight, rightX, color)
		} else {
			// 没有内椭圆，填充整个外椭圆
			drawScanlineSegment(dst, y, leftX, rightX, color)
		}
	}

	// 处理外椭圆左边缘的抗锯齿
	if outerLeft != nil {
		ellipseX := *outerLeft
		startX := int(math.Floor(ellipseX))
		endX := int(math.Ceil(ellipseX))

		// 绘制左边缘的抗锯齿
		for pixelX := startX; pixelX <= endX; pixelX++ {
			if pixelX >= dst.Bounds().Min.X && pixelX < dst.Bounds().Max.X {
				// 计算Alpha值，基于像素中心到椭圆边界的距离
				distance := math.Abs(float64(pixelX) + 0.5 - ellipseX)
				alpha := uint8(255 * (1.0 - distance))
				// 获取背景颜色
				bgColor := getPixelColor(dst, pixelX, y)
				// 混合前景色和背景色
				mixedColor := blendColors(color, bgColor, alpha)
				setPixel(dst, pixelX, y, mixedColor, 255)
			}
		}
	}

	// 处理外椭圆右边缘的抗锯齿
	if outerRight != nil {
		ellipseX := *outerRight
		startX := int(math.Floor(ellipseX))
		endX := int(math.Ceil(ellipseX))

		// 绘制右边缘的抗锯齿
		for pixelX := startX; pixelX <= endX; pixelX++ {
			if pixelX >= dst.Bounds().Min.X && pixelX < dst.Bounds().Max.X {
				// 计算Alpha值，基于像素中心到椭圆边界的距离
				distance := math.Abs(float64(pixelX) + 0.5 - ellipseX)
				alpha := uint8(255 * (1.0 - distance))
				// 获取背景颜色
				bgColor := getPixelColor(dst, pixelX, y)
				// 混合前景色和背景色
				mixedColor := blendColors(color, bgColor, alpha)
				setPixel(dst, pixelX, y, mixedColor, 255)
			}
		}
	}
}

// getPixelColor 获取像素的当前颜色
func getPixelColor(dst *image.RGBA, x, y int) color.RGBA {
	bounds := dst.Bounds()
	if x < bounds.Min.X || x >= bounds.Max.X || y < bounds.Min.Y || y >= bounds.Max.Y {
		return color.RGBA{255, 255, 255, 255} // 默认返回白色
	}
	// RGBA的像素索引计算
	idx := (y-dst.Rect.Min.Y)*dst.Stride + (x-dst.Rect.Min.X)*4
	return color.RGBA{
		R: dst.Pix[idx],
		G: dst.Pix[idx+1],
		B: dst.Pix[idx+2],
		A: dst.Pix[idx+3],
	}
}

// blendColors 混合前景色和背景色
func blendColors(foreground, background color.RGBA, alpha uint8) color.RGBA {
	// 计算前景色的权重
	fgAlpha := float64(alpha) / 255.0
	// 计算背景色的权重
	bgAlpha := 1.0 - fgAlpha

	// 混合RGB通道
	r := uint8(float64(foreground.R)*fgAlpha + float64(background.R)*bgAlpha)
	g := uint8(float64(foreground.G)*fgAlpha + float64(background.G)*bgAlpha)
	b := uint8(float64(foreground.B)*fgAlpha + float64(background.B)*bgAlpha)

	return color.RGBA{R: r, G: g, B: b, A: 255}
}

// drawScanlineSegment 绘制扫描线上的一段
func drawScanlineSegment(dst *image.RGBA, y int, startX, endX float64, color color.RGBA) {
	// 确保startX < endX
	if startX > endX {
		startX, endX = endX, startX
	}

	// 转换为整数坐标
	intStartX := int(math.Floor(startX))
	intEndX := int(math.Ceil(endX))

	// 绘制完全在内部的像素
	for x := intStartX + 1; x < intEndX; x++ {
		setPixel(dst, x, y, color, 255)
	}

	// 处理左边缘的抗锯齿
	if intStartX >= dst.Bounds().Min.X && intStartX < dst.Bounds().Max.X {
		alpha := uint8(255 * (1.0 - (startX - math.Floor(startX))))
		setPixel(dst, intStartX, y, color, alpha)
	}

	// 处理右边缘的抗锯齿
	if intEndX >= dst.Bounds().Min.X && intEndX < dst.Bounds().Max.X {
		alpha := uint8(255 * (endX - math.Floor(endX)))
		setPixel(dst, intEndX, y, color, alpha)
	}
}

// Sector 扇形：中心点 + X/Y轴半径 + 起始角度 + 终止角度 + 填充颜色 + 旋转角度
type Sector struct {
	CenterX, CenterY int        // 中心点坐标
	RadiusX, RadiusY int        // X轴和Y轴半径
	StartAngle       float64    // 起始角度（弧度）
	EndAngle         float64    // 终止角度（弧度）
	Color            color.RGBA // 填充颜色
	Rotation         float64    // 旋转角度（弧度）
}

// NewSector 创建扇形（参数校验）
func NewSector(centerX, centerY, radiusX, radiusY int, startAngle, endAngle, rotation float64, c color.RGBA) *Sector {
	if radiusX <= 0 {
		radiusX = 1
	}
	if radiusY <= 0 {
		radiusY = 1
	}
	// 确保起始角度小于终止角度
	if startAngle > endAngle {
		startAngle, endAngle = endAngle, startAngle
	}
	return &Sector{
		CenterX:    centerX,
		CenterY:    centerY,
		RadiusX:    radiusX,
		RadiusY:    radiusY,
		StartAngle: startAngle,
		EndAngle:   endAngle,
		Rotation:   rotation,
		Color:      c,
	}
}

func (s *Sector) GetWH() (int, int) {
	// todo 计算错误修正bug
	return s.CenterX + s.RadiusX, s.CenterY + s.RadiusY
}

// Edge 表示扫描线算法中的边
type Edge struct {
	YMin   int     // 边的最小y坐标
	YMax   int     // 边的最大y坐标
	X      float64 // 当前扫描线与边的交点x坐标
	DeltaX float64 // 扫描线y增加1时，x的增量
}

// Render 扇形渲染 - 使用扫描线算法
func (s *Sector) Render(src image.Image) image.Image {
	// 如果源图像为nil，创建一个新的RGBA图像
	var dst *image.RGBA
	if src == nil {
		// 创建一个足够大的图像来容纳整个扇形
		width := s.CenterX + s.RadiusX + 10
		height := s.CenterY + s.RadiusY + 10
		dst = image.NewRGBA(image.Rect(0, 0, width, height))
	} else {
		// 将源图像转换为RGBA格式
		bounds := src.Bounds()
		dst = image.NewRGBA(bounds)
		draw.Draw(dst, bounds, src, bounds.Min, draw.Src)
	}

	// 获取图像边界
	bounds := dst.Bounds()
	width := bounds.Max.X
	height := bounds.Max.Y

	// 计算扇形的边界框，确定需要扫描的y范围
	yMin, yMax := s.calculateYRange()

	// 为每个y坐标创建边表
	edgeTable := make(map[int][]Edge)

	// 添加扇形的两条半径边
	s.addRadiusEdges(edgeTable)

	// 添加扇形的圆弧边（使用多边形逼近）
	s.addArcEdges(edgeTable)

	// 初始化活性边表
	activeEdges := make([]Edge, 0)

	// 扫描线算法主循环
	for y := yMin; y <= yMax; y++ {
		// 如果y超出图像边界，跳过
		if y < 0 || y >= height {
			continue
		}

		// 将当前y对应的边添加到活性边表
		if edges, ok := edgeTable[y]; ok {
			activeEdges = append(activeEdges, edges...)
		}

		// 从活性边表中移除yMax <= y的边
		newActiveEdges := make([]Edge, 0)
		for _, edge := range activeEdges {
			if edge.YMax > y {
				newActiveEdges = append(newActiveEdges, edge)
			}
		}
		activeEdges = newActiveEdges

		// 更新活性边表中每条边的x坐标
		for i := range activeEdges {
			activeEdges[i].X += activeEdges[i].DeltaX
		}

		// 按x坐标排序活性边
		sortEdgesByX(activeEdges)

		// 填充扫描线
		s.fillScanLine(dst, y, activeEdges, width)
	}

	return dst
}

// calculateYRange 计算扇形的y坐标范围
func (s *Sector) calculateYRange() (int, int) {
	// 计算扇形四个极端点的y坐标
	points := []struct{ angle float64 }{
		{angle: s.StartAngle},
		{angle: s.EndAngle},
		{angle: s.StartAngle + math.Pi},
		{angle: s.EndAngle + math.Pi},
	}

	yCoords := make([]int, len(points))
	for i, p := range points {
		_, y := s.calculatePoint(p.angle)
		yCoords[i] = y
	}

	// 找到y的最小值和最大值
	yMin := yCoords[0]
	yMax := yCoords[0]
	for _, y := range yCoords[1:] {
		if y < yMin {
			yMin = y
		}
		if y > yMax {
			yMax = y
		}
	}

	return yMin, yMax
}

// calculatePoint 根据角度计算扇形上的点
func (s *Sector) calculatePoint(angle float64) (int, int) {
	// 计算椭圆上的点
	x := float64(s.RadiusX) * math.Cos(angle)
	y := float64(s.RadiusY) * math.Sin(angle)

	// 应用旋转变换
	rotatedX := x*math.Cos(s.Rotation) - y*math.Sin(s.Rotation)
	rotatedY := x*math.Sin(s.Rotation) + y*math.Cos(s.Rotation)

	// 转换为图像坐标
	return s.CenterX + int(math.Round(rotatedX)), s.CenterY + int(math.Round(rotatedY))
}

// addRadiusEdges 添加扇形的两条半径边到边表
func (s *Sector) addRadiusEdges(edgeTable map[int][]Edge) {
	// 计算起始半径的两个端点
	startX, startY := s.CenterX, s.CenterY
	endX1, endY1 := s.calculatePoint(s.StartAngle)

	// 添加起始半径边
	s.addEdge(edgeTable, startX, startY, endX1, endY1)

	// 计算终止半径的两个端点
	endX2, endY2 := s.calculatePoint(s.EndAngle)

	// 添加终止半径边
	s.addEdge(edgeTable, startX, startY, endX2, endY2)
}

// addArcEdges 使用多边形逼近添加扇形的圆弧边到边表
func (s *Sector) addArcEdges(edgeTable map[int][]Edge) {
	// 计算角度差
	angleDiff := s.EndAngle - s.StartAngle

	// 根据角度差确定逼近的步数（每10度一个点）
	steps := int(math.Ceil(angleDiff * 180 / math.Pi / 10))
	if steps < 2 {
		steps = 2
	}

	angleStep := angleDiff / float64(steps)

	// 生成圆弧上的点
	prevX, prevY := s.calculatePoint(s.StartAngle)
	for i := 1; i <= steps; i++ {
		angle := s.StartAngle + float64(i)*angleStep
		currX, currY := s.calculatePoint(angle)

		// 添加圆弧边
		s.addEdge(edgeTable, prevX, prevY, currX, currY)

		prevX, prevY = currX, currY
	}
}

// addEdge 添加一条边到边表
func (s *Sector) addEdge(edgeTable map[int][]Edge, x1, y1, x2, y2 int) {
	// 忽略水平边
	if y1 == y2 {
		return
	}

	// 确保y1 < y2
	if y1 > y2 {
		x1, x2 = x2, x1
		y1, y2 = y2, y1
	}

	// 计算边的斜率倒数（deltaX）
	deltaY := y2 - y1
	deltaX := float64(x2-x1) / float64(deltaY)

	// 创建边
	edge := Edge{
		YMin:   y1,
		YMax:   y2,
		X:      float64(x1),
		DeltaX: deltaX,
	}

	// 添加到边表
	edgeTable[y1] = append(edgeTable[y1], edge)
}

// sortEdgesByX 按x坐标排序边
func sortEdgesByX(edges []Edge) {
	// 使用冒泡排序，简单但效率不高，适合小规模数据
	for i := 0; i < len(edges); i++ {
		for j := i + 1; j < len(edges); j++ {
			if edges[i].X > edges[j].X {
				edges[i], edges[j] = edges[j], edges[i]
			}
		}
	}
}

// fillScanLine 填充扫描线
func (s *Sector) fillScanLine(dst *image.RGBA, y int, edges []Edge, width int) {
	// 两两配对边，填充之间的像素
	for i := 0; i < len(edges); i += 2 {
		// 确保有配对的边
		if i+1 >= len(edges) {
			break
		}

		// 获取左右两条边的x坐标
		leftX := edges[i].X
		rightX := edges[i+1].X

		// 转换为整数像素坐标
		xStart := int(math.Round(leftX))
		xEnd := int(math.Round(rightX))

		// 确保x坐标在图像范围内
		if xStart < 0 {
			xStart = 0
		}
		if xEnd >= width {
			xEnd = width - 1
		}

		// 填充从xStart到xEnd的像素
		for x := xStart; x <= xEnd; x++ {
			dst.SetRGBA(x, y, s.Color)
		}
	}
}

// isPointInSector 检查点是否在扇形内
func (s *Sector) isPointInSector(x, y int) bool {
	// 计算点到中心点的向量
	dx := x - s.CenterX
	dy := y - s.CenterY

	// 如果点就是中心点，返回true
	if dx == 0 && dy == 0 {
		return true
	}

	// 计算点到中心点的距离
	distance := math.Hypot(float64(dx), float64(dy))

	// 如果距离为0，返回true
	if distance == 0 {
		return true
	}

	// 计算点相对于中心点的角度
	angle := math.Atan2(float64(dy), float64(dx))

	// 应用逆旋转变换（将点旋转回未旋转的坐标系）
	rotatedAngle := angle - s.Rotation

	// 确保角度在[0, 2π)范围内
	for rotatedAngle < 0 {
		rotatedAngle += 2 * math.Pi
	}
	for rotatedAngle >= 2*math.Pi {
		rotatedAngle -= 2 * math.Pi
	}

	// 确保扇形的起始和终止角度在[0, 2π)范围内
	startAngle := s.StartAngle
	for startAngle < 0 {
		startAngle += 2 * math.Pi
	}
	endAngle := s.EndAngle
	for endAngle < 0 {
		endAngle += 2 * math.Pi
	}

	// 检查角度是否在扇形范围内
	var inAngleRange bool
	if startAngle <= endAngle {
		inAngleRange = rotatedAngle >= startAngle && rotatedAngle <= endAngle
	} else {
		inAngleRange = rotatedAngle >= startAngle || rotatedAngle <= endAngle
	}

	if !inAngleRange {
		return false
	}

	// 计算点在椭圆上的期望距离
	// 椭圆方程：(x^2/a^2) + (y^2/b^2) = 1
	// 其中a是RadiusX，b是RadiusY
	// 点(x,y)到中心的向量为(dx, dy)
	// 旋转后的向量为(dx*cosθ + dy*sinθ, -dx*sinθ + dy*cosθ)
	rotatedDx := float64(dx)*math.Cos(s.Rotation) + float64(dy)*math.Sin(s.Rotation)
	rotatedDy := -float64(dx)*math.Sin(s.Rotation) + float64(dy)*math.Cos(s.Rotation)

	// 计算点到中心的距离与椭圆上同方向点的距离的比值
	ellipseRatio := (rotatedDx*rotatedDx)/(float64(s.RadiusX)*float64(s.RadiusX)) +
		(rotatedDy*rotatedDy)/(float64(s.RadiusY)*float64(s.RadiusY))

	// 如果比值小于等于1，点在椭圆内
	return ellipseRatio <= 1.0
}

// Star 星形：中心点 + 外半径 + 内半径 + 角数 + 填充颜色 + 旋转角度
type Star struct {
	CenterX, CenterY int        // 中心点坐标
	OuterRadius      int        // 外半径
	InnerRadius      int        // 内半径
	Points           int        // 角数
	Color            color.RGBA // 填充颜色
	Rotation         float64    // 旋转角度（弧度）
}

func (s *Star) GetWH() (int, int) {
	// todo 计算错误修正bug
	return s.CenterX, s.CenterY
}

// NewStar 创建星形（参数校验）
func NewStar(centerX, centerY, outerRadius, innerRadius, points int, rotation float64, c color.RGBA) *Star {
	if outerRadius <= 0 {
		outerRadius = 1
	}
	if innerRadius <= 0 {
		innerRadius = 1
	}
	if points < 3 {
		points = 5 // 默认5角星
	}
	if innerRadius >= outerRadius {
		innerRadius = outerRadius / 2
	}
	return &Star{
		CenterX:     centerX,
		CenterY:     centerY,
		OuterRadius: outerRadius,
		InnerRadius: innerRadius,
		Points:      points,
		Rotation:    rotation,
		Color:       c,
	}
}

// Render 星形渲染 - 使用扫描线算法
func (s *Star) Render(src image.Image) image.Image {
	// 如果源图像为nil，创建一个新的RGBA图像
	var dst *image.RGBA
	if src == nil {
		// 创建一个足够大的图像来容纳整个星形
		width := s.CenterX + s.OuterRadius + 10
		height := s.CenterY + s.OuterRadius + 10
		dst = image.NewRGBA(image.Rect(0, 0, width, height))
	} else {
		// 将源图像转换为RGBA格式
		bounds := src.Bounds()
		dst = image.NewRGBA(bounds)
		draw.Draw(dst, bounds, src, bounds.Min, draw.Src)
	}

	// 获取图像边界
	bounds := dst.Bounds()
	width := bounds.Max.X
	height := bounds.Max.Y

	// 计算星形的顶点
	points := s.calculatePoints()

	// 计算星形的边界框，确定需要扫描的y范围
	yMin, yMax := s.calculateYRange(points)

	// 为每个y坐标创建边表
	edgeTable := make(map[int][]Edge)

	// 添加星形的边到边表
	s.addEdges(edgeTable, points)

	// 初始化活性边表
	activeEdges := make([]Edge, 0)

	// 扫描线算法主循环
	for y := yMin; y <= yMax; y++ {
		// 如果y超出图像边界，跳过
		if y < 0 || y >= height {
			continue
		}

		// 将当前y对应的边添加到活性边表
		if edges, ok := edgeTable[y]; ok {
			activeEdges = append(activeEdges, edges...)
		}

		// 从活性边表中移除yMax <= y的边
		newActiveEdges := make([]Edge, 0)
		for _, edge := range activeEdges {
			if edge.YMax > y {
				newActiveEdges = append(newActiveEdges, edge)
			}
		}
		activeEdges = newActiveEdges

		// 更新活性边表中每条边的x坐标
		for i := range activeEdges {
			activeEdges[i].X += activeEdges[i].DeltaX
		}

		// 按x坐标排序活性边
		sortEdgesByX(activeEdges)

		// 填充扫描线
		s.fillScanLine(dst, y, activeEdges, width)
	}

	return dst
}

// calculatePoints 计算星形的顶点坐标
func (s *Star) calculatePoints() []image.Point {
	points := make([]image.Point, 0, s.Points*2)

	for i := 0; i < s.Points*2; i++ {
		// 计算角度
		angle := float64(i)*math.Pi/float64(s.Points) + s.Rotation

		// 交替使用外半径和内半径
		radius := s.OuterRadius
		if i%2 == 1 {
			radius = s.InnerRadius
		}

		// 计算顶点坐标
		x := float64(s.CenterX) + float64(radius)*math.Cos(angle-math.Pi/2)
		y := float64(s.CenterY) + float64(radius)*math.Sin(angle-math.Pi/2)

		points = append(points, image.Point{X: int(math.Round(x)), Y: int(math.Round(y))})
	}

	return points
}

// calculateYRange 计算星形的y坐标范围
func (s *Star) calculateYRange(points []image.Point) (int, int) {
	if len(points) == 0 {
		return 0, 0
	}

	yMin := points[0].Y
	yMax := points[0].Y

	for _, p := range points[1:] {
		if p.Y < yMin {
			yMin = p.Y
		}
		if p.Y > yMax {
			yMax = p.Y
		}
	}

	return yMin, yMax
}

// addEdges 添加星形的边到边表
func (s *Star) addEdges(edgeTable map[int][]Edge, points []image.Point) {
	if len(points) < 3 {
		return
	}

	for i := 0; i < len(points); i++ {
		x1, y1 := points[i].X, points[i].Y
		x2, y2 := points[(i+1)%len(points)].X, points[(i+1)%len(points)].Y

		s.addEdge(edgeTable, x1, y1, x2, y2)
	}
}

// addEdge 添加一条边到边表
func (s *Star) addEdge(edgeTable map[int][]Edge, x1, y1, x2, y2 int) {
	// 忽略水平边
	if y1 == y2 {
		return
	}

	// 确保y1 < y2
	if y1 > y2 {
		x1, x2 = x2, x1
		y1, y2 = y2, y1
	}

	// 计算边的斜率倒数（deltaX）
	deltaY := y2 - y1
	deltaX := float64(x2-x1) / float64(deltaY)

	// 创建边
	edge := Edge{
		YMin:   y1,
		YMax:   y2,
		X:      float64(x1),
		DeltaX: deltaX,
	}

	// 添加到边表
	edgeTable[y1] = append(edgeTable[y1], edge)
}

// fillScanLine 填充扫描线
func (s *Star) fillScanLine(dst *image.RGBA, y int, edges []Edge, width int) {
	// 两两配对边，填充之间的像素
	for i := 0; i < len(edges); i += 2 {
		// 确保有配对的边
		if i+1 >= len(edges) {
			break
		}

		// 获取左右两条边的x坐标
		leftX := edges[i].X
		rightX := edges[i+1].X

		// 转换为整数像素坐标
		xStart := int(math.Round(leftX))
		xEnd := int(math.Round(rightX))

		// 确保x坐标在图像范围内
		if xStart < 0 {
			xStart = 0
		}
		if xEnd >= width {
			xEnd = width - 1
		}

		// 填充从xStart到xEnd的像素
		for x := xStart; x <= xEnd; x++ {
			dst.SetRGBA(x, y, s.Color)
		}
	}
}

// CurveType 定义曲线类型
type CurveType int

const (
	QuadraticBezier CurveType = iota // 二次贝塞尔曲线
	CubicBezier                      // 三次贝塞尔曲线
)

// Curve 曲线图形：包含自身的属性（起点、终点、控制点、颜色、线宽、曲线类型）
type Curve struct {
	X0, Y0     int        // 起点
	X1, Y1     int        // 终点
	Cp1x, Cp1y int        // 第一个控制点（二次/三次贝塞尔）
	Cp2x, Cp2y int        // 第二个控制点（仅三次贝塞尔）
	Color      color.RGBA // 自身颜色
	LineWidth  int        // 直线粗度（像素数，最小为1）
	CurveType  CurveType  // 曲线类型
}

// NewQuadraticBezier 创建二次贝塞尔曲线
func NewQuadraticBezier(x0, y0, x1, y1, cp1x, cp1y int, c color.RGBA, lineWidth int) *Curve {
	if lineWidth < 1 {
		lineWidth = 1
	}
	return &Curve{
		X0: x0, Y0: y0,
		X1: x1, Y1: y1,
		Cp1x: cp1x, Cp1y: cp1y,
		Color:     c,
		LineWidth: lineWidth,
		CurveType: QuadraticBezier,
	}
}

// NewCubicBezier 创建三次贝塞尔曲线
func NewCubicBezier(x0, y0, x1, y1, cp1x, cp1y, cp2x, cp2y int, c color.RGBA, lineWidth int) *Curve {
	if lineWidth < 1 {
		lineWidth = 1
	}
	return &Curve{
		X0: x0, Y0: y0,
		X1: x1, Y1: y1,
		Cp1x: cp1x, Cp1y: cp1y,
		Cp2x: cp2x, Cp2y: cp2y,
		Color:     c,
		LineWidth: lineWidth,
		CurveType: CubicBezier,
	}
}

func (c *Curve) GetWH() (int, int) {
	return c.X1, c.Y1
}

// Render 渲染曲线
func (c *Curve) Render(src image.Image) image.Image {
	dst := image.NewRGBA(image.Rect(0, 0, src.Bounds().Dx(), src.Bounds().Dy()))
	bgColor := color.RGBA{R: 0, G: 0, B: 0, A: 0}
	draw.Draw(dst, dst.Bounds(), &image.Uniform{C: bgColor}, image.Point{}, draw.Over)

	// 计算曲线的包围盒
	bounds := c.calculateBoundingBox()

	// 转换为图像整数坐标范围
	imgBounds := src.Bounds()
	startX := int(math.Max(float64(bounds.Min.X), float64(imgBounds.Min.X)))
	endX := int(math.Min(float64(bounds.Max.X), float64(imgBounds.Max.X)))
	startY := int(math.Max(float64(bounds.Min.Y), float64(imgBounds.Min.Y)))
	endY := int(math.Min(float64(bounds.Max.Y), float64(imgBounds.Max.Y)))

	// 预计算平方值，减少浮点运算
	halfWidth := float64(c.LineWidth) / 2.0
	halfWidthMin05Sq := math.Pow(halfWidth-0.5, 2)
	halfWidthAdd05Sq := math.Pow(halfWidth+0.5, 2)

	// 生成曲线上的采样点
	samplePoints := c.generateSamplePoints()
	if len(samplePoints) < 2 {
		// 没有足够的采样点，无法绘制曲线
		draw.Draw(dst, imgBounds, src, image.Point{}, draw.Over)
		return dst
	}

	// 遍历包围盒内的所有像素，逐像素计算透明度并绘制
	for y := startY; y <= endY; y++ {
		for x := startX; x <= endX; x++ {
			// 亚像素坐标 - 使用像素中心计算
			px := float64(x) + 0.5
			py := float64(y) + 0.5

			// 计算像素到曲线的最短距离
			minDistSq := c.calculateMinDistanceSq(px, py, samplePoints)
			minDist := math.Sqrt(minDistSq)

			// 计算像素的Alpha透明度
			var alpha uint8
			if minDistSq < halfWidthMin05Sq {
				alpha = 255
			} else if minDistSq > halfWidthAdd05Sq {
				alpha = 0
			} else {
				// 浮点精度优化，避免微小误差导致的Alpha抖动
				alphaVal := (halfWidth + 0.5 - minDist) * 255
				if alphaVal < 0 {
					alphaVal = 0
				} else if alphaVal > 255 {
					alphaVal = 255
				}
				alpha = uint8(alphaVal)
			}

			// 绘制像素（支持任意颜色的Alpha混合）
			if alpha > 0 {
				setPixel(dst, x, y, c.Color, alpha)
			}
		}
	}

	draw.Draw(dst, imgBounds, src, image.Point{}, draw.Over)

	return dst
}

// calculateBoundingBox 计算曲线的包围盒
func (c *Curve) calculateBoundingBox() image.Rectangle {
	// 获取所有控制点和端点
	points := []image.Point{
		{c.X0, c.Y0},
		{c.X1, c.Y1},
		{c.Cp1x, c.Cp1y},
	}

	if c.CurveType == CubicBezier {
		points = append(points, image.Point{c.Cp2x, c.Cp2y})
	}

	// 初始化边界
	minX, minY := points[0].X, points[0].Y
	maxX, maxY := points[0].X, points[0].Y

	// 找到最小和最大坐标
	for _, p := range points[1:] {
		if p.X < minX {
			minX = p.X
		}
		if p.X > maxX {
			maxX = p.X
		}
		if p.Y < minY {
			minY = p.Y
		}
		if p.Y > maxY {
			maxY = p.Y
		}
	}

	// 考虑线宽
	halfWidth := float64(c.LineWidth) / 2.0
	minX -= int(halfWidth) + 1
	maxX += int(halfWidth) + 1
	minY -= int(halfWidth) + 1
	maxY += int(halfWidth) + 1

	return image.Rect(minX, minY, maxX, maxY)
}

// generateSamplePoints 生成曲线上的采样点
func (c *Curve) generateSamplePoints() []image.Point {
	var points []image.Point

	// 根据曲线类型和复杂度确定采样点数
	numSamples := 100 + c.LineWidth*2
	if c.CurveType == CubicBezier {
		numSamples *= 2
	}

	// 生成采样点
	for i := 0; i <= numSamples; i++ {
		t := float64(i) / float64(numSamples)
		var x, y float64

		switch c.CurveType {
		case QuadraticBezier:
			// 二次贝塞尔曲线公式
			x = math.Pow(1-t, 2)*float64(c.X0) + 2*(1-t)*t*float64(c.Cp1x) + math.Pow(t, 2)*float64(c.X1)
			y = math.Pow(1-t, 2)*float64(c.Y0) + 2*(1-t)*t*float64(c.Cp1y) + math.Pow(t, 2)*float64(c.Y1)

		case CubicBezier:
			// 三次贝塞尔曲线公式
			x = math.Pow(1-t, 3)*float64(c.X0) + 3*math.Pow(1-t, 2)*t*float64(c.Cp1x) +
				3*(1-t)*math.Pow(t, 2)*float64(c.Cp2x) + math.Pow(t, 3)*float64(c.X1)
			y = math.Pow(1-t, 3)*float64(c.Y0) + 3*math.Pow(1-t, 2)*t*float64(c.Cp1y) +
				3*(1-t)*math.Pow(t, 2)*float64(c.Cp2y) + math.Pow(t, 3)*float64(c.Y1)
		}

		points = append(points, image.Point{int(x), int(y)})
	}

	return points
}

// calculateMinDistanceSq 计算点到曲线的最小距离平方
func (c *Curve) calculateMinDistanceSq(px, py float64, samplePoints []image.Point) float64 {
	minDistSq := math.MaxFloat64

	// 计算点到每个线段的距离
	for i := 0; i < len(samplePoints)-1; i++ {
		p0 := samplePoints[i]
		p1 := samplePoints[i+1]

		// 线段的起点和终点
		x0, y0 := float64(p0.X), float64(p0.Y)
		x1, y1 := float64(p1.X), float64(p1.Y)

		// 计算点到线段的距离
		distSq := distancePointToLineSegmentSq(px, py, x0, y0, x1, y1)

		// 更新最小距离
		if distSq < minDistSq {
			minDistSq = distSq
		}
	}

	return minDistSq
}

// distancePointToLineSegmentSq 计算点到线段的距离平方
func distancePointToLineSegmentSq(px, py, x0, y0, x1, y1 float64) float64 {
	// 线段的向量
	dx := x1 - x0
	dy := y1 - y0

	// 点到起点的向量
	tx := px - x0
	ty := py - y0

	// 计算点在直线上的投影参数
	t := (tx*dx + ty*dy) / (dx*dx + dy*dy)

	// 限制t在0到1之间
	if t < 0 {
		t = 0
	} else if t > 1 {
		t = 1
	}

	// 计算投影点
	projX := x0 + t*dx
	projY := y0 + t*dy

	// 计算点到投影点的距离平方
	dx2 := px - projX
	dy2 := py - projY

	return dx2*dx2 + dy2*dy2
}

func setPixel(img *image.RGBA, x, y int, srcColor color.RGBA, alpha uint8) {
	// 检查坐标是否在图像范围内
	bounds := img.Bounds()
	if x < bounds.Min.X || x >= bounds.Max.X || y < bounds.Min.Y || y >= bounds.Max.Y {
		return
	}

	// 获取底图像素的RGBA值
	offset := img.PixOffset(x, y)
	dstR := img.Pix[offset]
	dstG := img.Pix[offset+1]
	dstB := img.Pix[offset+2]
	dstA := img.Pix[offset+3]

	// 标准Alpha混合系数计算
	srcAlpha := float64(alpha) / 255.0 // 源颜色的透明度占比
	dstAlpha := 1.0 - srcAlpha         // 底图颜色的透明度占比

	// 对R/G/B三个通道做标准混合，支持任意颜色
	mixR := uint8(float64(srcColor.R)*srcAlpha + float64(dstR)*dstAlpha)
	mixG := uint8(float64(srcColor.G)*srcAlpha + float64(dstG)*dstAlpha)
	mixB := uint8(float64(srcColor.B)*srcAlpha + float64(dstB)*dstAlpha)
	// 最终Alpha取最大值，保证底图透明度不被覆盖
	mixA := uint8(math.Max(float64(dstA), float64(alpha)))

	// 设置混合后的像素值
	img.Pix[offset] = mixR
	img.Pix[offset+1] = mixG
	img.Pix[offset+2] = mixB
	img.Pix[offset+3] = mixA
}
