package complexer

import (
	"math"

	"github.com/ph34rd/powwow/pkg/stat"
)

const (
	complexityMin = 20
	complexityMax = 32
)

type Complexer interface {
	NextComplexity() uint8
}

type Impl struct {
	sampler stat.CPUSampler
}

func NewComplexer(sampler stat.CPUSampler) *Impl {
	return &Impl{sampler: sampler}
}

func (c *Impl) NextComplexity() uint8 {
	frac := c.sampler.GetCPUUsageFraction()
	rng := float64(complexityMax-complexityMin) / 0.8
	val := math.Floor(rng*frac) + complexityMin
	if val > complexityMax {
		return complexityMax
	}
	return uint8(val)
}
