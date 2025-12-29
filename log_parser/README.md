## log parser

### 概述

解析从日志平台下载的 .zip 日志压缩包
1. 支持单次解析 （时间戳 -> UTC+8 时间）
2. 支持监听固定文件夹，当有新的 .zip 添加进入时，进行解析
3. 支持日志合并
4. 支持配置模块，默认已添加 APP_PROXY,EVENT_TASK,PERSIST

### 使用

可使用以下命令查看，有详尽的使用方法介绍
```
./log_parser -h
./log_parser listen -h
```
