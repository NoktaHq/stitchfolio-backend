# Agent guidance for this app

This application is a **Go service** that uses **Pixie** (`github.com/loop-kar/pixie`) as its core framework.

## When working in this repo

1. **Use Pixie, don't reimplement**  
   Prefer Pixie's packages for config, DB, errors, HTTP responses, logging, middleware, storage, tasks, and utils. Do not add a parallel implementation of behavior that Pixie already provides.

2. **Follow Pixie conventions**  
   Use `errs.XError` and `errs.ErrorType` for errors; `response.Response` / `DataResponse` / `FileResponse` and `response.FormatAndSend` for HTTP; `config.LoadConfig` for app config; `db.DBTransactionManager.WithTransaction` and store `*gorm.DB` in context; `storage.StorageProvider` (or `CloudStorageProvider`) for file storage. See `.cursor/rules/pixie-app.mdc` for details.

3. **App-specific code**  
   Handlers, domain logic, and app config structs live here; infrastructure and shared conventions come from Pixie.

## Reference

- Pixie: `github.com/loop-kar/pixie`
- Consumer rules and templates: see Pixie repo `docs/cursor-for-consumers/` and `.cursor/rules/pixie-consumer.mdc`.
