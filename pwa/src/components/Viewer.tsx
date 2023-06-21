import { Box, Typography, Tooltip } from "@mui/material";
import { JsonViewer } from '@textea/json-viewer';

import ExpandIcon from '@mui/icons-material/Expand';
import StartIcon from "@mui/icons-material/Start";

import { theme } from '@/App';

import type { State } from '@/App';
import type { JsonViewerProps } from '@textea/json-viewer'
import type { SxProps } from "@mui/system";

type ViewerProp = {
  state: State,
  sx: SxProps,
}

const Viewer = ({ state, sx }: ViewerProp) => {
  const { activeTabIndex, data, } = state;

  const fileName = Object.keys(data).sort()[activeTabIndex]
  const rv = window.Gastly.FromSourceCode(fileName, data[fileName])

  const viewerOpts: JsonViewerProps<{[key: string]: unknown}> = {
    value: {},
    theme: "light",
    displayDataTypes: false,
    displaySize: false,
    quotesOnKeys: false,
    objectSortKeys: false,
    enableClipboard: false,
    highlightUpdates: true,
    onSelect(path, value) {
      console.log(path, value)
    }
  }
  if (rv.err || rv.errMsg) {
    viewerOpts.value = { error: rv.err, errorMessage: rv.errMsg }
  } else {
    viewerOpts.value = JSON.parse(rv.ast??"{}");
  }
  return (
    <Box sx={{ width: "100%", height: "100%", ...sx }} >
      <Box sx={{
        borderBottom: 1, borderColor: 'divider',
        display: { xs: 'flex' }, justifyContent: 'space-between', alignItems: 'center',
      }}>
        <Tooltip title="View source code">
          <StartIcon sx={{ display: { md: 'none' }, ml: 2, transform: 'rotate(180deg)' }} />
        </Tooltip>
        <Typography component="p" sx={{ flexGrow: 1, ml: 2, py: 1.5 }}>
          AST
        </Typography>
        <Tooltip title="Expand" placement="bottom">
          <ExpandIcon sx={{ mr: 2 }} />
        </Tooltip>
      </Box>

      <Box sx={{ pl: 2 }} >
        <JsonViewer {...viewerOpts} sx={{ height: `calc(100vh - ${theme.spacing(14)})`, overflowY: "scroll" }} />
      </Box>
    </Box>
  )
}


export default Viewer;
