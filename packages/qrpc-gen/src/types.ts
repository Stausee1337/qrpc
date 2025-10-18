
export interface IRType {
    formatSchema(): string;
    readonly children: IRType[]
}
export interface UserDefinedType extends IRType {
    name: string
}

class SingletonMixin {
    get children(): IRType[] {
        return [];
    }
}

// Types
export class UUIDType extends SingletonMixin implements IRType {
    formatSchema(): string {
        return '$s.UUID'
    }
}

export type NumberKind = "signed"|"unsigned"|"floating";
export class NumberType extends SingletonMixin implements IRType {
    kind: NumberKind = 'floating';
    formatSchema(): string {  
        switch(this.kind) {
            case "signed":
                return '$s.Int'
            case "unsigned":
                return '$s.UInt'
            case "floating":
                return '$s.Number'
        }
    }
}

export class BoolType extends SingletonMixin implements IRType {
    formatSchema(): string {
        return '$s.Boolean'
    }
}

export class StringType extends SingletonMixin implements IRType {
    formatSchema(): string {
        return '$s.String'
    }
}

export class ArrayType implements IRType {
    type: IRType = new EmptyType();

    get children(): IRType[] {
        return [this.type];
    }

    formatSchema(): string {
        return `$s.Array(${this.type.formatSchema()})`
    }
}
export class OptionalType implements IRType {
    type: IRType = new EmptyType();

    get children(): IRType[] {
        return [this.type];
    }

    formatSchema(): string {
        return `$s.Optional(${this.type.formatSchema()})`  
    }
}
export class UnionType implements IRType {
    types: IRType[] = [];

    get children(): IRType[] {
        return this.types;
    }

    formatSchema(): string {
        return `$s.Union(${this.types.map(t => t.formatSchema()).join(', ')})`  
        
    }
}
export class EmptyType extends SingletonMixin implements IRType {
    formatSchema(): string {
        throw 'empty type cannot be formatted into schema'
    }
}

export class EnumType implements UserDefinedType {
    name: string = ''
    variants: string[] = []

    get children(): IRType[] {
        return [];
    }

    formatSchema(): string {
        return this.name;
    }
}
export class RecordType implements UserDefinedType {
    name: string = ''
    fields: NamedType[] = []

    get children(): IRType[] {
        return [];
    }

    formatSchema(): string {
        return this.name;
    }
}

export class Service {
    name: string = '';
    operations: Operation[] = []
}

function isNumberLike(val: unknown): val is Number {
    return val !== '' && !isNaN(val as number);
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
            this.#kind = (["query", "mutation"] as const)[kind as (0|1)]
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

export class Position {
    start: number = 0
    end: number = 0
    lineno: number = 1
    column: number = 1
    file: string = ''
}

export class SourceError {
    pos: Position = new Position()
    message: string = ''

    renderToConsole() {
        console.error(
            `ERROR: ${this.pos.file}:${this.pos.lineno}:${this.pos.column}: ${this.message}`
        )
    }
}


