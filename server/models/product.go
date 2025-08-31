package models

import (
	"encoding/json"
	"server/config"

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

    return products, nil
}