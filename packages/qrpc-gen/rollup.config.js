// import resolve from '@rollup/plugin-node-resolve';
import typescript from '@rollup/plugin-typescript'
import copy from 'rollup-plugin-copy'

export default [
  {
    input: "src/main.ts",
    output: [
      // {
      //   file: "dist/main.cjs",
      //   format: "cjs"
      // },
      {
        file: "dist/main.js",
        format: "es"
      }
    ],
    plugins: [
        typescript(),
        copy({
            targets: [
                { src: 'src/main.wasm', dest: 'dist' }
            ]
        })
    ]
  }
]
