package elastic

import (
	"testing"

	"github.com/index0h/go-tracker/dao"
)

func Test_EventRepository_Interface(t *testing.T) {
	func(event dao.VisitRepositoryInterface) {}(&VisitRepository{})
}