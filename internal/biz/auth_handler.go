package biz

import "flynoob/bibirt-sock/pkg/websocket"

func (handler *ClientHandler) Auth(client *websocket.Client) {
	client.OnUpgrade(func(c *websocket.Client) error {
		tok := c.Req.Form.Get("_token")
		uuid, err := handler.authService.ConnUUID(tok)
		if err != nil {
			return err
		}

		client.Set("uuid", uuid)

		return nil
	})
}
