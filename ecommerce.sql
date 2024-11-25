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

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: banners; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.banners (
    id integer NOT NULL,
    image character varying(255) NOT NULL,
    title character varying(255) NOT NULL,
    subtitle character varying(255) NOT NULL,
    path_page character varying(255) NOT NULL
);


ALTER TABLE public.banners OWNER TO postgres;

--
-- Name: banners_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.banners_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.banners_id_seq OWNER TO postgres;

--
-- Name: banners_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.banners_id_seq OWNED BY public.banners.id;


--
-- Name: categories; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.categories (
    id integer NOT NULL,
    name character varying(50) NOT NULL,
    variant jsonb DEFAULT '{}'::jsonb NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.categories OWNER TO postgres;

--
-- Name: categories_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.categories_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.categories_id_seq OWNER TO postgres;

--
-- Name: categories_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.categories_id_seq OWNED BY public.categories.id;


--
-- Name: order_items; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.order_items (
    id integer NOT NULL,
    order_id integer,
    product_id integer,
    user_id integer,
    quantity integer NOT NULL,
    price numeric(10,2) NOT NULL,
    total numeric(10,2) NOT NULL
);


ALTER TABLE public.order_items OWNER TO postgres;

--
-- Name: order_items_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.order_items_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.order_items_id_seq OWNER TO postgres;

--
-- Name: order_items_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.order_items_id_seq OWNED BY public.order_items.id;


--
-- Name: orders; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.orders (
    id integer NOT NULL,
    user_id integer,
    total_amount numeric(10,2) NOT NULL,
    shipping_address character varying(255) DEFAULT NULL::character varying,
    shipping character varying(20) DEFAULT 'FREE'::character varying,
    status character varying(20) DEFAULT 'pending'::character varying,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT orders_status_check CHECK (((status)::text = ANY ((ARRAY['pending'::character varying, 'shipped'::character varying, 'completed'::character varying, 'canceled'::character varying])::text[])))
);


ALTER TABLE public.orders OWNER TO postgres;

--
-- Name: orders_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.orders_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.orders_id_seq OWNER TO postgres;

--
-- Name: orders_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.orders_id_seq OWNED BY public.orders.id;


--
-- Name: products; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.products (
    id integer NOT NULL,
    category_id integer,
    name character varying(100) NOT NULL,
    title character varying(100) NOT NULL,
    subtitle character varying(100) NOT NULL,
    images jsonb DEFAULT '[]'::jsonb NOT NULL,
    description text,
    price numeric(10,2) NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT images_is_array CHECK ((jsonb_typeof(images) = 'array'::text))
);


ALTER TABLE public.products OWNER TO postgres;

--
-- Name: products_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.products_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.products_id_seq OWNER TO postgres;

--
-- Name: products_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.products_id_seq OWNED BY public.products.id;


--
-- Name: products_with_flag; Type: VIEW; Schema: public; Owner: postgres
--

CREATE VIEW public.products_with_flag AS
 SELECT id,
    category_id,
    name,
    title,
    subtitle,
    images,
    description,
    price,
    created_at,
    ((CURRENT_TIMESTAMP - (created_at)::timestamp with time zone) < '30 days'::interval) AS flag
   FROM public.products;


ALTER VIEW public.products_with_flag OWNER TO postgres;

--
-- Name: ratings; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.ratings (
    id integer NOT NULL,
    order_id integer,
    product_id integer,
    user_id integer,
    rating smallint NOT NULL,
    review character varying(255),
    CONSTRAINT ratings_rating_check CHECK (((rating >= 1) AND (rating <= 5)))
);


ALTER TABLE public.ratings OWNER TO postgres;

--
-- Name: ratings_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.ratings_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.ratings_id_seq OWNER TO postgres;

--
-- Name: ratings_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.ratings_id_seq OWNED BY public.ratings.id;


--
-- Name: recomments; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.recomments (
    id integer NOT NULL,
    product_id integer
);


ALTER TABLE public.recomments OWNER TO postgres;

--
-- Name: recomments_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.recomments_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.recomments_id_seq OWNER TO postgres;

--
-- Name: recomments_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.recomments_id_seq OWNED BY public.recomments.id;


--
-- Name: sessions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.sessions (
    id integer NOT NULL,
    user_id integer,
    token text NOT NULL,
    expires_at timestamp with time zone NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.sessions OWNER TO postgres;

--
-- Name: sessions_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.sessions_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.sessions_id_seq OWNER TO postgres;

--
-- Name: sessions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.sessions_id_seq OWNED BY public.sessions.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id integer NOT NULL,
    name character varying(100) NOT NULL,
    email character varying(100),
    phone character varying(13),
    address jsonb DEFAULT '[]'::jsonb NOT NULL,
    password character varying(255) NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT address_is_array CHECK ((jsonb_typeof(address) = 'array'::text))
);


ALTER TABLE public.users OWNER TO postgres;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.users_id_seq OWNER TO postgres;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: weekly_promotions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.weekly_promotions (
    id integer NOT NULL,
    product_id integer,
    discount_percentage smallint NOT NULL,
    start_date date NOT NULL,
    end_date date NOT NULL,
    CONSTRAINT weekly_promotions_check CHECK ((end_date > start_date)),
    CONSTRAINT weekly_promotions_discount_percentage_check CHECK (((discount_percentage >= 1) AND (discount_percentage <= 100)))
);


ALTER TABLE public.weekly_promotions OWNER TO postgres;

--
-- Name: weekly_promotions_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.weekly_promotions_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.weekly_promotions_id_seq OWNER TO postgres;

--
-- Name: weekly_promotions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.weekly_promotions_id_seq OWNED BY public.weekly_promotions.id;


--
-- Name: wishlists; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.wishlists (
    id integer NOT NULL,
    user_id integer,
    product_id integer
);


ALTER TABLE public.wishlists OWNER TO postgres;

--
-- Name: wishlists_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.wishlists_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.wishlists_id_seq OWNER TO postgres;

--
-- Name: wishlists_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.wishlists_id_seq OWNED BY public.wishlists.id;


--
-- Name: banners id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.banners ALTER COLUMN id SET DEFAULT nextval('public.banners_id_seq'::regclass);


--
-- Name: categories id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.categories ALTER COLUMN id SET DEFAULT nextval('public.categories_id_seq'::regclass);


--
-- Name: order_items id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.order_items ALTER COLUMN id SET DEFAULT nextval('public.order_items_id_seq'::regclass);


--
-- Name: orders id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.orders ALTER COLUMN id SET DEFAULT nextval('public.orders_id_seq'::regclass);


--
-- Name: products id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.products ALTER COLUMN id SET DEFAULT nextval('public.products_id_seq'::regclass);


--
-- Name: ratings id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ratings ALTER COLUMN id SET DEFAULT nextval('public.ratings_id_seq'::regclass);


--
-- Name: recomments id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.recomments ALTER COLUMN id SET DEFAULT nextval('public.recomments_id_seq'::regclass);


--
-- Name: sessions id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sessions ALTER COLUMN id SET DEFAULT nextval('public.sessions_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Name: weekly_promotions id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.weekly_promotions ALTER COLUMN id SET DEFAULT nextval('public.weekly_promotions_id_seq'::regclass);


--
-- Name: wishlists id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.wishlists ALTER COLUMN id SET DEFAULT nextval('public.wishlists_id_seq'::regclass);


--
-- Data for Name: banners; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.banners (id, image, title, subtitle, path_page) FROM stdin;
1	banner1.jpg	Welcome to Our Store	Discover amazing products at great prices	/home
2	banner2.jpg	Holiday Specials	Get exclusive deals for the holiday season	/holiday-deals
3	banner3.jpg	New Arrivals	Check out the latest products in our collection	/new-arrivals
4	banner4.jpg	Flash Sale	Limited-time discounts on your favorite items	/flash-sale
5	banner5.jpg	Sign Up Today	Join now and get a special welcome offer	/signup
\.


--
-- Data for Name: categories; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.categories (id, name, variant, created_at) FROM stdin;
1	Electronics	{"color": ["red", "blue", "green"]}	2024-11-24 14:20:27.47971
2	Gaming	{"color": ["red", "blue", "green"]}	2024-11-24 14:20:27.47971
3	Wearables	{}	2024-11-24 14:20:27.47971
4	Home Office	{}	2024-11-24 14:20:27.47971
5	Fitness	{}	2024-11-24 14:20:27.47971
6	Fashion	{"size": ["S", "M", "L", "XL"], "color": ["black", "white", "red", "blue", "green"]}	2024-11-24 14:20:27.47971
7	Kitchenware	{"color": ["white", "black", "gray"], "material": ["plastic", "stainless steel", "ceramic"]}	2024-11-24 14:20:27.47971
8	Books	{}	2024-11-24 14:20:27.47971
9	Gadget	{"color": ["black", "white", "silver", "blue"], "storage": ["64GB", "128GB", "256GB", "512GB"]}	2024-11-24 14:20:27.47971
10	Beauty	{}	2024-11-24 14:20:27.47971
\.


--
-- Data for Name: order_items; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.order_items (id, order_id, product_id, user_id, quantity, price, total) FROM stdin;
1	1	1	1	2	900.00	1800.00
2	1	2	1	1	450.00	450.00
3	2	3	2	1	450.00	450.00
4	3	4	2	1	300.00	300.00
5	3	5	3	1	200.00	200.00
6	4	6	3	2	300.00	600.00
7	5	7	4	1	1000.00	1000.00
8	6	8	4	1	550.00	550.00
9	7	9	5	2	350.00	700.00
10	8	10	5	1	400.00	400.00
11	9	1	2	1	600.00	600.00
12	10	2	4	1	250.00	250.00
13	\N	15	4	1	45.00	45.00
\.


--
-- Data for Name: orders; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.orders (id, user_id, total_amount, shipping_address, shipping, status, created_at) FROM stdin;
1	1	1800.00	\N	FREE	completed	2024-11-24 14:20:27.47971
2	2	450.00	\N	FREE	pending	2024-11-24 14:20:27.47971
3	3	300.00	\N	FREE	completed	2024-11-24 14:20:27.47971
4	4	600.00	\N	FREE	canceled	2024-11-24 14:20:27.47971
5	5	1000.00	\N	FREE	completed	2024-11-24 14:20:27.47971
6	1	550.00	\N	FREE	pending	2024-11-24 14:20:27.47971
7	2	700.00	\N	FREE	completed	2024-11-24 14:20:27.47971
8	3	400.00	\N	FREE	completed	2024-11-24 14:20:27.47971
9	4	150.00	\N	FREE	completed	2024-11-24 14:20:27.47971
10	5	250.00	\N	FREE	pending	2024-11-24 14:20:27.47971
\.


--
-- Data for Name: products; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.products (id, category_id, name, title, subtitle, images, description, price, created_at) FROM stdin;
1	1	Gaming Laptop Pro	High-End Performance	For Advanced Gamers	["https://example.com/images/product1_1.jpg", "https://example.com/images/product1_2.jpg", "https://example.com/images/product1_3.jpg", "https://example.com/images/product1_4.jpg", "https://example.com/images/product1_5.jpg"]	Powerful gaming laptop with RGB keyboard	2000.00	2024-11-24 14:20:27.47971
2	1	Office Laptop	Work Anywhere	Portable and Lightweight	["https://example.com/images/product2_1.jpg", "https://example.com/images/product2_2.jpg", "https://example.com/images/product2_3.jpg", "https://example.com/images/product2_4.jpg", "https://example.com/images/product2_5.jpg"]	Lightweight laptop ideal for remote work	1200.00	2024-11-24 14:20:27.47971
3	2	Gaming Headset	Immersive Sound	7.1 Surround	["https://example.com/images/product3_1.jpg", "https://example.com/images/product3_2.jpg", "https://example.com/images/product3_3.jpg", "https://example.com/images/product3_4.jpg", "https://example.com/images/product3_5.jpg"]	7.1 surround sound gaming headset	100.00	2024-11-24 14:20:27.47971
4	2	Wireless Earbuds	Compact and Clear	Bluetooth 5.0	["https://example.com/images/product4_1.jpg", "https://example.com/images/product4_2.jpg", "https://example.com/images/product4_3.jpg", "https://example.com/images/product4_4.jpg", "https://example.com/images/product4_5.jpg"]	Compact wireless earbuds with high clarity	75.00	2024-11-24 14:20:27.47971
5	3	Fitness Smartwatch	Track Your Health	Waterproof and Durable	["https://example.com/images/product5_1.jpg", "https://example.com/images/product5_2.jpg", "https://example.com/images/product5_3.jpg", "https://example.com/images/product5_4.jpg", "https://example.com/images/product5_5.jpg"]	Smartwatch for health and fitness tracking	150.00	2024-11-24 14:20:27.47971
6	3	Luxury Smartwatch	Style and Function	Stainless Steel	["https://example.com/images/product6_1.jpg", "https://example.com/images/product6_2.jpg", "https://example.com/images/product6_3.jpg", "https://example.com/images/product6_4.jpg", "https://example.com/images/product6_5.jpg"]	Premium smartwatch with stainless steel frame	400.00	2024-11-24 14:20:27.47971
7	4	Ultra-Wide Monitor	Expand Your View	Professional 34-inch	["https://example.com/images/product7_1.jpg", "https://example.com/images/product7_2.jpg", "https://example.com/images/product7_3.jpg", "https://example.com/images/product7_4.jpg", "https://example.com/images/product7_5.jpg"]	34-inch ultra-wide monitor for multitasking	800.00	2024-11-24 14:20:27.47971
8	4	Curved Gaming Monitor	Immersive Display	27-inch	["https://example.com/images/product8_1.jpg", "https://example.com/images/product8_2.jpg", "https://example.com/images/product8_3.jpg", "https://example.com/images/product8_4.jpg", "https://example.com/images/product8_5.jpg"]	Curved monitor for immersive gaming	500.00	2024-11-24 14:20:27.47971
9	5	Ergonomic Chair	Sit Comfortably	For Long Hours	["https://example.com/images/product9_1.jpg", "https://example.com/images/product9_2.jpg", "https://example.com/images/product9_3.jpg", "https://example.com/images/product9_4.jpg", "https://example.com/images/product9_5.jpg"]	Ergonomic chair with lumbar support	300.00	2024-11-24 14:20:27.47971
10	5	Gaming Desk	Organize Your Setup	Sturdy and Spacious	["https://example.com/images/product10_1.jpg", "https://example.com/images/product10_2.jpg", "https://example.com/images/product10_3.jpg", "https://example.com/images/product10_4.jpg", "https://example.com/images/product10_5.jpg"]	Spacious desk for gaming and office use	250.00	2024-11-24 14:20:27.47971
11	6	Casual T-Shirt	Everyday Wear	Soft Cotton	["https://example.com/images/product11_1.jpg", "https://example.com/images/product11_2.jpg", "https://example.com/images/product11_3.jpg", "https://example.com/images/product11_4.jpg", "https://example.com/images/product11_5.jpg"]	Comfortable cotton t-shirt in various colors	20.00	2024-11-24 14:20:27.47971
12	6	Formal Shirt	Look Sharp	Perfect for Meetings	["https://example.com/images/product12_1.jpg", "https://example.com/images/product12_2.jpg", "https://example.com/images/product12_3.jpg", "https://example.com/images/product12_4.jpg", "https://example.com/images/product12_5.jpg"]	Formal shirt for business and events	40.00	2024-11-24 14:20:27.47971
13	6	Sports Jacket	Stay Warm	Lightweight Design	["https://example.com/images/product13_1.jpg", "https://example.com/images/product13_2.jpg", "https://example.com/images/product13_3.jpg", "https://example.com/images/product13_4.jpg", "https://example.com/images/product13_5.jpg"]	Lightweight sports jacket with breathable fabric	60.00	2024-11-24 14:20:27.47971
14	6	Boxer	Casual and Sporty	Comfortable Fit	["https://example.com/images/product14_1.jpg", "https://example.com/images/product14_2.jpg", "https://example.com/images/product14_3.jpg", "https://example.com/images/product14_4.jpg", "https://example.com/images/product14_5.jpg"]	Comfortable boxer for everyday wear	50.00	2024-11-24 14:20:27.47971
15	6	Jeans	Classic Fit	Durable and Stylish	["https://example.com/images/product15_1.jpg", "https://example.com/images/product15_2.jpg", "https://example.com/images/product15_3.jpg", "https://example.com/images/product15_4.jpg", "https://example.com/images/product15_5.jpg"]	Durable jeans in various washes	45.00	2024-11-24 14:20:27.47971
16	7	Cookware Set	Complete Kitchen	Non-Stick	["https://example.com/images/product16_1.jpg", "https://example.com/images/product16_2.jpg", "https://example.com/images/product16_3.jpg", "https://example.com/images/product16_4.jpg", "https://example.com/images/product16_5.jpg"]	Non-stick cookware set with multiple pieces	100.00	2024-11-24 14:20:27.47971
17	7	Blender	Make Smoothies	High-Speed Motor	["https://example.com/images/product17_1.jpg", "https://example.com/images/product17_2.jpg", "https://example.com/images/product17_3.jpg", "https://example.com/images/product17_4.jpg", "https://example.com/images/product17_5.jpg"]	High-speed blender for smoothies and more	80.00	2024-11-24 14:20:27.47971
18	7	Microwave Oven	Quick Meals	Compact Design	["https://example.com/images/product18_1.jpg", "https://example.com/images/product18_2.jpg", "https://example.com/images/product18_3.jpg", "https://example.com/images/product18_4.jpg", "https://example.com/images/product18_5.jpg"]	Compact microwave oven for quick meals	150.00	2024-11-24 14:20:27.47971
19	8	Novel: The Adventures	Thrilling Story	Paperback	["https://example.com/images/product19_1.jpg", "https://example.com/images/product19_2.jpg", "https://example.com/images/product19_3.jpg", "https://example.com/images/product19_4.jpg", "https://example.com/images/product19_5.jpg"]	Exciting novel with captivating plot	15.00	2024-11-24 14:20:27.47971
20	8	Programming Book	Learn Coding	JavaScript Guide	["https://example.com/images/product20_1.jpg", "https://example.com/images/product20_2.jpg", "https://example.com/images/product20_3.jpg", "https://example.com/images/product20_4.jpg", "https://example.com/images/product20_5.jpg"]	A beginner's guide to learning JavaScript	25.00	2024-11-24 14:20:27.47971
21	9	Fitness Tracker	Monitor Your Progress	Step Counter	["https://example.com/images/product21_1.jpg", "https://example.com/images/product21_2.jpg", "https://example.com/images/product21_3.jpg", "https://example.com/images/product21_4.jpg", "https://example.com/images/product21_5.jpg"]	Track your steps and calories burned	30.00	2024-11-24 14:20:27.47971
22	9	Yoga Mat	Comfort and Grip	Non-Slip Surface	["https://example.com/images/product22_1.jpg", "https://example.com/images/product22_2.jpg", "https://example.com/images/product22_3.jpg", "https://example.com/images/product22_4.jpg", "https://example.com/images/product22_5.jpg"]	Non-slip yoga mat for maximum comfort	40.00	2024-11-24 14:20:27.47971
23	9	Resistance Bands	Strengthen Muscles	Durable and Stretchable	["https://example.com/images/product23_1.jpg", "https://example.com/images/product23_2.jpg", "https://example.com/images/product23_3.jpg", "https://example.com/images/product23_4.jpg", "https://example.com/images/product23_5.jpg"]	Set of resistance bands for home workouts	25.00	2024-11-24 14:20:27.47971
24	10	LED Desk Lamp	Brighten Your Workspace	Adjustable Lighting	["https://example.com/images/product24_1.jpg", "https://example.com/images/product24_2.jpg", "https://example.com/images/product24_3.jpg", "https://example.com/images/product24_4.jpg", "https://example.com/images/product24_5.jpg"]	Adjustable LED desk lamp for study or work	45.00	2024-11-24 14:20:27.47971
25	10	Smart Table Lamp	Touch Control	Color Changing	["https://example.com/images/product25_1.jpg", "https://example.com/images/product25_2.jpg", "https://example.com/images/product25_3.jpg", "https://example.com/images/product25_4.jpg", "https://example.com/images/product25_5.jpg"]	Smart lamp with touch controls and color changing options	60.00	2024-11-24 14:20:27.47971
26	10	USB Charging Station	Power Your Devices	Multi-Port USB	["https://example.com/images/product26_1.jpg", "https://example.com/images/product26_2.jpg", "https://example.com/images/product26_3.jpg", "https://example.com/images/product26_4.jpg", "https://example.com/images/product26_5.jpg"]	Multi-port USB charging station for devices	35.00	2024-11-24 14:20:27.47971
27	1	Bluetooth Speaker	Portable Sound	Waterproof	["https://example.com/images/product27_1.jpg", "https://example.com/images/product27_2.jpg", "https://example.com/images/product27_3.jpg", "https://example.com/images/product27_4.jpg", "https://example.com/images/product27_5.jpg"]	Portable Bluetooth speaker with waterproof feature	80.00	2024-11-24 14:20:27.47971
28	1	Home Theater System	Surround Sound	4K Compatible	["https://example.com/images/product28_1.jpg", "https://example.com/images/product28_2.jpg", "https://example.com/images/product28_3.jpg", "https://example.com/images/product28_4.jpg", "https://example.com/images/product28_5.jpg"]	4K-compatible home theater system for immersive sound	500.00	2024-11-24 14:20:27.47971
29	1	Portable Air Conditioner	Cool Down	Compact and Efficient	["https://example.com/images/product29_1.jpg", "https://example.com/images/product29_2.jpg", "https://example.com/images/product29_3.jpg", "https://example.com/images/product29_4.jpg", "https://example.com/images/product29_5.jpg"]	Portable air conditioner for quick cooling	250.00	2024-11-24 14:20:27.47971
30	1	Tower Fan	Air Circulation	Energy Efficient	["https://example.com/images/product30_1.jpg", "https://example.com/images/product30_2.jpg", "https://example.com/images/product30_3.jpg", "https://example.com/images/product30_4.jpg", "https://example.com/images/product30_5.jpg"]	Energy-efficient tower fan for home use	150.00	2024-11-24 14:20:27.47971
\.


--
-- Data for Name: ratings; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.ratings (id, order_id, product_id, user_id, rating, review) FROM stdin;
1	1	1	1	5	Produk sangat bagus, performanya luar biasa!
2	2	1	2	4	Bagus, tapi baterainya cepat habis.
3	3	2	1	3	Kualitas suara oke, tapi agak berat.
4	4	2	3	5	Headphone ini benar-benar luar biasa!
5	5	3	2	4	Desainnya menarik, cukup akurat.
6	6	3	4	5	Cocok untuk tracking aktivitas harian.
7	7	4	5	4	Layarnya sangat tajam, cocok untuk desain.
8	8	4	1	3	Terlalu besar untuk meja saya.
9	9	5	3	5	Mouse ini nyaman digunakan.
10	10	5	2	4	Tombolnya enak, tapi kabelnya agak pendek.
11	1	6	4	5	Keyboard terbaik yang pernah saya beli!
12	2	6	1	5	Sangat nyaman untuk mengetik.
13	3	7	5	4	SSD sangat cepat, worth the price.
14	4	7	3	5	Transfer data cepat dan desain elegan.
15	5	8	2	5	Smartphone terbaik tahun ini!
16	6	8	4	4	Kamera bagus, tapi terlalu mahal.
17	7	9	1	3	Speaker bagus tapi kurang bass.
18	8	9	3	5	Bass mantap dan suara jernih.
19	9	10	5	5	Kamera tahan air ini sangat bagus.
20	10	10	2	4	Gambar bagus, tapi baterai kurang awet.
\.


--
-- Data for Name: recomments; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.recomments (id, product_id) FROM stdin;
\.


--
-- Data for Name: sessions; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.sessions (id, user_id, token, expires_at, created_at) FROM stdin;
1	4	71372d0b3a1789f22d3467e64edc03deeece77342cdc3e7440c93498cd68cb10	2024-11-25 14:43:28.49134+07	2024-11-24 14:43:28.491795
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users (id, name, email, phone, address, password, created_at, updated_at) FROM stdin;
1	John Doe	john.doe@example.com	1234567890	["123 Main St, Cityville", "456 Elm St, Townsville"]	password123	2024-11-24 14:20:27.47971	2024-11-24 14:20:27.47971
2	Jane Smith	jane.smith@example.com	0987654321	["456 Oak St, Townsville", "789 Cedar St, Uptown"]	password456	2024-11-24 14:20:27.47971	2024-11-24 14:20:27.47971
3	Alice Johnson	alice.johnson@example.com	1122334455	["789 Pine St, Villagetown", "321 Birch St, Downtown"]	password789	2024-11-24 14:20:27.47971	2024-11-24 14:20:27.47971
4	Bob Brown	bob.brown@example.com	1231231234	["101 Maple St, Citycenter", "222 Willow St, Greenfield"]	password101	2024-11-24 14:20:27.47971	2024-11-24 14:20:27.47971
5	Charlie Davis	charlie.davis@example.com	1456782345	["202 Birch St, Metropolis", "333 Oakwood St, Springfield"]	password202	2024-11-24 14:20:27.47971	2024-11-24 14:20:27.47971
\.


--
-- Data for Name: weekly_promotions; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.weekly_promotions (id, product_id, discount_percentage, start_date, end_date) FROM stdin;
1	1	15	2024-11-01	2024-11-07
2	2	20	2024-11-05	2024-11-10
3	3	10	2024-11-01	2024-11-14
4	4	25	2024-11-03	2024-11-10
5	5	5	2024-11-01	2024-11-10
6	6	30	2024-11-01	2024-11-30
7	7	15	2024-11-07	2024-11-14
8	8	10	2024-11-01	2024-11-07
9	9	5	2024-11-02	2024-11-09
10	10	12	2024-11-01	2024-11-15
\.


--
-- Data for Name: wishlists; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.wishlists (id, user_id, product_id) FROM stdin;
1	1	1
2	2	2
3	3	3
4	4	4
5	5	5
\.


--
-- Name: banners_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.banners_id_seq', 5, true);


--
-- Name: categories_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.categories_id_seq', 10, true);


--
-- Name: order_items_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.order_items_id_seq', 13, true);


--
-- Name: orders_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.orders_id_seq', 10, true);


--
-- Name: products_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.products_id_seq', 30, true);


--
-- Name: ratings_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.ratings_id_seq', 20, true);


--
-- Name: recomments_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.recomments_id_seq', 1, false);


--
-- Name: sessions_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.sessions_id_seq', 1, true);


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.users_id_seq', 5, true);


--
-- Name: weekly_promotions_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.weekly_promotions_id_seq', 10, true);


--
-- Name: wishlists_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.wishlists_id_seq', 5, true);


--
-- Name: banners banners_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.banners
    ADD CONSTRAINT banners_pkey PRIMARY KEY (id);


--
-- Name: categories categories_name_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.categories
    ADD CONSTRAINT categories_name_key UNIQUE (name);


--
-- Name: categories categories_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.categories
    ADD CONSTRAINT categories_pkey PRIMARY KEY (id);


--
-- Name: order_items order_items_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.order_items
    ADD CONSTRAINT order_items_pkey PRIMARY KEY (id);


--
-- Name: orders orders_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_pkey PRIMARY KEY (id);


--
-- Name: products products_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_pkey PRIMARY KEY (id);


--
-- Name: ratings ratings_order_id_product_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ratings
    ADD CONSTRAINT ratings_order_id_product_id_key UNIQUE (order_id, product_id);


--
-- Name: ratings ratings_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ratings
    ADD CONSTRAINT ratings_pkey PRIMARY KEY (id);


--
-- Name: recomments recomments_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.recomments
    ADD CONSTRAINT recomments_pkey PRIMARY KEY (id);


--
-- Name: sessions sessions_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sessions
    ADD CONSTRAINT sessions_pkey PRIMARY KEY (id);


--
-- Name: sessions sessions_token_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sessions
    ADD CONSTRAINT sessions_token_key UNIQUE (token);


--
-- Name: order_items unique_cart_item; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.order_items
    ADD CONSTRAINT unique_cart_item UNIQUE (user_id, product_id);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_phone_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_phone_key UNIQUE (phone);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: weekly_promotions weekly_promotions_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.weekly_promotions
    ADD CONSTRAINT weekly_promotions_pkey PRIMARY KEY (id);


--
-- Name: wishlists wishlists_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.wishlists
    ADD CONSTRAINT wishlists_pkey PRIMARY KEY (id);


--
-- Name: wishlists wishlists_user_id_product_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.wishlists
    ADD CONSTRAINT wishlists_user_id_product_id_key UNIQUE (user_id, product_id);


--
-- Name: order_items order_items_order_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.order_items
    ADD CONSTRAINT order_items_order_id_fkey FOREIGN KEY (order_id) REFERENCES public.orders(id);


--
-- Name: order_items order_items_product_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.order_items
    ADD CONSTRAINT order_items_product_id_fkey FOREIGN KEY (product_id) REFERENCES public.products(id);


--
-- Name: order_items order_items_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.order_items
    ADD CONSTRAINT order_items_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: orders orders_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: products products_category_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.categories(id) ON DELETE SET NULL;


--
-- Name: ratings ratings_order_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ratings
    ADD CONSTRAINT ratings_order_id_fkey FOREIGN KEY (order_id) REFERENCES public.orders(id);


--
-- Name: ratings ratings_product_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ratings
    ADD CONSTRAINT ratings_product_id_fkey FOREIGN KEY (product_id) REFERENCES public.products(id);


--
-- Name: ratings ratings_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ratings
    ADD CONSTRAINT ratings_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: recomments recomments_product_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.recomments
    ADD CONSTRAINT recomments_product_id_fkey FOREIGN KEY (product_id) REFERENCES public.products(id);


--
-- Name: sessions sessions_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sessions
    ADD CONSTRAINT sessions_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: weekly_promotions weekly_promotions_product_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.weekly_promotions
    ADD CONSTRAINT weekly_promotions_product_id_fkey FOREIGN KEY (product_id) REFERENCES public.products(id);


--
-- Name: wishlists wishlists_product_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.wishlists
    ADD CONSTRAINT wishlists_product_id_fkey FOREIGN KEY (product_id) REFERENCES public.products(id);


--
-- Name: wishlists wishlists_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.wishlists
    ADD CONSTRAINT wishlists_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

