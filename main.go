package main

import (
	"os"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)
				
/*
========================
MODEL
========================
*/
type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

/*
========================
IN-MEMORY DATABASE
========================
*/
var categories = []Category{
	{ID: 1, Name: "Makanan", Description: "Produk makanan"},
	{ID: 2, Name: "Minuman", Description: "Produk minuman"},
}

/*
========================
HELPER
========================
*/
func jsonResponse(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

/*
========================
HANDLERS
========================
*/

// GET /api/categories/{id}
func getCategoryByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		jsonResponse(w, 400, map[string]string{"error": "Invalid category ID"})
		return
	}

	for _, c := range categories {
		if c.ID == id {
			jsonResponse(w, 200, c)
			return
		}
	}

	jsonResponse(w, 404, map[string]string{"error": "Category not found"})
}

// PUT /api/categories/{id}
func updateCategory(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		jsonResponse(w, 400, map[string]string{"error": "Invalid category ID"})
		return
	}

	var updated Category
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		jsonResponse(w, 400, map[string]string{"error": "Invalid request body"})
		return
	}

	for i := range categories {
		if categories[i].ID == id {
			updated.ID = id
			categories[i] = updated
			jsonResponse(w, 200, updated)
			return
		}
	}

	jsonResponse(w, 404, map[string]string{"error": "Category not found"})
}

// DELETE /api/categories/{id}
func deleteCategory(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		jsonResponse(w, 400, map[string]string{"error": "Invalid category ID"})
		return
	}

	for i, c := range categories {
		if c.ID == id {
			categories = append(categories[:i], categories[i+1:]...)
			jsonResponse(w, 200, map[string]string{"message": "Category deleted"})
			return
		}
	}

	jsonResponse(w, 404, map[string]string{"error": "Category not found"})
}

/*
========================
MAIN
========================
*/
func main() {

	// Health check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 200, map[string]string{
			"status":  "OK",
			"message": "API Running",
		})
	})

	/*
	========================
	/api/categories
	========================
	GET  -> get all
	POST -> create
	*/
	http.HandleFunc("/api/categories", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {

		case "GET":
			jsonResponse(w, 200, categories)

		case "POST":
			var newCategory Category
			if err := json.NewDecoder(r.Body).Decode(&newCategory); err != nil {
				jsonResponse(w, 400, map[string]string{"error": "Invalid request body"})
				return
			}

			newCategory.ID = len(categories) + 1
			categories = append(categories, newCategory)

			jsonResponse(w, 201, newCategory)

		default:
			jsonResponse(w, 405, map[string]string{"error": "Method not allowed"})
		}
	})

	/*
	========================
	/api/categories/{id}
	========================
	GET    -> detail
	PUT    -> update
	DELETE -> delete
	*/
	http.HandleFunc("/api/categories/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {

		case "GET":
			getCategoryByID(w, r)

		case "PUT":
			updateCategory(w, r)

		case "DELETE":
			deleteCategory(w, r)

		default:
			jsonResponse(w, 405, map[string]string{"error": "Method not allowed"})
		}
	})

	/ ambil port dari environment variable
port := os.Getenv("PORT")
if port == "" {
    port = "8080" // fallback untuk local
}

fmt.Println("🚀 Server running at port " + port)
http.ListenAndServe(":" + port, nil)
}
