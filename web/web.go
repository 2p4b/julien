package web

import (
	"encoding/json"
	"fmt"
	"julien/driver"
	"julien/form"
	"julien/fs"
	"julien/julien"
	"julien/pager"
	"julien/template"
	jutils "julien/utils"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/session"

	"github.com/gofiber/utils/v2"
)

const POST_KEY string = "post"

const FORM_KEY string = "session"

const VIEWS_KEY string = "views"

const LAYOUTS_KEY string = "layouts"

const FORM_TIMEOUT int64 = 30

type Posted struct {
	Form      string `json:"form"`
	Doc       string `json:"doc"`
	Timestamp int64  `json:"timestamp"`
}

type Post struct {
	Form      *form.Form
	Data      *form.Doc
	Timestamp int64
}

type Web struct {
	config   *julien.Julien
	store    *session.Store
	site     *julien.Site
	forms    *form.Root
	content  *pager.Root
	template *template.Template
}

func SaveSession(sess *session.Session) {
	if err := sess.Save(); err != nil {
		log.Error(err)
	}
}

func IncludeData(ctx *fiber.Ctx, form form.Form, data map[string]interface{}) map[string]interface{} {
	includes := form.Get("includes")
	if includes == nil {
		return data
	}
	includesmap, ok := includes.(map[interface{}]interface{})
	if !ok {
		return data
	}

	for key, value := range includesmap {
		strkey, ok := key.(string)
		if !ok {
			continue
		}
		strfield, ok := value.(string)
		if !ok {
			continue
		}

		switch strfield {
		case "ip":
			data[strkey] = ctx.IP()

		case "hostname":
			data[strkey] = ctx.Hostname()

		case "User-Agent", "Referer":
			data[strkey] = ctx.Get(strfield, "")
		}
	}

	return data
}

func MakePageInfo(fm *form.Form) map[string]string {
	now := time.Now()
	nanosecond := strconv.Itoa(now.Nanosecond())
	milisecond := strconv.Itoa(now.Nanosecond() / 1000000)
	ts := now.Format("20060102150405") + nanosecond
	return map[string]string{
		"nanosecond": nanosecond,
		"milisecond": milisecond,
		"timestamp":  ts,
		"default":    ts,
		"name":       fm.Name(),
	}
}

func GetName(format string, info map[string]string, data map[string]interface{}) string {
	if format[0] == '$' {
		value, ok := info[format[1:]]
		if ok {
			return value
		}
		return "page[" + format + "]"
	} else if format[0] == '@' {
		value, ok := data[format[1:]].(string)
		if ok {
			return value
		}
		return "data[" + format + "]"
	} else {
		return format
	}

}

func MakeName(fm *form.Form, data map[string]interface{}) string {
	format, ok := fm.Get("name").(string)
	if !ok {
		format = "$default"
	}

	seperator, ok := fm.Get("seperator").(string)
	if !ok {
		seperator = "_"
	}

	parts := strings.Split(format, seperator)

	info := MakePageInfo(fm)
	for index, part := range parts {
		parts[index] = GetName(part, info, data)
	}

	cleaned := make([]string, 0)
	for _, part := range parts {
		if part != "" {
			cleaned = append(cleaned, part)
		}
	}

	return strings.Join(cleaned, seperator)
}

func collect(data map[string]interface{}, rules map[string]interface{}) map[string]interface{} {
	values := make(map[string]interface{}, 0)
	for key := range rules {
		value, ok := data[key]
		if ok {
			values[key] = interface{}(value)
		}
	}

	return values
}

func render(web *Web, ctx *fiber.Ctx, name string) error {
	forms := web.Forms()
	content := web.Content()
	code, cerr := strconv.Atoi(name)
	page, err := content.Find(name)
	if err != nil {
		if cerr == nil && code == 400 {
			return ctx.SendStatus(code)
		} else {
			return render(web, ctx, "404")
		}
	} else {
		if cerr == nil && (code < 400 || code > 451) && (code < 500 || code > 511) {
			code = fiber.StatusOK
		}
	}

	// Reroute index to parent dir
	if page.IsIndex() {
		return ctx.Redirect("/"+page.Path(), 302)
	}

	view := page.View()
	layout := page.Layout()

	// Get session from storage
	sess, err := web.Session(ctx)
	if err != nil {
		log.Error(err)
		return render(web, ctx, "500")
	}

	webformstr, ok := sess.Get(FORM_KEY).(string)
	if !ok {
		webformstr = "{}"
	} else {
	}

	formdata := &FormData{}
	err = json.Unmarshal([]byte(webformstr), &formdata)
	if err != nil {
		log.Error(err)
	}

	vparams := fiber.Map{
		"Ctx":      ctx,
		"Page":     page,
		"Site":     web.Site(),
		"Forms":    forms,
		"Pager":    content,
		"FormData": formdata,
		"Template": web.Template(),
	}

	postedstr, ok := sess.Get(POST_KEY).(string)
	if !ok {
		postedstr = "{}"
	}
	posted := Posted{Doc: "", Form: "", Timestamp: 0}
	err = json.Unmarshal([]byte(postedstr), &posted)
	if err != nil {
		log.Error(err)
		return render(web, ctx, "500")
	}

	if posted.Timestamp > 0 {
		fm, err := forms.Find(posted.Form)
		if err == nil {
			doc, err := fm.Find(posted.Doc)
			if err == nil && posted.Timestamp >= (time.Now().Unix()-(FORM_TIMEOUT)) {
				vparams["Post"] = Post{
					Form:      fm,
					Data:      doc,
					Timestamp: posted.Timestamp,
				}
			}
		}
	}

	// Done extracting valid form data from session
	// So we don't need this in the session
	// anymore so trash it
	sess.Delete(FORM_KEY)
	sess.Delete(POST_KEY)
	SaveSession(sess)

	view = path.Clean(path.Join(web.template.GetString(VIEWS_KEY, VIEWS_KEY), view))
	layout = path.Clean(path.Join(web.template.GetString(LAYOUTS_KEY, LAYOUTS_KEY), layout))

	if (code < 400 || code > 451) && (code < 500 || code > 511) {
		return ctx.Render(view, vparams, layout)
	} else {
		return ctx.Status(code).Render(view, vparams, layout)
	}
}

func mapcast(schema map[interface{}]interface{}) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	for key, value := range schema {
		strkey, ok := key.(string)
		if !ok {
			return nil, fmt.Errorf("cast: invalid form schema key: %s", key)
		}
		data[strkey] = value
	}
	return data, nil
}

func BaseUrl(link *url.URL) string {
	return link.Scheme + "://" + link.Host + link.Path
}

func New(config *julien.Julien, site *julien.Site) Web {
	store := session.New()
	yamler := &driver.Yaml{}

	Data := config.Data
	Forms := config.Forms
	Content := config.Content

	tmpl, err := template.Find(config.TemplatePath())
	if err != nil {
		log.Error(err)
		panic(err)
	}

	cdisk := fs.Mount(Content.Path, Content.Index, Content.Ext)
	ddisk := fs.Mount(Data.Path, Data.Index, Data.Ext)
	fdisk := fs.Mount(Forms.Path, Forms.Index, Forms.Ext)
	forms := form.Init(fdisk, ddisk, yamler)
	content := pager.Init(cdisk, yamler)
	return Web{
		config:   config,
		store:    store,
		forms:    &forms,
		content:  &content,
		site:     site,
		template: tmpl,
	}
}

func (web *Web) Site() *julien.Site {
	return web.site
}

func (web *Web) Logger() func(*fiber.Ctx) error {
	fmtstr := strings.Trim(web.config.Logger.Format, " ")
	logfile := strings.Trim(web.config.Logger.File, " ")

	fmtstr = strings.Trim(fmtstr, "\n")
	config := logger.Config{Format: fmtstr + "\n"}

	if logfile != "" {
		file, err := os.OpenFile(logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Error("error opening logfile: %v", err)
		} else {
			config.Output = file
		}
	}

	return logger.New(config)
}

func (web *Web) Template() *template.Template {
	return web.template
}

func (web *Web) Session(ctx *fiber.Ctx) (*session.Session, error) {
	return web.store.Get(ctx)
}

func (web *Web) TemplateName() string {
	return web.config.TemplateName()
}

func (web *Web) TemplatePath() string {
	return web.template.Path
}

func (web *Web) TemplateExt() string {
	return web.template.GetString("ext", "html")
}

func (web *Web) FormsPath() string {
	return web.config.FormsPath()
}

func (web *Web) ContentPath() string {
	return web.config.ContentPath()
}

func (web *Web) AllowedFiles() []string {
	return web.config.ContentAssets()
}

func (web *Web) StaticPath() string {
	return web.config.StaticPath()
}

func (web *Web) Start(addr string) {

	var app = fiber.New(fiber.Config{
		AppName: "Julien",
		Views:   web.template.Engine(true),
	})

	app.Use(idempotency.New())
	app.Use(cors.New())
	app.Use(limiter.New(limiter.Config{
		Max:        60,
		Expiration: 5,
	}))

	app.Use(csrf.New(csrf.Config{
		KeyLookup:         "cookie:csrf",
		CookieName:        "csrf",
		CookieSameSite:    "Lax",
		Expiration:        1 * time.Hour,
		KeyGenerator:      utils.UUIDv4,
		Session:           web.store,
		SessionKey:        "fiber.csrf.token",
		HandlerContextKey: "fiber.csrf.handler",
	}))

	app.Use(web.Logger())
	app.Use(compress.New())

	app.Static("/static", web.config.StaticPath(), fiber.Static{
		Compress:      true,
		ByteRange:     true,
		CacheDuration: 24 * 60 * 60 * time.Second,
	})

	app.Static("/public", web.TemplatePath()+"/public")

	app.Get("/metrics", monitor.New())

	app.Get("/*", func(c *fiber.Ctx) error {
		return web.RenderPage(c)
	})

	app.Post("/:form?", func(c *fiber.Ctx) error {
		return web.RenderForm(c)
	})

	app.Listen(addr)
}

func (web *Web) Forms() *form.Root {
	return web.forms
}

func (web *Web) Content() *pager.Root {
	return web.content
}

func (web *Web) RenderPage(ctx *fiber.Ctx) error {
	name := ctx.Params("*")
	ext := path.Ext(name)

	if ext != "" && jutils.ArrayIncludes(web.AllowedFiles(), ext[1:]) {
		filepath := path.Join(web.ContentPath(), ctx.Path())
		return ctx.SendFile(filepath, true)
	}
	return render(web, ctx, name)
}

func (web *Web) RenderForm(ctx *fiber.Ctx) error {
	var is_formdata = false
	var errormap = make(map[string][]string, 0)
	var validate = validator.New()
	var name = ctx.Params("form")
	forms := web.Forms()

	source := ctx.Get("Referer", "/")

	fm, err := forms.Find(name)
	if err != nil {
		// Form not found
		return render(web, ctx, "404")
	}

	// Get session from storage
	sess, err := web.Session(ctx)
	if err != nil {
		return render(web, ctx, "500")
	}

	data := make(map[string]interface{}, 0)

	rules, ok := fm.Get("schema").(map[interface{}]interface{})

	if !ok {
		return render(web, ctx, "500")
	}

	skrules, err := mapcast(rules)
	if err != nil {
		log.Error(err)
		return ctx.Redirect(source, 302)
	}

	if err := ctx.BodyParser(&data); err != nil {
		nullrep := utils.UUID()
		// Attempt to copy formdata fields in the schema
		for rkey := range skrules {
			val := ctx.FormValue(rkey, nullrep)
			if val != nullrep {
				is_formdata = true
				data[rkey] = interface{}(val)
			}
		}
	}

	values := collect(data, skrules)

	// Redirect if no value was submitted
	if len(values) == 0 {
		return ctx.Redirect(source, 302)
	}

	errors := validate.ValidateMap(values, skrules)
	for key, ferror := range errors {
		ferrmap := make([]string, 0)
		ferrors := ferror.(validator.ValidationErrors)
		for _, ferr := range ferrors {
			ferrmap = append(ferrmap, ferr.ActualTag())
		}
		errormap[key] = ferrmap
	}

	if len(errors) > 0 {
		if is_formdata {
			formdata := FormData{
				Name:      name,
				Data:      values,
				Errors:    errormap,
				Timestamp: time.Now().Unix(),
			}
			serialdata, _ := json.Marshal(formdata)
			sess.Set(FORM_KEY, string(serialdata))
			SaveSession(sess)
			return ctx.Redirect(source, 302)
		}
		return ctx.JSON(errormap)
	}

	values = IncludeData(ctx, *fm, values)
	filename := MakeName(fm, values)

	// Create new doc for form data
	// return 500 if this fails
	doc, err := fm.Compose(filename, make([]byte, 0))
	if err != nil {
		log.Error(err)
		return ctx.Redirect(source, 500)
	}

	content := ""
	content_key, ok := fm.Get("content").(string)
	if ok && content_key != "" {
		content, ok = values[content_key].(string)
		if !ok {
			content = ""
		} else {
			// Found content so, remove it
			// from the values key value
			// store to data is not duplicated
			delete(values, content_key)
		}
	}

	// Fill created doc with form data and save
	// return 500 if this fails
	doc.Fill(values)
	doc.Body(content)
	err = doc.Save()
	if err != nil {
		log.Error(err)
		return ctx.Redirect(source, 500)
	}

	// Record form subimission in session
	posted := Posted{
		Form:      name,
		Doc:       filename,
		Timestamp: time.Now().Unix(),
	}
	serialdata, _ := json.Marshal(posted)
	sess.Set(POST_KEY, string(serialdata))
	SaveSession(sess)

	// Get redirect path if everything went well
	// or post path for request
	repath, ok := fm.Get("redirect").(string)
	if !ok {
		repath = name
	}

	if is_formdata {
		return ctx.Redirect(repath, 302)
	} else {
		// return Empty json if the post request
		// was a json
		return ctx.JSON(make(map[string]string, 0))
	}

}
