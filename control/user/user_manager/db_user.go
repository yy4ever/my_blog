package user_manager

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
	"my_blog/common/conf"
	"my_blog/common/define/permission"
	"my_blog/common/entity"
	"my_blog/common/errors"
	"my_blog/common/models"
	errs "my_blog/control/user/errors"
	"strings"
)

const (
	ROLE_USER     = "User"
	ROLE_ADMIN    = "Admin"
	ROLE_SYSADMIN = "SysAdmin"
)

const (
	USER_STATE_ACTIVE = "ACTIVE"
	USER_STATE_INACTIVE = "INACTIVE"
	USER_STATE_LOCKED = "LOCKED"
)

var roles = map[string]int{
	ROLE_USER:  permission.FOLLOW | permission.COMMENT | permission.POST,
	ROLE_ADMIN: permission.FOLLOW | permission.COMMENT | permission.POST | permission.MODERATE,
	ROLE_SYSADMIN: permission.FOLLOW | permission.COMMENT | permission.POST | permission.MODERATE |
		permission.SYSADMIN | permission.REMOVE,
}

func InitData() {
	InsertRoles()
	InsertDefaultUsers()
}

func InsertRoles() {
	sql := "INSERT IGNORE role (name, permissions) VALUES %s"
	var roleValues []interface{}
	var placeholder []string
	for r, p := range roles {
		placeholder = append(placeholder, "(?, ?)")
		roleValues = append(roleValues, r)
		roleValues = append(roleValues, p)
	}
	sql = fmt.Sprintf(sql, strings.Join(placeholder, ","))
	ret, err := conf.DB.Exec(sql, roleValues...)
	if err != nil {
		fmt.Printf("Failed to insert roles. err: %s", err)
		return
	}
	fmt.Println(ret)
}

func InsertDefaultUsers() {
	if _, err := GetUserByName(ROLE_SYSADMIN); err == nil {
		return
	}
	roleSysadmin, _ := GetRoleByName(ROLE_SYSADMIN)
	sql := `INSERT IGNORE user (name, password_hash, email, role_id) 
			VALUES (:name, :password_hash, :email, :role_id)`
	args := map[string]interface{}{
		"name":          ROLE_SYSADMIN,
		"password_hash": "",
		"email":         "123@qq.com",
		"role_id":       roleSysadmin.ID,
	}
	_, err := conf.DB.NamedExec(sql, args)
	if err != nil {
		fmt.Printf("Db err: %s", err)
	}
}

func GetRoleByName(name string) (models.Role, error) {
	var role models.Role
	err := conf.DB.Get(&role, "SELECT * FROM role WHERE name = ?", name)
	if err != nil {
		conf.Log.Error("Db error, err: %s", err)
	}
	return role, err
}

func GetRoleByID(ID int) (models.Role, error) {
	var role models.Role
	err := conf.DB.Get(&role, "SELECT * FROM role WHERE id = ?", ID)
	if err != nil {
		conf.Log.Error("Db error, err: %s", err)
	}
	return role, err
}

func GetRoles() ([]models.Role, error) {
	var roles []models.Role
	err := conf.DB.Select(&roles, "SELECT * FROM role")
	if err != nil {
		fmt.Printf("Db error, err: %s", err)
	}
	return roles, err
}

func GetUserByName(name string) (models.User, error) {
	sql := "SELECT * FROM user WHERE name = ? LIMIT 1"
	var u models.User
	err := conf.DB.Get(&u, sql, name)
	if err != nil {
		fmt.Printf("Failed to get user, err: %v", err)
	}
	return u, err
}

func GetUserByID(ID int) (models.User, error) {
	sql := "SELECT * FROM user WHERE id = ?"
	var u models.User
	err := conf.DB.Get(&u, sql, ID)
	if err != nil {
		conf.Log.Error("Failed to get user, err: %v", err)
		return u, err
	}
	u.Role, err = GetRoleByID(u.RoleID)
	u.RoleName = u.Role.Name
	if err != nil {
		return u, err
	}
	return u, err
}

func UserCan(user *models.User, perm int) bool {
	return user.Role.Permissions & perm == perm
}

func GetUserByIDs(userIDs ...int) ([]models.User, error) {
	var users []models.User
	sql := "SELECT * FROM user WHERE id IN (?)"
	query, args, err := sqlx.In(sql, userIDs)
	if err != nil {
		conf.Log.Error("Db query error.\nsql: %s\nerr: %s", sql, err)
		return users, errors.New(0, "Internal error.")
	}
	query = conf.DB.Rebind(query)
	err = conf.DB.Select(&users, query, args...)
	if err != nil {
		conf.Log.Error("Failed to get user, err: %v", err)
		return users, err
	}
	return users, err
}


func HasUserPermission(userID, perms int) {

}


func GenPasswordHash(p string) string {
	pwd := []byte(p)
	pHash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		fmt.Printf("Failed to generate password hash, err: %s", err)
	}
	return string(pHash)
}


func CheckPasswordHash(pwdHash string, plainPwd string) bool {
	bytePlainPwd := []byte(plainPwd)
	bytePwdHash := []byte(pwdHash)
	err := bcrypt.CompareHashAndPassword(bytePwdHash, bytePlainPwd)
	if err != nil {
		return false
	}
	return true
}


func AddUser(user entity.UserRegister) (err error){
	var u models.User
	err = conf.DB.Get(&u, "SELECT * FROM user WHERE name = ? LIMIT 1", user.Name)
	if err == nil {
		return errs.UserIsExistedErr
	}
	role, err := GetRoleByName(user.Role)
	if err != nil {
		return err
	}
	sql := `INSERT INTO user
				(name, password_hash, email, role_id, uuid, status)
			VALUES (:name, :password_hash, :email, :role_id, REPLACE(UUID(), '-', ''), :status)`
	args := map[string]interface{}{
		"name": user.Name,
		"password_hash": GenPasswordHash(user.Password),
		"email": user.Email,
		"role_id": role.ID,
		"status": USER_STATE_ACTIVE,
	}
	_, err = conf.DB.NamedExec(sql, args)
	if err != nil {
		conf.Log.Error("Db insert error: %s", err)
	}
	return
}


func List(offset, limit string) (gin.H, error) {
	var users []models.User
	var args []interface{}
	limits := ""
	if offset != "" && limit != ""{
		limits = " LIMIT ?, ?"
		args = append(args, offset, limit)
	}
	sql := "SELECT * FROM user"
	if limits != "" {
		sql += limits
	}
	err := conf.DB.Select(&users, sql, args...)
	if err != nil {
		conf.Log.Error("Db list error: %s", err)
	}
	roles, _ := GetRoles()
	for i, u := range users {
		for _, r := range roles {
			if r.ID == u.RoleID {
				users[i].RoleName = r.Name
				break
			}
		}
	}
	var cnt int
	sql = "SELECT COUNT(*) AS cnt FROM user"
	conf.DB.QueryRow(sql).Scan(&cnt)
	return gin.H{"rows": users, "total": cnt}, err
}

func Follow(followerID, followedID int) error {
	if followerID == followedID {
		return fmt.Errorf("cannot follow yourself")
	}
	_, err := GetUserByID(followedID)
	if err != nil {
		conf.Log.Error("Failed to follow, unknown user (id %d)", followedID)
		return fmt.Errorf("failed to follow, unknown user (id %d)", followedID)
	}
	sql := "INSERT IGNORE follows (follower_id, followed_id) VALUES (:follower_id, :followed_id)"
	args := map[string]interface{}{
		"follower_id": followerID,
		"followed_id": followedID,
	}
	_, err = conf.DB.NamedExec(sql, args)
	if err != nil {
		conf.Log.Error("Db insert err.\nsql: %s\nerr: %s", sql, err)
		return errors.New(0, "Internal error.")
	}
	return nil
}

func UnFollow(followerID, followedID int) error {
	sql := "SELECT * FROM follows WHERE follower_id = ? AND followed_id = ?"
	var follower models.Follow
	err := conf.DB.Get(&follower, sql, followerID, followedID)
	if err != nil {
		conf.Log.Error("Unknown follow (follower_id %d, followed_id %d)", followerID, followedID)
		return err
	}
	sql = "DELETE FROM follows WHERE follower_id = ? AND followed_id = ?"
	_, err = conf.DB.Exec(sql, followerID, followedID)
	return err
}

// @param target int "follower or following"
func ListFollows(userID int, target string) (gin.H, error) {
	var follows []models.Follow
	var users []models.User
	var sql string
	if target == "follower" {
		sql = "SELECT * FROM follows WHERE follower_id = ?"
	} else if target == "following" {
		sql = "SELECT * FROM follows WHERE followed_id = ?"
	}
	err := conf.DB.Select(&follows, sql, userID)
	if err != nil {
		conf.Log.Error("Db query error.\nsql: %s\nerr: %s", sql, err)
		return nil, errors.New(0, "Internal error.")
	}
	var userIDs []int
	for _, v := range follows {
		userIDs = append(userIDs, v.FollowerID)
	}
	sql = "SELECT name, uuid FROM user WHERE id IN (?)"
	query, args, err := sqlx.In(sql, userIDs)
	if err != nil {
		conf.Log.Error("Db query error.\nsql: %s\nerr: %s", sql, err)
		return nil, errors.New(0, "Internal error.")
	}
	query = conf.DB.Rebind(query)
	err = conf.DB.Select(&users, query, args...)
	if err != nil {
		conf.Log.Error("Db query error.\nsql: %s\nerr: %s", sql, err)
		return nil, errors.New(0, "Internal error.")
	}
	var ret []interface{}
	for _, u := range users {
		ret = append(ret, gin.H{
			"name": u.Name,
			"uuid": u.Uuid,
		})
	}
	return gin.H{"rows": ret, "total": len(ret)}, err
}
