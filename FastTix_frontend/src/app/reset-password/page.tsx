"use client";

import React, { useState, useContext } from "react";
import { NotificationContext } from "../../context/NotificationContext";

const ResetPasswordPage: React.FC = () => {
  const { addNotification } = useContext(NotificationContext);

  const [token, setToken] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  const handleResetPassword = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setSuccess(null);
    setLoading(true);
  
    try {
      // Pakai POST, bukan PUT
      const response = await fetch(`http://localhost:8080/api/v1/reset-password/${token}`, {
        method: "POST", // ganti jadi POST
        headers: { "Content-Type": "application/json" },
        // Sesuaikan field JSON dengan yang dipakai di Postman:
        body: JSON.stringify({ password: newPassword }),
      });
  
      const data = await response.json();
  
      if (!response.ok) {
        if (data.meta && data.meta.message) {
          throw new Error(data.meta.message);
        } else {
          throw new Error("Gagal mengatur password baru.");
        }
      }
  
      if (data.meta?.message) {
        setSuccess(data.meta.message);
      } else {
        setSuccess("Password berhasil direset.");
      }
  
      setToken("");
      setNewPassword("");
    } catch (err: any) {
      setError(err.message || "Terjadi kesalahan saat reset password.");
    } finally {
      setLoading(false);
    }
  };
  

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-2xl font-bold mb-4">Atur Password Baru</h1>
      <p className="text-gray-600 mb-6">
        Masukkan token reset password yang Anda dapatkan melalui email, serta password baru Anda.
      </p>

      {/* Error & Success */}
      {error && <div className="mb-4 p-2 bg-red-200 text-red-800 rounded">{error}</div>}
      {success && <div className="mb-4 p-2 bg-green-200 text-green-800 rounded">{success}</div>}

      <form onSubmit={handleResetPassword} className="max-w-md mx-auto">
        <div className="mb-4">
          <label htmlFor="token" className="block text-gray-700 font-semibold">
            Token Reset
          </label>
          <input
            type="text"
            id="token"
            className="w-full p-2 border border-gray-300 rounded mt-1"
            value={token}
            onChange={(e) => setToken(e.target.value)}
            required
            placeholder="Masukkan token yang dikirim ke email Anda"
          />
        </div>

        <div className="mb-4">
          <label htmlFor="new_password" className="block text-gray-700 font-semibold">
            Password Baru
          </label>
          <input
            type="password"
            id="new_password"
            className="w-full p-2 border border-gray-300 rounded mt-1"
            value={newPassword}
            onChange={(e) => setNewPassword(e.target.value)}
            required
            placeholder="Masukkan password baru"
          />
        </div>

        <button
          type="submit"
          className={`w-full text-white p-2 rounded transition ${
            loading ? "bg-gray-400 cursor-not-allowed" : "bg-blue-500 hover:bg-blue-600"
          }`}
          disabled={loading}
        >
          {loading ? "Memproses..." : "Reset Password"}
        </button>
      </form>
    </div>
  );
};

export default ResetPasswordPage;
