package store

import (
	"context"
	"database/sql"
	"education-website/api/request"
	"education-website/api/response/qlda"
	"education-website/entity/sports"
	"fmt"
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

func (c *classManagementStore) GetSportsStore(ctx context.Context) ([]sports.SportsEntity, error) {
	log.Infof("Get sports store")

	sqlQuery := "SELECT SPORT_ID, SPORT_NAME, SPORT_URL FROM SPORTS"
	var entities []sports.SportsEntity
	err := c.db.SelectContext(ctx, &entities, sqlQuery)

	if err != nil {
		log.WithError(err).Errorf("Failed to retrieve sports from the database")
		return nil, err
	}

	return entities, nil
}

func (c *classManagementStore) UploadImageStore(imageFile []byte, sportId int, ctx context.Context) error {
	log.Infof("Upload image store")

	sqlQuery := "UPDATE SPORTS SET IMAGE = ? WHERE SPORT_ID = ?"
	_, err := c.db.ExecContext(ctx, sqlQuery, imageFile, sportId)
	if err != nil {
		log.WithError(err).Errorf("Failed to upload image to the database")
		return err
	}

	return nil
}

func (c *classManagementStore) CreateSchemaStore(ctx context.Context, createSchemaSql string, request request.CreateSchemaRequest) error {
	tx1, err := c.db.BeginTxx(ctx, nil)
	if err != nil {
		log.WithError(err).Errorf("Failed to begin transaction")
		return err
	}

	_, err = tx1.ExecContext(ctx, createSchemaSql)
	if err != nil {
		log.WithError(err).Errorf("Failed to create schema in the database")
		err := tx1.Rollback()
		if err != nil {
			return err
		}
		return err
	}

	if err := tx1.Commit(); err != nil {
		log.WithError(err).Errorf("Failed to commit transaction")
		return err
	}

	tx1, err = c.db.BeginTxx(ctx, nil)
	if err != nil {
		log.WithError(err).Errorf("Failed to begin transaction")
		return err
	}

	_, err = tx1.ExecContext(ctx, fmt.Sprintf("SET search_path TO %s", request.CityId))
	log.Infof("set search path to %s", request.CityId)
	if err != nil {
		log.WithError(err).Errorf("Failed to set search path")
		tx1.Rollback()
		return err
	}

	err = executeSqlFile(c.db.DB, request, tx1)
	if err != nil {
		log.WithError(err).Errorf("Failed to execute sql file")
		tx1.Rollback()
		return err
	}

	if err := tx1.Commit(); err != nil {
		log.WithError(err).Errorf("Failed to commit transaction")
		return err
	}
	return nil
}

func executeSqlFile(db *sql.DB, request request.CreateSchemaRequest, tx *sqlx.Tx) error {
	createSchemaSql := fmt.Sprintf("CREATE TABLE %s.BILL (BILL_ID SERIAL PRIMARY KEY, STATUS VARCHAR(50), PAYMENT_METHOD VARCHAR(50), BANK_ACCOUNT VARCHAR(100), REFUND_AMOUNT DECIMAL(10, 2));"+
		"CREATE TABLE %s.BOOKING (BOOKING_ID SERIAL PRIMARY KEY, TICKET_TYPE_ID INT, SUB_EVENT_ID INT, BILL_ID INT, USER_ID INT, BOOKING_STATUS VARCHAR(50), CUSTOMER_EMAIL VARCHAR(100), CUSTOMER_NAME VARCHAR(100), CUSTOMER_PHONE_NUMBER VARCHAR(15), BOOKING_QUANTITY INT);"+
		"CREATE TABLE %s.USER_MEMBERSHIP (USER_MEMBERSHIP_ID SERIAL PRIMARY KEY, USER_ID INT, MEMBERSHIP_RANK VARCHAR(50), POINTS INT);"+
		"CREATE TABLE %s.USER (USER_ID SERIAL PRIMARY KEY, USER_FULL_NAME VARCHAR(100), USER_EMAIL VARCHAR(100), USER_PHONE_NUMBER VARCHAR(15), TENANT_ID INT, LOCATION VARCHAR(100), ROLE VARCHAR(50));",
		request.CityId, request.CityId, request.CityId, request.CityId)

	_, err := db.Exec(createSchemaSql)
	if err != nil {
		tx.Rollback()
		log.WithError(err).Errorf("Failed to create table in schema in the database")
		return err
	}

	//createOwner := fmt.Sprintf("INSERT INTO %s.USER (USER_FULL_NAME, USER_PHONE_NUMBER, USER_EMAIL, TENANT_ID, LOCATION, ROLE) VALUES ($1, $2, $3, $4, $5, $6) RETURNING USER_ID",
	//	request.UserFullName,
	//	request.PhoneNumber,
	//	request.Email,
	//	request.CityId,
	//	request.Location,
	//	"USER")

	//var ownerID int
	//err = db.QueryRow(createOwner, request.OwnerName, request.PhoneNumber, request.Email).Scan(&ownerID)
	//if err != nil {
	//	tx.Rollback()
	//	log.WithError(err).Errorf("Failed to insert owner")
	//	return err
	//}
	//insertTenatReference := fmt.Sprintf("INSERT INTO sport.TENANT_REFERENCE (TENANT_ID, TENANT_NAME, OWNER_ID) VALUES ($1, $2, $3)")
	//_, err = db.Exec(insertTenatReference, ownerID, request.SchemaCode, ownerID)
	//if err != nil {
	//	tx.Rollback()
	//	log.WithError(err).Errorf("Failed to insert tenant reference")
	//	return err
	//}

	//log.Infof("New owner inserted with OWNER_ID: %d", ownerID)

	return err
}

func (c *classManagementStore) GetCityStore(ctx context.Context, cityId int) (string, error) {
	log.Infof("Get city store")

	sqlQuery := "SELECT public.CITY_NAME FROM CITY WHERE CITY_ID = ?"
	var cityName string
	err := c.db.GetContext(ctx, &cityName, sqlQuery, cityId)

	if err != nil {
		log.WithError(err).Errorf("Failed to retrieve city from the database")
		return "", err
	}

	return cityName, nil
}

func (c *classManagementStore) UploadDocumentStore(fileData []byte, fileTypeDoc string, fileName string) (*int64, error) {
	//var document struct {
	//	Filename string `json:"filename"`
	//	Type     string `json:"type"`
	//}
	//
	//err := json.Unmarshal(fileData, &document)
	//if err != nil {
	//	log.WithError(err).Errorf("Failed to unmarshal file data")
	//	return nil, err
	//}

	insert := "INSERT INTO document (filename, type) VALUES (?, ?)"
	res, err := c.db.Exec(insert, fileName, fileTypeDoc)
	if err != nil {
		log.WithError(err).Errorf("Failed to upload document to the database")
		return nil, err
	}
	documentID, err := res.LastInsertId()
	if err != nil {
		log.WithError(err).Errorf("Failed to retrieve last insert ID")
		return nil, err
	}
	return &documentID, nil
}

func (c *classManagementStore) UploadMetadataStore(rq qlda.AutoGenerated, documentId int64, ctx context.Context) error {
	documentQuery := "SELECT id FROM document WHERE id = ?"
	var existingDocumentID int64
	err := c.db.GetContext(ctx, &existingDocumentID, documentQuery, documentId)
	if err != nil {
		return err
	}

	sqlInsert := "INSERT INTO page (pnum, text, document_id) VALUES (?, ?, ?)"
	_, err = c.db.ExecContext(ctx, sqlInsert, rq.Page.Page, rq.Page.Text, documentId)
	if err != nil {
		log.WithError(err).Errorf("Failed to upload metadata to the database")
		return err
	}

	for _, mdTable := range rq.Page.Tables {
		sqlInsert := "INSERT INTO md_table(pnum, content, table_index, document_id) VALUES (?, ?, ?, ?)"
		_, err = c.db.ExecContext(ctx, sqlInsert, rq.Page.Page, mdTable.Content, mdTable.TableIndex, documentId)
		if err != nil {
			log.WithError(err).Errorf("Failed to upload metadata to the database")
			return err
		}
	}

	for _, image := range rq.Page.Images {
		sqlInsert = "INSERT INTO image (pnum, name, content, document_id) VALUES (?, ?, ?, ?)"
		_, err = c.db.ExecContext(ctx, sqlInsert, rq.Page.Page, image.Name, image.Content, documentId)
		if err != nil {
			log.WithError(err).Errorf("Failed to upload metadata to the database")
			return err
		}
	}
	return nil
}
