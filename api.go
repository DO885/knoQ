package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

// GetHello テスト用API
func GetHello(c echo.Context) error {
	id := getRequestUser(c)                                      // リクエストしてきたユーザーのtraQID取得
	return c.String(http.StatusOK, fmt.Sprintf("hello %s!", id)) // レスポンスを返す
}

// RoomsAPI

// PostRoom traPで確保した部屋情報を作成
func PostRoom(c echo.Context) error {
	r := new(Room)
	if err := c.Bind(r); err != nil {
		return err
	}

	if err := db.Create(&r).Error; err != nil {
		return err
	}
	return c.JSON(http.StatusOK, r)
}

// GetRooms traPで確保した部屋情報を取得
func GetRooms(c echo.Context) error {
	r := []Room{}
	begin := c.QueryParam("date_begin")
	end := c.QueryParam("date_end")

	r, err := findRoomsByTime(begin, end)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, r)
}

// DeleteRoom traPで確保した部屋情報を削除
func DeleteRoom(c echo.Context) error {
	r := new(Room)
	r.ID, _ = strconv.Atoi(c.Param("roomid"))

	if err := db.First(&r, r.ID).Error; err != nil {
		return c.String(http.StatusNotFound, "部屋が存在しない")
	}

	if err := db.Delete(&r).Error; err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	return c.NoContent(http.StatusOK)
}

// groupsAPI

// PostGroup グループを作成
func PostGroup(c echo.Context) error {
	g := new(Group)

	if err := c.Bind(&g); err != nil {
		return err
	}

	g.CreatedByRefer = getRequestUser(c)
	if err := db.Where("traq_id = ?", g.CreatedByRefer).First(&g.CreatedBy).Error; err != nil {
		return err
	}

	// メンバーがdbにいるか
	if err := checkMembers(g); err != nil {
		return c.String(http.StatusBadRequest, "正しくないメンバーが含まれている")
	}

	if err := db.Create(&g).Error; err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprint(err))
	}

	return c.JSON(http.StatusOK, g)
}

// GetGroups グループを取得
func GetGroups(c echo.Context) error {
	groups := []Group{}
	traqID := c.QueryParam("userid")

	groups, err := findGroupsBelong(traqID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, groups)
}

// DeleteGroup グループを削除
func DeleteGroup(c echo.Context) error {
	g := new(Group)
	g.ID, _ = strconv.Atoi(c.Param("groupid"))

	if err := db.First(&g, g.ID).Related(&g.Members, "Members").Error; err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	if err := db.Model(&g).Association("Members").Clear().Error; err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	if err := db.Delete(&g).Error; err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	return c.NoContent(http.StatusOK)
}

// UpdateGroup グループメンバー、グループ名を更新
func UpdateGroup(c echo.Context) error {
	g := new(Group)

	if err := c.Bind(g); err != nil {
		return err
	}

	// メンバーがdbにいるか
	if err := checkMembers(g); err != nil {
		return c.String(http.StatusBadRequest, "正しくないメンバーが含まれている")
	}

	g.ID, _ = strconv.Atoi(c.Param("groupid"))
	if err := db.First(&g, g.ID).Error; err != nil {
		return err
	}
	// 作成者を取得
	if err := db.Where("traq_id = ?", g.CreatedByRefer).First(&g.CreatedBy).Error; err != nil {
		return err
	}
	if getRequestUser(c) != g.CreatedByRefer {
		return echo.NewHTTPError(http.StatusForbidden, "作成者ではない")
	}

	// メンバーを置き換え
	if err := db.Model(&g).Association("Members").Replace(g.Members).Error; err != nil {
		return err
	}

	// グループ名を変更
	if err := db.Model(&g).Update("name", g.Name).Error; err != nil {
		return err
	}

	return c.JSON(http.StatusOK, g)
}

// resrvations API

// PostReservation 部屋の使用宣言を作成
func PostReservation(c echo.Context) error {
	rv := new(Reservation)

	if err := c.Bind(&rv); err != nil {
		return err
	}

	rv.CreatedByRefer = getRequestUser(c)
	if err := db.Where("traq_id = ?", rv.CreatedByRefer).First(&rv.CreatedBy).Error; err != nil {
		return err
	}

	// groupがあるか
	if err := checkGroup(rv.GroupID, &rv.Group); err != nil {
		return c.String(http.StatusBadRequest, "groupが存在しません"+fmt.Sprintln(rv.GroupID))
	}
	// roomがあるか
	if err := checkRoom(rv.RoomID, &rv.Room); err != nil {
		return c.String(http.StatusBadRequest, "roomが存在しません")
	}

	// dateを代入
	r := new(Room)
	if err := db.First(&r, rv.RoomID).Error; err != nil {
		return err
	}
	// r.Date = 2018-08-10T00:00:00+09:00
	rv.Date = r.Date[:10]
	rv.Room.Date = rv.Room.Date[:10]

	if err := db.Create(&rv).Error; err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, rv)
}

// GetReservations 部屋の使用宣言情報を取得
func GetReservations(c echo.Context) error {
	reservations := []Reservation{}

	traqID := c.QueryParam("userid")
	groupID := c.QueryParam("groupid")
	begin := c.QueryParam("date_begin")
	end := c.QueryParam("date_end")
	reservations, err := findRvs(traqID, groupID, begin, end)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, reservations)
}

// DeleteReservation 部屋の使用宣言を削除
func DeleteReservation(c echo.Context) error {
	rv := new(Reservation)
	rv.ID, _ = strconv.Atoi(c.Param("reservationid"))

	traQID := getRequestUser(c)
	belong, err := checkBelongToGroup(rv.ID, traQID)
	if err != nil {
		return err
	}
	if !belong {
		return c.String(http.StatusForbidden, "削除できるのは所属ユーザーのみです。")
	}

	if err := db.Delete(&rv).Error; err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	return c.NoContent(http.StatusOK)
}

// UpdateReservation 部屋、開始時刻、終了時刻を更新
func UpdateReservation(c echo.Context) error {
	rv := new(Reservation)

	if err := c.Bind(&rv); err != nil {
		return err
	}
	rv.ID, _ = strconv.Atoi(c.Param("reservationid"))

	traQID := getRequestUser(c)
	belong, err := checkBelongToGroup(rv.ID, traQID)
	if err != nil {
		fmt.Println("I could not check", err)
		return err
	}
	if !belong {
		return echo.NewHTTPError(http.StatusForbidden, "変更できるのは所属ユーザーのみです you are "+traQID)
	}

	// roomがあるか
	if err := checkRoom(rv.RoomID, &rv.Room); err != nil {
		return c.String(http.StatusBadRequest, "roomが存在しません")
	}

	// r.Date = 2018-08-10T00:00:00+09:00
	rv.Room.Date = rv.Room.Date[:10]
	rv.Date = rv.Room.Date

	// roomid, timestart, timeendのみを変更(roomidに伴ってdateの変更する)
	if err := db.Model(&rv).Update(Reservation{RoomID: rv.RoomID, Date: rv.Date, TimeStart: rv.TimeStart, TimeEnd: rv.TimeEnd}).Error; err != nil {
		fmt.Println("DB could not be updated")
		return err
	}

	if err := db.First(&rv, rv.ID).Error; err != nil {
		return err
	}

	if err := rv.Group.AddRelation(rv.GroupID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "GroupRelationを追加できませんでした")
	}

	if err := rv.AddCreatedBy(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "ReservationCreatedByを追加できませんでした")
	}

	return c.JSON(http.StatusOK, rv)
}
