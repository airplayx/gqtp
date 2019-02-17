package main

import (
	"fmt"
	. "github.com/bingoladen/gqtp/config"
	"github.com/bingoladen/gqtp/log"
	"github.com/gocolly/colly"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main() {
	var c = colly.NewCollector(
		colly.Async(true),
	)
	// 更新菜单
	var goSync sync.WaitGroup
	c.OnHTML("div[id=nav]", func(e *colly.HTMLElement) {
		li := e.DOM.Find("li")
		for i := 0; i < li.Length(); i++ {
			//一级菜单
			a := li.Eq(i).Find("div[class=sonnav] a")
			TaskPool := NewPool(2)
			for j := 0; j < a.Length(); j++ {
				//二级菜单
				link := a.Eq(j)
				href, _ := link.Attr("href")
				var valid = regexp.MustCompile("[0-9]")
				var ID = valid.FindAllString(href, -1)
				goSync.Add(1)
				TaskPool.Add()
				go deep2Menu(c.Clone(), strings.Join(ID, ""), &goSync, TaskPool)
			}
		}
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Starting...", r.URL.String())
	})
	c.OnError(func(r *colly.Response, err error) {
		log.Error("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})
	c.Visit(Url)
	c.Wait()
	goSync.Wait()
}

//二级内页
func deep2Menu(c *colly.Collector, parentId string, n *sync.WaitGroup, TaskPool *ConcurrentPool) {
	var goSyncMe sync.WaitGroup
	pageCount := deep2MenuPageCount(c.Clone(), Url+"list_"+parentId+Suffix)
	//fmt.Println(Url+"list_"+parentId+Suffix, pageCount)
	for i := 1; i <= pageCount; i++ {
		page := strconv.Itoa(i)
		//fmt.Println("当前分类："+parentId, " 页码："+page+"/", pageCount)
		c.OnHTML("div[id=mainbodypul]", func(e *colly.HTMLElement) {
			//当前页主内容
			div := e.DOM.Find("div[class!=listmainrowstag]")
			TaskPool := NewPool(5)
			for i := 0; i < div.Length(); i++ {
				a := div.Eq(i).Find("a").First()
				href, _ := a.Attr("href") //内页链接
				var valid = regexp.MustCompile("[0-9]")
				var id = valid.FindAllString(href, -1)
				if id != nil {
					sl := strings.Split(href, "/")
					var valid = regexp.MustCompile("[0-9]")
					var ID = valid.FindAllString(sl[len(sl)-1], -1)
					goSyncMe.Add(1)
					TaskPool.Add()
					go deep3Menu(c.Clone(), strings.Join(ID, ""), &goSyncMe, TaskPool)
				}
			}
		})
		c.OnError(func(r *colly.Response, err error) {
			log.Error("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
		})
		c.Visit(Url + "list_" + parentId + "_" + page + Suffix)
		c.Wait()
	}
	time.Sleep(time.Second)
	goSyncMe.Wait()
	defer n.Done()
	defer TaskPool.Done()
}

//二级内容总页码
func deep2MenuPageCount(c *colly.Collector, link string) int {
	var count = 1
	c.OnHTML("ul", func(e *colly.HTMLElement) {
		//总页码
		if href, ok := e.DOM.Find("a[class=end]").Attr("href"); ok {
			sl := strings.Split(href, "_")
			var valid = regexp.MustCompile("[0-9]")
			var countPage = valid.FindAllString(sl[len(sl)-1], -1)
			countInt, _ := strconv.Atoi(strings.Join(countPage, ""))
			count = countInt
		} else {
			if href, ok := e.DOM.Find("a[class=num]").Last().Attr("href"); ok {
				sl := strings.Split(href, "_")
				var valid = regexp.MustCompile("[0-9]")
				var countPage = valid.FindAllString(sl[len(sl)-1], -1)
				countInt, _ := strconv.Atoi(strings.Join(countPage, ""))
				count = countInt
			}
		}
	})
	c.OnError(func(r *colly.Response, err error) {
		log.Error("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})
	c.Visit(link)
	return count
}

//三级内页
func deep3Menu(c *colly.Collector, ID string, n *sync.WaitGroup, TaskPool *ConcurrentPool) {
	pageCount := deep3MenuPageCount(c.Clone(), Url+ID+Suffix)
	for i := 1; i <= pageCount; i++ {
		page := strconv.Itoa(i)
		c.OnHTML("h1[class=center]", func(e *colly.HTMLElement) {
			//当前页主内容
			if imgUrl, ok := e.DOM.Closest("div").Find("img").First().Attr("src"); ok {
				downFile(imgUrl)
			}
		})
		c.OnError(func(r *colly.Response, err error) {
			log.Error("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
		})
		c.Visit(Url + ID + "_" + page + Suffix)
		c.Wait()
	}
	time.Sleep(time.Second)
	defer n.Done()
	defer TaskPool.Done()
}

//三级内页总页码
func deep3MenuPageCount(c *colly.Collector, link string) int {
	var count = 1
	c.OnHTML("h1[class=center]", func(e *colly.HTMLElement) {
		//总页码
		sl := strings.Split(e.DOM.Text(), "/")
		var valid = regexp.MustCompile("[0-9]")
		var countPage = valid.FindAllString(sl[len(sl)-1], -1)
		countInt, _ := strconv.Atoi(strings.Join(countPage, ""))
		count = countInt
	})
	c.OnError(func(r *colly.Response, err error) {
		log.Error("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})
	c.Visit(link)
	return count
}

//下载文件
func downFile(imgUrl string) {
	defer fmt.Println("下载文件:" + imgUrl)
	res, err := http.Get(imgUrl)
	if err != nil {
		log.Error(err)
	}
	defer res.Body.Close()
	if err != nil {
		log.Error(err)
	}
	imgUrl = strings.Replace(imgUrl, FileUrl, "", -1)
	imgPath := strings.Split(imgUrl, "/")
	imgPath = imgPath[:len(imgPath)-1]
	os.MkdirAll(strings.Join(imgPath, "/"), os.ModePerm)
	var fh *os.File
	_, fErr := os.Stat(imgUrl)
	if fErr != nil {
		fh, err = os.Create(imgUrl)
		if err != nil {
			log.Error(fErr, err)
		}
	} else {
		fh, err = os.Open(imgUrl)
		if err != nil {
			log.Error(err)
		}
	}
	defer fh.Close()
	imgByte, _ := ioutil.ReadAll(res.Body)
	fh.Write(imgByte)
}
