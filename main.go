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
	var goSync sync.WaitGroup
	c.OnHTML("div[id=nav]", func(e *colly.HTMLElement) {
		li := e.DOM.Find("li")
		for i := 0; i < li.Length(); i++ {
			//fmt.Println(li.Eq(i).Find("a").Eq(0).Text()) //*一级菜单
			//o, _ := li.Eq(i).Find("a").Attr("href")
			//fmt.Println(strings.Join(regexp.MustCompile("[0-9]").FindAllString(o, -1), "")) //*一级菜单标识
			a := li.Eq(i).Find("div[class=sonnav] a")
			TaskPool := NewPool(2)
			for j := 0; j < a.Length(); j++ {
				//fmt.Println(a.Eq(j).Text()) //*二级菜单
				goSync.Add(1)
				TaskPool.Add()
				href, _ := a.Eq(j).Attr("href")
				id := strings.Join(regexp.MustCompile("[0-9]").FindAllString(href, -1), "")
				//fmt.Println(id) //*二级菜单标识
				go deep2Menu(id, &goSync, TaskPool)
			}
		}
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Starting...", r.URL.String())
	})
	c.OnError(func(r *colly.Response, err error) {
		log.Error("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err.Error())
	})
	c.Visit(Url)
	c.Wait()
	goSync.Wait()
}

//二级内页
func deep2Menu(parentId string, n *sync.WaitGroup, TaskPool *ConcurrentPool) {
	var goSyncMe sync.WaitGroup
	pageCount := deep2MenuPageCount(Url + "list_" + parentId + Suffix)
	c := colly.NewCollector()
	for i := 1; i <= pageCount; i++ {
		//fmt.Println(parentId) //*图片父级标识
		c.OnHTML("div[id=mainbodypul]", func(e *colly.HTMLElement) {
			div := e.DOM.Find("div[class!=listmainrowstag]")
			TaskPool := NewPool(5)
			for i := 0; i < div.Length(); i++ {
				//fmt.Println(div.Eq(i).Find("a").Eq(1).Text()) //*图片标题
				//img, _ := div.Eq(i).Find("img").Attr("data-original")
				//fmt.Println(img) //*图片缩略图
				//loc, _ := time.LoadLocation("Local")
				//tm, _ := time.ParseInLocation("2006-01-02 15:04:05", div.Eq(i).Find("p[class=listmiantimer]").Text(), loc)
				//fmt.Println(tm.Unix()) //*图片时间
				href, _ := div.Eq(i).Find("a").First().Attr("href")
				if regexp.MustCompile("[0-9]").FindAllString(href, -1) != nil {
					goSyncMe.Add(1)
					TaskPool.Add()
					sl := strings.Split(href, "/")
					id := strings.Join(regexp.MustCompile("[0-9]").FindAllString(sl[len(sl)-1], -1), "")
					//fmt.Println(id) //*图片标识
					go deep3Menu(id, &goSyncMe, TaskPool)
				}
			}
		})
		c.OnError(func(r *colly.Response, err error) {
			log.Error("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err.Error())
		})
		c.Visit(Url + "list_" + parentId + "_" + strconv.Itoa(i) + Suffix)
		c.Wait()
	}
	time.Sleep(time.Second)
	goSyncMe.Wait()
	defer n.Done()
	defer TaskPool.Done()
}

//二级内容总页码
func deep2MenuPageCount(link string) int {
	var count = 1
	c := colly.NewCollector()
	c.OnHTML("ul", func(e *colly.HTMLElement) {
		if href, ok := e.DOM.Find("a[class=end]").Attr("href"); ok {
			sl := strings.Split(href, "_")
			countInt, _ := strconv.Atoi(strings.Join(regexp.MustCompile("[0-9]").FindAllString(sl[len(sl)-1], -1), ""))
			count = countInt
		} else {
			if href, ok := e.DOM.Find("a[class=num]").Last().Attr("href"); ok {
				sl := strings.Split(href, "_")
				countInt, _ := strconv.Atoi(strings.Join(regexp.MustCompile("[0-9]").FindAllString(sl[len(sl)-1], -1), ""))
				count = countInt
			}
		}
	})
	c.OnError(func(r *colly.Response, err error) {
		log.Error("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err.Error())
	})
	c.Visit(link)
	return count
}

//三级内页
func deep3Menu(ID string, n *sync.WaitGroup, TaskPool *ConcurrentPool) {
	pageCount := deep3MenuPageCount(Url + ID + Suffix)
	for i := 1; i <= pageCount; i++ {
		c := colly.NewCollector()
		c.OnHTML("h1[class=center]", func(e *colly.HTMLElement) {
			if imgUrl, ok := e.DOM.Closest("div").Find("img").First().Attr("src"); ok {
				downFile(imgUrl)
			}
		})
		c.OnError(func(r *colly.Response, err error) {
			log.Error("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err.Error())
		})
		c.Visit(Url + ID + "_" + strconv.Itoa(i) + Suffix)
		c.Wait()
	}
	time.Sleep(time.Second)
	defer n.Done()
	defer TaskPool.Done()
}

//三级内页总页码
func deep3MenuPageCount(link string) int {
	var count = 1
	c := colly.NewCollector()
	c.OnHTML("h1[class=center]", func(e *colly.HTMLElement) {
		sl := strings.Split(e.DOM.Text(), "/")
		countInt, _ := strconv.Atoi(strings.Join(regexp.MustCompile("[0-9]").FindAllString(sl[len(sl)-1], -1), ""))
		count = countInt
	})
	c.OnError(func(r *colly.Response, err error) {
		log.Error("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err.Error())
	})
	c.Visit(link)
	return count
}

var sum = 0

//下载文件
func downFile(imgUrl string) {
	sum++
	defer fmt.Println(sum, "目标文件:"+imgUrl)
	res, err := http.Get(imgUrl)
	if err != nil {
		log.Error(err)
		return
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
			return
		}
	} else {
		fh, err = os.Open(imgUrl)
		if err != nil {
			log.Error(err)
			return
		}
	}
	imgByte, _ := ioutil.ReadAll(res.Body)
	fh.Write(imgByte)
	fh.Close()
	res.Body.Close()
}
