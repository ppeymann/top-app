package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	otpapp "github.com/ppeymann/top-app.git"
	"github.com/ppeymann/top-app.git/models"
	"github.com/ppeymann/top-app.git/server"
)

type handler struct {
	next models.UserService
}

// SignUp is handler for create New user
//
// @BasePath			 		/api/v1/user
// @Summary		 				Create New user
// @Description		 			Create New user
// @Tags 						user
// @Accept						json
// @Produce 					json
//
// @Param						input body models.MobileInput true "MobileInput"
// @Success 					200 {object} otpapp.BaseResult{result=models.TokenBundlerOutput}	"always return status 200 but body contains error"
// @Router						/api/v1/user/signup	[post]
func (h *handler) SignUp(ctx *gin.Context) {
	in := &models.MobileInput{}

	if err := ctx.ShouldBindJSON(in); err != nil {
		ctx.JSON(http.StatusBadRequest, otpapp.BaseResult{
			Errors: []string{otpapp.ProvideRequiredJsonBody},
		})

		return
	}

	result := h.next.Register(ctx, in.Mobile)
	ctx.JSON(result.Status, result)
}

// Login is handler for log in
//
// @BasePath			/api/v1/user
// @Summary				log in
// @Description			log in with specific mobile number
// @Tags				user
// @Accept				json
// @Produce				json
//
// @Params				input body models.MobileInput	true	"MobileInput"
// @Success				200 {object} vendora.BaseResult{result=models.TokenBundlerOutput} 	"always return status 200 but body contains error"
// @Router				/api/v1/user/login	[post]
func (h *handler) SignIn(ctx *gin.Context) {
	in := &models.MobileInput{}

	if err := ctx.ShouldBindJSON(in); err != nil {
		ctx.JSON(http.StatusBadRequest, otpapp.BaseResult{
			Errors: []string{otpapp.ProvideRequiredJsonBody},
		})

		return
	}

	result := h.next.Login(ctx, in.Mobile)
	ctx.JSON(result.Status, result)
}

// GetUser is handler for get information
//
// @BasePath			/api/v1/user
// @Summary				user info
// @Description			get user information
// @Tags				user
// @Accept				json
// @Produce				json
//
// @Success				200	{object}	otpapp.BaseResult{result=models.UserEntity}	"always return status 200 but body contains error"
// @Router				/api/v1/user	[get]
func (h *handler) GetUser(ctx *gin.Context) {
	result := h.next.GetUserByPhone(ctx)
	ctx.JSON(result.Status, result)
}

// OtpVerify is handler for verification one time password
//
// @BasePath			/api/v1/user
// @Summary				otp verification
// @Tags				user
// @Accept				json
// @Product				json
//
// @Params				input body models.OtpInput	true	"OtpInput"
// @Success				200	{object}	otpapp.BaseResult{result=models.UserEntity}
// @Router				/api/v1/user/otp	[post]
func (h *handler) OtpVerify(ctx *gin.Context) {
	in := &models.OtpInput{}

	if err := ctx.ShouldBindJSON(in); err != nil {
		ctx.JSON(http.StatusBadRequest, otpapp.BaseResult{
			Errors: []string{otpapp.ProvideRequiredJsonBody},
		})

		return
	}

	result := h.next.OtpVerify(in, ctx)
	ctx.JSON(result.Status, result)
}

// GetAllUsers is handler for get all user
//
// @BasePath			/api/v1/user
// @Summary				get all user
// @Description			get all user
// @Tags				user
// @Accept				json
// @Product				json
//
// @Success				200	{object}	otpapp.BaseResult{result=[]models.UserEntity}
// @Router				/api/v1/user/{offset}/{page}	[get]
// @Security			bearer
func (h *handler) GetAllUsers(ctx *gin.Context) {
	offset := server.GetPathOffset(ctx)
	page, err := server.GetInt64Path("page", ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, otpapp.BaseResult{
			Errors: []string{otpapp.ProvideRequiredParam},
		})

		return
	}

	result := h.next.GetAllUser(ctx, int32(page), int32(offset))
	ctx.JSON(result.Status, result)
}

func NewHandler(srv models.UserService, s *server.Server) models.UserHandler {
	handler := &handler{
		next: srv,
	}

	group := s.Router.Group("/api/v1/user")
	{
		group.POST("/signup", handler.SignUp)
		group.POST("/signin", handler.SignIn)
		group.POST("/otp", handler.OtpVerify)
	}

	group.Use(s.Authenticate())
	{
		group.GET("/:offset/:page", handler.GetAllUsers)
		group.GET("/", handler.GetUser)
	}

	return handler
}
