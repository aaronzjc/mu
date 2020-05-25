## 项目

![GithubCI](https://github.com/aaronzjc/mu/workflows/build%20&%20release/badge.svg?branch=master)
[![DroneCI](https://drone.memosa.cn/api/badges/aaronzjc/mu/status.svg)](https://drone.memosa.cn/aaronzjc/mu)

### 介绍

热榜聚合这个产品已经有很多很多了。我只是按照自己的想法写了个自己用的而已。

主要是为了学习Go开发，目前使用到的技术栈

+ Go, Vue.js
+ Goland
+ Drone CI, Portainer 
+ Docker, Docker Swarm

### 预览
[这里](https://github.com/aaronzjc/mu/tree/master/doc)

### 本地运行

```shell
# 1. 安装MySQL和Redis
# 2. 执行deploy/db.sql导入数据库结构
# 3. 修改配置`conf/app.json`里面的各项配置
# 4. 本地添加host: 127.0.0.1 mu.memosa.local。如果熟悉了可以改成其他
# 5. 修改`deploy/compose.yml`配置文件目录为你的`app.json`目录

$docker-compose -f deploy/compose.yml up
```

### 本地开发

参考`Makefile`。

[开发说明](doc/DEV.md)

### LICENSE

本项目采用[MIT](https://github.com/aaronzjc/mu/blob/dev/LICENSE)开源授权许可证。