<template>
  <div class="verify-email-page">
    <CodeRain :useConfigFile="true" />
    <div class="verify-email-container">
      <div class="verify-email-header">
        <h1 class="verify-email-title">邮箱验证</h1>
        <p class="verify-email-subtitle">{{ statusMessage }}</p>
      </div>

      <div v-if="verifying" class="verifying">
        <Icon :icon="LoadingIcon" :size="40" :is-loading="true" />
        <p>正在验证邮箱...</p>
      </div>

      <div v-else-if="isSuccess" class="success">
        <Icon :icon="CircleCheckIcon" :size="60" color="#00FF00" />
        <p>邮箱验证成功！</p>
        <Button type="primary" @click="goToLogin">前往登录</Button>
      </div>

      <div v-else class="error">
        <Icon :icon="CircleCloseIcon" :size="60" color="#ff4d4f" />
        <p>{{ errorMessage }}</p>
        <div class="actions">
          <Button @click="resendEmail">重新发送验证邮件</Button>
          <Button type="primary" @click="goToLogin">返回登录</Button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { Loading, CircleCheck, CircleClose } from '@element-plus/icons-vue'
import { Button, Icon } from '@/components/common/base'
import { CodeRain } from '@/components/common/effects'
import { userService } from '@/api/services/user.service'
import { success, error } from '@/components/common/base/Message'

const LoadingIcon = Loading
const CircleCheckIcon = CircleCheck
const CircleCloseIcon = CircleClose

const router = useRouter()
const route = useRoute()

const verifying = ref(true)
const isSuccess = ref(false)
const statusMessage = ref('正在验证您的邮箱...')
const errorMessage = ref('')
const email = ref<string>('')

onMounted(async () => {
  const token = route.query.token as string
  if (!token) {
    verifying.value = false
    errorMessage.value = '验证令牌无效'
    return
  }

  try {
    const response = await userService.verifyEmail({ token })
    if (response.data?.success) {
      verifying.value = false
      isSuccess.value = true
      statusMessage.value = '邮箱验证成功！'
      success('邮箱验证成功')
    } else {
      verifying.value = false
      errorMessage.value = response.data?.message || '验证失败'
      error(response.data?.message || '验证失败')
    }
  } catch (err: any) {
    verifying.value = false
    const message = err?.response?.data?.message || err?.message || '验证失败，请稍后重试'
    errorMessage.value = message
    error(message)
    console.error('Verify email error:', err)
  }
})

const resendEmail = async () => {
  if (!email.value) {
    error('请先输入邮箱地址')
    return
  }

  try {
    const response = await userService.resendVerificationEmail({ email: email.value })
    if (response.data?.success) {
      success('验证邮件已重新发送，请查收')
    } else {
      error(response.data?.message || '发送失败，请稍后重试')
    }
  } catch (err: any) {
    const errorMessage = err?.response?.data?.message || err?.message || '发送失败，请稍后重试'
    error(errorMessage)
    console.error('Resend verification email error:', err)
  }
}

const goToLogin = () => {
  router.push('/auth/login')
}
</script>

<style scoped lang="scss">
.verify-email-page {
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

.verify-email-container {
  position: relative;
  z-index: 10;
  width: 100%;
  max-width: 500px;
  background-color: rgba(0, 0, 0, 0.85);
  border: 1px solid rgba(0, 255, 0, 0.3);
  border-radius: 8px;
  box-shadow:
    0 0 40px rgba(0, 255, 0, 0.5),
    inset 0 0 10px rgba(0, 255, 0, 0.2);
  padding: 40px;
  backdrop-filter: blur(10px);
  text-align: center;
}

.verify-email-header {
  margin-bottom: 32px;
}

.verify-email-title {
  font-size: 28px;
  font-weight: 600;
  color: #00FF00;
  margin: 0 0 8px 0;
  text-shadow: 0 0 10px rgba(0, 255, 0, 0.5);
}

.verify-email-subtitle {
  font-size: 14px;
  color: rgba(0, 255, 0, 0.7);
  margin: 0;
}

.verifying,
.success,
.error {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 20px;
  padding: 20px 0;

  p {
    color: rgba(255, 255, 255, 0.9);
    font-size: 16px;
    margin: 0;
  }
}

.actions {
  display: flex;
  gap: 12px;
  margin-top: 10px;
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
</style>

