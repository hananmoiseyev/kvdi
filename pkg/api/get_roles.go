package api

import (
	"net/http"

	"github.com/tinyzimmer/kvdi/pkg/auth/types"
	"github.com/tinyzimmer/kvdi/pkg/util/apiutil"
	"github.com/tinyzimmer/kvdi/pkg/util/rethinkdb"
)

// swagger:route GET /api/roles Roles getRoles
// Retrieves a list of the authorization rolse in kVDI.
// responses:
//   200: rolesResponse
//   403: error
//   500: error
func (d *desktopAPI) GetRoles(w http.ResponseWriter, r *http.Request) {
	sess, err := rethinkdb.New(rethinkdb.RDBAddrForCR(d.vdiCluster))
	if err != nil {
		apiutil.ReturnAPIError(err, w)
		return
	}
	defer sess.Close()
	roles, err := sess.GetAllRoles()
	if err != nil {
		apiutil.ReturnAPIError(err, w)
		return
	}
	apiutil.WriteJSON(roles, w)
}

// swagger:operation GET /api/roles/{role} Roles getRole
// ---
// summary: Retrieve the specified role.
// description: Details include the grants, namespaces, and template patterns for the role.
// parameters:
// - name: role
//   in: path
//   description: The role to retrieve details about
//   type: string
//   required: true
// responses:
//   "200":
//     "$ref": "#/responses/roleResponse"
//   "403":
//     "$ref": "#/responses/error"
//   "404":
//     "$ref": "#/responses/error"
//   "500":
//     "$ref": "#/responses/error"
func (d *desktopAPI) GetRole(w http.ResponseWriter, r *http.Request) {
	sess, err := rethinkdb.New(rethinkdb.RDBAddrForCR(d.vdiCluster))
	if err != nil {
		apiutil.ReturnAPIError(err, w)
		return
	}
	defer sess.Close()
	role, err := sess.GetRole(getRoleFromRequest(r))
	if err != nil {
		apiutil.ReturnAPIError(err, w)
		return
	}
	apiutil.WriteJSON(role, w)
}

// A list of roles
// swagger:response rolesResponse
type swaggerRolesResponse struct {
	// in:body
	Body []types.Role
}

// A single role
// swagger:response roleResponse
type swaggerRoleResponse struct {
	// in:body
	Body types.Role
}