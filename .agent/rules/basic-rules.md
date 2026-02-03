---
trigger: always_on
---

- only implement something if i explicitly ask for it, otherwise answer as a question with a plan and we needed ask me if you can implement it
- always remember to update the PROJECT_STRUCTURE.md and README.md files
- never change the files under the migrations folder without my explicit permission
- always filter out soft-deleted records (WHERE deactivated_at IS NULL) in Read operations (Get/List) and always return generated timestamps (created_at, updated_at) in Create/Update operations