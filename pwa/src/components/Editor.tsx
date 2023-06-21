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
import type { ReactCodeMirrorRef } from '@uiw/react-codemirror';

const goExt = ".go";
const editorExt = [StreamLanguage.define(go)];

type EditorProp = {
  state: State,
  setState: React.Dispatch<React.SetStateAction<State>>
}


const Editor = ({ state, setState }: EditorProp) => {
  const { data, activeTabIndex } = state;
  const keys = Object.keys(data).sort();
  const editorRef = React.useRef<ReactCodeMirrorRef>(null);

  const saveEditorContent = () => {
    setState(state => {
      const content = editorRef.current?.view?.state?.doc?.toString();
      const key = keys[activeTabIndex]
      let newState = { ...state, data: { ...state.data, [key]: content ?? "" } }
      localStorage.setItem("fileData", JSON.stringify(newState.data));
      return newState
    });
  }


  // Autosave every two seconds
  React.useEffect(() => {
    const saveInterval = setInterval(saveEditorContent, 2000);
    return () => {
      clearInterval(saveInterval);
    };
  }, [state]);

  const handleTabChange = (_: React.SyntheticEvent, tabIndex: number) => {
    saveEditorContent()
    setState(state => ({ ...state, activeTabIndex: tabIndex }))
  };


  const [isDialogOpen, setDialogState] = React.useState(false);
  const [fileName, setFileName] = React.useState("");
  const fileNameExt = React.useMemo(() => {
    return fileName + goExt
  }, [fileName, goExt])
  const closeDialog = () => {
    setFileName("")
    setDialogState(false)
  }
  const openDialog = () => setDialogState(true)
  const newFileDialog = (
    <Dialog open={isDialogOpen} onClose={closeDialog}>
      <DialogTitle>New Go file</DialogTitle>
      <DialogContent>
        <TextField
          value={fileName}
          error={fileNameExt in data}
          helperText={fileNameExt in data ? `${fileNameExt} already exist` : undefined}
          autoFocus={true}
          margin="none"
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
        <Button disabled={fileNameExt in data || !fileName} onClick={() => {
          setState(state => ({ ...state, data: { ...state.data, [fileNameExt]: "package main\n" } }))
          closeDialog()
        }}>Create</Button>
      </DialogActions>
    </Dialog>
  )
  const tabs = (
    <Box sx={{ display: { xs: 'flex', justifyContent: 'space-between', alignItems: 'center' }, borderBottom: 1, borderColor: 'divider' }}>
      <Tabs value={activeTabIndex} onChange={handleTabChange} sx={{ alignItems: 'center' }}
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
  )

  return (
    <Box sx={{ width: "100%", height: "100%", overflow: "hidden" }}>
      {tabs}

      {
        keys.map((item, index) => (
          <TabPanel key={index} value={activeTabIndex} index={index}>
            <CodeMirror
              autoFocus
              height={`calc(100vh - ${theme.spacing(14)})`}
              extensions={editorExt}
              value={data[item]}
              // onChange={onEditorChange}
              // onChange={setEditorContent}
              indentWithTab={true}
              ref={editorRef}
            />
          </TabPanel>
        ))
      }

      {isDialogOpen && newFileDialog}
    </Box>
  )
}

export default Editor
