package game

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"igropoisk_backend/internal/logger"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
)

type SearchRepository interface {
	IndexGame(ctx context.Context, game *Game) error
	DeleteGame(ctx context.Context, id int) error
	SearchGames(ctx context.Context, query string) ([]Game, error)
	SyncWith(ctx context.Context, repo Repository) error
}

type ElasticRepository struct {
	es *elasticsearch.Client
}

func NewElasticRepository(es *elasticsearch.Client) SearchRepository {
	return &ElasticRepository{es: es}
}

func (r *ElasticRepository) IndexGame(ctx context.Context, game *Game) error {
	body, _ := json.Marshal(game)
	res, err := r.es.Index(
		"games",
		bytes.NewReader(body),
		//r.es.Index.WithDocumentID(fmt.Sprint(game.ID)),
	)
	defer res.Body.Close()

	var indexRes map[string]interface{}
	json.NewDecoder(res.Body).Decode(&indexRes)
	fmt.Println(indexRes)
	logger.Logger.Info("Game indexed",
		"game_json", string(body))
	return err
}

func (r *ElasticRepository) DeleteGame(ctx context.Context, id int) error {
	_, err := r.es.Delete("games", fmt.Sprint(id))
	return err
}

func (r *ElasticRepository) SearchGames(ctx context.Context, query string) ([]Game, error) {
	q := fmt.Sprintf(`{
	  "query": {
	    "multi_match": {
	      "query": "%s",
	      "fields": ["name^3", "description"]
	    }
	  }
	}`, query)

	res, err := r.es.Search(
		r.es.Search.WithContext(ctx),
		r.es.Search.WithIndex("games"),
		r.es.Search.WithBody(bytes.NewReader([]byte(q))),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var resp struct {
		Hits struct {
			Hits []struct {
				Source Game `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return nil, err
	}

	games := make([]Game, 0, len(resp.Hits.Hits))
	for _, hit := range resp.Hits.Hits {
		games = append(games, hit.Source)
	}
	return games, nil
}

func (r *ElasticRepository) SyncWith(ctx context.Context, repo Repository) error {
	query := `{"query": {"match_all": {}}}`

	res, err := r.es.DeleteByQuery(
		[]string{"games"},
		strings.NewReader(query),
		r.es.DeleteByQuery.WithContext(context.Background()),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if err != nil {
		return err
	}
	if res.IsError() {
		return errors.New(res.String())
	}

	games, err := repo.GetAllGames(ctx)
	if err != nil {
		return err
	}

	for _, v := range games {
		err := r.IndexGame(ctx, &v)
		if err != nil {
			return err
		}
	}
	return nil
}
