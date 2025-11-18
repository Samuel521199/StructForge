<template>
  <div class="register-page">
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
            placeholder="请输入密码（6-20个字符）"
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
          <el-checkbox v-model="agreeTerms">
            我已阅读并同意
            <el-link type="primary" :underline="false">《用户协议》</el-link>
            和
            <el-link type="primary" :underline="false">《隐私政策》</el-link>
          </el-checkbox>
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
            <el-link type="primary" @click="goToLogin">立即登录</el-link>
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
import { Form, FormItem, Input, Button } from '@/components/common/base'
import { useAuthStore } from '@/stores/modules/auth.store'
import { success, error } from '@/components/common/base/Message'
import type { FormInstance, FormRules } from 'element-plus'
import { userService } from '@/api/services/user.service'

const router = useRouter()
const authStore = useAuthStore()

const formRef = ref<FormInstance>()
const loading = ref(false)
const agreeTerms = ref(false)

const form = reactive({
  username: '',
  email: '',
  password: '',
  confirmPassword: '',
})

// 验证确认密码
const validateConfirmPassword = (rule: any, value: string, callback: any) => {
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
    { min: 6, max: 20, message: '密码长度为 6-20 个字符', trigger: 'blur' },
  ],
  confirmPassword: [
    { required: true, message: '请再次输入密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' },
  ],
}

const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    if (!agreeTerms.value) {
      error('请先同意用户协议和隐私政策')
      return
    }

    loading.value = true
    try {
      // 调用注册API
      const response = await userService.register({
        username: form.username,
        email: form.email,
        password: form.password,
      })

      if (response.data) {
        success('注册成功，请登录')
        // 注册成功后跳转到登录页
        router.push('/auth/login')
      } else {
        error('注册失败，请稍后重试')
      }
    } catch (err: any) {
      error(err?.message || '注册失败，请稍后重试')
      console.error('Register error:', err)
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
.register-page {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}

.register-container {
  width: 100%;
  max-width: 400px;
  background: var(--el-bg-color);
  border-radius: 8px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
  padding: 40px;
}

.register-header {
  text-align: center;
  margin-bottom: 32px;
}

.register-title {
  font-size: 28px;
  font-weight: 600;
  color: var(--el-text-color-primary);
  margin: 0 0 8px 0;
}

.register-subtitle {
  font-size: 14px;
  color: var(--el-text-color-regular);
  margin: 0;
}

.register-form {
  .register-button {
    width: 100%;
  }

  .register-footer {
    text-align: center;
    font-size: 14px;
    color: var(--el-text-color-regular);

    .el-link {
      margin-left: 4px;
    }
  }
}
</style>
