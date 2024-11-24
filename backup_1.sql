--
-- PostgreSQL database dump
--

-- Dumped from database version 16rc1
-- Dumped by pg_dump version 16rc1

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

--
-- Name: adminpack; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS adminpack WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION adminpack; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION adminpack IS 'administrative functions for PostgreSQL';


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: Categories; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."Categories" (
    id integer NOT NULL,
    name character varying NOT NULL
);


ALTER TABLE public."Categories" OWNER TO postgres;

--
-- Name: Categories_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public."Categories_id_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public."Categories_id_seq" OWNER TO postgres;

--
-- Name: Categories_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public."Categories_id_seq" OWNED BY public."Categories".id;


--
-- Name: Instruction; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."Instruction" (
    id integer NOT NULL,
    name character varying(50) NOT NULL,
    "discordNitro" boolean DEFAULT false,
    gender character varying(10) NOT NULL
);


ALTER TABLE public."Instruction" OWNER TO postgres;

--
-- Name: Instruction_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public."Instruction_id_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public."Instruction_id_seq" OWNER TO postgres;

--
-- Name: Instruction_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public."Instruction_id_seq" OWNED BY public."Instruction".id;


--
-- Name: Menus; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."Menus" (
    id integer NOT NULL,
    "CategoryId" integer,
    stock integer NOT NULL,
    price integer NOT NULL,
    "creatAt" character varying NOT NULL
);


ALTER TABLE public."Menus" OWNER TO postgres;

--
-- Name: Menus_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public."Menus_id_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public."Menus_id_seq" OWNER TO postgres;

--
-- Name: Menus_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public."Menus_id_seq" OWNED BY public."Menus".id;


--
-- Name: Categories id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Categories" ALTER COLUMN id SET DEFAULT nextval('public."Categories_id_seq"'::regclass);


--
-- Name: Instruction id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Instruction" ALTER COLUMN id SET DEFAULT nextval('public."Instruction_id_seq"'::regclass);


--
-- Name: Menus id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Menus" ALTER COLUMN id SET DEFAULT nextval('public."Menus_id_seq"'::regclass);


--
-- Data for Name: Categories; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."Categories" (id, name) FROM stdin;
\.


--
-- Data for Name: Instruction; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."Instruction" (id, name, "discordNitro", gender) FROM stdin;
\.


--
-- Data for Name: Menus; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."Menus" (id, "CategoryId", stock, price, "creatAt") FROM stdin;
\.


--
-- Name: Categories_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public."Categories_id_seq"', 1, false);


--
-- Name: Instruction_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public."Instruction_id_seq"', 1, false);


--
-- Name: Menus_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public."Menus_id_seq"', 1, false);


--
-- Name: Categories Categories_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Categories"
    ADD CONSTRAINT "Categories_pkey" PRIMARY KEY (id);


--
-- Name: Instruction Instruction_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Instruction"
    ADD CONSTRAINT "Instruction_pkey" PRIMARY KEY (id);


--
-- Name: Menus Menus_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Menus"
    ADD CONSTRAINT "Menus_pkey" PRIMARY KEY (id);


--
-- Name: Menus Menus_CategoryId_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Menus"
    ADD CONSTRAINT "Menus_CategoryId_fkey" FOREIGN KEY ("CategoryId") REFERENCES public."Categories"(id);


--
-- PostgreSQL database dump complete
--

