package plyio

func LordPly(srcfile string) (*BitMaps, error) {
	ply, err := NewPly(srcfile)
	if err != nil {
		return nil, err
	}

	points, err := NewPoints(ply, [3]int{1, 2, 0})
	if err != nil {
		return nil, err
	}

	bms := NewBitMaps(points)

	return bms, nil
}
