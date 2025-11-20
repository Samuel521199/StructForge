package biz

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"StructForge/backend/apps/user/internal/data"
	"StructForge/backend/common/email"
	"StructForge/backend/common/log"
)

var (
	ErrUserNotFound             = errors.New("用户不存在")
	ErrUserExists               = errors.New("用户已存在")
	ErrInvalidPassword          = errors.New("密码错误")
	ErrInvalidPasswordFormat    = errors.New("密码格式不符合要求：必须包含字母和数字，且包含大小写字母或特殊字符")
	ErrInvalidVerificationToken = errors.New("验证令牌无效或已过期")
	ErrEmailNotVerified         = errors.New("邮箱未验证")
	ErrUserBanned               = errors.New("用户已被封禁")
	ErrUserInactive             = errors.New("用户未激活")
)

// UserUseCase 用户业务逻辑
type UserUseCase struct {
	userRepo              data.UserRepo
	userProfileRepo       data.UserProfileRepo
	emailVerificationRepo data.EmailVerificationRepo
	emailService          email.EmailService
}

// NewUserUseCase 创建用户业务逻辑实例
func NewUserUseCase(
	userRepo data.UserRepo,
	userProfileRepo data.UserProfileRepo,
	emailVerificationRepo data.EmailVerificationRepo,
	emailService email.EmailService,
) *UserUseCase {
	return &UserUseCase{
		userRepo:              userRepo,
		userProfileRepo:       userProfileRepo,
		emailVerificationRepo: emailVerificationRepo,
		emailService:          emailService,
	}
}

// hashPassword 加密密码
func (uc *UserUseCase) hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// comparePassword 比较密码
func (uc *UserUseCase) comparePassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// validatePassword 验证密码复杂度
// 要求：必须包含字母+数字，且有大小写或特殊字符
func (uc *UserUseCase) validatePassword(password string) error {
	if len(password) < 6 || len(password) > 20 {
		return errors.New("密码长度必须在 6-20 个字符之间")
	}

	// 检查是否包含字母和数字
	hasLetter := false
	hasNumber := false
	hasUpperCase := false
	hasLowerCase := false
	hasSpecialChar := false

	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasLetter = true
			hasUpperCase = true
		case char >= 'a' && char <= 'z':
			hasLetter = true
			hasLowerCase = true
		case char >= '0' && char <= '9':
			hasNumber = true
		case char == '!' || char == '@' || char == '#' || char == '$' || char == '%' ||
			char == '^' || char == '&' || char == '*' || char == '(' || char == ')' ||
			char == '_' || char == '+' || char == '-' || char == '=' || char == '[' ||
			char == ']' || char == '{' || char == '}' || char == ';' || char == ':' ||
			char == '"' || char == '\\' || char == '|' || char == ',' || char == '.' ||
			char == '<' || char == '>' || char == '/' || char == '?':
			hasSpecialChar = true
		}
	}

	// 必须包含字母和数字
	if !hasLetter || !hasNumber {
		return errors.New("密码必须包含字母和数字")
	}

	// 必须有大小写或特殊字符
	if !hasUpperCase && !hasLowerCase && !hasSpecialChar {
		return errors.New("密码必须包含大小写字母或特殊字符")
	}

	// 如果只有小写字母，必须有特殊字符
	if hasLowerCase && !hasUpperCase && !hasSpecialChar {
		return errors.New("密码必须包含大写字母或特殊字符")
	}

	// 如果只有大写字母，必须有特殊字符
	if hasUpperCase && !hasLowerCase && !hasSpecialChar {
		return errors.New("密码必须包含小写字母或特殊字符")
	}

	return nil
}

// generateToken 生成验证令牌
func (uc *UserUseCase) generateToken() string {
	// 使用时间戳 + UUID 生成唯一令牌
	uuid := generateUUID()
	return time.Now().Format("20060102150405") + "-" + uuid
}

// Register 用户注册
func (uc *UserUseCase) Register(ctx context.Context, username, email, password string) (*data.User, error) {
	// 验证密码复杂度
	if err := uc.validatePassword(password); err != nil {
		return nil, err
	}

	// 检查用户名是否存在
	exists, err := uc.userRepo.UsernameExists(ctx, username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrUserExists
	}

	// 检查邮箱是否存在
	exists, err = uc.userRepo.EmailExists(ctx, email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrUserExists
	}

	// 加密密码
	passwordHash, err := uc.hashPassword(password)
	if err != nil {
		log.Error(ctx, "密码加密失败",
			log.ErrorField(err),
		)
		return nil, err
	}

	// 创建用户
	user := &data.User{
		Username:      username,
		Email:         email,
		PasswordHash:  passwordHash,
		EmailVerified: false,
		Status:        "inactive", // 注册后需要验证邮箱才能激活
	}

	if err := uc.userRepo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	// 创建用户资料
	profile := &data.UserProfile{
		UserID:   user.ID,
		Nickname: username, // 默认昵称为用户名
	}

	if err := uc.userProfileRepo.CreateProfile(ctx, profile); err != nil {
		log.Warn(ctx, "创建用户资料失败",
			log.ErrorField(err),
			log.Int64("user_id", user.ID),
		)
		// 不返回错误，因为用户已创建成功
	}

	// 生成邮箱验证令牌
	token := uc.generateToken()
	verification := &data.EmailVerification{
		UserID:    user.ID,
		Email:     email,
		Token:     token,
		Type:      "register",
		ExpiresAt: time.Now().Add(24 * time.Hour), // 24小时有效期
		Used:      false,
	}

	if err := uc.emailVerificationRepo.CreateVerification(ctx, verification); err != nil {
		log.Warn(ctx, "创建邮箱验证记录失败",
			log.ErrorField(err),
			log.Int64("user_id", user.ID),
		)
		// 不返回错误，因为用户已创建成功
	}

	// 发送验证邮件（异步）
	go func() {
		if err := uc.emailService.SendVerificationEmail(context.Background(), email, token); err != nil {
			log.Error(context.Background(), "发送验证邮件失败",
				log.ErrorField(err),
				log.String("email", email),
			)
		}
	}()

	log.Info(ctx, "用户注册成功",
		log.Int64("user_id", user.ID),
		log.String("username", username),
		log.String("email", email),
	)

	return user, nil
}

// Login 用户登录
func (uc *UserUseCase) Login(ctx context.Context, usernameOrEmail, password, ip string) (*data.User, error) {
	// 查询用户（支持用户名或邮箱登录）
	var user *data.User
	var err error

	// 先尝试按用户名查询
	user, err = uc.userRepo.GetUserByUsername(ctx, usernameOrEmail)
	if err != nil {
		return nil, err
	}

	// 如果用户名查询失败，尝试按邮箱查询
	if user == nil {
		user, err = uc.userRepo.GetUserByEmail(ctx, usernameOrEmail)
		if err != nil {
			return nil, err
		}
	}

	// 用户不存在
	if user == nil {
		return nil, ErrUserNotFound
	}

	// 验证密码
	if !uc.comparePassword(user.PasswordHash, password) {
		return nil, ErrInvalidPassword
	}

	// 检查用户状态
	if user.Status == "banned" {
		return nil, ErrUserBanned
	}

	if user.Status == "inactive" {
		// 可选：检查邮箱是否已验证
		// if !user.EmailVerified {
		// 	return nil, ErrEmailNotVerified
		// }
		// 如果允许未验证邮箱登录，则激活用户
		user.Status = "active"
		if err := uc.userRepo.UpdateUser(ctx, user); err != nil {
			log.Warn(ctx, "更新用户状态失败",
				log.ErrorField(err),
				log.Int64("user_id", user.ID),
			)
		}
	}

	// 更新最后登录信息
	if err := uc.userRepo.UpdateLastLogin(ctx, user.ID, ip); err != nil {
		log.Warn(ctx, "更新最后登录信息失败",
			log.ErrorField(err),
			log.Int64("user_id", user.ID),
		)
	}

	log.Info(ctx, "用户登录成功",
		log.Int64("user_id", user.ID),
		log.String("username", user.Username),
	)

	return user, nil
}

// VerifyEmail 验证邮箱
func (uc *UserUseCase) VerifyEmail(ctx context.Context, token string) error {
	// 查询验证记录
	verification, err := uc.emailVerificationRepo.GetVerificationByToken(ctx, token)
	if err != nil {
		return err
	}

	if verification == nil {
		return ErrInvalidVerificationToken
	}

	// 检查是否已使用
	if verification.Used {
		return ErrInvalidVerificationToken
	}

	// 检查是否过期
	if verification.ExpiresAt.Before(time.Now()) {
		return ErrInvalidVerificationToken
	}

	// 查询用户
	user, err := uc.userRepo.GetUserByID(ctx, verification.UserID)
	if err != nil {
		return err
	}

	if user == nil {
		return ErrUserNotFound
	}

	// 更新用户邮箱验证状态
	now := time.Now()
	user.EmailVerified = true
	user.EmailVerifiedAt = &now
	user.Status = "active" // 激活用户

	if err := uc.userRepo.UpdateUser(ctx, user); err != nil {
		return err
	}

	// 标记验证记录为已使用
	if err := uc.emailVerificationRepo.MarkAsUsed(ctx, verification.ID); err != nil {
		log.Warn(ctx, "标记验证记录为已使用失败",
			log.ErrorField(err),
			log.Int64("verification_id", verification.ID),
		)
	}

	log.Info(ctx, "邮箱验证成功",
		log.Int64("user_id", user.ID),
		log.String("email", user.Email),
	)

	return nil
}

// GetUser 获取用户信息
func (uc *UserUseCase) GetUser(ctx context.Context, userID int64) (*data.User, error) {
	user, err := uc.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}

// UpdateUser 更新用户信息
func (uc *UserUseCase) UpdateUser(ctx context.Context, userID int64, profile *data.UserProfile) error {
	// 检查用户是否存在
	user, err := uc.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	if user == nil {
		return ErrUserNotFound
	}

	// 查询现有资料
	existingProfile, err := uc.userProfileRepo.GetProfileByUserID(ctx, userID)
	if err != nil {
		return err
	}

	// 如果资料不存在，创建新资料
	if existingProfile == nil {
		profile.UserID = userID
		if err := uc.userProfileRepo.CreateProfile(ctx, profile); err != nil {
			return err
		}
	} else {
		// 更新现有资料
		profile.ID = existingProfile.ID
		profile.UserID = userID
		if err := uc.userProfileRepo.UpdateProfile(ctx, profile); err != nil {
			return err
		}
	}

	log.Info(ctx, "更新用户信息成功",
		log.Int64("user_id", userID),
	)

	return nil
}

// ChangePassword 修改密码
func (uc *UserUseCase) ChangePassword(ctx context.Context, userID int64, oldPassword, newPassword string) error {
	// 查询用户
	user, err := uc.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	if user == nil {
		return ErrUserNotFound
	}

	// 验证旧密码
	if !uc.comparePassword(user.PasswordHash, oldPassword) {
		return ErrInvalidPassword
	}

	// 加密新密码
	newPasswordHash, err := uc.hashPassword(newPassword)
	if err != nil {
		return err
	}

	// 更新密码
	user.PasswordHash = newPasswordHash
	if err := uc.userRepo.UpdateUser(ctx, user); err != nil {
		return err
	}

	log.Info(ctx, "修改密码成功",
		log.Int64("user_id", userID),
	)

	return nil
}

// ResendVerificationEmail 重新发送验证邮件
func (uc *UserUseCase) ResendVerificationEmail(ctx context.Context, email string) error {
	// 查询用户
	user, err := uc.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return err
	}

	if user == nil {
		return ErrUserNotFound
	}

	// 如果邮箱已验证，不需要重新发送
	if user.EmailVerified {
		return nil
	}

	// 生成新的验证令牌
	token := uc.generateToken()
	verification := &data.EmailVerification{
		UserID:    user.ID,
		Email:     email,
		Token:     token,
		Type:      "register",
		ExpiresAt: time.Now().Add(24 * time.Hour), // 24小时有效期
		Used:      false,
	}

	if err := uc.emailVerificationRepo.CreateVerification(ctx, verification); err != nil {
		return err
	}

	// 发送验证邮件（异步）
	go func() {
		if err := uc.emailService.SendVerificationEmail(context.Background(), email, token); err != nil {
			log.Error(context.Background(), "重新发送验证邮件失败",
				log.ErrorField(err),
				log.String("email", email),
			)
		}
	}()

	log.Info(ctx, "重新发送验证邮件成功",
		log.String("email", email),
	)

	return nil
}

// RequestPasswordReset 请求重置密码
func (uc *UserUseCase) RequestPasswordReset(ctx context.Context, email string) error {
	// 查询用户
	user, err := uc.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return err
	}

	if user == nil {
		// 为了安全，即使用户不存在也返回成功（防止邮箱枚举）
		return nil
	}

	// 生成重置密码令牌
	token := uc.generateToken()
	verification := &data.EmailVerification{
		UserID:    user.ID,
		Email:     email,
		Token:     token,
		Type:      "reset_password",
		ExpiresAt: time.Now().Add(1 * time.Hour), // 1小时有效期
		Used:      false,
	}

	if err := uc.emailVerificationRepo.CreateVerification(ctx, verification); err != nil {
		return err
	}

	// 发送重置密码邮件（异步）
	go func() {
		if err := uc.emailService.SendPasswordResetEmail(context.Background(), email, token); err != nil {
			log.Error(context.Background(), "发送重置密码邮件失败",
				log.ErrorField(err),
				log.String("email", email),
			)
		}
	}()

	log.Info(ctx, "发送重置密码邮件成功",
		log.String("email", email),
	)

	return nil
}

// ResetPassword 重置密码
func (uc *UserUseCase) ResetPassword(ctx context.Context, token, newPassword string) error {
	// 查询验证记录
	verification, err := uc.emailVerificationRepo.GetVerificationByToken(ctx, token)
	if err != nil {
		return err
	}

	if verification == nil {
		return ErrInvalidVerificationToken
	}

	// 检查类型是否为重置密码
	if verification.Type != "reset_password" {
		return ErrInvalidVerificationToken
	}

	// 检查是否已使用
	if verification.Used {
		return ErrInvalidVerificationToken
	}

	// 检查是否过期
	if verification.ExpiresAt.Before(time.Now()) {
		return ErrInvalidVerificationToken
	}

	// 查询用户
	user, err := uc.userRepo.GetUserByID(ctx, verification.UserID)
	if err != nil {
		return err
	}

	if user == nil {
		return ErrUserNotFound
	}

	// 加密新密码
	newPasswordHash, err := uc.hashPassword(newPassword)
	if err != nil {
		return err
	}

	// 更新密码
	user.PasswordHash = newPasswordHash
	if err := uc.userRepo.UpdateUser(ctx, user); err != nil {
		return err
	}

	// 标记验证记录为已使用
	if err := uc.emailVerificationRepo.MarkAsUsed(ctx, verification.ID); err != nil {
		log.Warn(ctx, "标记验证记录为已使用失败",
			log.ErrorField(err),
			log.Int64("verification_id", verification.ID),
		)
	}

	log.Info(ctx, "密码重置成功",
		log.Int64("user_id", user.ID),
		log.String("email", user.Email),
	)

	return nil
}

// generateUUID 生成UUID
func generateUUID() string {
	return uuid.New().String()
}
