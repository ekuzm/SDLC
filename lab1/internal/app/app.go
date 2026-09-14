package app

import (
	"fmt"

	"fyne.io/fyne/v2/app"

	"github.com/ekuzm/sdlc-lab1/internal/controller"
	"github.com/ekuzm/sdlc-lab1/internal/model"
	"github.com/ekuzm/sdlc-lab1/internal/view"
)

func Run() error {
	application := app.NewWithID("by.bsuir.sdlc.mortgage-calculator")
	mortgageModel := model.NewMortgage()
	mortgageController := controller.NewMortgage(mortgageModel)
	mortgageView := view.NewMortgage(application, mortgageController)

	mortgageModel.Subscribe(mortgageView)
	mortgageView.ShowAndRun()

	if application.Driver() == nil {
		return fmt.Errorf("application driver is not initialized")
	}

	return nil
}
