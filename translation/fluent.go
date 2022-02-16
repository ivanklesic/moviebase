package translation

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/lus/fluent.go/fluent"
	"github.com/lus/fluent.go/fluent/parser"
	"golang.org/x/text/language"
)

// maps ftl filenames to specific language tags
var supportedLanguages = map[string]language.Tag{
	"en": language.English,	
	"de": language.German,
}

var bundles = make(map[string]*fluent.Bundle)

type FluentErrors struct {
	ParserErrors []*parser.Error
	Errors       []error
}

func Translate(languageTag string, message string, contexts ...*fluent.FormatContext) (string, []error, error) {
	return bundles[languageTag].FormatMessage(message, contexts...)
}

// wrapper around the fluent package func
func WithVar(key string, value interface{}) *fluent.FormatContext{
	return fluent.WithVariable(key, value)
}

// wrapper around the fluent package func
func WithVars(vars map[string]interface{}) *fluent.FormatContext{
	return fluent.WithVariables(vars)
}

// wrapper around the fluent package func
func WithFunc(key string, fn fluent.Function) *fluent.FormatContext{
	return fluent.WithFunction(key, fn)
}

// wrapper around the fluent package func
func WithFuncs(fns map[string]fluent.Function) *fluent.FormatContext{
	return fluent.WithFunctions(fns)
}

// needs to be called to setup translations
func ParseTranslations() *FluentErrors {
	langFiles, err := matchTranslations("translation/fluent", ".ftl")

	if err != nil {
		return &FluentErrors{
			Errors: []error{err},
		}
	}

	for lang, translations := range langFiles {
		resource, parserErrors := fluent.NewResource(translations)
		if len(parserErrors) > 0 {
			return &FluentErrors{
				ParserErrors: parserErrors,
			}
		}		

		bundle := fluent.NewBundle(supportedLanguages[lang])

		if bundle != nil {
			errors := bundle.AddResource(resource)
			if len(errors) > 0 {
				return &FluentErrors{
					Errors: errors,
				}
			}
			bundles[lang] = bundle
		}		
	}
	return nil
}

// walks the translation/fluent dir, which contains all translation files, reads every file and saves it to a map
func matchTranslations(root, pattern string) (map[string]string, error) {
	var matches = make(map[string]string)	
	err := filepath.WalkDir(root, func(s string, info os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(info.Name()) == pattern {
			content, err := ioutil.ReadFile(s)
			if err != nil {
				return err
			}
			fileName := strings.TrimSuffix(info.Name(), ".ftl")
			if _, ok := supportedLanguages[fileName] ; ok {
				matches[fileName] = string(content)
			}			
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matches, nil
}
