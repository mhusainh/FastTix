// pages/transactions.tsx

"use client";

import React, { useEffect, useState, useContext } from "react";
import { useRouter } from "next/navigation";
import { AuthContext } from "../../context/AuthContext";
import ProtectedRoute from "../components/ProtectRoute";

interface Transaction {
  id: number;
  transaction_status: "success" | "pending" | "failed"; // Tambahkan status lainnya jika ada
  product_id: number;
  user_id: number;
  transaction_quantity: number;
  transaction_amount: number;
  transaction_type: "submission" | "ticket" | string; // Tambahkan tipe lainnya jika ada
  order_id: string;
  verification_token: string;
  check_in: number; // 1 untuk sudah check-in, 0 untuk belum
  created_at: string;
  updated_at: string;
}

const TransactionsPage: React.FC = () => {
  const router = useRouter();
  const { token, logout, authLoading, role } = useContext(AuthContext);

  // Pastikan hanya user biasa yang bisa mengakses
  useEffect(() => {
    if (!authLoading) {
      if (!token) {
        router.push("/login");
      } else {
        if (role === "Administrator") {
          router.push("/admin");
        }
      }
    }
  }, [authLoading, token, role, router]);

  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchTransactions = async () => {
    setLoading(true);
    setError(null);

    try {
      const response = await fetch("http://localhost:8080/api/v1/transactions/user", {
        headers: {
          Authorization: token ? `Bearer ${token}` : "",
        },
      });

      if (response.status === 401) {
        logout();
        return;
      }

      if (!response.ok) {
        const data = await response.json();
        throw new Error(data.meta?.message || "Gagal mengambil data transaksi.");
      }

      const result = await response.json();
      setTransactions(result.data);
    } catch (err: any) {
      setError(err.message || "Terjadi kesalahan saat mengambil transaksi.");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (token && role !== "Administrator") {
      fetchTransactions();
    }
  }, [token, role]);

  // Format tanggal
  const formatDate = (dateString: string) => {
    const options: Intl.DateTimeFormatOptions = { year: 'numeric', month: 'long', day: 'numeric', hour: '2-digit', minute: '2-digit' };
    return new Date(dateString).toLocaleDateString(undefined, options);
  };

  // Format jumlah uang
  const formatAmount = (amount: number) => {
    return new Intl.NumberFormat(undefined, { style: 'currency', currency: 'IDR' }).format(amount);
  };

  // Mendapatkan status transaksi dengan label
  const getStatusLabel = (status: Transaction['transaction_status']) => {
    switch (status) {
      case "success":
        return <span className="text-green-600 font-semibold">Sukses</span>;
      case "pending":
        return <span className="text-yellow-600 font-semibold">Pending</span>;
      case "failed":
        return <span className="text-red-600 font-semibold">Gagal</span>;
      default:
        return <span className="text-gray-600 font-semibold">{status}</span>;
    }
  };

  // Mendapatkan tipe transaksi dengan label
  const getTypeLabel = (type: Transaction['transaction_type']) => {
    switch (type) {
      case "submission":
        return "Pendaftaran";
      case "ticket":
        return "Tiket";
      default:
        return type.charAt(0).toUpperCase() + type.slice(1);
    }
  };

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-6">Transaksi Anda</h1>
      <div className="max-w-4xl mx-auto bg-white shadow-md rounded-lg p-6">
        {loading ? (
          <p>Loading transaksi...</p>
        ) : error ? (
          <p className="text-red-500">{error}</p>
        ) : transactions.length === 0 ? (
          <p>Anda belum melakukan transaksi.</p>
        ) : (
          <div className="overflow-x-auto">
            <table className="min-w-full table-auto">
              <thead>
                <tr>
                  <th className="px-4 py-2 border">ID Transaksi</th>
                  <th className="px-4 py-2 border">Status</th>
                  <th className="px-4 py-2 border">Tipe</th>
                  <th className="px-4 py-2 border">Jumlah</th>
                  <th className="px-4 py-2 border">Tanggal</th>
                  <th className="px-4 py-2 border">Check-In</th>
                </tr>
              </thead>
              <tbody>
                {transactions.map((txn) => (
                  <tr key={txn.id} className="text-center">
                    <td className="px-4 py-2 border">{txn.id}</td>
                    <td className="px-4 py-2 border">{getStatusLabel(txn.transaction_status)}</td>
                    <td className="px-4 py-2 border">{getTypeLabel(txn.transaction_type)}</td>
                    <td className="px-4 py-2 border">{formatAmount(txn.transaction_amount)}</td>
                    <td className="px-4 py-2 border">{formatDate(txn.created_at)}</td>
                    <td className="px-4 py-2 border">
                      {txn.check_in === 1 ? (
                        <span className="text-green-600 font-semibold">Sudah Check-In</span>
                      ) : (
                        <span className="text-gray-600 font-semibold">Belum Check-In</span>
                      )}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
        <div className="mt-6 flex justify-end">
          <button
            onClick={() => router.back()}
            className="bg-gray-500 text-white px-4 py-2 rounded hover:bg-gray-600 transition"
          >
            Kembali
          </button>
        </div>
      </div>
    </div>
  );
};

// Bungkus dengan ProtectedRoute agar hanya user yang terotentikasi bisa akses
const ProtectedTransactionsPage: React.FC = () => {
  return (
    <ProtectedRoute>
      <TransactionsPage />
    </ProtectedRoute>
  );
};

export default ProtectedTransactionsPage;
