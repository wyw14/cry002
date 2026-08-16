# Bug 复现

## Bug 是什么
分类移动返回失败后仍留下已改变的父子关系。

## 如何触发
运行：`go test ./tests -run TestClassificationMoveFailurePreservesTree -count=1`。

## 错误信息
测试报告 partial update left parent=empty。
