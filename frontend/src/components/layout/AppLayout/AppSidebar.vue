<template>
  <aside class="app-sidebar">
    <!-- Logo 区域 -->
    <div class="sidebar-logo">
      <div class="logo-icon">
        <Icon :icon="Connection" :size="28" />
      </div>
      <div class="logo-text">
        <span class="logo-title">StructForge</span>
        <span class="logo-subtitle">AI Workflow</span>
      </div>
    </div>

    <!-- 菜单导航 -->
    <nav class="sidebar-nav">
      <MenuItem
        v-for="route in menuRoutes"
        :key="route.path"
        :label="getRouteLabel(route)"
        :icon="getIcon(route.meta?.icon as string)"
        :icon-type="getIconType(route.meta?.icon as string)"
        :active="isActive(route)"
        @click="handleNavClick(route)"
      />
    </nav>

    <!-- 底部装饰 -->
    <div class="sidebar-footer">
      <div class="footer-decoration"></div>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { Icon, MenuItem } from '@/components/common/base'
import type { RouteRecordRaw } from 'vue-router'
import {
  House,
  Connection,
  User,
  Setting,
  Cpu,
  DataAnalysis,
} from '@element-plus/icons-vue'

const router = useRouter()
const route = useRoute()

// 从路由配置中获取菜单项（只获取一级菜单）
const menuRoutes = computed(() => {
  const routes = router.getRoutes()
  return routes.filter((r) => {
    return (
      r.meta &&
      !r.meta.hideInMenu &&
      r.meta.requiresAuth &&
      r.path !== '/' &&
      !r.path.startsWith('/auth') &&
      // 只显示父级路由，不显示子路由
      (!r.children || r.children.length === 0 || r.redirect)
    )
  })
})

const isActive = (routeItem: RouteRecordRaw) => {
  return route.path.startsWith(routeItem.path)
}

const handleNavClick = (routeItem: RouteRecordRaw) => {
  if (routeItem.redirect) {
    router.push(routeItem.redirect as string)
  } else if (routeItem.path) {
    router.push(routeItem.path)
  }
}

const getIcon = (iconName: string) => {
  const iconMap: Record<string, unknown> = {
    dashboard: DataAnalysis,
    workflow: Connection,
    user: User,
    system: Setting,
    ai: Cpu,
  }
  return iconMap[iconName] || House
}

const getIconType = (iconName: string): 'dashboard' | 'workflow' | 'user' | 'ai' | 'system' | 'default' => {
  const typeMap: Record<string, 'dashboard' | 'workflow' | 'user' | 'ai' | 'system' | 'default'> = {
    dashboard: 'dashboard',
    workflow: 'workflow',
    user: 'user',
    system: 'system',
    ai: 'ai',
  }
  return typeMap[iconName] || 'default'
}

const getRouteLabel = (route: RouteRecordRaw): string => {
  if (route.meta?.title) {
    return route.meta.title as string
  }
  if (typeof route.name === 'string') {
    return route.name
  }
  return route.path || 'Unknown'
}
</script>

<style scoped lang="scss">
.app-sidebar {
  width: 280px;
  // 深海蓝与星云紫渐变背景 + 噪点纹理
  background: 
    // 噪点纹理层（使用伪元素实现）
    linear-gradient(180deg, rgba(10, 10, 25, 0.98) 0%, rgba(20, 10, 35, 0.95) 50%, rgba(10, 15, 30, 0.98) 100%),
    // 星云紫渐变
    radial-gradient(ellipse at 0% 0%, rgba(147, 51, 234, 0.15) 0%, transparent 60%),
    // 深海蓝渐变
    radial-gradient(ellipse at 100% 100%, rgba(30, 58, 138, 0.2) 0%, transparent 60%),
    // 蓝绿光晕
    radial-gradient(ellipse at 50% 50%, rgba(0, 212, 255, 0.08) 0%, transparent 70%);
  border-right: 1px solid rgba(0, 212, 255, 0.3);
  overflow-y: auto;
  overflow-x: hidden;
  height: 100%;
  position: relative;
  // 玻璃拟物效果 - 25-40px blur
  backdrop-filter: blur(30px) saturate(180%);
  -webkit-backdrop-filter: blur(30px) saturate(180%);
  box-shadow: 
    4px 0 40px rgba(0, 0, 0, 0.4),
    inset -1px 0 0 rgba(0, 212, 255, 0.1);
  
  // 添加噪点纹理（使用伪元素）
  &::before {
    content: '';
    position: absolute;
    inset: 0;
    background-image: 
      repeating-linear-gradient(0deg, rgba(255, 255, 255, 0.02) 0px, transparent 1px, transparent 2px, rgba(255, 255, 255, 0.02) 3px),
      repeating-linear-gradient(90deg, rgba(255, 255, 255, 0.02) 0px, transparent 1px, transparent 2px, rgba(255, 255, 255, 0.02) 3px);
    pointer-events: none;
    opacity: 0.3;
    z-index: 0;
  }

  // 滚动条样式 - 霓虹光晕效果
  &::-webkit-scrollbar {
    width: 6px;
  }

  &::-webkit-scrollbar-track {
    background: rgba(0, 0, 0, 0.3);
    border-radius: 3px;
  }

  &::-webkit-scrollbar-thumb {
    background: linear-gradient(180deg, rgba(0, 212, 255, 0.5) 0%, rgba(147, 51, 234, 0.5) 100%);
    border-radius: 3px;
    box-shadow: 0 0 10px rgba(0, 212, 255, 0.4);

    &:hover {
      background: linear-gradient(180deg, rgba(0, 212, 255, 0.7) 0%, rgba(147, 51, 234, 0.7) 100%);
      box-shadow: 0 0 15px rgba(0, 212, 255, 0.6);
    }
  }

  // Logo 区域 - 高科技立体图形
  .sidebar-logo {
    display: flex;
    align-items: center;
    gap: 14px;
    padding: 28px 24px;
    border-bottom: 1px solid rgba(0, 212, 255, 0.3);
    margin-bottom: 20px;
    position: relative;
    z-index: 1;

    &::after {
      content: '';
      position: absolute;
      bottom: -1px;
      left: 0;
      width: 80px;
      height: 3px;
      background: linear-gradient(90deg, #00d4ff 0%, #00ff88 50%, transparent 100%);
      box-shadow: 0 0 10px rgba(0, 212, 255, 0.5);
      animation: logoLineShimmer 3s ease-in-out infinite;
    }

    .logo-icon {
      width: 52px;
      height: 52px;
      border-radius: 14px;
      background: linear-gradient(135deg, rgba(0, 212, 255, 0.35) 0%, rgba(0, 255, 136, 0.35) 100%);
      display: flex;
      align-items: center;
      justify-content: center;
      color: #00d4ff;
      box-shadow: 
        0 0 30px rgba(0, 212, 255, 0.4),
        0 4px 20px rgba(0, 212, 255, 0.3),
        inset 0 0 20px rgba(0, 212, 255, 0.2);
      position: relative;
      overflow: visible;

      // 外圈霓虹光晕
      &::before {
        content: '';
        position: absolute;
        inset: -4px;
        border-radius: 16px;
        padding: 2px;
        background: linear-gradient(135deg, #00d4ff, #00ff88, #b794f6);
        -webkit-mask: linear-gradient(#fff 0 0) content-box, linear-gradient(#fff 0 0);
        -webkit-mask-composite: xor;
        mask: linear-gradient(#fff 0 0) content-box, linear-gradient(#fff 0 0);
        mask-composite: exclude;
        opacity: 0.6;
        filter: blur(2px);
        animation: logoBorderGlow 3s ease-in-out infinite;
      }
      
      // 内部光效
      &::after {
        content: '';
        position: absolute;
        inset: 2px;
        border-radius: 12px;
        background: linear-gradient(135deg, rgba(0, 212, 255, 0.1) 0%, transparent 100%);
        pointer-events: none;
      }
    }

    .logo-text {
      display: flex;
      flex-direction: column;
      gap: 4px;

      .logo-title {
        font-size: 20px;
        font-weight: 800;
        background: linear-gradient(135deg, #ffffff 0%, #00d4ff 50%, #00ff88 100%);
        -webkit-background-clip: text;
        -webkit-text-fill-color: transparent;
        background-clip: text;
        letter-spacing: 1.5px;
        text-shadow: 0 0 20px rgba(0, 212, 255, 0.3);
        position: relative;
      }

      .logo-subtitle {
        font-size: 10px;
        color: rgba(0, 212, 255, 0.7);
        letter-spacing: 3px;
        text-transform: uppercase;
        font-weight: 600;
        text-shadow: 0 0 10px rgba(0, 212, 255, 0.4);
      }
    }
  }

  // 菜单导航
  .sidebar-nav {
    padding: 12px 0;
    flex: 1;
    position: relative;
    z-index: 1;
  }

  // 底部装饰
  .sidebar-footer {
    padding: 24px 20px;
    border-top: 1px solid rgba(0, 212, 255, 0.3);
    margin-top: auto;
    position: relative;
    z-index: 1;

    .footer-decoration {
      height: 3px;
      background: linear-gradient(90deg, 
        transparent 0%, 
        rgba(0, 212, 255, 0.4) 25%,
        rgba(0, 255, 136, 0.6) 50%,
        rgba(183, 148, 246, 0.4) 75%,
        transparent 100%
      );
      border-radius: 2px;
      box-shadow: 0 0 15px rgba(0, 212, 255, 0.3);
      animation: footerShimmer 4s ease-in-out infinite;
    }
  }
}

// Logo 下划线闪烁动画
@keyframes logoLineShimmer {
  0%, 100% {
    opacity: 0.6;
    width: 80px;
  }
  50% {
    opacity: 1;
    width: 120px;
  }
}

// Logo 边框光晕动画
@keyframes logoBorderGlow {
  0%, 100% {
    opacity: 0.5;
    filter: blur(2px);
  }
  50% {
    opacity: 0.8;
    filter: blur(4px);
  }
}

// 底部装饰闪烁动画
@keyframes footerShimmer {
  0%, 100% {
    opacity: 0.4;
    transform: scaleX(0.8);
  }
  50% {
    opacity: 1;
    transform: scaleX(1);
  }
}
</style>
