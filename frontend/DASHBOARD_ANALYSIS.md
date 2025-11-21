# 仪表盘功能分析与设计文档

## 📊 功能模块分析

基于参考图片和 StructForge 工作流平台特性，仪表盘应包含以下核心功能模块：

### 1. **实时执行状态监控** (Top Priority)

#### 1.1 当前执行状态卡片
- **功能**：显示当前正在运行和排队的工作流数量
- **数据指标**：
  - `Running`: 正在执行的工作流数量（实时更新）
  - `Queued`: 排队等待执行的工作流数量
  - `Failed`: 最近失败的工作流数量（可选）
- **显示方式**：
  - 大号数字显示（绿色高亮）
  - 实时刷新（WebSocket 或轮询）
  - 点击可跳转到执行列表
- **UI组件**：
  - `Card` 组件
  - `Icon` 组件（运行图标）
  - 数字动画效果

#### 1.2 执行历史列表
- **功能**：显示最近 200 条工作流执行记录
- **数据字段**：
  - 状态图标（成功/失败/运行中）
  - 工作流名称
  - 开始时间
  - 执行状态（成功/失败/运行中）
  - 执行时长
  - 执行 ID
  - 执行数据（JSON 格式，可展开）
- **显示方式**：
  - `Table` 组件
  - 分页显示（每页 8-20 条）
  - 状态标签（`StatusTag` 组件）
  - 可点击跳转到执行详情
- **交互功能**：
  - 筛选（按状态、工作流名称）
  - 排序（按时间、状态）
  - 搜索
  - 实时更新（运行中的任务显示动画）

### 2. **执行统计与分析** (High Priority)

#### 2.1 执行成功率统计（最近 24 小时）
- **功能**：按工作流统计执行成功率
- **数据指标**：
  - 工作流名称
  - 总执行次数
  - 成功次数
  - 失败次数
  - 成功率百分比
- **显示方式**：
  - `Table` 组件
  - 成功率用进度条或百分比显示
  - 失败次数用红色高亮
  - 成功次数用绿色显示
- **UI组件**：
  - `Table` 组件
  - `Progress` 组件（进度条）
  - `Tag` 组件（状态标签）

#### 2.2 工作流错误趋势图
- **功能**：显示最近 24 小时内的错误趋势
- **数据指标**：
  - X 轴：时间（小时）
  - Y 轴：错误数量
  - 数据点：每小时错误数量
- **显示方式**：
  - 折线图（Line Chart）
  - 绿色线条
  - 时间范围选择器（24小时/7天/30天）
- **UI组件**：
  - 图表库（ECharts 或 Chart.js）
  - `Card` 组件作为容器

#### 2.3 平均执行时长统计
- **功能**：按工作流显示平均、最小、最大执行时长
- **数据指标**：
  - 工作流名称
  - 平均执行时长
  - 最小执行时长
  - 最大执行时长
- **显示方式**：
  - 水平条形图（Horizontal Bar Chart）
  - 不同颜色表示不同指标（平均/最小/最大）
  - 时间单位自动转换（ms/s/min）
- **UI组件**：
  - 图表库（ECharts）
  - `Card` 组件
  - 图例（Legend）

### 3. **Webhook 活动监控** (Medium Priority)

#### 3.1 最近 Webhook 活动
- **功能**：显示最近的 Webhook 触发记录
- **数据字段**：
  - 工作流名称
  - Webhook 路径
  - HTTP 方法（GET/POST）
  - 执行时间
  - Webhook ID/状态
- **显示方式**：
  - `Table` 组件
  - 时间格式化显示
  - HTTP 方法标签（颜色区分）
- **UI组件**：
  - `Table` 组件
  - `Tag` 组件（HTTP 方法标签）

#### 3.2 Webhook 使用趋势
- **功能**：显示 Webhook 调用频率趋势
- **数据指标**：
  - X 轴：时间
  - Y 轴：调用次数
- **显示方式**：
  - 折线图或柱状图
  - 绿色主题
- **UI组件**：
  - 图表库
  - `Card` 组件

### 4. **快速操作与概览** (Medium Priority)

#### 4.1 工作流快速操作
- **功能**：快速创建工作流、查看工作流列表
- **显示方式**：
  - 大按钮卡片
  - "创建新工作流" 按钮
  - "查看所有工作流" 链接
- **UI组件**：
  - `Button` 组件
  - `Card` 组件

#### 4.2 最近工作流
- **功能**：显示最近创建或修改的工作流
- **数据字段**：
  - 工作流名称
  - 状态
  - 最后修改时间
- **显示方式**：
  - 列表或卡片
  - 点击跳转到工作流详情
- **UI组件**：
  - `Card` 组件
  - `Link` 组件

### 5. **系统指标** (Low Priority)

#### 5.1 系统性能指标
- **功能**：显示系统整体性能
- **数据指标**：
  - 总工作流数量
  - 活跃工作流数量
  - 总执行次数（最近 7 天）
  - 失败执行次数（最近 7 天）
  - 失败率
  - 平均执行时间
  - 节省时间（如果可计算）
- **显示方式**：
  - 指标卡片网格
  - 数字 + 标签
  - 趋势箭头（上升/下降）
- **UI组件**：
  - `Card` 组件
  - `Icon` 组件（趋势图标）

## 🎨 UI 设计要点

### 1. **布局结构**

```
┌─────────────────────────────────────────────────┐
│  Header (已实现)                                │
├─────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐             │
│  │ 执行状态卡片 │  │ 执行状态卡片 │  ...        │
│  └─────────────┘  └─────────────┘             │
├─────────────────────────────────────────────────┤
│  ┌──────────────────┐  ┌──────────────────┐  │
│  │ 执行历史列表      │  │ 执行成功率统计    │  │
│  │ (Table + 分页)   │  │ (Table)          │  │
│  └──────────────────┘  └──────────────────┘  │
├─────────────────────────────────────────────────┤
│  ┌──────────────────┐  ┌──────────────────┐  │
│  │ 平均执行时长      │  │ 错误趋势图        │  │
│  │ (水平条形图)      │  │ (折线图)         │  │
│  └──────────────────┘  └──────────────────┘  │
├─────────────────────────────────────────────────┤
│  ┌──────────────────┐  ┌──────────────────┐  │
│  │ 最近 Webhook     │  │ Webhook 趋势      │  │
│  │ (Table)          │  │ (折线图)         │  │
│  └──────────────────┘  └──────────────────┘  │
└─────────────────────────────────────────────────┘
```

### 2. **颜色主题**

- **主色调**：绿色（#00FF00）- 与登录页代码雨主题一致
- **成功状态**：绿色（#67c23a）
- **失败状态**：红色（#f56c6c）
- **警告状态**：黄色（#e6a23c）
- **信息状态**：蓝色（#409eff）
- **背景**：深色主题（#1a1a1a 或 #000000）
- **卡片背景**：半透明黑色（rgba(0, 0, 0, 0.85)）
- **边框**：绿色发光效果（rgba(0, 255, 0, 0.3)）

### 3. **视觉效果**

- **卡片阴影**：绿色光晕效果
  ```scss
  box-shadow: 
    0 0 40px rgba(0, 255, 0, 0.5),
    inset 0 0 10px rgba(0, 255, 0, 0.2);
  ```
- **数字动画**：数字变化时的平滑过渡
- **加载状态**：骨架屏（Skeleton）或加载动画
- **实时更新**：WebSocket 连接，数据实时刷新
- **响应式设计**：适配不同屏幕尺寸

### 4. **交互体验**

- **实时更新**：关键指标自动刷新（每 5-10 秒）
- **点击跳转**：卡片、表格行可点击跳转到详情页
- **筛选排序**：表格支持多条件筛选和排序
- **时间选择**：图表支持时间范围选择（24小时/7天/30天）
- **数据导出**：支持导出统计数据（CSV/JSON）
- **空状态**：无数据时显示友好的空状态提示

## 🧩 需要的 UI 组件

### 已实现的组件
- ✅ `Card` - 卡片容器
- ✅ `Table` - 表格
- ✅ `Button` - 按钮
- ✅ `Icon` - 图标
- ✅ `Tag` - 标签（需要检查）
- ✅ `Pagination` - 分页
- ✅ `Loading` - 加载状态
- ✅ `Empty` - 空状态
- ✅ `Message` - 消息提示

### 需要实现的组件
- ❌ `Progress` - 进度条（用于成功率显示）
- ❌ `Badge` - 徽章（用于数字显示）
- ❌ `Statistic` - 统计数字（大号数字显示）
- ❌ `Chart` - 图表容器（封装 ECharts）
- ❌ `Skeleton` - 骨架屏（加载占位）
- ❌ `TimeAgo` - 相对时间显示（"2 小时前"）

### 需要集成的第三方库
- **ECharts** 或 **Chart.js** - 图表库
  - 推荐：ECharts（功能强大，中文文档完善）
  - 安装：`npm install echarts vue-echarts`

## 📐 布局实现建议

### 1. 响应式网格布局

使用 CSS Grid 或 Flexbox 实现响应式布局：

```scss
.dashboard-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 24px;
  
  // 大屏幕：2列
  @media (min-width: 1200px) {
    grid-template-columns: repeat(2, 1fr);
  }
  
  // 超大屏幕：3列
  @media (min-width: 1600px) {
    grid-template-columns: repeat(3, 1fr);
  }
}
```

### 2. 卡片组件增强

为仪表盘创建专门的统计卡片组件：

```vue
<StatCard
  title="正在运行"
  value="2"
  icon="running"
  color="green"
  :trend="{ value: 10, direction: 'up' }"
/>
```

### 3. 图表组件封装

封装 ECharts 为通用组件：

```vue
<Chart
  type="line"
  :data="errorTrendData"
  :options="chartOptions"
  height="300px"
/>
```

## 🔄 数据获取策略

### 1. API 端点设计

```typescript
// 仪表盘统计数据
GET /api/v1/dashboard/stats
Response: {
  running: number
  queued: number
  totalExecutions: number
  failedExecutions: number
  failureRate: number
  avgExecutionTime: number
}

// 执行历史
GET /api/v1/dashboard/executions?page=1&pageSize=20
Response: {
  list: Execution[]
  total: number
}

// 执行成功率
GET /api/v1/dashboard/success-rate?period=24h
Response: {
  workflows: Array<{
    name: string
    total: number
    successful: number
    failed: number
    successRate: number
  }>
}

// 错误趋势
GET /api/v1/dashboard/error-trend?period=24h
Response: {
  data: Array<{
    time: string
    count: number
  }>
}

// 平均执行时长
GET /api/v1/dashboard/execution-duration
Response: {
  workflows: Array<{
    name: string
    avgDuration: number
    minDuration: number
    maxDuration: number
  }>
}

// Webhook 活动
GET /api/v1/dashboard/webhook-activity?limit=10
Response: {
  list: WebhookActivity[]
}
```

### 2. 实时更新策略

- **WebSocket**：用于实时执行状态更新
- **轮询**：用于统计数据更新（每 30 秒）
- **手动刷新**：提供刷新按钮

## 📝 实现优先级

### Phase 1: 核心功能（必须）
1. ✅ 执行状态卡片（Running/Queued）
2. ✅ 执行历史列表（Table + 分页）
3. ✅ 执行成功率统计（Table）
4. ✅ 基础布局和样式

### Phase 2: 数据可视化（重要）
5. ⏳ 错误趋势图（折线图）
6. ⏳ 平均执行时长图（水平条形图）
7. ⏳ Webhook 使用趋势图

### Phase 3: 增强功能（可选）
8. ⏳ Webhook 活动列表
9. ⏳ 系统性能指标
10. ⏳ 快速操作卡片
11. ⏳ 数据导出功能

## 🎯 技术实现要点

1. **状态管理**：使用 Pinia 管理仪表盘数据
2. **数据缓存**：使用缓存减少 API 调用
3. **错误处理**：统一的错误处理和提示
4. **加载状态**：骨架屏提升用户体验
5. **性能优化**：虚拟滚动（如果列表很长）
6. **可访问性**：ARIA 标签和键盘导航支持

## 📚 参考资源

- ECharts 官方文档：https://echarts.apache.org/
- Vue ECharts：https://github.com/ecomfe/vue-echarts
- Element Plus 统计组件参考
- n8n Dashboard 设计参考

