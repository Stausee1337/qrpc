
export interface QSchema<T> {
    isSchema(value: unknown): value is T;
    parse(value: unknown): T;
}

export type TypeOf<T> = T extends QSchema<(infer U)> ? U : never;

export * from './primitives.js'
export * from './user-defined.js'

