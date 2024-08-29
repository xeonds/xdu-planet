# xdu-planet

一个简单的RSS博客聚合站

使用`Golang+Vue3+ElementPlus`构建，借助GitHub Action实现自动更新Feed数据并生成页面，最终页面构建于`html`分支下，并由GitHub Pages呈现。

## Build

### Linux

按照下面的脚本，先构建前端，再构建后端。最终的产物就是`xdu-planet`。

运行之前先安装`nodejs, npm, pnpm, golang>=1.19`并配置好网络环境/镜像源。

```bash
git clone https://github.com/xdlinux/planet.git && cd planet
(cd frontend && pnpm i && pnpm run build)
go mod tidy && go build
```

### Windows

首先安装`golang`和`pnpm`。前者下个安装包就行，后者需要先安装npm，随后：

```bash
npm install -g pnpm
```

然后开始构建：

```bash
cd frontend
pnpm install
pnpm run build
cd ..
go build
```

最终得到`xdu-planet.exe`。

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

初次运行会产生一个空的配置文件`config.yml`，需要手动填写。配置文件格式如下：

```yaml
version: 1
feeds:
  - "https://xeonds.github.io/atom"
```

然后直接运行`xdu-planet`即可。打开浏览器访问`http://localhost:8192`即可看到聚合站。

进程会每隔15分钟更新一次Feed源，更新后的数据会存储在`db/`目录下。

## License

MIT License
