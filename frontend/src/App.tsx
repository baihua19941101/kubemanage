import { Navigate, Route, Routes } from "react-router-dom";
import ShellLayout from "./layout/ShellLayout";
import AuthAuditPage from "./pages/AuthAuditPage";
import ClusterPage from "./pages/ClusterPage";
import NamespacePage from "./pages/NamespacePage";
import ServiceDiscoveryPage from "./pages/ServiceDiscoveryPage";
import StoragePage from "./pages/StoragePage";
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
        <Route path="service-discovery" element={<Navigate to="/service-discovery/services" replace />} />
        <Route path="service-discovery/services" element={<ServiceDiscoveryPage initialMode="services" />} />
        <Route path="service-discovery/ingresses" element={<ServiceDiscoveryPage initialMode="ingresses" />} />
        <Route path="service-discovery/hpas" element={<ServiceDiscoveryPage initialMode="hpas" />} />
        <Route path="storage" element={<Navigate to="/storage/pvs" replace />} />
        <Route path="storage/pvs" element={<StoragePage initialMode="pvs" />} />
        <Route path="storage/pvcs" element={<StoragePage initialMode="pvcs" />} />
        <Route path="storage/storageclasses" element={<StoragePage initialMode="storageclasses" />} />
        <Route path="storage/configmaps" element={<StoragePage initialMode="configmaps" />} />
        <Route path="storage/secrets" element={<StoragePage initialMode="secrets" />} />
        <Route path="resources" element={<Navigate to="/service-discovery/services" replace />} />
        <Route path="auth-audit" element={<AuthAuditPage />} />
      </Route>
      <Route path="*" element={<Navigate to="/cluster" replace />} />
    </Routes>
  );
}
