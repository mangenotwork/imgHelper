package imgHelper

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"os"
)

// 几何绘制图层

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
	return nil
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
