<template>
  <header class="app-header">
    <div class="header-left">
      <div class="logo">
        <h1 class="logo-text">StructForge</h1>
      </div>
    </div>
    
    <div class="header-center">
      <slot name="center" />
    </div>
    
    <div class="header-right">
      <div class="user-info">
        <Avatar 
          v-if="userStore.user" 
          :src="userStore.user.avatar" 
          :text="userStore.username || 'U'"
          :size="32"
        />
        <div class="user-menu">
          <Dropdown>
            <template #trigger>
              <Button type="text" class="user-menu-btn">
                <span class="username">{{ userStore.username || '用户' }}</span>
                <Icon :icon="ArrowDown" :size="14" />
              </Button>
            </template>
            <template #dropdown>
              <div class="dropdown-menu">
                <div class="menu-item" @click="goToProfile">
                  <Icon :icon="User" :size="16" />
                  <span>个人资料</span>
                </div>
                <div class="menu-item" @click="goToSettings">
                  <Icon :icon="Setting" :size="16" />
                  <span>设置</span>
                </div>
                <div class="menu-divider" />
                <div class="menu-item" @click="handleLogout">
                  <Icon :icon="SwitchButton" :size="16" />
                  <span>退出登录</span>
                </div>
              </div>
            </template>
          </Dropdown>
        </div>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/modules/auth.store'
import { useUserStore } from '@/stores/modules/user.store'
import { Avatar, Button, Icon, Dropdown } from '@/components/common/base'
import { User, Setting, SwitchButton, ArrowDown } from '@element-plus/icons-vue'
import { success } from '@/components/common/base/Message'

const router = useRouter()
const authStore = useAuthStore()
const userStore = useUserStore()

const goToProfile = () => {
  router.push('/user/profile')
}

const goToSettings = () => {
  router.push('/user/settings')
}

const handleLogout = async () => {
  await authStore.logout()
  success('已退出登录')
  router.push('/auth/login')
}
</script>

<style scoped lang="scss">
@use '@/assets/styles/glassmorphism' as *;

.app-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 64px;
  padding: 0 32px;
  @include glassmorphism(30px, 180%, rgba(20, 20, 30, 0.4));
  border-bottom: 1px solid $glass-border-cyan;
  box-shadow: 
    0 2px 20px rgba(0, 0, 0, 0.3),
    0 0 30px rgba(0, 212, 255, 0.1) inset;
  position: relative;
  z-index: 10;
  
  // 底部光晕线
  &::after {
    content: '';
    position: absolute;
    bottom: -1px;
    left: 0;
    right: 0;
    height: 2px;
    background: linear-gradient(90deg, 
      transparent 0%,
      rgba(0, 212, 255, 0.5) 20%,
      rgba(0, 255, 136, 0.6) 50%,
      rgba(183, 148, 246, 0.5) 80%,
      transparent 100%
    );
    box-shadow: 0 0 15px rgba(0, 212, 255, 0.4);
  }

  .header-left {
    display: flex;
    align-items: center;
    flex-shrink: 0;

    .logo {
      display: flex;
      align-items: center;
      gap: 12px;

      .logo-text {
        margin: 0;
        font-size: 22px;
        font-weight: 700;
        background: linear-gradient(135deg, #ffffff 0%, #00d4ff 50%, #00ff88 100%);
        -webkit-background-clip: text;
        -webkit-text-fill-color: transparent;
        background-clip: text;
        text-shadow: 0 0 20px rgba(0, 212, 255, 0.3);
        letter-spacing: 1px;
      }
    }
  }

  .header-center {
    flex: 1;
    display: flex;
    justify-content: center;
    align-items: center;
  }

  .header-right {
    display: flex;
    align-items: center;
    gap: 16px;
    flex-shrink: 0;

    .user-info {
      display: flex;
      align-items: center;
      gap: 12px;

      .user-menu {
        .user-menu-btn {
          display: flex;
          align-items: center;
          gap: 8px;
          padding: 4px 12px;

          .username {
            font-size: 14px;
            color: rgba(255, 255, 255, 0.9);
            font-weight: 500;
          }
        }
      }
    }
  }
}

.dropdown-menu {
  min-width: 160px;
  padding: 4px 0;

  .menu-item {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 8px 16px;
    cursor: pointer;
    transition: background-color 0.2s;
    font-size: 14px;
    color: var(--el-text-color-primary);

    &:hover {
      background-color: var(--el-fill-color-light);
    }

    span {
      flex: 1;
    }
  }

  .menu-divider {
    height: 1px;
    background-color: var(--el-border-color-light);
    margin: 4px 0;
  }
}
</style>
