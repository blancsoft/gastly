import * as React from 'react';

import {
    Box,
    Divider,
    Typography,
    createTheme
} from '@mui/material';

import Navbar from '@/components/Navbar.tsx';
import Editor from "@/components/Editor.tsx";
import Viewer from "@/components/Viewer.tsx";

import '@fontsource/roboto/300.css';
import '@fontsource/roboto/400.css';
import '@fontsource/roboto/500.css';
import '@fontsource/roboto/700.css';

import "@/assets/wasm_exec.js"
import wasmURL from "@/assets/gastly.wasm.br?url"

import type { IGo } from "./types"


export const theme = createTheme();


const loadWasm = async (): Promise<IGo> => {

    const go = new window.Go(); // Defined in wasm_exec.js
    let instance: WebAssembly.Instance;

    if ('instantiateStreaming' in WebAssembly) {
        const wasmSrc = await WebAssembly.instantiateStreaming(await fetch(wasmURL), go.importObject);
        instance = wasmSrc.instance;
    } else {
        const resp: Response = await fetch(wasmURL);
        const buffer: ArrayBuffer = await resp.arrayBuffer();
        const wasmModule: WebAssembly.Module = await WebAssembly.compile(buffer);
        instance = await WebAssembly.instantiate(wasmModule, go.importObject);
    }
    go.run(instance);

    return go
}

const useWasm = () => {
    const [go, setGo] = React.useState<IGo | null>(null);

    React.useEffect(() => {
        loadWasm().then(go => setGo(go));
    }, [setGo]);

    return go
};


export type State = {
    showSearch: boolean,
    showAST: boolean,
    data: { [key: string]: string },
    activeTabIndex: number,
}

function App() {
    const go = useWasm();

    let fileData: { [p: string]: string } = {};
    try {
        fileData = JSON.parse(localStorage.getItem("fileData") ?? "{}")
    } finally {
        if (Object.keys(fileData).length == 0) {
            fileData = {
                'main.go': (
                    "// You can edit this code!\n" +
                    "// Click here and start typing.\n" +
                    "package main\n\n" +
                    'import "fmt"\n\n' +
                    "func main() {\n" +
                    '\tfmt.Println("Hello, 世界")\n' +
                    "}\n"
                ),
            }
        }
    }
    const [state, setState] = React.useState<State>({
        showSearch: false,
        showAST: false,
        data: fileData,
        activeTabIndex: 0,
    })

    const app = (
        <>
            <Navbar showSearch={false} /> {/* TODO: enable search after fixing */}
            <Box sx={{ width: "100%", height: "100%", overflow: "hidden", display: "flex" }}>
                <Editor state={state} setState={setState} />
                <Divider orientation="vertical" variant="middle" flexItem />
                <Viewer state={state} sx={{ display: { xs: 'none', md: 'block' } }} />
            </Box>
        </>
    )
    const loader = (
        <Typography>Loading WebAssembly module...</Typography>
    )

    return (
        <Box>
            {go ? app : loader}
        </Box>
    )
}

export default App;
