import {
  Box,
  CircularProgress,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Typography
} from "@mui/material";
import type { ReactNode } from "react";

type Column<T> = {
  key: string;
  header: string;
  width?: number | string;
  render: (row: T) => ReactNode;
};

type Props<T> = {
  columns: Column<T>[];
  rows: T[];
  loading?: boolean;
  rowKey: (row: T) => string;
  onRowClick?: (row: T) => void;
};

export default function ResourceTable<T>(props: Props<T>) {
  if (props.loading) {
    return (
      <Box sx={{ py: 6, textAlign: "center" }}>
        <CircularProgress />
      </Box>
    );
  }

  if (props.rows.length === 0) {
    return (
      <Box sx={{ py: 6, textAlign: "center" }}>
        <Typography color="text.secondary">暂无数据</Typography>
      </Box>
    );
  }

  return (
    <TableContainer>
      <Table size="small">
        <TableHead>
          <TableRow>
            {props.columns.map((c) => (
              <TableCell
                key={c.key}
                sx={{
                  width: c.width,
                  fontWeight: 700,
                  bgcolor: "#f2f6fc",
                  borderBottomColor: "#dde6f3"
                }}
              >
                {c.header}
              </TableCell>
            ))}
          </TableRow>
        </TableHead>
        <TableBody>
          {props.rows.map((row) => (
            <TableRow
              key={props.rowKey(row)}
              hover={Boolean(props.onRowClick)}
              onClick={() => props.onRowClick?.(row)}
              sx={{
                cursor: props.onRowClick ? "pointer" : "default",
                "&:nth-of-type(even)": { bgcolor: "#fbfdff" }
              }}
            >
              {props.columns.map((c) => (
                <TableCell key={c.key}>{c.render(row)}</TableCell>
              ))}
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </TableContainer>
  );
}
