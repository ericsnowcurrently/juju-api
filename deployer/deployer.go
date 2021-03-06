// Copyright 2012, 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package deployer

import (
	"github.com/juju/names"

	"github.com/juju/api/base"
	"github.com/juju/api/common"
	"github.com/juju/api/params"
)

const deployerFacade = "Deployer"

// State provides access to the deployer worker's idea of the state.
type State struct {
	caller base.Caller
	*common.APIAddresser
}

// NewState creates a new State instance that makes API calls
// through the given caller.
func NewState(caller base.Caller) *State {
	return &State{
		APIAddresser: common.NewAPIAddresser(deployerFacade, caller),
		caller:       caller,
	}

}

func (st *State) call(method string, params, result interface{}) error {
	return st.caller.Call(deployerFacade, "", method, params, result)
}

// unitLife returns the lifecycle state of the given unit.
func (st *State) unitLife(tag string) (params.Life, error) {
	return common.Life(st.caller, deployerFacade, tag)
}

// Unit returns the unit with the given tag.
func (st *State) Unit(unitTag string) (*Unit, error) {
	life, err := st.unitLife(unitTag)
	if err != nil {
		return nil, err
	}
	tag, err := names.ParseUnitTag(unitTag)
	if err != nil {
		return nil, err
	}
	return &Unit{
		tag:  tag,
		life: life,
		st:   st,
	}, nil
}

// Machine returns the machine with the given tag.
func (st *State) Machine(machineTag string) (*Machine, error) {
	tag, err := names.ParseMachineTag(machineTag)
	if err != nil {
		return nil, err
	}
	return &Machine{
		tag: tag,
		st:  st,
	}, nil
}

// StateAddresses returns the list of addresses used to connect to the state.
func (st *State) StateAddresses() ([]string, error) {
	var result params.StringsResult
	err := st.call("StateAddresses", nil, &result)
	if err != nil {
		return nil, err
	}
	return result.Result, nil
}

// ConnectionInfo returns all the address information that the deployer task
// needs in one call.
func (st *State) ConnectionInfo() (result params.DeployerConnectionValues, err error) {
	err = st.call("ConnectionInfo", nil, &result)
	return result, err
}
