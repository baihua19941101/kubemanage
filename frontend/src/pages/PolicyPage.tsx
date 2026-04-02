import { Alert, Button, Dialog, DialogActions, DialogContent, DialogTitle, Stack, TextField, Typography } from "@mui/material";
import { useEffect, useMemo, useState } from "react";
import DetailDrawer from "../components/framework/DetailDrawer";
import PageScaffold from "../components/framework/PageScaffold";
import ResourceTable from "../components/framework/ResourceTable";
import YamlDialog from "../components/framework/YamlDialog";
import { apiFetch, parseApiError } from "../lib/api";
import { useAuthStore } from "../stores/useAuthStore";
import { useResourceStore } from "../stores/useResourceStore";

export type PolicyMode = "limitranges" | "resourcequotas" | "networkpolicies";

type Props = {
  initialMode?: PolicyMode;
};

type LimitRangeItem = {
  name: string;
  namespace: string;
  limitsCount: number;
  defaultCpu: string;
  defaultMemory: string;
  age: string;
};

type ResourceQuotaItem = {
  name: string;
  namespace: string;
  hardPods: string;
  usedPods: string;
  hardCpu: string;
  usedCpu: string;
  hardMemory: string;
  usedMemory: string;
  hardPvcs: string;
  usedPvcs: string;
  age: string;
};

type NetworkPolicyItem = {
  name: string;
  namespace: string;
  podSelector: string;
  policyTypes: string;
  ingressRules: number;
  egressRules: number;
  age: string;
};

export default function PolicyPage({ initialMode = "limitranges" }: Props) {
  const limitRanges = useResourceStore((s) => s.limitRanges);
  const resourceQuotas = useResourceStore((s) => s.resourceQuotas);
  const networkPolicies = useResourceStore((s) => s.networkPolicies);
  const loading = useResourceStore((s) => s.loading);
  const storeError = useResourceStore((s) => s.error);
  const load = useResourceStore((s) => s.load);
  const canWorkloadWrite = useAuthStore((s) => s.canWorkloadWrite);

  const [mode, setMode] = useState<PolicyMode>(initialMode);
  const [keyword, setKeyword] = useState("");
  const [selectedName, setSelectedName] = useState("");
  const [yamlOpen, setYamlOpen] = useState(false);
  const [yamlTitle, setYamlTitle] = useState("");
  const [yamlText, setYamlText] = useState("");
  const [deleteOpen, setDeleteOpen] = useState(false);
  const [pageError, setPageError] = useState("");

  useEffect(() => {
    void load();
  }, [load]);

  useEffect(() => {
    setMode(initialMode);
    setSelectedName("");
    setPageError("");
  }, [initialMode]);

  const lowerKeyword = keyword.toLowerCase().trim();
  const filteredLimitRanges = useMemo(
    () => limitRanges.filter((item) => item.name.toLowerCase().includes(lowerKeyword)),
    [limitRanges, lowerKeyword]
  );
  const filteredResourceQuotas = useMemo(
    () => resourceQuotas.filter((item) => item.name.toLowerCase().includes(lowerKeyword)),
    [resourceQuotas, lowerKeyword]
  );
  const filteredNetworkPolicies = useMemo(
    () => networkPolicies.filter((item) => item.name.toLowerCase().includes(lowerKeyword)),
    [networkPolicies, lowerKeyword]
  );

  const selectedLimitRange = limitRanges.find((item) => item.name === selectedName) ?? null;
  const selectedResourceQuota = resourceQuotas.find((item) => item.name === selectedName) ?? null;
  const selectedNetworkPolicy = networkPolicies.find((item) => item.name === selectedName) ?? null;

  const limitRangeColumns = [
    { key: "name", header: "名称", render: (r: LimitRangeItem) => r.name },
    { key: "ns", header: "命名空间", render: (r: LimitRangeItem) => r.namespace },
    { key: "limits", header: "Limits", render: (r: LimitRangeItem) => r.limitsCount },
    { key: "cpu", header: "默认 CPU", render: (r: LimitRangeItem) => r.defaultCpu || "-" },
    { key: "memory", header: "默认内存", render: (r: LimitRangeItem) => r.defaultMemory || "-" }
  ];

  const resourceQuotaColumns = [
    { key: "name", header: "名称", render: (r: ResourceQuotaItem) => r.name },
    { key: "ns", header: "命名空间", render: (r: ResourceQuotaItem) => r.namespace },
    { key: "pods", header: "Pods(used/hard)", render: (r: ResourceQuotaItem) => `${r.usedPods || "0"}/${r.hardPods || "0"}` },
    { key: "cpu", header: "CPU(used/hard)", render: (r: ResourceQuotaItem) => `${r.usedCpu || "0"}/${r.hardCpu || "0"}` },
    { key: "memory", header: "Memory(used/hard)", render: (r: ResourceQuotaItem) => `${r.usedMemory || "0"}/${r.hardMemory || "0"}` }
  ];

  const networkPolicyColumns = [
    { key: "name", header: "名称", render: (r: NetworkPolicyItem) => r.name },
    { key: "ns", header: "命名空间", render: (r: NetworkPolicyItem) => r.namespace },
    { key: "selector", header: "Pod Selector", render: (r: NetworkPolicyItem) => r.podSelector || "<all>" },
    { key: "types", header: "PolicyTypes", render: (r: NetworkPolicyItem) => r.policyTypes || "-" },
    { key: "rules", header: "Rules(I/E)", render: (r: NetworkPolicyItem) => `${r.ingressRules}/${r.egressRules}` }
  ];

  const currentLabel =
    mode === "limitranges" ? "LimitRange" : mode === "resourcequotas" ? "ResourceQuota" : "NetworkPolicy";
  const endpointPrefix =
    mode === "limitranges" ? "limitranges" : mode === "resourcequotas" ? "resourcequotas" : "networkpolicies";

  async function openYaml() {
    if (!selectedName) return;
    setPageError("");
    try {
      const resp = await apiFetch(`/api/v1/${endpointPrefix}/${encodeURIComponent(selectedName)}/yaml`);
      if (!resp.ok) {
        throw await parseApiError(resp, `获取 ${currentLabel} YAML 失败`);
      }
      const text = await resp.text();
      setYamlTitle(`${currentLabel} YAML - ${selectedName}`);
      setYamlText(text);
      setYamlOpen(true);
    } catch (err) {
      setPageError(err instanceof Error ? err.message : `获取 ${currentLabel} YAML 失败`);
    }
  }

  async function saveYaml(content: string) {
    if (!selectedName) return;
    setPageError("");
    const resp = await apiFetch(`/api/v1/${endpointPrefix}/${encodeURIComponent(selectedName)}/yaml`, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        "X-Action-Confirm": "CONFIRM"
      },
      body: JSON.stringify({ yaml: content })
    });
    if (!resp.ok) {
      throw await parseApiError(resp, `保存 ${currentLabel} YAML 失败`);
    }
    setYamlOpen(false);
    await load();
  }

  async function confirmDelete() {
    if (!selectedName) return;
    setPageError("");
    const resp = await apiFetch(`/api/v1/${endpointPrefix}/${encodeURIComponent(selectedName)}`, {
      method: "DELETE",
      headers: {
        "X-Action-Confirm": "CONFIRM"
      }
    });
    if (!resp.ok) {
      throw await parseApiError(resp, `删除 ${currentLabel} 失败`);
    }
    setDeleteOpen(false);
    setSelectedName("");
    await load();
  }

  return (
    <>
      <PageScaffold
        title="Policy 管理"
        description="按资源管理 LimitRange / ResourceQuota / NetworkPolicy，支持详情与 YAML 查看下载"
        toolbar={
          <TextField
            size="small"
            label="按名称筛选"
            value={keyword}
            onChange={(e) => setKeyword(e.target.value)}
            sx={{ width: 280 }}
          />
        }
      >
        {(storeError || pageError) && <Alert severity="error" sx={{ m: 1.5 }}>{pageError || storeError}</Alert>}

        {mode === "limitranges" && (
          <ResourceTable
            loading={loading}
            rows={filteredLimitRanges}
            rowKey={(r) => r.name}
            columns={limitRangeColumns}
            onRowClick={(r) => setSelectedName(r.name)}
          />
        )}

        {mode === "resourcequotas" && (
          <ResourceTable
            loading={loading}
            rows={filteredResourceQuotas}
            rowKey={(r) => r.name}
            columns={resourceQuotaColumns}
            onRowClick={(r) => setSelectedName(r.name)}
          />
        )}

        {mode === "networkpolicies" && (
          <ResourceTable
            loading={loading}
            rows={filteredNetworkPolicies}
            rowKey={(r) => r.name}
            columns={networkPolicyColumns}
            onRowClick={(r) => setSelectedName(r.name)}
          />
        )}
      </PageScaffold>

      <DetailDrawer
        open={selectedName.length > 0}
        title={selectedName ? `${currentLabel} 详情 - ${selectedName}` : `${currentLabel} 详情`}
        onClose={() => setSelectedName("")}
      >
        {mode === "limitranges" && selectedLimitRange && (
          <Stack spacing={1}>
            <Typography variant="body2">名称：{selectedLimitRange.name}</Typography>
            <Typography variant="body2">命名空间：{selectedLimitRange.namespace}</Typography>
            <Typography variant="body2">Limit 条目：{selectedLimitRange.limitsCount}</Typography>
            <Typography variant="body2">默认 CPU：{selectedLimitRange.defaultCpu || "-"}</Typography>
            <Typography variant="body2">默认内存：{selectedLimitRange.defaultMemory || "-"}</Typography>
            <Typography variant="body2">Age：{selectedLimitRange.age}</Typography>
            <Stack direction="row" spacing={1} sx={{ pt: 1 }}>
              <Button size="small" onClick={openYaml}>{canWorkloadWrite() ? "查看/编辑 YAML" : "查看 YAML"}</Button>
              <Button size="small" component="a" href={`/api/v1/limitranges/${encodeURIComponent(selectedLimitRange.name)}/yaml/download`}>
                下载 YAML
              </Button>
              {canWorkloadWrite() && <Button size="small" color="error" onClick={() => setDeleteOpen(true)}>删除</Button>}
            </Stack>
          </Stack>
        )}

        {mode === "resourcequotas" && selectedResourceQuota && (
          <Stack spacing={1}>
            <Typography variant="body2">名称：{selectedResourceQuota.name}</Typography>
            <Typography variant="body2">命名空间：{selectedResourceQuota.namespace}</Typography>
            <Typography variant="body2">Pods(used/hard)：{selectedResourceQuota.usedPods || "0"}/{selectedResourceQuota.hardPods || "0"}</Typography>
            <Typography variant="body2">CPU(used/hard)：{selectedResourceQuota.usedCpu || "0"}/{selectedResourceQuota.hardCpu || "0"}</Typography>
            <Typography variant="body2">Memory(used/hard)：{selectedResourceQuota.usedMemory || "0"}/{selectedResourceQuota.hardMemory || "0"}</Typography>
            <Typography variant="body2">PVCs(used/hard)：{selectedResourceQuota.usedPvcs || "0"}/{selectedResourceQuota.hardPvcs || "0"}</Typography>
            <Typography variant="body2">Age：{selectedResourceQuota.age}</Typography>
            <Stack direction="row" spacing={1} sx={{ pt: 1 }}>
              <Button size="small" onClick={openYaml}>{canWorkloadWrite() ? "查看/编辑 YAML" : "查看 YAML"}</Button>
              <Button size="small" component="a" href={`/api/v1/resourcequotas/${encodeURIComponent(selectedResourceQuota.name)}/yaml/download`}>
                下载 YAML
              </Button>
              {canWorkloadWrite() && <Button size="small" color="error" onClick={() => setDeleteOpen(true)}>删除</Button>}
            </Stack>
          </Stack>
        )}

        {mode === "networkpolicies" && selectedNetworkPolicy && (
          <Stack spacing={1}>
            <Typography variant="body2">名称：{selectedNetworkPolicy.name}</Typography>
            <Typography variant="body2">命名空间：{selectedNetworkPolicy.namespace}</Typography>
            <Typography variant="body2">Pod Selector：{selectedNetworkPolicy.podSelector || "<all>"}</Typography>
            <Typography variant="body2">PolicyTypes：{selectedNetworkPolicy.policyTypes || "-"}</Typography>
            <Typography variant="body2">Ingress Rules：{selectedNetworkPolicy.ingressRules}</Typography>
            <Typography variant="body2">Egress Rules：{selectedNetworkPolicy.egressRules}</Typography>
            <Typography variant="body2">Age：{selectedNetworkPolicy.age}</Typography>
            <Stack direction="row" spacing={1} sx={{ pt: 1 }}>
              <Button size="small" onClick={openYaml}>{canWorkloadWrite() ? "查看/编辑 YAML" : "查看 YAML"}</Button>
              <Button size="small" component="a" href={`/api/v1/networkpolicies/${encodeURIComponent(selectedNetworkPolicy.name)}/yaml/download`}>
                下载 YAML
              </Button>
              {canWorkloadWrite() && <Button size="small" color="error" onClick={() => setDeleteOpen(true)}>删除</Button>}
            </Stack>
          </Stack>
        )}
      </DetailDrawer>

      <YamlDialog
        open={yamlOpen}
        title={yamlTitle}
        yaml={yamlText}
        onClose={() => setYamlOpen(false)}
        onSave={canWorkloadWrite() ? saveYaml : undefined}
      />

      <Dialog open={deleteOpen} onClose={() => setDeleteOpen(false)} fullWidth maxWidth="xs">
        <DialogTitle>删除确认</DialogTitle>
        <DialogContent>
          <Typography variant="body2">确认删除 {currentLabel}：{selectedName} ？该操作不可撤销。</Typography>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setDeleteOpen(false)}>取消</Button>
          <Button
            color="error"
            variant="contained"
            onClick={async () => {
              try {
                await confirmDelete();
              } catch (err) {
                setPageError(err instanceof Error ? err.message : `删除 ${currentLabel} 失败`);
              }
            }}
          >
            确认删除
          </Button>
        </DialogActions>
      </Dialog>
    </>
  );
}
