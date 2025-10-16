import * as $s from './schemas.js'

const User = $s.Record(
    "User",
    {
        uuid: $s.UUID,
        username: $s.String,
        verified: $s.Boolean,
    }
);

type User = $s.TypeOf<typeof User>;

function foo<const T extends readonly string[]>(...values: T): T[number] {
  // You can return anything here if you're only using it for typing,
  // but returning a value enforces correct runtime usage.
  return values[0];
}



