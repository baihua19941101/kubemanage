import {
  AppBar,
  Box,
  Container,
  Tab,
  Tabs,
  Toolbar,
  Typography
} from "@mui/material";
import { useState } from "react";
import ClusterPage from "./pages/ClusterPage";
import NamespacePage from "./pages/NamespacePage";
import WorkloadPage from "./pages/WorkloadPage";

export default function App() {
  const [tab, setTab] = useState("clusters");

  return (
    <Box>
      <AppBar position="static" color="primary">
        <Toolbar>
          <Typography variant="h6" component="div">
            kubeManage MVP
          </Typography>
        </Toolbar>
      </AppBar>
      <Container sx={{ py: 3 }}>
        <Tabs value={tab} onChange={(_, value) => setTab(value)} sx={{ mb: 2 }}>
          <Tab label="集群管理" value="clusters" />
          <Tab label="名称空间管理" value="namespaces" />
          <Tab label="工作负载管理" value="workloads" />
        </Tabs>
        {tab === "clusters" && <ClusterPage />}
        {tab === "namespaces" && <NamespacePage />}
        {tab === "workloads" && <WorkloadPage />}
      </Container>
    </Box>
  );
}
