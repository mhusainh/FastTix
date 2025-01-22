"use client";

import React, {
    useState,
    useEffect,
    useContext,
    FormEvent,
    ChangeEvent,
} from "react";
import { useRouter } from "next/navigation";
import { AuthContext } from "../../context/AuthContext";

interface Submission {
    id: number;
    product_name: string;
    product_address: string;
    product_time: string;
    product_date: string;
    product_price: number;
    product_description: string;
    product_category: string;
    product_quantity: number;
    product_status: string;
    product_image: string | null; // URL/path gambar (jika ada)
}

export default function BuatEventPage() {
    const { token, authLoading } = useContext(AuthContext);
    const router = useRouter();

    // ---------------------------------------------
    // 0) Cek apakah user ter-otentikasi
    // ---------------------------------------------
    useEffect(() => {
        // Tunggu dulu sampai authLoading = false
        if (!authLoading) {
            if (!token) {
                router.push("/login");
            }
        }
    }, [token, authLoading, router]);

    // ---------------------------------------------
    // 1) STATE & LOGIC: FORM PEMBUATAN EVENT
    // ---------------------------------------------
    const [createFormData, setCreateFormData] = useState({
        product_name: "",
        product_address: "",
        product_time: "",
        product_date: "",
        product_price: 0,
        product_description: "",
        product_category: "",
        product_quantity: 0,
    });

    // State untuk menampung file gambar saat "Buat Event"
    const [createImageFile, setCreateImageFile] = useState<File | null>(null);

    // State global: error, success, loading
    const [error, setError] = useState<string | null>(null);
    const [success, setSuccess] = useState<string | null>(null);
    const [loading, setLoading] = useState<boolean>(false);

    // Handle perubahan data teks pembuatan event
    const handleChangeCreate = (
        e: ChangeEvent<
            HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement
        >
    ) => {
        const { name, value } = e.target;
        setCreateFormData((prev) => ({
            ...prev,
            [name]:
                name === "product_price" || name === "product_quantity"
                    ? parseInt(value) || 0
                    : value,
        }));
    };

    // Handle perubahan file saat "Buat Event"
    const handleCreateImageFileChange = (e: ChangeEvent<HTMLInputElement>) => {
        if (e.target.files && e.target.files.length > 0) {
            setCreateImageFile(e.target.files[0]);
        }
    };

    // ---------------------------------------------
    // 1.1) VERIFIKASI PAYMENT (STATE)
    // ---------------------------------------------
    const [isVerificationStep, setIsVerificationStep] = useState(false);
    const [verificationCode, setVerificationCode] = useState("");

    // Submit form "Buat Event"
    const handleSubmitCreate = async (e: FormEvent) => {
        e.preventDefault();
        setError(null);
        setSuccess(null);

        try {
            setLoading(true);

            // 1) Buat event baru (tanpa gambar)
            const response = await fetch("http://localhost:8080/api/v1/submissions", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    Authorization: `Bearer ${token}`,
                },
                body: JSON.stringify(createFormData),
            });

            if (!response.ok) {
                const data = await response.json();
                throw new Error(data.meta?.message || "Gagal membuat event.");
            }

            const result = await response.json();
            const newSubmissionId = result.data?.id; // pastikan respons berisi ID

            setSuccess("Event berhasil dibuat!");

            // 2) Jika user memilih gambar saat create, langsung upload
            if (createImageFile && newSubmissionId) {
                const formData = new FormData();
                formData.append("image", createImageFile);

                const uploadResponse = await fetch(
                    `http://localhost:8080/api/v1/submissions/${newSubmissionId}/image`,
                    {
                        method: "PUT",
                        headers: {
                            Authorization: `Bearer ${token}`,
                        },
                        body: formData,
                    }
                );

                if (!uploadResponse.ok) {
                    const data = await uploadResponse.json();
                    // Boleh tangani error upload gambar di sini, 
                    // atau hanya console.error jika tidak mau memblok
                    throw new Error(data.meta?.message || "Gagal upload gambar.");
                }
            }

            // Reset form create
            setCreateFormData({
                product_name: "",
                product_address: "",
                product_time: "",
                product_date: "",
                product_price: 0,
                product_description: "",
                product_category: "",
                product_quantity: 0,
            });
            setCreateImageFile(null);

            // 3) Tampilkan langkah verifikasi payment
            setIsVerificationStep(true);

            // Refresh list submission
            await fetchAndSetSubmissions();
        } catch (err: any) {
            setError(err.message || "Terjadi kesalahan saat membuat event/gambar.");
        } finally {
            setLoading(false);
        }
    };

    // ---------------------------------------------
    // 1.2) HANDLE VERIFIKASI PAYMENT
    // ---------------------------------------------
    const handlePaymentVerification = async (e: FormEvent) => {
        e.preventDefault();
        setError(null);
        setSuccess(null);

        try {
            setLoading(true);

            // 1) Verifikasi kode (POST /api/v1/payment/verify) -> Contoh:
            const verifyResponse = await fetch(`http://localhost:8080/api/v1/payment/checkout/${verificationCode}`, {
       
                headers: {
                    "Content-Type": "application/json",
                    Authorization: `Bearer ${token}`,
                },
    
            });

            if (!verifyResponse.ok) {
                const data = await verifyResponse.json();
                throw new Error(data.meta?.message || "Verifikasi payment gagal.");
            }

            const checkoutData = await verifyResponse.json();

            // 3) Dapatkan redirectURL -> redirect ke Midtrans Snap
            const redirectURL = checkoutData.data?.redirectURL;
            if (!redirectURL) {
                throw new Error("Redirect URL tidak ditemukan di respons server.");
            }

            // 4) Redirect user ke Midtrans Snap
            window.location.href = redirectURL;
        } catch (err: any) {
            setError(err.message || "Terjadi kesalahan saat verifikasi payment.");
        } finally {
            setLoading(false);
        }
    };

    // ---------------------------------------------
    // 2) LIST SUBMISSIONS & ERROR/SUCCESS
    // ---------------------------------------------
    const [submissions, setSubmissions] = useState<Submission[]>([]);

    // GET data submissions user
    const fetchAndSetSubmissions = async () => {
        try {
            setLoading(true);
            const response = await fetch("http://localhost:8080/api/v1/submissions/user", {
                headers: {
                    Authorization: token ? `Bearer ${token}` : "",
                },
            });

            if (response.status === 401) {
                // Token invalid
                router.push("/login");
                return;
            }
            if (!response.ok) {
                throw new Error("Gagal mengambil data submissions.");
            }
            const result = await response.json();
            setSubmissions(result.data);
        } catch (err: any) {
            setError(err.message || "Terjadi kesalahan saat memuat data.");
        } finally {
            setLoading(false);
        }
    };

    // Ambil data event user saat halaman dimuat
    useEffect(() => {
        if (token) {
            fetchAndSetSubmissions();
        }
    }, [token]);

    // ---------------------------------------------
    // 3) EDIT EVENT (DATA TEKS) SAJA
    // ---------------------------------------------
    const [editSubmission, setEditSubmission] = useState<Submission | null>(null);
    const [editFormData, setEditFormData] = useState({
        id: 0,
        product_name: "",
        product_address: "",
        product_time: "",
        product_date: "",
        product_price: 0,
        product_description: "",
        product_category: "",
        product_quantity: 0,
    });

    // Klik "Edit" => isi form edit
    const handleEditClick = (submission: Submission) => {
        setEditSubmission(submission);

        setEditFormData({
            id: submission.id,
            product_name: submission.product_name,
            product_address: submission.product_address,
            product_time: submission.product_time,
            product_date: submission.product_date,
            product_price: submission.product_price,
            product_description: submission.product_description,
            product_category: submission.product_category,
            product_quantity: submission.product_quantity,
        });
        setError(null);
        setSuccess(null);
    };

    // Perubahan input form edit
    const handleEditFormChange = (
        e: ChangeEvent<
            HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement
        >
    ) => {
        const { name, value } = e.target;
        setEditFormData((prev) => ({
            ...prev,
            [name]:
                name === "product_price" || name === "product_quantity"
                    ? parseInt(value) || 0
                    : value,
        }));
    };

    // Submit edit (PUT /submissions/:id)
    const handleUpdateSubmission = async (e: FormEvent) => {
        e.preventDefault();
        if (!editSubmission) return;

        setLoading(true);
        setError(null);
        setSuccess(null);

        try {
            const response = await fetch(
                `http://localhost:8080/api/v1/submissions/${editSubmission.id}`,
                {
                    method: "PUT",
                    headers: {
                        "Content-Type": "application/json",
                        Authorization: `Bearer ${token}`,
                    },
                    body: JSON.stringify(editFormData),
                }
            );

            if (!response.ok) {
                const data = await response.json();
                throw new Error(data.meta?.message || "Gagal update submission.");
            }

            setSuccess("Event berhasil di-update!");

            // Update di state
            setSubmissions((prev) =>
                prev.map((sub) =>
                    sub.id === editSubmission.id
                        ? {
                              ...sub,
                              product_name: editFormData.product_name,
                              product_address: editFormData.product_address,
                              product_time: editFormData.product_time,
                              product_date: editFormData.product_date,
                              product_price: editFormData.product_price,
                              product_description: editFormData.product_description,
                              product_category: editFormData.product_category,
                              product_quantity: editFormData.product_quantity,
                          }
                        : sub
                )
            );

            setEditSubmission(null);
        } catch (err: any) {
            setError(err.message || "Terjadi kesalahan saat update.");
        } finally {
            setLoading(false);
        }
    };

    // Batal edit
    const handleCancelEdit = () => {
        setEditSubmission(null);
        setEditFormData({
            id: 0,
            product_name: "",
            product_address: "",
            product_time: "",
            product_date: "",
            product_price: 0,
            product_description: "",
            product_category: "",
            product_quantity: 0,
        });
        setError(null);
        setSuccess(null);
    };

    // ---------------------------------------------
    // 4) UPLOAD / REPLACE GAMBAR (SETELAH BUAT ATAU SAAT EDIT)
    // ---------------------------------------------
    const [editImageFile, setEditImageFile] = useState<File | null>(null);

    // Handle pilih file (saat edit)
    const handleFileChange = (e: ChangeEvent<HTMLInputElement>) => {
        if (e.target.files && e.target.files.length > 0) {
            setEditImageFile(e.target.files[0]);
        }
    };

    // PUT /submissions/:id/image
    const handleUploadImage = async (e: FormEvent) => {
        e.preventDefault();
        if (!editSubmission || !editImageFile) return;

        setLoading(true);
        setError(null);
        setSuccess(null);

        try {
            const formData = new FormData();
            formData.append("image", editImageFile);

            const response = await fetch(
                `http://localhost:8080/api/v1/submissions/${editSubmission.id}/image`,
                {
                    method: "PUT",
                    headers: {
                        Authorization: `Bearer ${token}`,
                    },
                    body: formData,
                }
            );

            if (!response.ok) {
                const data = await response.json();
                throw new Error(data.meta?.message || "Gagal upload gambar.");
            }

            setSuccess("Gambar berhasil diupload!");

            // Refresh data agar tampilan gambar update
            await fetchAndSetSubmissions();
        } catch (err: any) {
            setError(err.message || "Terjadi kesalahan saat upload gambar.");
        } finally {
            setLoading(false);
            setEditImageFile(null);
        }
    };

    // ---------------------------------------------
    // 5) DELETE SUBMISSION
    // ---------------------------------------------
    const handleDeleteClick = async (id: number) => {
        const confirmation = window.confirm(
            "Apakah Anda yakin ingin menghapus event ini?"
        );
        if (!confirmation) return;

        try {
            setLoading(true);
            setError(null);
            setSuccess(null);

            const response = await fetch(
                `http://localhost:8080/api/v1/submissions/${id}`,
                {
                    method: "DELETE",
                    headers: {
                        Authorization: `Bearer ${token}`,
                    },
                }
            );

            if (!response.ok) {
                const data = await response.json();
                throw new Error(data.meta?.message || "Gagal menghapus submission.");
            }

            setSubmissions((prev) => prev.filter((sub) => sub.id !== id));
            setSuccess("Event berhasil dihapus!");
        } catch (err: any) {
            setError(err.message || "Terjadi kesalahan saat menghapus event.");
        } finally {
            setLoading(false);
        }
    };

    // ---------------------------------------------
    // 6) RENDER COMPONENT
    // ---------------------------------------------
    const categoryOptions = [
        "Music",
        "StandUp",
        "Theatre",
        "Sports",
        "Workshops",
        "Conferences",
        "Exhibitions",
        "Festivals",
        "Other",
    ];

    return (
        <div className="container mx-auto px-4 py-8">
            {/* PESAN ERROR / SUKSES GLOBAL */}
            {error && <div className="mb-4 p-2 bg-red-200 text-red-800 rounded">{error}</div>}
            {success && <div className="mb-4 p-2 bg-green-200 text-green-800 rounded">{success}</div>}

            {/* FORM CREATE EVENT */}
            <h1 className="text-3xl font-bold mb-4 text-white">Buat Event</h1>
            <form onSubmit={handleSubmitCreate} className="max-w-3xl mx-auto mb-8">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    {/* Nama Event */}
                    <div>
                        <label htmlFor="product_name" className="block text-white">
                            Nama Event
                        </label>
                        <input
                            type="text"
                            id="product_name"
                            name="product_name"
                            className="w-full p-2 border border-gray-300 rounded mt-1 text-black"
                            value={createFormData.product_name}
                            onChange={handleChangeCreate}
                            required
                        />
                    </div>

                    {/* Alamat Event */}
                    <div>
                        <label htmlFor="product_address" className="block text-white">
                            Alamat Event
                        </label>
                        <input
                            type="text"
                            id="product_address"
                            name="product_address"
                            className="w-full p-2 border border-gray-300 rounded mt-1 text-black"
                            value={createFormData.product_address}
                            onChange={handleChangeCreate}
                            required
                        />
                    </div>

                    {/* Waktu Event */}
                    <div>
                        <label htmlFor="product_time" className="block text-white">
                            Waktu Event
                        </label>
                        <input
                            type="time"
                            id="product_time"
                            name="product_time"
                            className="w-full p-2 border border-gray-300 rounded mt-1 text-black"
                            value={createFormData.product_time}
                            onChange={handleChangeCreate}
                            required
                        />
                    </div>

                    {/* Tanggal Event */}
                    <div>
                        <label htmlFor="product_date" className="block text-white">
                            Tanggal Event
                        </label>
                        <input
                            type="date"
                            id="product_date"
                            name="product_date"
                            className="w-full p-2 border border-gray-300 rounded mt-1 text-black"
                            value={createFormData.product_date}
                            onChange={handleChangeCreate}
                            required
                        />
                    </div>

                    {/* Jumlah Tiket */}
                    <div>
                        <label htmlFor="product_quantity" className="block text-white">
                            Jumlah Tiket
                        </label>
                        <input
                            type="number"
                            id="product_quantity"
                            name="product_quantity"
                            className="w-full p-2 border border-gray-300 rounded mt-1 text-black"
                            value={createFormData.product_quantity || ""}
                            onChange={handleChangeCreate}
                            required
                        />
                    </div>

                    {/* Harga Tiket */}
                    <div>
                        <label htmlFor="product_price" className="block text-white">
                            Harga Tiket (Rp)
                        </label>
                        <input
                            type="number"
                            id="product_price"
                            name="product_price"
                            className="w-full p-2 border border-gray-300 rounded mt-1 text-black"
                            value={createFormData.product_price || ""}
                            onChange={handleChangeCreate}
                            required
                        />
                    </div>

                    {/* Deskripsi */}
                    <div className="col-span-2">
                        <label htmlFor="product_description" className="block text-white">
                            Deskripsi Event
                        </label>
                        <textarea
                            id="product_description"
                            name="product_description"
                            className="w-full p-2 border border-gray-300 rounded mt-1 text-black"
                            value={createFormData.product_description}
                            onChange={handleChangeCreate}
                            rows={4}
                            required
                        ></textarea>
                    </div>

                    {/* Kategori */}
                    <div className="col-span-2">
                        <label htmlFor="product_category" className="block text-white">
                            Kategori Event
                        </label>
                        <select
                            id="product_category"
                            name="product_category"
                            className="w-full p-2 border border-gray-300 rounded mt-1 text-black"
                            value={createFormData.product_category}
                            onChange={handleChangeCreate}
                            required
                        >
                            <option value="" disabled>
                                Pilih Kategori
                            </option>
                            {categoryOptions.map((option) => (
                                <option key={option} value={option}>
                                    {option}
                                </option>
                            ))}
                        </select>
                    </div>

                    {/* Pilih Gambar (Opsional) */}
                    {/* <div className="col-span-2">
                        <label className="block text-white">Upload Gambar (Opsional)</label>
                        <input
                            type="file"
                            accept="image/*"
                            onChange={handleCreateImageFileChange}
                            className="text-white"
                        />
                    </div> */}
                </div>

                <button
                    type="submit"
                    className={`w-full text-white p-2 rounded mt-4 transition ${
                        loading ? "bg-gray-400 cursor-not-allowed" : "bg-blue-500 hover:bg-blue-600"
                    }`}
                    disabled={loading}
                >
                    {loading ? "Sedang Membuat..." : "Buat Event"}
                </button>
            </form>

            {/* 
                STEP VERIFIKASI PAYMENT 
                => Muncul setelah event berhasil dibuat.
            */}
            {isVerificationStep && (
                <div className="bg-gray-800 p-4 rounded mb-8">
                    <h2 className="text-2xl font-bold text-white mb-2">
                        Verifikasi Payment
                    </h2>
                    <p className="text-white mb-4">
                        Masukkan 6 digit kode yang telah dikirim ke email/HP Anda.
                    </p>
                    <form onSubmit={handlePaymentVerification}>
                        <input
                            type="text"
                            maxLength={6}
                            className="p-2 border rounded mr-2 text-black"
                            placeholder="Kode Verifikasi"
                            value={verificationCode}
                            onChange={(e) => setVerificationCode(e.target.value)}
                            required
                        />
                        <button
                            type="submit"
                            className={`text-white p-2 rounded transition ${
                                loading
                                    ? "bg-gray-400 cursor-not-allowed"
                                    : "bg-green-500 hover:bg-green-600"
                            }`}
                            disabled={loading}
                        >
                            {loading ? "Verifying..." : "Verifikasi"}
                        </button>
                    </form>
                </div>
            )}

            {/* LIST SUBMISSIONS */}
            <div>
                <h2 className="text-2xl font-bold mb-4 text-white">Daftar Event Saya</h2>
                {loading && <p className="text-gray-300">Loading data...</p>}
                {submissions.length === 0 ? (
                    <p className="text-gray-200">Belum ada event.</p>
                ) : (
                    <ul className="space-y-4">
                        {submissions.map((submission) => (
                            <li
                                key={submission.id}
                                className="border p-4 rounded-lg shadow-sm bg-gray-100"
                            >
                                <h3 className="text-lg font-bold text-black">
                                    {submission.product_name}
                                </h3>

                                {/* Tampilkan gambar jika ada */}
                                {submission.product_image && (
                                    <img
                                        src={submission.product_image}
                                        alt="Event"
                                        className="w-full h-auto max-w-sm mt-2 mb-4"
                                    />
                                )}

                                <p className="text-black">Alamat: {submission.product_address}</p>
                                <p className="text-black">Tanggal: {submission.product_date}</p>
                                <p className="text-black">Waktu: {submission.product_time}</p>
                                <p className="text-black">
                                    Harga Tiket: Rp{" "}
                                    {submission.product_price.toLocaleString("id-ID")}
                                </p>
                                <p className="text-black">
                                    Jumlah Tiket: {submission.product_quantity}
                                </p>
                                <p className="text-black">Status: {submission.product_status}</p>

                                <div className="mt-2 flex gap-2">
                                    <button
                                        onClick={() => handleEditClick(submission)}
                                        className="bg-yellow-500 hover:bg-yellow-600 text-white px-3 py-1 rounded"
                                    >
                                        Edit
                                    </button>
                                    <button
                                        onClick={() => handleDeleteClick(submission.id)}
                                        className="bg-red-500 hover:bg-red-600 text-white px-3 py-1 rounded"
                                    >
                                        Hapus
                                    </button>
                                </div>
                            </li>
                        ))}
                    </ul>
                )}
            </div>

            {/* FORM EDIT (DATA TEKS) + UPLOAD GAMBAR */}
            {editSubmission && (
                <div className="mt-8 bg-gray-800 p-4 rounded">
                    <h2 className="text-2xl font-bold mb-4 text-white">
                        Edit Event ID: {editSubmission.id}
                    </h2>

                    {/* EDIT DATA TEKS */}
                    <form onSubmit={handleUpdateSubmission}>
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                            {/* Nama */}
                            <div>
                                <label htmlFor="edit_product_name" className="block text-white">
                                    Nama Event
                                </label>
                                <input
                                    type="text"
                                    id="edit_product_name"
                                    name="product_name"
                                    className="w-full p-2 border border-gray-300 rounded mt-1 text-black"
                                    value={editFormData.product_name}
                                    onChange={handleEditFormChange}
                                    required
                                />
                            </div>

                            {/* Alamat */}
                            <div>
                                <label
                                    htmlFor="edit_product_address"
                                    className="block text-white"
                                >
                                    Alamat Event
                                </label>
                                <input
                                    type="text"
                                    id="edit_product_address"
                                    name="product_address"
                                    className="w-full p-2 border border-gray-300 rounded mt-1 text-black"
                                    value={editFormData.product_address}
                                    onChange={handleEditFormChange}
                                    required
                                />
                            </div>

                            {/* Waktu */}
                            <div>
                                <label htmlFor="edit_product_time" className="block text-white">
                                    Waktu Event
                                </label>
                                <input
                                    type="time"
                                    id="edit_product_time"
                                    name="product_time"
                                    className="w-full p-2 border border-gray-300 rounded mt-1 text-black"
                                    value={editFormData.product_time}
                                    onChange={handleEditFormChange}
                                    required
                                />
                            </div>

                            {/* Tanggal */}
                            <div>
                                <label htmlFor="edit_product_date" className="block text-white">
                                    Tanggal Event
                                </label>
                                <input
                                    type="date"
                                    id="edit_product_date"
                                    name="product_date"
                                    className="w-full p-2 border border-gray-300 rounded mt-1 text-black"
                                    value={editFormData.product_date}
                                    onChange={handleEditFormChange}
                                    required
                                />
                            </div>

                            {/* Qty */}
                            <div>
                                <label
                                    htmlFor="edit_product_quantity"
                                    className="block text-white"
                                >
                                    Jumlah Tiket
                                </label>
                                <input
                                    type="number"
                                    id="edit_product_quantity"
                                    name="product_quantity"
                                    className="w-full p-2 border border-gray-300 rounded mt-1 text-black"
                                    value={editFormData.product_quantity}
                                    onChange={handleEditFormChange}
                                    required
                                />
                            </div>

                            {/* Harga */}
                            <div>
                                <label
                                    htmlFor="edit_product_price"
                                    className="block text-white"
                                >
                                    Harga Tiket (Rp)
                                </label>
                                <input
                                    type="number"
                                    id="edit_product_price"
                                    name="product_price"
                                    className="w-full p-2 border border-gray-300 rounded mt-1 text-black"
                                    value={editFormData.product_price}
                                    onChange={handleEditFormChange}
                                    required
                                />
                            </div>

                            {/* Deskripsi */}
                            <div className="col-span-2">
                                <label
                                    htmlFor="edit_product_description"
                                    className="block text-white"
                                >
                                    Deskripsi Event
                                </label>
                                <textarea
                                    id="edit_product_description"
                                    name="product_description"
                                    className="w-full p-2 border border-gray-300 rounded mt-1 text-black"
                                    value={editFormData.product_description}
                                    onChange={handleEditFormChange}
                                    rows={4}
                                    required
                                ></textarea>
                            </div>

                            {/* Kategori */}
                            <div className="col-span-2">
                                <label
                                    htmlFor="edit_product_category"
                                    className="block text-white"
                                >
                                    Kategori Event
                                </label>
                                <select
                                    id="edit_product_category"
                                    name="product_category"
                                    className="w-full p-2 border border-gray-300 rounded mt-1 text-black"
                                    value={editFormData.product_category}
                                    onChange={handleEditFormChange}
                                    required
                                >
                                    <option value="" disabled>
                                        Pilih Kategori
                                    </option>
                                    {categoryOptions.map((option) => (
                                        <option key={option} value={option}>
                                            {option}
                                        </option>
                                    ))}
                                </select>
                            </div>
                        </div>

                        <div className="mt-4 flex gap-2">
                            <button
                                type="submit"
                                className={`text-white p-2 rounded transition ${
                                    loading ? "bg-gray-400 cursor-not-allowed" : "bg-green-500 hover:bg-green-600"
                                }`}
                                disabled={loading}
                            >
                                {loading ? "Updating..." : "Update Event"}
                            </button>
                            <button
                                type="button"
                                onClick={handleCancelEdit}
                                className="text-white p-2 rounded bg-red-500 hover:bg-red-600"
                            >
                                Batal
                            </button>
                        </div>
                    </form>

                    {/* FORM UPLOAD GAMBAR (SETELAH EVENT EXIST) */}
                    <div className="mt-6">
                        <h3 className="text-xl font-bold text-white">
                            Upload/Replace Gambar Event
                        </h3>
                        <form onSubmit={handleUploadImage} className="mt-2">
                            <input
                                type="file"
                                onChange={handleFileChange}
                                accept="image/*"
                                className="text-white mb-2"
                            />
                            <button
                                type="submit"
                                className={`text-white p-2 rounded ${
                                    !editImageFile || loading
                                        ? "bg-gray-400 cursor-not-allowed"
                                        : "bg-blue-500 hover:bg-blue-600"
                                }`}
                                disabled={!editImageFile || loading}
                            >
                                {loading ? "Uploading..." : "Upload Gambar"}
                            </button>
                        </form>
                    </div>
                </div>
            )}
        </div>
    );
}
