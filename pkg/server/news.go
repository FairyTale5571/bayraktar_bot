package server

import (
	"github.com/fairytale5571/bayraktar_bot/pkg/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (r *Router) getNews() string {
	var news models.NewsArray

	rows, err := r.db.Query("SELECT `id`, `title`, `link`, `description`, `published` FROM newsfeed WHERE hasActive = 1 ORDER BY id DESC")
	defer rows.Close()
	if err != nil {
		r.logger.Errorf("cant get news: %v", err)
		return "[\"Undefined\"]"
	}
	for rows.Next() {
		var n models.News
		err = rows.Scan(&n.ID, &n.Title, &n.Link, &n.Description, &n.Published)
		if err != nil {
			r.logger.Errorf("cant get news: %v", err)
			return "[\"Undefined\"]"
		}
		news.News = append(news.News, n)
	}
	return news.MakeArmaArray()
}

func (r *Router) news(c *gin.Context) {
	if cached, err := r.cache.Get("newsFeed"); cached != "" && err == nil {
		c.String(http.StatusOK, cached)
		return
	}
	feed := r.getNews()
	c.String(http.StatusOK, feed)
	if err := r.cache.Set("newsFeed", feed); err != nil {
		return
	}
}