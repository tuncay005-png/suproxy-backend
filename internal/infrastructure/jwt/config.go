package jwt

import "github.com/suproxy/backend/internal/infrastructure/config"

func (m *Manager) Config() *config.JWTConfig {
	return m.config
}
