# Auto Company

**Auto Company** — управляемый набор правил, skills, специализированных AI-ролей и CLI для создания и развития production-oriented SaaS и мобильных продуктов.

Он не обещает «идеальное приложение из одной фразы». Вместо этого он превращает фразу вроде:

> Создай SaaS для управления лидами из сервисов бронирования.

в контролируемый delivery-процесс:

```text
идея
  → проверка проблемы и рынка
  → Feature Contract
  → варианты архитектуры и ADR
  → UX / information architecture / design direction
  → план вертикальными срезами
  → реализация
  → независимая проверка
  → evidence packet
  → решение владельца о merge/release
```

## Что уже входит

- Go CLI `autoco` для инициализации нового или существующего репозитория.
- Профили `saas-web`, `fullstack-saas`, `expo-mobile`.
- Claude Code plugin с семью skills.
- Четыре узких subagent-роли: product research, architecture research, UX critique, release review.
- Registry проверенных внешних наборов: Spec Kit, gstack, Superpowers, UI UX Pro Max, Interface Design, Expo и Callstack.
- Durable project memory в `.auto-company/`.
- Product, architecture, UX, implementation и release-evidence шаблоны.
- Security hook, блокирующий чтение типичных секретов и опасные команды.
- Token/context policy: skills загружаются по необходимости, шумный research выносится в отдельный контекст, review-циклы ограничены.
- CI и локальные проверки самого toolkit.

## Что Auto Company сознательно не делает

- Не устанавливает сторонние plugins и scripts без твоего явного решения.
- Не копирует чужие skills внутрь репозитория и не создаёт лицензионную путаницу.
- Не разрешает AI автоматически merge’ить PR или выпускать production.
- Не выбирает технологию только по популярности.
- Не считает зелёные unit-тесты достаточным доказательством production readiness.
- Не запускает десятки субагентов на простую задачу.

## Быстрый старт на Windows

### 1. Клонировать Auto Company

```powershell
git clone https://github.com/aleka7sk/auto-company.git C:\Users\aleka\projects\auto-company
cd C:\Users\aleka\projects\auto-company
```

### 2. Проверить toolkit

```powershell
go test ./...
go run ./cmd/autoco version
go run ./cmd/autoco profiles
```

Можно установить CLI глобально в `GOBIN`:

```powershell
go install ./cmd/autoco
```

После этого команда будет называться `autoco`, если `%USERPROFILE%\go\bin` находится в `PATH`.

### 3. Подключить Claude Code plugin

Для локальной проверки без установки marketplace:

```powershell
claude --plugin-dir C:\Users\aleka\projects\auto-company
```

Для постоянной установки внутри Claude Code:

```text
/plugin marketplace add aleka7sk/auto-company
/plugin install auto-company@auto-company-marketplace
```

Перезапусти Claude Code после установки.

### 4. Создать или подключить продукт

Новый SaaS:

```powershell
mkdir C:\Users\aleka\projects\booking-os
cd C:\Users\aleka\projects\auto-company

go run ./cmd/autoco init `
  --target C:\Users\aleka\projects\booking-os `
  --profile fullstack-saas `
  --name "Booking OS" `
  --idea "Владельцы недвижимости теряют лиды из разных каналов и поздно связываются с клиентами"
```

Существующий Expo / React Native проект:

```powershell
go run ./cmd/autoco init `
  --target C:\Users\aleka\projects\belcanto-product `
  --profile expo-mobile `
  --name "Belcanto Product" `
  --idea "Помочь действующим ученикам видеть уроки, практику и реальный прогресс"
```

Инициализатор не перезаписывает существующие product-артефакты без `--force`. В `CLAUDE.md`, `AGENTS.md` и `.gitignore` он изменяет только собственный managed block.

### 5. Запустить работу в целевом проекте

```powershell
cd C:\Users\aleka\projects\booking-os
claude --plugin-dir C:\Users\aleka\projects\auto-company
```

Затем:

```text
/auto-company:create-saas

Создай SaaS, который собирает первичные бронирования из разных каналов,
распределяет их между продажниками и помогает владельцу контролировать
скорость обработки и конверсию. Начни с discovery и не пиши код до
утверждения Feature Contract.
```

Можно сначала посмотреть сгенерированный стартовый prompt:

```powershell
autoco prompt --target .
```

## Основной workflow

### Gate 1 — Opportunity

Агент проверяет проблему, пользователя, текущие обходные пути, рынок, конкурентов, возражения против идеи и самый дешёвый способ проверить рискованные предположения.

Результат:

```text
.auto-company/product/product-brief.md
```

### Gate 2 — Product contract

Фиксируются роли, permissions, user flow, бизнес-правила, функциональные и нефункциональные требования, состояния UI, аналитика и acceptance criteria.

Результат:

```text
.auto-company/product/feature-contract.md
```

Только утверждённый Feature Contract является продуктовой истиной для реализации.

### Gate 3 — Architecture

Researcher изучает существующий код и первичные источники, сравнивает минимум три реалистичных варианта, описывает trade-offs, failure modes, стоимость эксплуатации и обратимость.

Результат:

```text
.auto-company/architecture/ADR-0001.md
```

### Gate 4 — Experience

Сначала проектируются user flow, information architecture и все состояния. Только после этого выбирается визуальное направление. Внешние UI skills могут предложить варианты, но не имеют права менять утверждённую продуктовую логику.

Результат:

```text
.auto-company/ux/ux-spec.md
```

### Gate 5 — Delivery

Работа режется на небольшие end-to-end slices, каждый из которых приносит наблюдаемую ценность, имеет тесты, demo evidence и rollback.

Результат:

```text
.auto-company/delivery/implementation-plan.md
```

### Gate 6 — Production readiness

Независимый reviewer проверяет cumulative diff, traceability, permissions, tenant isolation, migrations, failures, security, accessibility, observability, backup, rollout, rollback и реальный основной сценарий.

Результат:

```text
.auto-company/evidence/release-evidence.md
```

## Команды CLI

```text
autoco init [flags]          инициализировать целевой проект
autoco doctor [flags]        проверить локальные инструменты и структуру
autoco validate [flags]      проверить обязательные артефакты
autoco prompt [flags]        вывести стартовый prompt
autoco profiles              показать профили продукта
autoco integrations [flags]  показать отобранные внешние integrations
autoco version               показать версию
```

Примеры:

```powershell
autoco doctor --target C:\Users\aleka\projects\booking-os
autoco validate --target C:\Users\aleka\projects\booking-os
autoco integrations --profile fullstack-saas --agent claude
autoco integrations --profile expo-mobile --agent codex --json
```

## Профили

| Profile | Когда использовать | Особые gates |
|---|---|---|
| `saas-web` | B2B/B2C web SaaS и dashboards | tenant isolation, migrations, responsive UX, observability |
| `fullstack-saas` | Полный SaaS с frontend, API, data и operations | product evidence, authz, payments/webhooks при наличии, security, rollback |
| `expo-mobile` | Expo/React Native приложение | permissions, offline/degraded states, iOS/Android behavior, real-device smoke test |

Профиль не жёстко навязывает framework. Например, `fullstack-saas` рекомендует modular monolith как стартовую гипотезу, но окончательное решение должно быть подтверждено ADR.

## Внешние integrations

Auto Company работает самостоятельно. Внешние наборы подключаются точечно:

```powershell
autoco integrations --profile fullstack-saas --agent claude
```

Рекомендуемое разделение ответственности:

| Набор | Роль |
|---|---|
| GitHub Spec Kit | связность spec → plan → tasks → implementation |
| gstack | founder, engineering и design critique |
| Superpowers | TDD, debugging, review, worktrees, verification |
| UI UX Pro Max | visual exploration, UI patterns, charts и design-system suggestions |
| Interface Design | устойчивые design decisions для dashboards/apps |
| Expo Skills | официальный Expo/EAS workflow |
| Callstack Agent Skills | production React Native performance, testing и device QA |

Не включай несколько end-to-end orchestrators одновременно. `auto-company:create-saas` остаётся главным процессом, а внешние plugins используются как специалисты.

## Структура целевого проекта

```text
.auto-company/
├── manifest.json
├── profile.json
├── quality-gates.json
├── product/
│   ├── product-brief.md
│   └── feature-contract.md
├── architecture/
│   └── ADR-0001.md
├── ux/
│   └── ux-spec.md
├── delivery/
│   └── implementation-plan.md
├── evidence/
│   └── release-evidence.md
├── prompts/
│   └── start.md
└── runs/                 # локальные run artifacts, не коммитятся
```

## Production-ready в Auto Company

Термин означает не обещание отсутствия ошибок, а доказанный набор условий:

- требования связаны с кодом и тестами;
- доступ и границы данных проверены;
- миграции и backward compatibility продуманы;
- есть loading/empty/error/permission/offline states;
- критические flows протестированы на целевой среде;
- логи, метрики, alerts, backup и restore определены;
- rollout и rollback воспроизводимы;
- нет unresolved blocker/high findings;
- владелец отдельно подтверждает merge и production release.

Подробнее: [`docs/policies/production-readiness.md`](docs/policies/production-readiness.md).

## Безопасность

Plugin hook блокирует типичные опасные операции до выполнения, включая чтение `.env`/ключей, force push, destructive Git cleanup, PR merge, infrastructure destroy, database drop, package publication и production submission/deploy commands.

Это дополнительный барьер, а не sandbox. Проверяй permissions Claude/Codex, используй отдельные credentials, branch protection и preview environment. Подробнее: [`SECURITY.md`](SECURITY.md).

## Токены и контекст

Auto Company уменьшает расход за счёт:

- маленького постоянного operating contract;
- progressive loading skills;
- отдельного контекста для noisy research и review;
- максимум двух research agents по умолчанию;
- максимум двух fix/review cycles;
- передачи между этапами через короткие durable artifacts, а не весь chat history;
- запрета полной загрузки репозитория без необходимости;
- разных workflow для tiny, normal и risky задач.

Подробнее: [`docs/policies/token-and-context.md`](docs/policies/token-and-context.md).

## Состояние проекта

Текущая версия — **foundation v0.1**. Она уже предоставляет работающий CLI, plugin, templates, policies, guards и tests. Следующие версии добавят управляемый runtime, автоматическую трассировку requirement IDs, budget accounting и интеграцию с GitHub PR evidence.

Смотри [`docs/roadmap.md`](docs/roadmap.md).

## Лицензия

MIT. Сторонние integrations сохраняют собственные лицензии и не входят в распространяемый код Auto Company.
