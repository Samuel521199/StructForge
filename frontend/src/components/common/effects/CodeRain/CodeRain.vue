<template>
  <canvas ref="canvasRef" class="code-rain-canvas" />
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, reactive } from 'vue'
import type { CodeRainProps } from './types'
import { loadConfig, mergeConfig, defaultConfig, type CodeRainConfig } from './config'

// 日文片假名字符集（矩阵风格）
const KATAKANA_CHARS = 'アァカサタナハマヤラワガザダバパ0123456789!@#$%^&*()_+-=[]{};:"|,./<>?'

const props = withDefaults(defineProps<CodeRainProps>(), {
  fontSize: 15,
  fontFamily: 'Monospace',
  fontWeight: 'normal',
  color: '#00FF00',
  backgroundColor: '#000000',
  speed: 2.5,
  speedVariation: 0.6,
  density: 0.008,
  opacity: 0.9,
  fadeSpeed: 0.04,
  characters: KATAKANA_CHARS,
  minLength: 15,
  maxLength: 35,
  enableLayers: false,
  enableGlow: true,
  enableGlitch: false,
  glowIntensity: 0.8,
  configPath: undefined,
  useConfigFile: false,
})

// 合并后的配置（优先级：props > 配置文件 > 默认值）
const config = reactive<CodeRainConfig>({ ...defaultConfig })

const canvasRef = ref<HTMLCanvasElement | null>(null)
let animationFrameId: number = 0
let canvas: HTMLCanvasElement | null = null
let ctx: CanvasRenderingContext2D | null = null

// 动态变量
let canvasWidth: number = 0
let canvasHeight: number = 0
let columns: number = 0
let drops: number[] = [] // 每列代码雨的 Y 坐标
let columnXPositions: number[] = [] // 每列的 X 坐标

// 计算属性：从配置中获取常量
const COLUMN_SPACING_FACTOR = () => config.columnSpacingFactor ?? 2.5
const TRAIL_LENGTH = () => config.trailLength ?? 20

/**
 * 从字符集中随机获取一个字符
 */
const getRandomCharacter = (): string => {
  const chars = config.characters ?? props.characters ?? KATAKANA_CHARS
  return chars.charAt(Math.floor(Math.random() * chars.length))
}

/**
 * 重新计算 Canvas 尺寸和代码雨列数
 */
const resizeCanvas = () => {
  if (!canvas) return

  // 设置 Canvas 实际像素尺寸（解决模糊问题）
  const rect = canvas.getBoundingClientRect()
  canvasWidth = canvas.width = rect.width
  canvasHeight = canvas.height = rect.height

  // 计算列数
  const fontSize = config.fontSize ?? props.fontSize ?? 15
  const spacing = fontSize * COLUMN_SPACING_FACTOR()
  columns = Math.floor(canvasWidth / spacing)

  // 初始化/重置 drops 数组
  if (drops.length !== columns) {
    drops = new Array(columns).fill(0).map(() =>
      // 随机初始Y坐标，避免所有代码串同步
      -Math.random() * canvasHeight * 0.8
    )
    columnXPositions = new Array(columns).fill(0).map((_, i) => i * spacing)
  }
}

/**
 * 核心渲染函数：每一帧的绘制逻辑
 */
const draw = () => {
  if (!ctx || !canvas) return

  // 1. 残影清屏（电影级拖尾的关键）
  // 极低透明度的黑色，创建字符的残影和拖尾效果
  ctx.fillStyle = 'rgba(0, 0, 0, 0.04)'
  ctx.fillRect(0, 0, canvasWidth, canvasHeight)

  // 禁用滤镜，以防上次渲染中启用了模糊
  ctx.filter = 'none'

  // 从配置中获取参数
  const fontSize = config.fontSize ?? props.fontSize ?? 15
  const fontWeight = config.fontWeight ?? props.fontWeight ?? 'normal'
  const fontFamily = config.fontFamily ?? props.fontFamily ?? 'Monospace'
  const opacity = config.opacity ?? props.opacity ?? 0.9
  const enableGlow = config.enableGlow ?? props.enableGlow ?? true
  const glowIntensity = config.glowIntensity ?? props.glowIntensity ?? 0.8
  const speed = config.speed ?? props.speed ?? 2.5

  // 设置基础字体样式
  ctx.font = `${fontWeight} ${fontSize}px ${fontFamily}`

  // 遍历每一列的代码串
  for (let i = 0; i < drops.length; i++) {
    const x_pos = columnXPositions[i]
    let y = drops[i]

    // --- 2. 渲染尾迹（渐变 Alpha/Glow） ---
    const trailLength = TRAIL_LENGTH()
    for (let j = 1; j <= trailLength; j++) {
      const charY = y - j * fontSize // 计算当前尾迹字符的 Y 位置

      // 越远离头部 (j越大)，透明度越低，颜色越深
      const alpha = opacity * 0.7 * (1 - j / trailLength)

      // 尾迹颜色：柔和的绿色
      ctx.fillStyle = `rgba(0, 255, 0, ${alpha})`

      if (enableGlow) {
        ctx.shadowColor = '#00FF00'
        ctx.shadowBlur = 4 * glowIntensity // 较弱的光晕
      } else {
        ctx.shadowBlur = 0
        ctx.shadowColor = 'transparent'
      }

      // 随机获取字符
      const text = getRandomCharacter()
      ctx.fillText(text, x_pos, charY)
    }

    // --- 3. 渲染领头字符（最强发光头） ---
    // 领头字符的颜色和光晕设置
    ctx.fillStyle = '#FFFFFF' // 纯白高光

    if (enableGlow) {
      ctx.shadowColor = '#FFFFFF' // 白色光晕
      ctx.shadowBlur = 20 * glowIntensity // 强大的发光效果（电影感核心）
    } else {
      ctx.shadowBlur = 0
      ctx.shadowColor = 'transparent'
    }

    // 随机获取字符
    const leaderText = getRandomCharacter()
    ctx.fillText(leaderText, x_pos, y)

    // 重置阴影
    ctx.shadowBlur = 0
    ctx.shadowColor = 'transparent'

    // --- 4. 更新 Y 坐标和重置 ---
    // 每列的速度可以随机，这里使用基础速度
    const dropSpeed = speed * fontSize * 0.8
    drops[i] += dropSpeed

    // 如果代码串落到底部，随机重置到屏幕外顶部
    if (y > canvasHeight && Math.random() > 0.95) {
      drops[i] = -Math.random() * canvasHeight
    }
  }

  // 循环调用
  animationFrameId = requestAnimationFrame(draw)
}

// 初始化画布
const initCanvas = async () => {
  if (!canvasRef.value) return

  // 如果启用配置文件，先加载配置
  if (props.useConfigFile) {
    try {
      const fileConfig = await loadConfig(props.configPath)
      
      // 当使用配置文件时，过滤掉 props 中的默认值，让配置文件的值生效
      // 只有当 props 中显式传递了不同的值时，才使用 props 的值
      const filteredProps: Partial<CodeRainConfig> = {}
      
      // 对于 characters：如果等于默认值，则不使用（让配置文件的值生效）
      if (props.characters && props.characters !== KATAKANA_CHARS) {
        filteredProps.characters = props.characters
      }
      
      // 对于其他属性：如果值不等于默认值，则认为是显式传递的
      if (props.fontSize !== 15) filteredProps.fontSize = props.fontSize
      if (props.speed !== 2.5) filteredProps.speed = props.speed
      if (props.speedVariation !== 0.6) filteredProps.speedVariation = props.speedVariation
      if (props.density !== 0.008) filteredProps.density = props.density
      if (props.opacity !== 0.9) filteredProps.opacity = props.opacity
      if (props.fadeSpeed !== 0.04) filteredProps.fadeSpeed = props.fadeSpeed
      if (props.glowIntensity !== 0.8) filteredProps.glowIntensity = props.glowIntensity
      if (props.fontFamily !== 'Monospace') filteredProps.fontFamily = props.fontFamily
      if (props.fontWeight !== 'normal') filteredProps.fontWeight = props.fontWeight
      if (props.color !== '#00FF00') filteredProps.color = props.color
      if (props.backgroundColor !== '#000000') filteredProps.backgroundColor = props.backgroundColor
      if (props.minLength !== 15) filteredProps.minLength = props.minLength
      if (props.maxLength !== 35) filteredProps.maxLength = props.maxLength
      if (props.enableLayers !== false) filteredProps.enableLayers = props.enableLayers
      if (props.enableGlow !== true) filteredProps.enableGlow = props.enableGlow
      if (props.enableGlitch !== false) filteredProps.enableGlitch = props.enableGlitch
      
      // 合并配置（优先级：显式传递的 props > 配置文件 > 默认值）
      const mergedConfig = mergeConfig(filteredProps, fileConfig, defaultConfig)
      Object.assign(config, mergedConfig)
      
      // 调试：输出加载的配置（仅开发环境）
      if (import.meta.env.DEV) {
        console.log('[CodeRain] 配置加载完成:', {
          characters: config.characters?.substring(0, 30) + '...',
          charactersLength: config.characters?.length,
          fromFile: fileConfig.characters?.substring(0, 30) + '...',
          fromProps: props.characters?.substring(0, 30) + '...',
          isDefaultChars: props.characters === KATAKANA_CHARS,
          speed: config.speed,
          fontSize: config.fontSize
        })
      }
    } catch (error) {
      console.warn('[CodeRain] 配置文件加载失败，使用 props 和默认值:', error)
      // 如果配置文件加载失败，使用 props 和默认值
      Object.assign(config, mergeConfig(props, defaultConfig, defaultConfig))
    }
  } else {
    // 不使用配置文件，直接使用 props 和默认值
    Object.assign(config, mergeConfig(props, defaultConfig, defaultConfig))
  }

  canvas = canvasRef.value
  ctx = canvas.getContext('2d')
  if (!ctx) return

  resizeCanvas()
  window.addEventListener('resize', resizeCanvas)

  // 开始渲染循环
  draw()
}

onMounted(() => {
  initCanvas()
})

onUnmounted(() => {
  // 清理资源，避免内存泄漏
  window.removeEventListener('resize', resizeCanvas)
  if (animationFrameId) {
    cancelAnimationFrame(animationFrameId)
  }
})
</script>

<style scoped lang="scss">
.code-rain-canvas {
  /* 关键：确保 Canvas 覆盖整个视口，并在内容下方 */
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 1; /* 比 z-index 为 10 的登录框低 */
  background-color: #000; /* 确保背景是纯黑，提高对比度 */
  pointer-events: none;
  display: block;
}
</style>
