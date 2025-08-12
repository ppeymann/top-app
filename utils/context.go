package utils

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"github.com/ppeymann/top-app.git/auth"
)

const (
	ContextUserKey = "CONTEXT_USER"
	ContextRoleKey = "CONTEXT_ROLE"
)

// ErrUserPrincipalsNotFound is returned when UserPrincipals are not found in the context.
var ErrUserPrincipalsNotFound = errors.New("UserPrincipals not found in context")

func CatchClaims(ctx *gin.Context, claims *auth.Claims) error {
	user, ok := ctx.Get(ContextUserKey)
	if !ok {
		return errors.New("user not found in context")
	}

	err := mapstructure.Decode(user, claims)
	if err != nil {
		return errors.New("error parsing user object to claims")
	}

	return nil
}
