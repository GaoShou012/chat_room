package tenant_api

import (
	"errors"
	"github.com/gin-gonic/gin"
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