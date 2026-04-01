import { Alert, Box, Button, Card, CardContent, FormControl, InputLabel, MenuItem, Select, Stack, TextField, Typography } from "@mui/material";
import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { apiFetch } from "../lib/api";
import { useAuthStore } from "../stores/useAuthStore";

type AuthProviderItem = {
  name: string;
  type: string;
  isDefault: boolean;
};

export default function LoginPage() {
  const navigate = useNavigate();
  const login = useAuthStore((s) => s.login);
  const [username, setUsername] = useState("admin");
  const [password, setPassword] = useState("123456");
  const [provider, setProvider] = useState("local");
  const [providers, setProviders] = useState<AuthProviderItem[]>([{ name: "local", type: "local", isDefault: true }]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  useEffect(() => {
    void loadProviders();
  }, []);

  async function loadProviders() {
    try {
      const resp = await apiFetch("/api/v1/auth/providers/public");
      if (!resp.ok) {
        return;
      }
      const data = (await resp.json()) as { items?: AuthProviderItem[] };
      if (!data.items || data.items.length === 0) {
        return;
      }
      setProviders(data.items);
      const defaultProvider = data.items.find((item) => item.isDefault);
      if (defaultProvider) {
        setProvider(defaultProvider.name);
      } else {
        setProvider(data.items[0].name);
      }
    } catch {
      // fallback to local provider only
    }
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setLoading(true);
    setError("");
    try {
      const result = await login(username, password, provider);
      if (!result.ok) {
        setError(result.error || "登录失败");
        return;
      }
      navigate("/cluster", { replace: true });
    } finally {
      setLoading(false);
    }
  }

  return (
    <Box
      sx={{
        minHeight: "100vh",
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
        background: "linear-gradient(160deg, #eef4ff 0%, #f6fbff 55%, #ecf7f1 100%)",
        px: 2
      }}
    >
      <Card sx={{ width: "100%", maxWidth: 420, borderRadius: 3, boxShadow: "0 20px 60px rgba(21,101,192,0.12)" }}>
        <CardContent sx={{ p: 4 }}>
          <Stack component="form" spacing={2.2} onSubmit={handleSubmit}>
            <Typography variant="h5" sx={{ fontWeight: 800, color: "#123d77" }}>
              kubeManage 登录
            </Typography>
            <Typography variant="body2" color="text.secondary">
              首次启动默认管理员账号：admin / 123456
            </Typography>
            {error && <Alert severity="error">{error}</Alert>}
            <FormControl>
              <InputLabel id="provider-label">认证源</InputLabel>
              <Select
                labelId="provider-label"
                label="认证源"
                value={provider}
                onChange={(e) => setProvider(e.target.value)}
              >
                {providers.map((item) => (
                  <MenuItem key={item.name} value={item.name}>
                    {item.name} ({item.type})
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
            <TextField
              label="用户名"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              autoComplete="username"
              required
            />
            <TextField
              label="密码"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              autoComplete="current-password"
              required
            />
            <Button type="submit" variant="contained" size="large" disabled={loading}>
              {loading ? "登录中..." : "登录"}
            </Button>
          </Stack>
        </CardContent>
      </Card>
    </Box>
  );
}
