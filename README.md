![build & release](https://github.com/aaronzjc/mu/workflows/build%20&%20release/badge.svg)

## 项目

### 介绍

热榜聚合这个产品已经有很多很多了。我只是按照自己的想法写了个自己用的而已。

主要是为了学习Go开发，目前使用到的技术栈

+ Go, Vue.js，MySQL, Redis
+ Github Actions，Portainer 
+ Docker，Docker Swarm

### 预览
[这里](https://github.com/aaronzjc/mu/tree/master/doc)

### 本地运行

```shell
# 1. 安装MySQL和Redis
# 2. 执行deploy/db.sql导入数据库结构
# 3. 修改配置`conf/app.json`里面的各项配置
# 4. 容器启动参考deploy/stack.yml文件
```

### 本地开发

```shell
make api-dev # 启动接口
make commander-dev # 启动调度器
make agent-dev # 启动任务执行节点
```

[开发说明](doc/DEV.md)

### LICENSE

本项目采用[MIT](https://github.com/aaronzjc/mu/blob/dev/LICENSE)开源授权许可证。
