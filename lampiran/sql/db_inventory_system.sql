--
-- PostgreSQL database dump
--

-- Dumped from database version 18.1
-- Dumped by pg_dump version 18.1

-- Started on 2026-01-03 14:21:21

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- TOC entry 5123 (class 0 OID 0)
-- Dependencies: 4
-- Name: SCHEMA public; Type: COMMENT; Schema: -; Owner: pg_database_owner
--

COMMENT ON SCHEMA public IS 'standard public schema';


--
-- TOC entry 866 (class 1247 OID 17721)
-- Name: user_role; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.user_role AS ENUM (
    'super_admin',
    'admin',
    'staff'
);


ALTER TYPE public.user_role OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- TOC entry 223 (class 1259 OID 17767)
-- Name: categories; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.categories (
    id bigint NOT NULL,
    name character varying(100) NOT NULL,
    description text,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.categories OWNER TO postgres;

--
-- TOC entry 222 (class 1259 OID 17766)
-- Name: categories_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.categories_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.categories_id_seq OWNER TO postgres;

--
-- TOC entry 5124 (class 0 OID 0)
-- Dependencies: 222
-- Name: categories_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.categories_id_seq OWNED BY public.categories.id;


--
-- TOC entry 229 (class 1259 OID 17804)
-- Name: products; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.products (
    id bigint NOT NULL,
    sku character varying(50) NOT NULL,
    name character varying(150) NOT NULL,
    description text,
    stock integer DEFAULT 0 NOT NULL,
    price numeric(15,2) DEFAULT 0.00 NOT NULL,
    category_id bigint,
    warehouse_id bigint,
    rack_id bigint,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT products_stock_check CHECK ((stock >= 0))
);


ALTER TABLE public.products OWNER TO postgres;

--
-- TOC entry 228 (class 1259 OID 17803)
-- Name: products_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.products_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.products_id_seq OWNER TO postgres;

--
-- TOC entry 5125 (class 0 OID 0)
-- Dependencies: 228
-- Name: products_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.products_id_seq OWNED BY public.products.id;


--
-- TOC entry 227 (class 1259 OID 17793)
-- Name: racks; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.racks (
    id bigint NOT NULL,
    name character varying(50) NOT NULL,
    category character varying(50),
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.racks OWNER TO postgres;

--
-- TOC entry 226 (class 1259 OID 17792)
-- Name: racks_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.racks_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.racks_id_seq OWNER TO postgres;

--
-- TOC entry 5126 (class 0 OID 0)
-- Dependencies: 226
-- Name: racks_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.racks_id_seq OWNED BY public.racks.id;


--
-- TOC entry 233 (class 1259 OID 17857)
-- Name: sale_details; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.sale_details (
    id bigint NOT NULL,
    sale_id bigint NOT NULL,
    product_id bigint NOT NULL,
    quantity integer NOT NULL,
    unit_price numeric(15,2) NOT NULL,
    subtotal numeric(15,2) NOT NULL,
    CONSTRAINT sale_details_quantity_check CHECK ((quantity > 0))
);


ALTER TABLE public.sale_details OWNER TO postgres;

--
-- TOC entry 232 (class 1259 OID 17856)
-- Name: sale_details_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.sale_details_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.sale_details_id_seq OWNER TO postgres;

--
-- TOC entry 5127 (class 0 OID 0)
-- Dependencies: 232
-- Name: sale_details_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.sale_details_id_seq OWNED BY public.sale_details.id;


--
-- TOC entry 231 (class 1259 OID 17840)
-- Name: sales; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.sales (
    id bigint NOT NULL,
    user_id bigint,
    transaction_date timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    total_amount numeric(15,2) DEFAULT 0.00 NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.sales OWNER TO postgres;

--
-- TOC entry 230 (class 1259 OID 17839)
-- Name: sales_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.sales_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.sales_id_seq OWNER TO postgres;

--
-- TOC entry 5128 (class 0 OID 0)
-- Dependencies: 230
-- Name: sales_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.sales_id_seq OWNED BY public.sales.id;


--
-- TOC entry 221 (class 1259 OID 17746)
-- Name: sessions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.sessions (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id bigint NOT NULL,
    token uuid DEFAULT gen_random_uuid() NOT NULL,
    ip_address character varying(45),
    user_agent text,
    is_revoked boolean DEFAULT false,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    expired_at timestamp with time zone NOT NULL,
    revoked_at timestamp with time zone
);


ALTER TABLE public.sessions OWNER TO postgres;

--
-- TOC entry 220 (class 1259 OID 17728)
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id bigint NOT NULL,
    username character varying(50) NOT NULL,
    email character varying(100) NOT NULL,
    password_hash character varying(255) NOT NULL,
    role public.user_role DEFAULT 'staff'::public.user_role NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.users OWNER TO postgres;

--
-- TOC entry 219 (class 1259 OID 17727)
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.users_id_seq OWNER TO postgres;

--
-- TOC entry 5129 (class 0 OID 0)
-- Dependencies: 219
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- TOC entry 225 (class 1259 OID 17780)
-- Name: warehouses; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.warehouses (
    id bigint NOT NULL,
    name character varying(100) NOT NULL,
    location character varying(255),
    description text,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.warehouses OWNER TO postgres;

--
-- TOC entry 224 (class 1259 OID 17779)
-- Name: warehouses_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.warehouses_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.warehouses_id_seq OWNER TO postgres;

--
-- TOC entry 5130 (class 0 OID 0)
-- Dependencies: 224
-- Name: warehouses_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.warehouses_id_seq OWNED BY public.warehouses.id;


--
-- TOC entry 4901 (class 2604 OID 17770)
-- Name: categories id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.categories ALTER COLUMN id SET DEFAULT nextval('public.categories_id_seq'::regclass);


--
-- TOC entry 4910 (class 2604 OID 17807)
-- Name: products id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.products ALTER COLUMN id SET DEFAULT nextval('public.products_id_seq'::regclass);


--
-- TOC entry 4907 (class 2604 OID 17796)
-- Name: racks id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.racks ALTER COLUMN id SET DEFAULT nextval('public.racks_id_seq'::regclass);


--
-- TOC entry 4919 (class 2604 OID 17860)
-- Name: sale_details id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sale_details ALTER COLUMN id SET DEFAULT nextval('public.sale_details_id_seq'::regclass);


--
-- TOC entry 4915 (class 2604 OID 17843)
-- Name: sales id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sales ALTER COLUMN id SET DEFAULT nextval('public.sales_id_seq'::regclass);


--
-- TOC entry 4893 (class 2604 OID 17731)
-- Name: users id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- TOC entry 4904 (class 2604 OID 17783)
-- Name: warehouses id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.warehouses ALTER COLUMN id SET DEFAULT nextval('public.warehouses_id_seq'::regclass);


--
-- TOC entry 5107 (class 0 OID 17767)
-- Dependencies: 223
-- Data for Name: categories; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.categories (id, name, description, created_at, updated_at) VALUES (1, 'Elektronik', 'Barang-barang gadget dan elektronik rumah', '2026-01-01 12:22:55.081405+07', '2026-01-01 12:22:55.081405+07');
INSERT INTO public.categories (id, name, description, created_at, updated_at) VALUES (2, 'Furniture', 'Perabot rumah tangga dan perlengkapan kantor', '2026-01-03 13:19:33.150512+07', '2026-01-03 13:19:33.150512+07');
INSERT INTO public.categories (id, name, description, created_at, updated_at) VALUES (3, 'Fashion', 'Pakaian, sepatu, dan aksesoris fashion pria/wanita', '2026-01-03 13:19:39.817193+07', '2026-01-03 13:19:39.817193+07');
INSERT INTO public.categories (id, name, description, created_at, updated_at) VALUES (4, 'F&B', 'Makanan ringan, minuman kemasan, dan bahan pokok', '2026-01-03 13:19:44.404442+07', '2026-01-03 13:19:44.404442+07');
INSERT INTO public.categories (id, name, description, created_at, updated_at) VALUES (5, 'Stationery', 'Alat tulis kantor, kertas, dan perlengkapan sekolah', '2026-01-03 13:19:53.553754+07', '2026-01-03 13:19:53.553754+07');
INSERT INTO public.categories (id, name, description, created_at, updated_at) VALUES (6, 'Hiking & Shi', 'Alat perlengkapan untuk pergi mendaki gunung', '2026-01-03 13:20:30.452691+07', '2026-01-03 13:32:03.648323+07');
INSERT INTO public.categories (id, name, description, created_at, updated_at) VALUES (8, 'Nuclear', 'Experiment', '2026-01-03 13:59:21.672272+07', '2026-01-03 13:59:21.672272+07');


--
-- TOC entry 5113 (class 0 OID 17804)
-- Dependencies: 229
-- Data for Name: products; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.products (id, sku, name, description, stock, price, category_id, warehouse_id, rack_id, created_at, updated_at) VALUES (1, 'LAPTOP-001', 'Macbook Pro M2', 'Laptop mahal', 7, 20000000.00, 1, 1, 1, '2026-01-01 13:41:07.423908+07', '2026-01-01 13:51:17.966375+07');
INSERT INTO public.products (id, sku, name, description, stock, price, category_id, warehouse_id, rack_id, created_at, updated_at) VALUES (7, 'FUR-CHR-002', 'Kursi Ergonomis Herman Miller', 'Kursi kantor premium dengan support tulang belakang terbaik', 10, 12000000.00, 2, 1, 1, '2026-01-03 13:58:10.289691+07', '2026-01-03 13:58:10.289691+07');
INSERT INTO public.products (id, sku, name, description, stock, price, category_id, warehouse_id, rack_id, created_at, updated_at) VALUES (8, 'FSH-SHOE-003', 'Nike Air Jordan 1 High', 'Sepatu sneakers klasik warna merah hitam', 20, 2800000.00, 3, 1, 1, '2026-01-03 13:59:50.740302+07', '2026-01-03 13:59:50.740302+07');
INSERT INTO public.products (id, sku, name, description, stock, price, category_id, warehouse_id, rack_id, created_at, updated_at) VALUES (9, 'FNB-COF-004', 'Arabica Coffee Beans 1kg', 'Biji kopi asli dari dataran tinggi Gayo, Aceh', 50, 150000.00, 4, 1, 1, '2026-01-03 13:59:55.234506+07', '2026-01-03 13:59:55.234506+07');
INSERT INTO public.products (id, sku, name, description, stock, price, category_id, warehouse_id, rack_id, created_at, updated_at) VALUES (10, 'ATK-BOOK-005', 'Moleskine Classic Notebook', 'Buku catatan hardcover warna hitam, kertas polos', 100, 350000.00, 5, 1, 1, '2026-01-03 14:00:00.869323+07', '2026-01-03 14:00:00.869323+07');
INSERT INTO public.products (id, sku, name, description, stock, price, category_id, warehouse_id, rack_id, created_at, updated_at) VALUES (11, 'NUC-007', 'Nuclear Momentums', 'Experimental', 1, 1000000000.00, 6, 5, 1, '2026-01-03 14:03:14.410736+07', '2026-01-03 14:03:14.410736+07');
INSERT INTO public.products (id, sku, name, description, stock, price, category_id, warehouse_id, rack_id, created_at, updated_at) VALUES (6, 'LPT-GAM-001', 'Asus ROG Zephyrus G14', 'Laptop gaming compact dengan prosesor Ryzen 9 dan RTX 4060', 4, 24500000.00, 1, 1, 1, '2026-01-03 13:57:25.578178+07', '2026-01-03 14:05:45.849958+07');


--
-- TOC entry 5111 (class 0 OID 17793)
-- Dependencies: 227
-- Data for Name: racks; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.racks (id, name, category, created_at, updated_at) VALUES (1, 'Rak A', 'Elektronik', '2026-01-01 13:24:02.119216+07', '2026-01-01 13:24:02.119216+07');
INSERT INTO public.racks (id, name, category, created_at, updated_at) VALUES (2, 'Rak B', 'Peripheral', '2026-01-01 13:24:47.074201+07', '2026-01-01 13:24:47.074201+07');
INSERT INTO public.racks (id, name, category, created_at, updated_at) VALUES (4, 'Rak D', 'Hiking Shi', '2026-01-03 13:45:51.73253+07', '2026-01-03 13:45:51.73253+07');
INSERT INTO public.racks (id, name, category, created_at, updated_at) VALUES (5, 'Rak E', 'Food', '2026-01-03 13:46:13.133031+07', '2026-01-03 13:46:13.133031+07');
INSERT INTO public.racks (id, name, category, created_at, updated_at) VALUES (6, 'Rak F', 'Beverage', '2026-01-03 13:46:20.422828+07', '2026-01-03 13:46:20.422828+07');
INSERT INTO public.racks (id, name, category, created_at, updated_at) VALUES (3, 'Rak C', 'Sepatu fashion', '2026-01-03 13:45:10.948064+07', '2026-01-03 13:46:51.017261+07');


--
-- TOC entry 5117 (class 0 OID 17857)
-- Dependencies: 233
-- Data for Name: sale_details; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.sale_details (id, sale_id, product_id, quantity, unit_price, subtotal) VALUES (1, 1, 1, 3, 20000000.00, 60000000.00);
INSERT INTO public.sale_details (id, sale_id, product_id, quantity, unit_price, subtotal) VALUES (2, 3, 6, 1, 24500000.00, 24500000.00);


--
-- TOC entry 5115 (class 0 OID 17840)
-- Dependencies: 231
-- Data for Name: sales; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.sales (id, user_id, transaction_date, total_amount, created_at) VALUES (1, 3, '2026-01-01 13:51:17.966375+07', 60000000.00, '2026-01-01 13:51:17.966375+07');
INSERT INTO public.sales (id, user_id, transaction_date, total_amount, created_at) VALUES (3, 8, '2026-01-03 14:05:45.849958+07', 24500000.00, '2026-01-03 14:05:45.849958+07');


--
-- TOC entry 5105 (class 0 OID 17746)
-- Dependencies: 221
-- Data for Name: sessions; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.sessions (id, user_id, token, ip_address, user_agent, is_revoked, created_at, expired_at, revoked_at) VALUES ('05fba59a-fb96-467c-8fb5-0b743d3aa11d', 2, '8de66eeb-f33a-4c8a-a31b-7baad9437142', NULL, NULL, false, '2025-12-30 13:03:32.882226+07', '2025-12-31 13:03:32.881235+07', NULL);
INSERT INTO public.sessions (id, user_id, token, ip_address, user_agent, is_revoked, created_at, expired_at, revoked_at) VALUES ('dc121794-6a62-4963-b425-5f6c3bc6da4d', 2, 'a751951f-e345-4604-adea-6e66956cd731', NULL, NULL, false, '2025-12-30 13:40:51.677069+07', '2025-12-31 13:40:51.673347+07', NULL);
INSERT INTO public.sessions (id, user_id, token, ip_address, user_agent, is_revoked, created_at, expired_at, revoked_at) VALUES ('22b8cc6a-e6fa-4657-930c-3b0d468968aa', 3, '37f4a42c-757d-4827-82c0-a89ed1a3ba74', NULL, NULL, false, '2025-12-30 13:44:20.268478+07', '2025-12-31 13:44:20.267854+07', NULL);
INSERT INTO public.sessions (id, user_id, token, ip_address, user_agent, is_revoked, created_at, expired_at, revoked_at) VALUES ('c4858340-a605-40a4-9106-469b1e9b1191', 4, '083560eb-35b4-4efd-8863-8978acf510dc', NULL, NULL, false, '2025-12-31 13:53:49.545997+07', '2026-01-01 13:53:49.542045+07', NULL);
INSERT INTO public.sessions (id, user_id, token, ip_address, user_agent, is_revoked, created_at, expired_at, revoked_at) VALUES ('1b8c41bf-c67d-46f0-8468-5ebd5be9703a', 2, '4d4cf071-6764-405f-9594-24cac439e5b2', NULL, NULL, false, '2026-01-01 12:22:24.23645+07', '2026-01-02 12:22:24.232909+07', NULL);
INSERT INTO public.sessions (id, user_id, token, ip_address, user_agent, is_revoked, created_at, expired_at, revoked_at) VALUES ('ea6aef14-4821-4866-8fbc-58e234ad34ad', 3, '98626510-cee8-4365-ab11-4110e911a647', NULL, NULL, false, '2026-01-01 12:25:35.56168+07', '2026-01-02 12:25:35.561024+07', NULL);
INSERT INTO public.sessions (id, user_id, token, ip_address, user_agent, is_revoked, created_at, expired_at, revoked_at) VALUES ('7023f6f5-3770-47c5-9e89-c5eff4125ab0', 2, 'f369dd99-c090-426e-8fcc-6117500ba09c', NULL, NULL, false, '2026-01-01 12:51:09.605704+07', '2026-01-02 12:51:09.604378+07', NULL);
INSERT INTO public.sessions (id, user_id, token, ip_address, user_agent, is_revoked, created_at, expired_at, revoked_at) VALUES ('4d75c648-cb07-4058-a5f1-c2ba40599c1f', 2, 'd2a9e330-834b-41d8-8119-75560cb59182', NULL, NULL, false, '2026-01-01 12:51:17.405045+07', '2026-01-02 12:51:17.403977+07', NULL);
INSERT INTO public.sessions (id, user_id, token, ip_address, user_agent, is_revoked, created_at, expired_at, revoked_at) VALUES ('b35b7486-b4c4-4305-b262-6ae695e55e09', 3, 'dbe5f379-6878-47c2-b08a-7f982a44bb88', NULL, NULL, false, '2026-01-01 13:00:54.251921+07', '2026-01-02 13:00:54.250625+07', NULL);
INSERT INTO public.sessions (id, user_id, token, ip_address, user_agent, is_revoked, created_at, expired_at, revoked_at) VALUES ('07a65805-8422-4777-baf5-36df28613537', 3, 'bf8f8227-c8be-40e5-aec8-1c61c91391fb', NULL, NULL, false, '2026-01-01 13:01:20.757218+07', '2026-01-02 13:01:20.756862+07', NULL);
INSERT INTO public.sessions (id, user_id, token, ip_address, user_agent, is_revoked, created_at, expired_at, revoked_at) VALUES ('25757c19-5a1e-4694-b379-45642866c139', 2, 'c070de20-3e3b-4461-963f-58f9bf039173', NULL, NULL, false, '2026-01-01 13:21:41.576342+07', '2026-01-02 13:21:41.574583+07', NULL);
INSERT INTO public.sessions (id, user_id, token, ip_address, user_agent, is_revoked, created_at, expired_at, revoked_at) VALUES ('07fab886-7fe7-4562-bdcb-1f798477f5bf', 3, '3c768736-5769-4fb8-aa66-f53db07650ab', NULL, NULL, false, '2026-01-01 13:25:26.919009+07', '2026-01-02 13:25:26.917663+07', NULL);
INSERT INTO public.sessions (id, user_id, token, ip_address, user_agent, is_revoked, created_at, expired_at, revoked_at) VALUES ('e3ec002e-7320-4121-a736-fbbf0f238eb4', 2, 'e1e3bd95-72c1-4974-9f73-9adc46987057', NULL, NULL, false, '2026-01-02 20:32:11.037736+07', '2026-01-03 20:32:11.033009+07', NULL);
INSERT INTO public.sessions (id, user_id, token, ip_address, user_agent, is_revoked, created_at, expired_at, revoked_at) VALUES ('aea5af56-3d6d-4b26-a191-8249d34dcb92', 3, 'faf391f7-c9fa-44b5-a35c-09805dcd30b9', NULL, NULL, false, '2026-01-02 20:33:07.598534+07', '2026-01-03 20:33:07.597377+07', NULL);
INSERT INTO public.sessions (id, user_id, token, ip_address, user_agent, is_revoked, created_at, expired_at, revoked_at) VALUES ('b92f8db5-8c07-4a15-8f4f-50f1918b0290', 6, 'd5da5c03-bef6-4e5d-9916-18ccfcc0889e', NULL, NULL, false, '2026-01-03 12:58:27.75954+07', '2026-01-04 12:58:27.756366+07', NULL);
INSERT INTO public.sessions (id, user_id, token, ip_address, user_agent, is_revoked, created_at, expired_at, revoked_at) VALUES ('25c92201-214e-46d1-b4a5-e656f9ff8035', 8, '7e0d26e1-09f0-46cf-88d0-5b69674b96be', NULL, NULL, false, '2026-01-03 13:01:59.208752+07', '2026-01-04 13:01:59.20785+07', NULL);
INSERT INTO public.sessions (id, user_id, token, ip_address, user_agent, is_revoked, created_at, expired_at, revoked_at) VALUES ('1041a731-30fc-45ad-aa06-70854bd1ab69', 7, 'da60c594-6fa6-4444-8098-325e66eb5bde', NULL, NULL, true, '2026-01-03 13:00:54.604766+07', '2026-01-04 13:00:54.604266+07', NULL);
INSERT INTO public.sessions (id, user_id, token, ip_address, user_agent, is_revoked, created_at, expired_at, revoked_at) VALUES ('cae7f8c5-04fe-462b-b3d5-8c6460b06b91', 7, '3f20485c-75ce-49ae-a285-eb8fc2790fcb', NULL, NULL, false, '2026-01-03 13:04:53.05046+07', '2026-01-04 13:04:53.049632+07', NULL);


--
-- TOC entry 5104 (class 0 OID 17728)
-- Dependencies: 220
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.users (id, username, email, password_hash, role, created_at, updated_at) VALUES (1, 'SuperAdmin', 'superadmin@inventory.com', '$2a$10$2.d.w/M7.u7.t.k.h.0.u.e.r.0.h.a.s.h.e.d.p.a.s.s.w.o.r.d', 'super_admin', '2025-12-30 12:11:44.586383+07', '2025-12-30 12:11:44.586383+07');
INSERT INTO public.users (id, username, email, password_hash, role, created_at, updated_at) VALUES (2, 'babon1', 'babon1@mail.com', '$2a$10$NmDm9PjDYrgx5VdoGL.39u4KCmIHLmhEpU69s4OEXJZFgeaMNP9Zu', 'super_admin', '2025-12-30 13:01:24.849845+07', '2025-12-30 13:01:24.849845+07');
INSERT INTO public.users (id, username, email, password_hash, role, created_at, updated_at) VALUES (3, '9udang', 'udang@mail.com', '$2a$10$jFCSoa14q567j4PS3h5ch.nTZvNvLt59KlryDnm6tZsVO8H2sYC7u', 'staff', '2025-12-30 13:41:45.459416+07', '2025-12-30 13:41:45.459416+07');
INSERT INTO public.users (id, username, email, password_hash, role, created_at, updated_at) VALUES (4, 'jamal1231', 'jj@mail.com', '$2a$10$4uwTu0BU47OZu65ogP9Dn.U3T3tkrA7R5kimoF/VN23VfswwXqW6G', 'admin', '2025-12-31 13:53:10.183059+07', '2026-01-02 20:35:10.131284+07');
INSERT INTO public.users (id, username, email, password_hash, role, created_at, updated_at) VALUES (6, 'anakbos1', 'abos1@mail.com', '$2a$10$BqEl8JErwHMzO6H6sKI3t.9kBEE9DKixXgc0NVLSOGJ.XP9Ihcs6a', 'super_admin', '2026-01-03 12:52:16.554488+07', '2026-01-03 12:52:16.554488+07');
INSERT INTO public.users (id, username, email, password_hash, role, created_at, updated_at) VALUES (7, 'bawahanbos1', 'babos1@mail.com', '$2a$10$VOeUb5Y12dbdDMAEjp4J5..t8Pz8cNIuoQY6yW/WGkLINYd1GTm2y', 'admin', '2026-01-03 12:52:45.311277+07', '2026-01-03 12:52:45.311277+07');
INSERT INTO public.users (id, username, email, password_hash, role, created_at, updated_at) VALUES (8, 'staff_jam', 'jamers@mail.com', '$2a$10$XV8vsfjIAbE89kqoL9xKE..QsAyFksdHKoJ3n/NBoLG8JksN2YxgK', 'staff', '2026-01-03 12:53:08.123156+07', '2026-01-03 12:53:08.123156+07');


--
-- TOC entry 5109 (class 0 OID 17780)
-- Dependencies: 225
-- Data for Name: warehouses; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.warehouses (id, name, location, description, created_at, updated_at) VALUES (1, 'Warehouse A', 'Jl. Ujung Harapan Kel, RT.001/RW.015, Kaliabang Tengah', 'Tempat penitipan Barang-barang', '2026-01-01 12:58:51.281544+07', '2026-01-01 12:58:51.281544+07');
INSERT INTO public.warehouses (id, name, location, description, created_at, updated_at) VALUES (2, 'Warehouse B', 'Jl. Ujung Tanduk Kel, RT.003/RW.012, Harapan Jaya', 'Tempat penitipan Barang-barang', '2026-01-01 13:00:31.698896+07', '2026-01-01 13:00:31.698896+07');
INSERT INTO public.warehouses (id, name, location, description, created_at, updated_at) VALUES (3, 'Warehouse C', 'Jakarta Pusat', 'Cabang Jakarta', '2026-01-03 13:38:07.285183+07', '2026-01-03 13:38:07.285183+07');
INSERT INTO public.warehouses (id, name, location, description, created_at, updated_at) VALUES (4, 'Warehouse D', 'Aceh, Sumatra Utara', 'Cabang Sumatra bagian Utara', '2026-01-03 13:38:33.222884+07', '2026-01-03 13:38:33.222884+07');
INSERT INTO public.warehouses (id, name, location, description, created_at, updated_at) VALUES (6, 'Warehouse F', 'Surabaya, Jawa Timur', 'Cabang Jawa Timur', '2026-01-03 13:39:21.050752+07', '2026-01-03 13:39:21.050752+07');
INSERT INTO public.warehouses (id, name, location, description, created_at, updated_at) VALUES (5, 'Warehouse E', 'Padang, Sumatra Barat', 'Cabang Sumatra bagian Barat', '2026-01-03 13:38:49.710243+07', '2026-01-03 13:41:18.066554+07');


--
-- TOC entry 5131 (class 0 OID 0)
-- Dependencies: 222
-- Name: categories_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.categories_id_seq', 8, true);


--
-- TOC entry 5132 (class 0 OID 0)
-- Dependencies: 228
-- Name: products_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.products_id_seq', 11, true);


--
-- TOC entry 5133 (class 0 OID 0)
-- Dependencies: 226
-- Name: racks_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.racks_id_seq', 7, true);


--
-- TOC entry 5134 (class 0 OID 0)
-- Dependencies: 232
-- Name: sale_details_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.sale_details_id_seq', 2, true);


--
-- TOC entry 5135 (class 0 OID 0)
-- Dependencies: 230
-- Name: sales_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.sales_id_seq', 5, true);


--
-- TOC entry 5136 (class 0 OID 0)
-- Dependencies: 219
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.users_id_seq', 8, true);


--
-- TOC entry 5137 (class 0 OID 0)
-- Dependencies: 224
-- Name: warehouses_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.warehouses_id_seq', 7, true);


--
-- TOC entry 4933 (class 2606 OID 17778)
-- Name: categories categories_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.categories
    ADD CONSTRAINT categories_pkey PRIMARY KEY (id);


--
-- TOC entry 4941 (class 2606 OID 17821)
-- Name: products products_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_pkey PRIMARY KEY (id);


--
-- TOC entry 4943 (class 2606 OID 17823)
-- Name: products products_sku_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_sku_key UNIQUE (sku);


--
-- TOC entry 4937 (class 2606 OID 17802)
-- Name: racks racks_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.racks
    ADD CONSTRAINT racks_pkey PRIMARY KEY (id);


--
-- TOC entry 4948 (class 2606 OID 17869)
-- Name: sale_details sale_details_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sale_details
    ADD CONSTRAINT sale_details_pkey PRIMARY KEY (id);


--
-- TOC entry 4946 (class 2606 OID 17850)
-- Name: sales sales_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sales
    ADD CONSTRAINT sales_pkey PRIMARY KEY (id);


--
-- TOC entry 4931 (class 2606 OID 17760)
-- Name: sessions sessions_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sessions
    ADD CONSTRAINT sessions_pkey PRIMARY KEY (id);


--
-- TOC entry 4924 (class 2606 OID 17745)
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- TOC entry 4926 (class 2606 OID 17741)
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- TOC entry 4928 (class 2606 OID 17743)
-- Name: users users_username_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_username_key UNIQUE (username);


--
-- TOC entry 4935 (class 2606 OID 17791)
-- Name: warehouses warehouses_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.warehouses
    ADD CONSTRAINT warehouses_pkey PRIMARY KEY (id);


--
-- TOC entry 4938 (class 1259 OID 17883)
-- Name: idx_products_category; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_products_category ON public.products USING btree (category_id);


--
-- TOC entry 4939 (class 1259 OID 17882)
-- Name: idx_products_sku; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_products_sku ON public.products USING btree (sku);


--
-- TOC entry 4944 (class 1259 OID 17884)
-- Name: idx_sales_date; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_sales_date ON public.sales USING btree (transaction_date);


--
-- TOC entry 4929 (class 1259 OID 17881)
-- Name: idx_sessions_token; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_sessions_token ON public.sessions USING btree (token);


--
-- TOC entry 4922 (class 1259 OID 17880)
-- Name: idx_users_email; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_users_email ON public.users USING btree (email);


--
-- TOC entry 4950 (class 2606 OID 17824)
-- Name: products products_category_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.categories(id) ON DELETE RESTRICT;


--
-- TOC entry 4951 (class 2606 OID 17834)
-- Name: products products_rack_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_rack_id_fkey FOREIGN KEY (rack_id) REFERENCES public.racks(id) ON DELETE RESTRICT;


--
-- TOC entry 4952 (class 2606 OID 17829)
-- Name: products products_warehouse_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_warehouse_id_fkey FOREIGN KEY (warehouse_id) REFERENCES public.warehouses(id) ON DELETE RESTRICT;


--
-- TOC entry 4954 (class 2606 OID 17875)
-- Name: sale_details sale_details_product_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sale_details
    ADD CONSTRAINT sale_details_product_id_fkey FOREIGN KEY (product_id) REFERENCES public.products(id) ON DELETE RESTRICT;


--
-- TOC entry 4955 (class 2606 OID 17870)
-- Name: sale_details sale_details_sale_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sale_details
    ADD CONSTRAINT sale_details_sale_id_fkey FOREIGN KEY (sale_id) REFERENCES public.sales(id) ON DELETE CASCADE;


--
-- TOC entry 4953 (class 2606 OID 17851)
-- Name: sales sales_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sales
    ADD CONSTRAINT sales_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE SET NULL;


--
-- TOC entry 4949 (class 2606 OID 17761)
-- Name: sessions sessions_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sessions
    ADD CONSTRAINT sessions_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


-- Completed on 2026-01-03 14:21:21

--
-- PostgreSQL database dump complete
--
