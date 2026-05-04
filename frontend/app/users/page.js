"use client";

import { useEffect, useState } from "react";

export default function UsersPage() {

  const [users, setUsers] = useState([]);

  useEffect(() => {

    const token = localStorage.getItem("access_token");

    // se non loggato torna al login
    if (!token) {
      window.location.href = "/";
      return;
    }

    fetchUsers(token);

  }, []);

  const fetchUsers = async (token) => {

    try {

      const response = await fetch("http://localhost:8080/users/read", {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      // token scaduto/non valido
      if (!response.ok) {
        localStorage.removeItem("access_token");
        localStorage.removeItem("refresh_token");

        window.location.href = "/";
        return;
      }

      const data = await response.json();

      setUsers(data);

    } catch (error) {

      console.error("Errore fetch utenti:", error);

    }
  };

  const logout = () => {

    localStorage.removeItem("access_token");
    localStorage.removeItem("refresh_token");

    window.location.href = "/";

  };

  return (

    <main className="container py-5">

      <div className="d-flex justify-content-between align-items-center mb-5">

        <div>

          <h1 className="display-5 fw-bold">
            Gestione Utenti
          </h1>

          <p className="text-muted mb-0">
            Lista utenti registrati nel gestionale
          </p>

        </div>

        <button
          className="btn btn-danger"
          onClick={logout}
        >
          Logout
        </button>

      </div>

      <div className="table-responsive">

        <table className="table table-striped align-middle">

          <thead>

            <tr>
              <th>Username</th>
              <th>Email</th>
              <th>Nome</th>
              <th>Cognome</th>
              <th>Stato</th>
            </tr>

          </thead>

          <tbody>

            {users.map((user) => (

              <tr key={user.id}>

                <td>{user.username}</td>
                <td>{user.email}</td>
                <td>{user.first_name}</td>
                <td>{user.last_name}</td>
                <td>{user.account_status}</td>

              </tr>

            ))}

          </tbody>

        </table>

      </div>

    </main>
  );
}