package gs

import (
	"plata_currency_quotation/internal/persistence"
	"plata_currency_quotation/internal/service/currency-conversion"
)

var CurrencyConversionService currency_conversion.Interface

var Db persistence.Interface
