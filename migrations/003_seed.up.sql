-- SEED UP

INSERT INTO file_types (id, name) VALUES (1, 'Log') ON CONFLICT (id) DO NOTHING;
INSERT INTO file_types (id, name) VALUES (2, 'Image') ON CONFLICT (id) DO NOTHING;
INSERT INTO file_types (id, name) VALUES (3, 'Video') ON CONFLICT (id) DO NOTHING;
INSERT INTO file_types (id, name) VALUES (4, 'PDF') ON CONFLICT (id) DO NOTHING;

INSERT INTO run_statuses (id, name) VALUES (1, 'Running') ON CONFLICT (id) DO NOTHING;
INSERT INTO run_statuses (id, name) VALUES (2, 'Finished') ON CONFLICT (id) DO NOTHING;
INSERT INTO run_statuses (id, name) VALUES (3, 'Aborted') ON CONFLICT (id) DO NOTHING;

INSERT INTO result_statuses (id, name) VALUES (1, 'Passed') ON CONFLICT (id) DO NOTHING;
INSERT INTO result_statuses (id, name) VALUES (2, 'Failed') ON CONFLICT (id) DO NOTHING;
INSERT INTO result_statuses (id, name) VALUES (3, 'Skipped') ON CONFLICT (id) DO NOTHING;
INSERT INTO result_statuses (id, name) VALUES (4, 'Error') ON CONFLICT (id) DO NOTHING;
