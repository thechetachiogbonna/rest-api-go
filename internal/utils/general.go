package utils

import "time"

const FIFTEEN_MINUTES = time.Minute * 15
const THIRTY_DAYS = time.Hour * 24 * 30
const ISSUER = "Rest-API-Go"
const AUDIENCE = "user"
const REFRESH_PATH = "/auth/refresh"
