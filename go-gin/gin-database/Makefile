.PHONY: migrate serve run

migrate:
	go run ./migrate

serve:
	go run .

run: migrate
	@echo "Starting server..."
	@$(MAKE) serve