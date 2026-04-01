import {
  Box,
  Divider,
  Drawer,
  IconButton,
  Stack,
  Typography
} from "@mui/material";
import CloseIcon from "@mui/icons-material/Close";
import type { ReactNode } from "react";

type Props = {
  open: boolean;
  title: string;
  onClose: () => void;
  actions?: ReactNode;
  children: ReactNode;
};

export default function DetailDrawer(props: Props) {
  return (
    <Drawer
      anchor="right"
      open={props.open}
      onClose={props.onClose}
      PaperProps={{ sx: { width: { xs: "100%", sm: 420 } } }}
    >
      <Stack direction="row" justifyContent="space-between" alignItems="center" sx={{ p: 2 }}>
        <Typography variant="h6">{props.title}</Typography>
        <Stack direction="row" spacing={1} alignItems="center">
          {props.actions}
          <IconButton onClick={props.onClose}>
            <CloseIcon />
          </IconButton>
        </Stack>
      </Stack>
      <Divider />
      <Box sx={{ p: 2 }}>{props.children}</Box>
    </Drawer>
  );
}
