package app

import (
	"database/sql"
	"log"
	"net/http"
)

const (
	queryCreateUser = `
INSERT INTO public.users
(
  login,
  email,
  active,
  admin,
  quota
)
VALUES (
  $1,
  $2,
  DEFAULT,
  DEFAULT,
  $3
)
RETURNING id;
`
	queryGetUserInfo = `
SELECT
  id,
  login,
  email,
  active,
  admin,
  quota
FROM public.users
WHERE login = $1;
`
)

type UserInfo struct {
	Id     uint
	Login  string
	Email  string
	Active bool
	Admin  bool
	Quota  uint
}

func (p *Provider) createUser(name, email string) (UserInfo, error) {
	ui := UserInfo{}
	db, err := p.Database.Connect()
	if err != nil {
		return ui, err
	}
	query, err := db.Prepare(queryCreateUser)
	if err != nil {
		return ui, err
	}
	result, err := query.Exec(name, email, p.Docker.Quota)
	if err != nil {
		return ui, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return ui, err
	}
	ui.Id = uint(id)
	ui.Login = name
	ui.Email = email
	ui.Active = true
	ui.Admin = false
	ui.Quota = p.Docker.Quota
	return ui, nil
}

func (p *Provider) getUserInfo(name string) (UserInfo, error) {
	ui := UserInfo{}
	db, err := p.Database.Connect()
	if err != nil {
		return ui, err
	}
	query, err := db.Prepare(queryGetUserInfo)
	if err != nil {
		return ui, err
	}
	err = query.QueryRow(name).Scan(
		&ui.Id,
		&ui.Login,
		&ui.Email,
		&ui.Active,
		&ui.Admin,
		&ui.Quota,
	)
	if err != nil {
		return ui, err
	}
	return ui, nil
}

func (p *Provider) actionAuth(w http.ResponseWriter, r *http.Request) {
	var user UserInfo
	username := r.FormValue("login")
	password := r.FormValue("pass")
	ok, data, err := p.LDAPClient.Authenticate(username, password)
	if !ok || err != nil {
		log.Printf("%s auth failed. %s", username, err)
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	user, err = p.getUserInfo(username)
	if err != nil {
		if err == sql.ErrNoRows {
			var email string
			email = data["email"]
			user, err = p.createUser(username, email)
			if err != nil {
				log.Print(err)
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}
		} else {
			log.Print(err)
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
	}
	session := sessionInfo{
		UserId: user.Id,
	}
	err = p.updateSession(w, session)
	if err != nil {
		log.Print(err)
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
