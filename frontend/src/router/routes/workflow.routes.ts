import type { RouteRecordRaw } from 'vue-router'

/**
 * 工作流模块路由
 */
const workflowRoutes: RouteRecordRaw[] = [
  {
    path: '/workflow',
    name: 'Workflow',
    redirect: '/workflow/list',
    meta: {
      title: '工作流',
      icon: 'workflow',
      requiresAuth: true
    },
    children: [
      {
        path: 'list',
        name: 'WorkflowList',
        component: () => import('@/views/workflow/WorkflowList.vue'),
        meta: {
          title: '工作流列表',
          requiresAuth: true
        }
      },
      {
        path: 'editor/:id?',
        name: 'WorkflowEditor',
        component: () => import('@/views/workflow/WorkflowEditor.vue'),
        meta: {
          title: '工作流编辑器',
          requiresAuth: true
        }
      },
      {
        path: 'detail/:id',
        name: 'WorkflowDetail',
        component: () => import('@/views/workflow/WorkflowDetail.vue'),
        meta: {
          title: '工作流详情',
          requiresAuth: true
        }
      },
      {
        path: 'execution/:id',
        name: 'WorkflowExecution',
        component: () => import('@/views/workflow/WorkflowExecution.vue'),
        meta: {
          title: '执行详情',
          requiresAuth: true
        }
      }
    ]
  }
]

export default workflowRoutes

