package codec

import "sort"

func ComparePoint(a, b point) bool {
	if a.x == b.x && a.y == b.y {
		return true
	} else {
		return false
	}
}

// 閉曲線の内側なら true 線上は不定
func closedAreaDesicion(p point, cc chaincode) bool {
	ps := cc.Points
	wn := 0
	for i := 0; i < len(ps)-1; i++ {
		if (ps[i].y <= p.y) && (ps[i+1].y > p.y) {
			vt := (p.y - ps[i].y) / (ps[i+1].y - ps[i].y)
			if p.x < (ps[i].x + (vt * (ps[i+1].x - ps[i].x))) {
				wn++
			}
		} else if (ps[i].y > p.y) && (ps[i+1].y <= p.y) {
			vt := (p.y - ps[i].y) / (ps[i+1].y - ps[i].y)
			if p.x < (ps[i].x + (vt * (ps[i+1].x - ps[i].x))) {
				wn--
			}
		}
	}
	return wn != 0
}

func fillArea(img bitmap, p point, prev, value byte) {
	img[p.y][p.x] = value

	nearest4points := [4]point{
		{p.x, p.y - 1},
		{p.x - 1, p.y},
		{p.x + 1, p.y},
		{p.x, p.y + 1},
	}

	for _, p := range nearest4points {
		if validPointByte(p, img) && img[p.y][p.x] == prev {
			fillArea(img, p, prev, value)
		}
	}
}

func validPointInt(p point, img label) bool {
	if p.y < 0 || p.x < 0 || len(img) <= p.y || len(img[0]) <= p.x {
		return false
	}
	return true
}

func validPointByte(p point, img bitmap) bool {
	if p.y < 0 || p.x < 0 || len(img) <= p.y || len(img[0]) <= p.x {
		return false
	}
	return true
}

func (p point) in(points []point) bool {
	for _, point := range points {
		if p == point {
			return true
		}
	}
	return false
}

func uint2byte(uintSlice []uint) (byteSlice []byte) {
	byteSlice = make([]byte, len(uintSlice))
	for i, v := range uintSlice {
		byteSlice[i] = byte(v % 256)
	}
	return
}

func (ply Ply) Sort() {
	sort.Sort(ply)
}

func (ply Ply) Len() int { return len(ply) }

func (ply Ply) Swap(i, j int) { ply[i], ply[j] = ply[j], ply[i] }

func (ply Ply) Less(i, j int) bool {
	switch {
	case ply[i][0] < ply[j][0]:
		return true
	case ply[i][0] > ply[j][0]:
		return false
	default:
		switch {
		case ply[i][1] < ply[j][1]:
			return true
		case ply[i][1] > ply[j][1]:
			return false
		default:
			switch {
			case ply[i][2] < ply[j][2]:
				return true
			case ply[i][2] > ply[j][2]:
				return false
			default:
				return false
			}
		}
	}
}
