// src/pages/DocumentsPage.tsx
import React from "react";
import { useLocation } from "react-router-dom";
import { DocumentTable } from "@/components/DocumentTable";

export const DocumentsPage: React.FC = () => {
  const location = useLocation();
  const params = new URLSearchParams(location.search);
  const referenceId = params.get("ReferenceId");

  return (
    <div className="p-6">
      <DocumentTable referenceId={referenceId} />
    </div>
  );
};
