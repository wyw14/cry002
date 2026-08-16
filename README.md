# 案卷全生命周期管理系统

这是一个从 0 到 1 实现的 Go + Vue 3 全栈项目，覆盖建卷、材料目录、送审、归档、版本修订、借阅审批、借出、续借、归还、分类树、附件和审计闭环。后端是唯一业务服务，所有授权和状态校验均在服务端完成。

## 技术栈与目录

- Go 1.24、Gin、GORM、validator、Zap、Viper、bcrypt、JWT
- PostgreSQL 17；Vue 3、TypeScript、Vite、Pinia、Element Plus
- `cmd/server`：装配和优雅停机；`internal/domain`：状态机；`internal/application`：用例编排
- `internal/repository`：内存测试适配器与 PostgreSQL 事务实现；`internal/platform/storage`：受控附件目录
- `internal/transport/http`：`/api/v1` 路由；`migrations`：版本化 SQL；`api/openapi`：OpenAPI 3.0
- `tests`：并发借出、归档版本、附件完整性、取消导出和分类事务测试；`web`：前端

## 角色与权限

| 能力 | 管理员 | 档案员 | 审核员 | 部门负责人 | 借阅人 |
|---|---:|---:|---:|---:|---:|
| 用户、分类与统计 | 是 | 分类/统计 | 否 | 否 | 否 |
| 建卷与提交 | 是 | 是 | 审核只读 | 本部门 | 否 |
| 审核归档 | 是 | 是 | 是 | 否 | 否 |
| 归档版本修订 | 是 | 是 | 否 | 否 | 否 |
| 借阅申请 | 是 | 是 | 是 | 是 | 是 |
| 审批、借出、归还 | 是 | 是 | 否 | 否 | 本人归还 |
| 审计时间线 | 是 | 是 | 是 | 本部门 | 本人相关 |

前端按钮不构成授权；API 会重新加载用户和资源并校验角色。附件下载跟随案卷权限，文件名不能包含路径，物理路径必须保持在附件根目录中。

## 状态规则

- 案卷：`draft -> pending_review -> archived -> borrowed -> archived -> sealed -> destroyed`。审核退回进入 `draft`；销毁和恢复要求管理员与原因。
- 已归档或封存案卷的核心元数据不能原地覆盖。修订会递增版本，保存旧快照、差异字段、原因和操作人。
- 借阅：`pending -> approved/rejected -> checked_out -> returned`，逾期任务将 `checked_out` 标记为 `overdue`。同一案卷的借出使用数据库事务、行锁和部分唯一索引保证单一赢家。
- 分类不能移动到自身或后代；移动和同级排序在一个事务中提交，失败时树保持不变。

## 启动

```bash
cp .env.example .env
docker compose up --build
```

API 为 `http://localhost:8080`，前端为 `http://localhost:8081`；健康检查为 `/healthz` 和 `/readyz`。也可设置 `ALLOW_IN_MEMORY=true` 后执行 `go run ./cmd/server` 进行快速演示。

演示账号统一密码为 `Demo123456!`：

- `admin@archive.local` 管理员
- `archivist@archive.local` 档案员
- `auditor@archive.local` 审核员
- `manager@archive.local` 部门负责人
- `borrower@archive.local` 借阅人

初始化可重复执行：固定邮箱、部门代码、分类代码和案卷号受唯一约束保护，已存在数据不会被覆盖。

## API 示例

```bash
curl http://localhost:8080/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"archivist@archive.local","password":"Demo123456!"}'
```

列表统一支持 `page`、`page_size`、筛选和排序。错误响应始终包含稳定 `code`、可读 `message`、`fields` 和 `request_id`。完整契约见 `api/openapi/openapi.yaml`。

## 验证

```bash
go build ./...
go test ./...
go test -race ./...
go vet ./...
npm --prefix web test
npm --prefix web run build
```

关键测试名称：`TestConcurrentCheckoutSingleWinner`、`TestArchivedUpdateCreatesVersionAndPreservesSnapshot`、`TestAttachmentStorageRejectsTraversalAndPreservesHash`、`TestCanceledExportLeavesNoAudit`、`TestClassificationMoveFailurePreservesTree`。

日志为结构化 Zap 日志并携带 `request_id`，中间件提供超时、panic 恢复、CORS 和安全头。不得记录密码、令牌或附件正文；仓库不提交 `.env`、缓存、构建产物、依赖目录或运行期数据。
