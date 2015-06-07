package memory

import (
	"testing"

	eventPackage "github.com/index0h/go-tracker/components"
)

func Test_EventRepository_Interface(t *testing.T) {
	func(event eventPackage.EventRepositoryInterface) {}(&Repository{})
}
