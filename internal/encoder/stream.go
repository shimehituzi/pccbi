package encoder

func NewStream(cb contourBuffer) *Stream {
	startPoints := [][3]uint{}
	codes := []uint{}
	numCodesArray := []uint{}

	for f, contours := range cb {
		for _, contour := range contours {
			for _, chaincode := range contour {
				startPoint := [3]uint{uint(f), uint(chaincode.start.y), uint(chaincode.start.x)}
				startPoints = append(startPoints, startPoint)
				for _, code := range chaincode.code {
					codes = append(codes, uint(code))
				}
				codes = append(codes, 8)
			}
			numCodesArray = append(numCodesArray, uint(len(contour)))
		}
	}

	stream := Stream{
		StartPoints:   startPoints,
		Codes:         codes,
		NumCodesArray: numCodesArray,
	}

	return &stream
}
