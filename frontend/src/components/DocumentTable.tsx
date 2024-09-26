// src/components/DocumentTable.tsx
import React, { useEffect, useState } from "react";
import { getDocuments, uploadFileToDocument, patchVersion } from "../api/api";
import { Card } from "../components/ui/card";
import { Button } from "../components/ui/button";

interface Document {
  ID: number;
  Name: string;
  Description: string;
  RequirementID: number;
  Path: string;
  Status: string;
  Versions: DocumentVersion[];
}

interface DocumentVersion {
  ID: number;
  CreatedAt: string;
  UpdatedAt: string;
  DocumentID: number;
  Version: string;
  Path: string;
  Status: string;
}

const statuses = ["compliant", "non-compliant", "pending"];

export const DocumentTable: React.FC = () => {
  const [documents, setDocuments] = useState<Document[]>([]);
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [currentVersionId, setCurrentVersionId] = useState<number | null>(null);
  const [showAllVersions, setShowAllVersions] = useState<{
    [key: number]: boolean;
  }>({});

  const fetchDocuments = async () => {
    const response = await getDocuments();
    setDocuments(response.data.documents);
  };

  useEffect(() => {
    fetchDocuments();
  }, []);

  const handleStatusChange = async (versionId: number, status: string) => {
    const document = documents.find((doc) =>
      doc.Versions.some((version) => version.ID === versionId)
    );
    if (document) {
      const version = document.Versions.find(
        (version) => version.ID === versionId
      );
      if (version) {
        const updatedVersion = { ...version, Status: status };
        await patchVersion(versionId.toString(), updatedVersion);
        fetchDocuments();
      }
    }
  };

  const handleArchive = async (versionId: number) => {
    const document = documents.find((doc) =>
      doc.Versions.some((version) => version.ID === versionId)
    );
    if (document) {
      const version = document.Versions.find(
        (version) => version.ID === versionId
      );
      if (version) {
        const updatedVersion = { ...version, Archived: true };
        await patchVersion(versionId.toString(), updatedVersion);
        fetchDocuments();
      }
    }
  };

  const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    if (event.target.files && event.target.files.length > 0) {
      setSelectedFile(event.target.files[0]);
    }
  };

  const handleFileUpload = async () => {
    if (!selectedFile || currentVersionId === null) return;
    await uploadFileToDocument(currentVersionId.toString(), selectedFile);
    setSelectedFile(null);
    setCurrentVersionId(null);
    fetchDocuments();
  };

  const openUploadCard = (versionId: number) => {
    setCurrentVersionId(versionId);
  };

  const toggleShowAllVersions = (docId: number) => {
    setShowAllVersions((prev) => ({
      ...prev,
      [docId]: !prev[docId],
    }));
  };

  return (
    <div className="grid grid-cols-1 gap-4">
      {documents.map((doc) => (
        <Card key={doc.ID} className="p-4 shadow-md">
          <h3 className="text-lg font-semibold">{doc.Name}</h3>
          <p>Requirement ID: {doc.RequirementID}</p>
          <div className="mt-4">
            <h3 className="text-lg font-semibold">Document Versions</h3>
            {doc.Versions.slice()
              .reverse()
              .slice(0, showAllVersions[doc.ID] ? undefined : 1)
              .map((version) => (
                <div key={version.ID} className="mt-2 p-4 border rounded">
                  <p>V{version.Version}</p>
                  <p>Created on: {version.CreatedAt}</p>
                  <p>Last updated: {version.UpdatedAt}</p>
                  <div className="mt-2 flex flex-row">
                    <p>Status:</p>
                    <select
                      id={`status-${version.ID}`}
                      value={version.Status}
                      onChange={(e) =>
                        handleStatusChange(version.ID, e.target.value)
                      }
                      className="mb-2"
                    >
                      {statuses.map((status) => (
                        <option key={status} value={status}>
                          {status}
                        </option>
                      ))}
                    </select>
                  </div>
                  <p>Path: {version.Path}</p>
                  <div className="mt-2">
                    <Button onClick={() => openUploadCard(version.ID)}>
                      {version.Path ? "Edit File" : "Add File"}
                    </Button>
                  </div>
                  {currentVersionId === version.ID && (
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
                          onClick={() => setCurrentVersionId(null)}
                          className="ml-2"
                        >
                          Cancel
                        </Button>
                      </div>
                    </div>
                  )}
                  <Button
                    className="mt-2"
                    onClick={() => handleArchive(version.ID)}
                  >
                    Delete
                  </Button>
                </div>
              ))}
            <Button
              className="mt-2"
              onClick={() => toggleShowAllVersions(doc.ID)}
            >
              {showAllVersions[doc.ID] ? "Show Less" : "Show All"}
            </Button>
          </div>
        </Card>
      ))}
    </div>
  );
};
