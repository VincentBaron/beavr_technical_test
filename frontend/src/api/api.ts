// src/api/api.ts
import axios from "axios";

const api = axios.create({
  baseURL: "http://localhost:8080", // Update with your backend URL
});

export const getRequirements = () => api.get("/requirements");
export const getDocuments = async (params?: { ReferenceId?: number }) => {
  if (params) {
    return await api.get("/documents", { params });
  } else {
    return await api.get("/documents");
  }
};
export const patchDocument = (id: string, data: any) => {
  return api.patch(`/documents/${id}`, data);
};
export const patchVersion = (id: string, data: any) => {
  return api.patch(`/documents/versions/${id}`, data);
};
export const uploadFileToDocument = (id: string, file: File) => {
  const formData = new FormData();
  formData.append("file", file);

  return api.patch(`/documents/versions/${id}/upload-file`, formData, {
    headers: {
      "Content-Type": "multipart/form-data",
    },
  });
};
export const createVersion = async (docId: string) => {
  return await api.post(`/documents/${docId}/versions`);
};
