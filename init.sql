-- Seed Event types
INSERT INTO public.event_types(
	id,name)
	VALUES (1,'Password Created');
INSERT INTO public.event_types(
	id,name)
	VALUES (2,'Password Updated');
INSERT INTO public.event_types(
	id,name)
	VALUES (3,'Password Deleted');
INSERT INTO public.event_types(
	id, name)
	VALUES (4,'Secret Created');
INSERT INTO public.event_types(
	id,name)
	VALUES (5,'Secret Updated');
INSERT INTO public.event_types(
	id,name)
	VALUES (6, 'Secret Deleted');
INSERT INTO public.event_types(
	id,name)
	VALUES (7,'Identity Created');
INSERT INTO public.event_types(
	id,name)
	VALUES (8,'Identity Updated');
INSERT INTO public.event_types(
	id,name)
	VALUES (9,'Identity Deleted');
INSERT INTO public.event_types(
	id,name)
	VALUES (10,'Totp Created');
INSERT INTO public.event_types(
	id,name)
	VALUES (11,'Totp Updated');
INSERT INTO public.event_types(
	id,name)
	VALUES (12,'Totp Deleted');

-- Seed key types

INSERT INTO public.key_types(id, name, description, has_size)
	VALUES (1, 'ED25519', 'Edwards-curve Digital Signature Algorithm (EdDSA)', false);

INSERT INTO public.key_types(
	id, name, description, has_size)
	VALUES (2, 'RSA',  'Rivest–Shamir–Adleman (RSA)', true);

-- Seed Time Based Code Types
INSERT INTO public.time_based_code_types(
	id, name)
	VALUES (1, 'HOTP');
INSERT INTO public.time_based_code_types(
	id, name)
	VALUES (2, 'TOTP');
INSERT INTO public.time_based_code_types(
	id, name)
	VALUES (3, 'OTP');

-- Seed the algorithm table
INSERT INTO public.time_based_algorithms(
	id, name)
	VALUES (1, 'AlgorithmSHA1');
INSERT INTO public.time_based_algorithms(
	id, name)
	VALUES (2, 'AlgorithmSHA256');
INSERT INTO public.time_based_algorithms(
	id, name)
	VALUES (3, 'AlgorithmSHA512');
INSERT INTO public.time_based_algorithms(
	id, name)
	VALUES (4, 'AlgorithmMD5');
