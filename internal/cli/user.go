package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

const usersBaseURL = "v1/users"
const userSpecificURL = usersBaseURL + "/{user_email}"
const userInviteURL = usersBaseURL + "/invite"

var userReadFields = []string{
	"email",
	"status",
	"teams",
	"roles.name",
	"managedByScim",
	"inactivityTimeout",
}

func (c *PortClient) ReadUser(ctx context.Context, email string) (*User, int, error) {
	var body UserBody
	resp, err := c.Client.R().
		SetContext(ctx).
		SetHeader("Accept", "application/json").
		SetPathParam("user_email", email).
		SetQueryParamsFromValues(url.Values{
			"fields": userReadFields,
		}).
		SetResult(&body).
		Get(userSpecificURL)
	if err != nil {
		return nil, 0, err
	} else if resp.IsError() || !body.OK {
		return nil, resp.StatusCode(), fmt.Errorf("failed to read user, got: %s", resp.Body())
	}
	if body.User == nil {
		return nil, resp.StatusCode(), fmt.Errorf("port-api returned an invalid response: user is nil")
	}

	return body.User, resp.StatusCode(), nil
}

func (c *PortClient) InviteUser(ctx context.Context, invite *UserInviteRequest, notify bool) error {
	resp, err := c.Client.R().
		SetBody(invite).
		SetContext(ctx).
		SetQueryParam("notify", strconv.FormatBool(notify)).
		Post(userInviteURL)
	if err != nil {
		return err
	}

	var body PortBody
	if err := json.Unmarshal(resp.Body(), &body); err != nil {
		return err
	}
	if !body.OK {
		return fmt.Errorf("failed to invite user, got: %s", resp.Body())
	}

	return nil
}

func (c *PortClient) UpdateUser(ctx context.Context, email string, update *UserUpdate) (*User, error) {
	resp, err := c.Client.R().
		SetBody(update.ToPatchBody()).
		SetContext(ctx).
		SetPathParam("user_email", email).
		Patch(userSpecificURL)
	if err != nil {
		return nil, err
	}

	var body PortBody
	if err := json.Unmarshal(resp.Body(), &body); err != nil {
		return nil, err
	}
	if !body.OK {
		return nil, fmt.Errorf("failed to update user, got: %s", resp.Body())
	}

	user, _, err := c.ReadUser(ctx, email)
	return user, err
}

func (c *PortClient) DeleteUser(ctx context.Context, email string) error {
	resp, err := c.Client.R().
		SetContext(ctx).
		SetHeader("Accept", "application/json").
		SetPathParam("user_email", email).
		Delete(userSpecificURL)
	if err != nil {
		return err
	}

	var body PortBodyDelete
	if err := json.Unmarshal(resp.Body(), &body); err != nil {
		return err
	}
	if !body.Ok {
		return fmt.Errorf("failed to delete user, got: %s", string(resp.Body()))
	}

	return nil
}
