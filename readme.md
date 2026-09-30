# Alloy Determination API

REST API для сервиса определения состава сплавов по XRF-анализу.

## Технологии
- Go + Gin + GORM
- PostgreSQL
- MinIO (S3-совместимое хранилище)

## Таблицы БД

### `users`
| Поле | Тип | Описание |
|---|---|---|
| user_id | BIGINT PK | Идентификатор |
| user_login | VARCHAR(25) UNIQUE | Логин |
| user_password | VARCHAR(100) | Пароль |

### `alloys`
| Поле | Тип | Описание |
|---|---|---|
| alloy_id | BIGINT PK | Идентификатор |
| alloy_name | VARCHAR(100) | Название сплава |
| alloy_description | VARCHAR(200) | Описание |
| alloy_status | VARCHAR(15) | черновик / опубликован / удален |
| alloy_image_url | VARCHAR(255) | Имя файла картинки в MinIO |
| alloy_video_url | VARCHAR(255) | Имя файла видео в MinIO |
| alloy_energy_kev | NUMERIC(10,2) | Энергия пика, кэВ |
| alloy_intensity_cps | NUMERIC(12,2) | Интенсивность, cps |
| alloy_created_at | TIMESTAMPTZ | Дата создания |
| alloy_formed_at | TIMESTAMPTZ | Дата формирования |
| creator_id | BIGINT FK → users.user_id | Создатель |

### `likes`
| Поле | Тип | Описание |
|---|---|---|
| like_id | BIGINT PK | Идентификатор |
| user_id | BIGINT FK → users.user_id | Кто лайкнул |
| alloy_id | BIGINT FK → alloys.alloy_id | Что лайкнул |

## API методы

Базовый URL: `http://localhost:3000/api`

### Услуги (alloy)

| Метод | URL | Описание |
|---|---|---|
| GET | `/alloyCatalog?energy=X` | Список опубликованных с фильтром по энергии |
| GET | `/alloyFeed?alloy_id=N&next=true` | Лента опубликованных |
| GET | `/alloyDraft` | Черновик текущего пользователя |
| POST | `/alloyDraft` | Создать черновик + файлы (multipart: name, image, video) |
| PUT | `/alloyPublish` | Публикация (JSON) |
| DELETE | `/alloyDelete/:id` | Soft delete |
| POST | `/alloyLike` | Лайк 0/1 (JSON: alloy_id, like) |

### Пользователь (user)

| Метод | URL | Описание |
|---|---|---|
| POST | `/userRegister` | Регистрация |
| POST | `/userLogin` | Аутентификация (заглушка) |
| POST | `/userLogout` | Деавторизация (заглушка) |
