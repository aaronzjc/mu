### 目录结构

非严格标准目录结构，尽可能按照官方推荐的方式。

```text
├── Dockerfile      # 构建镜像
├── Makefile        # 构建命令
├── app.json        # 项目配置文件，数据库连接等
├── cmd             # 应用入口
│   ├── commander
│   └── agent
|   └── api
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

几个说明

```text
1. cmd为应用入口，保持逻辑尽可能简单  
2. app目录为应用实体，做初始化相关，例如，初始化配置，服务器等。
3. model是数据库访问层，除了util之外，只能被其他的包调用。
4. route是项目路由组件。除了自身和app包，可以调用internal下其他所有包
5. svc是封装的一些业务逻辑。可以调用model和util包  
6. util包只能调用自身。
```