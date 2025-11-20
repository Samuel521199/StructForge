<template>
  <div class="forgot-password-page">
    <CodeRain :useConfigFile="true" />
    <div class="forgot-password-container">
      <div class="forgot-password-header">
        <h1 class="forgot-password-title">忘记密码</h1>
        <p class="forgot-password-subtitle">请输入您的邮箱地址，我们将发送重置密码链接</p>
      </div>

      <Form
        ref="formRef"
        :model="form"
        :rules="rules"
        class="forgot-password-form"
        @submit.prevent="handleSubmit"
      >
        <FormItem label="邮箱" prop="email">
          <Input
            v-model="form.email"
            type="email"
            placeholder="请输入注册时使用的邮箱地址"
            size="large"
            :prefix-icon="Message"
            clearable
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
            发送重置链接
          </Button>
        </FormItem>

        <FormItem>
          <div class="forgot-password-footer">
            <Link type="primary" @click="goToLogin">返回登录</Link>
          </div>
        </FormItem>
      </Form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { Message } from '@element-plus/icons-vue'
import { Form, FormItem, Input, Button, Link } from '@/components/common/base'
import { CodeRain } from '@/components/common/effects'
import { userService } from '@/api/services/user.service'
import { success, error } from '@/components/common/base/Message'
import type { FormInstance, FormRules } from 'element-plus'

const router = useRouter()
const formRef = ref<FormInstance>()
const loading = ref(false)

const form = reactive({
  email: '',
})

const rules: FormRules = {
  email: [
    { required: true, message: '请输入邮箱地址', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' },
  ],
}

const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    loading.value = true
    try {
      const response = await userService.requestPasswordReset({
        email: form.email,
      })

      if (response.data?.success) {
        success(response.data.message || '重置密码邮件已发送，请查收')
        // 延迟后返回登录页
        setTimeout(() => {
          router.push('/auth/login')
        }, 2000)
      } else {
        error(response.data?.message || '发送失败，请稍后重试')
      }
    } catch (err: any) {
      const errorMessage = err?.response?.data?.message || err?.message || '发送失败，请稍后重试'
      error(errorMessage)
      console.error('Request password reset error:', err)
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
.forgot-password-page {
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

.forgot-password-container {
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

.forgot-password-header {
  text-align: center;
  margin-bottom: 32px;
}

.forgot-password-title {
  font-size: 28px;
  font-weight: 600;
  color: #00FF00;
  margin: 0 0 8px 0;
  text-shadow: 0 0 10px rgba(0, 255, 0, 0.5);
}

.forgot-password-subtitle {
  font-size: 14px;
  color: rgba(0, 255, 0, 0.7);
  margin: 0;
}

.forgot-password-form {
  .submit-button {
    width: 100%;
  }

  .forgot-password-footer {
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

