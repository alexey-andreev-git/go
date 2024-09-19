package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"what-to.com/internal/config"
	"what-to.com/internal/models/entity"
	"what-to.com/internal/resources"
	"what-to.com/internal/safemap"

	_ "github.com/lib/pq"
)

const (
	ENTITY_CREATE = `INSERT INTO entities (entity_name, entity_comment) VALUES ($1, $2) RETURNING entity_id`
	ENTITY_GET    = `SELECT entity_id, entity_name, entity_comment FROM entities`
	ENTITY_UPDATE = `UPDATE entities SET %s WHERE entity_id=$%d`
	ENTITY_DELETE = `DELETE FROM entities WHERE entity_id=$1`

	ENTITY_DATA_CREATE = `INSERT INTO entities_data (entities_data_entity, entities_data_order, entities_data_value) VALUES ($1, $2, $3)`
	ENTITY_DATA_GET    = `SELECT entities_data_entity, entities_data_order, entities_data_value FROM entities_data`
	ENTITY_DATA_UPDATE = `UPDATE entities_data SET %s WHERE entities_data_entity=$%d AND entities_data_order=$%d`
	ENTITY_DATA_DELETE = `DELETE FROM entities_data WHERE entities_data_entity=$1 AND entities_data_order=$2`

	ENTITY_DATA_REF_CREATE = `INSERT INTO entities_data_reference (entities_data_reference_entity, entities_data_reference_order, entities_data_reference_name, entities_data_reference_type, entities_data_reference_comment) VALUES ($1, $2, $3, $4, $5)`
	ENTITY_DATA_REF_GET    = `SELECT entities_data_reference_entity, entities_data_reference_order, entities_data_reference_name, entities_data_reference_type, entities_data_reference_comment FROM entities_data_reference`
	ENTITY_DATA_REF_UPDATE = `UPDATE entities_data_reference SET %s WHERE entities_data_reference_entity=$%d AND entities_data_reference_order=$%d`
	ENTITY_DATA_REF_DELETE = `DELETE FROM entities_data_reference WHERE entities_data_reference_entity=$1 AND entities_data_reference_order=$2`

	// ------------------------------------------------ fmt constants SELECT
	GET_QUERY_SELECT    = "SELECT entity_data_val_%s_value "
	GET_QUERY_FROM      = "FROM entities_data "
	GET_QUERY_LEFT_JOIN = "LEFT JOIN entities_data_val_%s ON entity_data_value = entity_data_val_%s_id "
	GET_QUERY_WHERE     = "WHERE entity_data_entity = entity_id AND entity_data_order=%d"
	// ------------------------------------------------ fmt constants SELECT

	// ------------------------------------------------ fmt constants INSERT
	CREATE_QUERY_INSERT    = "INSERT INTO entities_data_val_%s (entity_data_val_%s_value) VALUES ($%d) ON CONFLICT (entity_data_val_%s_value) DO NOTHING;"
	CREATE_QUERY_FROM      = "FROM entities_data "
	CREATE_QUERY_LEFT_JOIN = "LEFT JOIN entities_data_val_%s ON entity_data_value = entity_data_val_%s_id "
	CREATE_QUERY_WHERE     = "WHERE entity_data_entity = entity_id AND entity_data_order=%d) AS field_%d"
	// ------------------------------------------------ fmt constants INSERT
)

type (
	// DBConfig holds the database connection configuration
	DBConfig struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		DBName   string `json:"dbname"`
	}

	EntityQuerys struct {
		CreateFields string
		GetFields    string
		UpdateFields string
		DeleteFields string
	}

	// PgRepository is the PostgreSQL repository
	PgRepository struct {
		DB                    *sql.DB
		appConfig             *config.Config
		dbConfig              DBConfig
		mapRefEntities        *safemap.SafeMap
		mapRefEntitiesIds     *safemap.SafeMap
		mapRefEntitiesData    *safemap.SafeMap
		mapRefEntitiesDataIds *safemap.SafeMap
		mapEntitiesQrys       map[uint]EntityQuerys
	}
)

var entityFieldToColumn = map[string]string{
	"id":      "entity_id",
	"name":    "entity_name",
	"comment": "entity_comment",
}

var entityDataFieldToColumn = map[string]string{
	"entity": "entities_data_entity",
	"order":  "entities_data_order",
	"value":  "entities_data_value",
}

var entityDataRefFieldToColumn = map[string]string{
	"entity":  "entities_data_reference_entity",
	"order":   "entities_data_reference_order",
	"name":    "entities_data_reference_name",
	"type":    "entities_data_reference_type",
	"comment": "entities_data_reference_comment",
}

// NewPgRepository initializes a new PostgreSQL repository
func NewPgRepository(appConfig *config.Config) *PgRepository {
	r := &PgRepository{
		DB:        nil,
		appConfig: appConfig,
	}
	r.mapRefEntities = safemap.NewSafeMap()
	r.mapRefEntitiesIds = safemap.NewSafeMap()
	r.mapRefEntitiesData = safemap.NewSafeMap()
	r.mapRefEntitiesDataIds = safemap.NewSafeMap()
	r.mapEntitiesQrys = make(map[uint]EntityQuerys)
	r.SetRepoConfig(appConfig.GetConfig()["database"].(config.ConfigT))
	r.connectToDb()
	return r
}

// connectToDb connects to the PostgreSQL database
func (r *PgRepository) connectToDb() {
	// psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
	// 	r.dbConfig.Host, r.dbConfig.Port, r.dbConfig.User, r.dbConfig.Password, r.dbConfig.DBName)

	r.checkDB()

	db, err := sql.Open("postgres", r.GetRepoConfigStr())
	if err != nil {
		r.appConfig.GetLogger().Fatal("Connection to the database failed:", err)
	}

	err = db.Ping()
	if err != nil {
		r.appConfig.GetLogger().Fatal("Failed to execute ping on the database:", err)
	}

	r.appConfig.GetLogger().Info("PostgreSQL DB successfully connected!")
	r.DB = db
	r.UpdateDB()
}

// SetRepoConfig sets the DBConfig struct
func (r *PgRepository) SetRepoConfig(dbConfigP config.ConfigT) {
	r.dbConfig = DBConfig{
		Host:     dbConfigP["host"].(string),
		Port:     int(dbConfigP["port"].(int)),
		User:     dbConfigP["user"].(string),
		Password: dbConfigP["password"].(string),
		DBName:   dbConfigP["dbname"].(string),
	}
}

// GetRepoConfig returns the DBConfig
func (r *PgRepository) GetRepoConfig() DBConfig {
	return r.dbConfig
}

// GetRepoConfigStr returns the DBConfig as a string
func (r *PgRepository) GetRepoConfigStr() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		r.dbConfig.Host, r.dbConfig.Port, r.dbConfig.User, r.dbConfig.Password, r.dbConfig.DBName)
}

// checkDB checks if the database exists, and creates it if it doesn't
func (r *PgRepository) checkDB() bool {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable",
		r.dbConfig.Host, r.dbConfig.Port, r.dbConfig.User, r.dbConfig.Password)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		r.appConfig.GetLogger().Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	var exists int
	db.QueryRow("SELECT 1 FROM pg_database WHERE datname=$1", r.dbConfig.DBName).Scan(&exists)

	if exists == 0 {
		r.appConfig.GetLogger().Warn("Database does not exist. Creating...")
		_, err := db.Exec(fmt.Sprintf("CREATE DATABASE \"%s\"", r.dbConfig.DBName))
		if err != nil {
			r.appConfig.GetLogger().Fatal("Failed to create database:", err)
		}
		r.appConfig.GetLogger().Info("Database created successfully.")
	}

	return true
}

// UpdateDB updates the database schema
func (r *PgRepository) UpdateDB() {
	appRes := resources.NewAppSources()
	fn := r.appConfig.GetConfig()[config.KeyInitDbFileName].(string)
	data, err := appRes.GetRes().ReadFile(fn)
	if err != nil {
		r.appConfig.GetLogger().Fatal("File read error [%s] "+fn, err)
	}
	_, err = r.DB.Exec(string(data))
	if err != nil {
		r.appConfig.GetLogger().Fatal("Failed to update database:", err)
	}

	r.LoadEntitiesReference(context.Background())
	r.LoadEntitiesDataReference(context.Background())
	r.genSqlByEntityReference()
	thresult, err := r.FindAll(context.Background(), 1)
	if err != nil {
		r.appConfig.GetLogger().Fatal("Failed to get all entities:", err)
	}
	fmt.Println(thresult)
	thresult, err = r.FindByID(context.Background(), 1, 5)
	if err != nil {
		r.appConfig.GetLogger().Fatal("Failed to get entity by id:", err)
	}
	json, err := json.Marshal(thresult)
	if err != nil {
		r.appConfig.GetLogger().Fatal("Failed to marshal result:", err)
	}
	fmt.Println(string(json))
}

// CreateEntity creates a new entity in the database
func (r *PgRepository) CreateEntity(ent map[string]interface{}) (entity.Entity, error) {
	// createdEntity := entity.Entity{
	// 	Id:      0,
	// 	Name:    ent["Name"].(string),
	// 	Comment: ent["Comment"].(string),
	// }
	createdEntity := entity.Entity{
		Id:        0,
		Reference: ent["Reference"].(uint),
	}
	err := r.DB.QueryRow(ENTITY_CREATE, ent["Name"], ent["Comment"]).Scan(&createdEntity.Id)
	return createdEntity, err
}

// GetEntity retrieves an entity from the database based on the given filter
func (r *PgRepository) GetEntity(filter map[string]interface{}) (entity.Entity, error) {
	var conditions []string
	var args []interface{}
	i := 1
	for key, value := range filter {
		column, ok := entityFieldToColumn[strings.ToLower(key)]
		if !ok {
			return entity.Entity{}, fmt.Errorf("invalid field: %s", key)
		}
		conditions = append(conditions, fmt.Sprintf("%s = $%d", column, i))
		args = append(args, value)
		i++
	}
	query := ENTITY_GET
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return entity.Entity{}, err
	}
	defer rows.Close()

	var entities []entity.Entity
	for rows.Next() {
		var e entity.Entity
		// if err := rows.Scan(&e.Id, &e.Name, &e.Comment); err != nil {
		if err := rows.Scan(&e.Id, &e.Reference); err != nil {
			return entity.Entity{}, err
		}
		entities = append(entities, e)
	}
	if len(entities) > 1 {
		return entity.Entity{}, fmt.Errorf("more than one entity found")
	}
	if len(entities) == 0 {
		return entity.Entity{}, fmt.Errorf("entity not found")
	}
	if err = rows.Err(); err != nil {
		return entity.Entity{}, err
	}
	return entities[0], nil
}

// UpdateEntity updates an entity in the database
func (r *PgRepository) UpdateEntity(ent map[string]interface{}) (sql.Result, error) {
	idStr, ok := ent["Id"]
	if !ok {
		return nil, fmt.Errorf("missing entity id")
	}
	var setClauses []string
	var args []interface{}
	i := 1
	for key, value := range ent {
		if key == "Id" {
			continue
		}
		column, ok := entityFieldToColumn[strings.ToLower(key)]
		if !ok {
			return nil, fmt.Errorf("invalid field: %s", key)
		}
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", column, i))
		args = append(args, value)
		i++
	}

	if len(setClauses) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	query := fmt.Sprintf(ENTITY_UPDATE, strings.Join(setClauses, ", "), i)
	args = append(args, idStr)

	result, err := r.DB.Exec(query, args...)
	return result, err
}

// DeleteEntity deletes an entity from the database
func (r *PgRepository) DeleteEntity(ent map[string]interface{}) (sql.Result, error) {
	idStr, ok := ent["Id"]
	if !ok {
		return nil, fmt.Errorf("missing entity id")
	}
	result, err := r.DB.Exec(ENTITY_DELETE, idStr)
	return result, err
}

// CreateEntityData creates a new entity_data in the database
func (r *PgRepository) CreateEntityData(data map[string]interface{}) ([]entity.EntityData, error) {
	createdData := []entity.EntityData{
		{
			Entity: data["Entity"].(uint),
			Order:  data["Order"].(uint),
			Value:  data["Value"].(uint),
		},
	}
	response, err := r.DB.Exec(ENTITY_DATA_CREATE, data["Entity"], data["Order"], data["Value"])
	ent, rerr := response.LastInsertId()
	if rerr != nil {
		return createdData, rerr
	}
	createdData[0].Entity = uint(ent)
	return createdData, err
}

// GetEntityData retrieves an entity_data from the database based on the given filter
func (r *PgRepository) GetEntityData(filter map[string]interface{}) ([]entity.EntityData, error) {
	var conditions []string
	var args []interface{}
	i := 1
	for key, value := range filter {
		column, ok := entityDataFieldToColumn[strings.ToLower(key)]
		if !ok {
			return nil, fmt.Errorf("invalid field: %s", key)
		}
		conditions = append(conditions, fmt.Sprintf("%s = $%d", column, i))
		args = append(args, value)
		i++
	}
	query := ENTITY_DATA_GET
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	rows, err := r.DB.Query(query, args...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entityData []entity.EntityData
	for rows.Next() {
		var ed entity.EntityData
		if err := rows.Scan(&ed.Entity, &ed.Order, &ed.Value); err != nil {
			return nil, err
		}
		entityData = append(entityData, ed)
	}
	if len(entityData) == 0 {
		return nil, fmt.Errorf("entity_data not found")
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return entityData, nil
}

// UpdateEntityData updates an entity_data in the database
func (r *PgRepository) UpdateEntityData(data map[string]interface{}) (sql.Result, error) {
	entityStr, entityOk := data["Entity"]
	orderStr, orderOk := data["Order"]
	if !entityOk || !orderOk {
		return nil, fmt.Errorf("missing entity_data entity or order")
	}
	var setClauses []string
	var args []interface{}
	i := 1
	for key, value := range data {
		if key == "Entity" || key == "Order" {
			continue
		}
		column, ok := entityDataFieldToColumn[strings.ToLower(key)]
		if !ok {
			return nil, fmt.Errorf("invalid field: %s", key)
		}
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", column, i))
		args = append(args, value)
		i++
	}

	if len(setClauses) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	query := fmt.Sprintf(ENTITY_DATA_UPDATE, strings.Join(setClauses, ", "), i, i+1)
	args = append(args, entityStr, orderStr)

	result, err := r.DB.Exec(query, args...)
	return result, err
}

// DeleteEntityData deletes an entity_data from the database
func (r *PgRepository) DeleteEntityData(data map[string]interface{}) (sql.Result, error) {
	entityStr, entityOk := data["Entity"]
	orderStr, orderOk := data["Order"]
	if !entityOk || !orderOk {
		return nil, fmt.Errorf("missing entity_data entity or order")
	}
	result, err := r.DB.Exec(ENTITY_DATA_DELETE, entityStr, orderStr)
	return result, err
}

// CreateEntityDataRef creates a new entity_data_reference in the database
func (r *PgRepository) CreateEntityDataRef(ref map[string]interface{}) ([]entity.EntityDataReference, error) {
	createdRef := []entity.EntityDataReference{
		{
			Reference: ref["Reference"].(uint),
			Order:     ref["Order"].(uint),
			Name:      ref["Name"].(string),
			Type:      ref["Type"].(string),
			Comment:   ref["Comment"].(string),
		},
	}
	_, err := r.DB.Exec(ENTITY_DATA_REF_CREATE, ref["Entity"], ref["Order"], ref["Name"], ref["Type"], ref["Comment"])
	return createdRef, err
}

// GetEntityDataRef retrieves an entity_data_reference from the database based on the given filter
func (r *PgRepository) GetEntityDataRef(filter map[string]interface{}) ([]entity.EntityDataReference, error) {
	var conditions []string
	var args []interface{}
	i := 1
	for key, value := range filter {
		column, ok := entityDataRefFieldToColumn[strings.ToLower(key)]
		if !ok {
			return nil, fmt.Errorf("invalid field: %s", key)
		}
		conditions = append(conditions, fmt.Sprintf("%s = $%d", column, i))
		args = append(args, value)
		i++
	}
	query := ENTITY_DATA_REF_GET
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entityDataRefs []entity.EntityDataReference
	for rows.Next() {
		var edr entity.EntityDataReference
		if err := rows.Scan(&edr.Reference, &edr.Order, &edr.Name, &edr.Type, &edr.Comment); err != nil {
			return nil, err
		}
		entityDataRefs = append(entityDataRefs, edr)
	}
	if len(entityDataRefs) == 0 {
		return nil, fmt.Errorf("entity_data_reference not found")
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return entityDataRefs, nil
}

// UpdateEntityDataRef updates an entity_data_reference in the database
func (r *PgRepository) UpdateEntityDataRef(ref map[string]interface{}) (sql.Result, error) {
	entityStr, entityOk := ref["Entity"]
	orderStr, orderOk := ref["Order"]
	if !entityOk || !orderOk {
		return nil, fmt.Errorf("missing entity_data_reference entity or order")
	}
	var setClauses []string
	var args []interface{}
	i := 1
	for key, value := range ref {
		if key == "Entity" || key == "Order" {
			continue
		}
		column, ok := entityDataRefFieldToColumn[strings.ToLower(key)]
		if !ok {
			return nil, fmt.Errorf("invalid field: %s", key)
		}
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", column, i))
		args = append(args, value)
		i++
	}

	if len(setClauses) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	query := fmt.Sprintf(ENTITY_DATA_REF_UPDATE, strings.Join(setClauses, ", "), i, i+1)
	args = append(args, entityStr, orderStr)

	result, err := r.DB.Exec(query, args...)
	return result, err
}

// DeleteEntityDataRef deletes an entity_data_reference from the database
func (r *PgRepository) DeleteEntityDataRef(ref map[string]interface{}) (sql.Result, error) {
	entityStr, entityOk := ref["Entity"]
	orderStr, orderOk := ref["Order"]
	if !entityOk || !orderOk {
		return nil, fmt.Errorf("missing entity_data_reference entity or order")
	}
	result, err := r.DB.Exec(ENTITY_DATA_REF_DELETE, entityStr, orderStr)
	return result, err
}

func (r *PgRepository) Create(ctx context.Context, value interface{}) error {
	// value.(entity.Entity)
	// r.mapEntitiesQrys[]
	return nil
}

func (r *PgRepository) Update(ctx context.Context, value interface{}) error {
	return nil
}

func (r *PgRepository) Delete(ctx context.Context, value interface{}) error {
	return nil
}

func fillRowsData(rows *sql.Rows) (interface{}, error) {
	var out []interface{}

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	// Create a slice of interfaces to hold the values for each column
	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))

	for i := range columns {
		valuePtrs[i] = &values[i]
	}
	for rows.Next() {
		// Scan the result into the value pointers
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		// Store the result in the map
		ent := make(map[string]interface{})
		for i, col := range columns {
			ent[col] = values[i]
		}
		out = append(out, ent)
	}
	return out, nil
}

func (r *PgRepository) FindByID(ctx context.Context, entId uint, id uint) (interface{}, error) {
	// Get SQL query string from the map based on the ID of the entity's reference
	sql := r.mapEntitiesQrys[entId].GetFields + fmt.Sprintf(" WHERE entity_id = %d", id)

	// Execute the query
	rows, err := r.DB.QueryContext(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // Ensure rows are closed after use

	// Iterate over rows
	out, err := fillRowsData(rows)
	if err != nil {
		return nil, err
	}
	// Check for any error that might have occurred during iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func (r *PgRepository) FindAll(ctx context.Context, entId uint) (interface{}, error) {
	// Get SQL query string from the map based on the ID of the entity's reference
	sql := r.mapEntitiesQrys[entId].GetFields

	// Execute the query
	rows, err := r.DB.QueryContext(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // Ensure rows are closed after use

	// Iterate over rows
	out, err := fillRowsData(rows)
	if err != nil {
		return nil, err
	}
	// Check for any error that might have occurred during iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func (r *PgRepository) LoadEntitiesReference(ctx context.Context) error {
	cols := "entity_reference_id,entity_reference_name,entity_reference_comment"
	qry := fmt.Sprintf("SELECT %s FROM entities_reference", cols)
	rows, err := r.DB.QueryContext(ctx, qry)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var ref entity.EntityReference
		var refById entity.EntityReferenceById
		if err := rows.Scan(&ref.Id, &ref.Name, &ref.Comment); err != nil {
			return err
		}
		refById.Name = ref.Name
		refById.Comment = ref.Comment
		r.mapRefEntities.Set(ref.Id, refById)
		r.mapRefEntitiesIds.Set(ref.Name, ref.Id)
	}
	if err = rows.Err(); err != nil {
		return err
	}
	return nil
}

func (r *PgRepository) LoadEntitiesDataReference(ctx context.Context) error {
	cols := "entity_data_reference_entity_reference,entity_data_reference_order,entity_data_reference_name,entity_data_reference_type,entity_data_reference_comment"
	qry := fmt.Sprintf("SELECT %s FROM entities_data_reference", cols)
	rows, err := r.DB.QueryContext(ctx, qry)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var ref entity.EntityDataReference
		if err := rows.Scan(&ref.Reference, &ref.Order, &ref.Name, &ref.Type, &ref.Comment); err != nil {
			return err
		}
		refOrder := entity.EntityDataReferenceByOrder{
			Name:    ref.Name,
			Type:    ref.Type,
			Comment: ref.Comment,
		}
		rec, ok := r.mapRefEntitiesData.Get(ref.Reference)
		if !ok {
			rec = safemap.NewSafeMap()
			r.mapRefEntitiesData.Set(ref.Reference, rec)
		}
		rec.(*safemap.SafeMap).Set(ref.Order, refOrder)
	}
	if err = rows.Err(); err != nil {
		return err
	}
	return nil
}

func (r *PgRepository) genSqlByEntityReference() error {
	fmt.Println("Entities Reference IDs:")
	r.mapRefEntitiesIds.Range(func(key, value interface{}) bool {
		fmt.Println(key, value)
		fmt.Print("    ")
		fmt.Println(r.mapRefEntities.Get(value.(uint)))
		r.genSqlByEntityDataReference(value.(uint))
		return true
	})
	return nil
}

func (r *PgRepository) genSqlByEntityDataReference(entityId uint) error {
	fmt.Println("Entities Data Reference IDs:")
	prepQry := EntityQuerys{}

	value, ok := r.mapRefEntitiesData.Get(entityId)
	if ok {
		// fmt.Print("    ")
		// fmt.Println(entityId, ": --------v")
		getFields := []string{}
		createFields := []string{}
		value.(*safemap.SafeMap).Range(func(key, val interface{}) bool {
			t := val.(entity.EntityDataReferenceByOrder).Type
			getFields = append(getFields, r.prepFieldForGetQuery(key.(uint), t))
			createFields = append(createFields, r.prepFieldForCreateQueryInsert(key.(uint), t))
			// fmt.Print("    ", "    ")
			// fmt.Println(key, ": ", val)
			return true
		})
		prepQry.GetFields = r.prepSqlForGetQuery(strings.Join(getFields, ", "), entityId)
		prepQry.CreateFields = strings.Join(createFields, "\n")
		r.mapEntitiesQrys[entityId] = prepQry
		// fmt.Print("    ", "    ", "    ")
		// fmt.Println(r.mapEntitiesQrys[entityId])
	} else {
		fmt.Println(entityId, "not found")
	}
	return nil
}

func (r *PgRepository) prepFieldForGetQuery(key uint, t string) string {
	p1 := fmt.Sprintf(GET_QUERY_SELECT, t)
	p2 := fmt.Sprintf(GET_QUERY_LEFT_JOIN, t, t)
	if t == "int" {
		p1 = "SELECT entity_data_value "
		p2 = ""
	}
	field := "(" + p1 + GET_QUERY_FROM + p2 + GET_QUERY_WHERE + ") AS field_%d"
	return fmt.Sprintf(field, key, key)
}

func (r *PgRepository) prepSqlForGetQuery(fields string, id uint) string {
	return fmt.Sprintf("SELECT * FROM (SELECT *, %s FROM entities e WHERE entity_reference = %d) as ed", fields, id)
}

func (r *PgRepository) prepFieldForCreateQueryInsert(key uint, t string) string {
	if t != "int" {
		return fmt.Sprintf(CREATE_QUERY_INSERT, t, t, key, t)
	}
	return ""
}
