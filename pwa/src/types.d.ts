/// <reference types="./App.tsx" />
/// <reference types="./components/Viewer.tsx" />

// <reference path="go.d.ts"/>

type Result = {
    name?: string;
    ast?: string;
    dump?: string;
    source?: { [key: string]: string };
    err?: unknown;
    errMsg?: string;
};

interface IGo {
    argv: string[];
    env: Record<string, string>;
    exit: (code: number) => void;
    importObject: {
        go: {
            "runtime.wasmExit": (sp: number) => void;
            "runtime.wasmWrite": (sp: number) => void;
            "runtime.resetMemoryDataView": (sp: number) => void;
            "runtime.nanotime1": (sp: number) => void;
            "runtime.walltime": (sp: number) => void;
            "runtime.scheduleTimeoutEvent": (sp: number) => void;
            "runtime.clearTimeoutEvent": (sp: number) => void;
            "runtime.getRandomData": (sp: number) => void;
            "syscall/js.finalizeRef": (sp: number) => void;
            "syscall/js.stringVal": (sp: number) => void;
            "syscall/js.valueGet": (sp: number) => void;
            "syscall/js.valueSet": (sp: number) => void;
            "syscall/js.valueDelete": (sp: number) => void;
            "syscall/js.valueIndex": (sp: number) => void;
            "syscall/js.valueSetIndex": (sp: number) => void;
            "syscall/js.valueCall": (sp: number) => void;
            "syscall/js.valueInvoke": (sp: number) => void;
            "syscall/js.valueNew": (sp: number) => void;
            "syscall/js.valueLength": (sp: number) => void;
            "syscall/js.valuePrepareString": (sp: number) => void;
            "syscall/js.valueLoadString": (sp: number) => void;
            "syscall/js.valueInstanceOf": (sp: number) => void;
            "syscall/js.copyBytesToGo": (sp: number) => void;
            "syscall/js.copyBytesToJS": (sp: number) => void;
            "debug": (value: unknown) => void
        }
    };
    exited: boolean;
    mem: DataView;
    run(instance: WebAssembly.Instance): Promise<void>;
}


declare global {
    export interface Window {
        Go: new () => IGo;
        Gastly: {
            FromSourceCode: (fileName: string, sourceCode: string) => Result
            FromPackages: (pkgName: string[]) => Result[]
        }
    }
}

export {IGo};
