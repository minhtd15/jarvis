package courseClass

import (
	"context"
	"database/sql"
	batman "education-website"
	api_request "education-website/api/request"
	api_response "education-website/api/response"
	"education-website/entity/course_class"
	"fmt"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

type classService struct {
	classStore batman.ClassStore
}

type ClassServiceCfg struct {
	ClassStore batman.ClassStore
}

func NewClassService(cfg ClassServiceCfg) batman.ClassService {
	return classService{
		classStore: cfg.ClassStore,
	}
}

func (c classService) AddNewClass(request api_request.NewCourseRequest, ctx context.Context) error {
	log.Infof("Start insert new course to db")

	// get total sessions of type course
	totalSessions, courseTypeId, err := c.classStore.GetSessionsByCourseType(request.CourseType, ctx)
	if err != nil {
		log.WithError(err).Errorf("Unable to get total sessions of type class: %s", request.CourseType)
		return err
	}

	startDate, err := time.Parse("2006-01-02", request.StartDate)
	if err != nil {
		log.WithError(err).Errorf("Unable to parse string to date")
		return err
	}
	studyDays := []time.Weekday{}
	for _, v := range request.StudyDays {
		switch strings.ToUpper(v) {
		case "MONDAY":
			studyDays = append(studyDays, time.Monday)
		case "TUESDAY":
			studyDays = append(studyDays, time.Tuesday)
		case "WEDNESDAY":
			studyDays = append(studyDays, time.Wednesday)
		case "THURSDAY":
			studyDays = append(studyDays, time.Thursday)
		case "FRIDAY":
			studyDays = append(studyDays, time.Friday)
		case "SATURDAY":
			studyDays = append(studyDays, time.Saturday)
		case "SUNDAY":
			studyDays = append(studyDays, time.Sunday)
		}
	}

	schedule := c.generateWeeklySchedule(startDate, *totalSessions, studyDays)
	endDate := schedule[len(schedule)-1]

	entity := course_class.CourseEntity{
		CourseTypeId: *courseTypeId,
		MainTeacher:  request.MainTeacher,
		Room:         request.Room,
		StartDate: sql.NullString{
			String: startDate.Format("2006-01-02"),
			Valid:  true,
		},
		EndDate: sql.NullString{
			String: endDate.Format("2006-01-02"),
			Valid:  true,
		},
		StudyDays:     concatenateWeekdays(studyDays),
		TotalSessions: int64(*totalSessions),
	}

	// insert course and all the sessions of that course
	err = c.classStore.InsertNewCourseStore(entity, request, schedule, ctx)
	if err != nil {
		log.WithError(err).Errorf("Unable to insert new course to db")
		return err
	}

	return nil
}

func concatenateWeekdays(studyDays []time.Weekday) string {
	var weekdays []string

	// Chuyển đổi các giá trị time.Weekday thành chuỗi
	for _, day := range studyDays {
		weekdays = append(weekdays, day.String())
	}

	// Kết hợp các chuỗi thành một chuỗi duy nhất, phân cách bằng ","
	concatenatedDays := strings.Join(weekdays, ",")
	return concatenatedDays
}

func (c classService) generateWeeklySchedule(startDate time.Time, totalSessions int, studyDays []time.Weekday) []time.Time {
	schedule := []time.Time{}
	currentDate := startDate

	found := false
	for _, dayOfWeek := range studyDays {
		log.Infof("Current date weekday: %s", currentDate.Weekday())

		if currentDate.Weekday() == dayOfWeek {
			found = true
			break
		}
	}
	if !found {
		// Nếu startDate không thuộc studyDays, di chuyển ngày đến ngày tiếp theo thuộc studyDays
		for {
			currentDate = currentDate.AddDate(0, 0, 1)
			if containsWeekday(studyDays, currentDate.Weekday()) {
				break
			}
		}
	}

	for i := 0; i < totalSessions; {
		// Kiểm tra nếu currentDate nằm trong studyDays, bắt đầu từ currentDate
		if containsWeekday(studyDays, currentDate.Weekday()) {
			// Thêm currentDate vào schedule
			schedule = append(schedule, currentDate)
			i++
			// Di chuyển ngày tiếp theo
			currentDate = currentDate.AddDate(0, 0, 1)
		}

		// Di chuyển đến ngày tiếp theo
		currentDate = currentDate.AddDate(0, 0, 1)
	}

	return schedule
}

func containsWeekday(studyDays []time.Weekday, day time.Weekday) bool {
	for _, d := range studyDays {
		if d == day {
			return true
		}
	}
	return false
}

func (c classService) GetCourseInformationByClassName(request api_request.CourseInfoRequest, ctx context.Context) (*api_response.CourseInfoResponse, error) {
	log.Info("Get class information by class name")

	// get class information
	rs, err := c.classStore.GetCourseInformationStore(request, ctx)
	if err != nil {
		log.WithError(err).Errorf("Error getting clas information by class Name")
		return nil, err
	}

	var startDate, endDate time.Time
	if rs.StartDate.Valid {
		startDate, err = time.Parse("2006-01-02", rs.StartDate.String)
		if err != nil {
			log.WithError(err).Errorf("Cannot parse start time")
			return nil, err
		}
	}

	if rs.EndDate.Valid {
		endDate, err = time.Parse("2006-01-02", rs.EndDate.String)
		if err != nil {
			log.WithError(err).Errorf("Cannot parse start time")
			return nil, err
		}
	}

	var courseStatus string
	currentDate := time.Now()
	if currentDate.Before(startDate) || currentDate.Equal(startDate) {
		courseStatus = "INACTIVE"
	} else if currentDate.After(endDate) {
		courseStatus = "FINISHED"
	} else {
		courseStatus = "ACTIVE"
	}

	return &api_response.CourseInfoResponse{
		CourseId:      int64(rs.CourseId),
		CourseName:    rs.CourseName,
		MainTeacher:   rs.MainTeacher,
		Room:          int64(rs.Room),
		StartDate:     startDate.Format("2006-01-02"),
		EndDate:       endDate.Format("2006-01-02"),
		StudyDays:     rs.StudyDays,
		CourseStatus:  courseStatus,
		TotalSessions: rs.TotalSessions,
	}, nil
}

func (c classService) GetAllCourses(ctx context.Context) ([]api_response.CourseInfoResponse, error) {
	entities, err := c.classStore.GetAllCoursesStore(ctx)
	if err != nil {
		log.WithError(err).Errorf("Unable to get all courses store")
		return nil, err
	}

	rs := make([]api_response.CourseInfoResponse, 0)

	var startDate, endDate time.Time
	for _, v := range entities {
		if v.StartDate.Valid {
			startDate, err = time.Parse("2006-01-02", v.StartDate.String)
			if err != nil {
				log.WithError(err).Errorf("Cannot parse start time")
				return nil, err
			}
		}
		if v.EndDate.Valid {
			endDate, err = time.Parse("2006-01-02", v.EndDate.String)
			if err != nil {
				log.WithError(err).Errorf("Cannot parse start time")
				return nil, err
			}
		}

		currentDate := time.Now()
		var courseStatus string
		if currentDate.Before(startDate) || currentDate.Equal(startDate) {
			courseStatus = "INACTIVE"
		} else if currentDate.After(endDate) {
			courseStatus = "FINISHED"
		} else {
			courseStatus = "ACTIVE"
		}

		tmp := api_response.CourseInfoResponse{
			CourseId:      int64(v.CourseId),
			CourseName:    v.CourseName,
			MainTeacher:   v.MainTeacher,
			Room:          int64(v.Room),
			StartDate:     startDate.Format("2006-01-02"),
			EndDate:       endDate.Format("2006-01-02"),
			StudyDays:     v.StudyDays,
			CourseStatus:  courseStatus,
			TotalSessions: v.TotalSessions,
		}

		rs = append(rs, tmp)
	}

	return rs, nil
}

func (c classService) GetCourseType(ctx context.Context) (map[int]string, error) {
	courseTypeList, err := c.classStore.GetAllCourseType(ctx)
	if err != nil {
		log.WithError(err).Errorf("Unable to get course type from db")
		return nil, err
	}

	rs := make(map[int]string)
	for _, v := range courseTypeList {
		rs[v.CourseTypeId] = v.CourseCode
	}
	log.Infof("Successfully get course type from db")
	return rs, nil
}

func (c classService) GetFromToSchedule(fromDate string, toDate string, userId string, courseType map[int]string, ctx context.Context) ([]api_response.FromToScheduleResponse, error) {
	log.Infof("get %s classes from %s to %s", userId, fromDate, toDate)

	fromToScheduleEntity, err := c.classStore.GetClassFromToDateStore(fromDate, toDate, userId, ctx)
	if err != nil {
		log.WithError(err).Errorf("Error getting class from %s to %s for user %s", fromDate, toDate, userId)
		return nil, err
	}

	var rs []api_response.FromToScheduleResponse
	for _, v := range fromToScheduleEntity {
		tmp := api_response.FromToScheduleResponse{
			CourseId:   v.CourseId,
			CourseCode: courseType[int(v.CourseTypeId)],
			StartTime:  v.StartTime,
			EndTime:    v.EndTime,
			Date:       v.Date,
		}
		tmp.CourseName = fmt.Sprintf("%s%d", tmp.CourseCode, v.CourseId)
		rs = append(rs, tmp)
	}

	return rs, nil
}

func (c classService) GetClassInformationByClassId(classId string, ctx context.Context) {

}

func (c classService) DeleteCourseByCourseId(courseId string, ctx context.Context) error {
	return c.classStore.DeleteCourseById(courseId, ctx)
}

func (c classService) DeleteClassByClassId(rq api_request.DeleteClassInfo, ctx context.Context) error {
	return c.classStore.DeleteClassByIdStore(rq.ClassId, ctx)
}

func (c classService) GetAllSessionsByCourseIdService(courseId string, ctx context.Context) ([]api_response.ClassResponse, error) {
	log.Infof("Start service get sessions for course: %s", courseId)

	rs, err := c.classStore.GetAllSessionsByCourseIdStore(courseId, ctx)
	if err != nil {
		log.WithError(err).Errorf("Error getting sessions store for course %s", courseId)
		return nil, err
	}

	var result []api_response.ClassResponse
	for _, v := range rs {
		tmp := api_response.ClassResponse{
			ClassId:   v.ClassId,
			StartTime: v.StartTime,
			EndTime:   v.EndTime,
			Date:      v.Date.String,
			Room:      v.Room,
			TypeClass: v.TypeClass,
			Note:      v.Note.String,
		}

		intNumber, err := strconv.ParseInt(v.ClassId, 10, 64)
		if err != nil {
			// Handle the error if the conversion fails
			fmt.Println("Error converting string to int:", err)
			return nil, err
		}

		TaList, err := c.classStore.GetTaListInSessionStore(int(intNumber), ctx)
		if err != nil {
			log.WithError(err).Errorf("Error getting TA List store for classId %s", v.ClassId)
			return nil, err
		}
		tmp.TaList = TaList
		result = append(result, tmp)
	}

	return result, nil
}

func (c classService) FixCourseInformationService(rq api_request.ModifyCourseInformation, ctx context.Context) error {
	return c.classStore.FixCourseInformationStore(rq, ctx)
}

func (c classService) AddNoteService(noteRequest api_request.AddNoteRequest, add []string, delete []string, ctx context.Context) error {
	return c.classStore.AddNoteStore(noteRequest, add, delete, ctx)
}

func (c classService) GetTAListService(classId int, ctx context.Context) ([]string, error) {
	return c.classStore.GetTaListInSessionStore(classId, ctx)
}

func (c classService) GetCheckInHistoryByCourseId(courseId string, ctx context.Context) ([]api_response.CheckInHistory, error) {
	log.Infof("Start to get check in history")

	rs, err := c.classStore.GetCheckInHistoryByCourseIdStore(courseId, ctx)
	if err != nil {
		log.WithError(err).Errorf("Error getting check in history store for course %s", courseId)
		return nil, err
	}
	var result []api_response.CheckInHistory
	for _, v := range rs {
		tmp := api_response.CheckInHistory{
			UserId:      v.UserId,
			ClassId:     v.ClassId,
			CheckInTime: v.CheckInTime,
			Status:      v.Status,
		}
		result = append(result, tmp)
	}
	return result, nil
}
