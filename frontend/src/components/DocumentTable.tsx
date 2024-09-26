// src/components/DocumentTable.tsx
import React, { useEffect, useState } from "react";
import { getDocuments, patchDocument, uploadFileToDocument } from "../api/api";
import { Card } from "../components/ui/card";
import { Button } from "../components/ui/button";

interface Document {
  ID: number;
  Name: string;
  Description: string;
  RequirementID: number;
  Path: string;
  Status: string;
}

const statuses = ["compliant", "non-compliant", "pending"];

export const DocumentTable: React.FC = () => {
  const [documents, setDocuments] = useState<Document[]>([]);
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [currentDocId, setCurrentDocId] = useState<number | null>(null);

  const fetchDocuments = async () => {
    const response = await getDocuments();
    setDocuments(response.data.documents);
  };

  useEffect(() => {
    fetchDocuments();
  }, []);

  const handleStatusChange = async (id: number, status: string) => {
    const document = documents.find((doc) => doc.ID === id);
    if (document) {
      const updatedDocument = { ...document, Status: status };
      await patchDocument(id.toString(), updatedDocument);
      fetchDocuments();
    }
  };

  const handleArchive = async (id: number) => {
    const document = documents.find((doc) => doc.ID === id);
    if (document) {
      const updatedDocument = { ...document, Archived: true };
      await patchDocument(id.toString(), updatedDocument);
      fetchDocuments();
    }
  };

  const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    if (event.target.files && event.target.files.length > 0) {
      setSelectedFile(event.target.files[0]);
    }
  };

  const handleFileUpload = async () => {
    if (!selectedFile || currentDocId === null) return;
    await uploadFileToDocument(currentDocId.toString(), selectedFile);
    setSelectedFile(null);
    setCurrentDocId(null);
    fetchDocuments();
  };

  const openUploadCard = (docId: number) => {
    setCurrentDocId(docId);
  };

  return (
    <div className="grid grid-cols-1 gap-4">
      {documents.map((doc) => (
        <Card key={doc.ID} className="p-4 shadow-md">
          <h3 className="text-lg font-semibold">{doc.Name}</h3>
          <p>Requirement ID: {doc.RequirementID}</p>
          <div className="mt-2">
            <select
              id={`status-${doc.ID}`}
              value={doc.Status}
              onChange={(e) => handleStatusChange(doc.ID, e.target.value)}
              className="mb-2"
            >
              {statuses.map((status) => (
                <option key={status} value={status}>
                  {status}
                </option>
              ))}
            </select>
          </div>
          <div className="mt-2">
            {doc.Path ? (
              <>
                <p>Document: {doc.Path}</p>
                <Button onClick={() => openUploadCard(doc.ID)}>
                  Edit Document
                </Button>
              </>
            ) : (
              <Button onClick={() => openUploadCard(doc.ID)}>
                Add Document
              </Button>
            )}
          </div>
          {currentDocId === doc.ID && (
            <div className="mt-4 p-4 border rounded">
              <h3 className="text-lg font-semibold">Upload Document</h3>
              <div className="mt-2">
                <label htmlFor="file-upload" className="block mb-1">
                  Upload File:
                </label>
                <input
                  type="file"
                  id="file-upload"
                  onChange={handleFileChange}
                  className="mb-2"
                />
                <Button onClick={handleFileUpload}>Upload File</Button>
                <Button
                  variant="secondary"
                  onClick={() => setCurrentDocId(null)}
                  className="ml-2"
                >
                  Cancel
                </Button>
              </div>
            </div>
          )}
          <Button className="mt-2" onClick={() => handleArchive(doc.ID)}>
            Delete
          </Button>
        </Card>
      ))}
    </div>
  );
};
