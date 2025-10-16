import type { QSchema, TypeOf } from "./schemas.js";

export function Array<T extends QSchema<any>>(
    type: T
): QSchema<TypeOf<T>[]>  {
    function parseOrNull(value: unknown): TypeOf<T>[]|null {
        if (!globalThis.Array.isArray(value))
            return null;
        const result: TypeOf<T>[] = [];
        for (const item of value) {
            result.push(type.parse(item));
        }
        return result;
    }

    const schema: QSchema<TypeOf<T>[]> = {
        isSchema(value): value is TypeOf<T>[] {
            return (
                globalThis.Array.isArray(value) &&
                value.every(v => type.isSchema(v))
            );
        },
        parse(value) {
            const array = parseOrNull(value);
            if (array === null)
                throw `expected array, found ${typeof value}`
            return array;
        },
    };

    return Object.freeze(schema);
}

export function Optional<T extends QSchema<any>>(
    type: T
): QSchema<TypeOf<T>|undefined>  {
    const schema: QSchema<TypeOf<T>|undefined> = {
        isSchema(value): value is TypeOf<T>|undefined {
            if (value === undefined)
                return true;
            return type.isSchema(value);
        },
        parse(value) {
            if (value === undefined || value === null)
                return undefined;
            return type.parse(type);
        },
    };

    return Object.freeze(schema);
}

export function Union<T extends QSchema<any>[]>(
    ...types: T
): QSchema<TypeOf<T[number]>> {
    function parseOrNull(value: unknown): TypeOf<T[number]>|null {
        if (typeof value !== "object")
            return null;
        if (value === null)
            return null;

        for (let i = 0; i < types.length; i++) {
            const schema = types[i];
            const desc = Object.getOwnPropertyDescriptor(value, `variant${i + 1}`)
            if (desc !== undefined)
                return schema.parse(desc.value)
        }
        return null;
    }

    const schema: QSchema<TypeOf<T[number]>> = {
        isSchema(value): value is TypeOf<T[number]> {
            for (const type of types) {
                if (type.isSchema(value))
                    return true;
            } 
            return false;
        },
        parse(value) {
            const parsed = parseOrNull(value);
            if (parsed !== null)
                return parsed;
            throw `unexpected ${typeof value}`
        },
    };

    return Object.freeze(schema);
}

