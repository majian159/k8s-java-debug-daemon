# Java高CPU自动抓取堆栈信息

## 构建
### 二进制
```sh
# 为当前系统平台构建
make

# 指定目标系统, GOOS: linux darwin window freebsd
make GOOS=linux

```

### Docker镜像
```sh
make docker

# 自定义镜像tag
make docker IMAGE=test
```