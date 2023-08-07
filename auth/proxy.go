package auth

import (
	"github.com/filebrowser/filebrowser/v2/settings"
)

// MethodProxyAuth is used to identify no auth.
const MethodProxyAuth settings.AuthMethod = "proxy"

// ProxyAuth is a proxy implementation of an auther.
type ProxyAuth struct {
	Header string `json:"header"`
}

// LoginPage tells that proxy auth doesn't require a login page.
func (a ProxyAuth) LoginPage() bool {
	return false
}
