-- estensione necessaria per generare UUID (Universally Unique Identifier) automaticamente
CREATE EXTENSION IF NOT EXISTS "pgcrypto";




-- TABELLA RUOLI

CREATE TABLE roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    role_name VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);




-- TABELLA UTENTI

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,

    -- sicurezza password - memorizzata in hash nel db
    password_hash TEXT NOT NULL,

    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,

    birth_date DATE,

    role_id UUID NOT NULL,

    failed_login_attempts INT DEFAULT 0,

    account_status VARCHAR(20) DEFAULT 'active'
        CHECK (account_status IN ('active', 'blocked')),

    last_login TIMESTAMP,

    -- crud timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    -- campo per soft delete, permette di marcare un utente come "cancellato"
    -- senza rimuoverlo fisicamente dal database
    deleted_at TIMESTAMP NULL,

    CONSTRAINT fk_user_role
        FOREIGN KEY (role_id)
        REFERENCES roles(id)
);




-- TABELLA TIPI PRODOTTO

CREATE TABLE product_types (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type_name VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);




-- TABELLA PRODOTTI

CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    product_name VARCHAR(150) NOT NULL,
    description TEXT,

    quantity_available INT NOT NULL DEFAULT 0,
    unit_price NUMERIC(10,2) NOT NULL,

    minimum_stock_threshold INT NOT NULL DEFAULT 0,

    inserted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    deleted_at TIMESTAMP NULL,

    product_type_id UUID NOT NULL,
    created_by UUID NOT NULL,

    CONSTRAINT fk_product_type
        FOREIGN KEY (product_type_id)
        REFERENCES product_types(id),

    CONSTRAINT fk_product_creator
        FOREIGN KEY (created_by)
        REFERENCES users(id)
);




-- TABELLA MOVIMENTI MAGAZZINO->audit trail per tracciare le operazioni sui prodotti.

CREATE TABLE inventory_movements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    product_id UUID NOT NULL,

    movement_type VARCHAR(20) NOT NULL
        CHECK (movement_type IN ('load', 'unload')),

    quantity INT NOT NULL,

    movement_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    operated_by UUID NOT NULL,

    notes TEXT,

    CONSTRAINT fk_inventory_product
        FOREIGN KEY (product_id)
        REFERENCES products(id),

    CONSTRAINT fk_inventory_user
        FOREIGN KEY (operated_by)
        REFERENCES users(id)
);