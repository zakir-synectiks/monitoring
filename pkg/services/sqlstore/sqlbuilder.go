package sqlstore

import (
	"bytes"
	"strings"

	m "github.com/xformation/synectiks-monitoring/pkg/models"
)

type SqlBuilder struct {
	sql    bytes.Buffer
	params []interface{}
}

func (sb *SqlBuilder) Write(sql string, params ...interface{}) {
	sb.sql.WriteString(sql)

	if len(params) > 0 {
		sb.params = append(sb.params, params...)
	}
}

func (sb *SqlBuilder) GetSqlString() string {
	return sb.sql.String()
}

func (sb *SqlBuilder) AddParams(params ...interface{}) {
	sb.params = append(sb.params, params...)
}

func (sb *SqlBuilder) writeDashboardPermissionFilter(user *m.SignedInUser, permission m.PermissionType) {

	if user.OrgRole == m.ROLE_ADMIN {
		return
	}

	okRoles := []interface{}{user.OrgRole}

	if user.OrgRole == m.ROLE_EDITOR {
		okRoles = append(okRoles, m.ROLE_VIEWER)
	}

	falseStr := dialect.BooleanStr(false)

	sb.sql.WriteString(` AND
	(
		dashboard.id IN (
			SELECT distinct d.id AS DashboardId
			FROM dashboard AS d
			 	LEFT JOIN dashboard folder on folder.id = d.folder_id
			    LEFT JOIN dashboard_acl AS da ON
	 			da.dashboard_id = d.id OR
	 			da.dashboard_id = d.folder_id OR
	 			(
	 				-- include default permissions -->
					da.org_id = -1 AND (
					  (folder.id IS NOT NULL AND folder.has_acl = ` + falseStr + `) OR
					  (folder.id IS NULL AND d.has_acl = ` + falseStr + `)
					)
	 			)
				LEFT JOIN team_member as ugm on ugm.team_id = da.team_id
			WHERE
				d.org_id = ? AND
				da.permission >= ? AND
				(
					da.user_id = ? OR
					ugm.user_id = ? OR
					da.role IN (?` + strings.Repeat(",?", len(okRoles)-1) + `)
				)
		)
	)`)

	sb.params = append(sb.params, user.OrgId, permission, user.UserId, user.UserId)
	sb.params = append(sb.params, okRoles...)
}
