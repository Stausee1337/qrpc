import * as $s from './schemas.js'

const CommonGender = $s.Enum(
    "CommonGender",
    "MALE",
    "FEMALE",
    "NONBINARY"
);

const User = $s.Record(
    "User",
    {
        uuid: $s.UUID,
        username: $s.String,
        verified: $s.Boolean,
        gender: $s.Union(CommonGender, $s.Number)
    }
);
type User = $s.TypeOf<typeof User>;

const Division = $s.Record(
    "Division",
    {
        name: $s.String,
        description: $s.String,
        previousScores: $s.Array($s.UInt),
    }
);
type Division = $s.TypeOf<typeof Division>;

const Script = $s.Record(
    "Script",
    {
        uuid: $s.UUID,
        name: $s.String,
        divisions: $s.Array(Division)
    }
);
type Script = $s.TypeOf<typeof Script>;

function x(user: User) {
}

