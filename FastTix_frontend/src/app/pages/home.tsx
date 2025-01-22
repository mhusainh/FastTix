// src/app/pages/home.tsx

"use client"; // Add this directive at the top

import React, { useEffect, useState } from "react";
import EventCard from "../components/EventCard";

interface BackendEvent {
  id: number;
  product_name: string;
  product_address: string;
  product_time: string;
  product_date: string;
  product_price: number;
  product_description: string;
  product_category: string;
  product_quantity: number;
  product_type: string;
  product_status: string;
  user_id: number;
  order_id: string;
  created_at: string;
  updated_at: string;
  // If your backend provides an image URL, include it here
  image_url?: string;
  
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

const HomePage: React.FC = () => {
  const [events, setEvents] = useState<EventCardProps[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchEvents = async () => {
      try {
        const response = await fetch("http://localhost:8080/api/v1/tickets");
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }
        const result = await response.json();

        // Map backend data to EventCard props
        const mappedEvents: EventCardProps[] = result.data.map(
          (event: BackendEvent) => ({
            image: getImageUrl(event), // Function to determine image URL
            location: event.product_address,
            date: formatDate(event.product_date),
            title: event.product_name,
            price: formatPrice(event.product_price),
            status: mapStatus(event.product_status),
          })
        );

        setEvents(mappedEvents);
        setLoading(false);
      } catch (err: any) {
        setError(err.message || "An error occurred while fetching events.");
        setLoading(false);
      }
    };

    fetchEvents();
  }, []);

  // Helper function to format the date
  const formatDate = (isoDate: string): string => {
    const date = new Date(isoDate);
    const options: Intl.DateTimeFormatOptions = {
      day: "2-digit",
      month: "short",
      year: "2-digit",
    };
    return date.toLocaleDateString("id-ID", options); // Indonesian locale
  };

  // Helper function to format the price
  const formatPrice = (price: number): string => {
    return `Rp. ${price.toLocaleString("id-ID")}`;
  };

  // Helper function to map backend status to frontend status
  const mapStatus = (status: string): string => {
    switch (status.toLowerCase()) {
      case "available":
        return "Tiket Tersedia";
      case "sold_out":
        return "Tiket Tidak Tersedia";
      case "pending":
        return "Pending";
      default:
        return "Unknown Status";
    }
  };

  // Helper function to determine the image URL
  const getImageUrl = (event: BackendEvent): string => {
    // If your backend provides an image URL, use it
    if (event.image_url) {
      return event.image_url;
    }

    // Otherwise, assign a default image or map based on category/type
    // Ensure that /default-event.jpg exists in your public directory
    return "../tes.jpg";
  };

  if (loading) {
    return (
      <div className="container mx-auto px-4 py-8">
        <h2 className="text-2xl font-bold mb-6">Event Terbaru</h2>
        <p>Loading events...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="container mx-auto px-4 py-8">
        <h2 className="text-2xl font-bold mb-6">Event Terbaru</h2>
        <p className="text-red-500">Error: {error}</p>
      </div>
    );
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <h2 className="text-2xl font-bold mb-6">Event Terbaru</h2>
      {events.length === 0 ? (
        <p>No events available.</p>
      ) : (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
          {events.map((event, index) => (
            <EventCard key={index} {...event} />
          ))}
        </div>
      )}
    </div>
  );
};

export default HomePage;
