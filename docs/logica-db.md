# Logica Database

Un elenco delle tabelle presenti nel file init.sql, con le eventuali relazioni tra esse. 

note:
-Le primary keys delle tabelle sono UUID
-password salvate tramite hash (logica che verra' implementata in backend- bcrypt)
-audit trails movimenti magazzino.

## Tabella roles

Relazione one-to-many con Utenti.

Un ruolo può appartenere a più utenti, mentre ogni utente può avere un solo ruolo.

- Admin
- Operatore

Foreign Key:
users.role_id → roles.id

## Tabelle product, product_types

One-to-many -> una tipologia puo' contenere piu' prodotti, ogni prodotto ha una sola tipologia.
Foreign Key:

products.product_type_id → product_types.id

## Tabelle users,product

One to many -> ogni utente puo' creare piu' prodotti, ogni prodotto ha un solo utente creatore. 

Foreign Key:

products.created_by → users.id

# Tabelle inventory_movements, products

One-to-many-> un prodotto puo' avere piu movimenti, ogni movimento e' riconducibile a un solo prodotto.
Foreign Key:

inventory_movements.product_id → products.id

# Tabelle users, inventory_movements
One-to-many

Foreign Key:

inventory_movements.operated_by → users.id

# Tabelle users, access_logs
One-to-many
Foreign key:
access_logs.user_id → users.id

# Tabelle users, password_resets
One-to-many
Foreign key:
password_resets.user_id → users.id

# Tabella email_notifications

NO FK. Utilizziamo recipient_email per includere destinatari NON presenti negli users del database. recipient_email salva direttamente in DB l'indirizzo email del destinatario. 

--------------------------------------------------
\dt


                 List of tables
 Schema |        Name         | Type  |  Owner
--------+---------------------+-------+----------
 public | access_logs         | table | postgres
 public | email_notifications | table | postgres
 public | inventory_movements | table | postgres
 public | password_resets     | table | postgres
 public | product_types       | table | postgres
 public | products            | table | postgres
 public | roles               | table | postgres
 public | users               | table | postgres
(8 rows)