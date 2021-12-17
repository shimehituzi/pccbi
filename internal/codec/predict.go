package codec

const divisor = 8

func encDecorrelation(code []byte) (diff []byte) {
	diff = make([]byte, len(code))
	for i := range diff {
		diff[i] = (code[i] - naivePredict(code[:i])) % divisor
	}
	return
}

func decDecorrelation(diff []byte) (code []byte) {
	code = make([]byte, len(diff))
	for i := range code {
		code[i] = (diff[i] + naivePredict(code[:i])) % divisor
	}
	return
}

func naivePredict(data []byte) (pred byte) {
	if len(data) == 0 {
		return byte(0)
	}
	pred = data[len(data)-1]
	return
}
