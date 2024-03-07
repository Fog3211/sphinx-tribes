package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocking the database structure for testing
type MockDB struct {
	mock.Mock
	database
}

// // Implementing mock methods for GetUserRoles and GetOrganizationByUuid
// func (m *MockDB) GetUserRoles(uuid string, pubKeyFromAuth string) []UserRoles {
// 	args := m.Called(uuid, pubKeyFromAuth)
// 	return args.Get(0).([]UserRoles)
// }

// func (m *MockDB) GetOrganizationByUuid(uuid string) Organization {
// 	args := m.Called(uuid)
// 	return args.Get(0).(Organization)
// }

func TestRolesCheck(t *testing.T) {
	userRoles := []UserRoles{
		{Role: "ADD BOUNTY"},
	}

	// test returns true when a user has a role
	assert.True(t, RolesCheck(userRoles, "ADD BOUNTY"))

	// returns false when a use does not have a role
	assert.False(t, RolesCheck(userRoles, "DELETE BOUNTY"))
	assert.False(t, RolesCheck(userRoles, "DELETE BOUNTY2"))
}

func TestCheckUser(t *testing.T) {
	userRoles := []UserRoles{
		{OwnerPubKey: "userPublicKey"},
	}
	// if in the user roles, one of the owner_pubkey belongs to the user return true else return false
	assert.True(t, CheckUser(userRoles, "userPublicKey"))
	assert.False(t, CheckUser(userRoles, "anotherPublicKey"))
}

func TestUserHasAccess(t *testing.T) {
	mockDB := new(MockDB)

	org := Organization{OwnerPubKey: "adminKey"}
	mockDB.On("GetOrganizationByUuid", "orgUUID").Return(org)

	// Test for admin access
	assert.True(t, DB.UserHasAccess("adminKey", "orgUUID", "ANY ROLE"))

	// Test for non-admin with role access
	userRoles := []UserRoles{{Role: "VIEW REPORT"}}
	mockDB.On("GetUserRoles", "orgUUID", "userKey").Return(userRoles)
	assert.True(t, DB.UserHasAccess("userKey", "orgUUID", "VIEW REPORT"))

	// Test for non-admin without role access
	assert.False(t, DB.UserHasAccess("userKey", "orgUUID", "ADD BOUNTY"))
}

func TestUserHasManageBountyRoles(t *testing.T) {
	mockDB := new(MockDB)

	org := Organization{OwnerPubKey: "adminKey"}
	mockDB.On("GetOrganizationByUuid", "orgUUID").Return(org)

	// Test for organization admin
	assert.True(t, DB.UserHasManageBountyRoles("adminKey", "orgUUID"))

	// Test for user with all manage bounty roles
	userRoles := make([]UserRoles, len(ManageBountiesGroup))
	for i, role := range ManageBountiesGroup {
		userRoles[i] = UserRoles{Role: role}
	}
	mockDB.On("GetUserRoles", "orgUUID", "userKey").Return(userRoles)
	assert.True(t, DB.UserHasManageBountyRoles("userKey", "orgUUID"))

	// Test for user missing some manage bounty roles
	partialRoles := userRoles[:len(userRoles)-1]
	mockDB.On("GetUserRoles", "orgUUID", "userKey").Return(partialRoles)
	assert.False(t, DB.UserHasManageBountyRoles("userKey", "orgUUID"))
}
