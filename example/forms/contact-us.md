---
title: Contact Us
name: $timestamp
redirect: /contact-us

includes: 
    ip: ip
    hostname: hostname
    referer: Referer
    useragent: User-Agent

schema:
    name: required,min=1,max=32
    phone: required,min=1,max=32
    email: required,email
    company: required,min=1,max=255
    about: required,min=1,max=255
content: about
---

Contact Us Form