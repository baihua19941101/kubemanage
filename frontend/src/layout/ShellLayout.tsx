import {
  AppBar,
  Avatar,
  Box,
  Button,
  Breadcrumbs,
  Collapse,
  Divider,
  Drawer,
  IconButton,
  List,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  MenuItem,
  Select,
  Stack,
  Toolbar,
  Typography
} from "@mui/material";
import MenuIcon from "@mui/icons-material/Menu";
import ClusterIcon from "@mui/icons-material/Hub";
import WorkloadIcon from "@mui/icons-material/ViewQuilt";
import ConfigIcon from "@mui/icons-material/Tune";
import StorageIcon from "@mui/icons-material/Storage";
import SecurityIcon from "@mui/icons-material/Security";
import ExpandLessIcon from "@mui/icons-material/ExpandLess";
import ExpandMoreIcon from "@mui/icons-material/ExpandMore";
import ChevronRightIcon from "@mui/icons-material/ChevronRight";
import type { ReactNode } from "react";
import { useMemo, useState } from "react";
import { NavLink, Outlet, useLocation } from "react-router-dom";
import { useAuthStore } from "../stores/useAuthStore";

const drawerWidth = 256;

type NavChild = {
  label: string;
  path: string;
};

type NavSection = {
  key: string;
  label: string;
  icon: ReactNode;
  children: NavChild[];
};

const navSections: NavSection[] = [
  {
    key: "cluster",
    label: "Cluster",
    icon: <ClusterIcon fontSize="small" />,
    children: [
      { label: "集群管理", path: "/cluster" },
      { label: "名称空间", path: "/namespaces" },
      { label: "节点管理", path: "/nodes" }
    ]
  },
  {
    key: "workloads",
    label: "Workloads",
    icon: <WorkloadIcon fontSize="small" />,
    children: [
      { label: "Deployment", path: "/workloads/deployments" },
      { label: "Pod", path: "/workloads/pods" },
      { label: "StatefulSet", path: "/workloads/statefulsets" },
      { label: "DaemonSet", path: "/workloads/daemonsets" },
      { label: "Job", path: "/workloads/jobs" },
      { label: "CronJob", path: "/workloads/cronjobs" }
    ]
  },
  {
    key: "service-discovery",
    label: "Service Discovery",
    icon: <ConfigIcon fontSize="small" />,
    children: [
      { label: "Service", path: "/service-discovery/services" },
      { label: "Ingress", path: "/service-discovery/ingresses" },
      { label: "HPA", path: "/service-discovery/hpas" }
    ]
  },
  {
    key: "storage",
    label: "Storage",
    icon: <StorageIcon fontSize="small" />,
    children: [
      { label: "PersistentVolumes", path: "/storage/pvs" },
      { label: "PersistentVolumeClaims", path: "/storage/pvcs" },
      { label: "StorageClasses", path: "/storage/storageclasses" },
      { label: "ConfigMaps", path: "/storage/configmaps" },
      { label: "Secrets", path: "/storage/secrets" }
    ]
  },
  {
    key: "security",
    label: "Security",
    icon: <SecurityIcon fontSize="small" />,
    children: [{ label: "权限与审计", path: "/auth-audit" }]
  }
];

function makeInitialOpenState(pathname: string): Record<string, boolean> {
  const state: Record<string, boolean> = {};
  for (const section of navSections) {
    state[section.key] = section.children.some((child) => child.path === pathname);
  }
  if (!Object.values(state).some(Boolean)) {
    state.cluster = true;
  }
  return state;
}

export default function ShellLayout() {
  const [mobileOpen, setMobileOpen] = useState(false);
  const location = useLocation();
  const [openMap, setOpenMap] = useState<Record<string, boolean>>(() => makeInitialOpenState(location.pathname));
  const user = useAuthStore((s) => s.user);
  const role = useAuthStore((s) => s.role);
  const logout = useAuthStore((s) => s.logout);

  const current = useMemo(() => {
    for (const section of navSections) {
      const child = section.children.find((item) => item.path === location.pathname);
      if (child) {
        return { section, child };
      }
    }
    return null;
  }, [location.pathname]);

  const drawer = (
    <Box sx={{ height: "100%", bgcolor: "#f4f7fb" }}>
      <Stack sx={{ px: 2, py: 2.5 }}>
        <Typography variant="h6" sx={{ fontWeight: 800, color: "#123d77", lineHeight: 1.1 }}>
          kubeManage
        </Typography>
        <Typography variant="caption" color="text.secondary">
          Rancher-like Console
        </Typography>
      </Stack>
      <Divider />
      <List sx={{ px: 1, py: 1 }}>
        {navSections.map((section) => {
          const sectionActive = section.children.some((child) => child.path === location.pathname);
          const sectionOpen = openMap[section.key] ?? false;

          return (
            <Box key={section.key} sx={{ mb: 1 }}>
              <ListItemButton
                onClick={() => setOpenMap((prev) => ({ ...prev, [section.key]: !sectionOpen }))}
                selected={sectionActive}
                sx={{
                  borderRadius: 1.5,
                  mx: 0.5,
                  my: 0.2,
                  bgcolor: sectionActive ? "#dce9fb" : "transparent",
                  color: sectionActive ? "#0b3b75" : "inherit",
                  "&:hover": { bgcolor: sectionActive ? "#dce9fb" : "#eaf1fb" }
                }}
              >
                <ListItemIcon sx={{ minWidth: 30, color: "inherit" }}>{section.icon}</ListItemIcon>
                <ListItemText primary={section.label} />
                {sectionOpen ? <ExpandLessIcon fontSize="small" /> : <ExpandMoreIcon fontSize="small" />}
              </ListItemButton>

              <Collapse in={sectionOpen} timeout="auto" unmountOnExit>
                <List disablePadding sx={{ pl: 1.5 }}>
                  {section.children.map((child) => {
                    const active = location.pathname === child.path;
                    return (
                      <ListItemButton
                        key={child.path}
                        component={NavLink}
                        to={child.path}
                        selected={active}
                        sx={{
                          borderRadius: 1.5,
                          mx: 0.5,
                          my: 0.1,
                          pl: 4,
                          bgcolor: active ? "#dce9fb" : "transparent",
                          color: active ? "#0b3b75" : "inherit",
                          "&:hover": { bgcolor: active ? "#dce9fb" : "#eaf1fb" }
                        }}
                        onClick={() => setMobileOpen(false)}
                      >
                        <ListItemText primary={child.label} primaryTypographyProps={{ variant: "body2" }} />
                      </ListItemButton>
                    );
                  })}
                </List>
              </Collapse>
            </Box>
          );
        })}
      </List>
    </Box>
  );

  return (
    <Box sx={{ display: "flex", minHeight: "100vh", bgcolor: "#edf2f8" }}>
      <AppBar
        position="fixed"
        sx={{
          width: { sm: `calc(100% - ${drawerWidth}px)` },
          ml: { sm: `${drawerWidth}px` },
          bgcolor: "#ffffff",
          color: "#0f2f5d",
          borderBottom: "1px solid #dce4ef",
          boxShadow: "none"
        }}
      >
        <Toolbar sx={{ justifyContent: "space-between", minHeight: 64 }}>
          <Stack spacing={0.5}>
            <Stack direction="row" spacing={1.5} alignItems="center">
              <IconButton
                color="inherit"
                edge="start"
                onClick={() => setMobileOpen((v) => !v)}
                sx={{ display: { sm: "none" } }}
              >
                <MenuIcon />
              </IconButton>
              <Typography variant="subtitle1" sx={{ fontWeight: 700 }}>
                {current?.child.label || "控制台"}
              </Typography>
            </Stack>
            <Breadcrumbs separator={<ChevronRightIcon fontSize="small" />} sx={{ color: "text.secondary" }}>
              <Typography variant="caption" color="text.secondary">
                Cluster Explorer
              </Typography>
              <Typography variant="caption" color="text.primary">
                {current?.section.label || "General"}
              </Typography>
              <Typography variant="caption" color="text.primary">
                {current?.child.label || "Overview"}
              </Typography>
            </Breadcrumbs>
          </Stack>
          <Stack direction="row" spacing={2} alignItems="center">
            <Select size="small" value="default" sx={{ minWidth: 170 }}>
              <MenuItem value="default">default-cluster</MenuItem>
            </Select>
            <Stack direction="row" spacing={1} alignItems="center">
              <Avatar sx={{ width: 30, height: 30, bgcolor: "#1d4f91" }}>{(user || "U").slice(0, 1).toUpperCase()}</Avatar>
              <Stack spacing={0} sx={{ minWidth: 100 }}>
                <Typography variant="body2">{user || "unknown"}</Typography>
                <Typography variant="caption" color="text.secondary">
                  {role}
                </Typography>
              </Stack>
              <Button
                size="small"
                onClick={() => {
                  void logout();
                }}
              >
                退出
              </Button>
            </Stack>
          </Stack>
        </Toolbar>
      </AppBar>

      <Box component="nav" sx={{ width: { sm: drawerWidth }, flexShrink: { sm: 0 } }}>
        <Drawer
          variant="temporary"
          open={mobileOpen}
          onClose={() => setMobileOpen(false)}
          ModalProps={{ keepMounted: true }}
          sx={{
            display: { xs: "block", sm: "none" },
            "& .MuiDrawer-paper": { boxSizing: "border-box", width: drawerWidth }
          }}
        >
          {drawer}
        </Drawer>
        <Drawer
          variant="permanent"
          sx={{
            display: { xs: "none", sm: "block" },
            "& .MuiDrawer-paper": {
              boxSizing: "border-box",
              width: drawerWidth,
              borderRight: "1px solid #dbe3ef"
            }
          }}
          open
        >
          {drawer}
        </Drawer>
      </Box>

      <Box
        component="main"
        sx={{
          flexGrow: 1,
          width: { sm: `calc(100% - ${drawerWidth}px)` },
          px: { xs: 2, sm: 3 },
          py: 3
        }}
      >
        <Toolbar />
        <Outlet />
      </Box>
    </Box>
  );
}
