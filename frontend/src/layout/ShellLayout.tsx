import {
  AppBar,
  Avatar,
  Box,
  Divider,
  Drawer,
  IconButton,
  List,
  ListItemButton,
  ListItemText,
  MenuItem,
  Select,
  Stack,
  Toolbar,
  Typography
} from "@mui/material";
import MenuIcon from "@mui/icons-material/Menu";
import { useMemo, useState } from "react";
import { NavLink, Outlet, useLocation } from "react-router-dom";

const drawerWidth = 252;

type NavItem = {
  group: string;
  label: string;
  path: string;
};

const navItems: NavItem[] = [
  { group: "Cluster", label: "集群管理", path: "/cluster" },
  { group: "Workloads", label: "工作负载", path: "/workloads" },
  { group: "Configuration", label: "名称空间", path: "/namespaces" },
  { group: "Configuration", label: "服务与配置", path: "/resources" },
  { group: "Security", label: "权限与审计", path: "/auth-audit" }
];

export default function ShellLayout() {
  const [mobileOpen, setMobileOpen] = useState(false);
  const location = useLocation();

  const grouped = useMemo(() => {
    const map = new Map<string, NavItem[]>();
    for (const item of navItems) {
      if (!map.has(item.group)) {
        map.set(item.group, []);
      }
      map.get(item.group)!.push(item);
    }
    return Array.from(map.entries());
  }, []);

  const drawer = (
    <Box sx={{ height: "100%", bgcolor: "#f6f8fb" }}>
      <Stack sx={{ px: 2, py: 2 }}>
        <Typography variant="h6" sx={{ fontWeight: 700, color: "#123d77" }}>
          kubeManage
        </Typography>
        <Typography variant="caption" color="text.secondary">
          Rancher 风格控制台
        </Typography>
      </Stack>
      <Divider />
      <List sx={{ px: 1 }}>
        {grouped.map(([group, items]) => (
          <Box key={group} sx={{ mb: 1.5 }}>
            <Typography
              variant="overline"
              sx={{ px: 1.5, color: "text.secondary", letterSpacing: 0.8 }}
            >
              {group}
            </Typography>
            {items.map((item) => (
              <ListItemButton
                key={item.path}
                component={NavLink}
                to={item.path}
                selected={location.pathname === item.path}
                sx={{ borderRadius: 1.5, mx: 0.5, my: 0.2 }}
                onClick={() => setMobileOpen(false)}
              >
                <ListItemText primary={item.label} />
              </ListItemButton>
            ))}
          </Box>
        ))}
      </List>
    </Box>
  );

  return (
    <Box sx={{ display: "flex", minHeight: "100vh", bgcolor: "#eef2f7" }}>
      <AppBar
        position="fixed"
        sx={{
          width: { sm: `calc(100% - ${drawerWidth}px)` },
          ml: { sm: `${drawerWidth}px` },
          bgcolor: "#ffffff",
          color: "#0f2f5d",
          boxShadow: "0 1px 2px rgba(15,47,93,.08)"
        }}
      >
        <Toolbar sx={{ justifyContent: "space-between" }}>
          <Stack direction="row" spacing={1.5} alignItems="center">
            <IconButton
              color="inherit"
              edge="start"
              onClick={() => setMobileOpen((v) => !v)}
              sx={{ display: { sm: "none" } }}
            >
              <MenuIcon />
            </IconButton>
            <Typography variant="subtitle1" sx={{ fontWeight: 600 }}>
              {navItems.find((n) => n.path === location.pathname)?.label || "控制台"}
            </Typography>
          </Stack>
          <Stack direction="row" spacing={2} alignItems="center">
            <Select size="small" value="default" sx={{ minWidth: 160 }}>
              <MenuItem value="default">default-cluster</MenuItem>
            </Select>
            <Stack direction="row" spacing={1} alignItems="center">
              <Avatar sx={{ width: 30, height: 30 }}>A</Avatar>
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
