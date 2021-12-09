package codec

func EncStream(contour contour) *Stream {
	startPoints := [][3]uint{}
	codes := []uint{}
	numCodesArray := []uint{}

	for f, contourFrame := range contour {
		for _, contourSegment := range contourFrame {
			for _, chaincode := range contourSegment {
				startPoint := [3]uint{uint(f), uint(chaincode.start.y), uint(chaincode.start.x)}
				startPoints = append(startPoints, startPoint)
				for _, code := range chaincode.code {
					codes = append(codes, uint(code))
				}
				codes = append(codes, 8)
			}
			numCodesArray = append(numCodesArray, uint(len(contourSegment)))
		}
	}

	stream := Stream{
		StartPoints:   startPoints,
		Codes:         codes,
		NumCodesArray: numCodesArray,
	}

	return &stream
}
