package resources

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/beevik/etree"
	"github.com/onego-project/onego/errors"
)

// Resource structure contains XML data and main methods for open nebula resources
type Resource struct {
	XMLData *etree.Element
}

// Permissions structure represents permissions
type Permissions struct {
	User  PermissionGroup
	Group PermissionGroup
	Other PermissionGroup
}

// PermissionGroup structure to create permission structure
type PermissionGroup struct {
	Use    bool
	Manage bool
	Admin  bool
}

const invalidCode = -1

// constants for value conversion in resources methods
const base10 = 10
const bitSize64 = 64

// CreateResource constructs Resource by tag and ID
func CreateResource(tag string, id int) *Resource {
	doc := etree.NewDocument()

	el := doc.CreateElement(tag).CreateElement("ID")
	el.SetText(fmt.Sprintf("%d", id))

	return &Resource{doc.Root()}
}

// Attribute gets resource attribute founded on the path
func (r *Resource) Attribute(path string) (string, error) {
	return attributeFromElement(r.XMLData, path)
}

func (r *Resource) intAttribute(path string) (int, error) {
	return intAttributeFromElement(r.XMLData, path)
}

func attributeFromElement(e *etree.Element, path string) (string, error) {
	if e == nil {
		return "", &errors.XMLElementError{Path: path}
	}

	element := e.FindElement(path)
	if element == nil {
		return "", &errors.XMLElementError{Path: path}
	}
	return element.Text(), nil
}

func intAttributeFromElement(e *etree.Element, path string) (int, error) {
	attribute, err := attributeFromElement(e, path)
	if err != nil {
		return invalidCode, err
	}
	i, err := strconv.Atoi(attribute)
	if err != nil {
		return invalidCode, err
	}
	return i, nil
}

// ID gets resource id
func (r *Resource) ID() (int, error) {
	return r.intAttribute("ID")
}

// Name gets resource name
func (r *Resource) Name() (string, error) {
	return r.Attribute("NAME")
}

func (r *Resource) arrayOfIDs(tag string) ([]int, error) {
	elements := r.XMLData.FindElements(tag + "/ID")
	if len(elements) == 0 {
		return make([]int, 0), nil
	}

	object := make([]int, len(elements))

	for i, e := range elements {
		id, err := strconv.Atoi(e.Text())
		if err != nil {
			return nil, err
		}
		object[i] = id
	}
	return object, nil
}

func intToBool(i int) bool {
	return i == 1
}

func (r *Resource) createPermission(perm string) (*PermissionGroup, error) {
	permissionTypeArray := [3]string{"U", "M", "A"}

	var resPermTypeArray [len(permissionTypeArray)]int
	var err error

	for i, permType := range permissionTypeArray {
		resPermTypeArray[i], err = r.intAttribute("PERMISSIONS/" + perm + "_" + permType)
		if err != nil {
			return nil, err
		}
	}

	return &PermissionGroup{Use: intToBool(resPermTypeArray[0]), Manage: intToBool(resPermTypeArray[1]),
		Admin: intToBool(resPermTypeArray[2])}, nil
}

func (r *Resource) permissions() (*Permissions, error) {
	permissionGroupArray := [3]string{"OWNER", "GROUP", "OTHER"}

	var resPermGroupArray [len(permissionGroupArray)]*PermissionGroup
	var err error

	for i, permGroup := range permissionGroupArray {
		resPermGroupArray[i], err = r.createPermission(permGroup)
		if err != nil {
			return nil, err
		}
	}

	return &Permissions{User: *resPermGroupArray[0], Group: *resPermGroupArray[1], Other: *resPermGroupArray[2]}, nil
}

// parseStringsFromElement creates an array of a same length as an incoming array.
// If the tag (string) from incoming array is not found the function returns nil and error.
// Otherwise returns the array of strings and nil (no error).
func parseStringsFromElement(element *etree.Element, parseStrings []string) ([]string, error) {
	parsedStrings := make([]string, len(parseStrings))
	var err error

	for i, parseString := range parseStrings {
		parsedStrings[i], err = attributeFromElement(element, parseString)
		if err != nil {
			return nil, err
		}
	}

	return parsedStrings, nil
}

// parseStringsFromElementWithoutError creates an array of a same length as an incoming array.
// If the tag (string) from incoming array is found in a given element - the value is added to
// an outgoing array (on a same position as the tag in incoming array).
// If the tag (string) is not found in the given element - the outgoing value (on a same position
// as tag) is default. For string it is empty string ("").
func parseStringsFromElementWithoutError(element *etree.Element, parseStrings []string) []string {
	parsedStrings := make([]string, len(parseStrings))

	for i, parseString := range parseStrings {
		s, err := attributeFromElement(element, parseString)
		if err == nil {
			parsedStrings[i] = s
		}
	}

	return parsedStrings
}

// parseIntsFromElement creates an array of a same length as an incoming array.
// If the tag (string) from incoming array is not found the function returns nil and error.
// Otherwise returns the array of pointers to int and nil (no error).
func parseIntsFromElement(element *etree.Element, parseInts []string) ([]int, error) {
	parsedInts := make([]int, len(parseInts))
	var err error

	for i, parseInt := range parseInts {
		parsedInts[i], err = intAttributeFromElement(element, parseInt)
		if err != nil {
			return nil, err
		}
	}

	return parsedInts, nil
}

// parseIntsFromElementWithoutError creates an array of a same length as an incoming array.
// If the tag (string) from incoming array is found in a given element - the value is added to
// an outgoing array (on a same position as the tag in incoming array).
// If the tag (string) is not found in the given element - the outgoing value (on a same position
// as tag) is default. For int pointer it is nil.
func parseIntsFromElementWithoutError(element *etree.Element, ints []string) []*int {
	parsed := make([]*int, len(ints))

	for i, parseInt := range ints {
		p, err := intAttributeFromElement(element, parseInt)
		if err == nil {
			parsed[i] = &p
		}
	}

	return parsed
}

// parseTimesFromElement creates an array of a same length as an incoming array.
// If the tag (string) from incoming array is not found the function returns nil and error.
// Otherwise returns the array of pointers to time and nil (no error).
func parseTimesFromElement(element *etree.Element, times []string) ([]*time.Time, error) {
	parsed := make([]*time.Time, len(times))

	for i, t := range times {
		num, err := intAttributeFromElement(element, t)
		if err != nil {
			return nil, err
		}

		if num != 0 {
			time := time.Unix(int64(num), 0)
			parsed[i] = &time
		}
	}

	return parsed, nil
}

// parseTimesFromElementWithoutError creates an array of a same length as an incoming array.
// If the tag (string) from incoming array is found in a given element - the value is added to
// an outgoing array (on a same position as the tag in incoming array).
// If the tag (string) is not found in the given element - the outgoing value (on a same position
// as tag) is default. For time pointer it is nil.
func parseTimesFromElementWithoutError(element *etree.Element, times []string) []*time.Time {
	parsed := make([]*time.Time, len(times))

	for i, t := range times {
		num, err := intAttributeFromElement(element, t)
		if err == nil && num != 0 {
			time := time.Unix(int64(num), 0)
			parsed[i] = &time
		}
	}

	return parsed
}

// RenderInterfaceToXMLString renders structures to XML as string.
func RenderInterfaceToXMLString(r interface{}) (string, error) {
	buf := new(bytes.Buffer)

	enc := xml.NewEncoder(buf)

	if err := enc.Encode(r); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (r *Resource) registrationTime() (*time.Time, error) {
	return r.parseTime("REGTIME")
}

func (r *Resource) parseTime(path string) (*time.Time, error) {
	timeInt, err := r.intAttribute(path)
	if err != nil {
		return nil, err
	}

	var t *time.Time
	if timeInt != 0 {
		tt := time.Unix(int64(timeInt), 0)
		t = &tt
	}

	return t, nil
}

func findDiskTypeByValue(value string) (*DiskType, error) {
	for key, val := range DiskTypeMap {
		if val == value {
			return &key, nil
		}
	}
	return nil, fmt.Errorf("unable to find DiskType of value: %s", value)
}

func stringToBool(s string) bool {
	return s == "YES"
}

// parseIntsFromString creates an array of ints with same length as a number of elements separated by comma in string.
// The format of incoming string should be "1,2,3,5"
// If a value in incoming string is not int type the function returns nil and error.
// Otherwise returns the array of ints and nil (no error).
func parseIntsFromString(s string) ([]int, error) {
	str := strings.Split(s, ",")
	ints := make([]int, len(str))
	var err error
	for i, c := range str {
		ints[i], err = strconv.Atoi(strings.TrimSpace(c))
		if err != nil {
			return nil, err
		}
	}
	return ints, err
}
