package auth

import (
	"crypto/md5"
	"fmt"
	"strconv"
	"time"
)

type AuthToken struct {
	token      string
	createTime int64
	lifeTime   int64
}

func NewAuthToken(token string, createTime, lifeTime int64) *AuthToken {
	return &AuthToken{
		token:      token,
		createTime: createTime,
		lifeTime:   lifeTime,
	}
}

func (at *AuthToken) IsExpire() bool {
	unix := time.Now().Unix()
	if at.createTime >= unix || (at.createTime+at.lifeTime) < unix {
		return true
	}
	return false
}

func (at *AuthToken) GenerateToken(source string, createTime int64, secret string) AuthToken {

	return AuthToken{
		token: fmt.Sprintf("%x", md5.Sum([]byte(source+strconv.Itoa(int(createTime))+secret))),
	}
}

func (at *AuthToken) IsMatch(token string) bool {
	return at.token == token
}
