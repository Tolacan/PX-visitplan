package search

import (
	"PX-visitplan/models"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	elastic "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type ElasticSearchRepository struct {
	client *elastic.Client
}

func NewElastic(url string) (*ElasticSearchRepository, error) {
	client, err := elastic.NewClient(elastic.Config{
		Addresses: []string{url},
	})

	if err != nil {
		return nil, err
	}

	return &ElasticSearchRepository{client: client}, nil
}

func (r *ElasticSearchRepository) Close() {
	//
}

func (repo *ElasticSearchRepository) IndexVisitPlan(ctx context.Context, visitPlan models.VisitPlan) error {
	body, er := json.Marshal(visitPlan)
	if er != nil {
		return errors.New("Error al serializar el documento" + er.Error())
	}

	req := esapi.IndexRequest{
		Index:      "visitplans",
		DocumentID: visitPlan.Uuid,
		Body:       bytes.NewReader(body),
		Refresh:    "true",
	}

	// Realizar la solicitud de indexación
	res, err := req.Do(ctx, repo.client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// Verificar la respuesta de Elasticsearch
	if res.IsError() {
		return errors.New("Error al indexar el documento:" + res.String())
	} else {
		return nil
	}
}

func (repo *ElasticSearchRepository) UpdateVisitPlan(ctx context.Context, visitPlan models.VisitPlan) error {

	body, er := json.Marshal(visitPlan)
	if er != nil {
		return errors.New("Error al serializar el documento" + er.Error())
	}

	req := esapi.UpdateRequest{
		Index:      "visitplans",
		DocumentID: visitPlan.Uuid,
		Body:       bytes.NewReader([]byte(fmt.Sprintf(`{"doc":%s}`, body))),
		Refresh:    "true",
	}

	res, err := req.Do(ctx, repo.client)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	if res.IsError() {
		return errors.New("Failed to update document" + res.String())
	}

	return nil
}

func (repo *ElasticSearchRepository) SearchVisitPlan(ctx context.Context, query string) (results []models.VisitPlan, err error) {

	var buf bytes.Buffer
	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  query,
				"fields": []string{"nombre", "ruta", "responsables"},
			},
		},
	}

	if err = json.NewEncoder(&buf).Encode(searchQuery); err != nil {
		return nil, err
	}
	res, err := repo.client.Search(
		repo.client.Search.WithContext(ctx),
		repo.client.Search.WithIndex("visitplans"),
		repo.client.Search.WithBody(&buf),
		repo.client.Search.WithTrackTotalHits(true),
	)

	if err != nil {
		return nil, err
	}
	//Cierra la conexion con elasticSearch al terminar la busqueda
	defer func() {
		if err := res.Body.Close(); err != nil {
			results = nil
		}
	}()
	//Devuelve un error nuevo deacuerdo al error que encontro elasticsearch
	if res.IsError() {
		return nil, errors.New("elasticsearch error " + res.String())
	}

	//Decodifica a json y lo guarda eRes
	var eRes map[string]interface{}
	//Decodifica
	if err := json.NewDecoder(res.Body).Decode(&eRes); err != nil {
		return nil, err
	}

	var visitsPlans []models.VisitPlan
	for _, hit := range eRes["hits"].(map[string]interface{})["hits"].([]interface{}) {
		visitplan := models.VisitPlan{}
		source := hit.(map[string]interface{})["_source"]
		marshal, err := json.Marshal(source)
		if err != nil {
			return nil, err
		}
		//Decodifica
		if err := json.Unmarshal(marshal, &visitplan); err == nil {
			visitsPlans = append(visitsPlans, visitplan)
		}
	}
	return visitsPlans, nil
}
