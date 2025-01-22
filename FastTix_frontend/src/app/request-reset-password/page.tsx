"use client";

import React, { useState, useContext } from "react";
import { NotificationContext } from "../../context/NotificationContext";
import { useRouter } from "next/navigation";
const RequestResetPasswordPage: React.FC = () => {
  const { addNotification } = useContext(NotificationContext);
  const router = useRouter();
  const [email, setEmail] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  const handleRequestReset = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setSuccess(null);
    setLoading(true);

    try {
      const response = await fetch("http://localhost:8080/api/v1/request-reset-password", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email }),
      });

      const data = await response.json();

      if (!response.ok) {
        if (data.meta && data.meta.message) {
          throw new Error(data.meta.message);
        } else {
          throw new Error("Gagal mengirim permintaan reset password.");
        }
      }

      // Notifikasi
      if (data.meta && data.meta.message) {
        addNotification(data.meta.message, "success");
        setSuccess(data.meta.message);
      } else {
        setSuccess("Permintaan reset password berhasil dikirim.");
      }
      router.push("/reset-password");
      // Kosongkan email
      setEmail("");
    } catch (err: any) {
      setError(err.message || "Terjadi kesalahan saat request reset password.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-2xl font-bold mb-4">Reset Password</h1>
      <p className="text-gray-600 mb-6">
        Masukkan email Anda, kami akan mengirimkan token reset password ke email tersebut.
      </p>

      {/* Error & Success */}
      {error && <div className="mb-4 p-2 bg-red-200 text-red-800 rounded">{error}</div>}
      {success && <div className="mb-4 p-2 bg-green-200 text-green-800 rounded">{success}</div>}

      <form onSubmit={handleRequestReset} className="max-w-md mx-auto">
        <div className="mb-4">
          <label htmlFor="email" className="block text-gray-700 font-semibold">
            Email
          </label>
          <input
            type="email"
            id="email"
            className="w-full p-2 border border-gray-300 rounded mt-1"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
            placeholder="Masukkan email Anda"
          />
        </div>
        <button
          type="submit"
          className={`w-full text-white p-2 rounded transition ${
            loading ? "bg-gray-400 cursor-not-allowed" : "bg-blue-500 hover:bg-blue-600"
          }`}
          disabled={loading}
        >
          {loading ? "Memproses..." : "Kirim Permintaan Reset Password"}
        </button>
      </form>
    </div>
  );
};

export default RequestResetPasswordPage;
