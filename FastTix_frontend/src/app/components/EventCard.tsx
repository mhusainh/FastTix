// src/app/components/EventCard.tsx

import React from "react";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faMapMarkerAlt, faCalendarAlt } from "@fortawesome/free-solid-svg-icons";

interface EventCardProps {
  image: string;
  location: string;
  date: string;
  title: string;
  price: string;
  status: string;
  stock: number; // Stock quantity
}

const EventCard: React.FC<EventCardProps> = ({
  image,
  location,
  date,
  title,
  price,
  status,
  stock,
}) => {
  return (
    <div className="bg-white shadow-md rounded-lg overflow-hidden flex flex-col h-full">
      {/* Event Image */}
      <div className="h-60 overflow-hidden">
        <img
          src={image}
          alt={title}
          className="w-full h-full object-cover"
        />
      </div>

      {/* Event Details */}
      <div className="flex flex-col flex-grow p-2">
        {/* Location and Date */}
        <div className="flex flex-col items-start text-gray-500 text-sm mb-2">
          <div className="flex items-center">
            <FontAwesomeIcon icon={faMapMarkerAlt} className="mr-1" />
            <span>{location}</span>
          </div>
          <div className="flex items-center">
            <FontAwesomeIcon icon={faCalendarAlt} className="mr-1" />
            <span>{date}</span>
          </div>
        </div>

        {/* Title */}
        <h3 className="text-sm font-semibold text-gray-800">{title}</h3>

        <div className="flex-grow"></div>

        {/* Stock */}
        <p className="text-gray-500 text-sm mt-2">
          <strong>Stok:</strong> {stock}
        </p>

        {/* Price and Status */}
        <div className="flex items-center text-sm justify-between mt-2">
          <p className="text-red-500 font-bold">{price}</p>
          <span
            className={`${
              status === "Tiket Tersedia"
                ? "text-green-500"
                : status === "Tiket Tidak Tersedia"
                ? "text-red-500"
                : "text-yellow-500"
            } text-sm font-medium`}
          >
            {status}
          </span>
        </div>
      </div>
    </div>
  );
};

export default EventCard;
