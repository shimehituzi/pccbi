package bitmap

import (
	"image"
	"image/color"
)

// 2値画像を表す構造体
type BitMap [][]byte

// imgae.Image の InterFace を実装
func (bm BitMap) ColorModel() color.Model {
	return color.GrayModel
}

// imgae.Image の InterFace を実装
func (bm BitMap) Bounds() image.Rectangle {
	return image.Rect(0, 0, len(bm[0]), len(bm))
}

// imgae.Image の InterFace を実装
func (bm BitMap) At(x, y int) color.Color {
	rect := image.Rect(0, 0, len(bm[0]), len(bm))
	if !(image.Point{x, y}.In(rect)) {
		return color.Gray{}
	}
	if bm[y][x] == 1 {
		return color.Gray{0}
	} else {
		return color.Gray{255}
	}
}
