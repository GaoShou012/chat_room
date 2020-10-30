package tenant_api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
)

type Operator struct {
	TenantCode string
}

func SetOperator(ctx *gin.Context, op *Operator) {
	ctx.Set("operator", op)
}
func GetOperator(ctx *gin.Context) (*Operator, error) {
	val, exists := ctx.Get("operator")
	if !exists {
		return nil, errors.New("Lose the operator info\n")
	}
	operator, ok := val.(*Operator)
	if !ok {
		return nil, errors.New("Assert operator type failed\n")
	}
	return operator, nil
}

func (o *Operator) EncryptByJwt(key []byte) (string, error) {
	m := jwt.MapClaims{}
	data, err := json.Marshal(o)
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(data, &m)
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, m)
	return token.SignedString(key)
}
func (o *Operator) DecryptByJwt(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v\n", token.Header["alg"])
		}
		return []byte(TokenKey), nil
	})
	if err != nil {
		return err
	}

	err = mapstructure.WeakDecode(token.Claims.(jwt.MapClaims), o)
	if err != nil {
		return err
	}

	return nil
}
