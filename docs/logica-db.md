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



