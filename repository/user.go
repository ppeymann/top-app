package repository

import (
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/ppeymann/top-app.git/models"
	"github.com/ppeymann/top-app.git/utils"
	"gorm.io/gorm"
)

type userRepo struct {
	pg       *gorm.DB
	database string
	table    string
}

// Create implements models.UserRepository.
func (r *userRepo) Create(mobile string) (*models.UserEntity, error) {
	code := utils.RandNumberDigits(6)

	user := &models.UserEntity{
		Model:              gorm.Model{},
		Mobile:             mobile,
		Verification:       code,
		VerificationExpire: time.Now().Add(180 * time.Second).UTC().Unix(),
	}

	err := r.pg.Transaction(func(tx *gorm.DB) error {
		if res := r.Model().Create(user).Error; res != nil {
			str := res.(*pgconn.PgError).Message
			if strings.Contains(str, "duplicate key value") {
				return models.ErrAccountExist
			}

			return res
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}

// Find implements models.UserRepository.
func (r *userRepo) Find(mobile string) (*models.UserEntity, error) {
	user := &models.UserEntity{}
	err := r.Model().Where("mobile = ?", mobile).First(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil

}

// FindAllUser implements models.UserRepository.
func (r *userRepo) FindAllUser(page int32, limit int32) ([]models.UserEntity, error) {
	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 2
	}
	var totalRows int64
	r.Model().Count(&totalRows)

	offset := (page - 1) * limit
	var users []models.UserEntity
	err := r.pg.Order("id ASC").Limit(int(limit)).Offset(int(offset)).Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}

// FindByID implements models.UserRepository.
func (r *userRepo) FindByID(id uint) (*models.UserEntity, error) {
	user := &models.UserEntity{}
	err := r.Model().Where("id = ?", id).First(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

// SetOtp implements models.UserRepository.
func (r *userRepo) SetOtp(id uint, otp string, expire int64) error {
	user, err := r.FindByID(id)
	if err != nil {
		return err
	}

	user.Verification = otp
	user.VerificationExpire = expire

	return r.Update(user)
}

// Update implements models.UserRepository.
func (r *userRepo) Update(user *models.UserEntity) error {
	return r.pg.Save(user).Error
}

// Migrate implements models.UserRepository.
func (r *userRepo) Migrate() error {
	return r.pg.AutoMigrate(models.UserEntity{})
}

// Model implements models.UserRepository.
func (r *userRepo) Model() *gorm.DB {
	return r.pg.Model(&models.UserEntity{})
}

// Name implements models.UserRepository.
func (r *userRepo) Name() string {
	return r.table
}

func NewUserRepo(pg *gorm.DB, database string) models.UserRepository {
	return &userRepo{
		pg:       pg,
		database: database,
		table:    "user_entities",
	}
}
