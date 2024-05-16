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
INSERT INTO public.key_types(
	name, algorithm)
	VALUES ('ED25519', 'Edwards-curve Digital Signature Algorithm (EdDSA)');

INSERT INTO public.key_types(
	name, algorithm)
	VALUES ('RSA', 'Rivest–Shamir–Adleman (RSA)');

-- Seed Time Based Code Types
INSERT INTO public.time_based_code_types(
	id, name)
	VALUES (1, 'HTOP');
INSERT INTO public.time_based_code_types(
	id, name)
	VALUES (2, 'TOTP');
INSERT INTO public.time_based_code_types(
	id, name)
	VALUES (3, 'OTP');
