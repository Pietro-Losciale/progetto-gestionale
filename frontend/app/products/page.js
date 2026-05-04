"use client";

import { useEffect, useState } from "react";

export default function ProductsPage() {

  const [products, setProducts] = useState([]);

  // stato ricerca
  const [searchTerm, setSearchTerm] = useState("");

  // stato filtro tipo prodotto
  const [selectedType, setSelectedType] = useState("");

  // stato ordinamento
  const [sortBy, setSortBy] = useState("");

  // tipi prodotto
  const [productTypes, setProductTypes] = useState([]);

  useEffect(() => {

    const token = localStorage.getItem("access_token");

    // se non autenticato torna al login
    if (!token) {

      window.location.href = "/";
      return;

    }

    fetchProducts(token);
    fetchProductTypes(token);

  }, []);

  // FETCH PRODOTTI

  const fetchProducts = async (token) => {

    try {

      const response = await fetch(
        "http://localhost:8080/products/read",
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );

      // token non valido/scaduto
      if (!response.ok) {

        localStorage.removeItem("access_token");
        localStorage.removeItem("refresh_token");

        window.location.href = "/";
        return;

      }

      const data = await response.json();

      setProducts(data);

    } catch (error) {

      console.error("Errore fetch prodotti:", error);

    }
  };

  // FETCH PRODUCT TYPES

  const fetchProductTypes = async (token) => {

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

  // LOGOUT

  const logout = () => {

    localStorage.removeItem("access_token");
    localStorage.removeItem("refresh_token");

    window.location.href = "/";

  };

  // FILTRO + SEARCH + SORT

  const filteredProducts = products
    .filter((product) => {

      // filtro ricerca nome
      const matchesSearch = product.product_name
        .toLowerCase()
        .includes(searchTerm.toLowerCase());

      // filtro tipologia
      const matchesType =
        selectedType === "" ||
        product.product_type_id === selectedType;

      return matchesSearch && matchesType;

    })
    .sort((a, b) => {

      // ordinamento prezzo
      if (sortBy === "price") {
        return a.unit_price - b.unit_price;
      }

      // ordinamento quantità
      if (sortBy === "quantity") {
        return a.quantity_available - b.quantity_available;
      }

      return 0;

    });

  return (

    <main className="container py-5">

      {/* HEADER */}

      <div className="d-flex justify-content-between align-items-center mb-5">

        <div>

          <h1 className="display-5 fw-bold">
            Gestione Prodotti
          </h1>

          <p className="text-muted mb-0">
            Gestione inventario magazzino
          </p>

        </div>

        <button
          className="btn btn-danger"
          onClick={logout}
        >
          Logout
        </button>

      </div>

      {/* AZIONI */}

      <div className="row mb-4">

        {/* BUTTON NUOVO PRODOTTO */}

        <div className="col-12 col-lg-3 mb-3">

          <button
            className="btn btn-primary w-100"
            onClick={() => {
              window.location.href = "/products/create";
            }}
          >
            Nuovo prodotto
          </button>

        </div>

        {/* SEARCH BAR */}

        <div className="col-12 col-lg-3 mb-3">

          <input
            type="text"
            className="form-control"
            placeholder="Ricerca prodotto..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
          />

        </div>

        {/* FILTRO TIPO */}

        <div className="col-12 col-lg-3 mb-3">

          <select
            className="form-control"
            value={selectedType}
            onChange={(e) => setSelectedType(e.target.value)}
          >

            <option value="">
              Tutte le tipologie
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

        {/* ORDINAMENTO */}

        <div className="col-12 col-lg-3 mb-3">

          <select
            className="form-control"
            value={sortBy}
            onChange={(e) => setSortBy(e.target.value)}
          >

            <option value="">
              Ordinamento
            </option>

            <option value="price">
              Prezzo
            </option>

            <option value="quantity">
              Quantità
            </option>

          </select>

        </div>

      </div>

      {/* TABELLA PRODOTTI */}

      <div className="table-responsive">

        <table className="table table-striped align-middle">

          <thead>

            <tr>

              <th>Nome</th>
              <th>Descrizione</th>
              <th>Quantità</th>
              <th>Prezzo</th>
              <th>Soglia minima</th>
              <th>Low Stock</th>

            </tr>

          </thead>

          <tbody>

            {filteredProducts.map((product) => (

              <tr key={product.id}>

                <td>
                  {product.product_name}
                </td>

                <td>
                  {product.description}
                </td>

                <td>
                  {product.quantity_available}
                </td>

                <td>
                  € {product.unit_price}
                </td>

                <td>
                  {product.minimum_stock_threshold}
                </td>

                <td>

                  {product.quantity_available <= product.minimum_stock_threshold ? (

                    <span className="badge bg-danger">
                      ATTENZIONE
                    </span>

                  ) : (

                    <span className="badge bg-success">
                      OK
                    </span>

                  )}

                </td>

              </tr>

            ))}

          </tbody>

        </table>

      </div>

    </main>
  );
}