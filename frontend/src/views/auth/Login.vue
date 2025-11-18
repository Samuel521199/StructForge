<template>
  <div class="login-page">
    <CodeRain
      :color="'#00FF00'"
      :backgroundColor="'#000000'"
      :fontSize="14"
      :fontWeight="'bold'"
      :speed="2.5"
      :speedVariation="0.6"
      :density="0.003"
      :opacity="0.85"
      :fadeSpeed="0.08"
      :minLength="15"
      :maxLength="35"
      :enableLayers="true"
      :enableGlow="true"
      :enableGlitch="true"
      :glowIntensity="0.4"
    />
    <div class="login-container">
      <div class="login-header">
        <h1 class="login-title">欢迎回来</h1>
        <p class="login-subtitle">登录到 StructForge</p>
      </div>

      <Form
        ref="formRef"
        :model="form"
        :rules="rules"
        class="login-form"
        @submit.prevent="handleSubmit"
      >
        <FormItem label="用户名" prop="username">
          <Input
            v-model="form.username"
            placeholder="请输入用户名"
            size="large"
            :prefix-icon="User"
            clearable
          />
        </FormItem>

        <FormItem label="密码" prop="password">
          <Input
            v-model="form.password"
            type="password"
            placeholder="请输入密码"
            size="large"
            :prefix-icon="Lock"
            show-password
            clearable
            @keyup.enter="handleSubmit"
          />
        </FormItem>

        <FormItem>
          <div class="login-options">
            <el-checkbox v-model="rememberMe">记住我</el-checkbox>
            <el-link type="primary" :underline="false">忘记密码？</el-link>
          </div>
        </FormItem>

        <FormItem>
          <Button
            type="primary"
            size="large"
            :loading="loading"
            native-type="submit"
            class="login-button"
          >
            登录
          </Button>
        </FormItem>

        <FormItem>
          <div class="login-footer">
            <span>还没有账号？</span>
            <el-link type="primary" @click="goToRegister">立即注册</el-link>
          </div>
        </FormItem>
      </Form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { User, Lock } from '@element-plus/icons-vue'
import { Form, FormItem, Input, Button } from '@/components/common/base'
import { CodeRain } from '@/components/common/effects'
import { useAuthStore } from '@/stores/modules/auth.store'
import { success, error } from '@/components/common/base/Message'
import type { FormInstance } from 'element-plus'

const router = useRouter()
const authStore = useAuthStore()

const formRef = ref<FormInstance>()
const loading = ref(false)
const rememberMe = ref(false)

const form = reactive({
  username: '',
  password: '',
})

const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' as const },
    { min: 3, max: 20, message: '用户名长度为 3-20 个字符', trigger: 'blur' as const },
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' as const },
    { min: 6, max: 20, message: '密码长度为 6-20 个字符', trigger: 'blur' as const },
  ],
}

const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    loading.value = true
    try {
      const result = await authStore.login(form.username, form.password)
      if (result) {
        success('登录成功')
        // 跳转到之前访问的页面或首页
        const redirect = router.currentRoute.value.query.redirect as string
        router.push(redirect || '/')
      } else {
        error('登录失败，请检查用户名和密码')
      }
    } catch (err) {
      error('登录失败，请稍后重试')
      console.error('Login error:', err)
    } finally {
      loading.value = false
    }
  })
}

const goToRegister = () => {
  router.push('/auth/register')
}
</script>

<style scoped lang="scss">
.login-page {
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
  /* 确保背景是纯黑色，移除任何干扰 */
}

.login-container {
  position: relative;
  z-index: 2;
  width: 100%;
  max-width: 400px;
  background-color: rgba(0, 0, 0, 0.8);
  border: 1px solid rgba(0, 255, 0, 0.3);
  border-radius: 8px;
  box-shadow: 
    0 0 40px rgba(0, 255, 0, 0.5),
    inset 0 0 10px rgba(0, 255, 0, 0.2);
  padding: 40px;
  backdrop-filter: blur(10px);
}

.login-header {
  text-align: center;
  margin-bottom: 32px;
}

.login-title {
  font-size: 28px;
  font-weight: 600;
  color: #00FF00;
  margin: 0 0 8px 0;
  text-shadow: 0 0 10px rgba(0, 255, 0, 0.5);
}

.login-subtitle {
  font-size: 14px;
  color: rgba(0, 255, 0, 0.7);
  margin: 0;
}

.login-form {
  .login-options {
    display: flex;
    justify-content: space-between;
    align-items: center;
    width: 100%;
  }

  .login-button {
    width: 100%;
  }

  .login-footer {
    text-align: center;
    font-size: 14px;
    color: rgba(255, 255, 255, 0.7);

    .el-link {
      margin-left: 4px;
      color: #00FF00;
      
      &:hover {
        color: #00FF00;
        text-shadow: 0 0 5px rgba(0, 255, 0, 0.5);
      }
    }
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
  
  :deep(.el-checkbox__label) {
    color: rgba(255, 255, 255, 0.7);
  }
  
  :deep(.el-checkbox__input.is-checked .el-checkbox__inner) {
    background-color: #00FF00;
    border-color: #00FF00;
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
