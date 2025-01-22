"use client";

import React, { useEffect, useState, useContext, ChangeEvent } from "react";
import { useRouter } from "next/navigation";
import { AuthContext } from "../../context/AuthContext";
import ProtectedRoute from "../components/ProtectRoute";

// ------------------------------------------------
// Interface
// ------------------------------------------------
interface FullUserProfile {
    // Berisi semua field yang muncul di respons GET /users/:id
    id: number;
    full_name: string;
    gender: string;
    email: string;
    password: string;
    role: string;
    reset_password_token: string;
    verify_email_token: string;
    is_verified: number;
    created_at: string;
    updated_at: string;
}

interface ProductItem {
    id: number;
    product_name: string;
    product_address: string;
    product_image: string;
    product_time: string;   // "10:00:00"
    product_date: string;   // "2024-12-17T00:00:00+07:00"
    product_price: number;
    product_sold: number;
    product_description: string;
    product_category: string;
    product_quantity: number;
    product_type: string;
    product_status: string;
    user_id: number;
    order_id: string;
    created_at: string;
    updated_at: string;
}

interface TransactionItem {
    id: number;
    transaction_status: string;  // "success", "pending", ...
    product_id: number;
    user_id: number;
    transaction_quantity: number;
    transaction_amount: number;
    transaction_type: string;     // "submission" or "ticket"
    order_id: string;
    verification_token: string;
    check_in: number;             // 0 atau 1
    created_at: string;
    updated_at: string;
}

// ------------------------------------------------
// Komponen
// ------------------------------------------------
const AdminPage: React.FC = () => {
    const router = useRouter();
    const { token, logout, authLoading, role } = useContext(AuthContext);

    // Pilihan menu: "User" | "Product" | "Transaction"
    const [activeTab, setActiveTab] = useState<"user" | "product" | "transaction">("user");

    // ------------------------------------------------
    // State umum
    // ------------------------------------------------
    const [error, setError] = useState<string | null>(null);
    const [loading, setLoading] = useState<boolean>(true);

    // ------------------------------------------------
    // State untuk User
    // ------------------------------------------------
    // Daftar user ringkas (untuk list). 
    // Anda boleh menambahkan field lain jika /users memang kembalikan data lengkap 
    // atau gunakan FullUserProfile[] jika endpoint-nya sama.
    interface BasicUserInfo {

        id: number;
        full_name: string;
        gender: string;
        email: string;
        role: string;
        created_at: string;
        updated_at: string;
        is_verified: number;

    }
    const [users, setUsers] = useState<BasicUserInfo[]>([]);

    // Input search & hasil pencarian user
    const [searchUserID, setSearchUserID] = useState<string>("");
    const [searchedUser, setSearchedUser] = useState<FullUserProfile | null>(null);

    // ------------------------------------------------
    // State untuk Product
    // ------------------------------------------------
    const [products, setProducts] = useState<ProductItem[]>([]);
    const [searchProductID, setSearchProductID] = useState<string>("");
    const [searchedProduct, setSearchedProduct] = useState<ProductItem | null>(null);

    // ------------------------------------------------
    // State untuk Transaction
    // ------------------------------------------------
    const [transactions, setTransactions] = useState<TransactionItem[]>([]);
    const [searchTransID, setSearchTransID] = useState<string>("");
    const [searchedTrans, setSearchedTrans] = useState<TransactionItem | null>(null);

    // ------------------------------------------------
    // useEffect: cek role
    // ------------------------------------------------
    useEffect(() => {
        if (!authLoading) {
            if (!token) {
                router.push("/login");
            } else {
                if (role !== "Administrator") {
                    // Bukan admin => redirect
                    router.push("/");
                } else {
                    // Jika admin, default tab user => fetch user
                    fetchAllUsers();
                }
            }
        }
    }, [authLoading, token, role, router]);

    // ------------------------------------------------
    // Fungsi fetch data (Users, Products, Transactions)
    // ------------------------------------------------
    const fetchAllUsers = async () => {
        setLoading(true);
        setError(null);
        try {
            const resp = await fetch("http://localhost:8080/api/v1/users", {
                headers: { Authorization: token ? `Bearer ${token}` : "" },
            });
            if (resp.status === 401) {
                logout();
                return;
            }
            if (!resp.ok) {
                const data = await resp.json();
                throw new Error(data.meta?.message || "Gagal mengambil data users.");
            }
            const result = await resp.json();
            // Asumsikan /users hanya kembalikan data ringkas (ID, full_name, gender, email)
            setUsers(result.data || []);
        } catch (err: any) {
            setError(err.message || "Terjadi kesalahan fetch users.");
        } finally {
            setLoading(false);
        }
    };

    const fetchAllProducts = async () => {
        setLoading(true);
        setError(null);
        try {
            const resp = await fetch("http://localhost:8080/api/v1/products", {
                headers: { Authorization: token ? `Bearer ${token}` : "" },
            });
            if (resp.status === 401) {
                logout();
                return;
            }
            if (!resp.ok) {
                const data = await resp.json();
                throw new Error(data.meta?.message || "Gagal mengambil data products.");
            }
            const result = await resp.json();
            setProducts(result.data || []);
        } catch (err: any) {
            setError(err.message || "Terjadi kesalahan fetch products.");
        } finally {
            setLoading(false);
        }
    };

    const fetchAllTransactions = async () => {
        setLoading(true);
        setError(null);
        try {
            const resp = await fetch("http://localhost:8080/api/v1/transactions", {
                headers: { Authorization: token ? `Bearer ${token}` : "" },
            });
            if (resp.status === 401) {
                logout();
                return;
            }
            if (!resp.ok) {
                const data = await resp.json();
                throw new Error(data.meta?.message || "Gagal mengambil data transactions.");
            }
            const result = await resp.json();
            setTransactions(result.data || []);
        } catch (err: any) {
            setError(err.message || "Terjadi kesalahan fetch transactions.");
        } finally {
            setLoading(false);
        }
    };

    // ------------------------------------------------
    // Switch tab
    // ------------------------------------------------
    const handleTabClick = (tab: "user" | "product" | "transaction") => {
        setActiveTab(tab);
        setError(null);
        setSearchedUser(null);
        setSearchedProduct(null);
        setSearchedTrans(null);

        // Panggil fetch sesuai tab
        if (tab === "user") {
            fetchAllUsers();
        } else if (tab === "product") {
            fetchAllProducts();
        } else {
            fetchAllTransactions();
        }
    };

    // ------------------------------------------------
    // ADMIN: Hapus user
    // ------------------------------------------------
    const handleDeleteUser = async (userId: number) => {
        const sure = window.confirm(`Yakin hapus user ID: ${userId}?`);
        if (!sure) return;
        setLoading(true);
        setError(null);
        try {
            const resp = await fetch(`http://localhost:8080/api/v1/users/${userId}`, {
                method: "DELETE",
                headers: { Authorization: token ? `Bearer ${token}` : "" },
            });
            if (resp.status === 401) {
                logout();
                return;
            }
            if (!resp.ok) {
                const data = await resp.json();
                throw new Error(data.meta?.message || "Gagal menghapus user.");
            }
            setUsers((prev) => prev.filter((u) => u.id !== userId));
            if (searchedUser && searchedUser.id === userId) {
                setSearchedUser(null);
            }
        } catch (err: any) {
            setError(err.message || "Terjadi kesalahan hapus user.");
        } finally {
            setLoading(false);
        }
    };

    // ------------------------------------------------
    // ADMIN: Search user by ID (mengembalikan FullUserProfile)
    // ------------------------------------------------
    const handleSearchUserByID = async () => {
        if (!searchUserID.trim()) {
            setError("Mohon masukkan ID user untuk pencarian.");
            return;
        }
        setLoading(true);
        setError(null);
        setSearchedUser(null);
        try {
            const resp = await fetch(`http://localhost:8080/api/v1/users/${searchUserID}`, {
                headers: { Authorization: token ? `Bearer ${token}` : "" },
            });
            if (resp.status === 401) {
                logout();
                return;
            }
            if (resp.status === 404) {
                throw new Error(`User dengan ID ${searchUserID} tidak ditemukan.`);
            }
            if (!resp.ok) {
                throw new Error("Gagal mencari user.");
            }
            const result = await resp.json();
            // Asumsikan di endpoint ini kita dapat data lengkap (FullUserProfile)
            setSearchedUser(result.data);
        } catch (err: any) {
            setError(err.message || "Terjadi kesalahan saat mencari user.");
        } finally {
            setLoading(false);
        }
    };

    // ------------------------------------------------
    // ADMIN: Search product by ID
    // ------------------------------------------------
    const handleSearchProductByID = async () => {
        if (!searchProductID.trim()) {
            setError("Mohon masukkan ID product untuk pencarian.");
            return;
        }
        setLoading(true);
        setError(null);
        setSearchedProduct(null);
        try {
            const resp = await fetch(`http://localhost:8080/api/v1/products/${searchProductID}`, {
                headers: { Authorization: token ? `Bearer ${token}` : "" },
            });
            if (resp.status === 401) {
                logout();
                return;
            }
            if (resp.status === 404) {
                throw new Error(`Product dengan ID ${searchProductID} tidak ditemukan.`);
            }
            if (!resp.ok) {
                throw new Error("Gagal mencari product.");
            }
            const result = await resp.json();
            setSearchedProduct(result.data);
        } catch (err: any) {
            setError(err.message || "Terjadi kesalahan saat mencari product.");
        } finally {
            setLoading(false);
        }
    };

    // ------------------------------------------------
    // ADMIN: Search transaction by ID
    // ------------------------------------------------
    const handleSearchTransByID = async () => {
        if (!searchTransID.trim()) {
            setError("Mohon masukkan ID transaction untuk pencarian.");
            return;
        }
        setLoading(true);
        setError(null);
        setSearchedTrans(null);
        try {
            const resp = await fetch(`http://localhost:8080/api/v1/transactions/${searchTransID}`, {
                headers: { Authorization: token ? `Bearer ${token}` : "" },
            });
            if (resp.status === 401) {
                logout();
                return;
            }
            if (resp.status === 404) {
                throw new Error(`Transaction dengan ID ${searchTransID} tidak ditemukan.`);
            }
            if (!resp.ok) {
                throw new Error("Gagal mencari transaction.");
            }
            const result = await resp.json();
            setSearchedTrans(result.data);
        } catch (err: any) {
            setError(err.message || "Terjadi kesalahan saat mencari transaction.");
        } finally {
            setLoading(false);
        }
    };

    // Render: jika masih authLoading
    if (authLoading) {
        return <div className="container mx-auto px-4 py-8">Memeriksa autentikasi...</div>;
    }

    // Render: jika masih loading data
    if (loading) {
        return (
            <div className="container mx-auto px-4 py-8">
                <h1 className="text-2xl font-bold mb-4">Admin Panel</h1>
                <p>Loading data...</p>
            </div>
        );
    }

    // Render: main
    return (
        <div className="container mx-auto px-4 py-8 text-black">
            <h1 className="text-2xl font-bold mb-6">Admin Panel</h1>

            {/* Tabs Menu */}
            <div className="flex gap-4 mb-6">
                <button
                    onClick={() => handleTabClick("user")}
                    className={`px-4 py-2 rounded ${activeTab === "user" ? "bg-blue-500 text-white" : "bg-gray-300 text-black"
                        }`}
                >
                    User
                </button>
                <button
                    onClick={() => handleTabClick("product")}
                    className={`px-4 py-2 rounded ${activeTab === "product" ? "bg-blue-500 text-white" : "bg-gray-300 text-black"
                        }`}
                >
                    Product
                </button>
                <button
                    onClick={() => handleTabClick("transaction")}
                    className={`px-4 py-2 rounded ${activeTab === "transaction" ? "bg-blue-500 text-white" : "bg-gray-300 text-black"
                        }`}
                >
                    Transaction
                </button>
            </div>

            {/* Pesan Error */}
            {error && <p className="text-red-500 mb-4">{error}</p>}

            {/* ---------------------------------- */}
            {/* TAB: USER */}
            {/* ---------------------------------- */}
            {activeTab === "user" && (
                <div>
                    {/* Search user by ID */}
                    <div className="flex items-center gap-4 mb-6">
                        <div>
                            <label htmlFor="searchUserID" className="block font-semibold">
                                Cari User by ID:
                            </label>
                            <input
                                id="searchUserID"
                                type="text"
                                className="border p-2 rounded"
                                value={searchUserID}
                                onChange={(e: ChangeEvent<HTMLInputElement>) => setSearchUserID(e.target.value)}
                                placeholder="Masukkan ID user"
                            />
                        </div>
                        <button
                            onClick={handleSearchUserByID}
                            className="bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded mt-5"
                        >
                            Cari
                        </button>
                    </div>

                    {/* Hasil Pencarian: FullUserProfile */}
                    {searchedUser && (
                        <div className="mb-4 bg-green-100 p-4 rounded shadow-sm">
                            <h2 className="text-xl font-bold mb-2">Hasil Pencarian (Detail User)</h2>
                            <table className="table-auto w-full text-sm">
                                <tbody>
                                    <tr>
                                        <td className="font-semibold w-1/4">ID</td>
                                        <td>{searchedUser.id}</td>
                                    </tr>
                                    <tr>
                                        <td className="font-semibold">Full Name</td>
                                        <td>{searchedUser.full_name}</td>
                                    </tr>
                                    <tr>
                                        <td className="font-semibold">Gender</td>
                                        <td>{searchedUser.gender}</td>
                                    </tr>
                                    <tr>
                                        <td className="font-semibold">Email</td>
                                        <td>{searchedUser.email}</td>
                                    </tr>
                                    <tr>
                                        <td className="font-semibold">Password (hashed)</td>
                                        <td>{searchedUser.password}</td>
                                    </tr>
                                    <tr>
                                        <td className="font-semibold">Role</td>
                                        <td>{searchedUser.role}</td>
                                    </tr>
                                    <tr>
                                        <td className="font-semibold">Reset Password Token</td>
                                        <td>{searchedUser.reset_password_token}</td>
                                    </tr>
                                    <tr>
                                        <td className="font-semibold">Verify Email Token</td>
                                        <td>{searchedUser.verify_email_token}</td>
                                    </tr>
                                    <tr>
                                        <td className="font-semibold">Is Verified</td>
                                        <td>{searchedUser.is_verified}</td>
                                    </tr>
                                    <tr>
                                        <td className="font-semibold">Created At</td>
                                        <td>{searchedUser.created_at}</td>
                                    </tr>
                                    <tr>
                                        <td className="font-semibold">Updated At</td>
                                        <td>{searchedUser.updated_at}</td>
                                    </tr>
                                </tbody>
                            </table>
                            <button
                                onClick={() => handleDeleteUser(searchedUser.id)}
                                className="bg-red-500 text-white px-3 py-1 mt-3 rounded hover:bg-red-600"
                            >
                                Hapus User Ini
                            </button>
                        </div>
                    )}

                    {/* Tabel User (ringkas) */}
                    <table className="min-w-full bg-white border border-gray-300">
                        <thead className="bg-gray-200 border-b border-gray-300">
                            <tr>
                                <th>ID</th>
                                <th>Nama Lengkap</th>
                                <th>Gender</th>
                                <th>Email</th>
                                <th>Role</th>
                                <th>Is Verified</th>
                                <th>Created At</th>
                                <th>Updated At</th>
                                <th>Aksi</th>
                            </tr>
                        </thead>
                        <tbody>
                            {users.length === 0 ? (
                                <tr>
                                    <td colSpan={5} className="py-2 px-4 text-center">
                                        Belum ada data user.
                                    </td>
                                </tr>
                            ) : (
                                users.map((u) => (
                                    <tr key={u.id} className="border-b border-gray-200">
                                        <td className="py-2 px-4">{u.id}</td>
                                        <td className="py-2 px-4">{u.full_name}</td>
                                        <td className="py-2 px-4">{u.gender}</td>
                                        <td className="py-2 px-4">{u.email}</td>
                                        <td className="py-2 px-4">{u.role}</td>
                                        <td className="py-2 px-4">{u.is_verified}</td>
                                        <td className="py-2 px-4">{u.created_at}</td>
                                        <td className="py-2 px-4">{u.updated_at}</td>
                                        
                                        
                                        <td className="py-2 px-4">
                                            <button
                                                onClick={() => handleDeleteUser(u.id)}
                                                className="bg-red-500 text-white px-3 py-1 rounded hover:bg-red-600"
                                            >
                                                Hapus
                                            </button>
                                        </td>
                                    </tr>
                                ))
                            )}
                        </tbody>
                    </table>
                </div>
            )}

            {/* ---------------------------------- */}
            {/* TAB: PRODUCT */}
            {/* ---------------------------------- */}
            {activeTab === "product" && (
                <div>
                    {/* Search Product by ID */}
                    <div className="flex items-center gap-4 mb-6">
                        <div>
                            <label htmlFor="searchProductID" className="block font-semibold">
                                Cari Product by ID:
                            </label>
                            <input
                                id="searchProductID"
                                type="text"
                                className="border p-2 rounded"
                                value={searchProductID}
                                onChange={(e: ChangeEvent<HTMLInputElement>) => setSearchProductID(e.target.value)}
                                placeholder="Masukkan ID product"
                            />
                        </div>
                        <button
                            onClick={handleSearchProductByID}
                            className="bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded mt-5"
                        >
                            Cari
                        </button>
                    </div>

                    {/* Hasil Pencarian Product */}
                    {searchedProduct && (
                        <div className="mb-4 bg-green-100 p-4 rounded shadow-sm">
                            <h2 className="text-xl font-bold mb-2">Hasil Pencarian Product</h2>
                            <table className="table-auto w-full text-sm">
                                <tbody>
                                    <tr>
                                        <td className="font-semibold w-1/4">ID</td>
                                        <td>{searchedProduct.id}</td>
                                    </tr>
                                    <tr>
                                        <td className="font-semibold">Name</td>
                                        <td>{searchedProduct.product_name}</td>
                                    </tr>
                                    <tr>
                                        <td className="font-semibold">Address</td>
                                        <td>{searchedProduct.product_address}</td>
                                    </tr>
                                    <tr>
                                        <td className="font-semibold">Image</td>
                                        <td>{searchedProduct.product_image}</td>
                                    </tr>
                                    <tr>
                                        <td className="font-semibold">Time</td>
                                        <td>{searchedProduct.product_time}</td>
                                    </tr>
                                    <tr>
                                        <td className="font-semibold">Date</td>
                                        <td>{searchedProduct.product_date}</td>
                                    </tr>
                                    <tr>
                                        <td className="font-semibold">Price</td>
                                        <td>{searchedProduct.product_price}</td>
                                    </tr>
                                    <tr>
                                        <td className="font-semibold">Sold</td>
                                        <td>{searchedProduct.product_sold}</td>
                                    </tr>
                                    <tr>
                                        <td className="font-semibold">Description</td>
                                        <td>{searchedProduct.product_description}</td>
                                    </tr>
                                    <tr>
                                        <td className="font-semibold">Category</td>
                                        <td>{searchedProduct.product_category}</td>
                                    </tr>
                                    <tr>
                                        <td className="font-semibold">Quantity</td>
                                        <td>{searchedProduct.product_quantity}</td>
                                    </tr>
                                    <tr>
                                        <td className="font-semibold">Type</td>
                                        <td>{searchedProduct.product_type}</td>
                                    </tr>
                                    <tr>
                                        <td className="font-semibold">Status</td>
                                        <td>{searchedProduct.product_status}</td>
                                    </tr>
                                    <tr>
                                        <td className="font-semibold">User ID</td>
                                        <td>{searchedProduct.user_id}</td>
                                    </tr>
                                    <tr>
                                        <td className="font-semibold">Order ID</td>
                                        <td>{searchedProduct.order_id}</td>
                                    </tr>
                                    <tr>
                                        <td className="font-semibold">Created At</td>
                                        <td>{searchedProduct.created_at}</td>
                                    </tr>
                                    <tr>
                                        <td className="font-semibold">Updated At</td>
                                        <td>{searchedProduct.updated_at}</td>
                                    </tr>
                                </tbody>
                            </table>
                        </div>
                    )}

                    {/* Tabel Product */}
                    <table className="min-w-full bg-white border border-gray-300 text-sm">
                        <thead className="bg-gray-200 border-b border-gray-300">
                            <tr>
                                <th className="py-2 px-4 text-left">ID</th>
                                <th className="py-2 px-4 text-left">Name</th>
                                <th className="py-2 px-4 text-left">Address</th>
                                <th className="py-2 px-4 text-left">Image</th>
                                <th className="py-2 px-4 text-left">Time</th>
                                <th className="py-2 px-4 text-left">Date</th>
                                <th className="py-2 px-4 text-left">Price</th>
                                <th className="py-2 px-4 text-left">Sold</th>
                                <th className="py-2 px-4 text-left">Description</th>
                                <th className="py-2 px-4 text-left">Category</th>
                                <th className="py-2 px-4 text-left">Quantity</th>
                                <th className="py-2 px-4 text-left">Type</th>
                                <th className="py-2 px-4 text-left">Status</th>
                                <th className="py-2 px-4 text-left">User ID</th>
                                <th className="py-2 px-4 text-left">Order ID</th>
                                <th className="py-2 px-4 text-left">Created At</th>
                                <th className="py-2 px-4 text-left">Updated At</th>
                            </tr>
                        </thead>
                        <tbody>
                            {products.length === 0 ? (
                                <tr>
                                    <td colSpan={17} className="py-2 px-4 text-center">
                                        Belum ada data product.
                                    </td>
                                </tr>
                            ) : (
                                products.map((p) => (
                                    <tr key={p.id} className="border-b border-gray-200">
                                        <td className="py-2 px-4">{p.id}</td>
                                        <td className="py-2 px-4">{p.product_name}</td>
                                        <td className="py-2 px-4">{p.product_address}</td>
                                        <td className="py-2 px-4">{p.product_image}</td>
                                        <td className="py-2 px-4">{p.product_time}</td>
                                        <td className="py-2 px-4">{p.product_date}</td>
                                        <td className="py-2 px-4">{p.product_price}</td>
                                        <td className="py-2 px-4">{p.product_sold}</td>
                                        <td className="py-2 px-4">{p.product_description}</td>
                                        <td className="py-2 px-4">{p.product_category}</td>
                                        <td className="py-2 px-4">{p.product_quantity}</td>
                                        <td className="py-2 px-4">{p.product_type}</td>
                                        <td className="py-2 px-4">{p.product_status}</td>
                                        <td className="py-2 px-4">{p.user_id}</td>
                                        <td className="py-2 px-4">{p.order_id}</td>
                                        <td className="py-2 px-4">{p.created_at}</td>
                                        <td className="py-2 px-4">{p.updated_at}</td>
                                    </tr>
                                ))
                            )}
                        </tbody>
                    </table>
                </div>
            )}

            {/* ---------------------------------- */}
            {/* TAB: TRANSACTION */}
            {/* ---------------------------------- */}
            {activeTab === "transaction" && (
                <div>
                    {/* Search Transaction by ID */}
                    <div className="flex items-center gap-4 mb-6">
                        <div>
                            <label htmlFor="searchTransID" className="block font-semibold">
                                Cari Transaction by ID:
                            </label>
                            <input
                                id="searchTransID"
                                type="text"
                                className="border p-2 rounded"
                                value={searchTransID}
                                onChange={(e: ChangeEvent<HTMLInputElement>) => setSearchTransID(e.target.value)}
                                placeholder="Masukkan ID transaction"
                            />
                        </div>
                        <button
                            onClick={handleSearchTransByID}
                            className="bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded mt-5"
                        >
                            Cari
                        </button>
                    </div>

                    {searchedTrans && (
                        <div className="mb-4 bg-green-100 p-4 rounded shadow-sm">
                            <h2 className="text-xl font-bold mb-2">Hasil Pencarian Transaction</h2>
                            <table className="table-auto w-full text-sm">
                                <tbody>
                                    <tr>
                                        <td className="font-semibold w-1/4">ID</td>
                                        <td>{searchedTrans.id}</td>
                                    </tr>
                                    <tr>
                                        <td className="font-semibold">Status</td>
                                        <td>{searchedTrans.transaction_status}</td>
                                    </tr>
                                    <tr>
                                        <td className="font-semibold">Product ID</td>
                                        <td>{searchedTrans.product_id}</td>
                                    </tr>
                                    <tr>
                                        <td className="font-semibold">User ID</td>
                                        <td>{searchedTrans.user_id}</td>
                                    </tr>
                                    <tr>
                                        <td className="font-semibold">Quantity</td>
                                        <td>{searchedTrans.transaction_quantity}</td>
                                    </tr>
                                    <tr>
                                        <td className="font-semibold">Amount</td>
                                        <td>{searchedTrans.transaction_amount}</td>
                                    </tr>
                                    <tr>
                                        <td className="font-semibold">Type</td>
                                        <td>{searchedTrans.transaction_type}</td>
                                    </tr>
                                    <tr>
                                        <td className="font-semibold">Order ID</td>
                                        <td>{searchedTrans.order_id}</td>
                                    </tr>
                                    <tr>
                                        <td className="font-semibold">Verification Token</td>
                                        <td>{searchedTrans.verification_token}</td>
                                    </tr>
                                    <tr>
                                        <td className="font-semibold">Check In</td>
                                        <td>{searchedTrans.check_in}</td>
                                    </tr>
                                    <tr>
                                        <td className="font-semibold">Created At</td>
                                        <td>{searchedTrans.created_at}</td>
                                    </tr>
                                    <tr>
                                        <td className="font-semibold">Updated At</td>
                                        <td>{searchedTrans.updated_at}</td>
                                    </tr>
                                </tbody>
                            </table>
                        </div>
                    )}

                    {/* Tabel Transaction */}
                    <table className="min-w-full bg-white border border-gray-300 text-sm">
                        <thead className="bg-gray-200 border-b border-gray-300">
                            <tr>
                                <th className="py-2 px-4">ID</th>
                                <th className="py-2 px-4">Status</th>
                                <th className="py-2 px-4">Product ID</th>
                                <th className="py-2 px-4">User ID</th>
                                <th className="py-2 px-4">Quantity</th>
                                <th className="py-2 px-4">Amount</th>
                                <th className="py-2 px-4">Type</th>
                                <th className="py-2 px-4">Order ID</th>
                                <th className="py-2 px-4">Verification Token</th>
                                <th className="py-2 px-4">Check In</th>
                                <th className="py-2 px-4">Created At</th>
                                <th className="py-2 px-4">Updated At</th>
                            </tr>
                        </thead>
                        <tbody>
                            {transactions.length === 0 ? (
                                <tr>
                                    <td colSpan={12} className="py-2 px-4 text-center">
                                        Belum ada data transaction.
                                    </td>
                                </tr>
                            ) : (
                                transactions.map((t) => (
                                    <tr key={t.id} className="border-b border-gray-200">
                                        <td className="py-2 px-4">{t.id}</td>
                                        <td className="py-2 px-4">{t.transaction_status}</td>
                                        <td className="py-2 px-4">{t.product_id}</td>
                                        <td className="py-2 px-4">{t.user_id}</td>
                                        <td className="py-2 px-4">{t.transaction_quantity}</td>
                                        <td className="py-2 px-4">{t.transaction_amount}</td>
                                        <td className="py-2 px-4">{t.transaction_type}</td>
                                        <td className="py-2 px-4">{t.order_id}</td>
                                        <td className="py-2 px-4">{t.verification_token}</td>
                                        <td className="py-2 px-4">{t.check_in}</td>
                                        <td className="py-2 px-4">{t.created_at}</td>
                                        <td className="py-2 px-4">{t.updated_at}</td>
                                    </tr>
                                ))
                            )}
                        </tbody>
                    </table>
                </div>
            )}

            {/* Tombol Logout */}
            <div className="mt-6 flex justify-end">
                <button
                    onClick={logout}
                    className="bg-red-500 text-white px-4 py-2 rounded hover:bg-red-600 transition"
                >
                    Logout
                </button>
            </div>
        </div>
    );
};

// Bungkus AdminPage dengan ProtectedRoute
const ProtectedAdminPage: React.FC = () => {
    return (
        <ProtectedRoute>
            <AdminPage />
        </ProtectedRoute>
    );
};

export default ProtectedAdminPage;
