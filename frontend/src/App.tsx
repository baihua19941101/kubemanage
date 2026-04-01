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
        <Route path="workloads" element={<Navigate to="/workloads/deployments" replace />} />
        <Route path="workloads/deployments" element={<WorkloadPage initialMode="deployments" showModeSwitcher={false} />} />
        <Route path="workloads/pods" element={<WorkloadPage initialMode="pods" showModeSwitcher={false} />} />
        <Route path="workloads/statefulsets" element={<WorkloadPage initialMode="statefulsets" showModeSwitcher={false} />} />
        <Route path="workloads/daemonsets" element={<WorkloadPage initialMode="daemonsets" showModeSwitcher={false} />} />
        <Route path="workloads/jobs" element={<WorkloadPage initialMode="jobs" showModeSwitcher={false} />} />
        <Route path="workloads/cronjobs" element={<WorkloadPage initialMode="cronjobs" showModeSwitcher={false} />} />
        <Route path="resources" element={<ResourcePage />} />
        <Route path="auth-audit" element={<AuthAuditPage />} />
      </Route>
      <Route path="*" element={<Navigate to="/cluster" replace />} />
    </Routes>
  );
}
