import type { CurrencyInfo } from '@/core/currency.ts';
import { CurrencyType } from '@/core/currency.ts';

// ISO 4217
// Reference: https://www.six-group.com/dam/download/financial-information/data-center/iso-currrency/lists/list-one.xml
export const ALL_FIAT_CURRENCIES: Record<string, CurrencyInfo> = {
    'AED': { // UAE Dirham
        code: 'AED',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'Dh',
            plural: 'Dhs'
        },
        unit: 'Dirham'
    },
    'AFN': { // Afghani
        code: 'AFN',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'Af.',
            plural: 'Afs.'
        },
        unit: 'Afghani'
    },
    'ALL': { // Lek
        code: 'ALL',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'L'
        },
        unit: 'Lek'
    },
    'AMD': { // Armenian Dram
        code: 'AMD',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: '֏'
        },
        unit: 'Dram'
    },
    'ANG': { // Netherlands Antillean Guilder
        code: 'ANG',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'ƒ'
        },
        unit: 'Guilder'
    },
    'AOA': { // Kwanza
        code: 'AOA',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'Kz'
        },
        unit: 'Kwanza'
    },
    'ARS': { // Argentine Peso
        code: 'ARS',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: '$'
        },
        unit: 'Peso'
    },
    'AUD': { // Australian Dollar
        code: 'AUD',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: '$'
        },
        unit: 'Dollar'
    },
    'AWG': { // Aruban Florin
        code: 'AWG',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'Afl.'
        },
        unit: 'Florin'
    },
    'AZN': { // Azerbaijan Manat
        code: 'AZN',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: '₼'
        },
        unit: 'Manat'
    },
    'BAM': { // Convertible Mark
        code: 'BAM',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'KM'
        },
        unit: 'Mark'
    },
    'BBD': { // Barbados Dollar
        code: 'BBD',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: '$'
        },
        unit: 'Dollar'
    },
    'BDT': { // Taka
        code: 'BDT',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: '৳'
        },
        unit: 'Taka'
    },
    'BGN': { // Bulgarian Lev
        code: 'BGN',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'лв'
        },
        unit: 'Lev'
    },
    'BHD': { // Bahraini Dinar
        code: 'BHD',
        type: CurrencyType.Fiat,
        fraction: 3,
        symbol: {
            normal: 'BD'
        },
        unit: 'Dinar'
    },
    'BIF': { // Burundi Franc
        code: 'BIF',
        type: CurrencyType.Fiat,
        fraction: 0,
        symbol: {
            normal: 'FBu'
        },
        unit: 'Franc'
    },
    'BMD': { // Bermudian Dollar
        code: 'BMD',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: '$'
        },
        unit: 'Dollar'
    },
    'BND': { // Brunei Dollar
        code: 'BND',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: '$'
        },
        unit: 'Dollar'
    },
    'BOB': { // Boliviano
        code: 'BOB',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'Bs'
        },
        unit: 'Boliviano'
    },
    'BRL': { // Brazilian Real
        code: 'BRL',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'R$'
        },
        unit: 'Real'
    },
    'BSD': { // Bahamian Dollar
        code: 'BSD',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: '$'
        },
        unit: 'Dollar'
    },
    'BTN': { // Ngultrum
        code: 'BTN',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'Nu.'
        },
        unit: 'Ngultrum'
    },
    'BWP': { // Pula
        code: 'BWP',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'P'
        },
        unit: 'Pula'
    },
    'BYN': { // Belarusian Ruble
        code: 'BYN',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'Rbl',
            plural: 'Rbls'
        },
        unit: 'Ruble'
    },
    'BZD': { // Belize Dollar
        code: 'BZD',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: '$'
        },
        unit: 'Dollar'
    },
    'CAD': { // Canadian Dollar
        code: 'CAD',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: '$'
        },
        unit: 'Dollar'
    },
    'CDF': { // Congolese Franc
        code: 'CDF',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'FC'
        },
        unit: 'Franc'
    },
    'CHF': { // Swiss Franc
        code: 'CHF',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'CHF'
        },
        unit: 'Franc'
    },
    'CLP': { // Chilean Peso
        code: 'CLP',
        type: CurrencyType.Fiat,
        fraction: 0,
        symbol: {
            normal: '$'
        },
        unit: 'Peso'
    },
    'CNY': { // Yuan Renminbi
        code: 'CNY',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: '¥'
        },
        unit: 'Yuan'
    },
    'COP': { // Colombian Peso
        code: 'COP',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: '$'
        },
        unit: 'Peso'
    },
    'CRC': { // Costa Rican Colon
        code: 'CRC',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: '₡'
        },
        unit: 'Colon'
    },
    'CUC': { // Peso Convertible
        code: 'CUC',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: '$'
        },
        unit: 'Peso'
    },
    'CUP': { // Cuban Peso
        code: 'CUP',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: '$'
        },
        unit: 'Peso'
    },
    'CVE': { // Cabo Verde Escudo
        code: 'CVE',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: '$'
        },
        unit: 'Escudo'
    },
    'CZK': { // Czech Koruna
        code: 'CZK',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'Kč'
        },
        unit: 'Koruna'
    },
    'DJF': { // Djibouti Franc
        code: 'DJF',
        type: CurrencyType.Fiat,
        fraction: 0,
        symbol: {
            normal: 'Fdj'
        },
        unit: 'Franc'
    },
    'DKK': { // Danish Krone
        code: 'DKK',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'kr.'
        },
        unit: 'Krone'
    },
    'DOP': { // Dominican Peso
        code: 'DOP',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: '$'
        },
        unit: 'Peso'
    },
    'DZD': { // Algerian Dinar
        code: 'DZD',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'DA'
        },
        unit: 'Dinar'
    },
    'EGP': { // Egyptian Pound
        code: 'EGP',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: '£'
        },
        unit: 'Pound'
    },
    'ERN': { // Nakfa
        code: 'ERN',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'Nfk'
        },
        unit: 'Nakfa'
    },
    'ETB': { // Ethiopian Birr
        code: 'ETB',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'Br'
        },
        unit: 'Birr'
    },
    'EUR': { // Euro
        code: 'EUR',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: '€'
        },
        unit: 'Euro'
    },
    'FJD': { // Fiji Dollar
        code: 'FJD',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: '$'
        },
        unit: 'Dollar'
    },
    'FKP': { // Falkland Islands Pound
        code: 'FKP',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: '£'
        },
        unit: 'Pound'
    },
    'GBP': { // Pound Sterling
        code: 'GBP',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: '£'
        },
        unit: 'Pound'
    },
    'GEL': { // Lari
        code: 'GEL',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'ლ'
        },
        unit: 'Lari'
    },
    'GHS': { // Ghana Cedi
        code: 'GHS',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'GH₵'
        },
        unit: 'Cedi'
    },
    'GIP': { // Gibraltar Pound
        code: 'GIP',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: '£'
        },
        unit: 'Pound'
    },
    'GMD': { // Dalasi
        code: 'GMD',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'D'
        },
        unit: 'Dalasi'
    },
    'GNF': { // Guinean Franc
        code: 'GNF',
        type: CurrencyType.Fiat,
        fraction: 0,
        symbol: {
            normal: 'FG'
        },
        unit: 'Franc'
    },
    'GTQ': { // Quetzal
        code: 'GTQ',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'Q'
        },
        unit: 'Quetzal'
    },
    'GYD': { // Guyana Dollar
        code: 'GYD',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: '$'
        },
        unit: 'Dollar'
    },
    'HKD': { // Hong Kong Dollar
        code: 'HKD',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'HK$'
        },
        unit: 'Dollar'
    },
    'HNL': { // Lempira
        code: 'HNL',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'L'
        },
        unit: 'Lempira'
    },
    'HTG': { // Gourde
        code: 'HTG',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'G'
        },
        unit: 'Gourde'
    },
    'HUF': { // Forint
        code: 'HUF',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'Ft'
        },
        unit: 'Forint'
    },
    'IDR': { // Rupiah
        code: 'IDR',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'Rp'
        },
        unit: 'Rupiah'
    },
    'ILS': { // New Israeli Sheqel
        code: 'ILS',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: '₪'
        },
        unit: 'Shekel'
    },
    'INR': { // Indian Rupee
        code: 'INR',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: '₹'
        },
        unit: 'Rupee'
    },
    'IQD': { // Iraqi Dinar
        code: 'IQD',
        type: CurrencyType.Fiat,
        fraction: 3,
        symbol: {
            normal: 'ID'
        },
        unit: 'Dinar'
    },
    'IRR': { // Iranian Rial
        code: 'IRR',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'Rl',
            plural: 'Rls'
        },
        unit: 'Rial'
    },
    'ISK': { // Iceland Krona
        code: 'ISK',
        type: CurrencyType.Fiat,
        fraction: 0,
        symbol: {
            normal: 'kr'
        },
        unit: 'Krona'
    },
    'JMD': { // Jamaican Dollar
        code: 'JMD',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: '$'
        },
        unit: 'Dollar'
    },
    'JOD': { // Jordanian Dinar
        code: 'JOD',
        type: CurrencyType.Fiat,
        fraction: 3,
        symbol: {
            normal: 'د.أ'
        },
        unit: 'Dinar'
    },
    'JPY': { // Yen
        code: 'JPY',
        type: CurrencyType.Fiat,
        fraction: 0,
        symbol: {
            normal: '¥'
        },
        unit: 'Yen'
    },
    'KES': { // Kenyan Shilling
        code: 'KES',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: '/='
        },
        unit: 'Shilling'
    },
    'KGS': { // Som
        code: 'KGS',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: '⃀'
        },
        unit: 'Som'
    },
    'KHR': { // Riel
        code: 'KHR',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: '៛'
        },
        unit: 'Riel'
    },
    'KMF': { // Comorian Franc
        code: 'KMF',
        type: CurrencyType.Fiat,
        fraction: 0,
        symbol: {
            normal: 'CF'
        },
        unit: 'Franc'
    },
    'KPW': { // North Korean Won
        code: 'KPW',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: '₩'
        },
        unit: 'Won'
    },
    'KRW': { // Won
        code: 'KRW',
        type: CurrencyType.Fiat,
        fraction: 0,
        symbol: {
            normal: '₩'
        },
        unit: 'Won'
    },
    'KWD': { // Kuwaiti Dinar
        code: 'KWD',
        type: CurrencyType.Fiat,
        fraction: 3,
        symbol: {
            normal: 'KD'
        },
        unit: 'Dinar'
    },
    'KYD': { // Cayman Islands Dollar
        code: 'KYD',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: '$'
        },
        unit: 'Dollar'
    },
    'KZT': { // Tenge
        code: 'KZT',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: '₸'
        },
        unit: 'Tenge'
    },
    'LAK': { // Lao Kip
        code: 'LAK',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: '₭'
        },
        unit: 'Kip'
    },
    'LBP': { // Lebanese Pound
        code: 'LBP',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'LL'
        },
        unit: 'Pound'
    },
    'LKR': { // Sri Lanka Rupee
        code: 'LKR',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'රු'
        },
        unit: 'Rupee'
    },
    'LRD': { // Liberian Dollar
        code: 'LRD',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: '$'
        },
        unit: 'Dollar'
    },
    'LSL': { // Loti
        code: 'LSL',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'L',
            plural: 'M'
        },
        unit: 'Loti'
    },
    'LYD': { // Libyan Dinar
        code: 'LYD',
        type: CurrencyType.Fiat,
        fraction: 3,
        symbol: {
            normal: 'LD'
        },
        unit: 'Dinar'
    },
    'MAD': { // Moroccan Dirham
        code: 'MAD',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'DH'
        },
        unit: 'Dirham'
    },
    'MDL': { // Moldovan Leu
        code: 'MDL',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'L'
        },
        unit: 'Leu'
    },
    'MGA': { // Malagasy Ariary
        code: 'MGA',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'Ar'
        },
        unit: 'Ariary'
    },
    'MKD': { // Denar
        code: 'MKD',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'DEN'
        },
        unit: 'Denar'
    },
    'MMK': { // Kyat
        code: 'MMK',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'K',
            plural: 'Ks.'
        },
        unit: 'Kyat'
    },
    'MNT': { // Tugrik
        code: 'MNT',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: '₮'
        },
        unit: 'Tugrik'
    },
    'MOP': { // Pataca
        code: 'MOP',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: '$'
        },
        unit: 'Pataca'
    },
    'MRU': { // Ouguiya
        code: 'MRU',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'UM'
        },
        unit: 'Ouguiya'
    },
    'MUR': { // Mauritius Rupee
        code: 'MUR',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'Re.',
            plural: 'Rs.'
        },
        unit: 'Rupee'
    },
    'MVR': { // Rufiyaa
        code: 'MVR',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'Rf.'
        },
        unit: 'Rufiyaa'
    },
    'MWK': { // Malawi Kwacha
        code: 'MWK',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'K'
        },
        unit: 'Kwacha'
    },
    'MXN': { // Mexican Peso
        code: 'MXN',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: '$'
        },
        unit: 'Peso'
    },
    'MYR': { // Malaysian Ringgit
        code: 'MYR',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'RM'
        },
        unit: 'Ringgit'
    },
    'MZN': { // Mozambique Metical
        code: 'MZN',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'MT'
        },
        unit: 'Metical'
    },
    'NAD': { // Namibia Dollar
        code: 'NAD',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: '$'
        },
        unit: 'Dollar'
    },
    'NGN': { // Naira
        code: 'NGN',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: '₦'
        },
        unit: 'Naira'
    },
    'NIO': { // Cordoba Oro
        code: 'NIO',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'C$'
        },
        unit: 'Cordoba'
    },
    'NOK': { // Norwegian Krone
        code: 'NOK',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'kr'
        },
        unit: 'Krone'
    },
    'NPR': { // Nepalese Rupee
        code: 'NPR',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'रु'
        },
        unit: 'Rupee'
    },
    'NZD': { // New Zealand Dollar
        code: 'NZD',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: '$'
        },
        unit: 'Dollar'
    },
    'OMR': { // Rial Omani
        code: 'OMR',
        type: CurrencyType.Fiat,
        fraction: 3,
        symbol: {
            normal: 'R.O'
        },
        unit: 'Rial'
    },
    'PAB': { // Balboa
        code: 'PAB',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'B/.'
        },
        unit: 'Balboa'
    },
    'PEN': { // Sol
        code: 'PEN',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'S/'
        },
        unit: 'Sol'
    },
    'PGK': { // Kina
        code: 'PGK',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'K'
        },
        unit: 'Kina'
    },
    'PHP': { // Philippine Peso
        code: 'PHP',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: '₱'
        },
        unit: 'Peso'
    },
    'PKR': { // Pakistan Rupee
        code: 'PKR',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'Re.',
            plural: 'Rs.'
        },
        unit: 'Rupee'
    },
    'PLN': { // Zloty
        code: 'PLN',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'zł'
        },
        unit: 'Zloty'
    },
    'PYG': { // Guarani
        code: 'PYG',
        type: CurrencyType.Fiat,
        fraction: 0,
        symbol: {
            normal: '₲'
        },
        unit: 'Guarani'
    },
    'QAR': { // Qatari Rial
        code: 'QAR',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'QR'
        },
        unit: 'Rial'
    },
    'RON': { // Romanian Leu
        code: 'RON',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'L'
        },
        unit: 'Leu'
    },
    'RSD': { // Serbian Dinar
        code: 'RSD',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'din.'
        },
        unit: 'Dinar'
    },
    'RUB': { // Russian Ruble
        code: 'RUB',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: '₽'
        },
        unit: 'Ruble'
    },
    'RWF': { // Rwanda Franc
        code: 'RWF',
        type: CurrencyType.Fiat,
        fraction: 0,
        symbol: {
            normal: 'FRw'
        },
        unit: 'Franc'
    },
    'SAR': { // Saudi Riyal
        code: 'SAR',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'SAR'
        },
        unit: 'Riyal'
    },
    'SBD': { // Solomon Islands Dollar
        code: 'SBD',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: '$'
        },
        unit: 'Dollar'
    },
    'SCR': { // Seychelles Rupee
        code: 'SCR',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'Re.',
            plural: 'Rs.'
        },
        unit: 'Rupee'
    },
    'SDG': { // Sudanese Pound
        code: 'SDG',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'LS'
        },
        unit: 'Pound'
    },
    'SEK': { // Swedish Krona
        code: 'SEK',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'kr'
        },
        unit: 'Krona'
    },
    'SGD': { // Singapore Dollar
        code: 'SGD',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: '$'
        },
        unit: 'Dollar'
    },
    'SHP': { // Saint Helena Pound
        code: 'SHP',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: '£'
        },
        unit: 'Pound'
    },
    'SLE': { // Leone
        code: 'SLE',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'Le'
        },
        unit: 'Leone'
    },
    'SOS': { // Somali Shilling
        code: 'SOS',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'Sh.So.'
        },
        unit: 'Shilling'
    },
    'SRD': { // Surinam Dollar
        code: 'SRD',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: '$'
        },
        unit: 'Dollar'
    },
    'SSP': { // South Sudanese Pound
        code: 'SSP',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'SS£'
        },
        unit: 'Pound'
    },
    'STN': { // Dobra
        code: 'STN',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'Db'
        },
        unit: 'Dobra'
    },
    'SVC': { // El Salvador Colon
        code: 'SVC',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: '₡'
        },
        unit: 'Colon'
    },
    'SYP': { // Syrian Pound
        code: 'SYP',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'LS'
        },
        unit: 'Pound'
    },
    'SZL': { // Lilangeni
        code: 'SZL',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'E'
        },
        unit: 'Lilangeni'
    },
    'THB': { // Baht
        code: 'THB',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: '฿'
        },
        unit: 'Baht'
    },
    'TJS': { // Somoni
        code: 'TJS',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'SM'
        },
        unit: 'Somoni'
    },
    'TMT': { // Turkmenistan New Manat
        code: 'TMT',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'm'
        },
        unit: 'Manat'
    },
    'TND': { // Tunisian Dinar
        code: 'TND',
        type: CurrencyType.Fiat,
        fraction: 3,
        symbol: {
            normal: 'DT'
        },
        unit: 'Dinar'
    },
    'TOP': { // Pa’anga
        code: 'TOP',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'T$'
        },
        unit: 'Paanga'
    },
    'TRY': { // Turkish Lira
        code: 'TRY',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: '₺'
        },
        unit: 'Lira'
    },
    'TTD': { // Trinidad and Tobago Dollar
        code: 'TTD',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: '$'
        },
        unit: 'Dollar'
    },
    'TWD': { // New Taiwan Dollar
        code: 'TWD',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'NT$'
        },
        unit: 'Dollar'
    },
    'TZS': { // Tanzanian Shilling
        code: 'TZS',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: '/='
        },
        unit: 'Shilling'
    },
    'UAH': { // Hryvnia
        code: 'UAH',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: '₴'
        },
        unit: 'Hryvnia'
    },
    'UGX': { // Uganda Shilling
        code: 'UGX',
        type: CurrencyType.Fiat,
        fraction: 0,
        symbol: {
            normal: '/='
        },
        unit: 'Shilling'
    },
    'USD': { // US Dollar
        code: 'USD',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: '$'
        },
        unit: 'Dollar'
    },
    'UYU': { // Peso Uruguayo
        code: 'UYU',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: '$'
        },
        unit: 'Peso'
    },
    'UZS': { // Uzbekistan Sum
        code: 'UZS',
        type: CurrencyType.Fiat,
        fraction: 2,
        unit: 'Sum'
    },
    'VED': { // Bolívar Soberano
        code: 'VED',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'Bs.D'
        },
        unit: 'Bolivar'
    },
    'VES': { // Bolívar Soberano
        code: 'VES',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'Bs.S'
        },
        unit: 'Bolivar'
    },
    'VND': { // Dong
        code: 'VND',
        type: CurrencyType.Fiat,
        fraction: 0,
        symbol: {
            normal: '₫'
        },
        unit: 'Dong'
    },
    'VUV': { // Vatu
        code: 'VUV',
        type: CurrencyType.Fiat,
        fraction: 0,
        symbol: {
            normal: 'VT'
        },
        unit: 'Vatu'
    },
    'WST': { // Tala
        code: 'WST',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: '$'
        },
        unit: 'Tala'
    },
    'XAF': { // CFA Franc BEAC
        code: 'XAF',
        type: CurrencyType.Fiat,
        fraction: 0,
        symbol: {
            normal: 'F.CFA'
        },
        unit: 'Franc'
    },
    'XCD': { // East Caribbean Dollar
        code: 'XCD',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: '$'
        },
        unit: 'Dollar'
    },
    'XOF': { // CFA Franc BCEAO
        code: 'XOF',
        type: CurrencyType.Fiat,
        fraction: 0,
        symbol: {
            normal: 'F.CFA'
        },
        unit: 'Franc'
    },
    'XPF': { // CFP Franc
        code: 'XPF',
        type: CurrencyType.Fiat,
        fraction: 0,
        symbol: {
            normal: 'F'
        },
        unit: 'Franc'
    },
    'XSU': { // Sucre
        code: 'XSU',
        type: CurrencyType.Fiat,
        symbol: {
            normal: 'S/.'
        },
        unit: 'Sucre'
    },
    'YER': { // Yemeni Rial
        code: 'YER',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'YRl',
            plural: 'YRls'
        },
        unit: 'Rial'
    },
    'ZAR': { // Rand
        code: 'ZAR',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'R'
        },
        unit: 'Rand'
    },
    'ZMW': { // Zambian Kwacha
        code: 'ZMW',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'K'
        },
        unit: 'Kwacha'
    },
    'ZWG': { // Zimbabwe Gold
        code: 'ZWG',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: 'ZiG'
        },
        unit: 'ZiG'
    },
    'ZWL': { // Zimbabwe Dollar
        code: 'ZWL',
        type: CurrencyType.Fiat,
        fraction: 2,
        symbol: {
            normal: '$'
        },
        unit: 'Dollar'
    }
};

// Cryptocurrencies
export const ALL_CRYPTOCURRENCIES: Record<string, CurrencyInfo> = {
    'BTC': { // Bitcoin
        code: 'BTC',
        type: CurrencyType.Cryptocurrency,
        fraction: 8,
        symbol: {
            normal: '₿'
        },
        unit: 'Bitcoin'
    },
    'ETH': { // Ethereum
        code: 'ETH',
        type: CurrencyType.Cryptocurrency,
        fraction: 8,
        symbol: {
            normal: 'Ξ'
        },
        unit: 'Ethereum'
    },
    'BNB': { // Binance Coin
        code: 'BNB',
        type: CurrencyType.Cryptocurrency,
        fraction: 8,
        symbol: {
            normal: 'BNB'
        },
        unit: 'Binance Coin'
    },
    'SOL': { // Solana
        code: 'SOL',
        type: CurrencyType.Cryptocurrency,
        fraction: 8,
        symbol: {
            normal: '◎'
        },
        unit: 'Solana'
    },
    'ADA': { // Cardano
        code: 'ADA',
        type: CurrencyType.Cryptocurrency,
        fraction: 8,
        symbol: {
            normal: '₳'
        },
        unit: 'Cardano'
    },
    'XRP': { // Ripple
        code: 'XRP',
        type: CurrencyType.Cryptocurrency,
        fraction: 8,
        symbol: {
            normal: 'XRP'
        },
        unit: 'Ripple'
    },
    'DOT': { // Polkadot
        code: 'DOT',
        type: CurrencyType.Cryptocurrency,
        fraction: 8,
        symbol: {
            normal: 'DOT'
        },
        unit: 'Polkadot'
    },
    'DOGE': { // Dogecoin
        code: 'DOGE',
        type: CurrencyType.Cryptocurrency,
        fraction: 8,
        symbol: {
            normal: 'Ð'
        },
        unit: 'Dogecoin'
    },
    'MATIC': { // Polygon
        code: 'MATIC',
        type: CurrencyType.Cryptocurrency,
        fraction: 8,
        symbol: {
            normal: 'MATIC'
        },
        unit: 'Polygon'
    },
    'USDT': { // Tether
        code: 'USDT',
        type: CurrencyType.Cryptocurrency,
        fraction: 2,
        symbol: {
            normal: 'USDT'
        },
        unit: 'Tether'
    }
};

export const DEFAULT_CURRENCY_SYMBOL: string = '¤';
export const PARENT_ACCOUNT_CURRENCY_PLACEHOLDER: string = '---';

/**
 * Get currency information by currency code and optional type
 * @param currencyCode Currency code
 * @param type Currency type (optional, if provided only searches in the specified type, more efficient)
 * @returns Currency information or undefined if not found
 */
export function getCurrencyInfo(currencyCode: string, type?: CurrencyType): CurrencyInfo | undefined {
    // If type is specified, only search in the corresponding type
    if (type === CurrencyType.Fiat) {
        return ALL_FIAT_CURRENCIES[currencyCode];
    } else if (type === CurrencyType.Cryptocurrency) {
        return ALL_CRYPTOCURRENCIES[currencyCode];
    } else if (type === CurrencyType.Security) {
        // return ALL_SECURITIES[currencyCode];  // Future: when securities are added
        return undefined;
    }
    
    // If type is not specified, search in all types (not priority-based, but searching all possible locations)
    return ALL_FIAT_CURRENCIES[currencyCode] 
        || ALL_CRYPTOCURRENCIES[currencyCode]
        // || ALL_SECURITIES[currencyCode]  // Future: when securities are added
        || undefined;
}

/**
 * Get currency type by currency code
 * @param currencyCode Currency code
 * @returns Currency type or undefined if not found
 */
export function getCurrencyType(currencyCode: string): CurrencyType | undefined {
    if (ALL_FIAT_CURRENCIES[currencyCode]) {
        return CurrencyType.Fiat;
    } else if (ALL_CRYPTOCURRENCIES[currencyCode]) {
        return CurrencyType.Cryptocurrency;
    }
    // else if (ALL_SECURITIES[currencyCode]) {
    //     return CurrencyType.Security;
    // }
    return undefined;
}

/**
 * Get all currency codes (for iteration)
 * @returns Array of all currency codes
 */
export function getAllCurrencyCodes(): string[] {
    return [
        ...Object.keys(ALL_FIAT_CURRENCIES),
        ...Object.keys(ALL_CRYPTOCURRENCIES),
        // ...Object.keys(ALL_SECURITIES),  // Future: when securities are added
    ];
}

export const DEFAULT_CURRENCY_CODE: string = getCurrencyInfo('USD')?.code || 'USD';
