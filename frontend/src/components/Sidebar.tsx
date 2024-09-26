// src/components/Sidebar.tsx
import React from "react";
import { Link, useLocation } from "react-router-dom";
import { Button } from "@/components/ui/button";

export const Sidebar: React.FC = () => {
  const location = useLocation();

  return (
    <aside className="fixed h-screen w-64 bg-green-900 text-white p-6">
      <h2 className="text-xl font-semibold mb-4 flex items-center">
        {" "}
        🦫 | Kit
      </h2>
      <nav className="flex flex-col space-y-4">
        <Button
          variant="ghost"
          className={
            location.pathname === "/requirements" ? "bg-green-700" : ""
          }
        >
          <Link to="/requirements" className="text-white">
            Requirements
          </Link>
        </Button>
        <Button
          variant="ghost"
          className={location.pathname === "/documents" ? "bg-green-700" : ""}
        >
          <Link to="/documents" className="text-white">
            Documents
          </Link>
        </Button>
      </nav>
    </aside>
  );
};
