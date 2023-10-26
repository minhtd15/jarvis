package store

import (
	"context"
	"database/sql"
	api_request "education-website/api/request"
	"education-website/entity/course_class"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

type classManagementStore struct {
	db *sqlx.DB
}

type ClassManagementStoreCfg struct {
	Db *sqlx.DB
}

func NewClassManagementStore(classManagementStoreCfg ClassManagementStoreCfg) *classManagementStore {
	return &classManagementStore{
		db: classManagementStoreCfg.Db,
	}
}

func (c *classManagementStore) InsertNewClassStore(request api_request.NewClassRequest, ctx context.Context) error {
	return nil
}

func (c *classManagementStore) GetCourseInformationStore(request api_request.CourseInfoRequest, ctx context.Context) (course_class.CourseEntity, error) {
	log.Info("Get class information from database")

	entity := course_class.CourseEntity{}
	sqlQuery := "select C.* , CONCAT(CT.CODE, C.COURSE_ID) AS CLASS_NAME, CT.TOTAL_SESSIONS " +
		"FROM COURSE C join COURSE_TYPE CT " +
		"where c.COURSE_TYPE_ID = CT.COURSE_TYPE_ID AND C.COURSE_ID = ?"

	err := c.db.QueryRowxContext(ctx, sqlQuery, request.CourseId).Scan(&entity.CourseId, &entity.CourseTypeId, &entity.MainTeacher, &entity.Room, &entity.StartDate, &entity.EndDate, &entity.StudyDays, &entity.CourseName, &entity.TotalSessions)
	if err != nil {
		if err == sql.ErrNoRows {
			//log.WithError(err).Errorf("Cannot find user with class name: %s", userName)
			return entity, err
		}
		log.WithError(err).Errorf("Cannot get info from database for user: %s", request.CourseId)
		return entity, err
	}
	return entity, nil
}

func (c *classManagementStore) GetTeacherStore(request api_request.CourseInfoRequest, ctx context.Context) (string, []string, error) {
	log.Info("Get teacher information")

	var teacherName string
	sqlQuery := "SELECT CM.CLASS_ROLE from CLASS_MANAGER cm where cm.CLASS_ID = ? AND CM.CLASS_ROLE = 'Teacher'"
	err := c.db.QueryRowxContext(ctx, sqlQuery, request.CourseId).Scan(teacherName)
	if err != nil {
		if err == sql.ErrNoRows {
			log.WithError(err).Errorf("Cannot find Teacher with class name: %s", request.CourseId)
			return "", []string{}, err
		}
		log.WithError(err).Errorf("Cannot get info from database for user: %s", request.CourseId)
		return "", []string{}, err
	}

	var teachingAssistant []string
	sqlQuery = "SELECT CM.CLASS_ROLE from CLASS_MANAGER cm where cm.CLASS_ID = ? AND CM.CLASS_ROLE = 'TA'"
	err = c.db.QueryRowxContext(ctx, sqlQuery, request.CourseId).Scan(teachingAssistant)
	if err != nil {
		if err == sql.ErrNoRows {
			log.WithError(err).Errorf("Cannot find TA with class name: %s", request.CourseId)
			return "", []string{}, err
		}
		log.WithError(err).Errorf("Cannot get info from database for user: %s", request.CourseId)
		return "", []string{}, err
	}

	return teacherName, teachingAssistant, nil
}
