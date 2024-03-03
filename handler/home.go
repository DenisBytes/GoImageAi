package handler

import (
	"fmt"
	"net/http"

	// "com.github.denisbytes.goimageai/db"
	"com.github.denisbytes.goimageai/view/home"
)

func HandleHomeIndex(w http.ResponseWriter, r *http.Request) error {
	user := GetAuthenticatedUser(r)
	// account, err := db.GetAccountByUserID(user.ID)
	// if err != nil {
	// 	return err
	// }

	fmt.Printf("%+v\n", user.Account)
	return home.Index().Render(r.Context(), w)
}
