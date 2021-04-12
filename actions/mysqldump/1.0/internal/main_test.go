package main

import (
	"os"
	"regexp"
	"testing"

	"gopkg.in/stretchr/testify.v1/assert"
)

func init() {
	os.Setenv("UPLOADDIR", "/tmp/sqldump/uploaddir")
	os.Setenv("ACTION_MYSQL_HOST", "localhost")
	os.Setenv("ACTION_MYSQL_PORT", "3306")
	os.Setenv("ACTION_MYSQL_USERNAME", "")
	os.Setenv("ACTION_MYSQL_PASSWORD", "")
	os.Setenv("ACTION_MYSQL_DATABASE", "")
}

func TestDump(t *testing.T) {
	os.Setenv("ACTION_IGNORE_TABLES_REGEXP", `^(fdp|pmp|okr|pipeline|dice_repo_web|dice_notify_histories|dice_nexus_|qa_sonar|qa_test_records|s_|ps_activities|ps_runtime_instances|ci_v3_build|ps_tickets|cm_containers|ps_v2_deployments|uc_user_event_log).*$`)
	main()
}

func TestRegexp(t *testing.T) {
	r, err := regexp.CompilePOSIX(`^(fdp|pmp|okr|pipeline|dice_repo_web|dice_notify_histories|dice_nexus_|qa_sonar|qa_test_records|s_|ps_activities|ps_runtime_instances|ci_v3_build|ps_tickets|cm_containers|ps_v2_deployments|uc_user_event_log).*$`)
	assert.NoError(t, err)
	assert.True(t, r.MatchString("pmp_*"))
	assert.True(t, r.MatchString("dice_notify_histories"))
	assert.True(t, r.MatchString("pipeline_crons"))
	assert.True(t, r.MatchString("s_instances"))
	assert.False(t, r.MatchString("apipelines"))
}
