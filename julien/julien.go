package julien

import (
	"julien/driver"
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

type StaticMount struct {
	Path string `yaml:"path"`
	// Send static files bytes ranges in future fiber v3
	Bytes bool `yaml:"bytes"`
	// Compress static files in future fiber v3
	Compress bool `yaml:"compress"`
}

type Logger struct {
	Format string `yaml:"format"`
	File   string `yaml:"file"`
}

type MountPoint struct {
	Path   string   `yaml:"path"`
	Ext    string   `yaml:"ext"`
	Index  string   `yaml:"index"`
	Driver string   `yaml:"driver"`
	Assets []string `yaml:"assets"`
}

type Meta struct {
	Name     string `yaml:"name"`
	Property string `yaml:"property"`
	Content  string `yaml:"content"`
}

type Template struct {
	Path string `yaml:"path"`
	Name string `yaml:"name"`
}

type Julien struct {
	Data     MountPoint  `yaml:"data"`
	Forms    MountPoint  `yaml:"forms"`
	Content  MountPoint  `yaml:"content"`
	Template Template    `yaml:"template"`
	Static   StaticMount `yaml:"static"`
	Logger   Logger      `yaml:"logger"`
}

func (j *Julien) DataPath() string {
	return j.Data.Path
}

func (j *Julien) FormsPath() string {
	return j.Forms.Path
}

func (j *Julien) ContentPath() string {
	return j.Content.Path
}

func (j *Julien) StaticPath() string {
	return j.Static.Path
}

func (j *Julien) TemplatePath() string {
	return path.Join(j.Template.Path, j.Template.Name)
}

func (j *Julien) TemplateName() string {
	return j.Template.Name
}

func (j *Julien) ContentAssets() []string {
	return j.Content.Assets
}

func DefaultSite() Site {
	return Site{
		body: "",
		path: "",
		meta: make(map[string]interface{}, 0),
	}
}

func CreateDefaultTemplate(path string) Template {
	return Template{
		Name: "julien",
		Path: path,
	}
}

func CreateDefaultMount(path string) MountPoint {
	return MountPoint{
		Ext:    "md",
		Path:   path,
		Index:  "index",
		Driver: "yaml-md",
		Assets: []string{"png", "jpg", "jpeg", "gif", "mp4", "webm"},
	}
}

func CreateStaticMount(path string) StaticMount {
	return StaticMount{
		Path: path,
	}
}

func DefaultJulien() Julien {
	return Julien{
		Data:     CreateDefaultMount("data"),
		Forms:    CreateDefaultMount("forms"),
		Content:  CreateDefaultMount("content"),
		Template: CreateDefaultTemplate("templates"),
		Static:   CreateStaticMount("static"),
		Logger:   Logger{Format: "[${ip}]:${method} ${path} - ${status}"},
	}
}

func ReadFile(path string) string {
	cbytes, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return string(cbytes)
}

func LoadConfigFromStr(config string, j *Julien) {
	if merr := yaml.Unmarshal([]byte(config), j); merr != nil {
		panic(merr)
	}
}

func LoadSiteFromStr(config string, s *Site) {
	yamler := driver.Yaml{}
	frontstr, content, err := yamler.Parts(config)
	if err != nil {
		panic(err)
	}

	if merr := yaml.Unmarshal([]byte(frontstr), s); merr != nil {
		panic(merr)
	}
	s.body = content
}

func LoadConfig(path string, j *Julien) {
	strconfig := ReadFile(path)
	LoadConfigFromStr(strconfig, j)
}
