<template>
  <div ref="chartRef" class="sf-chart" :style="{ width: width, height: height }"></div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue'
import * as echarts from 'echarts'
import type { EChartsOption } from 'echarts'
import type { ChartProps } from './types'

const props = withDefaults(defineProps<ChartProps>(), {
  width: '100%',
  height: '400px',
  theme: 'dark',
})

const chartRef = ref<HTMLDivElement>()
let chartInstance: echarts.ECharts | null = null

const initChart = () => {
  if (!chartRef.value) return

  // 销毁旧实例
  if (chartInstance) {
    chartInstance.dispose()
  }

  // 创建新实例
  chartInstance = echarts.init(chartRef.value, props.theme)

  // 设置配置
  if (props.option) {
    chartInstance.setOption(props.option, true)
  }

  // 监听窗口大小变化
  const resizeHandler = () => {
    chartInstance?.resize()
  }
  window.addEventListener('resize', resizeHandler)

  // 清理函数
  onUnmounted(() => {
    window.removeEventListener('resize', resizeHandler)
    if (chartInstance) {
      chartInstance.dispose()
      chartInstance = null
    }
  })
}

// 监听配置变化
watch(
  () => props.option,
  (newOption) => {
    if (chartInstance && newOption) {
      chartInstance.setOption(newOption, true)
    }
  },
  { deep: true }
)

onMounted(() => {
  initChart()
})

// 暴露方法
defineExpose({
  getInstance: () => chartInstance,
  resize: () => chartInstance?.resize(),
})
</script>

<style scoped lang="scss">
.sf-chart {
  width: 100%;
  height: 100%;
}
</style>

