import './wasm_exec.js'
import fs from 'node:fs/promises'
import path from 'node:path'
import * as tys from './types.js'

export type FileDesc = {
    filename: string,
    contents: string
};

export interface Exports {
    analyzeSourceFiles(files: FileDesc[]): tys.AnalysisResult|tys.SourceError;
}

export async function loadAndInitGo(): Promise<Exports> {
    const go = new globalThis.Go();

    const wasmFile = path.join(import.meta.dirname, 'main.wasm')
    const buffer = await fs.readFile(wasmFile)
    const array = new Uint8Array(buffer)
    const source = await WebAssembly.instantiate(array, go.importObject)
    go.run(source.instance)

    return globalThis.$goInit(tys)
}

declare namespace globalThis {
    interface Go {
        run(instance: WebAssembly.Instance): Promise<unknown>;
        readonly importObject: WebAssembly.Imports;
    }
    const Go: {
        new(): Go;
    };
    function $goInit(constructors: Record<string, Function>): Exports;
}

