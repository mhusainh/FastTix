// components/ProfilePage.tsx

"use client";

import React, { useEffect, useState, useContext, ChangeEvent, FormEvent } from "react";
import { useRouter } from "next/navigation";
import { AuthContext } from "../../context/AuthContext";
import ProtectedRoute from "../components/ProtectRoute";

interface UserProfile {
  ID: number; // ID user
  full_name: string;
  gender: string;
  email: string;
}

// Halaman Profile ini HANYA untuk user biasa
const ProfilePage: React.FC = () => {
  const router = useRouter();
  const { token, logout, authLoading, role } = useContext(AuthContext);

  // Pastikan jika role = "Administrator", redirect (atau tampilkan pesan saja),
  // karena ini hanya untuk user biasa. Bisa juga diarahkan ke "/admin".
  useEffect(() => {
    if (!authLoading) {
      // Jika user belum login => ke /login
      if (!token) {
        router.push("/login");
      } else {
        // Jika ternyata rolenya "Administrator", kita lempar ke /admin
        if (role === "Administrator") {
          router.push("/admin");
        }
      }
    }
  }, [authLoading, token, role, router]);

  // State profil user
  const [profile, setProfile] = useState<UserProfile | null>(null);

  // Loading & Error
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // State untuk edit profile
  const [isEditing, setIsEditing] = useState(false);
  const [updateLoading, setUpdateLoading] = useState(false);
  const [updateError, setUpdateError] = useState<string | null>(null);
  const [updateSuccess, setUpdateSuccess] = useState<string | null>(null);

  // State form data
  const [formData, setFormData] = useState({
    full_name: "",
    gender: "",
    email: "",
  });

  // Ambil data profile user biasa
  const fetchUserProfile = async () => {
    setLoading(true);
    setError(null);

    try {
      const response = await fetch("http://localhost:8080/api/v1/users/profile", {
        headers: {
          Authorization: token ? `Bearer ${token}` : "",
        },
      });

      if (response.status === 401) {
        // Token invalid => logout
        logout();
        return;
      }

      if (!response.ok) {
        const data = await response.json();
        throw new Error(data.meta?.message || "Gagal mengambil data profil.");
      }

      const result = await response.json();
      setProfile(result.data);
      // Isi formData
      setFormData({
        full_name: result.data.full_name,
        gender: result.data.gender,
        email: result.data.email,
      });
    } catch (err: any) {
      setError(err.message || "Terjadi kesalahan saat mengambil profil.");
    } finally {
      setLoading(false);
    }
  };

  // useEffect => panggil fetchUserProfile
  useEffect(() => {
    if (token && role !== "Administrator") {
      fetchUserProfile();
    }
  }, [token, role]);

  // Handle input
  const handleChange = (e: ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    const { name, value } = e.target;
    setFormData((prev) => ({
      ...prev,
      [name]: value,
    }));
  };

  // Submit update profile
  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setUpdateError(null);
    setUpdateSuccess(null);

    // Validasi sederhana
    if (!formData.full_name.trim()) {
      setUpdateError("Nama lengkap tidak boleh kosong.");
      return;
    }
    if (!formData.email.trim()) {
      setUpdateError("Email tidak boleh kosong.");
      return;
    }
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(formData.email)) {
      setUpdateError("Format email tidak valid.");
      return;
    }

    setUpdateLoading(true);
    try {
      const response = await fetch("http://localhost:8080/api/v1/users/profile", {
        method: "PUT",
        headers: {
          Authorization: token ? `Bearer ${token}` : "",
          "Content-Type": "application/json",
        },
        body: JSON.stringify(formData),
      });

      if (response.status === 401) {
        logout();
        return;
      }

      if (!response.ok) {
        const data = await response.json();
        throw new Error(data.meta?.message || "Gagal memperbarui profil.");
      }

      const result = await response.json();
      setProfile(result.data);
      setUpdateSuccess("Profil berhasil diperbarui.");
      setIsEditing(false);
    } catch (err: any) {
      setUpdateError(err.message || "Terjadi kesalahan saat memperbarui profil.");
    } finally {
      setUpdateLoading(false);
    }
  };

  // Batal edit
  const handleCancelEdit = () => {
    setUpdateError(null);
    setUpdateSuccess(null);
    setIsEditing(false);
    if (profile) {
      setFormData({
        full_name: profile.full_name,
        gender: profile.gender,
        email: profile.email,
      });
    }
  };

  // Navigasi ke halaman transaksi
  const handleViewTransactions = () => {
    router.push("/transactions");
  };

  // Render
  if (authLoading) {
    return <div className="container mx-auto px-4 py-8">Memeriksa autentikasi...</div>;
  }

  if (loading) {
    return (
      <div className="container mx-auto px-4 py-8">
        <h1 className="text-3xl font-bold mb-4">Profile</h1>
        <p>Loading data profil...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="container mx-auto px-4 py-8">
        <h1 className="text-3xl font-bold mb-4">Profile</h1>
        <p className="text-red-500">{error}</p>
      </div>
    );
  }

  if (!profile) {
    return (
      <div className="container mx-auto px-4 py-8">
        <h1 className="text-3xl font-bold mb-4">Profile</h1>
        <p>Profil tidak ditemukan.</p>
      </div>
    );
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-6">Profile</h1>
      <div className="max-w-md mx-auto bg-white shadow-md rounded-lg p-6">
        {updateError && (
          <div className="mb-4 p-2 bg-red-200 text-red-800 rounded">
            {updateError}
          </div>
        )}
        {updateSuccess && (
          <div className="mb-4 p-2 bg-green-200 text-green-800 rounded">
            {updateSuccess}
          </div>
        )}

        {/* Tombol Edit Profile, Lihat Transaksi + Logout */}
        <div className="mb-6 flex items-center justify-between space-x-2">
          {!isEditing && (
            <button
              onClick={() => setIsEditing(true)}
              className="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600 transition"
            >
              Edit Profile
            </button>
          )}
          <button
            onClick={handleViewTransactions}
            className="bg-green-500 text-white px-4 py-2 rounded hover:bg-green-600 transition"
          >
            Lihat Transaksi
          </button>
          <button
            onClick={logout}
            className="bg-red-500 text-white px-4 py-2 rounded hover:bg-red-600 transition"
          >
            Logout
          </button>
        </div>

        {/* Form Update Profile */}
        {isEditing ? (
          <form onSubmit={handleSubmit}>
            <div className="mb-4">
              <label htmlFor="full_name" className="block text-gray-700 font-semibold">
                Nama Lengkap:
              </label>
              <input
                type="text"
                id="full_name"
                name="full_name"
                value={formData.full_name}
                onChange={handleChange}
                className="w-full p-2 border border-gray-300 rounded mt-1"
                required
              />
            </div>

            <div className="mb-4">
              <label htmlFor="gender" className="block text-gray-700 font-semibold">
                Jenis Kelamin:
              </label>
              <select
                id="gender"
                name="gender"
                value={formData.gender}
                onChange={handleChange}
                className="w-full p-2 border border-gray-300 rounded mt-1"
                required
              >
                <option value="" disabled>
                  Pilih Jenis Kelamin
                </option>
                <option value="male">Pria</option>
                <option value="female">Wanita</option>
              </select>
            </div>

            <div className="mb-6">
              <label htmlFor="email" className="block text-gray-700 font-semibold">
                Email:
              </label>
              <input
                type="email"
                id="email"
                name="email"
                value={formData.email}
                onChange={handleChange}
                className="w-full p-2 border border-gray-300 rounded mt-1"
                required
              />
            </div>

            <div className="flex justify-between">
              <button
                type="submit"
                className={`bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600 transition ${
                  updateLoading ? "bg-blue-300 cursor-not-allowed" : ""
                }`}
                disabled={updateLoading}
              >
                {updateLoading ? "Memperbarui..." : "Perbarui Profil"}
              </button>
              <button
                type="button"
                onClick={handleCancelEdit}
                className="bg-gray-500 text-white px-4 py-2 rounded hover:bg-gray-600 transition"
              >
                Batal
              </button>
            </div>
          </form>
        ) : (
          // Tampilkan data profil
          <div>
            <div className="mb-4">
              <label className="block text-gray-700 font-semibold">Nama Lengkap:</label>
              <p className="text-gray-600">{profile.full_name}</p>
            </div>
            <div className="mb-4">
              <label className="block text-gray-700 font-semibold">Jenis Kelamin:</label>
              <p className="text-gray-600">
                {profile.gender === "male" ? "Pria" : "Wanita"}
              </p>
            </div>
            <div className="mb-4">
              <label className="block text-gray-700 font-semibold">Email:</label>
              <p className="text-gray-600">{profile.email}</p>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

// Bungkus dengan ProtectedRoute agar hanya user yang terotentikasi bisa akses
const ProtectedProfilePage: React.FC = () => {
  return (
    <ProtectedRoute>
      <ProfilePage />
    </ProtectedRoute>
  );
};

export default ProtectedProfilePage;
