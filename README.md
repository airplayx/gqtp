# 多线程爬虫的简单实践
Go实现获取高清无水印图片网站的简单爬虫程序
![image](https://github.com/bingoladen/gqtp/blob/master/show.gif)
=====================
- 目标网站 [http://www.gqtp.com/](http://www.gqtp.com/) 
- 已编译windows版本可执行程序，在release下载体验。

<h4>依赖</h4> <code>go get -u github.com/gocolly/colly/...</code>
- 理论上可以成千上万个线程(或者说协程？)同时抓取，已测试目标服务器有点扛不住，但比普通服务器要友好。
- 比较稳定时候最多同时分类用**2**个线程，和内页**5**个线程跑。
- 之前用世界上最好的语言**PHP**获取目标分类多进程跑，单核双核服务器瞬间CPU**100%**，内存爆炸。

<h4>在线演示</h4> [http://2me.cc/](http://2me.cc/) 

<h4>注意</h4>
- 如果你想整个任务跑完，请准备**60G+**的硬盘！