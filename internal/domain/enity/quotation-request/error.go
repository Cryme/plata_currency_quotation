package quotation_request

import "errors"

var ErrSameCurrency = errors.New("base and quote currency cannot be the same")
