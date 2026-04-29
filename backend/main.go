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

	fmt.Println("Server avviato su http://localhost:8080")

	// avvio del server HTTP sulla porta 8080.
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Errore avvio server:", err)
	}
}
