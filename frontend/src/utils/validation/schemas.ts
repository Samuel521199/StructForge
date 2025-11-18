/**
 * 验证模式定义
 */

import { z } from 'zod'

export const loginSchema = z.object({
  username: z.string().min(1, '用户名不能为空'),
  password: z.string().min(6, '密码至少6位'),
})

export const workflowSchema = z.object({
  name: z.string().min(1, '工作流名称不能为空'),
  description: z.string().optional(),
})

