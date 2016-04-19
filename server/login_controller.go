package server

import (
	"database/sql"
	"github.com/osvaldshpengler/browsercalls/tools"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"net/http"
	"os"
)

type loginController struct{}

type loginPage struct {
	HasErrors     bool
	HasFormValues bool
	Error         string
	Username      string
	Password      string
}

func (l *loginController) handleLogin(rw http.ResponseWriter, r *http.Request) {
	stat := loginPage{}

	if "GET" == r.Method {
		t, err := template.ParseFiles(os.Getenv("BC_APP_PATH") + "frontend/login.html")
		if nil != err {
			tools.Log.Error(err)
		}
		t.Execute(rw, &stat)
		return
	}

	if "POST" != r.Method {
		http.Error(rw, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		tools.Log.Error(err)
	}

	stat.Username = r.PostFormValue("username")
	stat.Password = r.PostFormValue("password")
	stat.HasFormValues = true

	dba, err := tools.GetDbAccessor()
	if err != nil {
		tools.Log.Error(err)
	}
	var password string
	var id int
	err = dba.QueryRow("SELECT id, password FROM users WHERE username = $1", stat.Username).Scan(&id, &password)

	if sql.ErrNoRows == err {
		stat.HasErrors = true
		stat.Error = "пользователь с таким именем не найден"
		t, _ := template.ParseFiles(os.Getenv("BC_APP_PATH") + "frontend/login.html")
		t.Execute(rw, &stat)
		return
	}
	if nil != err {
		tools.Log.Error(err)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(password), []byte(stat.Password)); nil != err {
		stat.HasErrors = true
		stat.Error = "введен неверный пароль"
		t, _ := template.ParseFiles(os.Getenv("BC_APP_PATH") + "frontend/login.html")
		t.Execute(rw, &stat)
		return
	}

	u := &User{id, stat.Username, password}

	if err = initUserSession(rw, r, u); nil != err {
		tools.Log.Error(err)
	}

	http.Redirect(rw, r, "/", http.StatusFound)
}

func (l *loginController) handleRegister(rw http.ResponseWriter, r *http.Request) {
	if "POST" != r.Method {
		http.Error(rw, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		tools.Log.Error(err)
	}

	stat := loginPage{}

	stat.Username = r.PostFormValue("username")
	stat.Password = r.PostFormValue("password")
	stat.HasFormValues = true

	if len(stat.Password) < 6 {
		stat.HasErrors = true
		stat.Error = "пароль должен быть не короче 6 символов"
		t, _ := template.ParseFiles(os.Getenv("BC_APP_PATH") + "frontend/login.html")
		t.Execute(rw, &stat)
		return
	}

	password, err := bcrypt.GenerateFromPassword([]byte(stat.Password), bcrypt.DefaultCost)
	if nil != err {
		tools.Log.Error(err)
	}

	dba, err := tools.GetDbAccessor()
	if nil != err {
		tools.Log.Error(err)
	}
	var cnt int
	err = dba.QueryRow("SELECT COUNT(*) FROM users WHERE username = $1", stat.Username).Scan(&cnt)
	if nil != err {
		tools.Log.Error(err.Error())
	}

	if cnt > 0 {
		stat.HasErrors = true
		stat.Error = "пользователь с таким именем уже зарегистрирован"
		t, _ := template.ParseFiles(os.Getenv("BC_APP_PATH") + "frontend/login.html")
		t.Execute(rw, &stat)
		return
	}

	var id int
	err = dba.QueryRow("INSERT INTO users(username, password) VALUES ($1, $2) RETURNING id", stat.Username, password).Scan(&id)
	if nil != err {
		tools.Log.Error(err)
	}

	u := &User{id, stat.Username, string(password)}

	if err = initUserSession(rw, r, u); nil != err {
		tools.Log.Error(err)
	}

	http.Redirect(rw, r, "/", http.StatusFound)
}
