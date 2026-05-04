"use client";

import { useEffect, useState } from "react";

export default function LowStockPage() {

  const [products, setProducts] = useState([]);

  useEffect(() => {

    const token = localStorage.getItem("access_token");

    // se non autenticato torna al login
    if (!token) {
      window.location.href = "/";
      return;
    }

    fetchLowStockProducts(token);

  }, []);

  const fetchLowStockProducts = async (token) => {

    try {

      const response = await fetch(
        "http://localhost:8080/products/low-stock",
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

      console.error("Errore fetch low stock:", error);

    }
  };

  const logout = () => {

    localStorage.removeItem("access_token");
    localStorage.removeItem("refresh_token");

    window.location.href = "/";

  };

  return (

    <main className="container py-5">

      {/* HEADER */}

      <div className="d-flex justify-content-between align-items-center mb-5">

        <div>

          <h1 className="display-5 fw-bold text-danger">
            Prodotti Low Stock
          </h1>

          <p className="text-muted mb-0">
            Prodotti sotto la soglia minima di magazzino
          </p>

        </div>

        <button
          className="btn btn-danger"
          onClick={logout}
        >
          Logout
        </button>

      </div>

      {/* TABELLA LOW STOCK */}

      <div className="table-responsive">

        <table className="table table-striped align-middle">

          <thead>

            <tr>
              <th>Nome</th>
              <th>Descrizione</th>
              <th>Quantità</th>
              <th>Soglia Minima</th>
              <th>Stato</th>
            </tr>

          </thead>

          <tbody>

            {products.map((product) => (

              <tr key={product.id}>

                <td>{product.product_name}</td>

                <td>{product.description}</td>

                <td className="fw-bold text-danger">
                  {product.quantity_available}
                </td>

                <td>
                  {product.minimum_stock_threshold}
                </td>

                <td>

                  <span className="badge bg-danger">
                    LOW STOCK
                  </span>

                </td>

              </tr>

            ))}

          </tbody>

        </table>

      </div>

    </main>
  );
}