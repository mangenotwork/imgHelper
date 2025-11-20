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

	if gLayer.resource == nil {
		// todo 寻找最大的坐标范围创建透明背景
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

func (gLayer *GeometryLayer) AddShape(s Shape) {
	gLayer.shapes = append(gLayer.shapes, s)
}

type Shape interface {
	Render(src image.Image) image.Image
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
