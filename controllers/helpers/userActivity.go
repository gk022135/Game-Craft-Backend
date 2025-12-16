package helpers

import (
	"context"

	db "github.com/gk022135/Game-Craft-Backend/prisma/db"
)

func UserActivity(userId int, questionID int, userQuery string, isValid bool) error {
	client := db.NewClient()
	if err := client.Connect(); err != nil {
		return err
	}
	defer client.Disconnect()

	ctx := context.Background()

	_, err := client.UserActivityLog.CreateOne(
		db.UserActivityLog.Solution.Set(userQuery),
		db.UserActivityLog.User.Link(
			db.User.ID.Equals(userId),
		),
		db.UserActivityLog.Question.Link(
			db.QuestionRecords.ID.Equals(questionID),
		),
		db.UserActivityLog.IsValid.Set(isValid),
	).Exec(ctx)

	return err
}
