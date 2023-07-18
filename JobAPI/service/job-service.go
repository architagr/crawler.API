package service

import (
	"JobAPI/filters"
	"JobAPI/models"
	"JobAPI/repository"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

type IJobService interface {
	GetJobs(filter *models.JobFilter, pageSize, pageNumber int16) (*models.GetJobResponse, error)
	GetJobDetail(id string) (*models.JobDetails, error)
	GetCourses(keywords string) (*[]models.Courses, error)
}

type jobService struct {
	repo repository.IJobDetailsRepository
}

var jobServiceObj IJobService

func InitJobService(repoObj repository.IJobDetailsRepository) IJobService {
	if jobServiceObj == nil {
		jobServiceObj = &jobService{
			repo: repoObj,
		}
	}
	return jobServiceObj
}

func (svc *jobService) GetJobs(filterData *models.JobFilter, pageSize, pageNumber int16) (*models.GetJobResponse, error) {
	var filter filters.IFilter = nil
	_filter := bson.M{}
	if filterData != nil {
		if filterData.Location != "" {
			filter = filters.InitLocationFilter(filter, filters.AND, filters.EQUAL, filterData.Location)
		}
		if filterData.Keywords != "" {
			filter = filters.InitTitleFilter(filter, filters.OR, filters.EQUAL, filterData.Keywords)
			filter = filters.InitCompanynameFilter(filter, filters.OR, filters.EQUAL, filterData.Keywords)
		}
	}
	if filter != nil {
		_filter = filter.Build()
	}
	data, err := svc.repo.GetJob(_filter, pageSize, pageNumber)
	if err != nil {
		return nil, err
	}
	return &models.GetJobResponse{
		Jobs:       data,
		PageSize:   pageSize,
		PageNumber: pageNumber,
	}, nil
}

func (svc *jobService) GetJobDetail(id string) (*models.JobDetails, error) {
	data, err := svc.repo.GetJobDetail(id)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (svc *jobService) GetCourses(keywords string) (*[]models.Courses, error) {
	clientID := "gWCbMZWTHiLZuIflw6BwR9B6mhOyHLKPZip14tVb"
	clientSecret := "QsMNXAmnq6JDTzOqsukTm36tSvJtd2B0FnQ52ONljlRN4R9mTXJFJCUiLuOnEs7Jbru6WAYNfUEcwBaf0AnZSfbfLz4A4MrOzhco32Jmtr2bpJiE1z5zI7yl7bkylHvT"

	// Encode the client ID and client secret key as Base64
	authToken := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))

	// Construct the request URL
	url := "https://www.udemy.com/api-2.0/courses/?search=" + keywords

	// Create a new HTTP client
	client := &http.Client{}

	// Create a new GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}

	// Add the Base64-encoded token to the Authorization header
	req.Header.Add("Authorization", "Basic "+authToken)

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return nil, err
	}

	var result models.CourseResponse

	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return nil, err
	}
	return &result.Results, nil
}
