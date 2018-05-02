package flylsh

import (
	"math"
	"math/bits"
	"math/rand"
)

const (
	sizeOfByte    = 8
	defaultNGrams = 4
)

var DefaultFlyLSH = New(defaultNGrams, int(defaultNGrams*sizeOfByte), int(10*defaultNGrams*sizeOfByte), 1234567890, 0.1, 0.05)

type FlyLSH struct {
	nGrams      uint
	nORNs       int
	nKCs        int
	nORNSamples int
	ORNs        []float64
	PNs         []float64
	M           [][]byte
	KCs         []float64
	nForAPL     int
}

func New(nGrams, nORNs, nKCs int, MSeed int64, randProjectPercent, APLPercent float64) (h *FlyLSH) {
	h = &FlyLSH{
		nGrams: uint(nGrams),
		nORNs:  nORNs,
		nKCs:   nKCs,
		//ORNs:    make([]float64, nORNs),
		PNs:     make([]float64, nORNs),
		M:       make([][]byte, nKCs),
		KCs:     make([]float64, nKCs),
		nForAPL: int(math.Ceil(float64(nKCs) * APLPercent)),
	}
	// Random seed for generate random project matrix
	rand.Seed(MSeed)
	// Sample number of ORNs pre Kenyon cell
	h.nORNSamples = int(math.Ceil(randProjectPercent * float64(nORNs)))
	// Mark connection between ORN and KC in matrix
	for j := range h.M {
		h.M[j] = make([]byte, nORNs/sizeOfByte+1)
		// random sample ORNs for each KC
		for m := 0; m < h.nORNSamples; m++ {
			i := uint(rand.Intn(nORNs))
			h.M[j][i/sizeOfByte] |= 1 << (i & (sizeOfByte - 1))
		}
	}
	return
}

func Sum(s []byte) []byte {
	return DefaultFlyLSH.Sum(s)
}

func Compare(b0, b1 []byte) (d int) {
	return DefaultFlyLSH.Compare(b0, b1)
}

func (h *FlyLSH) Sum(s []byte) []byte {
	defer h.clear()
	vector := h.getVector(s)
	if len(vector) != h.nORNs {
		panic("vector length must equal length of ORNs")
	}
	h.ORNs = vector
	h.zeroCenterORNs()
	h.randomProject()
	h.aplWTA()
	return h.aplIndices()
}

func (h *FlyLSH) Compare(b0, b1 []byte) (d int) {
	for i := range b0 {
		b := b0[i] ^ b1[i]
		d += bits.OnesCount8(b)
	}
	return
}

func (h *FlyLSH) Reset() {
	for j := range h.M {
		for i := range h.M[j] {
			h.M[j][i] = 0
		}
		h.KCs[j] = 0
	}
}

func (h *FlyLSH) getVector(s []byte) (v []float64) {
	v = make([]float64, h.nGrams*sizeOfByte)
	for i := uint(0); i < uint(len(s))-h.nGrams+1; i++ {
		for k := uint(0); k < h.nGrams; k++ {
			for j := uint(0); j < sizeOfByte; j++ {
				if s[i+k]&(1<<j) > 0 {
					v[sizeOfByte*k+j]++
				}
			}
		}
	}
	return
}

func (h *FlyLSH) zeroCenterORNs() {
	var mean float64
	for _, v := range h.ORNs {
		mean += v
	}
	mean /= float64(len(h.ORNs))
	for i := range h.PNs {
		h.PNs[i] = h.ORNs[i] - mean
	}
	//log.Info(h.PNs)
}

func (h *FlyLSH) randomProject() {
	for j := range h.M {
		for i := 0; i < h.nORNs; i++ {
			if h.M[j][i/sizeOfByte]&(1<<(uint(i)&(sizeOfByte-1))) > 0 {
				h.KCs[j] += h.PNs[i]
			}
		}
		h.KCs[j] /= float64(h.nORNSamples)
		//log.Println(h.KCs[j])
	}
}

func (h *FlyLSH) aplWTA() {
	heap := make([]float64, h.nForAPL)
	for _, v := range h.KCs {
		if v > heap[0] {
			heap[0] = v
			for k := 0; k < h.nForAPL/2; {
				if heap[k] > heap[2*k] || heap[k] > heap[2*k+1] {
					if heap[2*k] < heap[2*k+1] {
						heap[2*k], heap[k] = heap[k], heap[2*k]
						k = 2 * k
					} else {
						heap[2*k+1], heap[k] = heap[k], heap[2*k+1]
						k = 2*k + 1
					}
				} else {
					break
				}
			}
		}
	}
	for i, v := range h.KCs {
		if v < heap[0] {
			h.KCs[i] = 0
		}
	}
	//log.Info(heap, h.KCs)
}

func (h *FlyLSH) aplIndices() []byte {
	sum := make([]byte, len(h.KCs)/sizeOfByte)
	for j, v := range h.KCs {
		if v > 0 {
			sum[j/sizeOfByte] |= 1 << (uint(j) & (sizeOfByte - 1))
		}
	}
	return sum
}

func (h *FlyLSH) clear() {
	for i := range h.PNs {
		h.PNs[i] = 0
	}
	for j := range h.KCs {
		h.KCs[j] = 0
	}
}
