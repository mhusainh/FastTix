// src/context/AuthContext.tsx
"use client";

import React, { createContext, useState, useEffect, ReactNode } from "react";
import { useRouter } from "next/navigation";

interface AuthContextProps {
  isAuthenticated: boolean;
  token: string | null;
  role: string | null;           // TAMBAHKAN
  authLoading: boolean;
  login: (token: string) => void;
  logout: () => void;
}

export const AuthContext = createContext<AuthContextProps>({
  isAuthenticated: false,
  token: null,
  role: null,                    // default null
  authLoading: true,
  login: () => {},
  logout: () => {},
});

export const AuthProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [token, setToken] = useState<string | null>(null);
  const [role, setRole] = useState<string | null>(null); // TAMBAHKAN
  const [authLoading, setAuthLoading] = useState<boolean>(true);
  const router = useRouter();

  // Ambil token dari localStorage (jika ada) saat pertama kali load
  useEffect(() => {
    const storedToken = localStorage.getItem("authToken");
    if (storedToken) {
      setToken(storedToken);

      // Setelah dapat token, parse role
      const decodedPayload = parseJwt(storedToken);
      setRole(decodedPayload.role || null);
    }
    setAuthLoading(false);
  }, []);

  // Fungsi bantu untuk parse JWT
  const parseJwt = (token: string) => {
    try {
      // potong token: header.payload.signature
      const base64Payload = token.split(".")[1];
      const payload = atob(base64Payload); // decode base64
      return JSON.parse(payload);
    } catch (error) {
      console.error("Error parsing JWT:", error);
      return {};
    }
  };

  const login = (newToken: string) => {
    localStorage.setItem("authToken", newToken);
    setToken(newToken);

    // Parse role
    const decodedPayload = parseJwt(newToken);
    setRole(decodedPayload.role || null);
  };

  const logout = () => {
    localStorage.removeItem("authToken");
    setToken(null);
    setRole(null);   // reset role juga
    router.push("/login");
  };

  return (
    <AuthContext.Provider
      value={{
        isAuthenticated: !!token && !authLoading,
        token,
        role,
        authLoading,
        login,
        logout,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};
