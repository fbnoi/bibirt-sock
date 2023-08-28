package biz

import "flynoob/bibirt-sock/pkg/websocket"

func Auth(client *websocket.Client) {
	client.OnUpgrade(func(c *websocket.Client) error {
		return nil
	})
}
