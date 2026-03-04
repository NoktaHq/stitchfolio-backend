package responseModel

import (
	"time"
)

type User struct {
	ID       uint `json:"id,omitempty"`
	IsActive bool `json:"isActive,omitempty"`

	FirstName           string     `json:"firstName,omitempty"`
	LastName            string     `json:"lastName,omitempty"`
	Extension           string     `json:"extension,omitempty"`
	PhoneNumber         string     `json:"phoneNumber,omitempty"`
	Email               string     `json:"email,omitempty"`
	Role                string     `json:"role,omitempty"`
	IsLoginDisabled     bool       `json:"isLoginDisabled,omitempty"`
	IsLoggedIn          bool       `json:"isLoggedIn,omitempty"`
	LastLoginTime       *time.Time `json:"lastLoginTime,omitempty"`
	LoginFailureCounter int16      `json:"loginFailureCounter,omitempty"`
	ResetPasswordString *string    `json:"resetPasswordString,omitempty"`
	Experience          string     `json:"experience,omitempty"`
	Department          string     `json:"department,omitempty"`

	UserChannelDetails []UserChannelDetail `json:"userChannelDetails,omitempty"`
}

type UserAutoComplete struct {
	UserID    uint   `json:"userId,omitempty"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
}

type UserChannelDetail struct {
	ID        uint   `json:"id,omitempty"`
	ChannelID uint   `json:"channelId,omitempty"`
	Name      string `json:"name,omitempty"`
}
