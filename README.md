# Personal Blog — Project Context

> File này ghi lại toàn bộ bối cảnh, quyết định, và tiến độ của dự án để bất kỳ đoạn chat AI mới nào cũng có thể đọc và nắm được toàn bộ tình trạng dự án mà không cần người dùng giải thích lại. Cập nhật file này mỗi khi có thay đổi lớn về tiến độ hoặc quyết định.

## Tổng quan

Blog cá nhân viết bằng Go thuần (server-rendered `html/template`, không dùng framework/JS framework), SQLite làm database, Tailwind CSS cho styling. Mục tiêu cuối: 1 người (admin) viết bài, khách chỉ đọc (không có tính năng public form/comment).

**Tech stack:**
- Go 1.25.3, module `github.com/letrongvu/blog`
- DB: SQLite qua `modernc.org/sqlite` (pure Go driver, không cần CGO)
- Migration: `golang-migrate/migrate/v4`, embed bằng `go:embed`
- Template: `html/template` chuẩn, tự viết renderer riêng (không dùng lib ngoài)
- Static/template assets: nhúng vào binary qua `go:embed` (production), hoặc đọc từ đĩa qua `os.DirFS` khi `APP_ENV=dev`
- CSS: Tailwind CSS v4, build qua `npx tailwindcss` (Node chỉ cần lúc dev/build CSS, không cần lúc chạy server)
- Auth: session cookie tự viết tay (không dùng lib), password hash bằng `bcrypt`

**Quyết định đã chốt với người dùng** (đừng hỏi lại, đã xác nhận):
- Deploy target: **VPS tự quản lý** (không dùng PaaS) — cần tự setup TLS/reverse proxy/systemd/backup.
- Đối tượng dùng: **chỉ 1 admin viết bài**, khách public chỉ đọc, không có form public nào — vì vậy CSRF được coi là ưu tiên thấp (không phải không cần, chỉ là không khẩn cấp).

## Cấu trúc thư mục chính

```
cmd/api/main.go              → entrypoint, wiring toàn bộ DI + routes
cmd/seed/main.go             → tool seed 1 user admin (go run ./cmd/seed <username> <password>)
internal/model/              → Post, User (struct thuần, không có logic)
internal/repository/         → sqlite CRUD cho Post/User/Session, sentinel errors (ErrPostNotFound...)
internal/services/           → PostService (validate + slugify), UserService (bcrypt + session token)
internal/handler/            → HTTP handlers, 1 struct Handler dùng chung, method theo domain
internal/middleware/         → RequireAuth (session cookie guard), Recover (panic recovery)
internal/view/               → template renderer tay viết (render.go), PageData wrapper (data.go)
migration/                   → SQL migration (users, posts, sessions), embed.go để nhúng
web/embed.go                 → chọn nguồn template/static: go:embed (prod) hoặc os.DirFS (dev)
web/templates/                → layouts/, partials/, pages/ (kể cả pages/admin/)
web/static/css/               → input.css (nguồn Tailwind), output.css (build ra, đã .gitignore? KHÔNG — output.css được commit vì cần cho embed prod)
package.json                  → devDependency tailwindcss + @tailwindcss/cli (chỉ dùng lúc build CSS)
```

## Đã hoàn thành

### 1. Data layer
- Migration: `users` (id, username unique, password_hash, created_at), `posts` (id, title, slug unique, content, published, created_at, updated_at, có index published+created_at), `sessions` (token PK, user_id FK, expires_at, created_at).
- Repository: CRUD đầy đủ cho Post (`ListPublished`, `ListAll`, `GetBySlug`, `GetByID`, `Create`, `Update`, `Delete`), User (`GetByUsername`, `Create`), Session (`Create`, `GetUserIDByToken`, `Delete`).
- Migration tự chạy lúc server start (`repository.Migrate` trong `main.go`).

### 2. Service layer
- `PostService`: validate title/content không rỗng, tự sinh slug (lowercase, non-alphanumeric → `-`), CRUD passthrough.
- `UserService`: `Register` (bcrypt hash), `Login` (verify password, sinh token random 32 byte base64, lưu session hết hạn 24h), `ValidateSession`, `Logout`.

### 3. HTTP layer — routes đã đăng ký (`cmd/api/main.go`)
Public:
- `GET /` → Home (list bài published)
- `GET /about` → About (tĩnh)
- `GET /posts/{slug}` → Post detail, 404 nếu không tìm thấy slug
- `GET /login`, `POST /login` → form + xử lý đăng nhập, set cookie `session` (HttpOnly, 24h)
- `POST /logout` → xoá session DB + xoá cookie
- `GET /static/` → serve static assets
- `/` (catch-all) → NotFound (404.html)

Admin (bọc `middleware.RequireAuth(userService)`):
- `GET /admin` → Dashboard (list toàn bộ post kể cả draft)
- `GET /admin/posts/new`, `POST /admin/posts` → tạo bài
- `GET /admin/posts/{id}/edit`, `POST /admin/posts/{id}` → sửa bài
- `POST /admin/posts/{id}/delete` → xoá bài

### 4. Middleware
- `RequireAuth`: đọc cookie `session` → `ValidateSession` → redirect `/login` nếu thiếu/invalid; lưu `userID` vào context.
- `Recover`: bọc toàn bộ mux, catch panic → 500.

### 5. Template rendering (`internal/view`)
- `PageData{ CurrentUser *model.User; Data any; Error string }` — **quan trọng**: `Data` dùng cho payload chính (post, []post...), `Error` dùng riêng cho message lỗi validate — KHÔNG dùng `Data` để chứa string lỗi (đã từng là bug, xem mục Bug đã sửa).
- `view.Init(fsys)` parse `layouts/*.html` + `partials/*.html` làm base, rồi clone riêng cho mỗi file trong `pages/*.html` + `pages/admin/*.html`.
- `view.Render(w, "page.html", data)` render qua template `"layout"`.
- **Dev hot-reload**: nếu `APP_ENV=dev`, `Render` gọi lại `Init` mỗi request (re-parse từ đĩa) — để sửa HTML/Tailwind class thấy ngay không cần restart server. Mặc định (không set env) dùng bản đã `go:embed`, không re-parse (đúng hành vi production).

### 6. Tailwind CSS
- `web/static/css/input.css`: `@import "tailwindcss";`
- Build: `npx tailwindcss -i web/static/css/input.css -o web/static/css/output.css [--watch|--minify]`
- `base.html` có `<head><link rel="stylesheet" href="/static/css/output.css"></head>`
- `home.html` hiện có test style trên `<h1>` (`text-4xl font-bold text-blue-600`) để xác minh pipeline hoạt động — **đây là style test, chưa phải thiết kế cuối cùng**.

### 7. Git hygiene
- `.gitignore`: `node_modules/`, `blog.db`.
- `blog.db` đã được `git rm --cached` — **đã xác nhận: chưa từng có commit nào chứa dữ liệu thật (user/session) của `blog.db`**, chỉ 1 commit cũ (`b27d7d6`) chứa bản schema-only. Không có secret bị lộ lên GitHub.

## Bug đã phát hiện & sửa (lịch sử, giải thích tại sao code hiện tại như vậy)

1. **`RequireAuth` redirect sai path** — từng redirect `/admin/login` (không tồn tại) thay vì `/login` → luôn ra 404 khi chưa đăng nhập vào `/admin`. Đã sửa ở `internal/middleware/middleware.go` (commit `7342090`).
2. **`home.html` dump raw `PageData`** — từng in `{{.}}` (toàn bộ struct) thay vì loop `.Data`. Đã sửa (commit `a5a67e1`).
3. **`PageData.Data` bị dùng lẫn 2 kiểu** — khi validate lỗi, `CreatePost`/`UpdatePost` gán `Data: err.Error()` (string), nhưng `post_form.html` giả định `.Data` luôn là `*model.Post` và gọi `.Data.ID` → template execution error, và lỗi này bị nuốt vì code không check giá trị trả về của `view.Render` → **kết quả: submit form rỗng ra trang trắng, không có thông báo lỗi**. Đã sửa bằng cách thêm field `Error` riêng trong `PageData` (commit `0df4c40`). `CreatePost`/`UpdatePost` giờ dùng `Error:`, không dùng `Data:` cho lỗi.

## Vấn đề nhỏ CÒN TỒN ĐỌNG — chưa fix, cần biết

1. **Session hết hạn không tự xoá khỏi DB** (không có cleanup job) — không critical, nhưng bảng `sessions` sẽ phình dần theo thời gian.

### Đã fix (trước đây nằm ở mục này)
- ~~checkbox `published` typo `cheked` → `checked` ở `post_form.html:15`~~ — đã fix.
- ~~`auth.go` (Login) dùng `Data:` thay vì `Error:` cho message lỗi~~ — đã fix, giờ dùng `Error:`.
- ~~`package.json:2` typo `personal=blog`~~ — đã fix, đúng `personal-blog`.
- ~~Chưa test full CRUD end-to-end qua browser thật~~ — **đã test xong, CRUD hoạt động đúng** (tạo bài → xuất hiện ở Home → sửa bỏ publish → biến mất khỏi Home → xoá).

## Việc CÒN LẠI để hoàn thành 100% dự án

### Nhóm 2 — Hạ tầng VPS (CHƯA BẮT ĐẦU, chỉ mới lên checklist, chưa code/config gì)
- [ ] Reverse proxy (nginx/Caddy) + TLS (Let's Encrypt) trước Go server — server hiện chạy HTTP thuần cổng `:8080`, không tự có HTTPS.
- [ ] Set `Secure: true` cho cookie session trong `auth.go` (`Login`) — **chỉ làm sau khi có HTTPS**, nếu set sớm mà server chưa có TLS thì cookie sẽ không hoạt động.
- [ ] Build production bằng `go build -o blog ./cmd/api` (không dùng `go run` trên VPS).
- [ ] systemd service để binary tự start lúc boot + tự restart khi crash.
- [ ] Backup `blog.db` định kỳ (cron), vì VPS tự quản lý không có backup tự động như PaaS.
- [ ] Đảm bảo **KHÔNG set `APP_ENV=dev`** trên VPS (nếu set, server sẽ tìm `web/templates`/`web/static` trên đĩa — không tồn tại nếu chỉ deploy binary → crash).
- [ ] Domain/DNS — chưa quyết, cần hỏi người dùng khi tới bước này.
- [ ] Firewall: chỉ mở 80/443 ra ngoài, không mở 8080 trực tiếp.

### Nhóm 3 — Bảo mật còn lại
- [x] `.gitignore` cho `blog.db`, `node_modules/` — xong.
- [ ] CSRF protection cho form admin (POST login/create/update/delete) — ưu tiên thấp (chỉ 1 admin dùng), nhưng nên làm trước khi thật sự public.

### Nhóm 4 — Polish nội dung
- [ ] Nội dung thật cho `about.html`, `404.html` (hiện chỉ placeholder 1 dòng).
- [ ] Viết test (hiện có 0 file `_test.go` trong toàn bộ project) — ít nhất nên có test cho `PostService`/`UserService` (validate, slugify, login logic).
- [ ] Thiết kế UI thật cho `home.html`, `post.html`, `dashboard.html`, `post_form.html`, `login.html` bằng Tailwind (hiện chỉ có 1 style test ở `home.html`, còn lại HTML thô không class).

### Nhóm 5 — Đa ngôn ngữ (i18n) — mới phát sinh, chưa bắt đầu
- [x] Quyết định: ưu tiên tiếng Anh trước (English-first), tiếng Việt làm sau.
- [ ] `<html lang="...">` ở `base.html` cần đổi thành `en` cho khớp ưu tiên hiện tại (tạm thời, chưa có cơ chế chọn locale động).
- [ ] Footer đã có 2 link `VI`/`EN` placeholder (`web/templates/partials/footer.html`, trỏ `href="/"`, chưa có logic thật) — cần thiết kế cơ chế xác định locale (path `/en/...`, query param, hoặc cookie) trước khi làm route thật.
- [ ] Nội dung song ngữ cho trang tĩnh (`about.html`, `404.html`...) và bài viết — chưa quyết định cách lưu trữ bản dịch (route/trang riêng theo locale? field riêng trong DB `posts`?).

### Nhóm nhỏ — fix bug tồn đọng (nên làm sớm, rẻ)
- [x] Sửa `cheked` → `checked`.
- [x] Đổi `Login` sang dùng `Error:` thay vì `Data:` cho đồng bộ.
- [x] Sửa `package.json` name typo.
- [x] Chạy smoke test CRUD đầy đủ.

## Cách chạy dev

**Chạy thường (không cần thấy live-reload):**
```
go run ./cmd/api
```

**Chạy có Tailwind hot-reload** (2 terminal song song):
```powershell
# Terminal 1
npx tailwindcss -i web/static/css/input.css -o web/static/css/output.css --watch

# Terminal 2
$env:APP_ENV = "dev"
go run ./cmd/api
```

**Seed user admin đầu tiên** (bắt buộc trước khi test login lần đầu trên máy mới, vì DB không seed sẵn ai):
```
go run ./cmd/seed <username> <password>
```

## Ghi chú vận hành quan trọng
- `DB_PATH` env var override đường dẫn DB (default `./blog.db`).
- `APP_ENV=dev` **chỉ dùng lúc dev cá nhân**, tuyệt đối không set trên production/VPS.
- Mọi thay đổi ở `web/templates/` hoặc `web/static/` khi **không** có `APP_ENV=dev` đều cần `go run`/`go build` lại mới thấy hiệu lực, vì bị `go:embed` đóng băng vào binary.
