CREATE TABLE public.users
(
    id integer NOT NULL DEFAULT nextval('users_id_seq'::regclass),
    email character varying COLLATE pg_catalog."default" NOT NULL,
    firstname character varying COLLATE pg_catalog."default",
    lastname character varying COLLATE pg_catalog."default",
    password character varying COLLATE pg_catalog."default",
    CONSTRAINT users_pkey PRIMARY KEY (email)
)
WITH (
    OIDS = FALSE
)
TABLESPACE pg_default;

ALTER TABLE public.users
    OWNER to postgres;