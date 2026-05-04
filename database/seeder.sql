-- RUOLI

INSERT INTO roles (
    id,
    role_name,
    description
)
VALUES
(
    'd51a6050-b4cd-4fbc-a0ea-d467638fdc43',
    'admin',
    'Amministratore con accesso completo al sistema'
),
(
    '6599b873-9d78-4823-9b4d-c7a00ed731d4',
    'operatore',
    'Operatore gestione magazzino'
);





-- UTENTI DEMO

INSERT INTO users (
    username,
    email,
    password_hash,
    first_name,
    last_name,
    role_id
)
VALUES
(
    'admin',
    'admin@gestionale.com',
    '$2a$10$S.wIY7fVpIIQb0BI7RQGh.ZI8jCH9WXXxQjq.oHVcaB6ZbUGJk6TW',
    'Admin',
    'Sistema',
    'd51a6050-b4cd-4fbc-a0ea-d467638fdc43'
),
(
    'operatore1',
    'operatore2_updated@gestionale.com',
    '$2a$10$GHHeCLMzzRjOXh5V08RHu.lDbJhYMsl96azA7XQDyS4sFMTeRBFeG',
    'Mario',
    'Rossi',
    '6599b873-9d78-4823-9b4d-c7a00ed731d4'
),
(
    'operatore2',
    'operatore2@gestionale.com',
    '$2a$10$v0yw2kosjeti1iRZfs0fAekyi1LW.FlWqxZdfaCs5MUgiR8yfI3ey',
    'Luigi',
    'Verdi',
    '6599b873-9d78-4823-9b4d-c7a00ed731d4'
);





-- CATEGORIE PRODOTTI

INSERT INTO product_types (
    id,
    type_name,
    description
)
VALUES
(
    '61460224-5793-428e-a664-dcdc3cb3ffaa',
    'Carta',
    'Risme e prodotti cartacei'
),
(
    '189d2c97-51cd-40e1-807a-c72db8d59711',
    'Toner',
    'Toner e cartucce'
),
(
    '30de4d30-a961-4efd-bc6d-505d355907d4',
    'Buste',
    'Prodotti per imbustamento'
);