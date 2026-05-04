// creazione server Go per il backend dell'applicazione.

package main

import (
	"fmt"
	"net/http"
)

// middleware CORS per comunicazione frontend Next.js <-> backend Go
func enableCORS(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// gestione richieste preflight OPTIONS
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// funzione per gestire la rotta principale del server, che risponde con un messaggio di conferma quando viene visitata.
// equivalente di una Route::get('/') in Laravel.
func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Backend Go attivo. API funzionante.")
}

// funzione main per avviare il server HTTP e Go.
func main() {

	connectDB()

	// rotte
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/login", LoginHandler)
	http.HandleFunc("/protected", AuthMiddleware(ProtectedHandler))

	http.HandleFunc("/inventory/movement", AuthMiddleware(CreateInventoryMovementHandler))
	http.HandleFunc("/inventory/movements", AuthMiddleware(GetInventoryMovementsHandler))
	http.HandleFunc("/access-logs", AuthMiddleware(GetAccessLogsHandler))
	http.HandleFunc("/products/low-stock", AuthMiddleware(GetLowStockProductsHandler))

	http.HandleFunc("/product-types/read", AuthMiddleware(GetProductTypesHandler))

	// rotte per la gestione degli utenti (CRUD)
	http.HandleFunc("/users/create", AdminOnlyMiddleware(CreateUserHandler))
	http.HandleFunc("/users/read", AdminOnlyMiddleware(GetUsersHandler))
	http.HandleFunc("/users/update", AdminOnlyMiddleware(UpdateUserHandler))
	http.HandleFunc("/users/delete", AdminOnlyMiddleware(DeleteUserHandler))

	// rotte per la gestione dei prodotti (CRUD)
	http.HandleFunc("/products", AuthMiddleware(CreateProductHandler))
	http.HandleFunc("/products/read", AuthMiddleware(GetProductsHandler))
	http.HandleFunc("/products/update", AuthMiddleware(UpdateProductHandler))
	http.HandleFunc("/products/delete", AuthMiddleware(DeleteProductHandler))

	fmt.Println("Server avviato su http://localhost:8080")

	// avvio server con middleware CORS
	err := http.ListenAndServe(":8080", enableCORS(http.DefaultServeMux))

	if err != nil {
		fmt.Println("Errore avvio server:", err)
	}
}
