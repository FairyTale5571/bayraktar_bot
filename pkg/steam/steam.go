package steam

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/fairytale5571/bayraktar_bot/pkg/logger"
	"github.com/fairytale5571/bayraktar_bot/pkg/models"
)

const (
	EndpointSteamCommunity         = "https://steamcommunity.com"
	EndpointSteamWorkshop          = EndpointSteamCommunity + "/sharedfiles/filedetails"
	EndpointSteamWorkshopChangelog = EndpointSteamWorkshop + "/changelog"
)

type Steam struct {
	cfg    *models.Config
	logger *logger.LoggerWrapper
}

func New(cfg *models.Config) *Steam {
	return &Steam{
		cfg:    cfg,
		logger: logger.New("steam"),
	}
}

func (s *Steam) workshopInfo(itemId string) *goquery.Document {
	res, err := http.Get(EndpointSteamWorkshop + "/?id=" + itemId)
	defer res.Body.Close() // nolint: errcheck
	if err != nil {
		log.Fatalf("cant get workshop info: %v", err)
		return nil
	}

	if res.StatusCode != http.StatusOK {
		log.Fatalf("cant get workshop info: %v", res.StatusCode)
		return nil
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatalf("cant get workshop info: %v", err)
		return nil
	}
	return doc

}

func (s *Steam) workshopChangelogs(itemId string) *goquery.Document {
	res, err := http.Get(EndpointSteamWorkshopChangelog + "/" + itemId)
	if err != nil {
		log.Fatalf("cant get workshop info: %v", err)
		return nil
	}
	defer res.Body.Close() // nolint: errcheck

	if res.StatusCode != http.StatusOK {
		log.Fatalf("cant get workshop info: %v", res.StatusCode)
		return nil
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatalf("cant get workshop info: %v", err)
		return nil
	}
	return doc

}

func (s *Steam) GetLatestUpdate(itemId string) (string, string) {
	doc := s.workshopChangelogs(itemId)
	var update, id string
	doc.Find(".workshopAnnouncement").EachWithBreak(
		func(i int, s *goquery.Selection) bool {
			text := s.Find("p")
			_id, _ := text.Attr("id")
			update = text.Text()
			id = _id
			return false
		})
	return update, id
}

func (s *Steam) GetItemTitle(itemId string) {
	doc := s.workshopInfo(itemId)
	sel := doc.Find("div.workshopItemTitle").Each(
		func(i int, s *goquery.Selection) {
			text := s.Text()
			fmt.Printf("%d: %s\n", i, text)
		},
	)
	sel.Text()
}

func (s *Steam) GetItemLogo(itemId string) {
	doc := s.workshopInfo(itemId)
	doc.Find("head link").Each(
		func(i int, s *goquery.Selection) {
			if val, exist := s.Attr("rel"); val == "image_src" && exist {
				href := s.AttrOr("href", "")
				fmt.Printf("%d: %s\n", i, href)
			}
		},
	)
}

func (s *Steam) GetItemSize(itemId string) {
	doc := s.workshopInfo(itemId)
	doc.Find(".detailsStatsContainerRight > div:nth-child(1)").Each(
		func(i int, s *goquery.Selection) {
			text := s.Text()
			fmt.Printf("%d: %s\n", i, text)
		},
	)
}
