--
-- PostgreSQL database dump
--

\restrict ZLPAwQe9s3ipeSYK6y9cCU1vzUKsCkVeI4whYRFOzd1WN9h2mCgy8oPVs7OCd98

-- Dumped from database version 15.18 (Debian 15.18-1.pgdg13+1)
-- Dumped by pg_dump version 15.18 (Debian 15.18-1.pgdg13+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: abilities; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.abilities (
    "group" character varying(64) NOT NULL,
    model character varying(255) NOT NULL,
    channel_id bigint NOT NULL,
    enabled boolean,
    priority bigint DEFAULT 0,
    weight bigint DEFAULT 0,
    tag text
);


ALTER TABLE public.abilities OWNER TO root;

--
-- Name: channels; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.channels (
    id bigint NOT NULL,
    type bigint DEFAULT 0,
    key text NOT NULL,
    open_ai_organization text,
    test_model text,
    status bigint DEFAULT 1,
    name text,
    weight bigint DEFAULT 0,
    created_time bigint,
    test_time bigint,
    response_time bigint,
    base_url text DEFAULT ''::text,
    other text,
    balance numeric,
    balance_updated_time bigint,
    models text,
    "group" character varying(64) DEFAULT 'default'::character varying,
    used_quota bigint DEFAULT 0,
    model_mapping text,
    status_code_mapping character varying(1024) DEFAULT ''::character varying,
    priority bigint DEFAULT 0,
    auto_ban bigint DEFAULT 1,
    other_info text,
    tag text,
    setting text,
    param_override text,
    header_override text,
    remark character varying(255),
    channel_info json,
    settings text
);


ALTER TABLE public.channels OWNER TO root;

--
-- Name: channels_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.channels_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.channels_id_seq OWNER TO root;

--
-- Name: channels_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.channels_id_seq OWNED BY public.channels.id;


--
-- Name: checkins; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.checkins (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    checkin_date character varying(10) NOT NULL,
    quota_awarded bigint NOT NULL,
    created_at bigint
);


ALTER TABLE public.checkins OWNER TO root;

--
-- Name: checkins_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.checkins_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.checkins_id_seq OWNER TO root;

--
-- Name: checkins_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.checkins_id_seq OWNED BY public.checkins.id;


--
-- Name: custom_oauth_providers; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.custom_oauth_providers (
    id bigint NOT NULL,
    name character varying(64) NOT NULL,
    slug character varying(64) NOT NULL,
    icon character varying(128) DEFAULT ''::character varying,
    enabled boolean DEFAULT false,
    client_id character varying(256),
    client_secret character varying(512),
    authorization_endpoint character varying(512),
    token_endpoint character varying(512),
    user_info_endpoint character varying(512),
    scopes character varying(256) DEFAULT 'openid profile email'::character varying,
    user_id_field character varying(128) DEFAULT 'sub'::character varying,
    username_field character varying(128) DEFAULT 'preferred_username'::character varying,
    display_name_field character varying(128) DEFAULT 'name'::character varying,
    email_field character varying(128) DEFAULT 'email'::character varying,
    well_known character varying(512),
    auth_style bigint DEFAULT 0,
    access_policy text,
    access_denied_message character varying(512),
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.custom_oauth_providers OWNER TO root;

--
-- Name: custom_oauth_providers_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.custom_oauth_providers_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.custom_oauth_providers_id_seq OWNER TO root;

--
-- Name: custom_oauth_providers_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.custom_oauth_providers_id_seq OWNED BY public.custom_oauth_providers.id;


--
-- Name: logs; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.logs (
    id bigint NOT NULL,
    user_id bigint,
    created_at bigint,
    type bigint,
    content text,
    username text DEFAULT ''::text,
    token_name text DEFAULT ''::text,
    model_name text DEFAULT ''::text,
    quota bigint DEFAULT 0,
    prompt_tokens bigint DEFAULT 0,
    completion_tokens bigint DEFAULT 0,
    use_time bigint DEFAULT 0,
    is_stream boolean,
    channel_id bigint,
    channel_name text,
    token_id bigint DEFAULT 0,
    "group" text,
    ip text DEFAULT ''::text,
    request_id character varying(64) DEFAULT ''::character varying,
    upstream_request_id character varying(128) DEFAULT ''::character varying,
    other text
);


ALTER TABLE public.logs OWNER TO root;

--
-- Name: logs_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.logs_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.logs_id_seq OWNER TO root;

--
-- Name: logs_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.logs_id_seq OWNED BY public.logs.id;


--
-- Name: midjourneys; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.midjourneys (
    id bigint NOT NULL,
    code bigint,
    user_id bigint,
    action character varying(40),
    mj_id text,
    prompt text,
    prompt_en text,
    description text,
    state text,
    submit_time bigint,
    start_time bigint,
    finish_time bigint,
    image_url text,
    video_url text,
    video_urls text,
    status character varying(20),
    progress character varying(30),
    fail_reason text,
    channel_id bigint,
    quota bigint,
    buttons text,
    properties text
);


ALTER TABLE public.midjourneys OWNER TO root;

--
-- Name: midjourneys_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.midjourneys_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.midjourneys_id_seq OWNER TO root;

--
-- Name: midjourneys_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.midjourneys_id_seq OWNED BY public.midjourneys.id;


--
-- Name: models; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.models (
    id bigint NOT NULL,
    model_name character varying(128) NOT NULL,
    description text,
    icon character varying(128),
    tags character varying(255),
    vendor_id bigint,
    endpoints text,
    status bigint DEFAULT 1,
    sync_official bigint DEFAULT 1,
    created_time bigint,
    updated_time bigint,
    deleted_at timestamp with time zone,
    name_rule bigint DEFAULT 0
);


ALTER TABLE public.models OWNER TO root;

--
-- Name: models_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.models_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.models_id_seq OWNER TO root;

--
-- Name: models_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.models_id_seq OWNED BY public.models.id;


--
-- Name: options; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.options (
    key text NOT NULL,
    value text
);


ALTER TABLE public.options OWNER TO root;

--
-- Name: passkey_credentials; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.passkey_credentials (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    credential_id character varying(512) NOT NULL,
    public_key text NOT NULL,
    attestation_type character varying(255),
    aa_guid character varying(512),
    sign_count bigint DEFAULT 0,
    clone_warning boolean,
    user_present boolean,
    user_verified boolean,
    backup_eligible boolean,
    backup_state boolean,
    transports text,
    attachment character varying(32),
    last_used_at timestamp with time zone,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


ALTER TABLE public.passkey_credentials OWNER TO root;

--
-- Name: passkey_credentials_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.passkey_credentials_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.passkey_credentials_id_seq OWNER TO root;

--
-- Name: passkey_credentials_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.passkey_credentials_id_seq OWNED BY public.passkey_credentials.id;


--
-- Name: perf_metrics; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.perf_metrics (
    id bigint NOT NULL,
    model_name character varying(128),
    "group" character varying(64),
    bucket_ts bigint,
    request_count bigint DEFAULT 0,
    success_count bigint DEFAULT 0,
    total_latency_ms bigint DEFAULT 0,
    ttft_sum_ms bigint DEFAULT 0,
    ttft_count bigint DEFAULT 0,
    output_tokens bigint DEFAULT 0,
    generation_ms bigint DEFAULT 0
);


ALTER TABLE public.perf_metrics OWNER TO root;

--
-- Name: perf_metrics_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.perf_metrics_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.perf_metrics_id_seq OWNER TO root;

--
-- Name: perf_metrics_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.perf_metrics_id_seq OWNED BY public.perf_metrics.id;


--
-- Name: prefill_groups; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.prefill_groups (
    id bigint NOT NULL,
    name character varying(64) NOT NULL,
    type character varying(32) NOT NULL,
    items json,
    description character varying(255),
    created_time bigint,
    updated_time bigint,
    deleted_at timestamp with time zone
);


ALTER TABLE public.prefill_groups OWNER TO root;

--
-- Name: prefill_groups_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.prefill_groups_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.prefill_groups_id_seq OWNER TO root;

--
-- Name: prefill_groups_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.prefill_groups_id_seq OWNED BY public.prefill_groups.id;


--
-- Name: quota_data; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.quota_data (
    id bigint NOT NULL,
    user_id bigint,
    username character varying(64) DEFAULT ''::character varying,
    model_name character varying(64) DEFAULT ''::character varying,
    created_at bigint,
    token_used bigint DEFAULT 0,
    count bigint DEFAULT 0,
    quota bigint DEFAULT 0
);


ALTER TABLE public.quota_data OWNER TO root;

--
-- Name: quota_data_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.quota_data_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.quota_data_id_seq OWNER TO root;

--
-- Name: quota_data_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.quota_data_id_seq OWNED BY public.quota_data.id;


--
-- Name: redemptions; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.redemptions (
    id bigint NOT NULL,
    user_id bigint,
    key character(32),
    status bigint DEFAULT 1,
    name text,
    quota bigint DEFAULT 100,
    created_time bigint,
    redeemed_time bigint,
    used_user_id bigint,
    deleted_at timestamp with time zone,
    expired_time bigint
);


ALTER TABLE public.redemptions OWNER TO root;

--
-- Name: redemptions_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.redemptions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.redemptions_id_seq OWNER TO root;

--
-- Name: redemptions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.redemptions_id_seq OWNED BY public.redemptions.id;


--
-- Name: setups; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.setups (
    id bigint NOT NULL,
    version character varying(50) NOT NULL,
    initialized_at bigint NOT NULL
);


ALTER TABLE public.setups OWNER TO root;

--
-- Name: setups_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.setups_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.setups_id_seq OWNER TO root;

--
-- Name: setups_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.setups_id_seq OWNED BY public.setups.id;


--
-- Name: subscription_orders; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.subscription_orders (
    id bigint NOT NULL,
    user_id bigint,
    plan_id bigint,
    money numeric,
    trade_no character varying(255),
    payment_method character varying(50),
    payment_provider character varying(50) DEFAULT ''::character varying,
    status text,
    create_time bigint,
    complete_time bigint,
    provider_payload text
);


ALTER TABLE public.subscription_orders OWNER TO root;

--
-- Name: subscription_orders_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.subscription_orders_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.subscription_orders_id_seq OWNER TO root;

--
-- Name: subscription_orders_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.subscription_orders_id_seq OWNED BY public.subscription_orders.id;


--
-- Name: subscription_plans; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.subscription_plans (
    id bigint NOT NULL,
    title character varying(128) NOT NULL,
    subtitle character varying(255) DEFAULT ''::character varying,
    price_amount numeric(10,6) DEFAULT 0.000000 NOT NULL,
    currency character varying(8) DEFAULT 'USD'::character varying NOT NULL,
    duration_unit character varying(16) DEFAULT 'month'::character varying NOT NULL,
    duration_value bigint DEFAULT 1 NOT NULL,
    custom_seconds bigint DEFAULT 0 NOT NULL,
    enabled boolean DEFAULT true,
    sort_order bigint DEFAULT 0,
    stripe_price_id character varying(128) DEFAULT ''::character varying,
    creem_product_id character varying(128) DEFAULT ''::character varying,
    max_purchase_per_user bigint DEFAULT 0,
    upgrade_group character varying(64) DEFAULT ''::character varying,
    total_amount bigint DEFAULT 0 NOT NULL,
    quota_reset_period character varying(16) DEFAULT 'never'::character varying,
    quota_reset_custom_seconds bigint DEFAULT 0,
    created_at bigint,
    updated_at bigint
);


ALTER TABLE public.subscription_plans OWNER TO root;

--
-- Name: subscription_plans_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.subscription_plans_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.subscription_plans_id_seq OWNER TO root;

--
-- Name: subscription_plans_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.subscription_plans_id_seq OWNED BY public.subscription_plans.id;


--
-- Name: subscription_pre_consume_records; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.subscription_pre_consume_records (
    id bigint NOT NULL,
    request_id character varying(64),
    user_id bigint,
    user_subscription_id bigint,
    pre_consumed bigint DEFAULT 0 NOT NULL,
    status character varying(32),
    created_at bigint,
    updated_at bigint
);


ALTER TABLE public.subscription_pre_consume_records OWNER TO root;

--
-- Name: subscription_pre_consume_records_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.subscription_pre_consume_records_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.subscription_pre_consume_records_id_seq OWNER TO root;

--
-- Name: subscription_pre_consume_records_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.subscription_pre_consume_records_id_seq OWNED BY public.subscription_pre_consume_records.id;


--
-- Name: tasks; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.tasks (
    id bigint NOT NULL,
    created_at bigint,
    updated_at bigint,
    task_id character varying(191),
    platform character varying(30),
    user_id bigint,
    "group" character varying(50),
    channel_id bigint,
    quota bigint,
    action character varying(40),
    status character varying(20),
    fail_reason text,
    submit_time bigint,
    start_time bigint,
    finish_time bigint,
    progress character varying(20),
    properties json,
    private_data json,
    data json
);


ALTER TABLE public.tasks OWNER TO root;

--
-- Name: tasks_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.tasks_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.tasks_id_seq OWNER TO root;

--
-- Name: tasks_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.tasks_id_seq OWNED BY public.tasks.id;


--
-- Name: tokens; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.tokens (
    id bigint NOT NULL,
    user_id bigint,
    key character varying(128),
    status bigint DEFAULT 1,
    name text,
    created_time bigint,
    accessed_time bigint,
    expired_time bigint DEFAULT '-1'::integer,
    remain_quota bigint DEFAULT 0,
    unlimited_quota boolean,
    model_limits_enabled boolean,
    model_limits text,
    allow_ips text DEFAULT ''::text,
    used_quota bigint DEFAULT 0,
    "group" text DEFAULT ''::text,
    cross_group_retry boolean,
    deleted_at timestamp with time zone
);


ALTER TABLE public.tokens OWNER TO root;

--
-- Name: tokens_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.tokens_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.tokens_id_seq OWNER TO root;

--
-- Name: tokens_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.tokens_id_seq OWNED BY public.tokens.id;


--
-- Name: top_ups; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.top_ups (
    id bigint NOT NULL,
    user_id bigint,
    amount bigint,
    money numeric,
    trade_no character varying(255),
    payment_method character varying(50),
    payment_provider character varying(50) DEFAULT ''::character varying,
    create_time bigint,
    complete_time bigint,
    status text
);


ALTER TABLE public.top_ups OWNER TO root;

--
-- Name: top_ups_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.top_ups_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.top_ups_id_seq OWNER TO root;

--
-- Name: top_ups_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.top_ups_id_seq OWNED BY public.top_ups.id;


--
-- Name: two_fa_backup_codes; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.two_fa_backup_codes (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    code_hash character varying(255) NOT NULL,
    is_used boolean,
    used_at timestamp with time zone,
    created_at timestamp with time zone,
    deleted_at timestamp with time zone
);


ALTER TABLE public.two_fa_backup_codes OWNER TO root;

--
-- Name: two_fa_backup_codes_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.two_fa_backup_codes_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.two_fa_backup_codes_id_seq OWNER TO root;

--
-- Name: two_fa_backup_codes_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.two_fa_backup_codes_id_seq OWNED BY public.two_fa_backup_codes.id;


--
-- Name: two_fas; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.two_fas (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    secret character varying(255) NOT NULL,
    is_enabled boolean,
    failed_attempts bigint DEFAULT 0,
    locked_until timestamp with time zone,
    last_used_at timestamp with time zone,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


ALTER TABLE public.two_fas OWNER TO root;

--
-- Name: two_fas_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.two_fas_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.two_fas_id_seq OWNER TO root;

--
-- Name: two_fas_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.two_fas_id_seq OWNED BY public.two_fas.id;


--
-- Name: user_oauth_bindings; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.user_oauth_bindings (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    provider_id bigint NOT NULL,
    provider_user_id character varying(256) NOT NULL,
    created_at timestamp with time zone
);


ALTER TABLE public.user_oauth_bindings OWNER TO root;

--
-- Name: user_oauth_bindings_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.user_oauth_bindings_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.user_oauth_bindings_id_seq OWNER TO root;

--
-- Name: user_oauth_bindings_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.user_oauth_bindings_id_seq OWNED BY public.user_oauth_bindings.id;


--
-- Name: user_subscriptions; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.user_subscriptions (
    id bigint NOT NULL,
    user_id bigint,
    plan_id bigint,
    amount_total bigint DEFAULT 0 NOT NULL,
    amount_used bigint DEFAULT 0 NOT NULL,
    start_time bigint,
    end_time bigint,
    status character varying(32),
    source character varying(32) DEFAULT 'order'::character varying,
    last_reset_time bigint DEFAULT 0,
    next_reset_time bigint DEFAULT 0,
    upgrade_group character varying(64) DEFAULT ''::character varying,
    prev_user_group character varying(64) DEFAULT ''::character varying,
    created_at bigint,
    updated_at bigint
);


ALTER TABLE public.user_subscriptions OWNER TO root;

--
-- Name: user_subscriptions_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.user_subscriptions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.user_subscriptions_id_seq OWNER TO root;

--
-- Name: user_subscriptions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.user_subscriptions_id_seq OWNED BY public.user_subscriptions.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.users (
    id bigint NOT NULL,
    username text,
    password text NOT NULL,
    display_name text,
    role bigint DEFAULT 1,
    status bigint DEFAULT 1,
    email text,
    github_id text,
    discord_id text,
    oidc_id text,
    wechat_id text,
    telegram_id text,
    access_token character(32),
    quota bigint DEFAULT 0,
    used_quota bigint DEFAULT 0,
    request_count bigint DEFAULT 0,
    "group" character varying(64) DEFAULT 'default'::character varying,
    aff_code character varying(32),
    aff_count bigint DEFAULT 0,
    aff_quota bigint DEFAULT 0,
    aff_history bigint DEFAULT 0,
    inviter_id bigint,
    deleted_at timestamp with time zone,
    linux_do_id text,
    setting text,
    remark character varying(255),
    stripe_customer character varying(64),
    created_at bigint,
    last_login_at bigint DEFAULT 0
);


ALTER TABLE public.users OWNER TO root;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.users_id_seq OWNER TO root;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: vendors; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.vendors (
    id bigint NOT NULL,
    name character varying(128) NOT NULL,
    description text,
    icon character varying(128),
    status bigint DEFAULT 1,
    created_time bigint,
    updated_time bigint,
    deleted_at timestamp with time zone
);


ALTER TABLE public.vendors OWNER TO root;

--
-- Name: vendors_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.vendors_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.vendors_id_seq OWNER TO root;

--
-- Name: vendors_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.vendors_id_seq OWNED BY public.vendors.id;


--
-- Name: channels id; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.channels ALTER COLUMN id SET DEFAULT nextval('public.channels_id_seq'::regclass);


--
-- Name: checkins id; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.checkins ALTER COLUMN id SET DEFAULT nextval('public.checkins_id_seq'::regclass);


--
-- Name: custom_oauth_providers id; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.custom_oauth_providers ALTER COLUMN id SET DEFAULT nextval('public.custom_oauth_providers_id_seq'::regclass);


--
-- Name: logs id; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.logs ALTER COLUMN id SET DEFAULT nextval('public.logs_id_seq'::regclass);


--
-- Name: midjourneys id; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.midjourneys ALTER COLUMN id SET DEFAULT nextval('public.midjourneys_id_seq'::regclass);


--
-- Name: models id; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.models ALTER COLUMN id SET DEFAULT nextval('public.models_id_seq'::regclass);


--
-- Name: passkey_credentials id; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.passkey_credentials ALTER COLUMN id SET DEFAULT nextval('public.passkey_credentials_id_seq'::regclass);


--
-- Name: perf_metrics id; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.perf_metrics ALTER COLUMN id SET DEFAULT nextval('public.perf_metrics_id_seq'::regclass);


--
-- Name: prefill_groups id; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.prefill_groups ALTER COLUMN id SET DEFAULT nextval('public.prefill_groups_id_seq'::regclass);


--
-- Name: quota_data id; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.quota_data ALTER COLUMN id SET DEFAULT nextval('public.quota_data_id_seq'::regclass);


--
-- Name: redemptions id; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.redemptions ALTER COLUMN id SET DEFAULT nextval('public.redemptions_id_seq'::regclass);


--
-- Name: setups id; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.setups ALTER COLUMN id SET DEFAULT nextval('public.setups_id_seq'::regclass);


--
-- Name: subscription_orders id; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.subscription_orders ALTER COLUMN id SET DEFAULT nextval('public.subscription_orders_id_seq'::regclass);


--
-- Name: subscription_plans id; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.subscription_plans ALTER COLUMN id SET DEFAULT nextval('public.subscription_plans_id_seq'::regclass);


--
-- Name: subscription_pre_consume_records id; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.subscription_pre_consume_records ALTER COLUMN id SET DEFAULT nextval('public.subscription_pre_consume_records_id_seq'::regclass);


--
-- Name: tasks id; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.tasks ALTER COLUMN id SET DEFAULT nextval('public.tasks_id_seq'::regclass);


--
-- Name: tokens id; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.tokens ALTER COLUMN id SET DEFAULT nextval('public.tokens_id_seq'::regclass);


--
-- Name: top_ups id; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.top_ups ALTER COLUMN id SET DEFAULT nextval('public.top_ups_id_seq'::regclass);


--
-- Name: two_fa_backup_codes id; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.two_fa_backup_codes ALTER COLUMN id SET DEFAULT nextval('public.two_fa_backup_codes_id_seq'::regclass);


--
-- Name: two_fas id; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.two_fas ALTER COLUMN id SET DEFAULT nextval('public.two_fas_id_seq'::regclass);


--
-- Name: user_oauth_bindings id; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.user_oauth_bindings ALTER COLUMN id SET DEFAULT nextval('public.user_oauth_bindings_id_seq'::regclass);


--
-- Name: user_subscriptions id; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.user_subscriptions ALTER COLUMN id SET DEFAULT nextval('public.user_subscriptions_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Name: vendors id; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.vendors ALTER COLUMN id SET DEFAULT nextval('public.vendors_id_seq'::regclass);


--
-- Data for Name: abilities; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.abilities ("group", model, channel_id, enabled, priority, weight, tag) FROM stdin;
default	gpt-5.4	1	t	0	0	
default	gpt-image-2	1	t	0	0	
\.


--
-- Data for Name: channels; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.channels (id, type, key, open_ai_organization, test_model, status, name, weight, created_time, test_time, response_time, base_url, other, balance, balance_updated_time, models, "group", used_quota, model_mapping, status_code_mapping, priority, auto_ban, other_info, tag, setting, param_override, header_override, remark, channel_info, settings) FROM stdin;
1	1	sk-1bd27cc7b4a267e79828b10f04435361a423a4cad7a91847797619ec96584c78			1	openai	0	1779422254	1779426031	4465	http://119.29.249.17:8080		0	0	gpt-5.4,gpt-image-2	default	13750072			0	1			{"force_format":false,"thinking_to_content":false,"proxy":"","pass_through_body_enabled":false,"system_prompt":"","system_prompt_override":false}		\N	\N	{"is_multi_key":false,"multi_key_size":0,"multi_key_status_list":null,"multi_key_polling_index":0,"multi_key_mode":"random"}	{"allow_service_tier":false,"disable_store":false,"allow_safety_identifier":false,"allow_include_obfuscation":false,"upstream_model_update_check_enabled":false,"upstream_model_update_auto_sync_enabled":false,"upstream_model_update_ignored_models":[],"upstream_model_update_last_detected_models":[],"upstream_model_update_last_check_time":0}
\.


--
-- Data for Name: checkins; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.checkins (id, user_id, checkin_date, quota_awarded, created_at) FROM stdin;
\.


--
-- Data for Name: custom_oauth_providers; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.custom_oauth_providers (id, name, slug, icon, enabled, client_id, client_secret, authorization_endpoint, token_endpoint, user_info_endpoint, scopes, user_id_field, username_field, display_name_field, email_field, well_known, auth_style, access_policy, access_denied_message, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: logs; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.logs (id, user_id, created_at, type, content, username, token_name, model_name, quota, prompt_tokens, completion_tokens, use_time, is_stream, channel_id, channel_name, token_id, "group", ip, request_id, upstream_request_id, other) FROM stdin;
1	1	1779425071	2	模型测试	xiaopihong	模型测试	gpt-5.4	63	18	11	2	f	1	\N	0	default				{"admin_info":{"use_channel":null},"cache_ratio":1,"cache_tokens":0,"completion_ratio":6,"frt":-1000,"group_ratio":1,"model_price":-1,"model_ratio":0.75,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","user_group_ratio":-1}
2	2	1779425293	3	管理员增加用户额度 ＄20.000000 额度	sso-admin			0	0	0	0	f	0	\N	0					{"admin_info":{"admin_id":1,"admin_username":"xiaopihong"}}
3	2	1779425321	2		sso-admin	自用	gpt-5.4	72	18	13	2	t	1	\N	1	default		202605220448391916616318268d9d6LZuN90jR		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":1,"cache_tokens":0,"completion_ratio":6,"frt":1759,"group_ratio":1,"model_price":-1,"model_ratio":0.75,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","stream_status":{"end_reason":"done","status":"ok"},"user_group_ratio":-1}
4	2	1779425493	5	status_code=502, Upstream request failed	sso-admin	自用	gpt-image-2	0	0	0	1	f	1	\N	1	default		202605220451321160959288268d9d6gP36Cnht		{"admin_info":{"use_channel":["1"]},"channel_id":1,"channel_name":"openai","channel_type":1,"error_code":"unknown_error","error_type":"openai_error","request_path":"/v1/images/generations","status_code":502}
5	1	1779425876	2	模型测试	xiaopihong	模型测试	gpt-5.4	63	18	11	3	f	1	\N	0	default				{"admin_info":{"use_channel":null},"cache_ratio":1,"cache_tokens":0,"completion_ratio":6,"frt":-1000,"group_ratio":1,"model_price":-1,"model_ratio":0.75,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","user_group_ratio":-1}
6	1	1779426004	2	模型测试	xiaopihong	模型测试	gpt-image-2	250000	41	196	75	f	1	\N	0	default				{"admin_info":{"use_channel":null},"cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/generations","user_group_ratio":-1}
7	1	1779426031	2	模型测试	xiaopihong	模型测试	gpt-5.4	63	18	11	4	t	1	\N	0	default				{"admin_info":{"use_channel":null},"cache_ratio":1,"cache_tokens":0,"completion_ratio":6,"frt":3421,"group_ratio":1,"model_price":-1,"model_ratio":0.75,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","stream_status":{"end_reason":"done","status":"ok"},"user_group_ratio":-1}
8	2	1779426244	2	大小 1024x1024, 品质 standard, 生成数量 1	sso-admin	自用	gpt-image-2	250000	132	7024	184	f	1	\N	1	default		202605220501004673782848268d9d6Y3Zgxgfi		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/generations","user_group_ratio":-1}
9	2	1779426360	5	status_code=502, Upstream request failed	sso-admin	自用	gpt-image-2	0	0	0	4	f	1	\N	1	default		202605220505562634141758268d9d6g40KI7sM		{"admin_info":{"use_channel":["1"]},"channel_id":1,"channel_name":"openai","channel_type":1,"error_code":"unknown_error","error_type":"openai_error","request_path":"/v1/images/generations","status_code":502}
10	2	1779426394	5	status_code=502, Upstream request failed	sso-admin	自用	gpt-image-2	0	0	0	0	f	1	\N	1	default		202605220506338446758048268d9d6aSfjrfEF		{"admin_info":{"use_channel":["1"]},"channel_id":1,"channel_name":"openai","channel_type":1,"error_code":"unknown_error","error_type":"openai_error","request_path":"/v1/images/generations","status_code":502}
11	2	1779426637	2	大小 1024x1024, 品质 standard, 生成数量 1	sso-admin	自用	gpt-image-2	250000	151	1756	68	f	1	\N	1	default		202605220509291349631498268d9d6rdfnr6Ni		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/generations","user_group_ratio":-1}
12	2	1779426902	2	大小 1024x1024, 品质 standard, 生成数量 1	sso-admin	自用	gpt-image-2	250000	162	1756	136	f	1	\N	1	default		20260522051246858025998268d9d6eoOmWHqK		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/generations","user_group_ratio":-1}
13	2	1779433023	5	status_code=502, Upstream request failed	sso-admin	自用	gpt-image-2	0	0	0	66	f	1	\N	1	default		20260522065557930629958268d9d6YiRr7KUy		{"admin_info":{"use_channel":["1"]},"channel_id":1,"channel_name":"openai","channel_type":1,"error_code":"unknown_error","error_type":"openai_error","request_path":"/v1/images/edits","status_code":502}
14	2	1779433128	2	大小 1024x1024, 品质 standard, 生成数量 1	sso-admin	自用	gpt-image-2	250000	1088	1756	88	f	1	\N	1	default		202605220657204737298958268d9d6aWXw9ugw		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"image":true,"image_output":928,"image_ratio":0,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
15	2	1779433565	2	大小 1024x1024, 品质 standard, 生成数量 1	sso-admin	自用	gpt-image-2	250000	1117	1756	116	f	1	\N	1	default		202605220704096637926588268d9d649cAHRAT		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"image":true,"image_output":928,"image_ratio":0,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
16	2	1779434034	2	大小 1024x1024, 品质 standard, 生成数量 1	sso-admin	自用	gpt-image-2	250000	1066	1756	88	f	1	\N	1	default		202605220712268680305658268d9d6o99P71cD		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"image":true,"image_output":928,"image_ratio":0,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
17	2	1779434563	2	大小 1536x1024, 品质 standard, 生成数量 1	sso-admin	自用	gpt-image-2	250000	1079	1372	72	f	1	\N	1	default		202605220721315057605798268d9d6FahRaSkq		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"image":true,"image_output":928,"image_ratio":0,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
18	2	1779434970	2	大小 1792x1024, 品质 standard, 生成数量 1	sso-admin	自用	gpt-image-2	250000	1128	1243	73	f	1	\N	1	default		202605220728178686838378268d9d6DgbJeEQC		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"image":true,"image_output":928,"image_ratio":0,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
19	2	1779435334	2	大小 1792x1024, 品质 standard, 生成数量 1	sso-admin	自用	gpt-image-2	250000	1157	1243	106	f	1	\N	1	default		202605220733487427439748268d9d6FAjPbk0S		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"image":true,"image_output":928,"image_ratio":0,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
20	2	1779435676	2	大小 1792x1024, 品质 standard, 生成数量 1	sso-admin	自用	gpt-image-2	250000	1123	1243	70	f	1	\N	1	default		202605220740063140319958268d9d6eZpu2e40		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"image":true,"image_output":928,"image_ratio":0,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
21	2	1779435812	5	status_code=502, Upstream request failed	sso-admin	自用	gpt-image-2	0	0	0	14	f	1	\N	1	default		202605220743188075698808268d9d6Qm52cOnU		{"admin_info":{"use_channel":["1"]},"channel_id":1,"channel_name":"openai","channel_type":1,"error_code":"unknown_error","error_type":"openai_error","request_path":"/v1/images/edits","status_code":502}
22	2	1779435912	2	大小 1792x1024, 品质 standard, 生成数量 1	sso-admin	自用	gpt-image-2	250000	1149	1243	90	f	1	\N	1	default		202605220743422772517328268d9d6wP5R6fFE		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"image":true,"image_output":928,"image_ratio":0,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
23	2	1779436427	2	大小 1792x1024, 品质 standard, 生成数量 1	sso-admin	自用	gpt-image-2	250000	1175	5063	130	f	1	\N	1	default		202605220751373394399148268d9d63fudG5xc		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"image":true,"image_output":928,"image_ratio":0,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
24	2	1779436696	2	大小 1792x1024, 品质 standard, 生成数量 1	sso-admin	自用	gpt-image-2	250000	1126	1243	95	f	1	\N	1	default		20260522075641210252478268d9d6DM58cdkh		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"image":true,"image_output":928,"image_ratio":0,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
25	2	1779436974	5	status_code=502, Upstream request failed	sso-admin	自用	gpt-image-2	0	0	0	88	f	1	\N	1	default		202605220801251709887668268d9d6657M1ItH		{"admin_info":{"use_channel":["1"]},"channel_id":1,"channel_name":"openai","channel_type":1,"error_code":"unknown_error","error_type":"openai_error","request_path":"/v1/images/edits","status_code":502}
26	2	1779437052	2	大小 1792x1024, 品质 standard, 生成数量 1	sso-admin	自用	gpt-image-2	250000	1109	1243	74	f	1	\N	1	default		202605220802582390236918268d9d6opEcN1Lc		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"image":true,"image_output":928,"image_ratio":0,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
27	2	1779437882	2	大小 1024x1024, 品质 standard, 生成数量 1	sso-admin	自用	gpt-image-2	250000	1140	1756	130	f	1	\N	1	default		202605220815528190647198268d9d6nFJx8Xau		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"image":true,"image_output":928,"image_ratio":0,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
28	2	1779438183	2	大小 1024x1024, 品质 standard, 生成数量 1	sso-admin	自用	gpt-image-2	250000	1232	1756	108	f	1	\N	1	default		20260522082115175247298268d9d6tsThERrb		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"image":true,"image_output":928,"image_ratio":0,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
29	2	1779438568	2	大小 1024x1024, 品质 standard, 生成数量 1	sso-admin	自用	gpt-image-2	250000	1202	1756	120	f	1	\N	1	default		202605220827284188953618268d9d6N61agVOH		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"image":true,"image_output":928,"image_ratio":0,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
30	2	1779444957	2	大小 1024x1024, 品质 standard, 生成数量 1	sso-admin	自用	gpt-image-2	250000	1276	1756	89	f	1	\N	1	default		202605221014285918490878268d9d6ZGKExmpi		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"image":true,"image_output":928,"image_ratio":0,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
31	2	1779445385	5	status_code=502, Upstream request failed	sso-admin	自用	gpt-image-2	0	0	0	85	f	1	\N	1	default		2026052210214089705008268d9d69pzV0gcY		{"admin_info":{"use_channel":["1"]},"channel_id":1,"channel_name":"openai","channel_type":1,"error_code":"unknown_error","error_type":"openai_error","request_path":"/v1/images/edits","status_code":502}
32	2	1779445602	2	大小 1024x1024, 品质 standard, 生成数量 1	sso-admin	自用	gpt-image-2	250000	1280	7024	202	f	1	\N	1	default		202605221023208484943188268d9d6jOjKXYku		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"image":true,"image_output":928,"image_ratio":0,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
33	2	1779446194	2	大小 1024x1024, 品质 standard, 生成数量 1	sso-admin	自用	gpt-image-2	250000	1152	1756	96	f	1	\N	1	default		202605221034586506213248268d9d6jBOXreVc		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"image":true,"image_output":928,"image_ratio":0,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
34	2	1779446826	2	大小 1024x1024, 品质 standard, 生成数量 1	sso-admin	自用	gpt-image-2	250000	1281	7024	237	f	1	\N	1	default		202605221043097048254158268d9d6U2igcyGr		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"image":true,"image_output":928,"image_ratio":0,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
35	2	1779447494	2	大小 1024x1024, 品质 standard, 生成数量 1	sso-admin	自用	gpt-image-2	250000	1286	7024	218	f	1	\N	1	default		202605221054367351797708268d9d6a7SczK30		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"image":true,"image_output":928,"image_ratio":0,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
36	2	1779462795	2	大小 1024x1024, 品质 standard, 生成数量 1	sso-admin	自用	gpt-image-2	250000	1256	7024	258	f	1	\N	1	default		202605221508575115769648268d9d6SYS3IvAL		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"image":true,"image_output":928,"image_ratio":0,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
37	2	1779463316	2	大小 1024x1024, 品质 standard, 生成数量 1	sso-admin	自用	gpt-image-2	250000	1295	7024	267	f	1	\N	1	default		202605221517291002638258268d9d61rBuvLjV		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"image":true,"image_output":928,"image_ratio":0,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
38	2	1779463716	2	大小 1024x1024, 品质 standard, 生成数量 1	sso-admin	自用	gpt-image-2	250000	1221	7024	210	f	1	\N	1	default		202605221525061026586058268d9d6QXDSo1Bu		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"image":true,"image_output":928,"image_ratio":0,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
39	2	1779464516	5	status_code=500, upstream error: do request failed	sso-admin	自用	gpt-image-2	0	0	0	360	f	1	\N	1	default		202605221535558512775108268d9d63WgdD2YZ		{"admin_info":{"use_channel":["1"]},"channel_id":1,"channel_name":"openai","channel_type":1,"error_code":"do_request_failed","error_type":"new_api_error","request_path":"/v1/images/edits","status_code":500}
40	2	1779464618	2	大小 1024x1024, 品质 standard, 生成数量 1	sso-admin	自用	gpt-image-2	250000	1173	1756	108	f	1	\N	1	default		202605221541503944702778268d9d6eLuuQ83H		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"image":true,"image_output":928,"image_ratio":0,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
41	2	1779465725	2	大小 1024x1024, 品质 standard, 生成数量 1	sso-admin	自用	gpt-image-2	250000	1272	1756	96	f	1	\N	1	default		202605221600294706348408268d9d6yr80LPfB		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"image":true,"image_output":1024,"image_ratio":0,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
42	2	1779466346	2	大小 1024x1024, 品质 standard, 生成数量 1	sso-admin	自用	gpt-image-2	250000	1287	1756	135	f	1	\N	1	default		202605221610117775379698268d9d6fGIzBUTR		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"image":true,"image_output":1024,"image_ratio":0,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
43	2	1779466726	2	大小 1024x1024, 品质 standard, 生成数量 1	sso-admin	自用	gpt-image-2	250000	1319	1756	90	f	1	\N	1	default		202605221617168068042328268d9d6aU1tmrFC		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"image":true,"image_output":1024,"image_ratio":0,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
44	2	1779467277	2	大小 1024x1024, 品质 standard, 生成数量 1	sso-admin	自用	gpt-image-2	250000	1220	1756	108	f	1	\N	1	default		202605221626094221633898268d9d6Te3Qohzp		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"image":true,"image_output":1024,"image_ratio":0,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
45	2	1779467949	2	大小 1024x1024, 品质 standard, 生成数量 1	sso-admin	自用	gpt-image-2	250000	1238	1756	89	f	1	\N	1	default		202605221637404011158848268d9d6VZhUYF1k		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"image":true,"image_output":1024,"image_ratio":0,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
46	2	1779468498	2	大小 1024x1024, 品质 standard, 生成数量 1	sso-admin	自用	gpt-image-2	250000	1299	7024	197	f	1	\N	1	default		202605221645016960864038268d9d6fS7Ho9qN		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"image":true,"image_output":1024,"image_ratio":0,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
47	2	1779468774	2	大小 1024x1024, 品质 standard, 生成数量 1	sso-admin	自用	gpt-image-2	250000	1246	1756	101	f	1	\N	1	default		20260522165113631963368268d9d6W7TcG0BW		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"image":true,"image_output":1024,"image_ratio":0,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
48	2	1779530463	5	status_code=502, Upstream request failed	sso-admin	默认	gpt-image-2	0	0	0	600	f	1	\N	4	default		202605230951031030768508268d9d64FdFfcX1		{"admin_info":{"use_channel":["1"]},"channel_id":1,"channel_name":"openai","channel_type":1,"error_code":"unknown_error","error_type":"openai_error","request_path":"/v1/images/edits","status_code":502}
49	2	1779530558	2	大小 1344x704, 品质 standard, 生成数量 1	sso-admin	默认	gpt-image-2	250000	1193	3536	158	f	1	\N	4	default		20260523100000265273698268d9d6tjvUUGWn		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"image":true,"image_output":924,"image_ratio":0,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
50	2	1779531717	2	大小 1344x704, 品质 standard, 生成数量 1	sso-admin	默认	gpt-image-2	250000	1218	884	85	f	1	\N	4	default		20260523102032925710358268d9d6Wz9a2l5W		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"image":true,"image_output":924,"image_ratio":0,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
51	2	1779532570	2	大小 1344x704, 品质 standard, 生成数量 1	sso-admin	默认	gpt-image-2	250000	1292	3536	138	f	1	\N	4	default		202605231033526644857738268d9d6cdplxnXI		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"image":true,"image_output":924,"image_ratio":0,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
52	2	1779534018	5	status_code=502, Upstream request failed	sso-admin	默认	gpt-image-2	0	0	0	96	f	1	\N	4	default		202605231058419461680168268d9d6BGMD9XbM		{"admin_info":{"use_channel":["1"]},"channel_id":1,"channel_name":"openai","channel_type":1,"error_code":"unknown_error","error_type":"openai_error","request_path":"/v1/images/edits","status_code":502}
53	2	1779534115	2	大小 1216x544, 品质 standard, 生成数量 1	sso-admin	默认	gpt-image-2	250000	903	671	70	f	1	\N	4	default		202605231100453005354048268d9d6JXeh19dw		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"image":true,"image_output":646,"image_ratio":0,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
54	2	1779534832	2	大小 1216x544, 品质 standard, 生成数量 1	sso-admin	默认	gpt-image-2	250000	853	671	49	f	1	\N	4	default		202605231113031507284328268d9d6AkrBsLtZ		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"image":true,"image_output":646,"image_ratio":0,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
55	2	1779535575	2	大小 1216x544, 品质 standard, 生成数量 1	sso-admin	默认	gpt-image-2	250000	986	671	87	f	1	\N	4	default		202605231124488379492688268d9d6wcJYi0UF		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"image":true,"image_output":646,"image_ratio":0,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
56	2	1779536267	3	管理员增加用户额度 ＄100.000000 额度	sso-admin			0	0	0	0	f	0	\N	0					{"admin_info":{"admin_id":1,"admin_username":"xiaopihong"}}
57	3	1779536278	3	管理员增加用户额度 ＄100.000000 额度	xph-admin			0	0	0	0	f	0	\N	0					{"admin_info":{"admin_id":1,"admin_username":"xiaopihong"}}
58	2	1779536597	2	大小 1216x544, 品质 standard, 生成数量 1	sso-admin	默认	gpt-image-2	250000	1004	671	75	f	1	\N	4	default		202605231142027973663578268d9d67wv8ANfV		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"image":true,"image_output":646,"image_ratio":0,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
59	2	1779537393	2	大小 1216x544, 品质 standard, 生成数量 1	sso-admin	默认	gpt-image-2	250000	1014	671	86	f	1	\N	4	default		202605231155077153436618268d9d6zGUqG9gU		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"image":true,"image_output":646,"image_ratio":0,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
60	2	1779537803	5	status_code=502, Upstream request failed	sso-admin	默认	gpt-image-2	0	0	0	103	f	1	\N	4	default		202605231201395925392458268d9d6nMj5PARk		{"admin_info":{"use_channel":["1"]},"channel_id":1,"channel_name":"openai","channel_type":1,"error_code":"unknown_error","error_type":"openai_error","request_path":"/v1/images/edits","status_code":502}
61	2	1779537895	2	大小 1216x544, 品质 standard, 生成数量 1	sso-admin	默认	gpt-image-2	250000	1092	671	61	f	1	\N	4	default		202605231203541783249178268d9d61uCUh3Qm		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"image":true,"image_output":646,"image_ratio":0,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
62	2	1779538487	2	大小 1216x544, 品质 standard, 生成数量 1	sso-admin	默认	gpt-image-2	250000	968	671	104	f	1	\N	4	default		202605231213037904217458268d9d6IEnVXDcq		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"image":true,"image_output":646,"image_ratio":0,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
63	2	1779539549	2	大小 1024x640, 品质 standard, 生成数量 1	sso-admin	默认	gpt-image-2	250000	967	956	121	f	1	\N	4	default		202605231230288967693928268d9d6YIVEUpsN		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"image":true,"image_output":640,"image_ratio":0,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
64	2	1779539908	2	大小 1216x544, 品质 standard, 生成数量 1	sso-admin	默认	gpt-image-2	250000	905	671	84	f	1	\N	4	default		202605231237046833503598268d9d6CBbymk6X		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"image":true,"image_output":646,"image_ratio":0,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
65	2	1779540407	2	大小 1024x640, 品质 standard, 生成数量 1	sso-admin	默认	gpt-image-2	250000	1030	956	90	f	1	\N	4	default		202605231245177829685168268d9d6F1S1aRzV		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"image":true,"image_output":640,"image_ratio":0,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
66	2	1779540546	2	大小 1216x544, 品质 standard, 生成数量 1	sso-admin	默认	gpt-image-2	250000	957	671	91	f	1	\N	4	default		202605231247357886322548268d9d6VbuLrTpj		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"image":true,"image_output":646,"image_ratio":0,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
67	2	1779541072	5	status_code=502, Upstream request failed	sso-admin	默认	gpt-image-2	0	0	0	2	f	1	\N	4	default		202605231257494373767918268d9d6FkeUffHV		{"admin_info":{"use_channel":["1"]},"channel_id":1,"channel_name":"openai","channel_type":1,"error_code":"unknown_error","error_type":"openai_error","request_path":"/v1/images/edits","status_code":502}
68	2	1779541097	5	status_code=502, Upstream request failed	sso-admin	默认	gpt-image-2	0	0	0	1	f	1	\N	4	default		202605231258156519460378268d9d6cJ4dHfS1		{"admin_info":{"use_channel":["1"]},"channel_id":1,"channel_name":"openai","channel_type":1,"error_code":"unknown_error","error_type":"openai_error","request_path":"/v1/images/edits","status_code":502}
69	2	1779541101	5	status_code=502, Upstream request failed	sso-admin	默认	gpt-image-2	0	0	0	1	f	1	\N	4	default		20260523125820993505308268d9d6rn1ZwLhb		{"admin_info":{"use_channel":["1"]},"channel_id":1,"channel_name":"openai","channel_type":1,"error_code":"unknown_error","error_type":"openai_error","request_path":"/v1/images/edits","status_code":502}
70	2	1779541139	5	status_code=502, Upstream request failed	sso-admin	默认	gpt-image-2	0	0	0	3	f	1	\N	4	default		202605231258563795409268268d9d6leyeVOR5		{"admin_info":{"use_channel":["1"]},"channel_id":1,"channel_name":"openai","channel_type":1,"error_code":"unknown_error","error_type":"openai_error","request_path":"/v1/images/edits","status_code":502}
71	2	1779541316	2	大小 1216x544, 品质 standard, 生成数量 1	sso-admin	默认	gpt-image-2	250000	946	671	121	f	1	\N	4	default		202605231259557100610458268d9d6kMcL8TDa		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"image":true,"image_output":646,"image_ratio":0,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
72	2	1779542212	5	status_code=502, Upstream request failed	sso-admin	默认	gpt-image-2	0	0	0	67	f	1	\N	4	default		202605231315446436702348268d9d6wB2BxV8b		{"admin_info":{"use_channel":["1"]},"channel_id":1,"channel_name":"openai","channel_type":1,"error_code":"unknown_error","error_type":"openai_error","request_path":"/v1/images/edits","status_code":502}
73	2	1779542421	5	status_code=502, Upstream request failed	sso-admin	默认	gpt-image-2	0	0	0	600	f	1	\N	4	default		202605231310201902845028268d9d6sJODgpIA		{"admin_info":{"use_channel":["1"]},"channel_id":1,"channel_name":"openai","channel_type":1,"error_code":"unknown_error","error_type":"openai_error","request_path":"/v1/images/edits","status_code":502}
74	2	1779542691	2	大小 1216x544, 品质 standard, 生成数量 1	sso-admin	默认	gpt-image-2	250000	1040	671	125	f	1	\N	4	default		202605231322464843677278268d9d6Ib8XGlr3		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"image":true,"image_output":646,"image_ratio":0,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
75	2	1779543256	2	大小 1216x544, 品质 standard, 生成数量 1	sso-admin	默认	gpt-image-2	250000	924	671	66	f	1	\N	4	default		202605231333106547064448268d9d6Si6REKQF		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"image":true,"image_output":646,"image_ratio":0,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
76	2	1779544900	2	大小 1216x544, 品质 standard, 生成数量 1	sso-admin	默认	gpt-image-2	250000	1035	671	81	f	1	\N	4	default		202605231400191263584938268d9d6BdlWQCf3		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"image":true,"image_output":646,"image_ratio":0,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
77	2	1779545530	2	大小 1216x544, 品质 standard, 生成数量 1	sso-admin	默认	gpt-image-2	250000	953	671	127	f	1	\N	4	default		202605231410039276794798268d9d6vgJSuVgU		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"image":true,"image_output":646,"image_ratio":0,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
78	2	1779545824	2	大小 1216x544, 品质 standard, 生成数量 1	sso-admin	默认	gpt-image-2	250000	861	671	75	f	1	\N	4	default		202605231415493553259178268d9d6XrFYKXHL		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"image":true,"image_output":646,"image_ratio":0,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
79	2	1779548898	2	大小 1216x544, 品质 standard, 生成数量 1	sso-admin	默认	gpt-image-2	250000	957	671	146	f	1	\N	4	default		202605231505522285439588268d9d6lmSzcmeo		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"image":true,"image_output":646,"image_ratio":0,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
80	2	1779550524	2	大小 1216x544, 品质 standard, 生成数量 1	sso-admin	默认	gpt-image-2	250000	1031	671	161	f	1	\N	4	default		202605231532431425850738268d9d6xa9nBoal		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"image":true,"image_output":646,"image_ratio":0,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
\.


--
-- Data for Name: midjourneys; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.midjourneys (id, code, user_id, action, mj_id, prompt, prompt_en, description, state, submit_time, start_time, finish_time, image_url, video_url, video_urls, status, progress, fail_reason, channel_id, quota, buttons, properties) FROM stdin;
\.


--
-- Data for Name: models; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.models (id, model_name, description, icon, tags, vendor_id, endpoints, status, sync_official, created_time, updated_time, deleted_at, name_rule) FROM stdin;
1	gpt-5.4		OpenAI		1	{\n  "openai": {\n    "path": "/v1/chat/completions",\n    "method": "POST"\n  }\n}	1	1	1779425439	1779425439	\N	0
2	gpt-image-2		OpenAI		1	{\n  "image-generation": {\n    "path": "/v1/images/generations",\n    "method": "POST"\n  }\n}	1	1	1779425458	1779426011	\N	0
\.


--
-- Data for Name: options; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.options (key, value) FROM stdin;
SelfUseModeEnabled	false
DemoSiteEnabled	false
ServerAddress	http://localhost:2999
ModelPrice	{\n  "gpt-image-2": 0.5,\n  "black-forest-labs/flux-1.1-pro": 0.04,\n  "dall-e-3": 0.04,\n  "gpt-4-gizmo-*": 0.1,\n  "gpt-4o-mini-tts": 0.3,\n  "imagen-3.0-generate-002": 0.03,\n  "mj_blend": 0.1,\n  "mj_custom_zoom": 0,\n  "mj_describe": 0.05,\n  "mj_edits": 0.1,\n  "mj_high_variation": 0.1,\n  "mj_imagine": 0.1,\n  "mj_inpaint": 0,\n  "mj_low_variation": 0.1,\n  "mj_modal": 0.1,\n  "mj_pan": 0.1,\n  "mj_reroll": 0.1,\n  "mj_shorten": 0.1,\n  "mj_upload": 0.05,\n  "mj_upscale": 0.05,\n  "mj_variation": 0.1,\n  "mj_video": 0.8,\n  "mj_zoom": 0.1,\n  "sora-2": 0.3,\n  "sora-2-pro": 0.5,\n  "suno_lyrics": 0.01,\n  "suno_music": 0.1,\n  "swap_face": 0.05,\n  "veo-3.0-fast-generate-001": 0.15,\n  "veo-3.0-generate-001": 0.4,\n  "veo-3.1-fast-generate-preview": 0.15,\n  "veo-3.1-generate-preview": 0.4\n}
ModelRatio	{\n  "gpt-5.4": 0.75,\n  "360GPT_S2_V9": 0.8572,\n  "360gpt-pro": 0.8572,\n  "360gpt-turbo": 0.0858,\n  "360gpt-turbo-responsibility-8k": 0.8572,\n  "360gpt2-pro": 0.8572,\n  "ada": 10,\n  "babbage": 10,\n  "babbage-002": 0.2,\n  "bge-large-en": 0.13698630137,\n  "bge-large-zh": 0.13698630137,\n  "BLOOMZ-7B": 0.27397260274,\n  "chatglm_lite": 0.1429,\n  "chatglm_pro": 0.7143,\n  "chatglm_std": 0.3572,\n  "chatglm_turbo": 0.3572,\n  "chatgpt-4o-latest": 2.5,\n  "claude-3-5-haiku-20241022": 0.5,\n  "claude-3-5-sonnet-20240620": 1.5,\n  "claude-3-5-sonnet-20241022": 1.5,\n  "claude-3-7-sonnet-20250219": 1.5,\n  "claude-3-7-sonnet-20250219-thinking": 1.5,\n  "claude-3-haiku-20240307": 0.125,\n  "claude-3-opus-20240229": 7.5,\n  "claude-3-sonnet-20240229": 1.5,\n  "claude-haiku-4-5-20251001": 0.5,\n  "claude-opus-4-1-20250805": 7.5,\n  "claude-opus-4-20250514": 7.5,\n  "claude-opus-4-5-20251101": 2.5,\n  "claude-opus-4-6": 2.5,\n  "claude-opus-4-6-high": 2.5,\n  "claude-opus-4-6-low": 2.5,\n  "claude-opus-4-6-max": 2.5,\n  "claude-opus-4-6-medium": 2.5,\n  "claude-opus-4-7": 2.5,\n  "claude-opus-4-7-high": 2.5,\n  "claude-opus-4-7-low": 2.5,\n  "claude-opus-4-7-max": 2.5,\n  "claude-opus-4-7-medium": 2.5,\n  "claude-opus-4-7-xhigh": 2.5,\n  "claude-sonnet-4-20250514": 1.5,\n  "claude-sonnet-4-5-20250929": 1.5,\n  "code-davinci-edit-001": 10,\n  "command": 0.5,\n  "command-light": 0.5,\n  "command-light-nightly": 0.5,\n  "command-nightly": 0.5,\n  "command-r": 0.25,\n  "command-r-08-2024": 0.075,\n  "command-r-plus": 1.5,\n  "command-r-plus-08-2024": 1.25,\n  "curie": 10,\n  "davinci": 10,\n  "davinci-002": 1,\n  "deepseek-ai/DeepSeek-R1": 0.8,\n  "deepseek-ai/DeepSeek-R1-0528": 0.8,\n  "deepseek-ai/DeepSeek-V3-0324": 0.8,\n  "deepseek-ai/DeepSeek-V3.1": 0.8,\n  "deepseek-chat": 0.135,\n  "deepseek-coder": 0.135,\n  "deepseek-reasoner": 0.275,\n  "embedding_s1_v1": 0.0715,\n  "embedding-bert-512-v1": 0.0715,\n  "Embedding-V1": 0.13698630137,\n  "ERNIE-3.5-4K-0205": 0.821917808219,\n  "ERNIE-3.5-8K": 0.821917808219,\n  "ERNIE-3.5-8K-0205": 1.643835616438,\n  "ERNIE-3.5-8K-1222": 0.821917808219,\n  "ERNIE-4.0-8K": 8.219178082192,\n  "ERNIE-Bot-8K": 1.643835616438,\n  "ERNIE-Lite-8K-0308": 0.205479452055,\n  "ERNIE-Lite-8K-0922": 0.547945205479,\n  "ERNIE-Speed-128K": 0.27397260274,\n  "ERNIE-Speed-8K": 0.27397260274,\n  "ERNIE-Tiny-8K": 0.068493150685,\n  "gemini-1.5-flash-latest": 0.075,\n  "gemini-1.5-pro-latest": 1.25,\n  "gemini-2.0-flash": 0.05,\n  "gemini-2.5-flash": 0.15,\n  "gemini-2.5-flash-lite-preview-06-17": 0.05,\n  "gemini-2.5-flash-lite-preview-thinking-*": 0.05,\n  "gemini-2.5-flash-preview-04-17": 0.075,\n  "gemini-2.5-flash-preview-04-17-nothinking": 0.075,\n  "gemini-2.5-flash-preview-04-17-thinking": 0.075,\n  "gemini-2.5-flash-preview-05-20": 0.075,\n  "gemini-2.5-flash-preview-05-20-nothinking": 0.075,\n  "gemini-2.5-flash-preview-05-20-thinking": 0.075,\n  "gemini-2.5-flash-thinking-*": 0.075,\n  "gemini-2.5-pro": 0.625,\n  "gemini-2.5-pro-exp-03-25": 0.625,\n  "gemini-2.5-pro-preview-03-25": 0.625,\n  "gemini-2.5-pro-thinking-*": 0.625,\n  "gemini-embedding-001": 0.075,\n  "gemini-robotics-er-1.5-preview": 0.15,\n  "glm-3-turbo": 0.3572,\n  "glm-4": 7.143,\n  "glm-4-0520": 6.849315068493,\n  "glm-4-air": 0.068493150685,\n  "glm-4-airx": 0.684931506849,\n  "glm-4-alltools": 6.849315068493,\n  "glm-4-flash": 0,\n  "glm-4-long": 0.068493150685,\n  "glm-4-plus": 3.424657534247,\n  "glm-4v": 3.424657534247,\n  "glm-4v-plus": 0.684931506849,\n  "gpt-3.5-turbo": 0.25,\n  "gpt-3.5-turbo-0125": 0.25,\n  "gpt-3.5-turbo-0613": 0.75,\n  "gpt-3.5-turbo-1106": 0.5,\n  "gpt-3.5-turbo-16k": 1.5,\n  "gpt-3.5-turbo-16k-0613": 1.5,\n  "gpt-3.5-turbo-instruct": 0.75,\n  "gpt-4": 15,\n  "gpt-4-0125-preview": 5,\n  "gpt-4-0613": 15,\n  "gpt-4-1106-preview": 5,\n  "gpt-4-1106-vision-preview": 5,\n  "gpt-4-32k": 30,\n  "gpt-4-32k-0613": 30,\n  "gpt-4-all": 15,\n  "gpt-4-turbo": 5,\n  "gpt-4-turbo-2024-04-09": 5,\n  "gpt-4-turbo-preview": 5,\n  "gpt-4-vision-preview": 5,\n  "gpt-4.1": 1,\n  "gpt-4.1-2025-04-14": 1,\n  "gpt-4.1-mini": 0.2,\n  "gpt-4.1-mini-2025-04-14": 0.2,\n  "gpt-4.1-nano": 0.05,\n  "gpt-4.1-nano-2025-04-14": 0.05,\n  "gpt-4.5-preview": 37.5,\n  "gpt-4.5-preview-2025-02-27": 37.5,\n  "gpt-4o": 1.25,\n  "gpt-4o-2024-05-13": 2.5,\n  "gpt-4o-2024-08-06": 1.25,\n  "gpt-4o-2024-11-20": 1.25,\n  "gpt-4o-all": 15,\n  "gpt-4o-audio-preview": 1.25,\n  "gpt-4o-audio-preview-2024-10-01": 1.25,\n  "gpt-4o-gizmo-*": 2.5,\n  "gpt-4o-mini": 0.075,\n  "gpt-4o-mini-2024-07-18": 0.075,\n  "gpt-4o-mini-realtime-preview": 0.3,\n  "gpt-4o-mini-realtime-preview-2024-12-17": 0.3,\n  "gpt-4o-realtime-preview": 2.5,\n  "gpt-4o-realtime-preview-2024-10-01": 2.5,\n  "gpt-4o-realtime-preview-2024-12-17": 2.5,\n  "gpt-5": 0.625,\n  "gpt-5-2025-08-07": 0.625,\n  "gpt-5-chat-latest": 0.625,\n  "gpt-5-mini": 0.125,\n  "gpt-5-mini-2025-08-07": 0.125,\n  "gpt-5-nano": 0.025,\n  "gpt-5-nano-2025-08-07": 0.025,\n  "gpt-image-1": 2.5,\n  "grok-2": 1,\n  "grok-2-vision": 1,\n  "grok-3-beta": 1.5,\n  "grok-3-fast-beta": 2.5,\n  "grok-3-mini-beta": 0.15,\n  "grok-3-mini-fast-beta": 0.3,\n  "grok-beta": 2.5,\n  "grok-vision-beta": 2.5,\n  "hunyuan": 7.143,\n  "llama-3-sonar-large-32k-chat": 0,\n  "llama-3-sonar-large-32k-online": 0,\n  "llama-3-sonar-small-32k-chat": 0.1,\n  "llama-3-sonar-small-32k-online": 0.1,\n  "NousResearch/Hermes-4-405B-FP8": 0.8,\n  "o1": 7.5,\n  "o1-2024-12-17": 7.5,\n  "o1-mini": 0.55,\n  "o1-mini-2024-09-12": 0.55,\n  "o1-preview": 7.5,\n  "o1-preview-2024-09-12": 7.5,\n  "o1-pro": 75,\n  "o1-pro-2025-03-19": 75,\n  "o3": 1,\n  "o3-2025-04-16": 1,\n  "o3-deep-research": 5,\n  "o3-deep-research-2025-06-26": 5,\n  "o3-mini": 0.55,\n  "o3-mini-2025-01-31": 0.55,\n  "o3-mini-2025-01-31-high": 0.55,\n  "o3-mini-2025-01-31-low": 0.55,\n  "o3-mini-2025-01-31-medium": 0.55,\n  "o3-mini-high": 0.55,\n  "o3-mini-low": 0.55,\n  "o3-mini-medium": 0.55,\n  "o3-pro": 10,\n  "o3-pro-2025-06-10": 10,\n  "o4-mini": 0.55,\n  "o4-mini-2025-04-16": 0.55,\n  "o4-mini-deep-research": 1,\n  "o4-mini-deep-research-2025-06-26": 1,\n  "openai/gpt-oss-120b": 0.5,\n  "PaLM-2": 1,\n  "qwen-plus": 10,\n  "qwen-turbo": 0.8572,\n  "Qwen/Qwen3-235B-A22B-Instruct-2507": 0.3,\n  "Qwen/Qwen3-235B-A22B-Thinking-2507": 0.6,\n  "Qwen/Qwen3-Coder-480B-A35B-Instruct-FP8": 0.8,\n  "semantic_similarity_s1_v1": 0.0715,\n  "SparkDesk-v1.1": 1.2858,\n  "SparkDesk-v2.1": 1.2858,\n  "SparkDesk-v3.1": 1.2858,\n  "SparkDesk-v3.5": 1.2858,\n  "SparkDesk-v4.0": 1.2858,\n  "tao-8k": 0.13698630137,\n  "text-ada-001": 0.2,\n  "text-babbage-001": 0.25,\n  "text-curie-001": 1,\n  "text-davinci-edit-001": 10,\n  "text-embedding-004": 0.001,\n  "text-embedding-3-large": 0.065,\n  "text-embedding-3-small": 0.01,\n  "text-embedding-ada-002": 0.05,\n  "text-embedding-v1": 0.05,\n  "text-moderation-latest": 0.1,\n  "text-moderation-stable": 0.1,\n  "text-search-ada-doc-001": 10,\n  "tts-1": 7.5,\n  "tts-1-1106": 7.5,\n  "tts-1-hd": 15,\n  "tts-1-hd-1106": 15,\n  "whisper-1": 15,\n  "yi-34b-chat-0205": 0.18,\n  "yi-34b-chat-200k": 0.864,\n  "yi-large": 1.369863013698,\n  "yi-large-preview": 1.369863013698,\n  "yi-large-rag": 1.712328767123,\n  "yi-large-rag-preview": 1.712328767123,\n  "yi-large-turbo": 0.821917808219,\n  "yi-medium": 0.171232876713,\n  "yi-medium-200k": 0.821917808219,\n  "yi-spark": 0.068493150685,\n  "yi-vision": 0.410958904109,\n  "yi-vl-plus": 0.432,\n  "zai-org/GLM-4.5-FP8": 0.8\n}
CompletionRatio	{\n  "gpt-4-all": 2,\n  "gpt-4o-gizmo-*": 3,\n  "gpt-image-1": 8\n}
CacheRatio	{\n  "claude-3-5-haiku-20241022": 0.1,\n  "claude-3-5-sonnet-20240620": 0.1,\n  "claude-3-5-sonnet-20241022": 0.1,\n  "claude-3-7-sonnet-20250219": 0.1,\n  "claude-3-7-sonnet-20250219-thinking": 0.1,\n  "claude-3-haiku-20240307": 0.1,\n  "claude-3-opus-20240229": 0.1,\n  "claude-3-sonnet-20240229": 0.1,\n  "claude-haiku-4-5-20251001": 0.1,\n  "claude-opus-4-1-20250805": 0.1,\n  "claude-opus-4-1-20250805-thinking": 0.1,\n  "claude-opus-4-20250514": 0.1,\n  "claude-opus-4-20250514-thinking": 0.1,\n  "claude-opus-4-5-20251101": 0.1,\n  "claude-opus-4-5-20251101-thinking": 0.1,\n  "claude-opus-4-6": 0.1,\n  "claude-opus-4-6-high": 0.1,\n  "claude-opus-4-6-low": 0.1,\n  "claude-opus-4-6-max": 0.1,\n  "claude-opus-4-6-medium": 0.1,\n  "claude-opus-4-6-thinking": 0.1,\n  "claude-opus-4-7": 0.1,\n  "claude-opus-4-7-high": 0.1,\n  "claude-opus-4-7-low": 0.1,\n  "claude-opus-4-7-max": 0.1,\n  "claude-opus-4-7-medium": 0.1,\n  "claude-opus-4-7-thinking": 0.1,\n  "claude-opus-4-7-xhigh": 0.1,\n  "claude-sonnet-4-20250514": 0.1,\n  "claude-sonnet-4-20250514-thinking": 0.1,\n  "claude-sonnet-4-5-20250929": 0.1,\n  "claude-sonnet-4-5-20250929-thinking": 0.1,\n  "deepseek-chat": 0.25,\n  "deepseek-coder": 0.25,\n  "deepseek-reasoner": 0.25,\n  "gemini-3-flash-preview": 0.1,\n  "gemini-3-pro-preview": 0.1,\n  "gemini-3.1-pro-preview": 0.1,\n  "gpt-4": 0.5,\n  "gpt-4.1": 0.25,\n  "gpt-4.1-mini": 0.25,\n  "gpt-4.1-nano": 0.25,\n  "gpt-4.5-preview": 0.5,\n  "gpt-4.5-preview-2025-02-27": 0.5,\n  "gpt-4o": 0.5,\n  "gpt-4o-2024-08-06": 0.5,\n  "gpt-4o-2024-11-20": 0.5,\n  "gpt-4o-mini": 0.5,\n  "gpt-4o-mini-2024-07-18": 0.5,\n  "gpt-4o-mini-realtime-preview": 0.5,\n  "gpt-4o-realtime-preview": 0.5,\n  "gpt-5": 0.1,\n  "gpt-5-2025-08-07": 0.1,\n  "gpt-5-chat-latest": 0.1,\n  "gpt-5-mini": 0.1,\n  "gpt-5-mini-2025-08-07": 0.1,\n  "gpt-5-nano": 0.1,\n  "gpt-5-nano-2025-08-07": 0.1,\n  "o1": 0.5,\n  "o1-2024-12-17": 0.5,\n  "o1-mini": 0.5,\n  "o1-mini-2024-09-12": 0.5,\n  "o1-preview": 0.5,\n  "o1-preview-2024-09-12": 0.5,\n  "o3-mini": 0.5,\n  "o3-mini-2025-01-31": 0.5\n}
CreateCacheRatio	{\n  "claude-3-5-haiku-20241022": 1.25,\n  "claude-3-5-sonnet-20240620": 1.25,\n  "claude-3-5-sonnet-20241022": 1.25,\n  "claude-3-7-sonnet-20250219": 1.25,\n  "claude-3-7-sonnet-20250219-thinking": 1.25,\n  "claude-3-haiku-20240307": 1.25,\n  "claude-3-opus-20240229": 1.25,\n  "claude-3-sonnet-20240229": 1.25,\n  "claude-haiku-4-5-20251001": 1.25,\n  "claude-opus-4-1-20250805": 1.25,\n  "claude-opus-4-1-20250805-thinking": 1.25,\n  "claude-opus-4-20250514": 1.25,\n  "claude-opus-4-20250514-thinking": 1.25,\n  "claude-opus-4-5-20251101": 1.25,\n  "claude-opus-4-5-20251101-thinking": 1.25,\n  "claude-opus-4-6": 1.25,\n  "claude-opus-4-6-high": 1.25,\n  "claude-opus-4-6-low": 1.25,\n  "claude-opus-4-6-max": 1.25,\n  "claude-opus-4-6-medium": 1.25,\n  "claude-opus-4-6-thinking": 1.25,\n  "claude-opus-4-7": 1.25,\n  "claude-opus-4-7-high": 1.25,\n  "claude-opus-4-7-low": 1.25,\n  "claude-opus-4-7-max": 1.25,\n  "claude-opus-4-7-medium": 1.25,\n  "claude-opus-4-7-thinking": 1.25,\n  "claude-opus-4-7-xhigh": 1.25,\n  "claude-sonnet-4-20250514": 1.25,\n  "claude-sonnet-4-20250514-thinking": 1.25,\n  "claude-sonnet-4-5-20250929": 1.25,\n  "claude-sonnet-4-5-20250929-thinking": 1.25\n}
ImageRatio	{\n  "gpt-image-1": 2\n}
AudioRatio	{\n  "gpt-4o-audio-preview": 16,\n  "gpt-4o-mini-audio-preview": 66.67,\n  "gpt-4o-mini-realtime-preview": 16.67,\n  "gpt-4o-realtime-preview": 8\n}
AudioCompletionRatio	{\n  "gpt-4o-mini-realtime": 2,\n  "gpt-4o-realtime": 2\n}
billing_setting.billing_mode	{}
billing_setting.billing_expr	{}
oidc.client_id	newapi
oidc.client_secret	Newapi1.
oidc.authorization_endpoint	http://localhost:3001
oidc.enabled	true
oidc.well_known	<nil>
oidc.token_endpoint	http://host.docker.internal:3000/oauth/oidc/token
oidc.user_info_endpoint	http://host.docker.internal:3000/oauth/userinfo
\.


--
-- Data for Name: passkey_credentials; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.passkey_credentials (id, user_id, credential_id, public_key, attestation_type, aa_guid, sign_count, clone_warning, user_present, user_verified, backup_eligible, backup_state, transports, attachment, last_used_at, created_at, updated_at, deleted_at) FROM stdin;
\.


--
-- Data for Name: perf_metrics; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.perf_metrics (id, model_name, "group", bucket_ts, request_count, success_count, total_latency_ms, ttft_sum_ms, ttft_count, output_tokens, generation_ms) FROM stdin;
1	gpt-5.4	default	1779422400	1	1	2609	1758	1	13	850
2	gpt-image-2	default	1779422400	1	0	1242	0	0	0	0
3	gpt-image-2	default	1779426000	5	3	393364	0	0	10536	388223
4	gpt-image-2	default	1779429600	2	1	154173	0	0	1756	87865
5	gpt-image-2	default	1779433200	10	9	854315	0	0	16162	840097
6	gpt-image-2	default	1779436800	5	4	521765	0	0	6511	432762
7	gpt-image-2	default	1779444000	6	5	927405	0	0	24584	842131
8	gpt-image-2	default	1779462000	5	4	1204506	0	0	22828	843829
9	gpt-image-2	default	1779465600	7	7	816963	0	0	17560	816963
10	gpt-image-2	default	1779530400	4	3	982602	0	0	7956	381911
11	gpt-image-2	default	1779534000	6	5	464027	0	0	3355	367284
12	gpt-image-2	default	1779537600	11	6	663435	0	0	4596	550274
13	gpt-image-2	default	1779541200	5	3	980648	0	0	2013	311915
14	gpt-image-2	default	1779544800	3	3	282839	0	0	2013	282839
\.


--
-- Data for Name: prefill_groups; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.prefill_groups (id, name, type, items, description, created_time, updated_time, deleted_at) FROM stdin;
\.


--
-- Data for Name: quota_data; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.quota_data (id, user_id, username, model_name, created_at, token_used, count, quota) FROM stdin;
2	2	sso-admin	gpt-5.4	1779422400	31	1	72
3	1	xiaopihong	gpt-image-2	1779426000	237	1	250000
4	1	xiaopihong	gpt-5.4	1779426000	29	1	63
1	1	xiaopihong	gpt-5.4	1779422400	58	2	126
5	2	sso-admin	gpt-image-2	1779426000	10981	3	750000
6	2	sso-admin	gpt-image-2	1779429600	2844	1	250000
7	2	sso-admin	gpt-image-2	1779433200	26282	9	2250000
8	2	sso-admin	gpt-image-2	1779436800	11194	4	1000000
9	2	sso-admin	gpt-image-2	1779444000	30859	5	1250000
10	2	sso-admin	gpt-image-2	1779462000	27773	4	1000000
11	2	sso-admin	gpt-image-2	1779465600	26441	7	1750000
12	2	sso-admin	gpt-image-2	1779530400	11659	3	750000
13	2	sso-admin	gpt-image-2	1779534000	8115	5	1250000
14	2	sso-admin	gpt-image-2	1779537600	10515	6	1500000
15	2	sso-admin	gpt-image-2	1779541200	4923	3	750000
16	2	sso-admin	gpt-image-2	1779544800	4862	3	750000
17	2	sso-admin	gpt-image-2	1779548400	3330	2	500000
\.


--
-- Data for Name: redemptions; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.redemptions (id, user_id, key, status, name, quota, created_time, redeemed_time, used_user_id, deleted_at, expired_time) FROM stdin;
\.


--
-- Data for Name: setups; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.setups (id, version, initialized_at) FROM stdin;
1		1779422084
\.


--
-- Data for Name: subscription_orders; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.subscription_orders (id, user_id, plan_id, money, trade_no, payment_method, payment_provider, status, create_time, complete_time, provider_payload) FROM stdin;
\.


--
-- Data for Name: subscription_plans; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.subscription_plans (id, title, subtitle, price_amount, currency, duration_unit, duration_value, custom_seconds, enabled, sort_order, stripe_price_id, creem_product_id, max_purchase_per_user, upgrade_group, total_amount, quota_reset_period, quota_reset_custom_seconds, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: subscription_pre_consume_records; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.subscription_pre_consume_records (id, request_id, user_id, user_subscription_id, pre_consumed, status, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: tasks; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.tasks (id, created_at, updated_at, task_id, platform, user_id, "group", channel_id, quota, action, status, fail_reason, submit_time, start_time, finish_time, progress, properties, private_data, data) FROM stdin;
\.


--
-- Data for Name: tokens; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.tokens (id, user_id, key, status, name, created_time, accessed_time, expired_time, remain_quota, unlimited_quota, model_limits_enabled, model_limits, allow_ips, used_quota, "group", cross_group_retry, deleted_at) FROM stdin;
3	2	fVbBVEnWg5jdHX1kfAcOD4Br3uo5RFVFk5jo9quGoIzAupDs	1	测试	1779504734	1779504734	-1	0	t	f			0		f	2026-05-23 03:29:12.588237+00
4	2	WZmKtrzwB2aLDK10MVd0aLgN8CVV0XLYECGcpbNPHbgkR3Up	1	默认	1779506955	1779550526	-1	-5500000	t	f			5500000		f	\N
1	2	IAyIeafIBwjqeHcCEYNXjZVe93SCl1wo37OxXAsKSBGepFpL	1	自用	1779425231	1779468674	-1	-8250072	t	f			8250072		f	2026-05-23 02:36:44.895133+00
2	2	Dly7NbE5AHT1kkGW2Pf24UMbWigMCErbeUd1nVkGSxCCFxJ6	1	自用	1779503814	1779503814	-1	0	t	f			0		f	2026-05-23 02:52:23.781177+00
\.


--
-- Data for Name: top_ups; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.top_ups (id, user_id, amount, money, trade_no, payment_method, payment_provider, create_time, complete_time, status) FROM stdin;
\.


--
-- Data for Name: two_fa_backup_codes; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.two_fa_backup_codes (id, user_id, code_hash, is_used, used_at, created_at, deleted_at) FROM stdin;
\.


--
-- Data for Name: two_fas; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.two_fas (id, user_id, secret, is_enabled, failed_attempts, locked_until, last_used_at, created_at, updated_at, deleted_at) FROM stdin;
\.


--
-- Data for Name: user_oauth_bindings; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.user_oauth_bindings (id, user_id, provider_id, provider_user_id, created_at) FROM stdin;
\.


--
-- Data for Name: user_subscriptions; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.user_subscriptions (id, user_id, plan_id, amount_total, amount_used, start_time, end_time, status, source, last_reset_time, next_reset_time, upgrade_group, prev_user_group, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.users (id, username, password, display_name, role, status, email, github_id, discord_id, oidc_id, wechat_id, telegram_id, access_token, quota, used_quota, request_count, "group", aff_code, aff_count, aff_quota, aff_history, inviter_id, deleted_at, linux_do_id, setting, remark, stripe_customer, created_at, last_login_at) FROM stdin;
1	xiaopihong	$2a$10$w.mSTu2Yxr02DQEJWr.dp.1o3MNjOQAWIyf568/v3XZ984xlizG52	Root User	100	1							Z7Gh1c2hrvUrar9M7AD+607PjAQzNg==	100000000	0	0	default		0	0	0	0	\N					1779422084	1779536248
3	xph-admin	$2a$10$eYkyI3yo0vBrPAkV0yPx8OH0Yk1JmpuegoQH0R2Jz3dbfHfhGSqLm		1	1	811258683@qq.com			2			0ac62a322d06508e36e7897142ea189d	50000000	0	0	default	IB83	0	0	0	0	\N		{"gotify_priority":0,"sidebar_modules":"{\\"chat\\":{\\"chat\\":true,\\"enabled\\":true,\\"playground\\":true},\\"console\\":{\\"detail\\":true,\\"enabled\\":true,\\"log\\":true,\\"midjourney\\":true,\\"task\\":true,\\"token\\":true},\\"personal\\":{\\"enabled\\":true,\\"personal\\":true,\\"topup\\":true}}"}			1779425248	0
2	sso-admin	$2a$10$bBKOR3DUPyPBiWtQkLXbiObUOGNEUWrSavgaE7OqGuBUYLlhJBY4W		1	1	811258682@qq.com			1			acf8d17ef1fbe54eeec437826c84679a	46249928	13750072	56	default	jJJI	0	0	0	0	\N		{"gotify_priority":0,"sidebar_modules":"{\\"chat\\":{\\"chat\\":true,\\"enabled\\":true,\\"playground\\":true},\\"console\\":{\\"detail\\":true,\\"enabled\\":true,\\"log\\":true,\\"midjourney\\":true,\\"task\\":true,\\"token\\":true},\\"personal\\":{\\"enabled\\":true,\\"personal\\":true,\\"topup\\":true}}"}			1779425212	1779506912
\.


--
-- Data for Name: vendors; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.vendors (id, name, description, icon, status, created_time, updated_time, deleted_at) FROM stdin;
1	OpenAI		OpenAI	1	1779425103	1779425103	\N
\.


--
-- Name: channels_id_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.channels_id_seq', 1, true);


--
-- Name: checkins_id_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.checkins_id_seq', 1, false);


--
-- Name: custom_oauth_providers_id_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.custom_oauth_providers_id_seq', 1, false);


--
-- Name: logs_id_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.logs_id_seq', 80, true);


--
-- Name: midjourneys_id_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.midjourneys_id_seq', 1, false);


--
-- Name: models_id_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.models_id_seq', 2, true);


--
-- Name: passkey_credentials_id_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.passkey_credentials_id_seq', 1, false);


--
-- Name: perf_metrics_id_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.perf_metrics_id_seq', 14, true);


--
-- Name: prefill_groups_id_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.prefill_groups_id_seq', 1, false);


--
-- Name: quota_data_id_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.quota_data_id_seq', 17, true);


--
-- Name: redemptions_id_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.redemptions_id_seq', 1, false);


--
-- Name: setups_id_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.setups_id_seq', 1, true);


--
-- Name: subscription_orders_id_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.subscription_orders_id_seq', 1, false);


--
-- Name: subscription_plans_id_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.subscription_plans_id_seq', 1, false);


--
-- Name: subscription_pre_consume_records_id_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.subscription_pre_consume_records_id_seq', 1, false);


--
-- Name: tasks_id_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.tasks_id_seq', 1, false);


--
-- Name: tokens_id_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.tokens_id_seq', 4, true);


--
-- Name: top_ups_id_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.top_ups_id_seq', 1, false);


--
-- Name: two_fa_backup_codes_id_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.two_fa_backup_codes_id_seq', 1, false);


--
-- Name: two_fas_id_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.two_fas_id_seq', 1, false);


--
-- Name: user_oauth_bindings_id_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.user_oauth_bindings_id_seq', 1, false);


--
-- Name: user_subscriptions_id_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.user_subscriptions_id_seq', 1, false);


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.users_id_seq', 3, true);


--
-- Name: vendors_id_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.vendors_id_seq', 1, true);


--
-- Name: abilities abilities_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.abilities
    ADD CONSTRAINT abilities_pkey PRIMARY KEY ("group", model, channel_id);


--
-- Name: channels channels_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.channels
    ADD CONSTRAINT channels_pkey PRIMARY KEY (id);


--
-- Name: checkins checkins_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.checkins
    ADD CONSTRAINT checkins_pkey PRIMARY KEY (id);


--
-- Name: custom_oauth_providers custom_oauth_providers_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.custom_oauth_providers
    ADD CONSTRAINT custom_oauth_providers_pkey PRIMARY KEY (id);


--
-- Name: prefill_groups idx_prefill_groups_name; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.prefill_groups
    ADD CONSTRAINT idx_prefill_groups_name UNIQUE (name);


--
-- Name: logs logs_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.logs
    ADD CONSTRAINT logs_pkey PRIMARY KEY (id);


--
-- Name: midjourneys midjourneys_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.midjourneys
    ADD CONSTRAINT midjourneys_pkey PRIMARY KEY (id);


--
-- Name: models models_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.models
    ADD CONSTRAINT models_pkey PRIMARY KEY (id);


--
-- Name: options options_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.options
    ADD CONSTRAINT options_pkey PRIMARY KEY (key);


--
-- Name: passkey_credentials passkey_credentials_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.passkey_credentials
    ADD CONSTRAINT passkey_credentials_pkey PRIMARY KEY (id);


--
-- Name: perf_metrics perf_metrics_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.perf_metrics
    ADD CONSTRAINT perf_metrics_pkey PRIMARY KEY (id);


--
-- Name: prefill_groups prefill_groups_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.prefill_groups
    ADD CONSTRAINT prefill_groups_pkey PRIMARY KEY (id);


--
-- Name: quota_data quota_data_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.quota_data
    ADD CONSTRAINT quota_data_pkey PRIMARY KEY (id);


--
-- Name: redemptions redemptions_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.redemptions
    ADD CONSTRAINT redemptions_pkey PRIMARY KEY (id);


--
-- Name: setups setups_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.setups
    ADD CONSTRAINT setups_pkey PRIMARY KEY (id);


--
-- Name: subscription_orders subscription_orders_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.subscription_orders
    ADD CONSTRAINT subscription_orders_pkey PRIMARY KEY (id);


--
-- Name: subscription_orders subscription_orders_trade_no_key; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.subscription_orders
    ADD CONSTRAINT subscription_orders_trade_no_key UNIQUE (trade_no);


--
-- Name: subscription_plans subscription_plans_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.subscription_plans
    ADD CONSTRAINT subscription_plans_pkey PRIMARY KEY (id);


--
-- Name: subscription_pre_consume_records subscription_pre_consume_records_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.subscription_pre_consume_records
    ADD CONSTRAINT subscription_pre_consume_records_pkey PRIMARY KEY (id);


--
-- Name: tasks tasks_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.tasks
    ADD CONSTRAINT tasks_pkey PRIMARY KEY (id);


--
-- Name: tokens tokens_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.tokens
    ADD CONSTRAINT tokens_pkey PRIMARY KEY (id);


--
-- Name: top_ups top_ups_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.top_ups
    ADD CONSTRAINT top_ups_pkey PRIMARY KEY (id);


--
-- Name: top_ups top_ups_trade_no_key; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.top_ups
    ADD CONSTRAINT top_ups_trade_no_key UNIQUE (trade_no);


--
-- Name: two_fa_backup_codes two_fa_backup_codes_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.two_fa_backup_codes
    ADD CONSTRAINT two_fa_backup_codes_pkey PRIMARY KEY (id);


--
-- Name: two_fas two_fas_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.two_fas
    ADD CONSTRAINT two_fas_pkey PRIMARY KEY (id);


--
-- Name: two_fas two_fas_user_id_key; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.two_fas
    ADD CONSTRAINT two_fas_user_id_key UNIQUE (user_id);


--
-- Name: user_oauth_bindings user_oauth_bindings_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.user_oauth_bindings
    ADD CONSTRAINT user_oauth_bindings_pkey PRIMARY KEY (id);


--
-- Name: user_subscriptions user_subscriptions_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.user_subscriptions
    ADD CONSTRAINT user_subscriptions_pkey PRIMARY KEY (id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: users users_username_key; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_username_key UNIQUE (username);


--
-- Name: vendors vendors_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.vendors
    ADD CONSTRAINT vendors_pkey PRIMARY KEY (id);


--
-- Name: idx_abilities_channel_id; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_abilities_channel_id ON public.abilities USING btree (channel_id);


--
-- Name: idx_abilities_priority; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_abilities_priority ON public.abilities USING btree (priority);


--
-- Name: idx_abilities_tag; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_abilities_tag ON public.abilities USING btree (tag);


--
-- Name: idx_abilities_weight; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_abilities_weight ON public.abilities USING btree (weight);


--
-- Name: idx_channels_name; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_channels_name ON public.channels USING btree (name);


--
-- Name: idx_channels_tag; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_channels_tag ON public.channels USING btree (tag);


--
-- Name: idx_created_at_id; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_created_at_id ON public.logs USING btree (id, created_at);


--
-- Name: idx_created_at_type; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_created_at_type ON public.logs USING btree (created_at, type);


--
-- Name: idx_custom_oauth_providers_slug; Type: INDEX; Schema: public; Owner: root
--

CREATE UNIQUE INDEX idx_custom_oauth_providers_slug ON public.custom_oauth_providers USING btree (slug);


--
-- Name: idx_logs_channel_id; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_logs_channel_id ON public.logs USING btree (channel_id);


--
-- Name: idx_logs_group; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_logs_group ON public.logs USING btree ("group");


--
-- Name: idx_logs_ip; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_logs_ip ON public.logs USING btree (ip);


--
-- Name: idx_logs_model_name; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_logs_model_name ON public.logs USING btree (model_name);


--
-- Name: idx_logs_request_id; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_logs_request_id ON public.logs USING btree (request_id);


--
-- Name: idx_logs_token_id; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_logs_token_id ON public.logs USING btree (token_id);


--
-- Name: idx_logs_token_name; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_logs_token_name ON public.logs USING btree (token_name);


--
-- Name: idx_logs_upstream_request_id; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_logs_upstream_request_id ON public.logs USING btree (upstream_request_id);


--
-- Name: idx_logs_user_id; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_logs_user_id ON public.logs USING btree (user_id);


--
-- Name: idx_logs_username; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_logs_username ON public.logs USING btree (username);


--
-- Name: idx_midjourneys_action; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_midjourneys_action ON public.midjourneys USING btree (action);


--
-- Name: idx_midjourneys_finish_time; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_midjourneys_finish_time ON public.midjourneys USING btree (finish_time);


--
-- Name: idx_midjourneys_mj_id; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_midjourneys_mj_id ON public.midjourneys USING btree (mj_id);


--
-- Name: idx_midjourneys_progress; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_midjourneys_progress ON public.midjourneys USING btree (progress);


--
-- Name: idx_midjourneys_start_time; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_midjourneys_start_time ON public.midjourneys USING btree (start_time);


--
-- Name: idx_midjourneys_status; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_midjourneys_status ON public.midjourneys USING btree (status);


--
-- Name: idx_midjourneys_submit_time; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_midjourneys_submit_time ON public.midjourneys USING btree (submit_time);


--
-- Name: idx_midjourneys_user_id; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_midjourneys_user_id ON public.midjourneys USING btree (user_id);


--
-- Name: idx_models_deleted_at; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_models_deleted_at ON public.models USING btree (deleted_at);


--
-- Name: idx_models_vendor_id; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_models_vendor_id ON public.models USING btree (vendor_id);


--
-- Name: idx_passkey_credentials_credential_id; Type: INDEX; Schema: public; Owner: root
--

CREATE UNIQUE INDEX idx_passkey_credentials_credential_id ON public.passkey_credentials USING btree (credential_id);


--
-- Name: idx_passkey_credentials_deleted_at; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_passkey_credentials_deleted_at ON public.passkey_credentials USING btree (deleted_at);


--
-- Name: idx_passkey_credentials_user_id; Type: INDEX; Schema: public; Owner: root
--

CREATE UNIQUE INDEX idx_passkey_credentials_user_id ON public.passkey_credentials USING btree (user_id);


--
-- Name: idx_perf_bucket_ts; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_perf_bucket_ts ON public.perf_metrics USING btree (bucket_ts);


--
-- Name: idx_perf_model_group_bucket; Type: INDEX; Schema: public; Owner: root
--

CREATE UNIQUE INDEX idx_perf_model_group_bucket ON public.perf_metrics USING btree (model_name, "group", bucket_ts);


--
-- Name: idx_prefill_groups_deleted_at; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_prefill_groups_deleted_at ON public.prefill_groups USING btree (deleted_at);


--
-- Name: idx_prefill_groups_type; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_prefill_groups_type ON public.prefill_groups USING btree (type);


--
-- Name: idx_qdt_created_at; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_qdt_created_at ON public.quota_data USING btree (created_at);


--
-- Name: idx_qdt_model_user_name; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_qdt_model_user_name ON public.quota_data USING btree (model_name, username);


--
-- Name: idx_quota_data_user_id; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_quota_data_user_id ON public.quota_data USING btree (user_id);


--
-- Name: idx_redemptions_deleted_at; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_redemptions_deleted_at ON public.redemptions USING btree (deleted_at);


--
-- Name: idx_redemptions_key; Type: INDEX; Schema: public; Owner: root
--

CREATE UNIQUE INDEX idx_redemptions_key ON public.redemptions USING btree (key);


--
-- Name: idx_redemptions_name; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_redemptions_name ON public.redemptions USING btree (name);


--
-- Name: idx_subscription_orders_plan_id; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_subscription_orders_plan_id ON public.subscription_orders USING btree (plan_id);


--
-- Name: idx_subscription_orders_trade_no; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_subscription_orders_trade_no ON public.subscription_orders USING btree (trade_no);


--
-- Name: idx_subscription_orders_user_id; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_subscription_orders_user_id ON public.subscription_orders USING btree (user_id);


--
-- Name: idx_subscription_pre_consume_records_request_id; Type: INDEX; Schema: public; Owner: root
--

CREATE UNIQUE INDEX idx_subscription_pre_consume_records_request_id ON public.subscription_pre_consume_records USING btree (request_id);


--
-- Name: idx_subscription_pre_consume_records_status; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_subscription_pre_consume_records_status ON public.subscription_pre_consume_records USING btree (status);


--
-- Name: idx_subscription_pre_consume_records_updated_at; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_subscription_pre_consume_records_updated_at ON public.subscription_pre_consume_records USING btree (updated_at);


--
-- Name: idx_subscription_pre_consume_records_user_id; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_subscription_pre_consume_records_user_id ON public.subscription_pre_consume_records USING btree (user_id);


--
-- Name: idx_subscription_pre_consume_records_user_subscription_id; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_subscription_pre_consume_records_user_subscription_id ON public.subscription_pre_consume_records USING btree (user_subscription_id);


--
-- Name: idx_tasks_action; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_tasks_action ON public.tasks USING btree (action);


--
-- Name: idx_tasks_channel_id; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_tasks_channel_id ON public.tasks USING btree (channel_id);


--
-- Name: idx_tasks_created_at; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_tasks_created_at ON public.tasks USING btree (created_at);


--
-- Name: idx_tasks_finish_time; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_tasks_finish_time ON public.tasks USING btree (finish_time);


--
-- Name: idx_tasks_platform; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_tasks_platform ON public.tasks USING btree (platform);


--
-- Name: idx_tasks_progress; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_tasks_progress ON public.tasks USING btree (progress);


--
-- Name: idx_tasks_start_time; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_tasks_start_time ON public.tasks USING btree (start_time);


--
-- Name: idx_tasks_status; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_tasks_status ON public.tasks USING btree (status);


--
-- Name: idx_tasks_submit_time; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_tasks_submit_time ON public.tasks USING btree (submit_time);


--
-- Name: idx_tasks_task_id; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_tasks_task_id ON public.tasks USING btree (task_id);


--
-- Name: idx_tasks_user_id; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_tasks_user_id ON public.tasks USING btree (user_id);


--
-- Name: idx_tokens_deleted_at; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_tokens_deleted_at ON public.tokens USING btree (deleted_at);


--
-- Name: idx_tokens_key; Type: INDEX; Schema: public; Owner: root
--

CREATE UNIQUE INDEX idx_tokens_key ON public.tokens USING btree (key);


--
-- Name: idx_tokens_name; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_tokens_name ON public.tokens USING btree (name);


--
-- Name: idx_tokens_user_id; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_tokens_user_id ON public.tokens USING btree (user_id);


--
-- Name: idx_top_ups_trade_no; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_top_ups_trade_no ON public.top_ups USING btree (trade_no);


--
-- Name: idx_top_ups_user_id; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_top_ups_user_id ON public.top_ups USING btree (user_id);


--
-- Name: idx_two_fa_backup_codes_deleted_at; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_two_fa_backup_codes_deleted_at ON public.two_fa_backup_codes USING btree (deleted_at);


--
-- Name: idx_two_fa_backup_codes_user_id; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_two_fa_backup_codes_user_id ON public.two_fa_backup_codes USING btree (user_id);


--
-- Name: idx_two_fas_deleted_at; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_two_fas_deleted_at ON public.two_fas USING btree (deleted_at);


--
-- Name: idx_two_fas_user_id; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_two_fas_user_id ON public.two_fas USING btree (user_id);


--
-- Name: idx_user_checkin_date; Type: INDEX; Schema: public; Owner: root
--

CREATE UNIQUE INDEX idx_user_checkin_date ON public.checkins USING btree (user_id, checkin_date);


--
-- Name: idx_user_id_id; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_user_id_id ON public.logs USING btree (user_id, id);


--
-- Name: idx_user_sub_active; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_user_sub_active ON public.user_subscriptions USING btree (user_id, status, end_time);


--
-- Name: idx_user_subscriptions_end_time; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_user_subscriptions_end_time ON public.user_subscriptions USING btree (end_time);


--
-- Name: idx_user_subscriptions_next_reset_time; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_user_subscriptions_next_reset_time ON public.user_subscriptions USING btree (next_reset_time);


--
-- Name: idx_user_subscriptions_plan_id; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_user_subscriptions_plan_id ON public.user_subscriptions USING btree (plan_id);


--
-- Name: idx_user_subscriptions_status; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_user_subscriptions_status ON public.user_subscriptions USING btree (status);


--
-- Name: idx_user_subscriptions_user_id; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_user_subscriptions_user_id ON public.user_subscriptions USING btree (user_id);


--
-- Name: idx_users_access_token; Type: INDEX; Schema: public; Owner: root
--

CREATE UNIQUE INDEX idx_users_access_token ON public.users USING btree (access_token);


--
-- Name: idx_users_aff_code; Type: INDEX; Schema: public; Owner: root
--

CREATE UNIQUE INDEX idx_users_aff_code ON public.users USING btree (aff_code);


--
-- Name: idx_users_deleted_at; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_users_deleted_at ON public.users USING btree (deleted_at);


--
-- Name: idx_users_discord_id; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_users_discord_id ON public.users USING btree (discord_id);


--
-- Name: idx_users_display_name; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_users_display_name ON public.users USING btree (display_name);


--
-- Name: idx_users_email; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_users_email ON public.users USING btree (email);


--
-- Name: idx_users_git_hub_id; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_users_git_hub_id ON public.users USING btree (github_id);


--
-- Name: idx_users_inviter_id; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_users_inviter_id ON public.users USING btree (inviter_id);


--
-- Name: idx_users_linux_do_id; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_users_linux_do_id ON public.users USING btree (linux_do_id);


--
-- Name: idx_users_oidc_id; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_users_oidc_id ON public.users USING btree (oidc_id);


--
-- Name: idx_users_stripe_customer; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_users_stripe_customer ON public.users USING btree (stripe_customer);


--
-- Name: idx_users_telegram_id; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_users_telegram_id ON public.users USING btree (telegram_id);


--
-- Name: idx_users_username; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_users_username ON public.users USING btree (username);


--
-- Name: idx_users_we_chat_id; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_users_we_chat_id ON public.users USING btree (wechat_id);


--
-- Name: idx_vendors_deleted_at; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_vendors_deleted_at ON public.vendors USING btree (deleted_at);


--
-- Name: index_username_model_name; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX index_username_model_name ON public.logs USING btree (model_name, username);


--
-- Name: uk_model_name_delete_at; Type: INDEX; Schema: public; Owner: root
--

CREATE UNIQUE INDEX uk_model_name_delete_at ON public.models USING btree (model_name, deleted_at);


--
-- Name: uk_prefill_name; Type: INDEX; Schema: public; Owner: root
--

CREATE UNIQUE INDEX uk_prefill_name ON public.prefill_groups USING btree (name) WHERE (deleted_at IS NULL);


--
-- Name: uk_vendor_name_delete_at; Type: INDEX; Schema: public; Owner: root
--

CREATE UNIQUE INDEX uk_vendor_name_delete_at ON public.vendors USING btree (name, deleted_at);


--
-- Name: ux_provider_userid; Type: INDEX; Schema: public; Owner: root
--

CREATE UNIQUE INDEX ux_provider_userid ON public.user_oauth_bindings USING btree (provider_id, provider_user_id);


--
-- Name: ux_user_provider; Type: INDEX; Schema: public; Owner: root
--

CREATE UNIQUE INDEX ux_user_provider ON public.user_oauth_bindings USING btree (user_id, provider_id);


--
-- PostgreSQL database dump complete
--

\unrestrict ZLPAwQe9s3ipeSYK6y9cCU1vzUKsCkVeI4whYRFOzd1WN9h2mCgy8oPVs7OCd98

