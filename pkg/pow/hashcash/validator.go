package hashcash

type ValidatorImpl struct {
	hFn HashFunc
}

func NewValidator(hFn HashFunc) *ValidatorImpl {
	return &ValidatorImpl{hFn: hFn}
}

func (v *ValidatorImpl) Validate(prefix, nonce []byte, c int) error {
	h := v.hFn()
	if c >= h.Size()*8 {
		return ErrHComplexity
	}
	h.Write(prefix)
	h.Write(nonce)
	if countLeadingZeroBits(h.Sum(nil)) < c {
		return ErrHNotValidated
	}
	return nil
}
