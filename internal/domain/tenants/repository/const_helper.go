package repository

import "appointment_management_system/internal/pkg/helper"

var helperGetDB = helper.GetDB

type contextKey string

const ginContextKey contextKey = "ginContext"
