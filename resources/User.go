package resources

import (
	"strconv"
	"time"

	"github.com/beevik/etree"
	"github.com/onego-project/onego/errors"
)

// User structure represents user and inherits XML data and methods from Resource structure
type User struct {
	Resource
}

// LoginToken structure represents login token of given user with attributes token, expiration time and egid
type LoginToken struct {
	Token          string
	ExpirationTime *time.Time
	EGID           int
}

// CreateUserWithID constructs User with id
func CreateUserWithID(id int) *User {
	return &User{*CreateResource("USER", id)}
}

// CreateUserFromXML constructs User with full xml data
func CreateUserFromXML(XMLdata *etree.Element) *User {
	return &User{Resource: Resource{XMLData: XMLdata}}
}

// Password method returns password for the given user
func (u *User) Password() (string, error) {
	return u.Attribute("PASSWORD")
}

// AuthDriver method returns authentication driver for the given user
func (u *User) AuthDriver() (string, error) {
	return u.Attribute("AUTH_DRIVER")
}

// MainGroup method returns main group ID of the given user
func (u *User) MainGroup() (int, error) {
	return u.intAttribute("GID")
}

// Groups method returns list of groups with main group at index 0
func (u *User) Groups() ([]int, error) {
	groups, err := u.arrayOfIDs("GROUPS")
	if err != nil {
		return nil, err
	}

	if len(groups) < 1 { // should never happen that user has no group
		return nil, &errors.XMLElementError{Path: "group"}
	}

	return groups, nil
}

// Enabled returns true for enabled
func (u *User) Enabled() (bool, error) {
	ret, err := u.intAttribute("ENABLED")
	if err != nil {
		return false, err
	}
	return intToBool(ret), nil
}

// LoginTokens returns list of login tokens for the given user
func (u *User) LoginTokens() ([]LoginToken, error) {
	elements := u.XMLData.FindElements("LOGIN_TOKEN")
	if len(elements) == 0 {
		return make([]LoginToken, 0), nil
	}

	loginTokens := make([]LoginToken, len(elements))

	for i, e := range elements {
		tokenElement := e.FindElement("TOKEN")
		if tokenElement == nil {
			return nil, &errors.XMLElementError{Path: "token in login token"}
		}

		expTimeElement := e.FindElement("EXPIRATION_TIME")
		if expTimeElement == nil {
			return nil, &errors.XMLElementError{Path: "expiration time in login token"}
		}

		expTime, err := strconv.ParseInt(expTimeElement.Text(), base10, bitSize64)
		if err != nil {
			return nil, err
		}
		var expirationTime *time.Time

		if expTime == -1 {
			expirationTime = nil
		} else {
			expT := time.Unix(expTime, 0)
			expirationTime = &expT
		}

		egidElement := e.FindElement("EGID")
		if egidElement == nil {
			return nil, &errors.XMLElementError{Path: "egid in login token"}
		}

		egid, err := strconv.Atoi(egidElement.Text())
		if err != nil {
			return nil, err
		}

		loginTokens[i] = LoginToken{Token: tokenElement.Text(), ExpirationTime: expirationTime, EGID: egid}
	}
	return loginTokens, nil
}
