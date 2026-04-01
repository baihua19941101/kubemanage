import {
  AppBar,
  Avatar,
  Box,
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
import SecurityIcon from "@mui/icons-material/Security";
import ExpandLessIcon from "@mui/icons-material/ExpandLess";
import ExpandMoreIcon from "@mui/icons-material/ExpandMore";
import ChevronRightIcon from "@mui/icons-material/ChevronRight";
import type { ReactNode } from "react";
import { useMemo, useState } from "react";
import { NavLink, Outlet, useLocation } from "react-router-dom";

const drawerWidth = 256;

type NavLeaf = {
  group: string;
  label: string;
  path: string;
  icon: ReactNode;
};

const navItems: NavLeaf[] = [
  { group: "Cluster", label: "集群管理", path: "/cluster", icon: <ClusterIcon fontSize="small" /> },
  { group: "Workloads", label: "Deployment", path: "/workloads/deployments", icon: <WorkloadIcon fontSize="small" /> },
  { group: "Workloads", label: "Pod", path: "/workloads/pods", icon: <WorkloadIcon fontSize="small" /> },
  { group: "Workloads", label: "StatefulSet", path: "/workloads/statefulsets", icon: <WorkloadIcon fontSize="small" /> },
  { group: "Workloads", label: "DaemonSet", path: "/workloads/daemonsets", icon: <WorkloadIcon fontSize="small" /> },
  { group: "Workloads", label: "Job", path: "/workloads/jobs", icon: <WorkloadIcon fontSize="small" /> },
  { group: "Workloads", label: "CronJob", path: "/workloads/cronjobs", icon: <WorkloadIcon fontSize="small" /> },
  { group: "Configuration", label: "名称空间", path: "/namespaces", icon: <ConfigIcon fontSize="small" /> },
  { group: "Configuration", label: "服务与配置", path: "/resources", icon: <ConfigIcon fontSize="small" /> },
  { group: "Security", label: "权限与审计", path: "/auth-audit", icon: <SecurityIcon fontSize="small" /> }
];

export default function ShellLayout() {
  const [mobileOpen, setMobileOpen] = useState(false);
  const [workloadsOpen, setWorkloadsOpen] = useState(true);
  const location = useLocation();

  const current = navItems.find((leaf) => leaf.path === location.pathname);

  const grouped = useMemo(() => {
    const map = new Map<string, NavLeaf[]>();
    for (const item of navItems) {
      const group = item.group;
      if (!map.has(group)) {
        map.set(group, []);
      }
      map.get(group)!.push(item);
    }
    return Array.from(map.entries());
  }, []);

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
        {grouped.map(([group, items]) => (
          <Box key={group} sx={{ mb: 1.5 }}>
            <Typography
              variant="overline"
              sx={{ px: 1.5, color: "text.secondary", letterSpacing: 0.8 }}
            >
              {group}
            </Typography>
            {group !== "Workloads" &&
              items.map((entry) => {
                const active = location.pathname === entry.path;
                return (
                  <ListItemButton
                    key={entry.path}
                    component={NavLink}
                    to={entry.path}
                    selected={active}
                    sx={{
                      borderRadius: 1.5,
                      mx: 0.5,
                      my: 0.2,
                      bgcolor: active ? "#dce9fb" : "transparent",
                      color: active ? "#0b3b75" : "inherit",
                      "&:hover": { bgcolor: active ? "#dce9fb" : "#eaf1fb" }
                    }}
                    onClick={() => setMobileOpen(false)}
                  >
                    <ListItemIcon sx={{ minWidth: 30, color: "inherit" }}>
                      {entry.icon}
                    </ListItemIcon>
                    <ListItemText primary={entry.label} />
                  </ListItemButton>
                );
              })}

            {group === "Workloads" && (
              <>
                <ListItemButton
                  selected={location.pathname.startsWith("/workloads/")}
                  onClick={() => setWorkloadsOpen((prev) => !prev)}
                  sx={{
                    borderRadius: 1.5,
                    mx: 0.5,
                    my: 0.2,
                    bgcolor: location.pathname.startsWith("/workloads/") ? "#dce9fb" : "transparent",
                    color: location.pathname.startsWith("/workloads/") ? "#0b3b75" : "inherit",
                    "&:hover": { bgcolor: location.pathname.startsWith("/workloads/") ? "#dce9fb" : "#eaf1fb" }
                  }}
                >
                  <ListItemIcon sx={{ minWidth: 30, color: "inherit" }}>
                    <WorkloadIcon fontSize="small" />
                  </ListItemIcon>
                  <ListItemText primary="工作负载" />
                  {workloadsOpen ? <ExpandLessIcon fontSize="small" /> : <ExpandMoreIcon fontSize="small" />}
                </ListItemButton>
                <Collapse in={workloadsOpen} timeout="auto" unmountOnExit>
                  <List disablePadding sx={{ pl: 1.5 }}>
                    {items.map((entry) => {
                      const active = location.pathname === entry.path;
                      return (
                        <ListItemButton
                          key={entry.path}
                          component={NavLink}
                          to={entry.path}
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
                          <ListItemText primary={entry.label} primaryTypographyProps={{ variant: "body2" }} />
                        </ListItemButton>
                      );
                    })}
                  </List>
                </Collapse>
              </>
            )}
          </Box>
        ))}
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
                {current?.label || "控制台"}
              </Typography>
            </Stack>
            <Breadcrumbs separator={<ChevronRightIcon fontSize="small" />} sx={{ color: "text.secondary" }}>
              <Typography variant="caption" color="text.secondary">
                Cluster Explorer
              </Typography>
              <Typography variant="caption" color="text.primary">
                {current?.group || "General"}
              </Typography>
              <Typography variant="caption" color="text.primary">
                {current?.label || "Overview"}
              </Typography>
            </Breadcrumbs>
          </Stack>
          <Stack direction="row" spacing={2} alignItems="center">
            <Select size="small" value="default" sx={{ minWidth: 170 }}>
              <MenuItem value="default">default-cluster</MenuItem>
            </Select>
            <Stack direction="row" spacing={1} alignItems="center">
              <Avatar sx={{ width: 30, height: 30, bgcolor: "#1d4f91" }}>A</Avatar>
              <Typography variant="body2">admin</Typography>
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
