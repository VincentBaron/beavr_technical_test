// src/components/RequirementTable.tsx
import React, { useEffect, useState } from "react";
import { getRequirements } from "../api/api";
import { Card } from "@/components/ui/card";

interface Requirement {
  ID: number;
  Name: string;
  Description: string;
  Status: string;
}

export const RequirementTable: React.FC = () => {
  const [requirements, setRequirements] = useState<Requirement[]>([]);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchRequirements = async () => {
      try {
        const response = await getRequirements();
        if (Array.isArray(response.data.requirements)) {
          setRequirements(response.data.requirements);
        } else {
          setError("Unexpected response format");
        }
      } catch (err) {
        setError("Failed to fetch requirements");
      }
    };
    fetchRequirements();
  }, []);

  if (error) {
    return <div className="text-red-500">{error}</div>;
  }

  return (
    <div className="grid grid-cols-1 gap-4">
      {requirements.map((req) => (
        <Card key={req.ID} className="p-4 shadow-md">
          <h3 className="text-lg font-semibold">{req.Name}</h3>
          <p>{req.Description}</p>
          <p
            className={`mt-2 text-sm ${req.Status === "compliant" ? "text-green-500" : "text-red-500"}`}
          >
            {req.Status}
          </p>
        </Card>
      ))}
    </div>
  );
};
