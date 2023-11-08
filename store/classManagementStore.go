package store

import (
	"context"
	"database/sql"
	api_request "education-website/api/request"
	"education-website/entity/course_class"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"time"
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

func (c *classManagementStore) InsertNewCourseStore(entity course_class.CourseEntity, rq api_request.NewCourseRequest, schedule []time.Time, ctx context.Context) error {
	log.Info("insert new class and course to db")

	tx, err := c.db.Begin()
	if err != nil {
		log.WithError(err).Errorf("Failed to begin transaction")
		return err
	}

	defer func() {
		if err != nil {
			log.WithError(err).Errorf("Rolling back transaction")
			tx.Rollback()
			return
		}
		err = tx.Commit()
		if err != nil {
			log.WithError(err).Errorf("Error committing transaction")
		}
	}()

	// Insert into COURSE table
	sqlCourse := "INSERT INTO COURSE (COURSE_TYPE_ID, MAIN_TEACHER, ROOM, START_DATE, END_DATE, STUDY_DAYS, LOCATION) VALUES (?, ?, ?, ?, ?, ?, ?)"
	_, err = tx.Exec(sqlCourse, entity.CourseTypeId, entity.MainTeacher, entity.Room, entity.StartDate, entity.EndDate, entity.StudyDays, rq.Location)
	if err != nil {
		log.WithError(err).Errorf("Failed to insert course into the database")
		return err
	}

	// Get the last inserted ID
	var courseID int64
	err = tx.QueryRow("SELECT LAST_INSERT_ID()").Scan(&courseID)
	if err != nil {
		log.WithError(err).Errorf("Failed to get last insert ID")
		return err
	}

	// Insert into CLASS table for each schedule entry
	sqlClass := "INSERT INTO CLASS (COURSE_ID, START_TIME, END_TIME, DATE, ROOM, TYPE_CLASS) VALUES (?, ?, ?, ?, ?, ?)"
	stmt, err := tx.Prepare(sqlClass)
	if err != nil {
		log.WithError(err).Errorf("Failed to prepare SQL statement for CLASS")
		return err
	}
	defer stmt.Close()

	for _, v := range schedule {
		_, err := stmt.Exec(courseID, rq.StartTime, rq.EndTime, v, entity.Room, rq.TypeCourseCode)
		if err != nil {
			log.WithError(err).Errorf("Failed to insert class into the database")
			return err
		}
	}

	var userId string
	err = tx.QueryRow("SELECT USER_ID FROM USER WHERE USERNAME = ?", entity.MainTeacher).Scan(&userId)
	if err != nil {
		log.WithError(err).Errorf("Failed to get last insert ID")
		return err
	}

	// insert default main teacher in every class of this course
	sqlQuery := "INSERT INTO CLASS_MANAGER VALUES (?, ?, ?, ?)"
	stmt, err = tx.Prepare(sqlQuery)
	if err != nil {
		log.WithError(err).Errorf("Failed to prepare SQL statement for CLASS")
		return err
	}
	defer stmt.Close()

	for _, v := range schedule {
		_, err := stmt.Exec(userId, courseID, "Teacher", v.Format("2006-01-02"))
		if err != nil {
			log.WithError(err).Errorf("Failed to insert main teacher into the database")
			return err
		}
	}

	return nil
}

func (c *classManagementStore) GetAllCoursesStore(ctx context.Context) ([]course_class.CourseEntity, error) {
	log.Infof("Get all courses")

	sqlQuery := "SELECT C.*, CONCAT(CT.CODE, COURSE_ID) AS COURSE_NAME FROM COURSE C join COURSE_TYPE CT ON C.COURSE_TYPE_ID = CT.COURSE_TYPE_ID"
	var entities []course_class.CourseEntity
	err := c.db.SelectContext(ctx, &entities, sqlQuery)

	if err != nil {
		log.WithError(err).Errorf("Failed to retrieve courses from the database")
		return nil, err
	}

	return entities, nil

}

func (c *classManagementStore) GetSessionsByCourseType(typeCourseCode string, ctx context.Context) (*int, *int, error) {
	log.Infof("Get total sessions of course %s", typeCourseCode)

	var totalSessions int
	var courseTypeId int
	sqlQuery := "SELECT TOTAL_SESSIONS, COURSE_TYPE_ID FROM COURSE_TYPE WHERE CODE = ?"

	err := c.db.QueryRowxContext(ctx, sqlQuery, typeCourseCode).Scan(&totalSessions, &courseTypeId)
	if err != nil {
		if err == sql.ErrNoRows {
			log.WithError(err).Errorf("No row in db for course %s", typeCourseCode)
			return nil, nil, err
		}
		log.WithError(err).Errorf("Cannot find total sessions for course %s", typeCourseCode)
		return nil, nil, err
	}

	return &totalSessions, &courseTypeId, err
}

func (c *classManagementStore) GetCourseInformationStore(request api_request.CourseInfoRequest, ctx context.Context) (course_class.CourseEntity, error) {
	log.Info("Get class information from database")

	entity := course_class.CourseEntity{}
	sqlQuery := "select C.* , CONCAT(CT.CODE, C.COURSE_ID) AS COURSE_NAME, CT.TOTAL_SESSIONS " +
		"FROM COURSE C join COURSE_TYPE CT " +
		"where c.COURSE_TYPE_ID = CT.COURSE_TYPE_ID AND C.COURSE_ID = ?"

	err := c.db.QueryRowxContext(ctx, sqlQuery, request.CourseId).Scan(&entity.CourseId, &entity.CourseTypeId, &entity.MainTeacher, &entity.Room, &entity.StartDate, &entity.EndDate, &entity.StartTime, &entity.EndTime, &entity.StudyDays, &entity.CourseName, &entity.Location, &entity.TotalSessions)
	if err != nil {
		if err == sql.ErrNoRows {
			log.WithError(err).Errorf("Cannot find course %s information", request.CourseId)
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

func (c *classManagementStore) GetAllCourseType(ctx context.Context) ([]course_class.CourseTypeEntity, error) {
	log.Infof("Get all courses type")

	sqlQuery := "SELECT * FROM COURSE_TYPE"
	var entities []course_class.CourseTypeEntity
	err := c.db.SelectContext(ctx, &entities, sqlQuery)

	if err != nil {
		log.WithError(err).Errorf("Failed to retrieve courses from the database")
		return nil, err
	}

	return entities, nil
}

func (c *classManagementStore) GetClassFromToDateStore(fromDate string, toDate string, userId string, ctx context.Context) ([]course_class.FromToScheduleEntity, error) {
	log.Infof("Get all classes for user %s from %s, to %s", userId, fromDate, toDate)
	sqlQuery := "SELECT c.COURSE_ID, co.COURSE_TYPE_ID, c.START_TIME, c.END_TIME, c.DATE " +
		"FROM CLASS c " +
		"JOIN CLASS_MANAGER cm ON c.COURSE_ID = cm.COURSE_ID " +
		"JOIN COURSE co ON c.COURSE_ID = co.COURSE_ID " +
		"WHERE cm.USER_ID = ? " +
		"AND c.DATE >= ? AND c.DATE <= ?"

	var rs []course_class.FromToScheduleEntity
	err := c.db.SelectContext(ctx, &rs, sqlQuery, userId, fromDate, toDate)
	if err != nil {
		log.WithError(err).Errorf("Failed to get class for user %s from the database", userId)
		return nil, err
	}
	return rs, nil
}
