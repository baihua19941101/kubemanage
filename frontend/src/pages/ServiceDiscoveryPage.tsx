import { Alert, Stack, TextField, Typography } from "@mui/material";
import { useEffect, useMemo, useState } from "react";
import DetailDrawer from "../components/framework/DetailDrawer";
import PageScaffold from "../components/framework/PageScaffold";
import ResourceTable from "../components/framework/ResourceTable";
import { useResourceStore } from "../stores/useResourceStore";

export type ServiceDiscoveryMode = "services" | "ingresses" | "hpas";

type Props = {
  initialMode?: ServiceDiscoveryMode;
};

type ServiceItem = {
  name: string;
  namespace: string;
  type: string;
  clusterIP: string;
  ports: string;
  pods: number;
  age: string;
};

type IngressItem = {
  name: string;
  namespace: string;
  className: string;
  hosts: string[];
  address: string;
  tls: boolean;
  age: string;
};

type HPAItem = {
  name: string;
  namespace: string;
  targetKind: string;
  targetName: string;
  minReplicas: number;
  maxReplicas: number;
  currentReplicas: number;
  targetCPUPercent: number;
  currentCPUPercent: number;
  age: string;
};

type HPATarget = {
  kind: string;
  name: string;
  namespace: string;
  currentReplicas: number;
  desiredReplicas: number;
};

export default function ServiceDiscoveryPage({ initialMode = "services" }: Props) {
  const services = useResourceStore((s) => s.services);
  const ingresses = useResourceStore((s) => s.ingresses);
  const hpas = useResourceStore((s) => s.hpas);
  const loading = useResourceStore((s) => s.loading);
  const error = useResourceStore((s) => s.error);
  const load = useResourceStore((s) => s.load);
  const getIngressServices = useResourceStore((s) => s.getIngressServices);
  const getHPATarget = useResourceStore((s) => s.getHPATarget);

  const [mode, setMode] = useState<ServiceDiscoveryMode>(initialMode);
  const [keyword, setKeyword] = useState("");
  const [selectedName, setSelectedName] = useState("");
  const [ingressServices, setIngressServices] = useState<ServiceItem[]>([]);
  const [hpaTarget, setHPATarget] = useState<HPATarget | null>(null);

  useEffect(() => {
    void load();
  }, [load]);

  useEffect(() => {
    setMode(initialMode);
    setSelectedName("");
    setIngressServices([]);
    setHPATarget(null);
  }, [initialMode]);

  const lowerKeyword = keyword.toLowerCase().trim();
  const filteredServices = useMemo(
    () => services.filter((item) => item.name.toLowerCase().includes(lowerKeyword)),
    [services, lowerKeyword]
  );
  const filteredIngresses = useMemo(
    () => ingresses.filter((item) => item.name.toLowerCase().includes(lowerKeyword)),
    [ingresses, lowerKeyword]
  );
  const filteredHPAs = useMemo(
    () => hpas.filter((item) => item.name.toLowerCase().includes(lowerKeyword)),
    [hpas, lowerKeyword]
  );

  const selectedService = services.find((item) => item.name === selectedName) ?? null;
  const selectedIngress = ingresses.find((item) => item.name === selectedName) ?? null;
  const selectedHPA = hpas.find((item) => item.name === selectedName) ?? null;

  const serviceColumns = [
    { key: "name", header: "名称", render: (r: ServiceItem) => r.name },
    { key: "ns", header: "命名空间", render: (r: ServiceItem) => r.namespace },
    { key: "type", header: "类型", render: (r: ServiceItem) => r.type },
    { key: "ports", header: "端口", render: (r: ServiceItem) => r.ports },
    { key: "age", header: "Age", render: (r: ServiceItem) => r.age }
  ];

  const ingressColumns = [
    { key: "name", header: "名称", render: (r: IngressItem) => r.name },
    { key: "ns", header: "命名空间", render: (r: IngressItem) => r.namespace },
    { key: "class", header: "Class", render: (r: IngressItem) => r.className },
    { key: "hosts", header: "Hosts", render: (r: IngressItem) => r.hosts.join(", ") },
    { key: "tls", header: "TLS", render: (r: IngressItem) => (r.tls ? "是" : "否") }
  ];

  const hpaColumns = [
    { key: "name", header: "名称", render: (r: HPAItem) => r.name },
    { key: "ns", header: "命名空间", render: (r: HPAItem) => r.namespace },
    { key: "target", header: "目标", render: (r: HPAItem) => `${r.targetKind}/${r.targetName}` },
    { key: "rep", header: "副本", render: (r: HPAItem) => `${r.currentReplicas} (${r.minReplicas}-${r.maxReplicas})` },
    { key: "cpu", header: "CPU", render: (r: HPAItem) => `${r.currentCPUPercent}% / ${r.targetCPUPercent}%` }
  ];

  const currentLabel = mode === "services" ? "Service" : mode === "ingresses" ? "Ingress" : "HPA";

  return (
    <>
      <PageScaffold
        title="服务发现"
        description="按资源管理 Service / Ingress / HPA，支持关联关系查看"
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

        {mode === "services" && (
          <ResourceTable
            loading={loading}
            rows={filteredServices}
            rowKey={(r) => r.name}
            columns={serviceColumns}
            onRowClick={(r) => setSelectedName(r.name)}
          />
        )}

        {mode === "ingresses" && (
          <ResourceTable
            loading={loading}
            rows={filteredIngresses}
            rowKey={(r) => r.name}
            columns={ingressColumns}
            onRowClick={async (r) => {
              setSelectedName(r.name);
              setIngressServices(await getIngressServices(r.name));
            }}
          />
        )}

        {mode === "hpas" && (
          <ResourceTable
            loading={loading}
            rows={filteredHPAs}
            rowKey={(r) => r.name}
            columns={hpaColumns}
            onRowClick={async (r) => {
              setSelectedName(r.name);
              setHPATarget(await getHPATarget(r.name));
            }}
          />
        )}
      </PageScaffold>

      <DetailDrawer
        open={selectedName.length > 0}
        title={selectedName ? `${currentLabel} 详情 - ${selectedName}` : `${currentLabel} 详情`}
        onClose={() => {
          setSelectedName("");
          setIngressServices([]);
          setHPATarget(null);
        }}
      >
        {mode === "services" && selectedService && (
          <Stack spacing={1}>
            <Typography variant="body2">名称：{selectedService.name}</Typography>
            <Typography variant="body2">命名空间：{selectedService.namespace}</Typography>
            <Typography variant="body2">类型：{selectedService.type}</Typography>
            <Typography variant="body2">ClusterIP：{selectedService.clusterIP}</Typography>
            <Typography variant="body2">端口：{selectedService.ports}</Typography>
            <Typography variant="body2">关联 Pod：{selectedService.pods}</Typography>
          </Stack>
        )}

        {mode === "ingresses" && selectedIngress && (
          <Stack spacing={1}>
            <Typography variant="body2">名称：{selectedIngress.name}</Typography>
            <Typography variant="body2">命名空间：{selectedIngress.namespace}</Typography>
            <Typography variant="body2">Class：{selectedIngress.className}</Typography>
            <Typography variant="body2">Hosts：{selectedIngress.hosts.join(", ")}</Typography>
            <Typography variant="body2">Address：{selectedIngress.address}</Typography>
            <Typography variant="body2">TLS：{selectedIngress.tls ? "是" : "否"}</Typography>
            <Typography variant="body2" color="text.secondary">关联 Service：</Typography>
            {ingressServices.map((item) => (
              <Typography key={item.name} variant="body2" color="text.secondary">
                {item.namespace}/{item.name} ({item.ports})
              </Typography>
            ))}
          </Stack>
        )}

        {mode === "hpas" && selectedHPA && (
          <Stack spacing={1}>
            <Typography variant="body2">名称：{selectedHPA.name}</Typography>
            <Typography variant="body2">命名空间：{selectedHPA.namespace}</Typography>
            <Typography variant="body2">目标：{selectedHPA.targetKind}/{selectedHPA.targetName}</Typography>
            <Typography variant="body2">副本范围：{selectedHPA.minReplicas} - {selectedHPA.maxReplicas}</Typography>
            <Typography variant="body2">当前副本：{selectedHPA.currentReplicas}</Typography>
            <Typography variant="body2">CPU：{selectedHPA.currentCPUPercent}% / {selectedHPA.targetCPUPercent}%</Typography>
            {hpaTarget && (
              <>
                <Typography variant="body2" color="text.secondary">关联目标：</Typography>
                <Typography variant="body2" color="text.secondary">
                  {hpaTarget.namespace}/{hpaTarget.kind}/{hpaTarget.name}
                </Typography>
                <Typography variant="body2" color="text.secondary">
                  当前/期望副本：{hpaTarget.currentReplicas}/{hpaTarget.desiredReplicas}
                </Typography>
              </>
            )}
          </Stack>
        )}
      </DetailDrawer>
    </>
  );
}
