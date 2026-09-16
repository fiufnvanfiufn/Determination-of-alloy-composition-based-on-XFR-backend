--
-- PostgreSQL database dump
--

\restrict bbSXgds7mqsU12kFtdwvjahD4KDp0fH3ZbtYCGGbwL6rZul8V0g6XzNfZrfQaS2

-- Dumped from database version 16.15 (Debian 16.15-1.pgdg13+2)
-- Dumped by pg_dump version 18.2 (Ubuntu 18.2-1.pgdg24.04+1)

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

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: alloys; Type: TABLE; Schema: public; Owner: myuser
--

CREATE TABLE public.alloys (
    id bigint NOT NULL,
    name character varying(100) NOT NULL,
    description character varying(200),
    status character varying(15) DEFAULT 'черновик'::character varying NOT NULL,
    img_url character varying(255),
    video_url character varying(255),
    energy_kev numeric(10,2),
    intensity_cps numeric(12,2),
    date_create timestamp with time zone NOT NULL,
    date_formed timestamp with time zone,
    creator_id bigint NOT NULL
);


ALTER TABLE public.alloys OWNER TO myuser;

--
-- Name: alloys_id_seq; Type: SEQUENCE; Schema: public; Owner: myuser
--

CREATE SEQUENCE public.alloys_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.alloys_id_seq OWNER TO myuser;

--
-- Name: alloys_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: myuser
--

ALTER SEQUENCE public.alloys_id_seq OWNED BY public.alloys.id;


--
-- Name: likes; Type: TABLE; Schema: public; Owner: myuser
--

CREATE TABLE public.likes (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    alloy_id bigint NOT NULL
);


ALTER TABLE public.likes OWNER TO myuser;

--
-- Name: likes_id_seq; Type: SEQUENCE; Schema: public; Owner: myuser
--

CREATE SEQUENCE public.likes_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.likes_id_seq OWNER TO myuser;

--
-- Name: likes_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: myuser
--

ALTER SEQUENCE public.likes_id_seq OWNED BY public.likes.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: myuser
--

CREATE TABLE public.users (
    id bigint NOT NULL,
    login character varying(25) NOT NULL,
    password character varying(100) NOT NULL
);


ALTER TABLE public.users OWNER TO myuser;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: myuser
--

CREATE SEQUENCE public.users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.users_id_seq OWNER TO myuser;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: myuser
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: alloys id; Type: DEFAULT; Schema: public; Owner: myuser
--

ALTER TABLE ONLY public.alloys ALTER COLUMN id SET DEFAULT nextval('public.alloys_id_seq'::regclass);


--
-- Name: likes id; Type: DEFAULT; Schema: public; Owner: myuser
--

ALTER TABLE ONLY public.likes ALTER COLUMN id SET DEFAULT nextval('public.likes_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: myuser
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Data for Name: alloys; Type: TABLE DATA; Schema: public; Owner: myuser
--

COPY public.alloys (id, name, description, status, img_url, video_url, energy_kev, intensity_cps, date_create, date_formed, creator_id) FROM stdin;
5	Оловянистая бронза	tmp	опубликован	http://localhost:9000/media/бронза_картинка.jpg	http://localhost:9000/media/бронза.mp4	8.04	7.21	2026-09-16 00:01:37.80018+00	2026-09-16 00:01:37.80018+00	1
7	Латунь	tmp	черновик	http://localhost:9000/media/514770_big.jpg	http://localhost:9000/media/Rome_coin.mp4	8.63	6.51	2026-09-16 00:01:37.80018+00	\N	2
8	Пьютер	Сплав на основе олова с добавлением меди и сурьмы. Свинец отсутствует, что типично для качественной кухонной утвари.	опубликован	http://localhost:9000/media/пьютер.webp	http://localhost:9000/media/пьютер_видео.mp4	25.27	20.94	2026-09-16 00:01:37.80018+00	2026-09-16 00:01:37.80018+00	3
9	Электрум	tmp	опубликован	http://localhost:9000/media/514770_big.jpg	http://localhost:9000/media/Rome_coin.mp4	9.71	2.62	2026-09-16 00:01:37.80018+00	2026-09-16 00:01:37.80018+00	3
10	Серебро (удалённая)	Пробная карточка, логически удалена пользователем	удален	http://localhost:9000/media/silver.jpg		22.16	12.34	2026-09-16 00:01:58.054563+00	2026-09-16 00:01:58.054563+00	1
11	1		удален			0.00	0.00	2026-09-16 00:27:16.037319+00	2026-09-16 00:45:20.508299+00	1
4	Мышьяковистая бронза	Ранний бронзовый век. Сплав Cu-As, предшественник оловянной бронзы. Высокий пик мышьяка.	удален	http://localhost:9000/media/мышьяковая_бронза.png	http://localhost:9000/media/мышьяковая_бронза.mp4	10.53	83.07	2026-09-16 00:01:37.80018+00	2026-09-16 00:01:37.80018+00	1
6	Свинец	Вислая печать. Четкий пик свинца без значительных примесей серебра или олова.	удален	http://localhost:9000/media/свинец.jpg	http://localhost:9000/media/свинец.mp4	10.55	5.72	2026-09-16 00:01:37.80018+00	2026-09-16 00:01:37.80018+00	2
13	2		черновик			10.50	5.40	2026-09-16 01:44:53.27121+00	\N	2
12	1	tmp	опубликован			10.63	5.45	2026-09-16 01:42:58.920531+00	2026-09-16 01:44:11.467614+00	3
\.


--
-- Data for Name: likes; Type: TABLE DATA; Schema: public; Owner: myuser
--

COPY public.likes (id, user_id, alloy_id) FROM stdin;
1	2	4
2	3	4
3	1	5
4	3	5
5	1	6
6	2	8
7	3	8
8	1	9
9	2	12
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: myuser
--

COPY public.users (id, login, password) FROM stdin;
1	ivanov	pass1
2	petrova	pass2
3	sidorov	pass3
\.


--
-- Name: alloys_id_seq; Type: SEQUENCE SET; Schema: public; Owner: myuser
--

SELECT pg_catalog.setval('public.alloys_id_seq', 13, true);


--
-- Name: likes_id_seq; Type: SEQUENCE SET; Schema: public; Owner: myuser
--

SELECT pg_catalog.setval('public.likes_id_seq', 9, true);


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: myuser
--

SELECT pg_catalog.setval('public.users_id_seq', 3, true);


--
-- Name: alloys alloys_pkey; Type: CONSTRAINT; Schema: public; Owner: myuser
--

ALTER TABLE ONLY public.alloys
    ADD CONSTRAINT alloys_pkey PRIMARY KEY (id);


--
-- Name: likes likes_pkey; Type: CONSTRAINT; Schema: public; Owner: myuser
--

ALTER TABLE ONLY public.likes
    ADD CONSTRAINT likes_pkey PRIMARY KEY (id);


--
-- Name: users uni_users_login; Type: CONSTRAINT; Schema: public; Owner: myuser
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT uni_users_login UNIQUE (login);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: myuser
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: alloys fk_alloys_creator; Type: FK CONSTRAINT; Schema: public; Owner: myuser
--

ALTER TABLE ONLY public.alloys
    ADD CONSTRAINT fk_alloys_creator FOREIGN KEY (creator_id) REFERENCES public.users(id);


--
-- Name: likes fk_likes_alloy; Type: FK CONSTRAINT; Schema: public; Owner: myuser
--

ALTER TABLE ONLY public.likes
    ADD CONSTRAINT fk_likes_alloy FOREIGN KEY (alloy_id) REFERENCES public.alloys(id);


--
-- Name: likes fk_likes_user; Type: FK CONSTRAINT; Schema: public; Owner: myuser
--

ALTER TABLE ONLY public.likes
    ADD CONSTRAINT fk_likes_user FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- PostgreSQL database dump complete
--

\unrestrict bbSXgds7mqsU12kFtdwvjahD4KDp0fH3ZbtYCGGbwL6rZul8V0g6XzNfZrfQaS2

