package user_service

import (
	"ess/model/address"
	"ess/model/user"
	"ess/service/address_service"
	"ess/utils/amap_base"
	"ess/utils/cache"
	"ess/utils/db"
	"ess/utils/logging"
	"ess/utils/setting"
	"regexp"

	"github.com/jinzhu/copier"
)

func QueryUserByPhone(phone string) user.User {
	// phone to uid
	uid := cache.GetOrCreate(cache.GetKey(cache.PhoneToId, phone), func() interface{} {
		return Phone2Id(phone)
	}).(int)

	res := cache.GetOrCreate(cache.GetKey(cache.UserInfo, uid), func() interface{} {
		cacheUser, err := GetUserById(uid)
		if err != nil {
			logging.InfoF("a error occur in query user: %v\n", err)
			return &user.User{
				UserId: -1,
			}
		}
		return cacheUser
	}).(*user.User)

	return *res
}

func QueryUserByName(name string) user.User {
	// name to uid
	uid := cache.GetOrCreate(cache.GetKey(cache.NameToId, name), func() interface{} {
		return Name2Id(name)
	}).(int)

	res := cache.GetOrCreate(cache.GetKey(cache.UserInfo, uid), func() interface{} {
		cacheUser, err := GetUserById(uid)
		if err != nil {
			logging.InfoF("a error occur in query user: %v\n", err)
			return &user.User{
				UserId: -1,
			}
		}
		return cacheUser
	}).(*user.User)

	return *res
}

func QueryUserById(uid int) user.User {
	if uid == setting.AdminSetting.UserId {
		return user.User{
			UserId:   setting.AdminSetting.UserId,
			UserName: setting.AdminSetting.Name,
			UserRole: user.SysAdmin,
		}
	}
	res := cache.GetOrCreate(cache.GetKey(cache.UserInfo, uid), func() interface{} {
		cacheUser, err := GetUserById(uid)
		if err != nil {
			return &user.User{
				UserId: -1,
			}
		}
		return cacheUser
	}).(*user.User)

	return *res
}

func CreateUser(user *user.User) error {
	return db.MysqlDB.Create(user).Error
}

func UpdateUser(user *user.User) error {
	err := db.MysqlDB.Model(user).Updates(user).Error
	if err == nil {
		if err := db.MysqlDB.First(&user).Error; err != nil {
			return err
		}
		cache.Set(cache.GetKey(cache.UserInfo, user.UserId), user)
	}
	return err
}

func CleanUserCache(user user.User) {
	cache.Remove(cache.GetKey(cache.UserInfo, user.UserId))
}

func Phone2Id(phone string) int {
	usr := user.User{}
	if err := db.MysqlDB.Where(&user.User{UserPhone: phone}).First(&usr).Error; err != nil {
		return -1
	}
	return usr.UserId
}

func Name2Id(name string) int {
	usr := user.User{}
	if err := db.MysqlDB.Where(&user.User{UserName: name}).First(&usr).Error; err != nil {
		return -1
	}
	return usr.UserId
}

func GetUserById(uid int) (*user.User, error) {
	user := user.User{UserId: uid}
	if err := db.MysqlDB.First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

const PhoneRegex = "^1(3\\d|4[5-9]|5[0-35-9]|6[2567]|7[0-8]|8\\d|9[0-35-9])\\d{8}$"

// 长度 6-18 字母开头 只能包含数字 字母下划线
const PasswordRegex = "^[a-zA-Z]\\w{5,17}$"

// 4-16 位字母,数字,汉字,下划线
const NameRegex = "^([\u4e00-\u9fa5]{2,4})|([A-Za-z0-9_]{4,16})|([a-zA-Z0-9_\u4e00-\u9fa5]{3,16})$"

func ValidUser(user user.UserCreateReq) (*address.Address, bool) {
	if len(user.UserName) == 0 || len(user.UserPhone) == 0 || len(user.UserSecret) == 0 {
		return nil, false
	}

	if m, _ := regexp.MatchString(NameRegex, user.UserName); !m {
		logging.InfoF("用户名格式错误 name:%s\n", user.UserName)
		return nil, false
	}

	if m, _ := regexp.MatchString(PhoneRegex, user.UserPhone); !m {
		logging.InfoF("手机号格式错误 phoneNum:%s\n", user.UserPhone)
		return nil, false
	}

	if m, _ := regexp.MatchString(PasswordRegex, user.UserSecret); !m {
		logging.InfoF("密码格式错误 secret:%s\n", user.UserSecret)
		return nil, false
	}

	var addr address.Address
	_ = copier.Copy(&addr, &user.UserAddress)
	addr.AddressUserId = 0 // placeholder
	err := amap_base.GetCoordination(&addr)
	if err != nil {
		logging.ErrorF("获取坐标失败(%+v): %v\n", addr, err)
		return &addr, false
	}

	return &addr, true
}

func CreateUserWithAddress(user *user.User, addr *address.Address) error {
	err := address_service.CreateAddress(addr)
	if err != nil {
		return err
	}
	user.UserDefaultAddressId = addr.AddressId

	tx := db.MysqlDB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		return err
	}

	addr.AddressUserId = user.UserId
	if err := tx.Model(addr).Updates(addr).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func DeleteUserById(uid int) error {
	return db.MysqlDB.Where(&user.User{UserId: uid}).Delete(&user.User{UserId: uid}).Error
}

func QueryAvailableRiders() ([]user.User, error) {
	ret := []user.User{}
	err := db.MysqlDB.Model(&user.User{}).Where(&user.User{UserRole: user.Rider, UserAvailable: true}).Find(&ret).Error
	return ret, err
}

func QueryUsersByRole(role user.Role) ([]user.User, error) {
	ret := []user.User{}
	err := db.MysqlDB.Model(&user.User{}).Where(&user.User{UserRole: role}).Find(&ret).Error
	return ret, err
}
