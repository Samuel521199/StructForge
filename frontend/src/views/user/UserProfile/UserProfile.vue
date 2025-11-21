<template>
  <div class="user-profile-page">
    <Card>
      <template #header>
        <div class="profile-header">
          <h2>个人中心</h2>
        </div>
      </template>

      <div class="profile-content">
        <!-- 用户信息卡片 -->
        <div class="profile-card">
          <div class="profile-avatar">
            <Avatar :size="100" :src="form.avatar" :text="userDisplayName" />
            <Upload
              class="avatar-upload"
              :action="uploadAction"
              :show-file-list="false"
              :headers="uploadHeaders"
              :before-upload="beforeAvatarUpload"
              @success="handleAvatarSuccess"
            >
              <Button
                type="primary"
                :icon="Camera"
                circle
                class="avatar-upload-btn"
              />
            </Upload>
          </div>

          <div class="profile-info">
            <h3>{{ userDisplayName }}</h3>
            <p class="profile-email">{{ user?.email }}</p>
            <div class="profile-meta">
              <span>注册时间：{{ formatDate(user?.createdAt) }}</span>
            </div>
          </div>
        </div>

        <!-- 个人信息表单 -->
        <Form
          ref="formRef"
          :model="form"
          :rules="rules"
          label-width="100px"
          class="profile-form"
        >
          <FormItem label="用户名" prop="username">
            <Input v-model="form.username" placeholder="请输入用户名" />
          </FormItem>

          <FormItem label="邮箱" prop="email">
            <Input v-model="form.email" type="email" placeholder="请输入邮箱" />
          </FormItem>

          <FormItem label="头像URL" prop="avatar">
            <Input v-model="form.avatar" placeholder="请输入头像URL" />
          </FormItem>

          <FormItem>
            <Button type="primary" :loading="loading" @click="handleSave">
              保存
            </Button>
            <Button style="margin-left: 12px" @click="handleReset">重置</Button>
          </FormItem>
        </Form>
      </div>
    </Card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { Camera } from '@element-plus/icons-vue'
import { Card, Form, FormItem, Input, Button, Avatar, Upload } from '@/components/common/base'
import { useUserStore } from '@/stores/modules/user.store'
import { useAuthStore } from '@/stores/modules/auth.store'
import { userService } from '@/api/services/user.service'
import { success, error } from '@/components/common/base/Message'
import type { FormInstance, FormRules } from 'element-plus'
import type { UploadProps } from '@/components/common/base/Upload/types'

const userStore = useUserStore()
const authStore = useAuthStore()

const formRef = ref<FormInstance>()
const loading = ref(false)

const user = computed(() => userStore.user)
const userDisplayName = computed(() => {
  if (!user.value) return '用户'
  return user.value.username || user.value.email || '用户'
})

const form = reactive({
  username: '',
  email: '',
  avatar: '',
  nickname: '',
  bio: '',
})

const rules: FormRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度为 3-20 个字符', trigger: 'blur' },
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' },
  ],
}

// 上传配置
const uploadAction = computed(() => {
  const baseURL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8000'
  return `${baseURL}/api/v1/users/avatar`
})

const uploadHeaders = computed(() => {
  return {
    Authorization: `Bearer ${authStore.token}`,
  }
})

// 初始化表单
onMounted(async () => {
  // 加载用户信息
  try {
    const response = await userService.getUserInfo()
    if (response.data) {
      userStore.updateUser(response.data)
    }
  } catch (err) {
    console.error('加载用户信息失败:', err)
  }

  if (user.value) {
    form.username = user.value.username || ''
    form.email = user.value.email || ''
    form.avatar = user.value.profile?.avatarUrl || user.value.avatar || ''
    form.nickname = user.value.profile?.nickname || ''
    form.bio = user.value.profile?.bio || ''
  }
})

const formatDate = (date?: string) => {
  if (!date) return '-'
  return new Date(date).toLocaleString('zh-CN')
}

// 头像上传前验证和处理
const beforeAvatarUpload: UploadProps['beforeUpload'] = async (file) => {
  const isJPG = file.type === 'image/jpeg' || file.type === 'image/png' || file.type === 'image/webp'
  const isLt5M = file.size / 1024 / 1024 < 5

  if (!isJPG) {
    error('头像图片只能是 JPG/PNG/WebP 格式!')
    return false
  }
  if (!isLt5M) {
    error('头像图片大小不能超过 5MB!')
    return false
  }

  // 压缩和裁剪图片
  try {
    const { processAvatar } = await import('@/utils/imageCompress')
    const processedFile = await processAvatar(file, {
      maxWidth: 512,
      maxHeight: 512,
      quality: 0.85,
      outputFormat: 'image/jpeg',
      cropToCircle: false, // 不强制裁剪为圆形，保持原图比例
    })
    
    // 替换原文件
    Object.defineProperty(processedFile, 'name', {
      writable: true,
      value: file.name,
    })
    
    return processedFile as any
  } catch (err) {
    console.error('图片处理失败:', err)
    error('图片处理失败，请重试')
    return false
  }
}

// 头像上传成功
interface AvatarUploadResponse {
  code?: number
  data?: {
    url?: string
  }
  message?: string
}

const handleAvatarSuccess = (response: AvatarUploadResponse) => {
  if (response?.code === 200 && response?.data?.url) {
    form.avatar = response.data.url
    userStore.updateUser({ avatar: response.data.url })
    success('头像上传成功')
  } else {
    error(response?.message || '头像上传失败')
  }
}

const handleSave = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    loading.value = true
    try {
      const response = await userService.updateUserInfo({
        avatar_url: form.avatar,
        nickname: form.nickname,
        bio: form.bio,
      } as any)

      if (response.data) {
        userStore.updateUser(response.data)
        success('保存成功')
      } else {
        error('保存失败')
      }
    } catch (err: any) {
      error(err?.message || '保存失败')
      console.error('Update user error:', err)
    } finally {
      loading.value = false
    }
  })
}

const handleReset = () => {
  if (user.value) {
    form.username = user.value.username || ''
    form.email = user.value.email || ''
    form.avatar = user.value.profile?.avatarUrl || user.value.avatar || ''
    form.nickname = user.value.profile?.nickname || ''
    form.bio = user.value.profile?.bio || ''
  }
  formRef.value?.clearValidate()
}
</script>

<style scoped lang="scss">
.user-profile-page {
  padding: 24px;
  max-width: 800px;
  margin: 0 auto;
}

.profile-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.profile-content {
  .profile-card {
    display: flex;
    align-items: center;
    gap: 24px;
    padding: 24px;
    margin-bottom: 24px;
    background: var(--el-bg-color-page);
    border-radius: 8px;
  }

  .profile-avatar {
    position: relative;
    flex-shrink: 0;

    .avatar-upload {
      position: absolute;
      bottom: 0;
      right: 0;
    }

    .avatar-upload-btn {
      width: 32px;
      height: 32px;
    }
  }

  .profile-info {
    flex: 1;

    h3 {
      margin: 0 0 8px 0;
      font-size: 20px;
      font-weight: 600;
      color: var(--el-text-color-primary);
    }

    .profile-email {
      margin: 0 0 8px 0;
      color: var(--el-text-color-regular);
    }

    .profile-meta {
      font-size: 14px;
      color: var(--el-text-color-secondary);
    }
  }

  .profile-form {
    margin-top: 24px;
  }
}
</style>
