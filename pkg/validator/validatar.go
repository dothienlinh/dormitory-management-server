package validator

import (
	"time"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var formats = []string{
	time.Layout,
	time.ANSIC,
	time.UnixDate,
	time.RubyDate,
	time.RFC822,
	time.RFC822Z,
	time.RFC850,
	time.RFC1123,
	time.RFC1123Z,
	time.RFC3339,
	time.RFC3339Nano,
	time.Kitchen,
	time.Stamp,
	time.StampMilli,
	time.StampMicro,
	time.StampNano,
	time.DateTime,
	time.DateOnly,
	time.TimeOnly,
	"2006/01/02",
}

func CustomValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("validdate", validateDate)
		v.RegisterValidation("gtefield", validateDateAfter)
	} else {
		panic("Validator engine is not of type *validator.Validate")
	}
}

func isDateTime(dateString string) bool {
	var err error
	for _, format := range formats {
		_, err = time.Parse(format, dateString)
		if err == nil {
			return true
		}
	}

	return false
}

func validateDate(fl validator.FieldLevel) bool {
	dateStr := fl.Field().String()

	return isDateTime(dateStr)
}

func validateDateAfter(fl validator.FieldLevel) bool {
	currentField := fl.Field().String()
	compareField := fl.Parent().FieldByName(fl.Param()).String()

	for _, format := range formats {
		currentDate, err1 := time.Parse(format, currentField)
		compareDate, err2 := time.Parse(format, compareField)

		if err1 != nil || err2 != nil {
			continue
		} else if currentDate.IsZero() || compareDate.IsZero() {
			return false
		}

		return currentDate.After(compareDate) || currentDate.Equal(compareDate)
	}

	return false
}
