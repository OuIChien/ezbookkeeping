// Maximum safe integer in JavaScript: 9007199254740991
// For fraction=6, this allows up to ~9,007,199,254.74 coins (about 90 billion coins)
// This is well within int64 max value: 9223372036854775807
export const TRANSACTION_MIN_AMOUNT: number = -Number.MAX_SAFE_INTEGER;
export const TRANSACTION_MAX_AMOUNT: number = Number.MAX_SAFE_INTEGER;
export const TRANSACTION_MAX_PICTURE_COUNT: number = 10;
