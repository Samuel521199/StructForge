<template>
  <div class="sf-dropdown">
    <div class="dropdown-trigger" @click="toggleDropdown">
      <slot name="trigger" />
    </div>
    <Transition name="dropdown">
      <div v-if="isOpen" class="dropdown-menu" @click.stop>
        <slot name="dropdown" />
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref, onUnmounted, watch } from 'vue'

const props = withDefaults(defineProps<{
  trigger?: 'click' | 'hover'
}>(), {
  trigger: 'click',
})

const isOpen = ref(false)

const toggleDropdown = () => {
  if (props.trigger === 'click') {
    isOpen.value = !isOpen.value
  }
}

// 监听点击事件，点击外部关闭
const handleDocumentClick = (event: MouseEvent) => {
  const target = event.target as Node
  const dropdownEl = document.querySelector('.sf-dropdown')
  if (dropdownEl && !dropdownEl.contains(target)) {
    isOpen.value = false
  }
}

// 当打开时添加监听器
watch(isOpen, (newVal) => {
  if (newVal) {
    setTimeout(() => {
      document.addEventListener('click', handleDocumentClick)
    }, 0)
  } else {
    document.removeEventListener('click', handleDocumentClick)
  }
})

onUnmounted(() => {
  isOpen.value = false
  document.removeEventListener('click', handleDocumentClick)
})
</script>

<style scoped lang="scss">
.sf-dropdown {
  position: relative;
  display: inline-block;

  .dropdown-trigger {
    cursor: pointer;
  }

  .dropdown-menu {
    position: absolute;
    top: 100%;
    right: 0;
    margin-top: 4px;
    background: #fff;
    border: 1px solid var(--el-border-color-light);
    border-radius: 4px;
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
    z-index: 1000;
    min-width: 120px;
  }
}

.dropdown-enter-active,
.dropdown-leave-active {
  transition: opacity 0.2s, transform 0.2s;
}

.dropdown-enter-from,
.dropdown-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}
</style>

