<template>
  <canvas ref="canvasRef" class="code-rain-canvas" />
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import type { CodeRainProps } from './types'

// 日文片假名字符集（矩阵风格）
const KATAKANA_CHARS = 'アイウエオカキクケコサシスセソタチツテトナニヌネノハヒフヘホマミムメモヤユヨラリルレロワヲン0123456789$%&*+-=[]{}|;:,.<>?/~`'

const props = withDefaults(defineProps<CodeRainProps>(), {
  fontSize: 14,
  fontFamily: 'monospace',
  fontWeight: 'bold',
  color: '#00ff41',
  backgroundColor: '#000000',
  speed: 2.5,
  speedVariation: 0.6,
  density: 0.008,
  opacity: 0.9,
  fadeSpeed: 0.04,
  characters: KATAKANA_CHARS,
  minLength: 20,
  maxLength: 40,
  enableLayers: true,
  enableGlow: true,
  enableGlitch: true,
  glowIntensity: 0.8,
})

const canvasRef = ref<HTMLCanvasElement | null>(null)
let animationFrameId: number | null = null
let canvas: HTMLCanvasElement | null = null
let ctx: CanvasRenderingContext2D | null = null

// 颜色脉冲状态
let colorPulsePhase = 0
let colorPulseSpeed = 0.002

// 屏幕抖动状态
let glitchShakeX = 0
let glitchShakeY = 0
let glitchActive = false
let glitchTimer = 0

// 代码雨滴类（支持分层）
class CodeDrop {
  x: number
  y: number
  length: number
  speed: number
  opacity: number
  characters: string[]
  layer: 'far' | 'mid' | 'near' // 所属层级
  baseColor: string // 基础颜色

  constructor(x: number, canvasHeight: number, layer: 'far' | 'mid' | 'near' = 'mid') {
    this.x = x
    // 让起始位置更分散：有些从屏幕上方很远开始，有些从屏幕中间开始
    // 使用加权随机，让更多雨滴从不同位置开始
    const startPosition = Math.random()
    if (startPosition < 0.3) {
      // 30% 从屏幕上方很远开始
      this.y = Math.random() * -canvasHeight * 3 - canvasHeight
    } else if (startPosition < 0.6) {
      // 30% 从屏幕上方开始
      this.y = Math.random() * -canvasHeight * 2
    } else {
      // 40% 从屏幕中间或上方开始（更自然）
      this.y = Math.random() * -canvasHeight * 1.5
    }
    
    this.layer = layer
    // 增加长度随机性：有些很短，有些很长，使用加权随机
    const lengthRandom = Math.random()
    if (lengthRandom < 0.2) {
      // 20% 很短
      this.length = Math.floor(
        Math.random() * (props.minLength * 0.5) + props.minLength * 0.3
      )
    } else if (lengthRandom < 0.5) {
      // 30% 中等偏短
      this.length = Math.floor(
        Math.random() * (props.minLength * 0.8) + props.minLength * 0.5
      )
    } else if (lengthRandom < 0.8) {
      // 30% 正常长度
      this.length = Math.floor(
        Math.random() * (props.maxLength - props.minLength) + props.minLength
      )
    } else {
      // 20% 很长
      this.length = Math.floor(
        Math.random() * (props.maxLength * 1.5 - props.maxLength) + props.maxLength
      )
    }
    // 确保长度至少为1
    this.length = Math.max(1, this.length)
    
    // 根据层级调整速度
    const baseSpeed = props.speed
    const variation = props.speedVariation
    let speedMultiplier = 1
    
    if (props.enableLayers) {
      switch (layer) {
        case 'far':
          speedMultiplier = 0.4 // 远景慢
          break
        case 'mid':
          speedMultiplier = 1.0 // 中景正常
          break
        case 'near':
          speedMultiplier = 1.3 // 近景快
          break
      }
    }
    
    const minSpeed = baseSpeed * speedMultiplier * (1 - variation)
    const maxSpeed = baseSpeed * speedMultiplier * (1 + variation)
    this.speed = Math.random() * (maxSpeed - minSpeed) + minSpeed
    
    // 根据层级调整透明度
    if (props.enableLayers) {
      switch (layer) {
        case 'far':
          this.opacity = Math.random() * 0.2 + 0.1 // 远景很淡
          break
        case 'mid':
          this.opacity = Math.random() * 0.3 + 0.5 // 中景正常
          break
        case 'near':
          this.opacity = Math.random() * 0.2 + 0.8 // 近景很亮
          break
      }
    } else {
      this.opacity = Math.random() * 0.5 + 0.5
    }
    
    this.baseColor = props.color
    this.characters = this.generateCharacters()
  }

  generateCharacters(): string[] {
    const chars: string[] = []
    const charSet = props.characters
    for (let i = 0; i < this.length; i++) {
      chars.push(
        charSet[Math.floor(Math.random() * charSet.length)]
      )
    }
    return chars
  }
  
  // 更新字符，模拟数据流的混乱感
  updateCharacters() {
    // 每帧有 5-10% 的概率随机替换字符，模仿数据流的混乱感
    if (Math.random() < 0.08) {
      const changeCount = Math.floor(this.characters.length * 0.1)
      for (let i = 0; i < changeCount; i++) {
        const randomIndex = Math.floor(Math.random() * this.characters.length)
        this.characters[randomIndex] = props.characters[
          Math.floor(Math.random() * props.characters.length)
        ]
      }
    }
  }

  update(canvasHeight: number) {
    // 参考博客实现：简洁的下落逻辑
    this.y += this.speed
    
    // 字符随机更新
    this.updateCharacters()
    
    // 颜色脉冲效果（仅中景和近景）
    if (props.enableGlitch && (this.layer === 'mid' || this.layer === 'near')) {
      const pulse = Math.sin(colorPulsePhase + this.x * 0.01) * 0.1
      // 轻微调整颜色亮度
      this.baseColor = adjustColorBrightness(props.color, pulse)
    }
    
    // 参考博客实现：当雨滴到达底部或随机重置时，重新开始
    // 参考博客：value >= height || value > 8888 * Math.random() ? 0 : value + 10
    if (this.y >= canvasHeight || this.y > 8888 * Math.random()) {
      // 随机重置位置到屏幕上方，避免同步
      this.y = -Math.random() * canvasHeight * 2
      
      // 重新随机长度（增加变化）
      const lengthRandom = Math.random()
      if (lengthRandom < 0.2) {
        this.length = Math.floor(
          Math.random() * (props.minLength * 0.5) + props.minLength * 0.3
        )
      } else if (lengthRandom < 0.5) {
        this.length = Math.floor(
          Math.random() * (props.minLength * 0.8) + props.minLength * 0.5
        )
      } else if (lengthRandom < 0.8) {
        this.length = Math.floor(
          Math.random() * (props.maxLength - props.minLength) + props.minLength
        )
      } else {
        this.length = Math.floor(
          Math.random() * (props.maxLength * 1.5 - props.maxLength) + props.maxLength
        )
      }
      this.length = Math.max(1, this.length)
      
      this.characters = this.generateCharacters()
      
      // 重新随机透明度
      if (props.enableLayers) {
        switch (this.layer) {
          case 'far':
            this.opacity = Math.random() * 0.2 + 0.1
            break
          case 'mid':
            this.opacity = Math.random() * 0.3 + 0.5
            break
          case 'near':
            this.opacity = Math.random() * 0.2 + 0.8
            break
        }
      } else {
        this.opacity = Math.random() * 0.5 + 0.5
      }
    }
  }

  draw(ctx: CanvasRenderingContext2D, fontSize: number, baseColor: string) {
    const effectiveColor = this.baseColor || baseColor
    
    // 根据层级调整字体大小
    let effectiveFontSize = fontSize
    if (props.enableLayers) {
      switch (this.layer) {
        case 'far':
          effectiveFontSize = fontSize * 0.6 // 远景小
          break
        case 'mid':
          effectiveFontSize = fontSize // 中景正常
          break
        case 'near':
          effectiveFontSize = fontSize * 1.2 // 近景大
          break
      }
    }
    
    ctx.font = `${props.fontWeight} ${effectiveFontSize}px ${props.fontFamily}`
    
    // 绘制每个字符，实现光晕和渐变尾迹
    this.characters.forEach((char, index) => {
      const position = index / (this.characters.length - 1 || 1)
      
      // 渐变尾迹：快速 Alpha 衰减，营造数据流动的"模糊美"
      // 使用更强的衰减曲线，让尾迹快速变透明
      const fadeFactor = Math.pow(1 - position, 2.5)
      // Alpha 从接近 0.8 快速衰减到接近 0，确保字符清晰可见
      const charOpacity = this.opacity * (0.1 + fadeFactor * 0.7)
      
      // 领头字符（第一个）使用白色/极亮绿色，并添加强烈光晕
      let charColor = effectiveColor
      let shadowBlur = 0
      let shadowColor = 'transparent'
      
      if (index === 0) {
        // 领头字符：纯白色，全场最亮的点，像数据包的头部
        charColor = '#FFFFFF'
        if (props.enableGlow) {
          // 强烈光晕，至少 25，营造霓虹灯效果
          shadowBlur = Math.max(25, 25 * props.glowIntensity)
          shadowColor = '#FFFFFF'
        }
        // 领头字符使用最高透明度
        ctx.globalAlpha = this.opacity
      } else {
        // 后续字符：高饱和度、高亮度的绿色，确保清晰可见
        // 使用 rgba 格式，alpha 随位置衰减
        const greenIntensity = Math.floor(255 * fadeFactor)
        const alpha = charOpacity
        charColor = `rgba(0, ${greenIntensity}, 0, ${alpha})`
        
        if (props.enableGlow) {
          // 尾迹光晕：5-8 的 shadowBlur，柔和的绿色光晕
          shadowBlur = Math.max(5, Math.min(8, 8 * props.glowIntensity * fadeFactor))
          shadowColor = '#00FF00'
        }
        ctx.globalAlpha = charOpacity
      }
      
      // 应用光晕效果
      if (props.enableGlow && shadowBlur > 0) {
        ctx.shadowBlur = shadowBlur
        ctx.shadowColor = shadowColor
      } else {
        ctx.shadowBlur = 0
        ctx.shadowColor = 'transparent'
      }
      
      ctx.fillStyle = charColor
      
      // 增大字符间距（Y轴）：至少是字体大小的 1.2 倍，确保字符独立可见
      const charSpacing = effectiveFontSize * 1.2
      ctx.fillText(
        char,
        this.x,
        this.y + index * charSpacing
      )
    })
    
    // 重置全局状态
    ctx.globalAlpha = 1
    ctx.shadowBlur = 0
    ctx.shadowColor = 'transparent'
  }
}

// 调整颜色亮度的辅助函数
function adjustColorBrightness(hex: string, amount: number): string {
  const num = parseInt(hex.replace('#', ''), 16)
  const r = Math.max(0, Math.min(255, ((num >> 16) & 0xff) + amount * 255))
  const g = Math.max(0, Math.min(255, ((num >> 8) & 0xff) + amount * 255))
  const b = Math.max(0, Math.min(255, (num & 0xff) + amount * 255))
  return `#${Math.round(r).toString(16).padStart(2, '0')}${Math.round(g).toString(16).padStart(2, '0')}${Math.round(b).toString(16).padStart(2, '0')}`
}

// 初始化画布
const initCanvas = () => {
  if (!canvasRef.value) return

  canvas = canvasRef.value
  ctx = canvas.getContext('2d')
  if (!ctx) return

  // 创建代码雨滴的函数（支持分层）
  const createDrops = (): CodeDrop[] => {
    if (!canvas) return []
    
    const drops: CodeDrop[] = []
    // 参考博客实现：固定列间距，简洁高效
    // 列间距约为字体大小的倍数，确保有足够的留白
    const columnWidth = props.fontSize * 1.2
    const columnCount = Math.ceil(canvas.width / columnWidth)
    // 根据密度计算实际列数，确保分布均匀
    const activeColumnCount = Math.max(1, Math.floor(columnCount * props.density * 100))
    
    // 参考博客实现：固定列间距分布，简洁高效
    // 创建固定间距的列，类似博客中的 index * 10 方式
    // 参考博客：arr = Array(Math.ceil(width / 10)).fill(0)
    const availableColumns: number[] = []
    for (let i = 0; i < columnCount; i++) {
      availableColumns.push(i)
    }
    
    // 辅助函数：创建指定层的雨滴
    const createLayerDrops = (count: number, layer: 'far' | 'mid' | 'near') => {
      const layerDrops: CodeDrop[] = []
      if (!canvas) return layerDrops
      
      // 随机打乱可用列
      const shuffledColumns = [...availableColumns].sort(() => Math.random() - 0.5)
      
      // 选择前 count 个列
      for (let i = 0; i < Math.min(count, shuffledColumns.length); i++) {
        const columnIndex = shuffledColumns[i]
        const x = columnIndex * columnWidth
        layerDrops.push(new CodeDrop(x, canvas.height, layer))
      }
      
      return layerDrops
    }
    
    // 如果启用分层，按比例分配各层
    if (props.enableLayers) {
      const farCount = Math.floor(activeColumnCount * 0.3) // 30% 远景
      const midCount = Math.floor(activeColumnCount * 0.5) // 50% 中景
      const nearCount = activeColumnCount - farCount - midCount // 剩余 近景
      
      // 各层均匀分布（允许列重叠，创造深度感）
      drops.push(...createLayerDrops(farCount, 'far'))
      drops.push(...createLayerDrops(midCount, 'mid'))
      drops.push(...createLayerDrops(nearCount, 'near'))
    } else {
      // 不启用分层，创建所有雨滴
      drops.push(...createLayerDrops(activeColumnCount, 'mid'))
    }
    
    return drops
  }

  let drops: CodeDrop[] = []
  
  const resizeCanvas = () => {
    if (!canvas) return
    const rect = canvas.getBoundingClientRect()
    canvas.width = rect.width || canvas.offsetWidth || window.innerWidth
    canvas.height = rect.height || canvas.offsetHeight || window.innerHeight
    drops = createDrops()
  }

  resizeCanvas()
  window.addEventListener('resize', resizeCanvas)
  
  drops = createDrops()

  // 动画循环
  let frameCount = 0
  const animate = () => {
    if (!ctx || !canvas) return

    const canvasWidth = canvas.width
    const canvasHeight = canvas.height

    // 更新颜色脉冲
    if (props.enableGlitch) {
      colorPulsePhase += colorPulseSpeed
    }

    // 屏幕抖动效果（周期性触发）
    if (props.enableGlitch) {
      glitchTimer++
      // 每 300-600 帧（约 5-10 秒）触发一次抖动
      if (glitchTimer > 300 + Math.random() * 300 && !glitchActive) {
        glitchActive = true
        glitchTimer = 0
      }
      
      if (glitchActive) {
        // 抖动持续 10-20 帧
        if (glitchTimer < 15) {
          glitchShakeX = (Math.random() - 0.5) * 2 // ±1 像素
          glitchShakeY = (Math.random() - 0.5) * 2
        } else {
          glitchActive = false
          glitchShakeX = 0
          glitchShakeY = 0
          glitchTimer = 0
        }
      }
    }

    // 保存上下文状态
    ctx.save()
    
    // 应用抖动偏移
    if (props.enableGlitch && glitchActive) {
      ctx.translate(glitchShakeX, glitchShakeY)
    }

    // 参考博客实现：残影效果，使用 rgba(0,0,0,0.05) 创建拖尾
    // 这是实现流动感和拖影的关键！
    ctx.fillStyle = 'rgba(0, 0, 0, 0.05)'
    ctx.fillRect(0, 0, canvasWidth, canvasHeight)
    
    // 禁用滤镜，确保清晰
    ctx.filter = 'none'

    // 如果启用分层，按层级顺序绘制（远景 -> 中景 -> 近景）
    if (props.enableLayers) {
      // 远景层（先绘制，在底层）
      drops.filter(d => d.layer === 'far').forEach((drop) => {
        drop.update(canvasHeight)
        drop.draw(ctx!, props.fontSize, props.color)
      })
      
      // 中景层
      drops.filter(d => d.layer === 'mid').forEach((drop) => {
        drop.update(canvasHeight)
        drop.draw(ctx!, props.fontSize, props.color)
      })
      
      // 近景层（最后绘制，在顶层）
      drops.filter(d => d.layer === 'near').forEach((drop) => {
        drop.update(canvasHeight)
        drop.draw(ctx!, props.fontSize, props.color)
      })
    } else {
      // 不启用分层，正常绘制
      drops.forEach((drop) => {
        drop.update(canvasHeight)
        drop.draw(ctx!, props.fontSize, props.color)
      })
    }

    // 恢复上下文状态
    ctx.restore()

    frameCount++
    animationFrameId = requestAnimationFrame(animate)
  }

  animate()

  // 清理函数
  return () => {
    window.removeEventListener('resize', resizeCanvas)
    if (animationFrameId) {
      cancelAnimationFrame(animationFrameId)
    }
  }
}

onMounted(() => {
  const cleanup = initCanvas()
  onUnmounted(() => {
    if (cleanup) cleanup()
  })
})

onUnmounted(() => {
  if (animationFrameId) {
    cancelAnimationFrame(animationFrameId)
  }
})
</script>

<style scoped lang="scss">
.code-rain-canvas {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 1;
  pointer-events: none;
  display: block;
  background: #000000;
}
</style>
