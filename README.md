## 项目

[![Build Status](https://drone.memosa.cn/api/badges/aaronzjc/mu/status.svg)](https://drone.memosa.cn/aaronzjc/mu)

### 介绍

热榜聚合这个产品已经有很多很多了。我只是按照自己的想法写了个自己用的而已。

主要是为了学习Go开发，目前使用到的技术栈

+ Go, Vue.js
+ Goland
+ Drone CI, Portainer, Grafana + Loki 
+ Docker, Docker Swarm

### 预览
[这里](https://github.com/aaronzjc/mu/tree/master/doc)

### 开发

```shell
# 1. 本地安装MySQL和Redis
# 2. 修改配置`app.json`和`.env.development`
# 3. 启动各个服务

$ make crawler-dev # 抓取节点
$ make mu-dev # 前端页面和任务调度
$ make vue-dev # 本地开发H5
```

[开发说明](doc/DEV.md)

### LICENSE

本项目采用[MIT](https://github.com/aaronzjc/mu/blob/dev/LICENSE)开源授权许可证。