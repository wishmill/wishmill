-- FUNCTIONS

CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- USERS
CREATE TABLE users (
    id integer NOT NULL,
    name character varying(128) NOT NULL,
    email character varying(128) NOT NULL,
    auth_provider character varying(256) NOT NULL,
    sub character varying(256),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE users ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);

ALTER TABLE ONLY users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

-- CONFIG
CREATE TABLE config (
    id integer NOT NULL,
    key character varying(64) NOT NULL,
    value character varying(512) NOT NULL
);

ALTER TABLE config ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME config_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);

ALTER TABLE ONLY config
    ADD CONSTRAINT config_pkey PRIMARY KEY (id);


-- WISHLIST
CREATE TABLE wishlists (
    id integer NOT NULL,
    name character varying(256),
    surprise_me boolean,
    allow_third_party_add boolean,
    user_id integer NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE wishlists ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME wishlists_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);

ALTER TABLE ONLY wishlists
    ADD CONSTRAINT wishlists_pkey PRIMARY KEY (id);

ALTER TABLE ONLY wishlists
    ADD CONSTRAINT users_fk FOREIGN KEY (user_id) REFERENCES users(id);

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON wishlists
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

-- WISHLISTITEM
CREATE TABLE wishlistitems (
    id integer NOT NULL,
    name character varying(256),
    url character varying(1024),
    priority character varying(32),
    wishlist_id integer NOT NULL,
    user_id integer NOT NULL,
    bought_by integer,
    bought_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE wishlistitems ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME wishlistitems_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);

ALTER TABLE ONLY wishlistitems
    ADD CONSTRAINT wishlistitems_pkey PRIMARY KEY (id);

ALTER TABLE ONLY wishlistitems
    ADD CONSTRAINT wishlists_fk FOREIGN KEY (wishlist_id) REFERENCES wishlists(id);

ALTER TABLE ONLY wishlistitems
    ADD CONSTRAINT users_fk FOREIGN KEY (user_id) REFERENCES users(id);

ALTER TABLE ONLY wishlistitems
    ADD CONSTRAINT users_bought_by_fk FOREIGN KEY (bought_by) REFERENCES users(id);

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON wishlistitems
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();