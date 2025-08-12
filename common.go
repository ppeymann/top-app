package otpapp

import (
	"encoding/json"
	"errors"

	"gorm.io/gorm"
)

type (
	AccountRole string

	// BaseRepository is abstract interface that all repositories must implement its methods
	BaseRepository interface {
		// Migrate runs AutoMigrate for expected repository model
		Migrate() error

		// Name repository associated table name
		Name() string

		// Model returns *gorm.DB instance for repository
		Model() *gorm.DB
	}

	// BaseResult a basic GoLang struct which includes the following fields: Success, Errors, Messages, ResultCount, Result
	// It is the unified response model for entire service api calls
	//
	// swagger:model BaseResult
	BaseResult struct {
		Status int `json:"-"`

		// Errors provides list off error that occurred in processing request
		Errors []string `json:"errors" mapstructure:"errors"`

		// ResultCount specified number of records that returned in result_count field expected result been array.
		ResultCount int64 `json:"result_count,omitempty" mapstructure:"result_count"`

		// Result single/array of any type (object/number/string/boolean) that returns as response
		Result interface{} `json:"result" mapstructure:"result"`
	}

	Error struct {
		Code        int    `json:"code" mapstructure:"code"`
		Description string `json:"description" mapstructure:"description"`
	}

	SearchInputDTO struct {
		Query string `json:"query"`
	}

	// ContextUser a basic GoLang struct which includes the following fields: ID, Roles, Permissions
	// It used as User object that holds required information to identifying user and relegated roles and permissions
	// It may be embedded into any Input model, or you may build your own model without it
	ContextUser struct {
		ID          string        `json:"id" mapstructure:"id"`
		Roles       []AccountRole `json:"roles" mapstructure:"roles"`
		Permissions []string      `json:"permissions" mapstructure:"permissions"`
	}

	BaseDocument struct {
		Id        uint  `json:"id" bson:"_id"`
		CreatedAt int64 `json:"created_at" bson:"created_at"`
		UpdatedAt int64 `json:"updated_at" bson:"updated_at"`
	}
)

var (
	ErrUnimplementedRequest = errors.New("request is not implemented")
	ErrUnhandled            = errors.New("an unhandled error occurred during processing the request")
	ErrNotFound             = errors.New("not found")
	ErrInternalServer       = errors.New("internal server error")
	ErrEntityAlreadyExist   = errors.New("entity with specified properties already exist")
	ErrUnAuthorization      = errors.New("UnAuthorization Error")
)

const (
	ContextUserKey          string = "CONTEXT_USER"
	UserSessionKey          string = "USER_SESSION"
	AuthorizationFailed     string = "authorization failed"
	ProvideRequiredParam    string = "please provide required params"
	ProvideRequiredJsonBody string = "please provide required JSON body"
)

type (
	DescribedFile struct {
		File  string `json:"file" bson:"file" mapstructure:"file"`
		Title string `json:"title" bson:"title" mapstructure:"title"`
	}

	SearchQuery struct {
		Query  string `json:"query"`
		Size   int    `json:"size"`
		Offset int    `json:"offset"`
	}
)

const (
	UserRole   string = "USER"
	AdminRole  string = "ADMIN"
	SellerRole string = "SELLER"
)

var AllRoles = []string{UserRole, AdminRole, SellerRole}
var SellerRoles = []string{UserRole, SellerRole}

// ToJson is method for parsing ContextUser to json string.
func (p *ContextUser) ToJson() (string, error) {
	js, err := json.Marshal(p)
	return string(js), err
}

// FromJson is method for parsing ContextUser from json string.
func (p *ContextUser) FromJson(val string) error {
	return json.Unmarshal([]byte(val), p)
}
