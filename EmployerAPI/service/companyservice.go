package service

import (
	customerrors "EmployerAPI/custom_errors"
	"EmployerAPI/filters"
	"EmployerAPI/logger"
	"EmployerAPI/models"
	"EmployerAPI/repository"
	"mime/multipart"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ICompanyService interface {
	SaveCompany(jobDetail *models.Company) (*models.Company, error)
	GetCompanies(filterData *models.SearchFilter) (*models.GetCompanyResponse, error)
	SaveImagetoAWS(_file multipart.File, attribute string, companyId, fileName, mimetype string, size int64) error
}

type companyService struct {
	repo      repository.ICompanyRepository
	logObj    logger.ILogger
	s3Service IS3Service
}

var companyServiceObj ICompanyService

func InitCompanyService(repoObj repository.ICompanyRepository, s3Service IS3Service, logObj logger.ILogger) ICompanyService {
	if companyServiceObj == nil {
		companyServiceObj = &companyService{
			repo:      repoObj,
			s3Service: s3Service,
			logObj:    logObj,
		}
	}
	return companyServiceObj
}

func (s *companyService) SaveCompany(companyDetail *models.Company) (*models.Company, error) {
	if companyDetail.Id != "" {
		update := bson.M{"$set": bson.M{
			"name":     companyDetail.Name,
			"email":    companyDetail.Email,
			"phone":    companyDetail.Phone,
			"website":  companyDetail.Website,
			"category": companyDetail.Category,
			"teamsize": companyDetail.TeamSize,
			"about":    companyDetail.About,
		}}
		err := s.repo.UpdateSingle(update, companyDetail.Id)
		if err != nil {
			s.logObj.Printf("Error while updating Employer company, error: %s\n", err.Error())
			return nil, err
		}
	} else {
		companyId, err := s.repo.AddSingle(*companyDetail)
		if err != nil {
			return nil, err
		}
		companyDetail.Id = companyId
	}
	return companyDetail, nil
}

func (s *companyService) SaveImagetoAWS(_file multipart.File, attribute string, companyId, fileName, mimetype string, size int64) error {

	update := bson.M{"$set": bson.M{
		attribute: fileName,
	}}

	err := s.repo.UpdateSingle(update, companyId)
	if err != nil {
		s.logObj.Printf("Error while updating userId: %+v, error: %s\n", companyId, err.Error())
		return &customerrors.UpdateUserException{}
	}
	err = s.s3Service.Put(_file, fileName, mimetype, size)
	if err != nil {
		return err
	}
	return nil
}

func (s *companyService) GetCompanies(filterData *models.SearchFilter) (*models.GetCompanyResponse, error) {
	var filter filters.IFilter = nil
	_filter := bson.M{}
	if filterData != nil {
		if filterData.EmployerId != "" {
			filter = filters.InitEmployerIdFilter(filter, filters.AND, filters.EQUAL, filterData.EmployerId)
		}
		if filterData.Id != "" {
			objectId, err := primitive.ObjectIDFromHex(filterData.Id)
			if err != nil {
				s.logObj.Printf("error while converting id to hex %s, error: %s\n", filterData.Id, err.Error())
				return nil, err
			}
			filter = filters.InitIdFilter(filter, filters.AND, filters.EQUAL, objectId)
		}
	}
	if filter != nil {
		_filter = filter.Build()
	}
	data, err := s.repo.Get(_filter, int64(filterData.PageSize), int64(filterData.PageNumber))
	if err != nil {
		return nil, err
	}
	if filterData.Id != "" && len(data) > 0 {
		data[0].LogoFileName, err = s.GetImageURL(data[0].LogoFileName)
		if err != nil {
			return nil, err
		}
		data[0].CoverFileName, err = s.GetImageURL(data[0].CoverFileName)
		if err != nil {
			return nil, err
		}
	}
	return &models.GetCompanyResponse{
		Company:    data,
		PageSize:   filterData.PageSize,
		PageNumber: filterData.PageNumber,
	}, nil
}

func (s *companyService) GetImageURL(filename string) (string, error) {
	return s.s3Service.GetPreSignerUrl(filename)

}
