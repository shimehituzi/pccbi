package refactoring

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

// FyneBitMap の InterFace の実装
func (bc *bitCube) GetImage(f int, _ int) image.Image {
	return bc.Data[f]
}

// FyneBitMap の InterFace の実装
func (bc *bitCube) GetLabelLength(int) int {
	return 0
}

// FyneBitMap の InterFace の実装
func (bc *bitCube) GetLength() DimensionLength {
	return DimensionLength{
		D0: bc.Length[0],
		D1: bc.Length[1],
		D2: bc.Length[2],
	}
}

// imgae.Image の InterFace を実装
func (bm bitmap) ColorModel() color.Model {
	return color.RGBAModel
}

// imgae.Image の InterFace を実装
func (bm bitmap) Bounds() image.Rectangle {
	return image.Rect(0, 0, len(bm[0]), len(bm))
}

// imgae.Image の InterFace を実装
func (bm bitmap) At(x, y int) color.Color {
	rect := image.Rect(0, 0, len(bm[0]), len(bm))
	if !(image.Point{x, y}.In(rect)) {
		return color.RGBA{0, 0, 0, 0}
	}
	if bm[y][x] == 1 {
		return color.RGBA{255, 255, 255, 255}
	} else {
		return color.RGBA{0, 0, 0, 0}
	}
}
