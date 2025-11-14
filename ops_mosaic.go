package imgHelper

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
)

// OpsMosaic 马赛克操作
// 参数:
// - rg 马赛克范围
// - blockSize 马赛克块大小
func OpsMosaic(rg RangeValue, blockSize int) func(ctx *CanvasContext) error {
	return func(ctx *CanvasContext) error {

		switch rg.Type() {

		case RangeRectType:
			x0, y0, x1, y1 := rg.(Range).Value()
			ctx.Dst = Mosaic(ctx.Dst, x0, y0, x1, y1, blockSize).(*image.RGBA)

		case RangeCircleType:
			rgObj := rg.(RangeCircle)
			ctx.Dst = MosaicCircle(ctx.Dst, rgObj.Cx, rgObj.Cy, rgObj.R, blockSize).(*image.RGBA)

		case RangeTriangleType:
			rgObj := rg.(RangeTriangle)
			ctx.Dst = MosaicTriangle(ctx.Dst, rgObj.X0, rgObj.Y0, rgObj.X1, rgObj.Y1, rgObj.X2, rgObj.Y2, blockSize).(*image.RGBA)

		case RangePolygonType:
			rgObj := rg.(RangePolygon)
			points := make([]int, 0)
			for _, v := range rgObj.Points {
				points = append(points, v.X)
				points = append(points, v.Y)
			}
			dst, err := MosaicPolygon(ctx.Dst, blockSize, points...)
			ctx.Dst = dst.(*image.RGBA)
			ctx.Err = errors.Join(ctx.Err, err)

		}

		return nil
	}
}

// Mosaic 马赛克
// 参数:
// - x0, y0, x1, y1  马赛克范围
// - blockSize 马赛克块大小
func Mosaic(src image.Image, x0, y0, x1, y1 int, blockSize int) image.Image {
	bounds := src.Bounds()
	drawImg := image.NewRGBA(bounds)
	draw.Draw(drawImg, bounds, src, bounds.Min, draw.Src)
	x, y := x0, y0
	width, height := x1-x0, y1-y0
	for i := y; i < y+height; i += blockSize {
		for j := x; j < x+width; j += blockSize {
			var r, g, b, a uint32
			count := 0
			for m := 0; m < blockSize; m++ {
				for n := 0; n < blockSize; n++ {
					if i+m < y+height && j+n < x+width {
						pr, pg, pb, pa := drawImg.At(j+n, i+m).RGBA()
						r += pr
						g += pg
						b += pb
						a += pa
						count++
					}
				}
			}
			if count > 0 {
				r /= uint32(count)
				g /= uint32(count)
				b /= uint32(count)
				a /= uint32(count)
				c := color.RGBA64{R: uint16(r), G: uint16(g), B: uint16(b), A: uint16(a)}
				for m := 0; m < blockSize; m++ {
					for n := 0; n < blockSize; n++ {
						if i+m < y+height && j+n < x+width {
							drawImg.Set(j+n, i+m, c)
						}
					}
				}
			}
		}
	}
	return drawImg
}

// MosaicCircle 圆形范围马赛克
// 参数:
// - cx, cy：圆心坐标
// - r：圆的半径
// - blockSize：马赛克块大小（块越大，模糊效果越强）
func MosaicCircle(src image.Image, cx, cy, r int, blockSize int) image.Image {
	bounds := src.Bounds()
	drawImg := image.NewRGBA(bounds)
	draw.Draw(drawImg, bounds, src, bounds.Min, draw.Src)

	minX := cx - r
	maxX := cx + r
	minY := cy - r
	maxY := cy + r
	if minX < bounds.Min.X {
		minX = bounds.Min.X
	}
	if maxX > bounds.Max.X {
		maxX = bounds.Max.X
	}
	if minY < bounds.Min.Y {
		minY = bounds.Min.Y
	}
	if maxY > bounds.Max.Y {
		maxY = bounds.Max.Y
	}

	for y := minY; y < maxY; y += blockSize {
		for x := minX; x < maxX; x += blockSize {
			var totalR, totalG, totalB, totalA uint32
			pixelCount := 0
			for m := 0; m < blockSize; m++ {
				for n := 0; n < blockSize; n++ {
					pixelX := x + n
					pixelY := y + m
					if pixelX >= bounds.Min.X && pixelX < bounds.Max.X &&
						pixelY >= bounds.Min.Y && pixelY < bounds.Max.Y {
						dx := pixelX - cx
						dy := pixelY - cy
						if dx*dx+dy*dy <= r*r {
							r, g, b, a := drawImg.At(pixelX, pixelY).RGBA()
							totalR += r
							totalG += g
							totalB += b
							totalA += a
							pixelCount++
						}
					}
				}
			}

			if pixelCount > 0 {
				avgR := totalR / uint32(pixelCount)
				avgG := totalG / uint32(pixelCount)
				avgB := totalB / uint32(pixelCount)
				avgA := totalA / uint32(pixelCount)
				avgColor := color.RGBA64{
					R: uint16(avgR),
					G: uint16(avgG),
					B: uint16(avgB),
					A: uint16(avgA),
				}

				for m := 0; m < blockSize; m++ {
					for n := 0; n < blockSize; n++ {
						pixelX := x + n
						pixelY := y + m
						if pixelX >= bounds.Min.X && pixelX < bounds.Max.X &&
							pixelY >= bounds.Min.Y && pixelY < bounds.Max.Y {
							dx := pixelX - cx
							dy := pixelY - cy
							if dx*dx+dy*dy <= r*r {
								drawImg.Set(pixelX, pixelY, avgColor)
							}
						}
					}
				}

			}
		}
	}

	return drawImg
}

// MosaicTriangle 三角形范围马赛克
// 参数:
// - x1,y1, x2,y2, x3,y3：三角形三个顶点坐标
// - blockSize：马赛克块大小（块越大，颗粒感越强）
func MosaicTriangle(src image.Image, x1, y1, x2, y2, x3, y3 int, blockSize int) image.Image {
	bounds := src.Bounds()
	drawImg := image.NewRGBA(bounds)
	draw.Draw(drawImg, bounds, src, bounds.Min, draw.Src)
	minX := minValue(x1, x2, x3)
	maxX := maxValue(x1, x2, x3)
	minY := minValue(y1, y2, y3)
	maxY := maxValue(y1, y2, y3)
	if minX < bounds.Min.X {
		minX = bounds.Min.X
	}
	if maxX > bounds.Max.X {
		maxX = bounds.Max.X
	}
	if minY < bounds.Min.Y {
		minY = bounds.Min.Y
	}
	if maxY > bounds.Max.Y {
		maxY = bounds.Max.Y
	}
	for y := minY; y < maxY; y += blockSize {
		for x := minX; x < maxX; x += blockSize {
			var totalR, totalG, totalB, totalA uint32
			pixelCount := 0
			for m := 0; m < blockSize; m++ {
				for n := 0; n < blockSize; n++ {
					pixelX := x + n
					pixelY := y + m
					if pixelX >= bounds.Min.X && pixelX < bounds.Max.X &&
						pixelY >= bounds.Min.Y && pixelY < bounds.Max.Y &&
						isPointInTriangle(pixelX, pixelY, x1, y1, x2, y2, x3, y3) {
						r, g, b, a := drawImg.At(pixelX, pixelY).RGBA()
						totalR += r
						totalG += g
						totalB += b
						totalA += a
						pixelCount++
					}
				}
			}
			if pixelCount > 0 {
				avgR := totalR / uint32(pixelCount)
				avgG := totalG / uint32(pixelCount)
				avgB := totalB / uint32(pixelCount)
				avgA := totalA / uint32(pixelCount)
				avgColor := color.RGBA64{
					R: uint16(avgR),
					G: uint16(avgG),
					B: uint16(avgB),
					A: uint16(avgA),
				}
				for m := 0; m < blockSize; m++ {
					for n := 0; n < blockSize; n++ {
						pixelX := x + n
						pixelY := y + m
						if pixelX >= bounds.Min.X && pixelX < bounds.Max.X &&
							pixelY >= bounds.Min.Y && pixelY < bounds.Max.Y &&
							isPointInTriangle(pixelX, pixelY, x1, y1, x2, y2, x3, y3) {
							drawImg.Set(pixelX, pixelY, avgColor)
						}
					}
				}
			}
		}
	}
	return drawImg
}

// MosaicPolygon 多边形范围马赛克
// 参数:
//
//	src：源图像
//	blockSize：马赛克块大小（块越大，颗粒感越强）
//	points：多边形顶点坐标，格式为 [x0,y0, x1,y1, ..., xn,yn]（至少3个顶点，长度≥6）
func MosaicPolygon(src image.Image, blockSize int, points ...int) (image.Image, error) {
	if len(points) < 6 || len(points)%2 != 0 || blockSize <= 0 {
		bounds := src.Bounds()
		dst := image.NewRGBA(bounds)
		draw.Draw(dst, bounds, src, bounds.Min, draw.Src)
		return dst, fmt.Errorf("至少需要三个顶点")
	}

	vertexCount := len(points) / 2
	minX, maxX := points[0], points[0]
	minY, maxY := points[1], points[1]
	vertices := make([][2]int, vertexCount)
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

	bounds := src.Bounds()
	srcMinX, srcMinY := bounds.Min.X, bounds.Min.Y
	srcMaxX, srcMaxY := bounds.Max.X, bounds.Max.Y
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

	drawImg := image.NewRGBA(bounds)
	draw.Draw(drawImg, bounds, src, bounds.Min, draw.Src)

	for y := minY; y < maxY; y += blockSize {
		for x := minX; x < maxX; x += blockSize {

			var totalR, totalG, totalB, totalA uint32
			var pixelCount = 0

			for m := 0; m < blockSize; m++ {
				for n := 0; n < blockSize; n++ {
					pixelX := x + n
					pixelY := y + m
					if pixelX >= srcMinX && pixelX < srcMaxX &&
						pixelY >= srcMinY && pixelY < srcMaxY &&
						isPointInPolygon(pixelX, pixelY, vertices) {
						r, g, b, a := drawImg.At(pixelX, pixelY).RGBA()
						totalR += r
						totalG += g
						totalB += b
						totalA += a
						pixelCount++
					}
				}
			}

			if pixelCount > 0 {
				avgR := totalR / uint32(pixelCount)
				avgG := totalG / uint32(pixelCount)
				avgB := totalB / uint32(pixelCount)
				avgA := totalA / uint32(pixelCount)
				avgColor := color.RGBA64{
					R: uint16(avgR),
					G: uint16(avgG),
					B: uint16(avgB),
					A: uint16(avgA),
				}

				for m := 0; m < blockSize; m++ {
					for n := 0; n < blockSize; n++ {
						pixelX := x + n
						pixelY := y + m
						if pixelX >= srcMinX && pixelX < srcMaxX &&
							pixelY >= srcMinY && pixelY < srcMaxY &&
							isPointInPolygon(pixelX, pixelY, vertices) {
							drawImg.Set(pixelX, pixelY, avgColor)
						}
					}
				}

			}

		}
	}

	return drawImg, nil
}
