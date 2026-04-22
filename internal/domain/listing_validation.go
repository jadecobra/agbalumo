package domain

import (
	"errors"
	"fmt"
	"time"
)

type validationRule struct {
	condition func(*Listing) bool
	err       string
}

var validationRules = []validationRule{
	{condition: func(l *Listing) bool { return l.City == "" }, err: "city is required"},
}

var lengthRules = []struct {
	field func(*Listing) string
	name  string
	limit int
}{
	{name: "title", field: func(l *Listing) string { return l.Title }, limit: 100},
	{name: "description", field: func(l *Listing) string { return l.Description }, limit: 2000},
	{name: "company name", field: func(l *Listing) string { return l.Company }, limit: 100},
	{name: "address", field: func(l *Listing) string { return l.Address }, limit: 200},
}

var jobFields = []struct {
	field func(*Listing) string
	err   string
}{
	{field: func(l *Listing) string { return l.Company }, err: "company name is required for job listings"},
	{field: func(l *Listing) string { return l.Description }, err: "description is required"},
	{field: func(l *Listing) string { return l.Skills }, err: "skills are required for job listings"},
	{field: func(l *Listing) string { return l.PayRange }, err: "compensation/pay range is required"},
	{field: func(l *Listing) string { return l.JobApplyURL }, err: "apply url is required"},
}

// Validate enforces domain rules for the Listing.
func (l *Listing) Validate() error {
	if err := l.validateOrigin(); err != nil {
		return err
	}
	if err := l.validateTypeRequirements(); err != nil {
		return err
	}
	if err := l.validateContact(); err != nil {
		return err
	}
	if err := l.applyRules(); err != nil {
		return err
	}
	return l.validateTypeSpecific()
}

// applyRules runs the validationRules, lengthRules, and (for Job) jobFields in sequence.
func (l *Listing) applyRules() error {
	for _, rule := range validationRules {
		if rule.condition(l) {
			return errors.New(rule.err)
		}
	}
	for _, rule := range lengthRules {
		if len(rule.field(l)) > rule.limit {
			return fmt.Errorf("%s cannot exceed %d characters", rule.name, rule.limit)
		}
	}
	if l.Type != Job {
		return nil
	}
	for _, f := range jobFields {
		if f.field(l) == "" {
			return errors.New(f.err)
		}
	}
	return nil
}

func (l *Listing) validateTypeRequirements() error {
	if (l.Type == Business || l.Type == Food) && l.Address == "" {
		return errors.New("address is required for business and food listings")
	}
	if l.HoursOfOperation != "" && !(l.Type == Business || l.Type == Service || l.Type == Food) {
		return errors.New("hours of operation not applicable for this listing type")
	}
	return nil
}

func (l *Listing) validateContact() error {
	if l.ContactEmail == "" && l.ContactWhatsApp == "" && l.ContactPhone == "" && l.WebsiteURL == "" {
		return ErrMissingContact
	}
	return nil
}

func (l *Listing) validateTypeSpecific() error {
	switch l.Type {
	case Request:
		return l.validateRequest()
	case Event:
		return l.validateEvent()
	case Job:
		return l.validateJob()
	}
	return nil
}

func (l *Listing) validateRequest() error {
	if !l.Deadline.IsZero() && l.Deadline.Before(time.Now().Add(-24*time.Hour)) {
		return errors.New("deadline cannot be in the past")
	}

	start := l.CreatedAt
	if start.IsZero() {
		start = time.Now()
	}

	if l.Deadline.After(start.Add(90 * 24 * time.Hour)) {
		return ErrInvalidDeadline
	}
	return nil
}

func (l *Listing) validateEvent() error {
	if l.EventStart.IsZero() {
		return errors.New("event start time is required")
	}
	if l.EventEnd.IsZero() {
		return errors.New("event end time is required")
	}
	if l.EventEnd.Before(l.EventStart) {
		return errors.New("event end time cannot be before start time")
	}
	return nil
}

func (l *Listing) validateJob() error {
	if l.JobStartDate.IsZero() {
		return errors.New("job start date is required")
	}
	if l.JobStartDate.Before(time.Now().Add(-24 * time.Hour)) {
		return errors.New("job start date cannot be in the past")
	}
	return nil
}
