// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	models "github.com/bruno-nascimento/another-go-challenge/models"
	mock "github.com/stretchr/testify/mock"
)

// TeamService is an autogenerated mock type for the TeamService type
type TeamService struct {
	mock.Mock
}

// FindTeams provides a mock function with given fields:
func (_m *TeamService) FindTeams() []models.Team {
	ret := _m.Called()

	var r0 []models.Team
	if rf, ok := ret.Get(0).(func() []models.Team); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Team)
		}
	}

	return r0
}

type mockConstructorTestingTNewTeamService interface {
	mock.TestingT
	Cleanup(func())
}

// NewTeamService creates a new instance of TeamService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewTeamService(t mockConstructorTestingTNewTeamService) *TeamService {
	mock := &TeamService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
