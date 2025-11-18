<template>
  <div class="sf-avatar-group" :class="[`is-${size}`, { 'is-stacked': stacked }]">
    <div
      v-for="(avatar, index) in displayAvatars"
      :key="index"
      class="sf-avatar-group__item"
      :style="{ zIndex: displayAvatars.length - index }"
      @click="handleAvatarClick(avatar, index)"
    >
      <el-avatar
        :size="avatarSize"
        :src="avatar.src"
        :icon="avatar.icon"
        :shape="avatar.shape || shape"
      >
        {{ avatar.text }}
      </el-avatar>
      <span v-if="avatar.label" class="sf-avatar-group__label">
        {{ avatar.label }}
      </span>
    </div>
    <div
      v-if="max && avatars.length > max"
      class="sf-avatar-group__more"
      :style="{ zIndex: 0 }"
      @click="handleMoreClick"
    >
      <el-avatar :size="avatarSize" :shape="shape">
        +{{ avatars.length - max }}
      </el-avatar>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { AvatarGroupProps, AvatarItem } from './types'

const props = withDefaults(defineProps<AvatarGroupProps>(), {
  size: 'default',
  shape: 'circle',
  stacked: false,
  max: 0,
})

const emit = defineEmits<{
  avatarClick: [avatar: AvatarItem, index: number]
  moreClick: []
}>()

// 头像大小映射
const sizeMap = {
  large: 40,
  default: 32,
  small: 24,
}

const avatarSize = computed(() => {
  return typeof props.size === 'number' ? props.size : sizeMap[props.size]
})

// 显示的头像列表
const displayAvatars = computed(() => {
  if (props.max && props.max > 0) {
    return props.avatars.slice(0, props.max)
  }
  return props.avatars
})

const handleAvatarClick = (avatar: AvatarItem, index: number) => {
  emit('avatarClick', avatar, index)
}

const handleMoreClick = () => {
  emit('moreClick')
}
</script>

<style scoped lang="scss">
.sf-avatar-group {
  display: flex;
  align-items: center;
  gap: 8px;

  &.is-stacked {
    .sf-avatar-group__item,
    .sf-avatar-group__more {
      margin-left: -8px;

      &:first-child {
        margin-left: 0;
      }
    }
  }

  &__item {
    position: relative;
    cursor: pointer;
    transition: transform 0.2s;

    &:hover {
      transform: translateY(-2px);
    }
  }

  &__label {
    position: absolute;
    bottom: -20px;
    left: 50%;
    transform: translateX(-50%);
    font-size: 12px;
    color: var(--el-text-color-regular);
    white-space: nowrap;
  }

  &__more {
    cursor: pointer;
    transition: transform 0.2s;

    &:hover {
      transform: translateY(-2px);
    }
  }
}
</style>

