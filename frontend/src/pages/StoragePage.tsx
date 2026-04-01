import { Alert, Stack, TextField, Typography } from "@mui/material";
import { useEffect, useMemo, useState } from "react";
import DetailDrawer from "../components/framework/DetailDrawer";
import PageScaffold from "../components/framework/PageScaffold";
import ResourceTable from "../components/framework/ResourceTable";
import { useResourceStore } from "../stores/useResourceStore";

export type StorageMode = "pvs" | "pvcs" | "storageclasses" | "configmaps" | "secrets";

type Props = {
  initialMode?: StorageMode;
};

type PVItem = {
  name: string;
  capacity: string;
  accessModes: string;
  reclaimPolicy: string;
  status: string;
  claimRef: string;
  storageClass: string;
  age: string;
};

type PVCItem = {
  name: string;
  namespace: string;
  status: string;
  volume: string;
  capacity: string;
  accessModes: string;
  storageClass: string;
  age: string;
};

type StorageClassItem = {
  name: string;
  provisioner: string;
  reclaimPolicy: string;
  volumeBindingMode: string;
  allowVolumeExpansion: boolean;
  age: string;
};

type ConfigMapItem = {
  name: string;
  namespace: string;
  dataCount: number;
  age: string;
};

type SecretItem = {
  name: string;
  namespace: string;
  type: string;
  data: Record<string, string>;
  age: string;
};

export default function StoragePage({ initialMode = "pvs" }: Props) {
  const pvs = useResourceStore((s) => s.pvs);
  const pvcs = useResourceStore((s) => s.pvcs);
  const storageClasses = useResourceStore((s) => s.storageClasses);
  const configMaps = useResourceStore((s) => s.configMaps);
  const secrets = useResourceStore((s) => s.secrets);
  const loading = useResourceStore((s) => s.loading);
  const error = useResourceStore((s) => s.error);
  const load = useResourceStore((s) => s.load);

  const [mode, setMode] = useState<StorageMode>(initialMode);
  const [keyword, setKeyword] = useState("");
  const [selectedName, setSelectedName] = useState("");

  useEffect(() => {
    void load();
  }, [load]);

  useEffect(() => {
    setMode(initialMode);
    setSelectedName("");
  }, [initialMode]);

  const lowerKeyword = keyword.toLowerCase().trim();

  const filteredPVs = useMemo(
    () => pvs.filter((item) => item.name.toLowerCase().includes(lowerKeyword)),
    [pvs, lowerKeyword]
  );
  const filteredPVCs = useMemo(
    () => pvcs.filter((item) => item.name.toLowerCase().includes(lowerKeyword)),
    [pvcs, lowerKeyword]
  );
  const filteredSCs = useMemo(
    () => storageClasses.filter((item) => item.name.toLowerCase().includes(lowerKeyword)),
    [storageClasses, lowerKeyword]
  );
  const filteredConfigMaps = useMemo(
    () => configMaps.filter((item) => item.name.toLowerCase().includes(lowerKeyword)),
    [configMaps, lowerKeyword]
  );
  const filteredSecrets = useMemo(
    () => secrets.filter((item) => item.name.toLowerCase().includes(lowerKeyword)),
    [secrets, lowerKeyword]
  );

  const selectedPV = pvs.find((item) => item.name === selectedName) ?? null;
  const selectedPVC = pvcs.find((item) => item.name === selectedName) ?? null;
  const selectedSC = storageClasses.find((item) => item.name === selectedName) ?? null;
  const selectedConfigMap = configMaps.find((item) => item.name === selectedName) ?? null;
  const selectedSecret = secrets.find((item) => item.name === selectedName) ?? null;

  const pvColumns = [
    { key: "name", header: "名称", render: (r: PVItem) => r.name },
    { key: "cap", header: "容量", render: (r: PVItem) => r.capacity },
    { key: "status", header: "状态", render: (r: PVItem) => r.status },
    { key: "claim", header: "Claim", render: (r: PVItem) => r.claimRef },
    { key: "sc", header: "StorageClass", render: (r: PVItem) => r.storageClass }
  ];

  const pvcColumns = [
    { key: "name", header: "名称", render: (r: PVCItem) => r.name },
    { key: "ns", header: "命名空间", render: (r: PVCItem) => r.namespace },
    { key: "status", header: "状态", render: (r: PVCItem) => r.status },
    { key: "vol", header: "Volume", render: (r: PVCItem) => r.volume },
    { key: "sc", header: "StorageClass", render: (r: PVCItem) => r.storageClass }
  ];

  const scColumns = [
    { key: "name", header: "名称", render: (r: StorageClassItem) => r.name },
    { key: "prov", header: "Provisioner", render: (r: StorageClassItem) => r.provisioner },
    { key: "reclaim", header: "Reclaim", render: (r: StorageClassItem) => r.reclaimPolicy },
    { key: "bind", header: "BindingMode", render: (r: StorageClassItem) => r.volumeBindingMode },
    { key: "expand", header: "扩容", render: (r: StorageClassItem) => (r.allowVolumeExpansion ? "是" : "否") }
  ];

  const configColumns = [
    { key: "name", header: "名称", render: (r: ConfigMapItem) => r.name },
    { key: "ns", header: "命名空间", render: (r: ConfigMapItem) => r.namespace },
    { key: "count", header: "数据项", render: (r: ConfigMapItem) => r.dataCount },
    { key: "age", header: "Age", render: (r: ConfigMapItem) => r.age }
  ];

  const secretColumns = [
    { key: "name", header: "名称", render: (r: SecretItem) => r.name },
    { key: "ns", header: "命名空间", render: (r: SecretItem) => r.namespace },
    { key: "type", header: "类型", render: (r: SecretItem) => r.type },
    { key: "age", header: "Age", render: (r: SecretItem) => r.age }
  ];

  const currentLabel =
    mode === "pvs"
      ? "PersistentVolumes"
      : mode === "pvcs"
        ? "PersistentVolumeClaims"
        : mode === "storageclasses"
          ? "StorageClasses"
          : mode === "configmaps"
            ? "ConfigMaps"
            : "Secrets";

  return (
    <>
      <PageScaffold
        title="存储管理"
        description="按资源管理 PV / PVC / StorageClass / ConfigMap / Secret"
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
        {error && <Alert severity="error" sx={{ m: 1.5 }}>{error}</Alert>}

        {mode === "pvs" && (
          <ResourceTable
            loading={loading}
            rows={filteredPVs}
            rowKey={(r) => r.name}
            columns={pvColumns}
            onRowClick={(r) => setSelectedName(r.name)}
          />
        )}

        {mode === "pvcs" && (
          <ResourceTable
            loading={loading}
            rows={filteredPVCs}
            rowKey={(r) => r.name}
            columns={pvcColumns}
            onRowClick={(r) => setSelectedName(r.name)}
          />
        )}

        {mode === "storageclasses" && (
          <ResourceTable
            loading={loading}
            rows={filteredSCs}
            rowKey={(r) => r.name}
            columns={scColumns}
            onRowClick={(r) => setSelectedName(r.name)}
          />
        )}

        {mode === "configmaps" && (
          <ResourceTable
            loading={loading}
            rows={filteredConfigMaps}
            rowKey={(r) => r.name}
            columns={configColumns}
            onRowClick={(r) => setSelectedName(r.name)}
          />
        )}

        {mode === "secrets" && (
          <ResourceTable
            loading={loading}
            rows={filteredSecrets}
            rowKey={(r) => r.name}
            columns={secretColumns}
            onRowClick={(r) => setSelectedName(r.name)}
          />
        )}
      </PageScaffold>

      <DetailDrawer
        open={selectedName.length > 0}
        title={selectedName ? `${currentLabel} 详情 - ${selectedName}` : `${currentLabel} 详情`}
        onClose={() => setSelectedName("")}
      >
        {mode === "pvs" && selectedPV && (
          <Stack spacing={1}>
            <Typography variant="body2">名称：{selectedPV.name}</Typography>
            <Typography variant="body2">容量：{selectedPV.capacity}</Typography>
            <Typography variant="body2">状态：{selectedPV.status}</Typography>
            <Typography variant="body2">访问模式：{selectedPV.accessModes}</Typography>
            <Typography variant="body2">ReclaimPolicy：{selectedPV.reclaimPolicy}</Typography>
            <Typography variant="body2">StorageClass：{selectedPV.storageClass}</Typography>
            <Typography variant="body2">ClaimRef：{selectedPV.claimRef}</Typography>
          </Stack>
        )}

        {mode === "pvcs" && selectedPVC && (
          <Stack spacing={1}>
            <Typography variant="body2">名称：{selectedPVC.name}</Typography>
            <Typography variant="body2">命名空间：{selectedPVC.namespace}</Typography>
            <Typography variant="body2">状态：{selectedPVC.status}</Typography>
            <Typography variant="body2">容量：{selectedPVC.capacity}</Typography>
            <Typography variant="body2">Volume：{selectedPVC.volume}</Typography>
            <Typography variant="body2">StorageClass：{selectedPVC.storageClass}</Typography>
            <Typography variant="body2">访问模式：{selectedPVC.accessModes}</Typography>
          </Stack>
        )}

        {mode === "storageclasses" && selectedSC && (
          <Stack spacing={1}>
            <Typography variant="body2">名称：{selectedSC.name}</Typography>
            <Typography variant="body2">Provisioner：{selectedSC.provisioner}</Typography>
            <Typography variant="body2">ReclaimPolicy：{selectedSC.reclaimPolicy}</Typography>
            <Typography variant="body2">BindingMode：{selectedSC.volumeBindingMode}</Typography>
            <Typography variant="body2">允许扩容：{selectedSC.allowVolumeExpansion ? "是" : "否"}</Typography>
          </Stack>
        )}

        {mode === "configmaps" && selectedConfigMap && (
          <Stack spacing={1}>
            <Typography variant="body2">名称：{selectedConfigMap.name}</Typography>
            <Typography variant="body2">命名空间：{selectedConfigMap.namespace}</Typography>
            <Typography variant="body2">数据项：{selectedConfigMap.dataCount}</Typography>
            <Typography variant="body2">Age：{selectedConfigMap.age}</Typography>
          </Stack>
        )}

        {mode === "secrets" && selectedSecret && (
          <Stack spacing={1}>
            <Typography variant="body2">名称：{selectedSecret.name}</Typography>
            <Typography variant="body2">命名空间：{selectedSecret.namespace}</Typography>
            <Typography variant="body2">类型：{selectedSecret.type}</Typography>
            <Typography variant="body2">Age：{selectedSecret.age}</Typography>
            <Typography variant="body2" color="text.secondary">脱敏数据：</Typography>
            {Object.entries(selectedSecret.data).map(([k, v]) => (
              <Typography key={k} variant="body2" color="text.secondary">
                {k}: {v}
              </Typography>
            ))}
          </Stack>
        )}
      </DetailDrawer>
    </>
  );
}
