package imgHelper

import (
	"image"
	"image/color"
	"math"
)

// SignedNumeric 定义类型约束所有数值
type SignedNumeric interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

func abs[T SignedNumeric](x T) T {
	if x < 0 {
		return -x
	}
	return x
}

func clamp[T SignedNumeric](v, min, max T) T {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

// 传入n个参数返回最大的值
func maxValue[T SignedNumeric](n ...T) T {
	if len(n) == 0 {
		return 0
	}
	maxN := n[0]
	for _, v := range n[1:] {
		if v > maxN {
			maxN = v
		}
	}
	return maxN
}

func minValue[T SignedNumeric](n ...T) T {
	if len(n) == 0 {
		return 0
	}
	minN := n[0]
	for _, v := range n[1:] {
		if v < minN {
			minN = v
		}
	}
	return minN
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

// 辅助函数：判断点(pX,pY)是否在三角形(x1,y1)-(x2,y2)-(x3,y3)内部（含边界）
// 原理：通过向量叉乘判断点与三条边的位置关系，若在同一侧则在内部
func isPointInTriangle[T SignedNumeric](pX, pY, x1, y1, x2, y2, x3, y3 T) bool {
	// 计算三个叉积（判断点在边的哪一侧）
	cross1 := (x2-x1)*(pY-y1) - (y2-y1)*(pX-x1)
	cross2 := (x3-x2)*(pY-y2) - (y3-y2)*(pX-x2)
	cross3 := (x1-x3)*(pY-y3) - (y1-y3)*(pX-x3)

	// 判断三个叉积的符号是否一致（均非正或均非负，含零）
	hasPositive := (cross1 > 0) || (cross2 > 0) || (cross3 > 0)
	hasNegative := (cross1 < 0) || (cross2 < 0) || (cross3 < 0)

	// 若不同时存在正负，则点在三角形内（含边界）
	return !(hasPositive && hasNegative)
}

// isPointInPolygon 用射线法判断点(pX,pY)是否在多边形内部（含边界）
// 原理：从点向右发射水平射线，统计与多边形边的交点数量，奇数则在内部，偶数则在外部
func isPointInPolygon(pX, pY int, vertices [][2]int) bool {
	n := len(vertices)
	if n < 3 {
		return false // 非多边形
	}
	inside := false

	for i := 0; i < n; i++ {
		j := (i + 1) % n // 下一个顶点索引（闭合多边形）
		vi := vertices[i]
		vj := vertices[j]
		xi, yi := vi[0], vi[1]
		xj, yj := vj[0], vj[1]

		// 检查点是否在当前边上（边界情况）
		if isPointOnLine(pX, pY, xi, yi, xj, yj) {
			return true
		}

		// 射线与边相交判断（仅处理边跨越射线y坐标的情况）
		if (yi > pY) != (yj > pY) {
			// 计算交点的x坐标（射线是水平向右的，y=pY）
			xIntersect := ((pY-yi)*(xj-xi))/(yj-yi) + xi
			// 若交点在点的右侧，则计数+1
			if pX < xIntersect {
				inside = !inside
			}
		}
	}

	return inside
}

// isPointOnLine 判断点(pX,pY)是否在直线段(x1,y1)-(x2,y2)上（边界处理）
func isPointOnLine(pX, pY, x1, y1, x2, y2 int) bool {
	// 点在直线的 bounding box 内
	if (pX < minValue(x1, x2) || pX > minValue(x1, x2)) ||
		(pY < minValue(y1, y2) || pY > minValue(y1, y2)) {
		return false
	}

	// 向量叉积为0（点在直线上）
	cross := (x2-x1)*(pY-y1) - (y2-y1)*(pX-x1)
	if cross != 0 {
		return false
	}

	return true
}

// RGBToHSV 将 RGB 颜色转换为 HSV 颜色
func RGBToHSV(r, g, b uint8) (float64, float64, float64) {
	rNorm := float64(r) / 255.0
	gNorm := float64(g) / 255.0
	bNorm := float64(b) / 255.0
	maxVal := math.Max(rNorm, math.Max(gNorm, bNorm))
	minVal := math.Min(rNorm, math.Min(gNorm, bNorm))
	delta := maxVal - minVal
	var h, s, v float64
	v = maxVal
	if delta == 0 {
		h = 0
	} else {
		s = delta / maxVal
		if maxVal == rNorm {
			h = math.Mod((gNorm-bNorm)/delta, 6)
		} else if maxVal == gNorm {
			h = (bNorm-rNorm)/delta + 2
		} else {
			h = (rNorm-gNorm)/delta + 4
		}
		h *= 60
		if h < 0 {
			h += 360
		}
	}
	return h, s, v
}

// HSVToRGB 将 HSV 颜色转换为 RGB 颜色
func HSVToRGB(h, s, v float64) (uint8, uint8, uint8) {
	c := v * s
	hPrime := h / 60
	x := c * (1 - math.Abs(math.Mod(hPrime, 2)-1))
	var r1, g1, b1 float64
	switch {
	case 0 <= hPrime && hPrime < 1:
		r1 = c
		g1 = x
		b1 = 0
	case 1 <= hPrime && hPrime < 2:
		r1 = x
		g1 = c
		b1 = 0
	case 2 <= hPrime && hPrime < 3:
		r1 = 0
		g1 = c
		b1 = x
	case 3 <= hPrime && hPrime < 4:
		r1 = 0
		g1 = x
		b1 = c
	case 4 <= hPrime && hPrime < 5:
		r1 = x
		g1 = 0
		b1 = c
	case 5 <= hPrime && hPrime < 6:
		r1 = c
		g1 = 0
		b1 = x
	}
	m := v - c
	r := uint8((r1 + m) * 255)
	g := uint8((g1 + m) * 255)
	b := uint8((b1 + m) * 255)
	return r, g, b
}

// 生成一维高斯核
func generateGaussianKernel(sigma float64) []float64 {
	size := int(math.Ceil(sigma * 3))
	kernel := make([]float64, 2*size+1)
	sum := 0.0
	for i := -size; i <= size; i++ {
		kernel[i+size] = math.Exp(-float64(i*i) / (2 * sigma * sigma))
		sum += kernel[i+size]
	}
	for i := range kernel {
		kernel[i] /= sum
	}
	return kernel
}

// 针对黑色文字+浅色背景，强制将文字转为前景（255），背景转为0
func binaryImgForText(img image.Image) [][]uint8 {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	bin := make([][]uint8, height)
	for y := 0; y < height; y++ {
		bin[y] = make([]uint8, width)
		for x := 0; x < width; x++ {
			// 转为灰度后，低于阈值的文字像素设为255，否则为0
			r, g, b, _ := img.At(x, y).RGBA()
			gray := uint8(0.299*float64(r>>8) + 0.587*float64(g>>8) + 0.114*float64(b>>8))
			if gray <= 50 {
				bin[y][x] = 255 // 文字（前景）
			} else {
				bin[y][x] = 0 // 背景
			}
		}
	}
	return bin
}

// 8邻域索引（顺时针：p2-p9，对应坐标偏移）
var neighbors = []image.Point{
	{0, -1},  // p2: (x, y-1)
	{1, -1},  // p3: (x+1, y-1)
	{1, 0},   // p4: (x+1, y)
	{1, 1},   // p5: (x+1, y+1)
	{0, 1},   // p6: (x, y+1)
	{-1, 1},  // p7: (x-1, y+1)
	{-1, 0},  // p8: (x-1, y)
	{-1, -1}, // p9: (x-1, y-1)
}

// countConnections 计算8邻域的连接数（衡量像素的"拐角"程度）
func countConnections(bin [][]uint8, x, y int) int {
	height, width := len(bin), len(bin[0])
	// 取p2-p9的二值（1=前景，0=背景）
	p := make([]int, 8)
	for i, n := range neighbors {
		nx, ny := x+n.X, y+n.Y
		if nx >= 0 && nx < width && ny >= 0 && ny < height && bin[ny][nx] == 255 {
			p[i] = 1
		}
	}
	// 连接数=相邻像素从0→1的次数（p9与p2相邻）
	count := 0
	for i := 0; i < 8; i++ {
		j := (i + 1) % 8
		if p[i] == 0 && p[j] == 1 {
			count++
		}
	}
	return count
}

// countForeground 计算8邻域中前景像素（255）的数量
func countForeground(bin [][]uint8, x, y int) int {
	count := 0
	height, width := len(bin), len(bin[0])
	for _, n := range neighbors {
		nx, ny := x+n.X, y+n.Y
		// 边界外视为背景
		if nx >= 0 && nx < width && ny >= 0 && ny < height && bin[ny][nx] == 255 {
			count++
		}
	}
	return count
}

// getNeighbor 获取指定邻域像素的二值（1=前景，0=背景）
func getNeighbor(bin [][]uint8, x, y, idx int) int {
	n := neighbors[idx]
	nx, ny := x+n.X, y+n.Y
	height, width := len(bin), len(bin[0])
	if nx >= 0 && nx < width && ny >= 0 && ny < height && bin[ny][nx] == 255 {
		return 1
	}
	return 0
}

// 克隆图像
func cloneImage(img image.Image) *image.RGBA {
	bounds := img.Bounds()
	clone := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			clone.Set(x, y, img.At(x, y))
		}
	}
	return clone
}

// 检查坐标是否在边界内
func inBounds(bounds image.Rectangle, x, y int) bool {
	return x >= bounds.Min.X && x < bounds.Max.X &&
		y >= bounds.Min.Y && y < bounds.Max.Y
}

// normalize 向量归一化
func normalize(x, y float64) (nx, ny float64) {
	mh := math.Hypot(x, y)
	if mh < 1e-6 {
		return 0, 0
	}
	return x / mh, y / mh
}

// image.Image to *image.NRGBA
func imageToNRGBA(src image.Image) *image.NRGBA {
	srcBounds := src.Bounds()
	dstBounds := srcBounds.Sub(srcBounds.Min)

	dst := image.NewNRGBA(dstBounds)

	dstMinX := dstBounds.Min.X
	dstMinY := dstBounds.Min.Y

	srcMinX := srcBounds.Min.X
	srcMinY := srcBounds.Min.Y
	srcMaxX := srcBounds.Max.X
	srcMaxY := srcBounds.Max.Y

	switch src0 := src.(type) {

	case *image.NRGBA:
		rowSize := srcBounds.Dx() * 4
		numRows := srcBounds.Dy()

		i0 := dst.PixOffset(dstMinX, dstMinY)
		j0 := src0.PixOffset(srcMinX, srcMinY)

		di := dst.Stride
		dj := src0.Stride

		for row := 0; row < numRows; row++ {
			copy(dst.Pix[i0:i0+rowSize], src0.Pix[j0:j0+rowSize])
			i0 += di
			j0 += dj
		}

	case *image.NRGBA64:
		i0 := dst.PixOffset(dstMinX, dstMinY)
		for y := srcMinY; y < srcMaxY; y, i0 = y+1, i0+dst.Stride {
			for x, i := srcMinX, i0; x < srcMaxX; x, i = x+1, i+4 {

				j := src0.PixOffset(x, y)

				dst.Pix[i+0] = src0.Pix[j+0]
				dst.Pix[i+1] = src0.Pix[j+2]
				dst.Pix[i+2] = src0.Pix[j+4]
				dst.Pix[i+3] = src0.Pix[j+6]

			}
		}

	case *image.RGBA:
		i0 := dst.PixOffset(dstMinX, dstMinY)
		for y := srcMinY; y < srcMaxY; y, i0 = y+1, i0+dst.Stride {
			for x, i := srcMinX, i0; x < srcMaxX; x, i = x+1, i+4 {

				j := src0.PixOffset(x, y)
				a := src0.Pix[j+3]
				dst.Pix[i+3] = a

				switch a {
				case 0:
					dst.Pix[i+0] = 0
					dst.Pix[i+1] = 0
					dst.Pix[i+2] = 0
				case 0xff:
					dst.Pix[i+0] = src0.Pix[j+0]
					dst.Pix[i+1] = src0.Pix[j+1]
					dst.Pix[i+2] = src0.Pix[j+2]
				default:
					dst.Pix[i+0] = uint8(uint16(src0.Pix[j+0]) * 0xff / uint16(a))
					dst.Pix[i+1] = uint8(uint16(src0.Pix[j+1]) * 0xff / uint16(a))
					dst.Pix[i+2] = uint8(uint16(src0.Pix[j+2]) * 0xff / uint16(a))
				}
			}
		}

	case *image.RGBA64:
		i0 := dst.PixOffset(dstMinX, dstMinY)
		for y := srcMinY; y < srcMaxY; y, i0 = y+1, i0+dst.Stride {
			for x, i := srcMinX, i0; x < srcMaxX; x, i = x+1, i+4 {

				j := src0.PixOffset(x, y)
				a := src0.Pix[j+6]
				dst.Pix[i+3] = a

				switch a {
				case 0:
					dst.Pix[i+0] = 0
					dst.Pix[i+1] = 0
					dst.Pix[i+2] = 0
				case 0xff:
					dst.Pix[i+0] = src0.Pix[j+0]
					dst.Pix[i+1] = src0.Pix[j+2]
					dst.Pix[i+2] = src0.Pix[j+4]
				default:
					dst.Pix[i+0] = uint8(uint16(src0.Pix[j+0]) * 0xff / uint16(a))
					dst.Pix[i+1] = uint8(uint16(src0.Pix[j+2]) * 0xff / uint16(a))
					dst.Pix[i+2] = uint8(uint16(src0.Pix[j+4]) * 0xff / uint16(a))
				}
			}
		}

	case *image.Gray:
		i0 := dst.PixOffset(dstMinX, dstMinY)
		for y := srcMinY; y < srcMaxY; y, i0 = y+1, i0+dst.Stride {
			for x, i := srcMinX, i0; x < srcMaxX; x, i = x+1, i+4 {

				j := src0.PixOffset(x, y)
				c := src0.Pix[j]
				dst.Pix[i+0] = c
				dst.Pix[i+1] = c
				dst.Pix[i+2] = c
				dst.Pix[i+3] = 0xff

			}
		}

	case *image.Gray16:
		i0 := dst.PixOffset(dstMinX, dstMinY)
		for y := srcMinY; y < srcMaxY; y, i0 = y+1, i0+dst.Stride {
			for x, i := srcMinX, i0; x < srcMaxX; x, i = x+1, i+4 {

				j := src0.PixOffset(x, y)
				c := src0.Pix[j]
				dst.Pix[i+0] = c
				dst.Pix[i+1] = c
				dst.Pix[i+2] = c
				dst.Pix[i+3] = 0xff

			}
		}

	case *image.YCbCr:
		i0 := dst.PixOffset(dstMinX, dstMinY)
		for y := srcMinY; y < srcMaxY; y, i0 = y+1, i0+dst.Stride {
			for x, i := srcMinX, i0; x < srcMaxX; x, i = x+1, i+4 {

				yj := src0.YOffset(x, y)
				cj := src0.COffset(x, y)
				r, g, b := color.YCbCrToRGB(src0.Y[yj], src0.Cb[cj], src0.Cr[cj])

				dst.Pix[i+0] = r
				dst.Pix[i+1] = g
				dst.Pix[i+2] = b
				dst.Pix[i+3] = 0xff

			}
		}

	default:
		i0 := dst.PixOffset(dstMinX, dstMinY)
		for y := srcMinY; y < srcMaxY; y, i0 = y+1, i0+dst.Stride {
			for x, i := srcMinX, i0; x < srcMaxX; x, i = x+1, i+4 {

				c := color.NRGBAModel.Convert(src.At(x, y)).(color.NRGBA)

				dst.Pix[i+0] = c.R
				dst.Pix[i+1] = c.G
				dst.Pix[i+2] = c.B
				dst.Pix[i+3] = c.A

			}
		}
	}

	return dst
}
