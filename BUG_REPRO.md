# Bug 复现

## Bug 是什么
已取消的 CSV 导出仍留下成功审计。

## 如何触发
运行：`go test ./tests -run TestCanceledExportLeavesNoAudit -count=1`。

## 错误信息
测试报告 canceled export wrote 1 audits。
