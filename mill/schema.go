package mill

type Schema struct{}

func (m *Schema) ID() string {
	return "/schema"
}

func (m *Schema) Encrypt() bool {
	return false
}

func (m *Schema) Pin() bool {
	return true
}

func (m *Schema) AcceptMedia(media string) error {
	return accepts([]string{"application/json"}, media)
}

func (m *Schema) Options() (string, error) {
	return "", nil
}

func (m *Schema) Mill(input []byte, name string) (*Result, error) {
	return &Result{File: input}, nil
}
