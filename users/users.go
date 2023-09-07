package users

import (
	"path/filepath"

	"github.com/spf13/afero"

	"github.com/filebrowser/filebrowser/v2/errors"
	"github.com/filebrowser/filebrowser/v2/files"
	"github.com/filebrowser/filebrowser/v2/rules"
)

// ViewMode describes a view mode.
type ViewMode string

const (
	ListViewMode   ViewMode = "list"
	MosaicViewMode ViewMode = "mosaic"
)

// User describes a user.
type User struct {
	ID           uint          `storm:"id,increment" json:"id"`
	Username     string        `storm:"unique" json:"username"`
	Password     string        `json:"password"`
	Scope        string        `json:"scope"`
	Locale       string        `json:"locale"`
	ViewMode     ViewMode      `json:"viewMode"`
	SingleClick  bool          `json:"singleClick"`
	Perm         Permissions   `json:"perm"`
	Sorting      files.Sorting `json:"sorting"`
	Fs           afero.Fs      `json:"-" yaml:"-"`
	Rules        []rules.Rule  `json:"rules"`
	HideDotfiles bool          `json:"hideDotfiles"`
	DateFormat   bool          `json:"dateFormat"`
}

type EncryptedCredentials struct {
	Iv            string `json:"iv"`
	EncryptedData string `json:"encryptedData"`
}

type DecryptedCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Type     string `json:"type"`
	OU       string `json:"OU"`
	Hostname string `json:"hostname"`
}

type UserInfo struct {
	Locale               string               `json:"locale"`
	ViewMode             ViewMode             `json:"viewMode"`
	SingleClick          bool                 `json:"singleClick"`
	Perm                 Permissions          `json:"perm"`
	HideDotfiles         bool                 `json:"hideDotfiles"`
	DateFormat           bool                 `json:"dateFormat"`
	Scope                string               `json:"scope"`
	EncryptedCredentials EncryptedCredentials `json:"credentials"`
}

type TokenStruct struct {
	Scope                string               `json:"scope"`
	Locale               string               `json:"locale"`
	ViewMode             ViewMode             `json:"viewMode"`
	Perm                 Permissions          `json:"perm"`
	Fs                   afero.Fs             `json:"-" yaml:"-"`
	HideDotfiles         bool                 `json:"hideDotfiles"`
	EncryptedCredentials EncryptedCredentials `json:"credentiald"`
}

// GetRules implements rules.Provider.
func (u *User) GetRules() []rules.Rule {
	return u.Rules
}

var checkableFields = []string{
	"Username",
	"Password",
	"Scope",
	"ViewMode",
	"Sorting",
	"Rules",
}

// Clean cleans up a user and verifies if all its fields
// are alright to be saved.
//
//nolint:gocyclo
func (u *User) Clean(baseScope string, fields ...string) error {
	if len(fields) == 0 {
		fields = checkableFields
	}

	for _, field := range fields {
		switch field {
		case "Username":
			if u.Username == "" {
				return errors.ErrEmptyUsername
			}
		case "Password":
			if u.Password == "" {
				return errors.ErrEmptyPassword
			}
		case "ViewMode":
			if u.ViewMode == "" {
				u.ViewMode = ListViewMode
			}
		case "Sorting":
			if u.Sorting.By == "" {
				u.Sorting.By = "name"
			}
		case "Rules":
			if u.Rules == nil {
				u.Rules = []rules.Rule{}
			}
		}
	}

	if u.Fs == nil {
		scope := u.Scope
		scope = filepath.Join(baseScope, filepath.Join("/", scope)) //nolint:gocritic
		u.Fs = afero.NewBasePathFs(afero.NewOsFs(), scope)
	}

	return nil
}

// FullPath gets the full path for a user's relative path.
func (u *User) FullPath(path string) string {
	return afero.FullBaseFsPath(u.Fs.(*afero.BasePathFs), path)
}
