package handlers

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/gofiber/fiber/v2"
	"io/ioutil"
	"log"
	"math/big"
	"my-rest-api/entities"
	"net/http"
	"strings"
)

func ApiHandler(ctx *fiber.Ctx) error {

	//var actualUserId string
	//instagramAccount := "kaimook.bnk48official"
	if ctx.Params("account") == "" {
		json, _ := json.Marshal("Please Provide an account")
		return ctx.Status(500).Send(json)
	}

	instagramAccount := ctx.Params("account")
	url := "https://i.instagram.com/api/v1/users/web_profile_info/?username=" + instagramAccount

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("user-agent", "Mozilla/5.0 (Linux; Android 9; GM1903 Build/PKQ1.190110.001; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/75.0.3770.143 Mobile Safari/537.36 Instagram 103.1.0.15.119 Android (28/9; 420dpi; 1080x2260; OnePlus; GM1903; OnePlus7; qcom; sv_SE; 164094539)")
	cookies, _ := ioutil.ReadFile("cookie.txt")
	if string(cookies) != "" {
		req.Header.Set("cookie", string(cookies))
	}
	res, err := client.Do(req)
	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}

	if res.StatusCode != 200 {
		myString := string(body)
		log.Println("API Error:", myString)
		return ctx.Status(500).SendString("An Error Occurred with the API")
	}

	data := &entities.ProfileInfo{}
	err2 := json.Unmarshal(body, &data)
	fmt.Printf("Results: %v\n", data)
	if err2 != nil {
		return err2
	}
	average := CalculateAverage(data.DataInfo.User)
	json, _ := json.Marshal(average)
	return ctx.Send(json)
}

func CalculateAverage(data entities.UserInfo) *big.Float {
	var sum float64 = 0
	for _, obj := range data.Media.Edges {
		// Engagement Rate = (likes + comments) / followers
		var engagementRate = float64(obj.Node.EdgeLikedBy.Count+obj.Node.EdgeMediaToComment.Count) / float64(data.EdgeFollowedBy.Count)
		log.Println("<!------------------------------> ")
		log.Println("EdgeLikedBy ->  ", obj.Node.EdgeLikedBy.Count)
		log.Println("EdgeMediaToComment -> ", obj.Node.EdgeMediaToComment.Count)
		log.Println("EdgeFollowedBy -> ", data.EdgeFollowedBy.Count)
		log.Println("engagement rate --> ", engagementRate)
		log.Println("<!------------------------------> ")
		sum += engagementRate
	}
	average := sum * 100 / float64(len(data.Media.Edges))
	f := new(big.Float).SetMode(big.AwayFromZero).SetFloat64(average)
	return f.SetPrec(8)
}

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
	var average *big.Float
	var hasError bool = false

	c.OnHTML("html", func(e *colly.HTMLElement) {

		dat := e.ChildText("body > script:first-of-type")
		jsonData := dat[strings.Index(dat, "{") : len(dat)-1]
		data := &entities.ProfileData{}
		err := json.Unmarshal([]byte(jsonData), data)
		if err != nil {
			hasError = true
			log.Println(err)
			return

		}
		average = CalculateAverage(data.EntryData.ProfilePage[0].Graphql.User)

	})

	c.OnError(func(r *colly.Response, e error) {
		log.Println("error:", e, r.Request.URL, string(r.Body))
	})

	var err = c.Visit("https://www.instagram.com/" + instagramAccount + "/")
	if err != nil {
		log.Println(err)
		return err
	}

	if hasError {
		json, _ := json.Marshal("The Account does not exist")
		return ctx.Status(500).Send(json)
	}

	json, _ := json.Marshal(average)
	return ctx.Send(json)
}
