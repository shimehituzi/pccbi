package labeling

import (
	"image"
	"image/color"

	"github.com/shimehituzi/pccbi/internal/bitmap"
)

// FyneBitMap の interface の実装
func (lbms *LabeledBitMaps) GetLength() bitmap.DimensionLength {
	return bitmap.DimensionLength{
		D0: len(*lbms),
		D1: len((*lbms)[0].Image),
		D2: len((*lbms)[0].Image[0]),
	}
}

// FyneBitMap の interface の実装
func (lbms *LabeledBitMaps) GetImage(f int) image.Image {
	return &((*lbms)[f])
}

// imgae.Image の InterFace を実装
func (lbm *LabeledBitMap) ColorModel() color.Model {
	return color.RGBAModel
}

// imgae.Image の InterFace を実装
func (lbm *LabeledBitMap) Bounds() image.Rectangle {
	return image.Rect(0, 0, len(lbm.Image[0]), len(lbm.Image))
}

// imgae.Image の InterFace を実装
func (lbm *LabeledBitMap) At(x, y int) color.Color {
	rect := image.Rect(0, 0, len(lbm.Image[0]), len(lbm.Image))
	if !(image.Point{x, y}.In(rect)) {
		return color.RGBA{0, 0, 0, 0}
	}
	if lbm.Image[y][x] == 1 {
		label := uint8(lbm.GetCounterLabel(x, y))
		l := uint8(label * 10)
		switch label % 6 {
		case 0:
			return color.RGBA{255, l, l, 255}
		case 1:
			return color.RGBA{l, 255, l, 255}
		case 2:
			return color.RGBA{l, l, 255, 255}
		case 3:
			return color.RGBA{255, l, 255, 255}
		case 4:
			return color.RGBA{255, 255, l, 255}
		case 5:
			return color.RGBA{l, 255, 255, 255}
		default:
			return color.RGBA{255, 255, 255, 255}
		}
	} else {
		return color.RGBA{0, 0, 0, 0}
	}
}

// ChainCode のラベルを返す
func (lbm *LabeledBitMap) GetCounterLabel(x, y int) int {
	for _, contour := range lbm.Contour {
		for _, point := range contour.ChainCode.Points {
			if point.X == x && point.Y == y {
				return contour.Label
			}
		}
	}
	panic("Label がありませんでした")
}
