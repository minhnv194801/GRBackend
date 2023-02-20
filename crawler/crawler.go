package crawler

import (
	"context"
	"fmt"
	"io"
	"log"
	"magna/model"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/chromedp"
)

const (
	numberOfPage  = 5
	searchUrl     = "https://www.nettruyenking.com/tim-truyen?status=2&sort=10&page="
	excludedImage = "https://u.ntcdntempv3.com/content/2022-11-23/638047952612608555.jpg"
)

var wg sync.WaitGroup

func newChromedp() (context.Context, context.CancelFunc) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.Flag("start-fullscreen", false),
	)

	allocCtx, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))

	wg.Add(1)
	go func() {
		defer wg.Done()
		crawlManga(ctx, searchUrl+"1")
	}()
	for i := 2; i <= numberOfPage; i++ {
		newTabCtx, _ := chromedp.NewContext(ctx)
		wg.Add(1)
		innerI := i
		go func() {
			defer wg.Done()
			crawlManga(newTabCtx, searchUrl+strconv.Itoa(innerI))
		}()
	}

	wg.Wait()

	return ctx, cancel
}

func crawlManga(ctx context.Context, url string) {
	task := chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.Sleep(5 * time.Second),
		chromedp.ActionFunc(getAllMangaData),
	}

	if err := chromedp.Run(ctx, task); err != nil {
		fmt.Println(err)
	}
}

func getAllMangaData(ctx context.Context) error {
	node, err := dom.GetDocument().Do(ctx)
	if err != nil {
		log.Fatal(err)
		return err
	}
	res, err := dom.GetOuterHTML().WithNodeID(node.NodeID).Do(ctx)
	if err != nil {
		log.Fatal(err)
		return err
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(res))
	if err != nil {
		log.Fatal(err)
		return err
	}

	doc.Find("div.item div.image a").Each(func(index int, info *goquery.Selection) {
		mangaUrl, _ := info.Attr("href")
		getMangaData(ctx, mangaUrl)
	})

	return nil
}

func getMangaData(ctx context.Context, mangaUrl string) error {
	chromedp.Navigate(mangaUrl).Do(ctx)
	chromedp.Sleep(5 * time.Second).Do(ctx)
	node, err := dom.GetDocument().Do(ctx)
	if err != nil {
		log.Fatal(err)
		return err
	}
	res, err := dom.GetOuterHTML().WithNodeID(node.NodeID).Do(ctx)
	if err != nil {
		log.Fatal(err)
		return err
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(res))
	if err != nil {
		log.Fatal(err)
		return err
	}

	titleNode := doc.Find("#item-detail > h1").First()
	title := titleNode.Text()
	fmt.Println("Name:", title)

	coverNode := doc.Find("#item-detail > div.detail-info > div > div.col-xs-4.col-image > img").First()
	cover, _ := coverNode.Attr("src")
	cover = "https:" + cover
	fmt.Println("Cover:", cover)

	alternateNameNode := doc.Find("#item-detail > div.detail-info > div > div.col-xs-8.col-info > ul > li.othername.row > h2").First()
	alternateNames := strings.Split(alternateNameNode.Text(), ";")
	for i := range alternateNames {
		alternateNames[i] = strings.TrimSpace(alternateNames[i])
	}
	fmt.Println("AlternateName:", alternateNames)

	authorNode := doc.Find("#item-detail > div.detail-info > div > div.col-xs-8.col-info > ul > li.author.row > p.col-xs-8 > a").First()
	author := authorNode.Text()
	if strings.Compare(author, "Đang cập nhật") == 0 {
		author = ""
	}
	fmt.Println("Author:", author)

	statusNode := doc.Find("#item-detail > div.detail-info > div > div.col-xs-8.col-info > ul > li.status.row > p.col-xs-8").First()
	var status model.Status
	if statusNode.Text() == "Đang tiến hành" {
		status = model.Ongoing
	} else {
		status = model.Finished
	}
	fmt.Println("Status:", status)

	tagsNode := doc.Find("#item-detail > div.detail-info > div > div.col-xs-8.col-info > ul > li.kind.row > p.col-xs-8").First()
	var tags []string
	tagsNode.Find("a[href]").Each(func(index int, info *goquery.Selection) {
		tags = append(tags, info.Text())
	})
	fmt.Println("Tags:", tags)

	descriptionNode := doc.Find("#item-detail > div.detail-content > p").First()
	description := descriptionNode.Text()
	fmt.Println("Description:", description)

	manga := new(model.Manga)
	manga.Name = title
	manga.Cover = cover
	manga.AlternateName = alternateNames
	if author != "" {
		manga.Author = append(manga.Author, author)
	}
	manga.Status = status
	manga.Tags = tags
	manga.Description = description

	chapterNodes := doc.Find("div > nav > ul").First()
	var chapters []*model.Chapter
	var chapterUrl []string
	chapterNodes.Find("a[href]").Each(func(index int, info *goquery.Selection) {
		chapter := new(model.Chapter)
		chapter.Name = info.Text()
		chapter.Cover = manga.Cover
		chapters = append(chapters, chapter)
		url, _ := info.Attr("href")
		chapterUrl = append(chapterUrl, url)
	})

	// If the manga is too big then ignore it
	if len(chapters) > 500 {
		return nil
	}

	fmt.Println(manga.InsertToDatabase())

	for i := len(chapters) - 1; i >= 0; i-- {
		chapter := chapters[i]
		chapter.Manga = manga.Id
		chapter.UpdateTime = uint(time.Now().Unix())
		chapter.InsertToDatabase()
		manga.UpdateChapter(chapter)
		filepath := "./Manga/" + manga.Id.Hex() + "/" + chapter.Id.Hex()
		alreadyCrawled, _ := exists(filepath)
		if alreadyCrawled {
			continue
		}
		err := os.MkdirAll(filepath, os.ModePerm)
		if err != nil {
			log.Println(err)
		}

		chromedp.Navigate(chapterUrl[i]).Do(ctx)
		chromedp.Sleep(5 * time.Second).Do(ctx)

		node, err := dom.GetDocument().Do(ctx)
		if err != nil {
			return err
		}
		res, err := dom.GetOuterHTML().WithNodeID(node.NodeID).Do(ctx)
		if err != nil {
			return err
		}
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(res))
		if err != nil {
			return err
		}

		doc.Find("div.page-chapter").Each(func(index int, info *goquery.Selection) {
			img := info.Find("img").First()
			imgUrl, exists := img.Attr("src")
			if !exists {
				log.Fatal("error")
			}
			imgUrl = "https:" + imgUrl
			if strings.Compare(imgUrl, excludedImage) != 0 {
				downloadImage(filepath+"/"+strconv.Itoa(index)+".jpg", imgUrl)
			}
		})
	}

	return nil
}

func downloadImage(filepath, url string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("referer", "https://www.nettruyenking.com/")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	//open a file for writing
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Use io.Copy to just dump the response body to the file. This supports huge files
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
