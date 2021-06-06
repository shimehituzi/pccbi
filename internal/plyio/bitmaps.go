package plyio

// 2値画像の集合として点群の位置情報を表す構造体
type BitMaps struct {
	width, height, numFrame int
	frameList               []int
	Data                    [][][]int8
}

// BitMaps のコンストラクタ
func NewBitMaps() *BitMaps {
	pbms := new(BitMaps)
	return pbms
}

// ply の構造体から BitMaps の構造体の形式で読み込む
func (bm *BitMaps) ReadPoints(points *Points) {
	var biasH, biasW int
	bm.numFrame = points.numFrame()
	bm.height, biasH = points.frameHeight()
	bm.width, biasW = points.frameWidth()

	bm.frameList = make([]int, bm.numFrame)
	bm.Data = make([][][]int8, bm.numFrame)
	for i := range bm.Data {
		bm.Data[i] = make([][]int8, bm.height)
		for j := range bm.Data[i] {
			bm.Data[i][j] = make([]int8, bm.width)
		}
	}

	pdata := points.data
	so := points.sortOrders

	for i := range points.data {
		bm.frameList[pdata[i][so[0]]-pdata[0][so[0]]] = pdata[i][so[0]]
		bm.Data[pdata[i][so[0]]-pdata[0][so[0]]][pdata[i][so[1]]-biasH][pdata[i][so[2]]-biasW] = 1
	}
}

// Frames の構造体から BitMaps の構造体の形式で読み込む
func (bm *BitMaps) ReadFrames(frames *Frames) {
	bm.width = frames.width
	bm.height = frames.height
	bm.numFrame = frames.numFrame

	bm.frameList = make([]int, bm.numFrame)
	bm.Data = make([][][]int8, bm.numFrame)
	for i := range bm.Data {
		bm.Data[i] = make([][]int8, bm.height)
		for j := range bm.Data[i] {
			bm.Data[i][j] = make([]int8, bm.width)
		}
	}

	for i := range frames.Data {
		bm.frameList[i] = frames.Data[i].Index
		data := frames.Data[i].Coordinates
		for j := range data {
			bm.Data[i][data[j][0]][data[j][1]] = 1
		}
	}
}
