package main

// getUser ユーザー情報を取得します
func getUser(id string) (*User, error) {
	user := User{}

	// DBに登録されていない場合(初めてアクセスした場合)はDBにレコードを作成する
	if err := db.FirstOrCreate(&user, &User{TRAQID: id}).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// changeUserToAdmin ユーザーの管理者権限の有無を変更します
func changeUserToAdmin(id string, isAdmin bool) error {
	// ユーザー取得
	user, err := getUser(id)
	if err != nil {
		return err
	}

	// 変更
	if err := db.Model(user).Update("admin", isAdmin).Error; err != nil {
		return err
	}
	return nil
}

// checkMembers グループのメンバーがdbにいるか
func checkMembers (group *Group) error{
	for _, u := range group.Members{
		if err := db.Where("traq_id = ?", u.TRAQID).First(&u).Error; err != nil{
			return err
		}
	}
	return nil
}

func checkGroup (groupID int) error {
	g := new(Group)
	g.ID = groupID
	if err := db.First(&g, g.ID).Error; err != nil{
		return err
	}
	return nil
}

func checkRoom (roomID int) error{
	r := new(Room)
	r.ID = roomID
	if err := db.First(&r, r.ID).Error; err != nil{
		return err
	}
	return nil
}