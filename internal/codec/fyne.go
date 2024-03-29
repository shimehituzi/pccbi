package codec

import (
	"image"
	"image/color"
)

type FyneBitMap interface {
	GetLength() DimensionLength
	GetLabelLength(int) int
	GetImage(int, int) image.Image
}

type DimensionLength struct {
	D0 int
	D1 int
	D2 int
}

func NewFyneContour(contour Contour, header *Header) labeledVoxel {
	fc := make(labeledVoxel, header.Length[0])
	for f := range fc {
		fc[f] = make(label, header.Length[1])
		for y := range fc[f] {
			fc[f][y] = make([]int, header.Length[2])
		}
	}

	for f := range contour {
		for l := range contour[f] {
			for c, cc := range contour[f][l] {
				if c == 0 {
					for _, point := range cc.getPoints() {
						fc[f][point.y][point.x] = l*2 + 1
					}
				} else {
					for _, point := range cc.getPoints() {
						fc[f][point.y][point.x] = l*2 + 2
					}
				}
			}
		}
	}

	return fc
}

func (lv labeledVoxel) GetLength() DimensionLength {
	return DimensionLength{
		D0: len(lv),
		D1: len(lv[0]),
		D2: len(lv[0][0]),
	}
}

func (lv labeledVoxel) GetLabelLength(int) int {
	return 0
}

func (lv labeledVoxel) GetImage(i, _ int) image.Image {
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
	switch value % 12 {
	case 1:
		return color.RGBA{255, 0, 0, 255}
	case 2:
		return color.RGBA{255, 125, 125, 255}
	case 3:
		return color.RGBA{0, 255, 0, 255}
	case 4:
		return color.RGBA{125, 255, 125, 255}
	case 5:
		return color.RGBA{0, 0, 255, 255}
	case 6:
		return color.RGBA{125, 125, 255, 255}
	case 7:
		return color.RGBA{255, 0, 255, 255}
	case 8:
		return color.RGBA{255, 125, 255, 255}
	case 9:
		return color.RGBA{0, 255, 255, 255}
	case 10:
		return color.RGBA{125, 255, 255, 255}
	case 11:
		return color.RGBA{255, 255, 0, 255}
	case 0:
		return color.RGBA{255, 255, 125, 255}
	default:
		return color.RGBA{0, 0, 0, 255}
	}
}

func (fs frames) GetLength() DimensionLength {
	return DimensionLength{
		D0: len(fs),
		D1: len(fs[0][0]),
		D2: len(fs[0][0][0]),
	}
}

func (fs frames) GetLabelLength(f int) int {
	return len(fs[f])
}

func (fs frames) GetImage(f, l int) image.Image {
	if l == 0 {
		return fs[f]
	}
	return fs[f][l-1]
}

func (f frame) ColorModel() color.Model {
	return color.RGBAModel
}

func (f frame) Bounds() image.Rectangle {
	return image.Rect(0, 0, len(f[0][0]), len(f[0]))
}

func (f frame) At(x, y int) color.Color {
	rect := image.Rect(0, 0, len(f[0][0]), len(f[0]))
	if !(image.Point{x, y}.In(rect)) {
		return color.RGBA{0, 0, 0, 0}
	}
	for l := range f {
		if f[l][y][x] == 1 {
			return color.RGBA{255, 255, 255, 255}
		}
	}
	return color.RGBA{0, 0, 0, 0}
}

func (s segment) ColorModel() color.Model {
	return color.RGBAModel
}

func (s segment) Bounds() image.Rectangle {
	return image.Rect(0, 0, len(s[0]), len(s))
}

func (s segment) At(x, y int) color.Color {
	rect := image.Rect(0, 0, len(s[0]), len(s))
	if !(image.Point{x, y}.In(rect)) {
		return color.RGBA{0, 0, 0, 0}
	}
	if s[y][x] == 1 {
		return color.RGBA{255, 255, 255, 255}
	} else {
		return color.RGBA{0, 0, 0, 0}
	}
}
