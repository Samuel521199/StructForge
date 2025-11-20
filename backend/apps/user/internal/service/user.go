package service

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	v1 "StructForge/backend/api/user/v1"
	"StructForge/backend/apps/user/internal/biz"
	"StructForge/backend/apps/user/internal/data"
	"StructForge/backend/common/log"
)

// UserService 用户服务实现
type UserService struct {
	v1.UnimplementedUserServiceServer

	uc     *biz.UserUseCase
	jwtMgr *biz.JWTManager
}

// NewUserService 创建用户服务实例
func NewUserService(uc *biz.UserUseCase, jwtMgr *biz.JWTManager) *UserService {
	return &UserService{
		uc:     uc,
		jwtMgr: jwtMgr,
	}
}

// Register 用户注册
func (s *UserService) Register(ctx context.Context, req *v1.RegisterRequest) (*v1.RegisterResponse, error) {
	// 参数验证
	if req.Username == "" {
		return &v1.RegisterResponse{
			Success: false,
			Message: "用户名不能为空",
		}, nil
	}
	if req.Email == "" {
		return &v1.RegisterResponse{
			Success: false,
			Message: "邮箱不能为空",
		}, nil
	}
	if req.Password == "" {
		return &v1.RegisterResponse{
			Success: false,
			Message: "密码不能为空",
		}, nil
	}

	// 调用业务逻辑
	user, err := s.uc.Register(ctx, req.Username, req.Email, req.Password)
	if err != nil {
		if err == biz.ErrUserExists {
			return &v1.RegisterResponse{
				Success: false,
				Message: "用户名或邮箱已存在",
			}, nil
		}
		// 密码格式错误
		if err == biz.ErrInvalidPasswordFormat ||
			(err.Error() != "" && (err.Error() == "密码格式不符合要求：必须包含字母和数字，且包含大小写字母或特殊字符" ||
				err.Error() == "密码长度必须在 6-20 个字符之间" ||
				err.Error() == "密码必须包含字母和数字" ||
				err.Error() == "密码必须包含大小写字母或特殊字符" ||
				err.Error() == "密码必须包含大写字母或特殊字符" ||
				err.Error() == "密码必须包含小写字母或特殊字符")) {
			return &v1.RegisterResponse{
				Success: false,
				Message: err.Error(),
			}, nil
		}
		log.Error(ctx, "用户注册失败",
			log.ErrorField(err),
			log.String("username", req.Username),
			log.String("email", req.Email),
		)
		return &v1.RegisterResponse{
			Success: false,
			Message: "注册失败，请稍后重试",
		}, nil
	}

	// 转换为 Protobuf 格式
	pbUser := s.toProtoUser(user)

	return &v1.RegisterResponse{
		Success: true,
		Message: "注册成功，请查收邮件验证",
		User:    pbUser,
	}, nil
}

// Login 用户登录
func (s *UserService) Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginResponse, error) {
	// 参数验证
	if req.Username == "" {
		return &v1.LoginResponse{
			Success: false,
			Message: "用户名或邮箱不能为空",
		}, nil
	}
	if req.Password == "" {
		return &v1.LoginResponse{
			Success: false,
			Message: "密码不能为空",
		}, nil
	}

	// 获取客户端IP（从context中获取，需要中间件注入）
	ip := getClientIP(ctx)

	// 调用业务逻辑
	user, err := s.uc.Login(ctx, req.Username, req.Password, ip)
	if err != nil {
		if err == biz.ErrUserNotFound {
			return &v1.LoginResponse{
				Success: false,
				Message: "用户名或密码错误",
			}, nil
		}
		if err == biz.ErrInvalidPassword {
			return &v1.LoginResponse{
				Success: false,
				Message: "用户名或密码错误",
			}, nil
		}
		if err == biz.ErrUserBanned {
			return &v1.LoginResponse{
				Success: false,
				Message: "账户已被封禁",
			}, nil
		}
		log.Error(ctx, "用户登录失败",
			log.ErrorField(err),
			log.String("username", req.Username),
		)
		return &v1.LoginResponse{
			Success: false,
			Message: "登录失败，请稍后重试",
		}, nil
	}

	// 生成 JWT Token
	token, err := s.jwtMgr.GenerateToken(user.ID, user.Username)
	if err != nil {
		log.Error(ctx, "生成Token失败",
			log.ErrorField(err),
			log.Int64("user_id", user.ID),
		)
		return &v1.LoginResponse{
			Success: false,
			Message: "登录失败，请稍后重试",
		}, nil
	}

	// 转换为 Protobuf 格式
	pbUser := s.toProtoUser(user)

	return &v1.LoginResponse{
		Success: true,
		Message: "登录成功",
		Token:   token,
		User:    pbUser,
	}, nil
}

// VerifyEmail 验证邮箱
func (s *UserService) VerifyEmail(ctx context.Context, req *v1.VerifyEmailRequest) (*v1.VerifyEmailResponse, error) {
	if req.Token == "" {
		return &v1.VerifyEmailResponse{
			Success: false,
			Message: "验证令牌不能为空",
		}, nil
	}

	err := s.uc.VerifyEmail(ctx, req.Token)
	if err != nil {
		if err == biz.ErrInvalidVerificationToken {
			return &v1.VerifyEmailResponse{
				Success: false,
				Message: "验证令牌无效或已过期",
			}, nil
		}
		log.Error(ctx, "邮箱验证失败",
			log.ErrorField(err),
			log.String("token", req.Token),
		)
		return &v1.VerifyEmailResponse{
			Success: false,
			Message: "验证失败，请稍后重试",
		}, nil
	}

	return &v1.VerifyEmailResponse{
		Success: true,
		Message: "邮箱验证成功",
	}, nil
}

// ResendVerificationEmail 重新发送验证邮件
func (s *UserService) ResendVerificationEmail(ctx context.Context, req *v1.ResendVerificationEmailRequest) (*v1.ResendVerificationEmailResponse, error) {
	if req.Email == "" {
		return &v1.ResendVerificationEmailResponse{
			Success: false,
			Message: "邮箱不能为空",
		}, nil
	}

	err := s.uc.ResendVerificationEmail(ctx, req.Email)
	if err != nil {
		if err == biz.ErrUserNotFound {
			return &v1.ResendVerificationEmailResponse{
				Success: false,
				Message: "该邮箱未注册",
			}, nil
		}
		log.Error(ctx, "重新发送验证邮件失败",
			log.ErrorField(err),
			log.String("email", req.Email),
		)
		return &v1.ResendVerificationEmailResponse{
			Success: false,
			Message: "发送失败，请稍后重试",
		}, nil
	}

	return &v1.ResendVerificationEmailResponse{
		Success: true,
		Message: "验证邮件已发送，请查收",
	}, nil
}

// GetUser 获取用户信息
func (s *UserService) GetUser(ctx context.Context, req *v1.GetUserRequest) (*v1.GetUserResponse, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "用户ID不能为空")
	}

	user, err := s.uc.GetUser(ctx, req.Id)
	if err != nil {
		if err == biz.ErrUserNotFound {
			return nil, status.Error(codes.NotFound, "用户不存在")
		}
		log.Error(ctx, "获取用户信息失败",
			log.ErrorField(err),
			log.Int64("user_id", req.Id),
		)
		return nil, status.Error(codes.Internal, "获取用户信息失败")
	}

	pbUser := s.toProtoUser(user)
	return &v1.GetUserResponse{
		User: pbUser,
	}, nil
}

// GetCurrentUser 获取当前用户信息
func (s *UserService) GetCurrentUser(ctx context.Context, req *v1.GetCurrentUserRequest) (*v1.GetCurrentUserResponse, error) {
	// 从 context 中获取用户ID（需要中间件注入）
	userID, ok := getUserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "未认证")
	}

	user, err := s.uc.GetUser(ctx, userID)
	if err != nil {
		if err == biz.ErrUserNotFound {
			return nil, status.Error(codes.NotFound, "用户不存在")
		}
		log.Error(ctx, "获取当前用户信息失败",
			log.ErrorField(err),
			log.Int64("user_id", userID),
		)
		return nil, status.Error(codes.Internal, "获取用户信息失败")
	}

	pbUser := s.toProtoUser(user)
	return &v1.GetCurrentUserResponse{
		User: pbUser,
	}, nil
}

// UpdateUser 更新用户信息
func (s *UserService) UpdateUser(ctx context.Context, req *v1.UpdateUserRequest) (*v1.UpdateUserResponse, error) {
	// 从 context 中获取用户ID
	userID, ok := getUserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "未认证")
	}

	// 转换为业务模型（使用 data 层的类型）
	profile := &data.UserProfile{
		Nickname:  req.Nickname,
		AvatarURL: req.AvatarUrl,
		Bio:       req.Bio,
		Phone:     req.Phone,
		Gender:    req.Gender,
		Location:  req.Location,
		Website:   req.Website,
	}

	// 解析生日
	if req.Birthday != "" {
		birthday, err := time.Parse("2006-01-02", req.Birthday)
		if err == nil {
			profile.Birthday = &birthday
		}
	}

	err := s.uc.UpdateUser(ctx, userID, profile)
	if err != nil {
		if err == biz.ErrUserNotFound {
			return &v1.UpdateUserResponse{
				Success: false,
				Message: "用户不存在",
			}, nil
		}
		log.Error(ctx, "更新用户信息失败",
			log.ErrorField(err),
			log.Int64("user_id", userID),
		)
		return &v1.UpdateUserResponse{
			Success: false,
			Message: "更新失败，请稍后重试",
		}, nil
	}

	// 重新获取用户信息
	user, err := s.uc.GetUser(ctx, userID)
	if err != nil {
		log.Error(ctx, "获取更新后的用户信息失败",
			log.ErrorField(err),
			log.Int64("user_id", userID),
		)
		return &v1.UpdateUserResponse{
			Success: false,
			Message: "更新成功，但获取用户信息失败",
		}, nil
	}

	pbUser := s.toProtoUser(user)
	return &v1.UpdateUserResponse{
		Success: true,
		Message: "更新成功",
		User:    pbUser,
	}, nil
}

// ChangePassword 修改密码
func (s *UserService) ChangePassword(ctx context.Context, req *v1.ChangePasswordRequest) (*v1.ChangePasswordResponse, error) {
	// 从 context 中获取用户ID
	userID, ok := getUserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "未认证")
	}

	if req.OldPassword == "" || req.NewPassword == "" {
		return &v1.ChangePasswordResponse{
			Success: false,
			Message: "密码不能为空",
		}, nil
	}

	err := s.uc.ChangePassword(ctx, userID, req.OldPassword, req.NewPassword)
	if err != nil {
		if err == biz.ErrInvalidPassword {
			return &v1.ChangePasswordResponse{
				Success: false,
				Message: "旧密码错误",
			}, nil
		}
		log.Error(ctx, "修改密码失败",
			log.ErrorField(err),
			log.Int64("user_id", userID),
		)
		return &v1.ChangePasswordResponse{
			Success: false,
			Message: "修改失败，请稍后重试",
		}, nil
	}

	return &v1.ChangePasswordResponse{
		Success: true,
		Message: "密码修改成功",
	}, nil
}

// RequestPasswordReset 请求重置密码
func (s *UserService) RequestPasswordReset(ctx context.Context, req *v1.RequestPasswordResetRequest) (*v1.RequestPasswordResetResponse, error) {
	if req.Email == "" {
		return &v1.RequestPasswordResetResponse{
			Success: false,
			Message: "邮箱不能为空",
		}, nil
	}

	err := s.uc.RequestPasswordReset(ctx, req.Email)
	if err != nil {
		log.Error(ctx, "请求重置密码失败",
			log.ErrorField(err),
			log.String("email", req.Email),
		)
		return &v1.RequestPasswordResetResponse{
			Success: false,
			Message: "发送失败，请稍后重试",
		}, nil
	}

	// 为了安全，无论用户是否存在都返回成功
	return &v1.RequestPasswordResetResponse{
		Success: true,
		Message: "如果该邮箱已注册，重置密码邮件已发送，请查收",
	}, nil
}

// ResetPassword 重置密码
func (s *UserService) ResetPassword(ctx context.Context, req *v1.ResetPasswordRequest) (*v1.ResetPasswordResponse, error) {
	if req.Token == "" || req.NewPassword == "" {
		return &v1.ResetPasswordResponse{
			Success: false,
			Message: "令牌和密码不能为空",
		}, nil
	}

	// 验证密码长度
	if len(req.NewPassword) < 6 || len(req.NewPassword) > 20 {
		return &v1.ResetPasswordResponse{
			Success: false,
			Message: "密码长度必须在 6-20 个字符之间",
		}, nil
	}

	err := s.uc.ResetPassword(ctx, req.Token, req.NewPassword)
	if err != nil {
		if err == biz.ErrInvalidVerificationToken {
			return &v1.ResetPasswordResponse{
				Success: false,
				Message: "重置令牌无效或已过期",
			}, nil
		}
		log.Error(ctx, "重置密码失败",
			log.ErrorField(err),
			log.String("token", req.Token),
		)
		return &v1.ResetPasswordResponse{
			Success: false,
			Message: "重置失败，请稍后重试",
		}, nil
	}

	return &v1.ResetPasswordResponse{
		Success: true,
		Message: "密码重置成功",
	}, nil
}

// toProtoUser 转换为 Protobuf User
func (s *UserService) toProtoUser(user *data.User) *v1.User {
	if user == nil {
		return nil
	}

	pbUser := &v1.User{
		Id:            user.ID,
		Username:      user.Username,
		Email:         user.Email,
		EmailVerified: user.EmailVerified,
		Status:        user.Status,
		CreatedAt:     timestamppb.New(user.CreatedAt),
		UpdatedAt:     timestamppb.New(user.UpdatedAt),
	}

	// 注意：proto 中只有 EmailVerified (bool)，没有 EmailVerifiedAt (timestamp)
	// 如果需要时间戳，需要在 proto 中添加该字段

	if user.LastLoginAt != nil {
		pbUser.LastLoginAt = timestamppb.New(*user.LastLoginAt)
	}

	if user.LastLoginIP != "" {
		pbUser.LastLoginIp = user.LastLoginIP
	}

	// 转换用户资料
	if user.Profile != nil {
		pbUser.Profile = s.toProtoUserProfile(user.Profile)
	}

	return pbUser
}

// toProtoUserProfile 转换为 Protobuf UserProfile
func (s *UserService) toProtoUserProfile(profile *data.UserProfile) *v1.UserProfile {
	if profile == nil {
		return nil
	}

	pbProfile := &v1.UserProfile{
		Id:        profile.ID,
		UserId:    profile.UserID,
		Nickname:  profile.Nickname,
		AvatarUrl: profile.AvatarURL,
		Bio:       profile.Bio,
		Phone:     profile.Phone,
		Gender:    profile.Gender,
		Location:  profile.Location,
		Website:   profile.Website,
		CreatedAt: timestamppb.New(profile.CreatedAt),
		UpdatedAt: timestamppb.New(profile.UpdatedAt),
	}

	if profile.Birthday != nil {
		pbProfile.Birthday = profile.Birthday.Format("2006-01-02")
	}

	return pbProfile
}

// getUserIDFromContext 从 context 中获取用户ID
func getUserIDFromContext(ctx context.Context) (int64, bool) {
	// TODO: 从 context 中获取用户ID（需要中间件注入）
	// 这里先返回一个示例值，实际应该从 JWT Token 中解析
	return 0, false
}

// getClientIP 从 context 中获取客户端IP
func getClientIP(ctx context.Context) string {
	// TODO: 从 context 中获取客户端IP（需要中间件注入）
	return "127.0.0.1"
}
