package api

import (
	"net/http"

	"github.com/tinyzimmer/kvdi/pkg/auth/types"
	"github.com/tinyzimmer/kvdi/pkg/util/apiutil"
)

// swagger:operation DELETE /api/roles/{role} Roles deleteRoleRequest
// ---
// summary: Delete the specified role.
// parameters:
// - name: role
//   in: path
//   description: The role to delete
//   type: string
//   required: true
// responses:
//   "200":
//     "$ref": "#/responses/boolResponse"
//   "400":
//     "$ref": "#/responses/error"
//   "403":
//     "$ref": "#/responses/error"
//   "500":
//     "$ref": "#/responses/error"
func (d *desktopAPI) DeleteRole(w http.ResponseWriter, r *http.Request) {
	role := getRoleFromRequest(r)
	sess, err := d.getDB()
	if err != nil {
		apiutil.ReturnAPIError(err, w)
		return
	}
	defer sess.Close()
	if err := sess.DeleteRole(&types.Role{Name: role}); err != nil {
		apiutil.ReturnAPIError(err, w)
		return
	}
	apiutil.WriteOK(w)
}
