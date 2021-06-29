package refactoring

import (
	"image"
	"image/color"
)

// Fyne の canvas で描画するときに必要なメソッド
type FyneBitMap interface {
	// フレームの枚数，縦幅，横幅
	GetLength() DimensionLength
	// 一枚画像を取り出すメソッド
	GetImage(int) image.Image
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
func (lpc *labeledPointCloud) GetImage(f int) image.Image {
	return lpc.frames[f]
}

// FyneBitMap の InterFace の実装
func (lpc *labeledPointCloud) GetLength() DimensionLength {
	return DimensionLength{
		D0: lpc.length[0],
		D1: lpc.length[1],
		D2: lpc.length[2],
	}
}

// imgae.Image の InterFace を実装
func (f frame) ColorModel() color.Model {
	return color.RGBAModel
}

// imgae.Image の InterFace を実装
func (f frame) Bounds() image.Rectangle {
	return image.Rect(0, 0, len(f.img[0]), len(f.img))
}

// imgae.Image の InterFace を実装
func (f frame) At(x, y int) color.Color {
	rect := image.Rect(0, 0, len(f.img[0]), len(f.img))
	if !(image.Point{x, y}.In(rect)) {
		return color.RGBA{0, 0, 0, 0}
	}
	switch f.img[y][x] {
	case 1:
		if f.getOuterCounter(x, y) != -1 {
			return color.RGBA{255, 255, 255, 255}
		} else {
			return color.RGBA{128, 128, 128, 128}
		}
	case 2:
		label := f.getInnerCounter(x, y)
		switch label % 6 {
		case 1:
			return color.RGBA{255, 0, 0, 255}
		case 2:
			return color.RGBA{0, 255, 0, 255}
		case 3:
			return color.RGBA{0, 0, 255, 255}
		case 4:
			return color.RGBA{255, 0, 255, 255}
		case 5:
			return color.RGBA{0, 255, 255, 255}
		case 0:
			return color.RGBA{255, 255, 0, 255}
		default:
			return color.RGBA{0, 0, 0, 255}
		}
	default:
		return color.RGBA{0, 0, 0, 0}
	}
}

// ChainCode のラベルを返す
func (f frame) getOuterCounter(x, y int) int {
	for _, contour := range f.contours {
		for _, point := range contour.outer.points {
			if point.x == x && point.y == y {
				return contour.label
			}
		}
	}
	return -1
}

// ChainCode のラベルを返す
func (f frame) getInnerCounter(x, y int) int {
	for _, contour := range f.contours {
		for _, inner := range contour.inner {
			for _, point := range inner.points {
				if point.x == x && point.y == y {
					return contour.label
				}
			}
		}
	}
	return -1
}

// FyneBitMap の InterFace の実装
func (bc *bitCube) GetImage(f int) image.Image {
	return bc.data[f]
}

// FyneBitMap の InterFace の実装
func (bc *bitCube) GetLength() DimensionLength {
	return DimensionLength{
		D0: bc.length[0],
		D1: bc.length[1],
		D2: bc.length[2],
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

// FyneBitMap の InterFace の実装
func (lbms labeledBitMaps) GetImage(f int) image.Image {
	return lbms[f]
}

// FyneBitMap の InterFace の実装
func (lbms labeledBitMaps) GetLength() DimensionLength {
	return DimensionLength{
		D0: len(lbms),
		D1: len(lbms[0]),
		D2: len(lbms[0][0]),
	}
}

// imgae.Image の InterFace を実装
func (lbm labeledBitMap) ColorModel() color.Model {
	return color.RGBAModel
}

// imgae.Image の InterFace を実装
func (lbm labeledBitMap) Bounds() image.Rectangle {
	return image.Rect(0, 0, len(lbm[0]), len(lbm))
}

// imgae.Image の InterFace を実装
func (lbm labeledBitMap) At(x, y int) color.Color {
	rect := image.Rect(0, 0, len(lbm[0]), len(lbm))
	if !(image.Point{x, y}.In(rect)) {
		return color.RGBA{0, 0, 0, 0}
	}
	value := lbm[y][x]
	if value < 0 {
		switch -value % 6 {
		case 1:
			return color.RGBA{255, 200, 200, 255}
		case 2:
			return color.RGBA{200, 255, 200, 255}
		case 3:
			return color.RGBA{200, 200, 255, 255}
		case 4:
			return color.RGBA{255, 200, 255, 255}
		case 5:
			return color.RGBA{200, 255, 255, 255}
		case 0:
			return color.RGBA{255, 255, 200, 255}
		default:
			return color.RGBA{0, 0, 0, 255}
		}
	}
	if value == 0 {
		return color.RGBA{0, 0, 0, 0}
	}
	switch value % 6 {
	case 1:
		return color.RGBA{255, 0, 0, 255}
	case 2:
		return color.RGBA{0, 255, 0, 255}
	case 3:
		return color.RGBA{0, 0, 255, 255}
	case 4:
		return color.RGBA{255, 0, 255, 255}
	case 5:
		return color.RGBA{0, 255, 255, 255}
	case 0:
		return color.RGBA{255, 255, 0, 255}
	default:
		return color.RGBA{0, 0, 0, 255}
	}
}
