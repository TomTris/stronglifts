# StrongLifts 5×5 Tracker

Self-hosted, single-binary workout tracker cho chương trình StrongLifts 5×5. Multi-user với Google login, leaderboard, PWA. Chạy trên máy tính hoặc VPS, dùng từ trình duyệt điện thoại.

---

## Mục lục

- [Chương trình 5×5 là gì](#chương-trình-5×5-là-gì)
- [Tính năng](#tính-năng)
- [Yêu cầu hệ thống](#yêu-cầu-hệ-thống)
- [Setup Google OAuth](#setup-google-oauth)
- [Cài đặt và chạy](#cài-đặt-và-chạy)
- [Hướng dẫn sử dụng](#hướng-dẫn-sử-dụng)
- [Cài trên điện thoại (PWA)](#cài-trên-điện-thoại-pwa)
- [Deploy lên server](#deploy-lên-server)
- [Cấu hình](#cấu-hình)
- [Thuật toán và công thức](#thuật-toán-và-công-thức)
- [API Reference](#api-reference)
- [Cấu trúc database](#cấu-trúc-database)
- [Tech stack](#tech-stack)
- [Cấu trúc project](#cấu-trúc-project)
- [Development](#development)

---

## Chương trình 5×5 là gì

StrongLifts 5×5 là chương trình tập sức mạnh dành cho người mới hoặc quay lại tập. Xây trên nguyên tắc **progressive overload** — tăng tạ đều đặn mỗi buổi.

### Cấu trúc

Mỗi tuần 3 buổi, xen kẽ ngày nghỉ. Có 2 workout luân phiên:

| Workout A | Workout B |
|-----------|-----------|
| Squat 5×5 | Squat 5×5 |
| Bench Press 5×5 | Overhead Press 5×5 |
| Barbell Row 5×5 | Deadlift 1×5 |

Lịch xoay vòng: Tuần 1 (A-B-A), Tuần 2 (B-A-B), ...

Squat mỗi buổi. Deadlift chỉ 1 set vì stress lớn lên hệ thần kinh.

### Quy tắc progression

- Đủ 5 rep × 5 set → buổi sau **+2.5 kg** (Deadlift +5 kg).
- Không đủ rep ở bất kỳ set nào → giữ nguyên weight. Vẫn làm đủ 5 set.
- Fail cùng weight **3 buổi liên tiếp** → **deload 10%**.

### Weight khởi đầu mặc định

| Bài tập | Khởi đầu |
|---------|----------|
| Squat | 20 kg (bar trống) |
| Bench Press | 20 kg |
| Overhead Press | 20 kg |
| Barbell Row | 30 kg |
| Deadlift | 40 kg |

Nếu đã từng tập, app tính starting weight từ khả năng hiện tại (xem Setup Wizard).

---

## Tính năng

### Google Login & Multi-user
Đăng nhập bằng Google account. Mỗi user có dữ liệu riêng. Không cần tạo tài khoản thủ công.

### Nickname & Leaderboard
Sau lần đăng nhập đầu tiên, chọn nickname (1-20 ký tự). Nickname hiển thị trên bảng xếp hạng. Ranking theo tổng Squat + Bench Press + Deadlift (working weight hiện tại).

### Setup Wizard
Lần đầu: hỏi đã từng tập chưa. Nếu rồi: form pre-fill giá trị mẫu, chỉnh theo khả năng thực → app tính 1RM (Epley) → đặt starting weight ở 50% 1RM. Nếu chưa: dùng weight mặc định.

### Workout Logging
Tap set → chọn số rep → Log Set. Ô xanh (đủ rep) hoặc vàng (thiếu). Rest timer tự động sau mỗi set.

### Warm-up Sets
Mỗi bài có warm-up mở sẵn: bar trống × 5, 40% × 5, 60% × 3, 80% × 2. Chỉ hiện set nào có weight > bar. Kèm plate calculator riêng.

### Rest Timer
Tự động sau mỗi set. 90s (đủ rep), 180s (3-4 rep), 300s (0-2 rep). Âm thanh khi hết giờ. Skip bất kỳ lúc nào.

### Plate Calculator
Visual đĩa mỗi bên bar, mã màu theo trọng lượng (20kg đỏ, 10kg xanh, 5kg vàng, 2.5kg xanh lá, 1.25kg xám).

### Auto-Progression & Deload
Finish Workout → tự kiểm tra pass/fail → tăng weight hoặc tăng fail counter → 3 fail liên tiếp → deload 10%.

### Workout Notes
Ghi chú tự do mỗi buổi, auto-save. Hiển thị trong History.

### Progress Charts
Đồ thị weight theo thời gian cho từng bài tập (Chart.js).

### Body Weight Tracking
Log cân nặng hàng ngày + chart.

### Export / Backup
Settings → Export → JSON chứa toàn bộ exercises, workouts, body weights.

### PWA
Installable trên Android/iOS. Cache offline cho giao diện.

---

## Yêu cầu hệ thống

### Build từ source
- Go 1.22+
- Node.js 18+ và npm (chỉ cần nếu sửa frontend)

### Chạy
- Bất kỳ máy nào chạy Go binary
- Trình duyệt hiện đại trên thiết bị client
- Kết nối internet (để load Google Sign-In SDK)

---

## Setup Google OAuth

Cần 1 lần, miễn phí, ~5 phút.

1. Vào [console.cloud.google.com](https://console.cloud.google.com)
2. Tạo project mới (hoặc dùng cái có sẵn)
3. **APIs & Services → Credentials → Create Credentials → OAuth client ID**
4. Application type: **Web application**
5. Authorized JavaScript origins, thêm:
   - `http://localhost:8080` (dev local)
   - `https://gym.yourdomain.com` (nếu deploy production)
6. Create → copy **Client ID** (dạng `xxxx.apps.googleusercontent.com`)

Nếu OAuth consent screen chưa cấu hình:
- APIs & Services → OAuth consent screen
- User Type: **External**
- Chỉ cần điền App name, User support email, Developer email
- Scopes: không cần thêm gì (app chỉ dùng basic profile)
- Publish app (hoặc thêm test users nếu muốn giới hạn)

---

## Cài đặt và chạy

### Bước 1: Cấu hình

```bash
unzip stronglifts.zip
cd stronglifts
cp .env.example .env
```

Mở `.env`, thay `GOOGLE_CLIENT_ID`:

```env
GOOGLE_CLIENT_ID=xxxx.apps.googleusercontent.com
PORT=8080
DB_PATH=stronglifts.db
```

### Bước 2: Build

```bash
go mod tidy
go build -o stronglifts .
```

### Bước 3: Chạy

```bash
./stronglifts
```

Mở `http://localhost:8080`. Server tự đọc `.env`.

Nếu không muốn dùng `.env`, truyền trực tiếp:

```bash
GOOGLE_CLIENT_ID="xxxx.apps.googleusercontent.com" ./stronglifts
```

Environment variables luôn override `.env`.

### Build toàn bộ từ source (nếu sửa frontend)

```bash
cd frontend && npm install && npm run build-only && cd ..
go mod tidy && go build -o stronglifts .
```

---

## Hướng dẫn sử dụng

### Lần đầu

1. Mở app → bấm **Sign in with Google**.
2. Chọn nickname (hiện trên leaderboard).
3. Setup Wizard: chọn "Đã tập" → chỉnh weight/reps → app tính starting weight. Hoặc "Chưa" → bar trống.
4. Về Home, bấm **Start Workout**.

### Buổi tập

1. Mở warm-up, tập warm-up sets theo hướng dẫn.
2. Tap ô set → chọn rep → **Log Set**.
3. Rest timer tự chạy. Chờ hoặc Skip.
4. Lặp lại cho tất cả exercises.
5. Ghi notes nếu muốn.
6. **Finish Workout** → weight tự cập nhật.

### Leaderboard

Tab **Rank** → bảng xếp hạng tất cả users. Sorted theo tổng Squat + Bench + Deadlift working weight. Hiển thị avatar, nickname, chi tiết từng bài, số buổi tập.

### Export

Settings → Export → download JSON backup.

---

## Cài trên điện thoại (PWA)

### Android (Chrome)

1. Mở `http://<server>:8080` trong Chrome.
2. Menu (⋮) → **"Add to Home screen"** hoặc **"Install app"**.
3. Icon xuất hiện trên launcher.

### iOS (Safari)

1. Mở URL trong Safari.
2. Share (□↑) → **"Add to Home Screen"**.

API calls cần kết nối tới server. Giao diện cache offline.

---

## Deploy lên server

### Cùng WiFi (đơn giản)

```bash
# Tìm IP
ip addr | grep "inet " | grep -v 127.0.0.1   # Linux/macOS
ipconfig                                        # Windows
```

Truy cập: `http://192.168.x.x:8080`

### VPS (24/7)

Upload binary + `.env`:

```bash
scp stronglifts .env user@server:/opt/stronglifts/
ssh user@server
cd /opt/stronglifts && ./stronglifts
```

### Systemd service

```ini
# /etc/systemd/system/stronglifts.service
[Unit]
Description=StrongLifts 5x5
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/stronglifts
ExecStart=/opt/stronglifts/stronglifts
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

```bash
sudo systemctl daemon-reload
sudo systemctl enable --now stronglifts
```

Server tự đọc `.env` từ `WorkingDirectory`.

### Reverse proxy (Nginx + HTTPS)

```nginx
server {
    listen 80;
    server_name gym.yourdomain.com;

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

Thêm HTTPS với Certbot: `sudo certbot --nginx -d gym.yourdomain.com`

**Quan trọng:** thêm domain HTTPS vào Google OAuth "Authorized JavaScript origins".

---

## Cấu hình

Qua `.env` file hoặc environment variables. Env vars override `.env`.

| Variable | Mặc định | Bắt buộc | Mô tả |
|----------|----------|----------|-------|
| `GOOGLE_CLIENT_ID` | — | Có | Google OAuth Client ID |
| `PORT` | `8080` | Không | HTTP port |
| `DB_PATH` | `stronglifts.db` | Không | Đường dẫn SQLite |

File `.env` hỗ trợ:
- Comments (`# ...`)
- Quotes (`"value"` hoặc `'value'`)
- Không override env vars đã set

---

## Thuật toán và công thức

### 1RM (Epley Formula)

```
1RM = Weight × (1 + Reps / 30)
```

### Starting Weight

```
Starting Weight = floor(1RM × 0.5 / increment) × increment
```

Minimum 20 kg. Increment: 2.5 kg (Deadlift 5.0 kg).

### Deload

```
New Weight = floor(Current Weight × 0.9 / 2.5) × 2.5
```

Minimum 20 kg. Fail counter reset.

### Warm-up Sets

| Weight | Reps | Điều kiện |
|--------|------|-----------|
| 20 kg (bar) | 5 | Luôn có |
| round(WW × 0.4) | 5 | Nếu > 20 kg và < WW |
| round(WW × 0.6) | 3 | Nếu > 20 kg và < WW |
| round(WW × 0.8) | 2 | Nếu > 20 kg và < WW |

### Rest Timer

| Kết quả set | Nghỉ |
|-------------|------|
| 5/5 rep | 90s |
| 3-4 rep | 180s |
| 0-2 rep | 300s |

### Plate Calculator

`Per Side = (Total - 20) / 2`. Greedy: 20, 10, 5, 2.5, 1.25 kg.

### Leaderboard

```
Total = Squat (current working weight) + Bench Press + Deadlift
```

Sorted descending. Chỉ hiện users có nickname.

---

## API Reference

Base URL: `http://localhost:8080/api`

### Public (no auth)

| Method | Endpoint | Mô tả |
|--------|----------|-------|
| GET | `/config` | `{"google_client_id": "..."}` |
| POST | `/auth/google` | Login. Body: `{"credential": "google_id_token"}` |
| GET | `/warmup/{weight}` | Warm-up sets |
| GET | `/plates/{weight}` | Plate calculator |
| GET | `/leaderboard` | Bảng xếp hạng |

### Auth required (session cookie)

| Method | Endpoint | Mô tả |
|--------|----------|-------|
| GET | `/auth/me` | User hiện tại |
| POST | `/auth/nickname` | Set nickname. Body: `{"nickname": "xxx"}` |
| POST | `/auth/logout` | Logout |
| GET | `/setup/status` | `{"setup_done": true/false}` |
| POST | `/setup/starting-weights` | Set starting weights |
| GET | `/exercises` | User's exercises + weights |
| PUT | `/exercises/{id}/weight` | Chỉnh weight thủ công |
| GET | `/workout/active` | Workout đang chạy |
| GET | `/workout/next-type` | "A" hoặc "B" |
| POST | `/workout/start` | Tạo workout mới |
| GET | `/workout/{id}` | Chi tiết workout |
| POST | `/workout/{id}/complete` | Hoàn thành, trigger progression |
| PUT | `/workout/{id}/notes` | Cập nhật notes |
| DELETE | `/workout/{id}` | Xoá workout |
| POST | `/set/{id}/complete` | Log set. Body: `{"reps": 5}` |
| GET | `/progress/{exerciseId}` | Lịch sử weight |
| GET | `/history?limit=50` | Danh sách workouts |
| GET | `/bodyweight?limit=100` | Lịch sử cân nặng |
| POST | `/bodyweight` | Log cân nặng. Body: `{"weight": 75.5}` |
| DELETE | `/bodyweight/{id}` | Xoá entry |
| GET | `/export` | Download JSON backup |

Auth qua cookie `session` (set tự động khi login, 30 ngày).

---

## Cấu trúc database

SQLite, single file. 7 tables:

### users

| Column | Type | Mô tả |
|--------|------|-------|
| id | INTEGER PK | |
| google_id | TEXT UNIQUE | Google subject ID |
| email | TEXT | |
| nickname | TEXT | Hiển thị trên leaderboard |
| avatar_url | TEXT | Google profile picture |
| created_at | TEXT | ISO 8601 |

### sessions

| Column | Type | Mô tả |
|--------|------|-------|
| token | TEXT PK | Random 64-char hex |
| user_id | INTEGER FK | → users.id |
| created_at | TEXT | |
| expires_at | TEXT | +30 ngày |

### exercises (template)

| Column | Type | Mô tả |
|--------|------|-------|
| id | INTEGER PK | |
| name | TEXT UNIQUE | Tên bài tập |
| default_weight | REAL | Weight mặc định cho user mới |
| increment | REAL | 2.5 hoặc 5.0 |

### user_exercises

| Column | Type | Mô tả |
|--------|------|-------|
| id | INTEGER PK | |
| user_id | INTEGER FK | → users.id |
| exercise_id | INTEGER FK | → exercises.id |
| current_weight | REAL | Weight hiện tại của user |
| fail_count | INTEGER | 0-2 |
| UNIQUE | (user_id, exercise_id) | |

### workouts

| Column | Type | Mô tả |
|--------|------|-------|
| id | INTEGER PK | |
| user_id | INTEGER FK | → users.id |
| type | TEXT | "A" hoặc "B" |
| date | TEXT | ISO 8601 |
| completed | INTEGER | 0/1 |
| notes | TEXT | |

### workout_exercises, workout_sets

Giống trước, linked qua workout_id.

### body_weights

| Column | Type | Mô tả |
|--------|------|-------|
| id | INTEGER PK | |
| user_id | INTEGER FK | → users.id |
| date | TEXT | YYYY-MM-DD |
| weight | REAL | kg |

### app_settings

Key-value. Keys: `setup_done_{userID}`.

---

## Tech stack

| Layer | Công nghệ |
|-------|-----------|
| Backend | Go 1.22+, net/http, embed |
| Auth | Google Identity Services (ID token), server-side session |
| Database | SQLite (modernc.org/sqlite, pure Go) |
| Frontend | Vue 3 + TypeScript + Pinia + Vue Router |
| Charts | Chart.js 4 |
| Build | Vite 5 |
| PWA | Service Worker + Web App Manifest |

Single binary ~15MB. Zero runtime dependencies.

---

## Cấu trúc project

```
stronglifts/
├── main.go                  # Entry point, .env loader, embed frontend
├── go.mod
├── .env.example             # Template config
├── .env                     # Your config (gitignored)
├── Makefile
├── README.md
├── db/db.go                 # Schema, migrations, all queries
├── models/models.go         # Shared structs
├── handlers/handlers.go     # HTTP handlers + auth middleware
└── frontend/
    ├── index.html
    ├── package.json
    ├── vite.config.ts
    ├── public/              # PWA assets
    ├── dist/                # Pre-built (embedded into binary)
    └── src/
        ├── main.ts          # Router, auth guard, PWA
        ├── App.vue          # Shell, bottom nav
        ├── api/index.ts     # HTTP client
        ├── stores/
        │   ├── auth.ts      # Login state, Google auth
        │   └── workout.ts   # Exercises, active workout
        ├── types/index.ts
        ├── views/
        │   ├── LoginView.vue
        │   ├── NicknameView.vue
        │   ├── HomeView.vue
        │   ├── SetupView.vue
        │   ├── WorkoutView.vue
        │   ├── ProgressView.vue
        │   ├── BodyWeightView.vue
        │   ├── HistoryView.vue
        │   ├── LeaderboardView.vue
        │   └── SettingsView.vue
        └── components/
            ├── SetLogger.vue
            ├── WarmupDisplay.vue
            ├── RestTimer.vue
            └── PlateCalculator.vue
```

---

## Development

### Dev mode (2 terminals)

```bash
# Terminal 1: Go backend
go run .

# Terminal 2: Vue hot reload
cd frontend && npm run dev
```

Vite proxy `/api/*` → `:8080`. Mở `http://localhost:5173`.

### Cross-compile

```bash
GOOS=linux GOARCH=arm64 go build -o stronglifts-arm .   # ARM VPS
GOOS=windows GOARCH=amd64 go build -o stronglifts.exe .  # Windows
```

### Rebuild frontend

```bash
cd frontend && npm install && npm run build-only
cd .. && go build -o stronglifts .
```

### Backup

Copy `stronglifts.db` hoặc GET `/api/export`.