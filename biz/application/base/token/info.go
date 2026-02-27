package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/xh-polaris/synapse4b/biz/conf"
	"github.com/xh-polaris/synapse4b/biz/types/cst"
)

type Info struct {
	RawToken    string
	BasicUserId string         `json:"basic_user_id"` // 用户ID
	UnitId      string         `json:"unit_id"`       // 学校ID
	Code        string         `json:"code"`          // 学号
	Phone       string         `json:"phone"`         // 手机号
	Email       string         `json:"email"`         // 邮箱
	LoginTime   int64          `json:"login_time"`    // 登录时间(秒时间戳)
	AuthType    string         `json:"auth_type"`     // 登录类型
	Extra       map[string]any `json:"extra"`         // 额外信息
}

func SignJWT(tokenConf *conf.Token, info *Info) (string, error) {
	now := time.Now().UTC()

	// 将ino中信息签名
	claims := jwt.MapClaims{
		"iat":                now.Unix(),                                                    // sign time
		"nbf":                now.Unix(),                                                    // validate time
		"exp":                now.Add(time.Duration(tokenConf.Expire) * time.Second).Unix(), // expire time
		cst.TokenBasicUserID: info.BasicUserId,
		cst.TokenUnitID:      info.UnitId,
		cst.TokenCode:        info.Code,
		cst.TokenPhone:       info.Phone,
		cst.TokenEmail:       info.Email,
		cst.TokenLoginTime:   info.LoginTime,
		cst.TokenAuthType:    info.AuthType,
		cst.TokenExtra:       info.Extra,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	sk, err := conf.GetSecretKey(tokenConf.SecretKey)
	if err != nil {
		return "", err
	}
	str, err := token.SignedString(sk)
	return str, err
}

func ParseJWT(tokenConf *conf.Token, str string) (*Info, error) {
	token, err := jwt.Parse(str, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return conf.GetPublicKey(tokenConf.PublicKey)
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		info := &Info{
			BasicUserId: claims[cst.TokenBasicUserID].(string),
			UnitId:      claims[cst.TokenUnitID].(string),
			Code:        claims[cst.TokenCode].(string),
			Phone:       claims[cst.TokenPhone].(string),
			Email:       claims[cst.TokenEmail].(string),
			LoginTime:   int64(claims[cst.TokenLoginTime].(float64)),
			AuthType:    claims[cst.TokenAuthType].(string),
			Extra:       claims[cst.TokenExtra].(map[string]any),
			RawToken:    str,
		}
		return info, nil
	}
	return nil, errors.New("invalid token")
}
