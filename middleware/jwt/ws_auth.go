package jwt

import (
	"errors"
	"strings"
)

func WsAuthInHeader(token string) (jwtUid int, err error) {
	if token == "" {
		return -1, errors.New("无效的Token")
	}
	tokenSplits := strings.SplitN(token, " ", 2)
	if len(tokenSplits) != 2 || tokenSplits[0] != "Bearer" {
		return -1, errors.New("无效的Token")
	}
	claims, err := parseToken(tokenSplits[1])
	if err != nil {
		return -1, errors.New("无效的Token")
	}
	return claims.UserId, nil
}
