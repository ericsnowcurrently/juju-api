// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package uniter

import (
	"fmt"

	"github.com/juju/charm"
	"github.com/juju/names"

	"github.com/juju/api/common"
	"github.com/juju/api/params"
	"github.com/juju/api/watcher"
)

// This module implements a subset of the interface provided by
// state.Service, as needed by the uniter API.

// Service represents the state of a service.
type Service struct {
	st   *State
	tag  names.Tag
	life params.Life
}

// Name returns the service name.
func (s *Service) Name() string {
	return s.tag.Id()
}

// String returns the service as a string.
func (s *Service) String() string {
	return s.Name()
}

// Watch returns a watcher for observing changes to a service.
func (s *Service) Watch() (watcher.NotifyWatcher, error) {
	return common.Watch(s.st.caller, uniterFacade, s.tag.String())
}

// WatchRelations returns a StringsWatcher that notifies of changes to
// the lifecycles of relations involving s.
func (s *Service) WatchRelations() (watcher.StringsWatcher, error) {
	var results params.StringsWatchResults
	args := params.Entities{
		Entities: []params.Entity{{Tag: s.tag.String()}},
	}
	err := s.st.call("WatchServiceRelations", args, &results)
	if err != nil {
		return nil, err
	}
	if len(results.Results) != 1 {
		return nil, fmt.Errorf("expected 1 result, got %d", len(results.Results))
	}
	result := results.Results[0]
	if result.Error != nil {
		return nil, result.Error
	}
	w := watcher.NewStringsWatcher(s.st.caller, result)
	return w, nil
}

// Life returns the service's current life state.
func (s *Service) Life() params.Life {
	return s.life
}

// Refresh refreshes the contents of the Service from the underlying
// state.
func (s *Service) Refresh() error {
	life, err := s.st.life(s.tag.String())
	if err != nil {
		return err
	}
	s.life = life
	return nil
}

// CharmURL returns the service's charm URL, and whether units should
// upgrade to the charm with that URL even if they are in an error
// state (force flag).
//
// NOTE: This differs from state.Service.CharmURL() by returning
// an error instead as well, because it needs to make an API call.
func (s *Service) CharmURL() (*charm.URL, bool, error) {
	var results params.StringBoolResults
	args := params.Entities{
		Entities: []params.Entity{{Tag: s.tag.String()}},
	}
	err := s.st.call("CharmURL", args, &results)
	if err != nil {
		return nil, false, err
	}
	if len(results.Results) != 1 {
		return nil, false, fmt.Errorf("expected 1 result, got %d", len(results.Results))
	}
	result := results.Results[0]
	if result.Error != nil {
		return nil, false, result.Error
	}
	if result.Result != "" {
		curl, err := charm.ParseURL(result.Result)
		if err != nil {
			return nil, false, err
		}
		return curl, result.Ok, nil
	}
	return nil, false, fmt.Errorf("%q has no charm url set", s.tag)
}

// TODO(dimitern) bug #1270795 2014-01-20
// Add a doc comment here.
func (s *Service) GetOwnerTag() (string, error) {
	var result params.StringResult
	args := params.Entities{
		Entities: []params.Entity{{Tag: s.tag.String()}},
	}
	err := s.st.call("GetOwnerTag", args, &result)
	if err != nil {
		return "", err
	}
	if result.Error != nil {
		return "", result.Error
	}
	return result.Result, nil
}
