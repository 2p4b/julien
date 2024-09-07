Julien, "NoCMS" built for folks like you and me. It gets the job done without any headache. 

#### Why create yet another cms?
I was trying to get a web page up over the weekend, with a simple contact us form, but it required using a php wordpress plugin and start maintaining a yet another php stack, or using hugo and signing up for a subscription service to collect the form data

Most of my clients want a simple webpage and maybe some react with a simple form 
without needing to signup for yet another subscription service or data collection site
just for collecting contact us form data or other simple forms on their respective sites

Though static site generators like [Hugo](https://gohugo.io/) with a mailto, will meet needs of most users, I find myself reaching out to other third party services when I needed a simple form, most of which are moving towards a subscription model or clients/users have to pay with their data(Google forms) and I can only subscribe to so many services before i go bankrupt.

#### What is NoCMS
__NoCMS__, just like other BS words like serverless, nosql, etc...

#### Why Julien?
Julien is the Lemur King in the Animation series [All Hail King Julien](https://www.dreamworks.com/shows/all-hail-king-julien) who sought freedom for himself and other lemurs in his kingdom though he wasn't very bright
on how he went about doing this, his intentions were noble.

#### Start serving requests

Run dev
```sh
go run . --config=example-julien.yaml --site=example/index.md
```

Run julien
```sh
julien --config=julien.yaml --site=index.md
```

- `--host=127.0.0.1` webserver host default `localhost`
- `--port=8080` webserver port default `1234`
- `--site=index.md` Site file  default 'index.md'
- `--config=julien.yaml` julien config file  default 'julien.yaml'


##### How to configure Julien
Julien keeps things clean and simple with YAML for configuration. There are five mount location. Here's the lowdown on the essential settings and how to bend them to your will:

```yaml
# julien.yaml
static: 
    path: static # Where your static files like images, CSS, and JS hang out

content:
    ext: "md" # File extension for your content files (Markdown is the way to go!)
    path: "content" # The lair of your website's content
    index: "index" # The default filename for index pages (e.g., index.md)
    driver: "yaml-md" # How Julien reads and writes your content (stick with the default for now)
    assets: [png, jpg, jpeg, gif] # Allowed image file types within your content directories

data: 
    ext: "md" # Same as content, but for storing submitted form data
    path: "data" # Where Julien stashes the goods from your forms
    index: "index" 
    driver: "yaml-md" 

forms: 
    ext: "md" # You guessed it, Markdown for forms too!
    path: "forms" # The command center for your website's forms
    index: "index" 
    driver: "yaml-md" 

template: 
    path: templates # The directory where your website's templates reside
    name: julien # The chosen one - the template that will bring your website to life
```

#### File Structure
Think of file structure as the blueprint for your website. Here's a breakdown of the key locations and what they do:

1) __content__: This is where you keep your site page contents usually a markdown document with frontmatter for per page configuration, assets are the static files allowed to be served from the page location
- 
    ```text
    content/
    ├── index.md 
    ├── articles/
    │   ├── index.md
    │   ├── intro-hail-julien.md
    │   └── all-hail-king-julien/
    │       ├── index.md
    │       └── kings_porait.jpg
    │
    └── contact-us.md
    ```

2) __forms__: This is where you define your website's forms. Each Markdown file in this directory represents a form.

3) __data__: This is where Julien stores the data submitted through your forms. Each form submission gets its own file in this directory 

4) __static__: This is the home for all your site static assets, like images, CSS files, and JavaScript files 

5) __template__: This is where you store your website's templates. Templates define the look and feel of your website and how your content is displayed 

That's it for locations



#### content files
the contents directory should have and index.md file which will serve as the contents for the home or root page at / 
get request path are direct mappings to the filenames and other directries within the content directory

##### content view 
views for specific pages can be defined in the frontmatter by the view key
if you need to handle the views for a list of files in a directory that can be defined in the directory index.md frontmatter as

```yaml
# contents/articles/index.md
title: Articles
view: articles
layout: main
page:
    view: article
    layout: main
```

The `page:` directive in a directory index.md file defines default view and layout for all other markdown documents and sub directories within the current directory. This can be overwritten by the frontmatter of individual page be defining them in the frontmatter of the page like

```yaml
# contents/about-us.md
view: about-us
layout: main
```
Each sub directory inherits the view and layout of its parent directory by default if the `page:` directive is not found in the parent index file 

#### forms
Forms in Julien are defined in Markdown files within the /forms directory. Here's an example of a simple contact form (contact.md)

```yaml
# forms/contact-us.md
---
title: Contact Us
driver: yaml-md # Keep it simple, stick with YAML
name: $timestamp # Each submission gets a unique timestamped filename
redirect: /thank-you # Redirect to a thank-you page after successful submission

includes: 
    ip: ip # Include the user's IP address in the submission data
    hostname: hostname # Include the user's hostname
    referer: Referer # Include the referring URL
    useragent: User-Agent # Include the user's browser information

schema:
    name: required,min=1,max=32 # Validation rules for the 'name' field
    email: required,email # Validation rules for the 'email' field
    message: required,min=10 # Validation rules for the 'message' field

content: message # Optional content field to be extracted from the form and used as the markdown document content
---

Contact me now!
Why wait when we could be building all sort of ideas
Don't wait!
Write your articles in markdown with an editor of your choice. place them in a folder and
thats all to it

```

- __title__: defines the title of the form
- __driver__: defines the driver used to dump the contents to the filesystem (stil WIP)
- __name__: the file name for each document `$timestamp` uses the request constant timestamp used as the filename
- __redirect__: The URL to redirect to after a successful form submission.
- __includes__: Request runtime values to include in the form data with the keys being the same key to be used and the value is a request value to be extracted and added to the form content before it is written to disk

- __schema__: This is the form schema used to validate the form request see [validator](https://github.com/go-playground/validator) for other rules just keep it simple and restricted to map validation rules ONLY and you will be fine

- __content__: Content key if defined the field will be extracted from the frontmatter and use as the content body in the markdown document and will NOT be in the frontmatter when dumped to disk


### Template
Each template directory must include the index.md file at its root with information about the 
template and configuration 
Here is sample directory structure of a template named julien:

```text
    templates/
    │
    ├── julien/
    │   ├── index.md
    │   ├── public/
    │   │   ├── styles.css
    │   │   └── react.js
    │   ├── partials/
    │   │   ├── header.html
    │   │   └── footer.html
    │   ├── views/
    │   │   ├── 404.html
    │   │   ├── 500.html
    │   │   ├── index.html
    │   │   ├── contact-us.html
    │   │   ├── article.html
    │   │   └── articles.html
    │   └── layouts/
    │       ├── main.html
    │       └── mobile.html
    │
    └── mytemplate/
```
Julien uses the four types of templating systems from which you can choose from
or migrate and existing project

1) __html__: [golang official templating system](https://pkg.go.dev/html/template) template library  and the default templating system 

2) __pongo__: [Pongo 2](https://www.schlachter.tech/solutions/pongo2-template-engine/) template engine in the [fiber plugin](https://docs.gofiber.io/template/django/) library. I use this engine myself because i already have lost of experience with django templating systems 


3) __mustache__: [Mustache templating system](https://mustache.github.io/mustache.5.html) template engine

4) __jet__: [ Jet templating system](https://github.com/CloudyKit/jet/wiki/3.-Jet-template-syntax) template engine

here is a sample template index.md file

```yaml
# index.md
---
name: Julien
author: 2p4b 
layouts:  layouts #Path to templates views          default to layouts
public: public #Path to templates public assets     default to public
views: views #Path to templates views               default to views
type: pogo #Template engine to use                  default to html
ext: html #Templates file extensions                default to html
---
Julien Example Template

```

#### Data variables
These variables share a common interface e.g `Page.Get("title")` or `Site.Get("title", "default")` to get 
frontmatter data and `Page.Content` for the document body same is true for all data variables

- `Site` Site document
- `Page` Current page document
- `Post` is the data submitted successfully by a form 
    - `Post.Form` Post Form document
    - `Post.Data` Post Data document
    - `Post.Timestamp` is the unix timestamp of the submission
- `Template` Current template index.md document

#### Dealing with form submission errors
The `FormData` variable is populated after a post request with validation errors based on 
the form document `schema` field

```yaml
# Example form schema /forms/contact-us.md
schema:
    name: required,min=1,max=32 # Validation rules for the 'name' field
    email: required,email # Validation rules for the 'email' field
    company: required,min=10 # Validation rules for the 'message' field
```
- `FormData` Submitted form values with errors on `POST` requests
    - `FormData.Name` Form document name 
    - `FormData.Get("company")` Get submitted form `company` value
    - `FormData.Get("company", "default")` Get submitted form `company` value or default value
    - `FormData.HasErrors("company")` returns boolean if submitted form `company` has errors
    - `FormData.HasErrors("company", "min")` returns boolean if submitted form `company` has error with tag min e.g For form schema `name: required,min=1,max=32`
    - `FormData.HasErrors("company", "min", "max")` returns boolean if submitted form `company` has error with tag min or max e.g For form schema `name: required,min=1,max=32`. can check multiple tags 
    - `FormData.GetErrors("company")` Get form `company` error tags
    - `Post.Timestamp` is the unix timestamp of the submission

#### Mounts variables
- `Forms` is the mount point for the forms
- `Pager` is the mount point for the content

Mounts variables share the common interface of `Forms.Find("contact-us")` or `Pager.Find("articles/my-article")` and returns  a pointer and error. I the document is not found nil/null is return with an error why. If the document is found a pointer to the document data is returned with error set to nil/null.

Another quick way to get a document if you don't need to know the error is to use the `Open` method like
`Forms.Open("contact-us")` or `Pager.Open("articles/my-article")`. The `Open` method returns a pointer to the document if
found or nil/null.

#### Todo
- [ ] implement sqlite driver
- [ ] Optional Frontend ui for the masses