import { AppBar, Box, Container, Toolbar, Typography } from "@mui/material";
import ClusterPage from "./pages/ClusterPage";

export default function App() {
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
        <ClusterPage />
      </Container>
    </Box>
  );
}
