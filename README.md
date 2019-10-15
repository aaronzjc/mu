## 项目

[![Build Status](http://drone.memosa.cn/api/badges/aaronzjc/mu/status.svg)](http://drone.memosa.cn/aaronzjc/crawler)

### 介绍

热榜聚合这个产品已经有很多很多了。我只是按照自己的想法写了个自己用的而已。

主要是为了学习Go开发，目前使用到的技术栈

+ Go, Vue.js
+ Goland
+ Drone CI, Portainer, Grafana + Loki 
+ Docker, Docker Swarm

### 预览
[看这里](https://github.com/aaronzjc/mu/tree/master/doc)

### 本地开发

```shell
# 本地安装MySQL和Redis

# 修改配置`app.json`和`.env.development`

# 启动各个服务
$ make crawler-dev # 抓取节点
$ make mu-dev # 前端页面和任务调度
$ make vue-dev # 本地开发H5
```

### 项目目录和一些开发约定

非严格标准目录结构，尽可能按照官方推荐的方式。

```text
├── Dockerfile      # 构建镜像
├── Makefile        # 构建命令
├── app.json        # 项目配置文件，数据库连接等
├── cmd             # 应用入口
│   ├── master
│   └── node
├── deploy          # 部署相关，SQL脚本，docker-compose文件
├── internal        # 项目核心业务逻辑和库，工具
│   ├── app         # 应用实体，被入口文件调用
│   ├── model       # 数据库访问操作层
│   ├── route       # 接口路由
│   ├── svc         # 封装的业务执行
│   └── util        # 通用工具库。数据库访问，日志，配置文件，缓存，oauth等。
├── public          # 最终构建的前端页面
└── web             # 前端项目，基于Vue。
```

项目开发中的几个约定

```text
1. cmd为应用入口，保持逻辑尽可能简单  
2. app目录为应用实体，做初始化相关，例如，初始化配置，服务器等。
3. model是数据库访问层，除了util之外，只能被其他的包调用。
4. route是项目路由组件。除了自身和app包，可以调用internal下其他所有包
5. svc是封装的一些业务逻辑。可以调用model和util包  
6. util包可以被任何其他的包调用。但是不能调用除自身以外的包。
```

### LICENSE

本项目采用[MIT](https://github.com/aaronzjc/crawler/blob/dev/LICENSE)开源授权许可证。