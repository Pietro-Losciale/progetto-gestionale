"use client";

import { useEffect, useState } from "react";
import Link from "next/link";

export default function DashboardPage() {

  const [usersCount, setUsersCount] = useState("--");
  const [productsCount, setProductsCount] = useState("--");
  const [lowStockCount, setLowStockCount] = useState("--");

  useEffect(() => {

    const token = localStorage.getItem("access_token");

    // se token assente torna al login
    if (!token) {
      window.location.href = "/";
      return;
    }

    fetchDashboardData(token);

  }, []);

  const fetchDashboardData = async (token) => {

    try {

      // recupero utenti
      const usersResponse = await fetch(
        "http://localhost:8080/users/read",
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );

      const usersData = await usersResponse.json();
      setUsersCount(usersData.length);

      // recupero prodotti
      const productsResponse = await fetch(
        "http://localhost:8080/products/read",
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );

      const productsData = await productsResponse.json();
      setProductsCount(productsData.length);

      // recupero low stock
      const lowStockResponse = await fetch(
        "http://localhost:8080/products/low-stock",
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );

      const lowStockData = await lowStockResponse.json();
      setLowStockCount(lowStockData.length);

    } catch (error) {

      console.error("Errore dashboard:", error);

    }
  };

  return (

    <main className="container py-5">

      {/* HEADER DASHBOARD */}

      <div className="d-flex justify-content-between align-items-center mb-5">

        <div>

          <h1 className="display-4 fw-bold mb-2">
            Dashboard Gestionale
          </h1>

          <p className="text-secondary mb-0">
            Pannello amministrativo gestione magazzino
          </p>

        </div>

        <button
          className="btn btn-danger"
          onClick={() => {

            localStorage.removeItem("access_token");
            localStorage.removeItem("refresh_token");

            window.location.href = "/";

          }}
        >
          Logout
        </button>

      </div>

      <div className="row">

{/* CARD UTENTI */}

        <div className="col-12 col-md-4 mb-4">

          <Link
            href="/users"
            className="text-decoration-none"
          >

            <article className="it-card rounded shadow-sm border h-100">

              <h3 className="it-card-title px-4 pt-4">
                Utenti
              </h3>

              <div className="it-card-body">

                <p className="display-4 fw-bold mb-2 text-dark">
                  {usersCount}
                </p>

                <p className="it-card-text text-dark">
                  Numero totale utenti registrati
                </p>

              </div>

            </article>

          </Link>

        </div>

{/* CARD PRODOTTI */}

          <div className="col-12 col-md-4 mb-4">

    <Link
      href="/products"
      className="text-decoration-none"
    >

      <article className="it-card rounded shadow-sm border h-100">

        <h3 className="it-card-title px-4 pt-4">
          Prodotti
        </h3>

        <div className="it-card-body">

          <p className="display-4 fw-bold mb-2 text-dark">
            {productsCount}
          </p>

          <p className="it-card-text text-dark">
            Prodotti presenti nel magazzino
          </p>

        </div>

      </article>

    </Link>

  </div>

{/* CARD LOW STOCK */}

       <div className="col-12 col-md-4 mb-4">

        <Link
          href="/low-stock"
          className="text-decoration-none"
        >

          <article className="it-card rounded shadow-sm border h-100">

            <h3 className="it-card-title px-4 pt-4">
              Low Stock
            </h3>

            <div className="it-card-body">

              <p className="display-4 fw-bold text-danger mb-2">
                {lowStockCount}
              </p>

              <p className="it-card-text text-dark">
                Prodotti sotto soglia minima
              </p>

            </div>

          </article>

        </Link>

      </div>

      </div>

    </main>
  );
}