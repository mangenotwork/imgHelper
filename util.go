package imgHelper

import "image/color"

// SignedNumeric 定义类型约束：支持所有有符号整数和浮点数（这些类型有正负，需要计算绝对值）
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
