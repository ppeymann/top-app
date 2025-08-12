package user

import (
	"github.com/gin-gonic/gin"
	otpapp "github.com/ppeymann/top-app.git"
	"github.com/ppeymann/top-app.git/models"
	validations "github.com/ppeymann/top-app.git/validation"
)

type validationService struct {
	next   models.UserService
	schema map[string][]byte
}

// GetAllUser implements models.UserService.
func (v *validationService) GetAllUser(ctx *gin.Context, page int32, limit int32) *otpapp.BaseResult {
	return v.next.GetAllUser(ctx, page, limit)
}

// GetUserByPhone implements models.UserService.
func (v *validationService) GetUserByPhone(ctx *gin.Context) *otpapp.BaseResult {
	return v.next.GetUserByPhone(ctx)
}

// Login implements models.UserService.
func (v *validationService) Login(ctx *gin.Context, mobile string) *otpapp.BaseResult {
	return v.next.Login(ctx, mobile)
}

// OtpVerify implements models.UserService.
func (v *validationService) OtpVerify(in *models.OtpInput, ctx *gin.Context) *otpapp.BaseResult {
	return v.next.OtpVerify(in, ctx)
}

// Register implements models.UserService.
func (v *validationService) Register(ctx *gin.Context, mobile string) *otpapp.BaseResult {
	return v.next.Register(ctx, mobile)
}

func NewValidationService(srv models.UserService, path string) (models.UserService, error) {
	schema := make(map[string][]byte)

	// Load the schema from the specified path
	err := validations.LoadSchema(path, schema)
	if err != nil {
		return nil, err
	}

	return &validationService{
		next:   srv,
		schema: schema,
	}, nil
}
