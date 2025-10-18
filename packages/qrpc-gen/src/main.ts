import fs from 'node:fs';
import { fileURLToPath } from 'node:url';
import { type FileDesc, loadAndInitGo } from './wasm-interface.js'
import { SourceError } from './types.js';
import { generateCode } from './codegen.js'

const { analyzeSourceFiles } = await loadAndInitGo();

function isMain(importMetaUrl: string): boolean {
    return process.argv[1] === fileURLToPath(importMetaUrl)
}


export function genSchemasFromFilesWithConfig() {
    console.log(process.argv)
}

function readFileToDesc(name: string): FileDesc {
    return {
        filename: name,
        contents: fs.readFileSync(name, { encoding: 'utf8' })
    };
}

function main() {
    const res = analyzeSourceFiles(
        process.argv.slice(2).map(readFileToDesc),
    )
    if (res instanceof SourceError) {
        res.renderToConsole()
        process.exit(1)
    }
    generateCode(res).then(console.log);
}

if (isMain(import.meta.url)) {
    main()
}

