package refactoring

import (
	"image"
	"image/color"
)

type FyneBitMap interface {
	GetLength() DimensionLength
	GetImage(int) image.Image
}

type DimensionLength struct {
	D0 int
	D1 int
	D2 int
}

func (lv labeledVoxel) GetLength() DimensionLength {
	return DimensionLength{
		D0: len(lv),
		D1: len(lv[0]),
		D2: len(lv[0][0]),
	}
}

func (lv labeledVoxel) GetImage(i int) image.Image {
	return lv[i]
}

func (l label) ColorModel() color.Model {
	return color.RGBAModel
}

func (l label) Bounds() image.Rectangle {
	return image.Rect(0, 0, len(l[0]), len(l))
}

func (l label) At(x, y int) color.Color {
	rect := image.Rect(0, 0, len(l[0]), len(l))
	if !(image.Point{x, y}.In(rect)) {
		return color.RGBA{0, 0, 0, 0}
	}
	value := uint8(l[y][x])
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
