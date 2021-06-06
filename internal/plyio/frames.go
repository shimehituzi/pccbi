package plyio

// １フレームの構造体
type Frame struct {
	Index       int
	Coordinates [][2]int
}

// フレーム全体の構造体
type Frames struct {
	width, height, numFrame int
	Data                    []Frame
}

// Frames のコンストラクタ
func NewFrames() *Frames {
	frames := new(Frames)
	return frames
}

// ply の構造体から Frames の構造体の形式で読み込む
func (frames *Frames) ReadPoints(points *Points) {
	so := points.sortOrders
	data := points.data

	var biasH, biasW int
	frames.numFrame = points.numFrame()
	frames.height, biasH = points.frameHeight()
	frames.width, biasW = points.frameWidth()

	frames.Data = make([]Frame, frames.numFrame)

	f := 0
	frames.Data[f].Index = data[0][so[0]]
	for i := range data {
		if i != 0 && data[i][so[0]] != data[i-1][so[0]] {
			f++
			frames.Data[f].Index = data[i][so[0]]
		}
		frames.Data[f].Coordinates = append(
			frames.Data[f].Coordinates,
			[2]int{data[i][so[1]] - biasH, data[i][so[2]] - biasW},
		)
	}
}

// BitMaps の構造体から Frames の構造体の形式で読み込む
func (frames *Frames) ReadBitmaps(bm *BitMaps) {
	frames.width = bm.width
	frames.height = bm.height
	frames.numFrame = bm.numFrame

	frames.Data = make([]Frame, frames.numFrame)

	for i := range bm.Data {
		frames.Data[i].Index = bm.frameList[i]
		for j := range bm.Data[i] {
			for k := range bm.Data[i][j] {
				if bm.Data[i][j][k] == 1 {
					frames.Data[i].Coordinates = append(frames.Data[i].Coordinates, [2]int{j, k})
				}
			}
		}
	}
}
