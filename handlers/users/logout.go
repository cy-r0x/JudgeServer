package users

import "net/http"

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	//TODO: add the token to blacklist table
}
