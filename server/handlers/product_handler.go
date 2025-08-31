package handlers

import (
	"encoding/json"
	"net/http"
	"server/models"
	"server/utils"
)

func GetAllProducts(w http.ResponseWriter, r *http.Request) {

	products, err := models.GetAllProducts()
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, "Products not found")
        return
	}

	// Return the products as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}