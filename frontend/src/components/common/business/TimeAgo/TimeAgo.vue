<template>
  <span class="sf-time-ago" :title="fullTime">
    {{ displayText }}
  </span>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import type { TimeAgoProps } from './types'

const props = withDefaults(defineProps<TimeAgoProps>(), {
  format: 'YYYY-MM-DD HH:mm:ss',
  updateInterval: 60000, // 默认1分钟更新一次
})

const now = ref(Date.now())

// 格式化时间
const formatTime = (timestamp: number | string | Date): string => {
  const date = typeof timestamp === 'number' 
    ? new Date(timestamp) 
    : typeof timestamp === 'string' 
    ? new Date(timestamp) 
    : timestamp

  if (isNaN(date.getTime())) {
    return '无效时间'
  }

  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  const seconds = String(date.getSeconds()).padStart(2, '0')

  return props.format
    .replace('YYYY', String(year))
    .replace('MM', month)
    .replace('DD', day)
    .replace('HH', hours)
    .replace('mm', minutes)
    .replace('ss', seconds)
}

// 完整时间显示
const fullTime = computed(() => {
  return formatTime(props.time)
})

// 相对时间文本
const relativeText = computed(() => {
  const time = typeof props.time === 'number' 
    ? props.time 
    : typeof props.time === 'string' 
    ? new Date(props.time).getTime() 
    : props.time.getTime()

  if (isNaN(time)) {
    return '无效时间'
  }

  const diff = now.value - time
  const seconds = Math.floor(diff / 1000)
  const minutes = Math.floor(seconds / 60)
  const hours = Math.floor(minutes / 60)
  const days = Math.floor(hours / 24)
  const months = Math.floor(days / 30)
  const years = Math.floor(days / 365)

  if (seconds < 0) {
    return '未来'
  }

  if (seconds < 60) {
    return '刚刚'
  }

  if (minutes < 60) {
    return `${minutes}分钟前`
  }

  if (hours < 24) {
    return `${hours}小时前`
  }

  if (days < 30) {
    return `${days}天前`
  }

  if (months < 12) {
    return `${months}个月前`
  }

  return `${years}年前`
})

// 显示文本
const displayText = computed(() => {
  if (props.showFullTime) {
    return fullTime.value
  }
  return relativeText.value
})

// 定时更新
let timer: number | null = null

onMounted(() => {
  if (props.updateInterval > 0) {
    timer = window.setInterval(() => {
      now.value = Date.now()
    }, props.updateInterval)
  }
})

onUnmounted(() => {
  if (timer) {
    clearInterval(timer)
  }
})
</script>

<style scoped lang="scss">
.sf-time-ago {
  color: var(--el-text-color-regular);
  font-size: 14px;
  cursor: help;
}
</style>

