package models

import (
	"fmt"
	"time"

	"github.com/Katsusan/centaur/internal/config"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

type User struct {
	gorm.Model
	Username          string    `gorm:"column:username"`
	RealName          string    `gorm:"column:realname"`
	Orgid             string    `gorm:"column:orgid"`
	Password          string    `gorm:"column:password"`
	Status            int       `gorm:"column:status"`
	Roles             UserRoles `gorm:"column:roles"`
	Expireddate       time.Time `gorm:"column:expireddate"`
	Logintime         time.Time `gorm:"column:logintime"`
	Loginip           string    `gorm:"column:loginip"`
	Lasttime          time.Time `gorm:"column:lasttime"`
	Lastip            string    `gorm:"column:lastip"`
	Skin              string    `gorm:"column:skin"`
	Langcode          string    `gorm:"column:langcode"`
	Sex               string    `gorm:"column:sex"`
	Birthday          time.Time `gorm:"column:birthday"`
	Idcard            string    `gorm:"column:idcard"`
	School            string    `gorm:"column:school"`
	Graduation        string    `gorm:"column:graduation"`
	Degree            string    `gorm:"column:degree"`
	Major             string    `gorm:"column:major"`
	Country           string    `gorm:"column:country"`
	Province          string    `gorm:"column:province"`
	City              string    `gorm:"column:city"`
	Address           string    `gorm:"column:address"`
	Postcode          string    `gorm:"column:postcode"`
	Phone             string    `gorm:"column:phone"`
	Fax               string    `gorm:"column:fax"`
	Mobile            string    `gorm:"column:mobile"`
	Email             string    `gorm:"column:email"`
	Remark            string    `gorm:"column:remark"`
	Creator           string    `gorm:"column:creator"`
	Modifier          string    `gorm:"column:modifier"`
	Usertype          string    `gorm:"column:usertype"`
	Postid            string    `gorm:"column:postid"`
	Isleader          bool      `gorm:"column:isleader;null;default(false)"`
	Expired           string    `gorm:"column:expired;null;default(0)"`
	Ipconfig          string    `gorm:"column:ipconfig"`
	EnglishName       string    `gorm:"column:english_name"`
	Nationality       string    `gorm:"column:nationality"`
	Employeeid        string    `gorm:"column:employeeid"`
	Entrydate         time.Time `gorm:"column:entrydate"`
	ResidenceAddr     string    `gorm:"column:residence_addres)"`
	ResidenceType     string    `gorm:"column:residence_type"`
	MaritalStatus     string    `gorm:"column:marital_status"`
	NativePlace       string    `gorm:"column:native_place"`
	WorkDate          time.Time `gorm:"column:work_date"`
	ContactWay        string    `gorm:"column:contact_way"`
	ContactPerson     string    `gorm:"column:contact_person"`
	ProfessionalTitle string    `gorm:"column:professional_title"`
	ComputerLevel     string    `gorm:"column:computer_level"`
	ComputerCert      string    `gorm:"column:computer_cert"`
	EnglishLevel      string    `gorm:"column:english_level"`
	EnglishCert       string    `gorm:"column:english_cert"`
	JapaneseLevel     string    `gorm:"column:japanese_level"`
	JapaneseCert      string    `gorm:"column:japanese_cert"`
	Speciality        string    `gorm:"column:speciality"`
	SpecialityCert    string    `gorm:"column:speciality_cert"`
	HobbySport        string    `gorm:"column:hobby_sport"`
	HobbyArt          string    `gorm:"column:hobby_art"`
	HobbyOther        string    `gorm:"column:hobby_other"`
	KeyUser           string    `gorm:"column:key_user"`
	WorkCard          string    `gorm:"column:work_card"`
	GuardCard         string    `gorm:"column:guard_card"`
	Computer          string    `gorm:"column:computer"`
	Ext               string    `gorm:"column:ext"`
	Msn               string    `gorm:"column:msn"`
	Rank              string    `gorm:"column:rank"`
}

type Profile struct {
	UserID      string
	UserName    string
	CompanyCode string
	LoginIP     string
	LastLogin   time.Time
}

type UserRole struct {
	gorm.Model
	UserID     string `gorm:"column:user_id"`
	RoleID     string `gorm:"column:role_id"`
	ExpireDate time.Time
}

func (role UserRole) TableName() string {
	return "user_role"
}

type UserRoles []*UserRole

type Users []*User

//用户查询条件参数
type UserQueryParam struct {
	UserName     string
	RealName     string
	UserNameLike string
	RealNameLike string
	Status       int
	RoleIDs      []string
}

type UserQueryOptions struct {
	PageParam    *PaginationParam //分页参数
	IncludeRoles bool             //包含角色
}

type UserQueryResult struct {
	Res     Users
	PageRes *PaginationResult
}

type UserPageShow struct {
	UserName    string
	RealName    string
	PhoneNumber string
	Email       string
	Status      int
	CreatedAt   time.Time
	Roles       []*Role
}

func (User) TableName() string {
	return "user_tb"
}

//UserQuery will return UserQueryResult by given UserQueryParam(UserName,RealName,etc...) and UserQueryOptions
func UserQuery(params UserQueryParam, opts ...UserQueryOptions) (*UserQueryResult, error) {
	db := config.GetGlobalConfig().DB()

	if v := params.UserName; v != "" {
		db = db.Where("username=?", v)
	}

	if v := params.UserNameLike; v != "" {
		db = db.Where("username LIKE ?", "%"+v+"%")
	}

	if v := params.RealNameLike; v != "" {
		db = db.Where("realname LIKE ?", "%"+v+"%")
	}

	if v := params.Status; v > 0 {
		db = db.Where("status=?", v)
	}

	if v := params.RoleIDs; len(v) > 0 {
		subQuery := db.Model(UserRole{}).Select("user_id").Where("role_in IN(?)", v).SubQuery()
		db = db.Where("record_id IN(?)", subQuery)
	}

	db = db.Order("id DESC")

	var opt UserQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	var list Users
	pageRes, err := WrapQuery(db, opt.PageParam, &list)
	if err != nil {
		log.Println("Wrap User Query failed(error=%v)", err)
		return nil, fmt.Errorf("WrapUserQueryError")
	}

	queryRes := &UserQueryResult{
		PageRes: pageRes,
		Res:     list,
	}

	return queryRes, nil
}

func WrapQuery(db *gorm.DB, param *PaginationParam, out interface{}) (*PaginationResult, error) {
	if param != nil && param.PageIndex >= 0 && param.PageSize > 0 {
		var count int
		res := db.Count(&count)
		if err := res.Error; err != nil {
			return nil, err
		} else if count == 0 {
			return &PaginationResult{
				Total: 0,
			}, nil
		}
		db = db.Offset((param.PageIndex - 1) * param.PageSize).Limit(param.PageSize)
	}

	res := db.Find(out)
	return nil, res.Error
}

//ToRoleIDs will traverse UserRoles and return aggregation of every userrole's roleid.
func (uRoles UserRoles) ToRoleIDs() []string {
	roleIDs := make([]string, len(uRoles))
	for i, urole := range uRoles {
		roleIDs[i] = urole.RoleID
	}
	return roleIDs
}

func (usrs Users) ToRoleIDs() []string {
	var roleIDs []string
	for _, u := range usrs {
		roleIDs = append(roleIDs, u.Roles.ToRoleIDs()...)
	}
	return roleIDs
}

func (usrs Users) ToPageShows() []*PaginationResult
