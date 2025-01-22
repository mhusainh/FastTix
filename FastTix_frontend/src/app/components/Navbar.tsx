// src/app/components/Navbar.tsx
"use client";

import React, { useContext, useState } from "react";
import Link from "next/link";
import { AuthContext } from "../../context/AuthContext";
import { useRouter } from "next/navigation";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faSearch, faFilter } from "@fortawesome/free-solid-svg-icons";

const Navbar: React.FC = () => {
  const { isAuthenticated, logout, role } = useContext(AuthContext); // AMBIL role
  const [searchTerm, setSearchTerm] = useState<string>("");
  const [categoryDropdown, setCategoryDropdown] = useState<boolean>(false);
  const router = useRouter();

  const categories = [
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

  // const handleLogout = () => {
  //   logout();
  //   router.push("/login"); // Redirect to login page after logout
  // };

  const handleSearch = async (e: React.FormEvent) => {
    e.preventDefault(); // Prevent form submission from reloading the page
    if (searchTerm.trim()) {
      const query = encodeURIComponent(searchTerm.trim());
      const apiUrl = `http://localhost:8080/api/v1/tickets?search=${query}`;

      try {
        const response = await fetch(apiUrl);
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }

        const result = await response.json();
        console.log("Search results:", result);
        router.push(`/search-results?query=${query}`); // Redirect to a search results page
      } catch (error) {
        console.error("Error during search:", error);
      }
    }
  };

  const handleCategorySelect = async (category: string) => {
    const query = encodeURIComponent(category);
    const apiUrl = `http://localhost:8080/api/v1/tickets?search=${query}`;

    try {
      const response = await fetch(apiUrl);
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const result = await response.json();
      console.log("Category search results:", result);
      router.push(`/search-results?query=${query}`); // Redirect to a search results page
    } catch (error) {
      console.error("Error during category search:", error);
    }
  };

  return (
    <div className="navbar-container bg-blue-900 text-white py-4 px-8 flex justify-between items-center shadow-md">
      {/* Logo */}
      <div className="flex items-center gap-4">
        <Link href="/" className="text-white text-2xl font-bold no-underline">
          Fast <span className="text-yellow-500">Tix</span>
        </Link>

        {/* Link Event */}
        <Link
          href="/event"
          className="text-white text-lg hover:text-gray-300 no-underline"
        >
          Event
        </Link>

        {/* Link Buat Event (HANYA MUNCUL JIKA BUKAN ADMIN) */}
        {role !== "Administrator" && (
          <Link
            href="/buat-event"
            className="text-white text-lg hover:text-gray-300 no-underline"
          >
            Buat Event
          </Link>
        )}
      </div>

      <form
        onSubmit={handleSearch}
        className="search-container flex items-center bg-white rounded-full px-4 py-2 shadow-inner w-1/3 relative"
      >
        <input
          type="text"
          placeholder="Cari event dan atraksi di sini ..."
          className="flex-grow outline-none text-black text-sm"
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
        />
        <button
          type="submit"
          className="ml-2 text-blue-500 hover:text-blue-700"
        >
          <FontAwesomeIcon icon={faSearch} />
        </button>

        <div className="relative ml-4">
          <button
            type="button"
            className="text-blue-500 hover:text-blue-700"
            onClick={() => setCategoryDropdown((prev) => !prev)}
          >
            <FontAwesomeIcon icon={faFilter} />
          </button>

          {categoryDropdown && (
            <div className="absolute right-0 mt-2 w-48 bg-white text-black shadow-md rounded-lg z-50">
              {categories.map((category) => (
                <div
                  key={category}
                  className="px-4 py-2 cursor-pointer hover:bg-gray-200"
                  onClick={() => {
                    handleCategorySelect(category);
                    setCategoryDropdown(false);
                  }}
                >
                  {category}
                </div>
              ))}
            </div>
          )}
        </div>
      </form>

      {/* Navigation Items (Login, Profile, Logout, dll) */}
      <div className="flex items-center gap-6">

        {role !== "Administrator" && isAuthenticated && (
          <Link
            href="/profile"
            className="text-white text-lg hover:text-gray-300 no-underline"
          >
            Profile
          </Link>
        )}

        {/* Jika admin => tampilkan link Admin Panel */}
        {role === "Administrator" && isAuthenticated && (
          <Link
            href="/admin"
            className="text-white text-lg hover:text-gray-300 no-underline"
          >
            Admin Panel
          </Link>
        )}
        {isAuthenticated ? (
          <>
            {/* <Link
              href="/profile"
              className="text-white text-lg hover:text-gray-300 no-underline"
            >
              Profile
            </Link> */}
            {/* 
              Contoh: misal Anda ingin menambahkan link "Lihat Submissions" 
              untuk SEMUA ROLE, atau hanya admin => modifikasi di sini
            */}
            <Link
              href="/lihat-event-submission" 
              className="text-white text-lg hover:text-gray-300 no-underline"
            >
              Lihat Event Submission
            </Link>

            <button
              onClick={logout}
              className="text-white text-lg hover:text-gray-300"
            >
              Logout
            </button>
          </>
        ) : (
          <>
            <Link
              href="/login"
              className="text-white text-lg hover:text-gray-300 no-underline"
            >
              Masuk
            </Link>
            <Link
              href="/register"
              className="text-white text-lg hover:text-gray-300 no-underline"
            >
              Daftar
            </Link>
          </>
        )}
      </div>
    </div>
  );
};

export default Navbar;