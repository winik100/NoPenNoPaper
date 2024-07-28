package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/winik100/NoPenNoPaper/internal/models"
	"github.com/winik100/NoPenNoPaper/internal/validators"
)

type characterCreateForm struct {
	Info                     models.CharacterInfo
	Attributes               models.CharacterAttributes
	Skills                   models.Skills
	CustomSkills             models.CustomSkills
	validators.FormValidator `schema:"-"`
}

type userForm struct {
	Name                     string
	Password                 string
	validators.FormValidator `schema:"-"`
}

type itemForm struct {
	CharacterId              int
	Name                     string
	Description              string
	Count                    int
	validators.FormValidator `schema:"-"`
}

type noteForm struct {
	Text                     string
	validators.FormValidator `schema:"-"`
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	userId := app.sessionManager.GetInt(r.Context(), string(authenticatedUserIdContextKey))
	if userId == 0 {
		data := app.newTemplateData(r)
		w.WriteHeader(http.StatusOK)
		app.render(w, r, "home.tmpl.html", data)
		return
	}

	role := app.sessionManager.GetString(r.Context(), "role")
	data := app.newTemplateData(r)
	if role == RoleGM {
		characters, err := app.characters.GetAll()
		if err != nil {
			app.serverError(w, r, err)
			return
		}
		data.Characters = characters
	} else {
		characters, err := app.characters.GetAllFrom(userId)
		if err != nil {
			app.serverError(w, r, err)
			return
		}
		data.Characters = characters
	}

	w.WriteHeader(http.StatusOK)
	app.render(w, r, "home.tmpl.html", data)
}

func (app *application) create(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	skills, err := app.characters.GetAvailableSkills()
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	data.Form = characterCreateForm{Skills: skills}

	w.WriteHeader(http.StatusOK)
	app.render(w, r, "create.tmpl.html", data)
}

func (app *application) createPost(w http.ResponseWriter, r *http.Request) {
	var form characterCreateForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusUnprocessableEntity)
		return
	}

	form.InfoChecks()
	form.AttributeChecks()

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		w.WriteHeader(http.StatusUnprocessableEntity)
		app.render(w, r, "create.tmpl.html", data)
		return
	}

	_, err = app.characters.Insert(models.Character{Info: form.Info, Attributes: form.Attributes, Skills: form.Skills, CustomSkills: form.CustomSkills},
		app.sessionManager.GetInt(r.Context(), string(authenticatedUserIdContextKey)))

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) signup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userForm{}
	w.WriteHeader(http.StatusOK)
	app.render(w, r, "signup.tmpl.html", data)
}

func (app *application) signupPost(w http.ResponseWriter, r *http.Request) {
	var form userForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validators.NotBlank(form.Name), "Name", "Dieses Feld kann nicht leer sein.")
	form.CheckField(validators.NotBlank(form.Password), "Password", "Dieses Feld kann nicht leer sein.")
	form.CheckField(validators.MinChars(form.Password, 8), "Password", "Passwort muss mindestens 8 Zeichen lang sein.")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		w.WriteHeader(http.StatusUnprocessableEntity)
		app.render(w, r, "signup.tmpl.html", data)
		return
	}

	err = app.users.Insert(form.Name, form.Password)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Erfolgreich registriert! Bitte einloggen.")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userForm{}
	w.WriteHeader(http.StatusOK)
	app.render(w, r, "login.tmpl.html", data)
}

func (app *application) loginPost(w http.ResponseWriter, r *http.Request) {
	var form userForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validators.NotBlank(form.Name), "Name", "Dieses Feld kann nicht leer sein.")
	form.CheckField(validators.NotBlank(form.Password), "Password", "Dieses Feld kann nicht leer sein.")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		w.WriteHeader(http.StatusUnprocessableEntity)
		app.render(w, r, "login.tmpl.html", data)
		return
	}

	id, err := app.users.Authenticate(form.Name, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddGenericError("Name und/oder Password sind falsch.")

			data := app.newTemplateData(r)
			data.Form = form
			w.WriteHeader(http.StatusUnprocessableEntity)
			app.render(w, r, "login.tmpl.html", data)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	app.sessionManager.Put(r.Context(), string(authenticatedUserIdContextKey), id)

	role, err := app.users.GetRole(id)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	app.sessionManager.Put(r.Context(), "role", role)
	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) logoutPost(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Remove(r.Context(), string(authenticatedUserIdContextKey))
	app.sessionManager.Put(r.Context(), "role", RoleAnon)
	app.sessionManager.Put(r.Context(), "flash", "Erfolgreich ausgeloggt!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) viewCharacter(w http.ResponseWriter, r *http.Request) {
	characterId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	character, err := app.characters.Get(characterId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Character = character
	app.sessionManager.Put(r.Context(), "characterId", characterId)
	w.WriteHeader(http.StatusOK)
	app.render(w, r, "character.tmpl.html", data)
}

func (app *application) Inc(w http.ResponseWriter, r *http.Request) {
	characterId := app.sessionManager.GetInt(r.Context(), "characterId")
	if characterId == 0 {
		http.NotFound(w, r)
		return
	}

	character, err := app.characters.Get(characterId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	stat := r.FormValue("inc")
	updated, err := app.characters.IncrementStat(character, stat)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	updatedStat := updated.Stats.CurrentAsMap()[stat]
	tmplStr := `<div id="{{.Stat}}" value="{{.Value}}">{{.Value}}</div>`

	t, err := template.New("inc").Parse(tmplStr)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	data := map[string]string{
		"Stat":  stat,
		"Value": strconv.Itoa(updatedStat),
	}
	w.WriteHeader(http.StatusOK)
	t.ExecuteTemplate(w, "inc", data)
}

func (app *application) Dec(w http.ResponseWriter, r *http.Request) {
	characterId := app.sessionManager.GetInt(r.Context(), "characterId")
	if characterId == 0 {
		http.NotFound(w, r)
		return
	}

	character, err := app.characters.Get(characterId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	stat := r.FormValue("dec")
	updated, err := app.characters.DecrementStat(character, stat)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	updatedStat := updated.Stats.CurrentAsMap()[stat]
	tmplStr := `<div id="{{.Stat}}" name="Stat" value="{{.Value}}">{{.Value}}</div>`

	t, err := template.New("dec").Parse(tmplStr)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	data := map[string]string{
		"Stat":  stat,
		"Value": strconv.Itoa(updatedStat),
	}
	w.WriteHeader(http.StatusOK)
	t.ExecuteTemplate(w, "dec", data)
}

func (app *application) addItem(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = itemForm{
		CharacterId: app.sessionManager.GetInt(r.Context(), "characterId"),
	}
	w.WriteHeader(http.StatusOK)
	app.render(w, r, "item.tmpl.html", data)
}

func (app *application) addItemPost(w http.ResponseWriter, r *http.Request) {
	characterId := app.sessionManager.GetInt(r.Context(), "characterId")
	if characterId == 0 {
		http.NotFound(w, r)
		return
	}

	var form itemForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusUnprocessableEntity)
		return
	}

	form.CheckField(validators.NotBlank(form.Name), "Name", "Dieses Feld kann nicht leer sein.")
	form.CheckField(validators.MaxChars(form.Name, 50), "Name", "Maximal 50 Zeichen erlaubt.")
	form.CheckField(validators.NotBlank(form.Description), "Description", "Dieses Feld kann nicht leer sein.")
	form.CheckField(validators.MaxChars(form.Description, 255), "Description", "Maximal 255 Zeichen erlaubt.")
	form.CheckField(form.Count > 0, "Count", "Die Anzahl muss positiv sein.")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		w.WriteHeader(http.StatusUnprocessableEntity)
		app.render(w, r, "item.tmpl.html", data)
		return
	}

	err = app.characters.AddItem(characterId, form.Name, form.Description, form.Count)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	redirect := fmt.Sprintf("/characters/%d", characterId)
	http.Redirect(w, r, redirect, http.StatusSeeOther)
}

func (app *application) deleteItemPost(w http.ResponseWriter, r *http.Request) {
	type deleteForm struct {
		ItemId int
	}

	var form deleteForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusUnprocessableEntity)
		return
	}

	err = app.characters.DeleteItem(form.ItemId)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	t, err := template.New("empty").Parse("")
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	t.ExecuteTemplate(w, "empty", form)
}

func (app *application) addNote(w http.ResponseWriter, r *http.Request) {
	characterId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	character, err := app.characters.Get(characterId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	tmplStr := `<form hx-post="/characters/{{.Character.ID}}/addNote" hx-target="this" hx-swap="outerHTML">
					<input type="hidden" name="csrf_token" value="{{.CSRFToken}}">
					<div>
						<label>Notiz:</label>
						<input type="text" name="Text" textarea>
					</div>
					<button type="submit">Hinzufügen</button>
					<button hx-get="/characters/{{.Character.ID}}">Abbrechen</button>
				</form>`

	t, err := template.New("addnote").Parse(tmplStr)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	data := app.newTemplateData(r)
	data.Character = character
	data.Form = noteForm{}
	w.WriteHeader(http.StatusOK)
	t.ExecuteTemplate(w, "addnote", data)
}

func (app *application) addNotePost(w http.ResponseWriter, r *http.Request) {
	characterId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	var form noteForm
	err = app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusUnprocessableEntity)
		return
	}
	noteId, err := app.characters.AddNote(characterId, form.Text)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	tmplStr := fmt.Sprintf(`<div id="note" hx-target="this" hx-swap="outerHTML">
                	<button hx-get="/characters/{{.Character.ID}}/addNote">Notiz hinzufügen</button>
            	</div>
				<form id="deleteNote" hx-post="/characters/{{.Character.ID}}/deleteNote" hx-target="this" hx-swap="outerHTML">
                <input type="hidden" name="csrf_token" value="{{.CSRFToken}}">
                <input type="hidden" name="NoteId" Value="%d">
                <li>%s    <button type="submit">löschen</button></li>
            </form>`, noteId, form.Text)

	t, err := template.New("button").Parse(tmplStr)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	character, err := app.characters.Get(characterId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Character = character
	w.WriteHeader(http.StatusOK)
	t.ExecuteTemplate(w, "button", data)
}

func (app *application) deleteNotePost(w http.ResponseWriter, r *http.Request) {
	type deleteForm struct {
		NoteId int
	}

	var form deleteForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusUnprocessableEntity)
		return
	}

	err = app.characters.DeleteNote(form.NoteId)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	t, err := template.New("empty").Parse("")
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	t.ExecuteTemplate(w, "empty", nil)
}

func (app *application) customSkillInput(w http.ResponseWriter, r *http.Request) {
	tmplStr := `<tr id="{{.Category}}">
					<td>
						<input type='hidden' name='CustomSkills.Category' value='{{.Category}}'>
						<label>{{.Category}}</label>
						<input type="text" name="CustomSkills.Name">
						<select name="CustomSkills.Value">
							<option value="{{.Default}}" selected>{{.Default}}</option>
							<option value="70">70</option>
							<option value="60">60</option>
							<option value="50">50</option>
							<option value="40">40</option>
						</select>
						<button id="cancel" hx-get="/cancel" hx-target="#{{.Category}}" hx-swap="outerHTML">Abbrechen</button>
					</td>
				</tr>`

	t, err := template.New("customSkillInput").Parse(tmplStr)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	category := r.URL.Query().Get("category")
	defaultValue := models.DefaultForCategory(category)

	data := map[string]string{
		"Category": category,
		"Default":  strconv.Itoa(defaultValue),
	}
	w.WriteHeader(http.StatusOK)
	t.ExecuteTemplate(w, "customSkillInput", data)
}

func (app *application) cancel(w http.ResponseWriter, r *http.Request) {
	t, err := template.New("cancel").Parse("")
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	t.ExecuteTemplate(w, "cancel", nil)
}
