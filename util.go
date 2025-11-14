package imgHelper

import "image/color"

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
func isPointInTriangle(pX, pY, x1, y1, x2, y2, x3, y3 int) bool {
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
