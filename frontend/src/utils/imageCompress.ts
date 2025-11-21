/**
 * 图片压缩和裁剪工具
 */

export interface CompressOptions {
  /** 最大宽度（像素） */
  maxWidth?: number
  /** 最大高度（像素） */
  maxHeight?: number
  /** 质量（0-1，仅适用于 JPEG） */
  quality?: number
  /** 输出格式 */
  outputFormat?: 'image/jpeg' | 'image/png' | 'image/webp'
}

/**
 * 压缩图片
 */
export async function compressImage(
  file: File,
  options: CompressOptions = {}
): Promise<File> {
  const {
    maxWidth = 512,
    maxHeight = 512,
    quality = 0.85,
    outputFormat = 'image/jpeg',
  } = options

  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = (e) => {
      const img = new Image()
      img.onload = () => {
        // 计算新尺寸（保持宽高比）
        let width = img.width
        let height = img.height

        if (width > maxWidth || height > maxHeight) {
          const ratio = Math.min(maxWidth / width, maxHeight / height)
          width = width * ratio
          height = height * ratio
        }

        // 创建 canvas
        const canvas = document.createElement('canvas')
        canvas.width = width
        canvas.height = height
        const ctx = canvas.getContext('2d')

        if (!ctx) {
          reject(new Error('无法创建 canvas context'))
          return
        }

        // 绘制图片
        ctx.drawImage(img, 0, 0, width, height)

        // 转换为 Blob
        canvas.toBlob(
          (blob) => {
            if (!blob) {
              reject(new Error('图片压缩失败'))
              return
            }

            // 创建新的 File 对象
            const compressedFile = new File(
              [blob],
              file.name,
              { type: outputFormat }
            )
            resolve(compressedFile)
          },
          outputFormat,
          quality
        )
      }
      img.onerror = () => {
        reject(new Error('图片加载失败'))
      }
      img.src = e.target?.result as string
    }
    reader.onerror = () => {
      reject(new Error('文件读取失败'))
    }
    reader.readAsDataURL(file)
  })
}

/**
 * 裁剪图片（圆形）
 */
export async function cropImageToCircle(
  file: File,
  size: number = 512
): Promise<File> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = (e) => {
      const img = new Image()
      img.onload = () => {
        // 创建 canvas
        const canvas = document.createElement('canvas')
        canvas.width = size
        canvas.height = size
        const ctx = canvas.getContext('2d')

        if (!ctx) {
          reject(new Error('无法创建 canvas context'))
          return
        }

        // 计算裁剪区域（保持宽高比，居中裁剪）
        const minSize = Math.min(img.width, img.height)
        const x = (img.width - minSize) / 2
        const y = (img.height - minSize) / 2

        // 绘制圆形裁剪路径
        ctx.beginPath()
        ctx.arc(size / 2, size / 2, size / 2, 0, Math.PI * 2)
        ctx.clip()

        // 绘制图片
        ctx.drawImage(img, x, y, minSize, minSize, 0, 0, size, size)

        // 转换为 Blob
        canvas.toBlob(
          (blob) => {
            if (!blob) {
              reject(new Error('图片裁剪失败'))
              return
            }

            // 创建新的 File 对象
            const croppedFile = new File(
              [blob],
              file.name,
              { type: 'image/png' }
            )
            resolve(croppedFile)
          },
          'image/png',
          0.95
        )
      }
      img.onerror = () => {
        reject(new Error('图片加载失败'))
      }
      img.src = e.target?.result as string
    }
    reader.onerror = () => {
      reject(new Error('文件读取失败'))
    }
    reader.readAsDataURL(file)
  })
}

/**
 * 压缩并裁剪头像
 */
export async function processAvatar(
  file: File,
  options: CompressOptions & { cropToCircle?: boolean; size?: number } = {}
): Promise<File> {
  const { cropToCircle = false, size = 512, ...compressOptions } = options

  // 先压缩
  let processedFile = await compressImage(file, {
    maxWidth: size,
    maxHeight: size,
    ...compressOptions,
  })

  // 如果需要裁剪为圆形
  if (cropToCircle) {
    processedFile = await cropImageToCircle(processedFile, size)
  }

  return processedFile
}

