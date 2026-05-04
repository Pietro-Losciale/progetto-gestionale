"use client";

import { useState } from "react";

export default function LoginPage() {

  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");

  const handleLogin = async (e) => {

    e.preventDefault();

    setError("");

    try {

      const response = await fetch("http://localhost:8080/login", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          username,
          password,
        }),
      });

      if (!response.ok) {
        setError("Credenziali non valide");
        return;
      }

      const data = await response.json();

      // salvataggio token
      localStorage.setItem("access_token", data.access_token);
      localStorage.setItem("refresh_token", data.refresh_token);

      console.log("LOGIN OK", data);

      window.location.href = "/dashboard";

    } catch (err) {

      console.error(err);
      setError("Errore connessione backend");

    }
  };

  return (
    <main className="container py-5">

      <div className="row justify-content-center">

        <div className="col-12 col-md-6 col-lg-4">

          <div className="card shadow-sm border-0 rounded-4">

            <div className="card-body p-4">

              <h1 className="h3 text-center mb-4">
                Gestionale Magazzino
              </h1>

              <form onSubmit={handleLogin}>

                <div className="mb-3">

                  <label className="form-label">
                    Username
                  </label>

                  <input
                    type="text"
                    className="form-control"
                    value={username}
                    onChange={(e) => setUsername(e.target.value)}
                    required
                  />

                </div>

                <div className="mb-4">

                  <label className="form-label">
                    Password
                  </label>

                  <input
                    type="password"
                    className="form-control"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    required
                  />

                </div>

                {error && (
                  <div className="alert alert-danger">
                    {error}
                  </div>
                )}

                <button
                  type="submit"
                  className="btn btn-primary w-100"
                >
                  Accedi
                </button>

              </form>

            </div>

          </div>

        </div>

      </div>

    </main>
  );
}