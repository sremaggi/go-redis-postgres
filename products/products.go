package products

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis"
	_ "github.com/lib/pq"
)

type productRow struct {
	ProdID      *string `json:"pId"`
	ProdEanID   *string `json:"pEId"`
	IDEstilo    *string `json:"iDE"`
	ProdDescTXT *string `json:"desc"`
	PrdCatID    *string `json:"cat"`
}

type JsonResponse struct {
	Data   []productRow `json:"data"`
	Source string       `json:"source"`
}

func GetProducts() (*JsonResponse, error) {

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "password",
		DB:       0,
	})

	pong, err := redisClient.Ping().Result()
	fmt.Println(pong, err)

	cachedProducts, err := redisClient.Get("products").Bytes()

	response := JsonResponse{}

	if err != nil {

		dbProducts, err := fetchFromDb()

		if err != nil {
			return nil, err
		}

		cachedProducts, err = json.Marshal(dbProducts)

		if err != nil {
			return nil, err
		}

		err = redisClient.Set("products", string(cachedProducts), 10*time.Minute).Err()

		if err != nil {
			return nil, err
		}

		response = JsonResponse{Data: dbProducts, Source: "PostgreSQL"}

		return &response, err
	}

	products := []productRow{}

	err = json.Unmarshal(cachedProducts, &products)

	if err != nil {
		return nil, err
	}

	response = JsonResponse{Data: products, Source: "Redis Cache"}

	return &response, nil
}

func fetchFromDb() ([]productRow, error) {

	dbUser := "ext_sebasremaggi"
	dbPassword := "password"
	dbName := "postgres"

	conString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", conString)

	if err != nil {
		return nil, err
	}

	queryString := `SELECT * FROM products`

	rows, err := db.Query(queryString)

	if err != nil {
		return nil, err
	}

	var records []productRow

	for rows.Next() {

		var p productRow

		err = rows.Scan(&p.ProdID, &p.ProdEanID, &p.ProdEanID, &p.ProdDescTXT, &p.PrdCatID)

		records = append(records, p)

		if err != nil {
			return nil, err
		}

	}

	return records, nil
}
