-- SEED DOWN

TRUNCATE TABLE
    test_artifacts,
    test_runs,
    file_types,
    result_statuses,
    run_statuses,
    test_suites,
    builds,
    components,
    user_groups,
    users
CASCADE;
