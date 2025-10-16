import type { QSchema, TypeOf } from "./schemas.js";

const $type: unique symbol = Symbol('$type')

type MapRecord<T extends Record<string, QSchema<any>>> = {
    [P in keyof T]: TypeOf<T[P]>;
};

export function Record<T extends Record<string, QSchema<any>>>(
    name: string,
    descriptor: T
): QSchema<MapRecord<T>> {
    const object = new Object();

    function parseOrNull(rawObject: unknown): MapRecord<T>|null {
        if (typeof rawObject !== "object")
            return null;
        if (rawObject === null)
            return null;

        const object: Record<string, any> = {};
        for (const [key, schema] of Object.entries(descriptor)) {
            const rawValue = (rawObject as Record<string, unknown>)[key];
            const value = schema.parse(rawValue);
            object[key] = value;
        }

        Object.defineProperty(
            object, $type,
            { enumerable: false, writable: false, value: object }
        )

        return object as MapRecord<T>;
    }

    const schema: QSchema<MapRecord<T>> = {
        isSchema(value): value is MapRecord<T> {
            if (typeof value !== "object")
                return false;
            if (value === null)
                return false;
            return (value as any)[$type] === object;
        },
        parse(value) {
            const parsedObject = parseOrNull(value) 
            if (parsedObject === null)
                throw `expected record ${name}, found ${typeof value}`
            return parsedObject
        },
        serialize(value) {
            const object: Record<string, any> = {};
            for (const [key, schema] of Object.entries(descriptor)) {
                object[key] = schema.serialize(value[key]);
            }
            return object;
        },
    };

    return Object.freeze(schema);
}

export function Enum<const T extends readonly string[]>(
    name: string,
    ...values: T
): QSchema<T[number]> {
    const variants = new Set<string>(values);
    const schema: QSchema<T[number]> = {
        isSchema(value): value is T[number] {
            if (typeof value !== "string")
                return false;
            return variants.has(value);
        },
        parse(value) {
            if (this.isSchema(value))
                return value;
            throw `expected enum ${name}, found ${typeof value}`
        },
        serialize(value) {
            return value;
        },
    };

    return Object.freeze(schema);
}

