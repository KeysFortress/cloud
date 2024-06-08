-- Seed Event types
INSERT INTO public.event_types(id, name)
    VALUES (1, 'Password Created')
    ON CONFLICT (id) DO NOTHING;
INSERT INTO public.event_types(id, name)
    VALUES (2, 'Password Updated')
    ON CONFLICT (id) DO NOTHING;
INSERT INTO public.event_types(id, name)
    VALUES (3, 'Password Deleted')
    ON CONFLICT (id) DO NOTHING;
INSERT INTO public.event_types(id, name)
    VALUES (4, 'Secret Created')
    ON CONFLICT (id) DO NOTHING;
INSERT INTO public.event_types(id, name)
    VALUES (5, 'Secret Updated')
    ON CONFLICT (id) DO NOTHING;
INSERT INTO public.event_types(id, name)
    VALUES (6, 'Secret Deleted')
    ON CONFLICT (id) DO NOTHING;
INSERT INTO public.event_types(id, name)
    VALUES (7, 'Identity Created')
    ON CONFLICT (id) DO NOTHING;
INSERT INTO public.event_types(id, name)
    VALUES (8, 'Identity Updated')
    ON CONFLICT (id) DO NOTHING;
INSERT INTO public.event_types(id, name)
    VALUES (9, 'Identity Deleted')
    ON CONFLICT (id) DO NOTHING;
INSERT INTO public.event_types(id, name)
    VALUES (10, 'Totp Created')
    ON CONFLICT (id) DO NOTHING;
INSERT INTO public.event_types(id, name)
    VALUES (11, 'Totp Updated')
    ON CONFLICT (id) DO NOTHING;
INSERT INTO public.event_types(id, name)
    VALUES (12, 'Totp Deleted')
    ON CONFLICT (id) DO NOTHING;

-- Seed key types
INSERT INTO public.key_types(id, name, description, has_size)
    VALUES (1, 'ED25519', 'Edwards-curve Digital Signature Algorithm (EdDSA)', false)
    ON CONFLICT (id) DO NOTHING;

INSERT INTO public.key_types(id, name, description, has_size)
    VALUES (2, 'RSA', 'Rivest–Shamir–Adleman (RSA)', true)
    ON CONFLICT (id) DO NOTHING;

-- Seed Time Based Code Types
INSERT INTO public.time_based_code_types(id, name)
    VALUES (1, 'HOTP')
    ON CONFLICT (id) DO NOTHING;
INSERT INTO public.time_based_code_types(id, name)
    VALUES (2, 'TOTP')
    ON CONFLICT (id) DO NOTHING;
INSERT INTO public.time_based_code_types(id, name)
    VALUES (3, 'OTP')
    ON CONFLICT (id) DO NOTHING;

-- Seed the algorithm table
INSERT INTO public.time_based_algorithms(id, name)
    VALUES (1, 'AlgorithmSHA1')
    ON CONFLICT (id) DO NOTHING;
INSERT INTO public.time_based_algorithms(id, name)
    VALUES (2, 'AlgorithmSHA256')
    ON CONFLICT (id) DO NOTHING;
INSERT INTO public.time_based_algorithms(id, name)
    VALUES (3, 'AlgorithmSHA512')
    ON CONFLICT (id) DO NOTHING;
INSERT INTO public.time_based_algorithms(id, name)
    VALUES (4, 'AlgorithmMD5')
    ON CONFLICT (id) DO NOTHING;

-- Seed the mfa method types
INSERT INTO public.mfa_method_types(id, name)
    VALUES (1, 'email')
    ON CONFLICT (id) DO NOTHING;
INSERT INTO public.mfa_method_types(id, name)
    VALUES (2, 'authenticator')
    ON CONFLICT (id) DO NOTHING;

INSERT INTO public.device_types(
	id, name)
	VALUES (1, 'Android');
INSERT INTO public.device_types(
	id, name)
	VALUES (2, 'IOS');
INSERT INTO public.device_types(
	id, name)
	VALUES (3, 'Windows');
INSERT INTO public.device_types(
	id, name)
	VALUES (4, 'Linux');
INSERT INTO public.device_types(
	id, name)
	VALUES (5, 'MacOS');
