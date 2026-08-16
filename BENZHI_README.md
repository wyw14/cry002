# 本地评测镜像

```bash
./build_benzhi_docker.sh cry002 linux/amd64
docker run --rm cry002 go test ./...
```

`golang:1.24` 官方镜像同时提供 amd64 与 arm64 版本；本仓库只要求实际验证当前机器平台。
