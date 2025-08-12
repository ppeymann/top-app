package user

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-kit/kit/metrics"
	otpapp "github.com/ppeymann/top-app.git"
	"github.com/ppeymann/top-app.git/models"
)

type instrumentingService struct {
	next           models.UserService
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
}

// GetAllUser implements models.UserService.
func (i *instrumentingService) GetAllUser(ctx *gin.Context, page int32, limit int32) *otpapp.BaseResult {
	defer func(begin time.Time) {
		i.requestCount.With("method", "GetAllUser").Add(1)
		i.requestLatency.With("method", "GetAllUser").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return i.next.GetAllUser(ctx, page, limit)
}

// GetUserByPhone implements models.UserService.
func (i *instrumentingService) GetUserByPhone(ctx *gin.Context) *otpapp.BaseResult {
	defer func(begin time.Time) {
		i.requestCount.With("method", "GetUserByPhone").Add(1)
		i.requestLatency.With("method", "GetUserByPhone").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return i.next.GetUserByPhone(ctx)
}

// Login implements models.UserService.
func (i *instrumentingService) Login(ctx *gin.Context, mobile string) *otpapp.BaseResult {
	defer func(begin time.Time) {
		i.requestCount.With("method", "Login").Add(1)
		i.requestLatency.With("method", "Login").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return i.next.Login(ctx, mobile)
}

// OtpVerify implements models.UserService.
func (i *instrumentingService) OtpVerify(in *models.OtpInput, ctx *gin.Context) *otpapp.BaseResult {
	defer func(begin time.Time) {
		i.requestCount.With("method", "OtpVerify").Add(1)
		i.requestLatency.With("method", "OtpVerify").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return i.next.OtpVerify(in, ctx)
}

// Register implements models.UserService.
func (i *instrumentingService) Register(ctx *gin.Context, mobile string) *otpapp.BaseResult {
	defer func(begin time.Time) {
		i.requestCount.With("method", "Register").Add(1)
		i.requestLatency.With("method", "Register").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return i.next.Register(ctx, mobile)
}

func NewInstrumentingService(requestCount metrics.Counter, requestLatency metrics.Histogram, srv models.UserService) models.UserService {
	return &instrumentingService{
		next:           srv,
		requestCount:   requestCount,
		requestLatency: requestLatency,
	}
}
