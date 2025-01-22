"use client";

import React, { useEffect, useState } from "react";
import EventCard from "../components/EventCard";
import { useSearchParams } from "next/navigation";

interface BackendEvent {
  id: number;
  product_name: string;
  product_address: string;
  product_date: string;
  product_price: number;
  product_status: string;
  product_image: string;
  product_description: string;
  product_quantity: number;
}

interface EventCardProps {
  image: string;
  location: string;
  date: string;
  title: string;
  price: string;
  status: string;
  stock: number;
}

const SearchResultsPage: React.FC = () => {
  const searchParams = useSearchParams();
  const query = searchParams.get("query") || ""; // Get the search query from the URL

  const [events, setEvents] = useState<EventCardProps[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchSearchResults = async () => {
      setLoading(true);
      setError(null);
      // console.log(query)
      try {
        const response = await fetch(
          `http://localhost:8080/api/v1/tickets?search=${encodeURIComponent(query)}`
        );
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }

        const result = await response.json();

        const mappedEvents: EventCardProps[] = result.data.map(
          (event: BackendEvent) => ({
            image: event.product_image || "/default-image.jpg",
            location: event.product_address,
            date: new Date(event.product_date).toLocaleDateString("id-ID", {
              day: "2-digit",
              month: "short",
              year: "numeric",
            }),
            title: event.product_name,
            price: `Rp ${event.product_price.toLocaleString("id-ID")}`,
            status: event.product_status === "available" ? "Tiket Tersedia" : "Tiket Tidak Tersedia",
          })
        );

        setEvents(mappedEvents);
      } catch (err: any) {
        setError(err.message || "Terjadi kesalahan saat mengambil data.");
      } finally {
        setLoading(false);
      }
    };

    if (query) {
      fetchSearchResults();
    }
  }, [query]);

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-4">Hasil Pencarian</h1>
      <p className="text-gray-600 mb-6">
        Menampilkan hasil untuk pencarian: <strong>{query}</strong>
      </p>

      {loading && <p>Loading...</p>}
      {error && <p className="text-red-500">{error}</p>}

      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
        {events.map((event, index) => (
          <EventCard key={index} {...event} />
        ))}
      </div>
    </div>
  );
};

export default SearchResultsPage;
