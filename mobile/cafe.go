package mobile

// RegisterCafe calls core RegisterCafe
func (m *Mobile) RegisterCafe(peerId string) error {
	if _, err := m.node.RegisterCafe(peerId); err != nil {
		return err
	}
	return nil
}

// CafeSessions calls core CafeSessions
func (m *Mobile) CafeSessions() (string, error) {
	items, err := m.node.CafeSessions()
	if err != nil {
		return "", err
	}
	if len(items) == 0 {
		return "[]", nil
	}
	return toJSON(items)
}

// CafeSession calls core CafeSession
func (m *Mobile) CafeSession(peerId string) (string, error) {
	session, err := m.node.CafeSession(peerId)
	if err != nil {
		return "", err
	}
	if session == nil {
		return "", nil
	}
	return toJSON(session)
}

// RefreshCafeSession calls core RefreshCafeSession
func (m *Mobile) RefreshCafeSession(peerId string) (string, error) {
	session, err := m.node.RefreshCafeSession(peerId)
	if err != nil {
		return "", err
	}
	return toJSON(session)
}

// DeegisterCafe calls core DeregisterCafe
func (m *Mobile) DeregisterCafe(peerId string) error {
	return m.node.DeregisterCafe(peerId)
}

// CheckCafeMessages calls core CheckCafeMessages
func (m *Mobile) CheckCafeMessages() error {
	return m.node.CheckCafeMessages()
}
