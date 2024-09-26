// src/components/Sidebar.tsx
import React from "react";
import { Link } from "react-router-dom";
import { Button } from "@/components/ui/button";

export const Sidebar: React.FC = () => {
  return (
    <aside className="h-screen w-64 bg-gray-900 text-white p-6">
      <h2 className="text-xl font-semibold mb-4">CSR Manager</h2>
      <nav className="flex flex-col space-y-4">
        <Button variant="ghost">
          <Link to="/requirements" className="text-white">
            Requirements
          </Link>
        </Button>
        <Button variant="ghost">
          <Link to="/documents" className="text-white">
            Documents
          </Link>
        </Button>
      </nav>
    </aside>
  );
};
