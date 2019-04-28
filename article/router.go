package article

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Retrieve(c *gin.Context) {
	author := c.Query("author")
	results := RetrieveArticles(author)

	c.JSON(http.StatusOK, gin.H{
		"articles": results,
	})
}

func Create(c *gin.Context) {
	var paramArticle Article
	c.Bind(&paramArticle)
	if newArticle, err := RegisterArticle(paramArticle); err == nil {
		c.JSON(http.StatusCreated, gin.H{
			"article" : newArticle,
		})
	} else {
		c.JSON(http.StatusInternalServerError, nil)
	}
}

func Update(c *gin.Context) {
	var paramArticle Article
	c.ShouldBindUri(&paramArticle)
	c.Bind(&paramArticle)

	if a, err := UpdateArticle(paramArticle); err == nil {
		c.JSON(http.StatusOK, gin.H{
			"article": a,
		})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error" : err.Error(),
		})
	}
}