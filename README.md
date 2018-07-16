# gqtp
Go实现获取高清无水印图片网站的简单爬虫程序

![image](https://github.com/bingoladen/gqtp/blob/master/show.gif)
=====================
图片来自 [http://www.gqtp.com/](http://www.gqtp.com/) 

可下载release中的windows可执行程序体验。

依赖有名的Go爬虫程序 [colly](http://github.com/gocolly/colly/)

理论上可以成千上万个线程(或者说协程？)同时抓取，已测试目标服务器有点扛不住，但比普通服务器要友好。

比较稳定时候最多同时分类用2个线程，和内页5个线程跑。

下载速度一般，之前做过PHP获取目标分类多进程跑，单核双核服务器瞬间CPU100%，内存爆炸。

在线体验地址 [http://2me.cc/](http://2me.cc/) 

如果你想整个任务跑完，请准备60G+的硬盘！