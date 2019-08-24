package service

import (
	"strings"
	"elastic-talk-search/model"
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/blevesearch/bleve"
	"github.com/rs/zerolog"
)

// SearchService prodiveds an interface to search in the list of all products
type SearchService struct {
	index  bleve.Index
	Logger *zerolog.Logger
}

// Search queries all Products matching the given string
func (search *SearchService) Search(queryValue string, page int) (uint64, []model.Product) {
	var searchRequest *bleve.SearchRequest
	if queryValue == "*" {
		query := bleve.NewMatchAllQuery()
		searchRequest = bleve.NewSearchRequest(query)
	} else {
		search.Logger.Info().Msg("SearchTerm " + queryValue)
		query := bleve.NewPrefixQuery(strings.ToLower(queryValue));
		searchRequest = bleve.NewSearchRequest(query)
	}
	searchRequest.From = page * searchRequest.Size
	searchRequest.Fields = []string{"id", "title", "imageUrl", "brand", "price", "discounted"}
	searchResult, _ := search.index.Search(searchRequest)
	results := make([]model.Product, len(searchResult.Hits))
	for index, hit := range searchResult.Hits {
		results[index] = model.Product{
			ID:         hit.ID,
			Title:      hit.Fields["title"].(string),
			ImageURL:   hit.Fields["imageUrl"].(string),
			Brand:      hit.Fields["brand"].(string),
			Price:      int(hit.Fields["price"].(float64)),
			Discounted: hit.Fields["discounted"].(bool),
		}
	}
	return searchResult.Total, results
}

// LoadProducts loads all Products into the search Service. It either opens an eisting index or creates a new one
func (search *SearchService) LoadProducts() {
	if _, err := os.Stat("index.bleve"); os.IsNotExist(err) {
		search.Logger.Info().Msg("Indexing Products")
		search.indexProducts()
	}
	search.Logger.Info().Msg("Opening Index")
	index, err := bleve.Open("index.bleve")
	if err != nil {
		search.Logger.Error().
			Err(err).
			Msg("Unable to open Index")
		panic(err)
	}
	search.index = index
	count, err := index.DocCount()
	if err != nil {
		search.Logger.Error().
			Err(err).
			Msg("Unable to read Index Document Count")
		panic(err)
	}
	search.Logger.Info().Msg("Index contains " + strconv.Itoa(int(count)) + " Documents")
}

func (search *SearchService) indexProducts() {
	mapping := bleve.NewIndexMapping()
	index, err := bleve.New("index.bleve", mapping)
	if err != nil {
		panic(err)
	}
	defer index.Close()
	jsonFile, err := os.Open("products.json")
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	result := model.ProductList{}
	err = json.Unmarshal([]byte(byteValue), &result)
	if err != nil {
		panic(err)
	}
	for _, item := range result.Products {
		err = index.Index(item.ID, item)
		if err != nil {
			panic(err)
		}
	}
	search.Logger.Info().Msg("Indexed " + strconv.Itoa(len(result.Products)) + " Products")
}
