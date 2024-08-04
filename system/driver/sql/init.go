package sql

import (
	"context"
	"fmt"

	"github.com/apito-io/buffers/protobuff"
	"github.com/apito-io/buffers/shared"
	_const "github.com/apito-io/databasedriver"
	"gorm.io/datatypes"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgreSQLDriver struct {
	Gorm             *gorm.DB
	DriverCredential *protobuff.DriverCredentials
}

func (p PostgreSQLDriver) GetSystemUserByUsername(ctx context.Context, username string) (*protobuff.SystemUser, error) {
	//TODO implement me
	panic("implement me")
}

func GetSystemSQLDriver(driverCredentials *protobuff.DriverCredentials) (*PostgreSQLDriver, error) {

	var gormDB *gorm.DB
	var err error

	switch driverCredentials.Engine {
	case _const.MySQLDriver:
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			driverCredentials.User, driverCredentials.Password, driverCredentials.Host, driverCredentials.Port, driverCredentials.Database)
		gormDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	case _const.PostgresSQLDriver:
		dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
			driverCredentials.Host, driverCredentials.Port, driverCredentials.User, driverCredentials.Password, driverCredentials.Database)
		gormDB, err = gorm.Open(postgres.New(postgres.Config{
			DSN: dsn,
		}), &gorm.Config{})
	}

	if err != nil {
		return nil, err
	}

	return &PostgreSQLDriver{Gorm: gormDB, DriverCredential: driverCredentials}, nil
}

func (p PostgreSQLDriver) SearchResource(ctx context.Context, param *shared.CommonSystemParams) (*shared.SearchResponse[any], error) {
	//TODO implement me
	panic("implement me")
}

func (p PostgreSQLDriver) FindOrganizationAdmin(ctx context.Context, orgId string) (*protobuff.SystemUser, error) {
	//TODO implement me
	panic("implement me")
}

func (p PostgreSQLDriver) GetOrganizations(ctx context.Context, userId string) (*shared.SearchResponse[protobuff.Organization], error) {
	//TODO implement me
	panic("implement me")
}

type UsagesTracking struct {
	ProjectID string `json:"project_id"`

	ApiCalls        uint32  `protobuf:"varint,1,opt,name=api_calls,json=apiCalls,proto3" json:"api_calls,omitempty" firestore:"api_calls,omitempty"`
	ApiBandwidth    float64 `protobuf:"fixed64,2,opt,name=api_bandwidth,json=apiBandwidth,proto3" json:"api_bandwidth,omitempty" firestore:"api_bandwidth,omitempty"`
	MediaStorage    float64 `protobuf:"fixed64,3,opt,name=media_storage,json=mediaStorage,proto3" json:"media_storage,omitempty" firestore:"media_storage,omitempty"`
	MediaBandwidth  float64 `protobuf:"fixed64,4,opt,name=media_bandwidth,json=mediaBandwidth,proto3" json:"media_bandwidth,omitempty" firestore:"media_bandwidth,omitempty"`
	NumberOfMedia   float64 `protobuf:"fixed64,5,opt,name=number_of_media,json=numberOfMedia,proto3" json:"number_of_media,omitempty" firestore:"number_of_media,omitempty"`
	NumberOfRecords float64 `protobuf:"fixed64,6,opt,name=number_of_records,json=numberOfRecords,proto3" json:"number_of_records,omitempty" firestore:"number_of_records,omitempty"`
}

type DriverCredentials struct {
	ProjectID string `json:"project_id"`

	Engine   string `protobuf:"bytes,1,opt,name=engine,proto3" json:"engine,omitempty"`
	Host     string `protobuf:"bytes,2,opt,name=host,proto3" json:"host,omitempty"`
	Port     string `protobuf:"bytes,3,opt,name=port,proto3" json:"port,omitempty"`
	User     string `protobuf:"bytes,4,opt,name=user,proto3" json:"user,omitempty"`
	Password string `protobuf:"bytes,5,opt,name=password,proto3" json:"password,omitempty"`
	Database string `protobuf:"bytes,6,opt,name=database,proto3" json:"database,omitempty"`
	// for firebase
	FirebaseProjectId             string `protobuf:"bytes,7,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`
	FirebaseProjectCredentialJson string `protobuf:"bytes,8,opt,name=project_credential_json,json=projectCredentialJson,proto3" json:"project_credential_json,omitempty"`
	// for dynamodb
	AccessKey string `protobuf:"bytes,9,opt,name=access_key,json=accessKey,proto3" json:"access_key,omitempty"`
	SecretKey string `protobuf:"bytes,10,opt,name=secret_key,json=secretKey,proto3" json:"secret_key,omitempty"`
}

type APIToken struct {
	ProjectID string `json:"project_id"`

	Name   string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty" firestore:"name,omitempty"`
	Token  string `protobuf:"bytes,2,opt,name=token,proto3" json:"token,omitempty" firestore:"token,omitempty"`
	Role   string `protobuf:"bytes,3,opt,name=role,proto3" json:"role,omitempty" firestore:"role,omitempty"`
	Expire string `protobuf:"bytes,4,opt,name=expire,proto3" json:"expire,omitempty" firestore:"expire,omitempty"`
}

type AddOnsDetails struct {
	ProjectID string `json:"project_id"`

	Locals             []string `protobuf:"bytes,1,rep,name=locals,proto3" json:"locals,omitempty" firestore:"locals,omitempty"`
	SystemGraphqlHooks bool     `protobuf:"varint,2,opt,name=system_graphql_hooks,json=systemGraphqlHooks,proto3" json:"system_graphql_hooks,omitempty" firestore:"system_graphql_hooks,omitempty"`
	RevisionHistory    bool     `protobuf:"varint,3,opt,name=revision_history,json=revisionHistory,proto3" json:"revision_history,omitempty" firestore:"revision_history,omitempty"`
	EnableAuth         bool     `protobuf:"bytes,4,opt,name=enable_auth,proto3" json:"enable_auth,omitempty" firestore:"enable_auth,omitempty"`
}

// Project user project
type Project struct {
	XKey               string             `protobuf:"bytes,1,opt,name=_key,json=Key,proto3" json:"_key,omitempty" firestore:"_key,omitempty"`
	Id                 string             `gorm:"primaryKey;autoIncrement:false" protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty" firestore:"id,omitempty"`
	ProjectName        string             `protobuf:"bytes,3,opt,name=project_name,json=projectName,proto3" json:"project_name,omitempty" firestore:"project_name,omitempty"`
	ProjectDescription string             `protobuf:"bytes,4,opt,name=project_description,json=projectDescription,proto3" json:"project_description,omitempty" firestore:"project_description,omitempty"`
	Schema             datatypes.JSONMap  `protobuf:"bytes,5,opt,name=schema,proto3" json:"schema,omitempty" firestore:"schema,omitempty"`
	CreatedAt          string             `protobuf:"bytes,6,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty" firestore:"created_at,omitempty"`
	UpdatedAt          string             `protobuf:"bytes,7,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty" firestore:"updated_at,omitempty"`
	ExpireAt           string             `protobuf:"bytes,8,opt,name=expire_at,json=expireAt,proto3" json:"expire_at,omitempty" firestore:"expire_at,omitempty"`
	Plugins            datatypes.JSONMap  `protobuf:"bytes,9,rep,name=plugins,proto3" json:"plugins,omitempty" firestore:"plugins,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Addons             *AddOnsDetails     `protobuf:"bytes,10,opt,name=addons,proto3" json:"addons,omitempty" firestore:"addons,omitempty"`
	Tokens             []*APIToken        `gorm:"foreignKey:ProjectID" protobuf:"bytes,11,rep,name=tokens,proto3" json:"tokens,omitempty" firestore:"tokens,omitempty"`
	Roles              datatypes.JSONMap  `protobuf:"bytes,12,rep,name=roles,proto3" json:"roles,omitempty" firestore:"roles,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Driver             *DriverCredentials `gorm:"foreignKey:ProjectID" protobuf:"bytes,13,opt,name=driver,proto3" json:"driver,omitempty" firestore:"driver,omitempty"`
	TempBanned         bool               `protobuf:"varint,14,opt,name=temp_banned,json=tempBanned,proto3" json:"temp_banned,omitempty" firestore:"temp_banned,omitempty"`
	Plan               string             `protobuf:"bytes,15,opt,name=plan,proto3" json:"plan,omitempty" firestore:"plan,omitempty"`
	TrialEnds          string             `protobuf:"bytes,16,opt,name=trial_ends,json=trialEnds,proto3" json:"trial_ends,omitempty" firestore:"trial_ends,omitempty"`
	FromExample        string             `protobuf:"bytes,17,opt,name=from_example,json=fromExample,proto3" json:"from_example,omitempty" firestore:"from_example,omitempty"`
	Limits             *UsagesTracking    `gorm:"foreignKey:ProjectID" protobuf:"bytes,18,opt,name=limits,proto3" json:"limits,omitempty" firestore:"limits,omitempty"`
}

func (p PostgreSQLDriver) RunMigration() error {

	err := p.Gorm.AutoMigrate(Project{}, UsagesTracking{}, DriverCredentials{}, APIToken{}, AddOnsDetails{})
	if err != nil {
		fmt.Println(err.Error())
	}

	return nil
}
