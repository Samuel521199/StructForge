# 创建前端目录结构的PowerShell脚本

$basePath = "src"

# 公共组件库目录
$commonComponents = @(
    "components/common/base/Badge",
    "components/common/base/Tag",
    "components/common/base/Tooltip",
    "components/common/base/Popover",
    "components/common/base/Dropdown",
    "components/common/base/Menu",
    "components/common/base/Tabs",
    "components/common/base/Pagination",
    "components/common/base/Loading",
    "components/common/base/Empty",
    "components/common/base/Message",
    "components/common/base/Notification",
    
    "components/common/data-display/DataTable",
    "components/common/data-display/DataList",
    "components/common/data-display/DataCard",
    "components/common/data-display/Statistic",
    "components/common/data-display/Progress",
    "components/common/data-display/Timeline",
    "components/common/data-display/Tree",
    
    "components/common/feedback/Alert",
    "components/common/feedback/Toast",
    "components/common/feedback/Confirm",
    "components/common/feedback/Drawer",
    "components/common/feedback/Modal",
    
    "components/common/form/FormField",
    "components/common/form/FormItem",
    "components/common/form/FormValidator",
    "components/common/form/DatePicker",
    "components/common/form/TimePicker",
    "components/common/form/Upload",
    "components/common/form/Editor",
    "components/common/form/CodeEditor",
    
    "components/common/navigation/Breadcrumb",
    "components/common/navigation/Steps",
    "components/common/navigation/Anchor",
    "components/common/navigation/BackTop",
    
    "components/common/layout/Container",
    "components/common/layout/Grid",
    "components/common/layout/Flex",
    "components/common/layout/Space",
    
    "components/common/business/FilterPanel",
    "components/common/business/ActionBar",
    "components/common/business/AvatarGroup",
    "components/common/business/TimeAgo",
    "components/common/business/CopyButton"
)

# 工作流组件目录
$workflowComponents = @(
    "components/workflow/editor/Canvas",
    "components/workflow/editor/NodePalette",
    "components/workflow/editor/PropertiesPanel",
    "components/workflow/editor/Toolbar",
    "components/workflow/editor/MiniMap",
    
    "components/workflow/nodes/trigger/WebhookNode",
    "components/workflow/nodes/trigger/TimerNode",
    "components/workflow/nodes/trigger/ManualNode",
    
    "components/workflow/nodes/ai/TextGenNode",
    "components/workflow/nodes/ai/ImageGenNode",
    "components/workflow/nodes/ai/CodeGenNode",
    
    "components/workflow/nodes/data/TransformNode",
    "components/workflow/nodes/data/FilterNode",
    "components/workflow/nodes/data/AggregateNode",
    
    "components/workflow/nodes/integration/HttpNode",
    "components/workflow/nodes/integration/DatabaseNode",
    "components/workflow/nodes/integration/FileNode",
    
    "components/workflow/nodes/control/ConditionNode",
    "components/workflow/nodes/control/LoopNode",
    "components/workflow/nodes/control/ParallelNode",
    
    "components/workflow/nodes/tool/CodeExecutorNode",
    "components/workflow/nodes/tool/ScriptNode",
    
    "components/workflow/execution/ExecutionMonitor",
    "components/workflow/execution/ExecutionLog",
    "components/workflow/execution/ExecutionStatus",
    "components/workflow/execution/ExecutionHistory",
    
    "components/workflow/utils/NodeConnector",
    "components/workflow/utils/NodeValidator",
    "components/workflow/utils/WorkflowExporter"
)

# 业务组件目录
$businessComponents = @(
    "components/business/user/UserAvatar",
    "components/business/user/UserProfile",
    "components/business/user/PermissionManager",
    
    "components/business/ai/ModelSelector",
    "components/business/ai/ModelConfig",
    "components/business/ai/ModelStatus",
    
    "components/business/system/SystemStatus",
    "components/business/system/SystemConfig"
)

# 布局组件目录
$layoutComponents = @(
    "components/layout/PageLayout",
    "components/layout/SectionLayout",
    "components/layout/GridLayout"
)

# 视图目录
$views = @(
    "views/ai/ModelList",
    "views/ai/ModelConfig",
    "views/user/UserProfile",
    "views/user/UserSettings",
    "views/system/SystemSettings",
    "views/system/SystemLogs"
)

# 创建所有目录
$allDirs = $commonComponents + $workflowComponents + $businessComponents + $layoutComponents + $views

foreach ($dir in $allDirs) {
    $fullPath = Join-Path $basePath $dir
    if (-not (Test-Path $fullPath)) {
        New-Item -ItemType Directory -Path $fullPath -Force | Out-Null
        Write-Host "Created: $fullPath" -ForegroundColor Green
    } else {
        Write-Host "Exists: $fullPath" -ForegroundColor Yellow
    }
}

Write-Host "`nDirectory structure creation completed!" -ForegroundColor Cyan

