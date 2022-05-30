package handlers

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/gofiber/fiber/v2"
	"io/ioutil"
	"log"
	"math"
	"my-rest-api/entities"
	"strings"
)

func Calculate(ctx *fiber.Ctx) error {

	if ctx.Params("account") == "" {
		json, _ := json.Marshal("Please Provide an account")
		return ctx.Status(500).Send(json)
	}

	instagramAccount := ctx.Params("account")

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"),
	)

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("X-Requested-With", "XMLHttpRequest")
		r.Headers.Set("Referer", "https://www.instagram.com/"+instagramAccount+"/")
		cookies, _ := ioutil.ReadFile("cookie.txt")
		if string(cookies) != "" {
			r.Headers.Set("cookie", string(cookies))
		}
		if r.Ctx.Get("gis") != "" {
			gis := fmt.Sprintf("%s:%s", r.Ctx.Get("gis"), r.Ctx.Get("variables"))
			h := md5.New()
			h.Write([]byte(gis))
			gisHash := fmt.Sprintf("%x", h.Sum(nil))
			r.Headers.Set("X-Instagram-GIS", gisHash)
		}
	})
	var average float64 = 0
	var hasError bool = false

	c.OnHTML("html", func(e *colly.HTMLElement) {

		dat := e.ChildText("body > script:first-of-type")
		jsonData := dat[strings.Index(dat, "{") : len(dat)-1]
		data := &entities.ProfileData{}
		err := json.Unmarshal([]byte(jsonData), data)
		if err != nil {
			hasError = true
			return
			//log.Fatal(err)
		}

		var sum float64 = 0
		page := data.EntryData.ProfilePage[0]
		for _, obj := range page.Graphql.User.Media.Edges {
			// Engagement Rate = (likes + comments) / followers
			var engagementRate = float64(obj.Node.EdgeLikedBy.Count+obj.Node.EdgeMediaToComment.Count) / float64(page.Graphql.User.EdgeFollowedBy.Count)
			log.Println("<!------------------------------> ")
			log.Println("EdgeLikedBy ->  ", obj.Node.EdgeLikedBy.Count)
			log.Println("EdgeMediaToComment -> ", obj.Node.EdgeMediaToComment.Count)
			log.Println("EdgeFollowedBy -> ", page.Graphql.User.EdgeFollowedBy.Count)
			log.Println("engagement rate --> ", engagementRate)
			log.Println("<!------------------------------> ")
			sum += engagementRate
		}
		average = sum * 100 / float64(len(page.Graphql.User.Media.Edges))
	})

	c.OnError(func(r *colly.Response, e error) {
		log.Println("error:", e, r.Request.URL, string(r.Body))
	})

	var err = c.Visit("https://www.instagram.com/" + instagramAccount + "/")
	if err != nil {
		return err
	}

	if hasError {
		json, _ := json.Marshal("The Account does not exist")
		return ctx.Status(500).Send(json)
	}

	json, _ := json.Marshal(Round(average, 0.005))
	return ctx.Send(json)
}

func Round(x, unit float64) float64 {
	return math.Round(x/unit) * unit
}
