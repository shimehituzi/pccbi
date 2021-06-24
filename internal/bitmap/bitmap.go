package bitmap

import (
	"image"
	"image/color"
)

// Fyne の canvas で描画するときに必要なメソッド
type FyneBitMap interface {
	// フレームの枚数，縦幅，横幅
	GetLength() DimensionLength
	// ラベリングされた個数を返すメソッド
	GetLabelLength(int) int
	// 一枚画像を取り出すメソッド
	GetImage(int, int) image.Image
}

// フレームの枚数，縦幅，横幅
type DimensionLength struct {
	// フレームの枚数
	D0 int
	// 縦幅
	D1 int
	// 横幅
	D2 int
}

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
		return color.Gray{0}
	}
	if bm[y][x] == 1 {
		return color.Gray{0}
	} else {
		return color.Gray{255}
	}
}
