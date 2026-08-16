# Bug 复现

## Bug 是什么
附件存储接受路径穿越键并丢失 SHA-256 完整性元数据，权限链也被绕过。

## 如何触发
运行：`go test ./tests -run TestAttachmentStorageRejectsTraversalAndPreservesHash -count=20`。

## 错误信息
测试报告 path traversal accepted 或 hash 与期望不一致。
