BEGIN;


CREATE TABLE IF NOT EXISTS public.account_identties
(
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    name character varying(100) COLLATE pg_catalog."default" NOT NULL,
    "publicKey" character varying COLLATE pg_catalog."default",
    "privateKey" character varying COLLATE pg_catalog."default",
    key_type_id integer,
    account_id uuid,
    created_at timestamp without time zone,
    updated_at timestamp without time zone,
    CONSTRAINT account_identties_pkey PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS public.account_passwords
(
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    name character varying(250) COLLATE pg_catalog."default" NOT NULL,
    content character varying COLLATE pg_catalog."default" NOT NULL,
    account_id uuid NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone,
    CONSTRAINT account_passwords_pkey PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS public.account_secrets
(
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    name character varying(100) COLLATE pg_catalog."default" NOT NULL,
    description character varying(300) COLLATE pg_catalog."default",
    content character varying COLLATE pg_catalog."default" NOT NULL,
    account_id uuid NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone,
    CONSTRAINT account_secrets_pkey PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS public.accounts
(
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    email character varying(200) COLLATE pg_catalog."default" NOT NULL,
    phone character varying(16) COLLATE pg_catalog."default" NOT NULL,
    subscription_id uuid NOT NULL,
    account_type_id uuid NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone,
    enabled boolean NOT NULL DEFAULT true,
    CONSTRAINT accounts_pkey PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS public.associated_account_access_keys
(
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    account_id uuid NOT NULL,
    access_key character varying(250) COLLATE pg_catalog."default" NOT NULL,
    created_at timestamp without time zone NOT NULL,
    CONSTRAINT associated_account_access_keys_pkey PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS public.key_types
(
    id integer NOT NULL,
    CONSTRAINT key_types_pkey PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS public.rac_devices
(
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    name character varying(100) COLLATE pg_catalog."default" NOT NULL,
    address character varying(40) COLLATE pg_catalog."default" NOT NULL,
    key character varying COLLATE pg_catalog."default" NOT NULL,
    created_by uuid NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone,
    updated_by uuid,
    CONSTRAINT rac_devices_pkey PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS public.subscriptions
(
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    created_at timestamp without time zone NOT NULL,
    next_billing_date timestamp without time zone NOT NULL,
    enabled boolean NOT NULL,
    CONSTRAINT subscriptions_pkey PRIMARY KEY (id)
);

ALTER TABLE IF EXISTS public.account_identties
    ADD CONSTRAINT "account_Identities_account_id" FOREIGN KEY (account_id)
    REFERENCES public.accounts (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION;


ALTER TABLE IF EXISTS public.account_identties
    ADD CONSTRAINT account_identities_key_type_id FOREIGN KEY (key_type_id)
    REFERENCES public.key_types (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID;


ALTER TABLE IF EXISTS public.account_passwords
    ADD CONSTRAINT account_passwords_account_id FOREIGN KEY (account_id)
    REFERENCES public.accounts (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION;


ALTER TABLE IF EXISTS public.account_secrets
    ADD CONSTRAINT account_secrets_account_id FOREIGN KEY (account_id)
    REFERENCES public.accounts (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION;


ALTER TABLE IF EXISTS public.accounts
    ADD CONSTRAINT accounts_subscription_id FOREIGN KEY (subscription_id)
    REFERENCES public.subscriptions (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID;


ALTER TABLE IF EXISTS public.associated_account_access_keys
    ADD CONSTRAINT associated_account_access_keys_account_id FOREIGN KEY (account_id)
    REFERENCES public.accounts (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION;


ALTER TABLE IF EXISTS public.rac_devices
    ADD CONSTRAINT rac_devices_created_by FOREIGN KEY (created_by)
    REFERENCES public.accounts (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION;


ALTER TABLE IF EXISTS public.rac_devices
    ADD CONSTRAINT rac_updated_by FOREIGN KEY (updated_by)
    REFERENCES public.accounts (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION;

END;
