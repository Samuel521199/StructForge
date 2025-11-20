<template>
  <div class="login-page">
    <!-- 
      使用配置文件：修改 config.json 文件即可调整效果参数
      配置文件位置：frontend/src/components/common/effects/CodeRain/config.json
      修改配置文件后，需要刷新页面才能看到效果变化
      
      如果需要覆盖配置文件中的特定参数，可以添加 props：
      <CodeRain :useConfigFile="true" :speed="3.0" :glowIntensity="0.9" />
    -->
    <CodeRain
      :useConfigFile="true"
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
            <Checkbox v-model="rememberMe">记住我</Checkbox>
            <Link type="primary" :underline="false" @click="goToForgotPassword">忘记密码？</Link>
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
            <Link type="primary" @click="goToRegister">立即注册</Link>
          </div>
        </FormItem>
      </Form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { User, Lock } from '@element-plus/icons-vue'
import { Form, FormItem, Input, Button, Checkbox, Link } from '@/components/common/base'
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

const goToForgotPassword = () => {
  router.push('/auth/forgot-password')
}

// 检查是否有注册成功的提示
onMounted(() => {
  const registered = router.currentRoute.value.query.registered
  if (registered === 'true') {
    success('注册成功！请查收邮件验证您的邮箱')
  }
})
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
  /* 确保在 Canvas 之上 */
  position: relative;
  z-index: 10; /* 高于 Canvas 的 z-index: 1 */
  width: 100%;
  max-width: 400px;
  /* 登录框主体背景微透明，让代码雨的光影透过来 */
  background-color: rgba(0, 0, 0, 0.85);
  border: 1px solid rgba(0, 255, 0, 0.3);
  border-radius: 8px;
  /* 增加绿色光晕，与代码雨主题呼应 */
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
