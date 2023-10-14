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
