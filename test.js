import fs from 'fs';
import './constructors.js'
import './wasm_exec.js'

const go = new Go();

const buffer = fs.readFileSync('./main.wasm')
const array = new Uint8Array(buffer)
const result = await WebAssembly.instantiate(array, go.importObject)
go.run(result.instance)

function readFile(name) {
    return {
        filename: name,
        contents: fs.readFileSync(name, { encoding: 'utf8' })
    };
}

const x = test(
    process.argv.slice(2).map(readFile),
)
console.dir(x, { depth: 7 })


