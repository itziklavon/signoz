package web_filters

import (
	"goapm/logger"
	"github.com/gofiber/fiber/v2"
	"regexp"
)

var log = logger.LOGGER
var xssFilters = map[string]*regexp.Regexp{
	"caseInsensitiveOnlyScriptEnd":   regexp.MustCompile("</script>"),
	"caseInsensitiveOnlyScriptStart": regexp.MustCompile("<script(.*?)>"),
	"caseInsensitive":                regexp.MustCompile("<script>(.*?)</script>"),
	"caseInsensitiveMultiLine":       regexp.MustCompile("src[\n]*=[\n]*\\'(.*?)\\'"),
	"caseInsensitiveMultiLineDotAll": regexp.MustCompile("src[\n]*=[\n]*\\\"(.*?)\\\\"),
	"caseInsensitiveEval":            regexp.MustCompile("eval\\((.*?)\\)"),
	"caseInsensitiveExpression":      regexp.MustCompile("expression\\((.*?)\\)"),
	"caseInsensitiveVbScript":        regexp.MustCompile("vbscript:"),
	"caseInsensitiveOnLoad":          regexp.MustCompile("onload(.*?)="),
}

type XssFilterResponse struct {
	Response string `json:"response"`
}

// NewXssFilter middleware which intercepts requests, check for malicious content, and  stop request if not valid
func NewXssFilter() fiber.Handler {

	return func(c *fiber.Ctx) error {
		body := c.Body()
		if body != nil && len(body) > 0 && !validateXss(body) {
			c.Response().StatusCode()
			return c.JSON(XssFilterResponse{
				Response: "invalid_input",
			})
		}
		c.Request().SetBody(body)
		return c.Next()
	}
}

func validateXss(msg []byte) bool {
	strBody := string(msg)
	for _, v := range xssFilters {
		if v.MatchString(strBody) {
			return false
		}
	}
	return true
}
