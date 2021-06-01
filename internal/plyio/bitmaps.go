package plyio

type BitMaps struct {
	width, height, frame int
	data                 [][][]int8
}

func NewBitMaps() *BitMaps {
	pbms := new(BitMaps)
	return pbms
}

func (bm *BitMaps) ReadPoints(points *Points) {
	var biasH, biasW int
	bm.frame = points.numFrame()
	bm.height, biasH = points.frameHeight()
	bm.width, biasW = points.frameWidth()

	bm.data = make([][][]int8, bm.frame)
	for i := range bm.data {
		bm.data[i] = make([][]int8, bm.height)
		for j := range bm.data[i] {
			bm.data[i][j] = make([]int8, bm.width)
		}
	}

	pdata := points.data
	so := points.sortOrders

	for i := range points.data {
		bm.data[pdata[i][so[0]]-pdata[0][so[0]]][pdata[i][so[1]]-biasH][pdata[i][so[2]]-biasW] = 1
	}
}
