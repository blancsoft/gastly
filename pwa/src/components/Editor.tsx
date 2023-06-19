import * as React from 'react';
import StartIcon from '@mui/icons-material/Start';
import AddIcon from '@mui/icons-material/Add';
import {
  Tabs, Tab, Box, Button, IconButton, Tooltip, TextField,
  InputAdornment, DialogTitle, Dialog, DialogActions, DialogContent
} from '@mui/material';

import CodeMirror from '@uiw/react-codemirror';
import { StreamLanguage } from '@codemirror/language';
import { go } from '@codemirror/legacy-modes/mode/go';

import { theme } from '@/App';
import TabPanel from '@/components/TabPanel';

import type { State } from '@/App';

type EditorProp = {
  state: State,
  setState:  React.Dispatch<React.SetStateAction<State>>
}


const goExt = ".go"
const Editor = ({ state, setState }: EditorProp) => {
  const { data, activeTabIndex } = state;
  const handleChange = (_: React.SyntheticEvent, tabIndex: number) => {
    setState({...state, activeTabIndex: tabIndex})
  };

  const keys = Object.keys(data).sort();
  const onEditorChange = React.useCallback((update: string) => {
    console.log('value:', update);
    const key = keys[activeTabIndex]
    setState({...state, data: {...state.data, [key]: update}})
  }, [state, activeTabIndex, keys, setState]);

  const [isDialogOpen, setDialogState] = React.useState(false);
  const [fileName, setFileName] = React.useState("");
  const closeDialog = () => {
    setFileName("")
    setDialogState(false)
  }
  const openDialog = () => setDialogState(true)
  const fileExists = () => (`${fileName}${goExt}` in data)
  const newFileDialog = (
      <Dialog fullWidth open={isDialogOpen} onClose={closeDialog}>
        <DialogTitle>New Go file</DialogTitle>
        <DialogContent>
          <TextField
              value={fileName}
              error={fileExists()}
              helperText={fileExists() ?`${fileName}${goExt} already exist`:undefined}
              autoFocus
              margin="dense"
              id="fileName"
              label="File Name"
              type="text"
              fullWidth
              variant="standard"
              InputProps={{
                endAdornment: <InputAdornment position="end">{goExt}</InputAdornment>,
              }}
              onChange={(event: React.ChangeEvent<HTMLInputElement>) => {
                setFileName(event.target.value);
              }}
          />
        </DialogContent>
        <DialogActions>
          <Button disabled={fileName in data || !fileName} onClick={() => {
            setState({...state, data: {...state.data, [fileName + goExt]: "package main\n"}})
            closeDialog()
          }}>Create</Button>
        </DialogActions>
      </Dialog>
  )


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
            <IconButton aria-label="Create new tab" onClick={openDialog} sx={{ mr: 1 }}><AddIcon /></IconButton>
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
              height={`calc(100vh - ${theme.spacing(14)})`}
              extensions={[StreamLanguage.define(go)]}
              value={data[item]}
              onChange={onEditorChange}
              indentWithTab={true}
            />
          </TabPanel>
        ))
      }

      { isDialogOpen && newFileDialog }
    </Box>
  )
}

export default Editor
