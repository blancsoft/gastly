import * as React from 'react';
// import brotli from 'brotli';
import dedent from "dedent";

import {
    Box,
    Divider,
    SxProps,
    Typography,
    createTheme, useMediaQuery
} from '@mui/material';

import Navbar from '@/components/Navbar.tsx';
import Editor from "@/components/Editor.tsx";
import Viewer from "@/components/Viewer.tsx";

import '@fontsource/roboto/300.css';
import '@fontsource/roboto/400.css';
import '@fontsource/roboto/500.css';
import '@fontsource/roboto/700.css';

import "@/assets/wasm_exec.js"
import wasmURL from "@/assets/gastly.wasm?url"

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
    const isMobile = useMediaQuery(theme.breakpoints.down('md'))

    const [state, setState] = React.useState<State>({
        showSearch: false,
        showAST: !isMobile,
        data: fileData,
        activeTabIndex: 0,
    })

    const app = (
        <>
            <Navbar showSearch={false} /> {/* TODO: enable search after fixing */}
            <Box sx={{ width: "100%", height: "100%", overflow: "hidden", display: "flex" }}>
                <Editor state={state} setState={setState} sx={{
                    display: {
                        xs: isMobile && !state.showAST ? 'block' : 'none',
                        md: 'block'
                    }
                }} />
                <Divider orientation="vertical" variant="middle" flexItem />
                <Viewer state={state} setState={setState} sx={{
                    display: {
                        xs: isMobile && state.showAST ? 'block' : 'none',
                        md: 'block'
                    }
                }} />
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
