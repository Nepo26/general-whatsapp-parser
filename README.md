# Whatsapp Scrapper
This project started from the need to scrape whatsapp web for a family competition and my need to create Go code.

## Todo List
### Necessary
- [x] Open Whatsapp
- [x] Wait QR Code
- [x] Wait Login
- [x] List Profiles
- [x] Go to Specific Group
- [x] List every Group Member and its info
- [ ] Gatter text messages from all members from the screen
- [ ] Scroll to get messages from all time
- [ ] Insert every message to a database
- [ ] Send every image and audio to an S3 API

### Not necessary but nice to have
- [ ] Beautiful and real time Front-end 

---
#  :book: Learnings
Things that I've learned so far.

## XPath
This is navigate xml documents but it seems like we're able to query html with it. And that really makes things easier.
We are able to look for specific classes, certain elements in the html, tags that contain some specific text. It's good.

As I'm not proficient in said language I use a [devhints cheatcheet](https://devhints.io/xpath) and the following article helped me sometimes:
- [XPath Contains: Text, Following Sibling & Ancestor in Selenium](https://www.guru99.com/using-contains-sbiling-ancestor-to-find-element-in-selenium.html#1)

## Whatsapp Web
I disliked the whole random class names of whatsapp, but it's understandable as it probably helps its developers not conflict class names and et cetera.
Looking through their html I saw that they used a lot the tag `data-testid` which is helping me navigate around the code so far.


## :hamster: Go libraries
As I'm not a go developer per se, I needed to familiarize myself with a few libraries and topics.

### time
A built-in library to better handle time and timestamps. 

Mainly used [this article](https://blog.boot.dev/golang/golang-date-time/) to learn how to use it.

### browserdb
The library that I'm using to connect to the Chrome Devtools. It's know to be limited compared to selenium but in the go world is the most maintaned
so far.

I'm using their examples a lot. To see what I'm able to do and to improve upon their code. I hope to someday colaborate with some code there.
The main examples that I look were:
- [Logic](https://github.com/chromedp/examples/blob/master/logic/main.go): Something that resembles more a poject and how to grow better.
- [Click](https://github.com/chromedp/examples/blob/master/click/main.go): To learn how to use the framework to click around.
- [Screenshot](https://github.com/chromedp/examples/blob/master/screenshot/main.go): So I can screenshot the QR code and more.
- [Text](https://github.com/chromedp/examples/blob/master/text/main.go): Learning how to extract text from elements.

## Go Language

### Enums
Go don't have enums natively, so we have to implement it ourselves.
- https://www.sohamkamani.com/golang/enums/

## Facebook Whatsapp API
- [Reference Docs](https://developers.facebook.com/docs/whatsapp/on-premises/reference/messages/)


