package common

import "errors"

var OutOfOrderMessage = errors.New("message arrived out of order")
var DuplicateMessage = errors.New("duplicate message")
