<template>
  <header class="sf-header">
    <div class="sf-header__container">
      <!-- Logo -->
      <div class="sf-header__logo" @click="goHome">
        <span class="sf-header__logo-text">StructForge</span>
      </div>

      <!-- 导航菜单（可选） -->
      <nav v-if="showNav" class="sf-header__nav">
        <slot name="nav" />
      </nav>

      <!-- 右侧操作区 -->
      <div class="sf-header__actions">
        <!-- 未登录状态 -->
        <template v-if="!isLoggedIn">
          <Button type="primary" @click="handleLogin">
            登录
          </Button>
        </template>

        <!-- 已登录状态 -->
        <template v-else>
          <!-- 用户信息 -->
          <div class="sf-header__user" @click="handleUserClick">
            <el-avatar
              v-if="user?.avatar"
              :src="user.avatar"
              :size="32"
            />
            <el-avatar v-else :size="32">
              {{ userDisplayName }}
            </el-avatar>
            <span class="sf-header__username">{{ userDisplayName }}</span>
            <el-icon class="sf-header__dropdown-icon">
              <ArrowDown />
            </el-icon>
          </div>

          <!-- 用户下拉菜单 -->
          <el-dropdown
            ref="dropdownRef"
            trigger="click"
            placement="bottom-end"
            @command="handleCommand"
          >
            <span class="sf-header__dropdown-trigger" />
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="profile">
                  <el-icon><User /></el-icon>
                  个人中心
                </el-dropdown-item>
                <el-dropdown-item command="settings">
                  <el-icon><Setting /></el-icon>
                  设置
                </el-dropdown-item>
                <el-dropdown-item divided command="logout">
                  <el-icon><SwitchButton /></el-icon>
                  退出登录
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </template>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ArrowDown, User, Setting, SwitchButton } from '@element-plus/icons-vue'
import { Button } from '@/components/common/base'
import { useAuthStore } from '@/stores/modules/auth.store'
import { useUserStore } from '@/stores/modules/user.store'

interface HeaderProps {
  showNav?: boolean
}

withDefaults(defineProps<HeaderProps>(), {
  showNav: false,
})

const router = useRouter()
const authStore = useAuthStore()
const userStore = useUserStore()

const dropdownRef = ref()

// 计算属性
const isLoggedIn = computed(() => userStore.isLoggedIn)
const user = computed(() => userStore.user)
const userDisplayName = computed(() => {
  if (!user.value) return '用户'
  return user.value.username || user.value.email || '用户'
})

// 方法
const goHome = () => {
  router.push('/')
}

const handleLogin = () => {
  router.push('/auth/login')
}

const handleUserClick = () => {
  // 点击用户信息区域，显示下拉菜单
  dropdownRef.value?.handleClick()
}

const handleCommand = (command: string) => {
  switch (command) {
    case 'profile':
      router.push('/user/profile')
      break
    case 'settings':
      router.push('/user/settings')
      break
    case 'logout':
      handleLogout()
      break
  }
}

const handleLogout = async () => {
  await authStore.logout()
  router.push('/auth/login')
}
</script>

<style scoped lang="scss">
.sf-header {
  height: 60px;
  background: var(--el-bg-color);
  border-bottom: 1px solid var(--el-border-color-light);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);

  &__container {
    display: flex;
    align-items: center;
    justify-content: space-between;
    height: 100%;
    padding: 0 24px;
    max-width: 100%;
  }

  &__logo {
    display: flex;
    align-items: center;
    cursor: pointer;
    user-select: none;

    &-text {
      font-size: 20px;
      font-weight: 600;
      color: var(--el-color-primary);
      background: linear-gradient(135deg, var(--el-color-primary), var(--el-color-primary-light-3));
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
      background-clip: text;
    }
  }

  &__nav {
    flex: 1;
    display: flex;
    align-items: center;
    margin-left: 40px;
  }

  &__actions {
    display: flex;
    align-items: center;
    gap: 16px;
  }

  &__user {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 4px 12px;
    border-radius: 20px;
    cursor: pointer;
    transition: background-color 0.2s;
    user-select: none;

    &:hover {
      background-color: var(--el-bg-color-page);
    }
  }

  &__username {
    font-size: 14px;
    color: var(--el-text-color-primary);
    max-width: 120px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  &__dropdown-icon {
    font-size: 12px;
    color: var(--el-text-color-secondary);
    transition: transform 0.2s;
  }

  &__dropdown-trigger {
    position: absolute;
    width: 0;
    height: 0;
    opacity: 0;
  }
}
</style>

