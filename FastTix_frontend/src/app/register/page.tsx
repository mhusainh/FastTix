// src/app/register/page.tsx

"use client"; // Tandai ini sebagai client component

import React, { useState , useContext} from "react";
import Link from "next/link";
import { NotificationContext } from "../../context/NotificationContext";

const RegisterPage: React.FC = () => {
   const { addNotification } = useContext(NotificationContext);
  const [fullName, setFullName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [gender, setGender] = useState(""); // Tambahkan state untuk gender
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);
  const [loading, setLoading] = useState<boolean>(false);

  const handleRegister = async (e: React.FormEvent) => {
    e.preventDefault();
  
    // Reset error dan success messages
    setError(null);
    setSuccess(null);
  
    // Validasi sederhana
    if (password !== confirmPassword) {
      setError("Password dan konfirmasi password tidak cocok.");
      return;
    }
  
    if (password.length < 6) {
      setError("Password harus terdiri dari minimal 6 karakter.");
      return;
    }
  
    if (gender !== "male" && gender !== "female") {
      setError("Silakan pilih jenis kelamin.");
      return;
    }
  
    setLoading(true);
  
    try {
      const response = await fetch("http://localhost:8080/api/v1/register", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          full_name: fullName, // Sesuaikan dengan nama field backend
          email,
          password,
          gender,
        }),
      });
  
      const data = await response.json();
  
      if (!response.ok) {
        // Ambil pesan error dari backend jika ada
        if (data.meta && data.meta.message) {
          throw new Error(data.meta.message);
        } else {
          throw new Error("Gagal Register. Silakan coba lagi.");
        }
      }
      if (data.meta && data.meta.message) {
        addNotification(data.meta.message, "success");
      } else {
        addNotification("Berhasil Register", "success");
      }
      setSuccess("Registrasi berhasil! Silakan masuk.");
      // Reset form
      setFullName("");
      setEmail("");
      setPassword("");
      setConfirmPassword("");
      setGender("");
    } catch (err: any) {
      setError(err.message || "Terjadi kesalahan saat registrasi.");
    } finally {
      setLoading(false);
    }
  };
  

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-4">Register</h1>
      <form onSubmit={handleRegister} className="max-w-md mx-auto">
        {error && (
          <div className="mb-4 p-2 bg-red-200 text-red-800 rounded">
            {error}
          </div>
        )}
        {success && (
          <div className="mb-4 p-2 bg-green-200 text-green-800 rounded">
            {success}
          </div>
        )}
        <div className="mb-4">
          <label htmlFor="fullName" className="block text-gray-700">
            Nama Lengkap
          </label>
          <input
            type="text"
            id="fullName"
            className="w-full p-2 border border-gray-300 rounded mt-1"
            value={fullName}
            onChange={(e) => setFullName(e.target.value)}
            required
            placeholder="Masukkan nama lengkap Anda"
          />
        </div>
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
        <div className="mb-4">
          <label htmlFor="confirmPassword" className="block text-gray-700">
            Konfirmasi Password
          </label>
          <input
            type="password"
            id="confirmPassword"
            className="w-full p-2 border border-gray-300 rounded mt-1"
            value={confirmPassword}
            onChange={(e) => setConfirmPassword(e.target.value)}
            required
            placeholder="Konfirmasi password Anda"
          />
        </div>
        <div className="mb-4">
          <label className="block text-gray-700">Jenis Kelamin</label>
          <div className="mt-1 flex items-center gap-4">
            <label className="inline-flex items-center">
              <input
                type="radio"
                name="gender"
                value="male"
                checked={gender === "male"}
                onChange={(e) => setGender(e.target.value)}
                className="form-radio h-4 w-4 text-blue-600"
                required
              />
              <span className="ml-2">Pria</span>
            </label>
            <label className="inline-flex items-center">
              <input
                type="radio"
                name="gender"
                value="female"
                checked={gender === "female"}
                onChange={(e) => setGender(e.target.value)}
                className="form-radio h-4 w-4 text-pink-600"
                required
              />
              <span className="ml-2">Wanita</span>
            </label>
          </div>
        </div>
        <button
          type="submit"
          className={`w-full text-white p-2 rounded transition ${
            loading
              ? "bg-gray-400 cursor-not-allowed"
              : "bg-green-500 hover:bg-green-600"
          }`}
          disabled={loading}
        >
          {loading ? "Sedang Mendaftar..." : "Daftar"}
        </button>
      </form>
      <p className="mt-4 text-center text-gray-600">
        Sudah punya akun?{" "}
        <Link href="/login" className="text-blue-500 hover:underline">
          Masuk di sini
        </Link>
        .
      </p>
    </div>
  );
};

export default RegisterPage;
