package debug

import (
	"errors"
	"fmt"
	"log"
	"runtime"
)

func GetFuncInfo(backStepsInStack int) string {
	_, fileName, lineNumber, ok := runtime.Caller(backStepsInStack)
	if !ok {
		log.Fatal(errors.New("caller returned ok = false"))
	}
	return fmt.Sprintf(" [line number %d in %s] ", lineNumber, fileName) //todo: add the func name in the future
}
