/**
 * 请求封装组合函数
 */

import { ref } from 'vue'
import type { AxiosError } from 'axios'

export function useRequest<T>(
  requestFn: () => Promise<T>
) {
  const loading = ref(false)
  const error = ref<Error | null>(null)
  const data = ref<T | null>(null)

  const execute = async () => {
    loading.value = true
    error.value = null
    try {
      const result = await requestFn()
      data.value = result
      return result
    } catch (err) {
      error.value = err as Error
      throw err
    } finally {
      loading.value = false
    }
  }

  return {
    loading,
    error,
    data,
    execute,
  }
}

