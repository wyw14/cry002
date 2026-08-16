# Bug 复现

## Bug 是什么
两个已批准借阅申请并发借出同一案卷时都可能成功。

## 如何触发
运行：`go test ./tests -run TestConcurrentCheckoutSingleWinner -count=20`。

## 错误信息
测试报告 successful checkouts=2, want 1。
