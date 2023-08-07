package auth

import (
	"github.com/filebrowser/filebrowser/v2/settings"
)

// MethodNoAuth is used to identify no auth.
const MethodNoAuth settings.AuthMethod = "noauth"

// NoAuth is no auth implementation of auther.
type NoAuth struct{}

// LoginPage tells that no auth doesn't require a login page.
func (a NoAuth) LoginPage() bool {
	return false
}
