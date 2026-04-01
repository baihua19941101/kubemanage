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
  icon?: ReactNode;
  parentLabel?: string;
};

type NavSection = {
  group: string;
  label: string;
  icon: ReactNode;
  children: NavLeaf[];
};

type NavItem =
  | { type: "link"; item: NavLeaf; icon: ReactNode }
  | { type: "section"; section: NavSection };

const navItems: NavItem[] = [
  {
    type: "link",
    icon: <ClusterIcon fontSize="small" />,
    item: { group: "Cluster", label: "集群管理", path: "/cluster" }
  },
  {
    type: "section",
    section: {
      group: "Workloads",
      label: "工作负载",
      icon: <WorkloadIcon fontSize="small" />,
      children: [
        { group: "Workloads", label: "Deployment", path: "/workloads/deployments", parentLabel: "工作负载" },
        { group: "Workloads", label: "Pod", path: "/workloads/pods", parentLabel: "工作负载" },
        { group: "Workloads", label: "StatefulSet", path: "/workloads/statefulsets", parentLabel: "工作负载" },
        { group: "Workloads", label: "DaemonSet", path: "/workloads/daemonsets", parentLabel: "工作负载" },
        { group: "Workloads", label: "Job", path: "/workloads/jobs", parentLabel: "工作负载" },
        { group: "Workloads", label: "CronJob", path: "/workloads/cronjobs", parentLabel: "工作负载" }
      ]
    }
  },
  {
    type: "link",
    icon: <ConfigIcon fontSize="small" />,
    item: { group: "Configuration", label: "名称空间", path: "/namespaces" }
  },
  {
    type: "link",
    icon: <ConfigIcon fontSize="small" />,
    item: { group: "Configuration", label: "服务与配置", path: "/resources" }
  },
  {
    type: "link",
    icon: <SecurityIcon fontSize="small" />,
    item: { group: "Security", label: "权限与审计", path: "/auth-audit" }
  }
];

export default function ShellLayout() {
  const [mobileOpen, setMobileOpen] = useState(false);
  const [workloadsOpen, setWorkloadsOpen] = useState(true);
  const location = useLocation();

  const allLeaves = useMemo(() => {
    const leaves: NavLeaf[] = [];
    for (const item of navItems) {
      if (item.type === "link") {
        leaves.push(item.item);
      } else {
        leaves.push(...item.section.children);
      }
    }
    return leaves;
  }, []);

  const current = allLeaves.find((leaf) => leaf.path === location.pathname);

  const grouped = useMemo(() => {
    const map = new Map<string, NavItem[]>();
    for (const item of navItems) {
      const group = item.type === "link" ? item.item.group : item.section.group;
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
            {items.map((entry, index) => {
              if (entry.type === "link") {
                const active = location.pathname === entry.item.path;
                return (
                  <ListItemButton
                    key={entry.item.path}
                    component={NavLink}
                    to={entry.item.path}
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
                    <ListItemText primary={entry.item.label} />
                  </ListItemButton>
                );
              }

              const sectionActive = entry.section.children.some((child) => child.path === location.pathname);
              return (
                <Box key={`${entry.section.group}-${index}`}>
                  <ListItemButton
                    selected={sectionActive}
                    onClick={() => setWorkloadsOpen((prev) => !prev)}
                    sx={{
                      borderRadius: 1.5,
                      mx: 0.5,
                      my: 0.2,
                      bgcolor: sectionActive ? "#dce9fb" : "transparent",
                      color: sectionActive ? "#0b3b75" : "inherit",
                      "&:hover": { bgcolor: sectionActive ? "#dce9fb" : "#eaf1fb" }
                    }}
                  >
                    <ListItemIcon sx={{ minWidth: 30, color: "inherit" }}>
                      {entry.section.icon}
                    </ListItemIcon>
                    <ListItemText primary={entry.section.label} />
                    {workloadsOpen ? <ExpandLessIcon fontSize="small" /> : <ExpandMoreIcon fontSize="small" />}
                  </ListItemButton>
                  <Collapse in={workloadsOpen} timeout="auto" unmountOnExit>
                    <List disablePadding sx={{ pl: 1.5 }}>
                      {entry.section.children.map((child) => {
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
              {current?.parentLabel && (
                <Typography variant="caption" color="text.primary">
                  {current.parentLabel}
                </Typography>
              )}
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
