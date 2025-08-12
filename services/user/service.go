package user

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	otpapp "github.com/ppeymann/top-app.git"
	"github.com/ppeymann/top-app.git/auth"
	"github.com/ppeymann/top-app.git/config"
	"github.com/ppeymann/top-app.git/env"
	"github.com/ppeymann/top-app.git/models"
	"github.com/ppeymann/top-app.git/utils"
)

type service struct {
	repo models.UserRepository
	conf *config.Configuration
}

// GetAllUser implements models.UserService.
func (s *service) GetAllUser(ctx *gin.Context, page int32, limit int32) *otpapp.BaseResult {
	users, err := s.repo.FindAllUser(page, limit)
	if err != nil {
		return &otpapp.BaseResult{
			Errors: []string{err.Error()},
			Status: http.StatusOK,
		}
	}

	return &otpapp.BaseResult{
		Status:      http.StatusOK,
		Result:      users,
		ResultCount: int64(len(users)),
	}
}

// GetUserByPhone implements models.UserService.
func (s *service) GetUserByPhone(ctx *gin.Context) *otpapp.BaseResult {
	claims, _ := otpapp.CheckAuth(ctx)
	user, err := s.repo.FindByID(claims.Subject)
	if err != nil {
		return &otpapp.BaseResult{
			Errors: []string{err.Error()},
			Status: http.StatusOK,
		}
	}

	return &otpapp.BaseResult{
		Status: http.StatusOK,
		Result: user,
	}
}

// Login implements models.UserService.
func (s *service) Login(ctx *gin.Context, mobile string) *otpapp.BaseResult {
	user := &models.UserEntity{}
	var err error

	user, err = s.repo.Find(mobile)
	if err != nil {
		return &otpapp.BaseResult{
			Status: http.StatusOK,
			Errors: []string{err.Error()},
		}
	}

	if user.Verification != "" && !user.IsVerificationExpired() {
		return &otpapp.BaseResult{
			Status: http.StatusOK,
			Errors: []string{"previous otp not expired, please wait a few minutes"},
		}
	}

	code := utils.RandNumberDigits(6)
	exp := time.Now().Add(180 * time.Second).UTC().Unix()

	err = s.repo.SetOtp(user.ID, code, exp)
	if err != nil {
		return &otpapp.BaseResult{
			Status: http.StatusOK,
			Errors: []string{err.Error()},
		}
	}

	return &otpapp.BaseResult{
		Status: http.StatusOK,
		Result: fmt.Sprintf("%s for : %s", code, mobile),
	}
}

// OtpVerify implements models.UserService.
func (s *service) OtpVerify(in *models.OtpInput, ctx *gin.Context) *otpapp.BaseResult {
	user := &models.UserEntity{}
	var err error

	user, err = s.repo.Find(in.Mobile)
	if err != nil {
		return &otpapp.BaseResult{
			Status: http.StatusOK,
			Errors: []string{err.Error()},
		}
	}

	if in.Verification != user.Verification {
		return &otpapp.BaseResult{
			Status: http.StatusOK,
			Errors: []string{"OTP is not correct"},
		}
	}

	if user.IsVerificationExpired() {
		return &otpapp.BaseResult{
			Status: http.StatusOK,
			Errors: []string{"OTP Expired"},
		}
	}

	user.Verification = ""
	user.VerificationExpire = time.Now().UTC().Unix()

	err = s.repo.Update(user)
	if err != nil {
		return &otpapp.BaseResult{
			Status: http.StatusOK,
			Errors: []string{err.Error()},
		}
	}

	paseto, err := auth.NewPasetoMaker(env.GetEnv("JWT", ""))
	if err != nil {
		return &otpapp.BaseResult{
			Status: http.StatusOK,
			Errors: []string{otpapp.ErrInternalServer.Error()},
		}
	}

	tokenClaims := &auth.Claims{
		Subject:   user.ID,
		Issuer:    s.conf.Jwt.Issuer,
		Audience:  s.conf.Jwt.Audience,
		IssuedAt:  time.Unix(s.conf.Jwt.RefreshExpire, 0),
		ExpiredAt: time.Now().Add(time.Duration(s.conf.Jwt.TokenExpire) * time.Minute).UTC(),
	}

	tokenStr, err := paseto.CreateToken(tokenClaims)
	if err != nil {
		return &otpapp.BaseResult{
			Status: http.StatusOK,
			Errors: []string{otpapp.ErrInternalServer.Error()},
		}
	}

	return &otpapp.BaseResult{
		Status: http.StatusOK,
		Result: models.TokenBundlerOutput{
			Token:  tokenStr,
			Expire: tokenClaims.ExpiredAt,
		},
	}
}

// Register implements models.UserService.
func (s *service) Register(ctx *gin.Context, mobile string) *otpapp.BaseResult {
	user, err := s.repo.Create(mobile)
	if err != nil {
		return &otpapp.BaseResult{
			Status: http.StatusOK,
			Errors: []string{err.Error()},
		}
	}

	return &otpapp.BaseResult{
		Status: http.StatusOK,
		Result: fmt.Sprintf("%s for : %s", user.Verification, mobile),
	}
}

func NewService(repo models.UserRepository, conf *config.Configuration) models.UserService {
	return &service{
		repo: repo,
		conf: conf,
	}
}
