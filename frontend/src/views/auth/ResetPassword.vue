<template>
  <div class="reset-password-page">
    <CodeRain :useConfigFile="true" />
    <div class="reset-password-container">
      <div class="reset-password-header">
        <h1 class="reset-password-title">重置密码</h1>
        <p class="reset-password-subtitle">请输入您的新密码</p>
      </div>

      <Form
        ref="formRef"
        :model="form"
        :rules="rules"
        class="reset-password-form"
        @submit.prevent="handleSubmit"
      >
        <FormItem label="新密码" prop="newPassword">
          <Input
            v-model="form.newPassword"
            type="password"
            placeholder="请输入新密码（6-20个字符）"
            size="large"
            :prefix-icon="Lock"
            show-password
            clearable
          />
        </FormItem>

        <FormItem label="确认密码" prop="confirmPassword">
          <Input
            v-model="form.confirmPassword"
            type="password"
            placeholder="请再次输入新密码"
            size="large"
            :prefix-icon="Lock"
            show-password
            clearable
            @keyup.enter="handleSubmit"
          />
        </FormItem>

        <FormItem>
          <Button
            type="primary"
            size="large"
            :loading="loading"
            native-type="submit"
            class="submit-button"
          >
            重置密码
          </Button>
        </FormItem>

        <FormItem>
          <div class="reset-password-footer">
            <Link type="primary" @click="goToLogin">返回登录</Link>
          </div>
        </FormItem>
      </Form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { Lock } from '@element-plus/icons-vue'
import { Form, FormItem, Input, Button, Link } from '@/components/common/base'
import { CodeRain } from '@/components/common/effects'
import { userService } from '@/api/services/user.service'
import { success, error } from '@/components/common/base/Message'
import type { FormInstance, FormRules } from 'element-plus'

const router = useRouter()
const route = useRoute()
const formRef = ref<FormInstance>()
const loading = ref(false)

const form = reactive({
  newPassword: '',
  confirmPassword: '',
})

// 验证确认密码
const validateConfirmPassword = (rule: any, value: string, callback: any) => {
  if (value !== form.newPassword) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

const rules: FormRules = {
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, max: 20, message: '密码长度为 6-20 个字符', trigger: 'blur' },
  ],
  confirmPassword: [
    { required: true, message: '请再次输入新密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' },
  ],
}

// 从 URL 获取 token
const token = ref<string>('')

onMounted(() => {
  const tokenParam = route.query.token as string
  if (!tokenParam) {
    error('重置令牌无效，请重新申请')
    setTimeout(() => {
      router.push('/auth/forgot-password')
    }, 2000)
    return
  }
  token.value = tokenParam
})

const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    if (!token.value) {
      error('重置令牌无效，请重新申请')
      return
    }

    loading.value = true
    try {
      const response = await userService.resetPassword({
        token: token.value,
        newPassword: form.newPassword,
      })

      if (response.data?.success) {
        success('密码重置成功，请使用新密码登录')
        setTimeout(() => {
          router.push('/auth/login')
        }, 2000)
      } else {
        error(response.data?.message || '重置失败，请稍后重试')
      }
    } catch (err: any) {
      const errorMessage = err?.response?.data?.message || err?.message || '重置失败，请稍后重试'
      error(errorMessage)
      console.error('Reset password error:', err)
    } finally {
      loading.value = false
    }
  })
}

const goToLogin = () => {
  router.push('/auth/login')
}
</script>

<style scoped lang="scss">
.reset-password-page {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  width: 100vw;
  background: #000000;
  padding: 20px;
  overflow: hidden;
  box-sizing: border-box;
}

.reset-password-container {
  position: relative;
  z-index: 10;
  width: 100%;
  max-width: 400px;
  background-color: rgba(0, 0, 0, 0.85);
  border: 1px solid rgba(0, 255, 0, 0.3);
  border-radius: 8px;
  box-shadow:
    0 0 40px rgba(0, 255, 0, 0.5),
    inset 0 0 10px rgba(0, 255, 0, 0.2);
  padding: 40px;
  backdrop-filter: blur(10px);
}

.reset-password-header {
  text-align: center;
  margin-bottom: 32px;
}

.reset-password-title {
  font-size: 28px;
  font-weight: 600;
  color: #00FF00;
  margin: 0 0 8px 0;
  text-shadow: 0 0 10px rgba(0, 255, 0, 0.5);
}

.reset-password-subtitle {
  font-size: 14px;
  color: rgba(0, 255, 0, 0.7);
  margin: 0;
}

.reset-password-form {
  .submit-button {
    width: 100%;
  }

  .reset-password-footer {
    text-align: center;
    font-size: 14px;

    // Link 组件样式已通过全局样式统一，无需额外样式
  }

  :deep(.el-form-item__label) {
    color: rgba(255, 255, 255, 0.9);
  }

  :deep(.el-input__wrapper) {
    background-color: rgba(0, 0, 0, 0.5);
    border: 1px solid rgba(0, 255, 0, 0.3);
    box-shadow: 0 0 5px rgba(0, 255, 0, 0.1);

    &:hover {
      border-color: rgba(0, 255, 0, 0.5);
      box-shadow: 0 0 8px rgba(0, 255, 0, 0.2);
    }

    &.is-focus {
      border-color: rgba(0, 255, 0, 0.6);
      box-shadow: 0 0 10px rgba(0, 255, 0, 0.3);
    }
  }

  :deep(.el-input__inner) {
    color: #00FF00;

    &::placeholder {
      color: rgba(0, 255, 0, 0.4);
    }
  }

  :deep(.el-button--primary) {
    background-color: #00FF00;
    border-color: #00FF00;
    color: #000000;
    font-weight: 600;

    &:hover {
      background-color: #00cc00;
      border-color: #00cc00;
      box-shadow: 0 0 15px rgba(0, 255, 0, 0.5);
    }
  }
}
</style>

