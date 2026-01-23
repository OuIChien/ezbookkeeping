import {
    type HiddenAmount,
    type NumberFormatOptions,
    NumeralSystem,
    DecimalSeparator,
    DigitGroupingSymbol
} from '@/core/numeral.ts';

import { DEFAULT_DECIMAL_NUMBER_COUNT, MAX_SUPPORTED_DECIMAL_NUMBER_COUNT, DISPLAY_HIDDEN_AMOUNT } from '@/consts/numeral.ts';

import { isDefined, isString, isNumber, replaceAll, removeAll } from './common.ts';

export function sumAmounts(amounts: number[]): number {
    let sum = 0;

    for (const amount of amounts) {
        sum += amount;
    }

    return sum;
}

export function appendDigitGroupingSymbolAndDecimalSeparator(textualNumber: string, options: NumberFormatOptions): string {
    if (!textualNumber) {
        return textualNumber;
    }

    const numeralSystem = options.numeralSystem || NumeralSystem.Default;
    const digitGroupingType = options.digitGrouping;
    const digitGroupingSymbol = options.digitGroupingSymbol || DigitGroupingSymbol.Default.symbol;
    const decimalSeparator = options.decimalSeparator || DecimalSeparator.Default.symbol;

    const negative = textualNumber.charAt(0) === '-';

    if (negative) {
        textualNumber = textualNumber.substring(1);
    }

    const integerChars: string[] = [];
    const decimalChars: string[] = [];
    let currentDecimalSeparator = '';

    let decimalNumberCount = options.decimalNumberCount;

    if (!isNumber(decimalNumberCount) || decimalNumberCount > MAX_SUPPORTED_DECIMAL_NUMBER_COUNT) {
        decimalNumberCount = DEFAULT_DECIMAL_NUMBER_COUNT;
    }

    if (textualNumber === DISPLAY_HIDDEN_AMOUNT) {
        if (decimalNumberCount > 0) {
            for (let i = 0; i < textualNumber.length; i++) {
                integerChars.push(textualNumber.charAt(i));
            }

            for (let i = 0; i < decimalNumberCount; i++) {
                decimalChars.push(textualNumber.charAt(0));
            }
        } else {
            for (let i = 0; i < textualNumber.length; i++) {
                integerChars.push(textualNumber.charAt(i));
            }
        }
    } else {
        for (let i = 0; i < textualNumber.length; i++) {
            const ch = textualNumber.charAt(i);

            if (!currentDecimalSeparator) {
                if (numeralSystem.isDigit(ch)) {
                    integerChars.push(ch);
                } else {
                    currentDecimalSeparator = ch;
                }
            } else {
                if (numeralSystem.isDigit(ch)) {
                    decimalChars.push(ch);
                } else {
                    throw new Error('Number \"' + textualNumber + '\" is not a valid textual number');
                }
            }
        }
    }

    let integer = '';

    if (digitGroupingType) {
        integer = digitGroupingType.format(integerChars, digitGroupingSymbol);
    } else {
        integer = integerChars.join('');
    }

    const decimals = decimalChars.join('');

    if (decimals) {
        textualNumber = `${integer}${decimalSeparator}${decimals}`;
    } else {
        textualNumber = integer;
    }

    if (negative) {
        textualNumber = `-${textualNumber}`;
    }

    return textualNumber;
}

export function parseAmount(str: string, options: NumberFormatOptions): number {
    if (!isString(str)) {
        return 0;
    }

    if (!str || str.length < 1) {
        return 0;
    }

    const negative = str.charAt(0) === '-';

    if (negative) {
        str = str.substring(1);
    }

    if (!str || str.length < 1) {
        return 0;
    }

    const sign = negative ? -1 : 1;

    const numeralSystem = options.numeralSystem || NumeralSystem.Default;
    const decimalSeparator = options.decimalSeparator || DecimalSeparator.Default.symbol;
    const digitGroupingSymbol = options.digitGroupingSymbol || DigitGroupingSymbol.Default.symbol;

    let decimalNumberCount = options.decimalNumberCount;

    if (!isNumber(decimalNumberCount) || decimalNumberCount > MAX_SUPPORTED_DECIMAL_NUMBER_COUNT) {
        decimalNumberCount = DEFAULT_DECIMAL_NUMBER_COUNT;
    }

    const multiplier = Math.pow(10, decimalNumberCount);

    if (str.indexOf(digitGroupingSymbol) >= 0) {
        str = removeAll(str, digitGroupingSymbol);
    }

    let decimalSeparatorPos = str.indexOf(decimalSeparator);

    if (decimalSeparatorPos < 0) {
        return sign * numeralSystem.parseInt(str) * multiplier;
    } else if (decimalSeparatorPos === 0) {
        str = numeralSystem.digitZero + str;
        decimalSeparatorPos++;
    }

    const integer = str.substring(0, decimalSeparatorPos);
    let decimals = str.substring(decimalSeparatorPos + 1, str.length);

    if (decimals.length < decimalNumberCount) {
        for (let i = decimals.length; i < decimalNumberCount; i++) {
            decimals += numeralSystem.digitZero;
        }
    } else if (decimals.length > decimalNumberCount) {
        decimals = decimals.substring(0, decimalNumberCount);
    }

    if (decimals.length < 1) {
        return sign * numeralSystem.parseInt(integer) * multiplier;
    } else {
        return sign * numeralSystem.parseInt(integer) * multiplier + sign * numeralSystem.parseInt(decimals);
    }
}

export function formatAmount(value: number, options: NumberFormatOptions): string {
    if (!isNumber(value) || !Number.isFinite(value)) {
        return '';
    }

    const numeralSystem = options.numeralSystem || NumeralSystem.Default;
    let textualNumber = numeralSystem.formatNumber(value);

    if (!textualNumber) {
        return textualNumber;
    }

    const negative = textualNumber.charAt(0) === '-';

    if (negative) {
        textualNumber = textualNumber.substring(1);
    }

    const digitGroupingType = options.digitGrouping;
    const digitGroupingSymbol = options.digitGroupingSymbol || DigitGroupingSymbol.Default.symbol;
    const decimalSeparator = options.decimalSeparator || DecimalSeparator.Default.symbol;
    let decimalNumberCount = options.decimalNumberCount;

    if (!isNumber(decimalNumberCount) || decimalNumberCount > MAX_SUPPORTED_DECIMAL_NUMBER_COUNT) {
        decimalNumberCount = DEFAULT_DECIMAL_NUMBER_COUNT;
    }

    let integer = numeralSystem.digitZero;
    let decimals = '';

    if (textualNumber.length > decimalNumberCount) {
        integer = textualNumber.substring(0, textualNumber.length - decimalNumberCount);
        decimals = textualNumber.substring(textualNumber.length - decimalNumberCount);
    } else {
        decimals = textualNumber;

        for (let i = textualNumber.length; i < decimalNumberCount; i++) {
            decimals = numeralSystem.digitZero + decimals;
        }
    }

    if (options.trimTailZero) {
        let lastNonZeroIndex = -1;

        for (let i = decimals.length - 1; i >= 0; i--) {
            if (decimals.charAt(i) !== numeralSystem.digitZero) {
                lastNonZeroIndex = i;
                break;
            }
        }

        if (lastNonZeroIndex >= 0) {
            decimals = decimals.substring(0, lastNonZeroIndex + 1);
        } else {
            decimals = '';
        }
    }

    if (integer && integer.length > 1 && digitGroupingType) {
        integer = digitGroupingType.format(integer.split(''), digitGroupingSymbol);
    }

    if (decimals) {
        textualNumber = `${integer}${decimalSeparator}${decimals}`;
    } else {
        textualNumber = integer;
    }

    if (negative) {
        textualNumber = `-${textualNumber}`;
    }

    return textualNumber;
}

export function formatHiddenAmount(value: HiddenAmount, options: NumberFormatOptions): string {
    return appendDigitGroupingSymbolAndDecimalSeparator(value, options);
}

export function formatNumber(value: number, options: NumberFormatOptions, precision?: number): string {
    const numeralSystem = options.numeralSystem || NumeralSystem.Default;

    if (isDefined(precision)) {
        const ratio = Math.pow(10, precision);
        const normalizedValue = Math.trunc(value * ratio);
        const textualValue = numeralSystem.formatNumber(normalizedValue / ratio);
        return appendDigitGroupingSymbolAndDecimalSeparator(textualValue, options);
    } else {
        const textualValue = numeralSystem.formatNumber(value);
        return appendDigitGroupingSymbolAndDecimalSeparator(textualValue, options);
    }
}

export function formatPercent(value: number, precision: number, lowPrecisionValue: string, options: NumberFormatOptions): string {
    const numeralSystem = options.numeralSystem || NumeralSystem.Default;
    const ratio = Math.pow(10, precision);
    const normalizedValue = Math.trunc(value * ratio);

    if (value > 0 && normalizedValue < 1 && lowPrecisionValue) {
        const systemDecimalSeparator = DecimalSeparator.Dot.symbol;
        const decimalSeparator = options.decimalSeparator || DecimalSeparator.Default.symbol;

        lowPrecisionValue = numeralSystem.replaceWesternArabicDigitsToLocalizedDigits(lowPrecisionValue);

        if (systemDecimalSeparator === decimalSeparator) {
            return lowPrecisionValue + '%';
        }

        return replaceAll(lowPrecisionValue, systemDecimalSeparator, decimalSeparator) + '%';
    }

    return formatNumber(value, options, precision) + '%';
}

export function getAmountWithDecimalNumberCount(amount: number, decimalNumberCount: number): number {
    if (decimalNumberCount === 0) {
        return Math.trunc(amount / 100) * 100;
    } else if (decimalNumberCount === 1) {
        return Math.trunc(amount / 10) * 10;
    }
    return amount;
}

export function formatExchangeRateAmount(exchangeRateAmount: number, options: NumberFormatOptions): string {
    const numeralSystem = options.numeralSystem || NumeralSystem.Default;
    const rateStr = numeralSystem.formatNumber(exchangeRateAmount);
    const decimalSeparator = DecimalSeparator.Dot.symbol;

    if (rateStr.indexOf(decimalSeparator) < 0) {
        return appendDigitGroupingSymbolAndDecimalSeparator(rateStr, options);
    } else {
        let firstNonZeroPos = 0;

        for (let i = 0; i < rateStr.length; i++) {
            if (rateStr.charAt(i) !== decimalSeparator && rateStr.charAt(i) !== numeralSystem.digitZero) {
                firstNonZeroPos = Math.min(i + 4, rateStr.length);
                break;
            }
        }

        const trimmedRateStr = rateStr.substring(0, Math.max(6, Math.max(firstNonZeroPos, rateStr.indexOf(decimalSeparator) + 2)));
        return appendDigitGroupingSymbolAndDecimalSeparator(trimmedRateStr, options);
    }
}

export function getAdaptiveDisplayAmountRate(amount1: number, amount2: number, options: NumberFormatOptions, fromExchangeRate?: { rate: string }, toExchangeRate?: { rate: string }): string | null {
    const numeralSystem = options.numeralSystem || NumeralSystem.Default;

    if (!amount1 || !amount2 || amount1 === amount2) {
        if (!fromExchangeRate || !fromExchangeRate.rate || !toExchangeRate || !toExchangeRate.rate) {
            return null;
        }

        amount1 = parseFloat(fromExchangeRate.rate);
        amount2 = parseFloat(toExchangeRate.rate);
    }

    if (amount1 > amount2) {
        const rate = amount1 / amount2;
        const displayRateStr = formatExchangeRateAmount(rate, options);
        return `${displayRateStr} : ${numeralSystem.getLocalizedDigit(1)}`;
    } else {
        const rate = amount2 / amount1;
        const displayRateStr = formatExchangeRateAmount(rate, options);
        return `${numeralSystem.getLocalizedDigit(1)} : ${displayRateStr}`;
    }
}

export function getExchangedAmountByRate(amount: number, fromRate: string, toRate: string): number | null {
    const exchangeRate = parseFloat(toRate) / parseFloat(fromRate);

    if (!isNumber(exchangeRate)) {
        return null;
    }

    return amount * exchangeRate;
}
