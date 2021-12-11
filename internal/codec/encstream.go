package codec

func EncStream(contour Contour) *Stream {
	startPoints := [][3]uint{}
	codes := []uint{}
	numCodesArray := []uint{}

	for f, contourFrame := range contour {
		for _, contourSegment := range contourFrame {
			for _, chaincode := range contourSegment {
				startPoint := [3]uint{uint(f), uint(chaincode.Start.y), uint(chaincode.Start.x)}
				startPoints = append(startPoints, startPoint)
				code := encDecorrelation(chaincode.Code)
				for _, v := range code {
					codes = append(codes, uint(v))
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
