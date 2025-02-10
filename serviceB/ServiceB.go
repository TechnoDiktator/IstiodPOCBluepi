package serviceb


import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/yourusername/IstiodPOCBluepi"

	_ "github.com/go-sql-driver/mysql"


)




var db *sql.DB

func connectDB() {
	var err error
	dsn := "root:password@tcp(localhost:3306)/testdb"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatalf("Database is unreachable: %v", err)
	}
	fmt.Println("Connected to MySQL successfully!")
}

func getProducts(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name, price FROM products")
	if err != nil {
		http.Error(w, "Failed to fetch products", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price); err != nil {
			http.Error(w, "Failed to parse products", http.StatusInternalServerError)
			return
		}
		products = append(products, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func main() {
	connectDB()
	http.HandleFunc("/products", getProducts)
	fmt.Println("Service B is running on port 8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

