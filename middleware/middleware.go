package middleware

import (
	"net/http"
	repo "room/repository"

	"github.com/labstack/echo"
)

const traQID = "traQID"

// TraQUserMiddleware traQユーザーか判定するミドルウェア
func TraQUserMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Request().Header.Get("X-Showcase-User")
		if len(id) == 0 || id == "-" {
			// test用
			id = "fuji"
		}
		c.Set(traQID, id)
		return next(c)
	}
}

// AdminUserMiddleware 管理者ユーザーか判定するミドルウェア
func AdminUserMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := GetRequestUser(c)
		if len(id) == 0 {
			return echo.NewHTTPError(http.StatusForbidden) // traQにログインが必要
		}

		// ユーザー情報を取得
		user, err := repo.GetUser(id)
		if err != nil {
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError) // データベースエラー
		}

		// 判定
		if !user.Admin {
			return echo.NewHTTPError(http.StatusForbidden) // 管理者ユーザーでは無いのでエラー
		}

		return next(c)
	}
}

// GetRequestUser リクエストユーザーのtraQIDを返します
func GetRequestUser(c echo.Context) string {
	return c.Get(traQID).(string)
}
