package controllers

import (
	"log"
	"net/http"

	"github.com/DhruvikDonga/tabsmooth/config/migrations"
	"github.com/DhruvikDonga/tabsmooth/config/sessions"
	"github.com/DhruvikDonga/tabsmooth/models"
	"github.com/DhruvikDonga/tabsmooth/views"
	"github.com/gorilla/csrf"
	session "github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

func UserSession(s *session.Session) models.Auth {
	val := s.Values["user"]
	var user = models.Auth{}
	user, ok := val.(models.Auth)
	if !ok {
		return models.Auth{
			Authenticated: false,
		}
	}
	return user

}

func RegisterPage(w http.ResponseWriter, r *http.Request) {
	session, _ := sessions.Get(r)
	flashes := session.Flashes()
	//log.Print("from get:-", session.Name())
	user := UserSession(session)
	//log.Println(session.Values["user"])
	if user.Authenticated == true {
		session.AddFlash("Already signed in!!")
		session.Save(r, w)

		http.Redirect(w, r, "/profile", http.StatusFound)
		return
	}
	data := map[string]interface{}{}
	data[csrf.TemplateTag] = csrf.TemplateField(r)

	session.Save(r, w)

	if len(flashes) > 0 {
		//log.Println(flash)
		data["Error"] = flashes[0]
	}
	//log.Println("from get:-", flashes)

	index = views.NewView("bootstrap", "views/register.gohtml")

	index.Render(w, data)

}
func LoginPage(w http.ResponseWriter, r *http.Request) {
	session, _ := sessions.Get(r)
	flashes := session.Flashes()
	//log.Print("from get:-", session.Name())
	user := UserSession(session)
	//log.Println(session.Values["user"])
	if user.Authenticated == true {
		session.AddFlash("Already signed in!!")
		session.Save(r, w)

		http.Redirect(w, r, "/profile", http.StatusFound)
		return
	}
	data := map[string]interface{}{}
	data[csrf.TemplateTag] = csrf.TemplateField(r)

	session.Save(r, w)

	if len(flashes) > 0 {
		//log.Println(flash)
		data["Error"] = flashes[0]
	}
	//log.Println("from get:-", flashes)

	index = views.NewView("bootstrap", "views/login.gohtml")

	index.Render(w, data)

}
func PostLogin(w http.ResponseWriter, r *http.Request) {
	session, _ := sessions.Get(r)
	if session.ID != "" {

	}
	var getuser models.Users

	migrations.DB.Model(&models.Users{}).Where("email=?", r.FormValue("email")).First(&getuser)

	if getuser.Email == "" {
		session.AddFlash("E-mail does not exsist")
		session.Save(r, w)

		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(getuser.Password), []byte(r.FormValue("password")))

	if err != nil {
		session.AddFlash("Wrong password")
		session.Save(r, w)

		http.Redirect(w, r, "/login", http.StatusFound)
		return

	}
	session.AddFlash("Successfully login ")

	auth := models.Auth{
		Name:          getuser.Name,
		Id:            getuser.ID,
		Authenticated: true,
	}
	session.Values["user"] = auth
	log.Println(session.Values["user"])
	session.Save(r, w)
	http.Redirect(w, r, "/createblog", http.StatusFound)

}

func PostLogout(w http.ResponseWriter, r *http.Request) {
	session, err := sessions.Get(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["user"] = models.Auth{}
	session.Options.MaxAge = -1

	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}
func PostRegister(w http.ResponseWriter, r *http.Request) {
	session, _ := sessions.Get(r)

	var user models.Users
	user.Name = r.FormValue("name")
	user.Email = r.FormValue("email")
	var check models.Users
	migrations.DB.Model(&models.Users{}).Where("email=?", user.Email).First(&check)

	if check.Name != "" {

		session.AddFlash("Use unique email ,email already taken")
		session.Save(r, w)

		http.Redirect(w, r, "/register", http.StatusFound)
		return
	}
	if r.FormValue("password") != r.FormValue("confirmpassword") {
		session.AddFlash("Passwords not matching" + r.FormValue("password") + r.FormValue("confirmpassword"))
		session.Save(r, w)

		http.Redirect(w, r, "/register", http.StatusFound)
		return
	}
	user.Email = r.FormValue("email")
	bytes, _ := bcrypt.GenerateFromPassword([]byte(r.FormValue("password")), 14)

	user.Password = string(bytes)
	migrations.DB.Create(&user)
	session.AddFlash("Successfully registered ")

	auth := models.Auth{
		Name:          user.Name,
		Id:            user.ID,
		Authenticated: true,
	}
	session.Values["user"] = auth

	err := session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/createblog", http.StatusFound)

}
func GetUserProfile(w http.ResponseWriter, r *http.Request) {
	session, _ := sessions.Get(r)
	flashes := session.Flashes()
	//log.Print("from get:-", session.Name())

	data := map[string]interface{}{}

	session.Save(r, w)
	var blogs []models.Blogs
	user := UserSession(session)
	//log.Println(session.Values["user"])
	if user.Authenticated == false {
		session.AddFlash("Not authorized sorry")
		session.Save(r, w)

		http.Redirect(w, r, "/blogs", http.StatusFound)
		return
	}
	migrations.DB.Model(&models.Blogs{}).Where("users_id=?", user.Id).Find(&blogs)
	//log.Println(blogs)

	var getuser []models.Users
	migrations.DB.First(&getuser, user.Id)

	data["blogs"] = len(blogs)
	data["User"] = getuser
	if len(flashes) > 0 {
		//log.Println(flash)
		data["Error"] = flashes[0]
	}
	//log.Println("from get:-", flashes)

	index = views.NewView("bootstrap", "views/profile.gohtml")
	data[csrf.TemplateTag] = csrf.TemplateField(r)

	index.Render(w, data)

}
func ResetProfile(w http.ResponseWriter, r *http.Request) {
	session, _ := sessions.Get(r)

	var user models.Users
	migrations.DB.Model(&models.Users{}).Where(session.ID).First(&user)

	user.Name = r.FormValue("name")
	var check models.Users
	if user.Email != r.FormValue("email") {
		migrations.DB.Model(&models.Users{}).Where("email=?", r.FormValue("email")).First(&check)

		if check.Name != "" {

			session.AddFlash("Use unique email ,email already taken")
			session.Save(r, w)

			http.Redirect(w, r, "/profile", http.StatusFound)
			return
		}
	}
	user.Email = r.FormValue("email")
	user.Linkedin = r.FormValue("linkedin")
	user.Facebook = r.FormValue("facebook")
	user.Instagram = r.FormValue("instagram")
	user.Youtube = r.FormValue("youtube")
	user.Twitter = r.FormValue("twitter")
	user.Reddit = r.FormValue("reddit")
	user.Personal = r.FormValue("website")

	migrations.DB.Save(&user)
	session.AddFlash("Successfully updated ")

	session.Save(r, w)
	http.Redirect(w, r, "/profile", http.StatusFound)

}
