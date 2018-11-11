package mill

type Blob struct{}

func (m *Blob) ID() string {
	return "/blob"
}

func (m *Blob) Encrypt() bool {
	return true
}

func (m *Blob) Pin() bool {
	return false
}

func (m *Blob) AcceptMedia(media string) error {
	return nil
}

func (m *Blob) Mill(input []byte, name string) (*Result, error) {
	return &Result{File: input}, nil
}
