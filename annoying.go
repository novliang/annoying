package annoying

import (
	"errors"
	"github.com/jinzhu/gorm"
)

var Db *gorm.DB

type ItemInter interface {
	GetName() string
}

type Role struct {
	*Item
}

type Permission struct {
	*Item
}

type Item struct {
	*AuthItem
}

type Assignment struct {
	*AuthAssignment
}

//Table
const (
	AuthAssignmentTable = "auth_assignment"
	AuthRuleTable       = "auth_rule"
	AuthItemTable       = "auth_item"
	AuthItemChildTable  = "auth_item_child"
)

//Default Role
var DefaultRoles []string

func LoadDbInstance(db *gorm.DB) {
	Db = db
}

type AssignmentsNameIndexed map[string]Assignment

func CheckAssess(userId string, permissionName string, params map[string]string) bool {

	assignments := GetAssignments(userId)
	if HasNoAssignments(assignments) {
		return false
	}

	return CheckAccessRecursive(userId, permissionName, params, assignments)
}

func CheckAccessRecursive(userId string, itemName string, params map[string]string, assignments *AssignmentsNameIndexed) bool {
	item := GetItem(itemName)

	if item == nil {
		return false
	}

	if _, ok := (*assignments)[itemName]; ok {
		return true;
	}

	for _, v := range DefaultRoles {
		if v == itemName {
			return true
		}
	}

	parents := []string{}

	Db.Find(&AuthItemChild{}).Pluck("parents", parents)

	for _, v := range parents {
		if CheckAccessRecursive(userId, v, params, assignments) {
			return true
		}
	}
	return false
}

func HasNoAssignments(assignments *AssignmentsNameIndexed) bool {
	return (assignments == nil || len(*assignments) == 0) && len(DefaultRoles) == 0
}

func GetAssignment(roleName string, userId string) (*Assignment) {
	if userId == "" {
		return nil
	}

	i := Assignment{}

	if Db.Where("user_id = ? and item_name = ?", userId, roleName).First(&i).RecordNotFound() {
		return nil
	}

	return &i
}

func GetAssignments(userId string) (*AssignmentsNameIndexed) {

	if userId == "" {
		return nil
	}

	assignmentsRaw := []Assignment{}
	Db.Where("user_id = ?", userId).Find(&assignmentsRaw)

	assignments := AssignmentsNameIndexed{}

	for _, v := range assignmentsRaw {
		assignments[v.ItemName] = v
	}
	return &assignments
}

func CanAddChild(parent ItemInter, child ItemInter) bool {
	return !DetectLoop(parent, child);
}

func AddChild(parent ItemInter, child ItemInter) error {
	if parent.GetName() == child.GetName() {
		return errors.New("Cannot add " + parent.GetName() + " as a child of itself.")
	}

	if _, o1 := parent.(*Permission); o1 {
		if _, o2 := child.(*Role); o2 {
			return errors.New("Cannot add a role as a child of a permission")
		}
	}

	if DetectLoop(parent, child) {
		return errors.New("Cannot add '" + child.GetName() + "' as a child of '" + parent.GetName() + "'. A loop has been detected.")
	}

	i := AuthItemChild{
		Parent: parent.GetName(),
		Child:  child.GetName(),
	}

	if err := Db.Create(i).Error; err != nil {
		return err
	}

	return nil
}

func RemoveChild(parent ItemInter, child ItemInter) error {
	if err := Db.Where("parent = ? and child = ?", parent.GetName(), child.GetName()).Delete(&AuthItemChild{}).Error; err != nil {
		return err
	}
	return nil
}

func RemoveChildren(parent ItemInter) error {
	if err := Db.Where("parent = ?", parent.GetName()).Delete(&AuthItemChild{}).Error; err != nil {
		return err
	}
	return nil
}

func HasChild(parent ItemInter, child ItemInter) bool {
	if !Db.Where("parent = ? and child = ? ", parent.GetName(), child.GetName()).First(&AuthItemChild{}).RecordNotFound() {
		return true
	}
	return false
}

func GetChildren(name string) map[string]interface{} {
	children := map[string]interface{}{}
	i := []Item{}
	Db.Raw("SELECT name, type, description,rule_name, data, created_at, updated_at from " + AuthItemTable + "," + AuthItemChildTable + " where parent = '" + name + "' and name = child").Scan(&i)
	for _, v := range i {
		children[v.Name] = PopulateItem(v)
	}
	return children
}

func PopulateItem(item Item) interface{} {
	if item.Type == TypeRole {
		role := new(Role)
		role.Name = item.Name
		role.Type = item.Type
		role.Description = item.Description
		role.Data = item.Data
		role.CreatedAt = item.CreatedAt
		role.UpdatedAt = item.UpdatedAt
		return role
	} else if item.Type == TypePermission {
		permission := new(Permission)
		permission.Name = item.Name
		permission.Type = item.Type
		permission.Description = item.Description
		permission.Data = item.Data
		permission.CreatedAt = item.CreatedAt
		permission.UpdatedAt = item.UpdatedAt
		return permission
	}
	return nil
}

func DetectLoop(parent ItemInter, child ItemInter) bool {
	if child.GetName() == parent.GetName() {
		return true;
	}

	for _, v := range GetChildren(child.GetName()) {
		if DetectLoop(parent, v.(ItemInter)) {
			return true;
		}
	}
	return false;
}

func Assign(role *Role, userId string) (*Assignment, error) {
	assignment := new(Assignment)
	assignment.UserId = userId
	assignment.ItemName = role.Name
	err := Db.Create(assignment).Error
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return assignment, nil
}

func Revoke(role *Role, userId string) bool {
	if userId == "" {
		return false
	}
	d := Db.Where("user_id = ? and item_name", userId, role.Name).Delete(AuthAssignment{})

	if d.Error != nil {
		return false
	}

	return d.RowsAffected > 0
}

func RevokeAll(userId string) bool {
	if userId == "" {
		return false
	}

	d := Db.Where("user_id = ?", userId).Delete(AuthAssignment{})

	if d.Error != nil {
		return false
	}

	return d.RowsAffected > 0
}

func RemoveAll() error {
	err := RemoveAllAssignments()
	if err != nil {
		return err
	}
	err = RemoveAllItemChild()
	if err != nil {
		return err
	}
	err = RemoveAllItems()
	if err != nil {
		return err
	}
	err = RemoveAllRules()
	if err != nil {
		return err
	}
	return nil
}

func RemoveAllPermissions() error {
	err := RemoveAllItems()
	if err != nil {
		return err
	}
	return err
}

func RemoveAllItemChild() error {
	err := Db.Exec("delete from " + AuthItemChildTable + " where 1").Error
	if err != nil {
		return err
	}
	return nil
}

func RemoveAllRoles() error {
	err := RemoveAllItems()
	if err != nil {
		return err
	}
	return nil
	//Todo 删除缓存
}

func RemoveAllItems() error {
	err := Db.Exec("delete from " + AuthItemTable + " where 1").Error
	if err != nil {
		return err
	}
	return nil
	//Todo 删除缓存
}

func RemoveAllRules() error {
	err := Db.Exec("delete from " + AuthRuleTable + " where 1").Error
	if err != nil {
		return err
	}
	return nil
}

func RemoveAllAssignments() error {
	err := Db.Exec("delete from " + AuthAssignmentTable + " where 1").Error
	if err != nil {
		return err
	}
	return nil
}

func GetUserIdsByRole(roleName string) (userIds []string, err error) {
	if roleName == "" {
		return []string{}, nil;
	}

	assignments := []AuthAssignment{};
	if e := Db.Select([]string{"user_id"}).Where("item_name = ?", roleName).Find(&assignments).Error; e != nil {
		return nil, errors.New(e.Error())
	}

	for _, v := range assignments {
		userIds = append(userIds, v.UserId)
	}
	return
}

func GetItem(name string) interface{} {
	if name == "" {
		return nil
	}
	i := Item{}

	if Db.Where("name = ?", name).First(&i).RecordNotFound() {
		return nil
	}
	return PopulateItem(i)
}

func CreateRole(name string) *Role {
	r := new(Role)
	r.Name = name
	return r
}

func CreatePermission(name string) *Permission {
	p := new(Permission)
	p.Name = name
	return p
}
