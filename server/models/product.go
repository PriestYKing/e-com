package models

import (
	"context"
	"encoding/json"
	"server/cache"
	"server/config"
	"time"

	"github.com/lib/pq"
)

type Product struct{
	ID int `json:"id"`
	Name string `json:"name"`
	ShortDescription string `json:"short_description"`
	Description string `json:"description"`
	Price float64 `json:"price"`
	Sizes []string `json:"sizes"`
	Colors []string `json:"colors"`
	Images map[string]string `json:"images"` // Key-value pairs for color/image path
}

func GetAllProducts() ([]Product, error) {
     ctx := context.Background()
      var cachedProducts []Product
    if found, err := cache.GetCachedProducts(ctx, &cachedProducts); err == nil && found {
        return cachedProducts, nil
    }
    var products []Product
    rows, err := config.DB.Query("SELECT id, name, short_description, description, price, sizes, colors, images FROM products")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var p Product
        var imagesRaw []byte
        if err := rows.Scan(&p.ID, &p.Name, &p.ShortDescription, &p.Description, &p.Price, pq.Array(&p.Sizes), pq.Array(&p.Colors), &imagesRaw); err != nil {
            return nil, err
        }
        if err := json.Unmarshal(imagesRaw, &p.Images); err != nil {
            return nil, err
        }
        products = append(products, p)
    }
go func() {
        cacheCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()
        cache.CacheProducts(cacheCtx, products)
    }()
    return products, nil
}