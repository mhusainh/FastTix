"use client";

import React, { useEffect, useState } from "react";
import EventCard from "../components/EventCard";

interface BackendEvent {
  id: number;
  product_name: string;
  product_address: string;
  product_date: string;
  product_price: number;
  product_status: string;
  image_url: string;
  product_description: string;
}

interface EventCardProps {
  image: string;
  location: string;
  date: string;
  title: string;
  price: string;
  status: string;
  id: number;
  description: string;
}

const EventPage: React.FC = () => {
  const [events, setEvents] = useState<EventCardProps[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [sortBy, setSortBy] = useState<string>("id");
  const [order, setOrder] = useState<string>("desc");
  const [page, setPage] = useState<number>(1);
  const [limit, setLimit] = useState<number>(8);
  const [selectedEvent, setSelectedEvent] = useState<EventCardProps | null>(
    null
  );
  const [transactionQuantity, setTransactionQuantity] = useState<number>(1);

  useEffect(() => {
    const fetchEvents = async () => {
      setLoading(true);
      setError(null);

      try {
        const response = await fetch(
          `http://localhost:8080/api/v1/submissions?sort=${sortBy}&order=${order}&page=${page}&limit=${limit}`
        );
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }
        const result = await response.json();

        const mappedEvents: EventCardProps[] = result.data.map(
          (event: BackendEvent) => ({
            id: event.id,
            image: event.image_url || "/default-image.jpg",
            location: event.product_address,
            date: formatDate(event.product_date),
            title: event.product_name,
            price: formatPrice(event.product_price),
            status: mapStatus(event.product_status),
            description: event.product_description,
          })
        );

        setEvents(mappedEvents);
      } catch (err: any) {
        setError(err.message || "Terjadi kesalahan saat mengambil data.");
      } finally {
        setLoading(false);
      }
    };

    fetchEvents();
  }, [sortBy, order, page, limit]);

  const formatDate = (isoDate: string): string => {
    const date = new Date(isoDate);
    const options: Intl.DateTimeFormatOptions = {
      day: "2-digit",
      month: "short",
      year: "numeric",
    };
    return date.toLocaleDateString("id-ID", options);
  };

  const formatPrice = (price: number): string => {
    return `Rp ${price.toLocaleString("id-ID")}`;
  };

  const mapStatus = (status: string): string => {
    switch (status.toLowerCase()) {
      case "available":
        return "Available";
      case "unpaid":
        return "Tiket Tidak Tersedia";
      default:
        return "Pending";
    }
  };

  const handleSortChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const [sortKey, sortOrder] = e.target.value.split("-");
    setSortBy(sortKey);
    setOrder(sortOrder);
    setPage(1);
  };

  const handlePageChange = (direction: "next" | "prev") => {
    setPage((prevPage) =>
      direction === "next" ? prevPage + 1 : Math.max(prevPage - 1, 1)
    );
  };

  const handleCardClick = (event: EventCardProps) => {
    setSelectedEvent(event);
  };

  const handleCheckout = async () => {
    if (!selectedEvent) return;

    try {
      const response = await fetch(
        `http://localhost:8080/api/v1/tickets/${selectedEvent.id}/checkout`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({ transaction_quantity: transactionQuantity }),
        }
      );

      if (!response.ok) {
        const result = await response.json();
        throw new Error(result.message || "Gagal melakukan checkout.");
      }

      alert("Checkout berhasil!");
      setSelectedEvent(null);
    } catch (err: any) {
      alert(err.message || "Terjadi kesalahan saat checkout.");
    }
  };

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-4 text-white">Event</h1>

      {/* Sort Dropdown */}
      <div className="mb-6 flex justify-between items-center">
        <div>
          <label htmlFor="sort" className="block text-gray-700 mb-2">
            Urutkan Berdasarkan
          </label>
          <select
            id="sort"
            className="p-2 border border-gray-300 rounded"
            onChange={handleSortChange}
          >
            <option value="id-desc">Terbaru</option>
            <option value="id-asc">Terlama</option>
            <option value="product_price-asc">Termurah</option>
            <option value="product_price-desc">Termahal</option>
            <option value="product_date-asc">Tanggal Terdekat</option>
            <option value="product_date-desc">Tanggal Terjauh</option>
          </select>
        </div>

        <div>
          <label htmlFor="limit" className="block text-gray-700 mb-2">
            Jumlah Per Halaman
          </label>
          <select
            id="limit"
            className="p-2 border border-gray-300 rounded"
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
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
        {events.map((event, index) => (
          <div
            key={index}
            onClick={() => handleCardClick(event)}
            className="cursor-pointer"
          >
            <EventCard {...event} />
          </div>
        ))}
      </div>

      {/* Pagination Controls */}
      <div className="mt-6 flex justify-between">
        <button
          onClick={() => handlePageChange("prev")}
          disabled={page === 1}
          className="px-4 py-2 bg-gray-300 text-gray-700 rounded disabled:opacity-50"
        >
          Previous
        </button>
        <span className="text-gray-700">Halaman {page}</span>
        <button
          onClick={() => handlePageChange("next")}
          className="px-4 py-2 bg-gray-300 text-gray-700 rounded"
        >
          Next
        </button>
      </div>

      {/* Popup Dialog */}
      {selectedEvent && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white p-6 rounded-lg w-1/3 shadow-lg">
            <h2 className="text-xl font-bold mb-4">{selectedEvent.title}</h2>
            <p className="mb-4">{selectedEvent.description}</p>
            <label className="block mb-2 text-gray-700">Jumlah Tiket:</label>
            <input
              type="number"
              min={1}
              value={transactionQuantity}
              onChange={(e) =>
                setTransactionQuantity(parseInt(e.target.value) || 1)
              }
              className="p-2 border border-gray-300 rounded w-full mb-4"
            />
            <button
              onClick={handleCheckout}
              className="bg-blue-500 text-white p-2 rounded w-full"
            >
              Checkout
            </button>
            <button
              onClick={() => setSelectedEvent(null)}
              className="bg-red-500 text-white p-2 rounded w-full mt-4"
            >
              Tutup
            </button>
          </div>
        </div>
      )}
    </div>
  );
};

export default EventPage;
