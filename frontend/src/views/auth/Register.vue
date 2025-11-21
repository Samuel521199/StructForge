<template>
  <div class="register-page auth-page">
    <CodeRain :useConfigFile="true" />
    <div class="register-container">
      <div class="register-header">
        <h1 class="register-title">创建账号</h1>
        <p class="register-subtitle">加入 StructForge</p>
      </div>

      <Form
        ref="formRef"
        :model="form"
        :rules="rules"
        class="register-form"
        @submit.prevent="handleSubmit"
      >
        <FormItem label="用户名" prop="username">
          <Input
            v-model="form.username"
            placeholder="请输入用户名（3-20个字符）"
            size="large"
            :prefix-icon="User"
            clearable
          />
        </FormItem>

        <FormItem label="邮箱" prop="email">
          <Input
            v-model="form.email"
            type="email"
            placeholder="请输入邮箱地址"
            size="large"
            :prefix-icon="Message"
            clearable
          />
        </FormItem>

        <FormItem label="密码" prop="password">
          <Input
            v-model="form.password"
            type="password"
            placeholder="密码：6-20个字符，包含字母+数字，且有大小写或特殊字符"
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
            placeholder="请再次输入密码"
            size="large"
            :prefix-icon="Lock"
            show-password
            clearable
            @keyup.enter="handleSubmit"
          />
        </FormItem>

        <FormItem>
          <div class="terms-agreement">
            <Checkbox v-model="agreeTerms" class="terms-checkbox">
              <span class="terms-text">
                <span class="terms-line">我已阅读并同意</span>
                <span class="terms-line">
                  <Link type="primary" :underline="false">《用户协议》</Link>
                  和
                  <Link type="primary" :underline="false">《隐私政策》</Link>
                </span>
              </span>
            </Checkbox>
          </div>
        </FormItem>

        <FormItem>
          <Button
            type="primary"
            size="large"
            :loading="loading"
            :disabled="!agreeTerms"
            native-type="submit"
            class="register-button"
          >
            注册
          </Button>
        </FormItem>

        <FormItem>
          <div class="register-footer">
            <span>已有账号？</span>
            <Link type="primary" @click="goToLogin">立即登录</Link>
          </div>
        </FormItem>
      </Form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { User, Lock, Message } from '@element-plus/icons-vue'
import { Form, FormItem, Input, Button, Checkbox, Link } from '@/components/common/base'
import { CodeRain } from '@/components/common/effects'
import { success, error } from '@/components/common/base/Message'
import type { FormInstance, FormRules } from 'element-plus'
import { userService, type RegisterResponse } from '@/api/services/user.service'

const router = useRouter()

const formRef = ref<FormInstance>()
const loading = ref(false)
const agreeTerms = ref(false)

const form = reactive({
  username: '',
  email: '',
  password: '',
  confirmPassword: '',
})

// 验证密码复杂度
const validatePassword = (_rule: any, value: string, callback: any) => {
  if (!value) {
    callback(new Error('请输入密码'))
    return
  }

  if (value.length < 6 || value.length > 20) {
    callback(new Error('密码长度为 6-20 个字符'))
    return
  }

  // 检查是否包含字母和数字
  const hasLetter = /[a-zA-Z]/.test(value)
  const hasNumber = /[0-9]/.test(value)

  if (!hasLetter || !hasNumber) {
    callback(new Error('密码必须包含字母和数字'))
    return
  }

  // 检查是否有大小写或特殊字符
  const hasUpperCase = /[A-Z]/.test(value)
  const hasLowerCase = /[a-z]/.test(value)
  const hasSpecialChar = /[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]/.test(value)

  if (!hasUpperCase && !hasLowerCase && !hasSpecialChar) {
    callback(new Error('密码必须包含大小写字母或特殊字符'))
    return
  }

  // 如果只有小写字母，必须有特殊字符
  if (hasLowerCase && !hasUpperCase && !hasSpecialChar) {
    callback(new Error('密码必须包含大写字母或特殊字符'))
    return
  }

  // 如果只有大写字母，必须有特殊字符
  if (hasUpperCase && !hasLowerCase && !hasSpecialChar) {
    callback(new Error('密码必须包含小写字母或特殊字符'))
    return
  }

  callback()
}

// 验证确认密码
const validateConfirmPassword = (_rule: any, value: string, callback: any) => {
  if (!value) {
    callback(new Error('请再次输入密码'))
    return
  }

  if (value !== form.password) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

const rules: FormRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度为 3-20 个字符', trigger: 'blur' },
    { pattern: /^[a-zA-Z0-9_]+$/, message: '用户名只能包含字母、数字和下划线', trigger: 'blur' },
  ],
  email: [
    { required: true, message: '请输入邮箱地址', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' },
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { validator: validatePassword, trigger: 'blur' },
  ],
  confirmPassword: [
    { required: true, message: '请再次输入密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' },
  ],
}

const handleSubmit = async () => {
  console.log('[Register] 注册按钮被点击')
  
  if (!formRef.value) {
    console.error('[Register] formRef 为空')
    return
  }

  console.log('[Register] 开始验证表单')
  
  try {
    // Element Plus 的 validate 方法返回 Promise
    await formRef.value.validate()
    console.log('[Register] 表单验证通过')
  } catch (validationError) {
    console.warn('[Register] 表单验证失败:', validationError)
    return
  }

  if (!agreeTerms.value) {
    console.warn('[Register] 未同意用户协议')
    error('请先同意用户协议和隐私政策')
    return
  }

  console.log('[Register] 准备发送注册请求:', {
    username: form.username,
    email: form.email,
    passwordLength: form.password.length,
  })

  loading.value = true
  try {
    // 调用注册API
    console.log('[Register] 调用 userService.register')
    const response = await userService.register({
      username: form.username,
      email: form.email,
      password: form.password,
    })

      console.log('[Register] 注册响应:', response)

      // Axios 响应拦截器已经返回 response.data，所以 response 就是 RegisterResponse
      // 包含 success、message、user 字段
      // 注意：userService.register 返回 Promise<ApiResponse<RegisterResponse>>
      // 但响应拦截器返回 response.data，所以 response 的类型实际上是 RegisterResponse
      const registerResponse = response as unknown as RegisterResponse
      
      if (registerResponse.success) {
        console.log('[Register] 注册成功')
        success('注册成功！请查收邮件验证您的邮箱', {
          duration: 3000,
          showClose: true,
        })
        // 注册成功后跳转到登录页，并提示验证邮箱
        setTimeout(() => {
          router.push('/auth/login?registered=true')
        }, 2000)
      } else {
        console.warn('[Register] 注册失败:', registerResponse.message)
        const errorMsg = registerResponse.message || '注册失败，请稍后重试'
        error(errorMsg, {
          duration: 5000,
          showClose: true,
        })
      }
  } catch (err: any) {
    console.error('[Register] 注册异常:', err)
    console.error('[Register] 错误详情:', {
      message: err?.message,
      response: err?.response,
      status: err?.response?.status,
      data: err?.response?.data,
    })
    
    // 处理不同类型的错误
    let errorMessage = '注册失败，请稍后重试'
    
    if (err?.response?.data) {
      // 后端返回的错误响应
      if (err.response.data.message) {
        errorMessage = err.response.data.message
      } else if (typeof err.response.data === 'string') {
        errorMessage = err.response.data
      }
    } else if (err?.message) {
      // 网络错误或其他错误
      if (err.message.includes('Network Error')) {
        errorMessage = '网络连接失败，请检查网络后重试'
      } else if (err.message.includes('timeout')) {
        errorMessage = '请求超时，请稍后重试'
      } else {
        errorMessage = err.message
      }
    }
    
    error(errorMessage, {
      duration: 5000,
      showClose: true,
    })
  } finally {
    loading.value = false
    console.log('[Register] 注册流程结束')
  }
}

const goToLogin = () => {
  router.push('/auth/login')
}
</script>

<style scoped lang="scss">
.register-page {
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

.register-container {
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

.register-header {
  text-align: center;
  margin-bottom: 32px;
}

.register-title {
  font-size: 28px;
  font-weight: 600;
  color: #00FF00;
  margin: 0 0 8px 0;
  text-shadow: 0 0 10px rgba(0, 255, 0, 0.5);
}

.register-subtitle {
  font-size: 14px;
  color: rgba(0, 255, 0, 0.7);
  margin: 0;
}

.register-form {
  .register-button {
    width: 100%;
  }

  .register-footer {
    text-align: center;
    font-size: 14px;
    color: rgba(85, 81, 81, 0.562);

    // Link 组件样式已通过全局样式统一，无需额外样式
  }

  :deep(.el-form-item__label) {
    color: rgba(145, 137, 137, 0.551);
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
    color: rgba(131, 128, 128, 0.582);
    line-height: 1.5;
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

  .terms-agreement {
    width: 100%;
  }

  .terms-checkbox {
    width: 100%;
    align-items: flex-start;

    :deep(.el-checkbox__label) {
      width: calc(100% - 20px);
      white-space: normal;
      word-wrap: break-word;
    }
  }

  .terms-text {
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  .terms-line {
    display: block;
    line-height: 1.5;
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

    &:disabled {
      background-color: rgba(0, 255, 0, 0.3);
      border-color: rgba(0, 255, 0, 0.3);
      color: rgba(0, 0, 0, 0.5);
    }
  }
}
</style>
