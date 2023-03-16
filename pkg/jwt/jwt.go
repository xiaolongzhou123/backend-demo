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
	jwt.StandardClaims
}

func NewClaims(uid int64, user, cn, mail string) *Claims {
	conf := pkg.Conf()
	t := time.Now()
	return &Claims{
		Uid:  uid,
		User: user,
		CN:   cn,
		Mail: mail,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: t.Add(time.Duration(conf.JwtExp) * time.Minute).Unix(),
			Issuer:    user,
			IssuedAt:  time.Now().Unix(),
		},
	}
}

//create
func CreateToken(user, CN, Mail string) (string, error) {
	conf := pkg.Conf()
	c := NewClaims(1, user, CN, Mail)
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return claims.SignedString([]byte(conf.JwtSecret))
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
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
