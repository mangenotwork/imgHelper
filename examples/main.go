package main

import (
	"github.com/mangenotwork/imgHelper"
	"image/color"
	"image/png"
	"log"
	"os"
)

func main() {
	//case1()
	//case2()
	//case3()
	//case4()
	//case5()
	//case6()
	//case7()
	//case8()
	//case9()
	//case10()
	//case11()
	//case12()
	//case13()
	//case14()
	//case15()
	//case16()
	//case17()
	//case18()
	//case19()
	//case20()
	//case21()
	//case22()
	//case23()
	//case24()
	//case25()
	//case26()
	//case27()
	//case28()
	//case31()
	//case32()
	//case33()
	//case34()
	case35()
}

// 创建一个画布
// imgHelper.NewCanvas(400, 400)
func case1() {
	out := "./case1.png"
	cas := imgHelper.NewCanvas(400, 400)
	_ = cas.SaveToFile(out)
}

// 创建一个红色画布
// imgHelper.NewColorCanvas(400, 400, color.RGBA{R: 255, G: 0, B: 0, A: 255})
func case2() {
	out := "./case2.png"
	cas := imgHelper.NewColorCanvas(400, 400, color.RGBA{R: 255, G: 0, B: 0, A: 255})
	_ = cas.SaveToFile(out)
}

// 创建一个图片背景的画布
// imgHelper.NewImgCanvas(bkImg)
func case3() {
	out := "./case3.png"
	bkImg, err := imgHelper.OpenImgFromLocalFile("./test.png")
	if err != nil {
		log.Fatal(err)
	}
	err = imgHelper.NewImgCanvas(bkImg).SaveToFile(out)
	if err != nil {
		log.Fatal(err)
	}
}

// 创建一个图片背景的画布自定义宽高
// imgHelper.NewImgCanvasFromSize(1000, 200, fileBody)
func case4() {
	out := "./case4.png"
	bkImg, err := imgHelper.OpenImgFromLocalFile("./test.png")
	if err != nil {
		log.Fatal(err)
	}
	err = imgHelper.NewImgCanvasFromSize(800, 200, bkImg).SaveToFile(out)
	if err != nil {
		log.Fatal(err)
	}
}

// 创建一个图层，将图层输出为图片文件
func case5() {
	imgLayer, err := imgHelper.ImgLayerFromLocalFile("./test.png", imgHelper.Range{})
	if err != nil {
		log.Fatal(err)
	}
	imgLayer.Ext(imgHelper.OpsScale(100, 100)).Save("./case5.png")
}

// 从网络链接url获取一个图像并创建图层进行缩放保存到本地
func case6() {
	out := "case6.png"
	imgUrl := "https://resource.ecosmos.vip/ProductShare/158414f4c6b6b33a.png"
	src, err := imgHelper.OpenImgFromHttpGet(imgUrl)
	if err != nil {
		panic(err)
	}
	imgLayer := imgHelper.ImgLayer{
		Resource: src,
	}
	_ = imgLayer.Scale(200, 200)
	imgLayer.Save(out)
}

// 画图的简单示例 创建一个画布读取一个图片放在创建的图层上进行缩放，最终输出图片
func case7() {
	cas := imgHelper.NewCanvas(400, 400) // 创建400*400的画布

	testImg, err := imgHelper.OpenImgFromLocalFile("./test.png") // 本地读取test.png图片
	if err != nil {
		log.Fatal(err)
	}

	imgLayer := &imgHelper.ImgLayer{ // 创建一个图层,从右上角x0y0的位置开始绘制
		Resource: imgHelper.Scale(testImg, 1000, 1000),
		X0:       100,
		Y0:       100,
	}
	_ = imgLayer.Scale(200, 200) //将图层整体缩放到200*200

	err = cas.AddLayer(imgLayer).SaveToFile("./case7_3.png") // 图层添加到画布并存储到本地
	if err != nil {
		log.Fatal(err)
	}

}

// 两张图进行加法
func case8() {
	bkImg, err := imgHelper.OpenImgFromLocalFile("./test.png")
	if err != nil {
		log.Fatal(err)
	}
	cas := imgHelper.NewImgCanvas(bkImg)
	if cas.Err != nil {
		log.Fatal(cas.Err)
	}
	imgLayer, err := imgHelper.ImgLayerFromLocalFile("./case6.png", imgHelper.Range{
		X0: 100,
		Y0: 100,
	})
	if err != nil {
		log.Fatal(err)
	}
	err = cas.Addition(imgLayer).SaveToFile("./case8.png")
	if err != nil {
		log.Fatal(cas.Err)
	}
}

// 自定义图层并画在画布上
func case9() {
	// 打开图片作为图层
	imgLayer, err := imgHelper.ImgLayerFromLocalFile("./case6.png", imgHelper.Range{
		X0: 100,
		Y0: 100,
	})
	if err != nil {
		log.Fatal(err)
	}
	// 图层进行伸缩到 200*200
	_ = imgLayer.Scale(200, 200)
	// 绘制到画布上并输出图片
	err = imgHelper.NewCanvas(400, 400).Ext(imgLayer.Draw).SaveToFile("./case9.png")
	if err != nil {
		log.Fatal(err)
	}
}

// 添加一个旋转图层并旋转90度输出图片
func case10() {
	err := imgHelper.CanvasFromLocalImg("./test.png").Ext(imgHelper.OpsRotate(66)).SaveToFile("./case10_1.png")
	err = imgHelper.CanvasFromLocalImg("./test.png").Ext(imgHelper.OpsRotate90()).SaveToFile("./case10_2.png")
	err = imgHelper.CanvasFromLocalImg("./test.png").Ext(imgHelper.OpsRotate180()).SaveToFile("./case10_3.png")
	err = imgHelper.CanvasFromLocalImg("./test.png").Ext(imgHelper.OpsRotate270()).SaveToFile("./case10_4.png")
	if err != nil {
		log.Fatal(err)
	}
}

// 添加一个图片图层进行旋转在画布100*100的位置
// 使用 Ext 执行 OpsRotate 方法
func case11() {
	imgLayer, err := imgHelper.ImgLayerFromLocalFile("./case6.png", imgHelper.Range{
		X0: 100,
		Y0: 100,
	})
	if err != nil {
		log.Fatal(err)
	}
	// 图层旋转 66度
	imgLayer.Ext(imgHelper.OpsRotate(66))
	err = imgHelper.NewCanvas(800, 800).AddLayer(imgLayer).SaveToFile("./case11.png")
	if err != nil {
		log.Fatal(err)
	}
}

// 添加一个图片图层进行旋转在画布100*100的位置
// 使用 Rotate 函数
func case12() {
	imgLayer, err := imgHelper.ImgLayerFromLocalFile("./case6.png", imgHelper.Range{
		X0: 100,
		Y0: 100,
	})
	if err != nil {
		log.Fatal(err)
	}
	// 图层旋转 66度
	imgLayer.Resource = imgHelper.Rotate(imgLayer.Resource, 66)
	err = imgHelper.NewCanvas(800, 800).AddLayer(imgLayer).SaveToFile("./case12.png")
	if err != nil {
		log.Fatal(err)
	}
}

// 打开一张图片作为画布，然后打开第二张图缩放到100*200，旋转66度，绘制到画布的x0,y0=100,100的位置
func case13() {
	// 打开一张图片作为图层
	imgLayer, err := imgHelper.ImgLayerFromLocalFile("./case6.png", imgHelper.Range{X0: 100, Y0: 100})
	if err != nil {
		log.Fatal(err)
	}
	// 图层缩放和旋转操作
	imgLayer.Ext(imgHelper.OpsScale(100, 200)).Ext(imgHelper.OpsRotate(66))
	// 图层画绘制到画布上并保存到图片文件
	err = imgHelper.CanvasFromLocalImg("./test.png").AddLayer(imgLayer).Ext(imgHelper.OpsRotate90()).SaveToFile("./case13.png")
	if err != nil {
		log.Fatal(err)
	}
}

// 打开一张图片作为画布，然后打开第二张图缩放到100*200，旋转66度，绘制到画布的x0,y0=100,100的位置 - 另一种写法
func case14() {
	// 打开一张图片作为图层
	imgLayer, err := imgHelper.ImgLayerFromLocalFile("./case6.png", imgHelper.Range{X0: 100, Y0: 100})
	if err != nil {
		log.Fatal(err)
	}
	// 将图层的图片执行缩放
	imgHelper.Scale(imgLayer.Resource, 100, 200)
	// 将图层的图片执行旋转
	imgHelper.Rotate(imgLayer.Resource, 66)
	// 图层画绘制到画布上并保存到图片文件
	err = imgHelper.CanvasFromLocalImg("./test.png").AddLayer(imgLayer).Ext(imgHelper.OpsRotate90()).SaveToFile("./case14.png")
	if err != nil {
		log.Fatal(err)
	}
}

// 只想调用该库的方法
func case15() {
	file, err := os.Open("./test.png")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	src, err := png.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	// 调用图像缩放方法
	dst := imgHelper.Scale(src, 100, 100)
	outputFile, err := os.Create("./case15.png")
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()
	err = png.Encode(outputFile, dst)
	if err != nil {
		log.Fatal(err)
	}
}

// 灰度处理
func case16() {
	imgLayer, err := imgHelper.ImgLayerFromLocalFile("./test.png", imgHelper.Range{})
	if err != nil {
		log.Fatal(err)
	}
	_ = imgLayer.Ext(imgHelper.OpsGray()).Save("./case16.png")
}

// 亮度调整
func case17() {
	imgLayer, err := imgHelper.ImgLayerFromLocalFile("./test.png", imgHelper.Range{})
	if err != nil {
		log.Fatal(err)
	}
	_ = imgLayer.Ext(imgHelper.OpsBrightness(40)).Save("./case17_2.png")
}

// 两张图进行减法
func case18() {
	bkImg, err := imgHelper.OpenImgFromLocalFile("./test.png")
	if err != nil {
		log.Fatal(err)
	}
	cas := imgHelper.NewImgCanvas(bkImg)
	if cas.Err != nil {
		log.Fatal(cas.Err)
	}
	imgLayer, err := imgHelper.ImgLayerFromLocalFile("./case16.png", imgHelper.Range{})
	if err != nil {
		log.Fatal(err)
	}
	err = cas.Subtraction(imgLayer, false).SaveToFile("./case18_2.png")
	if err != nil {
		log.Fatal(cas.Err)
	}
}

// 两张图进行乘法
func case19() {
	bkImg, err := imgHelper.OpenImgFromLocalFile("./test.png")
	if err != nil {
		log.Fatal(err)
	}
	cas := imgHelper.NewImgCanvas(bkImg)
	if cas.Err != nil {
		log.Fatal(cas.Err)
	}
	imgLayer, err := imgHelper.ImgLayerFromLocalFile("./case9.png", imgHelper.Range{})
	if err != nil {
		log.Fatal(err)
	}
	err = cas.Multiplication(imgLayer, false).SaveToFile("./case19_1.png")
	if err != nil {
		log.Fatal(cas.Err)
	}
}

// 两张图进行除法
func case20() {
	bkImg, err := imgHelper.OpenImgFromLocalFile("./test.png")
	if err != nil {
		log.Fatal(err)
	}
	cas := imgHelper.NewImgCanvas(bkImg)
	if cas.Err != nil {
		log.Fatal(cas.Err)
	}
	imgLayer, err := imgHelper.ImgLayerFromLocalFile("./case9.png", imgHelper.Range{})
	if err != nil {
		log.Fatal(err)
	}
	err = cas.Division(imgLayer, false).SaveToFile("./case20_2.png")
	if err != nil {
		log.Fatal(cas.Err)
	}
}

// 两张图进行与运算
func case21() {
	bkImg, err := imgHelper.OpenImgFromLocalFile("./test.png")
	if err != nil {
		log.Fatal(err)
	}
	cas := imgHelper.NewImgCanvas(bkImg)
	if cas.Err != nil {
		log.Fatal(cas.Err)
	}
	imgLayer, err := imgHelper.ImgLayerFromLocalFile("./case9.png", imgHelper.Range{})
	if err != nil {
		log.Fatal(err)
	}
	err = cas.AND(imgLayer).SaveToFile("./case21.png")
	if err != nil {
		log.Fatal(cas.Err)
	}
}

// 两张图进行或运算
func case22() {
	bkImg, err := imgHelper.OpenImgFromLocalFile("./test.png")
	if err != nil {
		log.Fatal(err)
	}
	cas := imgHelper.NewImgCanvas(bkImg)
	if cas.Err != nil {
		log.Fatal(cas.Err)
	}
	imgLayer, err := imgHelper.ImgLayerFromLocalFile("./case9.png", imgHelper.Range{})
	if err != nil {
		log.Fatal(err)
	}
	err = cas.OR(imgLayer).SaveToFile("./case22.png")
	if err != nil {
		log.Fatal(cas.Err)
	}
}

// 两张图进行异或运算
func case23() {
	bkImg, err := imgHelper.OpenImgFromLocalFile("./test.png")
	if err != nil {
		log.Fatal(err)
	}
	cas := imgHelper.NewImgCanvas(bkImg)
	if cas.Err != nil {
		log.Fatal(cas.Err)
	}
	imgLayer, err := imgHelper.ImgLayerFromLocalFile("./case9.png", imgHelper.Range{})
	if err != nil {
		log.Fatal(err)
	}
	err = cas.XOR(imgLayer).SaveToFile("./case23.png")
	if err != nil {
		log.Fatal(cas.Err)
	}
}

// 两张图进行非运算
func case24() {
	bkImg, err := imgHelper.OpenImgFromLocalFile("./test.png")
	if err != nil {
		log.Fatal(err)
	}
	cas := imgHelper.NewImgCanvas(bkImg)
	if cas.Err != nil {
		log.Fatal(cas.Err)
	}
	imgLayer, err := imgHelper.ImgLayerFromLocalFile("./case9.png", imgHelper.Range{})
	if err != nil {
		log.Fatal(err)
	}
	err = cas.NOT(imgLayer).SaveToFile("./case23.png")
	if err != nil {
		log.Fatal(cas.Err)
	}
}

// 图层平移操作
func case25() {
	imgLayer, err := imgHelper.ImgLayerFromLocalFile("./test.png", imgHelper.Range{})
	if err != nil {
		log.Fatal(err)
	}
	_ = imgLayer.Translation(100, 50).Save("./case25.png")
}

// 裁剪操作
func case26() {
	imgLayer, err := imgHelper.ImgLayerFromLocalFile("./test.png", imgHelper.Range{})
	if err != nil {
		log.Fatal(err)
	}
	//_ = imgLayer.Ext(imgHelper.OpsCrop(imgHelper.Range{100, 50, 300, 300})).Save("./case26.png")
	//_ = imgLayer.Ext(imgHelper.OpsCrop(imgHelper.RangeCircle{200, 200, 100})).Save("./case26_2.png")
	//_ = imgLayer.Ext(imgHelper.OpsCrop(imgHelper.RangeTriangle{200, 10, 20, 200, 340, 200})).Save("./case26_3.png")

	rgPolygon := imgHelper.RangePolygon{
		Points: []imgHelper.Point{
			{X: 0, Y: 0},
			{X: 100, Y: 10},
			{X: 120, Y: 40},
			{X: 80, Y: 60},
			{X: 40, Y: 100},
			{X: 20, Y: 10},
		},
	}
	_ = imgLayer.Ext(imgHelper.OpsCrop(rgPolygon)).Save("./case26_4.png")
}

// OpenImgFromBytes 方法例子
func case27() {
	file, err := os.ReadFile("./test.png")
	if err != nil {
		log.Fatal(err)
	}
	src, err := imgHelper.OpenImgFromBytes(file)
	if err != nil {
		log.Fatal(err)
	}
	_ = imgHelper.NewImgLayer(src, imgHelper.Range{}).Save("./case27.png")
}

// NewImgCanvasFromRange 例子
func case28() {
	file, err := os.ReadFile("./test.png")
	if err != nil {
		log.Fatal(err)
	}
	src, err := imgHelper.OpenImgFromBytes(file)
	if err != nil {
		log.Fatal(err)
	}
	_ = imgHelper.NewImgCanvasFromRange(imgHelper.Range{0, 0, 100, 200}, src).SaveToFile("./case28.png")
}

// OpsScaleNearestNeighbor
func case29() {
	imgLayer, err := imgHelper.ImgLayerFromLocalFile("./test.png", imgHelper.Range{})
	if err != nil {
		log.Fatal(err)
	}
	imgLayer.Ext(imgHelper.OpsScaleNearestNeighbor(100, 100)).Save("./case29.png")
}

// OpsScaleCatmullRom
func case30() {
	imgLayer, err := imgHelper.ImgLayerFromLocalFile("./test.png", imgHelper.Range{})
	if err != nil {
		log.Fatal(err)
	}
	_ = imgLayer.Ext(imgHelper.OpsScaleCatmullRom(100, 100)).Save("./case30.png")
}

// OpsTransposition 图像转置操作
func case31() {
	imgLayer, err := imgHelper.ImgLayerFromLocalFile("./test.png", imgHelper.Range{})
	if err != nil {
		log.Fatal(err)
	}
	_ = imgLayer.Ext(imgHelper.OpsTransposition()).Save("./case31.png")
}

// OpsMirrorHorizontal 镜像
func case32() {
	imgLayer, err := imgHelper.ImgLayerFromLocalFile("./test.png", imgHelper.Range{})
	if err != nil {
		log.Fatal(err)
	}
	_ = imgLayer.Ext(imgHelper.OpsMirrorHorizontal()).Save("./case32.png")
}

// OpsMirrorVertical  镜像
func case33() {
	imgLayer, err := imgHelper.ImgLayerFromLocalFile("./test.png", imgHelper.Range{})
	if err != nil {
		log.Fatal(err)
	}
	_ = imgLayer.Ext(imgHelper.OpsMirrorVertical()).Save("./case33.png")
}

// OpsBinaryImg 二值图操作
func case34() {
	imgLayer, err := imgHelper.ImgLayerFromLocalFile("./test.png", imgHelper.Range{})
	if err != nil {
		log.Fatal(err)
	}
	_ = imgLayer.Ext(imgHelper.OpsBinaryImg()).Save("./case34.png")
}

// OpsMosaic 马赛克操作
func case35() {
	imgLayer, err := imgHelper.ImgLayerFromLocalFile("./case6.png", imgHelper.Range{})
	if err != nil {
		log.Fatal(err)
	}
	//_ = imgLayer.Ext(imgHelper.OpsMosaic(imgHelper.Range{100, 100, 300, 300}, 20)).Save("./case35.png")
	//_ = imgLayer.Ext(imgHelper.OpsMosaic(imgHelper.RangeCircle{100, 100, 50}, 20)).Save("./case35_2.png")
	//_ = imgLayer.Ext(imgHelper.OpsMosaic(imgHelper.RangeTriangle{90, 50, 50, 200, 150, 200}, 20)).Save("./case35_3.png")
	rgPolygon := imgHelper.RangePolygon{
		Points: []imgHelper.Point{
			{X: 0, Y: 0},
			{X: 100, Y: 60},
			{X: 120, Y: 90},
			{X: 80, Y: 110},
			{X: 40, Y: 150},
			{X: 20, Y: 60},
		},
	}
	_ = imgLayer.Ext(imgHelper.OpsMosaic(rgPolygon, 20)).Save("./case35_4.png")
}
