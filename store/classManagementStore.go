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
	sqlCourse := "INSERT INTO COURSE (COURSE_TYPE_ID, MAIN_TEACHER, ROOM, START_DATE, END_DATE, START_TIME, END_TIME, STUDY_DAYS, LOCATION) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"
	_, err = tx.Exec(sqlCourse, entity.CourseTypeId, entity.MainTeacher, entity.Room, entity.StartDate, entity.EndDate, rq.StartTime, rq.EndTime, entity.StudyDays, rq.Location)
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

	var userId string
	err = tx.QueryRow("SELECT USER_ID FROM USER WHERE USERNAME = ?", entity.MainTeacher).Scan(&userId)
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

	// insert default main teacher in every class of this course
	sqlQuery := "INSERT INTO CLASS_MANAGER (USER_ID, COURSE_ID, CLASS_ROLE) VALUES (?, ?, ?)"
	tmp, err := tx.Prepare(sqlQuery)
	if err != nil {
		log.WithError(err).Errorf("Failed to prepare SQL statement for CLASS_MANAGER")
		return err
	}
	defer stmt.Close()

	for _, v := range schedule {
		rs, err := stmt.Exec(courseID, rq.StartTime, rq.EndTime, v, entity.Room, "1")
		if err != nil {
			log.WithError(err).Errorf("Failed to insert class into the database")
			return err
		}

		id, err := rs.LastInsertId()
		if err != nil {
			log.WithError(err).Errorf("Failed to get last insert ID")
			return err
		}

		_, err = tmp.Exec(userId, id, "Teacher")
		if err != nil {
			log.WithError(err).Errorf("Failed to insert main teacher into the database")
			return err
		}
	}

	return nil
}

func (c *classManagementStore) GetAllCoursesStore(ctx context.Context) ([]course_class.CourseEntity, error) {
	log.Infof("Get all courses")

	sqlQuery := "SELECT C.*, CONCAT(CT.CODE, COURSE_ID) AS COURSE_NAME, CT.TOTAL_SESSIONS FROM COURSE C join COURSE_TYPE CT ON C.COURSE_TYPE_ID = CT.COURSE_TYPE_ID"
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
		"where C.COURSE_TYPE_ID = CT.COURSE_TYPE_ID AND C.COURSE_ID = ?"

	err := c.db.QueryRowxContext(ctx, sqlQuery, request.CourseId).Scan(&entity.CourseId, &entity.CourseTypeId, &entity.MainTeacher, &entity.Room, &entity.StartDate, &entity.EndDate, &entity.StartTime, &entity.EndTime, &entity.StudyDays, &entity.Location, &entity.CourseName, &entity.TotalSessions)
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
	sqlQuery := "SELECT C.COURSE_ID, CO.COURSE_TYPE_ID, C.START_TIME, C.END_TIME, C.DATE " +
		"FROM CLASS C " +
		"JOIN CLASS_MANAGER CM ON C.CLASS_ID = CM.COURSE_ID " +
		"JOIN COURSE CO ON C.COURSE_ID = CO.COURSE_ID " +
		"WHERE CM.USER_ID = ? " +
		"AND C.DATE >= ? AND C.DATE <= ?"

	var rs []course_class.FromToScheduleEntity
	err := c.db.SelectContext(ctx, &rs, sqlQuery, userId, fromDate, toDate)
	if err != nil {
		log.WithError(err).Errorf("Failed to get class for user %s from the database", userId)
		return nil, err
	}
	return rs, nil
}

func (c *classManagementStore) DeleteCourseById(courseId string, ctx context.Context) error {
	log.Infof("Delete course %s store", courseId)

	// Start a transaction
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		log.Errorf("Error starting transaction: %v", err)
		return err
	}

	defer tx.Rollback()

	sqlQuery := "DELETE FROM COURSE WHERE COURSE_ID = ?"

	_, err = tx.ExecContext(ctx, sqlQuery, courseId)
	if err != nil {
		log.Errorf("Error deleting course: %v", err)
		return err
	}

	sqlQuery = "DELETE FROM CLASS WHERE COURSE_ID = ?"
	_, err = tx.ExecContext(ctx, sqlQuery, courseId)
	if err != nil {
		log.Errorf("Error deleting course: %v", err)
		return err
	}

	// Commit the transaction if everything is successful
	if err := tx.Commit(); err != nil {
		log.Errorf("Error committing transaction: %v", err)
		return err
	}

	log.Infof("Course with ID %s deleted successfully", courseId)
	return nil
}

func (c *classManagementStore) DeleteClassByIdStore(classId string, ctx context.Context) error {
	log.Infof("Start to delete class %s store", classId)
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		log.Errorf("Error starting transaction: %v", err)
		return err
	}

	defer tx.Rollback()

	sqlQuery := "DELETE FROM CLASS WHERE CLASS_ID = ?"
	_, err = tx.ExecContext(ctx, sqlQuery, classId)
	if err != nil {
		log.Errorf("Error deleting class: %v", err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Errorf("Error committing transaction: %v", err)
		return err
	}

	log.Infof("ClassiD with ID %s deleted successfully", classId)
	return nil
}

func (c *classManagementStore) GetAllSessionsByCourseIdStore(courseId string, ctx context.Context) ([]course_class.ClassEntity, error) {
	log.Infof("Get all sessions for course %s", courseId)
	sqlQuery := "SELECT * FROM CLASS WHERE COURSE_ID = ?"

	var rs []course_class.ClassEntity
	err := c.db.SelectContext(ctx, &rs, sqlQuery, courseId)
	if err != nil {
		log.WithError(err).Errorf("Failed to get sessions for course %s from the database", courseId)
		return nil, err
	}
	return rs, nil
}

func (c *classManagementStore) FixCourseInformationStore(rq api_request.ModifyCourseInformation, ctx context.Context) error {
	log.Infof("Fix course information %v", rq)

	// Begin a transaction
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		log.WithError(err).Errorf("Error starting transaction for fix course information %s", rq.CourseId)
		return err
	}

	// UPDATE IN COURSE TABLE
	sqlQuery := "UPDATE COURSE SET MAIN_TEACHER = ?, ROOM = ? WHERE COURSE_ID = ?"

	_, err = tx.ExecContext(ctx, sqlQuery, rq.Teacher, rq.Room, rq.CourseId)
	if err != nil {
		// Rollback the transaction if an error occurs
		tx.Rollback()
		log.WithError(err).Errorf("Error fix course %s information on TABLE COURSE", rq.CourseId)
		return err
	}

	// UPDATE IN CLASS TABLE
	sqlQuery = "UPDATE CLASS SET ROOM = ? WHERE COURSE_ID = ?"
	_, err = tx.ExecContext(ctx, sqlQuery, rq.Room, rq.CourseId)
	if err != nil {
		// Rollback the transaction if an error occurs
		tx.Rollback()
		log.WithError(err).Errorf("Error fix course %s information on TABLE CLASS", rq.CourseId)
		return err
	}

	// Commit the transaction if everything is successful
	err = tx.Commit()
	if err != nil {
		// Handle commit error if needed
		log.WithError(err).Errorf("Error committing transaction for course information update")
		return err
	}

	log.Infof("Successful update course information")
	return nil

}

func (c *classManagementStore) AddNoteStore(noteRequest api_request.AddNoteRequest, add []string, delete []string, ctx context.Context) error {
	log.Infof("Start add note store for class %s", noteRequest.ClassId)
	// Begin a transaction
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		log.WithError(err).Errorf("Error starting transaction for note class information %s", noteRequest)
		return err
	}

	// UPDATE IN COURSE TABLE
	sqlQuery := "UPDATE CLASS SET START_TIME = ?, END_TIME = ?, DATE = ?, NOTE = ? WHERE CLASS_ID = ?"

	_, err = tx.ExecContext(ctx, sqlQuery, noteRequest.StartTime, noteRequest.EndTime, noteRequest.Date, noteRequest.Note, noteRequest.ClassId)
	if err != nil {
		// Rollback the transaction if an error occurs
		tx.Rollback()
		log.WithError(err).Errorf("Error starting transaction for note class information %s", noteRequest)
		return err
	}

	sqlQuery = "INSERT INTO CLASS_MANAGER (USER_ID, COURSE_ID, CLASS_ROLE) VALUES (?, ?, ?)"

	for _, v := range add {
		_, err = tx.ExecContext(ctx, sqlQuery, v, noteRequest.ClassId, "TA")
		if err != nil {
			// Rollback the transaction if an error occurs
			tx.Rollback()
			log.WithError(err).Errorf("Error starting transaction for note class information %s", noteRequest)
			return err
		}
	}

	sqlQuery = "DELETE FROM CLASS_MANAGER WHERE USER_ID = ? AND COURSE_ID = ?"

	for _, v := range delete {
		_, err = tx.ExecContext(ctx, sqlQuery, v, noteRequest.ClassId)
		if err != nil {
			// Rollback the transaction if an error occurs
			tx.Rollback()
			log.WithError(err).Errorf("Error starting transaction for note class information %s", noteRequest)
			return err
		}
	}

	// Commit the transaction if everything is successful
	err = tx.Commit()
	if err != nil {
		// Handle commit error if needed
		log.WithError(err).Errorf("Error committing transaction for note class information")
		return err
	}

	log.Infof("Successful update note class information")
	return nil

}

func (c *classManagementStore) GetTaListInSessionStore(classId int, ctx context.Context) ([]string, error) {
	log.Infof("Start get TA List store for class %s", classId)

	sqlQuery := "SELECT CM.USER_ID FROM CLASS_MANAGER CM JOIN USER U ON CM.USER_ID = U.USER_ID WHERE CM.COURSE_ID = ? AND U.JOB_POSITION = 'TA'"
	var entities []string
	err := c.db.SelectContext(ctx, &entities, sqlQuery, classId)

	if err != nil {
		log.WithError(err).Errorf("Failed to retrieve courses from the database")
		return nil, err
	}

	return entities, nil
}

func (c *classManagementStore) GetCheckInHistoryByCourseIdStore(courseId string, ctx context.Context) ([]course_class.CheckInHistoryEntity, error) {
	log.Infof("Start get check in history List store for course %s", courseId)

	var entities []course_class.CheckInHistoryEntity
	sqlQuery := "SELECT A.USER_ID, A.CLASS_NAME, A.CHECKIN_TIME, A.STATUS FROM ATTENDANCE_HISTORY A JOIN CLASS C ON A.CLASS_NAME = C.CLASS_ID WHERE C.COURSE_ID = ?"
	args := []interface{}{courseId}

	rows, err := c.db.QueryxContext(ctx, sqlQuery, args...)
	if err != nil {
		log.WithError(err).Errorf("Cannot get check in history from the database for course: %s", courseId)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var entity course_class.CheckInHistoryEntity
		if err := rows.Scan(&entity.UserId, &entity.ClassId, &entity.CheckInTime, &entity.Status); err != nil {
			log.WithError(err).Errorf("Error scanning row: %s", err.Error())
			return nil, err
		}
		entities = append(entities, entity)
	}

	if err := rows.Err(); err != nil {
		log.WithError(err).Errorf("Error iterating rows: %s", err.Error())
		return nil, err
	}
	return entities, nil
}
