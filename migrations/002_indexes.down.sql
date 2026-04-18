-- INDEXES DOWN

DROP INDEX IF EXISTS idx_users_group_id;

DROP INDEX IF EXISTS idx_test_artifacts_file_type;
DROP INDEX IF EXISTS idx_test_artifacts_status_id;
DROP INDEX IF EXISTS idx_test_artifacts_run_id;

DROP INDEX IF EXISTS idx_test_runs_status_id;
DROP INDEX IF EXISTS idx_test_runs_suite_id;
DROP INDEX IF EXISTS idx_test_runs_build_id;
DROP INDEX IF EXISTS idx_test_runs_started_at;

DROP INDEX IF EXISTS idx_builds_component_id;
