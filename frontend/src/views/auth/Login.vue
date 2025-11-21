<template>
  <div class="login-page auth-page">
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
            @click="handleSubmit"
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

const handleSubmit = async (event?: Event) => {
  if (event) {
    event.preventDefault()
  }
  
  if (!formRef.value) {
    console.warn('Form ref is not available')
    return
  }

  try {
    await formRef.value.validate(async (valid) => {
      if (!valid) {
        console.warn('Form validation failed')
        return
      }

      loading.value = true
      console.log('Starting login with:', { username: form.username })
      
      try {
        const result = await authStore.login(form.username, form.password)
        console.log('Login result:', result)
        
        if (result) {
          success('登录成功')
          // 等待一下确保状态已更新
          await new Promise(resolve => setTimeout(resolve, 200))
          
          // 检查认证状态
          if (!authStore.isAuthenticated) {
            console.warn('登录成功但认证状态未更新，等待重试...')
            await new Promise(resolve => setTimeout(resolve, 300))
          }
          
          // 跳转到之前访问的页面或首页
          const redirect = router.currentRoute.value.query.redirect as string
          const targetPath = redirect || '/dashboard'
          console.log('Navigating to:', targetPath, 'isAuthenticated:', authStore.isAuthenticated)
          
          try {
            // 使用 push 而不是 replace，这样用户可以返回
            await router.push(targetPath)
            console.log('Navigation successful')
          } catch (err) {
            console.error('Navigation error:', err)
            // 如果路由跳转失败，可能是路由守卫阻止了，尝试刷新页面
            if (err instanceof Error && err.message.includes('Navigation')) {
              // 路由守卫可能阻止了导航，等待一下再试
              await new Promise(resolve => setTimeout(resolve, 500))
              window.location.href = targetPath
            }
          }
        } else {
          error('登录失败，请检查用户名和密码')
        }
      } catch (err: any) {
        console.error('Login error:', err)
        // 检查是否是网络错误
        if (err?.type === 'NETWORK' || err?.code === 'ERR_NETWORK' || err?.message?.includes('Network')) {
          error('无法连接到服务器，请确保后端服务正在运行')
        } else {
          error('登录失败，请稍后重试')
        }
      } finally {
        loading.value = false
      }
    })
  } catch (err) {
    console.error('Form validation error:', err)
    error('表单验证失败，请检查输入')
  }
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
  
  // 输入框包装器样式 - 深灰色半透明背景（最外层）
  // 使用最具体的选择器确保覆盖所有状态和 Element Plus 默认样式
  :deep(.el-form-item .el-input__wrapper),
  :deep(.el-form-item .el-input .el-input__wrapper),
  :deep(.el-input__wrapper),
  :deep(.el-input.is-disabled .el-input__wrapper),
  :deep(.el-input.is-error .el-input__wrapper),
  :deep(.el-input.is-success .el-input__wrapper),
  :deep(.el-form-item.is-error .el-input__wrapper),
  :deep(.el-form-item.is-success .el-input__wrapper) {
    background-color: rgba(30, 30, 30, 0.6) !important;
    background: rgba(30, 30, 30, 0.6) !important; // 兼容性写法
    border: 1px solid rgba(0, 255, 0, 0.3) !important;
    box-shadow: 0 0 5px rgba(0, 255, 0, 0.1) !important;
    backdrop-filter: blur(8px) !important;
    
    // 内部输入元素设为透明
    .el-input__inner {
      background-color: transparent !important;
      background: transparent !important;
    }
    
    // 前缀和后缀区域也设为透明
    .el-input__prefix,
    .el-input__suffix {
      background-color: transparent !important;
      background: transparent !important;
    }
    
    // 前缀和后缀内部元素
    .el-input__prefix-inner,
    .el-input__suffix-inner {
      background-color: transparent !important;
      background: transparent !important;
    }
    
    // 所有子元素都设为透明（除了包装器本身）
    > * {
      background-color: transparent !important;
      background: transparent !important;
    }
    
    &:hover {
      background-color: rgba(40, 40, 40, 0.7) !important;
      background: rgba(40, 40, 40, 0.7) !important;
      border-color: rgba(0, 255, 0, 0.5) !important;
      box-shadow: 0 0 8px rgba(0, 255, 0, 0.2) !important;
    }
    
    &.is-focus {
      background-color: rgba(50, 50, 50, 0.75) !important;
      background: rgba(50, 50, 50, 0.75) !important;
      border-color: rgba(0, 255, 0, 0.6) !important;
      box-shadow: 0 0 10px rgba(0, 255, 0, 0.3) !important;
    }
    
    // 当输入框有内容时（非 placeholder 状态）保持深灰色背景
    &:has(.el-input__inner:not(:placeholder-shown)) {
      background-color: rgba(30, 30, 30, 0.6) !important;
      background: rgba(30, 30, 30, 0.6) !important;
    }
    
    // 聚焦且有内容时
    &.is-focus:has(.el-input__inner:not(:placeholder-shown)) {
      background-color: rgba(50, 50, 50, 0.75) !important;
      background: rgba(50, 50, 50, 0.75) !important;
    }
    
    // 错误状态也使用深灰色背景
    &.is-error {
      background-color: rgba(30, 30, 30, 0.6) !important;
      background: rgba(30, 30, 30, 0.6) !important;
      border-color: rgba(255, 0, 0, 0.5) !important;
    }
  }
  
  // 输入框内部文字样式 - 覆盖所有状态
  :deep(.el-input__inner),
  :deep(.el-input.is-disabled .el-input__inner),
  :deep(.el-input.is-error .el-input__inner),
  :deep(.el-input.is-success .el-input__inner) {
    color: rgba(0, 255, 0, 0.9) !important;
    background-color: transparent !important;
    background: transparent !important;
    
    // 当输入框有内容时（非 placeholder 状态）
    &:not(:placeholder-shown) {
      background-color: transparent !important;
      background: transparent !important;
    }
    
    // 聚焦状态
    &:focus {
      background-color: transparent !important;
      background: transparent !important;
    }
    
    &::placeholder {
      color: rgba(0, 255, 0, 0.4) !important;
    }
  }
  
  // 确保输入框组件本身透明（但不影响包装器）
  :deep(.el-input),
  :deep(.el-input.is-disabled),
  :deep(.el-input.is-error),
  :deep(.el-input.is-success) {
    background-color: transparent !important;
    background: transparent !important;
  }
  
  // 覆盖所有可能的输入框状态
  :deep(.el-input.is-disabled .el-input__wrapper),
  :deep(.el-input.is-disabled .el-input__inner),
  :deep(.el-input.is-error .el-input__wrapper),
  :deep(.el-input.is-error .el-input__inner) {
    background-color: rgba(30, 30, 30, 0.6) !important;
    background: rgba(30, 30, 30, 0.6) !important;
  }
  
  // 确保输入框的所有内部容器都没有白色背景
  :deep(.el-input__container) {
    background-color: transparent !important;
    background: transparent !important;
  }
  
  // 强制覆盖所有可能的白色背景 - 使用最高优先级
  :deep(.el-input),
  :deep(.el-input *),
  :deep(.el-input__wrapper),
  :deep(.el-input__wrapper *),
  :deep(.el-input__inner),
  :deep(.el-input__inner *) {
    // 确保没有白色背景
    &[style*="background"],
    &[style*="background-color"] {
      background-color: transparent !important;
      background: transparent !important;
    }
  }
  
  // 特别处理包装器，确保它保持深灰色背景
  :deep(.el-input__wrapper) {
    &,
    &[style*="background"],
    &[style*="background-color"] {
      background-color: rgba(30, 30, 30, 0.6) !important;
      background: rgba(30, 30, 30, 0.6) !important;
    }
  }
  
  // 输入框前缀图标颜色
  :deep(.el-input__prefix) {
    .el-input__prefix-inner {
      color: rgba(0, 255, 0, 0.6) !important;
    }
  }
  
  // 输入框后缀图标颜色（清除按钮、显示密码按钮等）
  :deep(.el-input__suffix) {
    .el-input__suffix-inner {
      color: rgba(0, 255, 0, 0.6) !important;
      
      .el-input__clear {
        color: rgba(0, 255, 0, 0.6) !important;
        
        &:hover {
          color: rgba(0, 255, 0, 0.9) !important;
        }
      }
    }
  }
  
  :deep(.el-checkbox__label) {
    color: rgba(124, 121, 121, 0.575);
  }
  
  // 复选框未选中状态的背景透明度
  :deep(.el-checkbox__inner) {
    background-color: rgba(30, 30, 30, 0.6) !important;
    border-color: rgba(0, 255, 0, 0.3) !important;
    
    &:hover {
      background-color: rgba(40, 40, 40, 0.7) !important;
      border-color: rgba(0, 255, 0, 0.5) !important;
    }
  }
  
  // 复选框选中状态的样式
  :deep(.el-checkbox__input.is-checked .el-checkbox__inner) {
    background-color: #00FF00 !important;
    border-color: #00FF00 !important;
    
    &::after {
      border-color: #000000 !important;
    }
  }
  
  // 复选框禁用状态
  :deep(.el-checkbox.is-disabled .el-checkbox__inner) {
    background-color: rgba(0, 0, 0, 0.3) !important;
    border-color: rgba(0, 255, 0, 0.2) !important;
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
