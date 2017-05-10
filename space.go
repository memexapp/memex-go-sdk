package memex

import (
	"encoding/json"
	"fmt"
	"time"
)

// SpaceType represents semantic type of space
type SpaceType string

const (
	// Origin represents starting point space for user
	Origin SpaceType = "com.memex.origin"
	// Text space represents textual data
	Text SpaceType = "com.memex.media.text"
	// WebPage space represents web link URL
	WebPage SpaceType = "com.memex.media.webpage"
	// Image space represents image/diagram space
	Image SpaceType = "com.memex.media.image"
	// Collection is set/list of links to other spaces
	Collection SpaceType = "com.memex.media.collection"
)

// Space represents folder/text/everything
type Space struct {
	MUID            *string     `json:"muid,omitempty"`
	CreatedAt       *time.Time  `json:"created_at,omitempty"`
	UpdatedAt       *time.Time  `json:"updated_at,omitempty"`
	VisitedAt       *time.Time  `json:"visited_at,omitempty"`
	State           EntityState `json:"state"`
	OwnerID         *int64      `json:"owner_id,omitempty"`
	Caption         *string     `json:"tag_label,omitempty"`
	Color           *string     `json:"tag_color,omitempty"`
	TypeIdentifier  SpaceType   `json:"type_identifier"`
	Representations *[]Media    `json:"representations,omitempty"`
}

type spaceResponse struct {
	Space Space `json:"space"`
}

type spacesRequest struct {
	Spaces []*Space `json:"spaces"`
}

// RepresentationWithType returns representation with specified media type
func (space *Space) RepresentationWithType(mediaType MediaType) *Media {
	if space.Representations == nil {
		return nil
	}
	for _, media := range *space.Representations {
		if media.MediaType == mediaType {
			return &media
		}
	}
	return nil
}

// GetSpace returns space with representations
func (spaces *Spaces) GetSpace(muid string) (*Space, error) {
	path := fmt.Sprintf("/spaces/%v", muid)
	var responseObject spaceResponse
	_, requestError := spaces.perform("GET", path, nil, &responseObject)
	if requestError != nil {
		return nil, requestError
	}
	return &responseObject.Space, nil
}

// UpdateSpaces updates spaces
func (spaces *Spaces) UpdateSpaces(array []*Space, ownerID int64) error {
	message := &spacesRequest{
		Spaces: array,
	}
	body, serializationError := json.Marshal(message)
	if serializationError != nil {
		return serializationError
	}
	path := fmt.Sprintf("/spaces/multiple")
	var responseObject spaceResponse
	_, requestError := spaces.perform("POST", path, body, &responseObject)
	if requestError != nil {
		return requestError
	}
	return nil
}

// UpdateSpace updates single space
func (spaces *Spaces) UpdateSpace(space *Space) error {
	array := []*Space{space}
	if space.OwnerID == nil {
		return fmt.Errorf("Missing ownerID")
	}
	return spaces.UpdateSpaces(array, *space.OwnerID)
}