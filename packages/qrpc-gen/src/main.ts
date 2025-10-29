#!/usr/bin/env node
import fs from 'node:fs';
import cmd, { command, option, restPositionals } from 'cmd-ts'
import { type FileDesc, loadAndInitGo } from './wasm-interface.js'
import { SourceError } from './types.js';
import { generateCode } from './codegen.js'
import esMain from './es-main.js';

const { analyzeSourceFiles } = await loadAndInitGo();

function readFileToDesc(name: string): FileDesc {
    return {
        filename: name,
        contents: fs.readFileSync(name, { encoding: 'utf8' })
    };
}

export interface Config {
    inputFiles: string[];
    outputFile: string
}

export async function genSchemasFromFilesWithConfig(config: Config): Promise<number> {
    if (config.inputFiles.length === 0) {
        try {
            fs.writeFileSync(config.outputFile, '');
        } catch (e) {
            console.error(`${e}`)
            return 1;
        }
    }

    const filesWithContent: FileDesc[] = [];
    for (const filename of config.inputFiles) {
        try {
            filesWithContent.push(readFileToDesc(filename));
        } catch (e) {
            console.error(`${e}`)
            return 1;
        }
    }

    const res = analyzeSourceFiles(filesWithContent)
    if (res instanceof SourceError) {
        res.renderToConsole()
        return 1;
    }

    const code = await generateCode(res)
    try {
        fs.writeFileSync(config.outputFile, code);
    } catch (e) {
        console.error(`${e}`)
        return 1;
    }

    return 0;
}

const entrypoint = command({
    name: 'qrpc-gen',
    description: 'Generate .ts Schemas from .qrpc Files',
    version: '0.1.3',
    args: {
        inputFiles: restPositionals({ type: cmd.string, displayName: 'definitions' }),
        outputFile: option({ long: 'output', short: 'o', type: cmd.string }),
    },
    async handler(args) {
        const code = await genSchemasFromFilesWithConfig(args);
        process.exit(code);
    },
});

if (esMain(import.meta)) {
    cmd.run(entrypoint, process.argv.slice(2))
}

