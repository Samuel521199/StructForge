<template>
  <el-upload
    :action="action"
    :headers="headers"
    :data="data"
    :multiple="multiple"
    :accept="accept"
    :limit="limit"
    :file-list="fileList"
    :auto-upload="autoUpload"
    :drag="drag"
    :disabled="disabled"
    :show-file-list="showFileList"
    :list-type="listType"
    :on-success="handleSuccess"
    :on-error="handleError"
    :on-progress="handleProgress"
    :on-remove="handleRemove"
    :before-upload="beforeUpload"
    :before-remove="beforeRemove"
    :on-exceed="handleExceed"
    class="structforge-upload"
  >
    <slot>
      <Button v-if="!drag">
        点击上传
      </Button>
      <div v-else class="upload-drag-area">
        <div class="upload-text">将文件拖到此处，或<em>点击上传</em></div>
      </div>
    </slot>
    <template #tip v-if="tip">
      <div class="upload-tip">{{ tip }}</div>
    </template>
  </el-upload>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { Button } from '../Button'
import type { UploadProps, UploadEmits, UploadFile } from './types'

const props = withDefaults(defineProps<UploadProps>(), {
  multiple: false,
  autoUpload: true,
  drag: false,
  disabled: false,
  showFileList: true,
  listType: 'text',
  limit: 0,
})

const emit = defineEmits<UploadEmits>()

const fileList = ref<UploadFile[]>([])

const handleSuccess = (response: any, file: UploadFile, fileList: UploadFile[]) => {
  emit('success', response, file, fileList)
}

const handleError = (error: Error, file: UploadFile, fileList: UploadFile[]) => {
  emit('error', error, file, fileList)
}

const handleProgress = (event: ProgressEvent, file: UploadFile, fileList: UploadFile[]) => {
  emit('progress', event, file, fileList)
}

const handleRemove = (file: UploadFile, fileList: UploadFile[]) => {
  emit('remove', file, fileList)
}

const beforeUpload = (file: File) => {
  const result = props.beforeUpload?.(file)
  if (result === false) {
    return false
  }
  return true
}

const beforeRemove = (file: UploadFile, fileList: UploadFile[]) => {
  const result = props.beforeRemove?.(file, fileList)
  if (result === false) {
    return false
  }
  return true
}

const handleExceed = (files: File[], fileList: UploadFile[]) => {
  emit('exceed', files, fileList)
}

defineOptions({
  name: 'Upload',
})
</script>

<style scoped lang="scss">
.structforge-upload {
  .upload-drag-area {
    padding: 40px;
    text-align: center;
    border: 1px dashed var(--el-border-color);
    border-radius: 6px;
    cursor: pointer;
    transition: all 0.3s;

    &:hover {
      border-color: var(--el-color-primary);
    }

    .upload-text {
      margin-top: 16px;
      color: var(--el-text-color-regular);
      font-size: 14px;

      em {
        color: var(--el-color-primary);
        font-style: normal;
      }
    }
  }

  .upload-tip {
    margin-top: 8px;
    color: var(--el-text-color-secondary);
    font-size: 12px;
  }
}
</style>

