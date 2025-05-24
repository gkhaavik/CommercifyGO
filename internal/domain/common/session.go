package common

const (
	SessionCookieName     = "guest_session_id"
	CheckoutSessionCookie = "checkout_session_id"
	SessionCookieAge      = 86400 * 30 // 30 days in seconds
	CheckoutSessionMaxAge = 86400 * 7  // 7 days in seconds
)
