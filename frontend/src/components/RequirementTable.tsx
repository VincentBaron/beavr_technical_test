import React, { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { getRequirements } from "../api/api";
import { Card } from "@/components/ui/card";

interface Requirement {
  ID: number;
  Name: string;
  Description: string;
  Status: string;
  Documents: Document[];
}

interface Document {
  ID: number;
  Name: string;
  Status: string;
  RequirementID: number;
}

export const RequirementTable: React.FC = () => {
  const [requirements, setRequirements] = useState<Requirement[]>([]);
  const [error, setError] = useState<string | null>(null);
  const navigate = useNavigate();

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

  const handleRequirementClick = (requirementId: number) => {
    navigate(`/documents?ReferenceId=${requirementId}`);
  };

  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
      {requirements.map((req) => {
        const compliantDocuments = req.Documents.filter(
          (doc) => doc.Status === "compliant"
        ).length;
        const totalDocuments = req.Documents.length;

        return (
          <Card
            key={req.ID}
            className="p-4 shadow-md h-48 flex flex-col justify-between overflow-hidden relative"
          >
            <div className="overflow-hidden">
              <h3
                className="text-lg font-semibold truncate cursor-pointer"
                onClick={() => handleRequirementClick(req.ID)}
              >
                {req.Name}
              </h3>
              <p className="text-sm overflow-hidden text-ellipsis">
                {req.Description}
              </p>
            </div>
            <p
              className={`mt-2 text-sm ${
                req.Status === "compliant" ? "text-green-500" : "text-red-500"
              }`}
            >
              {req.Status}
            </p>
            <div
              className={`absolute bottom-2 right-2 text-white text-xs font-bold py-1 px-2 rounded-full ${
                req.Status === "compliant" ? "bg-green-500" : "bg-red-500"
              }`}
            >
              {compliantDocuments}/{totalDocuments}
            </div>
          </Card>
        );
      })}
    </div>
  );
};
