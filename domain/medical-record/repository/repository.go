package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/mrakhaf/halo-suster/domain/medical-record/interfaces"
	"github.com/mrakhaf/halo-suster/models/dto"
	"github.com/mrakhaf/halo-suster/models/entity"
	"github.com/mrakhaf/halo-suster/models/request"
	"github.com/mrakhaf/halo-suster/shared/utils"
)

type repoHandler struct {
	databaseDB *sql.DB
}

func NewRepository(databaseDB *sql.DB) interfaces.Repository {
	return &repoHandler{
		databaseDB: databaseDB,
	}
}

func (repo *repoHandler) SavePatient(req request.SavePatient) (data entity.Patient, err error) {

	data = entity.Patient{
		ID:                  utils.GenerateUUID(),
		IdentityNumber:      req.IdentityNumber,
		PhoneNumber:         req.PhoneNumber,
		Name:                req.Name,
		BirthDate:           req.BirthDate,
		Gender:              req.Gender,
		IdentityCardScanImg: req.IdentityCardScanImage,
		CreatedAt:           time.Now().Format("2006-01-02 15:04:05"),
	}

	query := fmt.Sprintf("INSERT INTO patient (id, identitynumber, name, birthdate, phonenumber, gender, identityscanimage, created_at) VALUES ('%s', %d, '%s', '%s', '%s', '%s', '%s', '%s')", data.ID, data.IdentityNumber, data.Name, data.BirthDate, data.PhoneNumber, data.Gender, data.IdentityCardScanImg, data.CreatedAt)

	_, err = repo.databaseDB.Exec(query)

	if err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint \"medical_record_identitynumber_key\"" {
			err = errors.New("Identity number already exist")
			return
		}

		err = errors.New("Save patient failed")
		return
	}

	return
}

func (repo *repoHandler) GetPatientByIdentity(identitynumber int) (data entity.Patient, err error) {
	query := fmt.Sprintf("SELECT * FROM patient WHERE identitynumber = %d", identitynumber)

	err = repo.databaseDB.QueryRow(query).Scan(&data.ID, &data.IdentityNumber, &data.PhoneNumber, &data.Name, &data.BirthDate, &data.Gender, &data.IdentityCardScanImg, &data.CreatedAt)

	if err != nil {
		fmt.Println(err.Error())
		if err == sql.ErrNoRows {
			err = errors.New("Identity not found")
			return
		}
		err = errors.New("Get data failed")
		return
	}

	return
}

func (repo *repoHandler) GetPatients(req request.GetPatientsParam) (data []entity.Patient, err error) {
	query := fmt.Sprintf("SELECT * FROM patient WHERE 1 = 1")

	if req.IdentityNumber != nil {
		query += fmt.Sprintf(" AND identitynumber = %d", *req.IdentityNumber)
	}

	if req.Name != nil {
		query += fmt.Sprintf(" AND LOWER(name) LIKE '%%%s%%'", *req.Name)
	}

	if req.PhoneNumber != nil {
		query += fmt.Sprintf(" AND phonenumber LIKE '%%%s%%'", *req.PhoneNumber)
	}

	if req.CreatedAt != nil {
		if *req.CreatedAt == "asc" {
			query += " ORDER BY created_at ASC"
		} else if *req.CreatedAt == "desc" {
			query += " ORDER BY created_at DESC"
		}
	}

	if req.Limit != nil {
		query += fmt.Sprintf(" LIMIT %d", *req.Limit)
	} else {
		query += fmt.Sprintf(" LIMIT 5")
	}

	if req.Offset != nil {
		query += fmt.Sprintf(" OFFSET %d", *req.Offset)
	} else {
		query += fmt.Sprintf(" OFFSET 0")
	}

	fmt.Println(query)
	rows, err := repo.databaseDB.Query(query)

	if err != nil {
		err = errors.New("Get patient failed")
		return
	}

	defer rows.Close()

	patients := entity.Patient{}

	for rows.Next() {
		err = rows.Scan(&patients.ID, &patients.IdentityNumber, &patients.Name, &patients.PhoneNumber, &patients.BirthDate, &patients.Gender, &patients.IdentityCardScanImg, &patients.CreatedAt)
		data = append(data, patients)
	}

	return
}

func (repo *repoHandler) SaveMedicalRecord(req request.SaveMedicalRecord, nip int) (data entity.MedicalRecord, err error) {
	data = entity.MedicalRecord{
		ID:             utils.GenerateUUID(),
		IdentityNumber: req.IdentityNumber,
		Symptoms:       req.Symptoms,
		Medications:    req.Medications,
		CreatedAt:      time.Now().Format("2006-01-02 15:04:05"),
	}

	query := fmt.Sprintf("INSERT INTO medical_record (id, identitynumber, symptoms, medications, nip, createdat) VALUES ('%s', %d, '%s', '%s', '%d', '%s')", data.ID, data.IdentityNumber, data.Symptoms, data.Medications, nip, data.CreatedAt)

	fmt.Println(query)
	_, err = repo.databaseDB.Exec(query)

	fmt.Println(err)
	if err != nil {
		err = errors.New("Save medical record failed")
		return
	}

	return
}

func (repo *repoHandler) GetMedicalRecords(req request.GetMedicalRecordsParam) (data []dto.MedicalRecordResponse, err error) {
	query := fmt.Sprintf("SELECT p.identitynumber, p.phonenumber, p.name, p.gender, p.birthdate, p.identityscanimage, mr.symptoms, " +
		"mr.medications, mr.createdat, u.nip, u.name, u.id FROM medical_record mr " +
		"JOIN patient p ON mr.identitynumber = p.identitynumber " +
		"JOIN users u ON mr.nip = u.nip WHERE 1 = 1")

	if req.IdentityNumber != nil {
		query += fmt.Sprintf(" AND p.identitynumber = %d", *req.IdentityNumber)
	}

	if req.UserId != nil {
		query += fmt.Sprintf(" AND u.id = '%s'", *req.UserId)
	}

	if req.Nip != nil {
		query += fmt.Sprintf(" AND u.nip = '%s'", *req.Nip)
	}

	if req.CreatedAt != nil {
		if *req.CreatedAt == "asc" {
			query += " ORDER BY createdat ASC"
		} else if *req.CreatedAt == "desc" {
			query += " ORDER BY createdat DESC"
		}
	}

	if req.Limit != nil {
		query += fmt.Sprintf(" LIMIT %d", *req.Limit)
	} else {
		query += fmt.Sprintf(" LIMIT 5")
	}

	if req.Offset != nil {
		query += fmt.Sprintf(" OFFSET %d", *req.Offset)
	} else {
		query += fmt.Sprintf(" OFFSET 0")
	}

	fmt.Println(query)
	rows, err := repo.databaseDB.Query(query)

	if err != nil {
		err = errors.New("Get medical record failed")
		return
	}

	defer rows.Close()

	patient := entity.Patient{}
	users := entity.Users{}
	medicalRecords := entity.MedicalRecord{}

	for rows.Next() {
		err = rows.Scan(&patient.IdentityNumber, &patient.Name, &patient.PhoneNumber, &patient.Gender, &patient.BirthDate, &patient.IdentityCardScanImg, &medicalRecords.Symptoms, &medicalRecords.Medications, &medicalRecords.CreatedAt, &users.NIP, &users.Name, &users.ID)

		data = append(data, dto.MedicalRecordResponse{
			IdentityDetail: dto.IdentityDetail{
				IdentityNumber: patient.IdentityNumber,
				PhoneNumber:    patient.PhoneNumber,
				Name:           patient.Name,
				BirthDate:      patient.BirthDate,
			},
			Symptoms:    medicalRecords.Symptoms,
			Medications: medicalRecords.Medications,
			CreatedAt:   medicalRecords.CreatedAt,
			CreatedBy: dto.CreatedBy{
				UserId: users.ID,
				Name:   users.Name,
				Nip:    users.NIP,
			},
		})
	}

	return
}
