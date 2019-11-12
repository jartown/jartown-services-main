package service

import (
	"fmt"

	common "github.com/singkorn/jartown-services-common"
	"github.com/singkorn/jartown-services-common/derror"
	"github.com/singkorn/jartown-services-main/database"
	"github.com/singkorn/jartown-services-main/model"
)

type ListReviewService interface {
	GetListReview(doer interface{}, id string) (model.ListReview, error)
}

type listReviewServiceImpl struct {
	db   database.ListReviewDatabase
	time common.TimeSource
}

func ListReviewServiceInit(db database.ListReviewDatabase) ListReviewService {
	return &listReviewServiceImpl{
		db:   db,
		time: &common.RealTime{},
	}
}

func (s *listReviewServiceImpl) GetListReview(doer interface{}, id string) (model.ListReview, error) {
	if id == "" {
		return model.ListReview{}, derror.ErrorDebug(common.ErrInputValidationFailed, "get listReview: blank ID")
	}
	listReview, err := s.db.GetListReview(id)
	if err != nil {
		return model.ListReview{}, derror.ErrorDebug(err, fmt.Sprintf("get listReview: listReview=%v", listReview))
	}
	return listReview, nil
}
