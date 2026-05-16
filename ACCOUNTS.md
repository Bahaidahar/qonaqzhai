# Demo Accounts

Тестовые аккаунты для qonaqzhai. Запусти backend + frontend, открой `http://localhost:3000`.

## Все пароли: `demo12345` (кроме admin)

| Email | Пароль | Роль | Статус | Что увидишь |
|-------|--------|------|--------|-------------|
| `customer@demo.kz` | `demo12345` | Customer | active | Чат с AI, каталог vendors, бронирования |
| `vendor1@demo.kz` | `demo12345` | Vendor | active | Профиль Rixos Almaty Ballroom (**approved**), inbox |
| `vendor2@demo.kz` | `demo12345` | Vendor | active | Профиль Studio Aitu Photo (**approved**), inbox |
| `vendor3@demo.kz` | `demo12345` | Vendor | active | Профиль Aizada Catering (**pending** — для admin demo) |
| `admin@qonaqzhai.kz` | `admin12345` | Admin | active | Модерация vendors, управление users, stats |

---

## Подробнее по аккаунтам

### Customer — `customer@demo.kz`
- Имя: **Aigerim Demo**
- Видит каталог с 2 одобренными vendors (Rixos, Studio Aitu)
- Может писать в чат, бронировать, отменять

### Vendor 1 — `vendor1@demo.kz`
- Название: **Rixos Almaty Ballroom**
- Категория: Venue
- Город: Almaty
- Цена от: 1 500 000 ₸
- Статус: **approved** (видим в каталоге)
- 1 фото

### Vendor 2 — `vendor2@demo.kz`
- Название: **Studio Aitu Photo**
- Категория: Photo & Video
- Город: Almaty
- Цена от: 450 000 ₸
- Статус: **approved**
- 1 фото

### Vendor 3 — `vendor3@demo.kz`
- Название: **Aizada Catering**
- Категория: Catering
- Город: Astana
- Цена от: 8 500 ₸
- Статус: **pending** — НЕ виден customer'ам пока admin не одобрит
- 1 фото

### Admin — `admin@qonaqzhai.kz`
- Pre-засеян автоматически при первом старте backend
- Логин: `admin@qonaqzhai.kz`
- Пароль: `admin12345`

---

## Сценарий "полный круг" для теста

1. Логин **admin** → `/admin` → Approve Aizada Catering
2. Logout → логин **customer@demo.kz** → `/vendors` → теперь видит 3 vendors → выбирает Aizada → жмёт **Request booking** (дата + гости)
3. Logout → логин **vendor3@demo.kz** → `/vendor/bookings` → видит запрос → **Accept**
4. Logout → логин **customer@demo.kz** → `/bookings` → статус **accepted**

## Сценарий "suspend user"

1. Логин **admin** → `/admin/users` → найти `customer@demo.kz` → **Suspend**
2. Logout → попытаться залогиниться `customer@demo.kz` → получишь **403 account suspended**
3. Логин admin → `/admin/users` → **Activate** → теперь снова работает

## Сценарий "chat AI"

Логин customer, в чате попробуй:
- `план тоя на 150 человек` → план + бюджет + vendors
- `бюджет на свадьбу` → блок с разбивкой
- `найди фотографа` → реальные approved vendors из БД
- `budget` / `vendor` / `photographer` — работает на en/ru/kz keywords

## Сценарий "vendor profile"

1. Логин **vendor1@demo.kz** → `/vendor` — профиль уже заполнен
2. Загрузи доп фото (jpg/png до 5MB)
3. Удали одну фотку (наведи курсор)
4. Поменяй описание → Save

---

## Reset базы (если что-то сломалось)

```bash
cd backend
# stop backend (Ctrl+C)
rm -f qonaqzhai.db qonaqzhai.db-wal qonaqzhai.db-shm
go run .
# в другом терминале:
python3 /tmp/seed_demo.py
```

Admin создастся автоматически при старте. Demo users — через seed скрипт.

---

## Что НЕ работает в демо

- Реальный Gemini AI (используется keyword-симулятор на backend)
- Email / SMS / push уведомления
- Оплата / эскроу
- Forgot password
- Mobile responsive (десктоп только)
- Производственный deploy (только локальный dev)
