# xdu-planet

一个简单的RSS博客聚合站

## Build

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

然后直接运行`xdu-planet.exe`即可。打开浏览器访问`http://localhost:8192`即可看到聚合站。

需要更新聚合内容时，无需退出，执行`xdu-planet.exe -fetch`即可。

## License

MIT License
