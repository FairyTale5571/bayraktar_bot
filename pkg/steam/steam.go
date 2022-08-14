package steam

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/fairytale5571/bayraktar_bot/pkg/logger"
	"github.com/fairytale5571/bayraktar_bot/pkg/models"
	"github.com/markbates/goth/providers/steam"
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
	defer res.Body.Close() // nolint: not needed
	if err != nil {
		s.logger.Fatalf("cant get workshop info: %v", err)
		return nil
	}

	if res.StatusCode != http.StatusOK {
		s.logger.Fatalf("cant get workshop info: %v", res.StatusCode)
		return nil
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		s.logger.Fatalf("cant get workshop info: %v", err)
		return nil
	}
	return doc
}

func (s *Steam) workshopChangelogs(itemId string) *goquery.Document {
	res, err := http.Get(EndpointSteamWorkshopChangelog + "/" + itemId)
	defer res.Body.Close() // nolint: not needed

	if err != nil {
		s.logger.Errorf("cant get workshop info: %v", err)
		return nil
	}

	if res.StatusCode != http.StatusOK {
		s.logger.Errorf("cant get workshop info: %v", res.StatusCode)
		return nil
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		s.logger.Errorf("cant get workshop info: %v", err)
		return nil
	}
	return doc
}

func (s *Steam) GetLatestUpdate(itemId string) (update, id string) {
	doc := s.workshopChangelogs(itemId)
	if doc == nil {
		return "", ""
	}
	doc.Find(".workshopAnnouncement").EachWithBreak(
		func(i int, s *goquery.Selection) bool {
			text := s.Find("p")
			_id, _ := text.Attr("id")
			html, _ := text.Html()
			update = strings.ReplaceAll(html, "<br/>", "\n")
			id = _id
			return false
		})
	return update, id
}

func (s *Steam) GetItemTitle(itemId string) {
	doc := s.workshopInfo(itemId)
	if doc == nil {
		return
	}

	sel := doc.Find("div.workshopItemTitle").Each(
		func(i int, s *goquery.Selection) {
			text := s.Text()
			fmt.Printf("%d: %s\n", i, text)
		},
	)
	sel.Text()
}

func (s *Steam) GetItemLogo(itemId string) (logo string) {
	doc := s.workshopInfo(itemId)
	if doc == nil {
		return ""
	}

	doc.Find("head link").Each(
		func(i int, s *goquery.Selection) {
			if val, exist := s.Attr("rel"); val == "image_src" && exist {
				logo = s.AttrOr("href", "")
			}
		},
	)
	return logo
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

func (s *Steam) GetAuthLink(guild, state string) string {
	client := steam.New(s.cfg.SteamKey, s.cfg.URL+"/auth/steam/?guild="+guild+"&state="+state)
	session, err := client.BeginAuth(state)
	if err != nil {
		s.logger.Errorf("cant get auth link: %v", err)
		return ""
	}

	url, err := session.GetAuthURL()
	if err != nil {
		s.logger.Errorf("cant get auth link: %v", err)
		return ""
	}
	return url
}
