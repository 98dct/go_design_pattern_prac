package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go-design-pattern-prac/auth"
	"net/http"
	"strconv"
)

func ApiAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		source := c.Request.Header.Get("source")
		timestamp := c.Request.Header.Get("times")
		sign := c.Request.Header.Get("signature")
		path := c.FullPath()
		method := c.Request.Method

		if err := NewVerifySign(source, method, path, timestamp, sign, int64(180)); err != nil {
			c.JSON(http.StatusOK, gin.H{"msg": "No Authorizated"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// NewVerifySign  验证验签信息
func NewVerifySign(source, method, path, timestamp, sign string, signExpire int64) error {
	if source == "" || path == "" || timestamp == "" || sign == "" {
		return errors.New("缺少header参数")
	}
	//验证source和path
	mysqlMgr := auth.NewMysqlRepository()
	secret, _ := mysqlMgr.GetPassword(source, method, path)

	// 验证过期时间
	tunix, _ := strconv.Atoi(timestamp)
	tokenMgr := auth.NewAuthToken(sign, int64(tunix), signExpire)

	if tokenMgr.IsExpire() {
		return errors.New("token expire")
	}
	//得到正确的sign供检验用
	tokenMgrNew := tokenMgr.GenerateToken(source, int64(tunix), secret)
	if !tokenMgrNew.IsMatch(sign) {
		return errors.New("signature不正确")
	}
	return nil
}
