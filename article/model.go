package article

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type Article struct {
	Id      uint64 `uri:"id" binding:"required"`
	Author  string `form:"author"`
	Content string `form:"content"`
	RegDate time.Time
	ModDate time.Time
}

var articles []Article
var seededRand *rand.Rand

func init() {
	seededRand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
}

func RetrieveArticles(author string) (results []Article) {
	for _, article := range articles {
		if article.Author == author {
			results = append(results, article)
		}
	}
	return
}

func RegisterArticle(article Article) (Article, error) {
	article.Id = seededRand.Uint64()
	article.RegDate = time.Now()
	article.ModDate = time.Now()
	articles = append(articles, article)
	return article, nil
}

func UpdateArticle(paramArticle Article) (Article, error) {
	for i, a := range articles {
		if a.Id == paramArticle.Id {
			paramArticle.RegDate = a.RegDate
			paramArticle.ModDate = time.Now()
			articles[i] = paramArticle
			return paramArticle, nil
		}
	}

	return paramArticle, errors.New(fmt.Sprintf("Cannot find article id %v", paramArticle.Id))
}
