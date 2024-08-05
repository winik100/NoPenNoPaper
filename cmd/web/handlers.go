package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"slices"
	"strconv"

	"github.com/justinas/nosurf"
	"github.com/winik100/NoPenNoPaper/internal/models"
	"github.com/winik100/NoPenNoPaper/internal/validators"
)

type characterCreateForm struct {
	Info                     models.CharacterInfo
	Attributes               models.CharacterAttributes
	Skills                   models.Skills
	SelectedSkills           []string
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

type itemEditForm struct {
	ItemId                   int
	Count                    int
	Direction                string
	validators.FormValidator `schema:"-"`
}

type noteForm struct {
	Text                     string
	validators.FormValidator `schema:"-"`
}

type skillEditForm struct {
	ID                       int
	Skill                    string
	NewValue                 int
	validators.FormValidator `schema:"-"`
}

type skillAddForm struct {
	ID           int
	AddableSkill string
	Value        int
}

type customSkillAddForm struct {
	ID                       int
	CustomSkill              string
	Category                 string
	Value                    int
	validators.FormValidator `schema:"-"`
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	userId := app.sessionManager.GetInt(r.Context(), authenticatedUserIdKey)
	if userId == 0 {
		data := app.newTemplateData(r)
		w.WriteHeader(http.StatusOK)
		app.render(w, r, "home.tmpl.html", data)
		return
	}

	role := app.sessionManager.GetString(r.Context(), roleKey)
	data := app.newTemplateData(r)
	if role == "gm" {
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

func (app *application) delete(w http.ResponseWriter, r *http.Request) {
	characterId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	tmplStr := `<form id="deleteCharacter" action="/characters/{{.Character.ID}}/delete" method="POST">
					<p>Sicher? Kann nicht rückgängig gemacht werden!</p>
					<input type="hidden" name="csrf_token" value="{{.CSRFToken}}">
					<input type="hidden" name="CharacterId" Value="{{.Character.ID}}">
					<button type="submit">OK</button>
					<button hx-get="/characters/{{.Character.ID}}" hx-target="#deleteCharacter" hx-select="#deleteCharacter">Abbrechen</button>
            	</form>`

	t, err := template.New("delete").Parse(tmplStr)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	character, err := app.characters.Get(characterId)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Character = character
	w.WriteHeader(http.StatusOK)
	t.ExecuteTemplate(w, "delete", data)
}

func (app *application) deletePost(w http.ResponseWriter, r *http.Request) {
	type deleteForm struct {
		CharacterId int
	}

	var form deleteForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusUnprocessableEntity)
		return
	}

	err = app.characters.Delete(form.CharacterId)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
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
		availableSkills, err := app.characters.GetAvailableSkills()
		if err != nil {
			app.serverError(w, r, err)
			return
		}
		form.Skills = mergeSkills(availableSkills, form.Skills)
		data.Form = form
		w.WriteHeader(http.StatusUnprocessableEntity)
		app.render(w, r, "create.tmpl.html", data)
		return
	}

	_, err = app.characters.Insert(models.Character{Info: form.Info, Attributes: form.Attributes, Skills: form.Skills, CustomSkills: form.CustomSkills},
		app.sessionManager.GetInt(r.Context(), authenticatedUserIdKey))

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

	_, err = app.users.Get(form.Name)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			err = app.users.Insert(form.Name, form.Password)
			if err != nil {
				app.serverError(w, r, err)
				return
			}
			app.sessionManager.Put(r.Context(), "flash", "Erfolgreich registriert! Bitte einloggen.")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		app.serverError(w, r, err)
		return
	}

	form.CheckField(false, "Name", "Dieser Name ist bereits vergeben.")
	data := app.newTemplateData(r)
	data.Form = form
	w.WriteHeader(http.StatusUnprocessableEntity)
	app.render(w, r, "signup.tmpl.html", data)
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
	app.sessionManager.Put(r.Context(), authenticatedUserIdKey, id)

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

	app.sessionManager.Remove(r.Context(), authenticatedUserIdKey)
	app.sessionManager.Remove(r.Context(), isAuthenticatedKey)
	app.sessionManager.Put(r.Context(), roleKey, "anonymous")
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

func (app *application) addSkill(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	character, err := app.characters.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}
	allSkills, err := app.characters.GetAvailableSkills()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	var addableSkills models.Skills
	for i, sk := range allSkills.Name {
		if !slices.Contains(character.Skills.Name, sk) {
			addableSkills.Name = append(addableSkills.Name, sk)
			addableSkills.Value = append(addableSkills.Value, allSkills.Value[i])
		}
	}

	tmplStr := `<form id="addSkillForm" hx-post="/characters/{{.ID}}/addSkill" hx-target="this" hx-swap="outerHTML">
				<input type="hidden" name="csrf_token" value="{{.CSRFToken}}">
				<input type="hidden" name="ID" value="{{.ID}}">
				<select name='AddableSkill'>
					{{range .AddableSkills.Name}}
						<option value='{{.}}'>{{.}}</option>
					{{end}}
				</select>
				<input type="number" name="Value">
				<button type="submit">OK</button>
				<button hx-get="/characters/{{.ID}}" hx-target="#addSkillForm" hx-swap="outerHTML" hx-select="#addSkill">Abbrechen</button>
				</form>`

	data := map[string]any{
		"ID":            id,
		"AddableSkills": addableSkills,
		"CSRFToken":     nosurf.Token(r),
	}

	t, err := template.New("addSkill").Parse(tmplStr)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	t.ExecuteTemplate(w, "addSkill", data)
}

func (app *application) addSkillPost(w http.ResponseWriter, r *http.Request) {
	var form skillAddForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	err = app.characters.AddSkill(form.ID, form.AddableSkill, form.Value)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	half := form.Value / 2
	fifth := form.Value / 5
	tmplStr := fmt.Sprintf(`<template>
							<tr hx-swap-oob="beforeend:#Skills">
								<th>{{.Form.AddableSkill}}</th>
								<td>
									<div id="Values{{.Form.AddableSkill}}">{{.Form.Value}} | %d | %d</div>
									<form id="edit{{.Form.AddableSkill}}" hx-get="/characters/{{.Form.ID}}/editSkill" hx-target="this" hx-swap="outerHTML">
										<input type="hidden" name="skill" value="{{.Form.AddableSkill}}">
										<input type="hidden" name="value" value="{{.Form.Value}}">
										<button type="submit">Bearbeiten</button>
									</form>
								</td>
							</tr>
							</template>
							<div id="addSkill" hx-target="this" hx-swap="outerHTML">
								<button hx-get="/characters/{{.Form.ID}}/addSkill">Fertigkeit hinzufügen</button>
							</div>`, half, fifth)

	t, err := template.New("addSkillDone").Parse(tmplStr)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Form = form
	w.WriteHeader(http.StatusOK)
	t.ExecuteTemplate(w, "addSkillDone", data)
}

func (app *application) editSkill(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	id := r.PathValue("id")
	skill := params.Get("skill")
	value, err := strconv.Atoi(params.Get("value"))
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	tmplStr := fmt.Sprintf(`<form id="editForm" hx-post="/characters/%s/editSkill" hx-target="this" hx-swap="outerHTML">
                <input type="hidden" name="csrf_token" value="{{.CSRFToken}}">
				<input type="hidden" name="Skill" value="%s">
                <input type="number" name="NewValue" value="%d">
				<button type="submit">OK</button>
				<button hx-get="/characters/%s" hx-target="#editForm" hx-swap="outerHTML" hx-select="#edit%s">Abbrechen</button>
            </form>`, id, skill, value, id, skill)

	t, err := template.New("editForm").Parse(tmplStr)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	w.WriteHeader(http.StatusOK)
	t.ExecuteTemplate(w, "editForm", data)
}

func (app *application) editSkillPost(w http.ResponseWriter, r *http.Request) {
	characterId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	var form skillEditForm
	err = app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	err = app.characters.EditSkill(characterId, form.Skill, form.NewValue)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	half := form.NewValue / 2
	fifth := form.NewValue / 5
	tmplStr := fmt.Sprintf(`<div hx-swap-oob="outerHTML:#Values{{.Form.Skill}}">{{.Form.NewValue}} | %d | %d</div>
							<form hx-get="/characters/{{.Form.ID}}/editSkill" hx-target="this" hx-swap="outerHTML">	
                            	<input type="hidden" name="skill" value="{{.Form.Skill}}">
                            	<input type="hidden" name="value" value="{{.Form.NewValue}}">
                            	<button type="submit">Bearbeiten</button>
                        	</form>`, half, fifth)

	t, err := template.New("editSkillDone").Parse(tmplStr)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Form = form
	w.WriteHeader(http.StatusOK)
	t.ExecuteTemplate(w, "editSkillDone", data)
}

func (app *application) addCustomSkill(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	tmplStr := `<form id="addCustomSkillForm" hx-post="/characters/{{.ID}}/addCustomSkill" hx-target="this" hx-swap="outerHTML">
					<input type="hidden" name="csrf_token" value="{{.CSRFToken}}">
					<input type="hidden" name="ID" value="{{.ID}}">
					<select name='Category'>
								<option value='' disabled selected>Wähle Kategorie</option>
								<option value='Muttersprache'>Muttersprache</option>
								<option value='Fremdsprache'>Fremdsprache</option>
								<option value='Handwerk'>Handwerk und Kunst</option>
								<option value='Naturwissenschaft'>Naturwissenschaft</option>
								<option value='Steuern'>Steuern</option>
								<option value='Überlebenskunst'>Überlebenskunst</option>
								<option value='Sonstiges'>Sonstiges</option>
							</select>
					<input type="text" name="CustomSkill">
					<input type="number" name="Value">
					<button type="submit">OK</button>
					<button hx-get="/characters/{{.ID}}" hx-target="#addCustomSkillForm" hx-swap="outerHTML" hx-select="#addCustomSkill">Abbrechen</button>
				</form>`

	data := map[string]any{
		"ID":        id,
		"CSRFToken": nosurf.Token(r),
	}

	t, err := template.New("addCustomSkill").Parse(tmplStr)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	t.ExecuteTemplate(w, "addCustomSkill", data)
}

func (app *application) addCustomSkillPost(w http.ResponseWriter, r *http.Request) {
	var form customSkillAddForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validators.NotBlank(form.CustomSkill), "Name", "Dieses Feld kann nicht leer sein.")
	form.CheckField(form.Value != 0, "Value", "Dieses Feld muss einen positiven Wert enthalten.")
	form.CheckField(models.DefaultForCategory(form.Category) != -1, "Category", "Es muss eine gültige Kategorie gewählt werden.")

	if !form.Valid() {
		tmplStr := `<form id="addCustomSkillForm" hx-post="/characters/{{.Form.ID}}/addCustomSkill" hx-target="this" hx-swap="outerHTML">
						<input type="hidden" name="csrf_token" value="{{.CSRFToken}}">
						<input type="hidden" name="ID" value="{{.Form.ID}}">
						<select name='Category'>
							<option value='' disabled selected>Wähle Kategorie</option>
							<option value='Muttersprache'>Muttersprache</option>
							<option value='Fremdsprache'>Fremdsprache</option>
							<option value='Handwerk'>Handwerk und Kunst</option>
							<option value='Naturwissenschaft'>Naturwissenschaft</option>
							<option value='Steuern'>Steuern</option>
							<option value='Überlebenskunst'>Überlebenskunst</option>
							<option value='Sonstiges'>Sonstiges</option>
						</select>
						{{with .Form.FieldErrors.Category}}<label class='error'>{{.}}</label>{{end}}
						<input type="text" name="CustomSkill">
						{{with .Form.FieldErrors.Name}}<label class='error'>{{.}}</label>{{end}}
						<input type="number" name="Value">
						{{with .Form.FieldErrors.Value}}<label class='error'>{{.}}</label>{{end}}
						<button type="submit">OK</button>
						<button hx-get="/characters/{{.Form.ID}}" hx-target="#addCustomSkillForm" hx-swap="outerHTML" hx-select="#addCustomSkill">Abbrechen</button>
					</form>`
		t, err := template.New("addCustomSkillFailed").Parse(tmplStr)
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		data := app.newTemplateData(r)
		data.Form = form
		w.WriteHeader(http.StatusOK)
		t.ExecuteTemplate(w, "addCustomSkillFailed", data)
		return
	}

	err = app.characters.AddCustomSkill(form.ID, form.CustomSkill, form.Category, form.Value)
	if err != nil {
		if errors.Is(err, models.ErrAlreadyHasSkill) {
			tmplStr := `<div id="addCustomSkill" hx-target="this" hx-swap="outerHTML">
                			<button hx-get="/characters/{{.Form.ID}}/addCustomSkill">Fertigkeit hinzufügen</button>
							<label class="error">Der Charaktere verfügt bereits über eine gleichnamige Fertigkeit.</label>
            			</div>`
			t, err := template.New("addCustomSkillDuplicate").Parse(tmplStr)
			if err != nil {
				app.serverError(w, r, err)
				return
			}

			data := app.newTemplateData(r)
			data.Form = form
			w.WriteHeader(http.StatusOK)
			t.ExecuteTemplate(w, "addCustomSkillDuplicate", data)
			return
		}
		app.serverError(w, r, err)
		return
	}

	half := form.Value / 2
	fifth := form.Value / 5
	tmplStr := fmt.Sprintf(`<template>
							<tr hx-swap-oob="beforeend:#CustomSkills">
								<th>{{.Form.CustomSkill}}</th>
								<td>
									<div id="Values{{.Form.CustomSkill}}" value="{{.Form.Value}}">{{.Form.Value}} | %d | %d</div>
									<form id="edit{{.Form.CustomSkill}}" hx-get="/characters/{{.Form.ID}}/editCustomSkill" hx-target="this" hx-swap="outerHTML">
										<input type="hidden" name="skill" value="{{.Form.CustomSkill}}">
										<input type="hidden" name="value" value="{{.Form.Value}}">
										<button type="submit">Bearbeiten</button>
									</form>
								</td>
							</tr>
							</template>
							<div id="addCustomSkill" hx-target="this" hx-swap="outerHTML">
								<button hx-get="/characters/{{.Form.ID}}/addCustomSkill">Fertigkeit hinzufügen</button>
							</div>`, half, fifth)

	t, err := template.New("addCustomSkillDone").Parse(tmplStr)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Form = form
	w.WriteHeader(http.StatusOK)
	t.ExecuteTemplate(w, "addCustomSkillDone", data)
}

func (app *application) editCustomSkill(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	id := r.PathValue("id")
	skill := params.Get("skill")
	value, err := strconv.Atoi(params.Get("value"))
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	tmplStr := fmt.Sprintf(`<form id="editForm" hx-post="/characters/%s/editCustomSkill" hx-target="this" hx-swap="outerHTML">
                <input type="hidden" name="csrf_token" value="{{.CSRFToken}}">
				<input type="hidden" name="ID" value="%s">
				<input type="hidden" name="Skill" value="%s">
                <input type="number" name="NewValue" value="%d">
				<button type="submit">OK</button>
				<button hx-get="/characters/%s" hx-target="#editForm" hx-swap="outerHTML" hx-select="#edit%s">Abbrechen</button>
            </form>`, id, id, skill, value, id, skill)

	t, err := template.New("editForm").Parse(tmplStr)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	w.WriteHeader(http.StatusOK)
	t.ExecuteTemplate(w, "editForm", data)
}

func (app *application) editCustomSkillPost(w http.ResponseWriter, r *http.Request) {
	characterId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	var form skillEditForm
	err = app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	half := form.NewValue / 2
	fifth := form.NewValue / 5
	tmplStr := fmt.Sprintf(`<div value="{{.Form.NewValue}}" hx-swap-oob="outerHTML:#Values{{.Form.Skill}}">{{.Form.NewValue}} | %d | %d</div>
							<form hx-get="/characters/{{.Form.ID}}/editCustomSkill" hx-target="this" hx-swap="outerHTML">	
                            	<input type="hidden" name="skill" value="{{.Form.Skill}}">
                            	<input type="hidden" name="value" value="{{.Form.NewValue}}">
                            	<button type="submit">Bearbeiten</button>
                        	</form>`, half, fifth)

	t, err := template.New("editCustomSkillDone").Parse(tmplStr)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = app.characters.EditCustomSkill(characterId, form.Skill, form.NewValue)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Form = form
	w.WriteHeader(http.StatusOK)
	t.ExecuteTemplate(w, "editCustomSkillDone", data)
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
						<button hx-get="/create" hx-target="#{{.Category}}" hx-swap="delete">Abbrechen</button>
					</td>
				</tr>`

	t, err := template.New("customSkillInput").Parse(tmplStr)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	category := r.URL.Query().Get("category")
	defaultValue := models.DefaultForCategory(category)

	data := map[string]any{
		"Category": category,
		"Default":  defaultValue,
	}
	w.WriteHeader(http.StatusOK)
	t.ExecuteTemplate(w, "customSkillInput", data)
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
	characterId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	data := app.newTemplateData(r)
	data.Form = itemForm{
		CharacterId: characterId,
	}
	w.WriteHeader(http.StatusOK)
	app.render(w, r, "item.tmpl.html", data)
}

func (app *application) addItemPost(w http.ResponseWriter, r *http.Request) {
	characterId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	var form itemForm
	err = app.decodePostForm(r, &form)
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

func (app *application) editItemCount(w http.ResponseWriter, r *http.Request) {
	characterId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	var form itemEditForm
	err = app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusUnprocessableEntity)
		return
	}

	tmplStr := `<div id="itemCount">
					{{if gt .Form.NewCount 1}}
                    <button type="submit" name="Direction" value="dec">-</button>
                    {{end}}
					<input type="hidden" name="ItemId" value="{{.Form.ItemId}}">
					<input type="hidden" name="Count" value="{{.Form.NewCount}}">
					{{.Form.NewCount}}
					<button type="submit" name="Direction" value="inc">+</button>
				</div>`

	t, err := template.New("editCountDone").Parse(tmplStr)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	var tmplData map[string]int
	switch form.Direction {
	case "inc":
		err = app.characters.EditItemCount(characterId, form.ItemId, form.Count+1)
		tmplData = map[string]int{
			"ItemId":   form.ItemId,
			"NewCount": form.Count + 1,
		}
	case "dec":
		err = app.characters.EditItemCount(characterId, form.ItemId, form.Count-1)
		tmplData = map[string]int{
			"ItemId":   form.ItemId,
			"NewCount": form.Count - 1,
		}
	}
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Form = tmplData
	w.WriteHeader(http.StatusOK)
	t.ExecuteTemplate(w, "editCountDone", data)
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

	tmplStr := `<form id="addNoteForm" hx-post="/characters/{{.Character.ID}}/addNote" hx-target="this" hx-swap="outerHTML">
					<input type="hidden" name="csrf_token" value="{{.CSRFToken}}">
					<div>
						<label>Notiz:</label>
						<input type="text" name="Text" textarea>
					</div>
					<button type="submit">Hinzufügen</button>
					<button hx-get="/characters/{{.Character.ID}}" hx-target="#addNoteForm" hx-swap="delete">Abbrechen</button>
				</form>`

	t, err := template.New("addNote").Parse(tmplStr)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	data := app.newTemplateData(r)
	data.Character = character
	data.Form = noteForm{}
	w.WriteHeader(http.StatusOK)
	t.ExecuteTemplate(w, "addNote", data)
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

	tmplStr := fmt.Sprintf(`<form id="deleteNote" hx-post="/characters/{{.Character.ID}}/deleteNote" hx-target="this" hx-swap="outerHTML">
								<input type="hidden" name="csrf_token" value="{{.CSRFToken}}">
								<input type="hidden" name="NoteId" Value="%d">
								<li>{{.Form.Text}}    <button type="submit">löschen</button></li>
							</form>`, noteId)

	t, err := template.New("addNoteDone").Parse(tmplStr)
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
	data.Form = form
	w.WriteHeader(http.StatusOK)
	t.ExecuteTemplate(w, "addNoteDone", data)
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

	t, err := template.New("deleteNoteDone").Parse("")
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	t.ExecuteTemplate(w, "deleteNoteDone", nil)
}
