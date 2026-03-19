package auth

import (
	"github.com/jadecobra/agbalumo/internal/config"
	"github.com/jadecobra/agbalumo/internal/domain"
)

// AuthDependencies encapsulates the dependencies required by AuthHandler.
type AuthDependencies struct {
	UserStore      domain.UserStore
	GoogleProvider GoogleProvider
	Config         *config.Config
}
