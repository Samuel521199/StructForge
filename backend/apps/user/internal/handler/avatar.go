package handler

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"StructForge/backend/common/log"

	"github.com/disintegration/imaging"
	"github.com/go-kratos/kratos/v2/transport/http"
)

const (
	// 头像存储目录
	avatarDir = "uploads/avatars"
	// 头像最大尺寸（像素）
	maxAvatarSize = 512
	// 头像质量（JPEG）
	jpegQuality = 85
	// 最大文件大小（2MB）
	maxFileSize = 2 * 1024 * 1024
)

// UploadAvatar 上传头像
func UploadAvatar(ctx http.Context) error {
	requestCtx := ctx.Request().Context()

	// 解析 multipart form
	err := ctx.Request().ParseMultipartForm(maxFileSize)
	if err != nil {
		log.Error(requestCtx, "解析表单失败", log.ErrorField(err))
		return ctx.JSON(400, map[string]interface{}{
			"code":    400,
			"message": "解析表单失败",
			"error":   err.Error(),
		})
	}

	// 获取文件
	file, header, err := ctx.Request().FormFile("file")
	if err != nil {
		log.Error(requestCtx, "获取文件失败", log.ErrorField(err))
		return ctx.JSON(400, map[string]interface{}{
			"code":    400,
			"message": "请选择文件",
			"error":   err.Error(),
		})
	}
	defer file.Close()

	// 验证文件大小
	if header.Size > maxFileSize {
		return ctx.JSON(400, map[string]interface{}{
			"code":    400,
			"message": fmt.Sprintf("文件大小不能超过 %dMB", maxFileSize/(1024*1024)),
		})
	}

	// 验证文件类型
	contentType := header.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return ctx.JSON(400, map[string]interface{}{
			"code":    400,
			"message": "只支持图片文件",
		})
	}

	// 读取文件数据
	fileData, err := io.ReadAll(file)
	if err != nil {
		log.Error(requestCtx, "读取文件失败", log.ErrorField(err))
		return ctx.JSON(500, map[string]interface{}{
			"code":    500,
			"message": "读取文件失败",
			"error":   err.Error(),
		})
	}

	// 解码图片
	img, format, err := image.Decode(strings.NewReader(string(fileData)))
	if err != nil {
		log.Error(requestCtx, "解码图片失败", log.ErrorField(err))
		return ctx.JSON(400, map[string]interface{}{
			"code":    400,
			"message": "无效的图片文件",
			"error":   err.Error(),
		})
	}

	// 只支持 JPEG 和 PNG
	if format != "jpeg" && format != "png" {
		return ctx.JSON(400, map[string]interface{}{
			"code":    400,
			"message": "只支持 JPEG 和 PNG 格式",
		})
	}

	// 调整图片大小（保持宽高比）
	img = imaging.Fit(img, maxAvatarSize, maxAvatarSize, imaging.Lanczos)

	// 创建存储目录
	if err := os.MkdirAll(avatarDir, 0755); err != nil {
		log.Error(requestCtx, "创建目录失败", log.ErrorField(err))
		return ctx.JSON(500, map[string]interface{}{
			"code":    500,
			"message": "创建存储目录失败",
			"error":   err.Error(),
		})
	}

	// 生成文件名（使用时间戳和随机字符串）
	fileName := fmt.Sprintf("%d_%s.%s", time.Now().Unix(), generateRandomString(8), format)
	filePath := filepath.Join(avatarDir, fileName)

	// 保存文件
	outputFile, err := os.Create(filePath)
	if err != nil {
		log.Error(requestCtx, "创建文件失败", log.ErrorField(err))
		return ctx.JSON(500, map[string]interface{}{
			"code":    500,
			"message": "保存文件失败",
			"error":   err.Error(),
		})
	}
	defer outputFile.Close()

	// 根据格式编码图片
	if format == "jpeg" {
		err = jpeg.Encode(outputFile, img, &jpeg.Options{Quality: jpegQuality})
	} else {
		err = png.Encode(outputFile, img)
	}
	if err != nil {
		log.Error(requestCtx, "编码图片失败", log.ErrorField(err))
		return ctx.JSON(500, map[string]interface{}{
			"code":    500,
			"message": "保存图片失败",
			"error":   err.Error(),
		})
	}

	// 生成访问 URL（这里使用相对路径，实际应该配置静态文件服务）
	avatarURL := fmt.Sprintf("/uploads/avatars/%s", fileName)

	log.Info(requestCtx, "头像上传成功",
		log.String("file_name", fileName),
		log.String("avatar_url", avatarURL),
	)

	// 返回成功响应
	return ctx.JSON(200, map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{
			"url": avatarURL,
		},
		"message": "上传成功",
	})
}

// generateRandomString 生成随机字符串
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	seed := time.Now().UnixNano()
	for i := range b {
		seed = seed*1103515245 + 12345
		b[i] = charset[seed%int64(len(charset))]
	}
	return string(b)
}
