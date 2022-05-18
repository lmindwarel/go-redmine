package redmine

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type timeEntriesResult struct {
	TimeEntries []TimeEntry `json:"time_entries"`
}

type timeEntryResult struct {
	TimeEntry TimeEntry `json:"time_entry"`
}

type timeEntryRequest struct {
	TimeEntry TimeEntry `json:"time_entry"`
}

type TimeEntry struct {
	Id         int     `json:"id"`
	Project    IdName  `json:"project"`
	Issue      Id      `json:"issue"`
	User       IdName  `json:"user"`
	Activity   IdName  `json:"activity"`
	Comments   string  `json:"comments"`
	Hours      float32 `json:"hours"`
	SpentOn    string  `json:"spent_on"`
	CreatedOn  string  `json:"created_on"`
	UpdatedOn  string  `json:"updated_on"`
	ActivityId int     `json:"activity_id,omitempty"`
	ProjectId  int     `json:"project_id,omitempty"`
}

func (c *Client) TimeEntries(projectId int, userId int, spentOn string) ([]TimeEntry, error) {
	res, err := c.Get(c.endpoint + "/projects/" + strconv.Itoa(projectId) + "/time_entries.json?user_id=" + strconv.Itoa(userId) + "&spent_on=" + spentOn + "&key=" + c.apikey + c.getPaginationClause())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var r timeEntriesResult
	if res.StatusCode == 404 {
		return nil, errors.New("Not Found")
	}
	if res.StatusCode != 200 {
		var er ErrorsResult
		err = decoder.Decode(&er)
		if err == nil {
			err = errors.New(strings.Join(er.Errors, "\n"))
		}
	} else {
		err = decoder.Decode(&r)
	}
	if err != nil {
		return nil, err
	}
	return r.TimeEntries, nil
}

func (c *Client) TimeEntry(id int) (*TimeEntry, error) {
	res, err := c.Get(c.endpoint + "/time_entries/" + strconv.Itoa(id) + ".json?key=" + c.apikey)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var r timeEntryResult
	if res.StatusCode == 404 {
		return nil, errors.New("Not Found")
	}
	if res.StatusCode != 200 {
		var er ErrorsResult
		err = decoder.Decode(&er)
		if err == nil {
			err = errors.New(strings.Join(er.Errors, "\n"))
		}
	} else {
		err = decoder.Decode(&r)
	}
	if err != nil {
		return nil, err
	}
	return &r.TimeEntry, nil
}

func (c *Client) CreateTimeEntry(timeEntry TimeEntry) (*TimeEntry, error) {
	var ir timeEntryRequest
	ir.TimeEntry = timeEntry
	s, err := json.Marshal(ir)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", c.endpoint+"/time_entries.json?key="+c.apikey, strings.NewReader(string(s)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var r timeEntryResult
	if res.StatusCode != 201 {
		var er ErrorsResult
		err = decoder.Decode(&er)
		if err == nil {
			err = errors.New(strings.Join(er.Errors, "\n"))
		}
	} else {
		err = decoder.Decode(&r)
	}
	if err != nil {
		return nil, err
	}
	return &r.TimeEntry, nil
}

func (c *Client) UpdateTimeEntry(timeEntry TimeEntry) error {
	var ir timeEntryRequest
	ir.TimeEntry = timeEntry
	s, err := json.Marshal(ir)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("PUT", c.endpoint+"/time_entries/"+strconv.Itoa(timeEntry.Id)+".json?key="+c.apikey, strings.NewReader(string(s)))
	if err != nil {
		return errors.Wrap(err, "failed to request redmine")
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := c.Do(req)
	if res.StatusCode == 404 {
		return errors.New("Not Found")
	}
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		decoder := json.NewDecoder(res.Body)
		var er ErrorsResult
		err = decoder.Decode(&er)
		if err == nil {
			err = errors.New(strings.Join(er.Errors, "\n"))
		}
	}
	if err != nil {
		return err
	}
	return err
}

func (c *Client) DeleteTimeEntry(id int) error {
	req, err := http.NewRequest(http.MethodDelete, c.endpoint+"/time_entries/"+strconv.Itoa(id)+".json?key="+c.apikey, strings.NewReader(""))
	if err != nil {
		return err
	}

	res, err := c.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to request redmine")
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return errors.New("Not Found")
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		decoder := json.NewDecoder(res.Body)
		var er ErrorsResult
		err = decoder.Decode(&er)
		if err == nil {
			return errors.New(strings.Join(er.Errors, "\n"))
		}

		return fmt.Errorf("bad api response status: %d", res.StatusCode)
	}

	return err
}
