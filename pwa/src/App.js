import * as React from 'react';
import Navbar from './components/Navbar';
import Editor from "./components/Editor";
import Viewer from "./components/Viewer";

import { Box, Divider } from '@mui/material';

import dedent from "dedent";

import '@fontsource/roboto/300.css';
import '@fontsource/roboto/400.css';
import '@fontsource/roboto/500.css';
import '@fontsource/roboto/700.css';

const WASM_URL = window.location.origin + '/gastly.wasm';


const loadWasm = async () => {
  const go = new window.Go(); // Defined in wasm_exec.js

  let wasm;
  if ('instantiateStreaming' in WebAssembly) {
    wasm = await WebAssembly.instantiateStreaming(await fetch(WASM_URL), go.importObject);
  } else {
    const resp = await fetch(WASM_URL)
    wasm = await WebAssembly.instantiate(resp.arrayBuffer(), go.importObject);
  }
  go.run(wasm.instance);
}

export const LoadWasm = (props) => {
  const [isLoading, setIsLoading] = React.useState(true);

  React.useEffect(() => {
    loadWasm().then(() => {
      setIsLoading(false);
    });
  }, []);

  if (isLoading) {
    return (
      <div>
        Loading WebAssembly module...
      </div>
    );
  } else {
    return <React.Fragment>{props.children}</React.Fragment>;
  }
};


function App() {
  const fileData = {
    'main.go': dedent`
    // You can edit this code!
    // Click here and start typing.
    package main
    
    import "fmt"
    
    func main() {
        fmt.Println("Hello, 世界")
    }
    `,
  }


  const [state, setState] = React.useState({
    showSearch: false,
    showAST: false,
    data: fileData,
    activeTabIndex: 0,
  })

  return (
    <LoadWasm>
      <Box>
        <Navbar showSearch={false} /> {/* TODO: enable search after fixing */}
        <Box sx={{ width: "100%", height: "100%", overflow: "hidden", display: "flex" }}>
          <Editor state={state} setState={setState} />
          <Divider orientation="vertical" variant="middle" flexItem />
          <Viewer state={state} sx={{ display: { xs: 'none', md: 'block' } }} />
        </Box>
      </Box>
    </LoadWasm>
  );
}

export default App;
