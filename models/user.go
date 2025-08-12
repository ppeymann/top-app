package models

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	otpapp "github.com/ppeymann/top-app.git"
	"gorm.io/gorm"
)

// Errors
var (
	ErrAccountExist     error = errors.New("account with specified params already exists")
	ErrSignInFailed     error = errors.New("account not found or password error")
	ErrPermissionDenied error = errors.New("specified role is not available for user")
	ErrAccountNotExist  error = errors.New("specified account does not exist")
)

type (
	// UserService represents method signatures for api user endpoint.
	// so any object that stratifying this interface can be used as user service for api endpoint.
	UserService interface {
		// Register a new user account.
		Register(ctx *gin.Context, mobile string) *otpapp.BaseResult

		// Login a user account.
		Login(ctx *gin.Context, mobile string) *otpapp.BaseResult

		// OtpVerify verifies the one time password for user account.
		OtpVerify(in *OtpInput, ctx *gin.Context) *otpapp.BaseResult

		// GetUserByPhone
		GetUserByPhone(ctx *gin.Context) *otpapp.BaseResult

		// GetAllUser
		GetAllUser(ctx *gin.Context, page, limit int32) *otpapp.BaseResult
	}

	// UserRepository represents method signatures for user domain repository.
	// so any object that stratifying this interface can be used as user domain repository.
	UserRepository interface {
		Create(mobile string) (*UserEntity, error)
		Find(mobile string) (*UserEntity, error)
		FindByID(id uint) (*UserEntity, error)
		SetOtp(id uint, otp string, expire int64) error
		Update(user *UserEntity) error
		FindAllUser(page, limit int32) ([]UserEntity, error)

		otpapp.BaseRepository
	}

	// UserHandler represents method signatures for user handlers.
	// so any object that stratifying this interface can be used as user handlers.
	UserHandler interface {
		SignUp(ctx *gin.Context)
		SignIn(ctx *gin.Context)
		OtpVerify(ctx *gin.Context)
		GetUser(ctx *gin.Context)
		GetAllUsers(ctx *gin.Context)
	}

	// UserEntity contains user info for stored on database
	//
	// swagger: model UserEntity
	UserEntity struct {
		gorm.Model

		// Mobile
		Mobile string `json:"mobile" gorm:"column:mobile;index;unique"`

		// Verification code for using as one time password for logging in to account
		Verification string `json:"-" gorm:"verification;index"`

		VerificationExpire int64 `json:"-" gorm:"verification_expire;index"`
	}

	OtpInput struct {
		// Mobile is the mobile number of user
		Mobile       string `json:"mobile"`
		Verification string `json:"verification"`
	}

	MobileInput struct {
		Mobile string `json:"mobile"`
	}

	// TokenBundlerOutput
	//
	// swagger: model TokenBundlerOutput
	TokenBundlerOutput struct {
		// Token is string that hashed by paseto
		Token string `json:"token"`

		// Refresh is string that for refresh old token
		Refresh string `json:"refresh"`

		// Expire is time for expire token
		Expire time.Time `json:"expire"`
	}
)

func (a *UserEntity) IsVerificationExpired() bool {
	if time.Unix(a.VerificationExpire, 0).UTC().Before(time.Now().UTC()) {
		return true
	}
	return false
}
