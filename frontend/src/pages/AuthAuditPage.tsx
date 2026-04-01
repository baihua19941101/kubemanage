import {
  Alert,
  Box,
  Button,
  Chip,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  FormControl,
  InputLabel,
  MenuItem,
  Select,
  Stack,
  TextField,
  Typography
} from "@mui/material";
import { useEffect, useMemo, useState } from "react";
import PageScaffold from "../components/framework/PageScaffold";
import ResourceTable from "../components/framework/ResourceTable";
import { apiFetch, parseApiError } from "../lib/api";
import { useAuthStore } from "../stores/useAuthStore";

type AuditItem = {
  time: string;
  user: string;
  role: string;
  method: string;
  path: string;
  statusCode: number;
};

type UserItem = {
  id: number;
  username: string;
  role: string;
  allowedNamespaces: string[];
  isActive: boolean;
  createdAt: string;
  updatedAt: string;
};

export default function AuthAuditPage() {
  const role = useAuthStore((s) => s.role);
  const setRole = useAuthStore((s) => s.setRole);
  const canAuditRead = useAuthStore((s) => s.canAuditRead);
  const canUserManage = useAuthStore((s) => s.canUserManage);
  const allowedNamespaces = useAuthStore((s) => s.allowedNamespaces);

  const [audits, setAudits] = useState<AuditItem[]>([]);
  const [users, setUsers] = useState<UserItem[]>([]);
  const [error, setError] = useState("");
  const [userError, setUserError] = useState("");
  const [loading, setLoading] = useState(false);
  const [usersLoading, setUsersLoading] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [passwordSubmitting, setPasswordSubmitting] = useState(false);
  const [userFilter, setUserFilter] = useState("");
  const [roleFilter, setRoleFilter] = useState("");
  const [pathFilter, setPathFilter] = useState("");
  const [methodFilter, setMethodFilter] = useState("");
  const [statusCodeFilter, setStatusCodeFilter] = useState("");
  const [limitFilter, setLimitFilter] = useState("50");

  const [newUsername, setNewUsername] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [newRole, setNewRole] = useState("readonly");
  const [newNamespaces, setNewNamespaces] = useState("dev");

  const [resetTarget, setResetTarget] = useState<UserItem | null>(null);
  const [resetPasswordValue, setResetPasswordValue] = useState("");

  useEffect(() => {
    void loadAudits();
    void loadUsers();
  }, [role]);

  async function loadAudits() {
    setError("");
    setAudits([]);
    if (!canAuditRead()) {
      return;
    }
    setLoading(true);
    try {
      const params = new URLSearchParams();
      if (userFilter.trim()) params.set("user", userFilter.trim());
      if (roleFilter.trim()) params.set("role", roleFilter.trim());
      if (pathFilter.trim()) params.set("path", pathFilter.trim());
      if (methodFilter.trim()) params.set("method", methodFilter.trim().toUpperCase());
      if (statusCodeFilter.trim()) params.set("statusCode", statusCodeFilter.trim());
      if (limitFilter.trim()) params.set("limit", limitFilter.trim());
      const query = params.toString();
      const resp = await apiFetch(`/api/v1/audits${query ? `?${query}` : ""}`);
      if (!resp.ok) {
        const err = await parseApiError(resp, "当前角色无权限查看审计日志");
        setError(err.message);
        return;
      }
      const data = (await resp.json()) as { items: AuditItem[] };
      setAudits(data.items);
    } finally {
      setLoading(false);
    }
  }

  async function loadUsers() {
    setUserError("");
    setUsers([]);
    if (!canUserManage()) {
      return;
    }
    setUsersLoading(true);
    try {
      const resp = await apiFetch("/api/v1/auth/users");
      if (!resp.ok) {
        const err = await parseApiError(resp, "加载用户列表失败");
        setUserError(err.message);
        return;
      }
      const data = (await resp.json()) as { items: UserItem[] };
      setUsers(data.items);
    } finally {
      setUsersLoading(false);
    }
  }

  function resetFilters() {
    setUserFilter("");
    setRoleFilter("");
    setPathFilter("");
    setMethodFilter("");
    setStatusCodeFilter("");
    setLimitFilter("50");
  }

  function parseNamespaces(raw: string) {
    return raw
      .split(",")
      .map((item) => item.trim())
      .filter(Boolean);
  }

  async function createUser() {
    if (!newUsername.trim() || !newPassword.trim()) {
      setUserError("用户名和密码不能为空");
      return;
    }
    setSubmitting(true);
    setUserError("");
    try {
      const resp = await apiFetch("/api/v1/auth/users", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          "X-Action-Confirm": "CONFIRM"
        },
        body: JSON.stringify({
          username: newUsername.trim(),
          password: newPassword,
          role: newRole,
          allowedNamespaces: parseNamespaces(newNamespaces)
        })
      });
      if (!resp.ok) {
        const err = await parseApiError(resp, "创建用户失败");
        setUserError(err.message);
        return;
      }
      setNewUsername("");
      setNewPassword("");
      setNewRole("readonly");
      setNewNamespaces("dev");
      await loadUsers();
    } finally {
      setSubmitting(false);
    }
  }

  async function updateUserStatus(item: UserItem, active: boolean) {
    setUserError("");
    try {
      const resp = await apiFetch(`/api/v1/auth/users/${encodeURIComponent(item.username)}/status`, {
        method: "PATCH",
        headers: {
          "Content-Type": "application/json",
          "X-Action-Confirm": "CONFIRM"
        },
        body: JSON.stringify({ isActive: active })
      });
      if (!resp.ok) {
        const err = await parseApiError(resp, `更新用户状态失败：${item.username}`);
        setUserError(err.message);
        return;
      }
      await loadUsers();
    } catch (err) {
      setUserError(err instanceof Error ? err.message : "更新用户状态失败");
    }
  }

  async function submitResetPassword() {
    if (!resetTarget) return;
    if (resetPasswordValue.length < 6) {
      setUserError("新密码长度不能小于 6");
      return;
    }
    setPasswordSubmitting(true);
    setUserError("");
    try {
      const resp = await apiFetch(`/api/v1/auth/users/${encodeURIComponent(resetTarget.username)}/reset-password`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          "X-Action-Confirm": "CONFIRM"
        },
        body: JSON.stringify({ password: resetPasswordValue })
      });
      if (!resp.ok) {
        const err = await parseApiError(resp, "重置密码失败");
        setUserError(err.message);
        return;
      }
      setResetTarget(null);
      setResetPasswordValue("");
    } finally {
      setPasswordSubmitting(false);
    }
  }

  const auditColumns = useMemo(
    () => [
      { key: "time", header: "时间", render: (r: AuditItem) => r.time },
      { key: "user", header: "用户", render: (r: AuditItem) => r.user },
      { key: "role", header: "角色", render: (r: AuditItem) => r.role },
      { key: "method", header: "方法", render: (r: AuditItem) => r.method },
      { key: "path", header: "路径", render: (r: AuditItem) => r.path },
      { key: "code", header: "状态码", render: (r: AuditItem) => r.statusCode }
    ],
    []
  );

  const userColumns = useMemo(
    () => [
      { key: "username", header: "用户名", render: (r: UserItem) => r.username },
      { key: "role", header: "角色", render: (r: UserItem) => r.role },
      {
        key: "namespaces",
        header: "授权命名空间",
        render: (r: UserItem) => (r.allowedNamespaces.length === 0 ? "-" : r.allowedNamespaces.join(", "))
      },
      {
        key: "status",
        header: "状态",
        render: (r: UserItem) => <Chip size="small" color={r.isActive ? "success" : "default"} label={r.isActive ? "active" : "disabled"} />
      },
      {
        key: "actions",
        header: "操作",
        render: (r: UserItem) => (
          <Stack direction="row" spacing={1}>
            <Button size="small" variant="outlined" onClick={() => void updateUserStatus(r, !r.isActive)}>
              {r.isActive ? "禁用" : "启用"}
            </Button>
            <Button size="small" variant="outlined" onClick={() => {
              setResetTarget(r);
              setResetPasswordValue("");
            }}>
              重置密码
            </Button>
          </Stack>
        )
      }
    ],
    []
  );

  return (
    <>
      <PageScaffold
        title="权限与审计"
        description="角色切换、授权范围说明、用户管理与审计筛选"
        actions={
          <FormControl sx={{ width: 220 }}>
            <InputLabel id="role-label">当前角色</InputLabel>
            <Select
              labelId="role-label"
              value={role}
              label="当前角色"
              onChange={(e) => setRole(e.target.value)}
            >
              <MenuItem value="readonly">readonly</MenuItem>
              <MenuItem value="standard-user">standard-user</MenuItem>
              <MenuItem value="admin">admin</MenuItem>
            </Select>
          </FormControl>
        }
      >
        <Alert severity="info" sx={{ m: 1.5 }}>
          当前角色 `{role}` 的可写命名空间：
          {allowedNamespaces().length === 0
            ? " 无"
            : allowedNamespaces()[0] === "*"
              ? " 全部"
              : ` ${allowedNamespaces().join(", ")}`}
        </Alert>

        <Box sx={{ px: 1.5, py: 1 }}>
          <Typography variant="subtitle1" sx={{ fontWeight: 700 }}>用户管理</Typography>
          <Typography variant="body2" color="text.secondary">仅 admin 可执行创建用户、启停与密码重置。</Typography>
        </Box>
        {!canUserManage() && (
          <Alert severity="info" sx={{ m: 1.5 }}>
            当前角色无用户管理权限。
          </Alert>
        )}
        {userError && (
          <Alert severity="error" sx={{ m: 1.5 }}>
            {userError}
          </Alert>
        )}
        {canUserManage() && (
          <Stack direction={{ xs: "column", md: "row" }} spacing={1.5} sx={{ px: 1.5, pb: 1.5 }} useFlexGap flexWrap="wrap">
            <TextField size="small" label="用户名" value={newUsername} onChange={(e) => setNewUsername(e.target.value)} sx={{ width: 180 }} />
            <TextField size="small" label="初始密码" value={newPassword} onChange={(e) => setNewPassword(e.target.value)} sx={{ width: 180 }} />
            <FormControl size="small" sx={{ width: 160 }}>
              <InputLabel id="new-role-label">角色</InputLabel>
              <Select labelId="new-role-label" value={newRole} label="角色" onChange={(e) => setNewRole(e.target.value)}>
                <MenuItem value="readonly">readonly</MenuItem>
                <MenuItem value="standard-user">standard-user</MenuItem>
                <MenuItem value="admin">admin</MenuItem>
              </Select>
            </FormControl>
            <TextField
              size="small"
              label="授权命名空间(逗号分隔)"
              value={newNamespaces}
              onChange={(e) => setNewNamespaces(e.target.value)}
              sx={{ width: 260 }}
            />
            <Button variant="contained" disabled={submitting} onClick={() => void createUser()}>
              创建用户
            </Button>
            <Button variant="outlined" onClick={() => void loadUsers()}>
              刷新用户
            </Button>
          </Stack>
        )}
        <ResourceTable
          loading={usersLoading}
          rows={users}
          rowKey={(r) => String(r.id)}
          columns={userColumns}
        />

        <Box sx={{ mt: 2, px: 1.5, py: 1, borderTop: "1px solid #d7e1ef" }}>
          <Typography variant="subtitle1" sx={{ fontWeight: 700 }}>审计日志</Typography>
          <Typography variant="body2" color="text.secondary">支持按用户、角色、方法、路径与状态码筛选。</Typography>
        </Box>
        {!canAuditRead() && (
          <Alert severity="info" sx={{ m: 1.5 }}>
            当前角色无审计日志查看权限。
          </Alert>
        )}
        {error && <Alert severity="error" sx={{ m: 1.5 }}>{error}</Alert>}
        {canAuditRead() && (
          <Stack direction="row" spacing={1.5} sx={{ px: 1.5, pb: 1.5 }} useFlexGap flexWrap="wrap">
            <TextField size="small" label="用户" value={userFilter} onChange={(e) => setUserFilter(e.target.value)} />
            <TextField size="small" label="角色" value={roleFilter} onChange={(e) => setRoleFilter(e.target.value)} sx={{ width: 130 }} />
            <TextField size="small" label="方法" value={methodFilter} onChange={(e) => setMethodFilter(e.target.value)} sx={{ width: 120 }} />
            <TextField size="small" label="路径包含" value={pathFilter} onChange={(e) => setPathFilter(e.target.value)} sx={{ width: 260 }} />
            <TextField size="small" label="状态码" value={statusCodeFilter} onChange={(e) => setStatusCodeFilter(e.target.value)} sx={{ width: 120 }} />
            <TextField size="small" label="数量限制" value={limitFilter} onChange={(e) => setLimitFilter(e.target.value)} sx={{ width: 120 }} />
            <Button variant="contained" onClick={() => void loadAudits()}>筛选</Button>
            <Button onClick={resetFilters}>清空</Button>
          </Stack>
        )}
        <ResourceTable
          loading={loading}
          rows={audits}
          rowKey={(r) => `${r.time}-${r.path}-${r.statusCode}`}
          columns={auditColumns}
        />
        {!loading && canAuditRead() && audits.length === 0 && (
          <Stack sx={{ p: 2 }}>
            <Typography color="text.secondary">暂无审计记录</Typography>
          </Stack>
        )}
      </PageScaffold>

      <Dialog open={Boolean(resetTarget)} onClose={() => setResetTarget(null)} fullWidth maxWidth="xs">
        <DialogTitle>重置密码</DialogTitle>
        <DialogContent>
          <Stack spacing={1.5} sx={{ mt: 1 }}>
            <Typography variant="body2" color="text.secondary">
              用户：{resetTarget?.username}
            </Typography>
            <TextField
              size="small"
              type="password"
              label="新密码（最少6位）"
              value={resetPasswordValue}
              onChange={(e) => setResetPasswordValue(e.target.value)}
              autoFocus
            />
          </Stack>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setResetTarget(null)}>取消</Button>
          <Button variant="contained" disabled={passwordSubmitting} onClick={() => void submitResetPassword()}>
            提交
          </Button>
        </DialogActions>
      </Dialog>
    </>
  );
}
