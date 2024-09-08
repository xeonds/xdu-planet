# xdu-planet

一个简单的RSS博客聚合站

使用`Golang+Vue3+ElementPlus`构建，借助GitHub Action实现自动更新Feed数据并生成页面，最终页面构建于`html`分支下，并由GitHub Pages呈现。

## Build

确保已经安装了golang sdk和nodejs后，克隆仓库。
```bash
cd xdu-planet
make frontend && make build/xdu-planet
```

最终得到的二进制文件在`build/`目录下。

## API 文档

- `GET /api/v1/feed` 获取所有文章
- `GET /api/v1/comment/:article_id` 获取某篇文章的评论
- `GET /api/v1/comment/reply_to/:comment_id` 获取某个评论的回复
- `POST /api/v1/comment/:article_id` 发表评论
  - `content` 评论内容
  - `user_id` 用户ID
  - `reply_to` 回复的评论ID（可选，仅当回复评论时传值）
- `DELETE /api/v1/comment/:comment_id` 举报评论
  - 举报后，评论自动转为`audit`状态，管理员可在后台审核
- `GET /api/v1/admin/comment/:filter` 获取评论列表
  - `filter` 可选，可选值为`ok`、`block`、`delete`、`audit`，分别表示已通过、已屏蔽、已删除、待审核
- `POST /api/v1/admin/comment/audit/:comment_id` 审核评论
  - `status` 审核状态，可选值为`ok`、`block`，分别表示通过、屏蔽

## 用法

### 作为服务部署

初次运行会产生一个空的配置文件`config.yml`，需要手动填写。参考配置如下：

```yaml
version: 1
databaseconfig:   # 数据库配置，现阶段用于支持评论系统
  type: "sqlite"  # 可选值为sqlite、mysql
  host: ""
  port: ""
  user: ""
  password: ""
  db: "planet.db"
  migrate: true
avalonguard:              # 评论审核系统配置
  enablegravetimer: true  # 是否启用：被举报评论超时自动转为屏蔽状态
  gravetimeout: 3600s     # 评论审核超时时间
  enablefilter: false     # 是否启用：评论内容过滤
  filter: []              # 过滤关键词列表
logfile: "admin.log"      # 管理员操作日志文件
admintoken: []            # 管理员token列表
feeds: []                 # RSS源列表
```

然后直接运行`xdu-planet`即可。打开浏览器访问`http://localhost:8192`即可看到聚合站。

进程会每隔15分钟更新一次Feed源，更新后的数据会存储在`db/`目录下。

### 作为命令行程序

更改完配置文件后，通过如下命令即可更新Feed数据：

```bash
./xdu-planet -fetch
```

抓取完毕后，会生成作为索引的`index.json`和包含文章正文的`db.json`，以及作为`index.json`索引指向的正文数据库`db/`。

### 作为静态站点部署

如果不需要评论系统，可以直接将生成的`build/`目录下的文件（除了数据库文件，可执行程序和配置文件）部署到静态服务器上。在需要更新时，只需要重新运行`xdu-planet -fetch`并将生成的文件部署到服务器上即可。

## License

MIT License
