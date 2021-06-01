package plyio

// １フレームの構造体
type frame struct {
	Index       int
	Coordinates [][2]int
}

// フレーム全体の構造体
type Frames struct {
	Width, Height, NumFrame int
	Data                    []frame
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
	frames.NumFrame = points.numFrame()
	frames.Height, biasH = points.frameHeight()
	frames.Width, biasW = points.frameWidth()

	frames.Data = make([]frame, frames.NumFrame)

	f := 0
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
