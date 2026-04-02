import { Alert, Box, Button, Dialog, DialogActions, DialogContent, DialogTitle, Stack, Typography } from "@mui/material";
import { useCallback, useEffect, useRef, useState } from "react";
import { FitAddon } from "@xterm/addon-fit";
import { Terminal } from "@xterm/xterm";
import "@xterm/xterm/css/xterm.css";

type Props = {
  open: boolean;
  title: string;
  createSession: () => Promise<{ wsPath: string }>;
  onClose: () => void;
};

type ConnState = "idle" | "connecting" | "connected" | "closed" | "error";

export default function TerminalDialog(props: Props) {
  const [mountEl, setMountEl] = useState<HTMLDivElement | null>(null);
  const termRef = useRef<Terminal | null>(null);
  const socketRef = useRef<WebSocket | null>(null);
  const connectRef = useRef<(() => void) | null>(null);
  const [connState, setConnState] = useState<ConnState>("idle");
  const [errorText, setErrorText] = useState("");

  const closeSocket = useCallback(() => {
    const ws = socketRef.current;
    socketRef.current = null;
    if (!ws) return;
    try {
      ws.close();
    } catch {
      // ignore close error
    }
  }, []);

  const connectSocket = useCallback(async () => {
    if (!props.open) return;
    const term = termRef.current;
    if (!term) return;
    closeSocket();
    setConnState("connecting");
    setErrorText("");

    let wsPath = "";
    try {
      const session = await props.createSession();
      wsPath = (session.wsPath || "").trim();
      if (!wsPath) {
        throw new Error("未获取到可用终端会话地址");
      }
    } catch (err) {
      const msg = err instanceof Error ? err.message : "创建终端会话失败";
      setConnState("error");
      setErrorText(msg);
      term.writeln(`\r\n[terminal] ${msg}`);
      return;
    }

    const url = new URL(wsPath, window.location.origin);
    url.protocol = window.location.protocol === "https:" ? "wss:" : "ws:";

    let ws: WebSocket;
    try {
      ws = new WebSocket(url.toString());
    } catch (err) {
      const msg = err instanceof Error ? err.message : "WebSocket 初始化失败";
      setConnState("error");
      setErrorText(msg);
      term.writeln(`\r\n[terminal] ${msg}`);
      return;
    }

    socketRef.current = ws;
    ws.onopen = () => {
      if (socketRef.current !== ws) return;
      setConnState("connected");
      term.focus();
      term.writeln("\r\n[terminal] connected");
    };
    ws.onmessage = (event) => {
      if (socketRef.current !== ws) return;
      if (typeof event.data === "string") {
        term.write(event.data);
      }
    };
    ws.onerror = () => {
      if (socketRef.current !== ws) return;
      setConnState("error");
      setErrorText("终端连接异常（会话可能已失效，请点击“重新连接”）");
      term.writeln("\r\n[terminal] websocket error");
    };
    ws.onclose = () => {
      if (socketRef.current === ws) {
        socketRef.current = null;
      }
      setConnState((prev) => (prev === "error" ? prev : "closed"));
      term.writeln("\r\n[terminal] disconnected");
    };
  }, [closeSocket, props.createSession, props.open]);

  connectRef.current = connectSocket;

  useEffect(() => {
    if (!props.open || !mountEl) return;
    const term = new Terminal({
      cursorBlink: true,
      fontSize: 13,
      lineHeight: 1.25,
      theme: {
        background: "#0f172a",
        foreground: "#e2e8f0"
      }
    });
    const fitAddon = new FitAddon();
    term.loadAddon(fitAddon);
    term.open(mountEl);
    fitAddon.fit();
    term.focus();
    term.writeln("[terminal] connecting...");

    const onResize = () => {
      fitAddon.fit();
    };
    window.addEventListener("resize", onResize);

    const dataDispose = term.onData((data) => {
      const ws = socketRef.current;
      if (!ws || ws.readyState !== WebSocket.OPEN) return;
      ws.send(data);
    });

    termRef.current = term;
    void connectSocket();

    return () => {
      dataDispose.dispose();
      window.removeEventListener("resize", onResize);
      closeSocket();
      termRef.current?.dispose();
      termRef.current = null;
      setConnState("idle");
      setErrorText("");
    };
  }, [closeSocket, connectSocket, mountEl, props.open]);

  function handleReconnect() {
    const term = termRef.current;
    if (term) {
      term.writeln("\r\n[terminal] reconnecting...");
    }
    void connectRef.current?.();
  }

  function handleClose() {
    closeSocket();
    props.onClose();
  }

  return (
    <Dialog open={props.open} onClose={handleClose} fullWidth maxWidth="lg">
      <DialogTitle>{props.title}</DialogTitle>
      <DialogContent>
        <Stack spacing={1.25} sx={{ mt: 1 }}>
          <Typography variant="body2" color="text.secondary">
            连接状态：{connState}
          </Typography>
          {errorText && <Alert severity="error">{errorText}</Alert>}
          <Box
            sx={{
              height: { xs: 320, md: 460 },
              borderRadius: 1,
              border: "1px solid",
              borderColor: "divider",
              overflow: "hidden",
              bgcolor: "#0f172a",
              p: 0.75
            }}
          >
            <Box ref={setMountEl} sx={{ height: "100%", width: "100%" }} />
          </Box>
        </Stack>
      </DialogContent>
      <DialogActions>
        <Button onClick={handleReconnect} disabled={connState === "connecting"}>
          {connState === "connecting" ? "连接中..." : "重新连接"}
        </Button>
        <Button onClick={handleClose}>关闭</Button>
      </DialogActions>
    </Dialog>
  );
}
