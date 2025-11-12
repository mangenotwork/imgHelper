package imgHelper

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
)

// OpenImgFromLocalFile 从本地文件读取图像
func OpenImgFromLocalFile(imgPath string) (image.Image, error) {
	imgFile, err := os.Open(imgPath)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = imgFile.Close()
	}()
	imgObj, err := png.Decode(imgFile)
	if err != nil {
		data, err := os.ReadFile(imgPath)
		data = findSOI(data)
		if data == nil {
			fmt.Println("未找到 JPEG 起始标记")
		}
		imgObj, err = jpeg.Decode(bytes.NewReader(data))
		if err != nil {
			return nil, err
		}
	}
	return imgObj, nil
}

func findSOI(data []byte) []byte {
	soi := []byte{0xFF, 0xD8}
	index := bytes.Index(data, soi)
	if index == -1 {
		return nil
	}
	return data[index:]
}

// OpenImgFromReader 从Reader读取图像
func OpenImgFromReader(rd io.Reader) (image.Image, error) {
	data, err := io.ReadAll(rd)
	if err != nil {
		return nil, err // 读取流失败（如网络中断、文件损坏）
	}
	imgObj, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		data = findSOI(data)
		if data == nil {
			fmt.Println("未找到 JPEG 起始标记")
		}
		imgObj, err = jpeg.Decode(bytes.NewReader(data))
		if err != nil {
			return nil, err
		}
	}
	return imgObj, nil
}

// OpenImgFromBytes 从Bytes读取图像
func OpenImgFromBytes(data []byte) (image.Image, error) {
	imgObj, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		data = findSOI(data)
		if data == nil {
			log.Println("未找到 JPEG 起始标记")
		}
		imgObj, err = jpeg.Decode(bytes.NewReader(data))
		if err != nil {
			return nil, err
		}
	}
	return imgObj, nil
}

// OpenImgFromHttpGet http get请求下载url图像
func OpenImgFromHttpGet(imgUrl string) (image.Image, error) {
	resp, err := http.Get(imgUrl)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		log.Printf("[imgHelper Warn] http get url = %s resp status is %s", imgUrl, resp.Status)
	}
	return OpenImgFromReader(resp.Body)
}
