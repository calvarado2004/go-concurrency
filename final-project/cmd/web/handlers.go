package main

import (
	"errors"
	"fmt"
	"github.com/calvarado2004/go-concurrency/data"
	"github.com/phpdave11/gofpdf"
	"github.com/phpdave11/gofpdf/contrib/gofpdi"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"time"
)

func (app *Config) HomePage(w http.ResponseWriter, r *http.Request) {

	app.render(w, r, "home.page.gohtml", nil)

}

func (app *Config) LoginPage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.gohtml", nil)
}

func (app *Config) PostLoginPage(w http.ResponseWriter, r *http.Request) {

	_ = app.Session.RenewToken(r.Context())

	// parse form post

	err := r.ParseForm()
	if err != nil {
		app.ErrorLog.Println(err)
	}

	// get email and password from post form

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	user, err := app.Models.User.GetByEmail(email)
	if err != nil {
		app.Session.Put(r.Context(), "error", "Invalid login credentials")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// check password
	validPassword, err := app.Models.User.PasswordMatches(password)
	if err != nil {
		app.Session.Put(r.Context(), "error", "Invalid login credentials")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if !validPassword {
		msg := Message{
			To:      email,
			Subject: "Failed login attempt",
			Data:    "Invalid login attempt",
		}
		app.sendEmail(msg)
		app.Session.Put(r.Context(), "error", "Invalid login credentials")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// ok, login
	app.Session.Put(r.Context(), "userID", user.ID)
	app.Session.Put(r.Context(), "user", user)

	app.Session.Put(r.Context(), "flash", "You've been logged in successfully")

	// redirect to home page
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func (app *Config) LogoutPage(w http.ResponseWriter, r *http.Request) {

	// clean up session
	_ = app.Session.Destroy(r.Context())
	_ = app.Session.RenewToken(r.Context())

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *Config) RegisterPage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "register.page.gohtml", nil)

}

func (app *Config) PostRegisterPage(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		app.ErrorLog.Println(err)
	}

	// TODO - validate data

	// create a user
	user := data.User{
		Email:     r.Form.Get("email"),
		FirstName: r.Form.Get("first-name"),
		LastName:  r.Form.Get("last-name"),
		Password:  r.Form.Get("password"),
		Active:    0,
		IsAdmin:   0,
	}

	// insert user into db
	_, err = app.Models.User.Insert(user)
	if err != nil {
		app.Session.Put(r.Context(), "error", "Failed to create user")
		app.InfoLog.Println(err)
		http.Redirect(w, r, "/register", http.StatusSeeOther)
	}

	// send an activation email
	url := fmt.Sprintf("%s/activate?email=%s", os.Getenv("DOMAIN"), user.Email)

	signedURL := GenerateTokenFromString(url)
	app.InfoLog.Println(signedURL)

	msg := Message{
		To:       user.Email,
		Subject:  "Activate your account",
		Template: "confirmation-email",
		Data:     template.HTML(signedURL),
	}

	app.sendEmail(msg)

	app.Session.Put(r.Context(), "flash", "You've been registered successfully")
	// subscribe the user to an account

}

func (app *Config) ActivateAccount(w http.ResponseWriter, r *http.Request) {

	// validate url
	url := r.RequestURI
	testURL := fmt.Sprintf("%s%s", os.Getenv("DOMAIN"), url)
	okay := VerifyToken(testURL)

	if !okay {
		app.Session.Put(r.Context(), "error", "Invalid token")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// activate the user
	user, err := app.Models.User.GetByEmail(r.URL.Query().Get("email"))
	if err != nil {
		app.Session.Put(r.Context(), "error", "No user found")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	user.Active = 1
	err = app.Models.User.Update(*user)
	if err != nil {
		app.Session.Put(r.Context(), "error", "Failed to activate user")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	app.Session.Put(r.Context(), "flash", "Your account has been activated")
	http.Redirect(w, r, "/login", http.StatusSeeOther)

}

func (app *Config) SubscribeToPlan(w http.ResponseWriter, r *http.Request) {

	// get the id of the plan
	id := r.URL.Query().Get("id")

	planID, err := strconv.Atoi(id)
	if err != nil {
		app.Session.Put(r.Context(), "error", "Invalid plan")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	app.InfoLog.Println("Plan id is", planID)

	// get the plan from the db
	plan, err := app.Models.Plan.GetOne(planID)
	if err != nil {
		app.Session.Put(r.Context(), "error", "Invalid plan")
		http.Redirect(w, r, "/members/plans", http.StatusSeeOther)
		return
	}

	// get the user from the session
	user, ok := app.Session.Get(r.Context(), "user").(data.User)
	if !ok {
		app.Session.Put(r.Context(), "error", "Please login")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// generate an invoice and email it
	app.Wait.Add(1)

	go func() {
		defer app.Wait.Done()
		invoice, err := app.getInvoice(user, plan)
		if err != nil {
			// send this to a channel
			app.ErrorChan <- err
		}

		msg := Message{
			To:       user.Email,
			Subject:  "Your Invoice",
			Template: "invoice",
			Data:     invoice,
		}

		app.sendEmail(msg)
	}()

	// generate a manual
	app.Wait.Add(1)
	go func() {
		defer app.Wait.Done()

		pdf := app.generateManual(user, plan)

		err := pdf.OutputFileAndClose(fmt.Sprintf("/tmp/%d_manual.pdf", user.ID))
		if err != nil {
			app.ErrorChan <- err
		}

		msg := Message{
			To:       user.Email,
			Subject:  "Your Manual",
			Template: "manual",
			Data:     "Yous user manual is attached",
			AttachmentMap: map[string]string{
				"Manual.pdf": fmt.Sprintf("/tmp/%d_manual.pdf", user.ID),
			},
		}

		app.sendEmail(msg)

		//test app error chan

		app.ErrorChan <- errors.New("some custom error sent to error chan")
	}()

	// subscribe the user to the plan
	err = app.Models.Plan.SubscribeUserToPlan(user, *plan)
	if err != nil {
		app.Session.Put(r.Context(), "error", "Failed to subscribe you to the plan")
		http.Redirect(w, r, "/members/plans", http.StatusSeeOther)
		return
	}

	userFromDB, err := app.Models.User.GetOne(user.ID)
	if err != nil {
		app.Session.Put(r.Context(), "error", "Error getting user from DB")
		http.Redirect(w, r, "/members/plans", http.StatusSeeOther)
		return
	}

	app.Session.Put(r.Context(), "user", userFromDB)

	// redirect to the dashboard
	app.Session.Put(r.Context(), "flash", "You've been subscribed to the plan")
	http.Redirect(w, r, "/members/plans", http.StatusSeeOther)

}

// dummy function to generate an invoice
func (app *Config) getInvoice(user data.User, plan *data.Plan) (string, error) {

	app.InfoLog.Println("amount is", plan.PlanAmountFormatted)

	return plan.PlanAmountFormatted, nil
}

// dummy function to generaete the manual
func (app *Config) generateManual(user data.User, plan *data.Plan) *gofpdf.Fpdf {

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetMargins(10, 13, 10)

	importer := gofpdi.NewImporter()

	time.Sleep(3 * time.Second)

	templatePDF := importer.ImportPage(pdf, "pdf/manual.pdf", 1, "/MediaBox")

	pdf.AddPage()

	importer.UseImportedTemplate(pdf, templatePDF, 0, 0, 210, 297)

	pdf.SetX(75)
	pdf.SetY(150)
	pdf.SetFont("Arial", "", 12)

	pdf.MultiCell(0, 4, fmt.Sprintf("%s %s", user.FirstName, user.LastName), "", "C", false)

	pdf.Ln(5)

	pdf.MultiCell(0, 4, fmt.Sprintf("%s User Guide", plan.PlanName), "", "C", false)

	return pdf
}

func (app *Config) ChooseSubscription(w http.ResponseWriter, r *http.Request) {

	plans, err := app.Models.Plan.GetAll()
	if err != nil {
		app.ErrorLog.Println(err)
	}

	dataMap := make(map[string]any)

	dataMap["plans"] = plans

	app.render(w, r, "plans.page.gohtml", &TemplateData{
		Data: dataMap,
	})

}
