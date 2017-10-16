/*
#!/usr/bin/env gorun
@author :yinzhengjie
Blog:http://www.cnblogs.com/yinzhengjie/tag/GO%E8%AF%AD%E8%A8%80%E7%9A%84%E8%BF%9B%E9%98%B6%E4%B9%8B%E8%B7%AF/
EMAIL:y1053419035@qq.com
*/

package main

import (
    "errors"
    "fmt"
    "github.com/PuerkitoBio/goquery"
    "log"
    "net/http"
    "net/url"
    "os"
    "strings"
    "path/filepath"
)

func fetch(url string) ([]string, error) { //改函数会拿到我们想要的图片的路径。
    var urls []string //定义一个空切片数组
    resp, err := http.Get(url)
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK {
        return nil, errors.New(resp.Status) //表示当出现错误是，返回空列表，并将错误状态返回。
    }
    doc, err := goquery.NewDocumentFromResponse(resp)
    if err != nil {
        log.Fatal(err)
    }
    doc.Find("img").Each(func(i int, s *goquery.Selection) {
        link, ok := s.Attr("src")
        if ok {
            urls = append(urls, link) //将过滤出来的图片路径都追加到urls的数组中去,最终返回给用户。
        } else {
            fmt.Println("抱歉，没有发现该路径。")
        }

    })
    return urls, nil
}

func Clean_urls(root_path string, picture_path []string) []string {
    var Absolute_path []string //定义一个绝对路径数组。
    url_info, err := url.Parse(root_path)
    if err != nil {
        log.Fatal(err)
    }
    Scheme := url_info.Scheme //获取到链接的协议
    //fmt.Println("使用的协议是：",Scheme)
    Host := url_info.Host //获取链接的主机名
    for _, souce_path := range picture_path {
        if strings.HasPrefix(souce_path, "https") { //如果当前当前路径是以“https”开头说明是绝对路径，因此我们给一行空代码，表示不执行任何操作，千万别写：“continue”，空着就好。

        } else if strings.HasPrefix(souce_path, "//") { //判断当前路径是否以“//”开头(说明包含主机名)
            souce_path = Scheme + ":" + souce_path //如果是就对其进行拼接操作。以下逻辑相同。
        } else if strings.HasPrefix(souce_path, "/") { //说明不包含主机名和协议，我们进行拼接即可。
            souce_path = Scheme + "://" + Host + souce_path
        } else {
            souce_path = filepath.Dir(root_path) + souce_path  //文件名称和用户输入的目录相拼接。
        }
        Absolute_path = append(Absolute_path, souce_path) //不管是否满足上面的条件，最终都会被追加到该数组中来。
    }
    return Absolute_path //最终返回处理后的每个链接的绝对路基。
}

func main() {
    root_path := os.Args[1]               //定义一个URl，也就是我们要爬的网站。
    picture_path, err := fetch(root_path) //“fetch”函数会帮我们拿到picture_path的路径，但是路径可能是相对路径或是绝对路径。不同意。
    if err != nil {
        log.Fatal(err)
    }

    Absolute_path := Clean_urls(root_path, picture_path) //“Clean_urls”函数会帮我们把picture_path的路径做一个统一，最终都拿到了绝对路径Absolute_path数组。

    for _, Picture_absolute_path := range Absolute_path {
        fmt.Println(Picture_absolute_path) //最终我们会得到一个图片的完整路径，我们可以对这个路径进行下载，压缩，加密等等操作。
    }
}
