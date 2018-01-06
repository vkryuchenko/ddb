package app

import (
	"fmt"
	"helpers"
	"log"
	"net/http"
	"time"
)

const (
	queryDropSession = `
DELETE FROM public.sessions
WHERE session_id = $1 ;
`
	querySessionById = `
SELECT
  user_id,
  expire,
  session_id
FROM public.sessions
WHERE session_id = $1;
`
	queryUpdateSession = `
INSERT INTO public.sessions
(
  user_id,
  expire,
  session_id
)
VALUES ($1, $2, $3)
 ON CONFLICT (user_id)
  DO UPDATE
    SET (expire, session_id) = ($2, $3);
`
)

type sessionInfo struct {
	UserId uint
	Expire time.Time
	UUID   string
}

func (p *Provider) sessionGetByUUID(uuid string) (sessionInfo, error) {
	var session sessionInfo

	db, err := p.Database.Connect()
	if err != nil {
		return session, err
	}
	defer db.Close()

	query, err := db.Prepare(querySessionById)
	if err != nil {
		return session, err
	}
	result := query.QueryRow(uuid)
	err = result.Scan(&session.UserId, &session.Expire, &session.UUID)
	if err != nil {
		return session, err
	}
	return session, nil
}

func (p *Provider) sessionValid(w http.ResponseWriter, r *http.Request) bool {
	cookie, err := r.Cookie(p.ApplicationName)
	if err != nil {
		return false
	}
	_, err = p.sessionCheck(cookie)
	if err != nil {
		p.sessionDrop(w, r)
		return false
	}
	return true
}

func (p *Provider) sessionUpdate(w http.ResponseWriter, s sessionInfo) error {
	uuid, err := helpers.NewUUID()
	if err != nil {
		return err
	}
	cookie := http.Cookie{
		Name:     p.ApplicationName,
		Value:    uuid,
		Expires:  time.Now().Add(8 * time.Hour),
		Path:     "/",
		HttpOnly: true,
	}
	s.Expire = cookie.Expires
	s.UUID = uuid

	db, err := p.Database.Connect()
	if err != nil {
		return err
	}
	defer db.Close()

	txn, err := db.Begin()
	if err != nil {
		return err
	}
	query, err := txn.Prepare(queryUpdateSession)
	if err != nil {
		txn.Rollback()
		return err
	}
	_, err = query.Exec(s.UserId, s.Expire, s.UUID)
	if err != nil {
		txn.Rollback()
		return err
	}
	err = txn.Commit()
	if err != nil {
		return err
	}

	w.Header().Set("Set-Cookie", cookie.String())
	return nil
}

func (p *Provider) sessionCheck(cookie *http.Cookie) (sessionInfo, error) {
	session, err := p.sessionGetByUUID(cookie.Value)
	if err != nil {
		return session, err
	}
	if session.Expire.Unix() <= time.Now().Unix() {
		return session, fmt.Errorf("session expired")
	}
	return session, nil
}

func (p *Provider) sessionDrop(w http.ResponseWriter, r *http.Request) {
	defer http.Redirect(w, r, "/login", http.StatusFound)
	cookie, err := r.Cookie(p.ApplicationName)
	if err != nil {
		log.Print(err)
		return
	}
	db, err := p.Database.Connect()
	if err != nil {
		log.Print(err)
		return
	}
	defer db.Close()
	query, err := db.Prepare(queryDropSession)
	if err != nil {
		log.Print(err)
		return
	}
	_, err = query.Exec(cookie.Value)
	if err != nil {
		log.Print(err)
		return
	}
}
