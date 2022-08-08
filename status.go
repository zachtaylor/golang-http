package http

import "net/http"

const (
	// StatusOK = 200
	StatusOK = http.StatusOK
	// StatusPartialContent = 206
	StatusPartialContent = http.StatusPartialContent
	// StatusMovedPermanently = 301
	StatusMovedPermanently = http.StatusMovedPermanently
	// StatusNotModified = 304
	StatusNotModified = http.StatusNotModified
	// StatusBadRequest = 400
	StatusBadRequest = http.StatusBadRequest
	// StatusUnauthorized = 401
	StatusUnauthorized = http.StatusUnauthorized
	// StatusPaymentRequired = 402
	StatusPaymentRequired = http.StatusPaymentRequired
	// StatusForbidden = 403
	StatusForbidden = http.StatusForbidden
	// StatusNotFound = 404
	StatusNotFound = http.StatusNotFound
	// StatusMethodNotAllowed = 405
	StatusMethodNotAllowed = http.StatusMethodNotAllowed
	// StatusConflict = 409
	StatusConflict = http.StatusConflict
	// StatusPreconditionFailed = 412
	StatusPreconditionFailed = http.StatusPreconditionFailed
	// StatusRequestedRangeNotSatisfiable = 416
	StatusRequestedRangeNotSatisfiable = http.StatusRequestedRangeNotSatisfiable
	// StatusInternalServerError = 500
	StatusInternalServerError = http.StatusInternalServerError
)
