// src/App.tsx
import React from "react";
import {
  BrowserRouter as Router,
  Routes,
  Route,
  Navigate,
} from "react-router-dom";
import { Sidebar } from "./components/Sidebar";
import { DocumentsPage } from "./pages/DocumentsPage";
import { RequirementsPage } from "./pages/RequirementsPage";

const App: React.FC = () => {
  return (
    <Router>
      <div className="flex min-h-screen">
        {/* Sidebar on the left */}
        <Sidebar />

        {/* Main content area */}
        <main className="flex-1 p-6 bg-gray-100">
          <Routes>
            {/* Default route redirects to the requirements page */}
            <Route path="/" element={<Navigate to="/requirements" />} />

            {/* Requirements Page Route */}
            <Route path="/requirements" element={<RequirementsPage />} />

            {/* Documents Page Route */}
            <Route path="/documents" element={<DocumentsPage />} />
          </Routes>
        </main>
      </div>
    </Router>
  );
};

export default App;
