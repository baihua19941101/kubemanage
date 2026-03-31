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
        </Tabs>
        {tab === "clusters" ? <ClusterPage /> : <NamespacePage />}
      </Container>
    </Box>
  );
}
