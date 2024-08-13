package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strconv"

	"github.com/winik100/NoPenNoPaper/internal/core"
	"github.com/winik100/NoPenNoPaper/internal/models"
	"github.com/winik100/NoPenNoPaper/internal/validators"
)

type characterCreateForm struct {
	Info                     core.CharacterInfo
	Attributes               core.CharacterAttributes
	Skills                   core.Skills
	SelectedSkills           []string
	CustomSkills             core.CustomSkills
	validators.FormValidator `schema:"-"`
}

type userForm struct {
	Name                     string
	Password                 string
	validators.FormValidator `schema:"-"`
}

type statEditForm struct {
	Name                     string
	Value                    int
	Direction                string
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
	CharacterId              int
	Text                     string
	validators.FormValidator `schema:"-"`
}

type skillEditForm struct {
	CharacterId              int
	Skill                    string
	NewValue                 int
	validators.FormValidator `schema:"-"`
}

type skillAddForm struct {
	CharacterId              int
	AddableSkill             string
	Value                    int
	validators.FormValidator `schema:"-"`
}

type customSkillAddForm struct {
	CharacterId              int
	CustomSkill              string
	Category                 string
	Value                    int
	validators.FormValidator `schema:"-"`
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	userId := app.sessionManager.GetInt(r.Context(), authenticatedUserIdKey)
	userName := app.sessionManager.GetString(r.Context(), authenticatedUserNameKey)
	if userId == 0 {
		data := app.newTemplateData(r)
		w.WriteHeader(http.StatusOK)
		app.render(w, r, "home.tmpl.html", data)
		return
	}

	role := app.sessionManager.GetString(r.Context(), roleKey)
	data := app.newTemplateData(r)
	data.User = core.User{ID: userId, Name: userName}
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
	characterId := app.sessionManager.GetInt(r.Context(), characterIdKey)

	tmplStr := `<form id="deleteCharacterForm" action="/characters/{{.Form.CharacterId}}/delete" method="POST">
					<p id="deleteCharacterMessage">Sicher? Kann nicht rückgängig gemacht werden!</p>
					<input type="hidden" name="csrf_token" value="{{.CSRFToken}}">
					<input type="hidden" name="CharacterId" Value="{{.Form.CharacterId}}">
					<button type="submit">OK</button>
					<button hx-get="/characters/{{.Form.CharacterId}}" hx-target="#deleteCharacterForm" hx-select="#deleteCharacter" hx-swap="outerHTML">Abbrechen</button>
            	</form>`

	data := app.newTemplateData(r)
	data.Form = map[string]any{
		"CharacterId": characterId,
	}
	w.WriteHeader(http.StatusOK)
	app.renderHtmx(w, r, "delete", tmplStr, data)
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
		form.Skills = core.MergeSkills(availableSkills, form.Skills)
		data.Form = form
		w.WriteHeader(http.StatusUnprocessableEntity)
		app.render(w, r, "create.tmpl.html", data)
		return
	}

	_, err = app.characters.Insert(core.Character{Info: form.Info, Attributes: form.Attributes, Skills: form.Skills, CustomSkills: form.CustomSkills},
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

	exists, err := app.users.Exists(form.Name)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	if exists {
		form.CheckField(false, "Name", "Dieser Name ist bereits vergeben.")
		data := app.newTemplateData(r)
		data.Form = form
		w.WriteHeader(http.StatusUnprocessableEntity)
		app.render(w, r, "signup.tmpl.html", data)
		return
	}

	_, err = app.users.Insert(form.Name, form.Password)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
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
	app.sessionManager.Put(r.Context(), authenticatedUserIdKey, id)
	app.sessionManager.Put(r.Context(), authenticatedUserNameKey, form.Name)

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
	app.sessionManager.Remove(r.Context(), authenticatedUserNameKey)
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
	app.sessionManager.Put(r.Context(), characterIdKey, characterId)
	w.WriteHeader(http.StatusOK)
	app.render(w, r, "character.tmpl.html", data)
}

func (app *application) addSkill(w http.ResponseWriter, r *http.Request) {
	characterId := app.sessionManager.GetInt(r.Context(), characterIdKey)
	character, err := app.characters.Get(characterId)
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

	var addableSkills core.Skills
	for i, sk := range allSkills.Name {
		if !slices.Contains(character.Skills.Name, sk) {
			addableSkills.Name = append(addableSkills.Name, sk)
			addableSkills.Value = append(addableSkills.Value, allSkills.Value[i])
		}
	}

	tmplStr := `<form id="addSkillForm" hx-post="/characters/{{.Form.CharacterId}}/addSkill" hx-target="this" hx-swap="outerHTML">
				<input type="hidden" name="csrf_token" value="{{.CSRFToken}}">
				<input type="hidden" name="CharacterId" value="{{.Form.CharacterId}}">
				<select name='AddableSkill'>
					{{range .Form.AddableSkills.Name}}
						<option value='{{.}}'>{{.}}</option>
					{{end}}
				</select><br>
				<input type="number" name="Value"><br>
				<button type="submit">OK</button>
				<button hx-get="/characters/{{.Form.CharacterId}}" hx-target="#addSkillForm" hx-swap="outerHTML" hx-select="#addSkill">Abbrechen</button>
				</form>`

	data := app.newTemplateData(r)
	data.Form = map[string]any{
		"CharacterId":   characterId,
		"AddableSkills": addableSkills,
	}

	w.WriteHeader(http.StatusOK)
	app.renderHtmx(w, r, "addSkill", tmplStr, data)
}

func (app *application) addSkillPost(w http.ResponseWriter, r *http.Request) {
	var form skillAddForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(form.Value != 0, "Value", "Dieses Feld muss einen positiven Wert enthalten.")

	if !form.Valid() {
		tmplStr := `<form id="addSkillForm" hx-post="/characters/{{.Form.CharacterId}}/addSkill" hx-target="this" hx-swap="outerHTML">
						<input type="hidden" name="csrf_token" value="{{.CSRFToken}}">
						<input type="hidden" name="CharacterId" value="{{.Form.CharacterId}}">
						<select name='AddableSkill'>
							{{range .AdditionalData.AddableSkills.Name}}
								<option value='{{.}}'>{{.}}</option>
							{{end}}
						</select><br>
						<input type="number" name="Value"><br>
						{{with .Form.FieldErrors.Value}}
							<label class='error'>{{.}}</label>
						{{end}}
						<button type="submit">OK</button>
						<button hx-get="/characters/{{.Form.CharacterId}}" hx-target="#addSkillForm" hx-swap="outerHTML" hx-select="#addSkill">Abbrechen</button>
					</form>`

		character, err := app.characters.Get(form.CharacterId)
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
		addableSkills := character.AddableSkills(allSkills)

		data := app.newTemplateData(r)
		data.Form = form
		data.AdditionalData = map[string]any{
			"AddableSkills": addableSkills}
		w.WriteHeader(http.StatusUnprocessableEntity)
		app.renderHtmx(w, r, "addSkillFailed", tmplStr, data)
		return
	}

	err = app.characters.AddSkill(form.CharacterId, form.AddableSkill, form.Value)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	half := half(form.Value)
	fifth := fifth(form.Value)
	tmplStr := fmt.Sprintf(`<template>
							<tr hx-swap-oob="beforeend:#Skills">
								<th>{{.Form.AddableSkill}}</th>
								<td>
									<div id="Values{{.Form.AddableSkill}}">{{.Form.Value}} | %d | %d</div>
									<form id="edit{{.Form.AddableSkill}}" hx-get="/characters/{{.Form.CharacterId}}/editSkill" hx-target="this" hx-swap="outerHTML">
										<input type="hidden" name="skill" value="{{.Form.AddableSkill}}">
										<input type="hidden" name="value" value="{{.Form.Value}}">
										<button type="submit">Bearbeiten</button>
									</form>
								</td>
							</tr>
							</template>
							<div id="addSkill" hx-target="this" hx-swap="outerHTML">
								<button hx-get="/characters/{{.Form.CharacterId}}/addSkill">Fertigkeit hinzufügen</button>
							</div>`, half, fifth)

	data := app.newTemplateData(r)
	data.Form = form
	w.WriteHeader(http.StatusOK)
	app.renderHtmx(w, r, "addSkillSuccess", tmplStr, data)
}

func (app *application) editSkill(w http.ResponseWriter, r *http.Request) {
	characterId := app.sessionManager.GetInt(r.Context(), characterIdKey)
	params := r.URL.Query()
	skill := params.Get("skill")
	value, err := strconv.Atoi(params.Get("value"))
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	trimmed := trim(skill)
	tmplStr := fmt.Sprintf(`<form id="editForm" hx-post="/characters/{{.Form.CharacterId}}/editSkill" hx-target="this" hx-swap="outerHTML">
                <input type="hidden" name="csrf_token" value="{{.CSRFToken}}">
				<input type="hidden" name="CharacterId" value="{{.Form.CharacterId}}">
				<input type="hidden" name="Skill" value="{{.Form.Skill}}">
                <input type="number" name="NewValue" value="{{.Form.Value}}">
				<button type="submit">OK</button>
				<button hx-get="/characters/{{.Form.CharacterId}}" hx-target="#editForm" hx-swap="outerHTML" hx-select="#edit%s">Abbrechen</button>
            </form>`, trimmed)

	data := app.newTemplateData(r)
	data.Form = map[string]any{
		"CharacterId": characterId,
		"Skill":       skill,
		"Value":       value,
	}
	w.WriteHeader(http.StatusOK)
	app.renderHtmx(w, r, "editSkillForm", tmplStr, data)
}

func (app *application) editSkillPost(w http.ResponseWriter, r *http.Request) {
	var form skillEditForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	err = app.characters.EditSkill(form.CharacterId, form.Skill, form.NewValue)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	half := half(form.NewValue)
	fifth := fifth(form.NewValue)
	trimmed := trim(form.Skill)
	tmplStr := fmt.Sprintf(`<div id="Values%s" hx-swap-oob="outerHTML:#Values%s">{{.Form.NewValue}} | %d | %d</div>
							<form hx-get="/characters/{{.Form.CharacterId}}/editSkill" hx-target="this" hx-swap="outerHTML">	
                            	<input type="hidden" name="skill" value="{{.Form.Skill}}">
                            	<input type="hidden" name="value" value="{{.Form.NewValue}}">
                            	<button type="submit">Bearbeiten</button>
                        	</form>`, trimmed, trimmed, half, fifth)

	data := app.newTemplateData(r)
	data.Form = form
	w.WriteHeader(http.StatusOK)
	app.renderHtmx(w, r, "editSkillSuccess", tmplStr, data)
}

func (app *application) addCustomSkill(w http.ResponseWriter, r *http.Request) {
	characterId := app.sessionManager.GetInt(r.Context(), characterIdKey)

	tmplStr := `<form id="addCustomSkillForm" hx-post="/characters/{{.Form.CharacterId}}/addCustomSkill" hx-target="this" hx-swap="outerHTML">
					<input type="hidden" name="csrf_token" value="{{.CSRFToken}}">
					<input type="hidden" name="CharacterId" value="{{.Form.CharacterId}}">
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
					<input type="number" name="Value"><br>
					<button type="submit">OK</button>
					<button hx-get="/characters/{{.Form.CharacterId}}" hx-target="#addCustomSkillForm" hx-swap="outerHTML" hx-select="#addCustomSkill">Abbrechen</button>
				</form>`

	data := app.newTemplateData(r)
	data.Form = map[string]any{
		"CharacterId": characterId,
	}

	w.WriteHeader(http.StatusOK)
	app.renderHtmx(w, r, "addCustomSkillForm", tmplStr, data)
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
		tmplStr := `<form id="addCustomSkillForm" hx-post="/characters/{{.Form.CharacterId}}/addCustomSkill" hx-target="this" hx-swap="outerHTML">
						<input type="hidden" name="csrf_token" value="{{.CSRFToken}}">
						<input type="hidden" name="CharacterId" value="{{.Form.CharacterId}}">
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
						<input type="text" name="CustomSkill" value="{{.Form.CustomSkill}}">
						{{with .Form.FieldErrors.Name}}<label class='error'>{{.}}</label>{{end}}
						<input type="number" name="Value" value="{{.Form.Value}}">
						{{with .Form.FieldErrors.Value}}<label class='error'>{{.}}</label>{{end}}
						<button type="submit">OK</button>
						<button hx-get="/characters/{{.Form.CharacterId}}" hx-target="#addCustomSkillForm" hx-swap="outerHTML" hx-select="#addCustomSkill">Abbrechen</button>
					</form>`

		data := app.newTemplateData(r)
		data.Form = form
		w.WriteHeader(http.StatusUnprocessableEntity)
		app.renderHtmx(w, r, "addCustomSkillInvalid", tmplStr, data)
		return
	}

	err = app.characters.AddCustomSkill(form.CharacterId, form.CustomSkill, form.Category, form.Value)
	if err != nil {
		if errors.Is(err, models.ErrAlreadyHasSkill) {
			tmplStr := `<div id="addCustomSkill" hx-target="this" hx-swap="outerHTML">
                			<button hx-get="/characters/{{.Form.CharacterId}}/addCustomSkill">Fertigkeit hinzufügen</button>
							<label class="error">Der Charaktere verfügt bereits über eine gleichnamige Fertigkeit.</label>
            			</div>`

			data := app.newTemplateData(r)
			data.Form = form
			w.WriteHeader(http.StatusOK)
			app.renderHtmx(w, r, "addCustomSkillDuplicate", tmplStr, data)
			return
		}
		app.serverError(w, r, err)
		return
	}

	half := half(form.Value)
	fifth := fifth(form.Value)
	tmplStr := fmt.Sprintf(`<template>
							<tr hx-swap-oob="beforeend:#CustomSkills">
								<th>{{.Form.CustomSkill}}</th>
								<td>
									<div id="Values{{.Form.CustomSkill}}" value="{{.Form.Value}}">{{.Form.Value}} | %d | %d</div>
									<form id="edit{{.Form.CustomSkill}}" hx-get="/characters/{{.Form.characterId}}/editCustomSkill" hx-target="this" hx-swap="outerHTML">
										<input type="hidden" name="skill" value="{{.Form.CustomSkill}}">
										<input type="hidden" name="value" value="{{.Form.Value}}">
										<button type="submit">Bearbeiten</button>
									</form>
								</td>
							</tr>
							</template>
							<div id="addCustomSkill" hx-target="this" hx-swap="outerHTML">
								<button hx-get="/characters/{{.Form.CharacterId}}/addCustomSkill">Fertigkeit hinzufügen</button>
							</div>`, half, fifth)

	data := app.newTemplateData(r)
	data.Form = form
	w.WriteHeader(http.StatusOK)
	app.renderHtmx(w, r, "addCustomSkillSuccess", tmplStr, data)
}

func (app *application) editCustomSkill(w http.ResponseWriter, r *http.Request) {
	characterId := app.sessionManager.GetInt(r.Context(), characterIdKey)
	params := r.URL.Query()
	skill := params.Get("skill")
	value, err := strconv.Atoi(params.Get("value"))
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	tmplStr := fmt.Sprintf(`<form id="editForm" hx-post="/characters/%d/editCustomSkill" hx-target="this" hx-swap="outerHTML">
                <input type="hidden" name="csrf_token" value="{{.CSRFToken}}">
				<input type="hidden" name="CharacterId" value="%d">
				<input type="hidden" name="Skill" value="%s">
                <input type="number" name="NewValue" value="%d">
				<button type="submit">OK</button>
				<button hx-get="/characters/%d" hx-target="#editForm" hx-swap="outerHTML" hx-select="#edit%s">Abbrechen</button>
            </form>`, characterId, characterId, skill, value, characterId, skill)

	data := app.newTemplateData(r)
	w.WriteHeader(http.StatusOK)
	app.renderHtmx(w, r, "editCustomSkillForm", tmplStr, data)
}

func (app *application) editCustomSkillPost(w http.ResponseWriter, r *http.Request) {
	var form skillEditForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	half := half(form.NewValue)
	fifth := fifth(form.NewValue)
	tmplStr := fmt.Sprintf(`<div value="{{.Form.NewValue}}" hx-swap-oob="outerHTML:#Values{{.Form.Skill}}">{{.Form.NewValue}} | %d | %d</div>
							<form hx-get="/characters/{{.Form.CharacterId}}/editCustomSkill" hx-target="this" hx-swap="outerHTML">	
                            	<input type="hidden" name="skill" value="{{.Form.Skill}}">
                            	<input type="hidden" name="value" value="{{.Form.NewValue}}">
                            	<button type="submit">Bearbeiten</button>
                        	</form>`, half, fifth)

	err = app.characters.EditCustomSkill(form.CharacterId, form.Skill, form.NewValue)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Form = form
	w.WriteHeader(http.StatusOK)
	app.renderHtmx(w, r, "editCustomSkillSuccess", tmplStr, data)
}

func (app *application) customSkillInput(w http.ResponseWriter, r *http.Request) {
	tmplStr := `<tr id="{{.Form.Category}}">
					<td>
						<input type='hidden' name='CustomSkills.Category' value='{{.Form.Category}}'>
						<label>{{.Form.Category}}</label>
						<input type="text" name="CustomSkills.Name">
						<select name="CustomSkills.Value">
							<option value="{{.Form.Default}}" selected>{{.Form.Default}}</option>
							<option value="70">70</option>
							<option value="60">60</option>
							<option value="50">50</option>
							<option value="40">40</option>
						</select>
						<button hx-get="/create" hx-target="#{{.Form.Category}}" hx-swap="delete">Abbrechen</button>
					</td>
				</tr>`

	category := r.URL.Query().Get("category")
	defaultValue := models.DefaultForCategory(category)

	data := app.newTemplateData(r)
	data.Form = map[string]any{
		"Category": category,
		"Default":  defaultValue,
	}
	w.WriteHeader(http.StatusOK)
	app.renderHtmx(w, r, "customSkillInput", tmplStr, data)
}

func (app *application) editStat(w http.ResponseWriter, r *http.Request) {
	characterId := app.sessionManager.GetInt(r.Context(), characterIdKey)

	var form statEditForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusUnprocessableEntity)
		return
	}

	tmplStr := `<div id="{{.Form.Stat}}">
					{{if gt .Form.NewValue 1}}
					<button type="submit" name="Direction" value="dec">-</button>
					{{end}}
					<input type="hidden" name="Name" value="{{.Form.Stat}}">
					<input type="hidden" name="Value" value="{{.Form.NewValue}}">
					{{.Form.NewValue}}
					{{if lt .Form.NewValue .Form.Max}}
					<button type="submit" name="Direction" value="inc">+</button>
					{{end}}
				</div>`

	character, err := app.characters.Get(characterId)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	max := character.Stats.GetStatMax(form.Name)

	var updated int
	data := app.newTemplateData(r)
	switch form.Direction {
	case "inc":
		updated, err = app.characters.IncrementStat(characterId, form.Name)
		data.Form = map[string]any{
			"Stat":     form.Name,
			"NewValue": updated,
			"Max":      max,
		}
	case "dec":
		updated, err = app.characters.DecrementStat(characterId, form.Name)
		data.Form = map[string]any{
			"Stat":     form.Name,
			"NewValue": updated,
			"Max":      max,
		}
	}
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	app.renderHtmx(w, r, "editStatSuccess", tmplStr, data)
}

func (app *application) addItem(w http.ResponseWriter, r *http.Request) {
	characterId := app.sessionManager.GetInt(r.Context(), characterIdKey)
	data := app.newTemplateData(r)
	data.Form = itemForm{
		CharacterId: characterId,
	}
	w.WriteHeader(http.StatusOK)
	app.render(w, r, "item.tmpl.html", data)
}

func (app *application) addItemPost(w http.ResponseWriter, r *http.Request) {
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

	fmt.Println(form.CharacterId, form.Name, form.Description, form.Count)

	err = app.characters.AddItem(form.CharacterId, form.Name, form.Description, form.Count)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	redirect := fmt.Sprintf("/characters/%d", form.CharacterId)
	http.Redirect(w, r, redirect, http.StatusSeeOther)
}

func (app *application) editItemCount(w http.ResponseWriter, r *http.Request) {
	characterId := app.sessionManager.GetInt(r.Context(), characterIdKey)

	var form itemEditForm
	err := app.decodePostForm(r, &form)
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

	data := app.newTemplateData(r)
	switch form.Direction {
	case "inc":
		err = app.characters.EditItemCount(characterId, form.ItemId, form.Count+1)
		data.Form = map[string]int{
			"ItemId":   form.ItemId,
			"NewCount": form.Count + 1,
		}
	case "dec":
		err = app.characters.EditItemCount(characterId, form.ItemId, form.Count-1)
		data.Form = map[string]int{
			"ItemId":   form.ItemId,
			"NewCount": form.Count - 1,
		}
	}
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	app.renderHtmx(w, r, "editItemCountSuccess", tmplStr, data)
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

	w.WriteHeader(http.StatusOK)
	app.renderHtmx(w, r, "empty", "", templateData{})
}

func (app *application) addNote(w http.ResponseWriter, r *http.Request) {
	characterId := app.sessionManager.GetInt(r.Context(), characterIdKey)

	tmplStr := `<form id="addNoteForm" hx-post="/characters/{{.Form.CharacterId}}/addNote" hx-target="this" hx-swap="outerHTML">
					<input type="hidden" name="csrf_token" value="{{.CSRFToken}}">
					<input type="hidden" name="CharacterId" value="{{.Form.CharacterId}}">
					<div>
						<label>Notiz:</label>
						<input type="text" name="Text" textarea>
					</div>
					<button type="submit">Hinzufügen</button>
					<button hx-get="/characters/{{.Form.CharacterId}}" hx-target="#addNoteForm" hx-swap="delete">Abbrechen</button>
				</form>`

	data := app.newTemplateData(r)
	data.Form = map[string]any{
		"CharacterId": characterId,
	}
	w.WriteHeader(http.StatusOK)
	app.renderHtmx(w, r, "addNoteForm", tmplStr, data)
}

func (app *application) addNotePost(w http.ResponseWriter, r *http.Request) {
	var form noteForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusUnprocessableEntity)
		return
	}
	noteId, err := app.characters.AddNote(form.CharacterId, form.Text)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	tmplStr := `<form id="deleteNote" hx-post="/characters/{{.Form.CharacterId}}/deleteNote" hx-target="this" hx-swap="outerHTML">
								<input type="hidden" name="csrf_token" value="{{.CSRFToken}}">
								<input type="hidden" name="NoteId" value="{{.Form.NoteId}}">
								<li>{{.Form.Text}}    <button type="submit">löschen</button></li>
							</form>`

	data := app.newTemplateData(r)
	data.Form = map[string]any{
		"CharacterId": form.CharacterId,
		"NoteId":      noteId,
		"Text":        form.Text,
	}
	w.WriteHeader(http.StatusOK)
	app.renderHtmx(w, r, "addNoteSuccess", tmplStr, data)
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

	w.WriteHeader(http.StatusOK)
	app.renderHtmx(w, r, "empty", "", templateData{})
}

type uploadForm struct {
	Title                    string
	UploadedBy               int
	validators.FormValidator `schema:"-"`
}

func (app *application) uploadMaterial(w http.ResponseWriter, r *http.Request) {
	userId := app.sessionManager.GetInt(r.Context(), authenticatedUserIdKey)

	data := app.newTemplateData(r)
	data.Form = uploadForm{
		UploadedBy: userId,
	}
	w.WriteHeader(http.StatusOK)
	app.render(w, r, "upload.tmpl.html", data)
}

func (app *application) uploadMaterialPost(w http.ResponseWriter, r *http.Request) {
	var form uploadForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusUnprocessableEntity)
		return
	}

	file, header, err := r.FormFile("FilePath")
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	defer file.Close()

	err = app.users.AddMaterial(form.Title, header.Filename, form.UploadedBy)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateFileName) {
			form.AddGenericError("Eine Datei mit diesem Namen existiert bereits.")

			data := app.newTemplateData(r)
			data.Form = form
			w.WriteHeader(http.StatusUnprocessableEntity)
			app.render(w, r, "upload.tmpl.html", data)
			return
		}
		app.serverError(w, r, err)
		return
	}

	path := filepath.Join("ui/static/img/uploads/", strconv.Itoa(form.UploadedBy))
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	f, err := os.Create(filepath.Join(path, header.Filename))
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	defer f.Close()

	_, err = io.Copy(f, file)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) materials(w http.ResponseWriter, r *http.Request) {
	userName := app.sessionManager.GetString(r.Context(), authenticatedUserNameKey)
	user, err := app.users.Get(userName)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.User = user
	w.WriteHeader(http.StatusOK)
	app.render(w, r, "materials.tmpl.html", data)
}
