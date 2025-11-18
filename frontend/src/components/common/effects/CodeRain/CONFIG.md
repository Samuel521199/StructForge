# CodeRain 配置文件使用指南

## 概述

CodeRain 组件支持通过 JSON 配置文件来控制所有效果参数，无需修改代码。配置文件提供了灵活的配置方式，支持动态加载和热更新。

## 配置优先级

配置参数的优先级顺序（从高到低）：
1. **Props**（组件属性）- 最高优先级
2. **配置文件**（JSON）- 中等优先级
3. **默认值**（代码中的默认值）- 最低优先级

## 使用方法

### 方式 1：使用默认配置文件

在组件目录下创建 `config.json` 文件，组件会自动加载：

```vue
<template>
  <CodeRain :useConfigFile="true" />
</template>
```

### 方式 2：指定自定义配置文件路径

```vue
<template>
  <CodeRain 
    :useConfigFile="true"
    :configPath="'/configs/code-rain.json'"
  />
</template>
```

### 方式 3：混合使用（推荐）

Props 用于覆盖配置文件中的特定参数：

```vue
<template>
  <CodeRain 
    :useConfigFile="true"
    :speed="3.0"
    :glowIntensity="0.9"
  />
</template>
```

## 配置文件格式

### 标准格式（带注释说明）

配置文件支持使用 `_comment` 字段添加注释说明，这些注释字段会在加载时自动过滤掉，不会影响配置解析：

```json
{
  "codeRain": {
    "_comment_fontSize": "字体大小（像素），建议范围：12-20",
    "fontSize": 15,
    
    "_comment_speed": "下落速度（像素/帧），值越大速度越快，建议范围：1.0-5.0",
    "speed": 2.5,
    
    "_comment_density": "雨滴密度（0-1），值越大密度越高，建议范围：0.003-0.015",
    "density": 0.008,
    
    "_comment_enableGlow": "是否启用光晕效果。true: 启用光晕，视觉效果更佳；false: 禁用光晕，性能更好",
    "enableGlow": true
  }
}
```

### 完整示例（带所有注释）

查看 `config.json` 文件，其中包含了所有参数的详细注释说明。注释字段以 `_comment` 开头，例如：
- `_comment_fontSize`: 说明 fontSize 参数的作用
- `_comment_speed`: 说明 speed 参数的作用
- `_comment`: 分组注释，用于区分不同类别的参数

### 参数说明

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `fontSize` | `number` | `15` | 字体大小（像素） |
| `fontFamily` | `string` | `"Monospace"` | 字体族 |
| `fontWeight` | `string \| number` | `"normal"` | 字体粗细 |
| `color` | `string` | `"#00FF00"` | 代码颜色（十六进制） |
| `backgroundColor` | `string` | `"#000000"` | 背景颜色（十六进制） |
| `speed` | `number` | `2.5` | 下落速度（像素/帧） |
| `speedVariation` | `number` | `0.6` | 速度变化范围（0-1） |
| `density` | `number` | `0.008` | 雨滴密度（0-1） |
| `opacity` | `number` | `0.9` | 透明度（0-1） |
| `fadeSpeed` | `number` | `0.04` | 拖尾淡出速度（0-1） |
| `minLength` | `number` | `15` | 代码串最小长度 |
| `maxLength` | `number` | `35` | 代码串最大长度 |
| `enableLayers` | `boolean` | `false` | 是否启用三层深度效果 |
| `enableGlow` | `boolean` | `true` | 是否启用光晕效果 |
| `enableGlitch` | `boolean` | `false` | 是否启用动态干扰 |
| `glowIntensity` | `number` | `0.8` | 光晕强度（0-1） |
| `columnSpacingFactor` | `number` | `2.5` | 列间距因子 |
| `trailLength` | `number` | `20` | 尾迹长度（字符数） |
| `characters` | `string` | 日文片假名 | 字符集 |

## 配置文件示例

### 示例 1：快速效果（高速度、高密度）

```json
{
  "codeRain": {
    "speed": 4.0,
    "density": 0.015,
    "opacity": 0.95,
    "glowIntensity": 0.9
  }
}
```

### 示例 2：慢速优雅效果（低速度、低密度）

```json
{
  "codeRain": {
    "speed": 1.5,
    "density": 0.003,
    "opacity": 0.7,
    "glowIntensity": 0.5,
    "trailLength": 30
  }
}
```

### 示例 3：自定义字符集

```json
{
  "codeRain": {
    "characters": "01",
    "fontSize": 12,
    "density": 0.01
  }
}
```

## 配置文件位置

### 默认位置

默认配置文件位置：`frontend/src/components/common/effects/CodeRain/config.json`

### 自定义位置

可以通过 `configPath` prop 指定任意路径：

```vue
<CodeRain 
  :useConfigFile="true"
  :configPath="'/public/configs/code-rain.json'"
/>
```

## 动态加载配置

配置文件支持动态加载，可以在运行时切换不同的配置文件：

```typescript
// 在组件中动态切换配置
const configPath = ref('/configs/fast.json')

// 切换配置
configPath.value = '/configs/slow.json'
```

## 注意事项

1. **配置文件格式**：必须是有效的 JSON 格式
2. **路径问题**：配置文件路径相对于项目根目录或使用绝对路径
3. **加载失败**：如果配置文件加载失败，组件会自动使用默认值
4. **性能考虑**：配置文件在组件挂载时加载一次，不会频繁读取

## 最佳实践

1. **开发环境**：使用 props 快速调整参数
2. **生产环境**：使用配置文件统一管理参数
3. **多场景**：为不同场景创建不同的配置文件
4. **版本控制**：将配置文件纳入版本控制，便于团队协作

## 故障排查

### 配置文件未生效

1. 检查 `useConfigFile` 是否为 `true`
2. 检查配置文件路径是否正确
3. 检查 JSON 格式是否有效
4. 查看浏览器控制台的错误信息

### 配置参数被覆盖

记住优先级顺序：**Props > 配置文件 > 默认值**

如果 props 中设置了参数，会覆盖配置文件中的相同参数。

