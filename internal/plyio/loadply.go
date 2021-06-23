package plyio

func LordPly(srcfile string, sortOrders SortOrders) (*BitMaps, error) {
	ply, err := NewPly(srcfile)
	if err != nil {
		return nil, err
	}

	points, err := NewPoints(ply, sortOrders)
	if err != nil {
		return nil, err
	}

	bms := NewBitMaps(points)

	return bms, nil
}
