/**
 * CodeRain 配置文件类型定义
 */
export interface CodeRainConfig {
  fontSize?: number
  fontFamily?: string
  fontWeight?: string | number
  color?: string
  backgroundColor?: string
  speed?: number
  speedVariation?: number
  density?: number
  opacity?: number
  fadeSpeed?: number
  minLength?: number
  maxLength?: number
  enableLayers?: boolean
  enableGlow?: boolean
  enableGlitch?: boolean
  glowIntensity?: number
  columnSpacingFactor?: number
  trailLength?: number
  characters?: string
}

export interface CodeRainConfigFile {
  codeRain: CodeRainConfig
}

/**
 * 默认配置
 */
export const defaultConfig: CodeRainConfig = {
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
  minLength: 15,
  maxLength: 35,
  enableLayers: false,
  enableGlow: true,
  enableGlitch: false,
  glowIntensity: 0.8,
  columnSpacingFactor: 2.5,
  trailLength: 20,
  characters: 'アァカサタナハマヤラワガザダバパ0123456789!@#$%^&*()_+-=[]{};:"|,./<>?',
}

/**
 * 去除 JSON 中的注释字段（_comment 开头的字段）
 * 使用更安全的方法：先尝试解析 JSON，然后过滤注释字段
 * @param jsonString 带注释的 JSON 字符串
 * @returns 去除注释后的 JSON 对象
 */
function removeJsonComments(jsonString: string): any {
  try {
    // 先尝试直接解析（如果 JSON 本身没有语法错误）
    const parsed = JSON.parse(jsonString)
    
    // 递归移除所有以 _comment 开头的字段
    const removeComments = (obj: any): any => {
      if (Array.isArray(obj)) {
        return obj.map(removeComments)
      } else if (obj !== null && typeof obj === 'object') {
        const cleaned: any = {}
        for (const [key, value] of Object.entries(obj)) {
          if (!key.startsWith('_comment')) {
            cleaned[key] = removeComments(value)
          }
        }
        return cleaned
      }
      return obj
    }
    
    return removeComments(parsed)
  } catch (error) {
    // 如果直接解析失败，尝试用正则移除注释字段后再解析
    let result = jsonString
    
    // 移除单行注释 //
    result = result.replace(/\/\/.*$/gm, '')
    
    // 移除多行注释 /* */
    result = result.replace(/\/\*[\s\S]*?\*\//g, '')
    
    // 移除以 _comment 开头的注释字段（更精确的正则）
    // 匹配整个键值对，包括引号、冒号、值和逗号
    result = result.replace(/"\s*_comment[^"]*"\s*:\s*"[^"]*"\s*,?\s*/g, '')
    result = result.replace(/"\s*_comment[^"]*"\s*:\s*[^,}\]]+\s*,?\s*/g, '')
    
    // 清理多余的逗号（避免 JSON 解析错误）
    result = result.replace(/,\s*}/g, '}')
    result = result.replace(/,\s*]/g, ']')
    result = result.replace(/,\s*,/g, ',') // 清理连续逗号
    
    // 再次尝试解析
    try {
      return JSON.parse(result)
    } catch (parseError) {
      // 如果还是失败，抛出原始错误
      throw new Error(`JSON 解析失败: ${parseError instanceof Error ? parseError.message : String(parseError)}`)
    }
  }
}

/**
 * 加载配置文件
 * @param configPath 配置文件路径（可选，默认为 ./config.json）
 * @returns 配置对象
 */
export async function loadConfig(configPath?: string): Promise<CodeRainConfig> {
  try {
    let configData: CodeRainConfigFile
    let jsonString: string
    
    if (configPath) {
      // 使用指定的配置文件路径
      const response = await fetch(configPath)
      if (!response.ok) {
        throw new Error(`配置文件加载失败: ${response.statusText}`)
      }
      jsonString = await response.text()
    } else {
      // 尝试从默认位置加载（使用 import 方式，Vite 会自动处理）
      try {
        // 方式 1：尝试直接 import（Vite 支持）
        const configModule = await import('./config.json?raw')
        jsonString = configModule.default
      } catch {
        // 方式 2：fallback 到 fetch
        const response = await fetch(new URL('./config.json', import.meta.url).href)
        if (!response.ok) {
          throw new Error(`默认配置文件加载失败: ${response.statusText}`)
        }
        jsonString = await response.text()
      }
    }
    
    // 去除注释后解析 JSON（removeJsonComments 已经返回解析后的对象）
    const cleanedData = removeJsonComments(jsonString)
    
    // 确保数据结构正确
    if (cleanedData && cleanedData.codeRain) {
      configData = { codeRain: cleanedData.codeRain }
    } else {
      throw new Error('配置文件格式错误：缺少 codeRain 字段')
    }
    
    // 调试：输出加载的配置（仅开发环境）
    if (import.meta.env.DEV) {
      console.log('[CodeRain] 配置文件解析成功:', {
        characters: configData.codeRain.characters,
        charactersLength: configData.codeRain.characters?.length,
        speed: configData.codeRain.speed,
        fontSize: configData.codeRain.fontSize,
        allKeys: Object.keys(configData.codeRain)
      })
    }
    
    const finalConfig = {
      ...defaultConfig,
      ...configData.codeRain,
    }
    
    // 调试：输出最终配置（仅开发环境）
    if (import.meta.env.DEV) {
      console.log('[CodeRain] 最终配置:', {
        characters: finalConfig.characters,
        charactersLength: finalConfig.characters?.length,
        isDefault: finalConfig.characters === defaultConfig.characters
      })
    }
    
    // 确保 characters 字段存在
    if (!finalConfig.characters) {
      console.warn('[CodeRain] 配置中缺少 characters 字段，使用默认值')
      finalConfig.characters = defaultConfig.characters
    }
    
    return finalConfig
  } catch (error) {
    console.warn('[CodeRain] 配置文件加载失败，使用默认配置:', error)
    return defaultConfig
  }
}

/**
 * 合并配置（优先级：props > 配置文件 > 默认值）
 */
export function mergeConfig(
  props: Partial<CodeRainConfig>,
  fileConfig: CodeRainConfig,
  defaults: CodeRainConfig = defaultConfig
): CodeRainConfig {
  return {
    ...defaults,
    ...fileConfig,
    ...props,
  }
}

