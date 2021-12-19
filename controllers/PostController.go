package controllers

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/DhruvikDonga/tabsmooth/config/migrations"
	"github.com/DhruvikDonga/tabsmooth/config/sessions"
	"github.com/DhruvikDonga/tabsmooth/models"
	"github.com/DhruvikDonga/tabsmooth/views"
	"github.com/gomarkdown/markdown"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
)

var createpost *views.View

func GetBlogs(w http.ResponseWriter, r *http.Request) {
	session, _ := sessions.Get(r)
	flashes := session.Flashes()
	//log.Print("from get:-", session.Name())

	data := map[string]interface{}{}
	data[csrf.TemplateTag] = csrf.TemplateField(r)

	session.Save(r, w)
	var blogs []models.Blogs

	migrations.DB.Model(&models.Blogs{}).Select([]string{"title", "blogdescription", "created_at", "banner", "slug"}).Find(&blogs)
	//log.Println(blogs)

	data["Content"] = blogs
	if len(flashes) > 0 {
		//log.Println(flash)
		data["Error"] = flashes[0]
	}
	//log.Println("from get:-", flashes)

	index = views.NewView("bootstrap", "views/Blogs.gohtml")

	index.Render(w, data)

}

func SingleBlog(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]
	data := map[string]interface{}{}

	var blog models.Blogs
	migrations.DB.First(&blog, "slug=?", slug)
	var user models.Users
	migrations.DB.First(&user, blog.UsersID)

	content := []byte(blog.Content)
	md := template.HTML(markdown.ToHTML(content, nil, nil))
	data["Content"] = md
	data["Date"] = blog.CreatedAt.Format("January 02, 2006")
	data["banner"] = blog.Banner
	data["title"] = blog.Title
	data["description"] = blog.Blogdescription

	data["Name"] = user.Name
	data["linkedin"] = user.Linkedin
	data["facebook"] = user.Facebook
	data["reddit"] = user.Reddit
	data["twitter"] = user.Twitter
	data["instagram"] = user.Instagram
	index = views.NewView("bootstrap", "views/singleblog.gohtml")

	index.Render(w, data)
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	session, _ := sessions.Get(r)
	flashes := session.Flashes()
	//log.Print("from get:-", session.Name())
	user := UserSession(session)
	//log.Println(session.Values["user"])
	if user.Authenticated == false {
		session.AddFlash("Not authorized sorry")
		session.Save(r, w)

		http.Redirect(w, r, "/blogs", http.StatusFound)
		return
	}
	data := map[string]interface{}{}
	data[csrf.TemplateTag] = csrf.TemplateField(r)
	data["username"] = user.Name
	session.Save(r, w)

	if len(flashes) > 0 {
		//log.Println(flash)
		data["Error"] = flashes[0]
	}
	//log.Println("from get:-", flashes)

	index = views.NewView("bootstrap", "views/postblog.gohtml")

	index.Render(w, data)

}

func GetUserBlogs(w http.ResponseWriter, r *http.Request) {
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

	data["blogs"] = blogs
	data["User"] = getuser
	if len(flashes) > 0 {
		//log.Println(flash)
		data["Error"] = flashes[0]
	}
	//log.Println("from get:-", flashes)

	index = views.NewView("bootstrap", "views/userblogs.gohtml")
	data[csrf.TemplateTag] = csrf.TemplateField(r)

	index.Render(w, data)

}
func PostBlog(w http.ResponseWriter, r *http.Request) {
	session, _ := sessions.Get(r)
	//log.Println(r.MultipartForm)

	file, handler, err := r.FormFile("banner")
	if err != nil {
		//log.Println("Error Retrieving the File")
		log.Println(err)
		return
	}
	defer file.Close()
	//log.Println(filepath.Ext(handler.Filename))
	var blog models.Blogs
	blog.Title = r.FormValue("title")
	blog.Slug = strings.ReplaceAll(blog.Title, " ", "-")

	newname := blog.Slug + "-banner" + filepath.Ext(handler.Filename)
	var getuser models.Users
	migrations.DB.Model(&models.Users{}).Where("name=?", UserSession(session).Name).First(&getuser)

	blog.UsersID = getuser.ID
	blog.Blogdescription = r.FormValue("description")
	blog.Content = r.FormValue("content")
	blog.Banner = "/" + newname
	//log.Println(blog)
	ext := filepath.Ext(handler.Filename)
	//log.Println(ext)
	var check models.Blogs
	migrations.DB.Model(&models.Blogs{}).Where("slug=?", blog.Slug).First(&check)
	if check.Slug != "" {

		session.AddFlash("Use unique slug ,title already taken")
		//log.Print("from post:-", session.Name())
		session.Save(r, w)

		http.Redirect(w, r, "/createblog", http.StatusFound)
		return
	}
	//image upload
	if ext != ".png" {

		session.AddFlash("Only upload png or jpeg files")
		//log.Print("from post:-", session.Name())
		session.Save(r, w)

		http.Redirect(w, r, "/createblog", http.StatusFound)
		return
	} else {
		outfile, err := os.Create("./static/image/" + newname)
		if err != nil {
			log.Println(err)
			return
		}
		defer outfile.Close()
		_, err = io.Copy(outfile, file)

		if err != nil {
			log.Println("Internal server error")
			return
		}
	}
	migrations.DB.Create(&blog)
	session.AddFlash("Blog created succesfully 📚")
	//log.Print("from post:-", session.Name())
	session.Save(r, w)

	http.Redirect(w, r, "/blogs", http.StatusSeeOther)

}
func EditBlogForm(w http.ResponseWriter, r *http.Request) {
	session, _ := sessions.Get(r)
	//log.Print("from get:-", session.Name())

	data := map[string]interface{}{}

	session.Save(r, w)
	user := UserSession(session)
	//log.Println(session.Values["user"])
	if user.Authenticated == false {
		session.AddFlash("Not authorized sorry")
		session.Save(r, w)

		http.Redirect(w, r, "/blogs", http.StatusFound)
		return
	}
	//log.Println(blogs)

	vars := mux.Vars(r)
	slug := vars["slug"]

	var blog models.Blogs
	migrations.DB.First(&blog, "slug=?", slug)

	data["Content"] = blog.Content
	data["blogid"] = blog.ID
	data["banner"] = blog.Banner
	data["title"] = blog.Title
	data["Description"] = blog.Blogdescription

	index = views.NewView("bootstrap", "views/editblog.gohtml")
	data[csrf.TemplateTag] = csrf.TemplateField(r)

	index.Render(w, data)

}
func EditBlog(w http.ResponseWriter, r *http.Request) {
	session, _ := sessions.Get(r)

	var editblog models.Blogs
	migrations.DB.Model(&models.Blogs{}).Where(r.FormValue("id")).First(&editblog)

	var check models.Blogs
	if editblog.Title != r.FormValue("title") {
		migrations.DB.Model(&models.Blogs{}).Where("title=?", r.FormValue("title")).First(&check)

		if check.Title != "" {

			session.AddFlash("Use unique title ,title already taken")
			session.Save(r, w)

			http.Redirect(w, r, "/editblog/"+editblog.Slug, http.StatusFound)
			return
		}
	}
	editblog.Title = r.FormValue("title")
	editblog.Slug = strings.ReplaceAll(editblog.Title, " ", "-")
	editblog.Blogdescription = r.FormValue("description")
	editblog.Content = r.FormValue("content")
	//banner stuff
	filethere := true
	file, handler, err := r.FormFile("banner")
	if err != nil {
		//log.Println("Error Retrieving the File")
		log.Println(err)
		filethere = false

	}
	if filethere == true {
		var newname string
		oldname := editblog.Banner[1:]
		if handler.Filename != "" {
			newname = editblog.Slug + "-banner" + filepath.Ext(handler.Filename)
			editblog.Banner = "/" + newname
		}
		defer file.Close()
		ext := filepath.Ext(handler.Filename)
		//log.Println(ext)
		//image upload
		if ext != ".png" {

			session.AddFlash("Only upload png or jpeg files")
			//log.Print("from post:-", session.Name())
			session.Save(r, w)

			http.Redirect(w, r, "/createblog", http.StatusFound)
			return
		} else {
			outfile, err := os.Create("./static/image/" + newname)
			if err != nil {
				log.Println(err)
				return
			}
			defer outfile.Close()
			_, err = io.Copy(outfile, file)

			if err != nil {
				log.Println("Internal server error")
				return
			}
			err = os.Remove("./static/image/" + oldname)
			if err != nil {
				log.Println(err)
				return
			}
		}
	}
	migrations.DB.Save(&editblog)
	session.AddFlash("Successfully updated ")

	session.Save(r, w)
	http.Redirect(w, r, "/your-blogs", http.StatusFound)

}
