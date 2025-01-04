package helpers

import (
	"fmt"
	"mime/multipart"
	"reflect"

	"github.com/gin-gonic/gin"
)

func PopulateDynamicStructWithFiles(c *gin.Context, obj interface{}) error {
	v := reflect.ValueOf(obj).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		if field.Type() == reflect.TypeOf((*multipart.FileHeader)(nil)) {
			formKey := fieldType.Tag.Get("form")

			if file, err := c.FormFile(formKey); err == nil {

				if field.CanSet() {
					field.Set(reflect.ValueOf(file))
				} else {
					fieldPointer := reflect.New(field.Type().Elem())
					fieldPointer.Elem().Set(reflect.ValueOf(file))
					field.Set(fieldPointer)
				}
			} else {
				return fmt.Errorf("no file provided for field %s", formKey)
			}
		}
	}

	return nil
}
