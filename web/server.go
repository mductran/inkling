package web

import (
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	query "search/internal/mango/aggregate"
	db "search/internal/mango/db"
	"search/internal/mango/hash"
	"sync"
)

var lock sync.Mutex

type QueryResult struct {
	PHash   string               `json:"phash,omitempty" xml:"phash" :"p_hash"`
	Results *map[string][]string `json:"results,omitempty" xml:"results" :"results"`
}

type HelloWorldResponse struct {
	Message string `json:"message,omitempty" xml:"message"`
}

func SearchHandler(c echo.Context, db *mongo.Database) error {
	lock.Lock()
	defer lock.Unlock()

	err := c.Request().ParseForm()
	if err != nil {
		return err
	}

	if c.Request().PostForm.Has("image-url") {
		url := c.FormValue("image-url")

		mat, err := hash.ReadImageFromURL(url)
		if err != nil {
			return err
		}

		output := make(map[string][]string)
		phash := hash.Phash(mat)

		resultSet, err := query.Query(db, phash)

		if err != nil {
			return err
		}
		for k, v := range *resultSet {
			output[k] = v
		}

		return c.JSON(http.StatusOK, QueryResult{phash, &output})
	} else {

		return c.JSON(http.StatusOK, "cannot read form data")
	}
}

func Serve() {
	client, err := db.Connect()
	if err != nil {
		panic(err)
	}
	// TODO: use env. var to choose database based on dev or production environment
	database := client.Database("madokami")

	server := echo.New()
	server.POST("/search", func(c echo.Context) error {
		if err := SearchHandler(c, database); err != nil {
			return err
		}
		return nil
	})
	server.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, &HelloWorldResponse{"hello client"})
	})
	server.Logger.Fatal(server.Start("localhost:5000"))
}
