# go-static-blog
Go 语言实现的静态博客生成器

## 使用

```shell
$ git clone https://github.com/imyhui/go-static-blog
$ cd go-static-blog
$ go build
```

```shell
$ ./go-static-blog -h
go-static-blog version: 1.0.0
Usage: go-static-blog [-g generate] [-s server] 

Options:
  -g	clean and generate
  -s	server on 8080
```

```golang
$ ./go-static-blog -g -s
```

访问 [http://localhost:8080](http://localhost:8080/)

## 界面
主页
![](http://ghost.andyhui.top/主页.png)
文章页面
![](http://ghost.andyhui.top/文章页面.png)
标签云
![](http://ghost.andyhui.top/标签云.png)



详细实现过程见[这里](http://andyhui.top/go-static-blog/)

## 后续优化

- - [x] ~~优化解析文章逻辑~~
  
  - - [x] ~~添加标签~~
- - [x] ~~优化页面，添加样式~~
- 添加页面
  - 归档页
  - - [x] ~~标签云~~
  - 关于页
- 添加命令行工具，如生成文章，部署等
- 添加持续集成服务，如[Travis CI](http://andyhui.top/go-static-blog/[https://travis-ci.org](https://travis-ci.org/)) 或者 [GitHub Actions](https://github.com/features/actions)