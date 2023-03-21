package jwt

import (
	"errors"
	"sso/pkg"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	User string `json:"user"`
	CN   string `json:"cn"`
	Mail string `json:"mail"`
	Uid  int64  `json:"uid"`
	Type bool   `json:"type"`
	jwt.StandardClaims
}

func NewClaims(uid, index int64, user, cn, mail string) *Claims {
	conf := pkg.Conf()
	t := time.Now()
	if index == 1 {
		return &Claims{
			Uid:  uid,
			User: user,
			CN:   cn,
			Mail: mail,
			Type: true,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: t.Add(time.Duration(conf.JwtExp) * time.Minute).Unix(),
				Issuer:    user,
				IssuedAt:  time.Now().Unix(),
			},
		}
	}
	return &Claims{
		Uid:  uid,
		User: user,
		CN:   cn,
		Mail: mail,
		Type: false,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: t.Add(time.Duration(conf.JwtRef) * time.Minute).Unix(),
			Issuer:    user,
			IssuedAt:  time.Now().Unix(),
		},
	}
}

//同时返回access_token和refresh_token
func CreateToken(uid int64, user, CN, Mail string) (string, string) {
	conf := pkg.Conf()
	claims_accces := NewClaims(uid, 1, user, CN, Mail)
	claims_refresh := NewClaims(uid, 2, user, CN, Mail)
	access := jwt.NewWithClaims(jwt.SigningMethodHS256, claims_accces)
	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, claims_refresh)
	str1, _ := access.SignedString([]byte(conf.JwtSecret))
	str2, _ := refresh.SignedString([]byte(conf.JwtSecret))
	return str1, str2
}
func ParseToken(tokenString string) (*Claims, error) {
	conf := pkg.Conf()
	// 调用解析函数
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 只返回加密的秘钥
		return []byte(conf.JwtSecret), nil
	})
	if err != nil {
		return nil, err
	}
	// 拿到 Claims 类型的结构
	// Claims 接口实现了 Valid() 方法， 对数据进行验证
	claims, ok := token.Claims.(*Claims)
	if ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
