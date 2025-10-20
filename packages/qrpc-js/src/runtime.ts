import type { QSchema, TypeOf } from "./schemas.js";
import { Record as createRecord } from './user-defined.js';

type MapRecord<T extends Record<string, QSchema<any>>> = {
    [P in keyof T]: TypeOf<T[P]>;
};

interface Operation<I extends Record<string, QSchema<any>>, R extends QSchema<any>> {
    readonly kind: "mutation"|"query";
    (input: MapRecord<I>): Promise<TypeOf<R>>;
}

interface PureOperation<R extends QSchema<any>> {
    readonly kind: "query";
    (): Promise<TypeOf<R>>;
}

interface VoidOperation<I extends Record<string, QSchema<any>>> {
    readonly kind: "mutation";
    (input: MapRecord<I>): Promise<void>;
}

export type Executor = (url: string, body: any) => Promise<Response>;
type Constructor<T> = new(executor: Executor) => T;

function overlayPrototype<T, P>(f: Constructor<T>, p: P): Constructor<T & P> {
    const result: Constructor<T & P> = f as any;
    Object.assign(result.prototype, p)
    return result;
}

export function service<O extends Record<string, Operation<any, any>>>(
    name: string,
    operations: O
): Constructor<O> {
    const $service: Constructor<$Service> = class extends $Service {
        constructor(executor: Executor) {
            super(name, executor);
        }
    };

    return overlayPrototype($service, operations)
}

function buildOperation<I extends Record<string, QSchema<any>>, R extends QSchema<any>>(
    kind: "mutation",
    name: string,
    inputTypes: I,
): VoidOperation<I>
function buildOperation<I extends Record<string, QSchema<any>>, R extends QSchema<any>>(
    kind: "mutation"|"query",
    name: string,
    inputTypes: I,
    outputSchema: R
): Operation<I, R>
function buildOperation<I extends Record<string, QSchema<any>>, R extends QSchema<any>>(
    kind: "mutation"|"query",
    name: string,
    inputTypes: I,
    outputSchema?: R
): Operation<I, R>|VoidOperation<I> {
    const inputSchema = createRecord(`${name}.Inputs`, inputTypes);
    async function executeOp(this: $Service|unknown, input: MapRecord<I>): Promise<TypeOf<R>> {
        if (!(this instanceof $Service))
            throw `${kind} ${name} was not executed on a service`
        const rawInput = inputSchema.serialize(input);
        const rawOutput = await this.$executeOperation(kind, name, rawInput)
        if (outputSchema === undefined)
            return;
        return outputSchema.parse(rawOutput)
    }
    executeOp.kind = kind;

    return executeOp;
}

export function query<R extends QSchema<any>>(
    name: string,
    outputSchema: R
): PureOperation<R>
export function query<I extends Record<string, QSchema<any>>, R extends QSchema<any>>(
    name: string,
    inputTypes: I,
    outputSchema: R
): Operation<I, R>
export function query<I extends Record<string, QSchema<any>>, R extends QSchema<any>>(
    name: string,
    secondArgument: I|R,
    outputSchema?: R
): Operation<I, R>|PureOperation<R> {
    if (outputSchema !== undefined)
        return buildOperation("query", name, secondArgument as I, outputSchema)
    const op = buildOperation("query", name, {}, secondArgument as R)
    function wrapOp(this: $Service|unknown, input: MapRecord<I>|undefined): Promise<TypeOf<R>> {
        return op.call(this, input ?? {});
    }
    wrapOp.kind = "query" as const;
    return wrapOp;
}

export function mutation<I extends Record<string, QSchema<any>>>(
    name: string,
    inputTypes: I,
): VoidOperation<I>;
export function mutation<I extends Record<string, QSchema<any>>, R extends QSchema<any>>(
    name: string,
    inputTypes: I,
    outputSchema: R
): Operation<I, R>
export function mutation<I extends Record<string, QSchema<any>>, R extends QSchema<any>>(
    name: string,
    inputTypes: I,
    outputSchema?: R
): Operation<I, R>|VoidOperation<I> {
    if (outputSchema === undefined)
        return buildOperation("mutation", name, inputTypes)
    return buildOperation("mutation", name, inputTypes, outputSchema)
}

abstract class $Service {
    public $executor: Executor;
    constructor(public $name: string, executor: Executor) {
        this.$executor = executor;
    }

    async $executeOperation(
        kind: "mutation"|"query",
        name: string,
        input: any
    ): Promise<unknown> {
        const response = await this.$executor(`/${kind}:${this.$name}.${name}`, input);
        const [data, error] = await this.#parseResponse(response);
        if (error !== undefined)
            throw `Unexpected API Response: ${error}`
        return data;
    }

    async #parseResponse(response: Response): Promise<[unknown, undefined]|[undefined, string]> {
        if (response.status !== 200)
            return [undefined, `${response.status}: ${response.statusText}`]
        const rawJSON = await response.json();
        if (typeof rawJSON !== "object") 
            return [undefined, `unexpected ${typeof rawJSON}`]
        if (rawJSON === null)
            return [undefined, `unexpected null`]
        const rawData = Object.getOwnPropertyDescriptor(rawJSON, 'data');
        if (rawData !== undefined)
            return [rawData.value, undefined];
        const rawError = Object.getOwnPropertyDescriptor(rawJSON, 'error');
        if (rawError !== undefined)
            return [undefined, String(rawError.value)];
        return [undefined, undefined]
    }
}

export function createSimpleExecutor(baseURL: string): Executor {
    return async (endpoint: string, body: any): Promise<Response> => {
        return await fetch(`${baseURL}${endpoint}`, { method: "POST", body })
    }
}


