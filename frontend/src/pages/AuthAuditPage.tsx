import {
  Alert,
  Card,
  CardContent,
  FormControl,
  InputLabel,
  MenuItem,
  Select,
  Stack,
  Typography
} from "@mui/material";
import { useEffect, useState } from "react";
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

  useEffect(() => {
    const load = async () => {
      setError("");
      setAudits([]);
      if (!canAuditRead()) {
        return;
      }
      const resp = await apiFetch("/api/v1/audits");
      if (!resp.ok) {
        setError("当前角色无权限查看审计日志");
        return;
      }
      const data = (await resp.json()) as { items: AuditItem[] };
      setAudits(data.items);
    };
    void load();
  }, [role, canAuditRead]);

  return (
    <Stack spacing={2}>
      <Typography variant="h5">权限与审计（MVP）</Typography>

      <FormControl sx={{ width: 240 }}>
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

      {!canAuditRead() && (
        <Alert severity="info">当前角色仅可查看基础资源，无审计日志查看权限。</Alert>
      )}
      {error && <Alert severity="error">{error}</Alert>}

      {audits.map((a, idx) => (
        <Card key={`${a.time}-${idx}`} variant="outlined">
          <CardContent>
            <Typography variant="body2">
              {a.time} | {a.user} | {a.role}
            </Typography>
            <Typography variant="body2" color="text.secondary">
              {a.method} {a.path} {"->"} {a.statusCode}
            </Typography>
          </CardContent>
        </Card>
      ))}
    </Stack>
  );
}
