<template>
  <div
    class="sf-menu-item"
    :class="{ active: active, disabled: disabled }"
    @click="handleClick"
  >
    <div class="menu-item-content">
      <!-- 图标容器 -->
      <div class="menu-icon-wrapper" :class="iconClass">
        <Icon v-if="icon" :icon="icon" :size="iconSize" />
        <img v-else-if="iconImage" :src="iconImage" :alt="label" class="icon-image" />
      </div>
      
      <!-- 文字标签 -->
      <span class="menu-label">{{ label }}</span>
      
      <!-- 徽章（可选） -->
      <Badge v-if="badge" :value="badge" class="menu-badge" />
      
      <!-- 右侧箭头（可选） -->
      <Icon v-if="showArrow" :icon="ArrowRight" :size="12" class="menu-arrow" />
    </div>
    
    <!-- 激活指示器 -->
    <div v-if="active" class="active-indicator"></div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Icon, Badge } from '@/components/common/base'
import { ArrowRight } from '@element-plus/icons-vue'
import type { MenuItemProps } from './types'

const props = withDefaults(defineProps<MenuItemProps>(), {
  iconSize: 20,
  showArrow: false,
  disabled: false,
})

const emit = defineEmits<{
  click: [event: MouseEvent]
}>()

const iconClass = computed(() => {
  if (props.iconColor) return ''
  return `icon-${props.iconType || 'default'}`
})

const handleClick = (event: MouseEvent) => {
  if (props.disabled) return
  emit('click', event)
}
</script>

<style scoped lang="scss">
.sf-menu-item {
  position: relative;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  margin: 4px 12px;
  border-radius: 12px;
  overflow: hidden;

  .menu-item-content {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 14px 16px;
    position: relative;
    z-index: 1;
  }

  .menu-icon-wrapper {
    width: 44px;
    height: 44px;
    border-radius: 12px;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
    transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
    position: relative;
    overflow: visible;

    // 背景光晕层
    &::before {
      content: '';
      position: absolute;
      inset: -2px;
      border-radius: 12px;
      background: currentColor;
      opacity: 0;
      filter: blur(8px);
      transition: all 0.4s ease;
      z-index: -1;
    }

    // 内部背景
    &::after {
      content: '';
      position: absolute;
      inset: 0;
      border-radius: 12px;
      background: currentColor;
      opacity: 0.15;
      transition: opacity 0.3s ease;
      z-index: 0;
    }

    // 图标本身
    :deep(svg),
    .icon-image {
      position: relative;
      z-index: 1;
      filter: drop-shadow(0 0 4px currentColor);
      transition: all 0.3s ease;
    }

    .icon-image {
      width: 100%;
      height: 100%;
      object-fit: contain;
    }
  }

  .menu-label {
    flex: 1;
    font-size: 14px;
    font-weight: 500;
    color: rgba(255, 255, 255, 0.85);
    transition: all 0.3s ease;
    white-space: nowrap;
  }

  .menu-badge {
    flex-shrink: 0;
  }

  .menu-arrow {
    flex-shrink: 0;
    color: rgba(255, 255, 255, 0.5);
    transition: all 0.3s ease;
  }

  // 不同图标类型的颜色和渐变
  .icon-default {
    color: rgba(255, 255, 255, 0.7);
    
    &::after {
      background: rgba(255, 255, 255, 0.1);
    }
  }

  .icon-dashboard {
    color: #00d4ff; // 青蓝色
    background: linear-gradient(135deg, rgba(0, 212, 255, 0.25) 0%, rgba(64, 158, 255, 0.25) 100%);
    box-shadow: 
      0 0 20px rgba(0, 212, 255, 0.3),
      inset 0 0 15px rgba(0, 212, 255, 0.1);
    
    &::after {
      background: linear-gradient(135deg, rgba(0, 212, 255, 0.2) 0%, rgba(64, 158, 255, 0.2) 100%);
    }
  }

  .icon-workflow {
    color: #00ff88; // 绿色
    background: linear-gradient(135deg, rgba(0, 255, 136, 0.25) 0%, rgba(103, 194, 58, 0.25) 100%);
    box-shadow: 
      0 0 20px rgba(0, 255, 136, 0.3),
      inset 0 0 15px rgba(0, 255, 136, 0.1);
    
    &::after {
      background: linear-gradient(135deg, rgba(0, 255, 136, 0.2) 0%, rgba(103, 194, 58, 0.2) 100%);
    }
  }

  .icon-user {
    color: #b794f6; // 紫红色
    background: linear-gradient(135deg, rgba(183, 148, 246, 0.25) 0%, rgba(219, 39, 119, 0.25) 100%);
    box-shadow: 
      0 0 20px rgba(183, 148, 246, 0.3),
      inset 0 0 15px rgba(183, 148, 246, 0.1);
    
    &::after {
      background: linear-gradient(135deg, rgba(183, 148, 246, 0.2) 0%, rgba(219, 39, 119, 0.2) 100%);
    }
  }

  .icon-ai {
    color: #ffb84d; // 橙金色
    background: linear-gradient(135deg, rgba(255, 184, 77, 0.25) 0%, rgba(255, 107, 107, 0.25) 100%);
    box-shadow: 
      0 0 20px rgba(255, 184, 77, 0.3),
      inset 0 0 15px rgba(255, 184, 77, 0.1);
    
    &::after {
      background: linear-gradient(135deg, rgba(255, 184, 77, 0.2) 0%, rgba(255, 107, 107, 0.2) 100%);
    }
  }

  .icon-system {
    color: #6c8eff; // 蓝紫色
    background: linear-gradient(135deg, rgba(108, 142, 255, 0.25) 0%, rgba(147, 51, 234, 0.25) 100%);
    box-shadow: 
      0 0 20px rgba(108, 142, 255, 0.3),
      inset 0 0 15px rgba(108, 142, 255, 0.1);
    
    &::after {
      background: linear-gradient(135deg, rgba(108, 142, 255, 0.2) 0%, rgba(147, 51, 234, 0.2) 100%);
    }
  }

  // 悬停效果
  &:hover:not(.disabled) {
    background: rgba(64, 158, 255, 0.08);
    transform: translateX(6px);
    box-shadow: 
      inset 0 0 30px rgba(64, 158, 255, 0.1),
      0 4px 20px rgba(0, 0, 0, 0.2);

    .menu-icon-wrapper {
      transform: scale(1.15) translateY(-2px);
      
      &::before {
        opacity: 0.4;
        filter: blur(12px);
      }
      
      &::after {
        opacity: 0.25;
      }
      
      :deep(svg),
      .icon-image {
        filter: drop-shadow(0 0 8px currentColor);
        transform: scale(1.1);
      }
    }

    .menu-label {
      color: rgba(255, 255, 255, 1);
      text-shadow: 0 0 10px rgba(255, 255, 255, 0.3);
    }

    .menu-arrow {
      color: rgba(255, 255, 255, 0.9);
      transform: translateX(6px);
    }
  }

  // 激活状态 - 柔和高亮蓝色光圈
  &.active {
    background: linear-gradient(90deg, rgba(0, 212, 255, 0.18) 0%, rgba(64, 158, 255, 0.08) 100%);
    border-left: 3px solid #00d4ff;
    box-shadow: 
      inset 0 0 30px rgba(0, 212, 255, 0.15),
      0 0 30px rgba(0, 212, 255, 0.3),
      0 4px 20px rgba(0, 0, 0, 0.3);
    transform: translateX(4px);

    .menu-icon-wrapper {
      background: linear-gradient(135deg, rgba(0, 212, 255, 0.35) 0%, rgba(64, 158, 255, 0.35) 100%);
      color: #00d4ff;
      box-shadow: 
        0 0 25px rgba(0, 212, 255, 0.5),
        0 4px 20px rgba(0, 212, 255, 0.3),
        inset 0 0 20px rgba(0, 212, 255, 0.2);
      animation: iconPulse 2s ease-in-out infinite;

      &::before {
        opacity: 0.5;
        filter: blur(15px);
      }
      
      &::after {
        opacity: 0.3;
      }
      
      :deep(svg),
      .icon-image {
        filter: drop-shadow(0 0 12px currentColor);
        animation: iconGlow 2s ease-in-out infinite;
      }
    }

    .menu-label {
      color: #ffffff;
      font-weight: 600;
      text-shadow: 
        0 0 10px rgba(0, 212, 255, 0.5),
        0 0 20px rgba(0, 212, 255, 0.3);
    }

    .active-indicator {
      position: absolute;
      right: 0;
      top: 50%;
      transform: translateY(-50%);
      width: 4px;
      height: 70%;
      background: linear-gradient(180deg, #00d4ff 0%, #00ff88 50%, #b794f6 100%);
      border-radius: 2px 0 0 2px;
      box-shadow: 
        0 0 15px rgba(0, 212, 255, 0.8),
        0 0 30px rgba(0, 212, 255, 0.4);
      animation: indicatorPulse 2s ease-in-out infinite;
    }
  }

  // 禁用状态
  &.disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
}

// 图标脉冲动画
@keyframes iconPulse {
  0%, 100% {
    box-shadow: 
      0 0 25px rgba(0, 212, 255, 0.5),
      0 4px 20px rgba(0, 212, 255, 0.3),
      inset 0 0 20px rgba(0, 212, 255, 0.2);
  }
  50% {
    box-shadow: 
      0 0 35px rgba(0, 212, 255, 0.7),
      0 4px 25px rgba(0, 212, 255, 0.5),
      inset 0 0 25px rgba(0, 212, 255, 0.3);
  }
}

// 图标光晕动画
@keyframes iconGlow {
  0%, 100% {
    filter: drop-shadow(0 0 12px currentColor);
  }
  50% {
    filter: drop-shadow(0 0 18px currentColor) drop-shadow(0 0 8px currentColor);
  }
}

// 指示器脉冲动画
@keyframes indicatorPulse {
  0%, 100% {
    opacity: 1;
    box-shadow: 
      0 0 15px rgba(0, 212, 255, 0.8),
      0 0 30px rgba(0, 212, 255, 0.4);
  }
  50% {
    opacity: 0.8;
    box-shadow: 
      0 0 20px rgba(0, 212, 255, 1),
      0 0 40px rgba(0, 212, 255, 0.6);
  }
}
</style>

