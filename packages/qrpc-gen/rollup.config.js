export default [
  {
    input: "src/main.ts",
    output: [
      {
        file: "dist/main.cjs",
        format: "cjs"
      },
      {
        file: "dist/main.js",
        format: "es"
      }
    ],
  }
]
