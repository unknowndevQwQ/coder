package codersdk

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/xerrors"
)

type CreateGroupRequest struct {
	Name string `json:"name"`
}

type Group struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	OrganizationID uuid.UUID `json:"organization_id"`
	Members        []User    `json:"members"`
}

func (c *Client) CreateGroup(ctx context.Context, orgID uuid.UUID, req CreateGroupRequest) (Group, error) {
	res, err := c.Request(ctx, http.MethodPost,
		fmt.Sprintf("/api/v2/organizations/%s/groups", orgID.String()),
		req,
	)
	if err != nil {
		return Group{}, xerrors.Errorf("make request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		return Group{}, readBodyAsError(res)
	}
	var resp Group
	return resp, json.NewDecoder(res.Body).Decode(&resp)
}

func (c *Client) GroupsByOrganization(ctx context.Context, orgID uuid.UUID) ([]Group, error) {
	res, err := c.Request(ctx, http.MethodGet,
		fmt.Sprintf("/api/v2/organizations/%s/groups", orgID.String()),
		nil,
	)
	if err != nil {
		return nil, xerrors.Errorf("make request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, readBodyAsError(res)
	}

	var groups []Group
	return groups, json.NewDecoder(res.Body).Decode(&groups)
}

func (c *Client) Group(ctx context.Context, group uuid.UUID) (Group, error) {
	res, err := c.Request(ctx, http.MethodGet,
		fmt.Sprintf("/api/v2/groups/%s", group.String()),
		nil,
	)
	if err != nil {
		return Group{}, xerrors.Errorf("make request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return Group{}, readBodyAsError(res)
	}
	var resp Group
	return resp, json.NewDecoder(res.Body).Decode(&resp)
}

type PatchGroupRequest struct {
	AddUsers    []string `json:"add_users"`
	RemoveUsers []string `json:"remove_users"`
	Name        string   `json:"name"`
}

func (c *Client) PatchGroup(ctx context.Context, group uuid.UUID, req PatchGroupRequest) (Group, error) {
	res, err := c.Request(ctx, http.MethodPatch,
		fmt.Sprintf("/api/v2/groups/%s", group.String()),
		req,
	)
	if err != nil {
		return Group{}, xerrors.Errorf("make request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return Group{}, readBodyAsError(res)
	}
	var resp Group
	return resp, json.NewDecoder(res.Body).Decode(&resp)
}

func (c *Client) DeleteGroup(ctx context.Context, group uuid.UUID) error {
	res, err := c.Request(ctx, http.MethodDelete,
		fmt.Sprintf("/api/v2/groups/%s", group.String()),
		nil,
	)
	if err != nil {
		return xerrors.Errorf("make request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return readBodyAsError(res)
	}
	return nil
}
