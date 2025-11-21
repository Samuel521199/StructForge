/**
 * Upload 组件类型定义
 */

export interface UploadFile {
  name: string
  url?: string
  status?: 'ready' | 'uploading' | 'success' | 'fail'
  uid?: number
  response?: any
  percentage?: number
}

export interface UploadRequestOptions {
  action: string
  file: File
  data?: Record<string, any>
  filename?: string
  headers?: Record<string, string>
  onProgress?: (event: ProgressEvent) => void
  onSuccess?: (response: any) => void
  onError?: (error: Error) => void
}

export interface UploadProps {
  /** 上传的地址 */
  action: string
  /** 设置上传的请求头部 */
  headers?: Record<string, string>
  /** 上传时附带的额外参数 */
  data?: Record<string, any>
  /** 是否支持多选文件 */
  multiple?: boolean
  /** 接受上传的文件类型 */
  accept?: string
  /** 最大允许上传个数 */
  limit?: number
  /** 文件列表 */
  fileList?: UploadFile[]
  /** 是否自动上传 */
  autoUpload?: boolean
  /** 是否启用拖拽上传 */
  drag?: boolean
  /** 是否禁用 */
  disabled?: boolean
  /** 是否显示已上传文件列表 */
  showFileList?: boolean
  /** 文件列表的类型 */
  listType?: 'text' | 'picture' | 'picture-card'
  /** 上传文件之前的钩子 */
  beforeUpload?: (file: File) => boolean | Promise<boolean>
  /** 删除文件之前的钩子 */
  beforeRemove?: (file: UploadFile, fileList: UploadFile[]) => boolean | Promise<boolean>
  /** 提示文字 */
  tip?: string
}

export interface UploadEmits {
  /** 文件上传成功时的钩子 */
  (e: 'success', response: any, file: UploadFile, fileList: UploadFile[]): void
  /** 文件上传失败时的钩子 */
  (e: 'error', error: Error, file: UploadFile, fileList: UploadFile[]): void
  /** 文件上传时的钩子 */
  (e: 'progress', event: ProgressEvent, file: UploadFile, fileList: UploadFile[]): void
  /** 文件列表移除文件时的钩子 */
  (e: 'remove', file: UploadFile, fileList: UploadFile[]): void
  /** 当超出限制时触发的钩子 */
  (e: 'exceed', files: File[], fileList: UploadFile[]): void
}

