import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  TextField
} from "@mui/material";
import { useEffect, useState } from "react";

type Props = {
  open: boolean;
  title: string;
  yaml: string;
  onClose: () => void;
  onSave?: (yaml: string) => Promise<void> | void;
};

export default function YamlDialog(props: Props) {
  const [value, setValue] = useState(props.yaml);

  useEffect(() => {
    setValue(props.yaml);
  }, [props.yaml, props.open]);

  return (
    <Dialog open={props.open} onClose={props.onClose} fullWidth maxWidth="md">
      <DialogTitle>{props.title}</DialogTitle>
      <DialogContent>
        <TextField
          multiline
          minRows={16}
          fullWidth
          value={value}
          onChange={(e) => setValue(e.target.value)}
          sx={{ mt: 1 }}
        />
      </DialogContent>
      <DialogActions>
        <Button onClick={props.onClose}>关闭</Button>
        {props.onSave && (
          <Button
            variant="contained"
            onClick={async () => {
              await props.onSave?.(value);
            }}
          >
            保存 YAML
          </Button>
        )}
      </DialogActions>
    </Dialog>
  );
}
