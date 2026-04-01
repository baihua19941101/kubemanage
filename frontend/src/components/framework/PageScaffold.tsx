import { Box, Paper, Stack, Typography } from "@mui/material";
import type { ReactNode } from "react";

type Props = {
  title: string;
  description?: string;
  actions?: ReactNode;
  toolbar?: ReactNode;
  children: ReactNode;
};

export default function PageScaffold(props: Props) {
  return (
    <Stack spacing={2}>
      <Paper variant="outlined" sx={{ p: 2 }}>
        <Stack
          direction={{ xs: "column", sm: "row" }}
          justifyContent="space-between"
          alignItems={{ xs: "flex-start", sm: "center" }}
          spacing={1.5}
        >
          <Box>
            <Typography variant="h5" sx={{ fontWeight: 700 }}>
              {props.title}
            </Typography>
            {props.description && (
              <Typography variant="body2" color="text.secondary">
                {props.description}
              </Typography>
            )}
          </Box>
          {props.actions}
        </Stack>
      </Paper>

      {props.toolbar && (
        <Paper variant="outlined" sx={{ p: 1.5 }}>
          {props.toolbar}
        </Paper>
      )}

      <Paper variant="outlined">{props.children}</Paper>
    </Stack>
  );
}
