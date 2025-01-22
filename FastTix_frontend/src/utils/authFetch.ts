// src/utils/authFetch.ts

import { useContext } from "react";
import { AuthContext } from "../context/AuthContext";
import { useRouter } from "next/navigation";

export const useAuthFetch = () => {
  const { token, logout } = useContext(AuthContext);
  const router = useRouter();

  const authFetch = async (url: string, options: RequestInit = {}) => {
    const headers: HeadersInit = options.headers || {};

    if (token) {
      headers["Authorization"] = `Bearer ${token}`;
    }

    try {
      const response = await fetch(url, { ...options, headers });

      if (response.status === 401) {
        // Token tidak valid atau kedaluwarsa
        logout();
        router.push("/login");
        return;
      }

      return response;
    } catch (error) {
      console.error("Fetch error:", error);
      throw error;
    }
  };

  return authFetch;
};
