// src/app/lihat-event-submission/page.tsx
"use client";

import React, { useEffect, useState, useContext } from "react";
import { AuthContext } from "../../context/AuthContext";

interface Submission {
  id: number;
  product_name: string;
  product_status: string;
  product_image: string | null; 
  product_category: string;
  product_quantity: number;
  product_description: string;
  product_price: number;
  product_date: string; // format: "YYYY-MM-DD"
  product_time: string; // format: "HH:mm:ss" (sesuaikan dengan backend)
}

const LihatEventSubmissionPage: React.FC = () => {
  const { token, role, authLoading, logout } = useContext(AuthContext);
  const [submissions, setSubmissions] = useState<Submission[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // ----------------------------------
  // Fungsi: ambil data submission
  // ----------------------------------
  const fetchSubmissions = async () => {
    setLoading(true);
    setError(null);

    try {
      const response = await fetch("http://localhost:8080/api/v1/submissions", {
        headers: {
          Authorization: token ? `Bearer ${token}` : "",
        },
      });

      if (response.status === 401) {
        logout();
        return;
      }

      if (!response.ok) {
        throw new Error("Gagal mengambil data submissions.");
      }

      const result = await response.json();
      setSubmissions(result.data);
    } catch (err: any) {
      setError(err.message || "Terjadi kesalahan.");
    } finally {
      setLoading(false);
    }
  };

  // ----------------------------------
  // useEffect: panggil fetchSubmissions
  // ----------------------------------
  useEffect(() => {
    if (token && !authLoading) {
      fetchSubmissions();
    }
  }, [token, authLoading]);

  // ----------------------------------
  // Fungsi: Approve / Reject
  // ----------------------------------
  const handleChangeStatus = async (id: number, status: "approve" | "reject") => {
    try {
      setLoading(true);
      setError(null);

      const response = await fetch(
        `http://localhost:8080/api/v1/submissions/${id}/${status}`,
        {
          method: "PUT",
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );

      if (response.status === 401) {
        logout();
        return;
      }

      if (!response.ok) {
        throw new Error(`Gagal mengubah status menjadi ${status}.`);
      }

      // Berhasil => refresh data
      await fetchSubmissions();
    } catch (err: any) {
      setError(err.message || "Terjadi kesalahan update status.");
    } finally {
      setLoading(false);
    }
  };

  // ----------------------------------
  // Fungsi: format date & time
  // (opsional, boleh inline)
  // ----------------------------------
  const formatDate = (isoDate: string) => {
    // Asumsikan backend kirim "YYYY-MM-DD"
    const dateObj = new Date(isoDate);
    return dateObj.toLocaleDateString("id-ID", {
      day: "2-digit",
      month: "long",
      year: "numeric",
    });
  };

  const formatTime = (timeString: string) => {
    // Asumsikan backend kirim "HH:mm:ss"
    // Kita buat object date tahun 1970
    const dateObj = new Date(`1970-01-01T${timeString}Z`);
    return dateObj.toLocaleTimeString("id-ID", {
      hour: "2-digit",
      minute: "2-digit",
    });
  };

  // ----------------------------------
  // Render
  // ----------------------------------
  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-6">Lihat Event Submission</h1>

      {loading && <p className="text-gray-700">Loading...</p>}
      {error && <p className="text-red-500">{error}</p>}

      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
        {submissions.map((sub) => (
          <div
            key={sub.id}
            className="bg-white rounded-lg overflow-hidden shadow-md flex flex-col"
          >
            {/* Bagian image */}
            <div className="h-48 w-full overflow-hidden">
              {sub.product_image ? (
                <img
                  src={sub.product_image}
                  alt={sub.product_name}
                  className="w-full h-full object-cover"
                />
              ) : (
                <img
                  src="/no-image-available.jpg"
                  alt="No Image"
                  className="w-full h-full object-cover"
                />
              )}
            </div>

            {/* Isi card */}
            <div className="p-4 flex flex-col flex-grow">
              <h2 className="text-lg font-bold mb-2">{sub.product_name}</h2>

              {/* Category, description, quantity, dll */}
              <p className="text-gray-600 text-sm mb-1">
                <span className="font-semibold">Kategori:</span> {sub.product_category}
              </p>
              <p className="text-gray-600 text-sm mb-1">
                <span className="font-semibold">Jumlah Tiket:</span> {sub.product_quantity}
              </p>
              <p className="text-gray-600 text-sm mb-1">
                <span className="font-semibold">Harga:</span>{" "}
                Rp {sub.product_price.toLocaleString("id-ID")}
              </p>
              <p className="text-gray-600 text-sm mb-1">
                <span className="font-semibold">Tanggal Event:</span>{" "}
                {formatDate(sub.product_date)}
              </p>
              <p className="text-gray-600 text-sm mb-1">
                <span className="font-semibold">Waktu Event:</span>{" "}
                {formatTime(sub.product_time)}
              </p>
              <p className="text-gray-600 text-sm mb-3">
                <span className="font-semibold">Deskripsi:</span>{" "}
                {sub.product_description}
              </p>

              <p className="text-gray-600 mb-4">
                Status: <span className="font-semibold">{sub.product_status}</span>
              </p>

              {/* Tombol hanya untuk admin */}
              {role === "Administrator" && (
                <div className="mt-auto flex items-center gap-2">
                  <button
                    onClick={() => handleChangeStatus(sub.id, "approve")}
                    className="bg-green-500 hover:bg-green-600 text-white px-3 py-1 rounded"
                    disabled={loading}
                  >
                    Approve
                  </button>
                  <button
                    onClick={() => handleChangeStatus(sub.id, "reject")}
                    className="bg-red-500 hover:bg-red-600 text-white px-3 py-1 rounded"
                    disabled={loading}
                  >
                    Reject
                  </button>
                </div>
              )}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default LihatEventSubmissionPage;
