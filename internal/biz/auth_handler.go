package biz

import "flynoob/bibirt-sock/pkg/websocket"

func (cu *ConnUseCase) Auth(client *websocket.Client) {
	client.OnUpgrade(func(c *websocket.Client) error {
		tok := c.Req.Form.Get("_token")
		uuid, err := cu.authService.ConnUUID(tok)
		if err != nil {
			return err
		}

		client.Set("uuid", uuid)

		return nil
	})
}
