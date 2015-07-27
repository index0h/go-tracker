package http

import (
	"github.com/index0h/go-tracker/entities"
	httpLib "net/http"
	"net/url"
	"regexp"
)

var regexTemplate *regexp.Regexp

func init() {
	regexTemplate = regexp.MustCompile(`^http\.(get|post)\..+`)
}

type Http struct {
	priority       int
	useEventFields bool
	useVisitFields bool
}

func New(priority int, useEventFields bool, useVisitFields bool) *Http {
	return &Http{priority: priority, useEventFields: useEventFields, useVisitFields: useVisitFields}
}

func (processor *Http) GetPriority() int {
	return processor.priority
}

func (processor *Http) Process(
	flash *entities.Flash,
	event *entities.Event,
	visit *entities.Visit,
) *entities.Flash {
	var (
		eventFields entities.Hash
		visitFields entities.Hash
		mustUpdate  bool
	)

	if processor.useEventFields {
		eventFields = flash.EventFields()
		startLen := len(eventFields)

		processor.processMap(eventFields)

		if len(eventFields) > startLen {
			mustUpdate = true
		}
	}

	if processor.useVisitFields {
		visitFields = flash.VisitFields()
		startLen := len(visitFields)

		processor.processMap(visitFields)

		if len(visitFields) > startLen {
			mustUpdate = true
		}
	}

	if mustUpdate {
		if !processor.useVisitFields {
			visitFields = flash.VisitFields()
		}

		if !processor.useEventFields {
			eventFields = flash.EventFields()
		}

		flash, _ = entities.NewFlashFromRaw(
			flash.FlashID(),
			flash.VisitID(),
			flash.EventID(),
			flash.Timestamp(),
			visitFields,
			eventFields,
		)
	}

	return flash
}

func (processor *Http) processMap(fields map[string]string) {
	var (
		requestURL  *url.URL
		queryValues url.Values
		err         error
	)

	for key, value := range fields {
		found := regexTemplate.FindStringSubmatch(key)
		if len(found) == 0 {
			continue
		}

		switch found[1] {
		case "get":
			if _, err = httpLib.Get(value); err != nil {
				fields[key+".error"] = err.Error()
			}
		case "post":
			if requestURL, err = url.Parse(value); err != nil {
				fields[key+".error"] = err.Error()
				break
			}

			if queryValues, err = url.ParseQuery(requestURL.RawQuery); err != nil {
				fields[key+".error"] = err.Error()
				break
			}

			requestURL.RawQuery = ""

			if _, err = httpLib.PostForm(requestURL.RequestURI(), queryValues); err != nil {
				fields[key+".error"] = err.Error()
			}
		}
	}
}
