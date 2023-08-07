package helper

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/afero"
	"os"
	"regexp"
	"strings"
)

func ValidateAndGetErrFields(obj any) []string {
	err := binding.Validator.ValidateStruct(obj)

	if err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			fields := make([]string, len(validationErrors))
			for index, e := range validationErrors {
				fields[index] = e.Field()
			}
			return fields
		}

		return []string{err.Error()}
	}

	return nil
}

func ParseConfigFileEnv(path string) ([]byte, error) {
	file, err := afero.ReadFile(afero.NewOsFs(), path)
	if err != nil {
		return nil, err
	}

	return ParseConfigContentEnv(file), nil
}

func ParseConfigContentEnv(b []byte) []byte {
	content := string(b)
	r, _ := regexp.Compile(`(\$\{[^$]+\})`)
	matches := r.FindAllString(content, -1)
	for _, val := range matches {
		tmpVal := strings.Replace(val, "${", "", 1)
		tmpVal = strings.Replace(tmpVal, "}", "", 1)
		tmpArr := strings.SplitN(tmpVal, "-", 2)
		envVal := os.Getenv(tmpArr[0])
		if envVal == "" && len(tmpArr) == 2 {
			envVal = tmpArr[1]
		}
		content = strings.Replace(content, val, envVal, 1)
	}
	return []byte(content)
}
