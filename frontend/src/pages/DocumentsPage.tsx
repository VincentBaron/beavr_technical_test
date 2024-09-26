// src/pages/DocumentsPage.tsx
import React from "react";
import { DocumentTable } from "@/components/DocumentTable";

export const DocumentsPage: React.FC = () => {
  return (
    <div className="p-6">
      <DocumentTable />
    </div>
  );
};
