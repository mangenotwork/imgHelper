package imgHelper

import (
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/draw"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

var fontObj *opentype.Font
var fontDataOnce sync.Once

// GetFontDefault 获取默认字体，全局的
func GetFontDefault() *opentype.Font {

	fontDataOnce.Do(func() {

		_, file, _, _ := runtime.Caller(0)
		absFile, _ := filepath.Abs(file)
		dir := filepath.Dir(absFile)
		log.Println("dir = ", dir)
		var err error
		fontFile, err := os.Open(dir + "/NotoSansSC-Regular.ttf")
		if err != nil {
			panic(err)
		}
		defer func() {
			_ = fontFile.Close()
		}()
		fontData, err := io.ReadAll(fontFile)
		if err != nil {
			panic(err)
		}
		fontObj, err = opentype.Parse(fontData)
		if err != nil {
			panic(err)
		}
	})
	return fontObj
}

// SetFont 自定义设置字体
func SetFont(fontPath string) (*opentype.Font, error) {
	if fontPath == "" {
		return nil, os.ErrInvalid
	}
	fontFile, err := os.Open(fontPath)
	if err != nil {
		return nil, err
	}
	defer func() { _ = fontFile.Close() }()
	fontData, err := io.ReadAll(fontFile)
	if err != nil {
		return nil, err
	}
	return opentype.Parse(fontData)
}

// TextLayer 图层 - 文字，在画布上绘制文字
type TextLayer struct {
	X0, Y0            int
	X1, Y1            int
	Str               string // 绘制的内容
	Size              float64
	DPI               float64
	Colour            color.Color
	MaxWidth          int            // 字的最大宽度
	FontGradient      bool           // 字体渐变
	FontGradientColor []color.Color  // 字体渐变颜色
	Font              *opentype.Font // 字体
	Align             Align          // 对齐方式
	// todo 字体阴影
	// todo 字体模糊（类似毛玻璃效果）
	// todo 斜体
	// todo 垂直绘制
}

// NewTextLayer 新建字体图层，默认字体，从左往右
func NewTextLayer(str string, size float64, X0, Y0 int, colour color.Color) *TextLayer {
	return &TextLayer{
		X0:           X0,
		Y0:           Y0,
		Str:          str,
		Size:         size,
		DPI:          FontFixed,
		Colour:       colour,
		FontGradient: false,
		Font:         GetFontDefault(),
		Align:        Left,
	}
}

// SetDPI 设置DPI
func (textLayer *TextLayer) SetDPI(dpi float64) *TextLayer {
	textLayer.DPI = dpi
	return textLayer
}

// SetFont 指定字体对象 只支持 *opentype.Font
func (textLayer *TextLayer) SetFont(font *opentype.Font) *TextLayer {
	textLayer.Font = font
	return textLayer
}

// SetFontFile 指定字体文件
func (textLayer *TextLayer) SetFontFile(fontPath string) (*TextLayer, error) {
	var err error
	textLayer.Font, err = SetFont(fontPath)
	return textLayer, err
}

// SetMaxWidth 设置字体的最大绘制宽度，超过部分已"..."代替
func (textLayer *TextLayer) SetMaxWidth(w int) *TextLayer {
	textLayer.MaxWidth = w
	return textLayer
}

// SetAlign 设置字体绘制方式，目前只支持水平方向
func (textLayer *TextLayer) SetAlign(a Align) *TextLayer {
	textLayer.Align = a
	return textLayer
}

// SetGradient 设置字体渐变
func (textLayer *TextLayer) SetGradient(cs []color.Color) *TextLayer {
	if len(cs) > 1 && cs[0] != cs[1] {
		textLayer.FontGradient = true
		textLayer.FontGradientColor = cs
	}
	return textLayer
}

type Align string

const (
	Left      Align = "left"
	Right     Align = "right"
	Center    Align = "center"
	FontFixed       = 64
)

func (textLayer *TextLayer) Draw(ctx *CanvasContext) error {
	ctxWidth := ctx.Dst.Bounds().Dx()
	face, err := opentype.NewFace(fontObj, &opentype.FaceOptions{
		Size:    textLayer.Size,
		DPI:     textLayer.DPI,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return err
	}

	if textLayer.MaxWidth > 0 {
		textLayer.Str = textLayer.textMaxWidth(textLayer.Str, textLayer.Size, textLayer.DPI, textLayer.Font, textLayer.MaxWidth)
	}

	textWidth := font.MeasureString(face, textLayer.Str).Ceil()

	// todo bug 右和居中没效果
	switch textLayer.Align {
	case Left: // 默认从左往右绘制文字
	case Right:
		textLayer.X0 = ctxWidth - (ctxWidth - (textLayer.X0 + textLayer.MaxWidth)) - textWidth
	case Center:
		textLayer.X0 = textLayer.X0 + ((textLayer.MaxWidth - textWidth) / 2) // 按坐标居中
	}

	textLayer.X1 = textLayer.X0 + textWidth
	xDot := fixed.Int26_6(textLayer.X0 * FontFixed)
	yDot := fixed.Int26_6((textLayer.Y0 + int(textLayer.Size)) * FontFixed)

	if textLayer.FontGradient && len(textLayer.FontGradientColor) == 2 && textLayer.FontGradientColor[0] != textLayer.FontGradientColor[1] {
		// 渐变色绘制
		strDrawer := textLayer.gradientDrawer(
			ctx.Dst,
			face,
			textLayer.FontGradientColor[0],
			textLayer.FontGradientColor[1],
		)

		currentX := textLayer.X0
		runes := []rune(textLayer.Str)

		for i, r := range runes {
			strDrawer.Drawer.Dot = fixed.Point26_6{
				X: fixed.Int26_6(currentX * FontFixed),
				Y: yDot,
			}
			strDrawer.DrawString(string(r), strDrawer.Drawer.Dot)
			charWidth := font.MeasureString(face, string(r)).Ceil()
			currentX += charWidth
			if i == len(runes)-1 {
				currentX += 1
			}
		}

	} else if !textLayer.FontGradient && textLayer.Colour != nil {

		cardNameDrawer := &font.Drawer{
			Dst:  ctx.Dst,
			Src:  image.NewUniform(textLayer.Colour),
			Face: face,
			Dot: fixed.Point26_6{
				X: xDot,
				Y: yDot,
			},
		}

		// 单字绘画 实现字间距
		currentX := textLayer.X0
		runes := []rune(textLayer.Str)

		for i, r := range runes {
			cardNameDrawer.Dot = fixed.Point26_6{
				X: fixed.Int26_6(currentX * FontFixed),
				Y: yDot,
			}
			cardNameDrawer.DrawString(string(r))
			charWidth := font.MeasureString(face, string(r)).Ceil()
			currentX += charWidth
			if i == len(runes)-1 {
				currentX += 1
			}
		}

	}
	return nil
}

func (textLayer *TextLayer) textMaxWidth(text string, fontSize, DPI float64, fontObj *opentype.Font, mixWidth int) string {
	textObj, err := opentype.NewFace(fontObj, &opentype.FaceOptions{
		Size:    fontSize,
		DPI:     DPI,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Println(err)
		return ""
	}

	textWidth := font.MeasureString(textObj, text).Ceil()
	if textWidth <= mixWidth {
		return text
	}

	ellipsis := "..."
	ellipsisWidth := font.MeasureString(textObj, ellipsis).Ceil()
	targetWidth := mixWidth - ellipsisWidth

	// 如果连省略号都放不下，直接绘制省略号
	if mixWidth <= ellipsisWidth {
		return ellipsis
	}

	start := 0
	texts := []rune(text)
	end := len(texts)
	fitText := ""
	counter := 0

	for start <= end {
		counter++
		mid := (start + end) / 2
		testText := string(texts[:mid])
		testWidth := font.MeasureString(textObj, testText).Ceil()
		if testWidth <= targetWidth {
			fitText = testText
			// 尝试更长文本
			start = mid + 1
		} else {
			// 尝试更短文本
			end = mid - 1
		}
	}

	return fitText + ellipsis
}

// gradientDrawer 结构体用于绘制渐变文字
type gradientDrawer struct {
	Drawer      *font.Drawer
	TopColor    color.Color
	BottomColor color.Color
}

// NewGradientDrawer 创建一个新的渐变文字绘制器
func (textLayer *TextLayer) gradientDrawer(dst draw.Image, face font.Face, topColor, bottomColor color.Color) *gradientDrawer {
	d := &font.Drawer{
		Dst:  dst,
		Src:  image.NewUniform(color.White), // 临时使用白色
		Face: face,
	}
	return &gradientDrawer{
		Drawer:      d,
		TopColor:    topColor,
		BottomColor: bottomColor,
	}
}

// DrawString 绘制带有上下渐变颜色的文字
func (gd *gradientDrawer) DrawString(s string, pt fixed.Point26_6) fixed.Point26_6 {
	// 获取文本的边界框
	bounds, _ := font.BoundString(gd.Drawer.Face, s)
	width := (bounds.Max.X - bounds.Min.X).Ceil()
	height := (bounds.Max.Y - bounds.Min.Y).Ceil()

	// 创建临时图像用于绘制文本
	temp := image.NewRGBA(image.Rect(0, 0, width, height))
	tempDrawer := *gd.Drawer
	tempDrawer.Dst = temp
	tempDrawer.Dot = fixed.Point26_6{X: 0, Y: -bounds.Min.Y}
	tempDrawer.DrawString(s)

	// 创建渐变图像
	gradient := image.NewRGBA(image.Rect(0, 0, width, height))

	// 将颜色转换为RGBA格式
	r1, g1, b1, a1 := gd.TopColor.RGBA()
	r2, g2, b2, a2 := gd.BottomColor.RGBA()

	// 计算每个通道的插值基数
	dr := float64(r2>>8) - float64(r1>>8)
	dg := float64(g2>>8) - float64(g1>>8)
	db := float64(b2>>8) - float64(b1>>8)
	da := float64(a2>>8) - float64(a1>>8)

	for y := 0; y < height; y++ {
		// 计算渐变因子 (0.0 到 1.0)
		factor := float64(y) / float64(height)

		// 线性插值计算当前行的颜色
		r := uint8(float64(r1>>8) + dr*factor)
		g := uint8(float64(g1>>8) + dg*factor)
		b := uint8(float64(b1>>8) + db*factor)
		a := uint8(float64(a1>>8) + da*factor)

		for x := 0; x < width; x++ {
			// 如果该像素在文本内 (非透明)，则应用渐变颜色
			if temp.RGBAAt(x, y).A > 0 {
				gradient.Set(x, y, color.RGBA{R: r, G: g, B: b, A: a})
			}
		}
	}

	// 将渐变文本绘制到目标图像
	draw.DrawMask(gd.Drawer.Dst,
		image.Rect(pt.X.Ceil(), pt.Y.Ceil()-height, pt.X.Ceil()+width, pt.Y.Ceil()),
		gradient, image.Point{0, 0},
		temp, image.Point{0, 0},
		draw.Over,
	)

	// 返回下一个绘制点
	return pt.Add(fixed.Point26_6{X: bounds.Max.X - bounds.Min.X, Y: 0})
}

func (textLayer *TextLayer) GetResource() image.Image {
	return nil
}

func (textLayer *TextLayer) Save(filePath string) error {
	return nil
}

func (textLayer *TextLayer) GetXY() (int, int) {
	return textLayer.X0, textLayer.Y0
}
