package memory

import (
	"testing"

	eventPackage "github.com/index0h/go-tracker/components"
	//	"github.com/index0h/go-tracker/entities"
	//	interfaceUUID "github.com/index0h/go-tracker/uuid"
	//	uuidDriver "github.com/index0h/go-tracker/uuid/driver"
	//	"github.com/index0h/go-tracker/entities"
	//	"github.com/stretchr/testify/assert"
)

func Test_EventRepository_Interface(t *testing.T) {
	func(event eventPackage.EventRepositoryInterface) {}(&Repository{})
}
