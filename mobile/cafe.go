package mobile

import (
	"github.com/textileio/textile-go/repo"
)

// CafeSessions is a wrapper around a list of sessions
type CafeSessions struct {
	Items []repo.CafeSession `json:"items"`
}

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
	sessions := &CafeSessions{Items: make([]repo.CafeSession, 0)}
	if len(items) > 0 {
		sessions.Items = items
	}
	return toJSON(sessions)
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

// CheckCafeMail calls core CheckCafeMessages
func (m *Mobile) CheckCafeMail() error {
	return m.node.CheckCafeMail()
}
