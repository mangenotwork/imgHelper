package imgHelper

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
)

// Crop 裁剪
func Crop(src image.Image, x0, y0, x1, y1 int) image.Image {
	cropRect := image.Rect(x0, y0, x1, y1)
	dst := image.NewRGBA(cropRect)
	draw.Draw(dst, dst.Bounds(), src, cropRect.Min, draw.Src)
	return dst
}

// OpsCrop 裁剪操作
// todo 三角形，多边形的锯齿太严重了
// todo 需要增加曲线多边形裁剪
func OpsCrop(rg RangeValue) func(ctx *CanvasContext) error {
	return func(ctx *CanvasContext) error {

		switch rg.Type() {
		case RangeRectType:
			rgObj := rg.(Range)
			ctx.Dst = Crop(ctx.Dst, rgObj.X0, rgObj.Y0, rgObj.X1, rgObj.Y1).(*image.RGBA)
		case RangeCircleType:
			rgObj := rg.(RangeCircle)
			ctx.Dst = CropCircle(ctx.Dst, rgObj.Cx, rgObj.Cy, rgObj.R).(*image.RGBA)
		case RangeTriangleType:
			rgObj := rg.(RangeTriangle)
			ctx.Dst = CropTriangle(ctx.Dst, rgObj.X0, rgObj.Y0, rgObj.X1, rgObj.Y1, rgObj.X2, rgObj.Y2).(*image.RGBA)
		case RangePolygonType:
			rgObj := rg.(RangePolygon)
			points := make([]int, 0)
			for _, v := range rgObj.Points {
				points = append(points, v.X)
				points = append(points, v.Y)
			}
			dst, err := CropPolygon(ctx.Dst, points...)
			ctx.Dst = dst.(*image.RGBA)
			ctx.Err = errors.Join(ctx.Err, err)
		}

		return nil
	}
}

// CropCircle 圆形裁剪：保留源图像中以(cx, cy)为圆心、r为半径的圆形区域，圆形外像素设为透明
// 参数：
//
//	cx, cy：圆心在源图像中的坐标
//	r：圆的半径
func CropCircle(src image.Image, cx, cy, r int) image.Image {
	// 边界矩形左上角：(cx - r, cy - r)，右下角：(cx + r, cy + r)
	boundMinX := cx - r
	boundMinY := cy - r
	boundMaxX := cx + r
	boundMaxY := cy + r

	// 越界处理
	srcBounds := src.Bounds()
	if boundMinX < srcBounds.Min.X {
		boundMinX = srcBounds.Min.X
	}
	if boundMinY < srcBounds.Min.Y {
		boundMinY = srcBounds.Min.Y
	}
	if boundMaxX > srcBounds.Max.X {
		boundMaxX = srcBounds.Max.X
	}
	if boundMaxY > srcBounds.Max.Y {
		boundMaxY = srcBounds.Max.Y
	}

	cropRect := image.Rect(boundMinX, boundMinY, boundMaxX, boundMaxY)
	dst := image.NewRGBA(cropRect)

	rSquared := r * r // 半径的平方（用于距离判断，避免开方运算）
	for y := boundMinY; y < boundMaxY; y++ {
		for x := boundMinX; x < boundMaxX; x++ {
			// 计算当前像素到圆心的距离的平方
			dx := x - cx
			dy := y - cy
			distanceSquared := dx*dx + dy*dy

			// 若距离 <= 半径，属于圆形内，保留源图像像素；否则设为透明
			if distanceSquared <= rSquared {

				r, g, b, a := src.At(x, y).RGBA()
				dst.SetRGBA(
					x, y,
					color.RGBA{
						R: uint8(r >> 8),
						G: uint8(g >> 8),
						B: uint8(b >> 8),
						A: uint8(a >> 8),
					},
				)
			} else {
				// 圆形外像素设为透明（Alpha=0）
				dst.SetRGBA(x, y, color.RGBA{R: 0, G: 0, B: 0, A: 0})
			}
		}
	}

	return dst
}

// CropTriangle 三角形裁剪：保留源图像中由三个顶点(x1,y1)、(x2,y2)、(x3,y3)围成的三角形区域，外部像素设为透明
// 参数：
//
//	src：源图像
//	x1,y1, x2,y2, x3,y3：三角形三个顶点在源图像中的坐标
func CropTriangle(src image.Image, x1, y1, x2, y2, x3, y3 int) image.Image {
	minX := minValue(x1, x2, x3)
	maxX := maxValue(x1, x2, x3)
	minY := minValue(y1, y2, y3)
	maxY := maxValue(y1, y2, y3)

	srcBounds := src.Bounds()
	srcMinX, srcMinY := srcBounds.Min.X, srcBounds.Min.Y
	srcMaxX, srcMaxY := srcBounds.Max.X, srcBounds.Max.Y

	if minX < srcMinX {
		minX = srcMinX
	}
	if maxX > srcMaxX {
		maxX = srcMaxX
	}
	if minY < srcMinY {
		minY = srcMinY
	}
	if maxY > srcMaxY {
		maxY = srcMaxY
	}

	// 目标图像的Bounds为 [minX, minY] 到 [maxX, maxY]，内部坐标以 (0,0) 对应 (minX, minY)
	cropRect := image.Rect(minX, minY, maxX, maxY)
	dst := image.NewRGBA(cropRect)
	dstWidth := maxX - minX  // 目标图像宽度
	dstHeight := maxY - minY // 目标图像高度

	// 遍历边界矩形内的每个像素，判断是否在三角形内
	for y := minY; y < maxY; y++ {
		for x := minX; x < maxX; x++ {
			// 计算当前像素在目标图像中的相对坐标（关键修复：转换为dst内部坐标）
			dx := x - minX
			dy := y - minY

			// 双重校验：确保相对坐标在目标图像范围内（避免极端情况越界）
			if dx < 0 || dx >= dstWidth || dy < 0 || dy >= dstHeight {
				continue
			}

			// 判断点是否在三角形内
			if isPointInTriangle(x, y, x1, y1, x2, y2, x3, y3) {
				// 三角形内：保留源图像像素（转换为8位通道）
				r, g, b, a := src.At(x, y).RGBA()
				dst.SetRGBA(
					dx, dy, // 使用相对坐标写入目标图像
					color.RGBA{
						R: uint8(r >> 8),
						G: uint8(g >> 8),
						B: uint8(b >> 8),
						A: uint8(a >> 8),
					},
				)
			} else {
				// 三角形外：设为透明（Alpha=0）
				dst.SetRGBA(dx, dy, color.RGBA{R: 0, G: 0, B: 0, A: 0})
			}
		}
	}

	return dst
}

// CropPolygon 多边形裁剪：保留源图像中由多个顶点围成的多边形区域，外部像素设为透明
// 参数：
//
//	src：源图像
//	points：多边形顶点坐标，格式为 [x0,y0, x1,y1, ..., xn,yn]（至少需3个顶点，即长度≥6）
//
// 返回：裁剪后的图像（失败时返回nil）
func CropPolygon(src image.Image, points ...int) (image.Image, error) {
	if len(points) < 6 || len(points)%2 != 0 {
		return nil, fmt.Errorf("至少需要三个顶点")
	}
	vertexCount := len(points) / 2

	minX, maxX := points[0], points[0]
	minY, maxY := points[1], points[1]
	vertices := make([][2]int, vertexCount) // 存储顶点坐标 (x,y)

	for i := 0; i < vertexCount; i++ {
		x := points[2*i]
		y := points[2*i+1]
		vertices[i] = [2]int{x, y}

		if x < minX {
			minX = x
		}
		if x > maxX {
			maxX = x
		}
		if y < minY {
			minY = y
		}
		if y > maxY {
			maxY = y
		}
	}

	// 边界矩形越界处理（限制在源图像范围内）
	srcBounds := src.Bounds()
	srcMinX, srcMinY := srcBounds.Min.X, srcBounds.Min.Y
	srcMaxX, srcMaxY := srcBounds.Max.X, srcBounds.Max.Y

	if minX < srcMinX {
		minX = srcMinX
	}
	if maxX > srcMaxX {
		maxX = srcMaxX
	}
	if minY < srcMinY {
		minY = srcMinY
	}
	if maxY > srcMaxY {
		maxY = srcMaxY
	}

	// 创建目标图像（大小为边界矩形的宽高）
	cropRect := image.Rect(minX, minY, maxX, maxY)
	dst := image.NewRGBA(cropRect)
	dstWidth := maxX - minX
	dstHeight := maxY - minY

	for y := minY; y < maxY; y++ {
		for x := minX; x < maxX; x++ {
			dx := x - minX
			dy := y - minY

			// 校验相对坐标有效性
			if dx < 0 || dx >= dstWidth || dy < 0 || dy >= dstHeight {
				continue
			}

			if isPointInPolygon(x, y, vertices) {
				// 多边形内：保留源图像像素
				r, g, b, a := src.At(x, y).RGBA()
				dst.SetRGBA(
					dx, dy,
					color.RGBA{
						R: uint8(r >> 8),
						G: uint8(g >> 8),
						B: uint8(b >> 8),
						A: uint8(a >> 8),
					},
				)
			} else {
				// 多边形外：设为透明
				dst.SetRGBA(dx, dy, color.RGBA{R: 0, G: 0, B: 0, A: 0})
			}
		}
	}

	return dst, nil
}
