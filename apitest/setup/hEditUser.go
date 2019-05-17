package setup

import (
	"github.com/romshark/dgraph_graphql_go/api/graph"
	"github.com/romshark/dgraph_graphql_go/api/graph/gqlmod"
	"github.com/romshark/dgraph_graphql_go/store"
	"github.com/stretchr/testify/require"
)

func (h Helper) editUser(
	assumedSuccess successAssumption,
	userID store.ID,
	editorID store.ID,
	newEmail *string,
	newPassword *string,
) (*gqlmod.User, *graph.ResponseError) {
	t := h.c.t

	var old struct {
		User *gqlmod.User `json:"user"`
	}
	require.NoError(t, h.ts.Debug().QueryVar(
		`query($userID: Identifier!) {
			user(id: $userID) {
				id
				creation
				displayName
				email
			}
		}`,
		map[string]interface{}{
			"userID": string(userID),
		},
		&old,
	))

	var result struct {
		EditUser *gqlmod.User `json:"editUser"`
	}
	err := h.c.QueryVar(
		`mutation (
			$user: Identifier!
			$editor: Identifier!
			$newEmail: String
			$newPassword: String
		) {
			editUser(
				user: $user
				editor: $editor
				newEmail: $newEmail
				newPassword: $newPassword
			) {
				id
				creation
				displayName
				email
			}
		}`,
		map[string]interface{}{
			"user":        string(userID),
			"editor":      string(editorID),
			"newEmail":    newEmail,
			"newPassword": newPassword,
		},
		&result,
	)

	if err := checkErr(t, assumedSuccess, err); err != nil {
		return nil, err
	}

	require.NotNil(t, result.EditUser)
	if old.User != nil {
		require.Equal(t, *old.User.ID, *result.EditUser.ID)
		if newEmail != nil {
			require.Equal(t, *newEmail, *result.EditUser.Email)
		} else {
			require.Equal(t, *old.User.Email, *result.EditUser.Email)
		}
		require.Equal(t, *old.User.DisplayName, *result.EditUser.DisplayName)
		require.Equal(t, *old.User.Creation, *result.EditUser.Creation)
	}

	return result.EditUser, nil
}

// EditUser helps editing a user
func (h Helper) EditUser(
	userID store.ID,
	editorID store.ID,
	newEmail *string,
	newPassword *string,
) (*gqlmod.User, *graph.ResponseError) {
	return h.editUser(
		potentialFailure,
		userID,
		editorID,
		newEmail,
		newPassword,
	)
}

// EditUser helps editing a user and assumes success
func (ok AssumeSuccess) EditUser(
	userID store.ID,
	editorID store.ID,
	newEmail *string,
	newPassword *string,
) *gqlmod.User {
	result, _ := ok.h.editUser(
		success,
		userID,
		editorID,
		newEmail,
		newPassword,
	)
	return result
}
