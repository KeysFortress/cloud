-- Seed Event types
INSERT INTO public.event_types(
	name)
	VALUES ('Password Created');
INSERT INTO public.event_types(
	name)
	VALUES ('Password Updated');
INSERT INTO public.event_types(
	name)
	VALUES ('Password Deleted');
INSERT INTO public.event_types(
	name)
	VALUES ('Secret Created');
INSERT INTO public.event_types(
	name)
	VALUES ('Secret Updated');
INSERT INTO public.event_types(
	name)
	VALUES ('Secret Deleted');
INSERT INTO public.event_types(
	name)
	VALUES ('Identity Created');
INSERT INTO public.event_types(
	name)
	VALUES ('Identity Updated');
INSERT INTO public.event_types(
	name)
	VALUES ('Identity Deleted');
INSERT INTO public.event_types(
	name)
	VALUES ('Totp Created');
INSERT INTO public.event_types(
	name)
	VALUES ('Totp Updated');
INSERT INTO public.event_types(
	name)
	VALUES ('Totp Deleted');

-- Seed key types
INSERT INTO public.key_types(
	name, algorithm)
	VALUES ('ED25519', 'Edwards-curve Digital Signature Algorithm (EdDSA)');

INSERT INTO public.key_types(
	name, algorithm)
	VALUES ('RSA', 'Rivest–Shamir–Adleman (RSA)');

