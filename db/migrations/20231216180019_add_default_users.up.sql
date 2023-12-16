INSERT INTO public.users (id, username, password, is_active, role, created_at, updated_at, created_by, updated_by, deleted_at)
VALUES
    (1, 'admin', '$2a$10$nH1zPLd.h1kjDN0wJbbhR.Gbqew.BUCoRpxhipb1hCTwjRPbGAMRS', true, 'admin', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'admin', 'admin', NULL),
    (2, 'owner', '$2a$10$nH1zPLd.h1kjDN0wJbbhR.Gbqew.BUCoRpxhipb1hCTwjRPbGAMRS', true, 'owner', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'admin', 'admin', NULL),
    (3, 'employee', '$2a$10$nH1zPLd.h1kjDN0wJbbhR.Gbqew.BUCoRpxhipb1hCTwjRPbGAMRS', true, 'employee', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'admin', 'admin', NULL),
    (4, 'developer', '$2a$10$nH1zPLd.h1kjDN0wJbbhR.Gbqew.BUCoRpxhipb1hCTwjRPbGAMRS', true, 'developer', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 'admin', 'admin', NULL);