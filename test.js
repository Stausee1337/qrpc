import fs from 'fs';
import './wasm_exec.js'

const go = new Go();

const buffer = fs.readFileSync('./main.wasm')
const array = new Uint8Array(buffer)
const result = await WebAssembly.instantiate(array, go.importObject)
go.run(result.instance)

const encoding = fs.readFileSync('./test.rpc', { encoding: 'utf8' })
const x = test(encoding, 'test.rpc')
console.dir(x, { depth: 11 })


