package data

import (
	"context"
	"time"

	"gorm.io/gorm"

	"StructForge/backend/common/log"
)

// User 用户模型
type User struct {
	ID              int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	Username        string     `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Email           string     `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	PasswordHash    string     `gorm:"type:varchar(255);not null;column:password_hash" json:"-"`
	EmailVerified   bool       `gorm:"default:false;column:email_verified" json:"email_verified"`
	EmailVerifiedAt *time.Time `gorm:"column:email_verified_at" json:"email_verified_at,omitempty"`
	Status          string     `gorm:"type:varchar(20);default:'active'" json:"status"`
	CreatedAt       time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	LastLoginAt     *time.Time `gorm:"column:last_login_at" json:"last_login_at,omitempty"`
	LastLoginIP     string     `gorm:"type:varchar(45);column:last_login_ip" json:"last_login_ip,omitempty"`

	// 关联
	Profile *UserProfile `gorm:"foreignKey:UserID" json:"profile,omitempty"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// UserProfile 用户资料模型
type UserProfile struct {
	ID        int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int64      `gorm:"uniqueIndex;not null;column:user_id" json:"user_id"`
	Nickname  string     `gorm:"type:varchar(100)" json:"nickname,omitempty"`
	AvatarURL string     `gorm:"type:varchar(500);column:avatar_url" json:"avatar_url,omitempty"`
	Bio       string     `gorm:"type:text" json:"bio,omitempty"`
	Phone     string     `gorm:"type:varchar(20)" json:"phone,omitempty"`
	Gender    string     `gorm:"type:varchar(10)" json:"gender,omitempty"`
	Birthday  *time.Time `gorm:"type:date" json:"birthday,omitempty"`
	Location  string     `gorm:"type:varchar(100)" json:"location,omitempty"`
	Website   string     `gorm:"type:varchar(255)" json:"website,omitempty"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	// 关联
	User *User `gorm:"foreignKey:UserID" json:"-"`
}

// TableName 指定表名
func (UserProfile) TableName() string {
	return "user_profiles"
}

// EmailVerification 邮箱验证模型
type EmailVerification struct {
	ID        int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int64      `gorm:"not null;column:user_id" json:"user_id"`
	Email     string     `gorm:"type:varchar(255);not null" json:"email"`
	Token     string     `gorm:"type:varchar(100);uniqueIndex;not null" json:"token"`
	Type      string     `gorm:"type:varchar(20);not null" json:"type"` // register, reset_password, change_email
	ExpiresAt time.Time  `gorm:"not null;column:expires_at" json:"expires_at"`
	Used      bool       `gorm:"default:false" json:"used"`
	UsedAt    *time.Time `gorm:"column:used_at" json:"used_at,omitempty"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
}

// TableName 指定表名
func (EmailVerification) TableName() string {
	return "email_verifications"
}

// UserRepo 用户数据访问接口
type UserRepo interface {
	// 创建用户
	CreateUser(ctx context.Context, user *User) error

	// 根据ID查询用户
	GetUserByID(ctx context.Context, id int64) (*User, error)

	// 根据用户名查询用户
	GetUserByUsername(ctx context.Context, username string) (*User, error)

	// 根据邮箱查询用户
	GetUserByEmail(ctx context.Context, email string) (*User, error)

	// 更新用户
	UpdateUser(ctx context.Context, user *User) error

	// 更新最后登录信息
	UpdateLastLogin(ctx context.Context, userID int64, ip string) error

	// 检查用户名是否存在
	UsernameExists(ctx context.Context, username string) (bool, error)

	// 检查邮箱是否存在
	EmailExists(ctx context.Context, email string) (bool, error)
}

// UserProfileRepo 用户资料数据访问接口
type UserProfileRepo interface {
	// 创建用户资料
	CreateProfile(ctx context.Context, profile *UserProfile) error

	// 根据用户ID查询资料
	GetProfileByUserID(ctx context.Context, userID int64) (*UserProfile, error)

	// 更新用户资料
	UpdateProfile(ctx context.Context, profile *UserProfile) error
}

// EmailVerificationRepo 邮箱验证数据访问接口
type EmailVerificationRepo interface {
	// 创建验证记录
	CreateVerification(ctx context.Context, verification *EmailVerification) error

	// 根据令牌查询验证记录
	GetVerificationByToken(ctx context.Context, token string) (*EmailVerification, error)

	// 标记验证记录为已使用
	MarkAsUsed(ctx context.Context, id int64) error

	// 删除过期验证记录
	DeleteExpired(ctx context.Context) error
}

// userRepo 用户数据访问实现
type userRepo struct {
	data *Data
	db   *gorm.DB
}

// NewUserRepo 创建用户数据访问实例
func NewUserRepo(data *Data) UserRepo {
	return &userRepo{
		data: data,
		db:   data.DB(),
	}
}

// CreateUser 创建用户
func (r *userRepo) CreateUser(ctx context.Context, user *User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		log.Error(ctx, "创建用户失败",
			log.ErrorField(err),
			log.String("username", user.Username),
			log.String("email", user.Email),
		)
		return err
	}
	return nil
}

// GetUserByID 根据ID查询用户
func (r *userRepo) GetUserByID(ctx context.Context, id int64) (*User, error) {
	var user User
	if err := r.db.WithContext(ctx).
		Preload("Profile").
		Where("id = ?", id).
		First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		log.Error(ctx, "查询用户失败",
			log.ErrorField(err),
			log.Int64("user_id", id),
		)
		return nil, err
	}
	return &user, nil
}

// GetUserByUsername 根据用户名查询用户
func (r *userRepo) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	var user User
	if err := r.db.WithContext(ctx).
		Preload("Profile").
		Where("username = ?", username).
		First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		log.Error(ctx, "查询用户失败",
			log.ErrorField(err),
			log.String("username", username),
		)
		return nil, err
	}
	return &user, nil
}

// GetUserByEmail 根据邮箱查询用户
func (r *userRepo) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	if err := r.db.WithContext(ctx).
		Preload("Profile").
		Where("email = ?", email).
		First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		log.Error(ctx, "查询用户失败",
			log.ErrorField(err),
			log.String("email", email),
		)
		return nil, err
	}
	return &user, nil
}

// UpdateUser 更新用户
func (r *userRepo) UpdateUser(ctx context.Context, user *User) error {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		log.Error(ctx, "更新用户失败",
			log.ErrorField(err),
			log.Int64("user_id", user.ID),
		)
		return err
	}
	return nil
}

// UpdateLastLogin 更新最后登录信息
func (r *userRepo) UpdateLastLogin(ctx context.Context, userID int64, ip string) error {
	now := time.Now()
	if err := r.db.WithContext(ctx).
		Model(&User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"last_login_at": now,
			"last_login_ip": ip,
		}).Error; err != nil {
		log.Error(ctx, "更新最后登录信息失败",
			log.ErrorField(err),
			log.Int64("user_id", userID),
		)
		return err
	}
	return nil
}

// UsernameExists 检查用户名是否存在
func (r *userRepo) UsernameExists(ctx context.Context, username string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&User{}).
		Where("username = ?", username).
		Count(&count).Error; err != nil {
		log.Error(ctx, "检查用户名是否存在失败",
			log.ErrorField(err),
			log.String("username", username),
		)
		return false, err
	}
	return count > 0, nil
}

// EmailExists 检查邮箱是否存在
func (r *userRepo) EmailExists(ctx context.Context, email string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&User{}).
		Where("email = ?", email).
		Count(&count).Error; err != nil {
		log.Error(ctx, "检查邮箱是否存在失败",
			log.ErrorField(err),
			log.String("email", email),
		)
		return false, err
	}
	return count > 0, nil
}

// userProfileRepo 用户资料数据访问实现
type userProfileRepo struct {
	data *Data
	db   *gorm.DB
}

// NewUserProfileRepo 创建用户资料数据访问实例
func NewUserProfileRepo(data *Data) UserProfileRepo {
	return &userProfileRepo{
		data: data,
		db:   data.DB(),
	}
}

// CreateProfile 创建用户资料
func (r *userProfileRepo) CreateProfile(ctx context.Context, profile *UserProfile) error {
	if err := r.db.WithContext(ctx).Create(profile).Error; err != nil {
		log.Error(ctx, "创建用户资料失败",
			log.ErrorField(err),
			log.Int64("user_id", profile.UserID),
		)
		return err
	}
	return nil
}

// GetProfileByUserID 根据用户ID查询资料
func (r *userProfileRepo) GetProfileByUserID(ctx context.Context, userID int64) (*UserProfile, error) {
	var profile UserProfile
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		First(&profile).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		log.Error(ctx, "查询用户资料失败",
			log.ErrorField(err),
			log.Int64("user_id", userID),
		)
		return nil, err
	}
	return &profile, nil
}

// UpdateProfile 更新用户资料
func (r *userProfileRepo) UpdateProfile(ctx context.Context, profile *UserProfile) error {
	if err := r.db.WithContext(ctx).Save(profile).Error; err != nil {
		log.Error(ctx, "更新用户资料失败",
			log.ErrorField(err),
			log.Int64("user_id", profile.UserID),
		)
		return err
	}
	return nil
}

// emailVerificationRepo 邮箱验证数据访问实现
type emailVerificationRepo struct {
	data *Data
	db   *gorm.DB
}

// NewEmailVerificationRepo 创建邮箱验证数据访问实例
func NewEmailVerificationRepo(data *Data) EmailVerificationRepo {
	return &emailVerificationRepo{
		data: data,
		db:   data.DB(),
	}
}

// CreateVerification 创建验证记录
func (r *emailVerificationRepo) CreateVerification(ctx context.Context, verification *EmailVerification) error {
	if err := r.db.WithContext(ctx).Create(verification).Error; err != nil {
		log.Error(ctx, "创建邮箱验证记录失败",
			log.ErrorField(err),
			log.Int64("user_id", verification.UserID),
		)
		return err
	}
	return nil
}

// GetVerificationByToken 根据令牌查询验证记录
func (r *emailVerificationRepo) GetVerificationByToken(ctx context.Context, token string) (*EmailVerification, error) {
	var verification EmailVerification
	if err := r.db.WithContext(ctx).
		Where("token = ? AND used = false AND expires_at > ?", token, time.Now()).
		First(&verification).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		log.Error(ctx, "查询邮箱验证记录失败",
			log.ErrorField(err),
			log.String("token", token),
		)
		return nil, err
	}
	return &verification, nil
}

// MarkAsUsed 标记验证记录为已使用
func (r *emailVerificationRepo) MarkAsUsed(ctx context.Context, id int64) error {
	now := time.Now()
	if err := r.db.WithContext(ctx).
		Model(&EmailVerification{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"used":    true,
			"used_at": now,
		}).Error; err != nil {
		log.Error(ctx, "标记邮箱验证记录为已使用失败",
			log.ErrorField(err),
			log.Int64("id", id),
		)
		return err
	}
	return nil
}

// DeleteExpired 删除过期验证记录
func (r *emailVerificationRepo) DeleteExpired(ctx context.Context) error {
	if err := r.db.WithContext(ctx).
		Where("expires_at < ?", time.Now()).
		Delete(&EmailVerification{}).Error; err != nil {
		log.Error(ctx, "删除过期邮箱验证记录失败",
			log.ErrorField(err),
		)
		return err
	}
	return nil
}
