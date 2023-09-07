import * as React from 'react';
import {  Box } from '@mui/material';

type TabPanelProps = {
  children: React.ReactNode,
  index: number,
  value: number,
  [key: string]: unknown;
}

const TabPanel = ({ children, value, index, ...other }: TabPanelProps) => {

  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={`editor-tabpanel-${index}`}
      aria-labelledby={`editor-tab-${index}`}
      {...other}
    >
      {value === index && (
        <Box sx={{}}>
          {children}
        </Box>
      )}
    </div>
  );
}

export default TabPanel;
