package imgHelper

import (
	"image"
	"image/color"
	"math"
)

// 仿射变换核心实现
func affineTransform(img image.Image, mat [6]float64) *image.RGBA {
	bounds := img.Bounds()
	dest := image.NewRGBA(bounds)
	a, b, c := mat[0], mat[1], mat[2]
	d, e, f := mat[3], mat[4], mat[5]
	det := a*e - b*d
	if det == 0 {
		return cloneImage(img)
	}
	invDet := 1.0 / det
	aPrime := e * invDet
	bPrime := -b * invDet
	cPrime := (-e*c + b*f) * invDet
	dPrime := -d * invDet
	ePrime := a * invDet
	fPrime := (d*c - a*f) * invDet
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			srcX := aPrime*float64(x) + bPrime*float64(y) + cPrime
			srcY := dPrime*float64(x) + ePrime*float64(y) + fPrime
			srcXInt := int(math.Round(srcX))
			srcYInt := int(math.Round(srcY))
			if inBounds(bounds, srcXInt, srcYInt) {
				dest.Set(x, y, img.At(srcXInt, srcYInt))
			} else {
				dest.Set(x, y, color.RGBA{R: 0, G: 0, B: 0, A: 255})
			}
		}
	}
	return dest
}

// RigidTransform 刚性变换（旋转、缩放、平移）
func RigidTransform(img image.Image, angle, scale, tx, ty float64) *image.RGBA {
	radian := angle * math.Pi / 180
	cos := math.Cos(radian)
	sin := math.Sin(radian)
	// 构建仿射变换矩阵参数
	mat := [6]float64{
		scale * cos, -scale * sin, tx,
		scale * sin, scale * cos, ty,
	}
	return affineTransform(img, mat)
}

// OpsRigidTransform 刚性变换（旋转、缩放、平移）
func OpsRigidTransform(angle, scale, tx, ty float64) func(ctx *CanvasContext) error {
	return func(ctx *CanvasContext) error {
		ctx.Dst = RigidTransform(ctx.Dst, angle, scale, tx, ty)
		return nil
	}
}

// AffineTransform 仿射变换
func AffineTransform(img image.Image, mat [6]float64) *image.RGBA {
	return affineTransform(img, mat)
}

// OpsAffineTransform 仿射变换
func OpsAffineTransform(mat [6]float64) func(ctx *CanvasContext) error {
	return func(ctx *CanvasContext) error {
		ctx.Dst = AffineTransform(ctx.Dst, mat)
		return nil
	}
}

// PerspectiveTransform 透视变换
func PerspectiveTransform(img image.Image, mat [9]float64) *image.RGBA {
	a, b, c := mat[0], mat[1], mat[2]
	d, e, f := mat[3], mat[4], mat[5]
	g, h, i := mat[6], mat[7], mat[8]
	det := a*(e*i-f*h) - b*(d*i-f*g) + c*(d*h-e*g)
	if det == 0 {
		return cloneImage(img)
	}
	invDet := 1.0 / det
	h11 := (e*i - f*h) * invDet
	h12 := (c*h - b*i) * invDet
	h13 := (b*f - c*e) * invDet
	h21 := (f*g - d*i) * invDet
	h22 := (a*i - c*g) * invDet
	h23 := (c*d - a*f) * invDet
	h31 := (d*h - e*g) * invDet
	h32 := (b*g - a*h) * invDet
	h33 := (a*e - b*d) * invDet
	bounds := img.Bounds()
	dest := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			xH := float64(x)
			yH := float64(y)
			wH := 1.0
			X := h11*xH + h12*yH + h13*wH
			Y := h21*xH + h22*yH + h23*wH
			W := h31*xH + h32*yH + h33*wH
			if W == 0 {
				dest.Set(x, y, color.RGBA{R: 0, G: 0, B: 0, A: 255})
				continue
			}
			srcX := X / W
			srcY := Y / W
			srcXInt := int(math.Round(srcX))
			srcYInt := int(math.Round(srcY))
			if inBounds(img.Bounds(), srcXInt, srcYInt) {
				dest.Set(x, y, img.At(srcXInt, srcYInt))
			} else {
				dest.Set(x, y, color.RGBA{R: 0, G: 0, B: 0, A: 255})
			}
		}
	}
	return dest
}

// OpsPerspectiveTransform 透视变换
func OpsPerspectiveTransform(mat [9]float64) func(ctx *CanvasContext) error {
	return func(ctx *CanvasContext) error {
		ctx.Dst = PerspectiveTransform(ctx.Dst, mat)
		return nil
	}
}

// 双线性插值用于根据浮点坐标计算像素值。
func bInterpolation(img *image.RGBA, x, y float64) color.RGBA {
	x1, y1 := int(x), int(y)
	x2, y2 := x1+1, y1+1
	if x2 >= img.Bounds().Max.X {
		x2 = img.Bounds().Max.X - 1
	}
	if y2 >= img.Bounds().Max.Y {
		y2 = img.Bounds().Max.Y - 1
	}

	q11 := img.RGBAAt(x1, y1)
	q12 := img.RGBAAt(x1, y2)
	q21 := img.RGBAAt(x2, y1)
	q22 := img.RGBAAt(x2, y2)

	xFrac, yFrac := x-float64(x1), y-float64(y1)

	r := uint8((1-xFrac)*(1-yFrac)*float64(q11.R) +
		xFrac*(1-yFrac)*float64(q21.R) +
		(1-xFrac)*yFrac*float64(q12.R) +
		xFrac*yFrac*float64(q22.R))
	g := uint8((1-xFrac)*(1-yFrac)*float64(q11.G) +
		xFrac*(1-yFrac)*float64(q21.G) +
		(1-xFrac)*yFrac*float64(q12.G) +
		xFrac*yFrac*float64(q22.G))
	b := uint8((1-xFrac)*(1-yFrac)*float64(q11.B) +
		xFrac*(1-yFrac)*float64(q21.B) +
		(1-xFrac)*yFrac*float64(q12.B) +
		xFrac*yFrac*float64(q22.B))

	return color.RGBA{R: r, G: g, B: b, A: 255}
}

// AffineTransform23 仿射变换通过 2x3 矩阵实现
func AffineTransform23(img *image.RGBA, matrix [2][3]float64) *image.RGBA {
	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// 计算逆变换
			det := matrix[0][0]*matrix[1][1] - matrix[0][1]*matrix[1][0]
			if det == 0 {
				continue
			}
			invM := [2][3]float64{
				{matrix[1][1] / det, -matrix[0][1] / det, (matrix[0][1]*matrix[1][2] - matrix[1][1]*matrix[0][2]) / det},
				{-matrix[1][0] / det, matrix[0][0] / det, (matrix[1][0]*matrix[0][2] - matrix[0][0]*matrix[1][2]) / det},
			}

			srcX := invM[0][0]*float64(x) + invM[0][1]*float64(y) + invM[0][2]
			srcY := invM[1][0]*float64(x) + invM[1][1]*float64(y) + invM[1][2]

			if srcX >= 0 && srcX < float64(bounds.Dx()) &&
				srcY >= 0 && srcY < float64(bounds.Dy()) {
				newImg.Set(x, y, bInterpolation(img, srcX, srcY))
			} else {
				newImg.Set(x, y, color.RGBA{R: 0, G: 0, B: 0, A: 255}) // 背景黑色
			}
		}
	}
	return newImg
}

// OpsAffineTransform23 仿射变换通过 2x3 矩阵实现
func OpsAffineTransform23(matrix [2][3]float64) func(ctx *CanvasContext) error {
	return func(ctx *CanvasContext) error {
		ctx.Dst = AffineTransform23(ctx.Dst, matrix)
		return nil
	}
}

// PerspectiveTransform33 透视变换通过 3x3 矩阵实现
func PerspectiveTransform33(img *image.RGBA, matrix [3][3]float64) *image.RGBA {
	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// 计算逆变换
			det := matrix[0][0]*(matrix[1][1]*matrix[2][2]-matrix[1][2]*matrix[2][1]) -
				matrix[0][1]*(matrix[1][0]*matrix[2][2]-matrix[1][2]*matrix[2][0]) +
				matrix[0][2]*(matrix[1][0]*matrix[2][1]-matrix[1][1]*matrix[2][0])
			if det == 0 {
				continue
			}
			invM := [3][3]float64{
				{(matrix[1][1]*matrix[2][2] - matrix[1][2]*matrix[2][1]) / det,
					(matrix[0][2]*matrix[2][1] - matrix[0][1]*matrix[2][2]) / det,
					(matrix[0][1]*matrix[1][2] - matrix[0][2]*matrix[1][1]) / det},
				{(matrix[1][2]*matrix[2][0] - matrix[1][0]*matrix[2][2]) / det,
					(matrix[0][0]*matrix[2][2] - matrix[0][2]*matrix[2][0]) / det,
					(matrix[0][2]*matrix[1][0] - matrix[0][0]*matrix[1][2]) / det},
				{(matrix[1][0]*matrix[2][1] - matrix[1][1]*matrix[2][0]) / det,
					(matrix[0][1]*matrix[2][0] - matrix[0][0]*matrix[2][1]) / det,
					(matrix[0][0]*matrix[1][1] - matrix[0][1]*matrix[1][0]) / det},
			}

			w := invM[2][0]*float64(x) + invM[2][1]*float64(y) + invM[2][2]
			srcX := (invM[0][0]*float64(x) + invM[0][1]*float64(y) + invM[0][2]) / w
			srcY := (invM[1][0]*float64(x) + invM[1][1]*float64(y) + invM[1][2]) / w

			if srcX >= 0 && srcX < float64(bounds.Dx()) &&
				srcY >= 0 && srcY < float64(bounds.Dy()) {
				newImg.Set(x, y, bInterpolation(img, srcX, srcY))
			} else {
				newImg.Set(x, y, color.RGBA{R: 0, G: 0, B: 0, A: 255}) // 背景黑色
			}
		}
	}
	return newImg
}

// OpsPerspectiveTransform33 透视变换通过 3x3 矩阵实现
func OpsPerspectiveTransform33(matrix [3][3]float64) func(ctx *CanvasContext) error {
	return func(ctx *CanvasContext) error {
		ctx.Dst = PerspectiveTransform33(ctx.Dst, matrix)
		return nil
	}
}
