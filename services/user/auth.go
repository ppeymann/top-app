package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	otpapp "github.com/ppeymann/top-app.git"
	"github.com/ppeymann/top-app.git/models"
)

type authService struct {
	next models.UserService
}

// GetAllUser implements models.UserService.
func (a *authService) GetAllUser(ctx *gin.Context, page int32, limit int32) *otpapp.BaseResult {
	_, err := otpapp.CheckAuth(ctx)
	if err != nil {
		return &otpapp.BaseResult{
			Errors: []string{err.Error()},
			Status: http.StatusOK,
		}
	}

	return a.next.GetAllUser(ctx, page, limit)
}

// GetUserByPhone implements models.UserService.
func (a *authService) GetUserByPhone(ctx *gin.Context) *otpapp.BaseResult {
	_, err := otpapp.CheckAuth(ctx)
	if err != nil {
		return &otpapp.BaseResult{
			Errors: []string{err.Error()},
			Status: http.StatusOK,
		}
	}

	return a.next.GetUserByPhone(ctx)
}

// OtpVerify implements models.UserService.
func (a *authService) OtpVerify(in *models.OtpInput, ctx *gin.Context) *otpapp.BaseResult {
	return a.next.OtpVerify(in, ctx)
}

// Register implements models.UserService.
func (a *authService) Register(ctx *gin.Context, mobile string) *otpapp.BaseResult {
	return a.next.Register(ctx, mobile)
}

// Login implements models.UserService.
func (a *authService) Login(ctx *gin.Context, mobile string) *otpapp.BaseResult {
	return a.next.Login(ctx, mobile)
}

func NewAuthService(srv models.UserService) models.UserService {
	return &authService{
		next: srv,
	}
}
