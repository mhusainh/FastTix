"use client";

import React, { useState, useContext } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { AuthContext } from "../../context/AuthContext";
import { NotificationContext } from "../../context/NotificationContext";

const LoginPage: React.FC = () => {
  const { addNotification } = useContext(NotificationContext);
  const router = useRouter();
  const { login } = useContext(AuthContext);

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState<boolean>(false);

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setLoading(true);

    try {
      const response = await fetch("http://localhost:8080/api/v1/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email, password }),
      });

      const data = await response.json();

      if (!response.ok) {
        if (data.meta && data.meta.message) {
          throw new Error(data.meta.message);
        } else {
          throw new Error("Gagal masuk. Silakan coba lagi.");
        }
      }

      if (data.meta && data.meta.message) {
        addNotification(data.meta.message, "success");
      } else {
        addNotification("Berhasil login", "success");
      }

      // Cek apakah respons memiliki token di dalam objek data
      if (data.data && data.data.token) {
        const token = data.data.token;

        // Simpan token di localStorage
        localStorage.setItem("authToken", token);

        // Simpan token di context
        login(token);

        // Redirect setelah login berhasil
        router.push("/");
      } else {
        throw new Error("Respons tidak sesuai dengan yang diharapkan.");
      }
    } catch (err: any) {
      setError(err.message || "Terjadi kesalahan saat login.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-4">Login</h1>
      <form onSubmit={handleLogin} className="max-w-md mx-auto">
        {error && (
          <div className="mb-4 p-2 bg-red-200 text-red-800 rounded">
            {error}
          </div>
        )}
        <div className="mb-4">
          <label htmlFor="email" className="block text-gray-700">
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
        <div className="mb-4">
          <label htmlFor="password" className="block text-gray-700">
            Password
          </label>
          <input
            type="password"
            id="password"
            className="w-full p-2 border border-gray-300 rounded mt-1"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
            placeholder="Masukkan password Anda"
          />
        </div>
        <button
          type="submit"
          className={`w-full text-white p-2 rounded transition ${
            loading ? "bg-gray-400 cursor-not-allowed" : "bg-blue-500 hover:bg-blue-600"
          }`}
          disabled={loading}
        >
          {loading ? "Sedang Masuk..." : "Masuk"}
        </button>
      </form>

      <p className="mt-4 text-center text-gray-600">
        Belum punya akun?{" "}
        <Link href="/register" className="text-blue-500 hover:underline">
          Daftar di sini
        </Link>
      </p>

      {/* Link ke halaman reset password */}
      <p className="mt-2 text-center">
        <Link href="/request-reset-password" className="text-red-500 hover:underline">
          Lupa password?
        </Link>
      </p>
    </div>
  );
};

export default LoginPage;
