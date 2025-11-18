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
            <el-avatar :size="100" :src="form.avatar">
              {{ userDisplayName }}
            </el-avatar>
            <el-upload
              class="avatar-upload"
              :action="uploadAction"
              :show-file-list="false"
              :on-success="handleAvatarSuccess"
              :before-upload="beforeAvatarUpload"
              :headers="uploadHeaders"
            >
              <el-button
                type="primary"
                :icon="Camera"
                circle
                class="avatar-upload-btn"
              />
            </el-upload>
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
import { Card, Form, FormItem, Input, Button } from '@/components/common/base'
import { useUserStore } from '@/stores/modules/user.store'
import { useAuthStore } from '@/stores/modules/auth.store'
import { userService } from '@/api/services/user.service'
import { success, error } from '@/components/common/base/Message'
import type { FormInstance, FormRules, UploadProps } from 'element-plus'
import { ElMessage } from 'element-plus'

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
  // TODO: 替换为实际的上传API地址
  return '/api/v1/users/avatar'
})

const uploadHeaders = computed(() => {
  return {
    Authorization: `Bearer ${authStore.token}`,
  }
})

// 初始化表单
onMounted(() => {
  if (user.value) {
    form.username = user.value.username
    form.email = user.value.email
    form.avatar = user.value.avatar || ''
  }
})

const formatDate = (date?: string) => {
  if (!date) return '-'
  return new Date(date).toLocaleString('zh-CN')
}

// 头像上传前验证
const beforeAvatarUpload: UploadProps['beforeUpload'] = (file) => {
  const isJPG = file.type === 'image/jpeg' || file.type === 'image/png'
  const isLt2M = file.size / 1024 / 1024 < 2

  if (!isJPG) {
    ElMessage.error('头像图片只能是 JPG/PNG 格式!')
    return false
  }
  if (!isLt2M) {
    ElMessage.error('头像图片大小不能超过 2MB!')
    return false
  }
  return true
}

// 头像上传成功
const handleAvatarSuccess: UploadProps['onSuccess'] = (response: any) => {
  if (response && response.data && response.data.url) {
    form.avatar = response.data.url
    userStore.updateUser({ avatar: response.data.url })
    success('头像上传成功')
  } else {
    error('头像上传失败')
  }
}

const handleSave = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    loading.value = true
    try {
      const response = await userService.updateUserInfo({
        username: form.username,
        email: form.email,
        avatar: form.avatar,
      })

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
    form.username = user.value.username
    form.email = user.value.email
    form.avatar = user.value.avatar || ''
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
