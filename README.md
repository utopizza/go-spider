# 概述
本项目名为 go-spider，一个小型的爬虫程序


# 代码目录
- cache: 缓存层，用于网页去重
- conf: 配置文件存放
- configs: 配置文件加载层
- data: 存放种子文件
- parser: 网页解析服务层
- spider: 爬虫服务封装
- storage: 网页存储服务层
- main: 服务启动入口

 
# 核心逻辑
spider包中的主要几个方法介绍：
 - NewSpiderService() 创建爬虫服务
 - Run() 运行爬虫，让爬虫对象开始工作
 - crawl(url string) 对目标url执行爬取任务
    - extractURLs(url string) 从目标url对应的网页中提取新url
        - parseWebpage() 解析网页
        - saveWebpage() 保存网页
        
其中crawl()是核心函数，爬虫对象会为每个url新开一个goroutine运行crawl方法。
为了控制并发度，每个goroutine在开始爬取时向爬虫对象请求一个token，爬取完后归还该token。
token数量即并发数，由配置文件指定。

crawl方法的主要工作逻辑：
 1. 调用extractURLs()方法，通过目标url请求页面，解析页面，提取出新url数组
 2. 遍历上述url数组，先判断url是否已经爬取过，是则忽略；否则记录下该url，然后开启一个go程执行crawl()方法爬取该url
 
在crawl方法中，还有一些网页爬取深度、爬取时间间隔、优雅退出等的控制。


# 编译
修改Makefile文件中的GO和GOROOT路径为本地相应路径，然后在本项目根目录执行make，编译通过后会在本地创建build目录


# 运行
进入上面build目录的bin目录：cd ./build/bin/

执行命令：./go-spider -c ../conf -l ../log

运行目录：
```
|- build
    |- bin
        |- go-spider         // 可执行二进制
    |- conf                   
        |- spider.conf         // 需手动放置或通过外部平台注入
    |- data
        |- url.data            // url种子文件    
    |- log                         
        |- log.txt             // 日志
    |- output                  // 存放抓取的web文件
```
