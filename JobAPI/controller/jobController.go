package controller

import (
	"JobAPI/models"
	"JobAPI/service"
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type IJobController interface {
	Test(c *fiber.Ctx) error
	GetJobs(c *fiber.Ctx) error
}

type jobController struct {
	service service.IJobService
}

var jobControllerObj IJobController

func InitJobController(jobServiceObj service.IJobService) IJobController {
	if jobControllerObj == nil {
		jobControllerObj = &jobController{
			service: jobServiceObj,
		}
	}
	return jobControllerObj
}

func (ctlr *jobController) Test(c *fiber.Ctx) error {
	name := c.Params("name")
	obj := fmt.Sprintf("Hello, %s!", name)
	return c.JSON(obj)
}

func (ctlr *jobController) GetJobs(c *fiber.Ctx) error {
	//var pageSize, pageNumber int64 = 10, 0 // todo: get this from querystring

	//get params from body
	filter := new(models.JobFilter)
	if err := c.BodyParser(filter); err != nil {
		return err
	}
	_filter, err := json.Marshal(filter)
	fmt.Println("filter: " + string(_filter))

	response, err := ctlr.service.GetJobs(filter, filter.PageSize, filter.PageNumber)
	if err != nil {
		return err
	}
	return c.JSON(response)
}
