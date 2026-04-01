import { Navigate, Route, Routes } from "react-router-dom";
import ShellLayout from "./layout/ShellLayout";
import AuthAuditPage from "./pages/AuthAuditPage";
import ClusterPage from "./pages/ClusterPage";
import NamespacePage from "./pages/NamespacePage";
import ResourcePage from "./pages/ResourcePage";
import WorkloadPage from "./pages/WorkloadPage";

export default function App() {
  return (
    <Routes>
      <Route path="/" element={<ShellLayout />}>
        <Route index element={<Navigate to="/cluster" replace />} />
        <Route path="cluster" element={<ClusterPage />} />
        <Route path="namespaces" element={<NamespacePage />} />
        <Route path="workloads" element={<WorkloadPage />} />
        <Route path="resources" element={<ResourcePage />} />
        <Route path="auth-audit" element={<AuthAuditPage />} />
      </Route>
      <Route path="*" element={<Navigate to="/cluster" replace />} />
    </Routes>
  );
}
