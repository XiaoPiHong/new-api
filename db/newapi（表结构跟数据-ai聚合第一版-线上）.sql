--
-- PostgreSQL database dump
--

\restrict pDsTkRoJ5UXhBiHvEfbQDC5GMpxTpEmF0ucff7LLvqFEhrsgLh36zJLjiXKMLko

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
default	gpt-image-2	5	t	0	0	
default	gpt-image-2_1k	7	t	0	0	
default	mimo-v2.5	2	t	0	0	
default	mimo-v2.5-pro	2	t	0	0	
default	gpt-5.4	1	t	0	0	
default	gpt-image-2	1	t	0	0	
default	gpt-5.5	1	t	0	0	
\.


--
-- Data for Name: channels; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.channels (id, type, key, open_ai_organization, test_model, status, name, weight, created_time, test_time, response_time, base_url, other, balance, balance_updated_time, models, "group", used_quota, model_mapping, status_code_mapping, priority, auto_ban, other_info, tag, setting, param_override, header_override, remark, channel_info, settings) FROM stdin;
5	1	sk-SPZbgGO5AtLFwVH2gYOhON2eGYetW_vR	\N		1	deepark公益	0	1779763592	1779767134	51506	https://image.deepark.tech		0	0	gpt-image-2	default	750000			0	1			{"force_format":false,"thinking_to_content":false,"proxy":"","pass_through_body_enabled":false,"system_prompt":"","system_prompt_override":false}		\N	\N	{"is_multi_key":false,"multi_key_size":0,"multi_key_status_list":null,"multi_key_polling_index":0,"multi_key_mode":"random"}	{"allow_service_tier":false,"disable_store":false,"allow_safety_identifier":false,"allow_include_obfuscation":false,"upstream_model_update_check_enabled":false,"upstream_model_update_auto_sync_enabled":false,"upstream_model_update_ignored_models":[],"upstream_model_update_last_detected_models":[],"upstream_model_update_last_check_time":0}
2	1	sk-clp1s15ujwjjb0n0w4u1rww0mxznribl4ttnh7ptxo7sxd9h	\N		1	mimo	0	1779635149	1779701933	1266	https://api.xiaomimimo.com		0	0	mimo-v2.5,mimo-v2.5-pro	default	1215			0	1			{"force_format":false,"thinking_to_content":false,"proxy":"","pass_through_body_enabled":false,"system_prompt":"","system_prompt_override":false}		\N	\N	{"is_multi_key":false,"multi_key_size":0,"multi_key_status_list":null,"multi_key_polling_index":0,"multi_key_mode":"random"}	{"allow_service_tier":false,"disable_store":false,"allow_safety_identifier":false,"allow_include_obfuscation":false,"upstream_model_update_check_enabled":false,"upstream_model_update_auto_sync_enabled":false,"upstream_model_update_ignored_models":[],"upstream_model_update_last_detected_models":[],"upstream_model_update_last_check_time":0}
1	1	sk-e1ed5d3323f9a9495aa726007241b9a61b874eb7eb72398c56eab3e4bbff859f			1	openai	0	1779290362	1779701949	2272	http://119.29.249.17:8080		0	0	gpt-5.4,gpt-image-2,gpt-5.5	default	5565060			0	1			{"force_format":false,"thinking_to_content":false,"proxy":"","pass_through_body_enabled":false,"system_prompt":"","system_prompt_override":false}		\N	\N	{"is_multi_key":false,"multi_key_size":0,"multi_key_status_list":null,"multi_key_polling_index":0,"multi_key_mode":"random"}	{"allow_service_tier":false,"disable_store":false,"allow_safety_identifier":false,"allow_include_obfuscation":false,"upstream_model_update_check_enabled":false,"upstream_model_update_auto_sync_enabled":false,"upstream_model_update_ignored_models":[],"upstream_model_update_last_detected_models":[],"upstream_model_update_last_check_time":0}
7	1	sk-mlP4jr9NeNazSFbI1ClUZ1rZPZPd0dpoYqfXFX7iI5kV5OP4			1	abrdns公益	0	1779839706	1779844522	56487	https://new-api.abrdns.com		0	0	gpt-image-2_1k	default	3000000			0	1			{"force_format":false,"thinking_to_content":false,"proxy":"","pass_through_body_enabled":false,"system_prompt":"","system_prompt_override":false}		\N	\N	{"is_multi_key":false,"multi_key_size":0,"multi_key_status_list":null,"multi_key_polling_index":0,"multi_key_mode":"random"}	{"allow_service_tier":false,"disable_store":false,"allow_safety_identifier":false,"allow_include_obfuscation":false,"upstream_model_update_check_enabled":false,"upstream_model_update_auto_sync_enabled":false,"upstream_model_update_ignored_models":[],"upstream_model_update_last_detected_models":[],"upstream_model_update_last_check_time":0}
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
1	1	1779290426	2	模型测试	xiaopihong	模型测试	gpt-5.4	8	18	11	3	f	1	\N	0	default				{"admin_info":{"use_channel":null},"cache_ratio":1,"cache_tokens":0,"completion_ratio":6,"frt":-1000,"group_ratio":1,"model_price":-1,"model_ratio":0.1,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","user_group_ratio":-1}
2	7	1779291044	4	新用户注册赠送 ＄2.000000 额度	sso-admin			0	0	0	0	f	0	\N	0					
3	7	1779291093	2		sso-admin	自用	gpt-5.4	10	18	13	4	t	1	\N	3	default		202605201531299632529678268d9d6ghtvJcf0		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":1,"cache_tokens":0,"completion_ratio":6,"frt":2312,"group_ratio":1,"model_price":-1,"model_ratio":0.1,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","stream_status":{"end_reason":"done","status":"ok"},"user_group_ratio":-1}
4	7	1779292106	2		sso-admin	自用	gpt-5.4	10	18	13	4	t	1	\N	3	default		202605201548223954959628268d9d6YcQxsnHe		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":1,"cache_tokens":0,"completion_ratio":6,"frt":3123,"group_ratio":1,"model_price":-1,"model_ratio":0.1,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","stream_status":{"end_reason":"done","status":"ok"},"user_group_ratio":-1}
5	7	1779292113	2		sso-admin	自用	gpt-5.4	35	40	51	2	t	1	\N	3	default		20260520154831384671268268d9d6ZNi9qx8p		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":1,"cache_tokens":0,"completion_ratio":6,"frt":873,"group_ratio":1,"model_price":-1,"model_ratio":0.1,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","stream_status":{"end_reason":"done","status":"ok"},"user_group_ratio":-1}
6	7	1779292558	2		sso-admin	自用	gpt-5.4	374	24	619	15	t	1	\N	3	default		202605201555432866331988268d9d6fnZcl2Pu		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":1,"cache_tokens":0,"completion_ratio":6,"frt":3005,"group_ratio":1,"model_price":-1,"model_ratio":0.1,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","stream_status":{"end_reason":"done","status":"ok"},"user_group_ratio":-1}
7	7	1779294219	2		sso-admin	自用	gpt-5.4	10	18	13	4	t	1	\N	3	default		202605201623358730717358268d9d6XpsddUia		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":1,"cache_tokens":0,"completion_ratio":6,"frt":2584,"group_ratio":1,"model_price":-1,"model_ratio":0.1,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","stream_status":{"end_reason":"done","status":"ok"},"user_group_ratio":-1}
8	9	1779318292	4	新用户注册赠送 ＄2.000000 额度	mmx			0	0	0	0	f	0	\N	0					
9	9	1779318488	2		mmx	111	gpt-5.4	182	24	300	9	t	1	\N	4	default		202605202307596046879308268d9d6ytaqFFKM		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":1,"cache_tokens":0,"completion_ratio":6,"frt":2479,"group_ratio":1,"model_price":-1,"model_ratio":0.1,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","stream_status":{"end_reason":"done","status":"ok"},"user_group_ratio":-1}
10	10	1779331795	4	新用户注册赠送 ＄2.000000 额度	chenruihua			0	0	0	0	f	0	\N	0					
11	10	1779332047	2		chenruihua	Common	gpt-5.4	248	26	409	11	t	1	\N	5	default		202605210253561530708938268d9d65e6bHOaA		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":1,"cache_tokens":0,"completion_ratio":6,"frt":3005,"group_ratio":1,"model_price":-1,"model_ratio":0.1,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","stream_status":{"end_reason":"done","status":"ok"},"user_group_ratio":-1}
12	10	1779332075	2		chenruihua	Common	gpt-5.4	359	454	523	12	t	1	\N	5	default		202605210254237979465258268d9d6InrzdYLD		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":1,"cache_tokens":0,"completion_ratio":6,"frt":1062,"group_ratio":1,"model_price":-1,"model_ratio":0.1,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","stream_status":{"end_reason":"done","status":"ok"},"user_group_ratio":-1}
13	10	1779332137	2		chenruihua	Common	gpt-5.4	449	985	584	14	t	1	\N	5	default		202605210255237557623748268d9d66F8ewLRX		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":1,"cache_tokens":0,"completion_ratio":6,"frt":993,"group_ratio":1,"model_price":-1,"model_ratio":0.1,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","stream_status":{"end_reason":"done","status":"ok"},"user_group_ratio":-1}
14	10	1779332214	2		chenruihua	Common	gpt-5.4	336	25	556	27	t	1	\N	5	default		202605210256276490869528268d9d6PtwPor63		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":1,"cache_tokens":0,"completion_ratio":6,"frt":3627,"group_ratio":1,"model_price":-1,"model_ratio":0.1,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","stream_status":{"end_reason":"done","status":"ok"},"user_group_ratio":-1}
15	11	1779332891	4	新用户注册赠送 ＄2.000000 额度	1061242019			0	0	0	0	f	0	\N	0					
16	10	1779332948	2		chenruihua	Common	gpt-5.4	143	27	234	8	t	1	\N	5	default		20260521030900211434648268d9d6hoPbbNym		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":1,"cache_tokens":0,"completion_ratio":6,"frt":3262,"group_ratio":1,"model_price":-1,"model_ratio":0.1,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","stream_status":{"end_reason":"done","status":"ok"},"user_group_ratio":-1}
17	10	1779332975	2		chenruihua	Common	gpt-5.4	95	279	112	4	t	1	\N	5	default		202605210309317327006058268d9d6fQPy5gUh		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":1,"cache_tokens":0,"completion_ratio":6,"frt":1248,"group_ratio":1,"model_price":-1,"model_ratio":0.1,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","stream_status":{"end_reason":"done","status":"ok"},"user_group_ratio":-1}
18	10	1779332995	2		chenruihua	Common	gpt-5.4	78	410	62	5	t	1	\N	5	default		202605210309509176122098268d9d6fRsetymz		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":1,"cache_tokens":0,"completion_ratio":6,"frt":957,"group_ratio":1,"model_price":-1,"model_ratio":0.1,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","stream_status":{"end_reason":"done","status":"ok"},"user_group_ratio":-1}
19	10	1779333007	2		chenruihua	Common	gpt-5.4	96	483	79	3	t	1	\N	5	default		202605210310047056701188268d9d69B5iDY6E		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":1,"cache_tokens":0,"completion_ratio":6,"frt":1854,"group_ratio":1,"model_price":-1,"model_ratio":0.1,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","stream_status":{"end_reason":"done","status":"ok"},"user_group_ratio":-1}
20	10	1779333029	2		chenruihua	Common	gpt-5.4	133	571	126	3	t	1	\N	5	default		202605210310262316005988268d9d6u0WD6fGm		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":1,"cache_tokens":0,"completion_ratio":6,"frt":878,"group_ratio":1,"model_price":-1,"model_ratio":0.1,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","stream_status":{"end_reason":"done","status":"ok"},"user_group_ratio":-1}
21	10	1779333053	2		chenruihua	Common	gpt-5.4	142	708	119	4	t	1	\N	5	default		202605210310497735635858268d9d6r6wJpU8x		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":1,"cache_tokens":0,"completion_ratio":6,"frt":951,"group_ratio":1,"model_price":-1,"model_ratio":0.1,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","stream_status":{"end_reason":"done","status":"ok"},"user_group_ratio":-1}
22	10	1779333075	2		chenruihua	Common	gpt-5.4	221	838	228	6	t	1	\N	5	default		202605210311093463202428268d9d6BSQh1DiV		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":1,"cache_tokens":0,"completion_ratio":6,"frt":1070,"group_ratio":1,"model_price":-1,"model_ratio":0.1,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","stream_status":{"end_reason":"done","status":"ok"},"user_group_ratio":-1}
23	7	1779333734	2		sso-admin	自用	gpt-5.4	10	18	13	7	t	1	\N	6	default		202605210322076211831978268d9d6GYwmGOpl		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":1,"cache_tokens":0,"completion_ratio":6,"frt":5166,"group_ratio":1,"model_price":-1,"model_ratio":0.1,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","stream_status":{"end_reason":"done","status":"ok"},"user_group_ratio":-1}
24	11	1779335743	2		1061242019	111	gpt-5.4	30	20	46	4	t	1	\N	7	default		20260521035539365808028268d9d6eLps0t99		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":1,"cache_tokens":0,"completion_ratio":6,"frt":2875,"group_ratio":1,"model_price":-1,"model_ratio":0.1,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","stream_status":{"end_reason":"done","status":"ok"},"user_group_ratio":-1}
25	11	1779335763	2		1061242019	111	gpt-5.4	56	73	81	3	t	1	\N	7	default		202605210356007137224038268d9d6ouwLBMQT		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":1,"cache_tokens":0,"completion_ratio":6,"frt":1020,"group_ratio":1,"model_price":-1,"model_ratio":0.1,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","stream_status":{"end_reason":"done","status":"ok"},"user_group_ratio":-1}
26	10	1779338892	2		chenruihua	Common	gpt-5.4	17	22	25	5	t	1	\N	5	default		202605210448079615680308268d9d6RDyb41TX		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":1,"cache_tokens":0,"completion_ratio":6,"frt":3005,"group_ratio":1,"model_price":-1,"model_ratio":0.1,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","stream_status":{"end_reason":"done","status":"ok"},"user_group_ratio":-1}
27	10	1779338905	2		chenruihua	Common	gpt-5.4	18	59	20	3	t	1	\N	5	default		202605210448225474332348268d9d6WS81BgbE		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":1,"cache_tokens":0,"completion_ratio":6,"frt":1905,"group_ratio":1,"model_price":-1,"model_ratio":0.1,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","stream_status":{"end_reason":"done","status":"ok"},"user_group_ratio":-1}
28	10	1779339018	2		chenruihua	Common	gpt-5.4	29	96	33	3	t	1	\N	5	default		202605210450158405745628268d9d6aGNYOv7l		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":1,"cache_tokens":0,"completion_ratio":6,"frt":1228,"group_ratio":1,"model_price":-1,"model_ratio":0.1,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","stream_status":{"end_reason":"done","status":"ok"},"user_group_ratio":-1}
29	10	1779339068	2		chenruihua	Common	gpt-5.4	49	140	59	4	t	1	\N	5	default		202605210451046914808148268d9d6oQyfDpxo		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":1,"cache_tokens":0,"completion_ratio":6,"frt":1012,"group_ratio":1,"model_price":-1,"model_ratio":0.1,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","stream_status":{"end_reason":"done","status":"ok"},"user_group_ratio":-1}
30	10	1779339137	2		chenruihua	Common	gpt-5.4	73	212	87	3	t	1	\N	5	default		20260521045214761214748268d9d6nms680z2		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":1,"cache_tokens":0,"completion_ratio":6,"frt":1188,"group_ratio":1,"model_price":-1,"model_ratio":0.1,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","stream_status":{"end_reason":"done","status":"ok"},"user_group_ratio":-1}
31	10	1779339207	2		chenruihua	Common	gpt-5.4	194	317	271	6	t	1	\N	5	default		202605210453212200946868268d9d6iViUBY5H		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":1,"cache_tokens":0,"completion_ratio":6,"frt":1012,"group_ratio":1,"model_price":-1,"model_ratio":0.1,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","stream_status":{"end_reason":"done","status":"ok"},"user_group_ratio":-1}
62	7	1779470132	5	status_code=502, Upstream request failed	sso-admin	自用	gpt-image-2	0	0	0	70	f	1	\N	6	default		202605221714225116362208268d9d6KSEbil1D		{"admin_info":{"use_channel":["1"]},"channel_id":1,"channel_name":"openai","channel_type":1,"error_code":"unknown_error","error_type":"openai_error","request_path":"/v1/images/edits","status_code":502}
32	10	1779339361	2		chenruihua	Common	gpt-5.4	111	598	86	4	t	1	\N	5	default		202605210455578567227908268d9d6JHfHRZfn		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":1,"cache_tokens":0,"completion_ratio":6,"frt":1279,"group_ratio":1,"model_price":-1,"model_ratio":0.1,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","stream_status":{"end_reason":"done","status":"ok"},"user_group_ratio":-1}
33	10	1779339419	2		chenruihua	Common	gpt-5.4	173	695	172	7	t	1	\N	5	default		202605210456529209173978268d9d6HfVBnmK5		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":1,"cache_tokens":0,"completion_ratio":6,"frt":1981,"group_ratio":1,"model_price":-1,"model_ratio":0.1,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","stream_status":{"end_reason":"done","status":"ok"},"user_group_ratio":-1}
34	10	1779339528	2		chenruihua	Common	gpt-5.4	192	879	174	6	t	1	\N	5	default		202605210458428945784168268d9d629WjlF5y		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":1,"cache_tokens":0,"completion_ratio":6,"frt":1727,"group_ratio":1,"model_price":-1,"model_ratio":0.1,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","stream_status":{"end_reason":"done","status":"ok"},"user_group_ratio":-1}
35	10	1779339622	2		chenruihua	Common	gpt-5.4	172	1067	108	5	t	1	\N	5	default		202605210500177908845468268d9d6xVLPg3NZ		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":1,"cache_tokens":0,"completion_ratio":6,"frt":1205,"group_ratio":1,"model_price":-1,"model_ratio":0.1,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","stream_status":{"end_reason":"done","status":"ok"},"user_group_ratio":-1}
36	10	1779339684	2		chenruihua	Common	gpt-5.4	225	1184	178	7	t	1	\N	5	default		202605210501179162789698268d9d6uIPKaQKI		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":1,"cache_tokens":0,"completion_ratio":6,"frt":1531,"group_ratio":1,"model_price":-1,"model_ratio":0.1,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","stream_status":{"end_reason":"done","status":"ok"},"user_group_ratio":-1}
37	10	1779339699	2		chenruihua	Common	gpt-5.4	196	1370	98	4	t	1	\N	5	default		202605210501355660693078268d9d6TQm2JNDx		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":1,"cache_tokens":0,"completion_ratio":6,"frt":1513,"group_ratio":1,"model_price":-1,"model_ratio":0.1,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","stream_status":{"end_reason":"done","status":"ok"},"user_group_ratio":-1}
38	10	1779339739	2		chenruihua	Common	gpt-5.4	211	1480	105	4	t	1	\N	5	default		20260521050215938719348268d9d6AkODwION		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":1,"cache_tokens":0,"completion_ratio":6,"frt":1538,"group_ratio":1,"model_price":-1,"model_ratio":0.1,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","stream_status":{"end_reason":"done","status":"ok"},"user_group_ratio":-1}
39	10	1779339844	2		chenruihua	Common	gpt-5.4	858	1605	1163	22	t	1	\N	5	default		202605210503421164030388268d9d6S0ZvQtfA		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":1,"cache_tokens":1280,"completion_ratio":6,"frt":1317,"group_ratio":1,"model_price":-1,"model_ratio":0.1,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","stream_status":{"end_reason":"done","status":"ok"},"user_group_ratio":-1}
40	10	1779340010	2		chenruihua	Common	gpt-5.4	315	2782	61	3	t	1	\N	5	default		20260521050647915848048268d9d6GMFYaWP5		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":1,"cache_tokens":1280,"completion_ratio":6,"frt":1421,"group_ratio":1,"model_price":-1,"model_ratio":0.1,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","stream_status":{"end_reason":"done","status":"ok"},"user_group_ratio":-1}
41	10	1779340040	2		chenruihua	Common	gpt-5.4	321	2854	60	3	t	1	\N	5	default		202605210507172223337518268d9d61oVzM7vE		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":1,"cache_tokens":2304,"completion_ratio":6,"frt":1216,"group_ratio":1,"model_price":-1,"model_ratio":0.1,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","stream_status":{"end_reason":"done","status":"ok"},"user_group_ratio":-1}
42	10	1779340100	2		chenruihua	Common	gpt-5.4	367	2927	123	5	t	1	\N	5	default		202605210508158030589418268d9d6GnKIZOiV		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":1,"cache_tokens":2304,"completion_ratio":6,"frt":1211,"group_ratio":1,"model_price":-1,"model_ratio":0.1,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","stream_status":{"end_reason":"done","status":"ok"},"user_group_ratio":-1}
43	10	1779340135	2		chenruihua	Common	gpt-5.4	354	3057	80	4	t	1	\N	5	default		202605210508517435483368268d9d61OZikSG7		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":1,"cache_tokens":2816,"completion_ratio":6,"frt":1287,"group_ratio":1,"model_price":-1,"model_ratio":0.1,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","stream_status":{"end_reason":"done","status":"ok"},"user_group_ratio":-1}
44	10	1779340177	2		chenruihua	Common	gpt-5.4	339	3147	41	3	t	1	\N	5	default		202605210509346802445448268d9d6FMlC5e4D		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":1,"cache_tokens":2816,"completion_ratio":6,"frt":1182,"group_ratio":1,"model_price":-1,"model_ratio":0.1,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","stream_status":{"end_reason":"done","status":"ok"},"user_group_ratio":-1}
45	10	1779340241	2		chenruihua	Common	gpt-5.4	372	3202	87	4	t	1	\N	5	default		202605210510377146086198268d9d6tp35zAMJ		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":1,"cache_tokens":2816,"completion_ratio":6,"frt":1410,"group_ratio":1,"model_price":-1,"model_ratio":0.1,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","stream_status":{"end_reason":"done","status":"ok"},"user_group_ratio":-1}
125	1	1779701933	2	模型测试	xiaopihong	模型测试	mimo-v2.5	132	248	16	1	f	2	\N	0	default				{"admin_info":{"use_channel":null},"cache_ratio":1,"cache_tokens":192,"completion_ratio":1,"frt":-1000,"group_ratio":1,"model_price":-1,"model_ratio":0.5,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","user_group_ratio":-1}
46	10	1779340379	2		chenruihua	Common	gpt-5.4	13	25	18	4	t	1	\N	5	default		202605210512556857800428268d9d6hYZVszQ8		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":1,"cache_tokens":0,"completion_ratio":6,"frt":3367,"group_ratio":1,"model_price":-1,"model_ratio":0.1,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","stream_status":{"end_reason":"done","status":"ok"},"user_group_ratio":-1}
47	10	1779340499	2		chenruihua	Common	gpt-5.4	30	20	47	21	t	1	\N	5	default		202605210514387495290468268d9d6zgxWwqmY		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":1,"cache_tokens":0,"completion_ratio":6,"frt":20588,"group_ratio":1,"model_price":-1,"model_ratio":0.1,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","stream_status":{"end_reason":"done","status":"ok"},"user_group_ratio":-1}
48	10	1779340549	2		chenruihua	Common	gpt-5.4	426	85	695	20	t	1	\N	5	default		202605210515295648701008268d9d6SS7Jov8S		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":1,"cache_tokens":0,"completion_ratio":6,"frt":931,"group_ratio":1,"model_price":-1,"model_ratio":0.1,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","stream_status":{"end_reason":"done","status":"ok"},"user_group_ratio":-1}
49	10	1779340699	2		chenruihua	Common	gpt-5.4	76	23	123	6	t	1	\N	5	default		202605210518133006471018268d9d6yJescpWu		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":1,"cache_tokens":0,"completion_ratio":6,"frt":3530,"group_ratio":1,"model_price":-1,"model_ratio":0.1,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","stream_status":{"end_reason":"done","status":"ok"},"user_group_ratio":-1}
50	10	1779340858	2		chenruihua	Common	gpt-5.4	208	158	321	9	t	1	\N	5	default		202605210520499122608178268d9d6Bd9ixnJx		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":1,"cache_tokens":0,"completion_ratio":6,"frt":1208,"group_ratio":1,"model_price":-1,"model_ratio":0.1,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","stream_status":{"end_reason":"done","status":"ok"},"user_group_ratio":-1}
51	7	1779421822	5	status_code=502, Upstream request failed	sso-admin	自用	gpt-image-2	0	0	0	5	f	1	\N	6	default		202605220350173087825588268d9d6j5aPePl4		{"admin_info":{"use_channel":["1"]},"channel_id":1,"channel_name":"openai","channel_type":1,"error_code":"unknown_error","error_type":"openai_error","request_path":"/v1/images/generations","status_code":502}
52	7	1779421839	5	status_code=502, Upstream request failed	sso-admin	自用	gpt-image-2	0	0	0	0	f	1	\N	6	default		202605220350389873577618268d9d6WtOjV0Xr		{"admin_info":{"use_channel":["1"]},"channel_id":1,"channel_name":"openai","channel_type":1,"error_code":"unknown_error","error_type":"openai_error","request_path":"/v1/images/generations","status_code":502}
53	7	1779421882	5	status_code=502, Upstream request failed	sso-admin	自用	gpt-image-2	0	0	0	0	f	1	\N	6	default		202605220351212507458238268d9d66qOe5zap		{"admin_info":{"use_channel":["1"]},"channel_id":1,"channel_name":"openai","channel_type":1,"error_code":"unknown_error","error_type":"openai_error","request_path":"/v1/images/generations","status_code":502}
54	10	1779437703	2	大小 1024x1024, 品质 standard, 生成数量 1	chenruihua	Common	gpt-image-2	250000	77	1756	106	f	1	\N	5	default		202605220813177019456048268d9d6Pyxg4yE4		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/generations","user_group_ratio":-1}
55	10	1779437985	2	大小 1024x1024, 品质 standard, 生成数量 1	chenruihua	Common	gpt-image-2	250000	118	1756	83	f	1	\N	5	default		202605220818227845002638268d9d6ac3GCZBa		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/generations","user_group_ratio":-1}
56	7	1779469586	2	大小 1024x1024, 品质 standard, 生成数量 1	sso-admin	自用	gpt-image-2	250000	1340	1756	77	f	1	\N	6	default		202605221705089818634678268d9d6kaNpaeIs		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"image":true,"image_output":1024,"image_ratio":0,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
57	7	1779469698	2	大小 1024x1024, 品质 standard, 生成数量 1	sso-admin	自用	gpt-image-2	250000	1236	1756	76	f	1	\N	6	default		2026052217070219752118268d9d6D4dTH27V		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"image":true,"image_output":1024,"image_ratio":0,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
58	3	1779469975	3	管理员增加用户额度 ＄30.000000 额度	xph-admin			0	0	0	0	f	0	\N	0					{"admin_info":{"admin_id":1,"admin_username":"xiaopihong"}}
59	10	1779469993	3	管理员增加用户额度 ＄20.000000 额度	chenruihua			0	0	0	0	f	0	\N	0					{"admin_info":{"admin_id":1,"admin_username":"xiaopihong"}}
60	7	1779470011	2	大小 1024x1024, 品质 standard, 生成数量 1	sso-admin	自用	gpt-image-2	250000	1347	1756	84	f	1	\N	6	default		202605221712079886580598268d9d6ZAvePkCY		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"image":true,"image_output":1024,"image_ratio":0,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
61	7	1779470038	3	管理员增加用户额度 ＄20.000000 额度	sso-admin			0	0	0	0	f	0	\N	0					{"admin_info":{"admin_id":1,"admin_username":"xiaopihong"}}
126	1	1779701943	2	模型测试	xiaopihong	模型测试	gpt-5.4	63	18	11	2	f	1	\N	0	default				{"admin_info":{"use_channel":null},"cache_ratio":1,"cache_tokens":0,"completion_ratio":6,"frt":-1000,"group_ratio":1,"model_price":-1,"model_ratio":0.75,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","user_group_ratio":-1}
63	7	1779470294	2	大小 1024x1024, 品质 standard, 生成数量 1	sso-admin	自用	gpt-image-2	250000	1265	1756	76	f	1	\N	6	default		20260522171658618443238268d9d64fDxAafM		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"image":true,"image_output":1024,"image_ratio":0,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
64	10	1779509581	5	status_code=502, Upstream request failed	chenruihua	Common	gpt-image-2	0	0	0	74	f	1	\N	5	default		202605230411468365007558268d9d64rqGWAYQ		{"admin_info":{"use_channel":["1"]},"channel_id":1,"channel_name":"openai","channel_type":1,"error_code":"unknown_error","error_type":"openai_error","request_path":"/v1/images/edits","status_code":502}
65	10	1779509879	2	大小 1024x1024, 品质 standard, 生成数量 1	chenruihua	Common	gpt-image-2	250000	1309	1756	75	f	1	\N	5	default		202605230416448156966808268d9d6k9pgdWEo		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"image":true,"image_output":1024,"image_ratio":0,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
66	7	1779552324	2	大小 1216x544, 品质 standard, 生成数量 1	sso-admin	自用	gpt-image-2	250000	1170	671	114	f	1	\N	6	default		202605231603309265796388268d9d6stQMZPZC		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"image":true,"image_output":646,"image_ratio":0,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
67	7	1779552461	2	大小 1216x544, 品质 standard, 生成数量 1	sso-admin	自用	gpt-image-2	250000	974	671	58	f	1	\N	6	default		202605231606438122477938268d9d6J4dlM9NS		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"image":true,"image_output":646,"image_ratio":0,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
68	7	1779552612	5	status_code=502, Upstream request failed	sso-admin	自用	gpt-image-2	0	0	0	85	f	1	\N	6	default		202605231608466376333158268d9d6g32SqkkA		{"admin_info":{"use_channel":["1"]},"channel_id":1,"channel_name":"openai","channel_type":1,"error_code":"unknown_error","error_type":"openai_error","request_path":"/v1/images/edits","status_code":502}
69	7	1779552840	2	大小 1216x544, 品质 standard, 生成数量 1	sso-admin	自用	gpt-image-2	250000	1051	671	58	f	1	\N	6	default		202605231613029207567628268d9d6CPAlblOO		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"image":true,"image_output":646,"image_ratio":0,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
70	7	1779589189	2	大小 1216x544, 品质 standard, 生成数量 1	sso-admin	自用	gpt-image-2	250000	1087	671	47	f	1	\N	6	default		202605240219024614218468268d9d6JIpqT0KI		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"image":true,"image_output":646,"image_ratio":0,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
71	7	1779589330	5	status_code=502, Upstream request failed	sso-admin	自用	gpt-image-2	0	0	0	39	f	1	\N	6	default		202605240221318212791628268d9d6HQ1V6pOG		{"admin_info":{"use_channel":["1"]},"channel_id":1,"channel_name":"openai","channel_type":1,"error_code":"unknown_error","error_type":"openai_error","request_path":"/v1/images/edits","status_code":502}
72	7	1779589422	2	大小 1216x544, 品质 standard, 生成数量 1	sso-admin	自用	gpt-image-2	250000	1025	2747	89	f	1	\N	6	default		202605240222139903656808268d9d6B85AmatZ		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"image":true,"image_output":646,"image_ratio":0,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
73	7	1779590856	2		sso-admin	自用	gpt-5.4	13389	260	2932	58	t	1	\N	6	default		202605240246384311178518268d9d6TvbaT2Q4		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":1,"cache_tokens":0,"completion_ratio":6,"frt":6317,"group_ratio":1,"model_price":-1,"model_ratio":0.75,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","stream_status":{"end_reason":"done","status":"ok"},"user_group_ratio":-1}
74	7	1779590906	2		sso-admin	自用	gpt-5.4	8861	3217	1433	29	t	1	\N	6	default		202605240247577866964048268d9d680LWRhtu		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":1,"cache_tokens":0,"completion_ratio":6,"frt":1814,"group_ratio":1,"model_price":-1,"model_ratio":0.75,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","stream_status":{"end_reason":"done","status":"ok"},"user_group_ratio":-1}
75	7	1779591168	2	大小 1216x544, 品质 standard, 生成数量 1	sso-admin	自用	gpt-image-2	250000	1371	671	58	f	1	\N	6	default		202605240251504317722798268d9d6JsJrUWGR		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"image":true,"image_output":646,"image_ratio":0,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
76	7	1779591344	2		sso-admin	自用	gpt-5.4	11452	4661	1768	36	t	1	\N	6	default		202605240255089767210028268d9d6z1wQKm5i		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":1,"cache_tokens":0,"completion_ratio":6,"frt":2268,"group_ratio":1,"model_price":-1,"model_ratio":0.75,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","stream_status":{"end_reason":"done","status":"ok"},"user_group_ratio":-1}
127	1	1779701949	2	模型测试	xiaopihong	模型测试	gpt-5.4	63	18	11	2	f	1	\N	0	default				{"admin_info":{"use_channel":null},"cache_ratio":1,"cache_tokens":0,"completion_ratio":6,"frt":-1000,"group_ratio":1,"model_price":-1,"model_ratio":0.75,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","user_group_ratio":-1}
77	7	1779591503	2		sso-admin	自用	gpt-5.4	13934	6441	2023	75	t	1	\N	6	default		202605240257085189883378268d9d6FbFRiMEi		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":1,"cache_tokens":0,"completion_ratio":6,"frt":3541,"group_ratio":1,"model_price":-1,"model_ratio":0.75,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","stream_status":{"end_reason":"done","status":"ok"},"user_group_ratio":-1}
78	7	1779591700	2		sso-admin	自用	gpt-5.4	774	24	168	7	t	1	\N	6	default		2026052403013340838108268d9d6ucdeK6aC		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":1,"cache_tokens":0,"completion_ratio":6,"frt":3062,"group_ratio":1,"model_price":-1,"model_ratio":0.75,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","stream_status":{"end_reason":"done","status":"ok"},"user_group_ratio":-1}
79	7	1779595730	2		sso-admin	自用	gpt-5.4	6005	26	1330	27	t	1	\N	6	default		202605240408239151615018268d9d64fBryTls		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":1,"cache_tokens":0,"completion_ratio":6,"frt":1423,"group_ratio":1,"model_price":-1,"model_ratio":0.75,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","stream_status":{"end_reason":"done","status":"ok"},"user_group_ratio":-1}
80	10	1779606017	5	status_code=502, Upstream request failed	chenruihua	Common	gpt-image-2	0	0	0	1	f	1	\N	5	default		202605240700159656513258268d9d63Arz8tfL		{"admin_info":{"use_channel":["1"]},"channel_id":1,"channel_name":"openai","channel_type":1,"error_code":"unknown_error","error_type":"openai_error","request_path":"/v1/images/edits","status_code":502}
81	10	1779606298	2	大小 1024x640, 品质 standard, 生成数量 1	chenruihua	Common	gpt-image-2	250000	1093	956	275	f	1	\N	5	default		202605240700235308176028268d9d6wa2vsRSn		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"image":true,"image_output":640,"image_ratio":0,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
82	10	1779606810	2	大小 1024x640, 品质 standard, 生成数量 1	chenruihua	Common	gpt-image-2	250000	959	3824	197	f	1	\N	5	default		20260524071013954859528268d9d6zJG4flEV		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"image":true,"image_output":640,"image_ratio":0,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
83	10	1779608205	2	大小 1216x544, 品质 standard, 生成数量 1	chenruihua	Common	gpt-image-2	250000	1225	671	56	f	1	\N	5	default		202605240735496058512408268d9d60yMCuocv		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"image":true,"image_output":646,"image_ratio":0,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
84	10	1779608710	2	大小 1216x544, 品质 standard, 生成数量 1	chenruihua	Common	gpt-image-2	250000	1061	671	49	f	1	\N	5	default		20260524074421286221848268d9d6byFqKM8h		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"image":true,"image_output":646,"image_ratio":0,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
85	7	1779610092	2	大小 1216x544, 品质 standard, 生成数量 1	sso-admin	自用	gpt-image-2	250000	939	671	86	f	1	\N	6	default		202605240806461936663658268d9d6IV6NiYfv		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"image":true,"image_output":646,"image_ratio":0,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
86	10	1779613116	2	大小 1216x544, 品质 standard, 生成数量 1	chenruihua	Common	gpt-image-2	250000	1029	671	56	f	1	\N	5	default		202605240857408197297818268d9d6AhGErvNu		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"image":true,"image_output":646,"image_ratio":0,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
87	1	1779623865	2	模型测试	xiaopihong	模型测试	gpt-5.4	63	18	11	1	f	1	\N	0	default				{"admin_info":{"use_channel":null},"cache_ratio":1,"cache_tokens":0,"completion_ratio":6,"frt":-1000,"group_ratio":1,"model_price":-1,"model_ratio":0.75,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","user_group_ratio":-1}
88	1	1779623874	2	模型测试	xiaopihong	模型测试	gpt-5.5	96	18	13	1	f	1	\N	0	default				{"admin_info":{"use_channel":null},"cache_ratio":1,"cache_tokens":0,"completion_ratio":6,"frt":-1000,"group_ratio":1,"model_price":-1,"model_ratio":1,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","user_group_ratio":-1}
89	1	1779623876	2	模型测试	xiaopihong	模型测试	gpt-5.4	63	18	11	3	f	1	\N	0	default				{"admin_info":{"use_channel":null},"cache_ratio":1,"cache_tokens":0,"completion_ratio":6,"frt":-1000,"group_ratio":1,"model_price":-1,"model_ratio":0.75,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","user_group_ratio":-1}
90	1	1779624110	2	模型测试	xiaopihong	模型测试	gpt-5.4	63	18	11	1	f	1	\N	0	default				{"admin_info":{"use_channel":null},"cache_ratio":1,"cache_tokens":0,"completion_ratio":6,"frt":-1000,"group_ratio":1,"model_price":-1,"model_ratio":0.75,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","user_group_ratio":-1}
91	1	1779624113	2	模型测试	xiaopihong	模型测试	gpt-5.5	96	18	13	5	f	1	\N	0	default				{"admin_info":{"use_channel":null},"cache_ratio":1,"cache_tokens":0,"completion_ratio":6,"frt":-1000,"group_ratio":1,"model_price":-1,"model_ratio":1,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","user_group_ratio":-1}
92	1	1779624391	5	status_code=400, The 'gpt-image-2' model is not supported when using Codex with a ChatGPT account.	xiaopihong	playground-default	gpt-image-2	0	0	0	1	t	1	\N	0	default		202605241206305252162708268d9d6EqrvRy05		{"admin_info":{"use_channel":["1"]},"channel_id":1,"channel_name":"openai","channel_type":1,"error_code":"unknown_error","error_type":"openai_error","request_path":"/pg/chat/completions","status_code":400}
93	1	1779624499	5	status_code=400, The 'gpt-image-2' model is not supported when using Codex with a ChatGPT account.	xiaopihong	playground-default	gpt-image-2	0	0	0	0	t	1	\N	0	default		202605241208187243218748268d9d6yXQDmkQp		{"admin_info":{"use_channel":["1"]},"channel_id":1,"channel_name":"openai","channel_type":1,"error_code":"unknown_error","error_type":"openai_error","request_path":"/pg/chat/completions","status_code":400}
94	1	1779624565	2	模型测试	xiaopihong	模型测试	gpt-5.4	63	18	11	1	f	1	\N	0	default				{"admin_info":{"use_channel":null},"cache_ratio":1,"cache_tokens":0,"completion_ratio":6,"frt":-1000,"group_ratio":1,"model_price":-1,"model_ratio":0.75,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","user_group_ratio":-1}
95	1	1779624577	2	模型测试	xiaopihong	模型测试	gpt-5.4	63	18	11	1	f	1	\N	0	default				{"admin_info":{"use_channel":null},"cache_ratio":1,"cache_tokens":0,"completion_ratio":6,"frt":-1000,"group_ratio":1,"model_price":-1,"model_ratio":0.75,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","user_group_ratio":-1}
96	1	1779624578	2	模型测试	xiaopihong	模型测试	gpt-5.5	96	18	13	2	f	1	\N	0	default				{"admin_info":{"use_channel":null},"cache_ratio":1,"cache_tokens":0,"completion_ratio":6,"frt":-1000,"group_ratio":1,"model_price":-1,"model_ratio":1,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","user_group_ratio":-1}
97	1	1779624926	2	模型测试	xiaopihong	模型测试	gpt-5.4	63	18	11	2	f	1	\N	0	default				{"admin_info":{"use_channel":null},"cache_ratio":1,"cache_tokens":0,"completion_ratio":6,"frt":-1000,"group_ratio":1,"model_price":-1,"model_ratio":0.75,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","user_group_ratio":-1}
98	1	1779624930	2	模型测试	xiaopihong	模型测试	gpt-5.4	63	18	11	1	f	1	\N	0	default				{"admin_info":{"use_channel":null},"cache_ratio":1,"cache_tokens":0,"completion_ratio":6,"frt":-1000,"group_ratio":1,"model_price":-1,"model_ratio":0.75,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","user_group_ratio":-1}
99	1	1779624930	2	模型测试	xiaopihong	模型测试	gpt-5.5	96	18	13	2	f	1	\N	0	default				{"admin_info":{"use_channel":null},"cache_ratio":1,"cache_tokens":0,"completion_ratio":6,"frt":-1000,"group_ratio":1,"model_price":-1,"model_ratio":1,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","user_group_ratio":-1}
100	1	1779624936	2	模型测试	xiaopihong	模型测试	gpt-5.4	63	18	11	1	f	1	\N	0	default				{"admin_info":{"use_channel":null},"cache_ratio":1,"cache_tokens":0,"completion_ratio":6,"frt":-1000,"group_ratio":1,"model_price":-1,"model_ratio":0.75,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","user_group_ratio":-1}
101	1	1779624997	2	模型测试	xiaopihong	模型测试	gpt-5.4	63	18	11	3	f	1	\N	0	default				{"admin_info":{"use_channel":null},"cache_ratio":1,"cache_tokens":0,"completion_ratio":6,"frt":-1000,"group_ratio":1,"model_price":-1,"model_ratio":0.75,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","user_group_ratio":-1}
102	7	1779627452	5	status_code=502, Upstream request failed	sso-admin	自用	gpt-image-2	0	0	0	3	f	1	\N	6	default		20260524125729558834658268d9d6ttOX5bU9		{"admin_info":{"use_channel":["1"]},"channel_id":1,"channel_name":"openai","channel_type":1,"error_code":"unknown_error","error_type":"openai_error","request_path":"/v1/images/generations","status_code":502}
103	1	1779634212	2	模型测试	xiaopihong	模型测试	gpt-5.5	96	18	13	3	f	1	\N	0	default				{"admin_info":{"use_channel":null},"cache_ratio":1,"cache_tokens":0,"completion_ratio":6,"frt":-1000,"group_ratio":1,"model_price":-1,"model_ratio":1,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","user_group_ratio":-1}
104	1	1779634212	2	模型测试	xiaopihong	模型测试	gpt-5.4	63	18	11	4	f	1	\N	0	default				{"admin_info":{"use_channel":null},"cache_ratio":1,"cache_tokens":0,"completion_ratio":6,"frt":-1000,"group_ratio":1,"model_price":-1,"model_ratio":0.75,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","user_group_ratio":-1}
105	1	1779634249	2	模型测试	xiaopihong	模型测试	gpt-5.4	63	18	11	1	f	1	\N	0	default				{"admin_info":{"use_channel":null},"cache_ratio":1,"cache_tokens":0,"completion_ratio":6,"frt":-1000,"group_ratio":1,"model_price":-1,"model_ratio":0.75,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","user_group_ratio":-1}
106	1	1779634249	2	模型测试	xiaopihong	模型测试	gpt-5.5	96	18	13	1	f	1	\N	0	default				{"admin_info":{"use_channel":null},"cache_ratio":1,"cache_tokens":0,"completion_ratio":6,"frt":-1000,"group_ratio":1,"model_price":-1,"model_ratio":1,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","user_group_ratio":-1}
107	1	1779635310	2	模型测试	xiaopihong	模型测试	mimo-v2.5	132	248	16	1	f	2	\N	0	default				{"admin_info":{"use_channel":null},"cache_ratio":1,"cache_tokens":192,"completion_ratio":1,"frt":-1000,"group_ratio":1,"model_price":-1,"model_ratio":0.5,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","user_group_ratio":-1}
108	1	1779636156	2	模型测试	xiaopihong	模型测试	gpt-5.4	63	18	11	1	f	1	\N	0	default				{"admin_info":{"use_channel":null},"cache_ratio":1,"cache_tokens":0,"completion_ratio":6,"frt":-1000,"group_ratio":1,"model_price":-1,"model_ratio":0.75,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","user_group_ratio":-1}
109	3	1779641363	2	大小 1216x544, 品质 standard, 生成数量 1	xph-admin	自用2	gpt-image-2	250000	533	671	53	f	1	\N	2	default		20260524164830746425998268d9d6YmkAaorC		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/generations","user_group_ratio":-1}
110	3	1779641910	5	status_code=502, Upstream request failed	xph-admin	自用2	gpt-image-2	0	0	0	153	f	1	\N	2	default		20260524165557410538788268d9d66IUBSZ5f		{"admin_info":{"use_channel":["1"]},"channel_id":1,"channel_name":"openai","channel_type":1,"error_code":"unknown_error","error_type":"openai_error","request_path":"/v1/images/edits","status_code":502}
111	3	1779642927	5	status_code=502, Upstream request failed	xph-admin	自用2	gpt-image-2	0	0	0	92	f	1	\N	2	default		202605241713551350517758268d9d6PiMqLpz4		{"admin_info":{"use_channel":["1"]},"channel_id":1,"channel_name":"openai","channel_type":1,"error_code":"unknown_error","error_type":"openai_error","request_path":"/v1/images/edits","status_code":502}
112	1	1779676165	2	模型测试	xiaopihong	模型测试	gpt-image-2	250000	0	0	50	f	1	\N	0	default				{"admin_info":{"use_channel":null},"cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/generations","user_group_ratio":-1}
113	7	1779676281	2	大小 1216x544, 品质 standard, 生成数量 1	sso-admin	自用	gpt-image-2	250000	1	0	49	f	1	\N	6	default		202605250230323903593278268d9d6AWoCulDh		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
114	7	1779676394	5	status_code=400, We’re so sorry, but the image we created may violate our guardrails concerning similarity to third-party content. If you think we got it wrong, please retry or edit your prompt.	sso-admin	自用	gpt-image-2	0	0	0	60	f	1	\N	6	default		202605250232133413550218268d9d6XDFtrfIO		{"admin_info":{"use_channel":["1"]},"channel_id":1,"channel_name":"openai","channel_type":1,"error_code":"image_generation_text_response","error_type":"openai_error","request_path":"/v1/images/edits","status_code":400}
115	7	1779676618	2	大小 1216x544, 品质 standard, 生成数量 1	sso-admin	自用	gpt-image-2	250000	1	0	50	f	1	\N	6	default		202605250236088918867418268d9d6hTBKVpm3		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/generations","user_group_ratio":-1}
116	7	1779685625	2		sso-admin	自用	gpt-5.4	72	18	13	8	t	1	\N	6	default		202605250506574562622288268d9d6teTvPrjW		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":1,"cache_tokens":0,"completion_ratio":6,"frt":6263,"group_ratio":1,"model_price":-1,"model_ratio":0.75,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","stream_status":{"end_reason":"done","status":"ok"},"user_group_ratio":-1}
117	7	1779685652	5	status_code=502, Upstream request failed	sso-admin	自用	gpt-image-2	0	0	0	5	f	1	\N	6	default		202605250507271989866858268d9d6Y0V5WIu9		{"admin_info":{"use_channel":["1"]},"channel_id":1,"channel_name":"openai","channel_type":1,"error_code":"unknown_error","error_type":"openai_error","request_path":"/v1/images/generations","status_code":502}
118	7	1779685808	5	status_code=502, Upstream request failed	sso-admin	自用	gpt-image-2	0	0	0	1	f	1	\N	6	default		202605250510067082671318268d9d6yZPaUGgi		{"admin_info":{"use_channel":["1"]},"channel_id":1,"channel_name":"openai","channel_type":1,"error_code":"unknown_error","error_type":"openai_error","request_path":"/v1/images/generations","status_code":502}
119	1	1779691712	2	模型测试	xiaopihong	模型测试	gpt-5.4	63	18	11	600	f	1	\N	0	default				{"admin_info":{"use_channel":null},"cache_ratio":1,"cache_tokens":0,"completion_ratio":6,"frt":-1000,"group_ratio":1,"model_price":-1,"model_ratio":0.75,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","user_group_ratio":-1}
120	11	1779699001	2		1061242019	111	mimo-v2.5	221	251	191	11	t	2	\N	7	default		202605250849507218216398268d9d6NHy4foq0		{"admin_info":{"use_channel":["2"]},"billing_source":"wallet","cache_ratio":1,"cache_tokens":192,"completion_ratio":1,"frt":8293,"group_ratio":1,"model_price":-1,"model_ratio":0.5,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","stream_status":{"end_reason":"done","status":"ok"},"user_group_ratio":-1}
121	11	1779699054	2		1061242019	111	gpt-5.5	1729	109	270	7	t	1	\N	7	default		20260525085047458884418268d9d6lMuSx7vG		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":1,"cache_tokens":0,"completion_ratio":6,"frt":2882,"group_ratio":1,"model_price":-1,"model_ratio":1,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","stream_status":{"end_reason":"done","status":"ok"},"user_group_ratio":-1}
122	11	1779699326	5	status_code=502, Upstream request failed	1061242019	111	gpt-5.5	0	0	0	0	t	1	\N	7	default		202605250855264398875798268d9d6d7EoAzy1		{"admin_info":{"use_channel":["1"]},"channel_id":1,"channel_name":"openai","channel_type":1,"error_code":"unknown_error","error_type":"openai_error","request_path":"/v1/chat/completions","status_code":502}
123	11	1779699367	2		1061242019	111	gpt-5.4	488	411	40	3	t	1	\N	7	default		202605250856042705577198268d9d6S0Habouh		{"admin_info":{"use_channel":["1"]},"billing_source":"wallet","cache_ratio":1,"cache_tokens":0,"completion_ratio":6,"frt":2157,"group_ratio":1,"model_price":-1,"model_ratio":0.75,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","stream_status":{"end_reason":"done","status":"ok"},"user_group_ratio":-1}
124	1	1779701930	2	模型测试	xiaopihong	模型测试	mimo-v2.5	132	248	16	1	f	2	\N	0	default				{"admin_info":{"use_channel":null},"cache_ratio":1,"cache_tokens":192,"completion_ratio":1,"frt":-1000,"group_ratio":1,"model_price":-1,"model_ratio":0.5,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","user_group_ratio":-1}
128	1	1779701949	2	模型测试	xiaopihong	模型测试	gpt-5.5	96	18	13	2	f	1	\N	0	default				{"admin_info":{"use_channel":null},"cache_ratio":1,"cache_tokens":0,"completion_ratio":6,"frt":-1000,"group_ratio":1,"model_price":-1,"model_ratio":1,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","user_group_ratio":-1}
129	1	1779702241	2	模型测试	xiaopihong	模型测试	grok-imagine-image-lite	17	5	64	7	f	4	\N	0	default			20260525094354350409047tPJdVLp9	{"admin_info":{"use_channel":null},"cache_ratio":1,"cache_tokens":0,"completion_ratio":1,"frt":-1000,"group_ratio":1,"model_price":-1,"model_ratio":0.25,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","user_group_ratio":-1}
130	7	1779704224	2	大小 1216x544, 品质 standard, 生成数量 1	sso-admin	自用	grok-imagine-image-lite	1	1	0	13	f	4	\N	6	default		202605251016519936826388268d9d6daw94Nrx	202605251016525479006XY5TlXRT	{"admin_info":{"use_channel":["4"]},"billing_source":"wallet","cache_ratio":1,"cache_tokens":0,"completion_ratio":1,"frt":-1000,"group_ratio":1,"model_price":-1,"model_ratio":0.25,"request_conversion":["openai_image"],"request_path":"/v1/images/generations","user_group_ratio":-1}
131	7	1779704314	5	status_code=502, Upstream request failed	sso-admin	自用	gpt-image-2	0	0	0	1	f	1	\N	6	default		202605251018328797824938268d9d6oWmrRd41		{"admin_info":{"use_channel":["1"]},"channel_id":1,"channel_name":"openai","channel_type":1,"error_code":"unknown_error","error_type":"openai_error","request_path":"/v1/images/edits","status_code":502}
132	7	1779704414	2	大小 1024x640, 品质 standard, 生成数量 1	sso-admin	自用	grok-imagine-image-lite	1	1	0	13	f	4	\N	6	default		202605251020013955463708268d9d6WsHRQJhv	20260525102001437696661yFhbqqxH	{"admin_info":{"use_channel":["4"]},"billing_source":"wallet","cache_ratio":1,"cache_tokens":0,"completion_ratio":1,"frt":-1000,"group_ratio":1,"model_price":-1,"model_ratio":0.25,"request_conversion":["openai_image"],"request_path":"/v1/images/generations","user_group_ratio":-1}
133	3	1779726808	2	大小 1216x544, 品质 standard, 生成数量 1	xph-admin	自用2	grok-imagine-image-lite	1	1	0	18	f	4	\N	2	default		202605251633103571043748268d9d64kaQiZ7f	20260525163310409390501VI7m81o2	{"admin_info":{"use_channel":["4"]},"billing_source":"wallet","cache_ratio":1,"cache_tokens":0,"completion_ratio":1,"frt":-1000,"group_ratio":1,"model_price":-1,"model_ratio":0.25,"request_conversion":["openai_image"],"request_path":"/v1/images/generations","user_group_ratio":-1}
134	3	1779731003	2	大小 1216x544, 品质 standard, 生成数量 1	xph-admin	自用2	grok-imagine-image-lite	1	1	0	15	f	4	\N	2	default		202605251743082275490058268d9d6C05Mipbp	20260525174308259115121z2dd65ab	{"admin_info":{"use_channel":["4"]},"billing_source":"wallet","cache_ratio":1,"cache_tokens":0,"completion_ratio":1,"frt":-1000,"group_ratio":1,"model_price":-1,"model_ratio":0.25,"request_conversion":["openai_image"],"request_path":"/v1/images/generations","user_group_ratio":-1}
135	3	1779731656	2	大小 1216x544, 品质 standard, 生成数量 1	xph-admin	自用2	grok-imagine-image-lite	1	1	0	14	f	4	\N	2	default		202605251754029118681948268d9d66RfrOuKc	20260525175402950658307BpNapPrx	{"admin_info":{"use_channel":["4"]},"billing_source":"wallet","cache_ratio":1,"cache_tokens":0,"completion_ratio":1,"frt":-1000,"group_ratio":1,"model_price":-1,"model_ratio":0.25,"request_conversion":["openai_image"],"request_path":"/v1/images/generations","user_group_ratio":-1}
136	3	1779731851	2	大小 1216x544, 品质 standard, 生成数量 1	xph-admin	自用2	grok-imagine-image-lite	1	1	0	12	f	4	\N	2	default		202605251757193704160378268d9d6sOJwcRAo	20260525175719391842991d4BODtcU	{"admin_info":{"use_channel":["4"]},"billing_source":"wallet","cache_ratio":1,"cache_tokens":0,"completion_ratio":1,"frt":-1000,"group_ratio":1,"model_price":-1,"model_ratio":0.25,"request_conversion":["openai_image"],"request_path":"/v1/images/generations","user_group_ratio":-1}
137	3	1779731936	2	大小 1216x544, 品质 standard, 生成数量 1	xph-admin	自用2	grok-imagine-image-lite	1	1	0	18	f	4	\N	2	default		20260525175838151987668268d9d6PvYIirZR	20260525175838265190493Gav9eek	{"admin_info":{"use_channel":["4"]},"billing_source":"wallet","cache_ratio":1,"cache_tokens":0,"completion_ratio":1,"frt":-1000,"group_ratio":1,"model_price":-1,"model_ratio":0.25,"request_conversion":["openai_image"],"request_path":"/v1/images/generations","user_group_ratio":-1}
138	7	1779732179	2	大小 1216x544, 品质 standard, 生成数量 1	sso-admin	自用	grok-imagine-image-lite	1	1	0	13	f	4	\N	6	default		202605251802462282826168268d9d6tYFQrXQm	20260525180246235187254PuhyrxW5	{"admin_info":{"use_channel":["4"]},"billing_source":"wallet","cache_ratio":1,"cache_tokens":0,"completion_ratio":1,"frt":-1000,"group_ratio":1,"model_price":-1,"model_ratio":0.25,"request_conversion":["openai_image"],"request_path":"/v1/images/generations","user_group_ratio":-1}
139	7	1779732324	2		sso-admin	自用	mimo-v2.5	273	253	293	5	t	2	\N	6	default		202605251805199046104348268d9d6Zg8UXISW		{"admin_info":{"use_channel":["2"]},"billing_source":"wallet","cache_ratio":1,"cache_tokens":192,"completion_ratio":1,"frt":1250,"group_ratio":1,"model_price":-1,"model_ratio":0.5,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","stream_status":{"end_reason":"done","status":"ok"},"user_group_ratio":-1}
140	7	1779732341	2		sso-admin	自用	mimo-v2.5	415	460	369	6	t	2	\N	6	default		20260525180535561787918268d9d6KQWJ3agf		{"admin_info":{"use_channel":["2"]},"billing_source":"wallet","cache_ratio":1,"cache_tokens":192,"completion_ratio":1,"frt":2150,"group_ratio":1,"model_price":-1,"model_ratio":0.5,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","stream_status":{"end_reason":"done","status":"ok"},"user_group_ratio":-1}
141	7	1779732371	2		sso-admin	自用	grok-imagine-image-lite	175	623	76	10	t	4	\N	6	default		202605251806016958473268268d9d6RIdzZLqH	202605251806017344414111Tehba2I	{"admin_info":{"use_channel":["4"]},"billing_source":"wallet","cache_ratio":1,"cache_tokens":0,"completion_ratio":1,"frt":8914,"group_ratio":1,"model_price":-1,"model_ratio":0.25,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","stream_status":{"end_reason":"done","status":"ok"},"user_group_ratio":-1}
142	7	1779732393	2		sso-admin	自用	grok-imagine-image-lite	192	687	82	10	t	4	\N	6	default		202605251806234367011908268d9d6UgE72608	20260525180623449658558RRvEdXPJ	{"admin_info":{"use_channel":["4"]},"billing_source":"wallet","cache_ratio":1,"cache_tokens":0,"completion_ratio":1,"frt":9992,"group_ratio":1,"model_price":-1,"model_ratio":0.25,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","stream_status":{"end_reason":"done","status":"ok"},"user_group_ratio":-1}
143	7	1779732439	2	大小 1024x1024, 品质 standard, 生成数量 1	sso-admin	自用	grok-imagine-image-lite	1	1	0	11	f	4	\N	6	default		202605251807086827460128268d9d6ZJrOq7kc	20260525180708723402781BHAvKs2Y	{"admin_info":{"use_channel":["4"]},"billing_source":"wallet","cache_ratio":1,"cache_tokens":0,"completion_ratio":1,"frt":-1000,"group_ratio":1,"model_price":-1,"model_ratio":0.25,"request_conversion":["openai_image"],"request_path":"/v1/images/generations","user_group_ratio":-1}
144	7	1779732642	2	大小 1216x544, 品质 standard, 生成数量 1	sso-admin	自用	grok-imagine-image-lite	1	1	0	20	f	4	\N	6	default		202605251810223856624408268d9d6fEcknyAk	20260525181022399121906GpRN6FSB	{"admin_info":{"use_channel":["4"]},"billing_source":"wallet","cache_ratio":1,"cache_tokens":0,"completion_ratio":1,"frt":-1000,"group_ratio":1,"model_price":-1,"model_ratio":0.25,"request_conversion":["openai_image"],"request_path":"/v1/images/generations","user_group_ratio":-1}
145	10	1779762032	5	status_code=400, 未指定模型名称，模型名称不能为空 (request id: 20260526022032803901333Az7oiftV)	chenruihua	Common	grok-imagine-image-lite	0	0	0	0	f	4	\N	5	default		202605260220327796487638268d9d65a7zYZU9	20260526022032803901333Az7oiftV	{"admin_info":{"use_channel":["4"]},"channel_id":4,"channel_name":"xai","channel_type":48,"error_code":"","error_type":"openai_error","request_path":"/v1/images/edits","status_code":400}
146	10	1779762040	5	status_code=400, 未指定模型名称，模型名称不能为空 (request id: 20260526022040599436145qLp5F5sk)	chenruihua	Common	grok-imagine-image-lite	0	0	0	0	f	4	\N	5	default		202605260220405922686218268d9d6rJsq9N2L	20260526022040599436145qLp5F5sk	{"admin_info":{"use_channel":["4"]},"channel_id":4,"channel_name":"xai","channel_type":48,"error_code":"","error_type":"openai_error","request_path":"/v1/images/edits","status_code":400}
147	10	1779762050	5	status_code=400, 未指定模型名称，模型名称不能为空 (request id: 2026052602205059937969FcswsFou)	chenruihua	Common	grok-imagine-image-lite	0	0	0	0	f	4	\N	5	default		20260526022050543534198268d9d6I80nLb5A	2026052602205059937969FcswsFou	{"admin_info":{"use_channel":["4"]},"channel_id":4,"channel_name":"xai","channel_type":48,"error_code":"","error_type":"openai_error","request_path":"/v1/images/edits","status_code":400}
148	7	1779762138	5	status_code=400, 未指定模型名称，模型名称不能为空 (request id: 20260526022218549014009rMOxvgu8)	sso-admin	自用	grok-imagine-image-lite	0	0	0	0	f	4	\N	6	default		202605260222185135640158268d9d6pR7Czpgd	20260526022218549014009rMOxvgu8	{"admin_info":{"use_channel":["4"]},"channel_id":4,"channel_name":"xai","channel_type":48,"error_code":"","error_type":"openai_error","request_path":"/v1/images/edits","status_code":400}
149	7	1779762175	2	大小 1024x640, 品质 standard, 生成数量 1	sso-admin	自用	grok-imagine-image-lite	1	1	0	12	f	4	\N	6	default		2026052602224345066678268d9d6hQFs4ImR	2026052602224312493983coQwofF6	{"admin_info":{"use_channel":["4"]},"billing_source":"wallet","cache_ratio":1,"cache_tokens":0,"completion_ratio":1,"frt":-1000,"group_ratio":1,"model_price":-1,"model_ratio":0.25,"request_conversion":["openai_image"],"request_path":"/v1/images/generations","user_group_ratio":-1}
150	10	1779762220	5	status_code=502, Upstream request failed	chenruihua	Common	gpt-image-2	0	0	0	2	f	1	\N	5	default		202605260223371933802148268d9d6tOn6Jxxk		{"admin_info":{"use_channel":["1"]},"channel_id":1,"channel_name":"openai","channel_type":1,"error_code":"unknown_error","error_type":"openai_error","request_path":"/v1/images/edits","status_code":502}
151	7	1779763710	5	status_code=500, upstream error: do request failed	sso-admin	自用	gpt-image-2	0	0	0	0	f	5	\N	6	default		20260526024830475699198268d9d6TLocPHh4		{"admin_info":{"use_channel":["5"]},"channel_id":5,"channel_name":"deepark公益","channel_type":20,"error_code":"do_request_failed","error_type":"new_api_error","request_path":"/v1/images/edits","status_code":500}
152	7	1779763848	2	大小 1216x544, 品质 standard, 生成数量 1	sso-admin	自用	gpt-image-2	250000	1	0	72	f	5	\N	6	default		202605260249364182527508268d9d6n5FnBuXQ		{"admin_info":{"use_channel":["5"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
153	1	1779767081	2	模型测试	xiaopihong	模型测试	grok-4.3-fast	4	7	17	2	f	6	\N	0	default			20260526034439290894022BbQO5yTb	{"admin_info":{"use_channel":null},"cache_ratio":1,"cache_tokens":0,"completion_ratio":1,"frt":-1000,"group_ratio":1,"model_price":-1,"model_ratio":0.15,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","user_group_ratio":-1}
154	1	1779767134	2	模型测试	xiaopihong	模型测试	gpt-image-2	250000	0	2025688	51	f	5	\N	0	default				{"admin_info":{"use_channel":null},"cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"model_price":0.5,"model_ratio":0,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","user_group_ratio":-1}
155	10	1779767767	3	管理员增加用户额度 ＄100.000000 额度	chenruihua			0	0	0	0	f	0	\N	0					{"admin_info":{"admin_id":1,"admin_username":"xiaopihong"}}
156	7	1779767787	2		sso-admin	自用	grok-4.3-fast	5	10	25	4	t	6	\N	6	vip		202605260356237241460968268d9d6UZJCJBZ7	20260526035624361075191NOR1yzYN	{"admin_info":{"use_channel":["6"]},"billing_source":"wallet","cache_ratio":1,"cache_tokens":0,"completion_ratio":1,"frt":3151,"group_ratio":1,"model_price":-1,"model_ratio":0.15,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","stream_status":{"end_reason":"done","status":"ok"},"user_group_ratio":-1}
157	7	1779767810	2		sso-admin	自用	grok-4.3-fast	29	46	145	14	t	6	\N	6	vip		202605260356366006285068268d9d6MUcuAOnZ	20260526035636844584932bTh530Mv	{"admin_info":{"use_channel":["6"]},"billing_source":"wallet","cache_ratio":1,"cache_tokens":0,"completion_ratio":1,"frt":4205,"group_ratio":1,"model_price":-1,"model_ratio":0.15,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","stream_status":{"end_reason":"done","status":"ok"},"user_group_ratio":-1}
158	7	1779768994	5	status_code=400, Field required	sso-admin	自用	grok-4.3-fast	0	0	0	0	f	6	\N	6	vip		202605260416341516426108268d9d6CxgYcf5i	20260526041634492475130whVhWoUg	{"admin_info":{"use_channel":["6"]},"channel_id":6,"channel_name":"marybrown公益","channel_type":20,"error_code":"invalid_value","error_type":"openai_error","request_path":"/v1/images/edits","status_code":400}
159	7	1779775286	2		sso-admin	自用	mimo-v2.5-pro	306	257	355	12	t	2	\N	6	vip		202605260601144062611328268d9d6KxeeJpVQ		{"admin_info":{"use_channel":["2"]},"billing_source":"wallet","cache_ratio":1,"cache_tokens":192,"completion_ratio":1,"frt":3140,"group_ratio":1,"model_price":-1,"model_ratio":0.5,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","stream_status":{"end_reason":"done","status":"ok"},"user_group_ratio":-1}
160	10	1779786050	5	status_code=502, Upstream request failed	chenruihua	Common	gpt-image-2	0	0	0	3	f	1	\N	5	default		202605260900477014280018268d9d65DIOdEFo		{"admin_info":{"use_channel":["1"]},"channel_id":1,"channel_name":"openai","channel_type":1,"error_code":"unknown_error","error_type":"openai_error","request_path":"/v1/images/edits","status_code":502}
161	10	1779786058	5	status_code=502, Upstream request failed	chenruihua	Common	gpt-image-2	0	0	0	2	f	1	\N	5	default		202605260900564798141678268d9d6KsWHlzc8		{"admin_info":{"use_channel":["1"]},"channel_id":1,"channel_name":"openai","channel_type":1,"error_code":"unknown_error","error_type":"openai_error","request_path":"/v1/images/edits","status_code":502}
162	10	1779786128	2	大小 1024x640, 品质 standard, 生成数量 1	chenruihua	Common	gpt-image-2	250000	1	0	61	f	5	\N	5	default		202605260901078968832698268d9d60LNpn5lp		{"admin_info":{"use_channel":["5"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
163	10	1779786611	5	status_code=502, Upstream request failed	chenruihua	Common	gpt-image-2	0	0	0	360	f	1	\N	5	default		202605260904115251559448268d9d6tOP8VQ8Y		{"admin_info":{"use_channel":["1"]},"channel_id":1,"channel_name":"openai","channel_type":1,"error_code":"unknown_error","error_type":"openai_error","request_path":"/v1/images/edits","status_code":502}
164	3	1779800906	2	大小 1216x544, 品质 standard, 生成数量 1	xph-admin	自用2	gpt-image-2	250000	1	0	98	f	5	\N	2	default		202605261306488217749768268d9d6y3n75qY2		{"admin_info":{"use_channel":["5"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
165	3	1779802529	5	status_code=502, Upstream request failed	xph-admin	自用2	gpt-image-2	0	0	0	2	f	1	\N	2	default		202605261335266592102368268d9d6MBbboeGU		{"admin_info":{"use_channel":["1"]},"channel_id":1,"channel_name":"openai","channel_type":1,"error_code":"unknown_error","error_type":"openai_error","request_path":"/v1/images/edits","status_code":502}
166	3	1779802538	5	status_code=502, Upstream request failed	xph-admin	自用2	gpt-image-2	0	0	0	1	f	1	\N	2	default		202605261335373246339198268d9d6jR9z02tO		{"admin_info":{"use_channel":["1"]},"channel_id":1,"channel_name":"openai","channel_type":1,"error_code":"unknown_error","error_type":"openai_error","request_path":"/v1/images/edits","status_code":502}
167	3	1779802543	5	status_code=500, bad response status code 500	xph-admin	自用2	gpt-image-2	0	0	0	1	f	5	\N	2	default		202605261335427789689568268d9d6dlfBmYKc		{"admin_info":{"use_channel":["5"]},"channel_id":5,"channel_name":"deepark公益","channel_type":1,"error_code":"bad_response_status_code","error_type":"openai_error","request_path":"/v1/images/edits","status_code":500}
168	7	1779808269	5	status_code=500, bad response status code 500	sso-admin	自用	gpt-image-2	0	0	0	1	f	5	\N	6	default		202605261511084915490758268d9d6fAcwZEGH		{"admin_info":{"use_channel":["5"]},"channel_id":5,"channel_name":"deepark公益","channel_type":1,"error_code":"bad_response_status_code","error_type":"openai_error","request_path":"/v1/images/edits","status_code":500}
169	7	1779810463	5	status_code=502, Upstream request failed	sso-admin	自用	gpt-image-2	0	0	0	2	f	1	\N	6	default		202605261547414864090048268d9d6F03oErHD		{"admin_info":{"use_channel":["1"]},"channel_id":1,"channel_name":"openai","channel_type":1,"error_code":"unknown_error","error_type":"openai_error","request_path":"/v1/images/edits","status_code":502}
170	7	1779837381	5	status_code=502, Upstream request failed	sso-admin	自用	gpt-image-2	0	0	0	2	f	1	\N	6	default		202605262316194448394858268d9d6o5auCf36		{"admin_info":{"use_channel":["1"]},"channel_id":1,"channel_name":"openai","channel_type":1,"error_code":"unknown_error","error_type":"openai_error","request_path":"/v1/images/edits","status_code":502}
171	7	1779840009	5	status_code=403, 无权访问 codex 分组 (request id: 202605270000097813100548268d9d6iDXloUNc)	sso-admin	自用	gpt-image-2_1k	0	0	0	1	f	7	\N	6	default		202605270000087391982878268d9d6NANlRFNw	202605270000097813100548268d9d6iDXloUNc	{"admin_info":{"use_channel":["7"]},"channel_id":7,"channel_name":"abrdns公益","channel_type":1,"error_code":"","error_type":"openai_error","request_path":"/v1/images/edits","status_code":403}
172	7	1779840146	2	大小 1216x544, 品质 standard, 生成数量 1	sso-admin	自用	gpt-image-2_1k	250000	1	0	73	f	7	\N	6	default		202605270001135380288338268d9d6e9jQvS6D	202605270001141540938918268d9d6SkxBPfqG	{"admin_info":{"use_channel":["7"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
173	7	1779840908	2	大小 1024x640, 品质 standard, 生成数量 1	sso-admin	自用	gpt-image-2_1k	250000	1	0	134	f	7	\N	6	default		202605270012545425535638268d9d6BHg42ah3	202605270012557529044898268d9d6PLCGQHd6	{"admin_info":{"use_channel":["7"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
174	7	1779843798	2	大小 1024x640, 品质 standard, 生成数量 1	sso-admin	自用	gpt-image-2_1k	250000	1	0	50	f	7	\N	6	default		202605270102281246009328268d9d6IL74CGXe	202605270102293051099018268d9d62Gz2qQKF	{"admin_info":{"use_channel":["7"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
175	1	1779844522	2	模型测试	xiaopihong	模型测试	gpt-image-2_1k	250000	9	2296747	56	f	7	\N	0	default			20260527011427249234968268d9d6Vi0HEHe5	{"admin_info":{"use_channel":null},"cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"model_price":0.5,"model_ratio":0,"request_conversion":["OpenAI Compatible"],"request_path":"/v1/chat/completions","user_group_ratio":-1}
176	7	1779845682	2	大小 1216x544, 品质 standard, 生成数量 1	sso-admin	自用	gpt-image-2_1k	250000	1	0	134	f	7	\N	6	default		202605270132281764883758268d9d6Hzwp3acc	202605270132293707348118268d9d6ZyzXGmY5	{"admin_info":{"use_channel":["7"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
177	7	1779846077	2	大小 1216x544, 品质 standard, 生成数量 1	sso-admin	自用	gpt-image-2_1k	250000	1	0	73	f	7	\N	6	default		202605270140049737779568268d9d6J4fpaRjj	202605270140059623525178268d9d6Vzy3CGof	{"admin_info":{"use_channel":["7"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
178	7	1779846650	2	大小 1216x544, 品质 standard, 生成数量 1	sso-admin	自用	gpt-image-2_1k	250000	1	0	105	f	7	\N	6	default		202605270149052108547888268d9d634EahJnA	202605270149062318721228268d9d6kSRWUt5v	{"admin_info":{"use_channel":["7"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
179	10	1779847661	2	大小 1024x640, 品质 standard, 生成数量 1	chenruihua	Common	gpt-image-2_1k	250000	1	0	58	f	7	\N	5	default		202605270206437987339508268d9d6y6DzImLL	202605270206449025988488268d9d6KSoy3BSg	{"admin_info":{"use_channel":["7"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
180	10	1779847925	2	大小 1216x544, 品质 standard, 生成数量 1	chenruihua	Common	gpt-image-2_1k	250000	1	0	56	f	7	\N	5	default		202605270211091381250428268d9d6jzE1TGoe	202605270211101230032058268d9d6Yar6jsjj	{"admin_info":{"use_channel":["7"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
181	7	1779849841	2	大小 1216x544, 品质 standard, 生成数量 1	sso-admin	自用	gpt-image-2_1k	250000	1	0	67	f	7	\N	6	default		202605270242543744552318268d9d6vgTyzfz7	202605270242554356267608268d9d63gPiv4Av	{"admin_info":{"use_channel":["7"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
182	10	1779851802	2	大小 1216x544, 品质 standard, 生成数量 1	chenruihua	Common	gpt-image-2_1k	250000	1	0	57	f	7	\N	5	default		202605270315451858994148268d9d6wK8bQzZV	202605270315461134182478268d9d6iplheIhD	{"admin_info":{"use_channel":["7"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/generations","user_group_ratio":-1}
183	7	1779856601	2	大小 1216x544, 品质 standard, 生成数量 1	sso-admin	自用	gpt-image-2_1k	250000	1	0	57	f	7	\N	6	default		202605270435442584068728268d9d6bnsXG3ma	202605270435462097556378268d9d6UaCJjcim	{"admin_info":{"use_channel":["7"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
184	7	1779862040	2	大小 1024x640, 品质 standard, 生成数量 1	sso-admin	自用	gpt-image-2_1k	250000	1	0	99	f	7	\N	6	default		202605270605416703995898268d9d69T5LSwrO	202605270605428327707198268d9d6JxdoD1Go	{"admin_info":{"use_channel":["7"]},"billing_source":"wallet","cache_ratio":0,"cache_tokens":0,"completion_ratio":0,"frt":-1000,"group_ratio":1,"model_price":0.5,"model_ratio":0,"request_conversion":["openai_image"],"request_path":"/v1/images/edits","user_group_ratio":-1}
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
1	gpt-5.4		OpenAI	常规模型	1	{\n  "openai": {\n    "path": "/v1/chat/completions",\n    "method": "POST"\n  }\n}	1	1	1779290480	1779290661	\N	0
3	gpt-5.5		OpenAI	常规模型	1	{\n  "openai": {\n    "path": "/v1/chat/completions",\n    "method": "POST"\n  }\n}	1	1	1779421611	1779421611	\N	0
7	grok-4.3-fast		Grok	常规模型	5	{\n  "openai": {\n    "path": "/v1/chat/completions",\n    "method": "POST"\n  }\n}	1	1	1779766747	1779777421	2026-05-26 08:29:09.882455+00	0
6	grok-imagine-image-lite		Grok	生图模型	3	{\n  "openai": {\n    "path": "/v1/chat/completions",\n    "method": "POST"\n  }\n}	1	1	1779696304	1779777449	2026-05-26 08:29:12.362422+00	0
2	gpt-image-2		OpenAI	生图模型	4	{\n  "image-generation": {\n    "path": "/v1/images/generations",\n    "method": "POST"\n  },\n  "image-edit": "{   \\"path\\": \\"/v1/images/edits\\",   \\"method\\": \\"POST\\" }"\n}	1	1	1779421328	1779784356	\N	0
4	mimo-v2.5			常规模型	2	{\n  "openai": {\n    "path": "/v1/chat/completions",\n    "method": "POST"\n  }\n}	1	1	1779635421	1779635455	\N	0
5	mimo-v2.5-pro			常规模型	2	{\n  "openai": {\n    "path": "/v1/chat/completions",\n    "method": "POST"\n  }\n}	1	1	1779635490	1779635490	\N	0
8	gpt-image-2_1k		OpenAI	生图模型	6	{\n  "openai": {\n    "path": "/v1/chat/completions",\n    "method": "POST"\n  },\n  "image-generation": {\n    "path": "/v1/images/generations",\n    "method": "POST"\n  }\n}	1	1	1779839959	1779839959	\N	0
\.


--
-- Data for Name: options; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.options (key, value) FROM stdin;
SelfUseModeEnabled	false
DemoSiteEnabled	false
SystemName	星谱汇
RegisterEnabled	false
PasswordRegisterEnabled	false
ModelRatio	{\n  "360GPT_S2_V9": 0.8572,\n  "360gpt-pro": 0.8572,\n  "360gpt-turbo": 0.0858,\n  "360gpt-turbo-responsibility-8k": 0.8572,\n  "360gpt2-pro": 0.8572,\n  "ada": 10,\n  "babbage": 10,\n  "babbage-002": 0.2,\n  "bge-large-en": 0.13698630137,\n  "bge-large-zh": 0.13698630137,\n  "BLOOMZ-7B": 0.27397260274,\n  "chatglm_lite": 0.1429,\n  "chatglm_pro": 0.7143,\n  "chatglm_std": 0.3572,\n  "chatglm_turbo": 0.3572,\n  "chatgpt-4o-latest": 2.5,\n  "claude-3-5-haiku-20241022": 0.5,\n  "claude-3-5-sonnet-20240620": 1.5,\n  "claude-3-5-sonnet-20241022": 1.5,\n  "claude-3-7-sonnet-20250219": 1.5,\n  "claude-3-7-sonnet-20250219-thinking": 1.5,\n  "claude-3-haiku-20240307": 0.125,\n  "claude-3-opus-20240229": 7.5,\n  "claude-3-sonnet-20240229": 1.5,\n  "claude-haiku-4-5-20251001": 0.5,\n  "claude-opus-4-1-20250805": 7.5,\n  "claude-opus-4-20250514": 7.5,\n  "claude-opus-4-5-20251101": 2.5,\n  "claude-opus-4-6": 2.5,\n  "claude-opus-4-6-high": 2.5,\n  "claude-opus-4-6-low": 2.5,\n  "claude-opus-4-6-max": 2.5,\n  "claude-opus-4-6-medium": 2.5,\n  "claude-opus-4-7": 2.5,\n  "claude-opus-4-7-high": 2.5,\n  "claude-opus-4-7-low": 2.5,\n  "claude-opus-4-7-max": 2.5,\n  "claude-opus-4-7-medium": 2.5,\n  "claude-opus-4-7-xhigh": 2.5,\n  "claude-sonnet-4-20250514": 1.5,\n  "claude-sonnet-4-5-20250929": 1.5,\n  "code-davinci-edit-001": 10,\n  "command": 0.5,\n  "command-light": 0.5,\n  "command-light-nightly": 0.5,\n  "command-nightly": 0.5,\n  "command-r": 0.25,\n  "command-r-08-2024": 0.075,\n  "command-r-plus": 1.5,\n  "command-r-plus-08-2024": 1.25,\n  "curie": 10,\n  "davinci": 10,\n  "davinci-002": 1,\n  "deepseek-ai/DeepSeek-R1": 0.8,\n  "deepseek-ai/DeepSeek-R1-0528": 0.8,\n  "deepseek-ai/DeepSeek-V3-0324": 0.8,\n  "deepseek-ai/DeepSeek-V3.1": 0.8,\n  "deepseek-chat": 0.135,\n  "deepseek-coder": 0.135,\n  "deepseek-reasoner": 0.275,\n  "embedding_s1_v1": 0.0715,\n  "embedding-bert-512-v1": 0.0715,\n  "Embedding-V1": 0.13698630137,\n  "ERNIE-3.5-4K-0205": 0.821917808219,\n  "ERNIE-3.5-8K": 0.821917808219,\n  "ERNIE-3.5-8K-0205": 1.643835616438,\n  "ERNIE-3.5-8K-1222": 0.821917808219,\n  "ERNIE-4.0-8K": 8.219178082192,\n  "ERNIE-Bot-8K": 1.643835616438,\n  "ERNIE-Lite-8K-0308": 0.205479452055,\n  "ERNIE-Lite-8K-0922": 0.547945205479,\n  "ERNIE-Speed-128K": 0.27397260274,\n  "ERNIE-Speed-8K": 0.27397260274,\n  "ERNIE-Tiny-8K": 0.068493150685,\n  "gemini-1.5-flash-latest": 0.075,\n  "gemini-1.5-pro-latest": 1.25,\n  "gemini-2.0-flash": 0.05,\n  "gemini-2.5-flash": 0.15,\n  "gemini-2.5-flash-lite-preview-06-17": 0.05,\n  "gemini-2.5-flash-lite-preview-thinking-*": 0.05,\n  "gemini-2.5-flash-preview-04-17": 0.075,\n  "gemini-2.5-flash-preview-04-17-nothinking": 0.075,\n  "gemini-2.5-flash-preview-04-17-thinking": 0.075,\n  "gemini-2.5-flash-preview-05-20": 0.075,\n  "gemini-2.5-flash-preview-05-20-nothinking": 0.075,\n  "gemini-2.5-flash-preview-05-20-thinking": 0.075,\n  "gemini-2.5-flash-thinking-*": 0.075,\n  "gemini-2.5-pro": 0.625,\n  "gemini-2.5-pro-exp-03-25": 0.625,\n  "gemini-2.5-pro-preview-03-25": 0.625,\n  "gemini-2.5-pro-thinking-*": 0.625,\n  "gemini-embedding-001": 0.075,\n  "gemini-robotics-er-1.5-preview": 0.15,\n  "glm-3-turbo": 0.3572,\n  "glm-4": 7.143,\n  "glm-4-0520": 6.849315068493,\n  "glm-4-air": 0.068493150685,\n  "glm-4-airx": 0.684931506849,\n  "glm-4-alltools": 6.849315068493,\n  "glm-4-flash": 0,\n  "glm-4-long": 0.068493150685,\n  "glm-4-plus": 3.424657534247,\n  "glm-4v": 3.424657534247,\n  "glm-4v-plus": 0.684931506849,\n  "gpt-3.5-turbo": 0.25,\n  "gpt-3.5-turbo-0125": 0.25,\n  "gpt-3.5-turbo-0613": 0.75,\n  "gpt-3.5-turbo-1106": 0.5,\n  "gpt-3.5-turbo-16k": 1.5,\n  "gpt-3.5-turbo-16k-0613": 1.5,\n  "gpt-3.5-turbo-instruct": 0.75,\n  "gpt-4": 15,\n  "gpt-4-0125-preview": 5,\n  "gpt-4-0613": 15,\n  "gpt-4-1106-preview": 5,\n  "gpt-4-1106-vision-preview": 5,\n  "gpt-4-32k": 30,\n  "gpt-4-32k-0613": 30,\n  "gpt-4-all": 15,\n  "gpt-4-turbo": 5,\n  "gpt-4-turbo-2024-04-09": 5,\n  "gpt-4-turbo-preview": 5,\n  "gpt-4-vision-preview": 5,\n  "gpt-4.1": 1,\n  "gpt-4.1-2025-04-14": 1,\n  "gpt-4.1-mini": 0.2,\n  "gpt-4.1-mini-2025-04-14": 0.2,\n  "gpt-4.1-nano": 0.05,\n  "gpt-4.1-nano-2025-04-14": 0.05,\n  "gpt-4.5-preview": 37.5,\n  "gpt-4.5-preview-2025-02-27": 37.5,\n  "gpt-4o": 1.25,\n  "gpt-4o-2024-05-13": 2.5,\n  "gpt-4o-2024-08-06": 1.25,\n  "gpt-4o-2024-11-20": 1.25,\n  "gpt-4o-all": 15,\n  "gpt-4o-audio-preview": 1.25,\n  "gpt-4o-audio-preview-2024-10-01": 1.25,\n  "gpt-4o-gizmo-*": 2.5,\n  "gpt-4o-mini": 0.075,\n  "gpt-4o-mini-2024-07-18": 0.075,\n  "gpt-4o-mini-realtime-preview": 0.3,\n  "gpt-4o-mini-realtime-preview-2024-12-17": 0.3,\n  "gpt-4o-realtime-preview": 2.5,\n  "gpt-4o-realtime-preview-2024-10-01": 2.5,\n  "gpt-4o-realtime-preview-2024-12-17": 2.5,\n  "gpt-5": 0.625,\n  "gpt-5-2025-08-07": 0.625,\n  "gpt-5-chat-latest": 0.625,\n  "gpt-5-mini": 0.125,\n  "gpt-5-mini-2025-08-07": 0.125,\n  "gpt-5-nano": 0.025,\n  "gpt-5-nano-2025-08-07": 0.025,\n  "gpt-5.4": 0.75,\n  "gpt-5.5": 1,\n  "gpt-image-1": 2.5,\n  "grok-2": 1,\n  "grok-2-vision": 1,\n  "grok-3-beta": 1.5,\n  "grok-3-fast-beta": 2.5,\n  "grok-3-mini-beta": 0.15,\n  "grok-3-mini-fast-beta": 0.3,\n  "grok-4.3-fast": 0.15,\n  "grok-beta": 2.5,\n  "grok-imagine-image-lite": 0.25,\n  "grok-vision-beta": 2.5,\n  "hunyuan": 7.143,\n  "llama-3-sonar-large-32k-chat": 0,\n  "llama-3-sonar-large-32k-online": 0,\n  "llama-3-sonar-small-32k-chat": 0.1,\n  "llama-3-sonar-small-32k-online": 0.1,\n  "mimo-v2.5": 0.5,\n  "mimo-v2.5-pro": 0.5,\n  "NousResearch/Hermes-4-405B-FP8": 0.8,\n  "o1": 7.5,\n  "o1-2024-12-17": 7.5,\n  "o1-mini": 0.55,\n  "o1-mini-2024-09-12": 0.55,\n  "o1-preview": 7.5,\n  "o1-preview-2024-09-12": 7.5,\n  "o1-pro": 75,\n  "o1-pro-2025-03-19": 75,\n  "o3": 1,\n  "o3-2025-04-16": 1,\n  "o3-deep-research": 5,\n  "o3-deep-research-2025-06-26": 5,\n  "o3-mini": 0.55,\n  "o3-mini-2025-01-31": 0.55,\n  "o3-mini-2025-01-31-high": 0.55,\n  "o3-mini-2025-01-31-low": 0.55,\n  "o3-mini-2025-01-31-medium": 0.55,\n  "o3-mini-high": 0.55,\n  "o3-mini-low": 0.55,\n  "o3-mini-medium": 0.55,\n  "o3-pro": 10,\n  "o3-pro-2025-06-10": 10,\n  "o4-mini": 0.55,\n  "o4-mini-2025-04-16": 0.55,\n  "o4-mini-deep-research": 1,\n  "o4-mini-deep-research-2025-06-26": 1,\n  "openai/gpt-oss-120b": 0.5,\n  "PaLM-2": 1,\n  "qwen-plus": 10,\n  "qwen-turbo": 0.8572,\n  "Qwen/Qwen3-235B-A22B-Instruct-2507": 0.3,\n  "Qwen/Qwen3-235B-A22B-Thinking-2507": 0.6,\n  "Qwen/Qwen3-Coder-480B-A35B-Instruct-FP8": 0.8,\n  "semantic_similarity_s1_v1": 0.0715,\n  "SparkDesk-v1.1": 1.2858,\n  "SparkDesk-v2.1": 1.2858,\n  "SparkDesk-v3.1": 1.2858,\n  "SparkDesk-v3.5": 1.2858,\n  "SparkDesk-v4.0": 1.2858,\n  "tao-8k": 0.13698630137,\n  "text-ada-001": 0.2,\n  "text-babbage-001": 0.25,\n  "text-curie-001": 1,\n  "text-davinci-edit-001": 10,\n  "text-embedding-004": 0.001,\n  "text-embedding-3-large": 0.065,\n  "text-embedding-3-small": 0.01,\n  "text-embedding-ada-002": 0.05,\n  "text-embedding-v1": 0.05,\n  "text-moderation-latest": 0.1,\n  "text-moderation-stable": 0.1,\n  "text-search-ada-doc-001": 10,\n  "tts-1": 7.5,\n  "tts-1-1106": 7.5,\n  "tts-1-hd": 15,\n  "tts-1-hd-1106": 15,\n  "whisper-1": 15,\n  "yi-34b-chat-0205": 0.18,\n  "yi-34b-chat-200k": 0.864,\n  "yi-large": 1.369863013698,\n  "yi-large-preview": 1.369863013698,\n  "yi-large-rag": 1.712328767123,\n  "yi-large-rag-preview": 1.712328767123,\n  "yi-large-turbo": 0.821917808219,\n  "yi-medium": 0.171232876713,\n  "yi-medium-200k": 0.821917808219,\n  "yi-spark": 0.068493150685,\n  "yi-vision": 0.410958904109,\n  "yi-vl-plus": 0.432,\n  "zai-org/GLM-4.5-FP8": 0.8\n}
CacheRatio	{\n  "claude-3-5-haiku-20241022": 0.1,\n  "claude-3-5-sonnet-20240620": 0.1,\n  "claude-3-5-sonnet-20241022": 0.1,\n  "claude-3-7-sonnet-20250219": 0.1,\n  "claude-3-7-sonnet-20250219-thinking": 0.1,\n  "claude-3-haiku-20240307": 0.1,\n  "claude-3-opus-20240229": 0.1,\n  "claude-3-sonnet-20240229": 0.1,\n  "claude-haiku-4-5-20251001": 0.1,\n  "claude-opus-4-1-20250805": 0.1,\n  "claude-opus-4-1-20250805-thinking": 0.1,\n  "claude-opus-4-20250514": 0.1,\n  "claude-opus-4-20250514-thinking": 0.1,\n  "claude-opus-4-5-20251101": 0.1,\n  "claude-opus-4-5-20251101-thinking": 0.1,\n  "claude-opus-4-6": 0.1,\n  "claude-opus-4-6-high": 0.1,\n  "claude-opus-4-6-low": 0.1,\n  "claude-opus-4-6-max": 0.1,\n  "claude-opus-4-6-medium": 0.1,\n  "claude-opus-4-6-thinking": 0.1,\n  "claude-opus-4-7": 0.1,\n  "claude-opus-4-7-high": 0.1,\n  "claude-opus-4-7-low": 0.1,\n  "claude-opus-4-7-max": 0.1,\n  "claude-opus-4-7-medium": 0.1,\n  "claude-opus-4-7-thinking": 0.1,\n  "claude-opus-4-7-xhigh": 0.1,\n  "claude-sonnet-4-20250514": 0.1,\n  "claude-sonnet-4-20250514-thinking": 0.1,\n  "claude-sonnet-4-5-20250929": 0.1,\n  "claude-sonnet-4-5-20250929-thinking": 0.1,\n  "deepseek-chat": 0.25,\n  "deepseek-coder": 0.25,\n  "deepseek-reasoner": 0.25,\n  "gemini-3-flash-preview": 0.1,\n  "gemini-3-pro-preview": 0.1,\n  "gemini-3.1-pro-preview": 0.1,\n  "gpt-4": 0.5,\n  "gpt-4.1": 0.25,\n  "gpt-4.1-mini": 0.25,\n  "gpt-4.1-nano": 0.25,\n  "gpt-4.5-preview": 0.5,\n  "gpt-4.5-preview-2025-02-27": 0.5,\n  "gpt-4o": 0.5,\n  "gpt-4o-2024-08-06": 0.5,\n  "gpt-4o-2024-11-20": 0.5,\n  "gpt-4o-mini": 0.5,\n  "gpt-4o-mini-2024-07-18": 0.5,\n  "gpt-4o-mini-realtime-preview": 0.5,\n  "gpt-4o-realtime-preview": 0.5,\n  "gpt-5": 0.1,\n  "gpt-5-2025-08-07": 0.1,\n  "gpt-5-chat-latest": 0.1,\n  "gpt-5-mini": 0.1,\n  "gpt-5-mini-2025-08-07": 0.1,\n  "gpt-5-nano": 0.1,\n  "gpt-5-nano-2025-08-07": 0.1,\n  "o1": 0.5,\n  "o1-2024-12-17": 0.5,\n  "o1-mini": 0.5,\n  "o1-mini-2024-09-12": 0.5,\n  "o1-preview": 0.5,\n  "o1-preview-2024-09-12": 0.5,\n  "o3-mini": 0.5,\n  "o3-mini-2025-01-31": 0.5\n}
ModelPrice	{\n  "gpt-image-2_1k": 0.5,\n  "black-forest-labs/flux-1.1-pro": 0.04,\n  "dall-e-3": 0.04,\n  "gpt-4-gizmo-*": 0.1,\n  "gpt-4o-mini-tts": 0.3,\n  "gpt-image-2": 0.5,\n  "imagen-3.0-generate-002": 0.03,\n  "mj_blend": 0.1,\n  "mj_custom_zoom": 0,\n  "mj_describe": 0.05,\n  "mj_edits": 0.1,\n  "mj_high_variation": 0.1,\n  "mj_imagine": 0.1,\n  "mj_inpaint": 0,\n  "mj_low_variation": 0.1,\n  "mj_modal": 0.1,\n  "mj_pan": 0.1,\n  "mj_reroll": 0.1,\n  "mj_shorten": 0.1,\n  "mj_upload": 0.05,\n  "mj_upscale": 0.05,\n  "mj_variation": 0.1,\n  "mj_video": 0.8,\n  "mj_zoom": 0.1,\n  "sora-2": 0.3,\n  "sora-2-pro": 0.5,\n  "suno_lyrics": 0.01,\n  "suno_music": 0.1,\n  "swap_face": 0.05,\n  "veo-3.0-fast-generate-001": 0.15,\n  "veo-3.0-generate-001": 0.4,\n  "veo-3.1-fast-generate-preview": 0.15,\n  "veo-3.1-generate-preview": 0.4\n}
ImageRatio	{\n  "gpt-image-1": 2\n}
oidc.client_id	newapi
oidc.client_secret	Newapi1.
oidc.enabled	true
AudioRatio	{\n  "gpt-4o-audio-preview": 16,\n  "gpt-4o-mini-audio-preview": 66.67,\n  "gpt-4o-mini-realtime-preview": 16.67,\n  "gpt-4o-realtime-preview": 8\n}
AudioCompletionRatio	{\n  "gpt-4o-mini-realtime": 2,\n  "gpt-4o-realtime": 2\n}
billing_setting.billing_mode	{}
billing_setting.billing_expr	{}
ServerAddress	https://gateway.xph.it.com
oidc.well_known	<nil>
oidc.authorization_endpoint	https://portal.xph.it.com
passkey.rp_display_name	星谱汇
passkey.rp_id	
passkey.user_verification	preferred
passkey.attachment_preference	
passkey.origins	
CompletionRatio	{\n  "gpt-4-all": 2,\n  "gpt-4o-gizmo-*": 3,\n  "gpt-image-1": 8\n}
oidc.token_endpoint	https://api.xph.it.com/oauth/oidc/token
oidc.user_info_endpoint	https://api.xph.it.com/oauth/userinfo
HeaderNavModules	{"home":true,"console":true,"pricing":{"enabled":true,"requireAuth":false},"docs":false,"about":false}
QuotaForNewUser	1000000
CreateCacheRatio	{\n  "claude-3-5-haiku-20241022": 1.25,\n  "claude-3-5-sonnet-20240620": 1.25,\n  "claude-3-5-sonnet-20241022": 1.25,\n  "claude-3-7-sonnet-20250219": 1.25,\n  "claude-3-7-sonnet-20250219-thinking": 1.25,\n  "claude-3-haiku-20240307": 1.25,\n  "claude-3-opus-20240229": 1.25,\n  "claude-3-sonnet-20240229": 1.25,\n  "claude-haiku-4-5-20251001": 1.25,\n  "claude-opus-4-1-20250805": 1.25,\n  "claude-opus-4-1-20250805-thinking": 1.25,\n  "claude-opus-4-20250514": 1.25,\n  "claude-opus-4-20250514-thinking": 1.25,\n  "claude-opus-4-5-20251101": 1.25,\n  "claude-opus-4-5-20251101-thinking": 1.25,\n  "claude-opus-4-6": 1.25,\n  "claude-opus-4-6-high": 1.25,\n  "claude-opus-4-6-low": 1.25,\n  "claude-opus-4-6-max": 1.25,\n  "claude-opus-4-6-medium": 1.25,\n  "claude-opus-4-6-thinking": 1.25,\n  "claude-opus-4-7": 1.25,\n  "claude-opus-4-7-high": 1.25,\n  "claude-opus-4-7-low": 1.25,\n  "claude-opus-4-7-max": 1.25,\n  "claude-opus-4-7-medium": 1.25,\n  "claude-opus-4-7-thinking": 1.25,\n  "claude-opus-4-7-xhigh": 1.25,\n  "claude-sonnet-4-20250514": 1.25,\n  "claude-sonnet-4-20250514-thinking": 1.25,\n  "claude-sonnet-4-5-20250929": 1.25,\n  "claude-sonnet-4-5-20250929-thinking": 1.25\n}
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
1	gpt-5.4	default	1779289200	4	4	24295	9312	4	696	14982
2	gpt-5.4	default	1779292800	1	1	3271	2584	1	13	687
3	gpt-5.4	default	1779318000	1	1	8685	2479	1	300	6205
4	gpt-5.4	default	1779328800	4	4	63800	8684	4	2072	55114
5	gpt-5.4	default	1779332400	10	10	47312	19275	10	1100	28032
6	gpt-5.4	default	1779336000	9	9	39049	14330	9	927	24713
7	gpt-5.4	default	1779339600	16	16	122116	44445	16	3308	77663
8	gpt-image-2	default	1779418800	3	0	7232	0	0	0	0
9	gpt-image-2	default	1779436800	2	2	188411	0	0	3512	188411
10	gpt-image-2	default	1779469200	5	4	384855	0	0	7024	314386
11	gpt-image-2	default	1779508800	2	1	149494	0	0	1756	75040
12	gpt-image-2	default	1779552000	4	3	314704	0	0	2013	229147
13	gpt-5.4	default	1779588000	4	4	196719	13939	4	8156	182778
14	gpt-image-2	default	1779588000	4	3	232504	0	0	4089	193360
15	gpt-5.4	default	1779591600	1	1	7504	3062	1	168	4442
16	gpt-5.4	default	1779595200	1	1	26148	1422	1	1330	24725
17	gpt-image-2	default	1779606000	5	4	579632	0	0	6122	577704
18	gpt-image-2	default	1779609600	2	2	142310	0	0	1342	142310
19	gpt-image-2	default	1779624000	3	0	5234	0	0	0	0
20	gpt-image-2	default	1779638400	2	1	207660	0	0	671	53874
21	gpt-image-2	default	1779642000	1	0	92397	0	0	0	0
22	gpt-image-2	default	1779674400	3	2	159562	0	0	0	0
23	gpt-image-2	default	1779685200	2	0	6591	0	0	0	0
24	gpt-5.4	default	1779685200	1	1	8048	6263	1	13	1784
25	gpt-5.4	default	1779696000	1	1	3568	2157	1	40	1410
26	mimo-v2.5	default	1779696000	1	1	10984	8293	1	191	2690
27	gpt-5.5	default	1779696000	2	1	8345	2882	1	270	5053
28	gpt-image-2	default	1779703200	1	0	1911	0	0	0	0
29	grok-imagine-image-lite	default	1779703200	2	2	25301	0	0	0	0
30	grok-imagine-image-lite	default	1779724800	1	1	18049	0	0	0	0
31	grok-imagine-image-lite	default	1779728400	4	4	60048	0	0	0	0
32	mimo-v2.5	default	1779732000	2	2	10771	3398	2	662	7372
33	grok-imagine-image-lite	default	1779732000	5	5	64166	18905	2	158	915
34	gpt-image-2	default	1779760800	3	1	74661	0	0	0	0
35	grok-imagine-image-lite	default	1779760800	5	1	12797	0	0	0	0
36	grok-4.3-fast	vip	1779764400	2	2	17599	7356	2	170	10243
37	grok-4.3-fast	vip	1779768000	1	0	475	0	0	0	0
38	mimo-v2.5-pro	vip	1779775200	1	1	11624	3139	1	355	8484
39	gpt-image-2	default	1779786000	4	1	426623	0	0	0	0
40	gpt-image-2	default	1779800400	4	1	103381	0	0	0	0
41	gpt-image-2	default	1779807600	2	0	3633	0	0	0	0
42	gpt-image-2	default	1779836400	1	0	2213	0	0	0	0
43	gpt-image-2_1k	default	1779840000	3	2	208769	0	0	0	0
44	gpt-image-2_1k	default	1779843600	4	4	361979	0	0	0	0
45	gpt-image-2_1k	default	1779847200	3	3	181797	0	0	0	0
46	gpt-image-2_1k	default	1779850800	1	1	57362	0	0	0	0
47	gpt-image-2_1k	default	1779854400	1	1	57028	0	0	0	0
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
1	1	xiaopihong	gpt-5.4	1779289200	29	1	8
44	3	xph-admin	grok-imagine-image-lite	1779728400	4	4	4
2	7	sso-admin	gpt-5.4	1779289200	796	4	429
3	7	sso-admin	gpt-5.4	1779292800	31	1	10
4	9	mmx	gpt-5.4	1779318000	324	1	182
5	10	chenruihua	gpt-5.4	1779328800	3562	4	1392
6	10	chenruihua	gpt-5.4	1779332400	4276	7	908
7	7	sso-admin	gpt-5.4	1779332400	31	1	10
8	11	1061242019	gpt-5.4	1779332400	220	2	86
9	10	chenruihua	gpt-5.4	1779336000	3945	9	856
46	7	sso-admin	mimo-v2.5	1779732000	1375	2	688
45	7	sso-admin	grok-imagine-image-lite	1779732000	1471	5	370
47	7	sso-admin	grok-imagine-image-lite	1779760800	1	1	1
10	10	chenruihua	gpt-5.4	1779339600	28294	16	4483
11	10	chenruihua	gpt-image-2	1779436800	3707	2	500000
48	7	sso-admin	gpt-image-2	1779760800	1	1	250000
49	1	xiaopihong	grok-4.3-fast	1779764400	24	1	4
12	7	sso-admin	gpt-image-2	1779469200	12212	4	1000000
13	10	chenruihua	gpt-image-2	1779508800	3065	1	250000
50	1	xiaopihong	gpt-image-2	1779764400	2025688	1	250000
14	7	sso-admin	gpt-image-2	1779552000	5208	3	750000
51	7	sso-admin	grok-4.3-fast	1779764400	226	2	34
52	7	sso-admin	mimo-v2.5-pro	1779775200	612	1	306
15	7	sso-admin	gpt-image-2	1779588000	7572	3	750000
16	7	sso-admin	gpt-5.4	1779588000	22735	4	47636
17	7	sso-admin	gpt-5.4	1779591600	192	1	774
18	7	sso-admin	gpt-5.4	1779595200	1356	1	6005
53	10	chenruihua	gpt-image-2	1779786000	1	1	250000
54	3	xph-admin	gpt-image-2	1779800400	1	1	250000
19	10	chenruihua	gpt-image-2	1779606000	10460	4	1000000
20	7	sso-admin	gpt-image-2	1779609600	1610	1	250000
21	10	chenruihua	gpt-image-2	1779609600	1700	1	250000
22	1	xiaopihong	gpt-5.4	1779620400	58	2	126
23	1	xiaopihong	gpt-5.5	1779620400	31	1	96
55	7	sso-admin	gpt-image-2_1k	1779840000	2	2	500000
25	1	xiaopihong	gpt-5.5	1779624000	93	3	288
24	1	xiaopihong	gpt-5.4	1779624000	203	7	441
26	1	xiaopihong	gpt-5.5	1779631200	62	2	192
27	1	xiaopihong	gpt-5.4	1779631200	58	2	126
28	1	xiaopihong	mimo-v2.5	1779634800	264	1	132
29	1	xiaopihong	gpt-5.4	1779634800	29	1	63
30	3	xph-admin	gpt-image-2	1779638400	1204	1	250000
31	1	xiaopihong	gpt-image-2	1779674400	0	1	250000
32	7	sso-admin	gpt-image-2	1779674400	2	2	500000
33	7	sso-admin	gpt-5.4	1779685200	31	1	72
34	1	xiaopihong	gpt-5.4	1779688800	29	1	63
35	11	1061242019	mimo-v2.5	1779696000	442	1	221
36	11	1061242019	gpt-5.5	1779696000	379	1	1729
37	11	1061242019	gpt-5.4	1779696000	451	1	488
38	1	xiaopihong	mimo-v2.5	1779699600	528	2	264
39	1	xiaopihong	gpt-5.4	1779699600	58	2	126
40	1	xiaopihong	gpt-5.5	1779699600	31	1	96
41	1	xiaopihong	grok-imagine-image-lite	1779699600	69	1	17
42	7	sso-admin	grok-imagine-image-lite	1779703200	2	2	2
43	3	xph-admin	grok-imagine-image-lite	1779724800	1	1	1
57	1	xiaopihong	gpt-image-2_1k	1779843600	2296756	1	250000
56	7	sso-admin	gpt-image-2_1k	1779843600	4	4	1000000
58	10	chenruihua	gpt-image-2_1k	1779847200	2	2	500000
59	7	sso-admin	gpt-image-2_1k	1779847200	1	1	250000
60	10	chenruihua	gpt-image-2_1k	1779850800	1	1	250000
61	7	sso-admin	gpt-image-2_1k	1779854400	1	1	250000
62	7	sso-admin	gpt-image-2_1k	1779861600	1	1	250000
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
1		1779202045
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
1	3	VdTJKaMa0B4GohmK3lF1dakJKS3WzAuka0s3q3vH0vUAD5sL	1	自用	1779259578	1779259578	-1	0	t	f			0		f	\N
4	9	9nbGEbFrKliwFlismTGSsJPthOcizi6Bd8REDe9B7UGqSfSr	1	111	1779318435	1779318491	-1	-182	t	f			182		f	\N
3	7	nURIJMuLw8BX0nXZNzTDqMiipwW2mac2jIb18GX5NfhPgbym	1	自用	1779291073	1779294223	-1	-439	t	f			439		f	2026-05-20 23:56:34.402818+00
2	3	ELDZImdz8soI1BohnCZ4jppf7wYDHzilzj3vj8NONQ1f38uV	1	自用2	1779261553	1779800908	-1	-500005	t	f			500005		f	\N
5	10	9ViIJTW8JlztLjaAoHlbMWBdnPRs1Q6oD0vtKeahmLeZw3oH	1	Common	1779331869	1779851807	-1	-3007639	t	f			3007639		f	\N
6	7	ZHUHLVwIY6hmwKrEcX9mujE8YrIFiGLA1co0oTdvtLeSnXgy	1	自用	1779333698	1779862040	-1	-5805898	t	f			5805898		f	\N
7	11	lA3PG12NziNhxkOh6ROSHOASdPbbhM8971FJagDSGN1aCYct	1	111	1779334247	1779699371	-1	-2524	t	f			2524		f	\N
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
5	normal2	$2a$10$FCk51v0WiWzc6m0UBeudyeb0eJPKWxdLdGR/ATiTPOf3Ho9/goUy2		1	1	260094892@qq.com			8			c5148ea020ed7a2bf6370690034dae4e	0	0	0	default	9dKp	0	0	0	0	\N		{"gotify_priority":0,"sidebar_modules":"{\\"chat\\":{\\"chat\\":true,\\"enabled\\":true,\\"playground\\":true},\\"console\\":{\\"detail\\":true,\\"enabled\\":true,\\"log\\":true,\\"midjourney\\":true,\\"task\\":true,\\"token\\":true},\\"personal\\":{\\"enabled\\":true,\\"personal\\":true,\\"topup\\":true}}"}			1779283600	0
9	mmx	$2a$10$MtLRnt6dNhvMbvNUIlw18.t1np6C8gxxeynGZSER/4m7ncpT7NQL.		1	1	1714385126@qq.com			10			fc83c12cce3bd47b7a846f031c016661	999818	182	1	default	sN5l	0	0	0	0	\N		{"gotify_priority":0,"sidebar_modules":"{\\"chat\\":{\\"chat\\":true,\\"enabled\\":true,\\"playground\\":true},\\"console\\":{\\"detail\\":true,\\"enabled\\":true,\\"log\\":true,\\"midjourney\\":true,\\"task\\":true,\\"token\\":true},\\"personal\\":{\\"enabled\\":true,\\"personal\\":true,\\"topup\\":true}}"}			1779318292	0
3	xph-admin	$2a$10$BQuVYx4IGHlr.7//MQ7EM.D0XDayB2cRk9xzBCsFUF0hjnOCr95IS		1	1	811258683@qq.com			2			0391c4bcc41e43e6e7e25081a3167154	14499995	500005	7	default	mC7e	0	0	0	0	\N		{"gotify_priority":0,"sidebar_modules":"{\\"chat\\":{\\"chat\\":true,\\"enabled\\":true,\\"playground\\":true},\\"console\\":{\\"detail\\":true,\\"enabled\\":true,\\"log\\":true,\\"midjourney\\":true,\\"task\\":true,\\"token\\":true},\\"personal\\":{\\"enabled\\":true,\\"personal\\":true,\\"topup\\":true}}"}			1779249548	1779290248
7	sso-admin	$2a$10$Lkiv7W20KJoXQRUFp2tw1OMTrYn3WbtggK0Y199T64vzuGiUwWB0m		1	1	811258682@qq.com			1			1fe06e55d204f37324c7f8fb556282ce	5193663	5806337	49	default	x59N	0	0	0	0	\N		{"gotify_priority":0,"sidebar_modules":"{\\"chat\\":{\\"chat\\":true,\\"enabled\\":true,\\"playground\\":true},\\"console\\":{\\"detail\\":true,\\"enabled\\":true,\\"log\\":true,\\"midjourney\\":true,\\"task\\":true,\\"token\\":true},\\"personal\\":{\\"enabled\\":true,\\"personal\\":true,\\"topup\\":true}}"}			1779291044	1779849108
1	xiaopihong	$2a$10$ccqBc3RuZUPqP.KRTJRun.EcJ/1q9TZj0Djo.1Mk/W/JpVu6lcQMu	Root User	100	1							/gmLP1sxvIlHRQVCwoM8n6Tqjslk+zL5	100000000	0	0	default	es6r	0	0	0	0	\N					1779202045	1779849024
11	1061242019	$2a$10$MbGROpksdWLh/IPZuiUe5Ozgh2DqT2.AJLq7ACI0Dj/6KBK/lGeRC		1	1	1061242019@qq.com			12			9af5ae71bab9a4fff8dff866a240c49f	997476	2524	5	default	iJtP	0	0	0	0	\N		{"gotify_priority":0,"sidebar_modules":"{\\"chat\\":{\\"chat\\":true,\\"enabled\\":true,\\"playground\\":true},\\"console\\":{\\"detail\\":true,\\"enabled\\":true,\\"log\\":true,\\"midjourney\\":true,\\"task\\":true,\\"token\\":true},\\"personal\\":{\\"enabled\\":true,\\"personal\\":true,\\"topup\\":true}}"}			1779332891	0
10	chenruihua	$2a$10$nlK7BKtuUiaNjhiybBQEFuCGX/rTCr/oMLb5/TifzOh3W37rtMEP.		1	1	chenruihua@2925.com			11			1cab384d39e4b86b853948696849c9e6	57992361	3007639	48	default	gQ60	0	0	0	0	\N		{"gotify_priority":0,"sidebar_modules":"{\\"chat\\":{\\"chat\\":true,\\"enabled\\":true,\\"playground\\":true},\\"console\\":{\\"detail\\":true,\\"enabled\\":true,\\"log\\":true,\\"midjourney\\":true,\\"task\\":true,\\"token\\":true},\\"personal\\":{\\"enabled\\":true,\\"personal\\":true,\\"topup\\":true}}"}			1779331795	1779332720
4	normal1	$2a$10$8KgUFeyd3EOYNBT5QRG85.uVLZ3b6MIiUi.JUDmkwqQeDxxKELaNG		1	1	3450670290@qq.com			6			7b2f2e6a1f5b9d8cf2f236a3da399d5d	0	0	0	default	hJ0o	0	0	0	0	\N		{"gotify_priority":0,"sidebar_modules":"{\\"chat\\":{\\"chat\\":true,\\"enabled\\":true,\\"playground\\":true},\\"console\\":{\\"detail\\":true,\\"enabled\\":true,\\"log\\":true,\\"midjourney\\":true,\\"task\\":true,\\"token\\":true},\\"personal\\":{\\"enabled\\":true,\\"personal\\":true,\\"topup\\":true}}"}			1779283417	0
\.


--
-- Data for Name: vendors; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.vendors (id, name, description, icon, status, created_time, updated_time, deleted_at) FROM stdin;
1	OpenAI		OpenAI	1	1779290610	1779290610	\N
2	MiMo		MiMo	1	0	1779635384	\N
4	deepark			1	1779763628	1779763628	\N
5	marybrown			1	1779766986	1779766986	2026-05-26 08:29:30.239787+00
3	xAI		XAI	1	1779696275	1779696275	2026-05-26 08:29:36.239896+00
6	abrdns			1	1779839867	1779839867	\N
\.


--
-- Name: channels_id_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.channels_id_seq', 7, true);


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

SELECT pg_catalog.setval('public.logs_id_seq', 184, true);


--
-- Name: midjourneys_id_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.midjourneys_id_seq', 1, false);


--
-- Name: models_id_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.models_id_seq', 8, true);


--
-- Name: passkey_credentials_id_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.passkey_credentials_id_seq', 1, false);


--
-- Name: perf_metrics_id_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.perf_metrics_id_seq', 47, true);


--
-- Name: prefill_groups_id_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.prefill_groups_id_seq', 1, false);


--
-- Name: quota_data_id_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.quota_data_id_seq', 62, true);


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

SELECT pg_catalog.setval('public.tokens_id_seq', 7, true);


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

SELECT pg_catalog.setval('public.users_id_seq', 11, true);


--
-- Name: vendors_id_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.vendors_id_seq', 6, true);


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

\unrestrict pDsTkRoJ5UXhBiHvEfbQDC5GMpxTpEmF0ucff7LLvqFEhrsgLh36zJLjiXKMLko

