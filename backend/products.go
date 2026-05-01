package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Product struct {
	ID                    string  `json:"id"`
	ProductName           string  `json:"product_name"`
	Description           string  `json:"description"`
	QuantityAvailable     int     `json:"quantity_available"`
	UnitPrice             float64 `json:"unit_price"`
	MinimumStockThreshold int     `json:"minimum_stock_threshold"`
	ProductTypeID         string  `json:"product_type_id"`
	CreatedBy             string  `json:"created_by"`
}

// struttura movimenti del magazzino
type InventoryMovement struct {
	ProductID    string `json:"product_id"`
	MovementType string `json:"movement_type"`
	Quantity     int    `json:"quantity"`
	OperatedBy   string `json:"operated_by"`
	Notes        string `json:"notes"`
}

// FUNZIONI CRUD

//CREATE

func CreateProductHandler(w http.ResponseWriter, r *http.Request) {

	// accetta solo POST
	if r.Method != http.MethodPost {
		http.Error(w, "Metodo non consentito", http.StatusMethodNotAllowed)
		return
	}

	var product Product

	// decode JSON ricevuto
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, "Errore lettura dati prodotto", http.StatusBadRequest)
		return
	}

	// query insert prodotto nel DB, con i campi necessari per la creazione di un nuovo prodotto. (equivalente di un Product::create() in Laravel)
	query := `
		INSERT INTO products (
			product_name,
			description,
			quantity_available,
			unit_price,
			minimum_stock_threshold,
			product_type_id,
			created_by
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err = DB.Exec(
		query,
		product.ProductName,
		product.Description,
		product.QuantityAvailable,
		product.UnitPrice,
		product.MinimumStockThreshold,
		product.ProductTypeID,
		product.CreatedBy,
	)

	if err != nil {
		http.Error(w, "Errore creazione prodotto", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Prodotto creato correttamente"))
}

//READ

func GetProductsHandler(w http.ResponseWriter, r *http.Request) {

	// accetta solo GET
	if r.Method != http.MethodGet {
		http.Error(w, "Metodo non consentito", http.StatusMethodNotAllowed)
		return
	}

	// query database
	rows, err := DB.Query(`
		SELECT
			id,
			product_name,
			description,
			quantity_available,
			unit_price,
			minimum_stock_threshold,
			product_type_id,
			created_by
		FROM products
		WHERE deleted_at IS NULL
	`)

	if err != nil {
		http.Error(w, "Errore recupero prodotti", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []Product

	// loop sulle righe
	for rows.Next() {

		var product Product

		err := rows.Scan(
			&product.ID,
			&product.ProductName,
			&product.Description,
			&product.QuantityAvailable,
			&product.UnitPrice,
			&product.MinimumStockThreshold,
			&product.ProductTypeID,
			&product.CreatedBy,
		)

		if err != nil {
			fmt.Println("Errore scansione prodotto:", err)
			continue
		}

		products = append(products, product)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// UPDATE

func UpdateProductHandler(w http.ResponseWriter, r *http.Request) {

	// accetta solo PUT
	if r.Method != http.MethodPut {
		http.Error(w, "Metodo non consentito", http.StatusMethodNotAllowed)
		return
	}

	var product Product

	// decode JSON
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, "Errore lettura dati prodotto", http.StatusBadRequest)
		return
	}

	// query update
	query := `
		UPDATE products
		SET
			product_name = $1,
			description = $2,
			quantity_available = $3,
			unit_price = $4,
			minimum_stock_threshold = $5,
			product_type_id = $6,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $7
	`

	_, err = DB.Exec(
		query,
		product.ProductName,
		product.Description,
		product.QuantityAvailable,
		product.UnitPrice,
		product.MinimumStockThreshold,
		product.ProductTypeID,
		product.ID,
	)

	if err != nil {
		http.Error(w, "Errore aggiornamento prodotto", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Prodotto aggiornato correttamente"))
}

// DELETE (soft delete)
func DeleteProductHandler(w http.ResponseWriter, r *http.Request) {

	// accetta solo DELETE
	if r.Method != http.MethodDelete {
		http.Error(w, "Metodo non consentito", http.StatusMethodNotAllowed)
		return
	}

	var product Product

	// decode JSON
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, "Errore lettura ID prodotto", http.StatusBadRequest)
		return
	}

	// soft delete
	query := `
		UPDATE products
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	_, err = DB.Exec(query, product.ID)

	if err != nil {
		http.Error(w, "Errore eliminazione prodotto", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Prodotto eliminato correttamente"))
}

//funzione gestione dei movimenti in magazzino

func CreateInventoryMovementHandler(w http.ResponseWriter, r *http.Request) {

	// accetta solo POST
	if r.Method != http.MethodPost {
		http.Error(w, "Metodo non consentito", http.StatusMethodNotAllowed)
		return
	}

	var movement InventoryMovement

	// decode JSON
	err := json.NewDecoder(r.Body).Decode(&movement)
	if err != nil {
		http.Error(w, "Errore lettura movimento", http.StatusBadRequest)
		return
	}

	// salvataggio movimento nel DB
	query := `
		INSERT INTO inventory_movements (
			product_id,
			movement_type,
			quantity,
			operated_by,
			notes
		)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err = DB.Exec(
		query,
		movement.ProductID,
		movement.MovementType,
		movement.Quantity,
		movement.OperatedBy,
		movement.Notes,
	)

	if err != nil {
		http.Error(w, "Errore creazione movimento", http.StatusInternalServerError)
		return
	}

	// aggiorna quantità prodotto
	if movement.MovementType == "load" {

		_, err = DB.Exec(`
			UPDATE products
			SET quantity_available = quantity_available + $1
			WHERE id = $2
		`, movement.Quantity, movement.ProductID)

	} else if movement.MovementType == "unload" {

		_, err = DB.Exec(`
			UPDATE products
			SET quantity_available = quantity_available - $1
			WHERE id = $2
		`, movement.Quantity, movement.ProductID)
	}

	if err != nil {
		http.Error(w, "Errore aggiornamento quantità prodotto", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Movimento inventario registrato correttamente"))
}

// storico movimenti magazzino
type InventoryMovementResponse struct {
	ID             string `json:"id"`
	ProductName    string `json:"product_name"`
	MovementType   string `json:"movement_type"`
	Quantity       int    `json:"quantity"`
	MovementDate   string `json:"movement_date"`
	OperatedByName string `json:"operated_by_name"`
	Notes          string `json:"notes"`
}

// GET INVENTORY MOVEMENTS

func GetInventoryMovementsHandler(w http.ResponseWriter, r *http.Request) {

	// accetta solo GET
	if r.Method != http.MethodGet {
		http.Error(w, "Metodo non consentito", http.StatusMethodNotAllowed)
		return
	}

	// query con JOIN
	rows, err := DB.Query(`
		SELECT
			im.id,
			p.product_name,
			im.movement_type,
			im.quantity,
			im.movement_date,
			u.username,
			im.notes
		FROM inventory_movements im
		JOIN products p
			ON im.product_id = p.id
		JOIN users u
			ON im.operated_by = u.id
		ORDER BY im.movement_date DESC
	`)

	if err != nil {
		http.Error(w, "Errore recupero movimenti inventario", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var movements []InventoryMovementResponse

	// loop risultati query
	for rows.Next() {

		var movement InventoryMovementResponse

		err := rows.Scan(
			&movement.ID,
			&movement.ProductName,
			&movement.MovementType,
			&movement.Quantity,
			&movement.MovementDate,
			&movement.OperatedByName,
			&movement.Notes,
		)

		if err != nil {
			fmt.Println("Errore scansione movimento:", err)
			continue
		}

		movements = append(movements, movement)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movements)
}
