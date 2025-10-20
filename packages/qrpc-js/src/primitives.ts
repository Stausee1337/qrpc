import type { QSchema } from "./schemas.js";

export interface QPrimitive<T> extends QSchema<T> {

}

type PrimitiveDescriptor<T> = {
    isPrimitive(value: unknown): value is T;
    parse(value: unknown): T;
};

function createQPrimitive<T>(
    descriptor: PrimitiveDescriptor<T>
): QPrimitive<T> {
    const schema: QPrimitive<T> = {
        isSchema(value) {
            return descriptor.isPrimitive(value);
        },
        parse(value) {
            return descriptor.parse(value);
        },
        serialize(value) {
            return value
        },
    };
    return Object.freeze(schema);
}

export const String = createQPrimitive<string>(
    {
        isPrimitive(value) {
            return typeof value === "string";
        },
        parse(value) {
            if (typeof value === "string")
                return value;
            throw `expected string, found ${typeof value}`
        },
    }
)

declare const brand: unique symbol;
export type UUID = string & { [brand]: UUID };

const UUID_REGEX =
  /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i;

export const UUID = createQPrimitive<UUID>(
    {
        isPrimitive(value): value is UUID {
            if (typeof value !== "string")
                return false;
            return UUID_REGEX.test(value);
        },
        parse(value) {
            if (this.isPrimitive(value))
                return value;
            throw `expected UUID, found ${typeof value}`   
        },
    }
)

export const Boolean = createQPrimitive<boolean>(
    {
        isPrimitive(value) {
            return typeof value === "boolean";
        },
        parse(value) {
            if (typeof value === "boolean")
                return value;
            throw `expected boolean, found ${typeof value}`
        },
    }
)

export const UInt = createQPrimitive<number>(
    {
        isPrimitive(value) {
            return typeof value === "number";
        },
        parse(value) {
            if (typeof value !== "number")
                throw `expected UInt, found ${typeof value}`
            else if (value < 0)
                throw `expected UInt, found ${value} (< 0)`
            else if (value > globalThis.Number.MAX_SAFE_INTEGER)
                throw `expected UInt, found ${value} (> MAX_INT)`
            else if (!globalThis.Number.isInteger(value))
                throw `expected UInt, found ${value} (float)`
            return value;
        },
    }
)

export const Int = createQPrimitive<number>(
    {
        isPrimitive(value) {
            return typeof value === "number";
        },
        parse(value) {
            if (typeof value !== "number")
                throw `expected Int, found ${typeof value}`
            else if (value < globalThis.Number.MIN_SAFE_INTEGER)
                throw `expected Int, found ${value} (< MIN_INT)`
            else if (value > globalThis.Number.MAX_SAFE_INTEGER)
                throw `expected Int, found ${value} (> MAX_INT)`
            else if (!globalThis.Number.isInteger(value))
                throw `expected Int, found ${value} (float)`
            return value;
        },
    }
)

export const Number = createQPrimitive<number>(
    {
        isPrimitive(value) {
            return typeof value === "number";
        },
        parse(value) {
            if (typeof value !== "number")
                throw `expected Number, found ${typeof value}`
            return value;
        },
    }
)
