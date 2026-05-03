// creazione server Go per il backend dell'applicazione.

package main

import (
	"fmt"
	"net/http"
)

// funzione per gestire la rotta principale del server, che risponde con un messaggio di conferma quando viene visitata.
// equivalente di una Route::get('/') in Laravel.)
func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Backend Go attivo. API funzionante.")
}

// funzione main per avviare il server HTTP e Go. (linguaggio compilato->avvio con il comando go run main.go)
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

	// avvio del server HTTP sulla porta 8080.
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Errore avvio server:", err)
	}
}
