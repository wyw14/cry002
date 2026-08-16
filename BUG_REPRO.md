# Bug 复现

## Bug 是什么
归档案卷元数据被原地覆盖，没有递增版本，也没有保存旧快照和差异。

## 如何触发
运行：`go test ./tests -run TestArchivedUpdateCreatesVersionAndPreservesSnapshot -count=20`。

## 错误信息
测试报告 updated case 版本错误或 version did not preserve old snapshot。
