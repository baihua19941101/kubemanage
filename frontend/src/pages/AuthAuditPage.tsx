import {
  Alert,
  FormControl,
  InputLabel,
  MenuItem,
  Button,
  Select,
  Stack,
  TextField,
  Typography
} from "@mui/material";
import { useEffect, useMemo, useState } from "react";
import PageScaffold from "../components/framework/PageScaffold";
import ResourceTable from "../components/framework/ResourceTable";
import { apiFetch } from "../lib/api";
import { useAuthStore } from "../stores/useAuthStore";

type AuditItem = {
  time: string;
  user: string;
  role: string;
  method: string;
  path: string;
  statusCode: number;
};

export default function AuthAuditPage() {
  const role = useAuthStore((s) => s.role);
  const setRole = useAuthStore((s) => s.setRole);
  const canAuditRead = useAuthStore((s) => s.canAuditRead);
  const allowedNamespaces = useAuthStore((s) => s.allowedNamespaces);
  const [audits, setAudits] = useState<AuditItem[]>([]);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  const [userFilter, setUserFilter] = useState("");
  const [roleFilter, setRoleFilter] = useState("");
  const [pathFilter, setPathFilter] = useState("");
  const [methodFilter, setMethodFilter] = useState("");
  const [statusCodeFilter, setStatusCodeFilter] = useState("");
  const [limitFilter, setLimitFilter] = useState("50");

  useEffect(() => {
    void loadAudits();
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
        setError("当前角色无权限查看审计日志");
        return;
      }
      const data = (await resp.json()) as { items: AuditItem[] };
      setAudits(data.items);
    } finally {
      setLoading(false);
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

  const columns = useMemo(
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

  return (
    <PageScaffold
      title="权限与审计"
      description="角色切换、命名空间写权限说明与审计筛选（admin 可查看审计）"
      actions={
        <FormControl sx={{ width: 220 }}>
          <InputLabel id="role-label">当前角色</InputLabel>
          <Select
            labelId="role-label"
            value={role}
            label="当前角色"
            onChange={(e) => setRole(e.target.value)}
          >
            <MenuItem value="viewer">viewer</MenuItem>
            <MenuItem value="operator">operator</MenuItem>
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
      {!canAuditRead() && (
        <Alert severity="info" sx={{ m: 1.5 }}>
          当前角色仅可查看基础资源，无审计日志查看权限。
        </Alert>
      )}
      {error && <Alert severity="error" sx={{ m: 1.5 }}>{error}</Alert>}
      {canAuditRead() && (
        <Stack direction="row" spacing={1.5} sx={{ px: 1.5, pb: 1.5 }} useFlexGap flexWrap="wrap">
          <TextField size="small" label="用户" value={userFilter} onChange={(e) => setUserFilter(e.target.value)} />
          <TextField size="small" label="角色" value={roleFilter} onChange={(e) => setRoleFilter(e.target.value)} sx={{ width: 120 }} />
          <TextField size="small" label="方法" value={methodFilter} onChange={(e) => setMethodFilter(e.target.value)} sx={{ width: 120 }} />
          <TextField size="small" label="路径包含" value={pathFilter} onChange={(e) => setPathFilter(e.target.value)} sx={{ width: 240 }} />
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
        columns={columns}
      />
      {!loading && canAuditRead() && audits.length === 0 && (
        <Stack sx={{ p: 2 }}>
          <Typography color="text.secondary">暂无审计记录</Typography>
        </Stack>
      )}
    </PageScaffold>
  );
}
