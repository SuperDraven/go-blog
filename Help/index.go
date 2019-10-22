package Help

import (
	"blog/conf"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type UserInfo struct {
	Name     string
	Password string
}

func Md5Encryption(password string) string {
	Md5Inst := md5.New()
	Md5Inst.Write([]byte(password))
	//Result := Md5Inst.Sum([]byte(""))
	return hex.EncodeToString(Md5Inst.Sum(nil))
}

func CreateToken(user *UserInfo) (tokenss string, err error) {
	//自定义claim

	claim := jwt.MapClaims{
		"name":     user.Name,
		"password": user.Password,
		"nbf":      time.Now().Unix(),
		"iat":      time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenss, err = token.SignedString([]byte(conf.LoadConf().JwtSecret))
	return
}

func secret() jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		return []byte(conf.LoadConf().JwtSecret), nil
	}
}

func ParseToken(tokenss string) (user *UserInfo, err error) {
	user = &UserInfo{}
	token, err := jwt.Parse(tokenss, secret())
	if err != nil {
		return
	}
	claim, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		err = errors.New("cannot convert claim to mapclaim")
		return
	}
	//验证token，如果token被修改过则为false
	if !token.Valid {
		err = errors.New("token is invalid")
		return
	}

	user.Name = claim["name"].(string)
	user.Password = claim["password"].(string)
	return
}
