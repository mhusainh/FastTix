// src/context/NotificationContext.tsx

"use client";

import React, { createContext, useState, ReactNode, useCallback } from "react";

interface Notification {
  id: number;
  message: string;
  type: "success" | "error" | "info";
}

interface NotificationContextProps {
  addNotification: (message: string, type: "success" | "error" | "info") => void;
}

export const NotificationContext = createContext<NotificationContextProps>({
  addNotification: () => {},
});

let notificationId = 0;

export const NotificationProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [notifications, setNotifications] = useState<Notification[]>([]);

  const addNotification = useCallback((message: string, type: "success" | "error" | "info") => {
    const id = notificationId++;
    setNotifications((prev) => [...prev, { id, message, type }]);

    setTimeout(() => {
      setNotifications((prev) => prev.filter((notif) => notif.id !== id));
    }, 5000); // Hapus notifikasi setelah 5 detik
  }, []);

  return (
    <NotificationContext.Provider value={{ addNotification }}>
      {children}
      <div className="fixed top-4 right-4 space-y-2 z-50">
        {notifications.map((notif) => (
          <div
            key={notif.id}
            className={`px-4 py-2 rounded shadow-md text-white ${
              notif.type === "success"
                ? "bg-green-500"
                : notif.type === "error"
                ? "bg-red-500"
                : "bg-blue-500"
            }`}
          >
            {notif.message}
          </div>
        ))}
      </div>
    </NotificationContext.Provider>
  );
};
