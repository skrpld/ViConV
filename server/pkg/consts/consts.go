package consts

import "time"

const (
	RefreshTokenExpiryTime = time.Hour * 24 * 7 //TODO: вынести в real-time config
	AccessTokenExpiryTime  = time.Hour * 2      //TODO: вынести в real-time config
	IssuedAtField          = "viconv"           //TODO: вынести в real-time config + rename

	CtxUserKey  = "user"
	NoChangeKey = "no_change"
)
