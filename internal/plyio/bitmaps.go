package plyio

// 2値画像の集合として点群の位置情報を表す構造体
type BitMaps struct {
	// ある次元の要素の値の幅
	Length [3]int
	// ある次元の要素の値のバイアス
	Bias [3]int
	// 2値画像の集合
	Data []BitMap
}

// BitMaps のコンストラクタ
func NewBitMaps() *BitMaps {
	pbms := new(BitMaps)
	return pbms
}

// ply の構造体から BitMaps の構造体の形式で読み込む
func (bms *BitMaps) ReadPoints(points *Points) {
	so := points.sortOrders
	for i := 0; i < 3; i++ {
		bms.Length[i], bms.Bias[i] = points.getLengthAndBias(so[i])
	}

	bms.Data = make([]BitMap, bms.Length[0])
	for i := range bms.Data {
		bms.Data[i] = make([][]byte, bms.Length[1])
		for j := range bms.Data[i] {
			bms.Data[i][j] = make([]byte, bms.Length[2])
		}
	}

	pdata := points.data
	for i := range points.data {
		dim0 := pdata[i][so[0]] - bms.Bias[0]
		dim1 := pdata[i][so[1]] - bms.Bias[1]
		dim2 := pdata[i][so[2]] - bms.Bias[2]
		bms.Data[dim0][dim1][dim2] = 1
	}
}
