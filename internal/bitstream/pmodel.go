package bitstream

type Pmodel struct {
	freq, cumfreq   []uint64
	totfreq, offset uint64
}

func (pm Pmodel) Freq() []uint64 {
	return pm.freq
}

func NewEncPmodel(val []uint, min, max uint) *Pmodel {
	freq := make([]uint64, max+1)
	for _, v := range val {
		freq[v]++
	}
	cumfreq := make([]uint64, max+2)
	cumfreq[0] = 0
	for i := range freq {
		cumfreq[i+1] = cumfreq[i] + freq[i]
	}
	offset := cumfreq[min]
	totfreq := cumfreq[max+1] - offset

	return &Pmodel{freq, cumfreq, totfreq, offset}
}

func NewDecPmodel(freq []uint64, min, max uint) *Pmodel {
	cumfreq := make([]uint64, max+2)
	cumfreq[0] = 0
	for i := range freq {
		cumfreq[i+1] = cumfreq[i] + freq[i]
	}
	offset := cumfreq[min]
	totfreq := cumfreq[max+1] - offset

	return &Pmodel{freq, cumfreq, totfreq, offset}
}
