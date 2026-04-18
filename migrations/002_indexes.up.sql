-- INDEXES UP

-- BUILDS
CREATE INDEX idx_builds_component_id ON builds(component_id);

-- TEST RUNS
CREATE INDEX idx_test_runs_build_id ON test_runs(build_id);
CREATE INDEX idx_test_runs_suite_id ON test_runs(suite_id);
CREATE INDEX idx_test_runs_status_id ON test_runs(status_id);
CREATE INDEX idx_test_runs_started_at ON test_runs(started_at);

-- TEST ARTIFACTS
CREATE INDEX idx_test_artifacts_run_id ON test_artifacts(run_id);
CREATE INDEX idx_test_artifacts_status_id ON test_artifacts(status_id);
CREATE INDEX idx_test_artifacts_file_type ON test_artifacts(file_type_id);

-- USERS
CREATE INDEX idx_users_group_id ON users(group_id);
