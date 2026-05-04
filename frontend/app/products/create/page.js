"use client";

import { useEffect, useState } from "react";

export default function CreateProductPage() {

const [productName, setProductName] = useState("");
const [description, setDescription] = useState("");
const [quantityAvailable, setQuantityAvailable] = useState("");
const [unitPrice, setUnitPrice] = useState("");
const [minimumStockThreshold, setMinimumStockThreshold] = useState("");

const [productTypes, setProductTypes] = useState([]);
const [selectedProductType, setSelectedProductType] = useState("");

const [successMessage, setSuccessMessage] = useState("");
const [errorMessage, setErrorMessage] = useState("");

useEffect(() => {

  fetchProductTypes();

}, []);

const fetchProductTypes = async () => {

  const token = localStorage.getItem("access_token");

  try {

    const response = await fetch(
      "http://localhost:8080/product-types/read",
      {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      }
    );

    const data = await response.json();

    setProductTypes(data);

  } catch (error) {

    console.error("Errore fetch product types:", error);

  }
};

const handleCreateProduct = async (e) => {

  e.preventDefault();

  setSuccessMessage("");
  setErrorMessage("");

  const token = localStorage.getItem("access_token");

  // se token assente torna al login
  if (!token) {

    window.location.href = "/";
    return;

  }

  try {

    const response = await fetch(
      "http://localhost:8080/products",
      {
        method: "POST",

        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },

        body: JSON.stringify({

          product_name: productName,
          description: description,
          quantity_available: parseInt(quantityAvailable),
          unit_price: parseFloat(unitPrice),
          minimum_stock_threshold: parseInt(minimumStockThreshold),

          product_type_id: selectedProductType,

          // temporaneo
          created_by: "1e5de66f-7f38-4373-87ef-6e7c88f92c84",

        }),
      }
    );

    if (!response.ok) {

      setErrorMessage("Errore creazione prodotto");
      return;

    }

    setSuccessMessage("Prodotto creato correttamente");

    setTimeout(() => {

      window.location.href = "/products";

    }, 4500);

  } catch (error) {

    console.error("Errore creazione prodotto:", error);

    setErrorMessage("Errore connessione server");

  }
};

return (

<main className="container py-5">

{/* HEADER */}

<div className="mb-5">

  <h1 className="display-5 fw-bold">
    Nuovo Prodotto
  </h1>

  <p className="text-muted">
    Inserimento nuovo prodotto nel magazzino
  </p>

</div>

{/* FORM */}

<div className="row">

  <div className="col-12 col-lg-8">

    <article className="it-card rounded shadow-sm border">

      <div className="it-card-body p-4">

        <form onSubmit={handleCreateProduct}>

          {/* NOTIFICA SUCCESSO */}

       

            {successMessage && (

              <div className="alert alert-success mb-4">

                {successMessage}

              </div>

            )}

          {/* NOTIFICA ERRORE */}

            {errorMessage && (

              <div className="alert alert-danger mb-4">

                {errorMessage}

              </div>

            )}

         
          {/* NOME PRODOTTO */}

          <div className="mb-4">

            <label className="form-label">
              Nome prodotto
            </label>

            <input
              type="text"
              className="form-control"
              value={productName}
              onChange={(e) => setProductName(e.target.value)}
              required
            />

          </div>

          {/* DESCRIZIONE */}

          <div className="mb-4">

            <label className="form-label">
              Descrizione
            </label>

            <textarea
              className="form-control"
              rows="3"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              required
            />

          </div>

          {/* TIPO PRODOTTO */}

          <div className="mb-4">

            <label className="form-label">
              Tipo prodotto
            </label>

            <select
              className="form-select bg-light border"
              style={{ cursor: "pointer" }}
              value={selectedProductType}
              onChange={(e) => setSelectedProductType(e.target.value)}
              required
            >

              <option value="">
                Seleziona tipologia
              </option>

              {productTypes.map((type) => (

                <option
                  key={type.id}
                  value={type.id}
                >
                  {type.type_name}
                </option>

              ))}

            </select>

          </div>

          {/* QUANTITA */}

          <div className="mb-4">

            <label className="form-label">
              Quantità disponibile
            </label>

            <input
              type="number"
              className="form-control"
              value={quantityAvailable}
              onChange={(e) => setQuantityAvailable(e.target.value)}
              required
            />

          </div>

          {/* PREZZO */}

          <div className="mb-4">

            <label className="form-label">
              Prezzo unitario
            </label>

            <input
              type="number"
              step="0.01"
              className="form-control"
              value={unitPrice}
              onChange={(e) => setUnitPrice(e.target.value)}
              required
            />

          </div>

          {/* SOGLIA MINIMA */}

          <div className="mb-4">

            <label className="form-label">
              Soglia minima magazzino
            </label>

            <input
              type="number"
              className="form-control"
              value={minimumStockThreshold}
              onChange={(e) => setMinimumStockThreshold(e.target.value)}
              required
            />

          </div>

          {/* BUTTON */}

          <button
            type="submit"
            className="btn btn-primary"
          >
            Crea prodotto
          </button>

        </form>

      </div>

    </article>

  </div>

</div>

</main>
);
}