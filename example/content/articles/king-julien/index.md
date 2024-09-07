---
title: King Julien Pongo
hero: king_julien_ai_gen.jpeg
view: article
---

### Template
Here is sample directory structure of a template named pongotemplate:

```text
templates/
└── pongotemplate/
    ├── index.md
    ├── public/
    │   ├── styles.css
    │   └── react.js
    ├── partials/
    │   ├── header.html
    │   └── footer.html
    ├── views/
    │   ├── 404.html
    │   ├── 500.html
    │   ├── index.html
    │   ├── contact-us.html
    │   ├── article.html
    │   └── articles.html
    └── layouts/
        ├── main.html
        └── mobile.html
 
```

Julien pongo [Pongo 2](https://www.schlachter.tech/solutions/pongo2-template-engine/) uses the template engine in the [fiber plugin](https://docs.gofiber.io/template/django/) library. I use this engine myself because i already have lost of experience with django templating systems 

here is a sample template index.md file
```yaml
# index.md
---
name: pongotemplate
author: 2p4b 
layouts:  layouts #Path to templates views          default to layouts
public: public #Path to templates public assets     default to public
views: views #Path to templates views               default to views
type: pogo #Template engine to use                  default to html
ext: html #Templates file extensions                default to html
---
Julien Example Pongo Template

```

#### Data variables
These variables share a common interface e.g `Page.Get("title")` or `Site.Get("title", "default")` to get 
frontmatter data and `Page.Content` for the document body same is true for all data variables

- `Site` Site document
- `Page` Current page document
- `Post` is the data submitted successfully by a form 
    - `Post.Form` Post Form document
    - `Post.Doc` Post Data document
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


