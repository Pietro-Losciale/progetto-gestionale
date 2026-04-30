## Progetto web avanzato 
Applicativo web gestionale full-stack con interfaccia amministrativa, che consente la gestione completa di:
-utenti
-ruoli
-inventario del negozio. 

## Obiettivo del progetto

Realizzare una piattaforma amministrativa sicura e strutturata, con particolare attenzione a:

- organizzazione dell’architettura applicativa full-stack
- focus su sicurezza e tracciamento delle operazioni all'interno del gestionale
- progettazione del database relazionale 
- gestione dei ruoli e delle autorizzazioni


## Linguaggi e tecnologie utilizzate
-Front-end: Next.js, Bootstrap Italia.
-Logica di Back-end:Go
-Database:PostgreSQL
-Autenticazione: JWT (Access+Refresh Token)
-requisiti di sicurezza:RBAC, CSRF, utilizzo del protocollo HTTPS, gestione dei refresh token con scadenza, Validazione input (mitigazione SQL Injection e Cross Site Scripting), logging centralizzato errori. 

# Comandi avvio progetto
**Avvio back-end:**

cd backend
go run .

**Avvio front-end:**
cd frontend
npm install
npm run dev

**Creazione DB:**
PostgreSQL:
CREATE DATABASE gestionale_db;

eseguire file.sql:
\i percorso/database/init.sql

**Variabili d'ambiente**
Creare file .env in /backend
Inserire: 
JWT_SECRET=gestionale_jwt_secret


# Stato attuale
-Creata la logica backend iniziale con Go

 -Creata struttura database con tabelle. Presente file.md in /docs con elenco delle tabelle e relazioni tra esse. UUID utilizzato per le primary keys delle tabelle.

 -creata la logica di autenticazione (auth.go). Chiamata API rest dal front-end, check su:
 esistenza hash password inserita e utente inserito nel db
 verifica di eventuale soft lock 
 soft lock dopo 5 tentativi errati di inserimento password.(per reset account:agire con reset amministrativo da database)

 -Implementato endpoint POST /login con autenticazione tramite bcrypt e generazione JWT (testato tramite cUrl)
 
 -Implementata logica di logging accessi per Audit trail: tentativi riusciti, falliti, timestamp ed indirizzo IP. 

 -Aggiornata generazione JWT con implementazione di ACCESS TOKEN (15 MIN) e REFRESH TOKEN (7GG)