package user

import (
	"coretest-go/service/user"
	"reflect"
	"testing"

	_ "github.com/lib/pq" // backend-db driver
	"github.com/jmoiron/sqlx" // backend-db wrapper extension
)

// setup
var (
	_masterDB string = "user=application password=mabufare dbname=testmaster sslmode=disable"
	_followerDB string = "user=application password=mabufare dbname=testfollower sslmode=disable"
	masterDB *sqlx.DB
	followerDB *sqlx.DB
	retrievedInsertedID [2]int64
)

// init
func init() {
	masterDB, _ = sqlx.Connect("postgres", _masterDB)
	followerDB, _ = sqlx.Connect("postgres", _followerDB)
}

// test New()
func TestNew(t *testing.T) {
	testConnectionInit := Resource{
		masterDB: masterDB,
		followerDB: followerDB,
	}
	type args struct {
		masterDB   *sqlx.DB
		followerDB *sqlx.DB
	}
	tests := []struct {
		name string
		args args
		want *Resource
	}{
		{
			name: "New User Resource",
			args: args{
				masterDB: testConnectionInit.masterDB,
				followerDB: testConnectionInit.followerDB,
			},
			want: &testConnectionInit,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.masterDB, tt.args.followerDB); (testConnectionInit.masterDB != tt.want.masterDB || testConnectionInit.followerDB != tt.want.followerDB) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

// test Create()
func TestResource_Create(t *testing.T) {
	type fields struct {
		masterDB   *sqlx.DB
		followerDB *sqlx.DB
	}
	type args struct {
		u user.User
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   user.User
	}{
		{
			name: "Add User 1 Active",
			fields: fields{
				masterDB: masterDB,
				followerDB: followerDB,
			},
			args: args {
				u: user.User{
					Name: "User 101",
					Email: "user101@tokopedia.com",
					IsActive: true,
					IdempotencyKey: "testIfAny",
				},
			},
			want: user.User {
				Name: "User 101",
				Email: "user101@tokopedia.com",
				IsActive: true,
				IdempotencyKey: "testIfAny",
			},
		},
		{
			name: "Add User 2 Not Active",
			fields: fields{
				masterDB: masterDB,
				followerDB: followerDB,
			},
			args: args {
				u: user.User{
					Name: "User 102",
					Email: "user102@tokopedia.com",
					IsActive: false,
					IdempotencyKey: "testIfAny",
				},
			},
			want: user.User {
				Name: "User 102",
				Email: "user102@tokopedia.com",
				IsActive: false,
				IdempotencyKey: "testIfAny",
			},
		},
	}
	for idx, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Resource{
				masterDB:   tt.fields.masterDB,
				followerDB: tt.fields.followerDB,
			}
			got := r.Create(tt.args.u)
			retrievedInsertedID[idx] = got.ID
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Resource.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

// test Get()
func TestResource_Get(t *testing.T) {
	type fields struct {
		masterDB   *sqlx.DB
		followerDB *sqlx.DB
	}
	type args struct {
		id int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    user.User
		wantErr bool
	}{
		{
			name: "Get User 101",
			fields: fields{
				masterDB: masterDB,
				followerDB: followerDB,
			},
			args: args{
				id: retrievedInsertedID[0],
			},
			want: user.User{
				ID: retrievedInsertedID[0],
				Name: "User 101",
				Email: "user101@tokopedia.com",
				IsActive: true,
				IdempotencyKey: "testIfAny",
			},
			wantErr: false,
		},
		{
			name: "Get User 102",
			fields: fields{
				masterDB: masterDB,
				followerDB: followerDB,
			},
			args: args{
				id: retrievedInsertedID[1],
			},
			want: user.User{
				ID: retrievedInsertedID[1],
				Name: "User 102",
				Email: "user102@tokopedia.com",
				IsActive: false,
				IdempotencyKey: "testIfAny",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Resource{
				masterDB:   tt.fields.masterDB,
				followerDB: tt.fields.followerDB,
			}
			got, err := r.Get(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Resource.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Resource.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

// test List()
func TestResource_GetList(t *testing.T) {
	type fields struct {
		masterDB   *sqlx.DB
		followerDB *sqlx.DB
	}
	type args struct {
		isPaginated    bool
		perPage        int32
		page           int32
		orderBy        string
		orderDirection string
		extraCondition string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []user.User
		wantErr bool
	}{
		{
			name: "Get User List that contain user 101 and 102",
			fields: fields{
				masterDB: masterDB,
				followerDB: followerDB,
			},
			args: args{
				isPaginated: false,
				perPage: -1,
				page: -1,
				orderBy: "",
				orderDirection: "",
				extraCondition: "",
			},
			want: []user.User {
				user.User{
					ID: retrievedInsertedID[0],
					Name: "User 101",
					Email: "user101@tokopedia.com",
					IsActive: true,
					IdempotencyKey: "testIfAny",
				},
				user.User{
					ID: retrievedInsertedID[1],
					Name: "User 102",
					Email: "user102@tokopedia.com",
					IsActive: false,
					IdempotencyKey: "testIfAny",
				},
			}, // subset of list that should be exist (at least 1 element)
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Resource{
				masterDB:   tt.fields.masterDB,
				followerDB: tt.fields.followerDB,
			}
			got, err := r.GetList(tt.args.isPaginated, tt.args.perPage, tt.args.page, tt.args.orderBy, tt.args.orderDirection, tt.args.extraCondition)
			if (err != nil) != tt.wantErr {
				t.Errorf("Resource.GetList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			isSubsetExist := false // at least 1
			for _, elmt := range got {
				for _, elmt2 := range tt.want {
					if reflect.DeepEqual(elmt, elmt2) { isSubsetExist = true }
				}
			}
			if !isSubsetExist {
				t.Errorf("Resource.GetList() = %v, want %v", got, tt.want)
			}
		})
	}
}

// test Update()
func TestResource_Update(t *testing.T) {
	type fields struct {
		masterDB   *sqlx.DB
		followerDB *sqlx.DB
	}
	type args struct {
		u user.User
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   user.User
	}{
		{
			name: "Update User 101 from Active to Inactive",
			fields: fields{
				masterDB: masterDB,
				followerDB: followerDB,
			},
			args: args {
				u: user.User{
					ID: retrievedInsertedID[0],
					Name: "User 101",
					Email: "user101@tokopedia.com",
					IsActive: false,
					IdempotencyKey: "testIfAny",
				},
			},
			want: user.User {
				ID: retrievedInsertedID[0],
				Name: "User 101",
				Email: "user101@tokopedia.com",
				IsActive: false,
				IdempotencyKey: "testIfAny",
			},
		},
		{
			name: "Add User 102 from Not Active to Active",
			fields: fields{
				masterDB: masterDB,
				followerDB: followerDB,
			},
			args: args {
				u: user.User{
					ID: retrievedInsertedID[1],
					Name: "User 102",
					Email: "user102@tokopedia.com",
					IsActive: true,
					IdempotencyKey: "testIfAny",
				},
			},
			want: user.User {
				ID: retrievedInsertedID[1],
				Name: "User 102",
				Email: "user102@tokopedia.com",
				IsActive: true,
				IdempotencyKey: "testIfAny",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Resource{
				masterDB:   tt.fields.masterDB,
				followerDB: tt.fields.followerDB,
			}
			if got := r.Update(tt.args.u); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Resource.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

// test Delete()
func TestResource_Delete(t *testing.T) {
	type fields struct {
		masterDB   *sqlx.DB
		followerDB *sqlx.DB
	}
	type args struct {
		id int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Delete User 101",
			fields: fields{
				masterDB: masterDB,
				followerDB: followerDB,
			},
			args: args{
				id: retrievedInsertedID[0],
			},
		},
		{
			name: "Delete User 102",
			fields: fields{
				masterDB: masterDB,
				followerDB: followerDB,
			},
			args: args{
				id: retrievedInsertedID[1],
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Resource{
				masterDB:   tt.fields.masterDB,
				followerDB: tt.fields.followerDB,
			}
			r.Delete(tt.args.id)
		})
	}
}
