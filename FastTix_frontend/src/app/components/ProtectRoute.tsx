// src/app/components/ProtectRoute.tsx
"use client";

import React, { ReactNode, useContext, useEffect } from "react";
import { useRouter } from "next/navigation";
import { AuthContext } from "../../context/AuthContext";

interface ProtectedRouteProps {
  children: ReactNode;
}

const ProtectedRoute: React.FC<ProtectedRouteProps> = ({ children }) => {
  const { isAuthenticated, authLoading } = useContext(AuthContext);
  const router = useRouter();

  useEffect(() => {
    // Jika authLoading masih true, jangan apa-apakan.
    if (!authLoading) {
      // authLoading false => cek token
      if (!isAuthenticated) {
        router.push("/login");
      }
    }
  }, [authLoading, isAuthenticated, router]);

  // Jika authLoading masih true, tampilkan "loading" dulu agar tidak langsung redirect
  if (authLoading) {
    return (
      <div className="container mx-auto px-4 py-8">
        <p>Memeriksa autentikasi...</p>
      </div>
    );
  }

  // Jika sudah tidak loading, tetapi user tidak autentik => 
  // (di useEffect sudah di-redirect, jadi di sini boleh return null)
  if (!isAuthenticated) {
    return null;
  }

  return <>{children}</>;
};

export default ProtectedRoute;
