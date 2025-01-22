// src/app/event/page.tsx

"use client";

import React, { useEffect, useState, useContext, useRef } from "react";
import EventCard from "../components/EventCard";
import { AuthContext } from "../../context/AuthContext";
import { useRouter } from "next/navigation";

interface BackendEvent {
  id: number;
  product_name: string;
  product_address: string;
  product_date: string;
  product_time: string; // Added time
  product_price: number;
  product_status: string;
  product_image: string;
  product_description: string;
  product_quantity: number; // Stock quantity
}

interface EventCardProps {
  image: string;
  location: string;
  date: string;
  time: string; // Added time
  title: string;
  price: string;
  status: string;
  id: number;
  description: string;
  stock: number; // Stock quantity
}

const EventPage: React.FC = () => {
  const { token, logout } = useContext(AuthContext);
  const router = useRouter();

  const [events, setEvents] = useState<EventCardProps[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [sortBy, setSortBy] = useState<string>("id");
  const [order, setOrder] = useState<string>("desc");
  const [page, setPage] = useState<number>(1);
  const [limit, setLimit] = useState<number>(8);
  // Added hasMore state to determine if more pages are available
  const [hasMore, setHasMore] = useState<boolean>(true);
  const [selectedEvent, setSelectedEvent] = useState<EventCardProps | null>(null);
  const [transactionQuantity, setTransactionQuantity] = useState<number>(1);
  const [showCheckout, setShowCheckout] = useState<boolean>(false);

  // Additional state for verification
  const [showVerification, setShowVerification] = useState<boolean>(false);
  const [verificationCode, setVerificationCode] = useState<string[]>(["", "", "", "", "", ""]);
  const [isProcessingCheckout, setIsProcessingCheckout] = useState<boolean>(false);
  const [isRedirecting, setIsRedirecting] = useState<boolean>(false); // Loading state during redirect
  const [verificationError, setVerificationError] = useState<string | null>(null); // Verification error message

  // Refs for verification inputs
  const verificationRefs = useRef<Array<HTMLInputElement | null>>([]);

  useEffect(() => {
    const fetchEvents = async () => {
      setLoading(true);
      setError(null);

      try {
        const response = await fetch(
          `http://localhost:8080/api/v1/tickets?sort=${sortBy}&order=${order}&page=${page}&limit=${limit}`,
          {
            headers: {
              "Authorization": token ? `Bearer ${token}` : "",
            },
          }
        );

        if (response.status === 401) {
          // Token is invalid or expired
          logout();
          router.push("/login");
          return;
        }

        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }

        const result = await response.json();
        const mappedEvents: EventCardProps[] = result.data.map(
          (event: BackendEvent) => ({
            id: event.id,
            image: event.product_image && event.product_image.trim() !== "" ? event.product_image : "../tes.jpg",
            location: event.product_address,
            date: formatDate(event.product_date),
            time: formatTime(event.product_time), // Map time
            title: event.product_name,
            price: formatPrice(event.product_price),
            status: event.product_status === "accepted" ? "Available" : event.product_status, // Status mapping
            description: event.product_description,
            stock: event.product_quantity,
          })
        );
        

        setEvents(mappedEvents);
        // Determine if there are more pages
        setHasMore(result.data.length === limit);
      } catch (err: any) {
        setError(err.message || "Terjadi kesalahan saat mengambil data.");
      } finally {
        setLoading(false);
      }
    };

    fetchEvents();
  }, [sortBy, order, page, limit, token, logout, router]);

  const formatDate = (isoDate: string): string => {
    const date = new Date(isoDate);
    return date.toLocaleDateString("id-ID", { day: "2-digit", month: "short", year: "numeric" });
  };

  const formatTime = (isoTime: string): string => {
    const date = new Date(`1970-01-01T${isoTime}Z`); // Assume product_time is a time string in ISO format
    return date.toLocaleTimeString("id-ID", { hour: "2-digit", minute: "2-digit" });
  };

  const formatPrice = (price: number): string => {
    return `Rp ${price.toLocaleString("id-ID")}`;
  };

  const handleSortChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const [sortKey, sortOrder] = e.target.value.split("-");
    setSortBy(sortKey);
    setOrder(sortOrder);
    setPage(1); // Reset to page 1 when sorting changes
  };

  const handlePageChange = (direction: "next" | "prev") => {
    setPage((prevPage) =>
      direction === "next" ? prevPage + 1 : Math.max(prevPage - 1, 1)
    );
  };

  const handleCardClick = (event: EventCardProps) => {
    setSelectedEvent(event);
    setTransactionQuantity(1); // Reset quantity when a new event is selected
  };

  const closePreviewModal = () => {
    setSelectedEvent(null);
  };

  const handlePesanSekarang = () => {
    setShowCheckout(true);
  };

  const closeCheckoutModal = () => {
    setShowCheckout(false);
    setSelectedEvent(null);
  };

  // Function to handle checkout
  const handleCheckout = async () => {
    if (!selectedEvent) return;
    setIsProcessingCheckout(true);

    try {
      const response = await fetch(
        `http://localhost:8080/api/v1/tickets/${selectedEvent.id}/checkout`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            "Authorization": token ? `Bearer ${token}` : "",
          },
          body: JSON.stringify({ transaction_quantity: transactionQuantity }),
        }
      );

      if (response.status === 401) {
        // Token is invalid or expired
        logout();
        router.push("/login");
        return;
      }

      if (!response.ok) {
        const result = await response.json();
        throw new Error(result.meta.message || "Gagal melakukan checkout.");
      }

      setIsProcessingCheckout(false);
      setShowCheckout(false);
      setShowVerification(true);
      setVerificationCode(["", "", "", "", "", ""]);
      setVerificationError(null);
      // Focus on the first input
      setTimeout(() => {
        verificationRefs.current[0]?.focus();
      }, 100);
    } catch (err: any) {
      alert(err.message || "Terjadi kesalahan saat checkout.");
      setIsProcessingCheckout(false);
    }
  };

  // Function to handle verification input
  const handleVerificationInput = (index: number, value: string) => {
    if (!/^[a-zA-Z0-9]?$/.test(value)) return; // Allow only letters and numbers

    const updatedCode = [...verificationCode];
    updatedCode[index] = value.toUpperCase().slice(0, 1); // Convert to uppercase
    setVerificationCode(updatedCode);
    setVerificationError(null); // Reset error on input

    if (value && index < verificationCode.length - 1) {
      verificationRefs.current[index + 1]?.focus();
    }
  };

  // Function to handle backspace and navigation in verification inputs
  const handleVerificationKeyDown = (e: React.KeyboardEvent<HTMLInputElement>, index: number) => {
    if (e.key === "Backspace" && !verificationCode[index] && index > 0) {
      const updatedCode = [...verificationCode];
      updatedCode[index - 1] = "";
      setVerificationCode(updatedCode);
      verificationRefs.current[index - 1]?.focus();
      e.preventDefault();
    }
  };

  // Function to handle verification
  const handleVerify = async () => {
    const verificationCodeStr = verificationCode.join("");

    if (verificationCodeStr.length !== 6) {
      setVerificationError("Silakan masukkan kode verifikasi 6 karakter (huruf/angka).");
      return;
    }

    setIsRedirecting(true); // Start loading indicator
    setVerificationError(null); // Reset error before verification

    try {
      const response = await fetch(
        `http://localhost:8080/api/v1/payment/checkout/${verificationCodeStr}`,
        {
          method: "GET",
          headers: {
            "Content-Type": "application/json",
            "Authorization": token ? `Bearer ${token}` : "",
          },
        }
      );

      if (response.status === 401) {
        // Token is invalid or expired
        logout();
        router.push("/login");
        return;
      }

      if (!response.ok) {
        const result = await response.json();
        throw new Error(result.meta.message || "Verifikasi gagal.");
      }

      const result = await response.json();

      if (result.data && result.data.redirectURL) {
        // Redirect to Midtrans URL
        window.location.href = result.data.redirectURL;
      } else {
        throw new Error("URL redirect pembayaran tidak tersedia.");
      }
    } catch (err: any) {
      setVerificationError(err.message || "Terjadi kesalahan saat verifikasi.");
    } finally {
      setIsRedirecting(false); // Stop loading indicator
      // Do not close the verification UI if failed
    }
  };

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-6">Event</h1>

      {/* Sorting and Pagination Controls */}
      <div className="mb-6 flex flex-col md:flex-row justify-between items-center gap-4">
        <div className="flex flex-col">
          <label htmlFor="sort" className="text-gray-700 mb-2">
            Urutkan Berdasarkan
          </label>
          <select
            id="sort"
            className="p-2 border border-gray-300 rounded text-gray-600"
            onChange={handleSortChange}
            value={`${sortBy}-${order}`}
          >
            <option value="id-desc">Terbaru</option>
            <option value="id-asc">Terlama</option>
            <option value="product_price-asc">Termurah</option>
            <option value="product_price-desc">Termahal</option>
            <option value="product_date-asc">Tanggal Terdekat</option>
            <option value="product_date-desc">Tanggal Terjauh</option>
          </select>
        </div>

        <div className="flex flex-col">
          <label htmlFor="limit" className="text-gray-700 mb-2">
            Jumlah Per Halaman
          </label>
          <select
            id="limit"
            className="p-2 border border-gray-300 rounded text-gray-600"
            value={limit}
            onChange={(e) => setLimit(parseInt(e.target.value))}
          >
            <option value={4}>4</option>
            <option value={8}>8</option>
            <option value={12}>12</option>
          </select>
        </div>
      </div>

      {/* Event Cards */}
      {loading ? (
        <div className="flex justify-center items-center h-64">
          <div className="loader ease-linear rounded-full border-4 border-t-4 border-gray-200 h-12 w-12"></div>
        </div>
      ) : error ? (
        <div className="text-red-500 text-center">{error}</div>
      ) : events.length === 0 ? (
        <div className="text-center text-gray-700">Tidak ada event untuk ditampilkan.</div>
      ) : (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
          {events.map((event) => (
            <div
              key={event.id}
              onClick={() => handleCardClick(event)}
              className="cursor-pointer transform hover:scale-105 transition-transform duration-200"
            >
              <EventCard {...event} />
            </div>
          ))}
        </div>
      )}

      {/* Pagination Controls */}
      <div className="mt-6 flex justify-center items-center gap-4">
        <button
          onClick={() => handlePageChange("prev")}
          disabled={page === 1}
          className={`px-4 py-2 bg-gray-300 text-gray-700 rounded ${
            page === 1 ? "opacity-50 cursor-not-allowed" : "hover:bg-gray-400"
          }`}
        >
          Previous
        </button>
        <span className="text-gray-700">
          Halaman {page}
        </span>
        <button
          onClick={() => handlePageChange("next")}
          disabled={!hasMore}
          className={`px-4 py-2 bg-gray-300 text-gray-700 rounded ${
            !hasMore ? "opacity-50 cursor-not-allowed" : "hover:bg-gray-400"
          }`}
        >
          Next
        </button>
      </div>

      {/* Event Preview Modal */}
      {selectedEvent && !showCheckout && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white p-6 rounded-lg w-full max-w-4xl shadow-lg relative flex flex-col md:flex-row">
            {/* Image */}
            <div className="md:w-1/3 pr-4 mb-4 md:mb-0">
              <img
                src={selectedEvent.image}
                alt={selectedEvent.title}
                className="rounded-lg w-full h-full object-cover"
              />
            </div>

            {/* Details */}
            <div className="md:w-2/3">
              <button
                onClick={closePreviewModal}
                className="absolute top-2 right-2 text-gray-500 hover:text-gray-800 text-2xl"
              >
                &times;
              </button>
              <h2 className="text-2xl font-bold mb-4">{selectedEvent.title}</h2>
              <p className="mb-4 text-gray-600">{selectedEvent.description}</p>
              <div className="mb-2">
                <p className="text-gray-700"><strong>Tanggal:</strong> {selectedEvent.date}</p>
                <p className="text-gray-700"><strong>Waktu:</strong> {selectedEvent.time}</p> {/* Display time */}
                <p className="text-gray-700"><strong>Lokasi:</strong> {selectedEvent.location}</p>
                <p className="text-gray-700"><strong>Harga:</strong> {selectedEvent.price}</p>
                <p className="text-gray-700"><strong>Stok:</strong> {selectedEvent.stock}</p>
              </div>
              <button
                onClick={handlePesanSekarang}
                className="mt-4 bg-blue-500 text-white py-2 px-4 rounded hover:bg-blue-600 transition-colors duration-200 w-full"
              >
                Pesan Sekarang
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Checkout Modal */}
      {showCheckout && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white p-6 rounded-lg w-full max-w-md shadow-lg relative">
            <button
              onClick={closeCheckoutModal}
              className="absolute top-2 right-2 text-gray-500 hover:text-gray-800 text-2xl"
            >
              &times;
            </button>
            <h2 className="text-2xl font-bold mb-4 text-gray-800">Checkout</h2>
            <div className="border-t border-b py-4 mb-4">
              <div className="mb-2">
                <p className="text-gray-700 font-semibold">Event Name:</p>
                <p className="text-gray-500">{selectedEvent?.title}</p>
              </div>
              <div className="mb-2">
                <p className="text-gray-700 font-semibold">Tanggal:</p>
                <p className="text-gray-500">{selectedEvent?.date}</p>
              </div>
              <div className="mb-2">
                <p className="text-gray-700 font-semibold">Waktu:</p> {/* Display time */}
                <p className="text-gray-500">{selectedEvent?.time}</p>
              </div>
              <div>
                <p className="text-gray-700 font-semibold">Lokasi:</p>
                <p className="text-gray-500">{selectedEvent?.location}</p>
              </div>
            </div>
            <div className="flex items-center gap-4 mb-4">
              <label htmlFor="quantity" className="block text-gray-700 font-semibold">Quantity:</label>
              <input
                type="number"
                id="quantity"
                min={1}
                max={selectedEvent?.stock}
                value={transactionQuantity}
                onChange={(e) =>
                  setTransactionQuantity(
                    Math.min(parseInt(e.target.value) || 1, selectedEvent?.stock || 1)
                  )
                }
                className="border p-2 w-20 text-center rounded"
              />
            </div>
            <div className="flex justify-between items-center mb-4">
              <span className="text-lg font-bold text-gray-800">Total Price:</span>
              <span className="text-lg font-bold text-blue-500">
                {selectedEvent && formatPrice(
                  parseInt(selectedEvent.price.replace(/[^\d]/g, "")) * transactionQuantity
                )}
              </span>
            </div>
            <button
              onClick={handleCheckout}
              disabled={isProcessingCheckout}
              className={`${
                isProcessingCheckout ? "bg-blue-300 cursor-not-allowed" : "bg-blue-500 hover:bg-blue-600"
              } text-white p-3 rounded w-full font-semibold transition-colors duration-200`}
            >
              {isProcessingCheckout ? "Processing..." : "Confirm and Checkout"}
            </button>
          </div>
        </div>
      )}

      {/* Verification Modal */}
      {showVerification && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 px-4">
          <div className="bg-white p-6 rounded-lg w-full max-w-md shadow-lg relative">
            <h2 className="text-2xl font-bold mb-4">Verifikasi Pembayaran</h2>
            <p className="mb-4 text-gray-700">Masukkan kode verifikasi 6 karakter (huruf/angka) yang Anda terima:</p>
            {verificationError && (
              <div className="mb-4 text-red-500 text-center">
                {verificationError}
              </div>
            )}
            <div className="flex gap-2 justify-center mb-4">
              {verificationCode.map((code, index) => (
                <input
                  key={index}
                  id={`verification-input-${index}`}
                  type="text"
                  inputMode="text"
                  value={code}
                  onChange={(e) => handleVerificationInput(index, e.target.value)}
                  onKeyDown={(e) => handleVerificationKeyDown(e, index)}
                  maxLength={1}
                  ref={(el) => (verificationRefs.current[index] = el)}
                  className="border p-2 w-12 text-center text-lg font-bold rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
                  placeholder="-"
                />
              ))}
            </div>
            {isRedirecting && (
              <p className="text-blue-500 mb-4 text-center">Mengalihkan ke halaman pembayaran...</p>
            )}
            <button
              onClick={handleVerify}
              disabled={isRedirecting}
              className={`${
                isRedirecting ? "bg-green-300 cursor-not-allowed" : "bg-green-500 hover:bg-green-600"
              } text-white p-2 rounded w-full font-semibold transition-colors duration-200`}
            >
              {isRedirecting ? "Redirecting..." : "Verifikasi"}
            </button>
          </div>
        </div>
      )}
    </div>
  );
};

export default EventPage;
