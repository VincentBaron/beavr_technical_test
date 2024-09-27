import React, { useEffect, useState } from "react";
import {
  getDocuments,
  uploadFileToDocument,
  patchVersion,
  patchDocument,
  createVersion as createNewVersion,
} from "../api/api";
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
  DocumentID: number;
  Version: string;
  Path: string;
  Status: string;
  Archived: boolean;
  CreatedAt: string;
  UpdatedAt: string;
}

const statuses = ["compliant", "non-compliant"];

export const DocumentTable: React.FC = () => {
  const [documents, setDocuments] = useState<Document[]>([]);
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [currentVersionId, setCurrentVersionId] = useState<number | null>(null);
  const [showAllVersions, setShowAllVersions] = useState<{
    [key: number]: boolean;
  }>({});
  const [showDocuments, setShowDocuments] = useState<{
    [key: number]: boolean;
  }>({});

  const fetchDocuments = async () => {
    const response = await getDocuments();
    setDocuments(response.data.documents);

    // Initialize showDocuments state with the first three RequirementIDs set to true
    const initialShowDocuments: { [key: number]: boolean } = {};
    const sortedDocs = response.data.documents.sort(
      (a: Document, b: Document) => a.RequirementID - b.RequirementID
    );
    sortedDocs.slice(0, 3).forEach((doc: Document) => {
      initialShowDocuments[doc.RequirementID] = true;
    });
    setShowDocuments(initialShowDocuments);
  };

  useEffect(() => {
    fetchDocuments();
  }, []);

  const handleStatusChange = async (documentId: number, status: string) => {
    const document = documents.find((doc) => doc.ID === documentId);
    if (document) {
      const updatedDocument = { ...document, Status: status };
      await patchDocument(documentId.toString(), updatedDocument);
      fetchDocuments();
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

  const handleCreateVersion = async (docId: number) => {
    await createNewVersion(docId.toString());
    fetchDocuments();
  };

  const toggleShowDocuments = (requirementId: number) => {
    setShowDocuments((prev) => ({
      ...prev,
      [requirementId]: !prev[requirementId],
    }));
  };

  // Sort documents by RequirementID
  const sortedDocuments = [...documents].sort(
    (a, b) => a.RequirementID - b.RequirementID
  );

  return (
    <div className="space-y-4">
      {sortedDocuments.reduce((acc, doc, index) => {
        const requirementId = doc.RequirementID;
        if (
          index === 0 ||
          requirementId !== sortedDocuments[index - 1].RequirementID
        ) {
          acc.push(
            <div
              key={`requirement-${requirementId}`}
              className="border p-4 rounded"
            >
              <div className="flex justify-between items-center mb-4">
                <h2 className="text-xl font-bold">
                  Requirement ID: {requirementId}
                </h2>
                <Button onClick={() => toggleShowDocuments(requirementId)}>
                  {showDocuments[requirementId]
                    ? "Hide Documents"
                    : "Show Documents"}
                </Button>
              </div>
              {showDocuments[requirementId] && (
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  {sortedDocuments
                    .filter((d) => d.RequirementID === requirementId)
                    .map((doc) => (
                      <Card key={doc.ID} className="p-4 shadow-md">
                        <h3 className="text-lg font-semibold">{doc.Name}</h3>
                        <p>Requirement ID: {doc.RequirementID}</p>
                        <div className="mt-2 flex flex-row items-center">
                          <select
                            id={`status-${doc.ID}`}
                            value={doc.Status}
                            onChange={(e) =>
                              handleStatusChange(doc.ID, e.target.value)
                            }
                            className="mb-2 ml-2"
                            style={{
                              color:
                                doc.Status === "compliant" ? "green" : "red",
                            }}
                          >
                            {statuses.map((status) => (
                              <option key={status} value={status}>
                                {status}
                              </option>
                            ))}
                          </select>
                        </div>
                        <div className="mt-4">
                          <h3 className="text-lg font-semibold">
                            Document Versions
                          </h3>
                          <Button onClick={() => handleCreateVersion(doc.ID)}>
                            Create Version
                          </Button>
                          {doc.Versions.filter((version) => !version.Archived)
                            .slice()
                            .reverse()
                            .slice(0, showAllVersions[doc.ID] ? undefined : 1)
                            .map((version) => (
                              <div
                                key={version.ID}
                                className="mt-2 p-4 border rounded"
                              >
                                <p>V{version.Version}</p>
                                <p>Created on: {version.CreatedAt}</p>
                                <p>Last updated: {version.UpdatedAt}</p>
                                <p>Path: {version.Path}</p>
                                <div className="mt-2">
                                  <Button
                                    onClick={() => openUploadCard(version.ID)}
                                  >
                                    {version.Path ? "Edit File" : "Add File"}
                                  </Button>
                                </div>
                                {currentVersionId === version.ID && (
                                  <div className="mt-4 p-4 border rounded">
                                    <h3 className="text-lg font-semibold">
                                      Upload Document
                                    </h3>
                                    <div className="mt-2">
                                      <label
                                        htmlFor="file-upload"
                                        className="block mb-1"
                                      >
                                        Upload File:
                                      </label>
                                      <input
                                        type="file"
                                        id="file-upload"
                                        onChange={handleFileChange}
                                        className="mb-2"
                                      />
                                      <Button onClick={handleFileUpload}>
                                        Upload File
                                      </Button>
                                      <Button
                                        variant="secondary"
                                        onClick={() =>
                                          setCurrentVersionId(null)
                                        }
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
              )}
            </div>
          );
        }
        return acc;
      }, [] as JSX.Element[])}
    </div>
  );
};
