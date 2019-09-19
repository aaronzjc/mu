## 项目

[![Build Status](http://drone.memosa.cn/api/badges/aaronzjc/crawler/status.svg)](http://drone.memosa.cn/aaronzjc/crawler)

### 介绍

热榜聚合这个产品已经有很多很多了。我只是按照自己的想法写了个自己用的而已。

主要是为了学习Go开发。 

目前使用到的技术栈

+ Go
+ Goland
+ Vue.js
+ Drone CI
+ Docker

### 安装

```shell
make crawler # 编译cron，抓取数据
make mu # 编译前端展示
```

### 目录介绍

非标准目录结构，尽可能按照官方推荐的方式。

```text
|- cmd              程序入口
    |- master       web程序
    |- node         后台程序
|- deploy           部署脚本
|- lib              私有库
|- public           编译后的前端页面
|- util             工具操作库，例如数据库，缓存等
|- web              前端代码
```