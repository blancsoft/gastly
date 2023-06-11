import * as React from "react";
import PropTypes from 'prop-types';
import { Box, Typography, Tooltip } from "@mui/material";
import { JsonViewer } from '@textea/json-viewer'

import ExpandIcon from '@mui/icons-material/Expand';
import StartIcon from "@mui/icons-material/Start";


const Viewer = ({ state, sx }) => {
  const { activeTabIndex, data, } = state;

  const fileName = Object.keys(data).sort()[activeTabIndex]
  const rv = window.Gastly.FromSourceCode(fileName, data[fileName])
  console.log("RV: ", rv)

  const viewerOpts = {
    theme: "light",
    displayDataTypes: false,
    displaySize: false,
    quotesOnKeys: false,
    objectSortKeys: false,
    enableClipboard: false,
    onSelect(path, value) {
      console.log(path, value)
    }
  }
  if (rv.err || rv.errMsg) {
    viewerOpts.value = { error: rv.err, errorMessage: rv.errMsg }
  } else {
    viewerOpts.value = JSON.parse(rv.ast)
  }
  return (
    <Box sx={{ width: "100%", height: "100%", overflow: "hidden", ...sx }} >
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

      <Box sx={{ pl: 2, pt: 2 }} >
        <JsonViewer {...viewerOpts} />
      </Box>
    </Box>
  )
}

Viewer.propTypes = {
  state: PropTypes.object.isRequired,
  sx: PropTypes.object,
};

export default Viewer;
