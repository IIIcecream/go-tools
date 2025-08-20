## password generator

### 概述

本工具可生成并持久化保存随机密码

### 编译

执行 `make release` 即可，在当前目录会生成 `password_gen-v1.0.0.tar.gz`

压缩包内容即 `output` 文件夹中三个架构下的可执行文件

### 使用

``` bash
./password_gen-darwin-arm64 --help

# 创建符合要求的密码
# 各参数说明见 `./password_gen-darwin-arm64 create --help`
./password_gen-darwin-arm64 create --key=test --length=12 --min=3 -l -n -s -u --save=true

# 查询 key 为 test 的持久化保存的密码
./password_gen-darwin-arm64 find --key=test

# 查询所有保存的密码
./password_gen-darwin-arm64 list
```