import * as React from 'react';
import PropTypes from 'prop-types';
import { Tooltip, Tabs, Tab, Box, IconButton } from '@mui/material';
import StartIcon from '@mui/icons-material/Start';
import AddIcon from '@mui/icons-material/Add';
import CodeMirror from '@uiw/react-codemirror';
import { StreamLanguage } from '@codemirror/language';
import { go } from '@codemirror/legacy-modes/mode/go';

import TabPanel from './TabPanel';


const Editor = ({ state, setState }) => {
  let { data, activeTabIndex } = state;
  const handleChange = (event, tabIndex) => {
    setState({...state, activeTab: tabIndex})
  };

  const keys = Object.keys(data).sort();
  const onEditorChange = React.useCallback((update, viewUpdate) => {
    console.log('value:', update);
    const key = keys[activeTabIndex]
    setState({...state, data: {...state.data, [key]: update}})
  }, [state, activeTabIndex, keys, setState]);


  return (
    <Box sx={{ width: "100%", height: "100%", overflow: "hidden" }}>
      <Box sx={{ display: { xs: 'flex', justifyContent: 'space-between', alignItems: 'center' }, borderBottom: 1, borderColor: 'divider' }}>
        <Tabs value={activeTabIndex} onChange={handleChange} sx={{ alignItems: 'center' }}
          variant="scrollable" scrollButtons="auto"
          aria-label="editor tabs">
          {
            keys.map((item, index) => (
              <Tab key={index} label={item} {...{
                id: `editor-tab-${index}`,
                'aria-controls': `editor-tabpanel-${index}`,
              }} />
            ))
          }

        </Tabs>
        <Box sx={{ display: { xs: 'flex' }, justifyContent: 'space-between', alignItems: 'center' }}>
          <Tooltip title="New tab">
            <IconButton aria-label="Create new tab" sx={{ mr: 1 }}><AddIcon /></IconButton>
          </Tooltip>
          <Tooltip title="View AST">
            <IconButton aria-label="View source AST" sx={{ display: { md: 'none' }, mr: 1 }}><StartIcon /></IconButton>
          </Tooltip>
        </Box>
      </Box>

      {
        keys.map((item, index) => (
          <TabPanel key={index} value={activeTabIndex} index={index}>
            <CodeMirror
              height="100%"
              extensions={[StreamLanguage.define(go)]}
              value={data[item]} onChange={onEditorChange} />
          </TabPanel>
        ))
      }
    </Box>
  )
}

Editor.propTypes = {
  state: PropTypes.object.isRequired,
  setState: PropTypes.func.isRequired,
};

export default Editor
