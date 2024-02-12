package api

// HTTPErrResponse used in
// responses from server when
// an error occurred.
type HTTPErrResponse struct {
	Error any `json:"error"`
}
