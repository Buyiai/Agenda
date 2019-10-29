# Agenda
@[toc]
# 1、概述
命令行实用程序并不是都象 cat、more、grep 是简单命令。go 项目管理程序，类似 java 项目管理 maven、Nodejs 项目管理程序 npm、git 命令行客户端、 docker 与 kubernetes 容器管理工具等等都是采用了较复杂的命令行。即一个实用程序同时支持多个子命令，每个子命令有各自独立的参数，命令之间可能存在共享的代码或逻辑，同时随着产品的发展，这些命令可能发生功能变化、添加新命令等。因此，符合 OCP 原则 的设计是至关重要的编程需求。

**任务目标**
1. 熟悉 go 命令行工具管理项目
2. 综合使用 go 的函数、数据结构与接口，编写一个简单命令行应用 agenda
3. 使用面向对象的思想设计程序，使得程序具有良好的结构命令，并能方便修改、扩展新的命令,不会影响其他命令的代码
4. 项目部署在 Github 上，合适多人协作，特别是代码归并
5. 支持日志（原则上不使用debug调试程序）

# 2、GO 命令
[GO命令](https://go-zh.org/cmd/go/) 的官方说明并不一定是最新版本。最新说明请使用命令 `go help` 获取。 [关于GO命令](https://go-zh.org/doc/articles/go_command.html)

必须了解的环境变量：**GOROOT，GOPATH**
[项目目录与 gopath](https://go-zh.org/cmd/go/#hdr-GOPATH_environment_variable)
## 2.1 go 命令的格式
使用：
```
go command [arguments]
```
版本（go 1.8）的命令有：
```
build       compile packages and dependencies
clean       remove object files
doc         show documentation for package or symbol
env         print Go environment information
bug         start a bug report
fix         run go tool fix on packages
fmt         run gofmt on package sources
generate    generate Go files by processing source
get         download and install packages and dependencies
install     compile and install packages and dependencies
list        list packages
run         compile and run Go program
test        test packages
tool        run specified go tool
version     print Go version
vet         run go tool vet on packages
```
## 2.2 go 命令分类
1. 环境显示：version、env
2. 构建流水线： clean、build、test、run、（publish/git）、get、install
3. 包管理： list、get、install
4. 杂项：fmt、vet、doc、tools …

具体命令格式与参数使用 `go help [topic]`
# 3、准备知识或资源
## 3.1 Golang 知识整理
这里推荐 time-track 的个人博客，它的学习轨迹与课程要求基本一致。以下是他语言学习的笔记，可用于语言快速浏览与参考：
+ [《Go程序设计语言》要点总结——程序结构](https://niyanchun.com/)
+ [《Go程序设计语言》要点总结——数据类](https://niyanchun.com/)
+ [《Go程序设计语言》要点总结——函数](https://niyanchun.com)
+ [《Go程序设计语言》要点总结——方法](https://niyanchun.com/)
+ [《Go程序设计语言》要点总结——接口](https://niyanchun.com/)

以上仅代表作者观点，部分内容是不准确的，请用批判的态度看待网上博客。 切记：
+ GO **不是面向对象(OOP)** 的。 所谓方法只是一种语法糖，它是特定类型上定义的操作（operation）
+ 指针是没有 nil 的，这可以避免一些尴尬。 `p.X` 与 `v.x` (p 指针， v 值) 在语义上是无区别的，但实现上是有区别的 `p.x` 是实现 c 语言 `p->x` 的语法糖
+ zero 值好重要
## 3.2 JSON 序列化与反序列化
参考：[JSON and Go](https://blog.go-zh.org/json-and-go)
json 包是内置支持的，文档位置：https://go-zh.org/pkg/encoding/json/
## 3.3 复杂命令行的处理
不要轻易“发明轮子”。为了实现 POSIX/GNU-风格参数处理，–flags，包括命令完成等支持，程序员们开发了无数第三方包，这些包可以在 godoc 找到。
+ pflag 包：https://godoc.org/github.com/spf13/pflag
+ cobra 包：https://godoc.org/github.com/spf13/cobra
+ goptions 包：https://godoc.org/github.com/voxelbrain/goptions
+ …
+ docker command 包：https://godoc.org/github.com/docker/cli/cli/command

[go dead project](https://www.xuebuyuan.com/1588520.html) 非常有用

这里我们选择 cobar 这个工具。

**tip: 安装 cobra**
+ 在 `$GOPATH/src/golang.org/x` 目录下用 `git clone` 下载 `sys` 和 `text` 项目
+ 使用命令 `go get -v github.com/spf13/cobra/cobra` 
+ 使用 `go install github.com/spf13/cobra/cobra`，安装后在 `$GOBIN` 下出现了 cobra 可执行程序。

**Cobra 的简单使用**
创建一个处理命令 `agenda register -uTestUser` 或 `agenda register --user=TestUser` 的小程序。

简要步骤如下：
```
cobra init
cobra add register
```
需要的文件就产生了。 你需要阅读 `main.go` 的 `main()` ；`root.go` 的 `Execute()`；最后修改 `register.go`，`init()` 添加：
```
registerCmd.Flags().StringP("user", "u", "Anonymous", "Help message for username")
```
`Run` 匿名回调函数中添加：
```
username, _ := cmd.Flags().GetString("user")
fmt.Println("register called by " + username)
```
测试命令：
```
$ go run main.go register --user=TestUser
register called by TestUser
```
参考文档：
+ [官方文档](https://github.com/spf13/cobra#overview) 推荐
+ [golang命令行库cobra的使用](https://www.cnblogs.com/borey/p/5715641.html) 中文翻译
# 4、agenda 开发项目
## 4.1 需求描述
+ 业务需求：
	+ 用户注册：
	1. 注册新用户时， 用户需设置一个唯一的用户名和一个密码。另外，还需登记邮箱及电话信息。
	2. 如果注册时提供的用户名已由其他用户使用，应反馈一个适当的出错信息；成功注册后，亦应反馈一个成功注册的信息。
	+ 用户登录：
	1. 用户使用用户名和密码登录 Agenda 系统。
	2. 用户名和密码同时正确则登录成功并反馈一个成功登录的信息。否则，登录失败并反馈一个失败登录的信息。
	+ 用户登出：
	1. 已登录的用户登出系统后，只能使用用户注册和用户登录功能。
+ 功能需求：设计一组命令完成 agenda 的管理，例如：
	+ agenda help ：列出命令说明
	+ agenda register -uUserName –password pass –email=a@xxx.com ：注册用户
	+ agenda help register ：列出 register 命令的描述
	+ agenda cm … : 创建一个会议
原则上一个命令对应一个业务功能
+ 持久化要求：
	+ 使用 json 存储 User 和 Meeting 实体
	+ 当前用户信息存储在 curUser.txt 中
+ 开发需求
	+ 完成两条命令
+ 项目目录
	+ cmd ：存放命令实现代码
	+ entity ：存放 User 和 Meeting 对象读写与处理逻辑
	+ 其他目录 ： 自由添加
+ 日志服务
	+ 使用 [log](https://go-zh.org/pkg/log/) 包记录命令执行情况

# 5、实验过程
## 5.1 安装使用 Cobra
+ 在 `$GOPATH/src/golang.org/x` 目录下用 `git clone` 下载 `sys` 和 `text` 项目。
```
git clone https://github.com/golang/text
git clone https://github.com/golang/sys
```
+ 安装 Cobra
```
go get -v github.com/spf13/cobra/cobra
go install github.com/spf13/cobra/cobra
```

+ 初始化并添加相应指令
```
cobra init --pkg-name Agenda
cobra add register
cobra add login
cobra add logout
```
![在这里插入图片描述](https://img-blog.csdnimg.cn/20191029164930820.png)
## 5.2 实现 Agenda 指令
本次实验一共完成了三条指令：用户注册、用户登录、用户登出。
+ 用户注册：
1. 注册新用户时， 用户需设置一个唯一的用户名和一个密码。另外，还需登记邮箱及电话信息。
2. 如果注册时提供的用户名已由其他用户使用，应反馈一个适当的出错信息；成功注册后，亦应反馈一个成功注册的信息。
```
register a new user, a unique username, a password, an email and a phone required

Usage:
  Agenda register [flags]

Flags:
  -c, --contact string    phone number
  -e, --email string      email address
  -h, --help              help for register
  -p, --password string   password
  -u, --username string   username
```
+ 用户登录：
1. 用户使用用户名和密码登录 Agenda 系统。
2. 用户名和密码同时正确则登录成功并反馈一个成功登录的信息。否则，登录失败并反馈一个失败登录的信息。
```
login a user, a username and a password required

Usage:
  Agenda login -u [username] -p [password] -e [email] -c [phone] [flags]

Flags:
  -h, --help              help for login
  -p, --password string   password
  -u, --username string   username
```

+ 用户登出：
1. 已登录的用户登出系统后，只能使用用户注册和用户登录功能。
```
logout a user

Usage:
  Agenda logout [flags]

Flags:
  -h, --help   help for logout
```

## 5.3 测试结果
查看 Agenda 的基本信息：
```
go run main.go
```
![在这里插入图片描述](https://img-blog.csdnimg.cn/20191029170653268.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FpYW9femhhbmc=,size_16,color_FFFFFF,t_70)

用户注册：
```
go run main.go register -u lzz -p 123 -e 123@qq.com -c 123
```
![在这里插入图片描述](https://img-blog.csdnimg.cn/20191029165812121.png)

注册成功后可以在 userList.txt 中看到用户的相关信息：
![在这里插入图片描述](https://img-blog.csdnimg.cn/20191029165954883.png)

用户登录：
```
go run main.go login -u lzz -p 123
```
![在这里插入图片描述](https://img-blog.csdnimg.cn/20191029165901326.png)

登录成功后可以在 curUser.txt 中看到已登录用户的相关信息：
![在这里插入图片描述](https://img-blog.csdnimg.cn/20191029170111273.png)

用户登出：
```
go run main.go logout
```
![在这里插入图片描述](https://img-blog.csdnimg.cn/2019102917014510.png)

在 log.log 文件中可以看到命令的执行情况：
![在这里插入图片描述](https://img-blog.csdnimg.cn/20191029170506593.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FpYW9femhhbmc=,size_16,color_FFFFFF,t_70)
# [项目地址](https://github.com/Buyiai/Agenda)
