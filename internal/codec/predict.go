package codec

const divisor = 8

func encDecorrelation(code []byte) []byte {
	return encDiff(code)
}

func decDecorrelation(code []byte) []byte {
	return decDiff(code)
}

func encDiff(code []byte) (diff []byte) {
	diff = make([]byte, len(code))
	for i := range code {
		if i == 0 {
			diff[i] = code[i]
		} else {
			diff[i] = (code[i] - code[i-1]) % divisor
		}
	}
	return
}

func decDiff(diff []byte) (code []byte) {
	code = make([]byte, len(diff))
	for i := range diff {
		if i == 0 {
			code[i] = diff[i]
		} else {
			code[i] = (diff[i] + code[i-1]) % divisor
		}
	}
	return
}
