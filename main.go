package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"  // or the Docker service name if running in another container
	port     = 5432         // default PostgreSQL port
	username = "test_user"  // as defined in docker-compose.yml
	password = "S3cret"     // as defined in docker-compose.yml
	dbname   = "mydatabase" // as defined in docker-compose.yml
)

var db *sql.DB

type Product struct {
	ID    int
	Name  string
	Price int
}

func main() {
	// Connection string
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, username, password, dbname)

	// Open a connection
	sdb, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	db = sdb

	defer db.Close()
	// Check the connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected!")

	err = createProduct(&Product{Name: "Go product", Price: 222})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Create Success")
}

func createProduct(product *Product) error {
	_, err := db.Exec(
		"INSERT INTO public.products(name, price) VALUES ($1, $2);",
		product.Name,
		product.Price,
	)

	return err
}

func getProduct(id int) (Product, error) {
	var p Product
	row := db.QueryRow(
		"SELECT id, name, price FROM products WHERE id=$1;",
		id,
	)

	err := row.Scan(&p.ID, &p.Name, &p.Price)

	if err != nil {
		return Product{}, err
	}

	return p, nil
}

func getProducts() ([]Product, error) {
	rows, err := db.Query("SELECT id, name, price from products;")

	if err != nil {
		return nil, err
	}

	var products []Product

	for rows.Next() {
		var p Product
		err := rows.Scan(&p.ID, &p.Name, &p.Price)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil

}

func updateProduct(id int, product *Product) (Product, error) {
	var p Product
	row := db.QueryRow(
		"UPDATE public.products SET name=$1, price=$2 WHERE id=$3 RETURNING id, name, price;",
		product.Name,
		product.Price,
		id,
	)

	err := row.Scan(&p.ID, &p.Name, &p.Price)

	if err != nil {
		return Product{}, err
	}

	return p, nil
}

func deleteProduct(id int) error {
	_, err := db.Exec(
		"DELETE FROM products WHERE id=$1;",
		id,
	)

	return err
}
