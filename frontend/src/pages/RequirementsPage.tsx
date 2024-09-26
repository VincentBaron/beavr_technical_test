// src/pages/RequirementsPage.tsx
import React from "react";
import { RequirementTable } from "@/components/RequirementTable";

export const RequirementsPage: React.FC = () => {
  return (
    <div className="p-6">
      <RequirementTable />
    </div>
  );
};
