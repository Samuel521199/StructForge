<template>
  <canvas ref="canvasRef" class="code-rain-canvas" />
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, reactive } from 'vue'
import type { CodeRainProps } from './types'
import { loadConfig, mergeConfig, defaultConfig, type CodeRainConfig } from './config'

// 日文片假名字符集（默认）
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
  minLength: 0.8, // 相对于屏幕高度的倍数
  maxLength: 1.5, // 相对于屏幕高度的倍数
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
let rows: number = 0 // 行数
let columns: number = 0 // 列数

/**
 * 格子数据结构（基于"更新-生命"模型）
 */
interface GridCell {
  char: string // 当前字符（激活后到消失前不改变）
  alpha: number // 亮度/透明度 (0.0 - 1.0)
  speed: number // 每列独立的下落速度
  trailLength: number // 流的生命周期长度（剩余长度）
}

/**
 * 列状态管理（用于随机激活控制）
 */
interface ColumnState {
  waitFrames: number // 等待激活的帧数（随机）
  isActive: boolean // 当前是否有激活的流
}

/**
 * 待激活任务（全局管理）
 */
interface PendingActivation {
  targetColumn: number // 目标列索引
  waitFrames: number // 剩余等待帧数
}

let grid: GridCell[][] = [] // 二维数组：grid[列][行]
let columnStates: ColumnState[] = [] // 每列的状态管理
let pendingActivations: PendingActivation[] = [] // 待激活任务列表（全局）

/**
 * 从字符集中随机获取一个字符
 */
const getRandomCharacter = (): string => {
  const chars = config.characters ?? props.characters ?? KATAKANA_CHARS
  return chars.charAt(Math.floor(Math.random() * chars.length))
}

/**
 * 激活一列新的流（生命周期开始）
 */
const startNewStream = (colIndex: number, minLength: number, maxLength: number) => {
  if (!grid[colIndex] || grid[colIndex].length === 0) {
    // 临时关闭日志输出以便调试注册功能
    // if (import.meta.env.DEV) {
    //   console.warn(`[CodeRain] ⚠️ 无法激活列 ${colIndex}：格子未初始化`)
    // }
    return
  }

  // 设置流的长度（生命）
  const trailLength = minLength + Math.floor(Math.random() * (maxLength - minLength))

  // 从顶部开始激活
  const headRow = 0
  grid[colIndex][headRow].trailLength = trailLength
  grid[colIndex][headRow].alpha = 1.0 // 头部最亮
  grid[colIndex][headRow].char = getRandomCharacter() // 头部字符（激活时确定，之后不改变）
  
  // 标记列为激活状态
  columnStates[colIndex].isActive = true
  
  // 临时关闭日志输出以便调试注册功能
  // if (import.meta.env.DEV) {
  //   console.log(`[CodeRain] ✨ 激活列 ${colIndex}，流长度: ${trailLength}`)
  // }
}

/**
 * 初始化格子系统
 */
const initGrid = () => {
  const fontSize = config.fontSize ?? props.fontSize ?? 15
  const columnSpacingFactor = config.columnSpacingFactor ?? 2.5
  const columnSpacing = fontSize * columnSpacingFactor

  // 重新计算行列数
  columns = Math.floor(canvasWidth / columnSpacing)
  rows = Math.floor(canvasHeight / fontSize)

  grid = []

  // 初始化所有格子和列状态
  columnStates = []
  pendingActivations = [] // 清空待激活列表
  
  for (let i = 0; i < columns; i++) {
    grid[i] = []
    for (let j = 0; j < rows; j++) {
      grid[i][j] = {
        char: getRandomCharacter(),
        alpha: 0, // 初始都是暗的
        speed: 0.3 + Math.random() * 0.4, // 随机速度（0.3-0.7）
        trailLength: 0, // 初始无流
      }
    }
    
    // 初始化列状态
    columnStates[i] = {
      waitFrames: 0,
      isActive: false,
    }
  }

  // 初始化时，随机激活几列
  // 生命周期长度基于屏幕高度计算（确保能贯穿屏幕）
  const minLengthRatio = config.minLength ?? props.minLength ?? 0.8
  const maxLengthRatio = config.maxLength ?? props.maxLength ?? 1.5
  const density = config.density ?? props.density ?? 0.008
  
  // 计算实际的生命周期长度（基于屏幕行数）
  const screenRows = rows // 屏幕可显示的行数
  const minLength = Math.max(10, Math.floor(screenRows * minLengthRatio)) // 至少 10 个字符
  const maxLength = Math.max(minLength + 5, Math.floor(screenRows * maxLengthRatio)) // 至少比 minLength 大 5
  
  let activatedCount = 0
  for (let i = 0; i < columns; i++) {
    if (Math.random() < density * 100) {
      // 根据密度概率激活一列
      startNewStream(i, minLength, maxLength)
      columnStates[i].isActive = true
      activatedCount++
    }
  }
  
  // 如果密度太低导致没有激活任何列，至少激活一列
  if (activatedCount === 0 && columns > 0) {
    const randomCol = Math.floor(Math.random() * columns)
    startNewStream(randomCol, minLength, maxLength)
    columnStates[randomCol].isActive = true
    
    // 临时关闭日志输出以便调试注册功能
    // if (import.meta.env.DEV) {
    //   console.log(`[CodeRain] 🎬 初始激活列 ${randomCol}（密度太低，强制激活）`)
    // }
  }
  // 临时关闭日志输出以便调试注册功能
  // if (import.meta.env.DEV && activatedCount > 0) {
  //   console.log(`[CodeRain] 🎬 初始激活了 ${activatedCount} 列`)
  // }
}

/**
 * 重新计算 Canvas 尺寸和格子系统
 */
const resizeCanvas = () => {
  if (!canvas) return

  // 设置 Canvas 实际像素尺寸（解决模糊问题）
  const rect = canvas.getBoundingClientRect()
  canvasWidth = canvas.width = rect.width
  canvasHeight = canvas.height = rect.height

  // 重新初始化格子
  initGrid()
}

/**
 * 核心渲染函数：每一帧的绘制逻辑（基于"点亮-更新-死亡"模型）
 */
const draw = () => {
  if (!ctx || !canvas) return

  // 1. 残影清屏（电影级拖尾的关键）
  // 极低透明度的黑色，创建字符的残影和拖尾效果，消除所有静止的背景干扰
  ctx.fillStyle = 'rgba(0, 0, 0, 0.04)'
  ctx.fillRect(0, 0, canvasWidth, canvasHeight)

  // 禁用滤镜
  ctx.filter = 'none'

  // 从配置中获取参数
  const fontSize = config.fontSize ?? props.fontSize ?? 15
  const fontWeight = config.fontWeight ?? props.fontWeight ?? 'normal'
  const fontFamily = config.fontFamily ?? props.fontFamily ?? 'Monospace'
  const opacity = config.opacity ?? props.opacity ?? 0.9
  const enableGlow = config.enableGlow ?? props.enableGlow ?? true
  const glowIntensity = config.glowIntensity ?? props.glowIntensity ?? 0.8
  const fadeSpeed = config.fadeSpeed ?? props.fadeSpeed ?? 0.04
  const columnSpacingFactor = config.columnSpacingFactor ?? 2.5
  const columnSpacing = fontSize * columnSpacingFactor
  const speed = config.speed ?? props.speed ?? 2.5
  // 生命周期长度基于屏幕高度计算（确保能贯穿屏幕）
  const minLengthRatio = config.minLength ?? props.minLength ?? 0.8
  const maxLengthRatio = config.maxLength ?? props.maxLength ?? 1.5
  const minWaitTime = config.minWaitTime ?? 0.1 // 最小等待时间（秒）
  const maxWaitTime = config.maxWaitTime ?? 5.0 // 最大等待时间（秒）
  
  // 计算实际的生命周期长度（基于屏幕行数）
  const screenRows = rows // 屏幕可显示的行数
  const minLength = Math.max(10, Math.floor(screenRows * minLengthRatio)) // 至少 10 个字符
  const maxLength = Math.max(minLength + 5, Math.floor(screenRows * maxLengthRatio)) // 至少比 minLength 大 5

  // 设置基础字体样式
  ctx.font = `${fontWeight} ${fontSize}px ${fontFamily}`

  // 2. 遍历所有列，处理激活等待和流的状态
  for (let col = 0; col < columns; col++) {
    const colState = columnStates[col]
    let hasActiveStream = false // 检查该列是否有激活的流

    // 遍历该列的所有格子（从上往下，符合流的移动方向）
    for (let row = 0; row < rows; row++) {
      const cell = grid[col][row]
      const x = col * columnSpacing
      const y = (row + 1) * fontSize // +1 因为 fillText 的 y 是基线位置

      // 检查是否有激活的流（trailLength > 0）
      if (cell.trailLength > 0) {
        hasActiveStream = true
        
        // A. 流的头部逻辑（点亮，字符不改变）
        // 头部特征：alpha >= 0.99 且 trailLength > 0
        if (cell.alpha >= 0.99) {
          // 头部：极亮绿色高光（电影感）
          ctx.fillStyle = `rgba(0, 255, 0, ${cell.alpha * opacity})` // 纯绿色，最亮
          if (enableGlow) {
            ctx.shadowColor = '#00FF00' // 绿色光晕
            ctx.shadowBlur = 20 * glowIntensity // 强大的发光效果
          } else {
            ctx.shadowBlur = 0
            ctx.shadowColor = 'transparent'
          }

          // 注意：字符在激活时已确定，这里不再更新

          // 头部向下移动（驱动流）
          // 提高移动速度，确保流能正常移动
          const moveSpeed = Math.min(1.0, cell.speed * speed * 0.3) // 速度控制（提高移动概率）
          if (row < rows - 1 && Math.random() < moveSpeed) {
            // 移动到下一行
            const nextRow = row + 1
            const nextCell = grid[col][nextRow]

            // 传递生命值
            nextCell.trailLength = cell.trailLength - 1
            nextCell.alpha = 1.0 // 新头部最亮
            nextCell.char = getRandomCharacter() // 新头部字符（激活时确定）

            // 当前头部变为尾迹（开始衰减）
            cell.alpha = 0.8
            cell.trailLength = 0 // 头部完成使命
          } else if (row >= rows - 1) {
            // 如果头部已经到达底部，生命值耗尽
            cell.trailLength = 0
            cell.alpha = 0.8 // 开始衰减
          }

          // 如果流的生命值耗尽，头部死亡
          if (cell.trailLength <= 0) {
            cell.alpha *= (1 - fadeSpeed * 2) // 快速衰减
            if (cell.alpha < 0.1) {
              cell.alpha = 0
              cell.trailLength = 0
            }
          }
        } else if (cell.alpha > 0) {
          // B. 尾迹逻辑：渐变衰减（死亡过程）
          // 注意：字符不改变，保持激活时的字符
          const decayFactor = 1 - fadeSpeed // 衰减速度
          cell.alpha *= decayFactor

          if (cell.alpha < 0.05) {
            cell.alpha = 0 // 彻底死亡
            cell.trailLength = 0
          } else {
            // 尾迹：渐变绿
            const green = Math.floor(255 * cell.alpha * opacity)
            ctx.fillStyle = `rgba(0, ${green}, 0, ${cell.alpha * opacity})`

            if (enableGlow) {
              ctx.shadowColor = `rgba(0, 255, 0, ${cell.alpha})`
              ctx.shadowBlur = 5 * glowIntensity * cell.alpha // 光晕随亮度衰减
            } else {
              ctx.shadowBlur = 0
              ctx.shadowColor = 'transparent'
            }
          }
        }

        // 绘制字符（如果可见）
        if (cell.alpha > 0) {
          ctx.fillText(cell.char, x, y)
        }
      }
    }

    // 更新列状态
    // 关键修复：只有当有激活的流（trailLength > 0）时才认为列是激活的
    // 尾迹（alpha > 0 但 trailLength = 0）不应该阻止列完成
    // 这样即使尾迹还在衰减，列也可以完成并触发新激活
    const wasActive = colState.isActive
    const isNowActive = hasActiveStream // 只检查是否有激活的流，不检查尾迹
    
    if (isNowActive) {
      // 有激活的流，标记为激活状态
      colState.isActive = true
    } else {
      // 没有激活的流（即使还有尾迹在衰减）
      if (wasActive && !isNowActive) {
        // 流刚完成（从有到无），随机选择一个未激活的列来激活
        colState.isActive = false
        
        // 临时关闭日志输出以便调试注册功能
        // if (import.meta.env.DEV) {
        //   // 检查是否还有可见的尾迹
        //   let hasAnyVisibleAlpha = false
        //   for (let checkRow = 0; checkRow < rows; checkRow++) {
        //     if (grid[col][checkRow].alpha > 0) {
        //       hasAnyVisibleAlpha = true
        //       break
        //     }
        //   }
        //   console.log(`[CodeRain] 📍 列 ${col} 完成（hasActiveStream: ${hasActiveStream}, hasAnyVisibleAlpha: ${hasAnyVisibleAlpha}）`)
        // }
        
        // 找到所有未激活的列（允许有尾迹，只要没有激活的流即可）
        const availableColumns: number[] = []
        for (let checkCol = 0; checkCol < columns; checkCol++) {
          if (checkCol === col) continue // 跳过当前列
          if (columnStates[checkCol].isActive) continue // 跳过已激活的列（有激活的流）
          
          // 检查该列是否有激活的流（trailLength > 0）
          let hasActiveStream = false
          for (let checkRow = 0; checkRow < rows; checkRow++) {
            if (grid[checkCol][checkRow].trailLength > 0) {
              hasActiveStream = true
              break
            }
          }
          
          // 只要没有激活的流，就可以激活（允许有尾迹）
          if (!hasActiveStream) {
            availableColumns.push(checkCol)
          }
        }
        
        // 临时关闭日志输出以便调试注册功能
        // if (import.meta.env.DEV) {
        //   console.log(`[CodeRain] 🔍 列 ${col} 完成，找到 ${availableColumns.length} 个可用列:`, availableColumns)
        // }
        
        // 如果有可用的列，随机选择一个并加入待激活列表
        if (availableColumns.length > 0) {
          const targetCol = availableColumns[Math.floor(Math.random() * availableColumns.length)]
          // 随机等待时间（从配置的范围，转换为帧数，假设 60fps）
          const fps = 60
          const waitTimeSeconds = minWaitTime + Math.random() * (maxWaitTime - minWaitTime)
          const waitFrames = Math.max(1, Math.floor(waitTimeSeconds * fps))
          
          pendingActivations.push({
            targetColumn: targetCol,
            waitFrames: waitFrames,
          })
          
          // 临时关闭日志输出以便调试注册功能
          // if (import.meta.env.DEV) {
          //   console.log(`[CodeRain] ⏰ 列 ${col} 完成，将在 ${waitTimeSeconds.toFixed(2)} 秒后激活列 ${targetCol}`)
          // }
        } else {
          // 如果没有可用列，尝试在当前列等待后重新激活（避免完全停止）
          const fps = 60
          const waitTimeSeconds = minWaitTime + Math.random() * (maxWaitTime - minWaitTime)
          const waitFrames = Math.max(1, Math.floor(waitTimeSeconds * fps))
          
          pendingActivations.push({
            targetColumn: col, // 如果没有其他可用列，等待后重新激活当前列
            waitFrames: waitFrames,
          })
          
          // 临时关闭日志输出以便调试注册功能
          // if (import.meta.env.DEV) {
          //   console.log(`[CodeRain] ⏰ 列 ${col} 完成，没有可用列，将在 ${waitTimeSeconds.toFixed(2)} 秒后重新激活当前列`)
          // }
        }
      } else {
        // 列一直处于非激活状态，不需要处理
        colState.isActive = false
      }
    }
  }

  // 处理待激活任务（全局）
  for (let i = pendingActivations.length - 1; i >= 0; i--) {
    const task = pendingActivations[i]
    task.waitFrames--
    
    if (task.waitFrames <= 0) {
      // 等待时间到了，检查目标列是否仍然可用
      const targetCol = task.targetColumn
      
      // 检查该列是否有激活的流（trailLength > 0）
      let hasActiveStream = false
      for (let row = 0; row < rows; row++) {
        if (grid[targetCol][row].trailLength > 0) {
          hasActiveStream = true
          break
        }
      }
      
      // 只要没有激活的流且未激活，就可以激活（允许有尾迹）
      if (!columnStates[targetCol].isActive && !hasActiveStream) {
        // 列可用，激活它
        startNewStream(targetCol, minLength, maxLength)
        
        // 临时关闭日志输出以便调试注册功能
        // if (import.meta.env.DEV) {
        //   console.log(`[CodeRain] ✅ 激活列 ${targetCol}`)
        // }
        
        // 移除已完成的任务
        pendingActivations.splice(i, 1)
      } else {
        // 如果列不可用（已激活或有激活的流），等待更长时间后再试
        const fps = 60
        const waitTimeSeconds = minWaitTime + Math.random() * (maxWaitTime - minWaitTime)
        task.waitFrames = Math.max(1, Math.floor(waitTimeSeconds * fps))
        
        // 临时关闭日志输出以便调试注册功能
        // if (import.meta.env.DEV) {
        //   // 检查是否有尾迹（用于调试）
        //   let hasAnyAlpha = false
        //   for (let row = 0; row < rows; row++) {
        //     if (grid[targetCol][row].alpha > 0) {
        //       hasAnyAlpha = true
        //       break
        //     }
        //   }
        //   console.log(`[CodeRain] ⏳ 列 ${targetCol} 不可用（isActive: ${columnStates[targetCol].isActive}, hasActiveStream: ${hasActiveStream}, hasAlpha: ${hasAnyAlpha}），等待 ${waitTimeSeconds.toFixed(2)} 秒后再试`)
        // }
        // 不删除任务，继续等待
      }
    }
  }
  
  // 如果没有任何激活的列和待激活任务，确保至少激活一列（防止完全停止）
  let hasAnyActiveColumn = false
  for (let col = 0; col < columns; col++) {
    if (columnStates[col].isActive) {
      hasAnyActiveColumn = true
      break
    }
  }
  
  if (!hasAnyActiveColumn && pendingActivations.length === 0) {
    // 没有任何激活的列和待激活任务，随机激活一列
    const availableColumns: number[] = []
    for (let col = 0; col < columns; col++) {
      if (columnStates[col].isActive) continue
      
      // 检查是否有激活的流（trailLength > 0）
      let hasActiveStream = false
      for (let row = 0; row < rows; row++) {
        if (grid[col][row].trailLength > 0) {
          hasActiveStream = true
          break
        }
      }
      
      // 只要没有激活的流，就可以激活（允许有尾迹）
      if (!hasActiveStream) {
        availableColumns.push(col)
      }
    }
    
    if (availableColumns.length > 0) {
      const targetCol = availableColumns[Math.floor(Math.random() * availableColumns.length)]
      startNewStream(targetCol, minLength, maxLength)
      
      // 临时关闭日志输出以便调试注册功能
      // if (import.meta.env.DEV) {
      //   console.log(`[CodeRain] 🔄 检测到所有列都停止，立即激活列 ${targetCol}（可用列: ${availableColumns.length}）`)
      // }
    } else {
      // 如果所有列都有激活的流，等待一小段时间后再检查
      // 临时关闭日志输出以便调试注册功能
      // if (import.meta.env.DEV) {
      //   console.warn(`[CodeRain] ⚠️ 所有列都有激活的流，无法立即激活，等待下一帧检查`)
      // }
    }
  }
  
  // 定期检查：如果长时间没有激活的列，强制激活一列（每 300 帧检查一次，约 5 秒）
  if (typeof (window as any).__codeRainLastCheck === 'undefined') {
    (window as any).__codeRainLastCheck = 0
  }
  (window as any).__codeRainLastCheck++
  
  if ((window as any).__codeRainLastCheck >= 300) {
    (window as any).__codeRainLastCheck = 0
    
    if (!hasAnyActiveColumn && pendingActivations.length === 0) {
      // 强制激活一列
      const allColumns: number[] = []
      for (let col = 0; col < columns; col++) {
        allColumns.push(col)
      }
      
      if (allColumns.length > 0) {
        const targetCol = allColumns[Math.floor(Math.random() * allColumns.length)]
        startNewStream(targetCol, minLength, maxLength)
        
        // 临时关闭日志输出以便调试注册功能
        // if (import.meta.env.DEV) {
        //   console.log(`[CodeRain] 🔄 定期检查：强制激活列 ${targetCol}`)
        // }
      }
    }
  }

  // 重置阴影
  ctx.shadowBlur = 0
  ctx.shadowColor = 'transparent'

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
      if (props.minLength !== 0.8) filteredProps.minLength = props.minLength
      if (props.maxLength !== 1.5) filteredProps.maxLength = props.maxLength
      if (props.enableLayers !== false) filteredProps.enableLayers = props.enableLayers
      if (props.enableGlow !== true) filteredProps.enableGlow = props.enableGlow
      if (props.enableGlitch !== false) filteredProps.enableGlitch = props.enableGlitch

      // 合并配置（优先级：显式传递的 props > 配置文件 > 默认值）
      const mergedConfig = mergeConfig(filteredProps, fileConfig, defaultConfig)
      Object.assign(config, mergedConfig)

      // 临时关闭日志输出以便调试注册功能
      // if (import.meta.env.DEV) {
      //   console.log('[CodeRain] 配置加载完成:', {
      //     characters: config.characters?.substring(0, 30) + '...',
      //     charactersLength: config.characters?.length,
      //     fromFile: fileConfig.characters?.substring(0, 30) + '...',
      //     fromProps: props.characters?.substring(0, 30) + '...',
      //     isDefaultChars: props.characters === KATAKANA_CHARS,
      //     speed: config.speed,
      //     fontSize: config.fontSize,
      //   })
      // }
    } catch (error) {
      // 临时关闭日志输出以便调试注册功能
      // console.warn('[CodeRain] 配置文件加载失败，使用 props 和默认值:', error)
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
  if (animationFrameId) {
    cancelAnimationFrame(animationFrameId)
  }
  window.removeEventListener('resize', resizeCanvas)
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
