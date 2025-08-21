package apis

import (
	"context"
	"errors"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/nanoteck137/pyrin"
	"github.com/nanoteck137/watchbook/core"
	"github.com/nanoteck137/watchbook/database"
	"github.com/nanoteck137/watchbook/kvstore"
	"github.com/nanoteck137/watchbook/types"
)

type Notification struct {
	Id     string `json:"id"`
	UserId string `json:"userId"`

	Type     types.NotificationType `json:"notificationType"`
	Title    string                 `json:"title"`
	Message  string                 `json:"message"`
	Metadata map[string]string      `json:"metadata"`
	IsRead   bool                   `json:"isRead"`
}

type GetNotifications struct {
	Page          types.Page     `json:"page"`
	Notifications []Notification `json:"notifications"`
}

type GetNotificationById struct {
	Notification
}

func ConvertDBNotification(c pyrin.Context, notification database.Notification) Notification {
	return Notification{
		Id:       notification.Id,
		UserId:   notification.UserId,
		Type:     notification.Type,
		Title:    notification.Title,
		Message:  notification.Message,
		Metadata: notification.Metadata,
		IsRead:   notification.IsRead > 0,
	}
}

func InstallNotificationHandlers(app core.App, group pyrin.Group) {
	group.Register(
		pyrin.ApiHandler{
			Name:         "GetNotifications",
			Method:       http.MethodGet,
			Path:         "/notifications",
			ResponseType: GetNotifications{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				q := c.Request().URL.Query()
				opts := getPageOptions(q)

				ctx := context.TODO()

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				filterStr := q.Get("filter")
				sortStr := q.Get("sort")
				notifications, page, err := app.DB().GetPagedNotifications(ctx, user.Id, filterStr, sortStr, opts)
				if err != nil {
					return nil, err
				}

				res := GetNotifications{
					Page:          page,
					Notifications: make([]Notification, len(notifications)),
				}

				for i, notification := range notifications {
					res.Notifications[i] = ConvertDBNotification(c, notification)
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetNotificationById",
			Method:       http.MethodGet,
			Path:         "/notifications/:id",
			ResponseType: GetNotificationById{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")

				ctx := context.TODO()

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				notification, err := app.DB().GetNotificationById(ctx, id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, NotificationNotFound()
					}

					return nil, err
				}

				if notification.UserId != user.Id {
					return nil, NotificationNotFound()
				}

				return GetNotificationById{
					Notification: ConvertDBNotification(c, notification),
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:   "MarkNotificationRead",
			Method: http.MethodPost,
			Path:   "/notifications/:id/read",
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")

				ctx := context.TODO()

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				notification, err := app.DB().GetNotificationById(ctx, id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, NotificationNotFound()
					}

					return nil, err
				}

				if notification.UserId != user.Id {
					return nil, NotificationNotFound()
				}

				err = app.DB().UpdateNotification(ctx, notification.Id, database.NotificationChanges{
					IsRead: database.Change[int]{
						Value:   1,
						Changed: true,
					},
				})
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},

		// pyrin.ApiHandler{
		// 	Name:         "CreateNotification",
		// 	Method:       http.MethodPost,
		// 	Path:         "/notifications",
		// 	ResponseType: CreateNotification{},
		// 	BodyType:     CreateNotificationBody{},
		// 	HandlerFunc: func(c pyrin.Context) (any, error) {
		// 		// TODO(patrik): Add admin check
		//
		// 		body, err := pyrin.Body[CreateNotificationBody](c)
		// 		if err != nil {
		// 			return nil, err
		// 		}
		//
		// 		ctx := context.Background()
		//
		// 		ty := types.NotificationType(body.NotificationType)
		//
		// 		id, err := app.DB().CreateNotification(ctx, database.CreateNotificationParams{
		// 			Type: ty,
		// 			Name: body.Name,
		// 		})
		// 		if err != nil {
		// 			return nil, err
		// 		}
		//
		// 		notificationDir := app.WorkDir().NotificationDirById(id)
		// 		dirs := []string{
		// 			notificationDir.String(),
		// 			notificationDir.Images(),
		// 		}
		//
		// 		for _, dir := range dirs {
		// 			err = os.Mkdir(dir, 0755)
		// 			if err != nil && !os.IsExist(err) {
		// 				return nil, err
		// 			}
		// 		}
		//
		// 		return CreateNotification{
		// 			Id: id,
		// 		}, nil
		// 	},
		// },

		// pyrin.ApiHandler{
		// 	Name:         "EditNotification",
		// 	Method:       http.MethodPatch,
		// 	Path:         "/notifications/:id",
		// 	ResponseType: nil,
		// 	BodyType:     EditNotificationBody{},
		// 	HandlerFunc: func(c pyrin.Context) (any, error) {
		// 		id := c.Param("id")
		//
		// 		// TODO(patrik): Add admin check
		//
		// 		body, err := pyrin.Body[EditNotificationBody](c)
		// 		if err != nil {
		// 			return nil, err
		// 		}
		//
		// 		ctx := context.Background()
		//
		// 		dbNotification, err := app.DB().GetNotificationById(ctx, nil, id)
		// 		if err != nil {
		// 			if errors.Is(err, database.ErrItemNotFound) {
		// 				return nil, NotificationNotFound()
		// 			}
		//
		// 			return nil, err
		// 		}
		//
		// 		changes := database.NotificationChanges{}
		//
		// 		if body.NotificationType != nil {
		// 			t := types.NotificationType(*body.NotificationType)
		//
		// 			changes.Type = database.Change[types.NotificationType]{
		// 				Value:   t,
		// 				Changed: t != dbNotification.Type,
		// 			}
		// 		}
		//
		// 		if body.Name != nil {
		// 			changes.Name = database.Change[string]{
		// 				Value:   *body.Name,
		// 				Changed: *body.Name != dbNotification.Name,
		// 			}
		// 		}
		//
		// 		err = app.DB().UpdateNotification(ctx, dbNotification.Id, changes)
		// 		if err != nil {
		// 			return nil, err
		// 		}
		//
		// 		return nil, nil
		// 	},
		// },

		pyrin.ApiHandler{
			Name:         "DeleteNotification",
			Method:       http.MethodDelete,
			Path:         "/notifications/:id",
			ResponseType: nil,
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				ctx := context.Background()

				dbNotification, err := app.DB().GetNotificationById(ctx, id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, NotificationNotFound()
					}

					return nil, err
				}

				if dbNotification.UserId != user.Id {
					return nil, NotificationNotFound()
				}

				err = app.DB().RemoveNotification(ctx, dbNotification.Id)
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "TestNotification",
			Method:       http.MethodPost,
			Path:         "/notifications/test",
			ResponseType: nil,
			HandlerFunc: func(c pyrin.Context) (any, error) {
				ctx := context.Background()

				users, err := app.DB().GetAllUsers(ctx)
				if err != nil {
					return nil, err
				}

				for _, user := range users {
					_, err := app.DB().CreateNotification(ctx, database.CreateNotificationParams{
						UserId:  user.Id,
						Type:    types.NotificationTypeGeneric,
						Title:   "Test Notification",
						Message: "Super cool message",
						Metadata: kvstore.Store{
							"lel": "wot",
						},
						IsRead:   0,
						DedupKey: strconv.Itoa(rand.Int()),
					})
					if err != nil {
						return nil, err
					}
				}

				return nil, nil
			},
		},
	)
}
