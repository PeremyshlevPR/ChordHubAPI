package opensearch

import (
	"bytes"
	"chords_app/internal/config"
	"chords_app/internal/models"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/opensearch-project/opensearch-go"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
)

func CreateOpenSearchClient(config *config.Opensearch) (*opensearch.Client, error) {
	return opensearch.NewClient(opensearch.Config{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Addresses: config.Addresses,
		Username:  config.Username,
		Password:  config.Password,
	})
}

type OpenSearchAdapter struct {
	client    *opensearch.Client
	indexName string
}

func NewOpenSearchAdapter(client *opensearch.Client, indexName string) *OpenSearchAdapter {
	return &OpenSearchAdapter{client, indexName}
}

type QueryResult struct {
	Score   float32
	ObjType string
	ObjId   uint
}

func (oa *OpenSearchAdapter) IndexSong(song *models.Song) error {
	body := map[string]interface{}{
		"id":          song.ID,
		"title":       song.Title,
		"description": song.Description,
		"content":     song.Content,
		"type":        "song",
	}
	doc_id := "song_" + strconv.FormatUint(uint64(song.ID), 10)
	return oa.indexDocument(doc_id, body)
}

func (oa *OpenSearchAdapter) IndexArtist(artist *models.Artist) error {
	body := map[string]interface{}{
		"id":          artist.ID,
		"name":        artist.Name,
		"description": artist.Description,
		"type":        "artist",
	}
	doc_id := "artist_" + strconv.FormatUint(uint64(artist.ID), 10)
	return oa.indexDocument(doc_id, body)
}

func (oa *OpenSearchAdapter) Search(query string) ([]QueryResult, error) {
	searchBody := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query": query,
				"fields": []string{
					"title^3",
					"name^3",
					"content",
					"description",
				},
				"fuzziness": "AUTO",
			},
		},
		"sort": []interface{}{
			map[string]interface{}{
				"_score": map[string]interface{}{
					"order": "desc",
				},
			},
		},
	}

	bodyBytes, err := json.Marshal(searchBody)
	if err != nil {
		return nil, err
	}

	searchRequest := opensearchapi.SearchRequest{
		Index: []string{oa.indexName},
		Body:  bytes.NewReader(bodyBytes),
	}

	response, err := searchRequest.Do(context.Background(), oa.client)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var searchResults struct {
		Hits struct {
			Hits []struct {
				Score  float32 `json:"_score"`
				Source struct {
					ID   uint   `json:"id"`
					Type string `json:"type"`
				} `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(response.Body).Decode(&searchResults); err != nil {
		return nil, err
	}

	results := make([]QueryResult, 0, len(searchResults.Hits.Hits))
	for _, hit := range searchResults.Hits.Hits {
		results = append(results, QueryResult{
			Score:   hit.Score,
			ObjType: hit.Source.Type,
			ObjId:   hit.Source.ID,
		})
	}

	return results, nil
}

func (oa *OpenSearchAdapter) indexDocument(id string, body map[string]interface{}) error {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return err
	}

	request := opensearchapi.IndexRequest{
		Index:      oa.indexName,
		DocumentID: id,
		Body:       bytes.NewReader(bodyBytes),
		Refresh:    "true",
	}

	response, err := request.Do(context.Background(), oa.client)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.IsError() {
		slog.Error("Error indexing document ID %s: %s", id, response.String())
		return fmt.Errorf("error indexing document ID %s", id)
	}

	return nil
}
