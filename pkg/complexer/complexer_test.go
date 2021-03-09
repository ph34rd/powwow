package complexer

import (
	"strconv"
	"testing"
)

type cpuSamplerMock struct {
	i      int
	values []float64
}

func (c *cpuSamplerMock) GetCPUUsageFraction() (v float64) {
	if c.i == len(c.values) {
		c.i = 0
	}
	v = c.values[c.i]
	c.i++
	return
}

func (c cpuSamplerMock) Shutdown() {}

func Test_ComplexerImpl(t *testing.T) {
	sampler := &cpuSamplerMock{values: []float64{0, 0.1, 0.2, 0.5, 0.8, 1}}
	cm := NewComplexer(sampler)
	expect := []uint8{20, 21, 23, 27, 32, 32}
	for i, test := range expect {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			res := cm.NextComplexity()
			if res != test {
				t.Fatalf("expected result: %d, got: %d", test, res)
			}
		})
	}
}
