import {
  Alert,
  Accordion,
  AccordionDetails,
  AccordionSummary,
  Box,
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  Stack,
  Typography
} from "@mui/material";
import ExpandMoreIcon from "@mui/icons-material/ExpandMore";
import UploadFileIcon from "@mui/icons-material/UploadFile";
import DownloadIcon from "@mui/icons-material/Download";
import RestartAltIcon from "@mui/icons-material/RestartAlt";
import DifferenceIcon from "@mui/icons-material/Difference";
import { useEffect, useMemo, useRef, useState } from "react";
import AceEditor from "react-ace";
import { parse } from "yaml";

import "ace-builds/src-noconflict/mode-yaml";
import "ace-builds/src-noconflict/theme-github";
import "ace-builds/src-noconflict/ext-language_tools";

type Props = {
  open: boolean;
  title: string;
  yaml: string;
  onClose: () => void;
  onSave?: (yaml: string) => Promise<void> | void;
  saving?: boolean;
  saveMeta?: {
    lastSavedAt?: string;
    lastRequestId?: string;
    history?: Array<{ at: string; requestId?: string }>;
  };
};

type ViewMode = "editor" | "tree";
type DiffLine = { kind: "added" | "removed" | "same"; text: string };

export default function YamlDialog(props: Props) {
  const [value, setValue] = useState(props.yaml);
  const [mode, setMode] = useState<ViewMode>("editor");
  const [expandedMap, setExpandedMap] = useState<Record<string, boolean>>({});
  const [showDiff, setShowDiff] = useState(false);
  const fileInputRef = useRef<HTMLInputElement | null>(null);

  useEffect(() => {
    setValue(props.yaml);
    setMode("editor");
    setExpandedMap({});
    setShowDiff(false);
  }, [props.yaml, props.open]);

  const parsed = useMemo(() => {
    try {
      return { data: parse(value) as unknown, error: "" };
    } catch (err) {
      return {
        data: null as unknown,
        error: err instanceof Error ? err.message : "YAML 解析失败"
      };
    }
  }, [value]);

  const dirty = useMemo(() => value !== props.yaml, [props.yaml, value]);

  const diffLines = useMemo(() => {
    const before = props.yaml.split("\n");
    const after = value.split("\n");
    const maxLen = Math.max(before.length, after.length);
    const rows: DiffLine[] = [];
    for (let i = 0; i < maxLen; i += 1) {
      const b = before[i];
      const a = after[i];
      if (b === a) {
        if (a !== undefined) {
          rows.push({ kind: "same", text: a });
        }
        continue;
      }
      if (b !== undefined) {
        rows.push({ kind: "removed", text: b });
      }
      if (a !== undefined) {
        rows.push({ kind: "added", text: a });
      }
    }
    return rows;
  }, [props.yaml, value]);

  const diffSummary = useMemo(() => {
    let added = 0;
    let removed = 0;
    diffLines.forEach((line) => {
      if (line.kind === "added") added += 1;
      if (line.kind === "removed") removed += 1;
    });
    return { added, removed };
  }, [diffLines]);

  function isContainer(node: unknown) {
    return Array.isArray(node) || (node !== null && typeof node === "object");
  }

  function containerCount(node: unknown) {
    if (Array.isArray(node)) {
      return node.length;
    }
    if (node !== null && typeof node === "object") {
      return Object.keys(node as Record<string, unknown>).length;
    }
    return 0;
  }

  function formatScalar(node: unknown) {
    if (node === null) return "null";
    if (node === undefined) return "undefined";
    if (typeof node === "string") return node.length === 0 ? '""' : node;
    return String(node);
  }

  function getExpanded(path: string, depth: number) {
    if (expandedMap[path] !== undefined) {
      return expandedMap[path];
    }
    return depth <= 1;
  }

  function setExpanded(path: string, next: boolean) {
    setExpandedMap((prev) => ({ ...prev, [path]: next }));
  }

  function collectContainerPaths(node: unknown, path: string, out: string[]) {
    if (!isContainer(node)) {
      return;
    }
    out.push(path);
    if (Array.isArray(node)) {
      node.forEach((item, idx) => collectContainerPaths(item, `${path}[${idx}]`, out));
      return;
    }
    Object.entries(node as Record<string, unknown>).forEach(([key, item]) => {
      collectContainerPaths(item, `${path}.${key}`, out);
    });
  }

  function expandOrCollapseAll(next: boolean) {
    const paths: string[] = [];
    collectContainerPaths(parsed.data, "$", paths);
    const nextMap: Record<string, boolean> = {};
    paths.forEach((path) => {
      nextMap[path] = next;
    });
    setExpandedMap(nextMap);
  }

  function renderNode(label: string, node: unknown, path: string, depth: number): React.JSX.Element {
    if (!isContainer(node)) {
      return (
        <Box
          key={path}
          sx={{
            display: "grid",
            gridTemplateColumns: "220px minmax(0, 1fr)",
            gap: 1.5,
            py: 0.75,
            px: 1,
            borderBottom: "1px solid",
            borderColor: "divider"
          }}
        >
          <Typography sx={{ fontFamily: "ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, Liberation Mono, monospace", color: "#3a5b88", fontSize: 13 }}>
            {label}
          </Typography>
          <Typography sx={{ fontFamily: "ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, Liberation Mono, monospace", fontSize: 13, whiteSpace: "pre-wrap", wordBreak: "break-word" }}>
            {formatScalar(node)}
          </Typography>
        </Box>
      );
    }

    const entries = Array.isArray(node)
      ? node.map((item, idx) => [`[${idx}]`, item] as const)
      : Object.entries(node as Record<string, unknown>);

    const expanded = getExpanded(path, depth);

    return (
      <Accordion
        key={path}
        expanded={expanded}
        onChange={(_, next) => setExpanded(path, next)}
        disableGutters
        elevation={0}
        sx={{ border: "1px solid", borderColor: "divider", borderRadius: 1, mb: 1, overflow: "hidden", "&:before": { display: "none" } }}
      >
        <AccordionSummary expandIcon={<ExpandMoreIcon />} sx={{ minHeight: 38, bgcolor: depth === 0 ? "#f4f7fb" : "#fafbfd", "& .MuiAccordionSummary-content": { my: 0.6 } }}>
          <Stack direction="row" spacing={1.5} alignItems="center">
            <Typography sx={{ fontFamily: "ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, Liberation Mono, monospace", fontSize: 13, color: "#123d77", fontWeight: 700 }}>
              {label}
            </Typography>
            <Typography variant="caption" color="text.secondary">
              {Array.isArray(node) ? `array(${containerCount(node)})` : `object(${containerCount(node)})`}
            </Typography>
          </Stack>
        </AccordionSummary>
        <AccordionDetails sx={{ p: 1 }}>
          {entries.length === 0 ? (
            <Typography variant="body2" color="text.secondary" sx={{ px: 1, py: 0.5 }}>
              (empty)
            </Typography>
          ) : (
            entries.map(([key, item]) => renderNode(key, item, `${path}.${key}`, depth + 1))
          )}
        </AccordionDetails>
      </Accordion>
    );
  }

  async function handleImportFile(file: File | null) {
    if (!file) return;
    const text = await file.text();
    setValue(text);
    setMode("editor");
  }

  function handleDownloadYaml() {
    const blob = new Blob([value], { type: "application/yaml;charset=utf-8" });
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = "resource.yaml";
    document.body.appendChild(a);
    a.click();
    a.remove();
    URL.revokeObjectURL(url);
  }

  function formatLocalTime(iso: string) {
    const dt = new Date(iso);
    if (Number.isNaN(dt.getTime())) return iso;
    return dt.toLocaleString();
  }

  return (
    <Dialog open={props.open} onClose={props.onClose} fullWidth maxWidth="lg">
      <DialogTitle>{props.title}</DialogTitle>
      <DialogContent>
        <Stack direction="row" spacing={1} flexWrap="wrap" sx={{ mt: 1, mb: 1 }}>
          <Button variant={mode === "editor" ? "contained" : "outlined"} size="small" onClick={() => setMode("editor")}>
            YAML 编辑器
          </Button>
          <Button variant={mode === "tree" ? "contained" : "outlined"} size="small" onClick={() => setMode("tree")}>
            结构视图
          </Button>
          <Button size="small" startIcon={<UploadFileIcon />} onClick={() => fileInputRef.current?.click()}>
            导入文件
          </Button>
          <Button size="small" startIcon={<DownloadIcon />} onClick={handleDownloadYaml}>
            下载
          </Button>
          <Button size="small" startIcon={<RestartAltIcon />} disabled={!dirty} onClick={() => setValue(props.yaml)}>
            还原
          </Button>
          <Button size="small" startIcon={<DifferenceIcon />} disabled={!dirty} onClick={() => setShowDiff((v) => !v)}>
            {showDiff ? "隐藏变更" : "查看变更"}
          </Button>
          {mode === "tree" && (
            <>
              <Button size="small" onClick={() => expandOrCollapseAll(true)}>全部展开</Button>
              <Button size="small" onClick={() => expandOrCollapseAll(false)}>全部折叠</Button>
            </>
          )}
        </Stack>

        <input
          ref={fileInputRef}
          type="file"
          accept=".yaml,.yml,text/yaml,text/plain"
          style={{ display: "none" }}
          onChange={async (e) => {
            await handleImportFile(e.target.files?.[0] ?? null);
            e.currentTarget.value = "";
          }}
        />

        {dirty && (
          <Typography variant="caption" color="text.secondary">
            已修改 | +{diffSummary.added} / -{diffSummary.removed}
          </Typography>
        )}

        {showDiff && dirty && (
          <Box sx={{ mt: 1, mb: 1, border: "1px solid", borderColor: "divider", borderRadius: 1, maxHeight: 180, overflow: "auto", bgcolor: "#fcfdff" }}>
            {diffLines.map((line, idx) => (
              <Box
                key={`${line.kind}-${idx}`}
                sx={{
                  px: 1,
                  py: 0.2,
                  fontFamily: "ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, Liberation Mono, monospace",
                  fontSize: 12,
                  whiteSpace: "pre-wrap",
                  bgcolor: line.kind === "added" ? "#edf7ed" : line.kind === "removed" ? "#fdecea" : "transparent",
                  color: line.kind === "added" ? "#1e6b2a" : line.kind === "removed" ? "#9f1d1d" : "text.primary"
                }}
              >
                {line.kind === "added" ? "+ " : line.kind === "removed" ? "- " : "  "}
                {line.text}
              </Box>
            ))}
          </Box>
        )}

        {mode === "editor" ? (
          <Box sx={{ mt: 1, border: "1px solid", borderColor: "divider", borderRadius: 1, overflow: "hidden" }}>
            <AceEditor
              mode="yaml"
              theme="github"
              name="yaml-editor"
              width="100%"
              height="62vh"
              value={value}
              onChange={(next) => setValue(next)}
              fontSize={13}
              showPrintMargin={false}
              showGutter
              highlightActiveLine
              setOptions={{
                useWorker: false,
                tabSize: 2,
                wrap: true,
                showLineNumbers: true,
                enableBasicAutocompletion: true,
                enableLiveAutocompletion: false,
                enableSnippets: false,
                foldStyle: "markbegin",
                displayIndentGuides: true,
                showFoldWidgets: true
              }}
              editorProps={{ $blockScrolling: true }}
            />
            {parsed.error && <Alert severity="warning" sx={{ m: 1 }}>当前 YAML 存在语法错误：{parsed.error}</Alert>}
          </Box>
        ) : (
          <Box sx={{ mt: 1, maxHeight: "62vh", overflow: "auto", pr: 0.5 }}>
            {parsed.error ? (
              <Alert severity="warning">YAML 解析失败，请先切换到“YAML 编辑器”修正。{parsed.error}</Alert>
            ) : (
              renderNode("root", parsed.data, "$", 0)
            )}
          </Box>
        )}
      </DialogContent>
      <DialogActions>
        <Box sx={{ mr: "auto", minWidth: 320 }}>
          {props.saveMeta?.lastSavedAt && (
            <Typography variant="caption" color="text.secondary" display="block">
              最近保存：{formatLocalTime(props.saveMeta.lastSavedAt)}
              {props.saveMeta.lastRequestId ? ` | requestId: ${props.saveMeta.lastRequestId}` : ""}
            </Typography>
          )}
          {(props.saveMeta?.history || []).length > 0 && (
            <Typography variant="caption" color="text.secondary" display="block">
              历史：
              {(props.saveMeta?.history || []).slice(0, 3).map((item, idx) => (
                <span key={`${item.at}-${idx}`}>
                  {idx > 0 ? " ; " : " "}
                  {formatLocalTime(item.at)}
                  {item.requestId ? ` (${item.requestId})` : ""}
                </span>
              ))}
            </Typography>
          )}
        </Box>
        <Button onClick={props.onClose}>关闭</Button>
        {props.onSave && (
          <Button
            variant="contained"
            disabled={props.saving}
            onClick={async () => {
              await props.onSave?.(value);
            }}
          >
            {props.saving ? "保存中..." : "保存 YAML"}
          </Button>
        )}
      </DialogActions>
    </Dialog>
  );
}
