package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func (r *Router) economy(c *gin.Context) {
	rows, err := r.db.Query("SELECT * FROM economy")
	defer rows.Close()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	type Economy struct {
		Resource         string
		Localize         string
		Price            int
		MaxPrice         int
		MinPrice         int
		DownPricePerItem float64
		RandomDownPrice  bool
		RandomMax        int
		RandomMin        int
		Illegal          bool
		Influenced       string
		LastUpdate       time.Time
	}
	var economies []Economy
	for rows.Next() {
		var e Economy
		err := rows.Scan(&e.Resource, &e.Localize, &e.Price, &e.MaxPrice, &e.MinPrice, &e.DownPricePerItem, &e.RandomDownPrice, &e.RandomMax, &e.RandomMin, &e.Illegal, &e.Influenced, &e.LastUpdate)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		economies = append(economies, e)
	}
	c.JSON(http.StatusOK, economies)
}
