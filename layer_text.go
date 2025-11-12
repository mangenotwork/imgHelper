package imgHelper

import "image/color"

// TextLayer 图层 - 文字，在画布上绘制文字
type TextLayer struct {
	Str               string // 绘制的内容
	Size              float64
	DPI               float64
	Colour            color.Color
	MaxWidth          int           // 字的最大宽度
	FontGradient      bool          // 字体渐变
	FontGradientColor []color.Color // 字体渐变颜色
	FontAlign         string        // 字体
}

func (textLayer *TextLayer) Draw(ctx *CanvasContext) error {
	return nil
}
