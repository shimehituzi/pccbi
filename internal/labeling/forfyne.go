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
func (lbms *LabeledBitMaps) GetLabelLength(f int) int {
	length := len((*lbms)[f].Segment)
	labelLength := 0
	if length != 0 {
		labelLength = (*lbms)[f].Segment[length-1].Label
	}
	return labelLength
}

// FyneBitMap の interface の実装
func (lbms *LabeledBitMaps) GetImage(f int, l int) image.Image {
	lbm := (*lbms)[f]
	if l == 0 {
		return &lbm
	}
	if len(lbm.Segment) < l {
		return &lbm
	}

	oneOfLabeled := new(LabeledBitMap)
	oneOfLabeled.Segment = []Segment{lbm.Segment[l-1]}
	oneOfLabeled.Image = make([][]byte, len(lbm.Image))
	for i := range lbm.Image {
		oneOfLabeled.Image[i] = make([]byte, len(lbm.Image[i]))
	}
	for _, contour := range oneOfLabeled.Segment[0].Contours {
		for _, point := range contour.ChainCode.Points {
			oneOfLabeled.Image[point.Y][point.X] = 1
		}
	}

	return oneOfLabeled
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
		label := lbm.GetCounterLabel(x, y)
		if label == 0 {
			return color.RGBA{255, 255, 255, 255}
			// return color.RGBA{0, 0, 0, 0}
		}
		l := uint8(label * 10)
		switch label % 6 {
		case 1:
			return color.RGBA{255, l, l, 255}
		case 2:
			return color.RGBA{l, 255, l, 255}
		case 3:
			return color.RGBA{l, l, 255, 255}
		case 4:
			return color.RGBA{255, l, 255, 255}
		case 5:
			return color.RGBA{255, 255, l, 255}
		case 0:
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
	for _, segment := range lbm.Segment {
		for i, contour := range segment.Contours {
			for _, point := range contour.ChainCode.Points {
				if point.X == x && point.Y == y {
					if i == 0 {
						return 0
					}
					return segment.Label
				}
			}
		}
	}
	panic("Label がありませんでした")
}
