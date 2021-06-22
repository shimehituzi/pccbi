package labeling

import (
	"image"
	"image/color"

	"github.com/shimehituzi/pccbi/internal/bitmap"
)

func (lbm *LabeledBitMap) GetLength() bitmap.DimensionLength {
	return bitmap.DimensionLength{
		D0: 0,
		D1: len(lbm.Image),
		D2: len(lbm.Image[0]),
	}
}

func (lbm *LabeledBitMap) GetImage(int) image.Image {
	return lbm
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
			return color.RGBA{0, 0, 0, 0}
		}
	} else {
		return color.RGBA{255, 255, 255, 255}
	}
}
