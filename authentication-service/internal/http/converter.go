package http

import (
    "net/http"

    "github.com/EmilioCliff/payment-polling-app/authentication-service/pkg"
)

func (s *HttpServer) convertPkgError(err *pkg.Error) int {
    switch err.Code {
    case pkg.ALREADY_EXISTS_ERROR:
        return http.StatusConflict
    case pkg.INTERNAL_ERROR:
        return http.StatusInternalServerError
    case pkg.INVALID_ERROR:
        return http.StatusBadRequest
    case pkg.NOT_FOUND_ERROR:
        return http.StatusNotFound
    case pkg.NOT_IMPLEMENTED_ERROR:
        return http.StatusNotImplemented
    case pkg.AUTHENTICATION_ERROR:
        return http.StatusUnauthorized
    default:
        return http.StatusInternalServerError
    }
}