import fs from 'fs';
import { type FileDesc, loadAndInitGo } from './wasm-interface.js'
const { analyzeSourceFiles } = await loadAndInitGo();

function readFile(name: string): FileDesc {
    return {
        filename: name,
        contents: fs.readFileSync(name, { encoding: 'utf8' })
    };
}

const res = analyzeSourceFiles(
    process.argv.slice(2).map(readFile),
)
console.dir(res, { depth: 7 })

export function genSchemasFromFilesWithConfig() {
    console.log(process.argv)
}


