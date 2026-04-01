import {
  Alert,
  FormControl,
  InputLabel,
  MenuItem,
  Select,
  Stack,
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
  const [audits, setAudits] = useState<AuditItem[]>([]);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    const load = async () => {
      setError("");
      setAudits([]);
      if (!canAuditRead()) {
        return;
      }
      setLoading(true);
      try {
        const resp = await apiFetch("/api/v1/audits");
        if (!resp.ok) {
          setError("当前角色无权限查看审计日志");
          return;
        }
        const data = (await resp.json()) as { items: AuditItem[] };
        setAudits(data.items);
      } finally {
        setLoading(false);
      }
    };
    void load();
  }, [role, canAuditRead]);

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
      description="角色切换与审计日志查看（admin 可查看）"
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
      {!canAuditRead() && (
        <Alert severity="info" sx={{ m: 1.5 }}>
          当前角色仅可查看基础资源，无审计日志查看权限。
        </Alert>
      )}
      {error && <Alert severity="error" sx={{ m: 1.5 }}>{error}</Alert>}
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
