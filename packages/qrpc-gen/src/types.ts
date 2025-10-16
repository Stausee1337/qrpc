
export interface IRType {
    formatSchema(): string;
}
export interface UserDefinedType extends IRType {
    name: string
}

// Types
export class UUIDType implements IRType {
    formatSchema(): string {
        return '$s.UUID'
    }
}

export type NumberKind = "signed"|"unsigned"|"floating";
export class NumberType implements IRType {
    kind: NumberKind = 'floating';
    formatSchema(): string {  
        switch(this.kind) {
            case "signed":
                return '$s.Int'
            case "unsigned":
                return '$s.Uint'
            case "floating":
                return '$s.Number'
        }
    }
}

export class StringType implements IRType {
    formatSchema(): string {
        return '$s.String'
    }
}

export class ArrayType implements IRType {
    type: IRType = new EmptyType();

    formatSchema(): string {
        return `$s.Array(${this.type.formatSchema()})`
    }
}
export class OptionalType implements IRType {
    type: IRType = new EmptyType();

    formatSchema(): string {
        return `$s.Optional(${this.type.formatSchema()})`  
    }
}
export class UnionType implements IRType {
    types: IRType[] = [];

    formatSchema(): string {
        return `$s.Union(${this.types.map(t => t.formatSchema()).join(', ')})`  
        
    }
}
export class EmptyType implements IRType {
    formatSchema(): string {
        throw 'empty type cannot be formatted into schema'
    }
}

export class EnumType implements UserDefinedType {
    name: string = ''

    formatSchema(): string {
        return this.name;
    }
}
export class RecordType implements UserDefinedType {
    name: string = ''

    formatSchema(): string {
        return this.name;
    }
}

export class Service {
    name: string = '';
    operations: Operation[] = []
}


function isNumberLike(val: unknown): val is Number {
  return val !== '' && !isNaN(val);
}

type OperationKind = "mutation"|"query"
export class Operation {
    name: string = ''
    inputTypes: NamedType[] = []
    resultType: IRType = new EmptyType()

    #kind: OperationKind = 'query'

    get kind(): string {
        return this.#kind;
    }

    set kind(kind: OperationKind|Number) {
        if (isNumberLike(kind))
            this.#kind = (["query", "mutation"] as const)[kind]
        else
            this.#kind = kind;
    }
}

export class NamedType {
    name: string = ''
    type: IRType = new EmptyType()
}

export class AnalysisResult {
    types: UserDefinedType[] = []
    services: Service[] = []
}

